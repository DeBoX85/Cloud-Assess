package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type fakeBatchClient struct {
	mu      sync.Mutex
	calls   []armBatchRequest
	handler func(armBatchRequest) (*http.Response, error)
}

func (f *fakeBatchClient) PostStream(_ context.Context, _ string, body io.ReadSeekCloser) (*http.Response, error) {
	defer body.Close()
	payload, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var request armBatchRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return nil, err
	}
	f.mu.Lock()
	f.calls = append(f.calls, request)
	f.mu.Unlock()
	return f.handler(request)
}

func jsonResponse(t *testing.T, value any) *http.Response {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(string(payload))),
		Header:     make(http.Header),
	}
}

func TestRecommendationCatalogMatchesReferenceContracts(t *testing.T) {
	if got := len(Recommendations()); got != 41 {
		t.Fatalf("dedicated diagnostics recommendations = %d, want 41", got)
	}
	if !Supports("Microsoft.Storage/storageAccounts") {
		t.Fatal("storage accounts should support diagnostic settings")
	}
	if !Supports("Microsoft.Compute/virtualMachines") {
		t.Fatal("virtual machines should remain in the reference support table")
	}
	if Supports("Microsoft.Compute/snapshots") {
		t.Fatal("snapshots should not be in the diagnostics support table")
	}

	storage, ok := RecommendationFor("Microsoft.Storage/storageAccounts")
	if !ok {
		t.Fatal("storage recommendation not found")
	}
	if storage.ID != "st-001" || storage.Impact != assessment.ImpactLow || storage.Category != CategoryMonitoringAndAlerting {
		t.Fatalf("unexpected storage recommendation: %+v", storage)
	}
	if storage.Source != Source || storage.ValidationMechanism != ValidationAzureResourceManager {
		t.Fatalf("unexpected storage provenance: source=%q validation=%q", storage.Source, storage.ValidationMechanism)
	}
	if _, ok := RecommendationFor("Microsoft.Compute/virtualMachines"); ok {
		t.Fatal("legacy VM support should not create a dedicated recommendation")
	}
}

func TestScanMatchesStorageDiagnosticSettingsBehavior(t *testing.T) {
	const (
		subscriptionID = "11111111-1111-1111-1111-111111111111"
		withDiagID     = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/withdiag"
		withoutDiagID  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/withoutdiag"
		vmID           = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm1"
		snapshotID     = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Compute/snapshots/snap1"
	)

	client := &fakeBatchClient{}
	client.handler = func(request armBatchRequest) (*http.Response, error) {
		responses := make([]armBatchResponseItem, 0, len(request.Requests))
		for _, item := range request.Requests {
			content := json.RawMessage(`{"value":[]}`)
			if strings.Contains(strings.ToLower(item.RelativeURL), "/withdiag/") {
				content = json.RawMessage(`{"value":[{"id":"` + withDiagID + `/providers/microsoft.insights/diagnosticSettings/default"}]}`)
			}
			responses = append(responses, armBatchResponseItem{HTTPStatusCode: http.StatusOK, Content: content})
		}
		return jsonResponse(t, armBatchResponse{Responses: responses}), nil
	}

	scanner := NewWithClient(client, "https://management.azure.com")
	result, err := scanner.Scan(context.Background(), []assessment.Resource{
		{ID: withDiagID, SubscriptionID: subscriptionID, ResourceGroup: "rg", Type: "Microsoft.Storage/storageAccounts", Name: "withdiag"},
		{ID: withoutDiagID, SubscriptionID: subscriptionID, ResourceGroup: "rg", Type: "Microsoft.Storage/storageAccounts", Name: "withoutdiag"},
		{ID: vmID, SubscriptionID: subscriptionID, ResourceGroup: "rg", Type: "Microsoft.Compute/virtualMachines", Name: "vm1"},
		{ID: snapshotID, SubscriptionID: subscriptionID, ResourceGroup: "rg", Type: "Microsoft.Compute/snapshots", Name: "snap1"},
	}, nil, map[string]string{subscriptionID: "Test Subscription"})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Findings) != 1 {
		t.Fatalf("findings = %d, want 1: %+v", len(result.Findings), result.Findings)
	}
	finding := result.Findings[0]
	if finding.RecommendationID != "st-001" || finding.ResourceID != withoutDiagID {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.SubscriptionName != "Test Subscription" || finding.ValidationMechanism != ValidationAzureResourceManager {
		t.Fatalf("unexpected normalized finding metadata: %+v", finding)
	}

	client.mu.Lock()
	defer client.mu.Unlock()
	if len(client.calls) != 1 {
		t.Fatalf("batch calls = %d, want 1", len(client.calls))
	}
	if got := len(client.calls[0].Requests); got != 3 {
		t.Fatalf("diagnostics batch resources = %d, want 3 supported resources", got)
	}
}

type denyFilter struct{ resourceID string }

func (f denyFilter) IsServiceExcluded(resourceID string) bool {
	return strings.EqualFold(resourceID, f.resourceID)
}

func TestScanHonorsDownstreamResourceFilter(t *testing.T) {
	const resourceID = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/withoutdiag"
	client := &fakeBatchClient{}
	client.handler = func(request armBatchRequest) (*http.Response, error) {
		responses := make([]armBatchResponseItem, len(request.Requests))
		for i := range responses {
			responses[i] = armBatchResponseItem{HTTPStatusCode: http.StatusOK, Content: json.RawMessage(`{"value":[]}`)}
		}
		return jsonResponse(t, armBatchResponse{Responses: responses}), nil
	}

	result, err := NewWithClient(client, "https://management.azure.com").Scan(
		context.Background(),
		[]assessment.Resource{{ID: resourceID, Type: "Microsoft.Storage/storageAccounts", Name: "withoutdiag"}},
		denyFilter{resourceID: resourceID},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("excluded resource produced findings: %+v", result.Findings)
	}
}

func TestDiagnosticsBatchesAtTwentyResources(t *testing.T) {
	client := &fakeBatchClient{}
	client.handler = func(request armBatchRequest) (*http.Response, error) {
		responses := make([]armBatchResponseItem, len(request.Requests))
		for i := range responses {
			responses[i] = armBatchResponseItem{HTTPStatusCode: http.StatusOK, Content: json.RawMessage(`{"value":[]}`)}
		}
		return jsonResponse(t, armBatchResponse{Responses: responses}), nil
	}

	resources := make([]assessment.Resource, 21)
	for i := range resources {
		resources[i] = assessment.Resource{
			ID:   "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st" + string(rune('a'+i)),
			Type: "Microsoft.Storage/storageAccounts",
			Name: "storage",
		}
	}

	result, err := NewWithClient(client, "https://management.azure.com").Scan(context.Background(), resources, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 21 {
		t.Fatalf("findings = %d, want 21", len(result.Findings))
	}

	client.mu.Lock()
	sizes := make([]int, 0, len(client.calls))
	for _, call := range client.calls {
		sizes = append(sizes, len(call.Requests))
	}
	client.mu.Unlock()
	sort.Ints(sizes)
	if len(sizes) != 2 || sizes[0] != 1 || sizes[1] != 20 {
		t.Fatalf("batch sizes = %v, want [1 20]", sizes)
	}
}

func TestBatchFailureReturnsErrorInsteadOfTerminatingProcess(t *testing.T) {
	client := &fakeBatchClient{handler: func(armBatchRequest) (*http.Response, error) {
		return nil, errors.New("boom")
	}}
	_, err := NewWithClient(client, "https://management.azure.com").Scan(
		context.Background(),
		[]assessment.Resource{{
			ID:   "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1",
			Type: "Microsoft.Storage/storageAccounts",
		}},
		nil,
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected propagated batch error, got %v", err)
	}
}

func TestNonSuccessSubrequestPreservesReferenceFindingAndAddsWarning(t *testing.T) {
	const resourceID = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	client := &fakeBatchClient{handler: func(request armBatchRequest) (*http.Response, error) {
		return jsonResponse(t, armBatchResponse{Responses: []armBatchResponseItem{{
			HTTPStatusCode: http.StatusForbidden,
			Content:        json.RawMessage(`{"error":{"code":"Forbidden"}}`),
		}}}), nil
	}}

	result, err := NewWithClient(client, "https://management.azure.com").Scan(
		context.Background(),
		[]assessment.Resource{{ID: resourceID, Type: "Microsoft.Storage/storageAccounts", Name: "st1"}},
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 || result.Findings[0].RecommendationID != "st-001" {
		t.Fatalf("reference-compatible finding not produced: %+v", result.Findings)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "diagnostics_subrequest_non_success" {
		t.Fatalf("expected explicit uncertainty warning, got %+v", result.Warnings)
	}
}

func TestMalformedDiagnosticSettingIDBecomesWarning(t *testing.T) {
	const resourceID = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	client := &fakeBatchClient{handler: func(request armBatchRequest) (*http.Response, error) {
		return jsonResponse(t, armBatchResponse{Responses: []armBatchResponseItem{{
			HTTPStatusCode: http.StatusOK,
			Content:        json.RawMessage(`{"value":[{"id":"not-a-diagnostic-setting-id"}]}`),
		}}}), nil
	}}

	result, err := NewWithClient(client, "https://management.azure.com").Scan(
		context.Background(),
		[]assessment.Resource{{ID: resourceID, Type: "Microsoft.Storage/storageAccounts", Name: "st1"}},
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "diagnostics_malformed_setting_id" {
		t.Fatalf("expected malformed-setting warning, got %+v", result.Warnings)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("malformed setting ID should not mark the resource compliant: %+v", result.Findings)
	}
}
