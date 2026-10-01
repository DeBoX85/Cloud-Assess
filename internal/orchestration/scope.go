package orchestration

import (
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func scopeIDs(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		if value = strings.ToLower(strings.TrimSpace(value)); value != "" {
			set[value] = true
		}
	}
	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func newScopeResolution(request preparedRequest) *assessment.ScopeResolution {
	selection := "accessible_subscriptions"
	if len(request.managementGroups) > 0 {
		selection = "management_groups"
	} else if len(request.subscriptions) > 0 {
		selection = "subscriptions"
	}
	return &assessment.ScopeResolution{
		Selection:                 selection,
		Status:                    "not_completed",
		RequestedSubscriptionIDs:  scopeIDs(request.subscriptions),
		RequestedManagementGroups: scopeIDs(request.managementGroups),
		IncludedSubscriptionIDs:   scopeIDs(request.filters.Assessment.Include.Subscriptions),
		ExcludedSubscriptionIDs:   scopeIDs(request.filters.Assessment.Exclude.Subscriptions),
		ResolvedSubscriptions:     []assessment.ScopeSubscription{},
		UnresolvedSubscriptionIDs: []string{},
	}
}

func resolveScope(scope *assessment.ScopeResolution, subscriptions map[string]string) {
	resolved := map[string]string{}
	for id, name := range subscriptions {
		id = strings.ToLower(strings.TrimSpace(id))
		// Keep duplicate-case inputs deterministic without modifying the query scope.
		if existing, ok := resolved[id]; !ok || name < existing {
			resolved[id] = name
		}
	}
	for _, id := range scopeIDs(keys(resolved)) {
		scope.ResolvedSubscriptions = append(scope.ResolvedSubscriptions, assessment.ScopeSubscription{SubscriptionID: id, SubscriptionName: resolved[id]})
	}
	scope.Status = "resolved"
	for _, id := range scope.RequestedSubscriptionIDs {
		if _, ok := resolved[id]; !ok {
			scope.UnresolvedSubscriptionIDs = append(scope.UnresolvedSubscriptionIDs, id)
		}
	}
	if len(scope.UnresolvedSubscriptionIDs) > 0 {
		scope.Status = "unresolved"
	}
}

func keys(values map[string]string) []string {
	ids := make([]string, 0, len(values))
	for id := range values {
		ids = append(ids, id)
	}
	return ids
}
