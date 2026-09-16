package defender

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
)

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

func TestScanStatusMapsFiltersWarnsAndSorts(t *testing.T) {
	const (
		subA = "11111111-1111-1111-1111-111111111111"
		subB = "22222222-2222-2222-2222-222222222222"
	)
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"SubscriptionId":"` + subA + `","SubscriptionName":"Subscription A","Name":"VirtualMachines","Tier":"Standard"}`),
		json.RawMessage(`{"SubscriptionId":"` + subA + `","SubscriptionName":"Subscription A","Name":"StorageAccounts","Tier":"Free"}`),
		json.RawMessage(`{"SubscriptionId":"` + subB + `","SubscriptionName":"Subscription B","Name":"SqlServers","Tier":"Standard"}`),
		json.RawMessage(`{"SubscriptionId":`),
	}}}
	filter := fakeFilter{excludedSubscriptions: map[string]bool{strings.ToLower(subB): true}}

	result, err := NewWithClient(graph).ScanStatus(
		context.Background(),
		map[string]string{subA: "ignored by status mapping", subB: "ignored too"},
		filter,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records = %d, want 2: %+v", len(result.Records), result.Records)
	}
	if result.Records[0].Name != "StorageAccounts" || result.Records[1].Name != "VirtualMachines" {
		t.Fatalf("unexpected deterministic order: %+v", result.Records)
	}
	if result.Records[0].SubscriptionName != "Subscription A" || result.Records[0].Tier != "Free" {
		t.Fatalf("unexpected status mapping: %+v", result.Records[0])
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "defender_status_malformed_arg_rows" {
		t.Fatalf("warnings = %+v", result.Warnings)
	}
	if !strings.Contains(graph.query, "microsoft.security/pricings") || !strings.Contains(graph.query, "subscriptionName = name") {
		t.Fatalf("status query lost source semantics: %s", graph.query)
	}
}

func TestScanRecommendationsMapsDeduplicatesFiltersWarnsAndSorts(t *testing.T) {
	const (
		subscriptionID = "11111111-1111-1111-1111-111111111111"
		resourceA      = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sta"
		resourceB      = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vmb"
		excluded       = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/skipme"
	)
	rowA := json.RawMessage(`{"SubscriptionId":"` + subscriptionID + `","ResourceGroupName":"rg","ResourceType":"Microsoft.Storage","ResourceName":"sta","Category":"Data","RecommendationSeverity":"High","RecommendationName":"Encrypt data","ActionDescription":"action","RemediationDescription":"remediate","AzPortalLink":"portal.azure.com/a","ResourceId":"` + resourceA + `"}`)
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"SubscriptionId":"` + subscriptionID + `","ResourceGroupName":"rg","ResourceType":"Microsoft.Compute","ResourceName":"vmb","Category":"Compute","RecommendationSeverity":"Medium","RecommendationName":"Protect VM","ActionDescription":"protect","RemediationDescription":"enable protection","AzPortalLink":"portal.azure.com/b","ResourceId":"` + resourceB + `"}`),
		rowA,
		rowA,
		json.RawMessage(`{"SubscriptionId":"` + subscriptionID + `","ResourceGroupName":"rg","ResourceType":"Microsoft.KeyVault","ResourceName":"skipme","Category":"Secrets","RecommendationSeverity":"High","RecommendationName":"Protect vault","AzPortalLink":"portal.azure.com/c","ResourceId":"` + excluded + `"}`),
		json.RawMessage(`{"SubscriptionId":`),
	}}}
	filter := fakeFilter{excludedResources: map[string]bool{strings.ToLower(excluded): true}}

	result, err := NewWithClient(graph).ScanRecommendations(
		context.Background(),
		map[string]string{subscriptionID: "Subscription A"},
		filter,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records = %d, want 2: %+v", len(result.Records), result.Records)
	}
	if result.Records[0].ResourceID != resourceB || result.Records[1].ResourceID != resourceA {
		t.Fatalf("unexpected deterministic order: %+v", result.Records)
	}
	storage := result.Records[1]
	if storage.SubscriptionName != "Subscription A" || storage.ResourceType != "Microsoft.Storage" || storage.ResourceGroup != "rg" {
		t.Fatalf("unexpected normalized recommendation: %+v", storage)
	}
	if storage.AzurePortalLink != "https://portal.azure.com/a" {
		t.Fatalf("portal link = %q, want source-compatible https prefix", storage.AzurePortalLink)
	}
	if storage.RecommendationName != "Encrypt data" || storage.RecommendationSeverity != "High" || storage.Category != "Data" {
		t.Fatalf("unexpected recommendation mapping: %+v", storage)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "defender_recommendations_malformed_arg_rows" {
		t.Fatalf("warnings = %+v", result.Warnings)
	}
	if !strings.Contains(graph.query, "properties.status.code == 'Unhealthy'") ||
		!strings.Contains(graph.query, "mvexpand Category = properties.metadata.categories") ||
		!strings.Contains(graph.query, "ResourceType = tostring(ResourceIdsplit[6])") ||
		!strings.Contains(graph.query, "| distinct SubscriptionId") {
		t.Fatalf("recommendations query lost source semantics: %s", graph.query)
	}
}

func TestRecommendationDedupKeyIncludesCategoryAndName(t *testing.T) {
	const resourceID = "/subscriptions/sub1/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"SubscriptionId":"sub1","ResourceGroupName":"rg","ResourceType":"Microsoft.Storage","ResourceName":"st1","Category":"Data","RecommendationSeverity":"High","RecommendationName":"Protect","AzPortalLink":"portal.azure.com/a","ResourceId":"` + resourceID + `"}`),
		json.RawMessage(`{"SubscriptionId":"sub1","ResourceGroupName":"rg","ResourceType":"Microsoft.Storage","ResourceName":"st1","Category":"Network","RecommendationSeverity":"High","RecommendationName":"Protect","AzPortalLink":"portal.azure.com/a","ResourceId":"` + resourceID + `"}`),
		json.RawMessage(`{"SubscriptionId":"sub1","ResourceGroupName":"rg","ResourceType":"Microsoft.Storage","ResourceName":"st1","Category":"Data","RecommendationSeverity":"High","RecommendationName":"Protect differently","AzPortalLink":"portal.azure.com/a","ResourceId":"` + resourceID + `"}`),
	}}}

	result, err := NewWithClient(graph).ScanRecommendations(context.Background(), map[string]string{"sub1": "Sub One"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 3 {
		t.Fatalf("records = %d, want 3 distinct composite-key records: %+v", len(result.Records), result.Records)
	}
}

func TestDefenderQueryFailuresAreReturned(t *testing.T) {
	statusGraph := &fakeGraph{err: errors.New("status failed")}
	_, err := NewWithClient(statusGraph).ScanStatus(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "status failed") {
		t.Fatalf("expected status query error, got %v", err)
	}

	recommendationGraph := &fakeGraph{err: errors.New("recommendations failed")}
	_, err = NewWithClient(recommendationGraph).ScanRecommendations(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "recommendations failed") {
		t.Fatalf("expected recommendations query error, got %v", err)
	}
}

func TestDefenderNilResultsAreErrors(t *testing.T) {
	graph := &fakeGraph{result: nil}
	if _, err := NewWithClient(graph).ScanStatus(context.Background(), nil, nil); err == nil {
		t.Fatal("expected nil Defender status ARG result to fail")
	}
	if _, err := NewWithClient(graph).ScanRecommendations(context.Background(), nil, nil); err == nil {
		t.Fatal("expected nil Defender recommendations ARG result to fail")
	}
}
