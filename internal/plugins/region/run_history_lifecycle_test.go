package region

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestRunHistoryOwnership(t *testing.T) {
	q := rhRequest()
	q.Currency = "NOK"
	scope := rhScope()
	in := rhCollection("NOK", "complete", "", []HistoricalCostMeter{{"m", 3}})
	got, e := AdmitRunHistory(context.Background(), scope, q, in)
	if e != nil || got == nil {
		t.Fatal("RH_OWNERSHIP_SETUP", e)
	}
	scope[rhA] = "mutated"
	in.Meters[0] = HistoricalCostMeter{"mutated", 999}
	in.Currency = "USD"
	in.Status = "unknown"
	got.Request.SubscriptionID = rhB
	got.Request.Currency = "USD"
	got.Status = "partial"
	got.FailureCode = "history_limit"
	got.Meters[0] = HistoricalCostMeter{"diagnostic", 777}
	v := got.Evidence()
	if !v.Admitted || v.Request.SubscriptionID != rhA || v.Request.Currency != "NOK" || v.Status != "complete" || v.FailureCode != "" || v.ObservedMeters != 1 || !v.History.Complete || !reflect.DeepEqual(v.History.Meters, []HistoricalCostMeter{{"m", 3}}) {
		t.Fatal("RH_DIAGNOSTIC_OWNERSHIP", v)
	}
	v.History.Meters[0] = HistoricalCostMeter{"accessor", 888}
	fresh := got.Evidence()
	if !reflect.DeepEqual(fresh.History.Meters, []HistoricalCostMeter{{"m", 3}}) || fresh.Request.Currency != "NOK" || !fresh.History.Complete {
		t.Fatal("RH_ACCESSOR_OWNERSHIP", fresh)
	}
}
func TestRunHistoryInputOwnership(t *testing.T) {
	in := rhCollection("USD", "complete", "", []HistoricalCostMeter{{"m", 3}})
	got, e := AdmitRunHistory(context.Background(), rhScope(), rhRequest(), in)
	if e != nil || got == nil {
		t.Fatal(e)
	}
	in.Meters[0].HistoricalCost = 999
	if !reflect.DeepEqual(got.Evidence().History.Meters, []HistoricalCostMeter{{"m", 3}}) {
		t.Fatal("RH_INPUT_OWNERSHIP")
	}
}
func TestRunHistoryForgery(t *testing.T) {
	for _, got := range []*RunHistoryBinding{nil, {}, {Request: rhRequest(), Status: "complete", Meters: []HistoricalCostMeter{{"m", 3}}}} {
		v := got.Evidence()
		if v.Admitted || !reflect.DeepEqual(v.Request, MeterHistoryRequest{}) || v.Status != "unknown" || v.FailureCode != "" || v.ObservedMeters != 0 || v.History.Complete || v.History.Meters == nil || len(v.History.Meters) != 0 {
			t.Fatal("RH_FORGED_ADMISSION", v)
		}
	}
}
func TestRunHistoryNilContext(t *testing.T) {
	got, e := AdmitRunHistory(nil, rhScope(), rhRequest(), rhCollection("USD", "complete", "", nil))
	rhRejected(t, got, e)
}
func TestRunHistoryCancellation(t *testing.T) {
	for _, kind := range []string{"cancel", "deadline"} {
		t.Run(kind, func(t *testing.T) {
			cause := error(context.Canceled)
			if kind == "deadline" {
				cause = context.DeadlineExceeded
			}
			base := &vmssCountingContext{Context: context.Background(), cause: cause}
			in := rhCollection("USD", "complete", "", []HistoricalCostMeter{{"a", 3}, {"b", 7}})
			got, e := AdmitRunHistory(base, rhScope(), rhRequest(), in)
			if e != nil || got == nil || base.calls < 1 || base.calls > 512 {
				t.Fatal("RH_CONTEXT_BASELINE", e, base.calls)
			}
			for at := 1; at <= base.calls; at++ {
				ctx := &vmssCountingContext{Context: context.Background(), terminalAt: at, cause: cause}
				r, e := AdmitRunHistory(ctx, rhScope(), rhRequest(), in)
				if r != nil || e != cause {
					t.Fatal("RH_CONTEXT_BOUNDARY", at, r, e)
				}
			}
			ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 1, cause: cause}
			q := rhRequest()
			q.Currency = "bad"
			r, e := AdmitRunHistory(ctx, nil, q, nil)
			if r != nil || e != cause {
				t.Fatal("RH_CONTEXT_PRECEDENCE", r, e)
			}
		})
	}
}
func TestRunHistoryConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	results := make(chan RunHistoryEvidence, 2)
	for _, id := range []string{rhA, rhB} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			q := rhRequest()
			q.SubscriptionID = id
			q.Currency = "USD"
			weight := float64(3)
			if id == rhB {
				q.Currency = "NOK"
				weight = 7
			}
			in := rhCollection(q.Currency, "complete", "", []HistoricalCostMeter{{"m", weight}})
			in.SubscriptionID = id
			r, e := AdmitRunHistory(context.Background(), rhScope(), q, in)
			if e != nil || r == nil {
				results <- RunHistoryEvidence{}
				return
			}
			v := r.Evidence()
			v.History.Meters[0].HistoricalCost = 999
			results <- r.Evidence()
		}(id)
	}
	wg.Wait()
	close(results)
	count := 0
	for v := range results {
		count++
		currency, weight := "USD", float64(3)
		if v.Request.SubscriptionID == rhB {
			currency, weight = "NOK", 7
		}
		if !v.Admitted || v.Request.Currency != currency || !reflect.DeepEqual(v.History.Meters, []HistoricalCostMeter{{"m", weight}}) {
			t.Fatal("RH_RUN_ISOLATION", v)
		}
	}
	if count != 2 {
		t.Fatal("RH_RUN_COUNT")
	}
}
func TestRunHistoryCurrencyIntegration(t *testing.T) {
	for _, scale := range []float64{1, 1000} {
		name := "original"
		if scale == 1000 {
			name = "scaled"
		}
		t.Run(name, func(t *testing.T) {
			scope := rhScope()
			history := map[string]CostHistoryEvidence{}
			for _, id := range []string{rhA, rhB} {
				q := rhRequest()
				q.SubscriptionID = id
				q.Currency = "USD"
				a, b := float64(9), float64(1)
				if id == rhB {
					q.Currency = "NOK"
					a, b = scale, 9*scale
				}
				in := rhCollection(q.Currency, "complete", "", []HistoricalCostMeter{{"m1", a}, {"m2", b}})
				in.SubscriptionID = id
				r, e := AdmitRunHistory(context.Background(), scope, q, in)
				if e != nil || r == nil || r.Evidence().Request.Currency != q.Currency {
					t.Fatal("RH_CURRENCY_LABEL", e)
				}
				history[id] = r.Evidence().History
			}
			comparisons := []Comparison{{SubscriptionID: rhA, SubscriptionName: "a", SourceRegion: "eastus", TargetRegion: "westeurope"}, {SubscriptionID: rhB, SubscriptionName: "b", SourceRegion: "eastus", TargetRegion: "westeurope"}}
			result, e := EnrichCost(context.Background(), scope, comparisons, &CostEnrichmentEvidence{History: history, PricingStatus: "complete", RegionPricing: map[string]map[string]float64{"m1": {"eastus": 1, "westeurope": 2}, "m2": {"eastus": 2, "westeurope": 1}}})
			if e != nil || result == nil || len(result.Comparisons) != 2 || result.Comparisons[0].AvgCostDifference != 85 || result.Comparisons[1].AvgCostDifference != -35 || !result.Comparisons[0].HasCostData || !result.Comparisons[1].HasCostData {
				t.Fatal("RH_RELATIVE_LITERAL", result, e)
			}
		})
	}
}
