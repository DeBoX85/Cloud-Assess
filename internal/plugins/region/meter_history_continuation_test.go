package region

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Isolate origin validation: exact subscription path/version, only host differs.
// The frozen foreign-origin case also has a wrong path and remains unchanged.
func TestMeterHistoryContinuationHostIsolation(t *testing.T) {
	calls := 0
	c, e := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(_ context.Context, u string, _ io.ReadSeekCloser, _ int64) ([]byte, *http.Response, error) {
		calls++
		if calls != 1 || u != historyFirst {
			t.Fatal("foreign continuation reached transport", calls, u)
		}
		link := strings.Replace(historyFirst, "management.azure.com", "foreign.example.test", 1) + "&$skiptoken=two"
		return []byte(historyPageWithNext(link)), &http.Response{StatusCode: 200}, nil
	}))
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.Collect(context.Background(), map[string]string{historyTestID: "Alpha"}, historyTestRequest())
	if e != nil || r == nil || calls != 1 {
		t.Fatal("continuation result", r, e, calls)
	}
	historyRequire(t, r, "partial", "history_continuation_invalid")
}

// A positive value halfway below the minimum isolates a weakened floor; the
// frozen 1e-13 case remains independently below both the true and halved floor.
func TestMeterHistoryPositiveFloorIsolation(t *testing.T) {
	b := strings.Replace(historyRows(1), `"USD",1]`, `"USD",5e-13]`, 1)
	r, calls := historyRunBodies(t, b)
	if calls != 1 {
		t.Fatal("floor admission calls", calls)
	}
	historyRequire(t, r, "unknown", "history_decode_invalid")
}

// Excess positive values also encounter the aggregation ceiling. A negative
// excess isolates absolute row-value validation before nonpositive omission.
func TestMeterHistoryNegativeValueCeilingIsolation(t *testing.T) {
	b := strings.Replace(historyRows(1), `"USD",1]`, `"USD",-1000000000001]`, 1)
	r, calls := historyRunBodies(t, b)
	if calls != 1 {
		t.Fatal("negative ceiling admission calls", calls)
	}
	historyRequire(t, r, "unknown", "history_decode_invalid")
}
