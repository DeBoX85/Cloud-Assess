package azure

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type mockCredential struct {
	token string
	err   error
}

func (m *mockCredential) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if m.err != nil {
		return azcore.AccessToken{}, m.err
	}
	return azcore.AccessToken{Token: m.token, ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func testHTTPOptions(transport policy.Transporter) *HTTPClientOptions {
	return &HTTPClientOptions{
		Timeout:          5 * time.Second,
		MaxRetries:       1,
		OperationTimeout: 10 * time.Second,
		Scope:            "https://management.azure.com/.default",
		Transport:        transport,
	}
}

func TestHTTPClientAddsBearerAuthentication(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization = %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	client := NewHTTPClient(&mockCredential{token: "test-token"}, testHTTPOptions(server.Client()))
	body, err := client.Get(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(body) != `{"result":"ok"}` {
		t.Fatalf("body = %q", string(body))
	}
}

func TestHTTPClientReturnsAzureResponseError(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"code":"InternalError","message":"failed"}}`))
	}))
	defer server.Close()

	client := NewHTTPClient(&mockCredential{token: "test-token"}, testHTTPOptions(server.Client()))
	_, err := client.Get(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected response error")
	}
	var responseError *azcore.ResponseError
	if !errors.As(err, &responseError) {
		t.Fatalf("error type = %T, want *azcore.ResponseError", err)
	}
	if responseError.StatusCode != http.StatusInternalServerError || responseError.ErrorCode != "InternalError" {
		t.Fatalf("unexpected response error: %+v", responseError)
	}
}

func TestDefaultHTTPClientOptions(t *testing.T) {
	options := DefaultHTTPClientOptions(45 * time.Second)
	if options.Timeout != 45*time.Second || options.MaxRetries != 5 || options.OperationTimeout != 450*time.Second {
		t.Fatalf("unexpected default options: %+v", options)
	}
	if options.Scope == "" {
		t.Fatal("default ARM scope is empty")
	}
}
