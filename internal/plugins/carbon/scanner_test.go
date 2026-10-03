package carbon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const carbonSub = "11111111-1111-4111-8111-111111111111"
const carbonOther = "22222222-2222-4222-8222-222222222222"
const dateJSON = `{"startDate":"2026-01-01","endDate":"2026-03-31"}`

// Same literal synthetic service items fed to the unchanged pinned source
// generator; expected rows/metadata come from its separately hashed captures.
const capturedItems = `[{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Compute/virtualMachines","latestMonthEmissions":80,"previousMonthEmissions":100,"monthlyEmissionsChangeValue":-20,"monthOverMonthEmissionsChangeRatio":999},{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Compute/virtualMachines","latestMonthEmissions":20,"previousMonthEmissions":0},{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Storage/storageAccounts","latestMonthEmissions":75,"previousMonthEmissions":0,"monthlyEmissionsChangeValue":0}]`
const simpleItem = `{"dataType":"ItemDetailsData","categoryType":"ResourceType","itemName":"Microsoft.Compute/virtualMachines","latestMonthEmissions":40,"previousMonthEmissions":100,"monthlyEmissionsChangeValue":-60}`
const allowedJSON = `[{"subscriptionId":"11111111-1111-4111-8111-111111111111","decision":"Allowed"}]`

type posterFunc func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)

func (f posterFunc) PostBounded(ctx context.Context, endpoint string, body io.ReadSeekCloser, limit int64) ([]byte, *http.Response, error) {
	return f(ctx, endpoint, body, limit)
}

type carbonTransportFunc func(*http.Request) (*http.Response, error)

func (f carbonTransportFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

type carbonCredential struct {
	scope string
	calls atomic.Int32
	t     *testing.T
}

func (c *carbonCredential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if !reflect.DeepEqual(options.Scopes, []string{c.scope}) {
		c.t.Errorf("wrong captured carbon audience: %v", options.Scopes)
	}
	c.calls.Add(1)
	return azcore.AccessToken{Token: "synthetic-carbon-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type trackedBody struct {
	*bytes.Reader
	closed atomic.Int32
	read   atomic.Int64
}

func (b *trackedBody) Read(p []byte) (int, error) {
	n, e := b.Reader.Read(p)
	b.read.Add(int64(n))
	return n, e
}
func (b *trackedBody) Close() error { b.closed.Add(1); return nil }
func successful(data string) ([]byte, *http.Response, error) {
	return []byte(data), &http.Response{StatusCode: 200}, nil
}
func pageJSON(items, decisions, token string) string {
	return `{"value":` + items + `,"subscriptionAccessDecisionList":` + decisions + `,"skipToken":` + fmt.Sprintf("%q", token) + `}`
}
func fixtureScanner(t *testing.T, report func(context.Context, []byte) ([]byte, *http.Response, error)) *Scanner {
	t.Helper()
	s, e := NewWithHTTPClient("https://management.fixture.invalid/", posterFunc(func(ctx context.Context, endpoint string, body io.ReadSeekCloser, limit int64) ([]byte, *http.Response, error) {
		if limit != MaxPageBytes {
			t.Error("unbounded response")
		}
		if strings.Contains(endpoint, "queryCarbonEmissionDataAvailableDateRange") {
			if body != nil {
				t.Error("date query must have no body")
			}
			return successful(dateJSON)
		}
		if endpoint != "https://management.fixture.invalid/providers/Microsoft.Carbon/carbonEmissionReports?api-version=2025-04-01" || body == nil {
			t.Fatal("wrong fixed report endpoint/body")
		}
		data, e := io.ReadAll(body)
		if e != nil {
			t.Fatal(e)
		}
		return report(ctx, data)
	}))
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func onlySub() map[string]string { return map[string]string{carbonSub: "fixture"} }
func checkCanonical(t *testing.T, table assessment.PluginTable) {
	t.Helper()
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil || table.Health.Records != len(table.Rows) {
		t.Fatal("invalid canonical carbon health/output")
	}
}

func TestCarbonAuthenticatedCloudBodyBatchingAndEverySourceCell(t *testing.T) {
	for _, mode := range []string{"all", "filtered", "batched", "empty"} {
		t.Run(mode, func(t *testing.T) {
			cred := &carbonCredential{scope: "https://management.usgovcloudapi.net/.default", t: t}
			calls, reports := 0, 0
			var bodies []*trackedBody
			options := azure.DefaultHTTPClientOptions(time.Second)
			options.MaxRetries = -1
			options.Scope = cred.scope
			options.Transport = carbonTransportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Method != "POST" || r.URL.Host != "management.usgovcloudapi.net" || r.URL.RawQuery != "api-version=2025-04-01" || r.Header.Get("Authorization") != "Bearer synthetic-carbon-canary" {
					t.Fatal("read-oriented selected-cloud auth contract")
				}
				data := dateJSON
				switch r.URL.Path {
				case "/providers/Microsoft.Carbon/queryCarbonEmissionDataAvailableDateRange":
					if r.Body != nil {
						b, _ := io.ReadAll(r.Body)
						if len(b) > 0 {
							t.Fatal("unexpected date body")
						}
					}
				case "/providers/Microsoft.Carbon/carbonEmissionReports":
					reports++
					raw, _ := io.ReadAll(r.Body)
					var q map[string]any
					if json.Unmarshal(raw, &q) != nil || len(q) != 8 {
						t.Fatal("unexpected query fields")
					}
					ids := q["subscriptionList"].([]any)
					wantCount := 1
					if mode == "batched" && reports == 1 {
						wantCount = 100
					}
					if len(ids) != wantCount || q["reportType"] != "ItemDetailsReport" || q["categoryType"] != "ResourceType" || q["orderBy"] != "LatestMonthEmissions" || q["sortDirection"] != "Desc" || q["pageSize"] != float64(1000) || !reflect.DeepEqual(q["dateRange"], map[string]any{"start": "2026-03-31", "end": "2026-03-31"}) || !reflect.DeepEqual(q["carbonScopeList"], []any{"Scope1", "Scope2", "Scope3"}) {
						t.Fatalf("source request contract: %s", raw)
					}
					var access []map[string]string
					for i, id := range ids {
						want := carbonSub
						if mode == "batched" {
							want = fmt.Sprintf("%08x-1111-4111-8111-111111111111", (reports-1)*100+i)
						}
						if id != want {
							t.Fatal("scope sorted/batched incorrectly")
						}
						access = append(access, map[string]string{"subscriptionId": want, "decision": "Allowed"})
					}
					a, _ := json.Marshal(access)
					data = pageJSON(capturedItems, string(a), "")
				default:
					t.Fatal("unexpected endpoint or provider registration")
				}
				b := &trackedBody{Reader: bytes.NewReader([]byte(data))}
				bodies = append(bodies, b)
				return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: b, Request: r}, nil
			})
			client := azure.NewHTTPClient(cred, options)
			s, e := NewWithHTTPClient("https://management.usgovcloudapi.net/", client)
			if e != nil {
				t.Fatal(e)
			}
			// Mutating caller options must not change an already-created pipeline.
			options.Scope = "https://wrong.invalid/.default"
			options.Transport = carbonTransportFunc(func(*http.Request) (*http.Response, error) { t.Fatal("mutable captured options"); return nil, nil })
			subs := onlySub()
			var filter Filter
			if mode == "filtered" {
				filter = typeFilter(func(v string) bool { return v != "Microsoft.Compute/virtualMachines" })
			}
			if mode == "batched" {
				subs = map[string]string{}
				for i := 0; i < 101; i++ {
					subs[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "fixture"
				}
			}
			if mode == "empty" {
				subs = map[string]string{}
			}
			table, e := s.Scan(context.Background(), subs, filter)
			if e != nil || table.Health.Status != assessment.StageCompleted {
				t.Fatalf("healthy source request: %+v %v", table.Health, e)
			}
			checkCanonical(t, table)
			rows := [][]string{table.Columns}
			for _, r := range table.Rows {
				if r.SubscriptionID != "" {
					t.Fatal("fabricated aggregate identity")
				}
				rows = append(rows, r.Cells)
			}
			if !reflect.DeepEqual(rows, capture(t, mode)) {
				t.Fatal("HTTP scan differs from actual source cells/metadata")
			}
			want := 2
			if mode == "batched" {
				want = 3
			}
			if mode == "empty" {
				want = 0
			}
			if calls != want || mode == "empty" && cred.calls.Load() != 0 {
				t.Fatal("wrong request/auth count")
			}
			for _, b := range bodies {
				if b.closed.Load() != 1 || b.read.Load() != int64(b.Size()) {
					t.Fatal("response body ownership")
				}
			}
		})
	}
}

func TestCarbonPagingOpaqueTokenAndLaterFailureRetention(t *testing.T) {
	for _, mode := range []string{"healthy", "cycle", "later403", "conflict", "cancel", "page-limit"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			s := fixtureScanner(t, func(c context.Context, body []byte) ([]byte, *http.Response, error) {
				calls++
				var q map[string]any
				_ = json.Unmarshal(body, &q)
				if calls == 1 && q["skipToken"] != nil || calls == 2 && q["skipToken"] != "https://untrusted.invalid/opaque?secret=canary" {
					t.Fatal("opaque token not confined to JSON")
				}
				if mode == "later403" && calls == 2 {
					return nil, &http.Response{StatusCode: 403}, fmt.Errorf("private denial %s", carbonSub)
				}
				if mode == "cancel" && calls == 2 {
					cancel()
					return nil, nil, c.Err()
				}
				decisions := allowedJSON
				if mode == "conflict" && calls == 2 {
					decisions = strings.ReplaceAll(allowedJSON, "Allowed", "Denied")
				}
				token := "https://untrusted.invalid/opaque?secret=canary"
				if calls == 2 && mode != "cycle" {
					token = ""
				}
				if mode == "page-limit" {
					token = fmt.Sprintf("token-%d", calls)
					if calls == 1 {
						token = "https://untrusted.invalid/opaque?secret=canary"
					}
				}
				return successful(pageJSON("["+simpleItem+"]", decisions, token))
			})
			table, e := s.Scan(ctx, onlySub(), nil)
			checkCanonical(t, table)
			want := "40.00"
			if mode == "healthy" {
				want = "80.00"
			}
			if mode == "page-limit" {
				want = "2560.00"
			}
			if len(table.Rows) != 1 || table.Rows[0].Cells[3] != want {
				t.Fatalf("valid prior-page sums: %+v", table.Rows)
			}
			if mode == "healthy" {
				if e != nil || table.Health.Status != assessment.StageCompleted {
					t.Fatal("safe paging failed")
				}
			} else {
				if e == nil || table.Health.Status != assessment.StageFailed || strings.Contains(e.Error(), carbonSub) || strings.Contains(e.Error(), "secret") {
					t.Fatal("incomplete/privacy health")
				}
			}
			if mode == "cycle" && (calls != 2 || table.Health.Error.Code != "carbon_token_cycle") {
				t.Fatal("cycle guard did not reject repeated token")
			}
			if mode == "cancel" && !errors.Is(e, context.Canceled) {
				t.Fatal("lost cancellation identity")
			}
			if mode == "page-limit" && (calls != 64 || table.Health.Error.Code != "carbon_page_limit") {
				t.Fatal("unbounded paging")
			}
		})
	}
}

func TestCarbonAccessDecisionHealthAndMalformedEnvelopes(t *testing.T) {
	for _, c := range []struct {
		name, access    string
		failed, warning bool
	}{
		{"allowed", allowedJSON, false, false},
		{"missing", "null", false, true},
		{"empty", "[]", false, true},
		{"denied", strings.ReplaceAll(allowedJSON, "Allowed", "Denied"), true, true},
		{"foreign", strings.ReplaceAll(allowedJSON, carbonSub, carbonOther), true, false},
		{"duplicate", "[" + strings.Trim(allowedJSON, "[]") + "," + strings.Trim(allowedJSON, "[]") + "]", true, false},
		{"alias", strings.ReplaceAll(allowedJSON, "subscriptionId", "SubscriptionId"), true, false},
		{"unknown", strings.ReplaceAll(allowedJSON, "Allowed", "Unknown"), true, false},
		{"shape", "{}", true, false},
		{"reason", strings.ReplaceAll(allowedJSON, `"decision":"Allowed"`, `"decision":"Allowed","denialReason":7`), true, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
				return successful(pageJSON("["+simpleItem+"]", c.access, ""))
			})
			table, e := s.Scan(context.Background(), onlySub(), nil)
			checkCanonical(t, table)
			if c.name == "denied" && (e == nil || len(table.Rows) != 1 || table.Health.Error == nil || table.Health.Error.Code != "carbon_access_denied") {
				t.Fatal("denied subscription masked by valid rows")
			}
			if (e != nil) != c.failed || (len(table.Health.Warnings) > 0) != c.warning {
				t.Fatalf("access health: %+v %v", table.Health, e)
			}
			if c.warning && !c.failed && table.Health.Status != assessment.StageCompletedWithWarnings {
				t.Fatal("missing access metadata certified complete")
			}
			if c.failed && !c.warning && len(table.Rows) != 0 {
				t.Fatal("ambiguous access page accepted")
			}
		})
	}
	for _, raw := range []string{`{}`, `{"value":null}`, `{"value":{}}`, `{"Value":[]}`, `{"value":[],"value":[]}`, `{"value":[],"error":{"message":"private-canary"}}`, `{"value":[],"skipToken":1}`, `{"value":[],"skipToken":"unsafe\ncanary"}`, `{"value":[],"skipToken":"` + strings.Repeat("x", MaxTokenBytes+1) + `"}`, `{"value":[],"subscriptionAccessDecisionList":[],"SubscriptionAccessDecisionList":[]}`, `[]`, `{"value":[]} trailing`, `{"value":[],"meta":"bad` + string([]byte{255}) + `"}`} {
		t.Run(fmt.Sprintf("envelope-%d", len(raw)), func(t *testing.T) {
			s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) { return successful(raw) })
			v, e := s.Scan(context.Background(), onlySub(), nil)
			if e == nil || v.Health.Status != assessment.StageFailed || len(v.Rows) != 0 {
				t.Fatal("invalid envelope marked healthy")
			}
		})
	}
	for _, items := range []string{"[]", "[" + simpleItem + ",null,{}," + strings.ReplaceAll(simpleItem, `"latestMonthEmissions":40`, `"latestMonthEmissions":"40"`) + "]", "[{}]"} {
		s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
			return successful(pageJSON(items, allowedJSON, ""))
		})
		v, e := s.Scan(context.Background(), onlySub(), nil)
		if e != nil {
			t.Fatal(e)
		}
		if items == "[]" {
			if v.Health.Status != assessment.StageCompleted || len(v.Rows) != 0 {
				t.Fatal("valid empty data")
			}
		} else if v.Health.Status != assessment.StageCompletedWithWarnings {
			t.Fatal("malformed rows hidden as healthy empty")
		}
	}
}

func TestCarbonDateAndScopePreflightWithoutCredentialRequests(t *testing.T) {
	for _, endpoint := range []string{"http://management.invalid", "https://user:password@management.invalid", "https://management.invalid/path", "https://management.invalid?", "https://management.invalid?x=1", "https://management.invalid#", "https://management.invalid/#", "https://management.invalid/#fragment", "https://management.invalid/%2f", "https:opaque"} {
		if _, e := NewWithHTTPClient(endpoint, posterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			t.Fatal("invalid endpoint sent credentials")
			return nil, nil, nil
		})); e == nil {
			t.Fatalf("invalid root: %s", endpoint)
		}
	}
	for _, scope := range []map[string]string{{"bad-id": "private"}, {"aaaaaaaa-1111-4111-8111-111111111111": "a", "AAAAAAAA-1111-4111-8111-111111111111": "b"}} {
		v, e := (*Scanner)(nil).Scan(context.Background(), scope, nil)
		if e == nil || v.Health.Error.Code != "carbon_scope_invalid" {
			t.Fatal("invalid/aliased scope bypassed preflight")
		}
	}
	tooMany := map[string]string{}
	for i := 0; i <= MaxSubscriptions; i++ {
		tooMany[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "fixture"
	}
	v, e := (*Scanner)(nil).Scan(context.Background(), tooMany, nil)
	if e == nil || v.Health.Error.Code != "carbon_scope_limit" {
		t.Fatal("unbounded scope")
	}
	v, e = (*Scanner)(nil).Scan(context.Background(), nil, nil)
	if e != nil || v.Health.Status != assessment.StageCompleted {
		t.Fatal("empty scope requires client or dates")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	v, e = (*Scanner)(nil).Scan(ctx, onlySub(), nil)
	if !errors.Is(e, context.Canceled) || v.Health.Error.Code != "carbon_cancelled" {
		t.Fatal("pre-cancelled scope")
	}
	for _, dates := range []string{`{}`, `{"startDate":null,"endDate":"2026-03-31"}`, `{"startDate":"2026-04-01","endDate":"2026-03-31"}`, `{"startDate":"2026-01-01","endDate":"2026-02-30"}`, `{"startDate":"2026-01-01","EndDate":"2026-03-31"}`, `{"startDate":"2026-01-01","endDate":"2026-03-31","endDate":"2026-03-31"}`} {
		calls := 0
		s, _ := NewWithHTTPClient("https://management.invalid", posterFunc(func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			return successful(dates)
		}))
		v, e := s.Scan(context.Background(), onlySub(), nil)
		if e == nil || calls != 1 || v.Health.Error.Code != "carbon_date_invalid" {
			t.Fatal("invalid dates reached report query")
		}
	}
}

func TestCarbonMalformedItemsAndLimits(t *testing.T) {
	for _, bad := range []string{`null`, `{}`, strings.ReplaceAll(simpleItem, "ItemDetailsData", "MonthlySummaryData"), strings.ReplaceAll(simpleItem, "ResourceType", "Resource"), strings.ReplaceAll(simpleItem, `"latestMonthEmissions":40`, `"latestMonthEmissions":1e999`), strings.ReplaceAll(simpleItem, `"latestMonthEmissions":40`, `"latestMonthEmissions":null`), strings.ReplaceAll(simpleItem, `"latestMonthEmissions":40`, `"LatestMonthEmissions":40`), strings.ReplaceAll(simpleItem, `"latestMonthEmissions":40`, `"latestMonthEmissions":40,"latestMonthEmissions":1`), strings.ReplaceAll(simpleItem, `"latestMonthEmissions":40`, `"latestMonthEmissions":40,"previousMonthEmissions":"1"`), strings.ReplaceAll(simpleItem, "Microsoft.Compute/virtualMachines", strings.Repeat("x", 513)), strings.ReplaceAll(simpleItem, "Microsoft.Compute/virtualMachines", `escaped\ud800`)} {
		s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
			return successful(pageJSON("["+simpleItem+","+bad+"]", allowedJSON, ""))
		})
		v, e := s.Scan(context.Background(), onlySub(), nil)
		if e != nil || len(v.Rows) != 1 || v.Rows[0].Cells[3] != "40.00" || len(v.Health.Warnings) != 1 || v.Health.Warnings[0].Message != "skipped 1 invalid carbon emission items" {
			t.Fatal("independent malformed accounting")
		}
	}
	rows := strings.TrimSuffix(strings.Repeat(simpleItem+",", 1001), ",")
	s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
		return successful(pageJSON("["+rows+"]", allowedJSON, ""))
	})
	v, e := s.Scan(context.Background(), onlySub(), nil)
	if e == nil || len(v.Rows) != 0 {
		t.Fatal("oversized item array accepted")
	}
	s = fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
		return successful(strings.Repeat("x", MaxPageBytes+1))
	})
	v, e = s.Scan(context.Background(), onlySub(), nil)
	if e == nil || len(v.Rows) != 0 {
		t.Fatal("injected client body limit bypass")
	}
}

func TestCarbonLaterBatchFailureAndConcurrentOwnedRuns(t *testing.T) {
	for _, failedBatch := range []int{1, 2} {
		calls := 0
		s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
			calls++
			if calls == failedBatch {
				return nil, nil, fmt.Errorf("provider-private-canary")
			}
			return successful(`{"value":` + capturedItems + `}`)
		})
		subs := map[string]string{}
		for i := 0; i < 101; i++ {
			subs[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "fixture"
		}
		v, e := s.Scan(context.Background(), subs, nil)
		if e == nil || calls != 2 || len(v.Rows) != 2 || v.Rows[0].Cells[3] != "100.00" || strings.Contains(e.Error(), "canary") {
			t.Fatal("valid later/prior batch or private error lost")
		}
	}
	s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
		return successful(pageJSON(capturedItems, allowedJSON, ""))
	})
	v, e := s.Scan(context.Background(), onlySub(), nil)
	if e != nil {
		t.Fatal(e)
	}
	v.Rows[0].Cells[3] = "mutated"
	v.Columns[0] = "mutated"
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, e := s.Scan(context.Background(), onlySub(), nil)
			if e != nil || v.Rows[0].Cells[3] != "100.00" || v.Columns[0] != "Period From" {
				t.Error("cross-run request/result state")
			}
			v.Rows[0].Cells[3] = "mutated"
		}()
	}
	wg.Wait()
}

func TestCarbonBoundedAuthenticatedRetriesErrorsAndClosure(t *testing.T) {
	for _, mode := range []string{"retry", "403", "oversize", "inflight-cancel", "nil", "redirect"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			var bodies []*trackedBody
			cred := &carbonCredential{scope: "https://management.fixture.invalid/.default", t: t}
			opts := azure.DefaultHTTPClientOptions(time.Second)
			opts.Scope = cred.scope
			opts.MaxRetries = 1
			opts.Transport = carbonTransportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				data := dateJSON
				status := 200
				header := http.Header{}
				if mode == "nil" {
					return nil, fmt.Errorf("private nil canary")
				}
				if mode == "inflight-cancel" {
					cancel()
					return nil, r.Context().Err()
				}
				if strings.HasSuffix(r.URL.Path, "carbonEmissionReports") {
					data = pageJSON("[]", allowedJSON, "")
				}
				if mode == "retry" && calls == 1 {
					status = 429
					data = `{"error":{"message":"private-retry-canary"}}`
					header.Set("x-ms-retry-after-ms", "1")
				}
				if mode == "403" {
					status = 403
					data = `{"error":{"code":"Denied","message":"private-denial-canary"}}`
				}
				if mode == "oversize" {
					data = strings.Repeat("x", MaxPageBytes+100)
				}
				if mode == "redirect" {
					status = 302
					data = ""
					header.Set("Location", "https://evil.invalid/private")
				}
				b := &trackedBody{Reader: bytes.NewReader([]byte(data))}
				bodies = append(bodies, b)
				return &http.Response{StatusCode: status, Header: header, Body: b, Request: r}, nil
			})
			s, e := NewWithHTTPClient("https://management.fixture.invalid", azure.NewHTTPClient(cred, opts))
			if e != nil {
				t.Fatal(e)
			}
			v, e := s.Scan(ctx, onlySub(), nil)
			checkCanonical(t, v)
			if mode == "retry" {
				if e != nil || calls != 3 || v.Health.Status != assessment.StageCompleted {
					t.Fatalf("bounded retry: %+v %v", v.Health, e)
				}
			} else if e == nil || strings.Contains(e.Error(), "canary") || v.Health.Status != assessment.StageFailed {
				t.Fatal("request failure hidden/leaked")
			}
			if mode == "inflight-cancel" && !errors.Is(e, context.Canceled) {
				t.Fatal("inflight cancellation identity")
			}
			for _, b := range bodies {
				if b.closed.Load() != 1 || b.read.Load() > MaxPageBytes+1 {
					t.Fatal("unclosed/unbounded attempt body")
				}
			}
			if mode == "redirect" && calls != 1 {
				t.Fatal("followed returned URL")
			}
		})
	}
}

func TestCarbonTotalResponseBudgetAndPartialAccess(t *testing.T) {
	// Eight full 2 MiB bodies plus date metadata exceed 16 MiB; reject body
	// eight without accepting its rows or attempting a second scope batch.
	calls := 0
	base := strings.TrimSuffix(pageJSON("["+simpleItem+"]", "null", "next"), "}") + `,"padding":"`
	padding := strings.Repeat("x", MaxPageBytes-len(base)-len(`"}`))
	s := fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
		calls++
		value := strings.Replace(base, `"next"`, fmt.Sprintf(`"page-%d"`, calls), 1)
		// Keep the exact bound despite token length changes.
		value += padding[:len(padding)-(len(value)-len(base))] + `"}`
		return successful(value)
	})
	subs := map[string]string{}
	for i := 0; i < 101; i++ {
		subs[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "fixture"
	}
	v, e := s.Scan(context.Background(), subs, nil)
	if e == nil || calls != 8 || len(v.Rows) != 1 || v.Rows[0].Cells[3] != "280.00" || v.Health.Error.Code != "carbon_response_limit" {
		t.Fatalf("total response budget/prior sums: %d %+v %v", calls, v.Health, e)
	}
	checkCanonical(t, v)
	for _, access := range []string{allowedJSON, strings.Replace(allowedJSON, `"decision":"Allowed"`, `"decision":"Denied","denialReason":"private-provider-token-canary"`, 1)} {
		s = fixtureScanner(t, func(context.Context, []byte) ([]byte, *http.Response, error) {
			return successful(pageJSON("[]", access, ""))
		})
		v, e = s.Scan(context.Background(), map[string]string{carbonSub: "A", carbonOther: "B"}, nil)
		if len(v.Rows) != 0 || len(v.Health.Warnings) == 0 || v.Health.Status == assessment.StageCompleted {
			t.Fatal("partial decision coverage certified complete empty data")
		}
		b, _ := json.Marshal(v)
		if strings.Contains(string(b), "canary") || strings.Contains(string(b), carbonSub) || strings.Contains(string(b), carbonOther) || e != nil && strings.Contains(e.Error(), "canary") {
			t.Fatal("provider access metadata leaked")
		}
	}
}

func TestCarbonCancellationBetweenBatchesDeadlineAndAfterResponse(t *testing.T) {
	for _, mode := range []string{"between", "deadline", "after-response"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if mode == "deadline" {
				var stop context.CancelFunc
				ctx, stop = context.WithTimeout(ctx, 10*time.Millisecond)
				defer stop()
			}
			calls := 0
			s := fixtureScanner(t, func(c context.Context, _ []byte) ([]byte, *http.Response, error) {
				calls++
				if mode == "deadline" {
					<-c.Done()
					return nil, nil, c.Err()
				}
				cancel()
				return successful(pageJSON("["+simpleItem+"]", "null", ""))
			})
			subs := onlySub()
			if mode == "between" {
				subs = map[string]string{}
				for i := 0; i < 101; i++ {
					subs[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "fixture"
				}
			}
			v, e := s.Scan(ctx, subs, nil)
			want := context.Canceled
			if mode == "deadline" {
				want = context.DeadlineExceeded
			}
			if !errors.Is(e, want) || calls != 1 || v.Health.Error.Code != "carbon_cancelled" {
				t.Fatal("cancellation identity/I/O boundary")
			}
			if mode != "deadline" && (len(v.Rows) != 1 || v.Rows[0].Cells[3] != "40.00") {
				t.Fatal("valid completed response discarded after cancellation")
			}
		})
	}
}
