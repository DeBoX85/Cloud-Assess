package findings

import (
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type RecommendationSummary struct {
	ID                string `json:"id"`
	Recommendation    string `json:"recommendation"`
	Source            string `json:"source"`
	Category          string `json:"category"`
	Impact            string `json:"impact"`
	ResourceType      string `json:"resourceType"`
	LongDescription   string `json:"longDescription,omitempty"`
	LearnURL          string `json:"learnUrl,omitempty"`
	ImpactedResources int    `json:"impactedResources"`
}

type Summary struct {
	Recommendations      []RecommendationSummary `json:"recommendations"`
	Resources            int                     `json:"resources"`
	ImpactedResources    int                     `json:"impactedResources"`
	ImpactedByImpact     map[string]int          `json:"impactedByImpact"`
	ImpactedByCategory   map[string]int          `json:"impactedByCategory"`
	RecommendationsFound int                     `json:"recommendationsFound"`
}

type findingKey struct {
	recommendationID string
	resourceID       string
}

func Build(definitions []assessment.RecommendationDefinition, results []assessment.Finding, resourceCount int) *Summary {
	records := make(map[string]RecommendationSummary)
	for _, definition := range definitions {
		if strings.EqualFold(definition.Category, assessment.CategorySLA) {
			continue
		}
		id := normalize(definition.ID)
		if id == "" {
			continue
		}
		learnURL := ""
		if len(definition.LearnMore) > 0 {
			learnURL = definition.LearnMore[0].URL
		}
		records[id] = RecommendationSummary{
			ID:              definition.ID,
			Recommendation:  definition.Recommendation,
			Source:          definition.Source,
			Category:        definition.Category,
			Impact:          definition.Impact,
			ResourceType:    definition.ResourceType,
			LongDescription: definition.LongDescription,
			LearnURL:        learnURL,
		}
	}

	counts := make(map[string]int)
	seen := make(map[findingKey]struct{}, len(results))
	for _, result := range results {
		if strings.EqualFold(result.Category, assessment.CategorySLA) {
			continue
		}
		id := normalize(result.RecommendationID)
		if id == "" {
			continue
		}
		if _, ok := records[id]; !ok {
			continue
		}
		key := findingKey{recommendationID: id, resourceID: normalize(result.ResourceID)}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		counts[id]++
	}

	summary := &Summary{
		Resources:          resourceCount,
		ImpactedByImpact:   map[string]int{},
		ImpactedByCategory: map[string]int{},
	}
	for id, record := range records {
		record.ImpactedResources = counts[id]
		summary.Recommendations = append(summary.Recommendations, record)
		if record.ImpactedResources == 0 {
			continue
		}
		summary.RecommendationsFound++
		summary.ImpactedResources += record.ImpactedResources
		summary.ImpactedByImpact[record.Impact] += record.ImpactedResources
		summary.ImpactedByCategory[record.Category] += record.ImpactedResources
	}

	sort.Slice(summary.Recommendations, func(i, j int) bool {
		if summary.Recommendations[i].ResourceType != summary.Recommendations[j].ResourceType {
			return summary.Recommendations[i].ResourceType < summary.Recommendations[j].ResourceType
		}
		return summary.Recommendations[i].ID < summary.Recommendations[j].ID
	})
	return summary
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
