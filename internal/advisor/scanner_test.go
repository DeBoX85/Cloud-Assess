package advisor

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
)

type fakeGetter struct {
	urls    []string
	pages   map[string][]byte
	errors  map[string]error
}

func (f *fakeGetter) Get(_ context.Context, url string) ([]byte, error) {
	f.urls = append(f.urls, url)
	if err := f.errors[url]; err != nil {
		return nil, err
	}
	body, ok := f.pages[url]
	if !ok {
		return nil, errors.New("unexpected URL: " + url)
	}
	return body, nil
}

func TestMetadataRecommendationTypesFollowsPagination(t *testing.T) {
	const endpoint = "https://management.azure.com"
	first := endpoint + "/providers/Microsoft.Advisor/metadata?api-version=2020-01-01"
	second := endpoint + "/next-page?api-version=2020-01-01"
	getter := &fakeGetter{pages: map[string][]byte{
		first: []byte(`{"value":[{"name":"category","properties":{"supportedValues":[{"id":"ignore","displayName":"Ignore"}]}},{"name":"recommendationType","properties":{"supportedValues":[{"id":"rec-1","displayName":"First recommendation"}]}}],"nextLink":"/next-page?api-version=2020-01-01"}`),
		second: []byte(`{"value":[{"name":"recommendationType","properties":{"supportedValues":[{"id":"rec-2","displayName":"Second recommendation"}]}}]}`),
	}}

	got, err := NewMetadataClient(getter, endpoint).RecommendationTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"rec-1": "First recommendation", "rec-2": "Second recommendation"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recommendation types = %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(getter.urls, []string{first, second}) {
		t.Fatalf("metadata URLs = %#v", getter.urls)
	}
}

type fakeMetadata struct {
	values map[string]string
	err    error
	calls  int
}

func (f *fakeMetadata) RecommendationTypes(context.Context) (map[string]string, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return f.values, nil
}

type fakeGraph struct {
	result        *arg.Result
	err           error
	calls         int
	query         string
	subscriptions map[string]string
}

func (f *fakeGraph) Query(_ context.Context, query string, subscriptions map[string]string, _ ...arg.QueryOptions) (*arg.Result, error) {
	f.calls++
	f.query = query
	f.subscriptions = subscriptions
	if f.err != nil {
		return nil, f.err
	}
	return f.result, nil
}

type fakeFilter struct {
	excludedSubscriptions map[string]bool
	excludedResources     map[string]bool
}

func (f fakeFilter) IsSubscriptionExcluded(subscriptionID string) bool {
	return f.excludedSubscriptions[strings.ToLower(subscriptionID)]
}

func (f fakeFilter) IsServiceExcluded(resourceID string) bool {
	return f.excludedResources[strings.ToLower(resourceID)]
}

func TestScannerNormalizesFiltersWarnsAndSorts(t *testing.T) {
	const (
		subA = "11111111-1111-1111-1111-111111111111"
		subB = "22222222-2222-2222-2222-222222222222"
		resA = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sta"
		resB = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vmb"
		resExcluded = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/skipme"
		resSubB = "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/subb"
	)
	rows := []json.RawMessage{
		json.RawMessage(`{"SubscriptionId":"` + subA + `","ResourceId":"` + resB + `","ImpactedValue":"vmb","Category":"Cost","Impact":"Medium","RecommendationTypeId":"rec-b"}`),
		json.RawMessage(`{"SubscriptionId":"` + subA + `","ResourceId":"` + resA + `","ImpactedValue":"sta","Category":"Reliability","Impact":"High","RecommendationTypeId":"rec-a"}`),
		json.RawMessage(`{"SubscriptionId":"` + subA + `","ResourceId":"` + resExcluded + `","ImpactedValue":"skipme","Category":"Security","Impact":"High","RecommendationTypeId":"rec-c"}`),
		json.RawMessage(`{"SubscriptionId":"` + subB + `","ResourceId":"` + resSubB + `","ImpactedValue":"subb","Category":"Cost","Impact":"Low","RecommendationTypeId":"rec-d"}`),
		json.RawMessage(`{"SubscriptionId":`),
	}
	graph := &fakeGraph{result: &arg.Result{Data: rows}}
	metadata := &fakeMetadata{values: map[string]string{
		"rec-a": "Recommendation A",
		"rec-b": "Recommendation B",
	}}
	filter := fakeFilter{
		excludedSubscriptions: map[string]bool{strings.ToLower(subB): true},
		excludedResources:     map[string]bool{strings.ToLower(resExcluded): true},
	}

	result, err := NewWithClients(graph, metadata).Scan(
		context.Background(),
		map[string]string{subA: "Subscription A", subB: "Subscription B"},
		filter,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records = %d, want 2: %+v", len(result.Records), result.Records)
	}
	if result.Records[0].RecommendationID != "rec-b" || result.Records[1].RecommendationID != "rec-a" {
		t.Fatalf("unexpected deterministic order: %+v", result.Records)
	}
	first := result.Records[0]
	if first.SubscriptionName != "Subscription A" || first.ResourceType != "Microsoft.Compute/virtualMachines" || first.ResourceName != "vmb" || first.Description != "Recommendation B" {
		t.Fatalf("unexpected normalized Advisor record: %+v", first)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "advisor_malformed_arg_rows" {
		t.Fatalf("warnings = %+v", result.Warnings)
	}
	if !strings.Contains(graph.query, "properties.suppressionIds") || !strings.Contains(graph.query, "summarize take_any(*) by ResourceId, RecommendationTypeId") {
		t.Fatalf("Advisor query lost source filtering/dedup semantics: %s", graph.query)
	}
}

func TestScannerPreservesEmptyDescriptionWhenMetadataMissing(t *testing.T) {
	const resourceID = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"SubscriptionId":"11111111-1111-1111-1111-111111111111","ResourceId":"` + resourceID + `","ImpactedValue":"st1","Category":"Cost","Impact":"Low","RecommendationTypeId":"unknown"}`),
	}}}
	result, err := NewWithClients(graph, &fakeMetadata{values: map[string]string{}}).Scan(context.Background(), map[string]string{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 || result.Records[0].Description != "" {
		t.Fatalf("missing metadata should preserve empty description: %+v", result.Records)
	}
}

func TestMetadataFailureStopsBeforeARGQuery(t *testing.T) {
	graph := &fakeGraph{result: &arg.Result{}}
	metadata := &fakeMetadata{err: errors.New("metadata failed")}
	_, err := NewWithClients(graph, metadata).Scan(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "metadata failed") {
		t.Fatalf("expected metadata error, got %v", err)
	}
	if graph.calls != 0 {
		t.Fatalf("ARG calls = %d, want 0 after metadata failure", graph.calls)
	}
}

func TestARGFailureIsReturned(t *testing.T) {
	graph := &fakeGraph{err: errors.New("ARG failed")}
	_, err := NewWithClients(graph, &fakeMetadata{values: map[string]string{}}).Scan(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "ARG failed") {
		t.Fatalf("expected ARG error, got %v", err)
	}
}
