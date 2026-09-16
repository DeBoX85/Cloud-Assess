package arg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
	top := int32(5000)
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
