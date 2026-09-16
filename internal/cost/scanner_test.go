package cost

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type fakePoster struct {
	mu      sync.Mutex
	calls   []postCall
	handler func(string, []byte) ([]byte, *http.Response, error)
}

type postCall struct {
	url  string
	body []byte
}

func (f *fakePoster) Post(_ context.Context, url string, body io.ReadSeekCloser) ([]byte, *http.Response, error) {
	defer body.Close()
	payload, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}
	f.mu.Lock()
	f.calls = append(f.calls, postCall{url: url, body: append([]byte(nil), payload...)})
	f.mu.Unlock()
	return f.handler(url, payload)
}

func TestScanAtPreservesCostQueryAndRowMapping(t *testing.T) {
	const subscriptionID = "11111111-1111-1111-1111-111111111111"
	poster := &fakePoster{}
	poster.handler = func(url string, body []byte) ([]byte, *http.Response, error) {
		if want := "https://management.azure.com/subscriptions/" + subscriptionID + "/providers/Microsoft.CostManagement/query?api-version=2021-10-01"; url != want {
			t.Fatalf("url = %q, want %q", url, want)
		}

		var request queryDefinition
		if err := json.Unmarshal(body, &request); err != nil {
			t.Fatal(err)
		}
		if request.Type != "ActualCost" || request.Timeframe != "Custom" {
			t.Fatalf("unexpected cost query type/timeframe: %+v", request)
		}
		if request.Dataset.Aggregation["TotalCost"].Name != "Cost" || request.Dataset.Aggregation["TotalCost"].Function != "Sum" {
			t.Fatalf("unexpected aggregation: %+v", request.Dataset.Aggregation)
		}
		if len(request.Dataset.Grouping) != 1 || request.Dataset.Grouping[0].Type != "Dimension" || request.Dataset.Grouping[0].Name != "ServiceName" {
			t.Fatalf("unexpected grouping: %+v", request.Dataset.Grouping)
		}

		return []byte(`{"properties":{"rows":[[12.34,"Storage","USD"],[5,"Compute","USD"]]}}`), &http.Response{StatusCode: http.StatusOK}, nil
	}

	now := time.Date(2026, 9, 16, 15, 0, 0, 0, time.FixedZone("CEST", 2*60*60))
	result, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{subscriptionID: "Production"},
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records = %d, want 2: %+v", len(result.Records), result.Records)
	}

	wantFrom := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	wantTo := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
	for _, record := range result.Records {
		if !record.From.Equal(wantFrom) || !record.To.Equal(wantTo) {
			t.Fatalf("unexpected cost period: %s - %s", record.From, record.To)
		}
		if record.SubscriptionID != subscriptionID || record.SubscriptionName != "Production" {
			t.Fatalf("unexpected subscription mapping: %+v", record)
		}
	}
	if result.Records[0].ServiceName != "Compute" || result.Records[0].Value != "5" || result.Records[0].Currency != "USD" {
		t.Fatalf("unexpected first row mapping/order: %+v", result.Records[0])
	}
	if result.Records[1].ServiceName != "Storage" || result.Records[1].Value != "12.34" {
		t.Fatalf("unexpected second row mapping/order: %+v", result.Records[1])
	}
}

func TestScanAtPopulatesSubscriptionNameUnlikePinnedStageBug(t *testing.T) {
	poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
		return []byte(`{"properties":{"rows":[[1,"Service","USD"]]}}`), &http.Response{StatusCode: http.StatusOK}, nil
	}}
	result, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{"sub-1": "Named Subscription"},
		time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 || result.Records[0].SubscriptionName != "Named Subscription" {
		t.Fatalf("subscription display name should be retained: %+v", result.Records)
	}
}

func TestSkippableAzureErrorWarnsAndOtherSubscriptionsContinue(t *testing.T) {
	poster := &fakePoster{}
	poster.handler = func(url string, _ []byte) ([]byte, *http.Response, error) {
		if strings.Contains(url, "/subscriptions/a/") {
			return nil, nil, &azcore.ResponseError{ErrorCode: "NotFound", StatusCode: http.StatusNotFound}
		}
		return []byte(`{"properties":{"rows":[[2,"Storage","USD"]]}}`), &http.Response{StatusCode: http.StatusOK}, nil
	}

	result, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{"a": "A", "b": "B"},
		time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 || result.Records[0].SubscriptionID != "b" {
		t.Fatalf("healthy subscription should still return costs: %+v", result.Records)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "cost_subscription_skipped" || !strings.Contains(result.Warnings[0].Message, "NotFound") {
		t.Fatalf("unexpected warnings: %+v", result.Warnings)
	}
}

func TestNonSkippableAzureErrorFailsStageContract(t *testing.T) {
	poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
		return nil, nil, &azcore.ResponseError{ErrorCode: "Forbidden", StatusCode: http.StatusForbidden}
	}}
	_, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{"a": "A"},
		time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC),
	)
	if err == nil || !strings.Contains(err.Error(), "Forbidden") {
		t.Fatalf("expected non-skippable Cost Management failure, got %v", err)
	}
}

func TestMalformedCostRowWarnsInsteadOfPanicking(t *testing.T) {
	poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
		return []byte(`{"properties":{"rows":[[1,"Storage"],[3,"Compute","EUR"]]}}`), &http.Response{StatusCode: http.StatusOK}, nil
	}}
	result, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{"a": "A"},
		time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 || result.Records[0].ServiceName != "Compute" {
		t.Fatalf("valid cost row should survive: %+v", result.Records)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "cost_malformed_row" {
		t.Fatalf("expected malformed-row warning, got %+v", result.Warnings)
	}
}

func TestEmptyCostResponseProducesEmptyDataset(t *testing.T) {
	poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
		return nil, &http.Response{StatusCode: http.StatusNoContent}, nil
	}}
	result, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{"a": "A"},
		time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 0 || len(result.Warnings) != 0 {
		t.Fatalf("empty response should produce an empty dataset: %+v", result)
	}
}

func TestTransportErrorIsReturned(t *testing.T) {
	poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
		return nil, nil, errors.New("network failed")
	}}
	_, err := NewWithClient(poster, "https://management.azure.com").scanAt(
		context.Background(),
		map[string]string{"a": "A"},
		time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC),
	)
	if err == nil || !strings.Contains(err.Error(), "network failed") {
		t.Fatalf("expected transport failure, got %v", err)
	}
}

func TestDefaultWorkerCeilingMatchesSource(t *testing.T) {
	scanner := NewWithClient(&fakePoster{}, "https://management.azure.com")
	if scanner.maxWorkers != 2 {
		t.Fatalf("Cost worker ceiling = %d, want 2", scanner.maxWorkers)
	}
}
