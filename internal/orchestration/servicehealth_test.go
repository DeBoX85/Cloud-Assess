package orchestration

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/servicehealth"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
)

type serviceGraph func(context.Context, arg.Request) (*arg.Response, error)

func (f serviceGraph) Do(ctx context.Context, r arg.Request) (*arg.Response, error) { return f(ctx, r) }

const serviceVM = `{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"eastus","targetResourceType":"microsoft.compute/virtualmachines","percentageOfTimeWithoutEvents":99.5,"events":2,"affectedResources":3}`
const serviceStorage = `{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"westus","targetResourceType":"microsoft.storage/storageaccounts","percentageOfTimeWithoutEvents":100,"events":0,"affectedResources":0}`

func TestServiceCoordinatorRepeatedConcurrentIsolation(t *testing.T) {
	ops := zoneOperations(t)
	ops.DiscoverSubscriptions = func(_ context.Context, ids []string, _ *config.Filters) (map[string]string, error) {
		if len(ids) != 1 {
			return nil, errors.New("unexpected discovery scope")
		}
		return map[string]string{ids[0]: "selected"}, nil
	}
	scanner := servicehealth.NewWithTransport(serviceGraph(func(_ context.Context, r arg.Request) (*arg.Response, error) {
		if len(r.Subscriptions) != 1 {
			return nil, errors.New("unexpected scope")
		}
		rows := []json.RawMessage{}
		for _, row := range []string{serviceVM, serviceStorage} {
			rows = append(rows, json.RawMessage(strings.ReplaceAll(row, zoneSub, r.Subscriptions[0])))
		}
		return &arg.Response{Data: rows}, nil
	}))
	ops.ScanServiceHealth = func(ctx context.Context, subscriptions map[string]string, filter *config.AssessmentFilter) (assessment.PluginTable, error) {
		return scanner.Scan(ctx, subscriptions, filter)
	}
	coordinator := NewCoordinator(ops)
	run := func(sub, key, want string) {
		request := Request{Subscriptions: []string{sub}, ScannerKeys: []string{key}, InternalPlugins: []string{"service-health"}, PluginOnly: true}
		for i := 0; i < 3; i++ {
			got, err := coordinator.Run(context.Background(), request)
			if err != nil || got == nil || len(got.PluginTables) != 1 || len(got.PluginTables[0].Rows) != 1 || got.PluginTables[0].Rows[0].SubscriptionID != sub || got.PluginTables[0].Rows[0].Cells[2] != want {
				t.Errorf("isolated scope/type: %+v %v", got, err)
				return
			}
			got.PluginTables[0].Rows[0].Cells[1] = "mutated returned value"
		}
	}
	run(zoneSub, "vm", "microsoft.compute/virtualmachines")
	var workers sync.WaitGroup
	workers.Add(2)
	go func() { defer workers.Done(); run(zoneSub, "vm", "microsoft.compute/virtualmachines") }()
	go func() {
		defer workers.Done()
		run("22222222-2222-4222-8222-222222222222", "st", "microsoft.storage/storageaccounts")
	}()
	workers.Wait()
}

func TestServiceCoordinatorSourceTypeScopeAndCopiedInputs(t *testing.T) {
	for _, include := range []bool{false, true} {
		t.Run(map[bool]string{false: "scanner", true: "include precedence"}[include], func(t *testing.T) {
			ops := zoneOperations(t)
			s := servicehealth.NewWithTransport(serviceGraph(func(_ context.Context, r arg.Request) (*arg.Response, error) {
				if r.Query != servicehealth.Query || strings.Join(r.Subscriptions, ",") != zoneSub {
					t.Fatal("selected query scope")
				}
				return &arg.Response{Data: []json.RawMessage{json.RawMessage(serviceVM), json.RawMessage(serviceStorage)}}, nil
			}))
			ops.ScanServiceHealth = func(ctx context.Context, subs map[string]string, filter *config.AssessmentFilter) (assessment.PluginTable, error) {
				v, e := s.Scan(ctx, subs, filter)
				subs[zoneSub] = "mutated"
				filter.SetAllowedResourceTypes(nil)
				return v, e
			}
			f := config.NewFilters()
			f.Assessment.Include.ResourceGroups = []string{"/subscriptions/" + zoneSub + "/resourceGroups/not-the-aggregate-region"}
			f.Assessment.Include.Tags = map[string]string{"Environment": "dev"}
			if include {
				f.Assessment.Include.ResourceTypes = []string{"st"}
			}
			f.RebuildIndexes()
			keys := []string{"vm"}
			if include {
				keys = nil
			}
			r, err := NewCoordinator(ops).Run(context.Background(), Request{Subscriptions: []string{zoneSub}, InternalPlugins: []string{"zone-mapping", "service-health", "service-health"}, PluginOnly: true, ScannerKeys: keys, Filters: f})
			want := "microsoft.compute/virtualmachines"
			if include {
				want = "microsoft.storage/storageaccounts"
			}
			if err != nil || r.Completeness != assessment.CompletenessComplete || len(r.PluginTables) != 2 || len(r.PluginTables[0].Rows) != 1 || r.PluginTables[0].Rows[0].Cells[2] != want || r.PluginTables[1].Rows[0].Cells[0] != "Dev" {
				t.Fatalf("source type scope/copy: %+v %v", r, err)
			}
			if f.Assessment.Include.Tags["Environment"] != "dev" {
				t.Fatal("caller filter changed")
			}
		})
	}
}

func TestServiceDualPartialCriticalCancellationAndInvalid(t *testing.T) {
	for _, mode := range []string{"partial", "critical", "cancel", "invalid"} {
		t.Run(mode, func(t *testing.T) {
			ops := zoneOperations(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			zoneCalls := 0
			oldZone := ops.ScanZoneMapping
			ops.ScanZoneMapping = func(c context.Context, s map[string]string) (zone.Result, error) { zoneCalls++; return oldZone(c, s) }
			s := servicehealth.NewWithTransport(serviceGraph(func(context.Context, arg.Request) (*arg.Response, error) {
				return &arg.Response{Data: []json.RawMessage{json.RawMessage(serviceVM)}}, nil
			}))
			ops.ScanServiceHealth = func(c context.Context, subs map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
				v, _ := s.Scan(c, subs, f)
				switch mode {
				case "partial":
					return v, errors.New("provider-secret-canary")
				case "cancel":
					cancel()
					return v, context.Canceled
				case "invalid":
					v.Metadata.Name = "unreviewed"
				}
				return v, nil
			}
			if mode == "critical" {
				ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
					return nil, errors.New("scope failed")
				}
			}
			r, err := NewCoordinator(ops).Run(ctx, Request{Subscriptions: []string{zoneSub}, InternalPlugins: []string{"zone-mapping", "service-health"}, PluginOnly: true})
			if r == nil || len(r.PluginTables) != 2 || r.Completeness == assessment.CompletenessComplete {
				t.Fatalf("requested partial tables: %+v %v", r, err)
			}
			b, _ := json.Marshal(r)
			if strings.Contains(string(b), "provider-secret-canary") {
				t.Fatal("provider error leaked")
			}
			if mode == "partial" && (zoneCalls != 1 || len(r.PluginTables[0].Rows) != 1 || len(r.PluginTables[1].Rows) != 1) {
				t.Fatal("non-context failure lost another adapter/rows")
			}
			if mode == "invalid" && (zoneCalls != 1 || len(r.PluginTables[0].Rows) != 0 || r.PluginTables[0].Health.Error.Code != "service_health_output_invalid") {
				t.Fatal("invalid output false success")
			}
			if mode == "cancel" && (!errors.Is(err, context.Canceled) || zoneCalls != 0 || r.PluginTables[1].Health.Status != assessment.StageSkipped) {
				t.Fatal("cancelled pending table lost")
			}
			if mode == "critical" && (err == nil || zoneCalls != 0 || r.PluginTables[0].Health.Status != assessment.StageSkipped || r.PluginTables[1].Health.Status != assessment.StageSkipped) {
				t.Fatal("critical pending tables lost")
			}
		})
	}
	ops := noOpOperations()
	ops.ScanServiceHealth = nil
	calls := 0
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) { calls++; return nil, nil }
	if _, err := NewCoordinator(ops).Run(context.Background(), Request{InternalPlugins: []string{"service-health"}, PluginOnly: true}); err == nil || calls != 0 {
		t.Fatal("missing selected operation reached scope")
	}
}
