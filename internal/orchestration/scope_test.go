package orchestration

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
)

type scopeLister struct{ subscriptions []discovery.Subscription }

func TestBlankExplicitScopeIsRejectedBeforeDiscovery(t *testing.T) {
	for _, request := range []Request{{Subscriptions: []string{" "}}, {ManagementGroups: []string{""}}} {
		operations := noOpOperations()
		operations.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
			t.Fatal("discovery ran for invalid scope")
			return nil, nil
		}
		operations.DiscoverManagementGroups = operations.DiscoverSubscriptions
		got, err := NewCoordinator(operations).Run(context.Background(), request)
		if err == nil || got != nil {
			t.Fatalf("blank scope result=%#v error=%v", got, err)
		}
	}
}

func (s scopeLister) ListSubscriptions(context.Context) ([]discovery.Subscription, error) {
	return s.subscriptions, nil
}

func TestExplicitScopeResolutionFailsBeforeInventoryWhenRequestedIDIsMissing(t *testing.T) {
	for _, tc := range []struct {
		name       string
		requested  []string
		available  []discovery.Subscription
		unresolved []string
		resolved   int
	}{
		{"absent", []string{"sub-1", "missing"}, []discovery.Subscription{{ID: "sub-1", State: "Enabled"}}, []string{"missing"}, 1},
		{"disabled", []string{"disabled"}, []discovery.Subscription{{ID: "disabled", State: "Disabled"}}, []string{"disabled"}, 0},
		{"deleted", []string{"deleted"}, []discovery.Subscription{{ID: "deleted", State: "Deleted"}}, []string{"deleted"}, 0},
		{"all missing", []string{"missing"}, nil, []string{"missing"}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			operations := noOpOperations()
			operations.DiscoverSubscriptions = func(ctx context.Context, requested []string, f *config.Filters) (map[string]string, error) {
				return discovery.DiscoverSubscriptions(ctx, scopeLister{tc.available}, requested, f)
			}
			operations.DiscoverResources = func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
				t.Fatal("resource query ran after unresolved explicit scope")
				return nil, nil
			}
			got, err := NewCoordinator(operations).Run(context.Background(), Request{Subscriptions: tc.requested, Stages: graphOnlyStages(t)})
			if err == nil || !strings.Contains(err.Error(), "scope_requested_subscription_unresolved") || got.Completeness != assessment.CompletenessFailed {
				t.Fatalf("resolution error=%v result=%#v", err, got)
			}
			if got.Scope.Status != "unresolved" || !reflect.DeepEqual(got.Scope.UnresolvedSubscriptionIDs, tc.unresolved) || len(got.Scope.ResolvedSubscriptions) != tc.resolved {
				t.Fatalf("wrong retained scope evidence: %#v", got.Scope)
			}
			if got.Stages[0].Status != assessment.StageFailed || got.Stages[0].Records != tc.resolved || got.Stages[1].Status != assessment.StageSkipped {
				t.Fatalf("scope/inventory status: %#v", got.Stages[:2])
			}
		})
	}
}

func TestScopeEvidenceRetainsZeroDataSubscriptionsAndFilterPrecedence(t *testing.T) {
	operations := noOpOperations()
	operations.DiscoverSubscriptions = func(ctx context.Context, ids []string, filters *config.Filters) (map[string]string, error) {
		return discovery.DiscoverSubscriptions(ctx, scopeLister{[]discovery.Subscription{{ID: "SUB-1", DisplayName: "One", State: "Enabled"}}}, ids, filters)
	}
	filters := config.NewFilters()
	filters.Assessment.Exclude.Subscriptions = []string{"sub-1"}
	got, err := NewCoordinator(operations).Run(context.Background(), Request{Subscriptions: []string{"SUB-1", "sub-1"}, Filters: filters, Stages: graphOnlyStages(t)})
	if err != nil || got.Scope.Status != "resolved" || len(got.Resources) != 0 || len(got.Scope.ResolvedSubscriptions) != 1 {
		t.Fatalf("zero-data scope: error=%v result=%#v", err, got)
	}
	if !reflect.DeepEqual(got.Scope.RequestedSubscriptionIDs, []string{"sub-1"}) || got.Scope.ResolvedSubscriptions[0].SubscriptionName != "One" || got.Scope.Selection != "subscriptions" {
		t.Fatalf("normalization or display name lost: %#v", got.Scope)
	}
	// An explicit include continues to override overlapping exclusions, as in source.
	if !reflect.DeepEqual(got.Scope.IncludedSubscriptionIDs, []string{"sub-1"}) || !reflect.DeepEqual(got.Scope.ExcludedSubscriptionIDs, []string{"sub-1"}) {
		t.Fatalf("filter intent lost: %#v", got.Scope)
	}
	got, err = NewCoordinator(operations).Run(context.Background(), Request{Filters: filters, Stages: graphOnlyStages(t)})
	if err != nil || got.Scope.Status != "resolved" || len(got.Scope.ResolvedSubscriptions) != 0 || len(got.Scope.UnresolvedSubscriptionIDs) != 0 {
		t.Fatalf("intentional exclusion became unresolved request: error=%v scope=%#v", err, got.Scope)
	}
	filters.Assessment.Include.Subscriptions = []string{"missing"}
	got, err = NewCoordinator(operations).Run(context.Background(), Request{Filters: filters, Stages: graphOnlyStages(t)})
	if err != nil || len(got.Scope.ResolvedSubscriptions) != 0 || len(got.Scope.UnresolvedSubscriptionIDs) != 0 || !reflect.DeepEqual(got.Scope.IncludedSubscriptionIDs, []string{"missing"}) {
		t.Fatalf("filter-only include became explicit CLI scope: error=%v scope=%#v", err, got.Scope)
	}
}

func TestScopeEvidenceDoesNotCertifyManagementGroupMembershipOrFailedListing(t *testing.T) {
	operations := noOpOperations()
	operations.DiscoverManagementGroups = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return map[string]string{}, nil
	}
	got, err := NewCoordinator(operations).Run(context.Background(), Request{ManagementGroups: []string{"Child"}, Stages: graphOnlyStages(t)})
	if err != nil || got.Scope.Selection != "management_groups" || got.Scope.Status != "resolved" || !reflect.DeepEqual(got.Scope.RequestedManagementGroups, []string{"child"}) || len(got.Scope.ResolvedSubscriptions) != 0 {
		t.Fatalf("empty MG scope: error=%v result=%#v", err, got)
	}
	operations.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return nil, errors.New("listing denied")
	}
	got, err = NewCoordinator(operations).Run(context.Background(), Request{Subscriptions: []string{"sub-1"}, Stages: graphOnlyStages(t)})
	if err == nil || got.Scope.Status != "not_completed" || !reflect.DeepEqual(got.Scope.RequestedSubscriptionIDs, []string{"sub-1"}) || len(got.Scope.UnresolvedSubscriptionIDs) != 0 {
		t.Fatalf("failed listing presented as resolved or classified missing: error=%v result=%#v", err, got)
	}
}
