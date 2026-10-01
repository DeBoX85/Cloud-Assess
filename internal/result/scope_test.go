package result

import (
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestBuildOwnsAndOrdersScopeEvidence(t *testing.T) {
	input := &assessment.ScopeResolution{
		RequestedSubscriptionIDs: []string{"b", "a"}, RequestedManagementGroups: []string{"b", "a"},
		IncludedSubscriptionIDs: []string{"b", "a"}, ExcludedSubscriptionIDs: []string{"b", "a"}, UnresolvedSubscriptionIDs: []string{"b", "a"},
		ResolvedSubscriptions: []assessment.ScopeSubscription{{SubscriptionID: "b", SubscriptionName: "B"}, {SubscriptionID: "a", SubscriptionName: "A"}},
	}
	got := Build(Input{Scope: input})
	if input.RequestedSubscriptionIDs[0] != "b" || input.ResolvedSubscriptions[0].SubscriptionID != "b" {
		t.Fatal("Build reordered caller data")
	}
	for _, ids := range [][]string{got.Scope.RequestedSubscriptionIDs, got.Scope.RequestedManagementGroups, got.Scope.IncludedSubscriptionIDs, got.Scope.ExcludedSubscriptionIDs, got.Scope.UnresolvedSubscriptionIDs} {
		if ids[0] != "a" {
			t.Fatal("scope IDs not canonical")
		}
		ids[0] = "changed"
	}
	got.Scope.ResolvedSubscriptions[0].SubscriptionName = "changed"
	for _, ids := range [][]string{input.RequestedSubscriptionIDs, input.RequestedManagementGroups, input.IncludedSubscriptionIDs, input.ExcludedSubscriptionIDs, input.UnresolvedSubscriptionIDs} {
		if ids[1] != "a" {
			t.Fatal("scope arrays alias caller")
		}
	}
	if input.ResolvedSubscriptions[1].SubscriptionName != "A" {
		t.Fatal("resolved subscriptions alias caller")
	}
}
