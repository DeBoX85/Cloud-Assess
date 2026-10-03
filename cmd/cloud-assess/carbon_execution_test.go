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
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/carbon"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/servicehealth"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/sqleol"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/xuri/excelize/v2"
)

// Literal service input used by the unchanged pinned source capture. Expected
// metadata and every cell below come from its separately hashed source output.
const cliCarbonItems = `[{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Compute/virtualMachines","latestMonthEmissions":80,"previousMonthEmissions":100,"monthlyEmissionsChangeValue":-20,"monthOverMonthEmissionsChangeRatio":999},{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Compute/virtualMachines","latestMonthEmissions":20,"previousMonthEmissions":0},{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Storage/storageAccounts","latestMonthEmissions":75,"previousMonthEmissions":0,"monthlyEmissionsChangeValue":0}]`

func TestCarbonCobraAuthenticatedPipelineSourceCellsReportsAndModes(t *testing.T) {
	for _, mode := range []string{"only", "scanner", "mixed", "all", "all-partial", "all-cancel", "empty", "malformed", "missing-access", "denied", "later-page", "cancel"} {
		for _, masked := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "-raw", true: "-masked"}[masked], func(t *testing.T) {
				capture := "source-all.json"
				if mode == "scanner" {
					capture = "source-filtered.json"
				}
				b, e := os.ReadFile("../../internal/plugins/carbon/testdata/" + capture)
				if e != nil {
					t.Fatal(e)
				}
				var source []struct {
					Metadata    struct{ Name, Version, Description, Author, License string }
					SheetName   string `json:"sheet_name"`
					Description string
					Table       [][]string
				}
				if json.Unmarshal(b, &source) != nil || len(source) != 1 {
					t.Fatal("source capture")
				}
				expected := source[0].Table
				if mode == "empty" || mode == "malformed" || mode == "all-partial" {
					expected = expected[:1]
				}
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				base := filepath.Join(t.TempDir(), "carbon")
				cred := &cliZoneCredential{t: t}
				var calls atomic.Int32
				opts := azure.DefaultHTTPClientOptions(time.Second)
				opts.Scope = "https://management.usgovcloudapi.net/.default"
				opts.MaxRetries = -1
				opts.Transport = cliZoneTransport(func(r *http.Request) (*http.Response, error) {
					n := calls.Add(1)
					prefix := "https://management.usgovcloudapi.net/providers/Microsoft.Carbon/"
					if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer synthetic-zone-canary" {
						t.Error("carbon authenticated read-oriented request")
					}
					body := `{"startDate":"2026-01-01","endDate":"2026-03-31"}`
					status := 200
					if n == 1 {
						if r.URL.String() != prefix+"queryCarbonEmissionDataAvailableDateRange?api-version=2025-04-01" || r.Body != nil {
							t.Error("carbon date request contract")
						}
					} else {
						if r.URL.String() != prefix+"carbonEmissionReports?api-version=2025-04-01" {
							t.Error("carbon fixed report endpoint")
						}
						var request map[string]any
						if json.NewDecoder(r.Body).Decode(&request) != nil {
							t.Fatal("carbon request body")
						}
						want := map[string]any{"reportType": "ItemDetailsReport", "subscriptionList": []any{cliZoneSub, cliSQLSubB}, "carbonScopeList": []any{"Scope1", "Scope2", "Scope3"}, "dateRange": map[string]any{"start": "2026-03-31", "end": "2026-03-31"}, "categoryType": "ResourceType", "orderBy": "LatestMonthEmissions", "sortDirection": "Desc", "pageSize": float64(1000)}
						if n == 3 {
							want["skipToken"] = "https://private-token.invalid/canary"
						}
						if !reflect.DeepEqual(request, want) {
							t.Errorf("carbon source body/scope: %#v", request)
						}
						decisions := `[{"subscriptionId":"11111111-1111-4111-8111-111111111111","decision":"Allowed"},{"subscriptionId":"22222222-2222-4222-8222-222222222222","decision":"Allowed"}]`
						items := cliCarbonItems
						token := ""
						switch mode {
						case "empty":
							items = `[]`
						case "malformed":
							items = `[{}]`
						case "missing-access":
							decisions = `null`
						case "denied":
							decisions = `[{"subscriptionId":"11111111-1111-4111-8111-111111111111","decision":"Allowed"},{"subscriptionId":"22222222-2222-4222-8222-222222222222","decision":"Denied","reason":"provider-secret-canary"}]`
						case "all-partial":
							status = 403
						case "later-page", "cancel", "all-cancel":
							if n == 2 {
								token = "https://private-token.invalid/canary"
							} else if mode == "later-page" {
								status = 403
							} else {
								cancel()
								return nil, context.Canceled
							}
						}
						body = `{"value":` + items + `,"subscriptionAccessDecisionList":` + decisions + `,"skipToken":"` + token + `"}`
						if status != 200 {
							body = `{"error":{"message":"provider-secret-canary"}}`
						}
					}
					return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				})
				scanner, e := carbon.NewWithHTTPClient("https://management.usgovcloudapi.net", azure.NewHTTPClient(cred, opts))
				if e != nil {
					t.Fatal(e)
				}
				ops := orchestration.Operations{}
				v := reflect.ValueOf(&ops).Elem()
				typ := v.Type()
				for i := 0; i < v.NumField(); i++ {
					name := typ.Field(i).Name
					f := v.Field(i)
					f.Set(reflect.MakeFunc(f.Type(), func([]reflect.Value) []reflect.Value {
						t.Fatalf("unexpected carbon-only operation %s", name)
						return nil
					}))
				}
				ops.DiscoverSubscriptions = func(_ context.Context, ids []string, _ *config.Filters) (map[string]string, error) {
					if !reflect.DeepEqual(ids, []string{cliZoneSub, cliSQLSubB}) {
						t.Error("CLI selected scope")
					}
					return map[string]string{cliZoneSub: "A", cliSQLSubB: "B"}, nil
				}
				order := []string{}
				ops.ScanCarbon = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
					order = append(order, "carbon-emissions")
					return scanner.Scan(c, s, f)
				}
				all := strings.HasPrefix(mode, "all")
				regular := all || mode == "mixed" || mode == "scanner"
				if all {
					service := servicehealth.NewWithTransport(cliSQLGraph(func(context.Context, arg.Request) (*arg.Response, error) {
						return &arg.Response{Data: []json.RawMessage{json.RawMessage(`{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"westus","targetResourceType":"microsoft.compute/virtualmachines","percentageOfTimeWithoutEvents":100,"events":0,"affectedResources":0}`)}}, nil
					}))
					sqlInput, e := os.ReadFile("../../internal/plugins/sqleol/testdata/source-input.json")
					if e != nil {
						t.Fatal(e)
					}
					var sqlRows []json.RawMessage
					if json.Unmarshal(sqlInput, &sqlRows) != nil {
						t.Fatal("SQL capture")
					}
					sql := sqleol.NewWithTransport(cliSQLGraph(func(context.Context, arg.Request) (*arg.Response, error) { return &arg.Response{Data: sqlRows}, nil }))
					ops.ScanServiceHealth = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
						order = append(order, "service-health")
						return service.Scan(c, s, f)
					}
					ops.ScanSQLEOL = func(c context.Context, s map[string]string, f *config.AssessmentFilter) (assessment.PluginTable, error) {
						order = append(order, "sql-eol")
						return sql.Scan(c, s, f)
					}
					ops.ScanZoneMapping = func(context.Context, map[string]string) (zone.Result, error) {
						order = append(order, "zone-mapping")
						return zone.Result{Rows: []zone.Row{{SubscriptionID: cliZoneSub, SubscriptionName: "A", Location: "westus", DisplayName: "West US", LogicalZone: "1", PhysicalZone: "westus-az1"}}}, nil
					}
				}
				inventoryCalls := 0
				if regular {
					ops.DiscoverResources = func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
						inventoryCalls++
						return &discovery.ResourceInventory{Included: []assessment.Resource{{ID: "/subscriptions/" + cliZoneSub + "/resourceGroups/unrelated/providers/Microsoft.Compute/virtualMachines/healthy", SubscriptionID: cliZoneSub, ResourceGroup: "unrelated", Type: "microsoft.compute/virtualmachines", Name: "healthy"}}}, nil
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
					filters.RebuildIndexes()
					outcome, e = runner.Run(c, app.ScanOptions{Assessment: assessmentRequest(f, filters, cfg), Outputs: app.OutputOptions{BaseName: f.outputName, JSON: f.json, CSV: f.csv, XLSX: f.xlsx, SARIF: f.sarif, RedactSubscriptionIDs: f.redactSubscriptionIDs}})
					return outcome.ExitCode, e
				})
				root.SetContext(context.WithValue(ctx, exitCodeContextKey{}, root.Context().Value(exitCodeContextKey{})))
				args := []string{"carbon-emissions", "--subscription-id", cliZoneSub + "," + cliSQLSubB, "--json", "--csv", "--sarif", "--output-name", base}
				if regular {
					args[0] = "scan"
					args = append(args, "--plugin", "carbon-emissions", "--stages=-diagnostics,-advisor,-defender")
				}
				if mode == "scanner" {
					args = append([]string{args[0], "vm"}, args[1:]...)
				}
				if all {
					for i, arg := range args {
						if arg == "--plugin" {
							args[i+1] = "zone-mapping,sql-eol,carbon-emissions,service-health,carbon-emissions"
						}
					}
				}
				if !masked {
					args = append(args, "--redact-subscription-ids=false")
				}
				root.SetArgs(args)
				e = root.Execute()
				wantExit := 0
				switch mode {
				case "denied", "all-partial", "later-page":
					wantExit = 3
				case "cancel", "all-cancel":
					wantExit = 1
					if !errors.Is(e, context.Canceled) {
						t.Fatal("carbon cancellation identity lost")
					}
				}
				if outcome.ExitCode != wantExit || *root.Context().Value(exitCodeContextKey{}).(*int) != wantExit || (e != nil) != (wantExit != 0) {
					t.Fatalf("carbon artifact-before-exit: %d %v", outcome.ExitCode, e)
				}
				b, e = os.ReadFile(base + ".json")
				if e != nil {
					t.Fatal(e)
				}
				if strings.Contains(string(b), "canary") || masked && (strings.Contains(string(b), cliZoneSub) || strings.Contains(string(b), cliSQLSubB)) {
					t.Fatal("carbon identity/provider/token leakage")
				}
				var report result.AssessmentResult
				if json.Unmarshal(b, &report) != nil || report.SchemaVersion != "1.1" {
					t.Fatal("carbon canonical schema")
				}
				count := 1
				if all {
					count = 4
				}
				if len(report.PluginTables) != count {
					t.Fatal("requested plugin tables missing")
				}
				table := report.PluginTables[0]
				meta := source[0].Metadata
				if table.Metadata.Name != meta.Name || table.Metadata.Version != meta.Version || table.Metadata.Description != meta.Description || table.Metadata.Author != meta.Author || table.Metadata.License != meta.License || table.Metadata.Type != "internal" || table.SheetName != source[0].SheetName || table.Description != source[0].Description || !reflect.DeepEqual(table.Columns, expected[0]) || len(table.Rows) != len(expected)-1 || table.Health.Records != len(expected)-1 {
					t.Fatal("carbon source metadata/headers/rows")
				}
				for i, row := range table.Rows {
					if row.SubscriptionID != "" || !reflect.DeepEqual(row.Cells, expected[i+1]) {
						t.Fatalf("every source carbon JSON cell: %+v", row)
					}
				}
				if regular && (inventoryCalls != 1 || len(report.Resources) != 1 || report.Resources[0].Name != "healthy") {
					t.Fatal("carbon failure lost regular inventory")
				}
				if all {
					want := []string{"carbon-emissions", "service-health", "sql-eol", "zone-mapping"}
					if mode == "all-cancel" {
						want = want[:1]
					}
					if !reflect.DeepEqual(order, want) {
						t.Fatalf("carbon source order/dedup/cancel: %v", order)
					}
					for i, name := range []string{"carbon-emissions", "service-health", "sql-eol", "zone-mapping"} {
						if report.PluginTables[i].Metadata.Name != name {
							t.Fatal("all tables order")
						}
						if i > 0 && mode == "all-cancel" {
							if report.PluginTables[i].Health.Status != assessment.StageSkipped {
								t.Fatal("cancel pending table lost")
							}
						} else if i > 0 && len(report.PluginTables[i].Rows) == 0 {
							t.Fatal("non-context failure lost other adapter rows")
						}
					}
				}
				wantCalls := int32(2)
				if mode == "later-page" || mode == "cancel" || mode == "all-cancel" {
					wantCalls = 3
				}
				if calls.Load() != wantCalls || cred.calls.Load() == 0 {
					t.Fatal("authenticated bounded carbon pipeline not observed")
				}
				csvExpected := make([][]string, len(expected))
				for i, row := range expected {
					csvExpected[i] = append([]string(nil), row...)
				}
				if len(expected) > 1 {
					csvExpected[1][6] = "'-20.00"
				}
				f, e := os.Open(base + ".plugin_carbon-emissions_emissions.csv")
				if e != nil {
					t.Fatal(e)
				}
				cells, e := csv.NewReader(f).ReadAll()
				f.Close()
				if e != nil || !reflect.DeepEqual(cells, csvExpected) {
					t.Fatalf("every carbon CSV cell: %v %v", cells, e)
				}
				book, e := excelize.OpenFile(base + ".xlsx")
				if e != nil {
					t.Fatal(e)
				}
				xls, e := book.GetRows("Carbon Emissions")
				book.Close()
				if e != nil || len(xls) < 4 || !reflect.DeepEqual(xls[0], []string{branding.Default().ReportTitle}) || !reflect.DeepEqual(xls[1], []string{"Carbon Emissions"}) || !reflect.DeepEqual(xls[3:], expected) {
					t.Fatalf("every carbon XLSX cell: %v %v", xls, e)
				}
				sarifBytes, e := os.ReadFile(base + ".sarif")
				if e != nil || strings.Contains(string(sarifBytes), "carbon-emissions") || strings.Contains(string(sarifBytes), "Latest Month Emissions") {
					t.Fatal("aggregate carbon table became invented SARIF findings")
				}
				if mode == "malformed" || mode == "missing-access" {
					if table.Health.Status != assessment.StageCompletedWithWarnings || report.Completeness != assessment.CompletenessCompleteWithWarnings {
						t.Fatal("carbon warning concealed")
					}
				}
				if wantExit == 3 && report.Completeness != assessment.CompletenessPartial {
					t.Fatal("carbon failed access/request false completeness")
				}
			})
		}
	}
}
