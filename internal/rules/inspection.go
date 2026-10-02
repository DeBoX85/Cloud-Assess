// Portions reproduce the pinned Microsoft Azure Quick Review rules inspection
// contract (MIT licensed). See NOTICE.md.

package rules

import (
	"sort"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// InspectionRow is the source-compatible offline rules projection. This is not
// an assessment result and intentionally contains no Azure inventory or queries.
type InspectionRow struct {
	Category         string `json:"category"`
	Impact           string `json:"impact"`
	LearnMoreURL     string `json:"learnMoreUrl"`
	Recommendation   string `json:"recommendation"`
	RecommendationID string `json:"recommendationId"`
	ResourceType     string `json:"resourceType"`
}

// Inspect selects scanner-supported embedded rules followed by Diagnostics for
// each type, with global ID replacement and sorted output as in the reference.
// resourceTypes must retain scanner order to preserve duplicate-ID precedence.
// It does not load external plugins, query Azure or apply operator scan filters.
func Inspect(catalog *Catalog, resourceTypes []string, diagnostics []assessment.RecommendationDefinition) []InspectionRow {
	byID := map[string]InspectionRow{}
	diagnosticsByType := map[string][]assessment.RecommendationDefinition{}
	for _, definition := range diagnostics {
		key := normalize(definition.ResourceType)
		diagnosticsByType[key] = append(diagnosticsByType[key], definition)
	}
	add := func(definition assessment.RecommendationDefinition) {
		row := InspectionRow{
			Category: definition.Category, Impact: string(definition.Impact),
			Recommendation: definition.Recommendation, RecommendationID: definition.ID,
			ResourceType: definition.ResourceType,
		}
		if len(definition.LearnMore) > 0 {
			row.LearnMoreURL = definition.LearnMore[0].URL
		}
		byID[definition.ID] = row
	}
	for _, resourceType := range resourceTypes {
		for _, definition := range catalog.ByResourceType(resourceType) {
			if IsExecutableEmbedded(definition, nil) {
				add(definition)
			}
		}
		for _, definition := range diagnosticsByType[normalize(resourceType)] {
			add(definition)
		}
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rows := make([]InspectionRow, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, byID[id])
	}
	return rows
}

// ResourceTypeCount includes every type in the embedded catalog, matching the
// reference rules heading. It is not a count of observed or scanned Azure types.
func (c *Catalog) ResourceTypeCount() int { return len(c.byResourceType) }
