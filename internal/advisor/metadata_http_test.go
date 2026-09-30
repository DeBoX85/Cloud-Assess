package advisor

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

type metadataCredential struct{}

func (metadataCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "synthetic-advisor-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type metadataTransport func(*http.Request) (*http.Response, error)

func (f metadataTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestMetadataHTTPRejectsUnsafePaginationBeforeSendingToken(t *testing.T) {
	for _, link := range []string{
		"https://foreign.invalid/page?sig=synthetic-secret", "http://management.azure.com/page",
		"https://management.azure.com:444/page", "https://user:secret@management.azure.com/page",
		"//foreign.invalid/page", "https://management.azure.com/page#secret", "https://[invalid/page",
	} {
		t.Run(link, func(t *testing.T) {
			calls := 0
			transport := metadataTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Method != http.MethodGet || r.Header.Get("Authorization") != "Bearer synthetic-advisor-canary" {
					t.Error("metadata request did not use authenticated GET")
				}
				body := `{"value":[]}`
				if calls == 1 {
					body = string(mustJSON(t, metadataListResult{NextLink: link}))
				}
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})
			client := azure.NewHTTPClient(metadataCredential{}, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, Scope: "https://management.azure.com/.default", Transport: transport})
			got, err := NewMetadataClient(client, "https://management.azure.com").RecommendationTypes(context.Background())
			if err == nil || got != nil || calls != 1 {
				t.Fatalf("unsafe continuation: error=%v calls=%d, want failure before second request", err, calls)
			}
			if strings.Contains(err.Error(), "synthetic-secret") || strings.Contains(err.Error(), "user:secret") {
				t.Fatal("rejection error echoed continuation credentials")
			}
		})
	}
}

func TestMetadataHTTPAllowsSameOriginPagination(t *testing.T) {
	for _, endpoint := range []string{"https://management.azure.com", "https://management.usgovcloudapi.net", "https://arm.example.invalid:8443"} {
		t.Run(endpoint, func(t *testing.T) {
			calls := 0
			transport := metadataTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				wantPath := "/providers/Microsoft.Advisor/metadata"
				body := string(mustJSON(t, metadataListResult{NextLink: endpoint + "/next-page?api-version=2020-01-01"}))
				if calls == 2 {
					wantPath = "/next-page"
					body = `{"value":[{"name":"recommendationType","properties":{"supportedValues":[{"id":"rec","displayName":"Description"}]}}]}`
				}
				if r.Method != http.MethodGet || r.URL.Path != wantPath || r.URL.Query().Get("api-version") != "2020-01-01" || r.Header.Get("Authorization") != "Bearer synthetic-advisor-canary" {
					t.Error("unexpected metadata HTTP contract")
				}
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})
			client := azure.NewHTTPClient(metadataCredential{}, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, Scope: endpoint + "/.default", Transport: transport})
			got, err := NewMetadataClient(client, endpoint).RecommendationTypes(context.Background())
			if err != nil || calls != 2 || got["rec"] != "Description" {
				t.Fatalf("pagination error=%v calls=%d result=%#v", err, calls, got)
			}
		})
	}
}
