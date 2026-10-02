package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/app"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/servicehealth"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/sqleol"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/xuri/excelize/v2"
)

const cliSQLSubB = "22222222-2222-4222-8222-222222222222"

type cliSQLGraph func(context.Context, arg.Request) (*arg.Response, error)

func (f cliSQLGraph) Do(c context.Context, r arg.Request) (*arg.Response, error) { return f(c, r) }

func TestSQLCobraAuthenticatedPipelineSourceCellsAllReportsAndModes(t *testing.T) {
	for _, mode := range []string{"only", "scanner", "mixed", "all", "all-partial", "all-cancel", "empty", "malformed", "envelope", "later-page", "cancel", "identity-text"} {
		for _, masked := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "-raw", true: "-masked"}[masked], func(t *testing.T) {
				input, err := os.ReadFile("../../internal/plugins/sqleol/testdata/source-input.json")
				if err != nil {
					t.Fatal(err)
				}
				var rows []json.RawMessage
				if json.Unmarshal(input, &rows) != nil {
					t.Fatal("source input")
				}
				capture, err := os.ReadFile("../../internal/plugins/sqleol/testdata/source-output.json")
				if err != nil {
					t.Fatal(err)
				}
				var source []struct{ Table [][]string }
				if json.Unmarshal(capture, &source) != nil || len(source) != 1 {
					t.Fatal("source capture")
				}
				expected := source[0].Table
				if mode == "identity-text" {
					var row map[string]any
					json.Unmarshal(rows[1], &row)
					row["Name"] = "sql-vm-one " + cliZoneSub
					rows[1], _ = json.Marshal(row)
					expected[2][2] = "sql-vm-one " + cliZoneSub
					if masked {
						expected[2][2] = "sql-vm-one xxxxxxxx-xxxx-xxxx-xxxx-xxxxx1111111"
					}
				}
				rowCount := 2
				if mode == "empty" || mode == "malformed" || mode == "envelope" || mode == "all-partial" || mode == "cancel" || mode == "all-cancel" {
					rowCount = 0
					expected = expected[:1]
				}
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				base := filepath.Join(t.TempDir(), "sql")
				cred := &cliZoneCredential{t: t}
				var httpCalls atomic.Int32
				opts := azure.DefaultHTTPClientOptions(time.Second)
				opts.Scope = "https://management.usgovcloudapi.net/.default"
				opts.MaxRetries = -1
				opts.Transport = cliZoneTransport(func(r *http.Request) (*http.Response, error) {
					n := httpCalls.Add(1)
					if r.Method != "POST" || r.URL.String() != "https://management.usgovcloudapi.net/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || r.Header.Get("Authorization") != "Bearer synthetic-zone-canary" {
						t.Error("SQL authenticated read-oriented request")
					}
					var request arg.Request
					if json.NewDecoder(r.Body).Decode(&request) != nil || request.Query != sqleol.Query || !reflect.DeepEqual(request.Subscriptions, []string{cliZoneSub, cliSQLSubB}) || request.Options == nil || request.Options.Top == nil || *request.Options.Top != 1000 || request.Options.AuthorizationScopeFilter != nil {
						t.Error("SQL source query/scope/body")
					}
					if mode == "cancel" || mode == "all-cancel" {
						cancel()
						return nil, context.Canceled
					}
					data, _ := json.Marshal(rows)
					body := `{"data":` + string(data) + `}`
					status := 200
					switch mode {
					case "empty":
						body = `{"data":[]}`
					case "malformed":
						body = `{"data":[{}]}`
					case "envelope":
						body = `{"data":[],"Data":[]}`
					case "all-partial":
						status = 403
						body = `{"error":{"message":"provider-secret-canary"}}`
					case "later-page":
						if n == 1 {
							body = `{"data":` + string(data) + `,"$skipToken":"second"}`
						} else {
							status = 403
							body = `{"error":{"message":"provider-secret-canary"}}`
							if request.Options.SkipToken == nil || *request.Options.SkipToken != "second" {
								t.Error("SQL safe continuation")
							}
						}
					}
					return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				})
				scanner, err := sqleol.NewWithHTTPClient("https://management.usgovcloudapi.net", azure.NewHTTPClient(cred, opts))
				if err != nil {
					t.Fatal(err)
				}
				ops := orchestration.Operations{}
				v := reflect.ValueOf(&ops).Elem()
				typ := v.Type()
				for i := 0; i < v.NumField(); i++ {
					name := typ.Field(i).Name
					field := v.Field(i)
					ft := field.Type()
					field.Set(reflect.MakeFunc(ft, func([]reflect.Value) []reflect.Value { t.Fatalf("unexpected SQL-only operation %s", name); return nil }))
				}
				ops.DiscoverSubscriptions = func(_ context.Context, ids []string, _ *config.Filters) (map[string]string, error) {
					if !reflect.DeepEqual(ids, []string{cliZoneSub, cliSQLSubB}) {
						t.Error("selected CLI scope")
					}
					return map[string]string{cliZoneSub: "A", cliSQLSubB: "B"}, nil
				}
				ops.ScanSQLEOL = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
					return scanner.Scan(c, s, f)
				}
				all := strings.HasPrefix(mode, "all")
				regular := all || mode == "mixed" || mode == "scanner"
				zoneCalls, serviceCalls, inventoryCalls := 0, 0, 0
				if all {
					service := servicehealth.NewWithTransport(cliSQLGraph(func(context.Context, arg.Request) (*arg.Response, error) {
						return &arg.Response{Data: []json.RawMessage{json.RawMessage(`{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"westus","targetResourceType":"microsoft.storage/storageaccounts","percentageOfTimeWithoutEvents":100,"events":0,"affectedResources":0}`)}}, nil
					}))
					ops.ScanServiceHealth = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
						serviceCalls++
						return service.Scan(c, s, f)
					}
					ops.ScanZoneMapping = func(context.Context, map[string]string) (zone.Result, error) {
						zoneCalls++
						return zone.Result{Rows: []zone.Row{{SubscriptionID: cliZoneSub, SubscriptionName: "A", Location: "westus", DisplayName: "West US", LogicalZone: "1", PhysicalZone: "westus-az1"}}}, nil
					}
				}
				if regular {
					ops.DiscoverResources = func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
						inventoryCalls++
						return &discovery.ResourceInventory{Included: []assessment.Resource{{ID: "/subscriptions/" + cliZoneSub + "/resourceGroups/unrelated/providers/Microsoft.Storage/storageAccounts/healthy", SubscriptionID: cliZoneSub, ResourceGroup: "unrelated", Type: "microsoft.storage/storageaccounts", Name: "healthy", Tags: map[string]string{"Environment": "unrelated"}}}}, nil
					}
					ops.LoadCatalog = func() (*rules.Catalog, error) { return rules.NewCatalog(), nil }
					ops.ExecuteGraph = func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
						return nil, nil, nil
					}
				}
				runner := app.NewRunner(orchestration.NewCoordinator(ops))
				var outcome app.Outcome
				root := newRootCommand(func(c context.Context, f scanFlags) (int, error) {
					cfg, e := configureStages(f)
					if e != nil {
						return 1, e
					}
					filters := config.NewFilters()
					filters.Assessment.Include.ResourceGroups = []string{"/subscriptions/" + cliZoneSub + "/resourceGroups/unrelated"}
					filters.Assessment.Include.Tags = map[string]string{"Environment": "unrelated"}
					filters.Assessment.Include.ResourceTypes = []string{"st"}
					filters.RebuildIndexes()
					outcome, e = runner.Run(c, app.ScanOptions{Assessment: assessmentRequest(f, filters, cfg), Outputs: app.OutputOptions{BaseName: f.outputName, JSON: f.json, CSV: f.csv, XLSX: f.xlsx, RedactSubscriptionIDs: f.redactSubscriptionIDs}})
					return outcome.ExitCode, e
				})
				root.SetContext(context.WithValue(ctx, exitCodeContextKey{}, root.Context().Value(exitCodeContextKey{})))
				args := []string{"sql-eol", "--subscription-id", cliZoneSub + "," + cliSQLSubB, "--json", "--csv", "--output-name", base}
				if regular {
					args[0] = "scan"
					args = append(args, "--plugin", "sql-eol", "--stages=-diagnostics,-advisor,-defender")
				}
				if mode == "scanner" {
					args = append([]string{args[0], "st"}, args[1:]...)
				}
				if all {
					for i, arg := range args {
						if arg == "--plugin" {
							args[i+1] = "zone-mapping,sql-eol,service-health,sql-eol"
						}
					}
				}
				if !masked {
					args = append(args, "--redact-subscription-ids=false")
				}
				root.SetArgs(args)
				err = root.Execute()
				wantExit := 0
				if mode == "all-partial" || mode == "envelope" || mode == "later-page" {
					wantExit = 3
				}
				if mode == "cancel" || mode == "all-cancel" {
					wantExit = 1
					if !errors.Is(err, context.Canceled) {
						t.Fatal("SQL context identity lost")
					}
				}
				if outcome.ExitCode != wantExit || *root.Context().Value(exitCodeContextKey{}).(*int) != wantExit || (err != nil) != (wantExit != 0) {
					t.Fatalf("SQL artifact-before-exit: %d %v", outcome.ExitCode, err)
				}
				b, e := os.ReadFile(base + ".json")
				if e != nil {
					t.Fatal(e)
				}
				if strings.Contains(string(b), "provider-secret-canary") || masked && (strings.Contains(string(b), cliZoneSub) || strings.Contains(string(b), cliSQLSubB)) {
					t.Fatal("SQL identity/provider leakage")
				}
				var report result.AssessmentResult
				if json.Unmarshal(b, &report) != nil || report.SchemaVersion != "1.1" {
					t.Fatal("SQL canonical schema")
				}
				tableIndex := 0
				wantTables := 1
				if all {
					tableIndex = 1
					wantTables = 3
				}
				if len(report.PluginTables) != wantTables {
					t.Fatal("SQL requested tables missing")
				}
				table := report.PluginTables[tableIndex]
				if table.Metadata.Name != "sql-eol" || len(table.Rows) != rowCount || table.Health.Records != rowCount || !reflect.DeepEqual(table.Columns, expected[0]) {
					t.Fatal("SQL identity/32 headers/rows")
				}
				for i, row := range table.Rows {
					if !reflect.DeepEqual(row.Cells, expected[i+1]) {
						t.Fatalf("every SQL source JSON cell/order: %+v", row)
					}
				}
				if regular && (inventoryCalls != 1 || len(report.Resources) != 1 || report.Resources[0].Name != "healthy") {
					t.Fatal("SQL failure lost healthy normal inventory")
				}
				if all && (serviceCalls != 1 || len(report.PluginTables[0].Rows) != 1 || report.PluginTables[2].Metadata.Name != "zone-mapping") {
					t.Fatal("SQL lost prior/other selected adapter")
				}
				if all && mode != "all-cancel" && (zoneCalls != 1 || len(report.PluginTables[2].Rows) != 1) {
					t.Fatal("SQL non-context failure lost zone")
				}
				if mode == "all-cancel" && (zoneCalls != 0 || report.PluginTables[2].Health.Status != assessment.StageSkipped) {
					t.Fatal("SQL cancellation lost pending zone")
				}
				if httpCalls.Load() == 0 || cred.calls.Load() == 0 || mode == "later-page" && httpCalls.Load() != 2 {
					t.Fatal("SQL authenticated pipeline not observed")
				}
				csvExpected := make([][]string, len(expected))
				for i, row := range expected {
					csvExpected[i] = append([]string(nil), row...)
				}
				// The negative saving is source text. CSV's accepted formula protection prefixes it; JSON/XLSX retain source text.
				if rowCount > 0 {
					csvExpected[1][30] = "'-12.34"
				}
				f, e := os.Open(base + ".plugin_sql-eol_lifecycle.csv")
				if e != nil {
					t.Fatal(e)
				}
				cells, e := csv.NewReader(f).ReadAll()
				f.Close()
				if e != nil || !reflect.DeepEqual(cells, csvExpected) {
					t.Fatalf("every SQL CSV cell: %+v %v", cells, e)
				}
				book, e := excelize.OpenFile(base + ".xlsx")
				if e != nil {
					t.Fatal(e)
				}
				xls, e := book.GetRows("SQL EOL")
				book.Close()
				if e != nil || len(xls) < 4 || !reflect.DeepEqual(xls[0], []string{branding.Default().ReportTitle}) || !reflect.DeepEqual(xls[1], []string{"SQL EOL"}) || !reflect.DeepEqual(xls[3:], expected) {
					t.Fatalf("every SQL XLSX cell: %+v %v", xls, e)
				}
				if mode == "malformed" && (table.Health.Status != assessment.StageCompletedWithWarnings || report.Completeness != assessment.CompletenessCompleteWithWarnings) {
					t.Fatal("SQL malformed health concealed")
				}
				if wantExit == 3 && report.Completeness != assessment.CompletenessPartial {
					t.Fatal("SQL failure marked complete")
				}
			})
		}
	}
}
