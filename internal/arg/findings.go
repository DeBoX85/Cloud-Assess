package arg

import (
	"encoding/json"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	azureid "github.com/DeBoX85/Cloud-Assess/internal/azure"
)

// FindingRow is the normalized shape projected by embedded assessment KQL.
type FindingRow struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Tags   json.RawMessage `json:"tags"`
	Param1 json.RawMessage `json:"param1"`
	Param2 json.RawMessage `json:"param2"`
	Param3 json.RawMessage `json:"param3"`
	Param4 json.RawMessage `json:"param4"`
	Param5 json.RawMessage `json:"param5"`
}

// FindingsFromRows converts decoded ARG rows into canonical findings.
// As in the reference implementation, encountering a row without an id stops processing
// the remaining rows for that recommendation.
func FindingsFromRows(
	definition assessment.RecommendationDefinition,
	rows []FindingRow,
	subscriptions map[string]string,
) []assessment.Finding {
	findings := make([]assessment.Finding, 0, len(rows))
	learnURL := ""
	if len(definition.LearnMore) > 0 {
		learnURL = definition.LearnMore[0].URL
	}

	for _, row := range rows {
		if row.ID == "" {
			break
		}

		subscriptionID := azureid.SubscriptionFromResourceID(row.ID)
		resourceType := azureid.ResourceTypeFromResourceID(row.ID)
		if resourceType == "" {
			resourceType = definition.ResourceType
		}

		findings = append(findings, assessment.Finding{
			RecommendationID:    definition.ID,
			Source:              definition.Source,
			ValidationMechanism: definition.ValidationMechanism,
			Category:            definition.Category,
			Impact:              definition.Impact,
			ResourceType:        resourceType,
			Recommendation:      definition.Recommendation,
			LongDescription:     definition.LongDescription,
			PotentialBenefits:   definition.PotentialBenefits,
			ResourceID:          row.ID,
			SubscriptionID:      subscriptionID,
			SubscriptionName:    subscriptions[subscriptionID],
			ResourceGroup:       azureid.ResourceGroupFromResourceID(row.ID),
			ResourceName:        row.Name,
			Tags:                RawMessageString(row.Tags),
			Parameters: []string{
				RawMessageString(row.Param1),
				RawMessageString(row.Param2),
				RawMessageString(row.Param3),
				RawMessageString(row.Param4),
				RawMessageString(row.Param5),
			},
			LearnMoreURL:        learnURL,
			AutomationAvailable: definition.AutomationAvailable,
		})
	}

	return findings
}
