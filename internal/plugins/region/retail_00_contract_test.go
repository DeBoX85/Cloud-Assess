package region

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

type retailTestGetter func(context.Context, string, int64) ([]byte, *http.Response, error)

func (f retailTestGetter) GetBounded(c context.Context, u string, n int64) ([]byte, *http.Response, error) {
	return f(c, u, n)
}

type retailLiteralCase struct {
	Name    string `json:"name"`
	Request struct {
		MeterIDs []string `json:"meter_ids"`
		Regions  []string `json:"regions"`
		Now      string   `json:"now"`
	} `json:"request"`
	Pages []struct {
		Status int    `json:"status"`
		Body   string `json:"body"`
	} `json:"pages"`
	WantStatus string                        `json:"want_status"`
	WantCode   string                        `json:"want_code"`
	WantCalls  int                           `json:"want_calls"`
	WantPrices map[string]map[string]float64 `json:"want_prices"`
}

// All literal outcomes precede implementation; no expected result calls target helpers.
func TestRetailPricingFrozenCases(t *testing.T) {
	b, e := os.ReadFile("testdata/retail-pricing-contract-v1.json")
	if e != nil {
		t.Fatal(e)
	}
	var contract struct {
		Format string              `json:"format"`
		Frozen bool                `json:"frozen_before_production"`
		Cases  []retailLiteralCase `json:"cases"`
	}
	if e = json.Unmarshal(b, &contract); e != nil || contract.Format != "retail-pricing-contract-v1" || !contract.Frozen || len(contract.Cases) != 52 {
		t.Fatal("frozen retail inventory", e)
	}
	seen := map[string]bool{}
	for _, v := range contract.Cases {
		if seen[v.Name] {
			t.Fatal("duplicate literal identity", v.Name)
		}
		seen[v.Name] = true
		t.Run(v.Name, func(t *testing.T) {
			calls := 0
			getter := retailTestGetter(func(ctx context.Context, raw string, limit int64) ([]byte, *http.Response, error) {
				if limit != 1<<20 {
					t.Fatal("per-page wire bound", limit)
				}
				u, e := url.Parse(raw)
				if e != nil || u.Scheme != "https" || u.Host != "prices.azure.com" || u.Path != "/api/retail/prices" || u.User != nil || u.Fragment != "" {
					t.Fatal("unexpected retail origin", raw)
				}
				query, e := url.ParseQuery(u.RawQuery)
				if e != nil || len(query) != 1 || len(query["$filter"]) != 1 {
					t.Fatal("initial query changed", query, e)
				}
				filter := query.Get("$filter")
				var expected string
				if calls == 0 {
					ids := append([]string{}, v.Request.MeterIDs...)
					sort.Strings(ids)
					clauses := []string{}
					for _, id := range ids {
						clauses = append(clauses, "meterId eq '"+strings.ReplaceAll(id, "'", "''")+"'")
					}
					expected = "currencyCode eq 'USD' and type eq 'Consumption' and isPrimaryMeterRegion eq true and (" + strings.Join(clauses, " or ") + ")"
				} else {
					expected = "currencyCode eq 'USD' and type eq 'Consumption' and ((meterName eq 'Meter' and productId eq 'Product' and skuName eq 'SKU'))"
				}
				if filter != expected {
					t.Fatalf("source filter changed\ngot %q\nwant %q", filter, expected)
				}
				if calls >= len(v.Pages) {
					t.Fatal("unexpected HTTP attempt", calls)
				}
				p := v.Pages[calls]
				calls++
				return []byte(p.Body), &http.Response{StatusCode: p.Status}, nil
			})
			c, e := NewRetailCollector(getter)
			if e != nil {
				t.Fatal(e)
			}
			now, e := time.Parse(time.RFC3339, v.Request.Now)
			if e != nil {
				t.Fatal(e)
			}
			q := RetailPricingRequest{MeterIDs: append([]string{}, v.Request.MeterIDs...), Regions: append([]string{}, v.Request.Regions...), Now: now}
			before, _ := json.Marshal(q)
			got, e := c.Collect(context.Background(), q)
			if v.WantCode == "request_error" {
				if e == nil || got != nil {
					t.Fatal("invalid request admitted", got, e)
				}
			} else {
				if e != nil || got == nil || got.Status != v.WantStatus || got.FailureCode != v.WantCode {
					t.Fatal("literal outcome", got, e, v.WantStatus, v.WantCode)
				}
				status, prices := got.Evidence()
				wantStatus := v.WantStatus
				if wantStatus == "unknown" {
					wantStatus = "unavailable"
				}
				if status != wantStatus || !reflect.DeepEqual(prices, v.WantPrices) {
					t.Fatal("complete-only owned pricing evidence", status, prices, v.WantPrices)
				}
				for id, regions := range prices {
					for region := range regions {
						prices[id][region] = 999
					}
				}
				_, again := got.Evidence()
				if !reflect.DeepEqual(again, v.WantPrices) {
					t.Fatal("evidence mutation changed collection")
				}
			}
			after, _ := json.Marshal(q)
			if string(after) != string(before) || calls != v.WantCalls {
				t.Fatal("request changed or HTTP count", calls, v.WantCalls)
			}
		})
	}
	if len(seen) != 52 {
		t.Fatal(fmt.Sprint("missing cases", len(seen)))
	}
}
