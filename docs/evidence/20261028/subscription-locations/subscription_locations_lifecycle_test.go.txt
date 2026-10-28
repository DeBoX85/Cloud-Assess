package region

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestSubscriptionLocationsOwnership(t *testing.T) {
	g := &locGetter{pages: []locPage{{Body: locBody(locWire("eastus", `[{"logicalZone":"1","physicalZone":"eastus-az2"}]`), ""), Status: 200}}}
	got, e := locCollect(t, g)
	if e != nil || got == nil {
		t.Fatal(e)
	}
	original := got.Evidence()
	got.SubscriptionID = "foreign"
	got.Status = "unknown"
	got.FailureCode = "forged"
	got.CountSource = "fallback"
	got.Regions[0] = "redirected"
	got.ZoneCounts["eastus"] = 99
	got.ZoneMappings["eastus"]["1"] = "redirected"
	if !reflect.DeepEqual(got.Evidence(), original) {
		t.Fatal("LOC_DIAGNOSTIC_ALIAS")
	}
	view := got.Evidence()
	view.Regions[0] = "mutated"
	view.ZoneCounts["eastus"] = 98
	view.ZoneMappings["eastus"]["1"] = "mutated"
	if !reflect.DeepEqual(got.Evidence(), original) {
		t.Fatal("LOC_ACCESSOR_ALIAS")
	}
	for _, c := range []*SubscriptionLocationCollection{nil, {}, {SubscriptionID: locSub, Status: "complete", CountSource: "api", Regions: []string{"eastus"}, ZoneCounts: map[string]int{"eastus": 3}}} {
		v := c.Evidence()
		if v.SubscriptionID != "" || v.Status != "unknown" || v.CountSource != "incomplete" || len(v.Regions) != 0 || len(v.ZoneCounts) != 0 || len(v.ZoneMappings) != 0 {
			t.Fatal("LOC_FABRICATED_ADMISSION", v)
		}
	}
}
func TestSubscriptionLocationsCallerAdmission(t *testing.T) {
	for _, name := range []string{"nil-context", "nil-collector", "zero-collector", "bad-id", "unselected", "bad-scope"} {
		t.Run(name, func(t *testing.T) {
			g := &locGetter{}
			c, e := NewSubscriptionLocationCollector("https://management.azure.com", g)
			if e != nil {
				t.Fatal(e)
			}
			ctx := context.Background()
			scope := map[string]string{locSub: "A"}
			q := SubscriptionLocationRequest{SubscriptionID: locSub}
			switch name {
			case "nil-context":
				ctx = nil
			case "nil-collector":
				c = nil
			case "zero-collector":
				c = &SubscriptionLocationCollector{}
			case "bad-id":
				q.SubscriptionID = "invalid"
			case "unselected":
				scope = map[string]string{}
			case "bad-scope":
				scope["not-a-uuid"] = "A"
			}
			got, e := c.Collect(ctx, scope, q)
			if e == nil || got != nil || g.calls() != 0 || strings.Contains(e.Error(), "SENSITIVE") {
				t.Fatal("LOC_CALLER_ADMISSION", got, e, g.calls())
			}
		})
	}
}
func TestSubscriptionLocationsConstructorAdmission(t *testing.T) {
	for _, origin := range []string{"http://management.azure.com", "https://u:p@management.azure.com", "https://management.azure.com/path", "https://management.azure.com?x=y", "https://management.azure.com#fragment", "https://management.azure.com:444", "https://management.azure.com/%2f", "https://bad_host"} {
		if c, e := NewSubscriptionLocationCollector(origin, &locGetter{}); e == nil || c != nil {
			t.Fatal("LOC_ORIGIN_ADMISSION", origin)
		}
	}
	var nilGetter *locGetter
	for _, g := range []RESTQuotaGetter{nil, nilGetter} {
		if c, e := NewSubscriptionLocationCollector("https://management.azure.com", g); e == nil || c != nil {
			t.Fatal("LOC_NIL_GETTER")
		}
	}
}
func TestSubscriptionLocationsCanonicalContext(t *testing.T) {
	for _, kind := range []string{"cancel", "deadline"} {
		t.Run(kind, func(t *testing.T) {
			g := &locGetter{pages: []locPage{{Error: kind}}}
			got, e := locCollect(t, g)
			want := error(context.Canceled)
			if kind == "deadline" {
				want = context.DeadlineExceeded
			}
			if got != nil || e != want || strings.Contains(e.Error(), "SENSITIVE") {
				t.Fatal("LOC_CONTEXT_SANITIZATION", got, e)
			}
		})
	}
}

// Deterministic terminal-context transitions exercise getter boundaries, not timers.
type locTerminalContext struct {
	context.Context
	mu       sync.RWMutex
	terminal error
	done     chan struct{}
}

func (c *locTerminalContext) Err() error            { c.mu.RLock(); defer c.mu.RUnlock(); return c.terminal }
func (c *locTerminalContext) Done() <-chan struct{} { return c.done }
func (c *locTerminalContext) mark(e error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.terminal == nil {
		c.terminal = e
		close(c.done)
	}
}
func TestSubscriptionLocationsCancellation(t *testing.T) {
	for _, kind := range []string{"cancel", "deadline"} {
		for at := 0; at < 4; at++ {
			t.Run(kind+"/"+string(rune('0'+at)), func(t *testing.T) {
				ctx := &locTerminalContext{Context: context.Background(), done: make(chan struct{})}
				terminal := error(context.Canceled)
				if kind == "deadline" {
					terminal = context.DeadlineExceeded
				}
				g := &locGetter{pages: []locPage{{Body: locBody(locWire("eastus", ""), locLink(1)), Status: 200}, {Body: locBody(locWire("westeurope", ""), locLink(2)), Status: 200}, {Body: locBody("", ""), Status: 200}}}
				if at == 0 {
					ctx.mark(terminal)
				} else {
					g.hook = func(i int) {
						if i == at-1 {
							ctx.mark(terminal)
						}
					}
				}
				c, e := NewSubscriptionLocationCollector("https://management.azure.com", g)
				if e != nil {
					t.Fatal(e)
				}
				got, e := c.Collect(ctx, map[string]string{locSub: "A"}, SubscriptionLocationRequest{SubscriptionID: locSub})
				if got != nil || !errors.Is(e, terminal) || g.calls() != at {
					t.Fatal("LOC_CANCEL_BOUNDARY", got, e, g.calls(), at)
				}
			})
		}
	}
}
func TestSubscriptionLocationsScopeSnapshot(t *testing.T) {
	scope := map[string]string{locSub: "A"}
	g := &locGetter{pages: []locPage{{Body: locBody(locWire("eastus", ""), locLink(1)), Status: 200}, {Body: locBody("", ""), Status: 200}}, hook: func(int) { delete(scope, locSub); scope["invalid"] = "SENSITIVE" }}
	c, e := NewSubscriptionLocationCollector("https://management.azure.com", g)
	if e != nil {
		t.Fatal(e)
	}
	got, e := c.Collect(context.Background(), scope, SubscriptionLocationRequest{SubscriptionID: strings.ToUpper(locSub)})
	if e != nil || got == nil || got.Evidence().Status != "complete" || got.Evidence().SubscriptionID != locSub || got.Evidence().ZoneCounts["eastus"] != 3 {
		t.Fatal("LOC_SCOPE_SNAPSHOT", got, e)
	}
}

type locConcurrentGetter struct{}

const locSubB = "22222222-2222-2222-2222-222222222222"

func (locConcurrentGetter) GetBoundedWithResponse(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
	physical := "A"
	if strings.Contains(u, locSubB) {
		physical = "B"
	}
	return []byte(locBody(locWire("eastus", `[{"logicalZone":"1","physicalZone":"`+physical+`"}]`), "")), &http.Response{StatusCode: 200}, nil
}
func TestSubscriptionLocationsConcurrent(t *testing.T) {
	c, e := NewSubscriptionLocationCollector("https://management.azure.com", locConcurrentGetter{})
	if e != nil {
		t.Fatal(e)
	}
	scope := map[string]string{locSub: "A", locSubB: "B"}
	type result struct {
		got      *SubscriptionLocationCollection
		err      error
		id, want string
	}
	ch := make(chan result, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		id, want := locSub, "A"
		if i%2 == 1 {
			id, want = locSubB, "B"
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := c.Collect(context.Background(), scope, SubscriptionLocationRequest{SubscriptionID: id})
			ch <- result{r, e, id, want}
		}()
	}
	wg.Wait()
	close(ch)
	for r := range ch {
		if r.err != nil || r.got == nil {
			t.Fatal("LOC_RUN_RESULT", r.err)
		}
		v := r.got.Evidence()
		if v.SubscriptionID != r.id || v.ZoneMappings["eastus"]["1"] != r.want {
			t.Fatal("LOC_RUN_STATE_LEAK", v, r.id, r.want)
		}
	}
}
func TestSubscriptionLocationsUnsafeContinuations(t *testing.T) {
	links := []string{
		"http://management.azure.com/subscriptions/" + locSub + "/locations?api-version=2022-12-01",
		"https://management.azure.com:444/subscriptions/" + locSub + "/locations?api-version=2022-12-01",
		"/subscriptions/" + locSubB + "/locations?api-version=2022-12-01",
		"/subscriptions/" + locSub + "/%6cocations?api-version=2022-12-01",
		"?api-version=2022-12-01&api-version=2022-12-01",
		"?API-version=2022-12-01", "?api-version=2020-01-01",
		"?api-version=2022-12-01&$skiptoken=a&$skiptoken=b",
		"?api-version=2022-12-01&$skiptoken=", "?api-version=2022-12-01&$skipToken=a",
		"?api-version=2022-12-01&extra=a", "?api-version=2022-12-01#fragment",
		"https://u:p@management.azure.com/subscriptions/" + locSub + "/locations?api-version=2022-12-01",
	}
	for i, link := range links {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			g := &locGetter{pages: []locPage{{Body: locBody("", link), Status: 200}}}
			r, e := locCollect(t, g)
			if e != nil || r == nil || r.Evidence().Status != "partial" || r.Evidence().FailureCode != "locations_unsafe_continuation" || g.calls() != 1 {
				t.Fatal("LOC_UNSAFE_CONTINUATION", i, r, e, g.calls())
			}
		})
	}
}

type locContextKey struct{}
type locContextGetter struct{ seen bool }

func (g *locContextGetter) GetBoundedWithResponse(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
	g.seen = ctx.Value(locContextKey{}) == "context-canary"
	return []byte(`{"value":[]}`), &http.Response{StatusCode: 200}, nil
}
func TestSubscriptionLocationsContextPropagation(t *testing.T) {
	g := &locContextGetter{}
	c, e := NewSubscriptionLocationCollector("https://management.azure.com", g)
	if e != nil {
		t.Fatal(e)
	}
	ctx := context.WithValue(context.Background(), locContextKey{}, "context-canary")
	r, e := c.Collect(ctx, map[string]string{locSub: "A"}, SubscriptionLocationRequest{SubscriptionID: locSub})
	if e != nil || r == nil || r.Evidence().Status != "complete" || !g.seen {
		t.Fatal("LOC_CONTEXT_PROPAGATION", r, e, g.seen)
	}
}
