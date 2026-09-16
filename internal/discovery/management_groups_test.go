package discovery

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

type fakeManagementGroupLister struct {
	subscriptions map[string][]Subscription
	descendants   map[string][]string
	subErr        map[string]error
	descErr       map[string]error
	visitedSubs   []string
	visitedDesc   []string
}

func (f *fakeManagementGroupLister) SubscriptionsUnderManagementGroup(_ context.Context, group string) ([]Subscription, error) {
	f.visitedSubs = append(f.visitedSubs, group)
	if err := f.subErr[group]; err != nil {
		return nil, err
	}
	return f.subscriptions[group], nil
}

func (f *fakeManagementGroupLister) DescendantManagementGroups(_ context.Context, group string) ([]string, error) {
	f.visitedDesc = append(f.visitedDesc, group)
	if err := f.descErr[group]; err != nil {
		return nil, err
	}
	return f.descendants[group], nil
}

func TestDiscoverManagementGroupSubscriptionsRecursesAndDeduplicates(t *testing.T) {
	lister := &fakeManagementGroupLister{
		subscriptions: map[string][]Subscription{
			"root":  {{ID: "sub-root", DisplayName: "Root", State: "Enabled"}},
			"child": {{ID: "sub-child", DisplayName: "Child", State: "Enabled"}},
			"leaf": {
				{ID: "sub-leaf", DisplayName: "Leaf", State: "Enabled"},
				{ID: "sub-deleted", DisplayName: "Deleted", State: "Deleted"},
			},
		},
		descendants: map[string][]string{
			"root":  {"child", "leaf"},
			"child": {"leaf"},
		},
	}

	got, err := DiscoverManagementGroupSubscriptions(context.Background(), lister, []string{"root"}, config.NewFilters())
	if err != nil {
		t.Fatalf("DiscoverManagementGroupSubscriptions() error = %v", err)
	}
	if len(got) != 3 || got["sub-root"] != "Root" || got["sub-child"] != "Child" || got["sub-leaf"] != "Leaf" {
		t.Fatalf("unexpected subscriptions: %#v", got)
	}
	if !reflect.DeepEqual(lister.visitedSubs, []string{"root", "child", "leaf"}) {
		t.Fatalf("visited groups = %#v, want each group once", lister.visitedSubs)
	}
}

func TestDiscoverManagementGroupSubscriptionsAppliesSubscriptionFilter(t *testing.T) {
	lister := &fakeManagementGroupLister{
		subscriptions: map[string][]Subscription{
			"root": {
				{ID: "sub-a", DisplayName: "Alpha", State: "Enabled"},
				{ID: "sub-b", DisplayName: "Beta", State: "Enabled"},
			},
		},
		descendants: map[string][]string{},
	}
	filters := config.NewFilters()
	filters.Assessment.Include.Subscriptions = []string{"sub-a"}
	filters.RebuildIndexes()

	got, err := DiscoverManagementGroupSubscriptions(context.Background(), lister, []string{"root"}, filters)
	if err != nil {
		t.Fatalf("DiscoverManagementGroupSubscriptions() error = %v", err)
	}
	if len(got) != 1 || got["sub-a"] != "Alpha" {
		t.Fatalf("unexpected filtered subscriptions: %#v", got)
	}
}

func TestDiscoverManagementGroupSubscriptionsReturnsSubscriptionError(t *testing.T) {
	boom := errors.New("subscriptions failed")
	lister := &fakeManagementGroupLister{subErr: map[string]error{"root": boom}}
	got, err := DiscoverManagementGroupSubscriptions(context.Background(), lister, []string{"root"}, nil)
	if err == nil || !errors.Is(err, boom) || got != nil {
		t.Fatalf("unexpected failure: got=%#v err=%v", got, err)
	}
}

func TestDiscoverManagementGroupSubscriptionsReturnsDescendantError(t *testing.T) {
	boom := errors.New("descendants failed")
	lister := &fakeManagementGroupLister{descErr: map[string]error{"root": boom}}
	got, err := DiscoverManagementGroupSubscriptions(context.Background(), lister, []string{"root"}, nil)
	if err == nil || !errors.Is(err, boom) || got != nil {
		t.Fatalf("unexpected failure: got=%#v err=%v", got, err)
	}
}
