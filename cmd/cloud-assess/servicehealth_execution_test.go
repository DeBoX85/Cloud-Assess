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
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/xuri/excelize/v2"
)

const cliServiceRow = `{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"eastus","targetResourceType":"microsoft.compute/virtualmachines","percentageOfTimeWithoutEvents":99.5,"events":2,"affectedResources":3}`

func TestServiceCobraAuthenticatedPipelineAllReportsAndModes(t *testing.T) {
	for _, mode := range []string{"only", "mixed", "dual", "dual-partial", "empty", "malformed", "envelope", "later-page", "cancel"} {
		for _, masked := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "-raw", true: "-masked"}[masked], func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				base := filepath.Join(t.TempDir(), "service")
				cred := &cliZoneCredential{t: t}
				var calls atomic.Int32
				opts := azure.DefaultHTTPClientOptions(time.Second)
				opts.Scope = "https://management.usgovcloudapi.net/.default"
				opts.MaxRetries = -1
				opts.Transport = cliZoneTransport(func(r *http.Request) (*http.Response, error) {
					n := calls.Add(1)
					if r.Method != "POST" || r.URL.String() != "https://management.usgovcloudapi.net/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || r.Header.Get("Authorization") != "Bearer synthetic-zone-canary" {
						t.Error("authenticated read-oriented endpoint")
					}
					var request arg.Request
					if json.NewDecoder(r.Body).Decode(&request) != nil || request.Query != servicehealth.Query || !reflect.DeepEqual(request.Subscriptions, []string{cliZoneSub}) || request.Options == nil || request.Options.AuthorizationScopeFilter != nil || *request.Options.Top != 1000 {
						t.Error("query scope/body contract")
					}
					if mode == "cancel" {
						cancel()
						return nil, context.Canceled
					}
					status := 200
					body := `{"data":[` + cliServiceRow + `]}`
					switch mode {
					case "empty":
						body = `{"data":[]}`
					case "malformed":
						body = `{"data":[{}]}`
					case "envelope":
						body = `{"data":[],"Data":[]}`
					case "dual-partial":
						status = 403
						body = `{"error":{"code":"Denied","message":"provider-secret-canary"}}`
					case "later-page":
						if n == 1 {
							body = `{"data":[` + cliServiceRow + `],"$skipToken":"second"}`
						} else {
							status = 403
							body = `{"error":{"code":"Denied","message":"provider-secret-canary"}}`
						}
					}
					if mode == "later-page" && n == 2 && (request.Options.SkipToken == nil || *request.Options.SkipToken != "second") {
						t.Error("opaque continuation")
					}
					return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				})
				scanner, err := servicehealth.NewWithHTTPClient("https://management.usgovcloudapi.net", azure.NewHTTPClient(cred, opts))
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
					field.Set(reflect.MakeFunc(ft, func([]reflect.Value) []reflect.Value { t.Fatalf("unexpected operation %s", name); return nil }))
				}
				ops.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
					return map[string]string{cliZoneSub: "Dev"}, nil
				}
				ops.ScanServiceHealth = func(c context.Context, subs map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
					return scanner.Scan(c, subs, f)
				}
				zoneCalls, graphCalls := 0, 0
				if strings.HasPrefix(mode, "dual") {
					ops.ScanZoneMapping = func(context.Context, map[string]string) (zone.Result, error) {
						zoneCalls++
						return zone.Result{Rows: []zone.Row{{SubscriptionID: cliZoneSub, SubscriptionName: "Dev", Location: "westus", DisplayName: "West US", LogicalZone: "1", PhysicalZone: "westus-az1"}}}, nil
					}
				}
				if mode == "mixed" || strings.HasPrefix(mode, "dual") {
					ops.DiscoverResources = func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
						return &discovery.ResourceInventory{}, nil
					}
					ops.LoadCatalog = func() (*rules.Catalog, error) { graphCalls++; return rules.NewCatalog(), nil }
					field := v.FieldByName("ExecuteGraph")
					ft := field.Type()
					field.Set(reflect.MakeFunc(ft, func([]reflect.Value) []reflect.Value {
						return []reflect.Value{reflect.Zero(ft.Out(0)), reflect.Zero(ft.Out(1)), reflect.Zero(ft.Out(2))}
					}))
				}
				runner := app.NewRunner(orchestration.NewCoordinator(ops))
				var outcome app.Outcome
				root := newRootCommand(func(c context.Context, f scanFlags) (int, error) {
					cfg, e := configureStages(f)
					if e != nil {
						return 1, e
					}
					filters := config.NewFilters()
					filters.Assessment.Include.Tags = map[string]string{"Environment": "dev"}
					filters.RebuildIndexes()
					outcome, e = runner.Run(c, app.ScanOptions{Assessment: assessmentRequest(f, filters, cfg), Outputs: app.OutputOptions{BaseName: f.outputName, JSON: f.json, CSV: f.csv, XLSX: f.xlsx, RedactSubscriptionIDs: f.redactSubscriptionIDs}})
					return outcome.ExitCode, e
				})
				root.SetContext(context.WithValue(ctx, exitCodeContextKey{}, root.Context().Value(exitCodeContextKey{})))
				args := []string{"service-health", "--subscription-id", cliZoneSub, "--resource-group", "unrelated-aggregate-rg", "--json", "--csv", "--output-name", base}
				if !masked {
					args = append(args, "--redact-subscription-ids=false")
				}
				if mode == "mixed" {
					args = []string{"scan", "vm", "--subscription-id", cliZoneSub, "--json", "--csv", "--output-name", base, "--plugin", "service-health", "--stages=-diagnostics,-advisor,-defender"}
					if !masked {
						args = append(args, "--redact-subscription-ids=false")
					}
				}
				if strings.HasPrefix(mode, "dual") {
					args[0] = "scan"
					args = append(args, "--plugin", "zone-mapping,service-health", "--stages=-diagnostics,-advisor,-defender")
				}
				root.SetArgs(args)
				err = root.Execute()
				wantExit := 0
				if mode == "dual-partial" || mode == "envelope" || mode == "later-page" {
					wantExit = 3
				}
				if mode == "cancel" {
					wantExit = 1
					if !errors.Is(err, context.Canceled) {
						t.Fatal("context identity lost")
					}
				}
				if outcome.ExitCode != wantExit || *root.Context().Value(exitCodeContextKey{}).(*int) != wantExit || (err != nil) != (wantExit != 0) {
					t.Fatalf("persisted process exit: %d %v", outcome.ExitCode, err)
				}
				if (mode == "mixed" || strings.HasPrefix(mode, "dual")) && graphCalls != 1 || strings.HasPrefix(mode, "dual") && zoneCalls != 1 {
					t.Fatal("other selected work lost")
				}
				b, e := os.ReadFile(base + ".json")
				if e != nil {
					t.Fatal(e)
				}
				if strings.Contains(string(b), "provider-secret-canary") || masked && strings.Contains(string(b), cliZoneSub) {
					t.Fatal("JSON provider/identity leak")
				}
				var report result.AssessmentResult
				if json.Unmarshal(b, &report) != nil || report.SchemaVersion != "1.1" {
					t.Fatal("canonical plugin assessment missing")
				}
				count := 1
				if strings.HasPrefix(mode, "dual") {
					count = 2
				}
				if len(report.PluginTables) != count {
					t.Fatal("requested tables missing")
				}
				table := report.PluginTables[0]
				rows := 0
				if mode == "only" || mode == "mixed" || mode == "dual" || mode == "later-page" {
					rows = 1
				}
				if len(table.Rows) != rows || table.Health.Records != rows {
					t.Fatal("source rows/retention changed")
				}
				if calls.Load() == 0 || cred.calls.Load() == 0 || mode == "later-page" && calls.Load() != 2 {
					t.Fatal("real authenticated transport not observed")
				}
				headers := []string{"Subscription ID", "Target Region", "Target Resource Type", "Percentage Without Events", "Events Count", "Affected Resources"}
				expected := [][]string{headers}
				if rows == 1 {
					id := cliZoneSub
					if masked {
						id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxx1111111"
					}
					expected = append(expected, []string{id, "eastus", "microsoft.compute/virtualmachines", "99.50%", "2", "3"})
					if !reflect.DeepEqual(table.Rows[0].Cells, expected[1]) {
						t.Fatalf("every canonical source cell: %+v", table.Rows[0])
					}
				}
				f, e := os.Open(base + ".plugin_service-health_availability.csv")
				if e != nil {
					t.Fatal(e)
				}
				cells, e := csv.NewReader(f).ReadAll()
				f.Close()
				if e != nil || !reflect.DeepEqual(cells, expected) {
					t.Fatalf("every CSV cell: %v %+v", e, cells)
				}
				book, e := excelize.OpenFile(base + ".xlsx")
				if e != nil {
					t.Fatal(e)
				}
				xls, e := book.GetRows(table.SheetName)
				book.Close()
				if e != nil || len(xls) < 4 || !reflect.DeepEqual(xls[0], []string{branding.Default().ReportTitle}) || !reflect.DeepEqual(xls[1], []string{table.SheetName}) || !reflect.DeepEqual(xls[3:], expected) {
					t.Fatalf("every XLSX cell: %v %+v", e, xls)
				}
				if mode == "malformed" && (table.Health.Status != assessment.StageCompletedWithWarnings || report.Completeness != assessment.CompletenessCompleteWithWarnings) {
					t.Fatal("malformed warning concealed")
				}
				if wantExit == 3 && report.Completeness != assessment.CompletenessPartial {
					t.Fatal("false complete failure")
				}
			})
		}
	}
}
