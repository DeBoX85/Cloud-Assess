package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	resource := assessment.Resource{
		ID:             "/subscriptions/" + subscriptionID + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store",
		SubscriptionID: subscriptionID,
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
	operations := operationsForEndToEndTest(t, catalog, resource, finding)
	coordinator := orchestration.NewCoordinator(operations)
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

	base := filepath.Join(t.TempDir(), "cloud-assess-e2e")
	runner := NewRunner(coordinator)
	var stdout bytes.Buffer
	runner.stdout = &stdout
	outcome, err := runner.Run(context.Background(), ScanOptions{
		Assessment: orchestration.Request{
			Subscriptions: []string{subscriptionID},
			ScannerKeys:   []string{"st"},
			Stages:        stageConfig,
		},
		Outputs: OutputOptions{
			BaseName:              base,
			XLSX:                  true,
			JSON:                  true,
			CSV:                   true,
			SARIF:                 true,
			Stdout:                true,
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
	var csvFiles int
	var maskedCSV bool
	for _, filename := range outcome.Files {
		if !strings.HasSuffix(filename, ".csv") {
			continue
		}
		csvFiles++
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.ToLower(string(data)), subscriptionID) {
			t.Fatalf("generated CSV %q leaked raw subscription ID", filename)
		}
		maskedCSV = maskedCSV || strings.Contains(string(data), redact.SubscriptionID(subscriptionID, true))
	}
	if csvFiles == 0 || !maskedCSV {
		t.Fatalf("generated files = %#v, want CSV tables containing the masked subscription ID", outcome.Files)
	}

	jsonBytes, err := os.ReadFile(base + ".json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(string(jsonBytes)), strings.ToLower(subscriptionID)) {
		t.Fatal("generated JSON leaked raw subscription ID")
	}
	masked := redact.SubscriptionID(subscriptionID, true)
	if strings.Contains(strings.ToLower(stdout.String()), subscriptionID) || !strings.Contains(stdout.String(), masked) {
		t.Fatal("stdout did not honor subscription ID redaction")
	}
	sarifBytes, err := os.ReadFile(base + ".sarif")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(string(sarifBytes)), strings.ToLower(resource.ID)) {
		t.Fatal("SARIF did not retain the resource identity needed for stable findings")
	}
	if !strings.Contains(string(jsonBytes), masked) {
		t.Fatalf("generated JSON does not contain masked subscription ID %q", masked)
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
	if value, err := workbook.GetCellValue("ImpactedResources", "H5"); err != nil || value != masked {
		t.Fatalf("ImpactedResources subscription = %q, err=%v, want %q", value, err, masked)
	}
	for _, sheet := range sheets {
		rows, err := workbook.GetRows(sheet)
		if err != nil {
			t.Fatal(err)
		}
		for _, row := range rows {
			for _, cell := range row {
				if strings.Contains(strings.ToLower(cell), subscriptionID) {
					t.Fatalf("XLSX sheet %q leaked raw subscription ID", sheet)
				}
			}
		}
	}
}

func TestCostPermissionFailurePersistsPartialReportWithHealthyFindings(t *testing.T) {
	const subscriptionID = "11111111-2222-3333-4444-555555555555"
	resource := assessment.Resource{
		ID:             "/subscriptions/" + subscriptionID + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store",
		SubscriptionID: subscriptionID, ResourceGroup: "rg", Type: "Microsoft.Storage/storageAccounts",
	}
	definition := assessment.RecommendationDefinition{
		ID: "storage-test-1", ResourceType: resource.Type, Query: "resources | where false", Source: rules.SourceCustom,
	}
	finding := assessment.Finding{RecommendationID: definition.ID, ResourceID: resource.ID, ResourceType: resource.Type}
	catalog := rules.NewCatalog()
	catalog.Add(definition)
	operations := operationsForEndToEndTest(t, catalog, resource, finding)
	operations.ScanCost = func(context.Context, map[string]string) (cost.Result, error) {
		return cost.Result{}, errors.New("cost access denied")
	}
	stageConfig := stages.NewDefault()
	for _, name := range []string{stages.Diagnostics, stages.Advisor, stages.Defender} {
		if err := stageConfig.Set(name, false); err != nil {
			t.Fatal(err)
		}
	}
	if err := stageConfig.Set(stages.Cost, true); err != nil {
		t.Fatal(err)
	}

	base := filepath.Join(t.TempDir(), "partial-cost")
	outcome, err := NewRunner(orchestration.NewCoordinator(operations)).Run(context.Background(), ScanOptions{
		Assessment: orchestration.Request{Subscriptions: []string{subscriptionID}, ScannerKeys: []string{"st"}, Stages: stageConfig},
		Outputs:    OutputOptions{BaseName: base, JSON: true},
	})
	if err == nil || outcome.ExitCode != ExitPartial {
		t.Fatalf("cost access error = %v, exit = %d; want partial exit", err, outcome.ExitCode)
	}
	data, err := os.ReadFile(base + ".json")
	if err != nil {
		t.Fatalf("partial report was not persisted: %v", err)
	}
	var report struct {
		Completeness string `json:"completeness"`
		Findings     []struct {
			RecommendationID string `json:"recommendationId"`
		} `json:"findings"`
		Stages []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
			Error  *struct {
				Code string `json:"code"`
			} `json:"error"`
		} `json:"stages"`
		Costs []json.RawMessage `json:"costs"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report.Completeness != string(assessment.CompletenessPartial) || len(report.Findings) != 1 || report.Findings[0].RecommendationID != definition.ID || len(report.Costs) != 0 {
		t.Fatalf("partial report lost healthy findings or fabricated costs: completeness=%s findings=%v costs=%d", report.Completeness, report.Findings, len(report.Costs))
	}
	for _, stage := range report.Stages {
		if stage.Name == stages.Cost {
			if stage.Status != string(assessment.StageFailed) || stage.Error == nil || stage.Error.Code != "stage_failed" {
				t.Fatalf("Cost access denial was not visible in report: %+v", stage)
			}
			return
		}
	}
	t.Fatal("partial report omitted Cost stage")
}

func TestGraphQueryFailurePersistsFailedReportAndStopsLaterStages(t *testing.T) {
	const subscriptionID = "11111111-2222-3333-4444-555555555555"
	resource := assessment.Resource{
		ID:             "/subscriptions/" + subscriptionID + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store",
		SubscriptionID: subscriptionID,
		ResourceGroup:  "rg",
		Type:           "Microsoft.Storage/storageAccounts",
	}
	catalog := rules.NewCatalog()
	catalog.Add(assessment.RecommendationDefinition{
		ID: "storage-test-1", ResourceType: resource.Type, Query: "resources | where false", Source: rules.SourceCustom,
	})
	operations := operationsForEndToEndTest(t, catalog, resource, assessment.Finding{})
	operations.ExecuteGraph = func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
		return nil, nil, errors.New("ARG query throttled")
	}
	advisorCalled := false
	operations.ScanAdvisor = func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
		advisorCalled = true
		return advisor.Result{}, nil
	}

	base := filepath.Join(t.TempDir(), "failed-graph")
	outcome, err := NewRunner(orchestration.NewCoordinator(operations)).Run(context.Background(), ScanOptions{
		Assessment: orchestration.Request{Subscriptions: []string{subscriptionID}, ScannerKeys: []string{"st"}, Stages: stages.NewDefault()},
		Outputs:    OutputOptions{BaseName: base, JSON: true},
	})
	if err == nil || !strings.Contains(err.Error(), "ARG query throttled") || outcome.ExitCode != ExitExecutionFail {
		t.Fatalf("Graph query error = %v, exit = %d; want execution failure", err, outcome.ExitCode)
	}
	if advisorCalled {
		t.Fatal("Advisor ran after critical Graph failure")
	}
	data, err := os.ReadFile(base + ".json")
	if err != nil {
		t.Fatalf("failed assessment report was not persisted: %v", err)
	}
	var report struct {
		Completeness string            `json:"completeness"`
		Resources    []json.RawMessage `json:"resources"`
		Findings     []json.RawMessage `json:"findings"`
		Stages       []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
			Error  *struct {
				Code string `json:"code"`
			} `json:"error"`
		} `json:"stages"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report.Completeness != string(assessment.CompletenessFailed) || len(report.Resources) != 1 || len(report.Findings) != 0 {
		t.Fatalf("failed report lost inventory or invented findings: completeness=%s resources=%d findings=%d", report.Completeness, len(report.Resources), len(report.Findings))
	}
	graphFailed, advisorSkipped := false, false
	for _, stage := range report.Stages {
		switch stage.Name {
		case stages.Graph:
			graphFailed = stage.Status == string(assessment.StageFailed) && stage.Error != nil && stage.Error.Code == "stage_failed"
		case stages.Advisor:
			advisorSkipped = stage.Status == string(assessment.StageSkipped)
		}
	}
	if !graphFailed || !advisorSkipped {
		t.Fatalf("Graph failure or skipped Advisor missing from persisted stage health: %+v", report.Stages)
	}
}

func operationsForEndToEndTest(t *testing.T, catalog *rules.Catalog, resource assessment.Resource, finding assessment.Finding) orchestration.Operations {
	return orchestration.Operations{
		DiscoverSubscriptions: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{resource.SubscriptionID: "Test Subscription"}, nil
		},
		DiscoverManagementGroups: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{resource.SubscriptionID: "Test Subscription"}, nil
		},
		DiscoverResources: func(_ context.Context, _ map[string]string, filters *config.Filters) (*discovery.ResourceInventory, error) {
			filters.Assessment.SetResourceScope(resource.ID, true)
			return &discovery.ResourceInventory{Included: []assessment.Resource{resource}}, nil
		},
		LoadCatalog: func() (*rules.Catalog, error) { return catalog, nil },
		ExecuteGraph: func(_ context.Context, definitions []assessment.RecommendationDefinition, _ map[string]string, _ *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
			found := false
			for _, definition := range definitions {
				if definition.ID == finding.RecommendationID {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("phase-two graph definitions did not include %s: %#v", finding.RecommendationID, definitions)
			}
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
