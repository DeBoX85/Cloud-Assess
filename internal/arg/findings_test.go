package arg

import (
	"encoding/json"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestFindingsFromRowsMapsReferenceFields(t *testing.T) {
	definition := assessment.RecommendationDefinition{
		ID:                  "rec-1",
		Recommendation:      "Fix widget",
		Category:            "Security",
		Impact:              "High",
		ResourceType:        "Microsoft.Test/widgets",
		LongDescription:     "Long guidance",
		PotentialBenefits:   "Benefit",
		AutomationAvailable: "true",
		Source:              "CUSTOM",
		ValidationMechanism: "Azure Resource Graph",
		LearnMore:           []assessment.LearnMoreLink{{Name: "Docs", URL: "https://example.test/docs"}},
	}
	id := "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	rows := []FindingRow{{
		ID:     id,
		Name:   "one",
		Tags:   json.RawMessage(`{"env":"prod"}`),
		Param1: json.RawMessage(`"hello"`),
		Param2: json.RawMessage(`42`),
		Param3: json.RawMessage(`true`),
		Param4: json.RawMessage(`["a"]`),
		Param5: json.RawMessage(`null`),
	}}

	findings := FindingsFromRows(definition, rows, map[string]string{
		"12345678-1234-1234-1234-123456789012": "Production",
	})
	if len(findings) != 1 {
		t.Fatalf("finding count = %d, want 1", len(findings))
	}
	finding := findings[0]
	if finding.RecommendationID != "rec-1" || finding.ResourceID != id || finding.ResourceName != "one" {
		t.Fatalf("unexpected finding identity: %+v", finding)
	}
	if finding.SubscriptionName != "Production" || finding.ResourceGroup != "rg" || finding.ResourceType != "Microsoft.Test/widgets" {
		t.Fatalf("unexpected Azure scope mapping: %+v", finding)
	}
	if finding.Tags != `{"env":"prod"}` {
		t.Fatalf("tags = %q", finding.Tags)
	}
	wantParams := []string{"hello", "42", "true", `["a"]`, ""}
	for i, want := range wantParams {
		if finding.Parameters[i] != want {
			t.Fatalf("parameter %d = %q, want %q", i+1, finding.Parameters[i], want)
		}
	}
	if finding.LearnMoreURL != "https://example.test/docs" {
		t.Fatalf("LearnMoreURL = %q", finding.LearnMoreURL)
	}
}

func TestFindingsFromRowsStopsAtMissingID(t *testing.T) {
	definition := assessment.RecommendationDefinition{ID: "rec", ResourceType: "Microsoft.Test/widgets"}
	rows := []FindingRow{
		{ID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/first", Name: "first"},
		{Name: "missing-id"},
		{ID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/third", Name: "third"},
	}
	findings := FindingsFromRows(definition, rows, nil)
	if len(findings) != 1 || findings[0].ResourceName != "first" {
		t.Fatalf("missing id should stop further processing: %#v", findings)
	}
}

func TestFindingsFromRowsFallsBackToDefinitionResourceType(t *testing.T) {
	definition := assessment.RecommendationDefinition{ID: "rec", ResourceType: "Specialized.Workload/Example"}
	rows := []FindingRow{{ID: "/subscriptions/sub/resourceGroups/rg", Name: "scope"}}
	findings := FindingsFromRows(definition, rows, nil)
	if len(findings) != 1 || findings[0].ResourceType != definition.ResourceType {
		t.Fatalf("resource type fallback failed: %#v", findings)
	}
}

func TestFindingsFromRowsAllowsMissingLearnMore(t *testing.T) {
	definition := assessment.RecommendationDefinition{ID: "rec", ResourceType: "Microsoft.Test/widgets"}
	rows := []FindingRow{{ID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"}}
	findings := FindingsFromRows(definition, rows, nil)
	if len(findings) != 1 || findings[0].LearnMoreURL != "" {
		t.Fatalf("missing learn-more link should be represented safely: %#v", findings)
	}
}
