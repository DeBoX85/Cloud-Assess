package aigov

import (
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

// These explicit transports never dial a socket. Tokens and origins are fake;
// no default HTTP transport or Azure session is used by request tests.
type requestTransport func(*http.Request) (*http.Response, error)

func (f requestTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

type requestCredential struct {
	calls atomic.Int64
	scope string
}

func (c *requestCredential) GetToken(ctx context.Context, o policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	if !reflect.DeepEqual(o.Scopes, []string{c.scope}) {
		return azcore.AccessToken{}, errors.New("wrong synthetic audience")
	}
	if e := ctx.Err(); e != nil {
		return azcore.AccessToken{}, e
	}
	return azcore.AccessToken{Token: "synthetic-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type closeBody struct {
	io.Reader
	closes *atomic.Int64
}

func (b closeBody) Close() error { b.closes.Add(1); return nil }
func response(r *http.Request, status int, body string, closes *atomic.Int64) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}}, Body: closeBody{strings.NewReader(body), closes}, Request: r}
}

type requestHarness struct {
	scanner                         *Scanner
	metricCredential, armCredential *requestCredential
	metrics, deployments, closed    atomic.Int64
}

func harness(t *testing.T, metric, deploy requestTransport) *requestHarness {
	t.Helper()
	h := &requestHarness{}
	h.metricCredential = &requestCredential{scope: "https://metrics.monitor.azure.com/.default"}
	h.armCredential = &requestCredential{scope: "https://arm.test/.default"}
	client := func(c *requestCredential, f requestTransport, calls *atomic.Int64) *azure.HTTPClient {
		o := azure.DefaultHTTPClientOptions(time.Second)
		o.MaxRetries = -1
		o.Scope = c.scope
		o.Transport = requestTransport(func(r *http.Request) (*http.Response, error) {
			calls.Add(1)
			if r.Header.Get("Authorization") != "Bearer synthetic-canary" {
				t.Error("missing synthetic bearer")
			}
			return f(r)
		})
		return azure.NewHTTPClient(c, o)
	}
	var e error
	h.scanner, e = NewWithClients("https://arm.test", client(h.metricCredential, metric, &h.metrics), client(h.armCredential, deploy, &h.deployments), func() time.Time { return time.Date(2026, 10, 3, 12, 34, 56, 987, time.FixedZone("fixture", 7200)) })
	if e != nil {
		t.Fatal(e)
	}
	return h
}
func locatedFixture() []LocatedAccount {
	return []LocatedAccount{{Account: fixtureAccounts()[0], Region: "WestEurope"}}
}
func hasWarning(table assessment.PluginTable, code string) bool {
	for _, w := range table.Health.Warnings {
		if w.Code == code {
			return true
		}
	}
	return false
}

func TestRequestSourceCellsAndExactContracts(t *testing.T) {
	metrics := fixtureBytes(t, "metrics-input")
	deployments := fixtureBytes(t, "deployments-input")
	var h *requestHarness
	h = harness(t, func(r *http.Request) (*http.Response, error) {
		q := r.URL.Query()
		want := map[string]string{"api-version": "2024-02-01", "metricnamespace": "Microsoft.CognitiveServices/accounts", "metricnames": "AzureOpenAIRequests", "aggregation": "Count", "interval": "PT1H", "filter": "StatusCode eq '*' and ModelDeploymentName eq '*' and ModelName eq '*'", "starttime": "2026-09-26T11:34:56Z", "endtime": "2026-10-03T10:34:56Z"}
		if r.Method != "POST" || r.URL.Scheme != "https" || r.URL.Host != "westeurope.metrics.monitor.azure.com" || r.URL.Path != "/subscriptions/"+fixtureSub+"/metrics:getBatch" || len(q) != len(want) {
			t.Error("wrong metrics destination/contract")
		}
		for k, v := range want {
			if q.Get(k) != v || len(q[k]) != 1 {
				t.Errorf("wrong metric field %s: %q", k, q[k])
			}
		}
		var body map[string][]string
		if e := json.NewDecoder(r.Body).Decode(&body); e != nil || len(body) != 1 || !reflect.DeepEqual(body["resourceids"], []string{fixtureID}) {
			t.Error("wrong metrics scope/body")
		}
		return response(r, 200, string(metrics), &h.closed), nil
	}, func(r *http.Request) (*http.Response, error) {
		if r.Method != "GET" || r.URL.Host != "arm.test" || r.URL.Path != fixtureID+"/deployments" || r.URL.RawQuery != "api-version=2025-06-01" || r.Body != nil {
			t.Error("wrong deployment contract")
		}
		return response(r, 200, string(deployments), &h.closed), nil
	})
	table, e := h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), nil)
	if e != nil || table.Health.Status != assessment.StageCompletedWithWarnings || !hasWarning(table, "ai_malformed_response") || !hasWarning(table, "ai_malformed_input") || h.metrics.Load() != 1 || h.deployments.Load() != 1 || h.closed.Load() != 2 {
		t.Fatalf("request result: %+v %v", table.Health, e)
	}
	valid(t, table)
	var capture []struct {
		Table [][]string `json:"table"`
	}
	if e := json.Unmarshal(fixtureBytes(t, "source-all"), &capture); e != nil {
		t.Fatal(e)
	}
	got := [][]string{}
	for _, r := range table.Rows {
		got = append(got, r.Cells)
	}
	if !reflect.DeepEqual(got, capture[0].Table[1:]) {
		t.Fatalf("request path differs from independently captured source cells: %v", got)
	}
	// A second result and the caller fixture must remain independently owned.
	table.Rows[0].Cells[2] = "mutated"
	again, e := h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), nil)
	if e != nil || again.Rows[0].Cells[2] != "fixture-ai" {
		t.Fatal("cross-run result mutation")
	}
}

func TestRequestEmptyPreflightAndCloud(t *testing.T) {
	var h *requestHarness
	h = harness(t, func(r *http.Request) (*http.Response, error) {
		return response(r, 200, metricWire(fixtureID, "[]"), &h.closed), nil
	}, func(*http.Request) (*http.Response, error) {
		t.Error("empty metrics enriched")
		return nil, errors.New("tripwire")
	})
	table, e := h.scanner.Scan(context.Background(), fixtureScope(), nil, nil)
	if e != nil || table.SheetName != "AI Throttling" || table.Health.Status != assessment.StageCompleted || h.metrics.Load() != 0 {
		t.Fatal("empty discovery")
	}
	table, e = h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), nil)
	if e != nil || table.SheetName != "AI Gov" || len(table.Rows) != 0 || table.Health.Status != assessment.StageCompleted || h.deployments.Load() != 0 {
		t.Fatal("empty series")
	}
	for _, a := range []LocatedAccount{{Account: fixtureAccounts()[0], Region: "https://evil.test"}, {Account: Account{ID: fixtureID + "/child"}, Region: "westus"}} {
		before := h.metrics.Load()
		table, e = h.scanner.Scan(context.Background(), fixtureScope(), []LocatedAccount{a}, nil)
		if e == nil || table.Health.Status != assessment.StageFailed || h.metrics.Load() != before {
			t.Fatal("invalid discovery reached auth")
		}
	}
	before := h.metrics.Load()
	table, e = h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), excludeFilter{})
	if e != nil || table.SheetName != "AI Throttling" || h.metrics.Load() != before {
		t.Fatal("excluded account reached auth")
	}
	for _, key := range []string{azure.EnvAzureCloud, azure.EnvAzureAuthorityHost, azure.EnvAzureResourceManagerEndpoint, azure.EnvAzureResourceManagerAudience} {
		t.Setenv(key, "")
	}
	credential := &requestCredential{}
	for _, name := range []string{"AzureGovernment", "AzureChina"} {
		t.Setenv(azure.EnvAzureCloud, name)
		if _, e := New(credential); e == nil {
			t.Fatal("unsupported cloud silently public")
		}
	}
	t.Setenv(azure.EnvAzureCloud, "AzurePublic")
	s, e := New(credential)
	if e != nil || s.origin.Host != "management.azure.com" || credential.calls.Load() != 0 {
		t.Fatal("public construction or premature auth")
	}
	t.Setenv(azure.EnvAzureAuthorityHost, "https://authority.test")
	t.Setenv(azure.EnvAzureResourceManagerEndpoint, "https://arm.test")
	t.Setenv(azure.EnvAzureResourceManagerAudience, "https://arm.test")
	if _, e := New(credential); e == nil || credential.calls.Load() != 0 {
		t.Fatal("custom cloud reached authentication")
	}
}

func TestRequestFailuresRetainMetricsAndHideProviderText(t *testing.T) {
	for _, tc := range []struct {
		name                       string
		metricStatus, deployStatus int
		metricBody, deployBody     string
		failed                     bool
		rows                       int
		warning                    string
	}{
		{"healthy", 200, 200, metricWire(fixtureID, oneSeries), `{"value":[]}`, false, 1, ""},
		{"missing coverage", 200, 200, `{"values":[]}`, `{"value":[]}`, true, 0, "ai_metrics_incomplete"},
		{"metric error", 200, 200, strings.Replace(metricWire(fixtureID, oneSeries), "Success", "Failed", 1), `{"value":[]}`, true, 0, "ai_metrics_incomplete"},
		{"metric 204", 204, 200, "", `{"value":[]}`, true, 0, "ai_metrics_incomplete"},
		{"metric denied", 403, 200, `{"error":{"code":"Denied","message":"private-provider-text"}}`, `{"value":[]}`, true, 0, "ai_metrics_incomplete"},
		{"deployment denied", 200, 403, metricWire(fixtureID, oneSeries), `{"error":{"code":"Denied","message":"private-provider-text"}}`, false, 1, "ai_deployment_requests_failed"},
		{"deployment malformed", 200, 200, metricWire(fixtureID, oneSeries), `{"value":[{"name":"fixture-deployment","sku":{"capacity":-1}}]}`, false, 1, "ai_malformed_response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var h *requestHarness
			h = harness(t, func(r *http.Request) (*http.Response, error) {
				return response(r, tc.metricStatus, tc.metricBody, &h.closed), nil
			}, func(r *http.Request) (*http.Response, error) {
				return response(r, tc.deployStatus, tc.deployBody, &h.closed), nil
			})
			table, e := h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), nil)
			valid(t, table)
			if (e != nil) != tc.failed || len(table.Rows) != tc.rows || tc.warning != "" && !hasWarning(table, tc.warning) {
				t.Fatalf("failure semantics: %+v %v", table.Health, e)
			}
			b, _ := json.Marshal(table)
			if strings.Contains(string(b), "private-provider-text") || e != nil && strings.Contains(e.Error(), "private-provider-text") {
				t.Fatal("raw provider error leaked")
			}
			if tc.rows == 1 && table.Rows[0].Cells[7] != "N/A" {
				t.Fatal("default enrichment changed")
			}
		})
	}
}

func TestDeploymentContinuationBeforeAuthentication(t *testing.T) {
	path := fixtureID + "/deployments"
	for _, next := range []string{"http://arm.test" + path + "?api-version=2025-06-01", "https://evil.test" + path + "?api-version=2025-06-01", "https://user@arm.test" + path + "?api-version=2025-06-01", "https://arm.test" + path + "?api-version=2025-06-01#", "https://arm.test" + path + "/child?api-version=2025-06-01", "https://arm.test" + strings.Replace(path, "deployments", "%64eployments", 1) + "?api-version=2025-06-01", "https://arm.test" + path + "?api-version=2024-01-01", "https://arm.test" + path + "?api-version=2025-06-01&API-VERSION=2025-06-01", "https://arm.test" + path + "?api-version=2025-06-01&skiptoken=x&skiptoken=y"} {
		t.Run(next, func(t *testing.T) {
			var h *requestHarness
			h = harness(t, func(r *http.Request) (*http.Response, error) {
				return response(r, 200, metricWire(fixtureID, oneSeries), &h.closed), nil
			}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Host != "arm.test" || r.URL.Path != path {
					t.Error("unsafe continuation reached authenticated transport")
				}
				b, _ := json.Marshal(map[string]any{"value": []any{map[string]any{"name": "fixture-deployment", "properties": map[string]any{"model": map[string]string{"version": "kept"}}}}, "nextLink": next})
				return response(r, 200, string(b), &h.closed), nil
			})
			table, e := h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), nil)
			if e != nil || len(table.Rows) != 1 || table.Rows[0].Cells[7] != "kept" || h.deployments.Load() != 1 || h.armCredential.calls.Load() != 1 || !hasWarning(table, "ai_deployment_requests_failed") {
				t.Fatalf("continuation/retention: %+v %v calls=%d", table.Health, e, h.deployments.Load())
			}
		})
	}
}

func TestDeploymentPagingCancellationAndCycle(t *testing.T) {
	for _, mode := range []string{"healthy", "denied", "cycle", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var h *requestHarness
			h = harness(t, func(r *http.Request) (*http.Response, error) {
				return response(r, 200, metricWire(fixtureID, oneSeries), &h.closed), nil
			}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Query().Get("skiptoken") == "" {
					return response(r, 200, `{"value":[{"name":"fixture-deployment","properties":{"model":{"version":"first"}}}],"nextLink":"https://arm.test`+fixtureID+`/deployments?skiptoken=two&api-version=2025-06-01"}`, &h.closed), nil
				}
				switch mode {
				case "denied":
					return response(r, 403, `{"error":{"code":"Denied"}}`, &h.closed), nil
				case "cancel":
					cancel()
					return nil, context.Canceled
				case "cycle":
					return response(r, 200, `{"value":[],"nextLink":"https://arm.test`+fixtureID+`/deployments?api-version=2025-06-01&skiptoken=two"}`, &h.closed), nil
				default:
					return response(r, 200, `{"value":[{"name":"fixture-deployment","properties":{"model":{"version":"last"}}}]}`, &h.closed), nil
				}
			})
			table, e := h.scanner.Scan(ctx, fixtureScope(), locatedFixture(), nil)
			valid(t, table)
			want := "first"
			if mode == "healthy" {
				want = "last"
			}
			if len(table.Rows) != 1 || table.Rows[0].Cells[7] != want || h.deployments.Load() != 2 {
				t.Fatalf("paging result: %+v %v", table, e)
			}
			if mode == "cancel" {
				if !errors.Is(e, context.Canceled) || table.Health.Status != assessment.StageFailed {
					t.Fatal("cancellation not retained")
				}
			} else if e != nil {
				t.Fatal(e)
			}
			if mode == "healthy" && table.Health.Status != assessment.StageCompleted || mode != "healthy" && !hasWarning(table, "ai_deployment_requests_failed") {
				t.Fatal("paging health concealed")
			}
		})
	}
}

func TestRequestGroupingOwnershipAndConcurrency(t *testing.T) {
	accounts := []LocatedAccount{}
	scope := fixtureScope()
	otherSub := "22222222-2222-4222-8222-222222222222"
	scope[otherSub] = "Other"
	for i := 0; i < 53; i++ {
		a := fixtureAccounts()[0]
		a.Name = fmt.Sprintf("account-%d", i)
		a.ID = strings.Replace(fixtureID, "fixture-ai", a.Name, 1)
		region := "westus"
		if i == 51 {
			region = "eastus"
		}
		if i == 52 {
			a.SubscriptionID = otherSub
			a.ID = strings.Replace(a.ID, fixtureSub, otherSub, 1)
		}
		accounts = append(accounts, LocatedAccount{a, region})
	}
	var mu sync.Mutex
	sizes := []int{}
	var h *requestHarness
	h = harness(t, func(r *http.Request) (*http.Response, error) {
		var b struct {
			IDs []string `json:"resourceids"`
		}
		if e := json.NewDecoder(r.Body).Decode(&b); e != nil {
			return nil, e
		}
		if len(b.IDs) > 50 {
			t.Error("oversized batch")
		}
		mu.Lock()
		sizes = append(sizes, len(b.IDs))
		mu.Unlock()
		rows := []map[string]any{}
		for _, id := range b.IDs {
			if !strings.Contains(id, r.URL.Path[:len(r.URL.Path)-len("/metrics:getBatch")]) {
				t.Error("cross-subscription batch")
			}
			rows = append(rows, map[string]any{"resourceid": id, "value": []any{}})
		}
		body, _ := json.Marshal(map[string]any{"values": rows})
		return response(r, 200, string(body), &h.closed), nil
	}, func(*http.Request) (*http.Response, error) {
		t.Error("empty accounts enriched")
		return nil, errors.New("tripwire")
	})
	before, _ := json.Marshal(accounts)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			table, e := h.scanner.Scan(context.Background(), scope, accounts, nil)
			if e != nil || table.Health.Status != assessment.StageCompleted || len(table.Rows) != 0 {
				t.Errorf("concurrent grouping: %+v %v", table.Health, e)
			}
		}()
	}
	wg.Wait()
	after, _ := json.Marshal(accounts)
	if string(before) != string(after) || len(sizes) != 12 {
		t.Fatal("caller mutated or incorrect batch count")
	}
	counts := map[int]int{}
	for _, n := range sizes {
		counts[n]++
	}
	if counts[50] != 3 || counts[1] != 9 {
		t.Fatalf("group partition %v", counts)
	}
}

func TestRequestAndResponseBudgets(t *testing.T) {
	b := &requestBudget{}
	for i := 0; i < MaxRequests; i++ {
		if e := b.before(context.Background()); e != nil {
			t.Fatal(e)
		}
	}
	if !errors.Is(b.before(context.Background()), errBudget) {
		t.Fatal("request bound missing")
	}
	b = &requestBudget{}
	raw := make([]byte, MaxPageBytes)
	for i := 0; i < MaxResponseBytes/MaxPageBytes; i++ {
		if e := b.after(raw); e != nil {
			t.Fatal(e)
		}
	}
	if !errors.Is(b.after([]byte{1}), errBudget) {
		t.Fatal("aggregate byte bound missing")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(b.before(ctx), context.Canceled) {
		t.Fatal("budget lost context")
	}
	var h *requestHarness
	h = harness(t, func(r *http.Request) (*http.Response, error) {
		return response(r, 200, metricWire(fixtureID, oneSeries), &h.closed), nil
	}, func(r *http.Request) (*http.Response, error) {
		page := h.deployments.Load()
		return response(r, 200, fmt.Sprintf(`{"value":[],"nextLink":"https://arm.test%s/deployments?api-version=2025-06-01&skiptoken=%d"}`, fixtureID, page), &h.closed), nil
	})
	table, e := h.scanner.Scan(context.Background(), fixtureScope(), locatedFixture(), nil)
	if e != nil || h.deployments.Load() != MaxPages || len(table.Rows) != 1 || !hasWarning(table, "ai_deployment_requests_failed") {
		t.Fatal("page bound/retention missing")
	}
}
