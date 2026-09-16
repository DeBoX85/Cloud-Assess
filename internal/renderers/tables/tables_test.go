package tables

import (
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/skus"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

const testSubscriptionID = "11111111-1111-1111-1111-111111111111"

func tableFixture() *result.AssessmentResult {
	storageID := "/subscriptions/" + testSubscriptionID + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	vmssID := "/subscriptions/" + testSubscriptionID + "/resourceGroups/rg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1"
	return result.Build(result.Input{
		GeneratedAt:  time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
		ScopeID:      "scope-1",
		Completeness: assessment.CompletenessCompleteWithWarnings,
		Stages: []assessment.StageExecution{
			{Name: stages.Graph, Status: assessment.StageCompletedWithWarnings, Records: 2, Warnings: []assessment.AssessmentWarning{{Code: "graph_warning", Message: "warning"}}},
			{Name: stages.Advisor, Status: assessment.StageSkipped},
			{Name: stages.Defender, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: "stage_failed", Message: "denied"}},
		},
		Recommendations: []assessment.RecommendationDefinition{
			{ID: "rec-1", Recommendation: "Secure storage", Category: "Security", Impact: "High", ResourceType: "microsoft.storage/storageaccounts", Source: "CUSTOM", LongDescription: "guidance", LearnMore: []assessment.LearnMoreLink{{URL: "https://example.test/rec"}}},
			{ID: "rec-na", Recommendation: "AKS rule", Category: "Reliability", Impact: "Medium", ResourceType: "microsoft.containerservice/managedclusters", Source: "APRL"},
		},
		Findings: []assessment.Finding{
			{RecommendationID: "rec-1", Source: "DIAGNOSTICS", ValidationMechanism: "Azure Resource Manager", Category: "Security", Impact: "High", ResourceType: "microsoft.storage/storageaccounts", Recommendation: "Secure storage", ResourceID: storageID, SubscriptionID: testSubscriptionID, SubscriptionName: "Sub One", ResourceGroup: "rg", ResourceName: "st1", Parameters: []string{"p1"}, LearnMoreURL: "https://example.test/rec"},
			{RecommendationID: "sla-1", Source: "CUSTOM", Category: assessment.CategorySLA, ResourceID: vmssID, SubscriptionID: testSubscriptionID, Parameters: []string{"99.95%"}},
		},
		Resources: []assessment.Resource{
			{ID: storageID, SubscriptionID: testSubscriptionID, ResourceGroup: "rg", Location: "westeurope", Type: "microsoft.storage/storageaccounts", Name: "st1"},
			{ID: vmssID, SubscriptionID: testSubscriptionID, ResourceGroup: "rg", Location: "westeurope", Type: "microsoft.compute/virtualmachinescalesets", Name: "vmss1", SKUName: "Standard_D2s_v5", SKUCapacity: 3},
		},
		ResourceTypes: []assessment.ResourceTypeCount{{SubscriptionID: testSubscriptionID, SubscriptionName: "Sub One", ResourceType: "microsoft.storage/storageaccounts", Count: 1}, {SubscriptionID: testSubscriptionID, SubscriptionName: "Sub One", ResourceType: "microsoft.compute/virtualmachinescalesets", Count: 1}},
	})
}

func TestBuildOrderAndAssessmentStatus(t *testing.T) {
	tables, err := Build(tableFixture(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Assessment Status", "Recommendations", "ImpactedResources", "ResourceTypes", "Inventory", "Advisor", "Azure Policy", "Arc SQL", "DefenderRecommendations", "Defender", "OutOfScope", "Costs"}
	if len(tables) != len(want) {
		t.Fatalf("got %d tables, want %d", len(tables), len(want))
	}
	for i := range want {
		if tables[i].SheetName != want[i] {
			t.Fatalf("table %d = %q, want %q", i, tables[i].SheetName, want[i])
		}
	}
	status := tables[0].Rows
	if status[1][0] != string(assessment.CompletenessCompleteWithWarnings) || status[1][4] != stages.Graph || status[1][8] != "graph_warning" {
		t.Fatalf("unexpected status row: %#v", status[1])
	}
}

func TestRecommendationsPreserveApplicabilityMeaning(t *testing.T) {
	tables, err := Build(tableFixture(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	rows := tables[1].Rows
	states := map[string]string{}
	for _, row := range rows[1:] {
		states[row[11]] = row[0]
	}
	if states["rec-1"] != "false" {
		t.Fatalf("rec-1 implemented = %q, want false", states["rec-1"])
	}
	if states["rec-na"] != "N/A" {
		t.Fatalf("rec-na implemented = %q, want N/A", states["rec-na"])
	}
}

func TestImpactedUsesRealValidationMechanismAndRedaction(t *testing.T) {
	tables, err := Build(tableFixture(), Options{RedactSubscriptionIDs: true})
	if err != nil {
		t.Fatal(err)
	}
	row := tables[2].Rows[1]
	if row[0] != "Azure Resource Manager" {
		t.Fatalf("validated using = %q", row[0])
	}
	if strings.Contains(row[7], "11111111-1111") || strings.Contains(row[11], "/subscriptions/11111111-1111") {
		t.Fatalf("subscription identifiers were not redacted: %#v", row)
	}
	if row[12] != "p1" || row[13] != "" || row[17] != "https://example.test/rec" {
		t.Fatalf("unexpected parameter/link projection: %#v", row)
	}
}

func TestInventoryPreservesSLAAndCapacityProjection(t *testing.T) {
	tables, err := Build(tableFixture(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	rows := tables[4].Rows
	var vmss []string
	for _, row := range rows[1:] {
		if row[4] == "vmss1" {
			vmss = row
			break
		}
	}
	if vmss == nil {
		t.Fatal("VMSS row not found")
	}
	if vmss[7] != skus.ComputeCapacity("Standard_D2s_v5", 3, true) {
		t.Fatalf("capacity = %q", vmss[7])
	}
	if vmss[9] != "99.95%" {
		t.Fatalf("SLA = %q, want 99.95%%", vmss[9])
	}
}

func TestShouldRenderStageGating(t *testing.T) {
	data := tableFixture()
	tables, err := Build(data, Options{})
	if err != nil {
		t.Fatal(err)
	}
	byKey := map[string]Table{}
	for _, table := range tables {
		byKey[table.Key] = table
	}
	if ShouldRender(data, byKey["advisor"]) {
		t.Fatal("disabled/skipped advisor table should not render")
	}
	if !ShouldRender(data, byKey["defender"]) {
		t.Fatal("failed requested defender table should render")
	}
	if !ShouldRender(data, byKey["assessmentStatus"]) || !ShouldRender(data, byKey["inventory"]) {
		t.Fatal("status and graph inventory should render")
	}
}
