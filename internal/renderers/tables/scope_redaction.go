package tables

import (
	"regexp"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

var statusScopeUUID = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

// Only identities present in this scope string are retained. Core-only reports
// therefore do not acquire the plugin masking identity limit or an estate-sized
// identity map. Contextual subscription paths are safe even without discovery.
func maskStatusScopeID(data *result.AssessmentResult, value string, enabled bool) string {
	if !enabled {
		return value
	}
	candidates := map[string]bool{}
	for _, id := range statusScopeUUID.FindAllString(value, -1) {
		candidates[strings.ToLower(id)] = false
	}
	if len(candidates) == 0 {
		return value
	}
	mark := func(id string) {
		key := strings.ToLower(id)
		if _, exists := candidates[key]; exists {
			candidates[key] = true
		}
	}
	for _, match := range pluginContextSubscriptionID.FindAllStringSubmatch(value, -1) {
		mark(match[1])
	}
	if data.Scope != nil {
		for _, group := range [][]string{data.Scope.RequestedSubscriptionIDs, data.Scope.IncludedSubscriptionIDs, data.Scope.ExcludedSubscriptionIDs, data.Scope.UnresolvedSubscriptionIDs} {
			for _, id := range group {
				mark(id)
			}
		}
		for _, sub := range data.Scope.ResolvedSubscriptions {
			mark(sub.SubscriptionID)
		}
	}
	for _, r := range data.Resources {
		mark(r.SubscriptionID)
	}
	for _, r := range data.OutOfScope {
		mark(r.SubscriptionID)
	}
	for _, r := range data.Findings {
		mark(r.SubscriptionID)
	}
	for _, r := range data.ResourceTypes {
		mark(r.SubscriptionID)
	}
	for _, r := range data.Advisor {
		mark(r.SubscriptionID)
	}
	for _, r := range data.Defender {
		mark(r.SubscriptionID)
	}
	for _, r := range data.DefenderRecommendations {
		mark(r.SubscriptionID)
	}
	for _, r := range data.AzurePolicy {
		mark(r.SubscriptionID)
	}
	for _, r := range data.ArcSQL {
		mark(r.SubscriptionID)
	}
	for _, r := range data.Costs {
		mark(r.SubscriptionID)
	}
	return statusScopeUUID.ReplaceAllStringFunc(value, func(id string) string {
		if candidates[strings.ToLower(id)] {
			return redact.SubscriptionID(id, true)
		}
		return id
	})
}

// Status rows are freshly projected; canonical result data is never changed.
func redactAssessmentStatusScope(data *result.AssessmentResult, rows [][]string, enabled bool, mask func(string) string) {
	value := maskStatusScopeID(data, mask(data.ScopeID), enabled)
	for _, row := range rows[1:] {
		row[2] = value
	}
}
