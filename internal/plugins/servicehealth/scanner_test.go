package servicehealth

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
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const subA = "11111111-1111-4111-8111-111111111111"
const subB = "22222222-2222-4222-8222-222222222222"
const rawA = `{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"eastus","targetResourceType":"microsoft.compute/virtualmachines","percentageOfTimeWithoutEvents":99.5,"events":2,"affectedResources":3}`
const rawB = `{"subscriptionId":"22222222-2222-4222-8222-222222222222","targetRegion":"westus","targetResourceType":"microsoft.storage/storageaccounts","percentageOfTimeWithoutEvents":100,"events":0,"affectedResources":0}`
const rawC = `{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"westus","targetResourceType":"microsoft.storage/storageaccounts","percentageOfTimeWithoutEvents":99.5,"events":4,"affectedResources":5}`

type transportFunc func(context.Context, arg.Request) (*arg.Response, error)

func (f transportFunc) Do(ctx context.Context, r arg.Request) (*arg.Response, error) {
	return f(ctx, r)
}

type filterFunc func(string) bool

func (f filterFunc) IsResourceTypeExcluded(s string) bool { return f(s) }
func response(values ...string) *arg.Response {
	r := &arg.Response{Data: []json.RawMessage{}}
	for _, v := range values {
		r.Data = append(r.Data, json.RawMessage(v))
	}
	return r
}
func selected() map[string]string { return map[string]string{subA: "A", subB: "B"} }

func TestActualSourceCaptureMetadataColumnsAndEveryCell(t *testing.T) {
	if fmt.Sprintf("%x", sha256.Sum256([]byte(Query))) != "b09ab97088e0c670d0199e83cac467aca7d6cae3c3fbae856d9d885660c36aff" {
		t.Fatal("pinned source query changed")
	}
	b, err := os.ReadFile("testdata/source-output.json")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(b)) != "b310fef3dd1eff03cf251e3e948fcebff275a25581966e48c609f0ee8b432549" {
		t.Fatal("source capture bytes changed")
	}
	var capture []struct {
		Metadata struct {
			Name, Version, Description, Author, License string
			Type                                        int
		}
		Sheet       string `json:"sheet_name"`
		Description string
		Table       [][]string
	}
	if err = json.Unmarshal(b, &capture); err != nil || len(capture) != 1 {
		t.Fatalf("capture: %v", err)
	}
	s := NewWithTransport(transportFunc(func(_ context.Context, r arg.Request) (*arg.Response, error) {
		if r.Query != Query || !reflect.DeepEqual(r.Subscriptions, []string{subA, subB}) || r.Options == nil || *r.Options.Top != 1000 || r.Options.AuthorizationScopeFilter != nil {
			t.Fatal("source query/scope/options contract")
		}
		return response(rawB, rawA, rawC), nil
	}))
	table, err := s.Scan(context.Background(), selected(), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := capture[0]
	m := w.Metadata
	wantMetadata := assessment.PluginMetadata{Name: m.Name, Version: m.Version, Description: m.Description, Author: m.Author, License: m.License, Type: "internal"}
	if m.Type != 1 || table.Metadata != wantMetadata || table.SheetName != w.Sheet || table.Description != w.Description || !reflect.DeepEqual(table.Columns, w.Table[0]) {
		t.Fatalf("source table metadata/columns: %+v", table)
	}
	got := [][]string{}
	for _, r := range table.Rows {
		if r.SubscriptionID != r.Cells[0] {
			t.Fatal("maskable identity")
		}
		got = append(got, r.Cells)
	}
	if !reflect.DeepEqual(got, w.Table[1:]) || table.Health.Status != assessment.StageCompleted || table.Health.Records != 3 {
		t.Fatalf("every source cell: %+v", table)
	}
	table.Rows[0].Cells[0] = "caller change"
	table.Columns[0] = "caller change"
	again, err := s.Scan(context.Background(), selected(), nil)
	if err != nil || again.Rows[0].Cells[0] != subA || again.Columns[0] != "Subscription ID" {
		t.Fatal("owned per-run values")
	}
}

func TestFilterEmptyMissingMalformedAndPrivacy(t *testing.T) {
	for _, tc := range []struct {
		name     string
		r        *arg.Response
		queryErr error
		filter   Filter
		rows     int
		status   assessment.StageStatus
		sheet    string
		wantErr  bool
	}{
		{"filter", response(rawB, rawA, rawC), nil, filterFunc(func(t string) bool { return strings.EqualFold(t, "microsoft.storage/storageaccounts") }), 1, assessment.StageCompleted, "Service Issues", false},
		{"empty", response(), nil, nil, 0, assessment.StageCompleted, "Service Issues", false},
		{"nil result", nil, nil, nil, 0, assessment.StageFailed, "Service Health Availability", true},
		{"nil data", &arg.Response{}, nil, nil, 0, assessment.StageFailed, "Service Health Availability", true},
		{"denied", nil, errors.New("provider secret https://evil.invalid token"), nil, 0, assessment.StageFailed, "Service Health Availability", true},
		{"partial malformed", response(rawA, `null`, `{}`, strings.Replace(rawB, "22222222-2222-4222-8222-222222222222", "33333333-3333-4333-8333-333333333333", 1)), nil, nil, 1, assessment.StageCompletedWithWarnings, "Service Issues", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) { return tc.r, tc.queryErr }))
			got, err := s.Scan(context.Background(), selected(), tc.filter)
			if len(got.Rows) != tc.rows || got.Health.Status != tc.status || got.SheetName != tc.sheet || (err != nil) != tc.wantErr {
				t.Fatalf("%+v %v", got, err)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "provider secret") || err != nil && strings.Contains(err.Error(), "provider secret") {
				t.Fatal("provider error exposed")
			}
			if assessment.ValidatePluginTables([]assessment.PluginTable{got}) != nil {
				t.Fatal("invalid canonical output")
			}
		})
	}
	for _, raw := range []string{`null`, `[]`, `{}`, strings.Replace(rawA, `"events":2`, `"events":null`, 1), strings.Replace(rawA, `"events":2`, `"events":2.2`, 1), strings.Replace(rawA, `"events":2`, `"events":-1`, 1), strings.Replace(rawA, `"events":2`, `"events":9223372036854775808`, 1), strings.Replace(rawA, `99.5`, `101`, 1), strings.Replace(rawA, `99.5`, `-1`, 1), strings.Replace(rawA, `99.5`, `1e309`, 1), strings.Replace(rawA, `"events":2`, `"events":2,"events":3`, 1), strings.Replace(rawA, `"events":2`, `"Events":2`, 1), strings.Replace(rawA, `"eastus"`, `"eastus\u001b"`, 1), rawA + `{}`, strings.Repeat("x", 8193), string([]byte{255})} {
		if _, err := decode(json.RawMessage(raw)); err == nil {
			t.Fatalf("malformed row accepted: %q", raw)
		}
	}
}

func TestCompletedPageRetentionCancellationAndBudgets(t *testing.T) {
	for _, mode := range []string{"later failure", "cycle", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			s := NewWithTransport(transportFunc(func(ctx context.Context, r arg.Request) (*arg.Response, error) {
				calls++
				if calls == 1 {
					token := "second"
					v := response(rawA)
					v.SkipToken = &token
					return v, nil
				}
				if mode == "cancel" {
					cancel()
					return nil, ctx.Err()
				}
				if mode == "cycle" {
					token := "second"
					v := response(rawB)
					v.SkipToken = &token
					return v, nil
				}
				return nil, errors.New("secret provider failure")
			}))
			got, err := s.Scan(ctx, selected(), nil)
			want := 1
			if mode == "cycle" {
				want = 2
			}
			if err == nil || got.Health.Status != assessment.StageFailed || len(got.Rows) != want || got.Rows[0].Cells[3] != "99.50%" || calls != 2 {
				t.Fatalf("completed-page retention: %+v %v calls=%d", got, err, calls)
			}
			if mode == "cancel" && !errors.Is(err, context.Canceled) {
				t.Fatal("context identity lost")
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	s := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) { calls++; return response(rawA), nil }))
	got, err := s.Scan(ctx, selected(), nil)
	if !errors.Is(err, context.Canceled) || calls != 0 || got.Health.Status != assessment.StageFailed {
		t.Fatal("cancelled scope preflight")
	}
	if _, err := s.Scan(context.Background(), map[string]string{"invalid": "x"}, nil); err == nil || calls != 0 {
		t.Fatal("invalid scope before transport")
	}
	for _, mode := range []string{"pages", "page rows", "bytes"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			s := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) {
				calls++
				v := response(rawA)
				switch mode {
				case "pages":
					token := fmt.Sprint(calls)
					v.SkipToken = &token
				case "page rows":
					v.Data = make([]json.RawMessage, 1001)
				case "bytes":
					v.Data = []json.RawMessage{json.RawMessage(strings.Repeat("x", (2<<20)+1))}
				}
				return v, nil
			}))
			got, err := s.Scan(context.Background(), selected(), nil)
			if err == nil || got.Health.Status != assessment.StageFailed {
				t.Fatal("budget accepted")
			}
			if mode == "pages" && (calls != 64 || len(got.Rows) != 64) {
				t.Fatalf("bounded retained pages: %d %d", calls, len(got.Rows))
			}
		})
	}
}

type credential struct{ scopes []string }

func (c *credential) GetToken(_ context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.scopes = append(c.scopes, opts.Scopes...)
	return azcore.AccessToken{Token: "service-health-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type httpTransport func(*http.Request) (*http.Response, error)

func (f httpTransport) Do(r *http.Request) (*http.Response, error) { return f(r) }

type trackedBody struct {
	io.Reader
	closed *atomic.Int64
}

func (b trackedBody) Close() error { b.closed.Add(1); return nil }

func TestAuthenticatedReadOnlyPOSTSovereignAudienceAndClosure(t *testing.T) {
	for _, status := range []int{200, 403, 429} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			cred := &credential{}
			var closes atomic.Int64
			opts := azure.DefaultHTTPClientOptions(time.Second)
			opts.MaxRetries = -1
			opts.Scope = "https://management.core.usgovcloudapi.net/.default"
			opts.Transport = httpTransport(func(r *http.Request) (*http.Response, error) {
				if r.Method != "POST" || r.URL.String() != "https://management.usgovcloudapi.net/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" || r.Header.Get("Authorization") != "Bearer service-health-canary" {
					t.Fatal("authenticated destination/method")
				}
				var request arg.Request
				if json.NewDecoder(r.Body).Decode(&request) != nil || request.Query != Query || !reflect.DeepEqual(request.Subscriptions, []string{subA, subB}) || request.Options.AuthorizationScopeFilter != nil {
					t.Fatal("read-oriented KQL envelope")
				}
				body := `{"data":[` + rawA + `]}`
				if status != 200 {
					body = `{"error":{"code":"Denied","message":"secret-provider-body"}}`
				}
				return &http.Response{StatusCode: status, Header: http.Header{}, Body: trackedBody{strings.NewReader(body), &closes}, Request: r}, nil
			})
			s, err := NewWithHTTPClient("https://management.usgovcloudapi.net", azure.NewHTTPClient(cred, opts))
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.Scan(context.Background(), selected(), nil)
			if (err != nil) != (status != 200) || closes.Load() != 1 || !reflect.DeepEqual(cred.scopes, []string{"https://management.core.usgovcloudapi.net/.default"}) {
				t.Fatalf("pipeline: %+v %v closes=%d scopes=%v", got, err, closes.Load(), cred.scopes)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "secret-provider-body") {
				t.Fatal("body privacy")
			}
		})
	}
	for _, endpoint := range []string{"http://evil.invalid", "https://user:secret@evil.invalid", "https://example.test/path", "https://example.test?next=x", "https://example.test#fragment", "https://example.test?", "https://example.test/%2f"} {
		if _, err := NewWithHTTPClient(endpoint, azure.NewHTTPClient(&credential{}, nil)); err == nil {
			t.Fatal("unsafe origin accepted")
		}
	}
}

func TestConcurrentExecutionIsolation(t *testing.T) {
	s := NewWithTransport(transportFunc(func(_ context.Context, r arg.Request) (*arg.Response, error) {
		if r.Subscriptions[0] == subA {
			return response(rawA), nil
		}
		return response(rawB), nil
	}))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := subA
			want := "99.50%"
			if i%2 != 0 {
				id = subB
				want = "100.00%"
			}
			got, err := s.Scan(context.Background(), map[string]string{id: "A"}, nil)
			if err != nil || len(got.Rows) != 1 || got.Rows[0].Cells[0] != id || got.Rows[0].Cells[3] != want {
				t.Errorf("isolation: %+v %v", got, err)
			}
		}(i)
	}
	wg.Wait()
}

func TestEnvelopeAmbiguityBeforeSharedDecoder(t *testing.T) {
	for _, data := range []string{`{"data":[],"data":[]}`, `{"Data":[]}`, `{"data":[],"DATA":[]}`, `{"data":[]} {}`, `[]`, `{}`, `{"data":[],"$SkipToken":"x"}`, string([]byte{255})} {
		if validateEnvelope([]byte(data)) == nil {
			t.Fatalf("ambiguous envelope accepted: %q", data)
		}
	}
	if validateEnvelope([]byte(`{"data":[],"count":0,"totalRecords":0,"resultTruncated":"false"}`)) != nil {
		t.Fatal("valid ARG metadata rejected")
	}
}

func TestWholeQueryDataAndScopeBudgets(t *testing.T) {
	calls := 0
	padded := rawA + strings.Repeat(" ", 8192-len(rawA))
	scanner := NewWithTransport(transportFunc(func(context.Context, arg.Request) (*arg.Response, error) {
		calls++
		r := response()
		for i := 0; i < 100; i++ {
			r.Data = append(r.Data, json.RawMessage(padded))
		}
		token := fmt.Sprint(calls)
		r.SkipToken = &token
		return r, nil
	}))
	table, err := scanner.Scan(context.Background(), selected(), nil)
	if err == nil || table.Health.Status != assessment.StageFailed || calls != 21 || len(table.Rows) != 2000 {
		t.Fatalf("whole-query byte guard/retention: %v calls=%d rows=%d", err, calls, len(table.Rows))
	}
	scope := map[string]string{}
	for i := 0; i < 3001; i++ {
		scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "A"
	}
	before := calls
	if _, err := scanner.Scan(context.Background(), scope, nil); err == nil || calls != before {
		t.Fatal("scope budget before transport")
	}
	if _, err := scanner.Scan(context.Background(), map[string]string{"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa": "A", "AAAAAAAA-AAAA-4AAA-8AAA-AAAAAAAAAAAA": "B"}, nil); err == nil || calls != before {
		t.Fatal("case-duplicate scope before transport")
	}
	if _, err := NewWithTransport(nil).Scan(context.Background(), selected(), nil); err == nil {
		t.Fatal("missing transport accepted")
	}
}
