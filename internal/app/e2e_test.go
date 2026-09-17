package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/advisor"
	"github.com/DeBoX85/Cloud-Assess/internal/arcsql"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/cost"
	"github.com/DeBoX85/Cloud-Assess/internal/defender"
	"github.com/DeBoX85/Cloud-Assess/internal/diagnostics"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
	"github.com/xuri/excelize/v2"
)

func TestEndToEndCoordinatorApplicationAndRenderers(t *testing.T) {
	resource := assessment.Resource{
		ID:             "/subscriptions/sub-1/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store",
		SubscriptionID: "sub-1",
		ResourceGroup:  "rg",
		Type:           "Microsoft.Storage/storageAccounts",
		Name:           "store",
		Location:       "norwayeast",
	}
	definition := assessment.RecommendationDefinition{
		ID:                  "storage-test-1",
		ResourceType:        resource.Type,
		Recommendation:      "Enable the test control",
		Category:            "Security",
		Impact:              assessment.ImpactHigh,
		Query:               "resources | where false",
		Source:              rules.SourceCustom,
		ValidationMechanism: rules.ValidationARG,
	}
	finding := assessment.Finding{
		RecommendationID:    definition.ID,
		ResourceID:          resource.ID,
		SubscriptionID:      resource.SubscriptionID,
		SubscriptionName:    "Test Subscription",
		ResourceGroup:       resource.ResourceGroup,
		ResourceName:        resource.Name,
		ResourceType:        resource.Type,
		Recommendation:      definition.Recommendation,
		Category:            definition.Category,
		Impact:              definition.Impact,
		Source:              definition.Source,
		ValidationMechanism: definition.ValidationMechanism,
	}

	catalog := rules.NewCatalog()
	catalog.Add(definition)
	operations := operationsForEndToEndTest(catalog, resource, finding)
	coordinator := orchestration.NewCoordinator(operations)
	stageConfig := stages.NewDefault()
	for _, name := range []string{stages.Diagnostics, stages.Advisor, stages.Defender, stages.DefenderRecommendations, stages.Policy, stages.Arc, stages.Cost, stages.Plugin} {
		if err := stageConfig.Set(name, false); err != nil {
			t.Fatal(err)
		}
	}

	base := filepath.Join(t.TempDir(), "cloud-assess-e2e")
	runner := NewRunner(coordinator)
	outcome, err := runner.Run(context.Background(), ScanOptions{
		Assessment: orchestration.Request{Subscriptions: []string{"sub-1"}, ScannerKeys: []string{"st"}, Stages: stageConfig},
		Outputs: OutputOptions{
			BaseName:              base,
			XLSX:                  true,
			JSON:                  true,
			RedactSubscriptionIDs: true,
			Version:               "test",
		},
	})
	if err != nil {
		t.Fatalf("end-to-end run failed: %v", err)
	}
	if outcome.ExitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want 0", outcome.ExitCode)
	}
	if len(outcome.Files) != 2 {
		t.Fatalf("generated files = %#v, want JSON + XLSX", outcome.Files)
	}

	jsonBytes, err := os.ReadFile(base + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		t.Fatalf("decode generated JSON: %v", err)
	}
	if payload["completeness"] != string(assessment.CompletenessComplete) {
		t.Fatalf("JSON completeness = %#v", payload["completeness"])
	}
	findings, ok := payload["findings"].([]any)
	if !ok || len(findings) != 1 {
		t.Fatalf("JSON findings = %#v, want one", payload["findings"])
	}

	workbook, err := excelize.OpenFile(base + ".xlsx")
	if err != nil {
		t.Fatalf("open generated XLSX: %v", err)
	}
	defer workbook.Close()
	sheets := workbook.GetSheetList()
	if len(sheets) < 2 || sheets[0] != "Assessment Status" || sheets[1] != "Recommendations" {
		t.Fatalf("sheet order = %#v", sheets)
	}
	if value, err := workbook.GetCellValue("Assessment Status", "A5"); err != nil || value != string(assessment.CompletenessComplete) {
		t.Fatalf("Assessment Status completeness = %q, err=%v", value, err)
	}
}

func operationsForEndToEndTest(catalog *rules.Catalog, resource assessment.Resource, finding assessment.Finding) orchestration.Operations {
	return orchestration.Operations{
		DiscoverSubscriptions: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{"sub-1": "Test Subscription"}, nil
		},
		DiscoverManagementGroups: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{"sub-1": "Test Subscription"}, nil
		},
		DiscoverResources: func(_ context.Context, _ map[string]string, filters *config.Filters) (*discovery.ResourceInventory, error) {
			filters.Assessment.SetResourceScope(resource.ID, true)
			return &discovery.ResourceInventory{Included: []assessment.Resource{resource}}, nil
		},
		LoadCatalog: func() (*rules.Catalog, error) { return catalog, nil },
		ExecuteGraph: func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
			return []assessment.Finding{finding}, nil, nil
		},
		ScanDiagnostics: func(context.Context, []assessment.Resource, *config.AssessmentFilter, map[string]string) (diagnostics.Result, error) {
			return diagnostics.Result{}, nil
		},
		ScanAdvisor: func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
			return advisor.Result{}, nil
		},
		ScanDefenderStatus: func(context.Context, map[string]string, *config.AssessmentFilter) (defender.StatusResult, error) {
			return defender.StatusResult{}, nil
		},
		ScanDefenderRecommendations: func(context.Context, map[string]string, *config.AssessmentFilter) (defender.RecommendationsResult, error) {
			return defender.RecommendationsResult{}, nil
		},
		ScanPolicy: func(context.Context, map[string]string, *config.AssessmentFilter) (policy.Result, error) {
			return policy.Result{}, nil
		},
		ScanArcSQL: func(context.Context, map[string]string, *config.AssessmentFilter) (arcsql.Result, error) {
			return arcsql.Result{}, nil
		},
		ScanCost: func(context.Context, map[string]string) (cost.Result, error) {
			return cost.Result{}, nil
		},
	}
}
