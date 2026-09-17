package orchestration

import (
	"context"
	"errors"
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
	"github.com/DeBoX85/Cloud-Assess/internal/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

func TestPrepareRequestPreservesScopeAndScannerSelectionSemantics(t *testing.T) {
	filters := config.NewFilters()
	filters.Assessment.Include.ResourceTypes = []string{"st"}
	filters.RebuildIndexes()

	prepared, err := prepareRequest(Request{
		Subscriptions:  []string{"sub-1"},
		ResourceGroups: []string{"rg-one"},
		Filters:        filters,
	})
	if err != nil {
		t.Fatalf("prepareRequest returned error: %v", err)
	}
	if len(prepared.selectedKeys) != 1 || prepared.selectedKeys[0] != "st" {
		t.Fatalf("selected keys = %#v, want [st]", prepared.selectedKeys)
	}
	if prepared.filters.Assessment.IsSubscriptionExcluded("sub-1") {
		t.Fatal("explicit subscription should be included")
	}
	if !prepared.filters.Assessment.IsSubscriptionExcluded("sub-2") {
		t.Fatal("other subscription should be outside explicit scan scope")
	}
	if prepared.filters.Assessment.IsResourceGroupExcluded("/subscriptions/sub-1/resourceGroups/rg-one") {
		t.Fatal("requested resource group should be included")
	}
	if !prepared.filters.Assessment.IsResourceGroupExcluded("/subscriptions/sub-1/resourceGroups/rg-two") {
		t.Fatal("other resource group should be outside requested RG scope")
	}
	if prepared.filters.Assessment.IsResourceTypeExcluded("Microsoft.Storage/storageAccounts") {
		t.Fatal("storage resource type should be allowed by selected st scanner")
	}
	if !prepared.filters.Assessment.IsResourceTypeExcluded("Microsoft.Compute/virtualMachines") {
		t.Fatal("VM resource type should not be allowed by selected st scanner")
	}
}

func TestPrepareRequestValidatesReferenceScopeRules(t *testing.T) {
	tests := []struct {
		name    string
		request Request
		want    string
	}{
		{
			name: "management group mixed with subscription",
			request: Request{
				ManagementGroups: []string{"mg"},
				Subscriptions:    []string{"sub"},
			},
			want: "management group name cannot be used",
		},
		{
			name: "resource group without subscription",
			request: Request{
				ResourceGroups: []string{"rg"},
			},
			want: "resource group name can only be used with a subscription ID",
		},
		{
			name: "resource group with multiple subscriptions",
			request: Request{
				Subscriptions:  []string{"one", "two"},
				ResourceGroups: []string{"rg"},
			},
			want: "resource group name can only be used with one subscription ID",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := prepareRequest(test.request)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestCoordinatorTwoPhaseGraphCatalogAndPruning(t *testing.T) {
	catalog := rules.NewCatalog()
	catalog.Add(assessment.RecommendationDefinition{ID: "st-1", ResourceType: "Microsoft.Storage/storageAccounts", Query: "storage", Source: rules.SourceAPRL})
	catalog.Add(assessment.RecommendationDefinition{ID: "redis-1", ResourceType: "Microsoft.Cache/Redis", Query: "redis", Source: rules.SourceAPRL})
	catalog.Add(assessment.RecommendationDefinition{ID: "redis-enterprise-1", ResourceType: "Microsoft.Cache/redisEnterprise", Query: "redisEnterprise", Source: rules.SourceAPRL})
	catalog.Add(assessment.RecommendationDefinition{ID: "generic-1", ResourceType: "Microsoft.Resources", Query: "generic", Source: rules.SourceCustom})

	var executed []assessment.RecommendationDefinition
	operations := noOpOperations()
	operations.LoadCatalog = func() (*rules.Catalog, error) { return catalog, nil }
	operations.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return map[string]string{"sub-1": "Subscription One"}, nil
	}
	operations.DiscoverResources = func(_ context.Context, _ map[string]string, filters *config.Filters) (*discovery.ResourceInventory, error) {
		storage := assessment.Resource{
			ID:             "/subscriptions/sub-1/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store",
			SubscriptionID: "sub-1",
			ResourceGroup:  "rg",
			Type:           "Microsoft.Storage/storageAccounts",
			Name:           "store",
		}
		if filters.Assessment.IsResourceExcluded(storage.ID, nil) {
			t.Fatal("storage should be inside selected scanner scope")
		}
		filters.Assessment.SetResourceScope(storage.ID, true)
		return &discovery.ResourceInventory{Included: []assessment.Resource{storage}}, nil
	}
	operations.ExecuteGraph = func(_ context.Context, definitions []assessment.RecommendationDefinition, _ map[string]string, _ *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
		executed = append([]assessment.RecommendationDefinition(nil), definitions...)
		return nil, nil, nil
	}

	stageConfig := graphOnlyStages(t)
	filters := config.NewFilters()
	filters.Assessment.Include.ResourceTypes = []string{"st", "redis"}
	filters.RebuildIndexes()

	assessmentResult, err := NewCoordinator(operations).Run(context.Background(), Request{Filters: filters, Stages: stageConfig})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	gotCatalog := recommendationIDs(assessmentResult.Recommendations)
	for _, id := range []string{"st-1", "redis-1", "redis-enterprise-1"} {
		if !gotCatalog[id] {
			t.Fatalf("phase-one catalog missing %s: %#v", id, gotCatalog)
		}
	}
	if gotCatalog["generic-1"] {
		t.Fatal("phase-one catalog should not gain generic resource recommendations when filter selected only st,redis")
	}

	gotExecuted := recommendationIDs(executed)
	if !gotExecuted["st-1"] || !gotExecuted["generic-1"] {
		t.Fatalf("phase-two definitions = %#v, want storage + generic", gotExecuted)
	}
	if gotExecuted["redis-1"] || gotExecuted["redis-enterprise-1"] {
		t.Fatalf("undeployed Redis definitions should be pruned: %#v", gotExecuted)
	}
}

func TestCoordinatorOptionalFailureProducesPartialAndContinues(t *testing.T) {
	operations := noOpOperations()
	advisorErr := errors.New("advisor unavailable")
	costCalled := false
	operations.ScanAdvisor = func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
		return advisor.Result{}, advisorErr
	}
	operations.ScanCost = func(context.Context, map[string]string) (cost.Result, error) {
		costCalled = true
		return cost.Result{Records: []assessment.CostRecord{{SubscriptionID: "sub-1", ServiceName: "Storage", Value: "1", Currency: "USD"}}}, nil
	}

	stageConfig := graphOnlyStages(t)
	if err := stageConfig.Set(stages.Advisor, true); err != nil {
		t.Fatal(err)
	}
	if err := stageConfig.Set(stages.Cost, true); err != nil {
		t.Fatal(err)
	}
	assessmentResult, err := NewCoordinator(operations).Run(context.Background(), Request{Stages: stageConfig})
	if err != nil {
		t.Fatalf("optional stage failure should not return execution error: %v", err)
	}
	if assessmentResult.Completeness != assessment.CompletenessPartial {
		t.Fatalf("completeness = %s, want partial", assessmentResult.Completeness)
	}
	if !costCalled {
		t.Fatal("later optional stage did not execute after Advisor failure")
	}
	if len(assessmentResult.Costs) != 1 {
		t.Fatalf("cost records = %d, want 1", len(assessmentResult.Costs))
	}
	if statusForStage(assessmentResult.Stages, stages.Advisor) != assessment.StageFailed {
		t.Fatal("Advisor stage was not recorded failed")
	}
	if statusForStage(assessmentResult.Stages, stages.Cost) != assessment.StageCompleted {
		t.Fatal("Cost stage was not recorded completed")
	}
}

func TestCoordinatorCriticalGraphFailureStopsLaterStages(t *testing.T) {
	operations := noOpOperations()
	operations.ExecuteGraph = func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
		return nil, nil, errors.New("graph failed")
	}
	advisorCalled := false
	operations.ScanAdvisor = func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
		advisorCalled = true
		return advisor.Result{}, nil
	}

	stageConfig := graphOnlyStages(t)
	if err := stageConfig.Set(stages.Advisor, true); err != nil {
		t.Fatal(err)
	}
	assessmentResult, err := NewCoordinator(operations).Run(context.Background(), Request{Stages: stageConfig})
	if err == nil {
		t.Fatal("critical Graph failure should return error")
	}
	if assessmentResult == nil || assessmentResult.Completeness != assessment.CompletenessFailed {
		t.Fatalf("result completeness = %#v, want failed", assessmentResult)
	}
	if advisorCalled {
		t.Fatal("Advisor should not run after critical Graph failure")
	}
	if statusForStage(assessmentResult.Stages, stages.Graph) != assessment.StageFailed {
		t.Fatal("Graph stage was not recorded failed")
	}
	if statusForStage(assessmentResult.Stages, stages.Advisor) != assessment.StageSkipped {
		t.Fatal("Advisor stage should be recorded skipped after critical failure")
	}
}

func TestCoordinatorWarningsProduceCompleteWithWarnings(t *testing.T) {
	operations := noOpOperations()
	operations.ExecuteGraph = func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
		return nil, []arg.RuleWarning{{RecommendationID: "rec", ResourceType: "Microsoft.Resources", Code: "unsupported_logical_table", Message: "unsupported"}}, nil
	}
	assessmentResult, err := NewCoordinator(operations).Run(context.Background(), Request{Stages: graphOnlyStages(t)})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if assessmentResult.Completeness != assessment.CompletenessCompleteWithWarnings {
		t.Fatalf("completeness = %s, want complete_with_warnings", assessmentResult.Completeness)
	}
	if statusForStage(assessmentResult.Stages, stages.Graph) != assessment.StageCompletedWithWarnings {
		t.Fatal("Graph stage should be completed_with_warnings")
	}
}

func TestCoordinatorDiagnosticsDefinitionsJoinPrimaryCatalog(t *testing.T) {
	operations := noOpOperations()
	operations.ScanDiagnostics = func(context.Context, []assessment.Resource, *config.AssessmentFilter, map[string]string) (diagnostics.Result, error) {
		definition := assessment.RecommendationDefinition{ID: "diag-1", ResourceType: "Microsoft.Storage/storageAccounts", Source: rules.SourceDiagnostics}
		return diagnostics.Result{
			Recommendations: []assessment.RecommendationDefinition{definition},
			Findings:        []assessment.Finding{{RecommendationID: "diag-1", ResourceID: "/subscriptions/sub-1/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store", ResourceType: definition.ResourceType, Source: rules.SourceDiagnostics}},
		}, nil
	}
	stageConfig := graphOnlyStages(t)
	if err := stageConfig.Set(stages.Diagnostics, true); err != nil {
		t.Fatal(err)
	}
	assessmentResult, err := NewCoordinator(operations).Run(context.Background(), Request{Stages: stageConfig})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !recommendationIDs(assessmentResult.Recommendations)["diag-1"] {
		t.Fatal("Diagnostics definition missing from canonical recommendation catalog")
	}
	if len(assessmentResult.Findings) != 1 || assessmentResult.Findings[0].RecommendationID != "diag-1" {
		t.Fatalf("Diagnostics finding missing: %#v", assessmentResult.Findings)
	}
}

func TestCoordinatorPreservesEmptyResolvedScope(t *testing.T) {
	operations := noOpOperations()
	operations.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return map[string]string{}, nil
	}
	assessmentResult, err := NewCoordinator(operations).Run(context.Background(), Request{Stages: graphOnlyStages(t)})
	if err != nil {
		t.Fatalf("empty source scope is not a validation error in the pinned reference: %v", err)
	}
	if assessmentResult.ScopeID == "" {
		t.Fatal("empty resolved scope should still produce deterministic scope ID")
	}
}

func noOpOperations() Operations {
	return Operations{
		DiscoverSubscriptions: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{"sub-1": "Subscription One"}, nil
		},
		DiscoverManagementGroups: func(context.Context, []string, *config.Filters) (map[string]string, error) {
			return map[string]string{"sub-1": "Subscription One"}, nil
		},
		DiscoverResources: func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
			return &discovery.ResourceInventory{}, nil
		},
		LoadCatalog: func() (*rules.Catalog, error) { return rules.NewCatalog(), nil },
		ExecuteGraph: func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
			return nil, nil, nil
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

func graphOnlyStages(t *testing.T) *stages.Config {
	t.Helper()
	stageConfig := stages.NewDefault()
	for _, name := range []string{stages.Diagnostics, stages.Advisor, stages.Defender, stages.DefenderRecommendations, stages.Policy, stages.Arc, stages.Cost, stages.Plugin} {
		if err := stageConfig.Set(name, false); err != nil {
			t.Fatal(err)
		}
	}
	return stageConfig
}

func recommendationIDs(values []assessment.RecommendationDefinition) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value.ID] = true
	}
	return out
}

func statusForStage(executions []assessment.StageExecution, name string) assessment.StageStatus {
	for _, execution := range executions {
		if execution.Name == name {
			return execution.Status
		}
	}
	return ""
}
