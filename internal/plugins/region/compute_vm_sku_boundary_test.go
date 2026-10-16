package region

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const vmReadyBody = `{"value":[{"resourceType":"virtualMachines","name":"s","locations":["eastus"],"restrictions":[]}]}`

func vmRequest() ComputeVMSKURequest {
	return ComputeVMSKURequest{SubscriptionID: vmTestID, Region: "eastus", ResourceType: vmTestType}
}
func vmEmptyCollection() *ComputeVMSKUCollection {
	return &ComputeVMSKUCollection{SubscriptionID: vmTestID, Region: "eastus", ResourceType: vmTestType, Status: "complete", Groups: map[string]ComputeVMSKUGroup{}}
}
func vmCollect(t *testing.T, body string) *ComputeVMSKUCollection {
	t.Helper()
	c, e := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
	if e != nil || r == nil {
		t.Fatal(r, e)
	}
	return r
}
func vmFault(t *testing.T, body, code string) {
	t.Helper()
	r := vmCollect(t, body)
	if r.Status != "unknown" || r.FailureCode != code || len(r.Groups) != 0 {
		t.Fatalf("expected atomic %s got %#v", code, r)
	}
}
func vmNext(n int) string {
	return "/subscriptions/" + vmTestID + "/providers/Microsoft.Compute/skus?%24filter=location+eq+%27eastus%27&%24skiptoken=" + fmt.Sprint(n) + "&api-version=2021-07-01"
}

func TestComputeVMSKUPageAndAggregateBounds(t *testing.T) {
	empty := `{"value":[]}`
	if got := vmCollect(t, empty+strings.Repeat(" ", 1048576-len(empty))); got.Status != "complete" {
		t.Fatal("exact page bytes", got)
	}
	vmFault(t, empty+strings.Repeat(" ", 1048577-len(empty)), "vm_sku_byte_limit")
	for _, extra := range []bool{false, true} {
		calls := 0
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			last := 8
			if extra {
				last = 9
			}
			body := empty
			if calls < last {
				body = `{"value":[],"nextLink":"` + vmNext(calls) + `"}`
			}
			size := 1048576
			if extra && calls == 1 {
				size -= len(empty) - 1
			}
			if extra && calls == 9 {
				size = len(empty)
			}
			return []byte(body + strings.Repeat(" ", size-len(body))), &http.Response{StatusCode: 200}, nil
		}))
		got, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		if e != nil || got == nil {
			t.Fatal(e)
		}
		if !extra && (got.Status != "complete" || calls != 8) || extra && (got.Status != "partial" || got.FailureCode != "vm_sku_aggregate_limit" || calls != 9 || len(got.Groups) != 0) {
			t.Fatal("8MiB/excess1", extra, got, calls)
		}
	}
	for _, more := range []bool{false, true} {
		calls := 0
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			body := empty
			if calls < 64 || more {
				body = `{"value":[],"nextLink":"` + vmNext(calls) + `"}`
			}
			return []byte(body), &http.Response{StatusCode: 200}, nil
		}))
		got, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		if e != nil || calls != 64 || !more && got.Status != "complete" || more && (got.Status != "partial" || got.FailureCode != "vm_sku_page_limit") {
			t.Fatal("64 page boundary", more, got, e, calls)
		}
	}
}

func TestComputeVMSKUGroupsAndScopeBounds(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		calls := 0
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			start := calls * 512
			calls++
			end := start + 512
			if end > n {
				end = n
			}
			rows := []string{}
			for i := start; i < end; i++ {
				rows = append(rows, fmt.Sprintf(`{"resourceType":"virtualMachines","name":"s%d","locations":[],"restrictions":[]}`, i))
			}
			next := ""
			if end < n {
				next = `,"nextLink":"` + vmNext(calls) + `"`
			}
			return []byte(`{"value":[` + strings.Join(rows, ",") + `]` + next + `}`), &http.Response{StatusCode: 200}, nil
		}))
		got, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		if e != nil || n == 8192 && (got.Status != "complete" || len(got.Groups) != 8192) || n == 8193 && (got.Status != "partial" || got.FailureCode != "vm_sku_work_limit" || len(got.Groups) != 8192) {
			t.Fatal("group boundary", n, got.Status, got.FailureCode, len(got.Groups), e)
		}
	}
	calls := 0
	c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(`{"value":[]}`), &http.Response{StatusCode: 200}, nil
	}))
	for _, n := range []int{1000, 1001} {
		scope := map[string]string{}
		for i := 0; i < n; i++ {
			scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "selected"
		}
		scope["00000000-1111-1111-1111-111111111111"] = "selected"
		q := vmRequest()
		q.SubscriptionID = "00000000-1111-1111-1111-111111111111"
		before := calls
		got, e := c.Collect(context.Background(), scope, q)
		if n == 1000 && (e != nil || got == nil || calls != before+1) || n == 1001 && (e == nil || got != nil || calls != before) {
			t.Fatal("scope limit", n, got, e, calls)
		}
	}
	for _, endpoint := range []string{"http://management.azure.com", "https://u@management.azure.com", "https://management.azure.com/path", "https://management.azure.com?x=1"} {
		if got, e := NewComputeVMSKUCollector(endpoint, c.transport.getter); got != nil || e == nil {
			t.Fatal("unsafe origin accepted", endpoint)
		}
	}
}

func TestComputeVMSKUJSONAndRawWorkBounds(t *testing.T) {
	for _, n := range []int{31, 32} {
		body := `{"value":[],"meta":` + strings.Repeat("[", n) + "0" + strings.Repeat("]", n) + "}"
		if n == 32 {
			vmFault(t, body, "vm_sku_invalid_response")
		} else if vmCollect(t, body).Status != "complete" {
			t.Fatal("depth32")
		}
	}
	for _, n := range []int{128, 129} {
		keys := []string{}
		for i := 0; i < n; i++ {
			keys = append(keys, fmt.Sprintf(`"k%d":null`, i))
		}
		body := `{"value":[],"meta":{` + strings.Join(keys, ",") + `}}`
		if n == 129 {
			vmFault(t, body, "vm_sku_invalid_response")
		} else if vmCollect(t, body).Status != "complete" {
			t.Fatal("128 keys")
		}
	}
	for _, n := range []int{65533, 65534} {
		body := `{"value":[],"meta":[` + strings.TrimSuffix(strings.Repeat("null,", n), ",") + `]}`
		if n == 65534 {
			vmFault(t, body, "vm_sku_invalid_response")
		} else if vmCollect(t, body).Status != "complete" {
			t.Fatal("65536 tokens")
		}
	}
	for _, n := range []int{1, 2} {
		body := `{"value":[],"meta":[` + strings.Repeat("{},", 32767) + strings.TrimSuffix(strings.Repeat("null,", n), ",") + `]}`
		if n == 2 {
			vmFault(t, body, "vm_sku_work_limit")
		} else if vmCollect(t, body).Status != "complete" {
			t.Fatal("65536 work")
		}
	}
	for _, n := range []int{512, 513} {
		body := `{"value":[],"meta":"` + strings.Repeat("x", n) + `"}`
		if n == 513 {
			vmFault(t, body, "vm_sku_invalid_response")
		} else if vmCollect(t, body).Status != "complete" {
			t.Fatal("512 scalar")
		}
	}
	for _, body := range []string{`{"value":[],"meta":"` + string([]byte{0xff}) + `"}`, `{"value":[],"meta":"\ud800"}`, `{"value":[],"meta":"\uffff"}`, `{"value":[],"meta":"\u0000"}`, strings.TrimSuffix(vmReadyBody, "]}") + `,null]}`} {
		vmFault(t, body, "vm_sku_invalid_response")
	}
}

func TestComputeVMSKUContinuationLimits(t *testing.T) {
	for _, n := range []int{4096, 4097} {
		calls := 0
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			body := `{"value":[]}`
			if calls == 1 {
				body = `{"value":[],"nextLink":"` + vmNext(1) + strings.Repeat("x", n-1) + `"}`
			}
			return []byte(body), &http.Response{StatusCode: 200}, nil
		}))
		got, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		if e != nil || n == 4096 && (got.Status != "complete" || calls != 2) || n == 4097 && (got.FailureCode != "vm_sku_continuation_invalid" || calls != 1) {
			t.Fatal("decoded token", n, got, e, calls)
		}
	}
	// Literal escaping exercises8192wire bytes with a<=4096decoded token.
	base := vmNext(1)
	for _, n := range []int{8192, 8193} {
		pad := n - len(base)
		escaped := strings.Repeat("%61", pad/3) + strings.Repeat("a", pad%3)
		link := base + escaped
		calls := 0
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			body := `{"value":[]}`
			if calls == 1 {
				body = `{"value":[],"nextLink":"` + link + `"}`
			}
			return []byte(body), &http.Response{StatusCode: 200}, nil
		}))
		got, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		if e != nil || n == 8192 && (got.Status != "complete" || calls != 2) || n == 8193 && (got.Status != "unknown" || got.FailureCode != "vm_sku_invalid_response" || calls != 1) {
			t.Fatal("wire link", n, got, e, calls)
		}
	}
	calls := 0
	c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(`{"value":[],"nextLink":"` + vmTestURL + `"}`), &http.Response{StatusCode: 200}, nil
	}))
	got, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
	if e != nil || calls != 1 || got.Status != "partial" || got.FailureCode != "vm_sku_cycle" {
		t.Fatal("cycle requested", got, e, calls)
	}
}

func TestComputeVMSKUProjectionBoundsAndCorruption(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		required := []string{}
		for i := 0; i < n; i++ {
			required = append(required, fmt.Sprintf("s%d", i))
		}
		got, e := vmEmptyCollection().Evidence(context.Background(), required)
		if n == 8192 && (e != nil || len(got.Values) != 8192) || n == 8193 && (e == nil || len(got.Values) != 0) {
			t.Fatal("required limit", n, got.Status, e)
		}
	}
	for _, extra := range []bool{false, true} {
		c := vmEmptyCollection()
		c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: []string{}}
		g := c.Groups["s"]
		for i := 1; i <= 65535; i++ {
			g.BlockedZones = append(g.BlockedZones, fmt.Sprint(i))
		}
		sort.Strings(g.BlockedZones)
		if extra {
			g.BlockedZones = append(g.BlockedZones, "99999")
		}
		c.Groups["s"] = g
		got, e := c.Evidence(context.Background(), nil)
		if !extra && (e != nil || got.Status != "complete") || extra && (e == nil || len(got.Values) != 0) {
			t.Fatal("65536 projection logical entries", extra, e)
		}
	}
	// Collection binding labels use exactly104bytes: UUID36+eastus6+resourceType34+complete8, plus actual resource type length checked below.
	binding := len(vmTestID) + len("eastus") + len(vmTestType) + len("complete")
	for _, extra := range []int{0, 1} {
		c := vmEmptyCollection()
		g := ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: []string{}}
		remaining := 16711680 - binding - 1
		for i := 0; remaining > 0; i++ {
			size := 512
			if remaining < size {
				size = remaining
			}
			zone := fmt.Sprintf("1%07d", i) + strings.Repeat("1", size-8)
			if remaining <= 512 && extra == 1 {
				zone += "1"
			}
			g.BlockedZones = append(g.BlockedZones, zone)
			remaining -= size
		}
		sort.Strings(g.BlockedZones)
		c.Groups["s"] = g
		got, e := c.Evidence(context.Background(), nil)
		if extra == 0 && (e != nil || got.Status != "complete") || extra == 1 && (e == nil || len(got.Values) != 0) {
			t.Fatal("projection text cap", extra, e)
		}
	}
	for _, change := range []func(*ComputeVMSKUCollection){func(c *ComputeVMSKUCollection) { c.Status = "unknown"; c.Groups["s"] = ComputeVMSKUGroup{Known: true} }, func(c *ComputeVMSKUCollection) { c.SubscriptionID = "bad" }, func(c *ComputeVMSKUCollection) { c.FailureCode = "opaque" }, func(c *ComputeVMSKUCollection) { c.Groups["S"] = ComputeVMSKUGroup{Known: true} }, func(c *ComputeVMSKUCollection) { c.Groups["s"] = ComputeVMSKUGroup{State: 0} }, func(c *ComputeVMSKUCollection) {
		c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: []string{"2", "1"}}
	}} {
		c := vmEmptyCollection()
		change(c)
		if got, e := c.Evidence(context.Background(), []string{"s"}); e == nil || len(got.Values) != 0 {
			t.Fatal("forged collection", c, got, e)
		}
	}
	for _, required := range [][]string{{"s", "S"}, {" "}} {
		if got, e := vmEmptyCollection().Evidence(context.Background(), required); e == nil || len(got.Values) != 0 {
			t.Fatal("bad requirements", got, e)
		}
	}
}

type vmWire struct {
	t             *testing.T
	status        int
	calls, closes atomic.Int32
}

func (w *vmWire) Do(r *http.Request) (*http.Response, error) {
	w.calls.Add(1)
	if r.Method != "GET" || r.URL.String() != vmTestURL || r.Body != nil && r.Body != http.NoBody || r.Header.Get("Authorization") != "Bearer OFFLINE_CANARY" {
		w.t.Error("actual request", r.Method, r.URL)
	}
	return &http.Response{StatusCode: w.status, Header: make(http.Header), Body: &diskCloseBody{Reader: io.Reader(strings.NewReader(vmReadyBody)), closes: &w.closes}, Request: r}, nil
}
func TestComputeVMSKUActualHTTPAndCancellation(t *testing.T) {
	for _, status := range []int{200, 403, 429, 500, 307} {
		cred := &diskCanaryCredential{}
		wire := &vmWire{t: t, status: status}
		client := azure.NewHTTPClient(cred, &azure.HTTPClientOptions{Timeout: time.Second, OperationTimeout: 5 * time.Second, MaxRetries: -1, Scope: "https://management.azure.com/.default", Transport: wire})
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", client)
		r, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		if e != nil || r == nil || wire.calls.Load() != 1 || wire.closes.Load() != 1 || cred.calls.Load() != 1 || !reflect.DeepEqual(cred.scopes, []string{"https://management.azure.com/.default"}) || status != 200 && r.FailureCode != "vm_sku_status_invalid" || status == 200 && r.Groups["s"].State != 0 {
			t.Fatal("actual HTTP", status, r, e)
		}
	}
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		ctx := &storageTerminalContext{Context: context.Background(), done: make(chan struct{})}
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			ctx.finish(cause)
			return []byte(vmReadyBody), &http.Response{StatusCode: 403}, errors.New("opaque canary")
		}))
		if r, e := c.Collect(ctx, map[string]string{vmTestID: "selected"}, vmRequest()); r != nil || !errors.Is(e, cause) {
			t.Fatal("context priority", r, e)
		}
	}
	for _, status := range []int{200, 403} {
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return []byte(vmReadyBody), &http.Response{StatusCode: status}, errors.New("opaque canary")
		}))
		r, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
		code := "vm_sku_request_failed"
		if status == 403 {
			code = "vm_sku_status_invalid"
		}
		if e != nil || r.FailureCode != code {
			t.Fatal("status/error pair", r, e)
		}
	}
}

func TestComputeVMSKUOwnershipAndLateContext(t *testing.T) {
	body := []byte(`{"value":[{"resourceType":"virtualMachines","name":"s","locations":["eastus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["2","1"]}}]}]}`)
	before := append([]byte{}, body...)
	c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return body, &http.Response{StatusCode: 200}, nil
	}))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
			if e != nil || r == nil {
				t.Error(r, e)
				return
			}
			required := []string{"s"}
			a, e := r.Evidence(context.Background(), required)
			if e != nil || !reflect.DeepEqual(a.Values["s"].BlockedZones, []string{"1", "2"}) {
				t.Error(a, e)
				return
			}
			a.Values["s"].BlockedZones[0] = "99"
			b, e := r.Evidence(context.Background(), required)
			if e != nil || b.Values["s"].BlockedZones[0] != "1" || required[0] != "s" {
				t.Error("shared output", b, e)
			}
		}()
	}
	wg.Wait()
	if !reflect.DeepEqual(body, before) {
		t.Fatal("borrowed response mutated")
	}
	for _, collection := range []*ComputeVMSKUCollection{nil, vmEmptyCollection()} {
		if collection != nil {
			collection.Status = "unknown"
		}
		limit := 2
		if collection != nil {
			limit = 6
		}
		ctx := &storageCountingContext{Context: context.Background(), limit: limit}
		got, e := collection.Evidence(ctx, nil)
		if len(got.Values) != 0 || !errors.Is(e, context.Canceled) {
			t.Fatal("late unknown cancellation", got, e)
		}
	}
	ss := vmEmptyCollection()
	ss.ResourceType = "microsoft.compute/virtualmachinescalesets"
	ctx := &storageCountingContext{Context: context.Background(), limit: 6}
	if got, e := ss.Evidence(ctx, nil); len(got.Values) != 0 || !errors.Is(e, context.Canceled) {
		t.Fatal("late unsupported cancellation", got, e)
	}
}
