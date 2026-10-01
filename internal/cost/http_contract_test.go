package cost

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

type costContractCredential struct{}

func (costContractCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "synthetic-cost-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type costContractTransport func(*http.Request) (*http.Response, error)

func (f costContractTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestCostAuthenticatedHTTPQueryContract(t *testing.T) {
	calls := 0
	transport := costContractTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.Scheme != "https" || r.URL.Host != "management.azure.com" || r.URL.Path != "/subscriptions/sub-fixture/providers/Microsoft.CostManagement/query" || r.URL.Query().Get("api-version") != "2021-10-01" {
			t.Error("Cost request departed from subscription query endpoint")
		}
		if r.Header.Get("Authorization") != "Bearer synthetic-cost-canary" || r.Header.Get("Content-Type") != "application/json" {
			t.Error("Cost request authentication/content type missing")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["type"] != "ActualCost" || body["timeframe"] != "Custom" {
			t.Errorf("unexpected query type/timeframe: %#v", body)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"properties":{"rows":[[12,"Storage","USD"]]}}`)), Request: r}, nil
	})
	client := azure.NewHTTPClient(costContractCredential{}, &azure.HTTPClientOptions{Timeout: 2 * time.Second, MaxRetries: -1, Scope: "https://management.azure.com/.default", Transport: transport})
	got, err := NewWithClient(client, "https://management.azure.com").scanAt(context.Background(), map[string]string{"sub-fixture": "Fixture"}, time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC))
	if err != nil || calls != 1 || len(got.Records) != 1 || len(got.Warnings) != 0 || got.Records[0].ServiceName != "Storage" {
		t.Fatalf("Cost query result: err=%v calls=%d result=%#v", err, calls, got)
	}
}
