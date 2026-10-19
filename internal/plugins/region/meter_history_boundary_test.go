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

const historyTestID = "11111111-1111-1111-1111-111111111111"
const historyEmpty = `{"properties":{"columns":[],"rows":[]}}`
const historyFirst = "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CostManagement/query?api-version=2021-10-01"

type historyPosterFunc func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)

func (f historyPosterFunc) PostBounded(c context.Context, u string, b io.ReadSeekCloser, n int64) ([]byte, *http.Response, error) {
	return f(c, u, b, n)
}
func historyTestRequest() MeterHistoryRequest {
	return MeterHistoryRequest{SubscriptionID: historyTestID, Months: 1, Now: time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC), Currency: "USD"}
}
func historyRunBodies(t *testing.T, bodies ...string) (*MeterHistoryCollection, int) {
	t.Helper()
	calls := 0
	c, e := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(_ context.Context, u string, _ io.ReadSeekCloser, n int64) ([]byte, *http.Response, error) {
		if calls >= len(bodies) || n != 1048576 {
			t.Fatal("unexpected boundary HTTP call", calls, n, u)
		}
		b := bodies[calls]
		calls++
		return []byte(b), &http.Response{StatusCode: 200}, nil
	}))
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.Collect(context.Background(), map[string]string{historyTestID: "Alpha"}, historyTestRequest())
	if e != nil || r == nil {
		t.Fatal("boundary collection", r, e)
	}
	return r, calls
}
func historyRequire(t *testing.T, r *MeterHistoryCollection, status, code string) {
	t.Helper()
	if r.Status != status || r.FailureCode != code || r.SubscriptionID != historyTestID || r.Currency != "USD" {
		t.Fatalf("explicit history health: got=%+v want=%s/%s", r, status, code)
	}
	if status != "complete" && (r.Evidence().Complete || len(r.Evidence().Meters) != 0) {
		t.Fatal("incomplete diagnostic prefix became cost evidence")
	}
}
func historyPad(n int, next string) string {
	head := `{"properties":{"columns":[],"rows":[]`
	if next != "" {
		b, _ := json.Marshal(next)
		head += `,"nextLink":` + string(b)
	}
	head += `},"padding":[`
	tail := `]}`
	remaining := n - len(head) - len(tail)
	k := (remaining + 1 + 8194) / 8195
	payload := remaining - (3*k - 1)
	parts := make([]string, k)
	for i := range parts {
		v := min(payload, 8192)
		parts[i] = `"` + strings.Repeat("x", v) + `"`
		payload -= v
	}
	return head + strings.Join(parts, ",") + tail
}
func historyPageWithNext(next string) string {
	return strings.TrimSuffix(historyEmpty, "}}") + `,"nextLink":"` + next + `"}}`
}
func TestMeterHistoryWireBound(t *testing.T) {
	for _, n := range []int{1048576, 1048577} {
		b := historyPad(n, "")
		if len(b) != n {
			t.Fatal("independent wire fixture", len(b), n)
		}
		r, _ := historyRunBodies(t, b)
		if n == 1048576 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
}
func TestMeterHistoryTotalWireBound(t *testing.T) {
	for _, extra := range []int{0, 1} {
		var bodies []string
		for i := 0; i < 8; i++ {
			next := ""
			if i < 7 || extra == 1 {
				next = fmt.Sprintf("?api-version=2021-10-01&$skiptoken=%d", i+2)
			}
			n := 1048576
			if i == 7 && extra == 1 {
				n -= len(historyEmpty) - 1
			}
			bodies = append(bodies, historyPad(n, next))
		}
		if extra == 1 {
			bodies = append(bodies, historyEmpty)
		}
		total := 0
		for _, b := range bodies {
			total += len(b)
		}
		if total != 8388608+extra {
			t.Fatal("independent cumulative wire fixture", total)
		}
		r, calls := historyRunBodies(t, bodies...)
		if calls != len(bodies) {
			t.Fatal("cumulative wire fixture masked", calls)
		}
		if extra == 0 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "partial", "history_limit")
		}
	}
}
func TestMeterHistoryPageBound(t *testing.T) {
	for _, n := range []int{64, 65} {
		bodies := make([]string, n)
		for i := range bodies {
			bodies[i] = historyEmpty
			if i < n-1 {
				bodies[i] = historyPageWithNext(fmt.Sprintf("?api-version=2021-10-01&$skiptoken=%d", i+2))
			}
		}
		r, calls := historyRunBodies(t, bodies...)
		if calls != 64 {
			t.Fatal("page bound must precede request65", calls)
		}
		if n == 64 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "partial", "history_limit")
		}
	}
}
func TestMeterHistoryRawWorkBound(t *testing.T) {
	for _, extra := range []int{0, 1} {
		// Five containers plus65531nulls is exactly65536values, independently
		// reachable below the wire, string, key and depth bounds.
		b := `{"properties":{"columns":[],"rows":[]},"extra":[` + strings.Repeat("null,", 65530) + "null" + strings.Repeat(",null", extra) + `]}`
		if len(b) > 1048576 {
			t.Fatal("masked raw-work fixture")
		}
		r, _ := historyRunBodies(t, b)
		if extra == 0 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
}
func TestMeterHistoryDepthBound(t *testing.T) {
	for _, n := range []int{31, 32} {
		b := `{"properties":{"columns":[],"rows":[]},"extra":` + strings.Repeat("[", n) + "null" + strings.Repeat("]", n) + "}"
		r, _ := historyRunBodies(t, b)
		if n == 31 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
}
func historyRows(n int) string {
	var rows []string
	for i := 0; i < n; i++ {
		rows = append(rows, fmt.Sprintf(`["/subscriptions/%s/x","m%04d","USD",1]`, historyTestID, i))
	}
	return `{"properties":{"columns":[{"name":"ResourceId","type":"String"},{"name":"ResourceGuid","type":"String"},{"name":"Currency","type":"String"},{"name":"PreTaxCost","type":"Number"}],"rows":[` + strings.Join(rows, ",") + `]}}`
}
func TestMeterHistoryRowBound(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		b := historyRows(n)
		if len(b) > 1048576 {
			t.Fatal("row bound masked by wire")
		}
		r, _ := historyRunBodies(t, b)
		if n == 8192 {
			historyRequire(t, r, "complete", "")
			if len(r.Meters) != 8192 || r.Meters[0].MeterID != "m0000" || r.Meters[8191].MeterID != "m8191" {
				t.Fatal("exact distinct-meter output")
			}
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
}
func TestMeterHistoryMetadataBounds(t *testing.T) {
	for _, n := range []int{128, 129} {
		var fields []string
		for i := 0; i < n; i++ {
			fields = append(fields, fmt.Sprintf(`"k%d":null`, i))
		}
		r, _ := historyRunBodies(t, strings.TrimSuffix(historyEmpty, "}")+`,"extra":{`+strings.Join(fields, ",")+`}}`)
		if n == 128 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
	for _, n := range []int{8192, 8193} {
		r, _ := historyRunBodies(t, strings.TrimSuffix(historyEmpty, "}")+`,"extra":"`+strings.Repeat("x", n)+`"}`)
		if n == 8192 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
	for _, extra := range []string{`"K":1,"\u212a":2`, `"extra":1e999`, `"extra":"\u0000"`, `"extra":"` + string([]byte{255}) + `"`} {
		r, _ := historyRunBodies(t, strings.TrimSuffix(historyEmpty, "}")+`,`+extra+`}`)
		historyRequire(t, r, "unknown", "history_decode_invalid")
	}
}
func TestMeterHistoryPrevalidationAndTerminalContext(t *testing.T) {
	var calls atomic.Int32
	p := historyPosterFunc(func(ctx context.Context, _ string, _ io.ReadSeekCloser, _ int64) ([]byte, *http.Response, error) {
		calls.Add(1)
		return []byte(historyEmpty), &http.Response{StatusCode: 200}, nil
	})
	for _, endpoint := range []string{"http://management.azure.com", "https://user@management.azure.com", "https://management.azure.com/path", "https://management.azure.com?x=y", "https://management.azure.com#x", "https://management.azure.com:444", "https://management.azure.com/%2f"} {
		if c, e := NewMeterHistoryCollector(endpoint, p); e == nil || c != nil {
			t.Fatal("unsafe constructor admitted", endpoint)
		}
	}
	if c, e := NewMeterHistoryCollector("https://management.azure.com", nil); e == nil || c != nil {
		t.Fatal("nil poster admitted")
	}
	c, _ := NewMeterHistoryCollector("https://management.azure.com", p)
	for _, scope := range []map[string]string{nil, {historyTestID: ""}, {strings.ToUpper("abcdefab-1234-4234-8234-123456789abc"): "A", "abcdefab-1234-4234-8234-123456789abc": "B"}} {
		r, e := c.Collect(context.Background(), scope, historyTestRequest())
		if r != nil || e == nil || calls.Load() != 0 {
			t.Fatal("invalid scope reached transport", r, e, calls.Load())
		}
	}
	var nilCollector *MeterHistoryCollector
	if r, e := nilCollector.Collect(context.Background(), map[string]string{historyTestID: "Alpha"}, historyTestRequest()); r != nil || e == nil {
		t.Fatal("nil collector panic/admission")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if r, e := c.Collect(ctx, map[string]string{historyTestID: "Alpha"}, historyTestRequest()); r != nil || !errors.Is(e, context.Canceled) || calls.Load() != 0 {
		t.Fatal("cancelled caller reached transport", r, e)
	}
	ctx, cancel = context.WithCancel(context.Background())
	late, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		cancel()
		return []byte(historyEmpty), &http.Response{StatusCode: 200}, nil
	}))
	if r, e := late.Collect(ctx, map[string]string{historyTestID: "Alpha"}, historyTestRequest()); r != nil || !errors.Is(e, context.Canceled) {
		t.Fatal("post-HTTP cancellation lost", r, e)
	}
}

type historyWire func(*http.Request) (*http.Response, error)

func (f historyWire) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestMeterHistoryActualBoundedPOST(t *testing.T) {
	for _, status := range []int{200, 403, 307} {
		var calls, closes atomic.Int32
		cred := &diskCanaryCredential{}
		wire := historyWire(func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			if req.Method != "POST" || req.URL.String() != historyFirst || req.Header.Get("Authorization") != "Bearer OFFLINE_CANARY" || req.Header.Get("Content-Type") != "application/json" {
				t.Fatal("actual bounded history POST destination/auth/method")
			}
			b, e := io.ReadAll(req.Body)
			var got any
			if e != nil || json.Unmarshal(b, &got) != nil || !reflect.DeepEqual(got, historyReadOracle(t).RequestBody) {
				t.Fatal("actual complete history request body", e)
			}
			return &http.Response{StatusCode: status, Header: http.Header{}, Request: req, Body: &diskCloseBody{Reader: strings.NewReader(historyEmpty), closes: &closes}}, nil
		})
		client := azure.NewHTTPClient(cred, &azure.HTTPClientOptions{Timeout: 10 * time.Second, MaxRetries: -1, OperationTimeout: 30 * time.Second, Scope: "https://audience.example.test/.default", Transport: wire})
		c, _ := NewMeterHistoryCollector("https://management.azure.com", client)
		r, e := c.Collect(context.Background(), map[string]string{historyTestID: "Alpha"}, historyTestRequest())
		if e != nil || r == nil || calls.Load() != 1 || closes.Load() != 1 || cred.calls.Load() != 1 || !reflect.DeepEqual(cred.scopes, []string{"https://audience.example.test/.default"}) {
			t.Fatal("actual POST closure/audience/redirect", r, e, calls.Load(), closes.Load(), cred.scopes)
		}
		if status == 200 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_response_invalid")
		}
		q := historyTestRequest()
		q.Months = 13
		if r, e = c.Collect(context.Background(), map[string]string{historyTestID: "Alpha"}, q); r != nil || e == nil || calls.Load() != 1 || cred.calls.Load() != 1 {
			t.Fatal("invalid history authenticated", r, e)
		}
	}
}
func TestMeterHistoryConcurrentOwnedEvidence(t *testing.T) {
	body := historyRows(2)
	c, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			r, e := c.Collect(context.Background(), map[string]string{historyTestID: "Alpha"}, historyTestRequest())
			if e != nil || r == nil || r.Status != "complete" || len(r.Meters) != 2 || r.Meters[0].MeterID != "m0000" {
				t.Error("concurrent owned collection", r, e)
				return
			}
			evidence := r.Evidence()
			r.Meters[0].MeterID = "caller-mutation"
			if evidence.Meters[0].MeterID != "m0000" {
				t.Error("evidence aliases collection")
			}
		}()
	}
	group.Wait()
}

func TestMeterHistoryCumulativeRawWorkBound(t *testing.T) {
	for _, extra := range []int{0, 1} {
		first := `{"properties":{"columns":[],"rows":[],"nextLink":"?api-version=2021-10-01&$skiptoken=two"},"extra":[` + strings.Repeat("null,", 32761) + `null]}`
		last := `{"properties":{"columns":[],"rows":[]},"extra":[` + strings.Repeat("null,", 32762) + "null" + strings.Repeat(",null", extra) + `]}`
		r, calls := historyRunBodies(t, first, last)
		if calls != 2 {
			t.Fatal("cumulative work not reached", calls)
		}
		if extra == 0 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "partial", "history_limit")
		}
	}
}
func TestMeterHistoryColumnBound(t *testing.T) {
	for _, n := range []int{32, 33} {
		cols := []string{`{"name":"ResourceId","type":"String"}`, `{"name":"ResourceGuid","type":"String"}`, `{"name":"Currency","type":"String"}`, `{"name":"PreTaxCost","type":"Number"}`}
		for i := 4; i < n; i++ {
			cols = append(cols, fmt.Sprintf(`{"name":"extra%d","type":"String"}`, i))
		}
		r, _ := historyRunBodies(t, `{"properties":{"columns":[`+strings.Join(cols, ",")+`],"rows":[]}}`)
		if n == 32 {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "unknown", "history_limit")
		}
	}
}
