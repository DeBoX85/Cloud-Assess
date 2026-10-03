package orchestration

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/carbon"
)

func TestCarbonCoordinatorTypeFiltersOwnedScopeRepeatedConcurrent(t *testing.T) {
	scope := map[string]string{zoneSub: "Dev"}
	for _, mode := range []string{"all", "scanner", "include"} {
		t.Run(mode, func(t *testing.T) {
			ops := zoneOperations(t)
			ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) { return scope, nil }
			ops.ScanCarbon = func(ctx context.Context, subs map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
				if subs[zoneSub] != "Dev" {
					t.Error("caller scope changed")
				}
				a, b, c := 100.0, 100.0, -20.0
				d := 75.0
				v, e := carbon.Project(ctx, "2026-01-01", "2026-03-31", []carbon.Item{{ResourceType: "Microsoft.Compute/virtualMachines", Latest: &a, Previous: &b, Change: &c}, {ResourceType: "Microsoft.Storage/storageAccounts", Latest: &d}}, f)
				subs[zoneSub] = "changed"
				f.Include.Tags["Environment"] = "changed"
				return v, e
			}
			f := config.NewFilters()
			f.Assessment.Include.Tags = map[string]string{"Environment": "unrelated"}
			f.Assessment.Include.ResourceGroups = []string{"/subscriptions/" + zoneSub + "/resourceGroups/unrelated"}
			f.Assessment.Exclude.Recommendations = []string{"unrelated"}
			keys := []string(nil)
			want := "Microsoft.Compute/virtualMachines,Microsoft.Storage/storageAccounts"
			switch mode {
			case "scanner":
				keys = []string{"vm"}
				want = "Microsoft.Compute/virtualMachines"
			case "include":
				f.Assessment.Include.ResourceTypes = []string{"vm"}
				want = "Microsoft.Compute/virtualMachines"
			}
			f.RebuildIndexes()
			c := NewCoordinator(ops)
			run := func() {
				got, e := c.Run(context.Background(), Request{Subscriptions: []string{zoneSub}, InternalPlugins: []string{"zone-mapping", "carbon-emissions", "carbon-emissions"}, PluginOnly: true, Filters: f, ScannerKeys: keys})
				if e != nil || got == nil || got.Completeness != assessment.CompletenessComplete || len(got.PluginTables) != 2 {
					t.Errorf("carbon type-filter run: %v", e)
					return
				}
				names := []string{}
				for _, r := range got.PluginTables[0].Rows {
					names = append(names, r.Cells[2])
					if r.SubscriptionID != "" {
						t.Error("aggregate identity invented")
					}
				}
				if strings.Join(names, ",") != want || got.PluginTables[1].Rows[0].Cells[0] != "Dev" {
					t.Error("source type filters/copied scope")
				}
				got.PluginTables[0].Columns[0] = "changed"
				got.PluginTables[0].Rows[0].Cells[3] = "changed"
			}
			run()
			run()
			var wg sync.WaitGroup
			wg.Add(2)
			for i := 0; i < 2; i++ {
				go func() { defer wg.Done(); run() }()
			}
			wg.Wait()
			if scope[zoneSub] != "Dev" || f.Assessment.Include.Tags["Environment"] != "unrelated" {
				t.Fatal("caller ownership lost")
			}
		})
	}
}

func TestCarbonCoordinatorCriticalAndMissingOperation(t *testing.T) {
	ops := zoneOperations(t)
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return nil, errors.New("scope failed")
	}
	got, e := NewCoordinator(ops).Run(context.Background(), Request{InternalPlugins: []string{"zone-mapping", "sql-eol", "carbon-emissions", "service-health"}, PluginOnly: true})
	if e == nil || got == nil || len(got.PluginTables) != 4 {
		t.Fatal("critical requested tables lost")
	}
	for i, name := range []string{"carbon-emissions", "service-health", "sql-eol", "zone-mapping"} {
		if got.PluginTables[i].Metadata.Name != name || got.PluginTables[i].Health.Status != assessment.StageSkipped {
			t.Fatal("critical pending metadata/order lost")
		}
	}
	b, _ := json.Marshal(got)
	if !strings.Contains(string(b), "Period From") {
		t.Fatal("pending carbon headers lost")
	}
	ops.ScanCarbon = nil
	called := false
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		called = true
		return nil, nil
	}
	if _, e := NewCoordinator(ops).Run(context.Background(), Request{InternalPlugins: []string{"carbon-emissions"}, PluginOnly: true}); e == nil || called {
		t.Fatal("unconfigured selected carbon operation reached scope")
	}
}
