package region

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Decoder-equivalent known field aliases must not be treated as unknown
// metadata and silently hide a continuation or required property.
func TestMeterHistoryKnownFieldAliases(t *testing.T) {
	base := historyRows(1)
	for _, property := range []string{"Properties", "propertieſ"} {
		b := strings.Replace(base, `"properties"`, `"`+property+`"`, 1)
		for _, pair := range [][2]string{{`"columns"`, `"Columns"`}, {`"rows"`, `"Rows"`}, {`"name"`, `"Name"`}, {`"type"`, `"Type"`}} {
			b = strings.ReplaceAll(b, pair[0], pair[1])
		}
		r, calls := historyRunBodies(t, b)
		historyRequire(t, r, "complete", "")
		if calls != 1 || len(r.Meters) != 1 || r.Meters[0].MeterID != "m0000" {
			t.Fatal("known field alias lost owned history", r, calls)
		}
		b = strings.TrimSuffix(b, "}}") + `,"NextLink":"https://evil.example/query?api-version=2021-10-01"}}`
		r, calls = historyRunBodies(t, b)
		historyRequire(t, r, "partial", "history_continuation_invalid")
		if calls != 1 || len(r.Meters) != 1 {
			t.Fatal("continuation alias silently completed or reached unsafe HTTP", r, calls)
		}
	}
}

func TestMeterHistoryRawPositiveUnderflow(t *testing.T) {
	for _, v := range []string{"1e-999", "0e-999", "-1e-999"} {
		b := strings.Replace(historyRows(1), `,"USD",1]`, `,"USD",`+v+`]`, 1)
		r, _ := historyRunBodies(t, b)
		if v == "1e-999" {
			historyRequire(t, r, "unknown", "history_decode_invalid")
		} else {
			historyRequire(t, r, "complete", "")
			if len(r.Meters) != 0 {
				t.Fatal("nonpositive amount became a weight", r)
			}
		}
	}
}

func TestMeterHistoryOtherScopeCollision(t *testing.T) {
	calls := 0
	c, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(historyEmpty), &http.Response{StatusCode: 200}, nil
	}))
	scope := map[string]string{historyTestID: "Alpha", "ABCDEFAB-1234-4234-8234-123456789ABC": "A", "abcdefab-1234-4234-8234-123456789abc": "B"}
	r, e := c.Collect(context.Background(), scope, historyTestRequest())
	if r != nil || e == nil || calls != 0 {
		t.Fatal("foreign selected-scope collision reached transport", r, e, calls)
	}
}
func TestMeterHistoryForeignZeroWeight(t *testing.T) {
	b := strings.Replace(historyRows(1), historyTestID, "22222222-2222-2222-2222-222222222222", 1)
	b = strings.Replace(b, `,"USD",1]`, `,"USD",0]`, 1)
	r, _ := historyRunBodies(t, b)
	historyRequire(t, r, "unknown", "history_decode_invalid")
}
func TestMeterHistoryParentContextPrecedence(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		cancel()
		return nil, &http.Response{StatusCode: 403}, fmt.Errorf("provider secret: %w", context.DeadlineExceeded)
	}))
	r, e := c.Collect(ctx, map[string]string{historyTestID: "Alpha"}, historyTestRequest())
	if r != nil || !errors.Is(e, context.Canceled) || errors.Is(e, context.DeadlineExceeded) {
		t.Fatal("parent terminal cause lost to wrapped provider/status", r, e)
	}
}
func TestMeterHistoryCompleteEvidenceIntegration(t *testing.T) {
	for _, complete := range []bool{true, false} {
		first := historyRows(2)
		if !complete {
			first = strings.TrimSuffix(first, "}}") + `,"nextLink":"?api-version=2021-10-01&$skiptoken=two"}}`
		}
		calls := 0
		c, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			if calls == 1 {
				return []byte(first), &http.Response{StatusCode: 200}, nil
			}
			if calls != 2 {
				t.Fatal("unexpected integration continuation", calls)
			}
			return []byte(strings.Replace(historyRows(1), `"USD"`, `"EUR"`, 1)), &http.Response{StatusCode: 200}, nil
		}))
		scope := map[string]string{historyTestID: "Alpha"}
		r, e := c.Collect(context.Background(), scope, historyTestRequest())
		if e != nil || r == nil {
			t.Fatal("integration collection", r, e)
		}
		if complete {
			historyRequire(t, r, "complete", "")
		} else {
			historyRequire(t, r, "partial", "history_decode_invalid")
		}
		input, _ := costRuntimeInput()
		input[0].SubscriptionID = historyTestID
		input[0].SubscriptionName = "Alpha"
		before := historyOracleJSON(t, input)
		evidence := &CostEnrichmentEvidence{PricingStatus: "complete", History: map[string]CostHistoryEvidence{historyTestID: r.Evidence()}, RegionPricing: map[string]map[string]float64{"m0000": {"eastus": 1, "westeurope": 2}, "m0001": {"eastus": 1, "westeurope": 2}}}
		got, e := EnrichCost(context.Background(), scope, input, evidence)
		if e != nil || got == nil {
			t.Fatal("collector-to-cost calculation", got, e)
		}
		want := append([]Comparison{}, input...)
		want[0].AvgCostDifference = 0
		want[0].HasCostData = false
		if complete {
			want[0].AvgCostDifference = 100
			want[0].HasCostData = true
			costHealth(t, got, 1)
		} else {
			costHealth(t, got, 1, "cost_history_partial")
		}
		if !reflect.DeepEqual(got.Comparisons, want) || !reflect.DeepEqual(before, historyOracleJSON(t, input)) {
			t.Fatal("incomplete prefix became cost data or changed unrelated input", got.Comparisons, want)
		}
	}
}
