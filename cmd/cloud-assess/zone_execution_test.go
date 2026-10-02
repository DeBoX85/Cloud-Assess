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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/app"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/xuri/excelize/v2"
)

const cliZoneSub = "11111111-1111-4111-8111-111111111111"

type cliZoneTransport func(*http.Request) (*http.Response, error)

func (f cliZoneTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

type cliZoneCredential struct {
	t     *testing.T
	calls atomic.Int32
}

func (c *cliZoneCredential) GetToken(_ context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	if !reflect.DeepEqual(opts.Scopes, []string{"https://management.usgovcloudapi.net/.default"}) {
		c.t.Error("wrong selected-cloud audience")
	}
	return azcore.AccessToken{Token: "synthetic-zone-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestZoneCobraAuthenticatedCoordinatorApplicationReports(t *testing.T) {
	for _, mode := range []string{"only", "mixed", "partial", "empty", "malformed", "later-page", "cancel"} {
		for _, masked := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "-raw", true: "-masked"}[masked], func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				base := filepath.Join(t.TempDir(), "zone")
				cred := &cliZoneCredential{t: t}
				var calls atomic.Int32
				options := azure.DefaultHTTPClientOptions(time.Second)
				options.Scope = "https://management.usgovcloudapi.net/.default"
				options.MaxRetries = -1
				options.Transport = cliZoneTransport(func(req *http.Request) (*http.Response, error) {
					n := calls.Add(1)
					expected := "https://management.usgovcloudapi.net/subscriptions/" + cliZoneSub + "/locations?api-version=2022-12-01"
					if n == 2 && mode == "later-page" {
						expected += "&page=2"
					}
					if req.Method != "GET" || req.URL.String() != expected || req.Header.Get("Authorization") != "Bearer synthetic-zone-canary" || req.Body != nil {
						t.Error("authenticated read request boundary changed")
					}
					if mode == "cancel" {
						cancel()
						return nil, context.Canceled
					}
					code := 200
					body := `{"value":[{"name":"westus","displayName":"West US","availabilityZoneMappings":[{"logicalZone":"1","physicalZone":"westus-az1"}]}]}`
					if mode == "empty" {
						body = `{"value":[]}`
					}
					if mode == "malformed" {
						body = `{"value":[null]}`
					}
					if mode == "partial" || mode == "later-page" && n == 2 {
						code = 403
						body = `{"error":{"code":"Denied","message":"provider-secret-canary"}}`
					}
					if mode == "later-page" && n == 1 {
						body = strings.TrimSuffix(body, "}") + `,"nextLink":"` + expected + `&page=2"}`
					}
					return &http.Response{StatusCode: code, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
				})
				scanner, err := zone.NewScanner("https://management.usgovcloudapi.net", azure.NewHTTPClient(cred, options))
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
				ops.ScanZoneMapping = scanner.Scan
				graphCalls := 0
				if mode == "mixed" {
					ops.DiscoverResources = func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
						return &discovery.ResourceInventory{}, nil
					}
					ops.LoadCatalog = func() (*rules.Catalog, error) { graphCalls++; return rules.NewCatalog(), nil }
					// The ordinary Graph executor is allowed, others remain tripwires.
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
					outcome, e = runner.Run(c, app.ScanOptions{Assessment: assessmentRequest(f, config.NewFilters(), cfg), Outputs: app.OutputOptions{BaseName: f.outputName, JSON: f.json, CSV: f.csv, XLSX: f.xlsx, RedactSubscriptionIDs: f.redactSubscriptionIDs}})
					return outcome.ExitCode, e
				})
				root.SetContext(context.WithValue(ctx, exitCodeContextKey{}, root.Context().Value(exitCodeContextKey{})))
				args := []string{"zone-mapping", "--subscription-id", cliZoneSub, "--json", "--csv", "--output-name", base}
				if !masked {
					args = append(args, "--redact-subscription-ids=false")
				}
				if mode == "mixed" {
					args[0] = "scan"
					args = append(args, "--plugin", "zone-mapping", "--stages=-diagnostics,-advisor,-defender")
				}
				root.SetArgs(args)
				err = root.Execute()
				failed := mode == "partial" || mode == "malformed" || mode == "later-page" || mode == "cancel"
				if !failed && err != nil || failed && err == nil {
					t.Fatalf("exit/error mismatch %s: %v", mode, err)
				}
				wantExit := 0
				if failed {
					wantExit = 3
				}
				if mode == "cancel" {
					wantExit = 1
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("cancellation lost: %v", err)
					}
				}
				if outcome.ExitCode != wantExit || *root.Context().Value(exitCodeContextKey{}).(*int) != wantExit {
					t.Fatalf("process mapped exit %d want %d", outcome.ExitCode, wantExit)
				}
				if mode == "mixed" && graphCalls != 1 {
					t.Fatal("mixed Graph lost")
				}
				data, e := os.ReadFile(base + ".json")
				if e != nil {
					t.Fatal(e)
				}
				if strings.Contains(string(data), "provider-secret-canary") {
					t.Fatal("raw provider error leaked")
				}
				var report result.AssessmentResult
				if e := json.Unmarshal(data, &report); e != nil {
					t.Fatal(e)
				}
				if len(report.PluginTables) != 1 || report.SchemaVersion != "1.1" {
					t.Fatal("requested plugin missing")
				}
				rows := 0
				if mode == "only" || mode == "mixed" || mode == "later-page" {
					rows = 1
				}
				if len(report.PluginTables[0].Rows) != rows || report.PluginTables[0].Health.Records != rows {
					t.Fatal("row retention/count changed")
				}
				if mode == "later-page" && calls.Load() != 2 {
					t.Fatal("continuation not executed")
				}
				if cred.calls.Load() == 0 {
					t.Fatal("synthetic pipeline never authenticated")
				}
				if masked && strings.Contains(string(data), cliZoneSub) {
					t.Fatal("JSON identity leaked")
				}
				tableFile, e := os.Open(base + ".plugin_zone-mapping_zones.csv")
				if e != nil {
					t.Fatal(e)
				}
				cells, e := csv.NewReader(tableFile).ReadAll()
				tableFile.Close()
				if e != nil {
					t.Fatal(e)
				}
				if len(cells) != rows+1 || !reflect.DeepEqual(cells[0], []string{"Subscription", "Location", "Display Name", "Logical Zone", "Physical Zone"}) {
					t.Fatal("CSV rows/headers changed")
				}
				workbook, e := excelize.OpenFile(base + ".xlsx")
				if e != nil {
					t.Fatal(e)
				}
				excelRows, e := workbook.GetRows("Zone Mapping")
				workbook.Close()
				if e != nil || len(excelRows) < 4 || !reflect.DeepEqual(excelRows[0], []string{branding.Default().ReportTitle}) || !reflect.DeepEqual(excelRows[1], []string{"Zone Mapping"}) || !reflect.DeepEqual(cells, excelRows[3:]) {
					t.Fatalf("XLSX/CSV cells differ: %v %#v %#v", e, cells, excelRows)
				}
				if rows == 1 && !reflect.DeepEqual(cells[1], []string{"Dev", "westus", "West US", "1", "westus-az1"}) {
					t.Fatal("source cells changed")
				}
				if failed && mode != "cancel" && report.Completeness != assessment.CompletenessPartial {
					t.Fatal("false complete failure")
				}
			})
		}
	}
}
