package result

import (
	"sort"
	"strings"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/findings"
)

const SchemaVersion = "1.0"

// AssessmentResult is the canonical in-memory assessment contract consumed by renderers.
// It intentionally separates primary recommendation findings from auxiliary Azure datasets.
type AssessmentResult struct {
	SchemaVersion           string                                `json:"schemaVersion"`
	GeneratedAt             time.Time                             `json:"generatedAt"`
	ScopeID                 string                                `json:"scopeId"`
	Completeness            assessment.Completeness               `json:"completeness"`
	Stages                  []assessment.StageExecution           `json:"stages"`
	Recommendations         []assessment.RecommendationDefinition `json:"recommendations"`
	Findings                []assessment.Finding                  `json:"findings"`
	Summary                 *findings.Summary                     `json:"summary"`
	Resources               []assessment.Resource                 `json:"resources"`
	OutOfScope              []assessment.Resource                 `json:"outOfScope"`
	ResourceTypes           []assessment.ResourceTypeCount        `json:"resourceTypes"`
	Advisor                 []assessment.AdvisorRecommendation    `json:"advisor"`
	Defender                []assessment.DefenderPlanStatus       `json:"defender"`
	DefenderRecommendations []assessment.DefenderRecommendation   `json:"defenderRecommendations"`
	AzurePolicy             []assessment.PolicyNonCompliance      `json:"azurePolicy"`
	ArcSQL                  []assessment.ArcSQLRecord             `json:"arcSQL"`
	Costs                   []assessment.CostRecord               `json:"costs"`
}

type Input struct {
	GeneratedAt             time.Time
	ScopeID                 string
	Completeness            assessment.Completeness
	Stages                  []assessment.StageExecution
	Recommendations         []assessment.RecommendationDefinition
	Findings                []assessment.Finding
	Resources               []assessment.Resource
	OutOfScope              []assessment.Resource
	ResourceTypes           []assessment.ResourceTypeCount
	Advisor                 []assessment.AdvisorRecommendation
	Defender                []assessment.DefenderPlanStatus
	DefenderRecommendations []assessment.DefenderRecommendation
	AzurePolicy             []assessment.PolicyNonCompliance
	ArcSQL                  []assessment.ArcSQLRecord
	Costs                   []assessment.CostRecord
}

// Build creates an isolated, deterministically ordered result without mutating caller-owned data.
func Build(input Input) *AssessmentResult {
	result := &AssessmentResult{
		SchemaVersion:           SchemaVersion,
		GeneratedAt:             input.GeneratedAt.UTC(),
		ScopeID:                 input.ScopeID,
		Completeness:            input.Completeness,
		Stages:                  cloneStages(input.Stages),
		Recommendations:         cloneRecommendations(input.Recommendations),
		Findings:                cloneFindings(input.Findings),
		Resources:               cloneResources(input.Resources),
		OutOfScope:              cloneResources(input.OutOfScope),
		ResourceTypes:           append([]assessment.ResourceTypeCount(nil), input.ResourceTypes...),
		Advisor:                 append([]assessment.AdvisorRecommendation(nil), input.Advisor...),
		Defender:                append([]assessment.DefenderPlanStatus(nil), input.Defender...),
		DefenderRecommendations: append([]assessment.DefenderRecommendation(nil), input.DefenderRecommendations...),
		AzurePolicy:             append([]assessment.PolicyNonCompliance(nil), input.AzurePolicy...),
		ArcSQL:                  append([]assessment.ArcSQLRecord(nil), input.ArcSQL...),
		Costs:                   append([]assessment.CostRecord(nil), input.Costs...),
	}

	sortCanonical(result)
	result.Summary = findings.Build(result.Recommendations, result.Findings, len(result.Resources))
	return result
}

func sortCanonical(result *AssessmentResult) {
	sort.Slice(result.Recommendations, func(i, j int) bool {
		left, right := result.Recommendations[i], result.Recommendations[j]
		if compare(left.ResourceType, right.ResourceType) != 0 {
			return compare(left.ResourceType, right.ResourceType) < 0
		}
		if compare(left.ID, right.ID) != 0 {
			return compare(left.ID, right.ID) < 0
		}
		return compare(left.Source, right.Source) < 0
	})
	sort.Slice(result.Findings, func(i, j int) bool {
		left, right := result.Findings[i], result.Findings[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		if compare(left.ResourceID, right.ResourceID) != 0 {
			return compare(left.ResourceID, right.ResourceID) < 0
		}
		if compare(left.RecommendationID, right.RecommendationID) != 0 {
			return compare(left.RecommendationID, right.RecommendationID) < 0
		}
		return compare(left.Category, right.Category) < 0
	})
	sortResources(result.Resources)
	sortResources(result.OutOfScope)
	sort.Slice(result.ResourceTypes, func(i, j int) bool {
		left, right := result.ResourceTypes[i], result.ResourceTypes[j]
		if compare(left.SubscriptionName, right.SubscriptionName) != 0 {
			return compare(left.SubscriptionName, right.SubscriptionName) < 0
		}
		if compare(left.ResourceType, right.ResourceType) != 0 {
			return compare(left.ResourceType, right.ResourceType) < 0
		}
		return compare(left.SubscriptionID, right.SubscriptionID) < 0
	})
	sort.Slice(result.Advisor, func(i, j int) bool {
		left, right := result.Advisor[i], result.Advisor[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		if compare(left.ResourceID, right.ResourceID) != 0 {
			return compare(left.ResourceID, right.ResourceID) < 0
		}
		return compare(left.RecommendationID, right.RecommendationID) < 0
	})
	sort.Slice(result.Defender, func(i, j int) bool {
		left, right := result.Defender[i], result.Defender[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		return compare(left.Name, right.Name) < 0
	})
	sort.Slice(result.DefenderRecommendations, func(i, j int) bool {
		left, right := result.DefenderRecommendations[i], result.DefenderRecommendations[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		if compare(left.ResourceID, right.ResourceID) != 0 {
			return compare(left.ResourceID, right.ResourceID) < 0
		}
		if compare(left.Category, right.Category) != 0 {
			return compare(left.Category, right.Category) < 0
		}
		return compare(left.RecommendationName, right.RecommendationName) < 0
	})
	sort.Slice(result.AzurePolicy, func(i, j int) bool {
		left, right := result.AzurePolicy[i], result.AzurePolicy[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		if compare(left.ResourceID, right.ResourceID) != 0 {
			return compare(left.ResourceID, right.ResourceID) < 0
		}
		return compare(left.PolicyDefinitionID, right.PolicyDefinitionID) < 0
	})
	sort.Slice(result.ArcSQL, func(i, j int) bool {
		left, right := result.ArcSQL[i], result.ArcSQL[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		return compare(left.SQLInstance, right.SQLInstance) < 0
	})
	sort.Slice(result.Costs, func(i, j int) bool {
		left, right := result.Costs[i], result.Costs[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		if compare(left.ServiceName, right.ServiceName) != 0 {
			return compare(left.ServiceName, right.ServiceName) < 0
		}
		if compare(left.Currency, right.Currency) != 0 {
			return compare(left.Currency, right.Currency) < 0
		}
		return compare(left.Value, right.Value) < 0
	})
}

func sortResources(resources []assessment.Resource) {
	sort.Slice(resources, func(i, j int) bool {
		left, right := resources[i], resources[j]
		if compare(left.SubscriptionID, right.SubscriptionID) != 0 {
			return compare(left.SubscriptionID, right.SubscriptionID) < 0
		}
		if compare(left.ResourceGroup, right.ResourceGroup) != 0 {
			return compare(left.ResourceGroup, right.ResourceGroup) < 0
		}
		if compare(left.Type, right.Type) != 0 {
			return compare(left.Type, right.Type) < 0
		}
		if compare(left.Name, right.Name) != 0 {
			return compare(left.Name, right.Name) < 0
		}
		return compare(left.ID, right.ID) < 0
	})
}

func cloneStages(values []assessment.StageExecution) []assessment.StageExecution {
	out := append([]assessment.StageExecution(nil), values...)
	for i := range out {
		out[i].Warnings = append([]assessment.AssessmentWarning(nil), values[i].Warnings...)
		if values[i].Error != nil {
			errorCopy := *values[i].Error
			out[i].Error = &errorCopy
		}
	}
	return out
}

func cloneRecommendations(values []assessment.RecommendationDefinition) []assessment.RecommendationDefinition {
	out := append([]assessment.RecommendationDefinition(nil), values...)
	for i := range out {
		out[i].Tags = append([]string(nil), values[i].Tags...)
		out[i].LearnMore = append([]assessment.LearnMoreLink(nil), values[i].LearnMore...)
	}
	return out
}

func cloneFindings(values []assessment.Finding) []assessment.Finding {
	out := append([]assessment.Finding(nil), values...)
	for i := range out {
		out[i].Parameters = append([]string(nil), values[i].Parameters...)
	}
	return out
}

func cloneResources(values []assessment.Resource) []assessment.Resource {
	out := append([]assessment.Resource(nil), values...)
	for i := range out {
		if values[i].Tags == nil {
			continue
		}
		out[i].Tags = make(map[string]string, len(values[i].Tags))
		for key, value := range values[i].Tags {
			out[i].Tags[key] = value
		}
	}
	return out
}

func compare(left, right string) int {
	left = strings.ToLower(strings.TrimSpace(left))
	right = strings.ToLower(strings.TrimSpace(right))
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}
