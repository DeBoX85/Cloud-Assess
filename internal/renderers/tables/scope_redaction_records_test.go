package tables

import (
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

func TestAssessmentStatusScopeKnownCoreRecords(t *testing.T) {
	cases := []struct {
		name     string
		populate func(*result.Input)
	}{
		{"resources", func(i *result.Input) { i.Resources = []assessment.Resource{{SubscriptionID: scopeRedactionID}} }},
		{"out-of-scope", func(i *result.Input) { i.OutOfScope = []assessment.Resource{{SubscriptionID: scopeRedactionID}} }},
		{"findings", func(i *result.Input) {
			i.Findings = []assessment.Finding{{RecommendationID: "scope", SubscriptionID: scopeRedactionID, ResourceID: "/subscriptions/" + scopeRedactionID + "/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/vm", Category: "Security", Impact: "Low", Source: "CUSTOM"}}
		}},
		{"resource-types", func(i *result.Input) {
			i.ResourceTypes = []assessment.ResourceTypeCount{{SubscriptionID: scopeRedactionID}}
		}},
		{"advisor", func(i *result.Input) {
			i.Advisor = []assessment.AdvisorRecommendation{{SubscriptionID: scopeRedactionID}}
		}},
		{"defender", func(i *result.Input) {
			i.Defender = []assessment.DefenderPlanStatus{{SubscriptionID: scopeRedactionID}}
		}},
		{"defender-recommendations", func(i *result.Input) {
			i.DefenderRecommendations = []assessment.DefenderRecommendation{{SubscriptionID: scopeRedactionID}}
		}},
		{"policy", func(i *result.Input) {
			i.AzurePolicy = []assessment.PolicyNonCompliance{{SubscriptionID: scopeRedactionID}}
		}},
		{"arc-sql", func(i *result.Input) { i.ArcSQL = []assessment.ArcSQLRecord{{SubscriptionID: scopeRedactionID}} }},
		{"costs", func(i *result.Input) { i.Costs = []assessment.CostRecord{{SubscriptionID: scopeRedactionID}} }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input := result.Input{ScopeID: "opaque:" + scopeRedactionID + ":eastus", Completeness: assessment.CompletenessComplete, Recommendations: []assessment.RecommendationDefinition{{ID: "scope", Source: "CUSTOM", Category: "Security", Impact: "Low", ResourceType: "microsoft.compute/virtualmachines"}}}
			c.populate(&input)
			data := result.Build(input)
			rows := scopeStatusRows(t, data, true)
			if rows[1][2] == input.ScopeID {
				t.Fatal("known core record subscription in scope must mask")
			}
		})
	}
}
