package discovery

import (
	"context"
	"errors"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

type fakeSubscriptionLister struct {
	subscriptions []Subscription
	err           error
}

func (f *fakeSubscriptionLister) ListSubscriptions(context.Context) ([]Subscription, error) {
	return f.subscriptions, f.err
}

func TestDiscoverSubscriptionsAllAccessibleActive(t *testing.T) {
	lister := &fakeSubscriptionLister{subscriptions: []Subscription{
		{ID: "a", DisplayName: "Alpha", State: "Enabled"},
		{ID: "b", DisplayName: "Beta", State: "Disabled"},
		{ID: "c", DisplayName: "Gamma", State: "Deleted"},
		{ID: "d", DisplayName: "Delta", State: "Warned"},
	}}
	got, err := DiscoverSubscriptions(context.Background(), lister, nil, config.NewFilters())
	if err != nil {
		t.Fatalf("DiscoverSubscriptions() error = %v", err)
	}
	if len(got) != 2 || got["a"] != "Alpha" || got["d"] != "Delta" {
		t.Fatalf("unexpected subscriptions: %#v", got)
	}
}

func TestDiscoverSubscriptionsExplicitSubsetAndFilters(t *testing.T) {
	lister := &fakeSubscriptionLister{subscriptions: []Subscription{
		{ID: "A", DisplayName: "Alpha", State: "Enabled"},
		{ID: "B", DisplayName: "Beta", State: "Enabled"},
		{ID: "C", DisplayName: "Gamma", State: "Enabled"},
	}}
	filters := config.NewFilters()
	filters.Assessment.Exclude.Subscriptions = []string{"b"}
	filters.RebuildIndexes()

	got, err := DiscoverSubscriptions(context.Background(), lister, []string{"a", "B"}, filters)
	if err != nil {
		t.Fatalf("DiscoverSubscriptions() error = %v", err)
	}
	if len(got) != 1 || got["A"] != "Alpha" {
		t.Fatalf("unexpected explicit subscription result: %#v", got)
	}
}

func TestDiscoverSubscriptionsReturnsListerError(t *testing.T) {
	boom := errors.New("list failed")
	got, err := DiscoverSubscriptions(context.Background(), &fakeSubscriptionLister{err: boom}, nil, nil)
	if err == nil || !errors.Is(err, boom) || got != nil {
		t.Fatalf("unexpected failure result: got=%#v err=%v", got, err)
	}
}

func TestScopeIDIsDeterministicAndScopeSensitive(t *testing.T) {
	one := ScopeID(map[string]string{"b": "Beta", "a": "Alpha"})
	two := ScopeID(map[string]string{"a": "Other Name", "b": "Different Name"})
	three := ScopeID(map[string]string{"a": "Alpha"})
	if one != two {
		t.Fatalf("scope ID should depend only on sorted subscription IDs: %q != %q", one, two)
	}
	if one == three {
		t.Fatal("different subscription scope should produce different scope ID")
	}
	if len(one) != 64 {
		t.Fatalf("scope ID length = %d, want 64 hex chars", len(one))
	}
}
