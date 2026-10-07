package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

type quotaGetterFunc func(context.Context, string, int64) ([]byte, *http.Response, error)

func (f quotaGetterFunc) GetBoundedWithResponse(c context.Context, u string, n int64) ([]byte, *http.Response, error) {
	return f(c, u, n)
}
func quotaCollector(t *testing.T, g quotaGetterFunc) *RESTQuotaCollector {
	t.Helper()
	c, e := NewRESTQuotaCollector("https://management.azure.com", g)
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func quotaPage(body string) ([]byte, *http.Response, error) {
	return []byte(body), &http.Response{StatusCode: 200}, nil
}
func collectBody(t *testing.T, body string) *RESTQuotaCollection {
	t.Helper()
	c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { return quotaPage(body) })
	r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
	if e != nil {
		t.Fatal(e)
	}
	return r
}

func TestRESTQuotaCollectorLiteralProviders(t *testing.T) {
	for _, f := range []struct{ kind, provider, version string }{{"Network", "Microsoft.Network", "2022-07-01"}, {"SQL", "Microsoft.Sql", "2021-11-01"}, {"App Service", "Microsoft.Web", "2023-01-01"}, {"Storage", "Microsoft.Storage", "2023-01-01"}} {
		t.Run(f.kind, func(t *testing.T) {
			calls := 0
			c := quotaCollector(t, func(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
				calls++
				want := "https://management.azure.com/subscriptions/" + quotaTestID + "/providers/" + f.provider + "/locations/eastus/usages?api-version=" + f.version
				if u != want || n != 1048576 {
					t.Fatalf("literal request: %s/%d", u, n)
				}
				return quotaPage(`{"value":[{"name":"one","currentValue":85,"limit":100},{"name":{"value":"two","localizedValue":"Display"},"currentValue":100,"limit":100},{"name":"unknown","limit":100},{"name":"null","currentValue":null,"limit":null}]}`)
			})
			r, e := c.Collect(context.Background(), quotaScope(), quotaRequest(f.kind))
			if e != nil {
				t.Fatal(e)
			}
			want := &RESTQuotaCollection{Evidence: QuotaEvidence{Request: quotaRequest(f.kind), Status: "complete", Usages: []QuotaUsage{{ResourceName: "one", LocalizedName: "one", Current: 85, Limit: 100, CurrentKnown: true, LimitKnown: true}, {ResourceName: "two", LocalizedName: "Display", Current: 100, Limit: 100, CurrentKnown: true, LimitKnown: true}, {ResourceName: "unknown", LocalizedName: "unknown", Limit: 100, LimitKnown: true}, {ResourceName: "null", LocalizedName: "null"}}}}
			if calls != 1 || !reflect.DeepEqual(r, want) {
				t.Fatalf("literal evidence: %#v", r)
			}
			calc, e := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{quotaRequest(f.kind)}, []QuotaEvidence{r.Evidence})
			if e != nil || len(calc.Rows) != 2 || calc.Health.Status != assessment.StageCompletedWithWarnings {
				t.Fatalf("calculation health: %#v %v", calc, e)
			}
			cells := [][]string{{"Selected", "eastus", f.kind, "one", "85", "100", "15", "15.0%", "OK"}, {"Selected", "eastus", f.kind, "Display", "100", "100", "0", "0.0%", "At/Over Limit"}}
			for i, row := range calc.Table.Rows {
				if !reflect.DeepEqual(row.Cells, cells[i]) {
					t.Fatalf("literal cells: %#v", row)
				}
			}
		})
	}
	if r := collectBody(t, `{"value":[]}`); !reflect.DeepEqual(r, &RESTQuotaCollection{Evidence: QuotaEvidence{Request: quotaRequest("Network"), Status: "complete", Usages: []QuotaUsage{}}}) {
		t.Fatalf("complete empty: %#v", r)
	}
}

func TestRESTQuotaCollectorRequestAdmission(t *testing.T) {
	g := quotaGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		t.Fatal("invalid request reached getter")
		return nil, nil, nil
	})
	for _, origin := range []string{"http://management.azure.com", "https://a..b", "https://a.-b", "https://a.b-", "https://" + strings.Repeat("x", 64) + ".example", "https://user@management.azure.com", "https://management.azure.com/x", "https://management.azure.com?", "https://management.azure.com#", "https://management.azure.com:444", "https://månagement.azure.com", "https://management.azure.com/%73", "https://management.azure.com\\evil", "https://management.azure.com?x=1", "https://"} {
		if _, e := NewRESTQuotaCollector(origin, g); e == nil {
			t.Fatalf("origin admitted: %q", origin)
		}
	}
	if _, e := NewRESTQuotaCollector("https://management.azure.com", nil); e == nil {
		t.Fatal("nil getter")
	}
	c := quotaCollector(t, g)
	for _, req := range []QuotaRequest{{quotaTestID, "eastus", "VM"}, {quotaTestID, "EastUS", "Network"}, {"bbbbbbbb-bbbb-cccc-dddd-eeeeeeeeeeee", "eastus", "Network"}, {quotaTestID, "eastus/evil", "Network"}, {"invalid", "eastus", "Network"}} {
		if r, e := c.Collect(context.Background(), quotaScope(), req); e == nil || r != nil {
			t.Fatalf("request admitted: %#v", req)
		}
	}
	for _, scope := range []map[string]string{{quotaTestID: ""}, {"invalid": "Name"}, {quotaTestID: "One", strings.ToUpper(quotaTestID): "Two"}, {quotaTestID: strings.Repeat("x", 513)}} {
		if r, e := c.Collect(context.Background(), scope, quotaRequest("Network")); e == nil || r != nil {
			t.Fatalf("scope admitted: %#v", scope)
		}
	}
}

func TestRESTQuotaCollectorContinuation(t *testing.T) {
	path := "/subscriptions/" + quotaTestID + "/providers/Microsoft.Network/locations/eastus/usages"
	for _, next := range []string{"https://evil.example" + path + "?api-version=2022-07-01", "http://management.azure.com" + path + "?api-version=2022-07-01", path + "?api-version=2023-01-01", path + "?api-version=2022-07-01&api-version=2022-07-01", path + "?api-version=2022-07-01&other=x", path + "?api-version=2022-07-01&$skiptoken=", path + "?api-version=2022-07-01#", strings.Replace(path, "/subscriptions/", "/other/", 1) + "?api-version=2022-07-01", strings.Replace(path, "eastus", "%65astus", 1) + "?api-version=2022-07-01", "https://user@management.azure.com" + path + "?api-version=2022-07-01"} {
		calls := 0
		b, _ := json.Marshal(map[string]any{"value": []any{map[string]any{"name": "prefix", "currentValue": 1, "limit": 10}}, "nextLink": next})
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return quotaPage(string(b))
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
		if e != nil || calls != 1 || r.Evidence.Status != "partial" || len(r.Evidence.Usages) != 1 {
			t.Fatalf("unsafe continuation reached auth/lost prefix: %q %#v %v calls%d", next, r, e, calls)
		}
	}
	calls := 0
	c := quotaCollector(t, func(_ context.Context, u string, _ int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 1 {
			return quotaPage(`{"value":[],"nextLink":"?api-version=2022-07-01&$skiptoken=two"}`)
		}
		if !strings.HasSuffix(u, "?%24skiptoken=two&api-version=2022-07-01") {
			t.Fatalf("relative continuation: %s", u)
		}
		return quotaPage(`{"value":[{"name":"two","currentValue":1,"limit":10}]}`)
	})
	r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
	if e != nil || calls != 2 || r.Evidence.Status != "complete" || len(r.Evidence.Usages) != 1 {
		t.Fatalf("safe pagination: %#v %v", r, e)
	}
	calls = 0
	c = quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return quotaPage(`{"value":[],"nextLink":"?api-version=2022-07-01"}`)
	})
	r, e = c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
	if e != nil || calls != 1 || r.FailureCode != "quota_continuation_cycle" || r.Evidence.Status != "partial" {
		t.Fatalf("cycle: %#v %v", r, e)
	}
}

func TestRESTQuotaCollectorFailureHealth(t *testing.T) {
	for _, status := range []int{201, 204, 302, 400, 401, 403, 404, 405, 429, 500} {
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return []byte(`{"value":[]}`), &http.Response{StatusCode: status}, nil
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
		want := "unknown"
		if status == 404 || status == 405 {
			want = "unsupported"
		}
		if e != nil || r.Evidence.Status != want || r.FailureCode != "quota_status_invalid" {
			t.Fatalf("non200 with no transport error%d: %#v %v", status, r, e)
		}
	}
	if r := collectBody(t, `{"value":[{"name":"prefix"}],"nextLink":"`+strings.Repeat("x", 8193)+`"}`); r.Evidence.Status != "unknown" || len(r.Evidence.Usages) != 0 || r.FailureCode != "quota_invalid_response" {
		t.Fatalf("oversized nextLink invalidates entire current page: %#v", r)
	}
	for _, status := range []int{201, 204, 302, 400, 401, 403, 404, 405, 429, 500} {
		for _, prefix := range []bool{false, true} {
			calls := 0
			c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
				calls++
				if prefix && calls == 1 {
					return quotaPage(`{"value":[],"nextLink":"?api-version=2022-07-01&$skiptoken=later"}`)
				}
				return nil, &http.Response{StatusCode: status}, errors.New("secret provider body/URL/token")
			})
			r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
			want := "unknown"
			if status == 404 || status == 405 {
				want = "unsupported"
			}
			if prefix {
				want = "partial"
			}
			if e != nil || r.Evidence.Status != want || r.FailureCode != "quota_request_failed" || len(r.Evidence.Usages) != 0 {
				t.Fatalf("status %d prefix%v: %#v %v", status, prefix, r, e)
			}
		}
	}
	r := collectBody(t, `{"value":[],"nextLink":"?api-version=2022-07-01&$skiptoken=again"}`)
	if r.Evidence.Status != "partial" {
		t.Fatal("empty validated page failure became unknown")
	}
	c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { return nil, nil, nil })
	r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
	if e != nil || r.Evidence.Status != "unknown" {
		t.Fatalf("missing status: %#v %v", r, e)
	}
	for _, body := range []string{`{"value":[null]}`, `{}`, `{"value":null}`, `{"value":{}}`, `{"value":[],"VALUE":[]}`, `{"value":[{"name":"x","Name":"y"}]}`, `{"value":[]} {}`, `{"value":[{"name":"x","currentValue":1.5}]}`, `{"value":[{"name":"x","currentValue":9223372036854775808}]}`, `{"value":[{"name":"x","currentValue":-1}]}`, `{"value":[{"name":"x","limit":1000000000001}]}`, `{"value":[{"name":"\ud800"}]}`, `{"value":[{"name":"\u0000"}]}`, `{"value":[{"name":"x","extra":"\ud800"}]}`, `{"value":[{"name":[]}]}`, `{"value":[],"nextLink":12}`, `{"value":[],"extra":"` + strings.Repeat("x", 513) + `"}`, string([]byte{'{', '"', 0xff, '"', ':', '0', '}'})} {
		if r := collectBody(t, body); r.Evidence.Status != "unknown" || len(r.Evidence.Usages) != 0 {
			t.Fatalf("malformed page accepted: %q %#v", body, r)
		}
	}
	for _, body := range []string{`{"value":[{"name":"x"},{"name":"x"}]}`, `{"value":[{"name":"x"},{"name":"x","currentValue":null}]}`} {
		if r := collectBody(t, body); r.FailureCode != "quota_identity_duplicate" || len(r.Evidence.Usages) != 0 {
			t.Fatalf("duplicate page admitted: %#v", r)
		}
	}
	duplicateCalls := 0
	dc := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		duplicateCalls++
		if duplicateCalls == 1 {
			return quotaPage(`{"value":[{"name":"same","currentValue":1,"limit":10}],"nextLink":"?api-version=2022-07-01&$skiptoken=two"}`)
		}
		return quotaPage(`{"value":[{"name":"same","currentValue":null}]}`)
	})
	dr, de := dc.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
	if de != nil || dr.FailureCode != "quota_identity_duplicate" || dr.Evidence.Status != "partial" || len(dr.Evidence.Usages) != 1 {
		t.Fatalf("crosspage unknown duplicate: %#v %v", dr, de)
	}
	calls := 0
	c = quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 1 {
			return quotaPage(`{"value":[{"name":"prefix","currentValue":1,"limit":10}],"nextLink":"?api-version=2022-07-01&$skiptoken=two"}`)
		}
		return quotaPage(`{"value":[{"name":"would-leak","limit":10},{"name":"bad","currentValue":-1}]}`)
	})
	r, e = c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
	if e != nil || r.Evidence.Status != "partial" || len(r.Evidence.Usages) != 1 || r.Evidence.Usages[0].ResourceName != "prefix" {
		t.Fatalf("page atomicity: %#v %v", r, e)
	}
	calc, e := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{quotaRequest("Network")}, []QuotaEvidence{r.Evidence})
	if e != nil || calc.Health.Status != assessment.StageCompletedWithWarnings || calc.Health.Records != 1 {
		t.Fatalf("partial calculation: %#v %v", calc, e)
	}
}

func TestRESTQuotaCollectorBudgets(t *testing.T) {
	for _, count := range []int{8192, 8193} {
		calls := 0
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			start, end := 0, 4096
			if calls == 2 {
				start, end = 4096, count
			}
			rows := make([]map[string]any, 0, end-start)
			for i := start; i < end; i++ {
				rows = append(rows, map[string]any{"name": fmt.Sprintf("CustomDomains%d", i)})
			}
			page := map[string]any{"value": rows}
			if calls == 1 {
				page["nextLink"] = "?api-version=2022-07-01&$skiptoken=two"
			}
			b, _ := json.Marshal(page)
			return quotaPage(string(b))
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
		if e != nil || calls != 2 || count == 8192 && (r.Evidence.Status != "complete" || len(r.Evidence.Usages) != 8192) || count == 8193 && (r.Evidence.Status != "partial" || r.FailureCode != "quota_row_limit" || len(r.Evidence.Usages) != 4096) {
			t.Fatalf("aggregate raw row boundary%d: %#v %v", count, r, e)
		}
	}
	for _, n := range []int{65533, 65534} {
		body := `{"value":[],"extra":[` + strings.Repeat("0,", n-1) + `0]}`
		r := collectBody(t, body)
		if (n == 65533) != (r.Evidence.Status == "complete") {
			t.Fatalf("token boundary%d: %#v", n, r)
		}
	}
	for _, n := range []int{127, 128} {
		fields := []string{`"value":[]`}
		for i := 0; i < n; i++ {
			fields = append(fields, fmt.Sprintf(`"extra%d":0`, i))
		}
		r := collectBody(t, "{"+strings.Join(fields, ",")+"}")
		if (n == 127) != (r.Evidence.Status == "complete") {
			t.Fatalf("object key boundary%d: %#v", n, r)
		}
	}
	for _, n := range []int{8192, 8193} {
		items := make([]map[string]any, n)
		for i := range items {
			items[i] = map[string]any{"name": fmt.Sprintf("counter%d", i), "currentValue": 1, "limit": 10}
		}
		b, _ := json.Marshal(map[string]any{"value": items})
		r := collectBody(t, string(b))
		if n == 8192 && (r.Evidence.Status != "complete" || len(r.Evidence.Usages) != 8192) || n == 8193 && r.Evidence.Status != "unknown" {
			t.Fatalf("raw row boundary%d: %#v", n, r)
		}
	}
	for _, n := range []int{512, 513} {
		r := collectBody(t, `{"value":[{"name":"`+strings.Repeat("x", n)+`","limit":0}]}`)
		if (n == 512) != (r.Evidence.Status == "complete") {
			t.Fatalf("skipped label boundary%d: %#v", n, r)
		}
	}
	for _, n := range []int{1048576, 1048577} {
		body := `{"value":[]}` + strings.Repeat(" ", n-len(`{"value":[]}`))
		r := collectBody(t, body)
		if (n == 1048576) != (r.Evidence.Status == "complete") {
			t.Fatalf("page byte boundary%d: %#v", n, r)
		}
	}
	for _, complete := range []bool{false, true} {
		calls := 0
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			body := fmt.Sprintf(`{"value":[],"nextLink":"?api-version=2022-07-01&$skiptoken=%d"}`, calls)
			if complete && calls == 8 {
				body = `{"value":[]}`
			}
			if calls <= 8 {
				body += strings.Repeat(" ", 1048576-len(body))
			}
			return quotaPage(body)
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
		if e != nil || complete && (r.Evidence.Status != "complete" || calls != 8) || !complete && (r.FailureCode != "quota_byte_limit" || calls != 9 || r.Evidence.Status != "partial") {
			t.Fatalf("total8MiB: calls%d %#v %v", calls, r, e)
		}
	}
	for _, complete := range []bool{false, true} {
		calls := 0
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			if complete && calls == 64 {
				return quotaPage(`{"value":[]}`)
			}
			return quotaPage(fmt.Sprintf(`{"value":[],"nextLink":"?api-version=2022-07-01&$skiptoken=%d"}`, calls))
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
		if e != nil || calls != 64 || complete && r.Evidence.Status != "complete" || !complete && (r.Evidence.Status != "partial" || r.FailureCode != "quota_page_limit") {
			t.Fatalf("page64: calls%d %#v %v", calls, r, e)
		}
	}
	for _, n := range []int{31, 32} {
		body := `{"value":[],"extra":` + strings.Repeat("[", n) + `0` + strings.Repeat("]", n) + `}`
		r := collectBody(t, body)
		if (n == 31) != (r.Evidence.Status == "complete") {
			t.Fatalf("depth boundary%d: %#v", n, r)
		}
	}
}

func TestRESTQuotaCollectorRetainedSource(t *testing.T) {
	raw, e := os.ReadFile("testdata/source-quota-inputs.json")
	if e != nil {
		t.Fatal(e)
	}
	var fixture struct {
		Cases map[string]struct{ Pages []struct{ Body string } }
	}
	if e = json.Unmarshal(raw, &fixture); e != nil {
		t.Fatal(e)
	}
	body := fixture.Cases["names"].Pages[0].Body
	c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { return quotaPage(body) })
	r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("SQL"))
	if e != nil {
		t.Fatal(e)
	}
	want := []QuotaUsage{{ResourceName: "Servers", LocalizedName: "SQL Servers", Current: 9, Limit: 10, CurrentKnown: true, LimitKnown: true}, {ResourceName: "ElasticPools", LocalizedName: "ElasticPools", Current: 1, Limit: 10, CurrentKnown: true, LimitKnown: true}, {Current: 1, Limit: 10, CurrentKnown: true, LimitKnown: true}}
	if r.Evidence.Status != "complete" || !reflect.DeepEqual(r.Evidence.Usages, want) {
		t.Fatalf("retained raw source names/presence: %#v", r)
	}
	calc, e := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{quotaRequest("SQL")}, []QuotaEvidence{r.Evidence})
	if e != nil || len(calc.Rows) != 2 || calc.Health.Status != assessment.StageCompletedWithWarnings {
		t.Fatalf("retained source integration: %#v %v", calc, e)
	}
	if !reflect.DeepEqual(calc.Table.Rows[0].Cells, []string{"Selected", "eastus", "SQL", "SQL Servers", "9", "10", "1", "10.0%", "Near Limit"}) || !reflect.DeepEqual(calc.Table.Rows[1].Cells, []string{"Selected", "eastus", "SQL", "ElasticPools", "1", "10", "9", "90.0%", "OK"}) {
		t.Fatalf("retained complete source cells: %#v", calc.Table.Rows)
	}
	sourceID := "11111111-1111-1111-1111-111111111111"
	pages := fixture.Cases["paged"].Pages
	calls := 0
	c = quotaCollector(t, func(_ context.Context, u string, _ int64) ([]byte, *http.Response, error) {
		if calls >= len(pages) {
			t.Fatal("unexpected source request")
		}
		want := "https://management.azure.com/subscriptions/" + sourceID + "/providers/Microsoft.Web/locations/eastus/usages?api-version=2023-01-01"
		if calls == 1 {
			want = "https://management.azure.com/subscriptions/" + sourceID + "/providers/Microsoft.Web/locations/eastus/usages?%24skiptoken=synthetic-second&api-version=2023-01-01"
		}
		if u != want {
			t.Fatalf("actual retained continuation: %s", u)
		}
		body := pages[calls].Body
		calls++
		return quotaPage(body)
	})
	req := QuotaRequest{sourceID, "eastus", "App Service"}
	r, e = c.Collect(context.Background(), map[string]string{sourceID: "Source fixture"}, req)
	wantPaged := &RESTQuotaCollection{Evidence: QuotaEvidence{Request: req, Status: "complete", Usages: []QuotaUsage{{ResourceName: "first", LocalizedName: "first", Current: 1, Limit: 10, CurrentKnown: true, LimitKnown: true}, {ResourceName: "second", LocalizedName: "second", Current: 9, Limit: 10, CurrentKnown: true, LimitKnown: true}}}}
	if e != nil || calls != 2 || !reflect.DeepEqual(r, wantPaged) {
		t.Fatalf("actual retained paged evidence: %#v %v", r, e)
	}
}

func TestRESTQuotaCollectorIsolationCancellation(t *testing.T) {
	c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return quotaPage(`{"value":[{"name":"counter","currentValue":1,"limit":10}]}`)
	})
	probe := &availabilityCancelContext{Context: context.Background(), at: 100000}
	if _, e := c.Collect(probe, quotaScope(), quotaRequest("Network")); e != nil {
		t.Fatal(e)
	}
	for at := 1; at <= probe.calls; at++ {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		r, e := c.Collect(ctx, quotaScope(), quotaRequest("Network"))
		if r != nil || !errors.Is(e, context.Canceled) {
			t.Fatalf("deterministic cancellation at%d: %#v %v", at, r, e)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
			if e != nil || r.Evidence.Usages[0].ResourceName != "counter" {
				t.Errorf("isolation: %#v %v", r, e)
				return
			}
			r.Evidence.Usages[0].ResourceName = "changed"
		}()
	}
	wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if r, e := c.Collect(ctx, quotaScope(), quotaRequest("Network")); r != nil || !errors.Is(e, context.Canceled) {
		t.Fatalf("pre-cancellation: %#v %v", r, e)
	}
	ctx, cancel = context.WithCancel(context.Background())
	c = quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		cancel()
		return quotaPage(`{"value":[]}`)
	})
	if r, e := c.Collect(ctx, quotaScope(), quotaRequest("Network")); r != nil || !errors.Is(e, context.Canceled) {
		t.Fatalf("late cancellation: %#v %v", r, e)
	}
	c = quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return nil, nil, context.DeadlineExceeded
	})
	if r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network")); r != nil || !errors.Is(e, context.DeadlineExceeded) {
		t.Fatalf("getter deadline: %#v %v", r, e)
	}
}

type quotaCredential struct {
	t     *testing.T
	scope string
}

func (c quotaCredential) GetToken(_ context.Context, o policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if !reflect.DeepEqual(o.Scopes, []string{c.scope}) {
		c.t.Fatalf("ARM audience: %#v", o.Scopes)
	}
	return azcore.AccessToken{Token: "synthetic-quota-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type quotaTransport func(*http.Request) (*http.Response, error)

func (f quotaTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

type quotaClosedBody struct {
	io.Reader
	closed *bool
}

func (b quotaClosedBody) Close() error { *b.closed = true; return nil }

func TestRESTQuotaCollectorAuthenticatedClient(t *testing.T) {
	for _, status := range []int{200, 404, 405, 403, 302} {
		closed := false
		calls := 0
		client := azure.NewHTTPClient(quotaCredential{t, "https://management.usgovcloudapi.net/.default"}, &azure.HTTPClientOptions{MaxRetries: -1, Scope: "https://management.usgovcloudapi.net/.default", Transport: quotaTransport(func(r *http.Request) (*http.Response, error) {
			calls++
			if r.Method != "GET" || r.Body != nil || r.URL.String() != "https://management.usgovcloudapi.net/subscriptions/"+quotaTestID+"/providers/Microsoft.Network/locations/eastus/usages?api-version=2022-07-01" || r.Header.Get("Authorization") != "Bearer synthetic-quota-token" {
				t.Fatalf("authenticated request: %#v", r)
			}
			return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: quotaClosedBody{strings.NewReader(`{"value":[]}`), &closed}, Request: r}, nil
		})})
		c, e := NewRESTQuotaCollector("https://management.usgovcloudapi.net", client)
		if e != nil {
			t.Fatal(e)
		}
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("Network"))
		want := "unknown"
		if status == 200 {
			want = "complete"
		}
		if status == 404 || status == 405 {
			want = "unsupported"
		}
		if e != nil || r.Evidence.Status != want || calls != 1 || !closed {
			t.Fatalf("status/body ownership%d: %#v %v calls%d closed%v", status, r, e, calls, closed)
		}
	}
}
