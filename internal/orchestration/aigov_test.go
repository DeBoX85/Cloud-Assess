package orchestration

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/aigov"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

func TestAICoordinatorRecordedTagScopeOwned(t *testing.T) {
	for _, key := range []string{azure.EnvAzureCloud, azure.EnvAzureAuthorityHost, azure.EnvAzureResourceManagerEndpoint, azure.EnvAzureResourceManagerAudience} {
		t.Setenv(key, "")
	}
	id := "/subscriptions/" + zoneSub + "/resourceGroups/fixture-rg/providers/Microsoft.CognitiveServices/accounts/fixture-ai"
	denied := strings.Replace(id, "fixture-ai", "excluded-ai", 1)
	unknown := strings.Replace(id, "fixture-ai", "unknown-ai", 1)
	scope := map[string]string{zoneSub: "Dev"}
	f := config.NewFilters()
	f.Assessment.Include.Tags = map[string]string{"Environment": "dev"}
	f.RebuildIndexes()
	ops := zoneOperations(t)
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) { return scope, nil }
	ops.DiscoverResources = func(_ context.Context, _ map[string]string, filters *config.Filters) (*discovery.ResourceInventory, error) {
		filters.Assessment.SetResourceScope(id, true)
		filters.Assessment.SetResourceScope(denied, false)
		return &discovery.ResourceInventory{Included: []assessment.Resource{{ID: id, SubscriptionID: zoneSub, ResourceGroup: "fixture-rg", Name: "fixture-ai", Type: "Microsoft.CognitiveServices/accounts"}}}, nil
	}
	ops.LoadCatalog = func() (*rules.Catalog, error) { return rules.NewCatalog(), nil }
	ops.ExecuteGraph = func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
		return nil, nil, nil
	}
	ops.ScanAIGovernance = func(ctx context.Context, subs map[string]string, filter *config.AssessmentFilter) (assessment.PluginTable, error) {
		if subs[zoneSub] != "Dev" || filter.IsServiceExcluded(strings.ToUpper(id)) || !filter.IsServiceExcluded(denied) || !filter.IsServiceExcluded(unknown) {
			t.Error("recorded included/excluded/unknown/inherited tag scope lost")
		}
		if !filter.IsServiceExcluded("/subscriptions/" + zoneSub + "/resourceGroups/fixture-rg/providers/Microsoft.Compute/virtualMachines/other") {
			t.Error("selected scanner scope widened")
		}
		filter.SetResourceScope(id, false)
		filter.Include.Tags["Environment"] = "changed"
		subs[zoneSub] = "changed"
		nested := filter.Clone()
		nested.SetResourceScope(denied, true)
		if !filter.IsServiceExcluded(denied) {
			t.Error("snapshot decisions share mutable map")
		}
		hour := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
		count := 4.0
		return aigov.Project(ctx, map[string]string{zoneSub: "Dev"}, []aigov.Account{{ID: id, SubscriptionID: zoneSub, ResourceGroup: "fixture-rg", Name: "fixture-ai", Kind: "OpenAI", SKU: "S0"}}, []aigov.Point{{ResourceID: id, Timestamp: &hour, Count: &count}}, map[string]aigov.DeploymentSet{id: {}}, nil)
	}
	cfg := stages.NewDefault()
	_ = cfg.Apply([]string{"-diagnostics,-advisor,-defender"})
	c := NewCoordinator(ops)
	run := func() {
		r, err := c.Run(context.Background(), Request{Subscriptions: []string{zoneSub}, Filters: f, Stages: cfg, ScannerKeys: []string{"aif"}, InternalPlugins: []string{"zone-mapping", "ai-gov", "ai-gov"}})
		if err != nil || r == nil || r.Completeness != assessment.CompletenessComplete || len(r.PluginTables) != 2 || len(r.PluginTables[0].Rows) != 1 || r.PluginTables[0].Rows[0].Cells[15] != "4" || r.PluginTables[1].Rows[0].Cells[0] != "Dev" {
			t.Errorf("owned AI/independent tables: %v", err)
			return
		}
		r.PluginTables[0].Rows[0].Cells[2] = "mutated"
		r.PluginTables[0].Columns[0] = "mutated"
	}
	run()
	run()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); run() }()
	}
	wg.Wait()
	if scope[zoneSub] != "Dev" || f.Assessment.Include.Tags["Environment"] != "dev" || !f.Assessment.IsServiceExcluded(id) {
		t.Fatal("AI changed caller scope/config/tag decisions")
	}
}

func TestAICoordinatorPreflightBeforeScopeAndCriticalHeaders(t *testing.T) {
	for _, key := range []string{azure.EnvAzureCloud, azure.EnvAzureAuthorityHost, azure.EnvAzureResourceManagerEndpoint, azure.EnvAzureResourceManagerAudience} {
		t.Setenv(key, "")
	}
	ops := zoneOperations(t)
	called := false
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		called = true
		return nil, errors.New("scope failed")
	}
	ops.ScanAIGovernance = nil
	request := Request{PluginOnly: true, InternalPlugins: []string{"ai-gov", "zone-mapping"}}
	if _, err := NewCoordinator(ops).Run(context.Background(), request); err == nil || called {
		t.Fatal("missing selected AI operation reached scope")
	}
	ops.ScanAIGovernance = func(context.Context, map[string]string, *config.AssessmentFilter) (assessment.PluginTable, error) {
		t.Fatal("critical scope reached AI")
		return aigov.PendingTable(), nil
	}
	t.Setenv(azure.EnvAzureCloud, "AzureChina")
	if _, err := NewCoordinator(ops).Run(context.Background(), request); err == nil || called {
		t.Fatal("unsupported AI cloud reached scope authentication")
	}
	t.Setenv(azure.EnvAzureCloud, "")
	r, err := NewCoordinator(ops).Run(context.Background(), request)
	if err == nil || r == nil || r.Completeness != assessment.CompletenessFailed || len(r.PluginTables) != 2 || r.PluginTables[0].Metadata.Name != "ai-gov" || r.PluginTables[0].Health.Status != assessment.StageSkipped || len(r.PluginTables[0].Columns) != 16 || r.PluginTables[0].SheetName != "AI Throttling" {
		t.Fatal("critical scope lost pending AI source headers")
	}
}
