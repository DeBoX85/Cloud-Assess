package region

import (
	"context"
	"errors"
	"fmt"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

const sqlMIEmptyBody = `{"name":"eastus","status":"Available","supportedManagedInstanceVersions":[]}`
const sqlMIReadyBody = `{"name":"eastus","status":"Available","supportedManagedInstanceVersions":[{"name":"12.0","status":"Available","supportedEditions":[{"name":"GeneralPurpose","status":"Available","supportedFamilies":[{"name":"Gen5","status":"Available","sku":"Explicit_Hardware42","supportedVcoresValues":[{"name":"c8","value":8,"status":"Available","standaloneSupported":false,"instancePoolSupported":true}]}]}]}]}`
const sqlMIKey = "12.0/generalpurpose/gen5/c8"

func sqlMIInput() (map[string]string, SQLCapabilityRequest, *SQLCapabilityCollection) {
	return map[string]string{sqlTestID: "selected"}, SQLCapabilityRequest{SubscriptionID: sqlTestID, Region: "eastus", ResourceType: "microsoft.sql/managedinstances"}, &SQLCapabilityCollection{SubscriptionID: sqlTestID, Region: "eastus", ResourceType: "microsoft.sql/managedinstances", Status: "complete", Groups: map[string]SQLCapabilityGroup{}}
}
func sqlMICollectBody(t *testing.T, body string) *SQLCapabilityCollection {
	t.Helper()
	scope, q, _ := sqlMIInput()
	calls := 0
	c, e := NewSQLManagedInstanceSKUCollector("https://management.azure.com", storageGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
		calls++
		if calls != 1 || u != sqlTestURL || limit != 1048576 {
			t.Fatal("one bounded MI GET", calls, u, limit)
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
func TestSQLMIWireBound(t *testing.T) {
	for _, n := range []int{1048576, 1048577} {
		r := sqlMICollectBody(t, sqlMIEmptyBody+strings.Repeat(" ", n-len(sqlMIEmptyBody)))
		if n == 1048576 {
			if r.Status != "complete" {
				t.Fatal("exact MI wire", r)
			}
		} else {
			sqlRequireFailure(t, r, "byte_limit")
		}
	}
}
func TestSQLMIRawWorkBound(t *testing.T) {
	for _, extra := range []int{0, 1} {
		body := strings.TrimSuffix(sqlMIEmptyBody, "}") + `,"padding":[` + strings.Repeat("{},", 32767) + "null" + strings.Repeat(",null", extra) + "]}"
		if len(body) > 1048576 {
			t.Fatal("masked MI work fixture")
		}
		r := sqlMICollectBody(t, body)
		if extra == 0 {
			if r.Status != "complete" {
				t.Fatal("exact raw MI work", r)
			}
		} else {
			sqlRequireFailure(t, r, "work_limit")
		}
	}
}
func TestSQLMILeafBound(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		leaves := make([]string, n)
		for i := range leaves {
			leaves[i] = fmt.Sprintf(`{"name":"c%04d"}`, i)
		}
		body := `{"name":"eastus","supportedManagedInstanceVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedFamilies":[{"name":"f","sku":"s","supportedVcoresValues":[` + strings.Join(leaves, ",") + `]}]}]}]}`
		r := sqlMICollectBody(t, body)
		if n == 8192 {
			if r.Status != "complete" || len(r.Groups) != 8192 {
				t.Fatal("exact independent MI leaf bound", r.Status, len(r.Groups), r.FailureCode)
			}
		} else {
			sqlRequireFailure(t, r, "row_limit")
		}
	}
}
func TestSQLMIStrictWire(t *testing.T) {
	for name, body := range map[string]string{
		"depth":             strings.TrimSuffix(sqlMIEmptyBody, "}") + `,"padding":` + strings.Repeat("[", 34) + "0" + strings.Repeat("]", 34) + "}",
		"tokens":            strings.TrimSuffix(sqlMIEmptyBody, "}") + `,"padding":[` + strings.Repeat("null,", 65536) + "null]}",
		"label":             strings.TrimSuffix(sqlMIEmptyBody, "}") + `,"padding":"` + strings.Repeat("a", 513) + `"}`,
		"folded-child":      strings.Replace(sqlMIReadyBody, `"value":8`, `"Value":8`, 1),
		"trailing":          sqlMIReadyBody + "{}",
		"invalid-utf8":      sqlMIEmptyBody + string([]byte{0xff}),
		"bad-unselected-db": strings.TrimSuffix(sqlMIEmptyBody, "}") + `,"supportedServerVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedServiceLevelObjectives":[{"name":"s","sku":{"name":"s","size":null}}]}]}]}`,
		"zero-value":        strings.Replace(sqlMIReadyBody, `"value":8`, `"value":0`, 1),
		"bad-value":         strings.Replace(sqlMIReadyBody, `"value":8`, `"value":2147483648`, 1),
		"bad-sku-ascii":     strings.Replace(sqlMIReadyBody, "Explicit_Hardware42", "BAD/SKU", 1),
		"bad-mode":          strings.Replace(sqlMIReadyBody, `"standaloneSupported":false`, `"standaloneSupported":"false"`, 1),
	} {
		t.Run(name, func(t *testing.T) { sqlRequireFailure(t, sqlMICollectBody(t, body), "invalid_response") })
	}
}
func TestSQLMIAuthorityAndModes(t *testing.T) {
	scope, q, _ := sqlMIInput()
	r := sqlMICollectBody(t, sqlMIReadyBody)
	g := r.Groups[sqlMIKey]
	want := SQLCapabilityGroup{Known: true, State: 0, Version: "12.0", Edition: "generalpurpose", Name: "c8", Family: "gen5", SKUName: "explicit_hardware42", SKUTier: "generalpurpose", Capacity: 8, Standalone: "unsupported", InstancePool: "supported"}
	if r.Status != "complete" || !reflect.DeepEqual(g, want) {
		t.Fatal("no MI SKU/name reconstruction", g, want)
	}
	out, e := r.Catalog(context.Background(), scope, q, []string{sqlMIKey})
	if e != nil || out.DeploymentEligibility != "not_assessed" || len(out.Warnings) != 1 || out.Values[sqlMIKey] != want {
		t.Fatal("catalog-only MI modes", out, e)
	}
	for _, mode := range []string{"standaloneSupported", "instancePoolSupported"} {
		body := sqlMIReadyBody
		value := "false"
		if mode == "instancePoolSupported" {
			value = "true"
		}
		body = strings.Replace(body, `,"`+mode+`":`+value, "", 1)
		r := sqlMICollectBody(t, body)
		g := r.Groups[sqlMIKey]
		flag := g.Standalone
		if mode == "instancePoolSupported" {
			flag = g.InstancePool
		}
		if !g.Known || g.State != 0 || flag != "unknown" {
			t.Fatal("missing mode is unknown without invented eligibility", mode, g)
		}
	}
	both := strings.Replace(sqlMIReadyBody, `"instancePoolSupported":true`, `"instancePoolSupported":false`, 1)
	g = sqlMICollectBody(t, both).Groups[sqlMIKey]
	if !g.Known || g.State != 0 || g.Standalone != "unsupported" || g.InstancePool != "unsupported" {
		t.Fatal("no silent mode selection", g)
	}
	for _, missing := range []string{"sku", "value"} {
		body := sqlMIReadyBody
		if missing == "sku" {
			body = strings.Replace(body, `"sku":"Explicit_Hardware42",`, "", 1)
		} else {
			body = strings.Replace(body, `"value":8,`, "", 1)
		}
		r := sqlMICollectBody(t, body)
		g := r.Groups[sqlMIKey]
		if g.Known || g.State != 2 || missing == "sku" && g.SKUName != "" || missing == "value" && g.Capacity != 0 {
			t.Fatal("missing authoritative field, no fallback", missing, g)
		}
		out, e := r.Catalog(context.Background(), scope, q, []string{sqlMIKey})
		if e != nil || out.Status != "unknown" || len(out.Values) != 0 {
			t.Fatal("required incomplete MI leaf", out, e)
		}
	}
}
func TestSQLMIStatusAndIdentity(t *testing.T) {
	for _, mode := range []string{"unknown-beats-disabled", "family-disabled"} {
		body := sqlMIReadyBody
		if mode == "unknown-beats-disabled" {
			body = strings.Replace(body, `"name":"eastus","status":"Available"`, `"name":"eastus","status":"Disabled"`, 1)
			body = strings.Replace(body, `"name":"Gen5","status":"Available"`, `"name":"Gen5","status":"Future"`, 1)
		} else {
			body = strings.Replace(body, `"name":"Gen5","status":"Available"`, `"name":"Gen5","status":"Disabled"`, 1)
		}
		r := sqlMICollectBody(t, body)
		g := r.Groups[sqlMIKey]
		if r.Status != "complete" || g.State != 2 || g.Known != (mode == "family-disabled") {
			t.Fatal("MI governing family status", mode, g)
		}
	}
	duplicate := `{"name":"eastus","supportedManagedInstanceVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedFamilies":[{"name":"f","supportedVcoresValues":[{"name":"c1"}]},{"name":"F","supportedVcoresValues":[{"name":"c2"}]}]}]}]}`
	sqlRequireFailure(t, sqlMICollectBody(t, duplicate), "collision")
}
func TestSQLMIFamilyPreauth(t *testing.T) {
	for _, mode := range []string{"constructor-db-family", "unselected-empty-label", "nil-collector", "nil-getter", "bad-origin"} {
		scope, q, _ := sqlMIInput()
		calls := 0
		get := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return []byte(sqlMIEmptyBody), &http.Response{StatusCode: 200}, nil
		})
		collector, e := NewSQLManagedInstanceSKUCollector("https://management.azure.com", get)
		if e != nil {
			t.Fatal(e)
		}
		switch mode {
		case "constructor-db-family":
			q.ResourceType = "microsoft.sql/servers/databases"
		case "unselected-empty-label":
			scope["bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"] = ""
		case "nil-collector":
			collector = nil
		case "nil-getter":
			if c, e := NewSQLManagedInstanceSKUCollector("https://management.azure.com", nil); c != nil || e == nil {
				t.Fatal("MI getter validation")
			}
			continue
		case "bad-origin":
			if c, e := NewSQLManagedInstanceSKUCollector("https://private.invalid/path", get); c != nil || e == nil || e.Error() != "invalid SQL capability origin or getter" {
				t.Fatal("MI origin validation", c, e)
			}
			continue
		}
		r, e := collector.Collect(context.Background(), scope, q)
		if r != nil || e == nil || e.Error() != "SQL capability request rejected [sql_capability_request_invalid]" || calls != 0 {
			t.Fatal("MI preauth family/scope", mode, r, e, calls)
		}
	}
}
func TestSQLMIFixedLateContext(t *testing.T) {
	// Manually frozen before MI reversals: shared Catalog nil9/complete21; empty
	// MI Collect entry/scope/pre/postGET4 + JSON4 + raw10 + shape5 + decode1 + finish1=25.
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		scope, q, c := sqlMIInput()
		for _, test := range []struct {
			name  string
			at    int
			value *SQLCapabilityCollection
		}{{"nil", 9, nil}, {"complete", 21, c}} {
			t.Run(test.name+cause.Error(), func(t *testing.T) {
				ctx := &vmssCountingContext{Context: context.Background(), terminalAt: test.at, cause: cause}
				out, e := test.value.Catalog(ctx, scope, q, nil)
				if e != cause || !reflect.DeepEqual(out, SQLCatalogEvidence{}) {
					t.Fatal("fixed MI catalog final context", ctx.calls, out, e)
				}
			})
		}
		collector, _ := NewSQLManagedInstanceSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return []byte(sqlMIEmptyBody), &http.Response{StatusCode: 200}, nil
		}))
		ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 25, cause: cause}
		r, e := collector.Collect(ctx, scope, q)
		if r != nil || e != cause {
			t.Fatal("fixed MI collector final context", ctx.calls, r, e)
		}
	}
}
func TestSQLMIActualHTTPClosure(t *testing.T) {
	for _, mode := range []string{"success", "403", "429", "500", "307", "oversize", "read-error", "cancel"} {
		scope, q, _ := sqlMIInput()
		ctx, cancel := context.WithCancel(context.Background())
		credential := &diskCanaryCredential{}
		wire := &sqlHTTPTransport{t: t, status: 200, body: sqlMIReadyBody, mode: mode, cancel: cancel}
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
		collector, _ := NewSQLManagedInstanceSKUCollector("https://management.azure.com", client)
		r, e := collector.Collect(ctx, scope, q)
		cancel()
		if wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 || !reflect.DeepEqual(credential.scopes, []string{"https://management.azure.com/.default"}) {
			t.Fatal("MI HTTP closure/audience", mode)
		}
		if mode == "cancel" {
			if r != nil || !errors.Is(e, context.Canceled) {
				t.Fatal("MI HTTP cancellation", r, e)
			}
		} else if r == nil || e != nil {
			t.Fatal(mode, r, e)
		} else if mode == "success" {
			if r.Status != "complete" || r.Groups[sqlMIKey].State != 0 {
				t.Fatal(mode, r)
			}
		} else if mode == "read-error" || mode == "oversize" {
			sqlRequireFailure(t, r, "request_failed")
		} else {
			sqlRequireFailure(t, r, "status_invalid")
		}
	}
}
func TestSQLMIOwnedConcurrent(t *testing.T) {
	scope, q, _ := sqlMIInput()
	collector, _ := NewSQLManagedInstanceSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(sqlMIReadyBody), &http.Response{StatusCode: 200}, nil
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
				t.Error("MI concurrent collection", e)
				return
			}
			out, e := r.Catalog(context.Background(), scope, q, []string{sqlMIKey})
			if e != nil || len(out.Values) != 1 {
				t.Error("MI concurrent projection", e)
				return
			}
			delete(out.Values, sqlMIKey)
			out.Warnings[0].Message = "mutated"
			delete(r.Groups, sqlMIKey)
		}()
	}
	wg.Wait()
	if !reflect.DeepEqual(before, vmssSnapshot(t, first)) {
		t.Fatal("shared MI ownership")
	}
}
func TestSQLMIAliasCopy(t *testing.T) {
	scope, q, _ := sqlMIInput()
	q.LocationAliases = map[string]string{"east us": "eastus"}
	collector, _ := NewSQLManagedInstanceSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		q.LocationAliases["east us"] = "westus"
		return []byte(strings.Replace(sqlMIEmptyBody, "eastus", "East US", 1)), &http.Response{StatusCode: 200}, nil
	}))
	r, e := collector.Collect(context.Background(), scope, q)
	if e != nil || r == nil || r.Status != "complete" {
		t.Fatal("MI alias copy", r, e)
	}
}

func TestSQLMIUnselectedRawRows(t *testing.T) {
	rows := make([]string, 9000)
	for i := range rows {
		rows[i] = fmt.Sprintf(`{"name":"s%04d"}`, i)
	}
	body := strings.TrimSuffix(sqlMIEmptyBody, "}") + `,"supportedServerVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedServiceLevelObjectives":[` + strings.Join(rows, ",") + `]}]}]}`
	r := sqlMICollectBody(t, body)
	if r.Status != "complete" || len(r.Groups) != 0 {
		t.Fatal("unselected Database rows are bounded/validated but not MI leaves", r.Status, r.FailureCode)
	}
}
func TestSQLDatabaseUnselectedMIRawRows(t *testing.T) {
	rows := make([]string, 9000)
	for i := range rows {
		rows[i] = fmt.Sprintf(`{"name":"c%04d"}`, i)
	}
	body := strings.TrimSuffix(sqlEmptyBody, "}") + `,"supportedManagedInstanceVersions":[{"name":"v","supportedEditions":[{"name":"e","supportedFamilies":[{"name":"f","supportedVcoresValues":[` + strings.Join(rows, ",") + `]}]}]}]}`
	r := sqlCollectBody(t, body)
	if r.Status != "complete" || len(r.Groups) != 0 {
		t.Fatal("unselected MI rows must preserve Database leaf semantics", r.Status, r.FailureCode)
	}
}
