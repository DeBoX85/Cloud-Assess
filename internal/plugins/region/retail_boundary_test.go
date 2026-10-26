package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func retailFixture(t *testing.T) string {
	t.Helper()
	return `{"BillingCurrency":"USD","Items":[{"currencyCode":"USD","tierMinimumUnits":0,"retailPrice":2,"unitPrice":2,"armRegionName":"eastus","effectiveStartDate":"2026-01-01T00:00:00Z","meterId":"m1","meterName":"Meter","productId":"Product","skuName":"SKU","unitOfMeasure":"1 Hour","type":"Consumption","isPrimaryMeterRegion":true}],"NextPageLink":""}`
}
func retailQ() RetailPricingRequest {
	return RetailPricingRequest{MeterIDs: []string{"m1"}, Regions: []string{"eastus"}, Now: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)}
}
func retailRun(t *testing.T, f func(int, string) ([]byte, *http.Response, error)) (*RetailPricingCollection, error, int) {
	t.Helper()
	n := 0
	c, e := NewRetailCollector(retailTestGetter(func(_ context.Context, u string, _ int64) ([]byte, *http.Response, error) { n++; return f(n, u) }))
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.Collect(context.Background(), retailQ())
	return r, e, n
}
func TestRetailContinuationOwnership(t *testing.T) {
	var firstFilter string
	r, e, n := retailRun(t, func(n int, raw string) ([]byte, *http.Response, error) {
		u, _ := url.Parse(raw)
		q := u.Query()
		body := retailFixture(t)
		if n == 1 {
			firstFilter = q.Get("$filter")
			next := "https://prices.azure.com:443/api/retail/prices?" + url.Values{"$filter": {firstFilter}, "$skip": {"1000"}}.Encode()
			var p map[string]any
			json.Unmarshal([]byte(body), &p)
			p["NextPageLink"] = next
			b, _ := json.Marshal(p)
			body = string(b)
		}
		if n == 2 && (q.Get("$filter") != firstFilter || q.Get("$skip") != "1000") {
			t.Fatal("continuation query lost")
		}
		if n > 3 {
			t.Fatal("unexpected pagination")
		}
		return []byte(body), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r == nil || r.Status != "complete" || n != 3 {
		t.Fatal("complete continuation", r, e, n)
	}
}
func TestRetailContinuationLoop(t *testing.T) {
	r, e, n := retailRun(t, func(n int, raw string) ([]byte, *http.Response, error) {
		u, _ := url.Parse(raw)
		q := u.Query()
		q.Set("$skip", "1")
		u.RawQuery = q.Encode()
		var p map[string]any
		json.Unmarshal([]byte(retailFixture(t)), &p)
		p["NextPageLink"] = u.String()
		b, _ := json.Marshal(p)
		return b, &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r == nil || r.Status != "partial" || r.FailureCode != "retail_continuation_invalid" || n != 2 {
		t.Fatal("loop admitted", r, e, n)
	}
}
func TestRetailPageBound(t *testing.T) {
	r, e, n := retailRun(t, func(n int, raw string) ([]byte, *http.Response, error) {
		u, _ := url.Parse(raw)
		q := u.Query()
		q.Set("$skip", fmt.Sprint(n))
		u.RawQuery = q.Encode()
		var p map[string]any
		json.Unmarshal([]byte(retailFixture(t)), &p)
		p["NextPageLink"] = u.String()
		b, _ := json.Marshal(p)
		return b, &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r == nil || r.FailureCode != "retail_limit" || n != 64 {
		t.Fatal("page cap", r, e, n)
	}
}
func TestRetailWireBound(t *testing.T) {
	r, e, n := retailRun(t, func(_ int, _ string) ([]byte, *http.Response, error) {
		return []byte(strings.Repeat(" ", (1<<20)+1)), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r == nil || r.Status != "unknown" || r.FailureCode != "retail_limit" || n != 1 {
		t.Fatal("wire cap", r, e, n)
	}
}
func TestRetailTypedAndAliasAdmission(t *testing.T) {
	base := retailFixture(t)
	for _, pair := range [][2]string{{`"tierMinimumUnits":0`, `"tierMinimumUnits":null`}, {`"unitPrice":2`, `"unitPrice":"2"`}, {`"isPrimaryMeterRegion":true`, `"isPrimaryMeterRegion":null`}, {`"retailPrice":2`, `"retailPrice":1e-999`}, {`"Items":`, `"Items":[],"items":`}, {`"NextPageLink":""`, `"NextPageLink":"","nextpagelink":"https://evil.example/"`}} {
		r, e, _ := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
			return []byte(strings.Replace(base, pair[0], pair[1], 1)), &http.Response{StatusCode: 200}, nil
		})
		if e != nil || r == nil || r.FailureCode != "retail_decode_invalid" {
			t.Fatal("typed/alias field admitted", pair, r, e)
		}
	}
	r, e, _ := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
		return []byte(strings.Replace(base, `"Items"`, `"items"`, 1)), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r.Status != "complete" {
		t.Fatal("known field case alias lost", r, e)
	}
}
func TestRetailCancellationPrecedence(t *testing.T) {
	for _, wrapped := range []error{context.Canceled, context.DeadlineExceeded} {
		c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return nil, &http.Response{StatusCode: 403}, fmt.Errorf("provider secret: %w", wrapped)
		}))
		r, e := c.Collect(context.Background(), retailQ())
		if r != nil || !errors.Is(e, wrapped) || strings.Contains(e.Error(), "secret") {
			t.Fatal("terminal cause", r, e)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		cancel()
		return []byte(retailFixture(t)), &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(ctx, retailQ())
	if r != nil || !errors.Is(e, context.Canceled) {
		t.Fatal("late cancel lost", r, e)
	}
}
func TestRetailConcurrentOwnedResults(t *testing.T) {
	c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(retailFixture(t)), &http.Response{StatusCode: 200}, nil
	}))
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := c.Collect(context.Background(), retailQ())
			if e != nil || r.Status != "complete" {
				t.Error("concurrent result", r, e)
				return
			}
			if len(r.Meters) != 1 || len(r.Prices) != 1 || len(r.RawItems) != 2 {
				t.Error("diagnostic identity lost", r)
				return
			}
			r.Meters[0].Item.MeterName = "changed"
			r.Prices[0].Item.MeterName = "changed"
			r.RawItems[0].MeterName = "changed"
			s, m := r.Evidence()
			if s != "complete" || !reflect.DeepEqual(m, map[string]map[string]float64{"m1": {"eastus": 2}}) {
				t.Error("diagnostic mutation changed admission", s, m)
			}
		}()
	}
	wg.Wait()
}
func TestRetailCompleteCostIntegration(t *testing.T) {
	for _, partial := range []bool{false, true} {
		c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			body := retailFixture(t)
			if partial {
				body = strings.Replace(body, `"retailPrice":2`, `"retailPrice":0`, 1)
				body = strings.Replace(body, `"unitPrice":2`, `"unitPrice":0`, 1)
			}
			return []byte(body), &http.Response{StatusCode: 200}, nil
		}))
		r, e := c.Collect(context.Background(), retailQ())
		if e != nil {
			t.Fatal(e)
		}
		status, prices := r.Evidence()
		input, _ := costRuntimeInput()
		scope := map[string]string{input[0].SubscriptionID: input[0].SubscriptionName}
		input[0].TargetRegion = "eastus"
		history := map[string]CostHistoryEvidence{input[0].SubscriptionID: {Complete: true, Meters: []HistoricalCostMeter{{MeterID: "m1", HistoricalCost: 1}}}}
		got, e := EnrichCost(context.Background(), scope, input, &CostEnrichmentEvidence{History: history, PricingStatus: status, RegionPricing: prices})
		if e != nil || got == nil {
			t.Fatal(got, e)
		}
		if got.Comparisons[0].HasCostData == partial {
			t.Fatal("partial retail became eligible", got)
		}
	}
}

type retailRoundTripper func(*http.Request) (*http.Response, error)

func (f retailRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type retailClosedBody struct {
	io.Reader
	closed *bool
}

func (b retailClosedBody) Close() error { *b.closed = true; return nil }
func TestRetailDefaultTransportIsolation(t *testing.T) {
	for _, status := range []int{200, 302, 403} {
		c := NewRetailHTTPClient()
		calls := 0
		closed := false
		c.client.Transport = retailRoundTripper(func(r *http.Request) (*http.Response, error) {
			calls++
			if r.Method != "GET" || r.Header.Get("Authorization") != "" || r.Header.Get("Cookie") != "" || r.Header.Get("Proxy-Authorization") != "" {
				t.Fatal("credential or method leak")
			}
			return &http.Response{StatusCode: status, Header: http.Header{"Location": []string{"https://evil.example/"}}, Body: retailClosedBody{strings.NewReader("{}"), &closed}, Request: r}, nil
		})
		b, _, e := c.GetBounded(context.Background(), "https://prices.azure.com/api/retail/prices?$filter=x", 16)
		if calls != 1 || !closed || (status == 200 && (e != nil || string(b) != "{}")) || (status != 200 && e == nil) {
			t.Fatal("transport redirect/status/body handling", status, calls, closed, e)
		}
		if c.client.Jar != nil || c.client.Timeout != 30*time.Second {
			t.Fatal("default state/timeout")
		}
	}
}
func TestRetailDefaultTransportBounds(t *testing.T) {
	c := NewRetailHTTPClient()
	calls := 0
	closed := false
	c.client.Transport = retailRoundTripper(func(r *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: 200, Body: retailClosedBody{strings.NewReader(strings.Repeat("x", 18)), &closed}, Request: r}, nil
	})
	b, _, e := c.GetBounded(context.Background(), "https://prices.azure.com/api/retail/prices?$filter=x", 16)
	if e == nil || b != nil || calls != 1 || !closed {
		t.Fatal("bounded response escaped", b, e, calls, closed)
	}
	for _, u := range []string{"http://prices.azure.com/api/retail/prices", "https://evil.example/api/retail/prices", "https://user@prices.azure.com/api/retail/prices", "https://prices.azure.com/other", "https://prices.azure.com/api/retail/prices#x"} {
		_, _, e = c.GetBounded(context.Background(), u, 16)
		if e == nil || calls != 1 {
			t.Fatal("invalid origin reached transport", u, e, calls)
		}
	}
}

func TestRetailRequestCellBound(t *testing.T) {
	calls := 0
	c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) { calls++; return nil, nil, nil }))
	q := retailQ()
	q.MeterIDs = make([]string, 2049)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%04d", i)
	}
	q.Regions = make([]string, 32)
	for i := range q.Regions {
		q.Regions[i] = fmt.Sprintf("region%d", i)
	}
	r, e := c.Collect(context.Background(), q)
	if r != nil || e == nil || calls != 0 {
		t.Fatal("cell product not bounded before I/O", r, e, calls)
	}
}
func TestRetailQueryEscapingAndPacking(t *testing.T) {
	q := retailQ()
	q.MeterIDs = []string{"m'1"}
	calls := 0
	c, _ := NewRetailCollector(retailTestGetter(func(_ context.Context, raw string, _ int64) ([]byte, *http.Response, error) {
		calls++
		u, _ := url.Parse(raw)
		if !strings.Contains(u.Query().Get("$filter"), "meterId eq 'm''1'") {
			t.Fatal("OData quote escaped incorrectly", u.Query())
		}
		return []byte(`{"Items":[]}`), &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(context.Background(), q)
	if e != nil || r.Status != "partial" || calls != 1 {
		t.Fatal(r, e, calls)
	}
	q = retailQ()
	q.MeterIDs = make([]string, 31)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%02d", i)
	}
	calls = 0
	c, _ = NewRetailCollector(retailTestGetter(func(_ context.Context, raw string, _ int64) ([]byte, *http.Response, error) {
		calls++
		u, _ := url.Parse(raw)
		filter := u.Query().Get("$filter")
		count := strings.Count(filter, "meterId eq '")
		want := 15
		if calls == 3 {
			want = 1
		}
		if len(raw) > 8192 || count != want {
			t.Fatal("metadata packing", calls, count, len(raw))
		}
		return []byte(`{"Items":[]}`), &http.Response{StatusCode: 200}, nil
	}))
	r, e = c.Collect(context.Background(), q)
	if e != nil || r.Status != "partial" || calls != 3 {
		t.Fatal(r, e, calls)
	}
}
func TestRetailPageAtomicAdmission(t *testing.T) {
	r, e, n := retailRun(t, func(n int, raw string) ([]byte, *http.Response, error) {
		body := retailFixture(t)
		if n == 2 {
			body = strings.Replace(body, `}],"NextPageLink"`, `},null],"NextPageLink"`, 1)
		}
		return []byte(body), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r == nil || n != 2 || r.Status != "partial" || r.FailureCode != "retail_decode_invalid" || len(r.RawItems) != 1 || len(r.Prices) != 0 {
		t.Fatal("invalid page retained prefix", r, e, n)
	}
}
func TestRetailJSONDepthAndTextBounds(t *testing.T) {
	for _, suffix := range []string{`,"extra":"` + strings.Repeat("x", 8193) + `"}`, `,"extra":` + strings.Repeat("[", 34) + `0` + strings.Repeat("]", 34) + `}`} {
		body := strings.TrimSuffix(retailFixture(t), "}") + suffix
		r, e, n := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
			return []byte(body), &http.Response{StatusCode: 200}, nil
		})
		if e != nil || r == nil || r.FailureCode != "retail_limit" || n != 1 {
			t.Fatal("JSON depth/text bound", r, e, n)
		}
	}
}
func TestRetailEqualDateIdentityConflict(t *testing.T) {
	r, e, n := retailRun(t, func(n int, _ string) ([]byte, *http.Response, error) {
		body := retailFixture(t)
		if n == 2 {
			var p map[string]any
			json.Unmarshal([]byte(body), &p)
			items := p["Items"].([]any)
			other := map[string]any{}
			for k, v := range items[0].(map[string]any) {
				other[k] = v
			}
			other["meterId"] = "different-target"
			p["Items"] = append(items, other)
			b, _ := json.Marshal(p)
			body = string(b)
		}
		return []byte(body), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r == nil || n != 2 || r.Status != "partial" || r.FailureCode != "retail_incomplete" {
		t.Fatal("identity tie accepted", r, e, n)
	}
}
