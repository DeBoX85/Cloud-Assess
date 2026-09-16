package excel

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
	"github.com/xuri/excelize/v2"
)

func excelFixture() *result.AssessmentResult {
	subscriptionID := "11111111-1111-1111-1111-111111111111"
	resourceID := "/subscriptions/" + subscriptionID + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	return result.Build(result.Input{
		GeneratedAt:  time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
		ScopeID:      "scope-xlsx",
		Completeness: assessment.CompletenessCompleteWithWarnings,
		Stages: []assessment.StageExecution{
			{Name: stages.Graph, Status: assessment.StageCompletedWithWarnings, Records: 1, Warnings: []assessment.AssessmentWarning{{Code: "sample_warning", Message: "warning"}}},
			{Name: stages.Advisor, Status: assessment.StageSkipped},
			{Name: stages.Defender, Status: assessment.StageCompleted, Records: 1},
			{Name: stages.Cost, Status: assessment.StageCompleted, Records: 1},
		},
		Recommendations: []assessment.RecommendationDefinition{
			{ID: "rec-1", Recommendation: "Secure storage", Category: "Security", Impact: "High", ResourceType: "microsoft.storage/storageaccounts", Source: "CUSTOM", LearnMore: []assessment.LearnMoreLink{{URL: "https://example.test/rec"}}},
		},
		Findings: []assessment.Finding{
			{RecommendationID: "rec-1", Recommendation: "Secure storage", Category: "Security", Impact: "High", ResourceType: "microsoft.storage/storageaccounts", Source: "CUSTOM", ValidationMechanism: "Azure Resource Graph", ResourceID: resourceID, SubscriptionID: subscriptionID, SubscriptionName: "Sub One", ResourceGroup: "rg", ResourceName: "st1", LearnMoreURL: "https://example.test/rec"},
		},
		Resources: []assessment.Resource{{ID: resourceID, SubscriptionID: subscriptionID, ResourceGroup: "rg", Location: "westeurope", Type: "microsoft.storage/storageaccounts", Name: "st1"}},
		ResourceTypes: []assessment.ResourceTypeCount{{SubscriptionID: subscriptionID, SubscriptionName: "Sub One", ResourceType: "microsoft.storage/storageaccounts", Count: 1}},
		Defender: []assessment.DefenderPlanStatus{{SubscriptionID: subscriptionID, SubscriptionName: "Sub One", Name: "VirtualMachines", Tier: "Standard"}},
		Costs: []assessment.CostRecord{{From: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), To: time.Date(2026, 8, 31, 23, 59, 59, 0, time.UTC), SubscriptionID: subscriptionID, SubscriptionName: "Sub One", ServiceName: "Storage", Value: "12.34", Currency: "NOK"}},
	})
}

func TestWriteFileCreatesExpectedWorkbook(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "assessment.xlsx")
	if err := WriteFile(excelFixture(), filename, tables.Options{RedactSubscriptionIDs: true}); err != nil {
		t.Fatal(err)
	}

	book, err := excelize.OpenFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer book.Close()

	wantSheets := []string{"Assessment Status", "Recommendations", "ImpactedResources", "ResourceTypes", "Inventory", "Defender", "OutOfScope", "Costs"}
	if got := book.GetSheetList(); !reflect.DeepEqual(got, wantSheets) {
		t.Fatalf("sheets = %#v, want %#v", got, wantSheets)
	}
	if got, _ := book.GetCellValue("Assessment Status", "A1"); got != "Azure Cloud Assessment" {
		t.Fatalf("report title = %q", got)
	}
	if got, _ := book.GetCellValue("Assessment Status", "A4"); got != "Assessment Completeness" {
		t.Fatalf("status header = %q", got)
	}
	if got, _ := book.GetCellValue("Recommendations", "A4"); got != "Implemented" {
		t.Fatalf("recommendations header = %q", got)
	}
	if formula, _ := book.GetCellFormula("Recommendations", "K5"); formula != `HYPERLINK("https://example.test/rec","https://example.test/rec")` {
		t.Fatalf("hyperlink formula = %q", formula)
	}
	if got, _ := book.GetCellValue("ImpactedResources", "H5"); got == "11111111-1111-1111-1111-111111111111" {
		t.Fatal("subscription ID was not redacted")
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("XLSX mode = %o, want 600", got)
	}
}

func TestWriteFileValidatesInputs(t *testing.T) {
	if err := WriteFile(nil, "x.xlsx", tables.Options{}); err == nil {
		t.Fatal("expected nil result error")
	}
	if err := WriteFile(excelFixture(), "", tables.Options{}); err == nil {
		t.Fatal("expected empty filename error")
	}
	if err := validateInventorySize(maxInventoryRows); err != nil {
		t.Fatalf("limit should be allowed: %v", err)
	}
	if err := validateInventorySize(maxInventoryRows + 1); err == nil {
		t.Fatal("expected XLSX inventory limit error")
	}
}

func TestHyperlinkFormulaEscapesQuotes(t *testing.T) {
	got := hyperlinkFormula(`https://example.test/?q="x"`)
	want := `HYPERLINK("https://example.test/?q=""x""","https://example.test/?q=""x""")`
	if got != want {
		t.Fatalf("formula = %q, want %q", got, want)
	}
}
