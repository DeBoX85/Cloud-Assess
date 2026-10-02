package orchestration

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

const zoneSub = "11111111-1111-4111-8111-111111111111"

type zoneGetter struct {
	t         *testing.T
	malformed bool
}

func (g zoneGetter) GetBounded(ctx context.Context, url string, limit int64) ([]byte, error) {
	if url != "https://management.example.test/subscriptions/"+zoneSub+"/locations?api-version=2022-12-01" || limit != zone.MaxPageBytes {
		g.t.Fatalf("request boundary: %s %d", url, limit)
	}
	if g.malformed {
		return []byte(`{"value":[null]}`), nil
	}
	return []byte(`{"value":[{"name":"westus","displayName":"West US","availabilityZoneMappings":[{"logicalZone":"1","physicalZone":"westus-az1"}]}]}`), nil
}

func zoneOperations(t *testing.T) Operations {
	t.Helper()
	ops := noOpOperations()
	// All regular operations are tripwires. Keep full validation contract configured.
	v := reflect.ValueOf(&ops).Elem()
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if strings.HasPrefix(typ.Field(i).Name, "DiscoverSubscriptions") || typ.Field(i).Name == "ScanZoneMapping" {
			continue
		}
		name := typ.Field(i).Name
		f := v.Field(i)
		ft := f.Type()
		f.Set(reflect.MakeFunc(ft, func([]reflect.Value) []reflect.Value {
			t.Fatalf("unexpected regular operation %s in plugin-only mode", name)
			out := make([]reflect.Value, ft.NumOut())
			for j := range out {
				out[j] = reflect.Zero(ft.Out(j))
			}
			return out
		}))
	}
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return map[string]string{zoneSub: "Dev"}, nil
	}
	scanner, err := zone.NewScanner("https://management.example.test", zoneGetter{t: t})
	if err != nil {
		t.Fatal(err)
	}
	ops.ScanZoneMapping = scanner.Scan
	return ops
}

func TestZoneCoordinatorOnlyRealAdapterAndCallerIsolation(t *testing.T) {
	ops := zoneOperations(t)
	cfg := stages.NewPluginOnly()
	names := []string{plugins.ZoneMapping, plugins.ZoneMapping}
	scopeMap := map[string]string{zoneSub: "Dev"}
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		_ = cfg.Set(stages.Graph, true)
		names[0] = "changed-after-prepare"
		return scopeMap, nil
	}
	scanner, err := zone.NewScanner("https://management.example.test", zoneGetter{t: t})
	if err != nil {
		t.Fatal(err)
	}
	ops.ScanZoneMapping = func(ctx context.Context, subs map[string]string) (zone.Result, error) {
		v, e := scanner.Scan(ctx, subs)
		subs[zoneSub] = "changed"
		return v, e
	}
	r, err := NewCoordinator(ops).Run(context.Background(), Request{Subscriptions: []string{zoneSub}, Stages: cfg, InternalPlugins: names, PluginOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if r.SchemaVersion != "1.1" || r.Completeness != assessment.CompletenessComplete || len(r.PluginTables) != 1 || len(r.PluginTables[0].Rows) != 1 || scopeMap[zoneSub] != "Dev" {
		t.Fatalf("result/ownership: %#v %v", r, scopeMap)
	}
	if !reflect.DeepEqual(r.PluginTables[0].Rows[0].Cells, []string{"Dev", "westus", "West US", "1", "westus-az1"}) {
		t.Fatal("real adapter cells changed")
	}
	for _, s := range r.Stages {
		if s.Name != "scope" && s.Name != "plugin" && s.Status != assessment.StageSkipped {
			t.Fatalf("unexpected execution: %#v", s)
		}
	}
}

func TestZoneCoordinatorPartialCriticalCancellationAndInvalid(t *testing.T) {
	for _, kind := range []string{"partial", "critical", "cancel", "invalid", "empty", "malformed"} {
		t.Run(kind, func(t *testing.T) {
			ops := zoneOperations(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "critical" {
				ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
					return nil, errors.New("scope denied")
				}
			} else if kind == "malformed" {
				scanner, _ := zone.NewScanner("https://management.example.test", zoneGetter{t: t, malformed: true})
				ops.ScanZoneMapping = scanner.Scan
			} else {
				ops.ScanZoneMapping = func(context.Context, map[string]string) (zone.Result, error) {
					switch kind {
					case "empty":
						return zone.Result{}, nil
					case "invalid":
						return zone.Result{Rows: []zone.Row{{SubscriptionID: "bad", SubscriptionName: "secret-canary"}}}, nil
					case "cancel":
						cancel()
						return zone.Result{Rows: []zone.Row{{SubscriptionID: zoneSub, SubscriptionName: "Dev"}}}, context.Canceled
					default:
						return zone.Result{Rows: []zone.Row{{SubscriptionID: zoneSub, SubscriptionName: "Dev"}}, Failures: []zone.Failure{{SubscriptionID: "22222222-2222-4222-8222-222222222222", Code: "zone_request_failed"}}}, nil
					}
				}
			}
			r, err := NewCoordinator(ops).Run(ctx, Request{Subscriptions: []string{zoneSub}, InternalPlugins: []string{plugins.ZoneMapping}, PluginOnly: true})
			if r == nil || len(r.PluginTables) != 1 {
				t.Fatalf("requested table missing: %v", err)
			}
			table := r.PluginTables[0]
			switch kind {
			case "empty":
				if err != nil || r.Completeness != assessment.CompletenessComplete || table.Health.Status != assessment.StageCompleted || len(table.Rows) != 0 {
					t.Fatal("empty mistaken for failure")
				}
			case "critical":
				if err == nil || r.Completeness != assessment.CompletenessFailed || table.Health.Status != assessment.StageSkipped || len(table.Columns) != 5 {
					t.Fatal("critical failure lost pending table")
				}
			case "cancel":
				if !errors.Is(err, context.Canceled) || r.Completeness != assessment.CompletenessFailed || len(table.Rows) != 1 {
					t.Fatal("cancellation/rows lost")
				}
			default:
				if err != nil || r.Completeness != assessment.CompletenessPartial || table.Health.Status != assessment.StageFailed {
					t.Fatalf("false complete: %s %v %#v", kind, err, r)
				}
			}
			if e := r.ValidatePluginExtension(); e != nil {
				t.Fatal(e)
			}
		})
	}
}

func TestZoneCoordinatorPreflightAndConcurrentRuns(t *testing.T) {
	unnamed := stages.NewPluginOnly()
	regions := stages.NewPluginOnly()
	_ = regions.ApplyParams([]string{"plugin.target-regions=westus"})
	for _, req := range []Request{{InternalPlugins: []string{"carbon-emissions"}}, {PluginOnly: true, Stages: unnamed}, {PluginOnly: true, InternalPlugins: []string{plugins.ZoneMapping}, Stages: stages.NewDefault()}, {InternalPlugins: []string{plugins.ZoneMapping}, Stages: regions}} {
		if _, err := prepareRequest(req); err == nil {
			t.Fatal("invalid plugin request accepted")
		}
	}
	base := stages.NewDefault()
	req := Request{Stages: base, InternalPlugins: []string{plugins.ZoneMapping}}
	p, err := prepareRequest(req)
	if err != nil || base.IsEnabled(stages.Plugin) || !p.stages.IsEnabled(stages.Plugin) || !p.stages.IsEnabled(stages.Graph) {
		t.Fatal("mixed selection or caller config changed")
	}
	req.PluginOnly = true
	req.Stages = nil
	coord := NewCoordinator(zoneOperations(t))
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := coord.Run(context.Background(), req)
			if e != nil || len(r.PluginTables) != 1 || len(r.PluginTables[0].Rows) != 1 {
				t.Errorf("isolated run failed: %v", e)
			}
		}()
	}
	wg.Wait()
}

func TestZoneCoordinatorMixedRetainsRegularExecution(t *testing.T) {
	ops := noOpOperations()
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return map[string]string{zoneSub: "Dev"}, nil
	}
	scanner, err := zone.NewScanner("https://management.example.test", zoneGetter{t: t})
	if err != nil {
		t.Fatal(err)
	}
	ops.ScanZoneMapping = scanner.Scan
	graphCalls := 0
	ops.ExecuteGraph = func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
		graphCalls++
		return nil, nil, nil
	}
	r, err := NewCoordinator(ops).Run(context.Background(), Request{Stages: graphOnlyStages(t), InternalPlugins: []string{plugins.ZoneMapping}})
	if err != nil || graphCalls != 1 || len(r.PluginTables) != 1 || len(r.PluginTables[0].Rows) != 1 || r.Completeness != assessment.CompletenessComplete {
		t.Fatalf("mixed regular execution lost: %v %d %#v", err, graphCalls, r)
	}
}

func TestZoneRequestOptionsAreOwnedAndUnconfiguredOperationIsPreflight(t *testing.T) {
	cfg := stages.NewPluginOnly()
	_ = cfg.ApplyParams([]string{"plugin.target-regions="})
	p, err := prepareRequest(Request{Stages: cfg, PluginOnly: true, InternalPlugins: []string{plugins.ZoneMapping}})
	if err != nil {
		t.Fatal(err)
	}
	_ = cfg.ApplyParams([]string{"plugin.target-regions=changed"})
	if p.stages.Options(stages.Plugin)["target-regions"] != "" {
		t.Fatal("caller option map leaked into prepared request")
	}
	ops := noOpOperations()
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		t.Fatal("discovery occurred before missing selected operation preflight")
		return nil, nil
	}
	if r, err := NewCoordinator(ops).Run(context.Background(), Request{InternalPlugins: []string{plugins.ZoneMapping}}); err == nil || r != nil {
		t.Fatal("missing selected operation accepted")
	}
}
