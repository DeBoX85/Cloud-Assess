package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const storageReadyBody = `{"value":[{"resourceType":"storageAccounts","name":"s","kind":"Storage","locations":["eastus"],"restrictions":[]}]}`

func storageCollectBody(t *testing.T, body string) *StorageSKUCollection {
	t.Helper()
	getter := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(body), &http.Response{StatusCode: 200}, nil
	})
	c, err := NewStorageSKUCollector("https://management.azure.com", getter)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Collect(context.Background(), map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus"})
	if err != nil || got == nil {
		t.Fatalf("unexpected request error: %v", err)
	}
	return got
}

func storageAssertFault(t *testing.T, body, code string) {
	t.Helper()
	got := storageCollectBody(t, body)
	if got.Status != "unknown" || got.FailureCode != code || len(got.Groups) != 0 || len(got.UnqualifiedNames) != 0 {
		t.Fatalf("fault=%s result=%#v; expected atomic empty unknown", code, got)
	}
	e, err := got.Evidence(context.Background(), map[string][]string{"s": {"Storage"}})
	if err != nil || e.Status != "unknown" || len(e.Values) != 0 {
		t.Fatalf("fault admitted evidence: %#v / %v", e, err)
	}
}

func TestStorageSKUBodyAndRowsBounds(t *testing.T) {
	empty := `{"value":[]}`
	if got := storageCollectBody(t, empty+strings.Repeat(" ", 1048576-len(empty))); got.Status != "complete" {
		t.Fatal("exact raw bound rejected", got)
	}
	storageAssertFault(t, empty+strings.Repeat(" ", 1048577-len(empty)), "storage_sku_byte_limit")
	for _, n := range []int{8192, 8193} {
		body := `{"value":[` + strings.TrimSuffix(strings.Repeat(`{"name":"s"},`, n), ",") + `]}`
		if n == 8193 {
			storageAssertFault(t, body, "storage_sku_row_limit")
		} else if got := storageCollectBody(t, body); got.Status != "complete" || !got.UnqualifiedNames["s"] {
			t.Fatal("exact row bound rejected", got)
		}
	}
	rows := make([]string, 8192)
	for i := range rows {
		rows[i] = fmt.Sprintf(`{"resourceType":"storageAccounts","name":"s%d","kind":"Storage","locations":[],"restrictions":[]}`, i)
	}
	if got := storageCollectBody(t, `{"value":[`+strings.Join(rows, ",")+`]}`); got.Status != "complete" || len(got.Groups) != 8192 {
		t.Fatal("exact distinct group bound rejected", got.Status, got.FailureCode, len(got.Groups))
	}
}

func TestStorageSKUJSONBounds(t *testing.T) {
	for _, n := range []int{31, 32} {
		body := `{"value":[],"metadata":` + strings.Repeat("[", n) + "0" + strings.Repeat("]", n) + "}"
		if n == 32 {
			storageAssertFault(t, body, "storage_sku_invalid_response")
		} else if got := storageCollectBody(t, body); got.Status != "complete" {
			t.Fatal("depth32 rejected", got)
		}
	}
	for _, n := range []int{128, 129} {
		fields := make([]string, n)
		for i := range fields {
			fields[i] = fmt.Sprintf(`"k%d":null`, i)
		}
		body := `{"value":[],"metadata":{` + strings.Join(fields, ",") + `}}`
		if n == 129 {
			storageAssertFault(t, body, "storage_sku_invalid_response")
		} else if got := storageCollectBody(t, body); got.Status != "complete" {
			t.Fatal("128properties rejected", got)
		}
	}
	for _, n := range []int{65533, 65534} {
		body := `{"value":[],"metadata":[` + strings.TrimSuffix(strings.Repeat("null,", n), ",") + `]}`
		if n == 65534 {
			storageAssertFault(t, body, "storage_sku_invalid_response")
		} else if got := storageCollectBody(t, body); got.Status != "complete" {
			t.Fatal("65536tokens rejected", got)
		}
	}
	for _, nulls := range []int{1, 2} {
		body := `{"value":[],"metadata":[` + strings.Repeat("{},", 32767) + strings.TrimSuffix(strings.Repeat("null,", nulls), ",") + `]}`
		if nulls == 2 {
			storageAssertFault(t, body, "storage_sku_work_limit")
		} else if got := storageCollectBody(t, body); got.Status != "complete" {
			t.Fatal("65536object/element work rejected", got)
		}
	}
	for _, n := range []int{512, 513} {
		body := `{"value":[],"metadata":"` + strings.Repeat("x", n) + `"}`
		if n == 513 {
			storageAssertFault(t, body, "storage_sku_invalid_response")
		} else if got := storageCollectBody(t, body); got.Status != "complete" {
			t.Fatal("512byte label rejected", got)
		}
	}
}

func TestStorageSKUStrictAtomicAdmission(t *testing.T) {
	for _, body := range []string{
		`{"value":[],"metadata":"` + string([]byte{0xff}) + `"}`,
		`{"value":[],"metadata":"\u0000"}`,
		`{"value":[],"metadata":"\ud800"}`,
		`{"value":[],"metadata":"\uffff"}`,
		`{"value":[{"name":"s","resourceType":"queues","locations":null}]}`,
		`{"value":[{"name":"s","kind":null}]}`,
		`{"value":[{"name":"s","tier":null}]}`,
		`{"value":[{"name":"s","capabilities":null}]}`,
		`{"value":[{"name":"s","restrictions":null}]}`,
		`{"value":[{"name":"s","capabilities":[{"name":null}]}]}`,
		`{"value":[{"name":"s","restrictions":[{"values":[null]}]}]}`,
	} {
		storageAssertFault(t, body, "storage_sku_invalid_response")
	}
	storageAssertFault(t, strings.TrimSuffix(storageReadyBody, "]}")+`,{"name":5}]}`, "storage_sku_invalid_response")
}

func TestStorageSKURequestBounds(t *testing.T) {
	calls := 0
	getter := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(`{"value":[]}`), &http.Response{StatusCode: 200}, nil
	})
	for _, endpoint := range []string{"http://management.azure.com", "https://user@management.azure.com", "https://management.azure.com/extra", "https://management.azure.com?x=1"} {
		if c, err := NewStorageSKUCollector(endpoint, getter); c != nil || err == nil {
			t.Fatal("unsafe constructor accepted", endpoint)
		}
	}
	if c, err := NewStorageSKUCollector("https://management.azure.com", nil); c != nil || err == nil {
		t.Fatal("nil getter accepted")
	}
	c, _ := NewStorageSKUCollector("https://management.azure.com", getter)
	for _, n := range []int{256, 257} {
		aliases := map[string]string{}
		for i := 0; i < n; i++ {
			aliases[fmt.Sprintf("display label %d", i)] = "eastus"
		}
		before := calls
		got, err := c.Collect(context.Background(), map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus", LocationAliases: aliases})
		if n == 256 && (got == nil || err != nil || calls != before+1) || n == 257 && (got != nil || err == nil || calls != before) {
			t.Fatal("alias boundary failed", n, got, err, calls-before)
		}
	}
	for _, n := range []int{1000, 1001} {
		scope := map[string]string{storageTestID: "selected"}
		for i := 1; i < n; i++ {
			scope[fmt.Sprintf("%08x-0000-0000-0000-000000000000", i)] = "selected"
		}
		before := calls
		got, err := c.Collect(context.Background(), scope, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus"})
		if n == 1000 && (got == nil || err != nil || calls != before+1) || n == 1001 && (got != nil || err == nil || calls != before) {
			t.Fatal("scope boundary failed", n, got, err, calls-before)
		}
	}
	for _, aliases := range []map[string]string{{"East US": "eastus"}, {"eastus2(stage)": "eastus2"}, {"east us": "EastUS"}, {"eastus": "westus"}} {
		before := calls
		if got, err := c.Collect(context.Background(), map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus", LocationAliases: aliases}); got != nil || err == nil || calls != before {
			t.Fatal("invalid alias acquired credentials", got, err)
		}
	}
}

func storageEmptyCollection() *StorageSKUCollection {
	return &StorageSKUCollection{SubscriptionID: storageTestID, Region: "eastus", Status: "complete", Groups: map[string]map[string]StorageSKUGroup{}, UnqualifiedNames: map[string]bool{}}
}

func TestStorageSKUProjectionBounds(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		required := map[string][]string{}
		for i := 0; i < n; i++ {
			required[fmt.Sprintf("s%d", i)] = []string{"Storage"}
		}
		got, err := storageEmptyCollection().Evidence(context.Background(), required)
		if n == 8192 && (err != nil || got.Status != "complete" || len(got.Values) != 8192) || n == 8193 && (err == nil || len(got.Values) != 0) {
			t.Fatal("requested SKU boundary failed", n, got.Status, err)
		}
	}
	for _, n := range []int{16, 17} {
		kinds := []string{}
		for i := 0; i < n; i++ {
			kinds = append(kinds, fmt.Sprintf("kind%d", i))
		}
		got, err := storageEmptyCollection().Evidence(context.Background(), map[string][]string{"s": kinds})
		if n == 16 && (err != nil || got.Status != "complete") || n == 17 && (err == nil || len(got.Values) != 0) {
			t.Fatal("required kinds boundary failed", n, err)
		}
	}
	for _, excess := range []bool{false, true} {
		required := map[string][]string{}
		for i := 0; i < 8192; i++ {
			required[fmt.Sprintf("s%d", i)] = []string{"k0", "k1", "k2", "k3", "k4", "k5", "k6"}
		}
		if excess {
			required["s0"] = append(required["s0"], "k7")
		}
		got, err := storageEmptyCollection().Evidence(context.Background(), required)
		if !excess && (err != nil || got.Status != "complete" || len(got.Values) != 8192) || excess && (err == nil || !strings.Contains(err.Error(), "[storage_sku_projection_limit]")) {
			t.Fatal("65536entry boundary failed", excess, got.Status, err)
		}
	}
	// Collection identity/status contributes50bytes. 2040*(1SKU+15kind)
	// labels of512bytes minus50 is exactly the independent16711680byte cap.
	for _, excess := range []bool{false, true} {
		required := map[string][]string{}
		for i := 0; i < 2040; i++ {
			name := fmt.Sprintf("s%04d", i) + strings.Repeat("x", 507)
			kinds := []string{}
			for k := 0; k < 15; k++ {
				length := 512
				if i == 0 && k == 0 {
					length = 462
					if excess {
						length++
					}
				}
				kinds = append(kinds, fmt.Sprintf("k%02d", k)+strings.Repeat("x", length-3))
			}
			required[name] = kinds
		}
		got, err := storageEmptyCollection().Evidence(context.Background(), required)
		if !excess && (err != nil || len(got.Values) != 2040) || excess && (err == nil || !strings.Contains(err.Error(), "[storage_sku_projection_limit]")) {
			t.Fatal("16711680byte text boundary failed", excess, got.Status, err)
		}
	}
}

func TestStorageSKUProjectionRejectsCorruption(t *testing.T) {
	for _, change := range []func(*StorageSKUCollection){
		func(c *StorageSKUCollection) { c.Status = "fabricated" },
		func(c *StorageSKUCollection) { c.SubscriptionID = "bad" },
		func(c *StorageSKUCollection) { c.Region = "global" },
		func(c *StorageSKUCollection) { c.Groups["S"] = map[string]StorageSKUGroup{} },
		func(c *StorageSKUCollection) {
			c.Groups["s"] = map[string]StorageSKUGroup{"storage": {Known: false, State: 0}}
		},
		func(c *StorageSKUCollection) {
			c.Groups["s"] = map[string]StorageSKUGroup{"storage": {Known: true, State: 3}}
		},
		func(c *StorageSKUCollection) { c.UnqualifiedNames["s"] = false },
	} {
		c := storageEmptyCollection()
		change(c)
		if got, err := c.Evidence(context.Background(), map[string][]string{"s": {"Storage"}}); err == nil || len(got.Values) != 0 {
			t.Fatal("corrupt collection admitted", c, got, err)
		}
	}
	for _, required := range []map[string][]string{{"s": {"Storage"}, "S": {"Storage"}}, {" ": {"Storage"}}, {"s": {"bad kind"}}} {
		if got, err := storageEmptyCollection().Evidence(context.Background(), required); err == nil || len(got.Values) != 0 {
			t.Fatal("malformed requirements admitted", got, err)
		}
	}
}

func TestStorageSKUOwnershipAndConcurrency(t *testing.T) {
	body := []byte(storageReadyBody)
	before := append([]byte{}, body...)
	getter := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return body, &http.Response{StatusCode: 200}, nil
	})
	c, _ := NewStorageSKUCollector("https://management.azure.com", getter)
	var wait sync.WaitGroup
	for i := 0; i < 8; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			got, err := c.Collect(context.Background(), map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus"})
			if err != nil || got == nil || got.Groups["s"]["storage"].State != 0 {
				t.Error("concurrent call state lost", got, err)
				return
			}
			required := map[string][]string{"s": {"Storage"}}
			e, err := got.Evidence(context.Background(), required)
			if err != nil || e.Values["s"].State != 0 || required["s"][0] != "Storage" {
				t.Error("projection mutated input", e, err)
			}
			e.Values["s"] = SKUAvailability{State: 2, BlockedZones: []string{"canary"}}
			if got.Groups["s"]["storage"].State != 0 {
				t.Error("output aliases collection")
			}
			got.Groups["s"]["storage"] = StorageSKUGroup{Known: true, State: 2}
		}()
	}
	wait.Wait()
	if !reflect.DeepEqual(body, before) {
		t.Fatal("decoder mutated caller body")
	}
}

type storageCanaryTransport struct {
	t             *testing.T
	status        int
	calls, closes atomic.Int32
}

func (f *storageCanaryTransport) Do(request *http.Request) (*http.Response, error) {
	f.calls.Add(1)
	if request.Method != "GET" || request.URL.String() != "https://management.azure.com/subscriptions/"+storageTestID+"/providers/Microsoft.Storage/skus?api-version=2023-01-01" || request.Body != nil && request.Body != http.NoBody || request.Header.Get("Authorization") != "Bearer OFFLINE_CANARY" {
		f.t.Error("unexpected authenticated Storage request", request.Method, request.URL)
	}
	return &http.Response{StatusCode: f.status, Header: make(http.Header), Body: &diskCloseBody{Reader: io.Reader(strings.NewReader(storageReadyBody)), closes: &f.closes}, Request: request}, nil
}

func TestStorageSKUActualHTTPClient(t *testing.T) {
	for _, status := range []int{200, 403, 429, 500, 307} {
		credential := &diskCanaryCredential{}
		wire := &storageCanaryTransport{t: t, status: status}
		client := azure.NewHTTPClient(credential, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, OperationTimeout: 5 * time.Second, Scope: "https://management.azure.com/.default", Transport: wire})
		c, _ := NewStorageSKUCollector("https://management.azure.com", client)
		got, err := c.Collect(context.Background(), map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus"})
		if err != nil || got == nil || wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 {
			t.Fatal("actual bounded client lifecycle mismatch", status, got, err)
		}
		if status != 200 && (got.Status != "unknown" || got.FailureCode != "storage_sku_status_invalid") {
			t.Fatal("actual status/error classification lost", status, got)
		}
		if status == 200 && (got.Status != "complete" || got.Groups["s"]["storage"].State != 0) || !reflect.DeepEqual(credential.scopes, []string{"https://management.azure.com/.default"}) {
			t.Fatal("actual successful Storage response/audience lost", got)
		}
	}
}

type storageCountingContext struct {
	context.Context
	calls, limit int
}

func (c *storageCountingContext) Err() error {
	c.calls++
	if c.calls >= c.limit {
		return context.Canceled
	}
	return nil
}

func TestStorageSKUCancellationDuringWork(t *testing.T) {
	body := `{"value":[` + strings.TrimSuffix(strings.Repeat(`{"name":"s"},`, 8192), ",") + `]}`
	getter := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(body), &http.Response{StatusCode: 200}, nil
	})
	c, _ := NewStorageSKUCollector("https://management.azure.com", getter)
	ctx := &storageCountingContext{Context: context.Background(), limit: 1000}
	if got, err := c.Collect(ctx, map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus"}); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatal("decode cancellation admitted a collection", got, err)
	}
	required := map[string][]string{}
	for i := 0; i < 8192; i++ {
		required[fmt.Sprintf("s%d", i)] = []string{"Storage"}
	}
	ctx = &storageCountingContext{Context: context.Background(), limit: 1000}
	if got, err := storageEmptyCollection().Evidence(ctx, required); len(got.Values) != 0 || !errors.Is(err, context.Canceled) {
		t.Fatal("projection cancellation admitted evidence", got, err)
	}
}

func TestStorageSKUCancellation(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		ctx := &storageTerminalContext{Context: context.Background(), done: make(chan struct{})}
		getter := storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			ctx.finish(cause)
			return []byte(storageReadyBody), &http.Response{StatusCode: 403}, errors.New("PRIVATE transport canary")
		})
		c, _ := NewStorageSKUCollector("https://management.azure.com", getter)
		if got, err := c.Collect(ctx, map[string]string{storageTestID: "selected"}, StorageSKURequest{SubscriptionID: storageTestID, Region: "eastus"}); got != nil || !errors.Is(err, cause) {
			t.Fatal("terminal context lost priority", got, err)
		}
		if got, err := storageEmptyCollection().Evidence(ctx, map[string][]string{"s": {"Storage"}}); len(got.Values) != 0 || !errors.Is(err, cause) {
			t.Fatal("projection ignored cancellation", got, err)
		}
	}
}

func TestStorageSKUInputOutputUnchanged(t *testing.T) {
	required := map[string][]string{"s": {"Storage", "storage"}}
	original, _ := json.Marshal(required)
	got := storageCollectBody(t, storageReadyBody)
	e, err := got.Evidence(context.Background(), required)
	after, _ := json.Marshal(required)
	if err != nil || string(after) != string(original) || e.Status != "complete" || len(e.Values) != 1 || e.Values["s"].State != 0 {
		t.Fatal("deduplicated kind projection mutated requirements", e, err)
	}
}
