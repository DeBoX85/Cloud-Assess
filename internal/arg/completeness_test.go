package arg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestExplicitTokenlessTruncationFailsWithoutPartialResult(t *testing.T) {
	for _, previous := range []bool{false, true} {
		t.Run(fmt.Sprint(previous), func(t *testing.T) {
			responses := []*Response{}
			if previous {
				responses = append(responses, &Response{Data: []json.RawMessage{json.RawMessage(`{"id":"earlier"}`)}, SkipToken: strptr("next")})
			}
			// Exercise the actual response decoder, not only a synthetic Response field.
			poster := &fakePoster{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"data":[{"id":"last"}],"resultTruncated":"true","count":1,"totalRecords":42}`))}}
			final, err := NewHTTPTransportWithClient(poster, "https://example.test/graph").Do(context.Background(), Request{})
			if err != nil {
				t.Fatal(err)
			}
			responses = append(responses, final)
			transport := &fakeTransport{responses: responses}
			result, err := NewClient(transport).Query(context.Background(), "resources", map[string]string{"sub": "fixture"})
			if err == nil || !strings.Contains(err.Error(), "truncated without a continuation token") || result != nil || len(transport.requests) != len(responses) {
				t.Fatalf("result=%v error=%v calls=%d", result, err, len(transport.requests))
			}
		})
	}
}

type contractPagerPoster struct {
	calls    int
	requests []Request
}

func (f *contractPagerPoster) PostStream(_ context.Context, _ string, body io.ReadSeekCloser) (*http.Response, error) {
	var request Request
	if err := json.NewDecoder(body).Decode(&request); err != nil {
		return nil, err
	}
	f.calls++
	f.requests = append(f.requests, request)
	if f.calls > 4 {
		return nil, fmt.Errorf("fixture runaway tripwire")
	}
	if request.Options == nil || request.Options.Top == nil || *request.Options.Top < 1 || *request.Options.Top > 1000 {
		return nil, fmt.Errorf("documented ARG page-size maximum exceeded")
	}
	offset := 0
	if request.Options.SkipToken != nil {
		var err error
		offset, err = strconv.Atoi(*request.Options.SkipToken)
		if err != nil {
			return nil, err
		}
	}
	if offset < 0 || offset > 2501 {
		return nil, fmt.Errorf("fixture offset invalid")
	}
	end := min(offset+int(*request.Options.Top), 2501)
	rows := []map[string]string{}
	for i := offset; i < end; i++ {
		rows = append(rows, map[string]string{"id": fmt.Sprintf("fixture-%04d", i)})
	}
	payload := map[string]any{"data": rows, "count": len(rows), "totalRecords": 2501, "resultTruncated": "false"}
	if end < 2501 {
		payload["$skipToken"] = strconv.Itoa(end)
		payload["resultTruncated"] = "true"
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(encoded)))}, nil
}
func TestQueryUsesDocumentedPageSizeAndCollectsEveryPage(t *testing.T) {
	poster := &contractPagerPoster{}
	result, err := NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(context.Background(), "resources", map[string]string{"sub": "fixture"})
	if err != nil || result == nil || len(result.Data) != 2501 || poster.calls != 3 {
		t.Fatalf("error=%v result=%v calls=%d", err, result, poster.calls)
	}
	for i, row := range result.Data {
		var decoded map[string]string
		if err := json.Unmarshal(row, &decoded); err != nil {
			t.Fatal(err)
		}
		if decoded["id"] != fmt.Sprintf("fixture-%04d", i) {
			t.Fatalf("row %d lost or duplicated", i)
		}
	}
	for _, request := range poster.requests {
		if request.Query != "resources" || len(request.Subscriptions) != 1 || request.Subscriptions[0] != "sub" || *request.Options.Top != 1000 {
			t.Fatal("query/scope/page size changed")
		}
	}
}

func TestQueryTruncationMetadataContractAndEmptyCompatibility(t *testing.T) {
	for _, test := range []struct {
		name, payload string
		wantError     bool
	}{
		{"complete-empty", `{"data":[],"resultTruncated":"false"}`, false},
		{"legacy-empty", `{"data":[]}`, false},
		{"null-legacy", `{"data":[],"resultTruncated":null}`, false},
		{"deliberate-kql-limit", `{"data":[{"id":"limited"}],"resultTruncated":"false","count":1,"totalRecords":1}`, false},
		{"truncated-empty", `{"data":[],"resultTruncated":"true"}`, true},
		{"empty-token", `{"data":[],"resultTruncated":"true","$skipToken":""}`, true},
		{"invalid-enum", `{"data":[],"resultTruncated":"secret-invalid"}`, true},
		{"wrong-json-type", `{"data":[],"resultTruncated":true}`, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			poster := &fakePoster{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(test.payload))}}
			query := "resources"
			if test.name == "deliberate-kql-limit" {
				query += " | limit 1"
			}
			result, err := NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(context.Background(), query, map[string]string{"sub": "fixture"})
			if test.wantError {
				if err == nil || result != nil {
					t.Fatalf("result=%v error=%v", result, err)
				}
				if strings.Contains(err.Error(), "secret-invalid") {
					t.Fatal("raw metadata leaked")
				}
			} else if err != nil || result == nil {
				t.Fatalf("result=%v error=%v", result, err)
			}
			if !test.wantError {
				wantRows := 0
				if test.name == "deliberate-kql-limit" {
					wantRows = 1
				}
				if len(result.Data) != wantRows {
					t.Fatalf("rows=%d want=%d", len(result.Data), wantRows)
				}
				var request Request
				if err := json.Unmarshal(poster.body, &request); err != nil || request.Query != query {
					t.Fatal("query was rewritten")
				}
			}

		})
	}
}
