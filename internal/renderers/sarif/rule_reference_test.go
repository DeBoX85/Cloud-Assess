package sarif

import (
	"encoding/json"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

// The canonical summary joins recommendation identities without case, while
// SARIF references must identify the exact descriptor emitted by this report.
func TestRuleReferencesUsePublishedDescriptorIdentity(t *testing.T) {
	const lower = "abcdefab-1234-4234-8234-123456789abc"
	const upper = "ABCDEFAB-1234-4234-8234-123456789ABC"
	for _, tc := range []struct{ name, definition, finding string }{
		{"same identity control", lower, lower},
		{"upper finding", lower, upper},
		{"upper descriptor", upper, lower},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := result.Build(result.Input{
				Recommendations: []assessment.RecommendationDefinition{{ID: tc.definition, Recommendation: "Synthetic rule", Category: "Security", Impact: "High"}},
				Findings:        []assessment.Finding{{RecommendationID: tc.finding, Recommendation: "Synthetic rule", Category: "Security", Impact: "High", ResourceID: "/subscriptions/11111111-1111-4111-8111-111111111111/resourceGroups/synthetic/providers/Microsoft.Test/widgets/one"}},
			})
			encoded, err := Marshal(data, "synthetic")
			if err != nil {
				t.Fatal(err)
			}
			var report struct {
				Runs []struct {
					Tool struct {
						Driver struct {
							Rules []struct {
								ID string `json:"id"`
							} `json:"rules"`
						} `json:"driver"`
					} `json:"tool"`
					Results []struct {
						RuleID string `json:"ruleId"`
					} `json:"results"`
				} `json:"runs"`
			}
			if err := json.Unmarshal(encoded, &report); err != nil {
				t.Fatal(err)
			}
			if len(report.Runs) != 1 || len(report.Runs[0].Tool.Driver.Rules) != 1 || len(report.Runs[0].Results) != 1 {
				t.Fatalf("expected one descriptor and finding: %s", encoded)
			}
			run := report.Runs[0]
			if run.Tool.Driver.Rules[0].ID != tc.definition || run.Results[0].RuleID != tc.definition {
				t.Fatalf("SARIF result does not reference its published descriptor: descriptor=%q result=%q want=%q", run.Tool.Driver.Rules[0].ID, run.Results[0].RuleID, tc.definition)
			}
			if data.Findings[0].RecommendationID != tc.finding || data.Recommendations[0].ID != tc.definition {
				t.Fatal("SARIF rendering changed canonical recommendation identity")
			}
		})
	}
}
