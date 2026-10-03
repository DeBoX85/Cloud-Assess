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

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

const discoveryRow = `{"id":"/subscriptions/11111111-1111-4111-8111-111111111111/resourceGroups/fixture-rg/providers/Microsoft.CognitiveServices/accounts/fixture-ai","subscriptionId":"11111111-1111-4111-8111-111111111111","resourceGroup":"fixture-rg","location":"WestEurope","type":"microsoft.cognitiveservices/ACCOUNTS","name":"fixture-ai","kind":"OpenAI","sku_name":"S0","sku_tier":"Standard"}`
const literalDiscoveryQuery = `resources
| where type =~ "Microsoft.CognitiveServices/accounts"
| where isempty(kind) or kind contains "openai" or kind contains "aiservices"
| project id, subscriptionId, resourceGroup, location, type, name, sku.name, sku.tier, kind
| order by subscriptionId, resourceGroup`

type discoveryPoster func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)

func (f discoveryPoster) PostBounded(c context.Context, u string, b io.ReadSeekCloser, n int64) ([]byte, *http.Response, error) {
	return f(c, u, b, n)
}

type discoveryFilter func(string) bool

func (f discoveryFilter) IsServiceExcluded(id string) bool { return f(id) }

func discoveryWire(rows string, count, total int, token string) []byte {
	continuation := ""
	if token != "" {
		encoded, _ := json.Marshal(token)
		continuation = `,"$skipToken":` + string(encoded)
	}
	return []byte(fmt.Sprintf(`{"count":%d,"totalRecords":%d,"resultTruncated":"false","data":[%s]%s}`, count, total, rows, continuation))
}
func fixtureDiscovery(t *testing.T, poster discoveryPoster) *Discovery {
	t.Helper()
	d, err := NewDiscoveryWithClient("https://arm.test", poster)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
func runDiscoveryBody(t *testing.T, body []byte, filter Filter) (DiscoveryResult, error) {
	t.Helper()
	d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		return body, nil, nil
	})
	return d.Discover(context.Background(), fixtureScope(), filter)
}

func TestDiscoveryLiteralSourceMappingAndOwnership(t *testing.T) {
	d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		return discoveryWire(discoveryRow, 1, 1, ""), nil, nil
	})
	scope := fixtureScope()
	out, err := d.Discover(context.Background(), scope, nil)
	want := []LocatedAccount{{Account: Account{ID: fixtureID, SubscriptionID: fixtureSub, ResourceGroup: "fixture-rg", Name: "fixture-ai", Kind: "OpenAI", SKU: "S0"}, Region: "westeurope"}}
	if err != nil || out.Health.Status != assessment.StageCompleted || out.Health.Records != 1 || !reflect.DeepEqual(out.Accounts, want) {
		t.Fatalf("literal mapping: %+v %v", out, err)
	}
	out.Accounts[0].Name = "mutated"
	again, err := d.Discover(context.Background(), scope, nil)
	if err != nil || !reflect.DeepEqual(again.Accounts, want) || scope[fixtureSub] != "Synthetic subscription" {
		t.Fatal("discovery ownership changed caller or subsequent result")
	}
	for _, row := range []string{strings.ReplaceAll(strings.ReplaceAll(discoveryRow, `"OpenAI"`, `null`), `"S0"`, `null`), strings.ReplaceAll(strings.ReplaceAll(discoveryRow, `"kind":"OpenAI",`, ""), `"sku_name":"S0",`, "")} {
		out, err := runDiscoveryBody(t, discoveryWire(row, 1, 1, ""), nil)
		if err != nil || len(out.Accounts) != 1 || out.Accounts[0].Kind != "" || out.Accounts[0].SKU != "" {
			t.Fatal("missing/null source display defaults changed")
		}
	}
}

func TestDiscoveryProductionARMConfinement(t *testing.T) {
	for _, key := range []string{azure.EnvAzureCloud, azure.EnvAzureAuthorityHost, azure.EnvAzureResourceManagerEndpoint, azure.EnvAzureResourceManagerAudience} {
		t.Setenv(key, "")
	}
	credential := &requestCredential{scope: "https://management.core.windows.net/.default"}
	var calls, closes atomic.Int64
	d, err := newDiscoveryWithTransport(credential, requestTransport(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		if r.Method != "POST" || r.URL.String() != "https://management.azure.com/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || r.Header.Get("Authorization") != "Bearer synthetic-canary" {
			t.Fatal("Graph request escaped production ARM confinement")
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body) != 3 || string(body["subscriptions"]) != `["11111111-1111-4111-8111-111111111111"]` || string(body["options"]) != `{"$top":1000,"resultFormat":"objectArray"}` {
			t.Fatal("Graph scope/options widened")
		}
		var query string
		if json.Unmarshal(body["query"], &query) != nil || query != literalDiscoveryQuery {
			t.Fatal("pinned source query changed")
		}
		return response(r, 200, string(discoveryWire(discoveryRow, 1, 1, "")), &closes), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	for _, scope := range []map[string]string{nil, {"invalid": "label"}, {fixtureSub: strings.Repeat("x", 513)}} {
		out, err := d.Discover(context.Background(), scope, nil)
		if (len(scope) != 0 && (err == nil || out.Health.Status != assessment.StageFailed)) || credential.calls.Load() != 0 || calls.Load() != 0 {
			t.Fatal("invalid/empty scope authenticated")
		}
	}
	out, err := d.Discover(context.Background(), fixtureScope(), nil)
	if err != nil || len(out.Accounts) != 1 || calls.Load() != 1 || closes.Load() != 1 || credential.calls.Load() != 1 {
		t.Fatalf("production discovery/body ownership: %+v %v", out, err)
	}
}

func TestDiscoveryCloudGuardBeforeAuthentication(t *testing.T) {
	keys := []string{azure.EnvAzureAuthorityHost, azure.EnvAzureResourceManagerEndpoint, azure.EnvAzureResourceManagerAudience}
	values := []string{"https://login.microsoftonline.com/", "https://management.azure.com", "https://management.core.windows.net"}
	for _, changed := range []int{-1, 0, 1, 2} {
		t.Run(fmt.Sprintf("complete-public-or-custom-%d", changed), func(t *testing.T) {
			t.Setenv(azure.EnvAzureCloud, "AzurePublic")
			for i, key := range keys {
				value := values[i]
				if i == changed {
					value = "https://unverified.test"
				}
				t.Setenv(key, value)
			}
			c := &requestCredential{}
			d, err := NewDiscovery(c)
			if changed == -1 && (err != nil || d == nil || d.origin.Host != "management.azure.com") {
				t.Fatal("complete public discovery rejected")
			}
			if changed >= 0 && (err == nil || d != nil) {
				t.Fatal("complete custom discovery cloud accepted")
			}
			if c.calls.Load() != 0 {
				t.Fatal("discovery constructor authenticated")
			}
		})
	}
	for mask := 1; mask < 7; mask++ {
		t.Run(fmt.Sprint(mask), func(t *testing.T) {
			t.Setenv(azure.EnvAzureCloud, "public")
			for i, key := range keys {
				value := ""
				if mask&(1<<i) != 0 {
					value = values[i]
				}
				t.Setenv(key, value)
			}
			c := &requestCredential{}
			d, err := NewDiscovery(c)
			if d != nil || err == nil || c.calls.Load() != 0 {
				t.Fatal("partial discovery cloud authenticated")
			}
		})
	}
	for _, name := range []string{"AzureGovernment", "AzureChina", "custom", "unknown"} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(azure.EnvAzureCloud, name)
			c := &requestCredential{}
			if d, err := NewDiscovery(c); d != nil || err == nil || c.calls.Load() != 0 {
				t.Fatal("unverified cloud accepted")
			}
		})
	}
}

func subscriptionFixture(n int) map[string]string {
	scope := map[string]string{}
	for i := 0; i < n; i++ {
		scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "Synthetic"
	}
	return scope
}

func TestDiscoveryFalseTokenFollowsFixedBatch(t *testing.T) {
	calls := 0
	var first []string
	d := fixtureDiscovery(t, func(_ context.Context, endpoint string, body io.ReadSeekCloser, limit int64) ([]byte, *http.Response, error) {
		calls++
		var r discoveryRequest
		if json.NewDecoder(body).Decode(&r) != nil || endpoint != "https://arm.test/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || limit != 2<<20 || r.Query != literalDiscoveryQuery || r.Options.Top != 1000 || r.Options.ResultFormat != "objectArray" {
			t.Fatal("fixed discovery request changed")
		}
		switch calls {
		case 1:
			if len(r.Subscriptions) != 300 || r.Subscriptions[0] != "00000000-1111-4111-8111-111111111111" || r.Subscriptions[299] != "0000012b-1111-4111-8111-111111111111" || r.Options.SkipToken != "" {
				t.Fatal("300 subscription first batch")
			}
			first = append([]string(nil), r.Subscriptions...)
			return discoveryWire("", 0, 0, "https://foreign.test/private-token"), nil, nil
		case 2:
			if !reflect.DeepEqual(first, r.Subscriptions) || r.Options.SkipToken != "https://foreign.test/private-token" {
				t.Fatal("continuation widened batch")
			}
		case 3:
			if !reflect.DeepEqual(r.Subscriptions, []string{"0000012c-1111-4111-8111-111111111111"}) || r.Options.SkipToken != "" {
				t.Fatal("second batch membership/token")
			}
		default:
			t.Fatal("unexpected extra request")
		}
		return discoveryWire("", 0, 0, ""), nil, nil
	})
	out, err := d.Discover(context.Background(), subscriptionFixture(301), nil)
	if err != nil || out.Health.Status != assessment.StageCompleted || calls != 3 {
		t.Fatal("false-with-token stopped paging or lost second batch")
	}
}

func TestDiscoveryBatchOwnershipBeforeFilter(t *testing.T) {
	scope := subscriptionFixture(301)
	foreign := strings.ReplaceAll(discoveryRow, fixtureSub, "0000012c-1111-4111-8111-111111111111")
	calls, filters := 0, 0
	d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 1 {
			return discoveryWire(foreign, 1, 1, ""), nil, nil
		}
		return discoveryWire("", 0, 0, ""), nil, nil
	})
	out, err := d.Discover(context.Background(), scope, discoveryFilter(func(string) bool { filters++; return false }))
	if err != nil || len(out.Accounts) != 0 || filters != 0 || out.Health.Status != assessment.StageCompletedWithWarnings || len(out.Health.Warnings) != 1 {
		t.Fatal("foreign batch account reached filter or lost warning")
	}
}

func TestDiscoveryEnvelopeRejectedBeforeAccounts(t *testing.T) {
	good := string(discoveryWire(discoveryRow, 1, 1, ""))
	bad := []string{"null", good + `{}`, strings.Replace(good, `"count":1`, `"count":null`, 1), strings.Replace(good, `"count":1`, `"count":"1"`, 1), strings.Replace(good, `"count":1`, `"count":1.5`, 1), strings.Replace(good, `"count":1`, `"count":-1`, 1), strings.Replace(good, `"count":1`, `"count":2`, 1), strings.Replace(good, `"count":1`, `"Count":1`, 1), strings.Replace(good, `"count":1`, `"count":1,"COUNT":1`, 1), strings.Replace(good, `"totalRecords":1`, `"totalRecords":0`, 1), strings.Replace(good, `"resultTruncated":"false"`, `"resultTruncated":false`, 1), strings.Replace(good, `"resultTruncated":"false"`, `"resultTruncated":"true"`, 1), strings.Replace(good, `"data":[`, `"error":{},"data":[`, 1), strings.Replace(good, `"subscriptionId"`, `"SubscriptionId"`, 1), string(discoveryWire(`null`, 1, 1, "")), strings.Replace(good, `"data":[`, `"$skipToken":null,"data":[`, 1), strings.Replace(good, `"data":[`, `"$skipToken":"","data":[`, 1), strings.Replace(good, `"data":[`, `"$SkipToken":"x","data":[`, 1), strings.Replace(good, `"data":[`, `"$skipToken":"\n","data":[`, 1), strings.Replace(good, `"data":[`, `"$skipToken":"`+strings.Repeat("a", 4097)+`","data":[`, 1), good + string([]byte{0xff})}
	for i, raw := range bad {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			filters := 0
			out, err := runDiscoveryBody(t, []byte(raw), discoveryFilter(func(string) bool { filters++; return false }))
			if err == nil || out.Health.Status != assessment.StageFailed || len(out.Accounts) != 0 || filters != 0 || out.Health.Error == nil {
				t.Fatal("invalid envelope admitted accounts")
			}
		})
	}
}

func TestDiscoveryInvalidSiblingRowsAndDuplicate(t *testing.T) {
	rows := []string{discoveryRow, discoveryRow, strings.ReplaceAll(discoveryRow, fixtureSub, "22222222-2222-4222-8222-222222222222"), strings.Replace(discoveryRow, `accounts/fixture-ai"`, `accounts/fixture-ai/child"`, 1), strings.Replace(discoveryRow, `"name":"fixture-ai"`, `"name":"contradictory"`, 1), strings.Replace(discoveryRow, `"WestEurope"`, `"https://foreign.test"`, 1), strings.Replace(discoveryRow, `"OpenAI"`, `"Vision"`, 1), strings.Replace(discoveryRow, `"sku_tier":"Standard"`, `"sku_tier":42`, 1), strings.Replace(discoveryRow, `"fixture-rg"`, `".."`, 1), strings.Replace(discoveryRow, `"S0"`, `"\u001b"`, 1)}
	filters := 0
	out, err := runDiscoveryBody(t, discoveryWire(strings.Join(rows, ","), len(rows), len(rows), ""), discoveryFilter(func(string) bool { filters++; return false }))
	if err != nil || len(out.Accounts) != 1 || filters != 1 || out.Health.Status != assessment.StageCompletedWithWarnings || out.Health.Warnings[0].Message != "skipped 9 invalid or duplicate AI account rows" {
		t.Fatalf("invalid siblings: %+v %v", out, err)
	}
}

func TestDiscoveryLaterFailureRetainsPrefix(t *testing.T) {
	for _, failure := range []string{"denied", "metadata", "total", "cycle", "cancel", "body"} {
		t.Run(failure, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
				calls++
				if calls == 1 {
					return discoveryWire(discoveryRow, 1, 2, "private-token"), nil, nil
				}
				switch failure {
				case "denied":
					return nil, nil, errors.New("private-provider-text")
				case "metadata":
					return []byte(`{"data":[]}`), nil, nil
				case "total":
					return discoveryWire("", 0, 3, ""), nil, nil
				case "cycle":
					return discoveryWire("", 0, 2, "private-token"), nil, nil
				case "cancel":
					cancel()
					return nil, nil, ctx.Err()
				default:
					return []byte(strings.Repeat("x", MaxPageBytes+1)), nil, nil
				}
			})
			out, err := d.Discover(ctx, fixtureScope(), nil)
			if err == nil || out.Health.Status != assessment.StageFailed || out.Health.Records != 1 || len(out.Accounts) != 1 || out.Accounts[0].ID != fixtureID || calls != 2 {
				t.Fatal("later failure discarded valid discovery prefix")
			}
			encoded, _ := json.Marshal(out.Health)
			if strings.Contains(string(encoded)+err.Error(), "private-") {
				t.Fatal("provider/token data leaked through health/error")
			}
			if failure == "cancel" && !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation identity lost")
			}
		})
	}
}

func TestDiscoveryActualAssessmentFilters(t *testing.T) {
	for _, name := range []string{"unknown-include", "known-include", "known-exclude", "unknown-exclude", "resource", "subscription", "rg", "scanner"} {
		t.Run(name, func(t *testing.T) {
			f := config.NewFilters()
			switch name {
			case "unknown-include", "known-include", "known-exclude":
				f.Assessment.Include.Tags["team"] = "synthetic"
			case "unknown-exclude":
				f.Assessment.Exclude.Tags["team"] = "synthetic"
			case "resource":
				f.Assessment.Exclude.Resources = []string{strings.ToUpper(fixtureID)}
			case "subscription":
				f.Assessment.Exclude.Subscriptions = []string{fixtureSub}
			case "rg":
				f.Assessment.Exclude.ResourceGroups = []string{"/subscriptions/" + fixtureSub + "/resourceGroups/fixture-rg"}
			}
			f.RebuildIndexes()
			f.Assessment.SetAllowedResourceTypes([]string{"Microsoft.CognitiveServices/accounts"})
			if name == "scanner" {
				f.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Storage/storageAccounts"})
			}
			if name == "known-include" {
				f.Assessment.SetResourceScope(fixtureID, true)
			}
			if name == "known-exclude" {
				f.Assessment.SetResourceScope(fixtureID, false)
			}
			out, err := runDiscoveryBody(t, discoveryWire(discoveryRow, 1, 1, ""), f.Assessment)
			want := 0
			if name == "known-include" || name == "unknown-exclude" {
				want = 1
			}
			if err != nil || len(out.Accounts) != want || out.Health.Status != assessment.StageCompleted {
				t.Fatalf("recorded filter %s: %+v %v", name, out, err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		return discoveryWire(discoveryRow, 1, 1, ""), nil, nil
	})
	out, err := d.Discover(ctx, fixtureScope(), discoveryFilter(func(string) bool { cancel(); return false }))
	if !errors.Is(err, context.Canceled) || len(out.Accounts) != 0 {
		t.Fatal("filter cancellation admitted account")
	}
}

func TestDiscoveryConcurrentRuns(t *testing.T) {
	d := fixtureDiscovery(t, func(_ context.Context, _ string, b io.ReadSeekCloser, _ int64) ([]byte, *http.Response, error) {
		var request discoveryRequest
		if json.NewDecoder(b).Decode(&request) != nil {
			t.Fatal("request JSON")
		}
		row := strings.ReplaceAll(discoveryRow, fixtureSub, request.Subscriptions[0])
		return discoveryWire(row, 1, 1, ""), nil, nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Go(func() {
			id := fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)
			out, err := d.Discover(context.Background(), map[string]string{id: "synthetic"}, nil)
			if err != nil || len(out.Accounts) != 1 || out.Accounts[0].SubscriptionID != id {
				t.Error("disjoint run state leaked")
			}
		})
	}
	wg.Wait()
}

func TestDiscoveryScopeAndWorkBudgets(t *testing.T) {
	for _, n := range []int{4096, 4097} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			return discoveryWire("", 0, 0, ""), nil, nil
		})
		out, err := d.Discover(context.Background(), subscriptionFixture(n), nil)
		if n == 4096 && (err != nil || calls != 14 || out.Health.Status != assessment.StageCompleted) {
			t.Fatal("allowed scope boundary")
		}
		if n == 4097 && (err == nil || calls != 0 || out.Health.Error.Code != "ai_discovery_scope_limit") {
			t.Fatal("excess scope authenticated")
		}
	}
	for _, scope := range []map[string]string{{"aaaaaaaa-1111-4111-8111-111111111111": "one", "AAAAAAAA-1111-4111-8111-111111111111": "two"}, {fixtureSub: strings.Repeat("a", 513)}} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			return nil, nil, nil
		})
		if _, err := d.Discover(context.Background(), scope, nil); err == nil || calls != 0 {
			t.Fatal("invalid owned scope authenticated")
		}
	}
	for _, n := range []int{32, 33} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			token := ""
			if calls < n {
				token = fmt.Sprint(calls)
			}
			return discoveryWire("", 0, 0, token), nil, nil
		})
		out, err := d.Discover(context.Background(), fixtureScope(), nil)
		if n == 32 && (err != nil || calls != 32) {
			t.Fatal("allowed page boundary")
		}
		if n == 33 && (err == nil || calls != 32 || out.Health.Error.Code != "ai_discovery_request_limit") {
			t.Fatal("excess page request sent")
		}
	}
	for _, subscriptions := range []int{600, 601} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			token := ""
			if calls%32 != 0 {
				token = fmt.Sprint(calls)
			}
			return discoveryWire("", 0, 0, token), nil, nil
		})
		out, err := d.Discover(context.Background(), subscriptionFixture(subscriptions), nil)
		if subscriptions == 600 && (err != nil || calls != 64) {
			t.Fatal("allowed total request boundary")
		}
		if subscriptions == 601 && (err == nil || calls != 64 || out.Health.Error.Code != "ai_discovery_request_limit") {
			t.Fatal("excess total request sent")
		}
	}
	for _, n := range []int{4096, 4097} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			count := 1000
			if calls == 5 {
				count = n - 4000
			}
			token := ""
			if calls < 5 {
				token = fmt.Sprint(calls)
			}
			rows := strings.TrimSuffix(strings.Repeat(discoveryRow+",", count), ",")
			return discoveryWire(rows, count, n, token), nil, nil
		})
		out, err := d.Discover(context.Background(), fixtureScope(), excludeFilter{})
		if len(out.Accounts) != 0 || calls != 5 {
			t.Fatal("excluded-row fixture did not finish boundary requests")
		}
		if n == 4096 && err != nil {
			t.Fatal("allowed raw row boundary")
		}
		if n == 4097 && (err == nil || out.Health.Error.Code != "ai_discovery_row_limit") {
			t.Fatal("excluded rows bypassed raw budget")
		}
	}
}

func TestDiscoverySuccessfulBodyBudgetAndTwoBatchPrefix(t *testing.T) {
	for _, pages := range []int{8, 9} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			token := ""
			if calls < pages {
				token = fmt.Sprint(calls)
			}
			rows := ""
			count := 0
			if calls == 1 {
				rows = discoveryRow
				count = 1
			}
			body := discoveryWire(rows, count, 1, token)
			body = append(body, []byte(strings.Repeat(" ", (2<<20)-len(body)))...)
			return body, nil, nil
		})
		out, err := d.Discover(context.Background(), fixtureScope(), nil)
		if len(out.Accounts) != 1 || out.Accounts[0].ID != fixtureID {
			t.Fatal("body budget discarded accepted prefix")
		}
		if pages == 8 && (err != nil || calls != 8) {
			t.Fatal("allowed 16 MiB successful body boundary")
		}
		if pages == 9 && (err == nil || calls != 9 || out.Health.Error.Code != "ai_discovery_body_limit") {
			t.Fatal("successful byte aggregate bypassed")
		}
	}
	for _, n := range []int{2 << 20, (2 << 20) + 1} {
		base := discoveryWire("", 0, 0, "")
		body := append(base, []byte(strings.Repeat(" ", n-len(base)))...)
		out, err := runDiscoveryBody(t, body, nil)
		if n == 2<<20 && err != nil {
			t.Fatal("allowed per-attempt body boundary")
		}
		if n > (2<<20) && (err == nil || out.Health.Error.Code != "ai_discovery_body_limit") {
			t.Fatal("trusted oversized body bypassed adapter check")
		}
	}
	calls := 0
	d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 1 {
			return discoveryWire(strings.ReplaceAll(discoveryRow, fixtureSub, "00000000-1111-4111-8111-111111111111"), 1, 1, ""), nil, nil
		}
		return nil, nil, errors.New("denied second batch")
	})
	out, err := d.Discover(context.Background(), subscriptionFixture(301), nil)
	if err == nil || calls != 2 || len(out.Accounts) != 1 || out.Accounts[0].SubscriptionID != "00000000-1111-4111-8111-111111111111" || out.Health.Status != assessment.StageFailed {
		t.Fatal("two-batch failure discarded accepted first batch")
	}
}

func TestDiscoveryTextAndDeadlineBudgets(t *testing.T) {
	scope := subscriptionFixture(4096)
	for id := range scope {
		scope[id] = strings.Repeat("a", 512)
	}
	row := strings.ReplaceAll(discoveryRow, fixtureSub, "00000000-1111-4111-8111-111111111111")
	row = strings.ReplaceAll(row, "fixture-rg", strings.Repeat("r", 512))
	row = strings.ReplaceAll(row, "fixture-ai", strings.Repeat("n", 512))
	row = strings.ReplaceAll(row, "OpenAI", "OpenAI"+strings.Repeat("k", 506))
	row = strings.ReplaceAll(row, "Standard", strings.Repeat("t", 512))
	row = strings.ReplaceAll(row, "S0", strings.Repeat("s", 512))
	calls := 0
	d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
		calls++
		rows := strings.TrimSuffix(strings.Repeat(row+",", 350), ",")
		return discoveryWire(rows, 350, 4096, fmt.Sprint(calls)), nil, nil
	})
	out, err := d.Discover(context.Background(), scope, nil)
	if err == nil || out.Health.Error.Code != "ai_discovery_text_limit" || len(out.Accounts) != 1 {
		t.Fatalf("scope/projected/ignored tier text not charged: code=%+v calls=%d error=%v", out.Health.Error, calls, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	before := calls
	out, err = d.Discover(ctx, fixtureScope(), nil)
	if !errors.Is(err, context.Canceled) || calls != before || out.Health.Status != assessment.StageFailed {
		t.Fatal("canceled discovery made request")
	}
	d = fixtureDiscovery(t, func(ctx context.Context, _ string, _ io.ReadSeekCloser, _ int64) ([]byte, *http.Response, error) {
		deadline, ok := ctx.Deadline()
		if !ok || deadline.Sub(time.Now()) > 5*time.Minute {
			t.Error("missing five-minute discovery deadline")
		}
		return discoveryWire("", 0, 0, ""), nil, nil
	})
	if _, err := d.Discover(context.Background(), fixtureScope(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoveryRejectedProjectedTextCountsAtExactBoundary(t *testing.T) {
	// The wire fixture is independently sized from its literal strings. All tier
	// values exceed the label limit and are rejected, but their decoded bytes
	// still count. Scope labels and body continuations count once each.
	scope := subscriptionFixture(4096)
	for id := range scope {
		scope[id] = strings.Repeat("a", 512)
	}
	var values map[string]string
	if json.Unmarshal([]byte(discoveryRow), &values) != nil {
		t.Fatal("literal row")
	}
	baseText := 0
	for key, value := range values {
		if key != "sku_tier" {
			baseText += len(value)
		}
	}
	tokenText := 0
	for i := 1; i < 14; i++ {
		tokenText += len(fmt.Sprint(i))
	}
	lastTier := (16 << 20) - 4096*(36+512) - 14*baseText - 13*(1<<20) - tokenText
	if lastTier <= 512 || lastTier > 1<<20 {
		t.Fatal("independent fixture sizing")
	}
	for _, excess := range []int{0, 1} {
		calls := 0
		d := fixtureDiscovery(t, func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error) {
			calls++
			if calls > 14 {
				return discoveryWire("", 0, 0, ""), nil, nil
			}
			tierSize := 1 << 20
			token := fmt.Sprint(calls)
			if calls == 14 {
				tierSize = lastTier + excess
				token = ""
			}
			row := strings.Replace(discoveryRow, `"Standard"`, `"`+strings.Repeat("x", tierSize)+`"`, 1)
			return discoveryWire(row, 1, 14, token), nil, nil
		})
		out, err := d.Discover(context.Background(), scope, nil)
		if excess == 0 && (err != nil || calls != 27 || out.Health.Status != assessment.StageCompletedWithWarnings || out.Health.Warnings[0].Message != "skipped 14 invalid or duplicate AI account rows") {
			t.Fatalf("allowed decoded text boundary: calls%d %+v %v", calls, out.Health, err)
		}
		if excess == 1 && (err == nil || calls != 14 || out.Health.Error.Code != "ai_discovery_text_limit") {
			t.Fatal("rejected projected text bypassed one-byte excess budget")
		}
	}
	rows := strings.TrimSuffix(strings.Repeat(discoveryRow+",", 1001), ",")
	out, err := runDiscoveryBody(t, discoveryWire(rows, 1001, 1001, ""), nil)
	if err == nil || out.Health.Error.Code != "ai_discovery_page_invalid" || len(out.Accounts) != 0 {
		t.Fatal("excess object-array page admitted accounts")
	}
}

func FuzzDiscoveryEnvelope(f *testing.F) {
	for _, body := range [][]byte{discoveryWire(discoveryRow, 1, 1, ""), discoveryWire("", 0, 0, "opaque"), []byte(`{"count":0,"data":null}`)} {
		f.Add(body)
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > MaxPageBytes+1 {
			t.Skip()
		}
		page, err := decodeDiscoveryPage(context.Background(), raw)
		if err == nil && (len(page.rows) > 1000 || page.total < int64(len(page.rows)) || len(page.token) > 4096) {
			t.Fatal("accepted envelope exceeds literal limits")
		}
	})
}
