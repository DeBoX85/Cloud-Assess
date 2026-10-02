package sqleol

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

const subA = "11111111-1111-4111-8111-111111111111"
const subB = "22222222-2222-4222-8222-222222222222"

type transportFunc func(context.Context, arg.Request) (*arg.Response, error)

func (f transportFunc) Do(c context.Context, r arg.Request) (*arg.Response, error) { return f(c, r) }

type filterFunc func(string) bool

func (f filterFunc) IsSubscriptionExcluded(id string) bool { return f(id) }
func selected() map[string]string {
	return map[string]string{subA: "discovered name is not output", subB: "B"}
}
func sourceRows(t *testing.T) []json.RawMessage {
	t.Helper()
	b, err := os.ReadFile("testdata/source-input.json")
	var rows []json.RawMessage
	if err != nil || fmt.Sprintf("%x", sha256.Sum256(b)) != "2e82c9d3d106ab20236c95a76f0b7c1d692eba439e67b3107da99d15e83821f0" || json.Unmarshal(b, &rows) != nil || len(rows) != 2 {
		t.Fatal("source fixture input")
	}
	return rows
}

func TestActualSourceEveryColumnMetadataFilterEmptyAndOrder(t *testing.T) {
	if fmt.Sprintf("%x", sha256.Sum256([]byte(Query))) != "f027cb07f082913cf5e83a7d3d77d337ceaa479b46393d3a55a61a2adbe893db" {
		t.Fatal("pinned source query bytes changed")
	}
	rows := sourceRows(t)
	scanner := NewWithTransport(transportFunc(func(_ context.Context, request arg.Request) (*arg.Response, error) {
		if request.Query != Query || !reflect.DeepEqual(request.Subscriptions, []string{subA, subB}) || request.Options == nil || request.Options.Top == nil || *request.Options.Top != 1000 || request.Options.AuthorizationScopeFilter != nil {
			t.Fatal("selected query contract")
		}
		return &arg.Response{Data: rows}, nil
	}))
	for _, mode := range []string{"output", "filtered", "empty"} {
		t.Run(mode, func(t *testing.T) {
			hash := map[string]string{"output": "881d489015558566b780eb53173b9db55aae992dafee6092c3005dd33cf4078d", "filtered": "149a153b603fed8c41e0d37ba3dde2f8cd194915f962c97188a031e62a8c9622", "empty": "b2561841d95d8ecca57018c2f8d92bc7758c6c5046de0688de4984dc30e0e82f"}[mode]
			b, err := os.ReadFile("testdata/source-" + mode + ".json")
			if err != nil || fmt.Sprintf("%x", sha256.Sum256(b)) != hash {
				t.Fatal("source capture bytes changed")
			}
			var capture []struct {
				Metadata struct {
					Name, Version, Description, Author, License string
					Type                                        int
					ColumnMetadata                              []struct{ Name string }
				}
				SheetName   string `json:"sheet_name"`
				Description string
				Table       [][]string
			}
			if json.Unmarshal(b, &capture) != nil || len(capture) != 1 {
				t.Fatal("source capture shape")
			}
			subs := selected()
			var filter Filter
			if mode == "filtered" {
				filter = filterFunc(func(id string) bool { return id != subA })
			}
			if mode == "empty" {
				subs = map[string]string{}
			}
			got, err := scanner.Scan(context.Background(), subs, filter)
			want := capture[0]
			m := want.Metadata
			if err != nil || got.Health.Status != assessment.StageCompleted || got.Metadata != (assessment.PluginMetadata{Name: m.Name, Version: m.Version, Description: m.Description, Author: m.Author, License: m.License, Type: "internal"}) || m.Type != 1 || got.SheetName != want.SheetName || got.Description != want.Description || !reflect.DeepEqual(got.Columns, want.Table[0]) || len(got.Columns) != 32 || len(got.Rows) != len(want.Table)-1 {
				t.Fatalf("complete source metadata/table: %+v %v", got, err)
			}
			for i, row := range got.Rows {
				if !reflect.DeepEqual(row.Cells, want.Table[i+1]) {
					t.Fatalf("all source cells/order row %d: %+v", i, row)
				}
			}
			for i, column := range m.ColumnMetadata {
				if got.Columns[i] != column.Name {
					t.Fatal("source declared headers")
				}
			}
			if mode == "output" && (got.Rows[0].SubscriptionID != subB || got.Rows[0].Cells[4] != "" || got.Rows[0].Cells[12] != "" || got.Rows[1].Cells[0] != "A") {
				t.Fatal("null/order/display subscription projection")
			}
		})
	}
	filters := config.NewFilters()
	filters.Assessment.Include.ResourceGroups = []string{"unrelated"}
	filters.Assessment.Include.Tags = map[string]string{"Environment": "absent"}
	filters.Assessment.SetAllowedResourceTypes([]string{"microsoft.storage/storageaccounts"})
	got, err := scanner.Scan(context.Background(), selected(), filters.Assessment)
	if err != nil || len(got.Rows) != 2 {
		t.Fatal("non-subscription filters affected source SQL rows")
	}
}

func TestMalformedRowsStrictStringsScopeAndPreservedFormatting(t *testing.T) {
	rows := sourceRows(t)
	for _, kind := range []string{"missing", "duplicate", "alias", "unknown", "numeric", "control", "text limit", "row limit", "foreign", "invalid UTF8", "not object", "trailing"} {
		t.Run(kind, func(t *testing.T) {
			var values map[string]any
			json.Unmarshal(rows[1], &values)
			switch kind {
			case "missing":
				delete(values, "Edition")
			case "alias":
				values["subscriptionid"] = values["SubscriptionId"]
				delete(values, "SubscriptionId")
			case "unknown":
				values["unreviewed"] = "value"
			case "numeric":
				values["vCores"] = 8
			case "control":
				values["Name"] = "bad\nname"
			case "text limit":
				values["MigrationRecommendation"] = strings.Repeat("x", 4097)
			case "row limit":
				for _, key := range []string{"MigrationRecommendation", "ESUCostBasis", "Subscription", "Name", "Location", "Edition", "EOLStatus", "CloudType", "SQLVersion", "ServiceType", "SQLLicenseType", "ESUApplicable", "ESUEnabled", "MigrationTargetTier", "SQLMIMigrationVerdict", "ArcServerName", "ConsolidationRatio"} {
					values[key] = strings.Repeat("x", 4096)
				}
			case "foreign":
				values["SubscriptionId"] = "33333333-3333-4333-8333-333333333333"
			}
			bad, _ := json.Marshal(values)
			switch kind {
			case "duplicate":
				bad = append([]byte(`{"Name":"duplicate",`), bad[1:]...)
			case "invalid UTF8":
				bad = []byte{'{', '"', 'x', '"', ':', '"', 0xff, '"', '}'}
			case "not object":
				bad = []byte(`[]`)
			case "trailing":
				bad = append(bad, []byte(` {}`)...)
			}
			s := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) {
				return &arg.Response{Data: []json.RawMessage{rows[1], bad}}, nil
			}))
			got, err := s.Scan(context.Background(), selected(), nil)
			if err != nil || len(got.Rows) != 1 || got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Message != "skipped 1 invalid sql-eol rows" {
				t.Fatalf("malformed row truth/retention: %+v %v", got, err)
			}
		})
	}
	var values map[string]any
	json.Unmarshal(rows[1], &values)
	values["EstSQLMIMonthlySaving"] = "006702.40"
	values["ESUStartDate"] = nil
	b, _ := json.Marshal(values)
	r, err := decode(b)
	if err != nil || r.Cells[30] != "006702.40" || r.Cells[12] != "" {
		t.Fatal("source string/null formatting changed")
	}
}

func TestPartialMissingContinuationCancellationAndBounds(t *testing.T) {
	rows := sourceRows(t)
	for _, mode := range []string{"empty", "all malformed", "missing", "denied", "later", "cycle", "pages", "truncated", "cancel", "row page limit", "byte page limit"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			s := NewWithTransport(transportFunc(func(_ context.Context, request arg.Request) (*arg.Response, error) {
				calls++
				response := &arg.Response{Data: []json.RawMessage{rows[1]}}
				switch mode {
				case "empty":
					response.Data = []json.RawMessage{}
				case "all malformed":
					response.Data = []json.RawMessage{json.RawMessage(`{}`)}
				case "missing":
					response.Data = nil
				case "denied":
					return nil, errors.New("secret-canary provider body")
				case "later":
					if calls == 2 {
						return nil, errors.New("secret-canary provider body")
					}
					token := "second"
					response.SkipToken = &token
				case "cycle":
					token := "repeat"
					response.SkipToken = &token
				case "pages":
					token := fmt.Sprint(calls)
					response.SkipToken = &token
				case "truncated":
					truncated := "true"
					response.ResultTruncated = &truncated
				case "cancel":
					cancel()
					return nil, context.Canceled
				case "row page limit":
					response.Data = make([]json.RawMessage, 1001)
				case "byte page limit":
					response.Data = []json.RawMessage{json.RawMessage(strings.Repeat("x", MaxPageBytes+1))}
				}
				return response, nil
			}))
			got, err := s.Scan(ctx, selected(), nil)
			if mode == "empty" {
				if err != nil || got.Health.Status != assessment.StageCompleted || len(got.Rows) != 0 {
					t.Fatal("confirmed empty not healthy")
				}
				return
			}
			if mode == "all malformed" {
				if err != nil || got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Rows) != 0 || len(got.Health.Warnings) != 1 {
					t.Fatal("all malformed input concealed")
				}
				return
			}
			if err == nil || got.Health.Status != assessment.StageFailed {
				t.Fatalf("incomplete query false success: %+v %v", got, err)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "secret-canary") || strings.Contains(err.Error(), "secret-canary") {
				t.Fatal("provider error leaked")
			}
			if mode == "later" && (calls != 2 || len(got.Rows) != 1) {
				t.Fatal("completed page rows lost")
			}
			if mode == "pages" && (calls != 64 || len(got.Rows) != 64) {
				t.Fatalf("bounded retained pages: %d %d", calls, len(got.Rows))
			}
			if mode == "cancel" && !errors.Is(err, context.Canceled) {
				t.Fatal("context identity lost")
			}
		})
	}
	calls := 0
	s := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) {
		calls++
		return &arg.Response{Data: rows}, nil
	}))
	// Use a UUID containing letters to exercise case-insensitive duplicate detection.
	letterID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	for _, scope := range []map[string]string{{"bad": "A"}, {letterID: "A", strings.ToUpper(letterID): "duplicate"}} {
		if _, err := s.Scan(context.Background(), scope, nil); err == nil {
			t.Fatal("invalid scope accepted")
		}
	}
	tooMany := map[string]string{}
	for i := 0; i < 3001; i++ {
		tooMany[fmt.Sprintf("%08x-aaaa-4aaa-8aaa-aaaaaaaaaaaa", i)] = "name"
	}
	if _, err := s.Scan(context.Background(), tooMany, nil); err == nil || calls != 0 {
		t.Fatal("scope guard reached transport")
	}
	b := &budgetTransport{transport: transportFunc(func(context.Context, arg.Request) (*arg.Response, error) {
		return &arg.Response{Data: []json.RawMessage{rows[1]}}, nil
	}), data: []json.RawMessage{}, bytes: MaxDataBytes}
	if _, err := b.Do(context.Background(), arg.Request{}); err == nil || len(b.data) != 0 {
		t.Fatal("whole data byte limit")
	}
}

type credential struct{ t *testing.T }

func (c credential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if !reflect.DeepEqual(options.Scopes, []string{"https://management.usgovcloudapi.net/.default"}) {
		c.t.Error("selected sovereign audience")
	}
	return azcore.AccessToken{Token: "synthetic-sql-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type httpFunc func(*http.Request) (*http.Response, error)

func (f httpFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

type ownedBody struct {
	io.Reader
	closed *bool
}

func (b ownedBody) Close() error { *b.closed = true; return nil }

func TestAuthenticatedQueryEnvelopeAndOwnedTransport(t *testing.T) {
	rows := sourceRows(t)
	closed := false
	calls := 0
	opts := azure.DefaultHTTPClientOptions(time.Second)
	opts.MaxRetries = -1
	opts.Scope = "https://management.usgovcloudapi.net/.default"
	opts.Transport = httpFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "POST" || r.URL.String() != "https://management.usgovcloudapi.net/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || r.Header.Get("Authorization") != "Bearer synthetic-sql-token" {
			t.Fatal("authenticated read-oriented destination")
		}
		var request arg.Request
		if json.NewDecoder(r.Body).Decode(&request) != nil || request.Query != Query || !reflect.DeepEqual(request.Subscriptions, []string{subA, subB}) || request.Options == nil || request.Options.Top == nil || *request.Options.Top != 1000 || request.Options.ResultFormat != "objectArray" || request.Options.AuthorizationScopeFilter != nil {
			t.Fatal("query/body selected scope")
		}
		data, _ := json.Marshal(rows)
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: ownedBody{strings.NewReader(`{"data":` + string(data) + `}`), &closed}, Request: r}, nil
	})
	client := azure.NewHTTPClient(credential{t}, opts)
	s, err := NewWithHTTPClient("https://management.usgovcloudapi.net", client)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Scan(context.Background(), selected(), nil)
	if err != nil || len(got.Rows) != 2 || calls != 1 || !closed {
		t.Fatalf("real bounded authenticated request/closure: %+v %v", got, err)
	}
	for _, endpoint := range []string{"http://management.invalid", "https://user@management.invalid", "https://management.invalid/path", "https://management.invalid?", "https://management.invalid#fragment", "https://management.invalid#"} {
		if _, err := NewWithHTTPClient(endpoint, client); err == nil {
			t.Fatal("unsafe root origin accepted")
		}
	}
	for _, raw := range []string{`{"data":[],"Data":[]}`, `{"data":[],"data":[]}`, `{"data":[]} {}`, `[]`, `{"count":0}`} {
		if validateEnvelope([]byte(raw)) == nil {
			t.Fatal("ambiguous envelope accepted")
		}
	}
	if validateEnvelope([]byte(`{"data":[],"count":0,"resultTruncated":"false"}`)) != nil {
		t.Fatal("unique metadata rejected")
	}
}

func TestRepeatedConcurrentOwnedExecution(t *testing.T) {
	rows := sourceRows(t)
	s := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) { return &arg.Response{Data: rows}, nil }))
	var workers sync.WaitGroup
	for _, id := range []string{subA, subB} {
		workers.Add(1)
		go func(id string) {
			defer workers.Done()
			for i := 0; i < 3; i++ {
				got, err := s.Scan(context.Background(), selected(), filterFunc(func(rowID string) bool { return rowID != id }))
				if err != nil || len(got.Rows) != 1 || got.Rows[0].SubscriptionID != id {
					t.Error("isolated filtering/rows")
					return
				}
				got.Rows[0].Cells[2] = "mutated output"
				got.Columns[0] = "mutated header"
			}
		}(id)
	}
	workers.Wait()
}

type countedBody struct {
	io.Reader
	readBytes *int
	closed    *bool
}

func (b countedBody) Read(p []byte) (int, error) {
	n, err := b.Reader.Read(p)
	*b.readBytes += n
	return n, err
}
func (b countedBody) Close() error { *b.closed = true; return nil }

func TestAuthenticatedDenialThrottleAndAttemptBodyBound(t *testing.T) {
	for _, mode := range []string{"denied", "throttled", "oversized"} {
		t.Run(mode, func(t *testing.T) {
			calls, readBytes := 0, 0
			closed := false
			opts := azure.DefaultHTTPClientOptions(time.Second)
			opts.MaxRetries = -1
			opts.Scope = "https://management.usgovcloudapi.net/.default"
			opts.Transport = httpFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Method != "POST" || r.URL.String() != "https://management.usgovcloudapi.net/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || r.Header.Get("Authorization") != "Bearer synthetic-sql-token" {
					t.Error("denial/body guard destination")
				}
				status := 403
				body := `{"error":{"code":"Denied","message":"provider-secret-canary"}}`
				if mode == "throttled" {
					status = 429
				}
				if mode == "oversized" {
					status = 200
					body = strings.Repeat("x", MaxPageBytes+32)
				}
				return &http.Response{StatusCode: status, Header: http.Header{}, Body: countedBody{strings.NewReader(body), &readBytes, &closed}, Request: r}, nil
			})
			s, err := NewWithHTTPClient("https://management.usgovcloudapi.net", azure.NewHTTPClient(credential{t}, opts))
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.Scan(context.Background(), selected(), nil)
			if err == nil || got.Health.Status != assessment.StageFailed || calls != 1 || !closed || len(got.Rows) != 0 {
				t.Fatalf("authenticated failed query/closure: %+v %v", got, err)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "provider-secret-canary") || strings.Contains(err.Error(), "provider-secret-canary") {
				t.Fatal("provider text leakage")
			}
			if mode == "oversized" && readBytes != MaxPageBytes+1 {
				t.Fatalf("attempt body guard read %d", readBytes)
			}
		})
	}
}
