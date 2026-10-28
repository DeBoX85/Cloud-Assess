package region

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRetailEvidenceCannotBeForged(t *testing.T) {
	r, e, _ := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
		return []byte(`{"Items":[]}`), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r.Status != "partial" {
		t.Fatal(r, e)
	}
	r.Status = "complete"
	r.FailureCode = ""
	r.Prices = []RetailPrice{{OriginalMeterID: "m1", Item: RetailPriceItem{ArmRegionName: "eastus", RetailPrice: 999}}}
	s, p := r.Evidence()
	if s != "partial" || len(p) != 0 {
		t.Fatal("mutable diagnostics forged evidence", s, p)
	}
	var nilResult *RetailPricingCollection
	s, p = nilResult.Evidence()
	if s != "unavailable" || p == nil || len(p) != 0 {
		t.Fatal(s, p)
	}
}

func TestRetailMetadataConflictAndFuture(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, lower := range []bool{false, true} {
			for _, future := range []bool{false, true} {
				var page map[string]any
				json.Unmarshal([]byte(retailFixture(t)), &page)
				row := page["Items"].([]any)[0].(map[string]any)
				copyRow := func() map[string]any {
					out := map[string]any{}
					for k, v := range row {
						out[k] = v
					}
					return out
				}
				other := copyRow()
				other["serviceId"] = "conflict"
				items := []any{row, other}
				if future {
					other["effectiveStartDate"] = "2027-01-01T00:00:00Z"
				}
				if lower {
					row["tierMinimumUnits"] = 1.
					other["tierMinimumUnits"] = 1.
					low := copyRow()
					low["tierMinimumUnits"] = 0.
					items = append(items, low)
				}
				if reverse {
					for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
						items[i], items[j] = items[j], items[i]
					}
				}
				page["Items"] = items
				body, _ := json.Marshal(page)
				r, e, n := retailRun(t, func(n int, _ string) ([]byte, *http.Response, error) {
					if n == 1 {
						return body, &http.Response{StatusCode: 200}, nil
					}
					return []byte(retailFixture(t)), &http.Response{StatusCode: 200}, nil
				})
				want := "partial"
				calls := 1
				if lower || future {
					want = "complete"
					calls = 2
				}
				if e != nil || r.Status != want || n != calls {
					t.Fatal("metadata conflict/date/minimum tier", reverse, lower, future, r, e, n)
				}
			}
		}
	}
}

func TestRetailContinuationMatrix(t *testing.T) {
	filter := "currencyCode eq 'USD' and type eq 'Consumption'"
	initial := "https://prices.azure.com/api/retail/prices?" + url.Values{"$filter": {filter}}.Encode()
	valid := "?" + url.Values{"$filter": {filter}, "$skip": {"001"}}.Encode()
	for _, raw := range []string{valid, "https://prices.azure.com:443/api/retail/prices" + valid, "/api/retail/prices" + valid} {
		out, n, ok := retailNext(raw, initial, filter, 0)
		if !ok || n != 1 || !strings.HasPrefix(out, "https://prices.azure.com/api/retail/prices?") {
			t.Fatal("valid reference", raw, out, n, ok)
		}
	}
	for _, suffix := range []string{"&extra=x", "&api-version=x", "&%24skip=2", "&%24filter=x", "#", "#fragment"} {
		if _, _, ok := retailNext(valid+suffix, initial, filter, 0); ok {
			t.Fatal("query smuggling", suffix)
		}
	}
	for _, raw := range []string{"?%24skip=1", "?%24filter=x&%24skip=1", "https://prices.azure.com/api/retail/%70rices" + valid, "https://prices.azure.com:0443/api/retail/prices" + valid, "https://prices.azure.com.evil/api/retail/prices" + valid, "https://user@prices.azure.com/api/retail/prices" + valid, "?", valid + "\n"} {
		if _, _, ok := retailNext(raw, initial, filter, 0); ok {
			t.Fatal("invalid continuation admitted", raw)
		}
	}
	for _, skip := range []string{"0", "-1", "+1", "1.0", "1e3", "18446744073709551616", ""} {
		raw := "?" + url.Values{"$filter": {filter}, "$skip": {skip}}.Encode()
		if _, _, ok := retailNext(raw, initial, filter, 0); ok {
			t.Fatal("invalid skip", skip)
		}
	}
	if _, _, ok := retailNext(valid, initial, filter, 1); ok {
		t.Fatal("skip failed to advance")
	}
}

func TestRetailWholeJSONAdmission(t *testing.T) {
	for _, extra := range []string{`"x":1e999`, `"x":-1e-999`, `"x":1e-999`, `"x":"\ud800"`, `"x":"\ufffd"`, `"x":"\ufdd0"`, `"x":"\u0000"`, `"S":1,"ſ":2`, `"x":true,"X":false`} {
		body := strings.TrimSuffix(retailFixture(t), "}") + "," + extra + "}"
		r, e, _ := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
			return []byte(body), &http.Response{StatusCode: 200}, nil
		})
		if e != nil || r.FailureCode != "retail_decode_invalid" || len(r.RawItems) != 0 {
			t.Fatal("whole JSON admission", extra, r, e)
		}
	}
	for _, valid := range []string{`{"x":-0e999}`, `{"x":0e-999}`, `{"x":1e-300}`} {
		work := 0
		if code := retailJSON(context.Background(), []byte(valid), &work); code != "" {
			t.Fatal(valid, code)
		}
	}
	work := 0
	if retailJSON(context.Background(), []byte{'{', '"', 'x', '"', ':', '"', 255, '"', '}'}, &work) != "retail_decode_invalid" {
		t.Fatal("invalid UTF8 admitted")
	}
	for _, cell := range []string{strings.Repeat("x", 513), "bad\ncell"} {
		body := strings.Replace(retailFixture(t), `"meterName":"Meter"`, `"meterName":`+strconvJSON(cell), 1)
		r, e, _ := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
			return []byte(body), &http.Response{StatusCode: 200}, nil
		})
		if e != nil || r.FailureCode != "retail_decode_invalid" {
			t.Fatal("known cell bound", r, e)
		}
	}
}
func strconvJSON(s string) string { b, _ := json.Marshal(s); return string(b) }

func TestRetailJSONExactBudgets(t *testing.T) {
	ctx := context.Background()
	for _, v := range []struct {
		body   string
		start  int
		want   string
		tokens int
	}{
		{`[]`, 1048574, "", 1048576}, {`[]`, 1048575, "retail_limit", 1048576},
		{`{"x":0}`, 1048572, "", 1048576}, {`{"x":0}`, 1048573, "retail_limit", 1048576},
		{strings.Repeat("[", 32) + "0" + strings.Repeat("]", 32), 0, "", 65},
		{strings.Repeat("[", 33) + "0" + strings.Repeat("]", 33), 0, "retail_limit", 33},
		{strconvJSON(strings.Repeat("x", 8192)), 0, "", 1}, {strconvJSON(strings.Repeat("x", 8193)), 0, "retail_limit", 1},
	} {
		work := v.start
		if code := retailJSON(ctx, []byte(v.body), &work); code != v.want || work != v.tokens {
			t.Fatal("exact token/depth/string bound", v.start, v.want, code, work)
		}
	}
	for _, n := range []int{128, 129} {
		fields := []string{}
		for i := range n {
			fields = append(fields, fmt.Sprintf(`"k%d":0`, i))
		}
		work := 0
		want := ""
		if n == 129 {
			want = "retail_limit"
		}
		if code := retailJSON(ctx, []byte("{"+strings.Join(fields, ",")+"}"), &work); code != want {
			t.Fatal("object key bound", n, code)
		}
	}
}

func TestRetailRunRequestBudget(t *testing.T) {
	calls := 0
	c, _ := NewRetailCollector(retailTestGetter(func(_ context.Context, raw string, _ int64) ([]byte, *http.Response, error) {
		calls++
		u, _ := url.Parse(raw)
		q := u.Query()
		skip := q.Get("$skip")
		next := ""
		if skip != "7" {
			n := 1
			if skip != "" {
				fmt.Sscan(skip, &n)
				n++
			}
			q.Set("$skip", fmt.Sprint(n))
			next = "?" + q.Encode()
		}
		return []byte(`{"Items":[],"NextPageLink":` + strconvJSON(next) + `}`), &http.Response{StatusCode: 200}, nil
	}))
	q := retailQ()
	q.MeterIDs = make([]string, 8192)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%04d", i)
	}
	r, e := c.Collect(context.Background(), q)
	if e != nil || r.Status != "partial" || r.FailureCode != "retail_limit" || calls != 4096 || r.Counts.Requests != 4096 {
		t.Fatal("global request cap", r, e, calls)
	}
}

func TestRetailRunByteBudget(t *testing.T) {
	r, e, n := retailRun(t, func(n int, raw string) ([]byte, *http.Response, error) {
		u, _ := url.Parse(raw)
		q := u.Query()
		q.Set("$skip", fmt.Sprint(n))
		body := `{"Items":[],"NextPageLink":` + strconvJSON("?"+q.Encode()) + `}`
		return []byte(body + strings.Repeat(" ", (1<<20)-len(body))), &http.Response{StatusCode: 200}, nil
	})
	if e != nil || r.FailureCode != "retail_limit" || n != 32 || r.Counts.ReceivedBytes != 32<<20 || r.Counts.Requests != 32 {
		t.Fatal("global byte budget", r, e, n)
	}
}

func TestRetailURLShrinkAndPricingPacking(t *testing.T) {
	q := retailQ()
	q.MeterIDs = make([]string, 16)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("%02d", i) + strings.Repeat("'", 510)
	}
	seen := []string{}
	c, _ := NewRetailCollector(retailTestGetter(func(_ context.Context, raw string, _ int64) ([]byte, *http.Response, error) {
		if len(raw) > 8192 {
			t.Fatal("URL length", len(raw))
		}
		u, _ := url.Parse(raw)
		filter := u.Query().Get("$filter")
		if strings.Count(filter, "meterId eq '") != 2 {
			t.Fatal("URL must shrink safely", len(raw), filter)
		}
		seen = append(seen, filter)
		return []byte(`{"Items":[]}`), &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(context.Background(), q)
	if e != nil || r.Status != "partial" || len(seen) != 8 {
		t.Fatal(r, e, len(seen))
	}
	q = retailQ()
	q.MeterIDs = make([]string, 11)
	items := []any{}
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%02d", i)
		var page map[string]any
		json.Unmarshal([]byte(retailFixture(t)), &page)
		row := page["Items"].([]any)[0].(map[string]any)
		row["meterId"] = q.MeterIDs[i]
		row["meterName"] = fmt.Sprintf("Meter%02d", i)
		items = append(items, row)
	}
	body, _ := json.Marshal(map[string]any{"Items": items})
	calls := 0
	c, _ = NewRetailCollector(retailTestGetter(func(_ context.Context, raw string, _ int64) ([]byte, *http.Response, error) {
		calls++
		if calls > 1 {
			u, _ := url.Parse(raw)
			want := 5
			if calls == 4 {
				want = 1
			}
			if n := strings.Count(u.Query().Get("$filter"), "meterName eq '"); n != want {
				t.Fatal("five tuple batches", calls, n)
			}
		}
		return body, &http.Response{StatusCode: 200}, nil
	}))
	r, e = c.Collect(context.Background(), q)
	if e != nil || r.Status != "complete" || calls != 4 {
		t.Fatal(r, e, calls)
	}
}

func TestRetailLaterCancellationAndSafeFailures(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		calls := 0
		c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			if calls == 2 {
				return []byte("provider secret"), &http.Response{StatusCode: 500}, fmt.Errorf("secret: %w", cause)
			}
			return []byte(retailFixture(t)), &http.Response{StatusCode: 200}, nil
		}))
		r, e := c.Collect(context.Background(), retailQ())
		if r != nil || e != cause || calls != 2 {
			t.Fatal("later terminal cause", r, e, calls)
		}
	}
	for _, status := range []int{0, 204, 302, 429, 500} {
		r, e, n := retailRun(t, func(int, string) ([]byte, *http.Response, error) {
			return []byte("secret provider body"), &http.Response{StatusCode: status}, nil
		})
		if e != nil || r.Status != "unknown" || r.FailureCode != "retail_request_failed" || n != 1 || len(r.RawItems) != 0 {
			t.Fatal("safe status", r, e, n)
		}
	}
	if c, e := NewRetailCollector(retailTestGetter(nil)); c != nil || e == nil {
		t.Fatal("typed nil getter")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, _ := NewRetailCollector(retailTestGetter(nil))
	_ = c
	c, _ = NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		t.Fatal("canceled request reached getter")
		return nil, nil, nil
	}))
	r, e := c.Collect(ctx, retailQ())
	if r != nil || e != context.Canceled {
		t.Fatal(r, e)
	}
}

func TestRetailExactRawAndBindingBudgets(t *testing.T) {
	// The JSON token ceiling usually dominates 65,536 full rows. Direct decoder
	// boundary checks certify the independent raw counter without claiming that
	// such a page is reachable through the stricter wire/token path.
	for _, start := range []int{65535, 65536} {
		tokens := 0
		count := start
		items, _, code := decodeRetailPage(context.Background(), []byte(retailFixture(t)), &tokens, &count, true)
		if start == 65535 {
			if code != "" || count != 65536 || len(items) != 1 {
				t.Fatal("raw exact boundary", code, count)
			}
		} else if code != "retail_limit" || len(items) != 0 {
			t.Fatal("raw cap", code, count)
		}
	}
	run := &retailCollectionRun{r: &RetailPricingCollection{Counts: RetailCounts{Bindings: 65535}}}
	if run.bind("m1", RetailPriceItem{}, "pricing") != "" || run.bind("m2", RetailPriceItem{}, "pricing") != "retail_limit" || len(run.r.Bindings) != 1 {
		t.Fatal("binding bound")
	}
	run.r.Counts.MatchingWork = 1048575
	if run.work() != "" || run.work() != "retail_limit" || run.r.Counts.MatchingWork != 1048576 {
		t.Fatal("matching bound")
	}
}

func TestRetailConflictSupersededAndDeterministic(t *testing.T) {
	var previous *RetailPricingCollection
	for _, reverse := range []bool{false, true} {
		r, e, _ := retailRun(t, func(n int, _ string) ([]byte, *http.Response, error) {
			if n == 1 {
				return []byte(retailFixture(t)), &http.Response{StatusCode: 200}, nil
			}
			var page map[string]any
			json.Unmarshal([]byte(retailFixture(t)), &page)
			first := page["Items"].([]any)[0].(map[string]any)
			second := map[string]any{}
			third := map[string]any{}
			for k, v := range first {
				second[k] = v
				third[k] = v
			}
			second["meterId"] = "target2"
			third["effectiveStartDate"] = "2026-02-01T00:00:00Z"
			items := []any{first, second, third}
			if reverse {
				items[0], items[2] = items[2], items[0]
			}
			page["Items"] = items
			body, _ := json.Marshal(page)
			return body, &http.Response{StatusCode: 200}, nil
		})
		if e != nil || r.Status != "complete" || len(r.Prices) != 1 || r.Prices[0].Item.EffectiveStartDate != "2026-02-01T00:00:00Z" {
			t.Fatal("later unambiguous supersession", r, e)
		}
		if previous != nil && !reflect.DeepEqual(previous, r) {
			t.Fatal("arrival order changed diagnostics/evidence")
		}
		previous = r
	}
}

func TestRetailRequestExactBoundaries(t *testing.T) {
	q := retailQ()
	q.MeterIDs = make([]string, 2048)
	q.Regions = make([]string, 32)
	for i := range q.MeterIDs {
		q.MeterIDs[i] = fmt.Sprintf("m%d", i)
	}
	for i := range q.Regions {
		q.Regions[i] = fmt.Sprintf("region%d", i)
	}
	if _, _, ok := retailRequest(q); !ok {
		t.Fatal("65,536 cells rejected")
	}
	q.MeterIDs = append(q.MeterIDs, "extra")
	if _, _, ok := retailRequest(q); ok {
		t.Fatal("cell overflow admitted")
	}
	q = retailQ()
	q.MeterIDs = []string{strings.Repeat("x", 512)}
	if _, _, ok := retailRequest(q); !ok {
		t.Fatal("512 byte identity rejected")
	}
	q.MeterIDs[0] += "x"
	if _, _, ok := retailRequest(q); ok {
		t.Fatal("513 byte identity admitted")
	}
	q = retailQ()
	q.Now = time.Time{}
	if _, _, ok := retailRequest(q); ok {
		t.Fatal("zero clock")
	}
}
