package azure

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type countedBody struct {
	io.Reader
	read, closes *atomic.Int64
}

func (b countedBody) Read(p []byte) (int, error) {
	n, e := b.Reader.Read(p)
	b.read.Add(int64(n))
	return n, e
}
func (b countedBody) Close() error { b.closes.Add(1); return nil }

func TestBoundedGETSuccessErrorsAndClosure(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		limit      int64
		wantError  bool
	}{
		{"exact", "12345678", 200, 8, false},
		{"oversized success", strings.Repeat("x", 10000), 200, 8, true},
		{"oversized error", strings.Repeat("x", 10000), 403, 8, true},
		{"Azure error", `{"error":{"code":"Denied"}}`, 403, 128, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var read, closes atomic.Int64
			options := testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != "GET" || r.Header.Get("Authorization") != "Bearer canary" {
					t.Fatal("request contract")
				}
				return &http.Response{StatusCode: tc.status, Header: http.Header{"Content-Type": {"application/json"}}, Body: countedBody{strings.NewReader(tc.body), &read, &closes}, Request: r}, nil
			}))
			options.MaxRetries = -1
			body, err := NewHTTPClient(&mockCredential{token: "canary"}, options).GetBounded(context.Background(), "https://example.test", tc.limit)
			if (err != nil) != tc.wantError || read.Load() > tc.limit+1 || closes.Load() != 1 {
				t.Fatalf("error=%v bytes=%d closes=%d", err, read.Load(), closes.Load())
			}
			if !tc.wantError && string(body) != tc.body {
				t.Fatal("body mismatch")
			}
			if tc.name == "Azure error" {
				var azerr *azcore.ResponseError
				if !errors.As(err, &azerr) || azerr.ErrorCode != "Denied" {
					t.Fatalf("response classification: %v", err)
				}
			}
		})
	}
}

func TestBoundedGETRetryBodyAndOperationDeadline(t *testing.T) {
	var read, closes, calls atomic.Int64
	options := testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": {"30"}}, Body: countedBody{strings.NewReader(strings.Repeat("x", 10000)), &read, &closes}, Request: r}, nil
	}))
	options.OperationTimeout = 50 * time.Millisecond
	_, err := NewHTTPClient(&mockCredential{token: "canary"}, options).GetBounded(context.Background(), "https://example.test", 8)
	if !errors.Is(err, context.DeadlineExceeded) || calls.Load() != 1 || read.Load() > 9 || closes.Load() != 1 {
		t.Fatalf("retry bound: err=%v calls=%d bytes=%d closes=%d", err, calls.Load(), read.Load(), closes.Load())
	}

	var bodyRead atomic.Bool
	options = testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: contextResponseBody{r.Context(), &bodyRead}, Request: r}, nil
	}))
	options.OperationTimeout = 50 * time.Millisecond
	_, err = NewHTTPClient(&mockCredential{token: "canary"}, options).GetBounded(context.Background(), "https://example.test", 8)
	if !errors.Is(err, context.DeadlineExceeded) || !bodyRead.Load() {
		t.Fatalf("body deadline: %v", err)
	}
}

func TestBoundedGETInvalidLimitBeforeTransport(t *testing.T) {
	options := testHTTPOptions(operationTransportFunc(func(*http.Request) (*http.Response, error) { t.Fatal("transport invoked"); return nil, nil }))
	client := NewHTTPClient(&mockCredential{token: "canary"}, options)
	for _, limit := range []int64{0, -1, 1<<63 - 1} {
		if _, err := client.GetBounded(context.Background(), "https://example.test", limit); err == nil {
			t.Fatal("invalid limit accepted")
		}
	}
}

func TestBoundedGETExposesStatusAndClosedOwnedBody(t *testing.T) {
	for _, status := range []int{200, 202, 204} {
		var read, closes atomic.Int64
		options := testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: status, Header: http.Header{"X-Fixture": {"kept"}}, Body: countedBody{strings.NewReader("owned"), &read, &closes}, Request: r}, nil
		}))
		options.MaxRetries = -1
		b, r, e := NewHTTPClient(&mockCredential{token: "canary"}, options).GetBoundedWithResponse(context.Background(), "https://example.test", 10)
		if e != nil || r == nil || r.StatusCode != status || r.Header.Get("X-Fixture") != "kept" || string(b) != "owned" || closes.Load() != 1 {
			t.Fatalf("status/body closure %v %v", r, e)
		}
	}
}
