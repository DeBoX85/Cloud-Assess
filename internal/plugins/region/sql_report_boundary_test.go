package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const sqlReportDBReady = "{\"name\":\"eastus\",\"status\":\"Available\",\"supportedServerVersions\":[{\"name\":\"12\",\"status\":\"Available\",\"supportedEditions\":[{\"name\":\"basic\",\"status\":\"Available\",\"supportedServiceLevelObjectives\":[{\"name\":\"o1\",\"status\":\"Available\",\"sku\":{\"name\":\"same_sku\",\"tier\":\"basic\",\"capacity\":1}}]}]}]}"
const sqlReportMIReady = "{\"name\":\"eastus\",\"status\":\"Available\",\"supportedManagedInstanceVersions\":[{\"name\":\"12\",\"status\":\"Available\",\"supportedEditions\":[{\"name\":\"GeneralPurpose\",\"status\":\"Available\",\"supportedFamilies\":[{\"name\":\"Gen5\",\"status\":\"Available\",\"sku\":\"explicit_sku\",\"supportedVcoresValues\":[{\"name\":\"c8\",\"value\":8,\"status\":\"Available\",\"standaloneSupported\":false,\"instancePoolSupported\":false}]}]}]}]}"

func sqlReportTestInput(family string) (map[string]string, SQLCapabilityRequest, string, string) {
	kind, key, body := "microsoft.sql/servers/databases", "12/basic/o1", sqlReportDBReady
	if family == "mi" {
		kind, key, body = "microsoft.sql/managedinstances", "12/generalpurpose/gen5/c8", sqlReportMIReady
	}
	return map[string]string{sqlReportOracleID: "selected"}, SQLCapabilityRequest{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: kind}, key, body
}
func sqlReportTestCollector(t *testing.T, family string, get RESTQuotaGetter) *SQLCapabilityCollector {
	t.Helper()
	var c *SQLCapabilityCollector
	var e error
	if family == "mi" {
		c, e = NewSQLManagedInstanceSKUCollector("https://management.azure.com", get)
	} else {
		c, e = NewSQLDatabaseSKUCollector("https://management.azure.com", get)
	}
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func sqlReportRequireFailure(t *testing.T, out assessment.PluginTable, e error, code string) {
	t.Helper()
	if e == nil || out.Rows == nil || len(out.Rows) != 0 || out.Health.Status != assessment.StageFailed || out.Health.Records != 0 || out.Health.Error == nil || out.Health.Error.Code != code || out.Health.Error.Message != "SQL capability report could not be safely projected" || out.Health.Warnings != nil {
		t.Fatal("atomic safe report failure", out, e)
	}
	if err := assessment.ValidatePluginTables([]assessment.PluginTable{out}); err != nil {
		t.Fatal(err)
	}
	if code != "sql_capability_report_cancelled" && e.Error() != "SQL capability report rejected ["+code+"]" {
		t.Fatal("fixed report error", e)
	}
}
func TestSQLReportOutputBoundary(t *testing.T) {
	for _, family := range []string{"db", "mi"} {
		for _, over := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", family, over), func(t *testing.T) {
				scope, q, c, required := sqlReportBudgetFixture(t, family, over)
				got, e := ProjectSQLCapabilityTable(context.Background(), scope, q, c, required)
				if over {
					sqlReportRequireFailure(t, got, e, "sql_capability_report_limit")
				} else if e != nil || len(got.Rows) != 8192 || got.Health.Records != 8192 || sqlReportOracleText(got) != 8388608 {
					t.Fatal("independent exact report output cap", len(got.Rows), sqlReportOracleText(got), e)
				}
			})
		}
	}
}
func TestSQLReportPreflightNoGET(t *testing.T) {
	for _, mode := range []string{"required-limit", "required-invalid", "required-collision", "aliases-limit", "invalid-alias", "selected-limit", "empty-scope", "wrong-family", "nil-collector"} {
		t.Run(mode, func(t *testing.T) {
			scope, q, key, body := sqlReportTestInput("db")
			required := []string{key}
			calls := 0
			col := sqlReportTestCollector(t, "db", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
				calls++
				return []byte(body), &http.Response{StatusCode: 200}, nil
			}))
			code := "sql_capability_report_invalid"
			switch mode {
			case "required-limit":
				required = make([]string, 8193)
				code = "sql_capability_report_limit"
			case "required-invalid":
				required = []string{"same_sku"}
			case "required-collision":
				required = []string{key, strings.ToUpper(key)}
			case "aliases-limit":
				q.LocationAliases = map[string]string{}
				for i := 0; i < 257; i++ {
					q.LocationAliases[fmt.Sprintf("reg ion%d", i)] = fmt.Sprintf("region%d", i)
				}
				code = "sql_capability_report_limit"
			case "invalid-alias":
				q.LocationAliases = map[string]string{"global": "eastus"}
			case "selected-limit":
				for i := 0; i < 1000; i++ {
					scope[fmt.Sprintf("%08x-2222-4222-8222-bbbbbbbbbbbb", i)] = "selected"
				}
				code = "sql_capability_report_limit"
			case "empty-scope":
				scope = map[string]string{}
			case "wrong-family":
				q.ResourceType = "microsoft.sql/managedinstances"
				required = []string{"12/generalpurpose/gen5/c8"}
			case "nil-collector":
				col = nil
			}
			out, e := CollectSQLCapabilityTable(context.Background(), col, scope, q, required)
			sqlReportRequireFailure(t, out, e, code)
			if calls != 0 {
				t.Fatal("report preflight must precede GET", calls)
			}
		})
	}
}
func TestSQLReportInheritedExactCounts(t *testing.T) {
	for _, which := range []string{"selected", "aliases", "label"} {
		for _, over := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", which, over), func(t *testing.T) {
				scope, q, _, _ := sqlReportTestInput("db")
				bound := 1000
				switch which {
				case "selected":
					n := bound
					if over {
						n++
					}
					for i := 1; i < n; i++ {
						scope[fmt.Sprintf("%08x-2222-4222-8222-bbbbbbbbbbbb", i)] = "selected"
					}
				case "aliases":
					bound = 256
					n := bound
					if over {
						n++
					}
					q.LocationAliases = map[string]string{}
					for i := 0; i < n; i++ {
						q.LocationAliases[fmt.Sprintf("reg ion%d", i)] = fmt.Sprintf("region%d", i)
					}
				case "label":
					bound = 512
					n := bound
					if over {
						n++
					}
					scope[sqlReportOracleID] = strings.Repeat("a", n)
				}
				out, e := ProjectSQLCapabilityTable(context.Background(), scope, q, nil, nil)
				if over {
					sqlReportRequireFailure(t, out, e, "sql_capability_report_limit")
				} else if e != nil || out.Health.Status != assessment.StageCompletedWithWarnings || len(out.Health.Warnings) != 2 || out.Rows == nil || len(out.Rows) != 0 {
					t.Fatal("exact inherited count/label before unknown", out, e)
				}
			})
		}
	}
}
func TestSQLReportFixedLateContext(t *testing.T) {
	// Manual pre-implementation freeze: Project=Catalog21+finish1=22.
	// Collect=entry1+nil preflight9+collector25+Project22=57. No calibration.
	for _, family := range []string{"db", "mi"} {
		scope, q, _, _ := sqlReportTestInput(family)
		c := &SQLCapabilityCollection{SubscriptionID: q.SubscriptionID, Region: q.Region, ResourceType: q.ResourceType, Status: "complete", Groups: map[string]SQLCapabilityGroup{}}
		body := "{\"name\":\"eastus\",\"status\":\"Available\",\"supportedServerVersions\":[]}"
		if family == "mi" {
			body = "{\"name\":\"eastus\",\"status\":\"Available\",\"supportedManagedInstanceVersions\":[]}"
		}
		for _, which := range []string{"project", "collect"} {
			at := 22
			if which == "collect" {
				at = 57
			}
			for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
				t.Run(family+which+cause.Error(), func(t *testing.T) {
					getter := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
						return []byte(body), &http.Response{StatusCode: 200}, nil
					})
					collector := sqlReportTestCollector(t, family, getter)
					ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 1 << 30, cause: cause}
					var out assessment.PluginTable
					var e error
					if which == "project" {
						out, e = ProjectSQLCapabilityTable(ctx, scope, q, c, nil)
					} else {
						out, e = CollectSQLCapabilityTable(ctx, collector, scope, q, nil)
					}
					if e != nil || ctx.calls != at {
						t.Fatal("fixed manually derived report schedule", which, ctx.calls, at, out, e)
					}
					ctx = &vmssCountingContext{Context: context.Background(), terminalAt: at, cause: cause}
					if which == "project" {
						out, e = ProjectSQLCapabilityTable(ctx, scope, q, c, nil)
					} else {
						out, e = CollectSQLCapabilityTable(ctx, collector, scope, q, nil)
					}
					sqlReportRequireFailure(t, out, e, "sql_capability_report_cancelled")
					if e != cause {
						t.Fatal("canonical final report context", e, cause, ctx.calls)
					}
				})
			}
		}
	}
}
func TestSQLReportWrappedGetterContext(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		scope, q, key, _ := sqlReportTestInput("db")
		calls := 0
		col := sqlReportTestCollector(t, "db", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return []byte("provider canary"), &http.Response{StatusCode: 503}, fmt.Errorf("provider canary: %w", cause)
		}))
		out, e := CollectSQLCapabilityTable(context.Background(), col, scope, q, []string{key})
		sqlReportRequireFailure(t, out, e, "sql_capability_report_cancelled")
		if e != cause || calls != 1 {
			t.Fatal("wrapped terminal priority", out, e, calls)
		}
	}
}
func TestSQLReportTransportSnapshot(t *testing.T) {
	scope, q, key, body := sqlReportTestInput("db")
	q.LocationAliases = map[string]string{"east us": "eastus"}
	required := []string{key}
	body = strings.Replace(body, "eastus", "East US", 1)
	calls := 0
	col := sqlReportTestCollector(t, "db", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		delete(scope, sqlReportOracleID)
		q.LocationAliases["east us"] = "westus"
		required[0] = "12/basic/other"
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	out, e := CollectSQLCapabilityTable(context.Background(), col, scope, q, required)
	if e != nil || calls != 1 || !reflect.DeepEqual(out, sqlReportOracleCases()[2].want) {
		t.Fatal("validated request snapshot before getter", out, e, calls)
	}
}
func TestSQLReportCanonicalBinding(t *testing.T) {
	c := sqlReportOracleCases()[2]
	c.request.SubscriptionID = strings.ToUpper(c.request.SubscriptionID)
	c.request.ResourceType = strings.ToUpper(c.request.ResourceType)
	c.selected = map[string]string{strings.ToUpper(sqlReportOracleID): "selected"}
	out, e := ProjectSQLCapabilityTable(context.Background(), c.selected, c.request, c.collection, c.required)
	if e != nil || !reflect.DeepEqual(out, c.want) {
		t.Fatal("canonical report binding", out, e)
	}
}
func TestSQLReportConcurrentOwnership(t *testing.T) {
	c := sqlReportOracleCases()[2]
	before := sqlReportInputSnapshot(t, c)
	first, e := ProjectSQLCapabilityTable(context.Background(), c.selected, c.request, c.collection, c.required)
	if e != nil {
		t.Fatal(e)
	}
	var calls atomic.Int32
	col := sqlReportTestCollector(t, "db", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls.Add(1)
		return []byte(sqlReportDBReady), &http.Response{StatusCode: 200}, nil
	}))
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, collect := range []bool{false, true} {
				var out assessment.PluginTable
				var err error
				if collect {
					out, err = CollectSQLCapabilityTable(context.Background(), col, c.selected, c.request, c.required)
				} else {
					out, err = ProjectSQLCapabilityTable(context.Background(), c.selected, c.request, c.collection, c.required)
				}
				if err != nil || !reflect.DeepEqual(out, c.want) {
					t.Error("concurrent report", out, err)
					return
				}
				out.Columns[0] = "mutated"
				out.Rows[0].Cells[0] = "mutated"
				out.Health.Warnings[0].Message = "mutated"
			}
		}()
	}
	wg.Wait()
	if calls.Load() != 12 || !reflect.DeepEqual(first, c.want) || !reflect.DeepEqual(before, sqlReportInputSnapshot(t, c)) {
		t.Fatal("owned report state", calls.Load(), first)
	}
}
func TestSQLReportActualHTTPClosure(t *testing.T) {
	for _, family := range []string{"db", "mi"} {
		for _, mode := range []string{"success", "403", "429", "500", "307", "oversize", "read-error", "cancel"} {
			t.Run(family+mode, func(t *testing.T) {
				scope, q, key, body := sqlReportTestInput(family)
				scope = map[string]string{sqlTestID: "selected"}
				q.SubscriptionID = sqlTestID
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				credential := &diskCanaryCredential{}
				wire := &sqlHTTPTransport{t: t, status: 200, body: body, mode: mode, cancel: cancel}
				switch mode {
				case "403":
					wire.status = 403
				case "429":
					wire.status = 429
				case "500":
					wire.status = 500
				case "307":
					wire.status = 307
				case "oversize":
					wire.body = strings.Repeat("x", 1048577)
				}
				client := azure.NewHTTPClient(credential, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, OperationTimeout: 5 * time.Second, Scope: "https://management.azure.com/.default", Transport: wire})
				col := sqlReportTestCollector(t, family, client)
				out, e := CollectSQLCapabilityTable(ctx, col, scope, q, []string{key})
				if wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 || !reflect.DeepEqual(credential.scopes, []string{"https://management.azure.com/.default"}) {
					t.Fatal("actual report HTTP closure/audience", mode)
				}
				if mode == "cancel" {
					sqlReportRequireFailure(t, out, e, "sql_capability_report_cancelled")
					if !errors.Is(e, context.Canceled) {
						t.Fatal(e)
					}
				} else if e != nil {
					t.Fatal(out, e)
				} else if mode == "success" {
					if len(out.Rows) != 1 || out.Rows[0].SubscriptionID != sqlTestID {
						t.Fatal(out)
					}
				} else if len(out.Rows) != 0 || len(out.Health.Warnings) != 2 || out.Health.Warnings[1].Code != "sql_capability_catalog_unknown" {
					t.Fatal("unknown HTTP report has no prefix", out)
				}
				raw, err := json.Marshal(out)
				if err != nil || strings.Contains(string(raw), "provider canary") {
					t.Fatal("safe report serialization", err)
				}
			})
		}
	}
}
