package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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
	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
	"github.com/xuri/excelize/v2"
)

func TestEndToEndCoordinatorApplicationAndRenderers(t *testing.T) {
	const subscriptionID = "11111111-2222-3333-4444-555555555555"
	const resourceType = "microsoft.storage/storageaccounts"
	resourceID := "/subscriptions/" + subscriptionID + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"

	catalog := rules.NewCatalog()
	catalog.Add(assessment.RecommendationDefinition{
		ID:                  "test-001",
		Recommendation:      "Use secure storage settings",
		Category:            "Security",
		Impact:              assessment.ImpactHigh,
		ResourceType:        resourceType,
		Source:              rules.SourceCustom,
		Query:               "resources | where false",
		ValidationMechanism: rules.ValidationARG,
	})

	operations := orchestration.Operations{
		DiscoverSubscriptions: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{subscriptionID: "Test Subscription"}, nil
		},
		DiscoverManagementGroups: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{subscriptionID: "Test Subscription"}, nil
		},
		DiscoverResources: func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
			return &discovery.ResourceInventory{Included: []assessment.Resource{{
				ID:             resourceID,
				SubscriptionID: subscriptionID,
				ResourceGroup:  "rg",
				Location:       "norwayeast",
				Type:           resourceType,
				Name:           "st1",
			}}}, nil
		},
		LoadCatalog: func() (*rules.Catalog, error) {
			return catalog, nil
		},
		ExecuteGraph: func(_ context.Context, definitions []assessment.RecommendationDefinition, _ map[string]string, _ *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
			found := false
			for _, definition := range definitions {
				if definition.ID == "test-001" {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("phase-two graph definitions did not include test-001: %#v", definitions)
			}
			return []assessment.Finding{{
				RecommendationID:    "test-001",
				Source:              rules.SourceCustom,
				ValidationMechanism: rules.ValidationARG,
				Category:            "Security",
				Impact:              assessment.ImpactHigh,
				ResourceType:        resourceType,
				Recommendation:      "Use secure storage settings",
				ResourceID:          resourceID,
				SubscriptionID:      subscriptionID,
				SubscriptionName:    "Test Subscription",
				ResourceGroup:       "rg",
				ResourceName:        "st1",
			}}, nil, nil
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

	stageConfig := stages.NewDefault()
	for _, name := range []string{
		stages.Diagnostics,
		stages.Advisor,
		stages.Defender,
		stages.DefenderRecommendations,
		stages.Policy,
		stages.Arc,
		stages.Cost,
		stages.Plugin,
	} {
		if err := stageConfig.Set(name, false); err != nil {
			t.Fatal(err)
		}
	}

	base := filepath.Join(t.TempDir(), "assessment")
	runner := NewRunner(orchestration.NewCoordinator(operations))
	outcome, err := runner.Run(context.Background(), ScanOptions{
		Assessment: orchestration.Request{
			Filters: config.NewFilters(),
			Stages:  stageConfig,
		},
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
		t.Fatalf("exit code = %d, want %d", outcome.ExitCode, ExitSuccess)
	}
	if len(outcome.Files) != 2 {
		t.Fatalf("generated files = %#v, want XLSX and JSON", outcome.Files)
	}

	jsonContent, err := os.ReadFile(base + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(string(jsonContent)), strings.ToLower(subscriptionID)) {
		t.Fatal("end-to-end JSON output leaked raw subscription ID")
	}
	masked := redact.SubscriptionID(subscriptionID, true)
	if !strings.Contains(string(jsonContent), masked) {
		t.Fatalf("JSON output does not contain masked subscription ID %q", masked)
	}

	workbook, err := excelize.OpenFile(base + ".xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer workbook.Close()

	sheets := workbook.GetSheetList()
	if len(sheets) == 0 || sheets[0] != "Assessment Status" {
		t.Fatalf("first worksheet = %#v, want Assessment Status", sheets)
	}
	impactedFound := false
	for _, sheet := range sheets {
		if sheet == "ImpactedResources" {
			impactedFound = true
			break
		}
	}
	if !impactedFound {
		t.Fatalf("ImpactedResources sheet missing: %#v", sheets)
	}
	subscriptionCell, err := workbook.GetCellValue("ImpactedResources", "H5")
	if err != nil {
		t.Fatal(err)
	}
	if subscriptionCell != masked {
		t.Fatalf("ImpactedResources H5 = %q, want masked subscription ID %q", subscriptionCell, masked)
	}
}
