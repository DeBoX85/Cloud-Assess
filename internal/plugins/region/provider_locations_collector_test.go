package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

func providerCollector(t *testing.T, getter quotaGetterFunc) *ProviderLocationCollector {
	t.Helper()
	c, err := NewProviderLocationCollector("https://management.azure.com", getter)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func providerCollect(t *testing.T, body string) *ProviderLocationCollection {
	t.Helper()
	c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return quotaPage(body)
	})
	got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
	if err != nil {
		t.Fatal(err)
	}
	return got
}

const providerBasic = `{"value":[{"id":"/subscriptions/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/providers/Microsoft.Compute","namespace":"Microsoft.Compute","resourceTypes":[{"resourceType":"VirtualMachines","locations":["East US","WEST EUROPE","East US"]},{"resourceType":"virtualMachines/extensions","locations":["global"]},{"resourceType":"empty","locations":[]}]}]}`

func TestProviderCollectorLiteral(t *testing.T) {
	// Build the optional ID from a literal fixture identity, not production output.
	body := strings.ReplaceAll(providerBasic, "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", quotaTestID)
	calls := 0
	c := providerCollector(t, func(_ context.Context, endpoint string, limit int64) ([]byte, *http.Response, error) {
		calls++
		if endpoint != "https://management.azure.com/subscriptions/"+quotaTestID+"/providers?api-version=2021-04-01" || limit != 1048576 {
			t.Fatalf("provider literal request: %s/%d", endpoint, limit)
		}
		return quotaPage(body)
	})
	got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{strings.ToUpper(quotaTestID)})
	want := &ProviderLocationCollection{
		SubscriptionID: quotaTestID, Status: "complete",
		Locations: map[string]map[string]bool{
			"microsoft.compute/virtualmachines": {"eastus": true, "westeurope": true},
			"microsoft.compute/virtualmachines/extensions": {"global": true},
			"microsoft.compute/empty": {},
		},
	}
	if err != nil || calls != 1 || !reflect.DeepEqual(got, want) {
		t.Fatalf("provider literal evidence: %#v %v", got, err)
	}
	if got := providerCollect(t, `{"value":[]}`); !reflect.DeepEqual(got, &ProviderLocationCollection{SubscriptionID: quotaTestID, Status: "complete", Locations: map[string]map[string]bool{}}) {
		t.Fatalf("provider complete empty: %#v", got)
	}
	if got := providerCollect(t, `{"value":[{"namespace":"Microsoft.Compute","resourceTypes":[]},{"namespace":"microsoft.compute","resourceTypes":[{"resourceType":"widgets","locations":[]}]}]}`); got.Status != "complete" || !reflect.DeepEqual(got.Locations, map[string]map[string]bool{"microsoft.compute/widgets": {}}) {
		t.Fatalf("disjoint namespace declarations: %#v", got)
	}
}

func TestProviderCollectorAdmission(t *testing.T) {
	getter := quotaGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		t.Fatal("invalid provider request reached authentication boundary")
		return nil, nil, nil
	})
	for _, endpoint := range []string{"http://management.azure.com", "https://a..b", "https://user@management.azure.com", "https://management.azure.com/x", "https://management.azure.com?x=1", "https://management.azure.com#", "https://management.azure.com:", "https://management.azure.com:444", "https://månagement.azure.com"} {
		if _, err := NewProviderLocationCollector(endpoint, getter); err == nil {
			t.Fatal("provider origin admitted", endpoint)
		}
	}
	if _, err := NewProviderLocationCollector("https://management.azure.com", nil); err == nil {
		t.Fatal("provider nil getter")
	}
	c := providerCollector(t, getter)
	for _, id := range []string{"bad", "bbbbbbbb-bbbb-cccc-dddd-eeeeeeeeeeee", quotaTestID + "/providers"} {
		if got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{id}); err == nil || got != nil {
			t.Fatalf("provider scope admitted: %#v %v", got, err)
		}
	}
	for _, scope := range []map[string]string{{quotaTestID: ""}, {"bad": "name"}, {quotaTestID: "one", strings.ToUpper(quotaTestID): "two"}, {quotaTestID: "bad\ufdd0"}} {
		if got, err := c.Collect(context.Background(), scope, ProviderLocationRequest{quotaTestID}); err == nil || got != nil {
			t.Fatalf("provider invalid scope: %#v %v", got, err)
		}
	}
}

func TestProviderCollectorShapes(t *testing.T) {
	for _, body := range []string{
		`{}`, `null`, `[]`, `{"value":null}`, `{"value":{}}`, `{"value":[null]}`,
		`{"value":[{"namespace":"P"}]}`, `{"value":[{"namespace":"P","resourceTypes":null}]}`,
		`{"value":[{"namespace":null,"resourceTypes":[]}]}`,
		`{"value":[{"namespace":"P","resourceTypes":[null]}]}`,
		`{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"t"}]}]}`,
		`{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"t","locations":null}]}]}`,
		`{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"t","locations":[null]}]}]}`,
		`{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"t","locations":[1]}]}]}`,
		`{"value":[],"nextLink":null}`, `{"value":[],"nextLink":123}`,
		`{"value":[],"VALUE":[]}`, `{"value":[],"other":{"K":1,"\u212a":2}}`,
		`{"value":[],"x":"\ud800"}`, `{"value":[],"x":"\ufdd0"}`,
		`{"value":[]} trailing`,
	} {
		got := providerCollect(t, body)
		if got.Status != "unknown" || got.FailureCode == "" || len(got.Locations) != 0 {
			t.Fatalf("provider malformed healthy: %s %#v", body, got)
		}
	}
	for _, change := range []string{" P", "P ", "K", "P/Other", ".bad", ""} {
		body, _ := json.Marshal(map[string]any{"value": []any{map[string]any{"namespace": change, "resourceTypes": []any{}}}})
		if got := providerCollect(t, string(body)); got.Status != "unknown" {
			t.Fatalf("provider namespace admitted %q", change)
		}
	}
	for _, change := range []string{"t//child", "t/../child", "t/", " t", "t?", "ſ", ""} {
		body, _ := json.Marshal(map[string]any{"value": []any{map[string]any{"namespace": "P", "resourceTypes": []any{map[string]any{"resourceType": change, "locations": []any{}}}}}})
		if got := providerCollect(t, string(body)); got.Status != "unknown" {
			t.Fatalf("provider type admitted %q", change)
		}
	}
	for _, location := range []string{"", "East\tUS", "East\u00a0US", "europe", "unassigned", "eaſtus", "eastus/other"} {
		body, _ := json.Marshal(map[string]any{"value": []any{map[string]any{"namespace": "P", "resourceTypes": []any{map[string]any{"resourceType": "t", "locations": []string{location}}}}}})
		if got := providerCollect(t, string(body)); got.Status != "unknown" || (location != "East\tUS" && got.FailureCode != "provider_location_invalid") {
			t.Fatalf("provider location admitted %q %#v", location, got)
		}
	}
}

func TestProviderCollectorIDs(t *testing.T) {
	for _, id := range []any{nil, "", 1, "/subscriptions/" + quotaTestID + "/providers/Other", "/subscriptions/bbbbbbbb-bbbb-cccc-dddd-eeeeeeeeeeee/providers/P", "/ſubscriptions/" + quotaTestID + "/providers/P", "/subscriptions/" + quotaTestID + "/providers/K", "/subscriptions/" + quotaTestID + "/providers/P/child"} {
		body, _ := json.Marshal(map[string]any{"value": []any{map[string]any{"id": id, "namespace": "P", "resourceTypes": []any{}}}})
		if got := providerCollect(t, string(body)); got.Status != "unknown" || got.FailureCode != "provider_identity_invalid" {
			t.Fatalf("provider returned ID admitted: %v %#v", id, got)
		}
	}
	body := `{"value":[{"id":"/SUBSCRIPTIONS/` + strings.ToUpper(quotaTestID) + `/PROVIDERS/p","namespace":"P","resourceTypes":[]}]}`
	if got := providerCollect(t, body); got.Status != "complete" {
		t.Fatalf("provider valid uppercase ID rejected: %#v", got)
	}
}

const providerPrefix = `{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"one","locations":["East US"]}]}]}`

func TestProviderCollectorAtomicHealth(t *testing.T) {
	bad := `{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"new","locations":["global"]},{"resourceType":"ONE","locations":[]}]}]}`
	for _, body := range []string{
		bad,
		`{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"one","locations":[]},{"resourceType":"ONE","locations":["eastus"]}]}]}`,
	} {
		got := providerCollect(t, body)
		if got.Status != "unknown" || got.FailureCode != "provider_type_collision" || len(got.Locations) != 0 {
			t.Fatalf("provider collision not atomic: %#v", got)
		}
	}
	calls := 0
	c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 1 {
			return quotaPage(strings.TrimSuffix(providerPrefix, "}") + `,"nextLink":"?api-version=2021-04-01&$skiptoken=next"}`)
		}
		return quotaPage(bad)
	})
	got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
	want := &ProviderLocationCollection{SubscriptionID: quotaTestID, Status: "partial", FailureCode: "provider_type_collision", Locations: map[string]map[string]bool{"p/one": {"eastus": true}}}
	if err != nil || calls != 2 || !reflect.DeepEqual(got, want) {
		t.Fatalf("provider previous page lost or partial bad page escaped: %#v %v", got, err)
	}
	for _, status := range []int{204, 302, 400, 403, 404, 405, 429, 500} {
		for _, late := range []bool{false, true} {
			calls = 0
			c = providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
				calls++
				if late && calls == 1 {
					return quotaPage(`{"value":[],"nextLink":"?api-version=2021-04-01&$skiptoken=two"}`)
				}
				return []byte("secret URL token response"), &http.Response{StatusCode: status}, nil
			})
			got, err = c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
			wantStatus, wantCalls := "unknown", 1
			if late {
				wantStatus, wantCalls = "partial", 2
			}
			if err != nil || got.Status != wantStatus || got.FailureCode != "provider_request_failed" || calls != wantCalls || len(got.Locations) != 0 {
				t.Fatalf("provider status health %d: %#v %v", status, got, err)
			}
		}
	}
	c = providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(providerPrefix), &http.Response{StatusCode: 200}, errors.New("secret token URL")
	})
	got, err = c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
	if err != nil || got.Status != "unknown" || got.FailureCode != "provider_request_failed" || len(got.Locations) != 0 {
		t.Fatalf("provider getter error fabricated success: %#v %v", got, err)
	}
}

func TestProviderCollectorContinuation(t *testing.T) {
	path := "/subscriptions/" + quotaTestID + "/providers"
	for _, next := range []string{
		"https://evil.example" + path + "?api-version=2021-04-01",
		"http://management.azure.com" + path + "?api-version=2021-04-01",
		path + "?api-version=2020-01-01", path + "?api-version=2021-04-01&api-version=2021-04-01",
		path + "?API-VERSION=2021-04-01", path + "?api-version=2021-04-01&$expand=metadata",
		path + "?api-version=2021-04-01&$skiptoken=", path + "?api-version=2021-04-01#",
		strings.Replace(path, quotaTestID, "bbbbbbbb-bbbb-cccc-dddd-eeeeeeeeeeee", 1) + "?api-version=2021-04-01",
		strings.Replace(path, "providers", "%70roviders", 1) + "?api-version=2021-04-01",
		"https://user@management.azure.com" + path + "?api-version=2021-04-01",
	} {
		calls := 0
		body, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": next})
		c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return quotaPage(string(body))
		})
		got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
		if err != nil || calls != 1 || got.Status != "partial" || got.FailureCode != "provider_unsafe_continuation" {
			t.Fatalf("provider unsafe next reached getter: %q %#v %v", next, got, err)
		}
	}
	for _, next := range []string{"?api-version=2021-04-01", "https://MANAGEMENT.AZURE.COM" + path + "?api-version=2021-04-01"} {
		calls := 0
		body, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": next})
		c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { calls++; return quotaPage(string(body)) })
		got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
		if err != nil || calls != 1 || got.FailureCode != "provider_continuation_cycle" || got.Status != "partial" {
			t.Fatalf("provider canonical cycle: %#v %v", got, err)
		}
	}
	calls := 0
	c := providerCollector(t, func(_ context.Context, endpoint string, _ int64) ([]byte, *http.Response, error) {
		calls++
		if calls == 1 {
			return quotaPage(`{"value":[],"nextLink":"?api-version=2021-04-01&$skiptoken=a%2Bb%26c"}`)
		}
		if endpoint != "https://management.azure.com"+path+"?%24skiptoken=a%2Bb%26c&api-version=2021-04-01" {
			t.Fatal("provider opaque token meaning changed", endpoint)
		}
		return quotaPage(providerPrefix)
	})
	got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
	if err != nil || calls != 2 || got.Status != "complete" || !reflect.DeepEqual(got.Locations, map[string]map[string]bool{"p/one": {"eastus": true}}) {
		t.Fatalf("provider safe pages: %#v %v", got, err)
	}
}

func TestProviderCollectorBounds(t *testing.T) {
	for _, n := range []int{1048576, 1048577} {
		body := `{"value":[]}` + strings.Repeat(" ", n-len(`{"value":[]}`))
		got := providerCollect(t, body)
		if n == 1048576 && got.Status != "complete" || n == 1048577 && (got.Status != "unknown" || got.FailureCode != "provider_byte_limit") {
			t.Fatalf("provider page byte boundary%d: %#v", n, got)
		}
	}
	for _, pageCount := range []int{8, 9, 64, 65} {
		calls := 0
		c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			body := `{"value":[]}`
			if calls < pageCount {
				body = fmt.Sprintf(`{"value":[],"nextLink":"?api-version=2021-04-01&$skiptoken=%d"}`, calls)
			}
			if pageCount <= 9 {
				body += strings.Repeat(" ", 1048576-len(body))
			}
			return quotaPage(body)
		})
		got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
		want, code, wantCalls := "complete", "", pageCount
		if pageCount == 9 { want, code = "partial", "provider_byte_limit" }
		if pageCount == 65 { want, code, wantCalls = "partial", "provider_page_limit", 64 }
		if err != nil || got.Status != want || got.FailureCode != code || calls != wantCalls {
			t.Fatalf("provider aggregate/page limit%d: %#v %v calls%d", pageCount, got, err, calls)
		}
	}
	for _, extra := range []bool{false, true} {
		calls := 0
		c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			rows := strings.TrimSuffix(strings.Repeat(`{"namespace":"P","resourceTypes":[]},`, 8192), ",")
			body := `{"value":[` + rows + `]}`
			if calls < 8 || extra && calls == 8 {
				body = strings.TrimSuffix(body, "}") + fmt.Sprintf(`,"nextLink":"?api-version=2021-04-01&$skiptoken=%d"}`, calls)
			}
			if calls == 9 { body = `{"value":[{"namespace":"P","resourceTypes":[]}]}` }
			return quotaPage(body)
		})
		got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
		if err != nil || !extra && got.Status != "complete" || extra && (got.Status != "partial" || got.FailureCode != "provider_work_limit") {
			t.Fatalf("provider raw work exact/excess: %#v %v", got, err)
		}
	}
	for _, n := range []int{512, 513} {
		body := `{"value":[{"namespace":"P","resourceTypes":[{"resourceType":"t","locations":["` + strings.Repeat(" ", n-6) + `eastus"]}]}]}`
		got := providerCollect(t, body)
		if n == 512 && got.Status != "complete" || n == 513 && got.Status != "unknown" {
			t.Fatalf("provider raw label boundary%d: %#v", n, got)
		}
	}
	// The whole JSON guard applies even to metadata not projected into locations.
	for _, body := range []string{
		`{"value":[],"metadata":[` + strings.TrimSuffix(strings.Repeat("0,", 65536), ",") + `]}`,
		`{"value":[],"metadata":` + strings.Repeat("[", 33) + "0" + strings.Repeat("]", 33) + "}",
	} {
		if got := providerCollect(t, body); got.Status != "unknown" || got.FailureCode != "provider_invalid_response" {
			t.Fatalf("provider whole JSON bound ignored: %#v", got)
		}
	}
}

func TestProviderCollectorCancellationOwnership(t *testing.T) {
	c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { return quotaPage(providerPrefix) })
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
			if err != nil || got.Status != "complete" || !got.Locations["p/one"]["eastus"] { t.Error("provider concurrent state", got, err); return }
			got.Locations["p/one"]["eastus"] = false
			delete(got.Locations, "p/one")
		}()
	}
	wg.Wait()
	got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
	if err != nil || !got.Locations["p/one"]["eastus"] { t.Fatalf("provider returned alias escaped: %#v %v", got, err) }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got, err := c.Collect(ctx, quotaScope(), ProviderLocationRequest{quotaTestID}); got != nil || !errors.Is(err, context.Canceled) { t.Fatalf("provider precancel: %#v %v", got, err) }
	ctx, cancel = context.WithCancel(context.Background())
	c = providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { cancel(); return quotaPage(providerPrefix) })
	if got, err := c.Collect(ctx, quotaScope(), ProviderLocationRequest{quotaTestID}); got != nil || !errors.Is(err, context.Canceled) { t.Fatalf("provider terminal cancellation: %#v %v", got, err) }
	c = providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { return nil, nil, context.DeadlineExceeded })
	if got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID}); got != nil || !errors.Is(err, context.DeadlineExceeded) { t.Fatalf("provider deadline: %#v %v", got, err) }
}

func TestProviderCollectorAuthenticatedClient(t *testing.T) {
	for _, status := range []int{200, 302, 403, 404, 405} {
		closed, calls := false, 0
		client := azure.NewHTTPClient(quotaCredential{t, "https://management.usgovcloudapi.net/.default"}, &azure.HTTPClientOptions{
			MaxRetries: -1, Scope: "https://management.usgovcloudapi.net/.default",
			Transport: quotaTransport(func(request *http.Request) (*http.Response, error) {
				calls++
				if request.Method != "GET" || request.Body != nil || request.URL.String() != "https://management.usgovcloudapi.net/subscriptions/"+quotaTestID+"/providers?api-version=2021-04-01" || request.Header.Get("Authorization") != "Bearer synthetic-quota-token" {
					t.Fatalf("provider authenticated contract: %#v", request)
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}, "Location": {"https://evil.example"}}, Body: quotaClosedBody{strings.NewReader(`{"value":[]}`), &closed}, Request: request}, nil
			}),
		})
		c, err := NewProviderLocationCollector("https://management.usgovcloudapi.net", client)
		if err != nil { t.Fatal(err) }
		got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
		want := "unknown"
		if status == 200 { want = "complete" }
		if err != nil || got.Status != want || calls != 1 || !closed { t.Fatalf("provider bounded transport%d: %#v %v calls%d closed%v", status, got, err, calls, closed) }
	}
}

func TestProviderCollectorAvailabilityHealth(t *testing.T) {
	for _, status := range []string{"complete", "partial", "unknown"} {
		scope, inventory, request, evidence := availabilityInputs()
		counts := inventory.Subscriptions[availabilityID]
		counts.ResourceTypesByRegion["eastus"] = map[string]int{"p/one": 1, "p/two": 2}
		inventory.Subscriptions[availabilityID] = counts
		evidence.Locations = map[string]map[string]bool{"p/one": {"global": true}}
		evidence.ProvidersComplete = status == "complete"
		got, err := CalculateAvailability(context.Background(), scope, inventory, request, evidence)
		if err != nil { t.Fatal(err) }
		want := Comparison{SubscriptionID: availabilityID, SubscriptionName: "selected", SourceRegion: "eastus", TargetRegion: "westeurope", SourceResourceTypeCount: 2, AvailableTypes: 1, UnavailableTypes: 1, AvailabilityPercent: 50, MissingResourceTypes: []string{"p/two"}, MissingSKUs: []string{}, RestrictedSKUs: []string{}, ZoneRestrictedSKUs: []string{}, SourceZoneCount: 3, TargetZoneCount: 2}
		wantHealth := assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: 1}
		if status != "complete" { wantHealth.Status = assessment.StageCompletedWithWarnings; wantHealth.Warnings = []assessment.AssessmentWarning{{Code: "availability_providers_partial", Message: "region availability evidence is incomplete"}} }
		if !reflect.DeepEqual(got.Comparison, want) || !reflect.DeepEqual(got.Health, wantHealth) { t.Fatalf("provider calculation boundary%s: %#v", status, got) }
	}
}

func TestProviderCollectorDecoderBoundaries(t *testing.T) {
	for _, count := range []int{65533, 65534} {
		body := `{"value":[],"metadata":[` + strings.TrimSuffix(strings.Repeat("0,", count), ",") + `]}`
		got := providerCollect(t, body)
		if count == 65533 && got.Status != "complete" || count == 65534 && got.Status != "unknown" { t.Fatalf("provider JSON token exact/excess%d: %#v", count, got) }
	}
	for _, depth := range []int{31, 32} {
		body := `{"value":[],"metadata":` + strings.Repeat("[", depth) + "0" + strings.Repeat("]", depth) + "}"
		got := providerCollect(t, body)
		if depth == 31 && got.Status != "complete" || depth == 32 && got.Status != "unknown" { t.Fatalf("provider JSON depth exact/excess%d: %#v", depth, got) }
	}
	for _, count := range []int{127, 128} {
		body := `{"value":[]`
		for i := 0; i < count; i++ { body += fmt.Sprintf(`,"k%d":0`, i) }
		body += "}"
		got := providerCollect(t, body)
		if count == 127 && got.Status != "complete" || count == 128 && got.Status != "unknown" { t.Fatalf("provider JSON keys exact/excess%d: %#v", count, got) }
	}
	for _, count := range []int{4096, 4097} {
		calls := 0
		c := providerCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			if calls == 1 {
				body, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": "?api-version=2021-04-01&$skiptoken=" + strings.Repeat("a", count)})
				return quotaPage(string(body))
			}
			return quotaPage(`{"value":[]}`)
		})
		got, err := c.Collect(context.Background(), quotaScope(), ProviderLocationRequest{quotaTestID})
		if err != nil || count == 4096 && (got.Status != "complete" || calls != 2) || count == 4097 && (got.Status != "partial" || calls != 1 || got.FailureCode != "provider_unsafe_continuation") { t.Fatalf("provider skiptoken exact/excess%d: %#v %v calls%d", count, got, err, calls) }
	}
	// Collect's byte cap dominates this text cap. Check the decoder's budget
	// seam independently, without claiming an estate or load validation.
	body := []byte(`{"value":[{"namespace":"P","resourceTypes":[]}]}`)
	for _, prior := range []int{16711679, 16711680} {
		_, code := decodeProviderPage(context.Background(), body, quotaTestID, nil, 0, prior)
		if prior == 16711679 && code != "" || prior == 16711680 && code == "" { t.Fatalf("provider decoded-text exact/excess%d: %s", prior, code) }
	}
}
