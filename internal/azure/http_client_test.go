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
	sdklog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
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

func TestHTTPClientRecoversTransientFailureWithoutDuplicatingResult(t *testing.T) {
	var calls atomic.Int32
	transport := operationTransportFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet || request.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected authenticated request: %s", request.Method)
		}
		status, payload := http.StatusServiceUnavailable, `{"error":{"code":"Unavailable"}}`
		if calls.Add(1) == 2 {
			status, payload = http.StatusOK, `{"data":[{"id":"one"}]}`
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(payload)), Request: request}, nil
	})
	client := NewHTTPClient(&mockCredential{token: "test-token"}, testHTTPOptions(transport))
	ctx := policy.WithRetryOptions(context.Background(), policy.RetryOptions{MaxRetries: 2, RetryDelay: time.Millisecond, MaxRetryDelay: time.Millisecond})
	data, err := client.Get(ctx, "https://example.test")
	if err != nil || string(data) != `{"data":[{"id":"one"}]}` || calls.Load() != 2 {
		t.Fatalf("retry recovery data=%q error=%v calls=%d", data, err, calls.Load())
	}
}

// Retry count uses production defaults. A server millisecond retry hint avoids
// changing the configured policy or spending minutes in exponential backoff.
func TestHTTPClientDefaultRetryBudgetAndNonRetriableStatus(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		recoverOn int
		wantCalls int
	}{
		{"throttling exhausted", 429, 0, 6},
		{"server failure exhausted", 503, 0, 6},
		{"last attempt succeeds", 429, 6, 6},
		{"permission denied is not retried", 403, 0, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			transport := operationTransportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				status := tc.status
				body := `{"error":{"code":"FixtureFailure","message":"synthetic"}}`
				if calls == tc.recoverOn {
					status = 200
					body = `{"ok":true}`
				}
				return &http.Response{StatusCode: status, Header: http.Header{"X-Ms-Retry-After-Ms": []string{"1"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})
			options := DefaultHTTPClientOptions(2 * time.Second)
			options.Transport = transport
			client := NewHTTPClient(&mockCredential{token: "synthetic-retry-token"}, options)
			body, err := client.Get(context.Background(), "https://management.azure.com/providers/Fixture/read")
			if calls != tc.wantCalls {
				t.Fatalf("calls=%d want %d", calls, tc.wantCalls)
			}
			if tc.recoverOn > 0 {
				if err != nil || string(body) != `{"ok":true}` {
					t.Fatalf("recovery body=%q err=%v", body, err)
				}
			} else if err == nil {
				t.Fatal("exhausted/denied retrieval represented as success")
			}
		})
	}
}

func TestHTTPClientSDKLogsAndResponseErrorDoNotExposeBearerHeader(t *testing.T) {
	const canary = "synthetic-secret-bearer-canary"
	var messages []string
	sdklog.SetEvents(sdklog.EventRequest, sdklog.EventResponse, sdklog.EventResponseError, sdklog.EventRetryPolicy)
	sdklog.SetListener(func(_ sdklog.Event, message string) { messages = append(messages, message) })
	defer func() { sdklog.SetListener(nil); sdklog.SetEvents() }()
	transport := operationTransportFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer "+canary {
			t.Error("fixture did not exercise bearer authentication")
		}
		return &http.Response{StatusCode: 403, Header: make(http.Header), Request: r, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"AuthorizationFailed","message":"resource-id-is-sensitive"}}`))}, nil
	})
	options := DefaultHTTPClientOptions(time.Second)
	options.Transport = transport
	_, err := NewHTTPClient(&mockCredential{token: canary}, options).Get(context.Background(), "https://management.azure.com/providers/Fixture/read")
	if err == nil || len(messages) == 0 {
		t.Fatal("fixture must exercise response error and enabled SDK logs")
	}
	if strings.Contains(err.Error(), canary) || strings.Contains(strings.Join(messages, "\n"), canary) {
		t.Fatal("bearer header exposed in error/SDK diagnostic output")
	}
	// Azure error content is not general anonymization: retain and document this boundary.
	if !strings.Contains(err.Error(), "resource-id-is-sensitive") {
		t.Fatal("fixture no longer demonstrates raw upstream error-message boundary")
	}
}

func TestHTTPClientDefaultTransportDoesNotFollowRedirect(t *testing.T) {
	var redirected atomic.Int32
	foreign := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirected.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer foreign.Close()
	first := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer redirect-canary" {
			t.Error("initial fixture request was not authenticated")
		}
		http.Redirect(w, r, foreign.URL+"/outside", http.StatusTemporaryRedirect)
	}))
	defer first.Close()
	// Trust only the local fixture's certificate without changing production trust.
	previous := sharedTransport
	fixtureTransport := first.Client().Transport.(*http.Transport).Clone()
	sharedTransport = fixtureTransport
	defer func() { sharedTransport = previous; fixtureTransport.CloseIdleConnections() }()
	options := DefaultHTTPClientOptions(time.Second)
	options.MaxRetries = -1
	_, err := NewHTTPClient(&mockCredential{token: "redirect-canary"}, options).Get(context.Background(), first.URL)
	var responseError *azcore.ResponseError
	if !errors.As(err, &responseError) || responseError.StatusCode != http.StatusTemporaryRedirect || redirected.Load() != 0 {
		t.Fatalf("redirect escaped request boundary: status error=%v forwarded requests=%d", err, redirected.Load())
	}
}
