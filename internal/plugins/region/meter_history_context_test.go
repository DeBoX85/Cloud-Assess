package region

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

// Single-caller synthetic context. The schedules below are manually derived
// from the unmodified implementation and frozen before compiling reversals.
type historyNthContext struct {
	context.Context
	cancel  context.CancelFunc
	enabled bool
	calls   int
	nth     int
}

func (c *historyNthContext) Err() error {
	if c.enabled {
		c.calls++
		if c.calls == c.nth {
			c.cancel()
		}
	}
	return c.Context.Err()
}
func historyNth(n int) *historyNthContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &historyNthContext{Context: ctx, cancel: cancel, nth: n}
}

func TestMeterHistoryQueryInvalidRequest(t *testing.T) {
	for _, q := range []MeterHistoryRequest{
		{SubscriptionID: "bad", Months: 1, Now: historyTestRequest().Now, Currency: "USD"},
		{SubscriptionID: historyTestID, Months: 0, Now: historyTestRequest().Now, Currency: "USD"},
		{SubscriptionID: historyTestID, Months: 13, Now: historyTestRequest().Now, Currency: "USD"},
		{SubscriptionID: historyTestID, Months: 1, Now: historyTestRequest().Now, Currency: "usd"},
		{SubscriptionID: historyTestID, Months: 1, Currency: "USD"},
		{SubscriptionID: historyTestID, Months: 1, Now: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC), Currency: "USD"},
	} {
		b, e := MeterHistoryQuery(context.Background(), q)
		if b != nil || e == nil || e.Error() != "region meter history rejected [history_request_invalid]" {
			t.Fatal("invalid standalone query admitted", q, string(b), e)
		}
	}
}

func TestMeterHistoryQueryLateContext(t *testing.T) {
	// Entry Err1, final post-marshal Err2. No helper/intermediate Err calls.
	ctx := historyNth(2)
	defer ctx.cancel()
	ctx.enabled = true
	b, e := MeterHistoryQuery(ctx, historyTestRequest())
	if b != nil || !errors.Is(e, context.Canceled) || ctx.calls != 2 {
		t.Fatal("query post-marshal cancellation lost", string(b), e, ctx.calls)
	}
}

func TestMeterHistoryCollectLateContext(t *testing.T) {
	// Enable immediately after POST: terminal Err1; JSON root/properties/
	// columns/rows Err2-5; final post-sort/result Err6. No rows or nextLink.
	ctx := historyNth(6)
	defer ctx.cancel()
	c, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		ctx.enabled = true
		return []byte(historyEmpty), &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(ctx, map[string]string{historyTestID: "Alpha"}, historyTestRequest())
	if r != nil || !errors.Is(e, context.Canceled) || ctx.calls != 6 {
		t.Fatal("collector final-result cancellation lost", r, e, ctx.calls)
	}
}

func TestMeterHistoryFailureLateContext(t *testing.T) {
	// Enable immediately after POST: terminal Err1, fail's final Err2.
	ctx := historyNth(2)
	defer ctx.cancel()
	c, _ := NewMeterHistoryCollector("https://management.azure.com", historyPosterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		ctx.enabled = true
		return nil, &http.Response{StatusCode: 403}, errors.New("provider secret")
	}))
	r, e := c.Collect(ctx, map[string]string{historyTestID: "Alpha"}, historyTestRequest())
	if r != nil || !errors.Is(e, context.Canceled) || ctx.calls != 2 {
		t.Fatal("failure finalisation lost caller cancellation", r, e, ctx.calls)
	}
}
