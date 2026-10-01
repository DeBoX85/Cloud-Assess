package arg

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

type argTestCredential struct{}

func (argTestCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type fakePoster struct {
	url  string
	body []byte
	resp *http.Response
	err  error
}

func (f *fakePoster) PostStream(_ context.Context, url string, body io.ReadSeekCloser) (*http.Response, error) {
	f.url = url
	if body != nil {
		f.body, _ = io.ReadAll(body)
	}
	return f.resp, f.err
}

func TestHTTPTransportSerializesRequestAndParsesResponseMetadata(t *testing.T) {
	poster := &fakePoster{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"x-ms-user-quota-remaining":    []string{"12"},
			"x-ms-user-quota-resets-after": []string{"01:02:03"},
		},
		Body: io.NopCloser(strings.NewReader(`{"data":[{"id":"one"}],"$skipToken":"next"}`)),
	}}
	transport := NewHTTPTransportWithClient(poster, "https://example.test/graph")
	top := int32(1000)
	response, err := transport.Do(context.Background(), Request{
		Subscriptions: []string{"sub"},
		Query:         "resources",
		Options:       &RequestOptions{ResultFormat: "objectArray", Top: &top},
	})
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if poster.url != "https://example.test/graph" {
		t.Fatalf("URL = %q", poster.url)
	}
	var request Request
	if err := json.Unmarshal(poster.body, &request); err != nil {
		t.Fatalf("request body is not valid JSON: %v", err)
	}
	if request.Query != "resources" || len(request.Subscriptions) != 1 || request.Subscriptions[0] != "sub" {
		t.Fatalf("unexpected serialized request: %+v", request)
	}
	if response.Quota != 12 || response.RetryAfter != time.Hour+2*time.Minute+3*time.Second {
		t.Fatalf("unexpected response metadata: %+v", response)
	}
	if response.SkipToken == nil || *response.SkipToken != "next" || len(response.Data) != 1 {
		t.Fatalf("unexpected ARG response: %+v", response)
	}
}

func TestHTTPTransportRejectsMalformedQuotaHeaders(t *testing.T) {
	poster := &fakePoster{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"x-ms-user-quota-remaining": []string{"not-a-number"}},
		Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
	}}
	transport := NewHTTPTransportWithClient(poster, "https://example.test/graph")
	if _, err := transport.Do(context.Background(), Request{}); err == nil {
		t.Fatal("expected malformed quota header error")
	}
}

func TestHTTPTransportRejectsMissingDataInsteadOfReportingEmptySuccess(t *testing.T) {
	for _, payload := range []string{`{}`, `{"data":null}`, `{"count":0}`} {
		t.Run(payload, func(t *testing.T) {
			poster := &fakePoster{resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(payload)),
			}}
			result, err := NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(
				context.Background(), "resources", map[string]string{"sub": "name"},
			)
			if err == nil || result != nil || !strings.Contains(err.Error(), "missing or null data array") {
				t.Fatalf("result = %#v, error = %v; want failed query, not empty success", result, err)
			}
		})
	}
	poster := &fakePoster{resp: &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
	}}
	result, err := NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(
		context.Background(), "resources", map[string]string{"sub": "name"},
	)
	if err != nil || result == nil || len(result.Data) != 0 {
		t.Fatalf("valid empty result = %#v, error = %v", result, err)
	}
}

func TestQueryDistinguishesThrottlingFromValidEmptyARGResponse(t *testing.T) {
	for _, test := range []struct {
		name    string
		status  int
		payload string
	}{
		{name: "throttled", status: http.StatusTooManyRequests, payload: `{"error":{"code":"TooManyRequests","message":"retry later"}}`},
		{name: "empty", status: http.StatusOK, payload: `{"data":[]}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				calls.Add(1)
				if request.Method != http.MethodPost || request.URL.Path != "/" || request.Header.Get("Authorization") != "Bearer test-token" {
					t.Errorf("unexpected ARG request method or authentication: %s", request.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.status)
				_, _ = io.WriteString(w, test.payload)
			}))
			defer server.Close()

			httpClient := azure.NewHTTPClient(argTestCredential{}, &azure.HTTPClientOptions{
				Timeout:    2 * time.Second,
				MaxRetries: -1, // The Azure SDK uses its default retry count for zero.
				Scope:      "https://management.azure.com/.default",
				Transport:  server.Client(),
			})
			result, err := NewClient(NewHTTPTransportWithClient(httpClient, server.URL)).Query(
				context.Background(), "resources", map[string]string{"sub": "name"},
			)
			if got := calls.Load(); got != 1 {
				t.Fatalf("ARG request count = %d, want 1", got)
			}
			if test.status == http.StatusOK {
				if err != nil || result == nil || len(result.Data) != 0 {
					t.Fatalf("valid empty result = %#v, error = %v", result, err)
				}
				return
			}
			var responseError *azcore.ResponseError
			if result != nil || !errors.As(err, &responseError) || responseError.StatusCode != http.StatusTooManyRequests {
				t.Fatalf("throttled result = %#v, error = %v; want Azure HTTP 429 failure", result, err)
			}
		})
	}
}

func TestParseResetAfter(t *testing.T) {
	got, err := parseResetAfter("00:05:30")
	if err != nil {
		t.Fatalf("parseResetAfter() error = %v", err)
	}
	if got != 5*time.Minute+30*time.Second {
		t.Fatalf("duration = %v", got)
	}
	if _, err := parseResetAfter("invalid"); err == nil {
		t.Fatal("expected invalid reset-after value to fail")
	}
}

func TestQueryInterruptedDuringHTTPRequestDoesNotReturnEmptySuccess(t *testing.T) {
	for _, name := range []string{"canceled", "deadline"} {
		t.Run(name, func(t *testing.T) {
			started := make(chan struct{})
			release := make(chan struct{})
			var calls atomic.Int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if calls.Add(1) == 1 {
					close(started)
				}
				select {
				case <-request.Context().Done():
				case <-release:
				}
			}))
			defer server.Close()
			defer close(release)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			httpClient := azure.NewHTTPClient(argTestCredential{}, &azure.HTTPClientOptions{
				Timeout:    10 * time.Second,
				MaxRetries: -1,
				Scope:      "https://management.azure.com/.default",
				Transport:  server.Client(),
			})
			type outcome struct {
				result *Result
				err    error
			}
			finished := make(chan outcome, 1)
			go func() {
				result, err := NewClient(NewHTTPTransportWithClient(httpClient, server.URL)).Query(ctx, "resources", map[string]string{"sub": "name"})
				finished <- outcome{result: result, err: err}
			}()
			select {
			case <-started:
			case <-ctx.Done():
				t.Fatal("request did not reach the endpoint before deadline")
			}
			want := context.DeadlineExceeded
			if name == "canceled" {
				want = context.Canceled
				cancel()
			}
			select {
			case got := <-finished:
				if got.result != nil || !errors.Is(got.err, want) {
					t.Fatalf("result = %#v, error = %v; want nil result and %v", got.result, got.err, want)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("query did not stop after context interruption")
			}
			if calls.Load() != 1 {
				t.Fatalf("requests = %d, want one interrupted request", calls.Load())
			}
		})
	}
}

func TestResourceGraphProductionEndpointIsQueryOnly(t *testing.T) {
	endpoint := resourceGraphEndpoint()
	if !strings.HasSuffix(endpoint, "/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01") {
		t.Fatalf("unexpected production query endpoint: %s", endpoint)
	}
}
