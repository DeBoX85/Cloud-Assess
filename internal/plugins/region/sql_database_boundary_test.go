package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const sqlTestID = "abcdefab-1234-4234-8234-123456789abc"
const sqlTestURL = "https://management.azure.com/subscriptions/abcdefab-1234-4234-8234-123456789abc/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
const sqlEmptyBody = `{"name":"eastus","status":"Available","supportedServerVersions":[]}`
const sqlReadyBody = `{"name":"eastus","status":"Available","supportedServerVersions":[{"name":"12.0","status":"Available","supportedEditions":[{"name":"Standard","status":"Available","supportedServiceLevelObjectives":[{"name":"S0","status":"Available","sku":{"name":"S0","tier":"Standard","capacity":10}}]}]}]}`

func sqlBoundaryInput() (map[string]string, SQLCapabilityRequest, *SQLCapabilityCollection) {
	return map[string]string{sqlTestID: "selected"}, SQLCapabilityRequest{SubscriptionID: sqlTestID, Region: "eastus", ResourceType: "microsoft.sql/servers/databases"}, &SQLCapabilityCollection{SubscriptionID: sqlTestID, Region: "eastus", ResourceType: "microsoft.sql/servers/databases", Status: "complete", Groups: map[string]SQLCapabilityGroup{}}
}
func sqlCollectBody(t *testing.T, body string) *SQLCapabilityCollection {
	t.Helper()
	scope, q, _ := sqlBoundaryInput()
	calls := 0
	c, e := NewSQLDatabaseSKUCollector("https://management.azure.com", storageGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
		calls++
		if calls != 1 || u != sqlTestURL || limit != 1048576 {
			t.Fatal("one bounded SQL GET", calls, u, limit)
		}
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	if e != nil {
		t.Fatal(e)
	}
	out, e := c.Collect(context.Background(), scope, q)
	if e != nil || out == nil || calls != 1 {
		t.Fatal(out, e, calls)
	}
	return out
}
func sqlRequireFailure(t *testing.T, r *SQLCapabilityCollection, code string) {
	t.Helper()
	if r.Status != "unknown" || len(r.Groups) != 0 || r.FailureCode != "sql_capability_"+code {
		t.Fatal("fixed whole-body rejection required", r.Status, len(r.Groups), r.FailureCode, code)
	}
}
func TestSQLDatabaseWireBound(t *testing.T) {
	for _, n := range []int{1048576, 1048577} {
		r := sqlCollectBody(t, sqlEmptyBody+strings.Repeat(" ", n-len(sqlEmptyBody)))
		if n == 1048576 {
			if r.Status != "complete" {
				t.Fatal("exact wire", r)
			}
		} else {
			sqlRequireFailure(t, r, "byte_limit")
		}
	}
}
func TestSQLDatabaseRawWorkBound(t *testing.T) {
	for _, extra := range []int{0, 1} {
		body := strings.TrimSuffix(sqlEmptyBody, "}") + `,"padding":[` + strings.Repeat("{},", 32767) + "null" + strings.Repeat(",null", extra) + "]}"
		if len(body) > 1048576 {
			t.Fatal("masked work fixture")
		}
		r := sqlCollectBody(t, body)
		if extra == 0 {
			if r.Status != "complete" {
				t.Fatal("exact raw work", r)
			}
		} else {
			sqlRequireFailure(t, r, "work_limit")
		}
	}
}
func TestSQLDatabaseLeafBound(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		leaves := make([]string, n)
		for i := range leaves {
			leaves[i] = fmt.Sprintf(`{"name":"s%04d"}`, i)
		}
		body := `{"name":"eastus","supportedServerVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedServiceLevelObjectives":[` + strings.Join(leaves, ",") + `]}]}]}`
		r := sqlCollectBody(t, body)
		if n == 8192 {
			if r.Status != "complete" || len(r.Groups) != 8192 {
				t.Fatal("exact independent leaf bound", r.Status, len(r.Groups), r.FailureCode)
			}
		} else {
			sqlRequireFailure(t, r, "row_limit")
		}
	}
}
func TestSQLDatabaseStrictWireBoundaries(t *testing.T) {
	for name, body := range map[string]string{
		"depth":                       strings.TrimSuffix(sqlEmptyBody, "}") + `,"padding":` + strings.Repeat("[", 34) + "0" + strings.Repeat("]", 34) + "}",
		"tokens":                      strings.TrimSuffix(sqlEmptyBody, "}") + `,"padding":[` + strings.Repeat("null,", 65536) + "null]}",
		"label":                       strings.TrimSuffix(sqlEmptyBody, "}") + `,"padding":"` + strings.Repeat("a", 513) + `"}`,
		"folded-child":                strings.Replace(sqlReadyBody, `"capacity":10`, `"Capacity":10`, 1),
		"trailing":                    sqlReadyBody + "{}",
		"invalid-utf8":                sqlEmptyBody + string([]byte{0xff}),
		"malformed-unselected":        strings.TrimSuffix(sqlEmptyBody, "}") + `,"supportedManagedInstanceVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedFamilies":[{"name":"f","supportedVcoresValues":[{"name":"c","value":2147483648}]}]}]}]}`,
		"malformed-unselected-double": strings.TrimSuffix(sqlEmptyBody, "}") + `,"supportedManagedInstanceVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedStorageCapabilities":[{"status":null}]}]}]}`,
	} {
		t.Run(name, func(t *testing.T) { sqlRequireFailure(t, sqlCollectBody(t, body), "invalid_response") })
	}
}
func TestSQLDatabaseStatusAndIdentity(t *testing.T) {
	for name, body := range map[string]string{"unknown-beats-disabled": strings.Replace(strings.Replace(sqlReadyBody, `"status":"Available"`, `"status":"Disabled"`, 1), `"name":"S0","status":"Available"`, `"name":"S0","status":"Future"`, 1), "disabled-ancestor": strings.Replace(sqlReadyBody, `"status":"Available"`, `"status":"Disabled"`, 1)} {
		r := sqlCollectBody(t, body)
		g := r.Groups["12.0/standard/s0"]
		if r.Status != "complete" || g.State != 2 || g.Known != (name == "disabled-ancestor") {
			t.Fatal("ancestor status contract", name, r)
		}
	}
	duplicate := `{"name":"eastus","supportedServerVersions":[{"name":"v","supportedEditions":[{"name":"e1","supportedServiceLevelObjectives":[]}]},{"name":"V","supportedEditions":[{"name":"e2","supportedServiceLevelObjectives":[]}]}]}`
	sqlRequireFailure(t, sqlCollectBody(t, duplicate), "collision")
	malformed := strings.Replace(strings.TrimSuffix(sqlReadyBody, "}"), `"capacity":10`, `"capacity":10`, 1) + `,"supportedManagedInstanceVersions":[{"name":"v","status":false}]}`
	sqlRequireFailure(t, sqlCollectBody(t, malformed), "invalid_response")
}
func TestSQLDatabaseScopeAliasPreauth(t *testing.T) {
	for _, mode := range []string{"nil-collector", "bad-origin", "nil-getter", "empty-selected", "label", "id", "alias", "family", "region"} {
		scope, q, _ := sqlBoundaryInput()
		calls := 0
		get := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return []byte(sqlEmptyBody), &http.Response{StatusCode: 200}, nil
		})
		c, e := NewSQLDatabaseSKUCollector("https://management.azure.com", get)
		if e != nil {
			t.Fatal(e)
		}
		switch mode {
		case "nil-collector":
			c = nil
		case "bad-origin":
			if x, e := NewSQLDatabaseSKUCollector("https://private.invalid/path", get); x != nil || e == nil || e.Error() != "invalid SQL capability origin or getter" {
				t.Fatal("origin rejection", x, e)
			}
			continue
		case "nil-getter":
			if x, e := NewSQLDatabaseSKUCollector("https://management.azure.com", nil); x != nil || e == nil {
				t.Fatal("getter rejection")
			}
			continue
		case "empty-selected":
			scope = map[string]string{}
		case "label":
			scope[sqlTestID] = ""
		case "id":
			q.SubscriptionID = "bad"
		case "alias":
			q.LocationAliases = map[string]string{"east us": "westus"}
		case "family":
			q.ResourceType = "microsoft.sql/managedinstances"
		case "region":
			q.Region = "EASTUS"
		}
		r, e := c.Collect(context.Background(), scope, q)
		if r != nil || e == nil || e.Error() != "SQL capability request rejected [sql_capability_request_invalid]" || calls != 0 {
			t.Fatal("preauth request boundary", mode, r, e, calls)
		}
	}
}
func TestSQLDatabaseFixedLateContext(t *testing.T) {
	// Manually fixed BEFORE compiling reversals: nil Catalog entry1, selected3,
	// scope1, expected binding3, finish1=9. Complete empty Catalog adds5 collection
	// strings and7 output strings, final=21. Empty Collect entry/scope/pre/postGET4,
	// strict JSON4, raw walk10, whole-shape5, decode finish1, collector finish1=25.
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		scope, q, c := sqlBoundaryInput()
		for _, test := range []struct {
			name  string
			at    int
			value *SQLCapabilityCollection
		}{{"nil-catalog", 9, nil}, {"complete-catalog", 21, c}} {
			t.Run(test.name+cause.Error(), func(t *testing.T) {
				ctx := &vmssCountingContext{Context: context.Background(), terminalAt: test.at, cause: cause}
				r, e := test.value.Catalog(ctx, scope, q, nil)
				if e != cause || !reflect.DeepEqual(r, SQLCatalogEvidence{}) {
					t.Fatal("fixed catalog final-return priority", ctx.calls, r, e)
				}
			})
		}
		collector, _ := NewSQLDatabaseSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return []byte(sqlEmptyBody), &http.Response{StatusCode: 200}, nil
		}))
		ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 25, cause: cause}
		r, e := collector.Collect(ctx, scope, q)
		if r != nil || e != cause {
			t.Fatal("fixed collector final-return priority", ctx.calls, r, e)
		}
	}
}

type sqlHTTPTransport struct {
	t             *testing.T
	status        int
	body          string
	mode          string
	cancel        context.CancelFunc
	calls, closes atomic.Int32
}

func (c *sqlHTTPTransport) Do(req *http.Request) (*http.Response, error) {
	c.calls.Add(1)
	if req.Method != "GET" || req.URL.String() != sqlTestURL || req.Header.Get("Authorization") != "Bearer OFFLINE_CANARY" || req.Body != nil && req.Body != http.NoBody {
		c.t.Error("SQL HTTP boundary")
	}
	var reader io.Reader = strings.NewReader(c.body)
	if c.mode == "read-error" {
		reader = cognitiveErrorReader{}
	}
	if c.mode == "cancel" {
		reader = cognitiveCancelReader{cancel: c.cancel}
	}
	return &http.Response{StatusCode: c.status, Header: make(http.Header), Body: &diskCloseBody{Reader: reader, closes: &c.closes}, Request: req}, nil
}
func TestSQLDatabaseActualHTTPClosure(t *testing.T) {
	for _, mode := range []string{"success", "403", "429", "500", "307", "oversize", "read-error", "cancel"} {
		scope, q, _ := sqlBoundaryInput()
		ctx, cancel := context.WithCancel(context.Background())
		credential := &diskCanaryCredential{}
		wire := &sqlHTTPTransport{t: t, status: 200, body: sqlReadyBody, mode: mode, cancel: cancel}
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
		collector, _ := NewSQLDatabaseSKUCollector("https://management.azure.com", client)
		r, e := collector.Collect(ctx, scope, q)
		cancel()
		if wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 || !reflect.DeepEqual(credential.scopes, []string{"https://management.azure.com/.default"}) {
			t.Fatal("actual HTTP closure/audience", mode, wire.calls.Load(), wire.closes.Load())
		}
		if mode == "cancel" {
			if r != nil || !errors.Is(e, context.Canceled) {
				t.Fatal("HTTP cancellation", r, e)
			}
		} else if r == nil || e != nil {
			t.Fatal(mode, r, e)
		} else if mode == "success" {
			if r.Status != "complete" || r.Groups["12.0/standard/s0"].State != 0 {
				t.Fatal(mode, r)
			}
		} else if mode == "read-error" || mode == "oversize" {
			sqlRequireFailure(t, r, "request_failed")
		} else {
			sqlRequireFailure(t, r, "status_invalid")
		}
	}
}
func TestSQLDatabaseOwnedConcurrent(t *testing.T) {
	scope, q, _ := sqlBoundaryInput()
	q.LocationAliases = map[string]string{"east us": "eastus"}
	collector, _ := NewSQLDatabaseSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(sqlReadyBody), &http.Response{StatusCode: 200}, nil
	}))
	first, e := collector.Collect(context.Background(), scope, q)
	if e != nil {
		t.Fatal(e)
	}
	before := vmssSnapshot(t, first)
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := collector.Collect(context.Background(), scope, q)
			if e != nil || !reflect.DeepEqual(r, first) {
				t.Error("concurrent collection", e)
				return
			}
			out, e := r.Catalog(context.Background(), scope, q, []string{"12.0/standard/s0"})
			if e != nil || len(out.Values) != 1 {
				t.Error("concurrent projection", e)
				return
			}
			delete(out.Values, "12.0/standard/s0")
			out.Warnings[0].Message = "mutated"
			delete(r.Groups, "12.0/standard/s0")
		}()
	}
	wg.Wait()
	if !reflect.DeepEqual(before, vmssSnapshot(t, first)) {
		t.Fatal("shared output")
	}
}
func TestSQLDatabaseAliasCopy(t *testing.T) {
	scope, q, _ := sqlBoundaryInput()
	q.LocationAliases = map[string]string{"east us": "eastus"}
	collector, _ := NewSQLDatabaseSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		q.LocationAliases["east us"] = "westus"
		return []byte(strings.Replace(sqlEmptyBody, "eastus", "East US", 1)), &http.Response{StatusCode: 200}, nil
	}))
	r, e := collector.Collect(context.Background(), scope, q)
	if e != nil || r == nil || r.Status != "complete" {
		t.Fatal("request alias copy", r, e)
	}
}
func TestSQLDatabaseShapeRegistry(t *testing.T) {
	var primary struct {
		Definitions map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"definitions"`
	}
	if json.Unmarshal(vmRead(t, "sql-primary-2021-11-01-schema.json"), &primary) != nil || len(primary.Definitions) != 26 || len(sqlCapabilityShapes) != 27 {
		t.Fatal("pinned shape inventory")
	}
	for name, def := range primary.Definitions {
		fields, ok := sqlCapabilityShapes[name]
		if !ok || len(fields) != len(def.Properties) {
			t.Fatal("incomplete pinned shape", name)
		}
		for key := range def.Properties {
			if _, ok := fields[key]; !ok {
				t.Fatal("missing pinned field", name, key)
			}
			raw := json.RawMessage(fmt.Sprintf(`{%q:null}`, key))
			if sqlCapabilityShape(context.Background(), raw, name) {
				t.Fatal("known optional null accepted", name, key)
			}
		}
	}
	if sqlCapabilityShape(context.Background(), json.RawMessage(`{}`), "Sku") {
		t.Fatal("required SKU name")
	}
}

func TestSQLDatabaseUnselectedEmptyLabel(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	scope["bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"] = ""
	calls := 0
	collector, _ := NewSQLDatabaseSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(sqlEmptyBody), &http.Response{StatusCode: 200}, nil
	}))
	r, e := collector.Collect(context.Background(), scope, q)
	if r != nil || e == nil || calls != 0 {
		t.Fatal("unselected empty label at request", r, e, calls)
	}
	out, e := c.Catalog(context.Background(), scope, q, nil)
	if e == nil || !reflect.DeepEqual(out, SQLCatalogEvidence{}) {
		t.Fatal("unselected empty label at projection", out, e)
	}
}
func TestSQLDatabaseProjectionBinding(t *testing.T) {
	scope, q, _ := sqlBoundaryInput()
	q.ResourceType = "microsoft.sql/managedinstances"
	var c *SQLCapabilityCollection
	out, e := c.Catalog(context.Background(), scope, q, nil)
	if e == nil || !reflect.DeepEqual(out, SQLCatalogEvidence{}) {
		t.Fatal("nil catalog still validates family", out, e)
	}
}
