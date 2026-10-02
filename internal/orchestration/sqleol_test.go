package orchestration

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/servicehealth"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/sqleol"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
)

func realSQLScanner(t *testing.T) *sqleol.Scanner {
	t.Helper()
	b, err := os.ReadFile("../plugins/sqleol/testdata/source-input.json")
	if err != nil {
		t.Fatal(err)
	}
	var rows []json.RawMessage
	if json.Unmarshal(b, &rows) != nil {
		t.Fatal("source input")
	}
	return sqleol.NewWithTransport(serviceGraph(func(_ context.Context, request arg.Request) (*arg.Response, error) {
		if request.Query != sqleol.Query || strings.Join(request.Subscriptions, ",") != zoneSub {
			t.Error("selected SQL scope/query")
		}
		return &arg.Response{Data: []json.RawMessage{rows[1]}}, nil
	}))
}

func TestSQLCoordinatorSourceFiltersCopiedRepeatedConcurrent(t *testing.T) {
	ops := zoneOperations(t)
	s := realSQLScanner(t)
	ops.ScanSQLEOL = func(ctx context.Context, subs map[string]string, filter *config.AssessmentFilter) (assessment.PluginTable, error) {
		v, e := s.Scan(ctx, subs, filter)
		subs[zoneSub] = "mutated"
		filter.Include.Tags["Environment"] = "mutated"
		return v, e
	}
	f := config.NewFilters()
	f.Assessment.Include.Tags = map[string]string{"Environment": "unrelated"}
	f.Assessment.Include.ResourceGroups = []string{"/subscriptions/" + zoneSub + "/resourceGroups/unrelated"}
	f.Assessment.Include.ResourceTypes = []string{"st"}
	f.RebuildIndexes()
	coordinator := NewCoordinator(ops)
	run := func() {
		got, err := coordinator.Run(context.Background(), Request{Subscriptions: []string{zoneSub}, InternalPlugins: []string{"sql-eol", "zone-mapping"}, PluginOnly: true, ScannerKeys: []string{"st"}, Filters: f})
		if err != nil || got == nil || got.Completeness != assessment.CompletenessComplete || len(got.PluginTables) != 2 || len(got.PluginTables[0].Rows) != 1 || got.PluginTables[0].Rows[0].Cells[2] != "sql-vm-one" || got.PluginTables[1].Rows[0].Cells[0] != "Dev" {
			t.Errorf("SQL source filters/ownership: %+v %v", got, err)
			return
		}
		got.PluginTables[0].Rows[0].Cells[2] = "mutated returned value"
	}
	run()
	run()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); run() }()
	go func() { defer wg.Done(); run() }()
	wg.Wait()
	if f.Assessment.Include.Tags["Environment"] != "unrelated" {
		t.Fatal("caller filters mutated")
	}
}

func TestSQLAllAdaptersPartialCriticalCancelInvalidAndOptionalPreflight(t *testing.T) {
	for _, mode := range []string{"partial", "critical", "cancel", "invalid"} {
		t.Run(mode, func(t *testing.T) {
			ops := zoneOperations(t)
			sql := realSQLScanner(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			zoneCalls := 0
			oldZone := ops.ScanZoneMapping
			ops.ScanZoneMapping = func(c context.Context, s map[string]string) (zone.Result, error) { zoneCalls++; return oldZone(c, s) }
			service := servicehealth.NewWithTransport(serviceGraph(func(context.Context, arg.Request) (*arg.Response, error) {
				return &arg.Response{Data: []json.RawMessage{json.RawMessage(serviceVM)}}, nil
			}))
			ops.ScanServiceHealth = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
				return service.Scan(c, s, f)
			}
			ops.ScanSQLEOL = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
				value, err := sql.Scan(c, s, f)
				switch mode {
				case "partial":
					return value, errors.New("provider-secret-canary")
				case "cancel":
					cancel()
					return value, context.Canceled
				case "invalid":
					value.Metadata.Description = "unreviewed"
				}
				return value, err
			}
			if mode == "critical" {
				ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
					return nil, errors.New("scope failed")
				}
			}
			got, err := NewCoordinator(ops).Run(ctx, Request{Subscriptions: []string{zoneSub}, InternalPlugins: []string{"zone-mapping", "sql-eol", "service-health"}, PluginOnly: true})
			if got == nil || len(got.PluginTables) != 3 || got.PluginTables[0].Metadata.Name != "service-health" || got.PluginTables[1].Metadata.Name != "sql-eol" || got.PluginTables[2].Metadata.Name != "zone-mapping" || got.Completeness == assessment.CompletenessComplete {
				t.Fatalf("requested ordered partial tables: %+v %v", got, err)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "provider-secret-canary") {
				t.Fatal("SQL provider text leaked")
			}
			if mode == "partial" && (zoneCalls != 1 || len(got.PluginTables[0].Rows) != 1 || len(got.PluginTables[1].Rows) != 1 || len(got.PluginTables[2].Rows) != 1) {
				t.Fatal("SQL failure lost other rows")
			}
			if mode == "invalid" && (zoneCalls != 1 || len(got.PluginTables[1].Rows) != 0 || got.PluginTables[1].Health.Error.Code != "sql_eol_output_invalid") {
				t.Fatal("invalid SQL false success")
			}
			if mode == "cancel" && (!errors.Is(err, context.Canceled) || zoneCalls != 0 || len(got.PluginTables[0].Rows) != 1 || got.PluginTables[2].Health.Status != assessment.StageSkipped) {
				t.Fatal("SQL context lost prior/pending output")
			}
			if mode == "critical" && (err == nil || zoneCalls != 0 || got.PluginTables[0].Health.Status != assessment.StageSkipped || got.PluginTables[1].Health.Status != assessment.StageSkipped || got.PluginTables[2].Health.Status != assessment.StageSkipped) {
				t.Fatal("critical failure lost pending tables")
			}
		})
	}
	ops := noOpOperations()
	ops.ScanSQLEOL = nil
	scopeCalls := 0
	ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		scopeCalls++
		return nil, nil
	}
	if _, err := NewCoordinator(ops).Run(context.Background(), Request{InternalPlugins: []string{"sql-eol"}, PluginOnly: true}); err == nil || scopeCalls != 0 {
		t.Fatal("missing selected SQL operation reached scope")
	}
}
