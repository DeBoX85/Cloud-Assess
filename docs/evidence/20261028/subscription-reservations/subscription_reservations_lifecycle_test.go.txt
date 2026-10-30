package region

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

type bulkSignalContext struct {
	context.Context
	signal chan struct{}
	cause  error
}

func (c *bulkSignalContext) Done() <-chan struct{} { return c.signal }
func (c *bulkSignalContext) Err() error {
	select {
	case <-c.signal:
		return c.cause
	default:
		return nil
	}
}

func TestSubscriptionReservationsCancellation(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		for at := 0; at <= 7; at++ {
			t.Run(fmt.Sprintf("%s/%d", cause, at), func(t *testing.T) {
				ctx := &bulkSignalContext{Context: context.Background(), signal: make(chan struct{}), cause: cause}
				if at == 0 {
					close(ctx.signal)
				}
				calls := 0
				get := bulkScenarioGetter(t, "three")
				c := bulkCollector(t, func(callCtx context.Context, u string, n int64) ([]byte, *http.Response, error) {
					calls++
					b, r, e := get(callCtx, u, n)
					if calls == at {
						close(ctx.signal)
					}
					return b, r, e
				})
				result, e := c.CollectSubscription(ctx, map[string]string{bulkSubID: "Selected"}, SubscriptionReservationRequest{SubscriptionID: bulkSubID})
				if result != nil || !errors.Is(e, cause) || e.Error() != cause.Error() || calls != at {
					t.Fatalf("BULK_CANCEL_TERMINAL: at=%d calls=%d result=%#v error=%v", at, calls, result, e)
				}
			})
		}
	}
	calls := 0
	get := bulkScenarioGetter(t, "three")
	c := bulkCollector(t, func(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 2 {
			return nil, nil, fmt.Errorf("DO_NOT_ECHO_PROVIDER_TEXT: %w", context.Canceled)
		}
		return get(ctx, u, n)
	})
	result, e := c.CollectSubscription(context.Background(), map[string]string{bulkSubID: "Selected"}, SubscriptionReservationRequest{SubscriptionID: bulkSubID})
	if result != nil || e != context.Canceled || calls != 2 || strings.Contains(e.Error(), "DO_NOT_ECHO") {
		t.Fatalf("BULK_CONTEXT_SANITIZATION: %#v %v %d", result, e, calls)
	}
}

func TestSubscriptionReservationsCallerAdmission(t *testing.T) {
	calls := 0
	c := bulkCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return bulkHTTP(bulkPage(nil, ""), 200)
	})
	tests := []struct {
		name  string
		ctx   context.Context
		c     *ReservationCollector
		scope map[string]string
		id    string
		code  string
	}{
		{"nil-context", nil, c, map[string]string{bulkSubID: "Selected"}, bulkSubID, "request_invalid"},
		{"nil-collector", context.Background(), nil, map[string]string{bulkSubID: "Selected"}, bulkSubID, "request_invalid"},
		{"zero-collector", context.Background(), &ReservationCollector{}, map[string]string{bulkSubID: "Selected"}, bulkSubID, "request_invalid"},
		{"bad-id", context.Background(), c, map[string]string{bulkSubID: "Selected"}, "DO_NOT_ECHO_BAD_ID", "request_invalid"},
		{"unselected", context.Background(), c, map[string]string{bulkSubID: "Selected"}, bulkForeignID, "request_invalid"},
		{"bad-scope", context.Background(), c, map[string]string{"DO_NOT_ECHO_BAD_SCOPE": "Selected"}, bulkSubID, "scope_invalid"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, e := tc.c.CollectSubscription(tc.ctx, tc.scope, SubscriptionReservationRequest{SubscriptionID: tc.id})
			if r != nil || e == nil || !strings.Contains(e.Error(), "["+tc.code+"]") || strings.Contains(e.Error(), "DO_NOT_ECHO") || calls != 0 {
				t.Fatalf("BULK_CALLER_ADMISSION: %#v %v calls=%d", r, e, calls)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, e := c.CollectSubscription(ctx, map[string]string{"bad": "bad"}, SubscriptionReservationRequest{SubscriptionID: "bad"})
	if r != nil || e != context.Canceled || calls != 0 {
		t.Fatal("BULK_CANCEL_PRIORITY", r, e, calls)
	}
	good, e := c.CollectSubscription(context.Background(), map[string]string{strings.ToUpper(bulkSubID): "Selected"}, SubscriptionReservationRequest{SubscriptionID: strings.ToUpper(bulkSubID)})
	if e != nil || good == nil || good.SubscriptionID != bulkSubID || good.Evidence().SubscriptionID != bulkSubID || calls != 1 {
		t.Fatal("BULK_SCOPE_NORMALIZATION", good, e, calls)
	}
}

func TestSubscriptionReservationsAdmissionOwnership(t *testing.T) {
	for _, c := range []*SubscriptionReservationCollection{nil, {}, {SubscriptionID: bulkSubID, Status: "complete", Regions: []string{"eastus"}}} {
		e := c.Evidence()
		if e.Status != "unknown" || e.SubscriptionID != "" || len(e.Requests) != 0 || len(e.Regional) != 0 {
			t.Fatalf("BULK_FABRICATED_ADMISSION: %#v", e)
		}
	}
	c, _ := bulkCollect(t, "three")
	want := c.Evidence()
	c.SubscriptionID = bulkForeignID
	c.Status = "unknown"
	c.FailureCode = "DO_NOT_ECHO"
	c.Regions[0] = "westus"
	c.Reservations[0].SKU = "changed"
	if !reflect.DeepEqual(c.Evidence(), want) {
		t.Fatal("BULK_DIAGNOSTIC_ALIAS", c.Evidence(), want)
	}
	view := c.Evidence()
	view.Requests[0].Region = "westus"
	view.Regional[0].Status = "unknown"
	view.Regional[0].Reservations[0].SKU = "changed"
	if !reflect.DeepEqual(c.Evidence(), want) {
		t.Fatal("BULK_ACCESSOR_ALIAS", c.Evidence(), want)
	}
	failed, _ := bulkCollect(t, "initial-failure")
	failed.Status = "complete"
	failed.SubscriptionID = bulkForeignID
	failed.FailureCode = ""
	if e := failed.Evidence(); e.Status != "unknown" || e.SubscriptionID != bulkSubID || e.FailureCode != "reservation_request_failed" || len(e.Requests) != 0 {
		t.Fatal("BULK_PREFIX_PROMOTION", e)
	}
	again, _ := bulkCollect(t, "three")
	if !reflect.DeepEqual(again.Evidence(), want) {
		t.Fatal("BULK_RUN_STATE_LEAK")
	}
}

func TestSubscriptionReservationsConcurrent(t *testing.T) {
	get := bulkScenarioGetter(t, "three")
	c := bulkCollector(t, get)
	var wg sync.WaitGroup
	results := make([]*SubscriptionReservationCollection, 2)
	errs := make([]error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = c.CollectSubscription(context.Background(), map[string]string{bulkSubID: "Selected"}, SubscriptionReservationRequest{SubscriptionID: bulkSubID})
		}(i)
	}
	wg.Wait()
	if errs[0] != nil || errs[1] != nil || results[0] == nil || results[1] == nil || !reflect.DeepEqual(results[0].Evidence(), results[1].Evidence()) {
		t.Fatal("BULK_CONCURRENT_RESULTS", errs)
	}
	results[0].Reservations[0].SKU = "changed"
	if results[1].Reservations[0].SKU != "S" || results[0].Evidence().Regional[0].Reservations[0].SKU != "S" {
		t.Fatal("BULK_CONCURRENT_ALIAS")
	}
}

func TestSubscriptionReservationsRegionalCompatibility(t *testing.T) {
	calls := 0
	get := bulkScenarioGetter(t, "three")
	c := bulkCollector(t, func(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
		calls++
		return get(ctx, u, n)
	})
	result, e := c.Collect(context.Background(), map[string]string{bulkSubID: "Selected"}, ReservationRequest{SubscriptionID: bulkSubID, Region: "eastus"})
	if e != nil || result == nil || result.Evidence.Status != "complete" || len(result.Evidence.Reservations) != 1 || result.Evidence.Reservations[0].Region != "eastus" || calls != 3 {
		t.Fatal("BULK_REGIONAL_COMPATIBILITY", result, e, calls)
	}
	result, e = c.Collect(context.Background(), map[string]string{bulkSubID: "Selected"}, ReservationRequest{SubscriptionID: bulkSubID, Region: ""})
	if result != nil || e == nil || calls != 3 {
		t.Fatal("BULK_NO_EMPTY_REGION_BYPASS", result, e, calls)
	}
}
