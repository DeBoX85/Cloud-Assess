package region

import (
	"context"
	"errors"
	"fmt"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const cognitiveReadyBody = `{"value":[{"resourceType":"accounts","name":"s","kind":"OpenAI","locations":["eastus"],"restrictions":[]}]}`
const cognitiveFirstURL = "https://management.azure.com/subscriptions/abcdefab-1234-4234-8234-123456789abc/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"

func cognitivePageLink(n int) string { return cognitiveFirstURL + "&$skiptoken=" + fmt.Sprint(n) }
func cognitiveCollectPages(t *testing.T, pages []string) (*CognitiveSKUCollection, int) {
	t.Helper()
	scope, q, _ := cognitiveBoundaryInput()
	calls := 0
	c, e := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
		if limit != 1048576 || calls >= len(pages) {
			t.Fatal("unexpected bounded call", calls, limit)
		}
		calls++
		return []byte(pages[calls-1]), &http.Response{StatusCode: 200}, nil
	}))
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.Collect(context.Background(), scope, q)
	if e != nil || r == nil {
		t.Fatal("unexpected request error", r, e)
	}
	return r, calls
}
func cognitiveChain(body string, n int) string {
	return strings.TrimSuffix(body, "}") + `,"nextLink":"` + cognitivePageLink(n) + `"}`
}
func TestCognitiveSKUByteAndPageBounds(t *testing.T) {
	empty := `{"value":[]}`
	for _, n := range []int{1048576, 1048577} {
		r, _ := cognitiveCollectPages(t, []string{empty + strings.Repeat(" ", n-len(empty))})
		if n == 1048576 {
			if r.Status != "complete" {
				t.Fatal("exact page bound", r)
			}
		} else if r.FailureCode != "cognitive_sku_byte_limit" {
			t.Fatal("page overflow", r)
		}
	}
	for _, extra := range []int{0, 1} {
		pages := []string{}
		for i := 0; i < 9; i++ {
			b := empty
			if i < 8 {
				b = cognitiveChain(b, i+1)
			}
			size := 1048576
			if i == 7 {
				size = 1048576 - len(empty)
			}
			if i == 8 {
				size = len(empty) + extra
			}
			pages = append(pages, b+strings.Repeat(" ", size-len(b)))
		}
		r, calls := cognitiveCollectPages(t, pages)
		if extra == 0 {
			if r.Status != "complete" || calls != 9 {
				t.Fatal("exact aggregate", r, calls)
			}
		} else if r.Status != "partial" || r.FailureCode != "cognitive_sku_aggregate_limit" || calls != 9 {
			t.Fatal("aggregate overrun", r, calls)
		}
	}
	for _, n := range []int{64, 65} {
		pages := []string{}
		for i := 0; i < n; i++ {
			b := empty
			if i < n-1 {
				b = cognitiveChain(b, i+1)
			}
			pages = append(pages, b)
		}
		r, calls := cognitiveCollectPages(t, pages)
		if n == 64 {
			if r.Status != "complete" || calls != 64 {
				t.Fatal("exact pages", r, calls)
			}
		} else if calls != 64 || r.FailureCode != "cognitive_sku_page_limit" {
			t.Fatal("page limit", r, calls)
		}
	}
}
func TestCognitiveSKURowAndWorkBounds(t *testing.T) {
	for _, n := range []int{8192, 8193} {
		pages := []string{}
		for i := 0; i < n; i += 256 {
			rows := []string{}
			end := i + 256
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				rows = append(rows, fmt.Sprintf(`{"resourceType":"accounts","name":"s%d","kind":"OpenAI","locations":["eastus"],"restrictions":[]}`, j))
			}
			b := `{"value":[` + strings.Join(rows, ",") + `]}`
			if end < n {
				b = cognitiveChain(b, len(pages)+1)
			}
			pages = append(pages, b)
		}
		r, _ := cognitiveCollectPages(t, pages)
		if n == 8192 {
			if r.Status != "complete" || len(r.Groups) != 8192 {
				t.Fatal("exact raw rows/groups", r.Status, r.FailureCode, len(r.Groups))
			}
		} else if r.Status != "partial" || r.FailureCode != "cognitive_sku_row_limit" || len(r.Groups) != 8192 {
			t.Fatal("raw row overflow", r.Status, r.FailureCode, len(r.Groups))
		}
	}
	for _, extra := range []int{0, 1} {
		padding := func(n int) string {
			return `{"value":[],"padding":[` + strings.TrimSuffix(strings.Repeat("null,", n), ",") + `]}`
		}
		pages := []string{cognitiveChain(padding(32767), 1), padding(32767 + extra)}
		r, _ := cognitiveCollectPages(t, pages)
		if extra == 0 {
			if r.Status != "complete" {
				t.Fatal("exact cumulative work", r)
			}
		} else if r.Status != "partial" || r.FailureCode != "cognitive_sku_work_limit" {
			t.Fatal("cumulative work overflow", r)
		}
	}
}
func TestCognitiveSKUScopeAliasPreauth(t *testing.T) {
	for _, field := range []string{"scope", "aliases", "label"} {
		for _, over := range []bool{false, true} {
			scope, q, _ := cognitiveBoundaryInput()
			wantCalls := 1
			switch field {
			case "scope":
				n := 1000
				if over {
					n++
				}
				for i := 1; i < n; i++ {
					scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "selected"
				}
			case "aliases":
				n := 256
				if over {
					n++
				}
				q.LocationAliases = map[string]string{}
				for i := 0; i < n; i++ {
					s := fmt.Sprintf("test%d", i)
					q.LocationAliases[s] = s
				}
			case "label":
				n := 512
				if over {
					n++
				}
				scope[cognitiveBoundaryID] = strings.Repeat("x", n)
			}
			if over {
				wantCalls = 0
			}
			calls := 0
			c, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
				calls++
				return []byte(cognitiveReadyBody), &http.Response{StatusCode: 200}, nil
			}))
			r, e := c.Collect(context.Background(), scope, q)
			if calls != wantCalls || !over && (e != nil || r.Status != "complete") || over && (r != nil || e == nil || e.Error() != "Cognitive SKU request rejected [cognitive_sku_request_invalid]") {
				t.Fatal("preauth exact/over", field, over, r, e, calls)
			}
			// Projection independently admits exact alias count and rejects one-over.
			if field == "aliases" {
				_, _, direct := cognitiveBoundaryInput()
				out, e := direct.Catalog(context.Background(), scope, q, nil)
				if over {
					cognitiveRequireLimit(t, out, e)
				} else if e != nil || out.Status != "complete" {
					t.Fatal("exact aliases", out, e)
				}
			}
		}
	}
	for _, endpoint := range []string{"http://management.azure.com", "https://user:private@example.invalid", "https://management.azure.com/path"} {
		if c, e := NewCognitiveSKUCollector(endpoint, storageGetterFunc(nil)); c != nil || e == nil || e.Error() != "invalid Cognitive SKU origin or getter" {
			t.Fatal("sanitized constructor", c, e)
		}
	}
	if c, e := NewCognitiveSKUCollector("https://management.azure.com", nil); c != nil || e == nil || e.Error() != "invalid Cognitive SKU origin or getter" {
		t.Fatal("nil getter", c, e)
	}
}
func TestCognitiveSKUContinuationBounds(t *testing.T) {
	for _, n := range []int{4096, 4097} {
		link := cognitiveFirstURL + "&$skiptoken=" + strings.Repeat("x", n)
		r, calls := cognitiveCollectPages(t, []string{strings.TrimSuffix(cognitiveReadyBody, "}") + `,"nextLink":"` + link + `"}`, `{"value":[]}`})
		if n == 4096 {
			if r.Status != "complete" || calls != 2 {
				t.Fatal("exact token", r, calls)
			}
		} else if r.FailureCode != "cognitive_sku_continuation_invalid" || calls != 1 {
			t.Fatal("token overflow", r, calls)
		}
	}
	// Exact raw URL and one-over, decoded token stays below4096 through escapes.
	for _, n := range []int{8192, 8193} {
		prefix := cognitiveFirstURL + "&$skiptoken="
		remaining := n - len(prefix)
		token := strings.Repeat("%41", remaining/3) + strings.Repeat("b", remaining%3)
		link := prefix + token
		if _, e := url.Parse(link); e != nil {
			t.Fatal(e)
		}
		r, calls := cognitiveCollectPages(t, []string{strings.TrimSuffix(cognitiveReadyBody, "}") + `,"nextLink":"` + link + `"}`, `{"value":[]}`})
		if n == 8192 {
			if r.Status != "complete" || calls != 2 {
				t.Fatal("exact URL", r, calls)
			}
		} else if r.FailureCode != "cognitive_sku_invalid_response" || r.Status != "unknown" || calls != 1 {
			t.Fatal("URL overflow", r, calls)
		}
	}
}

type cognitiveHTTPTransport struct {
	t             *testing.T
	status        int
	body          string
	readErr       bool
	cancel        context.CancelFunc
	calls, closes atomic.Int32
}

func (c *cognitiveHTTPTransport) Do(req *http.Request) (*http.Response, error) {
	c.calls.Add(1)
	if req.Method != "GET" || req.URL.String() != cognitiveFirstURL || req.Header.Get("Authorization") != "Bearer OFFLINE_CANARY" || req.Body != nil && req.Body != http.NoBody {
		c.t.Error("HTTP request boundary", req.Method, req.URL)
	}
	var reader io.Reader = strings.NewReader(c.body)
	if c.readErr {
		reader = cognitiveErrorReader{}
	}
	if c.cancel != nil {
		reader = cognitiveCancelReader{cancel: c.cancel}
	}
	return &http.Response{StatusCode: c.status, Header: make(http.Header), Body: &diskCloseBody{Reader: reader, closes: &c.closes}, Request: req}, nil
}

type cognitiveErrorReader struct{}

func (cognitiveErrorReader) Read([]byte) (int, error) { return 0, errors.New("private body failure") }

type cognitiveCancelReader struct{ cancel context.CancelFunc }

func (c cognitiveCancelReader) Read([]byte) (int, error) { c.cancel(); return 0, context.Canceled }
func TestCognitiveSKUActualHTTPClosure(t *testing.T) {
	for _, mode := range []string{"success", "403", "429", "500", "307", "oversize", "read-error", "cancel"} {
		scope, q, _ := cognitiveBoundaryInput()
		ctx, cancel := context.WithCancel(context.Background())
		credential := &diskCanaryCredential{}
		wire := &cognitiveHTTPTransport{t: t, status: 200, body: cognitiveReadyBody}
		switch mode {
		case "403":
			wire.status = 403
		case "429":
			wire.status = 429
		case "500":
			wire.status = 500
		case "307":
			wire.status = 307
		case "oversize":
			wire.body = strings.Repeat("x", 1048577)
		case "read-error":
			wire.readErr = true
		case "cancel":
			wire.cancel = cancel
		}
		client := azure.NewHTTPClient(credential, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, OperationTimeout: 5 * time.Second, Scope: "https://management.azure.com/.default", Transport: wire})
		c, _ := NewCognitiveSKUCollector("https://management.azure.com", client)
		r, e := c.Collect(ctx, scope, q)
		cancel()
		if wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 || !reflect.DeepEqual(credential.scopes, []string{"https://management.azure.com/.default"}) {
			t.Fatal("HTTP closure/audience", mode, wire.calls.Load(), wire.closes.Load())
		}
		if mode == "cancel" {
			if r != nil || !errors.Is(e, context.Canceled) {
				t.Fatal("HTTP cancel", r, e)
			}
		} else if e != nil || r == nil {
			t.Fatal(mode, r, e)
		} else if mode == "success" {
			if r.Status != "complete" || r.Groups["s"].State != 0 {
				t.Fatal(mode, r)
			}
		} else if mode == "oversize" || mode == "read-error" {
			if r.FailureCode != "cognitive_sku_request_failed" {
				t.Fatal(mode, r)
			}
		} else if r.FailureCode != "cognitive_sku_status_invalid" {
			t.Fatal(mode, r)
		}
	}
}
func TestCognitiveSKUDecodeCancellation(t *testing.T) {
	scope, q, _ := cognitiveBoundaryInput()
	body := `{"value":[],"ignored":[` + strings.TrimSuffix(strings.Repeat(`{"x":"value"},`, 4000), ",") + `]}`
	c, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	ctx := &storageCountingContext{Context: context.Background(), limit: 1000}
	if r, e := c.Collect(ctx, scope, q); r != nil || !errors.Is(e, context.Canceled) {
		t.Fatal("decode cancellation", r, e)
	}
}

func TestCognitiveSKUFinalReturnContext(t *testing.T) {
	// Fixed12 observations for empty-page completion: entry1, scope1, loop1,
	// post-get1, strict JSON2, raw budget4, admission1, final return1.
	// Do not derive this endpoint from the implementation under mutation.
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		scope, q, _ := cognitiveBoundaryInput()
		c, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return []byte(`{"value":[]}`), &http.Response{StatusCode: 200}, nil
		}))
		ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 12, cause: cause}
		r, e := c.Collect(ctx, scope, q)
		if r != nil || e != cause {
			t.Fatal("fixed collector final-return priority", r, e, ctx.calls)
		}
	}
}
func TestCognitiveSKUAtomicRegionalState(t *testing.T) {
	prior := map[string]cognitiveFlags{"s": {listed: true, zones: map[string]bool{"1": true}, scopes: map[string]bool{"eastus": true}}}
	want := map[string]cognitiveFlags{"s": {listed: true, zones: map[string]bool{"1": true}, scopes: map[string]bool{"eastus": true}}}
	_, q, _ := cognitiveBoundaryInput()
	body := `{"value":[{"resourceType":"accounts","name":"s","kind":"OpenAI","locations":["westus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["2"]}}]},null]}`
	if _, code := decodeCognitivePage(context.Background(), []byte(body), q, prior, 1, 0, 0); code != "cognitive_sku_invalid_response" || !reflect.DeepEqual(prior, want) {
		t.Fatal("rejected page mutated regional prefix", code, prior)
	}
}
func TestCognitiveSKUTransportURLBound(t *testing.T) {
	c, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(nil))
	first, _ := url.Parse(cognitiveFirstURL)
	for _, n := range []int{8192, 8193} {
		prefix := cognitiveFirstURL + "&$skiptoken="
		remaining := n - len(prefix)
		link := prefix + strings.Repeat("%41", remaining/3) + strings.Repeat("b", remaining%3)
		_, ok := c.transport.continuation(first, "2024-10-01", link)
		if ok != (n == 8192) {
			t.Fatal("direct transport raw URL bound", n, ok)
		}
	}
}

func TestCognitiveSKUAliasProvenance(t *testing.T) {
	// Display-key mismatch cannot also be rejected by physical-key rebinding.
	scope, q, c := cognitiveBoundaryInput()
	q.LocationAliases = map[string]string{"east us": "westus"}
	calls := 0
	collector, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(cognitiveReadyBody), &http.Response{StatusCode: 200}, nil
	}))
	r, e := collector.Collect(context.Background(), scope, q)
	if r != nil || e == nil || e.Error() != "Cognitive SKU request rejected [cognitive_sku_request_invalid]" || calls != 0 {
		t.Fatal("unverified display alias admitted", r, e, calls)
	}
	out, e := c.Catalog(context.Background(), scope, q, nil)
	if !reflect.DeepEqual(out, CognitiveCatalogEvidence{}) || e == nil || e.Error() != "Cognitive SKU catalog projection rejected [cognitive_sku_projection_invalid]" {
		t.Fatal("projection display alias", out, e)
	}
}
func TestCognitiveSKUFilteredRawRowBound(t *testing.T) {
	// Raw rows include foreign kinds; isolate row admission from the group cap.
	for _, n := range []int{8192, 8193} {
		pages := []string{}
		for i := 0; i < n; i += 256 {
			end := i + 256
			if end > n {
				end = n
			}
			rows := strings.TrimSuffix(strings.Repeat(`{"resourceType":"accounts","name":"s","kind":"Face","locations":["eastus"],"restrictions":[]},`, end-i), ",")
			b := `{"value":[` + rows + `]}`
			if end < n {
				b = cognitiveChain(b, len(pages)+1)
			}
			pages = append(pages, b)
		}
		r, _ := cognitiveCollectPages(t, pages)
		if len(r.Groups) != 0 || n == 8192 && r.Status != "complete" || n == 8193 && (r.Status != "partial" || r.FailureCode != "cognitive_sku_row_limit") {
			t.Fatal("filtered raw row bound", n, r.Status, r.FailureCode, len(r.Groups))
		}
	}
}

func TestCognitiveSKUUnresolvedLocationLabels(t *testing.T) {
	for _, label := range []string{"", " eastus", "東日本"} {
		body := strings.Replace(cognitiveReadyBody, `["eastus"]`, `["`+label+`"]`, 1)
		r, _ := cognitiveCollectPages(t, []string{body})
		if r.Status != "complete" || r.Groups["s"].Known || r.Groups["s"].State != 2 {
			t.Fatal("valid unresolved location must remain unknown", label, r)
		}
	}
}

func TestCognitiveSKUAliasLabelBounds(t *testing.T) {
	for _, n := range []int{512, 513} {
		scope, q, c := cognitiveBoundaryInput()
		key := "e" + strings.Repeat(" ", n-6) + "astus"
		q.LocationAliases = map[string]string{key: "eastus"}
		calls := 0
		collector, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return []byte(cognitiveReadyBody), &http.Response{StatusCode: 200}, nil
		}))
		r, e := collector.Collect(context.Background(), scope, q)
		if n == 512 {
			if r == nil || r.Status != "complete" || e != nil || calls != 1 {
				t.Fatal("exact alias label", r, e, calls)
			}
		} else if r != nil || e == nil || calls != 0 {
			t.Fatal("overlong alias preauth", r, e, calls)
		}
		out, e := c.Catalog(context.Background(), scope, q, nil)
		if n == 512 {
			if e != nil || out.Status != "complete" {
				t.Fatal("exact projection alias label", out, e)
			}
		} else {
			cognitiveRequireLimit(t, out, e)
		}
	}
}
func TestCognitiveSKUOwnedConcurrentCollection(t *testing.T) {
	scope, q, _ := cognitiveBoundaryInput()
	q.LocationAliases = map[string]string{"east us": "eastus"}
	before := vmssSnapshot(t, []any{scope, q})
	body := []byte(`{"value":[{"resourceType":"accounts","name":"s","kind":"OpenAI","locations":["East US"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["1","2"]}}]}]}`)
	c, _ := NewCognitiveSKUCollector("https://management.azure.com", storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return body, &http.Response{StatusCode: 200}, nil
	}))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a, e := c.Collect(context.Background(), scope, q)
			if e != nil || a.Status != "complete" || !reflect.DeepEqual(a.Groups["s"].BlockedZones, []string{"1", "2"}) {
				t.Error("initial collection", e)
				return
			}
			a.Groups["s"].BlockedZones[0] = "99"
			delete(a.Groups, "s")
			b, e := c.Collect(context.Background(), scope, q)
			if e != nil || !reflect.DeepEqual(b.Groups["s"].BlockedZones, []string{"1", "2"}) {
				t.Error("shared collection", e)
			}
		}()
	}
	wg.Wait()
	if !reflect.DeepEqual(before, vmssSnapshot(t, []any{scope, q})) {
		t.Fatal("shared scope/aliases")
	}
}
