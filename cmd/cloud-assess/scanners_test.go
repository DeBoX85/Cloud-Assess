package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
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
)

func TestScannerCommandsMatchPinnedExecutable(t *testing.T) {
	data, err := os.ReadFile("testdata/scanner-keys-reference.json")
	if err != nil {
		t.Fatal(err)
	}
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand(nil)
	scan, _, err := root.Find([]string{"scan"})
	if err != nil {
		t.Fatal(err)
	}
	var actual []string
	for _, child := range scan.Commands() {
		actual = append(actual, child.Name())
	}
	if !reflect.DeepEqual(actual, keys) {
		t.Fatalf("scanner commands = %v, reference = %v", actual, keys)
	}
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			var got scanFlags
			root := newRootCommand(func(_ context.Context, f scanFlags) (int, error) { got = f; return 2, nil })
			root.SetArgs([]string{"scan", "--subscription-id", "sub", key, "--resource-group", "rg", "--json", "--xlsx=false", "--filters", "filters.yml"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.scannerKeys, []string{key}) || !reflect.DeepEqual(got.subscriptions, []string{"sub"}) || !reflect.DeepEqual(got.resourceGroups, []string{"rg"}) || !got.json || got.xlsx || got.filtersFile != "filters.yml" {
				t.Fatalf("mapped flags: %+v", got)
			}
			if *root.Context().Value(exitCodeContextKey{}).(*int) != 2 {
				t.Fatal("exit code lost")
			}
		})
	}
}

func TestScannerCommandRejectsExtraArguments(t *testing.T) {
	called := false
	root := newRootCommand(func(context.Context, scanFlags) (int, error) { called = true; return 0, nil })
	root.SetArgs([]string{"scan", "st", "extra"})
	if err := root.Execute(); err == nil || called {
		t.Fatalf("extra argument reached executor: called=%v", called)
	}
	root = newRootCommand(func(context.Context, scanFlags) (int, error) { return 2, errors.New("quality gate failed") })
	root.SetArgs([]string{"scan", "vm"})
	if err := root.Execute(); err == nil || *root.Context().Value(exitCodeContextKey{}).(*int) != 2 {
		t.Fatal("scanner error/exit lost")
	}
}

func TestScannerCommandThroughProductionCoordinator(t *testing.T) {
	// Fabricated source-contract fixture: single-key commands override resourceTypes
	// selection; generic scan respects that selection. No credentials or Azure requests.
	for _, tc := range []struct {
		name     string
		args     []string
		wantType string
	}{
		{"storage overrides vm filter", []string{"scan", "st"}, "Microsoft.Storage/storageAccounts"},
		{"generic retains vm filter", []string{"scan"}, "Microsoft.Compute/virtualMachines"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			filters := config.NewFilters()
			filters.Assessment.Include.ResourceTypes = []string{"vm"}
			filters.RebuildIndexes()
			stageConfig := stages.NewDefault()
			if err := stageConfig.Apply([]string{"-diagnostics", "-advisor", "-defender"}); err != nil {
				t.Fatal(err)
			}
			catalog := rules.NewCatalog()
			for _, resourceType := range []string{"Microsoft.Storage/storageAccounts", "Microsoft.Compute/virtualMachines"} {
				catalog.Add(assessment.RecommendationDefinition{ID: resourceType, ResourceType: resourceType, Query: "synthetic", Source: rules.SourceAPRL})
			}
			ops := orchestration.Operations{
				DiscoverManagementGroups: func(context.Context, []string, *config.Filters) (map[string]string, error) {
					t.Fatal("unexpected management-group discovery")
					return nil, nil
				},
				ScanDiagnostics: func(context.Context, []assessment.Resource, *config.AssessmentFilter, map[string]string) (diagnostics.Result, error) {
					t.Fatal("disabled diagnostics executed")
					return diagnostics.Result{}, nil
				},
				ScanAdvisor: func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
					t.Fatal("disabled advisor executed")
					return advisor.Result{}, nil
				},
				ScanDefenderStatus: func(context.Context, map[string]string, *config.AssessmentFilter) (defender.StatusResult, error) {
					t.Fatal("disabled defender executed")
					return defender.StatusResult{}, nil
				},
				ScanDefenderRecommendations: func(context.Context, map[string]string, *config.AssessmentFilter) (defender.RecommendationsResult, error) {
					t.Fatal("disabled recommendations executed")
					return defender.RecommendationsResult{}, nil
				},
				ScanPolicy: func(context.Context, map[string]string, *config.AssessmentFilter) (policy.Result, error) {
					t.Fatal("disabled policy executed")
					return policy.Result{}, nil
				},
				ScanArcSQL: func(context.Context, map[string]string, *config.AssessmentFilter) (arcsql.Result, error) {
					t.Fatal("disabled arc executed")
					return arcsql.Result{}, nil
				},
				ScanCost: func(context.Context, map[string]string) (cost.Result, error) {
					t.Fatal("disabled cost executed")
					return cost.Result{}, nil
				},

				DiscoverSubscriptions: func(context.Context, []string, *config.Filters) (map[string]string, error) {
					return map[string]string{"sub": "Subscription"}, nil
				},
				LoadCatalog: func() (*rules.Catalog, error) { return catalog, nil },
				DiscoverResources: func(_ context.Context, _ map[string]string, f *config.Filters) (*discovery.ResourceInventory, error) {
					inventory := &discovery.ResourceInventory{}
					for _, resourceType := range []string{"Microsoft.Storage/storageAccounts", "Microsoft.Compute/virtualMachines"} {
						resource := assessment.Resource{ID: "/subscriptions/sub/resourceGroups/rg/providers/" + resourceType + "/item", SubscriptionID: "sub", ResourceGroup: "rg", Type: resourceType}
						included := !f.Assessment.IsResourceTypeExcluded(resourceType)
						f.Assessment.SetResourceScope(resource.ID, included)
						if included {
							inventory.Included = append(inventory.Included, resource)
						} else {
							inventory.Excluded = append(inventory.Excluded, resource)
						}
					}
					return inventory, nil
				},
				ExecuteGraph: func(_ context.Context, definitions []assessment.RecommendationDefinition, _ map[string]string, _ *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
					if len(definitions) != 1 || definitions[0].ResourceType != tc.wantType {
						t.Fatalf("executed wrong scanner definitions: %v", definitions)
					}
					return nil, nil, nil
				},
			}
			executed := false
			root := newRootCommand(func(ctx context.Context, flags scanFlags) (int, error) {
				executed = true
				result, err := orchestration.NewCoordinator(ops).Run(ctx, assessmentRequest(flags, filters, stageConfig))
				if err != nil {
					return 1, err
				}
				if len(result.Resources) != 1 || result.Resources[0].Type != tc.wantType {
					t.Fatalf("wrong selected inventory: %v", result.Resources)
				}
				if len(result.OutOfScope) != 1 || result.OutOfScope[0].Type == tc.wantType {
					t.Fatalf("wrong excluded inventory: %v", result.OutOfScope)
				}
				return 0, nil
			})
			root.SetArgs(tc.args)
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if !executed {
				t.Fatal("executor not called")
			}
		})
	}
}

func TestScannerSelectionDoesNotLeakIntoGenericCommand(t *testing.T) {
	var selections [][]string
	root := newRootCommand(func(_ context.Context, flags scanFlags) (int, error) {
		selections = append(selections, append([]string(nil), flags.scannerKeys...))
		return 0, nil
	})
	root.SetArgs([]string{"scan", "st"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	root.SetArgs([]string{"scan"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if len(selections) != 2 || !reflect.DeepEqual(selections[0], []string{"st"}) || len(selections[1]) != 0 {
		t.Fatalf("scanner selection leaked: %v", selections)
	}
}
