package result

import (
	"sort"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func cloneScope(input *assessment.ScopeResolution) *assessment.ScopeResolution {
	if input == nil {
		return nil
	}
	out := *input
	clone := func(ids []string) []string {
		result := append([]string{}, ids...)
		sort.Strings(result)
		return result
	}
	out.RequestedSubscriptionIDs = clone(input.RequestedSubscriptionIDs)
	out.RequestedManagementGroups = clone(input.RequestedManagementGroups)
	out.IncludedSubscriptionIDs = clone(input.IncludedSubscriptionIDs)
	out.ExcludedSubscriptionIDs = clone(input.ExcludedSubscriptionIDs)
	out.UnresolvedSubscriptionIDs = clone(input.UnresolvedSubscriptionIDs)
	out.ResolvedSubscriptions = append([]assessment.ScopeSubscription{}, input.ResolvedSubscriptions...)
	sort.Slice(out.ResolvedSubscriptions, func(i, j int) bool {
		return out.ResolvedSubscriptions[i].SubscriptionID < out.ResolvedSubscriptions[j].SubscriptionID
	})
	return &out
}
