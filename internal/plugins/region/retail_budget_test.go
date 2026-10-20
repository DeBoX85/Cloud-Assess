package region

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestRetailDefaultTransportPolicy(t *testing.T) {
	c := NewRetailHTTPClient()
	tr, ok := c.client.Transport.(*http.Transport)
	if !ok || !tr.DisableKeepAlives || !tr.DisableCompression || tr.Proxy != nil || tr.Protocols == nil || !tr.Protocols.HTTP1() || tr.Protocols.HTTP2() || tr.Protocols.UnencryptedHTTP2() || tr.ResponseHeaderTimeout != 30*time.Second || tr.TLSHandshakeTimeout != 30*time.Second || tr.MaxResponseHeaderBytes != 1<<20 || c.client.Jar != nil || c.client.Timeout != 30*time.Second {
		t.Fatal("isolated one-attempt bounded transport", tr)
	}
	if c.client.CheckRedirect(nil, nil) != http.ErrUseLastResponse {
		t.Fatal("redirect policy")
	}
}

func TestRetailPricingBoolRequired(t *testing.T) {
	for _, replacement := range []string{`null`, `"false"`, `0`} {
		r, e, n := retailRun(t, func(n int, _ string) ([]byte, *http.Response, error) {
			body := retailFixture(t)
			if n == 2 {
				body = strings.Replace(body, `"isPrimaryMeterRegion":true`, `"isPrimaryMeterRegion":`+replacement, 1)
			}
			return []byte(body), &http.Response{StatusCode: 200}, nil
		})
		if e != nil || r.Status != "partial" || r.FailureCode != "retail_decode_invalid" || n != 2 || len(r.Prices) != 0 || len(r.RawItems) != 1 {
			t.Fatal("pricing boolean coercion", replacement, r, e, n)
		}
	}
}

func TestRetailUnrepresentableTuple(t *testing.T) {
	r, e, n := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
		var page map[string]any
		json.Unmarshal([]byte(retailFixture(t)), &page)
		row := page["Items"].([]any)[0].(map[string]any)
		for _, key := range []string{"meterName", "productId", "skuName"} {
			row[key] = strings.Repeat("'", 512)
		}
		body, _ := json.Marshal(page)
		return body, &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r.Status != "partial" || r.FailureCode != "retail_limit" || n != 1 || len(r.Meters) != 1 {
		t.Fatal("single tuple URL cannot be truncated", r, e, n)
	}
}

func TestRetailRequestMisuseBeforeTransport(t *testing.T) {
	calls := 0
	c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) { calls++; return nil, nil, nil }))
	for _, q := range []RetailPricingRequest{
		{MeterIDs: []string{"m1"}, Regions: []string{"global"}, Now: retailQ().Now},
		{MeterIDs: []string{"m1"}, Regions: []string{"europe"}, Now: retailQ().Now},
		{MeterIDs: []string{"m1"}, Regions: []string{"east us"}, Now: retailQ().Now},
		{MeterIDs: []string{"m1"}, Regions: []string{"Eastus"}, Now: retailQ().Now},
		{MeterIDs: []string{"m\n1"}, Regions: []string{"eastus"}, Now: retailQ().Now},
		{MeterIDs: []string{string([]byte{255})}, Regions: []string{"eastus"}, Now: retailQ().Now},
	} {
		r, e := c.Collect(context.Background(), q)
		if r != nil || e == nil || e.Error() != "retail request failed" || calls != 0 {
			t.Fatal("invalid request reached transport", r, e, calls)
		}
	}
	q := retailQ()
	q.MeterIDs = make([]string, 8193)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%d", i)
	}
	if r, e := c.Collect(context.Background(), q); r != nil || e == nil || calls != 0 {
		t.Fatal("meter count cap", r, e, calls)
	}
}

func TestRetailRemainingWireBudget(t *testing.T) {
	// A partially consumed run must reduce the next getter bound before I/O.
	calls := 0
	c, _ := NewRetailCollector(retailTestGetter(func(_ context.Context, _ string, limit int64) ([]byte, *http.Response, error) {
		calls++
		if limit != 7 {
			t.Fatal("remaining wire limit", limit)
		}
		return nil, &http.Response{StatusCode: 200}, errRetailLimit
	}))
	run := &retailCollectionRun{ctx: context.Background(), c: c, r: &RetailPricingCollection{Counts: RetailCounts{ReceivedBytes: (32 << 20) - 7}}}
	code, e := run.chain(retailInitial([]string{"meterId eq 'm1'"}, true), true, nil)
	if e != nil || code != "retail_limit" || calls != 1 || run.r.Counts.ReceivedBytes != (32<<20)+1 {
		t.Fatal("overflow detection accounting", code, e, calls, run.r.Counts)
	}
	// Error bodies count even though no diagnostic/provider text escapes.
	r, e, _ := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
		return []byte("secret"), &http.Response{StatusCode: 500}, fmt.Errorf("provider secret")
	})
	if e != nil || r.FailureCode != "retail_request_failed" || r.Counts.ReceivedBytes != 6 || len(r.RawItems) != 0 {
		t.Fatal("error byte accounting", r, e)
	}
}

func TestRetailRunTokenBudget(t *testing.T) {
	zeros := strings.Repeat("0,", 249999) + "0"
	r, e, n := retailRun(t, func(n int, raw string) ([]byte, *http.Response, error) {
		u, _ := url.Parse(raw)
		q := u.Query()
		q.Set("$skip", fmt.Sprint(n))
		return []byte(`{"Items":[],"NextPageLink":` + strconvJSON("?"+q.Encode()) + `,"extra":[` + zeros + `]}`), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r.FailureCode != "retail_limit" || r.Status != "partial" || n != 5 || r.Counts.JSONTokens != 1048576 || r.Counts.RawItems != 0 {
		t.Fatal("global token budget", r, e, n)
	}
}

func TestRetailRunBindingBudget(t *testing.T) {
	q := retailQ()
	q.MeterIDs = make([]string, 8192)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%04d", i)
	}
	var fixture map[string]any
	json.Unmarshal([]byte(retailFixture(t)), &fixture)
	original := fixture["Items"].([]any)[0].(map[string]any)
	calls := 0
	c, _ := NewRetailCollector(retailTestGetter(func(_ context.Context, raw string, _ int64) ([]byte, *http.Response, error) {
		calls++
		u, _ := url.Parse(raw)
		filter := u.Query().Get("$filter")
		items := []any{}
		if strings.Contains(filter, "isPrimaryMeterRegion eq true") {
			for _, clause := range strings.Split(filter, "meterId eq '")[1:] {
				id := strings.SplitN(clause, "'", 2)[0]
				row := map[string]any{}
				for k, v := range original {
					row[k] = v
				}
				row["meterId"] = id
				items = append(items, row)
			}
		} else {
			for range 8 {
				items = append(items, original)
			}
		}
		body, _ := json.Marshal(map[string]any{"Items": items})
		return body, &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(context.Background(), q)
	if e != nil || r.Status != "partial" || r.FailureCode != "retail_limit" || calls != 548 || r.Counts.Bindings != 65536 || len(r.Bindings) != 65536 || len(r.Meters) != 8192 {
		t.Fatal("global shared-tuple binding budget", r, e, calls)
	}
	s, prices := r.Evidence()
	if s != "partial" || len(prices) != 0 {
		t.Fatal("limited prefix became cost evidence", s, len(prices))
	}
}
