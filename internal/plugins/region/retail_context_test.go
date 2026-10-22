package region

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// Single caller only. Schedules are manually derived and frozen before controls.
type retailNthContext struct {
	context.Context
	cancel     context.CancelFunc
	enabled    bool
	calls, nth int
}

func (c *retailNthContext) Err() error {
	if c.enabled {
		c.calls++
		if c.calls == c.nth {
			c.cancel()
		}
	}
	return c.Context.Err()
}
func retailNth(n int) *retailNthContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &retailNthContext{Context: ctx, cancel: cancel, nth: n}
}

func TestRetailFinalContext(t *testing.T) {
	// Enable after GET: transport cause1, JSON root2/Items3, decoded cause4,
	// admitted cause5, finish entry6, post-sort7 (plus error retrieval8).
	ctx := retailNth(7)
	defer ctx.cancel()
	c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		ctx.enabled = true
		return []byte(`{"Items":[]}`), &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(ctx, retailQ())
	if r != nil || !errors.Is(e, context.Canceled) || ctx.calls != 8 {
		t.Fatal("final cancellation lost", r, e, ctx.calls)
	}
}

func TestRetailFailureFinalContext(t *testing.T) {
	// Enable after GET: terminal cause1, finish entry2, error retrieval3.
	ctx := retailNth(2)
	defer ctx.cancel()
	c, _ := NewRetailCollector(retailTestGetter(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		ctx.enabled = true
		return nil, &http.Response{StatusCode: 403}, errors.New("provider secret")
	}))
	r, e := c.Collect(ctx, retailQ())
	if r != nil || !errors.Is(e, context.Canceled) || ctx.calls != 3 {
		t.Fatal("failure final cancellation lost", r, e, ctx.calls)
	}
}
