package azure

import (
	"context"
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

type operationTransportFunc func(*http.Request) (*http.Response, error)

func (f operationTransportFunc) Do(request *http.Request) (*http.Response, error) {
	return f(request)
}

type contextResponseBody struct {
	ctx  context.Context
	read *atomic.Bool
}

func (b contextResponseBody) Read([]byte) (int, error) {
	b.read.Store(true)
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (contextResponseBody) Close() error { return nil }

func TestHTTPClientOperationTimeoutBoundsRetryWait(t *testing.T) {
	var calls atomic.Int32
	transport := operationTransportFunc(func(request *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Header:     http.Header{"Retry-After": []string{"30"}},
			Body:       io.NopCloser(strings.NewReader(`{"error":{"code":"Unavailable"}}`)),
			Request:    request,
		}, nil
	})
	options := testHTTPOptions(transport)
	options.OperationTimeout = 100 * time.Millisecond
	options.MaxRetries = 5
	client := NewHTTPClient(&mockCredential{token: "test-token"}, options)
	finished := make(chan error, 1)
	go func() {
		_, err := client.Get(context.Background(), "https://example.test")
		finished <- err
	}()
	select {
	case err := <-finished:
		if !errors.Is(err, context.DeadlineExceeded) || calls.Load() != 1 {
			t.Fatalf("error=%v calls=%d; want operation deadline during retry wait", err, calls.Load())
		}
	case <-time.After(3 * time.Second):
		t.Fatal("operation timeout did not bound retry wait")
	}
}

func TestHTTPClientOperationTimeoutBoundsResponseBody(t *testing.T) {
	for _, method := range []string{"GET", "POST", "STREAM"} {
		t.Run(method, func(t *testing.T) {
			var read atomic.Bool
			transport := operationTransportFunc(func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       contextResponseBody{ctx: request.Context(), read: &read},
					Request:    request,
				}, nil
			})
			options := testHTTPOptions(transport)
			options.OperationTimeout = 50 * time.Millisecond
			options.MaxRetries = -1
			client := NewHTTPClient(&mockCredential{token: "test-token"}, options)
			finished := make(chan error, 1)
			go func() {
				var err error
				switch method {
				case "GET":
					_, err = client.Get(context.Background(), "https://example.test")
				case "POST":
					_, _, err = client.Post(context.Background(), "https://example.test", nil)
				case "STREAM":
					var response *http.Response
					response, err = client.PostStream(context.Background(), "https://example.test", nil)
					if err == nil {
						_, err = io.ReadAll(response.Body)
						_ = response.Body.Close()
					}
				}
				finished <- err
			}()
			select {
			case err := <-finished:
				if !read.Load() || !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("body read=%v error=%v; want operation deadline", read.Load(), err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("operation timeout did not bound body read")
			}
		})
	}
}

func TestHTTPClientOperationContextPreservesCallerLimits(t *testing.T) {
	parent, stop := context.WithTimeout(context.Background(), time.Minute)
	defer stop()
	client := &HTTPClient{operationTimeout: time.Hour}
	ctx, cancel := client.operationContext(parent)
	defer cancel()
	want, _ := parent.Deadline()
	got, ok := ctx.Deadline()
	if !ok || !got.Equal(want) {
		t.Fatalf("operation extended caller deadline: %v, want %v", got, want)
	}
	stop()
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("caller cancellation was not propagated: %v", ctx.Err())
	}
	for _, timeout := range []time.Duration{0, -time.Second} {
		client.operationTimeout = timeout
		ctx, cancel := client.operationContext(context.Background())
		if _, ok := ctx.Deadline(); ok {
			t.Fatalf("non-positive timeout %v added a deadline", timeout)
		}
		cancel()
	}
}

func TestHTTPClientPostStreamKeepsContextUntilBodyRelease(t *testing.T) {
	for _, action := range []string{"eof", "close"} {
		t.Run(action, func(t *testing.T) {
			var operationCtx context.Context
			transport := operationTransportFunc(func(request *http.Request) (*http.Response, error) {
				operationCtx = request.Context()
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("ok")), Request: request}, nil
			})
			options := testHTTPOptions(transport)
			options.Timeout = 0 // Avoid a per-attempt child canceled after SDK buffering.
			client := NewHTTPClient(&mockCredential{token: "test-token"}, options)
			response, err := client.PostStream(context.Background(), "https://example.test", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if operationCtx.Err() != nil {
				t.Fatalf("stream context canceled before body release: %v", operationCtx.Err())
			}
			if action == "eof" {
				data, err := io.ReadAll(response.Body)
				if err != nil || string(data) != "ok" {
					t.Fatalf("stream body=%q error=%v", data, err)
				}
			} else if err := response.Body.Close(); err != nil {
				t.Fatal(err)
			}
			if !errors.Is(operationCtx.Err(), context.Canceled) {
				t.Fatalf("operation context retained after body release: %v", operationCtx.Err())
			}
		})
	}
}
