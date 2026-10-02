package azure

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestBoundedPOSTAttemptBodiesHeadersAndClosure(t *testing.T) {
	for _, tc := range []struct {
		status  int
		body    string
		limit   int64
		wantErr bool
	}{
		{200, "12345678", 8, false}, {200, strings.Repeat("x", 10000), 8, true}, {403, strings.Repeat("x", 10000), 8, true},
	} {
		var read, closes atomic.Int64
		opts := testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(r.Body)
			if err != nil || string(body) != "read-only-query" || r.Method != "POST" || r.Header.Get("Authorization") != "Bearer canary" || r.Header.Get("Content-Type") != "application/json" {
				t.Fatal("bounded POST request")
			}
			return &http.Response{StatusCode: tc.status, Header: http.Header{"X-Synthetic": {"header"}}, Body: countedBody{strings.NewReader(tc.body), &read, &closes}, Request: r}, nil
		}))
		opts.MaxRetries = -1
		body, r, err := NewHTTPClient(&mockCredential{token: "canary"}, opts).PostBounded(context.Background(), "https://example.test", NopReadSeekCloser{Reader: bytes.NewReader([]byte("read-only-query"))}, tc.limit)
		if (err != nil) != tc.wantErr || read.Load() > 9 || closes.Load() != 1 {
			t.Fatalf("bounded POST: err=%v read=%d closes=%d", err, read.Load(), closes.Load())
		}
		if !tc.wantErr && (string(body) != tc.body || r.Header.Get("X-Synthetic") != "header") {
			t.Fatal("response body/headers")
		}
	}
	var read, closes, calls atomic.Int64
	opts := testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": {"30"}}, Body: countedBody{strings.NewReader(strings.Repeat("x", 10000)), &read, &closes}, Request: r}, nil
	}))
	opts.OperationTimeout = 2 * time.Second
	_, _, err := NewHTTPClient(&mockCredential{token: "canary"}, opts).PostBounded(context.Background(), "https://example.test", NopReadSeekCloser{Reader: bytes.NewReader([]byte("query"))}, 8)
	if err == nil || !strings.Contains(err.Error(), "response exceeds byte limit") || calls.Load() != 1 || read.Load() > 9 || closes.Load() != 1 {
		t.Fatalf("retry bound: %v read=%d closes=%d calls=%d", err, read.Load(), closes.Load(), calls.Load())
	}

	read.Store(0)
	closes.Store(0)
	calls.Store(0)
	opts = testHTTPOptions(operationTransportFunc(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": {"30"}}, Body: countedBody{strings.NewReader("small"), &read, &closes}, Request: r}, nil
	}))
	opts.OperationTimeout = 50 * time.Millisecond
	_, _, err = NewHTTPClient(&mockCredential{token: "canary"}, opts).PostBounded(context.Background(), "https://example.test", NopReadSeekCloser{Reader: bytes.NewReader([]byte("query"))}, 8)
	if !errors.Is(err, context.DeadlineExceeded) || calls.Load() != 1 || read.Load() > 9 || closes.Load() != 1 {
		t.Fatalf("small-body retry deadline: %v read=%d closes=%d calls=%d", err, read.Load(), closes.Load(), calls.Load())
	}
	client := NewHTTPClient(&mockCredential{token: "canary"}, testHTTPOptions(operationTransportFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("invalid limit reached transport")
		return nil, nil
	})))
	for _, limit := range []int64{0, -1, 1<<63 - 1} {
		if _, _, err := client.PostBounded(context.Background(), "https://example.test", nil, limit); err == nil {
			t.Fatal("invalid POST limit accepted")
		}
	}
}
