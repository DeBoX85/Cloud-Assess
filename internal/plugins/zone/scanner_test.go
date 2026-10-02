package zone

import (
	"context"
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
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const subA = "11111111-1111-4111-8111-111111111111"
const subB = "22222222-2222-4222-8222-222222222222"
const origin = "https://management.azure.com"
const literalPage = `{"value":[{"name":"eastus","displayName":"East US","availabilityZoneMappings":[{"logicalZone":"2","physicalZone":"eastus-az3"},{"logicalZone":"1","physicalZone":"eastus-az1"}]},{"name":"westus","displayName":"West US"}]}`

func TestActualPinnedSourceParserCapture(t *testing.T) {
	// Captured by executing unchanged parseZoneMappings at the source pin in
	// docs/ZONE_MAPPING.md, not generated using the target decoder.
	data, err := os.ReadFile("testdata/source-rows.json")
	if err != nil {
		t.Fatal(err)
	}
	var want []Row
	if err = json.Unmarshal(data, &want); err != nil {
		t.Fatal(err)
	}
	got, next, err := decodePage([]byte(literalPage), subscription{subA, "A"})
	if err != nil || next != "" || !reflect.DeepEqual(got, want) {
		t.Fatalf("complete source parser rows differ: got=%+v want=%+v err=%v", got, want, err)
	}
}

type getterFunc func(context.Context, string, int64) ([]byte, error)

func (f getterFunc) GetBounded(c context.Context, u string, n int64) ([]byte, error) {
	return f(c, u, n)
}
func scanner(t *testing.T, f getterFunc) *Scanner {
	t.Helper()
	s, err := NewScanner(origin, f)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
func initial(id string) string {
	return origin + "/subscriptions/" + id + "/locations?api-version=2022-12-01"
}
func withNext(body, next string) []byte {
	return []byte(strings.TrimSuffix(body, "}") + `,"nextLink":` + strconvJSON(next) + `}`)
}
func strconvJSON(s string) string { b, _ := json.Marshal(s); return string(b) }

func TestSourceRowsPaginationOwnershipAndIsolation(t *testing.T) {
	s := scanner(t, func(_ context.Context, u string, n int64) ([]byte, error) {
		if n != 1048576 {
			t.Fatal("literal body ceiling changed")
		}
		if strings.Contains(u, "skiptoken=second") {
			return []byte(`{"value":[{"name":"north","availabilityZoneMappings":[{"logicalZone":"3","physicalZone":"north-az2"}]}]}`), nil
		}
		return withNext(literalPage, "?api-version=2022-12-01&skiptoken=second"), nil
	})
	scope := map[string]string{subA: "A"}
	want := []Row{{subA, "A", "eastus", "East US", "1", "eastus-az1"}, {subA, "A", "eastus", "East US", "2", "eastus-az3"}, {subA, "A", "north", "", "3", "north-az2"}}
	for i := 0; i < 2; i++ {
		got, err := s.Scan(context.Background(), scope)
		if err != nil || len(got.Failures) != 0 || !reflect.DeepEqual(got.Rows, want) {
			t.Fatalf("complete literal rows: %+v %v", got, err)
		}
		got.Rows[0].Location = "changed by caller"
	}
	if scope[subA] != "A" {
		t.Fatal("caller scope mutated")
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, e := s.Scan(context.Background(), map[string]string{subB: "B"})
			if e != nil || len(got.Rows) != 3 || got.Rows[0].SubscriptionID != subB {
				t.Error("concurrent run leaked scope")
			}
		}()
	}
	wg.Wait()
}

func TestPartialFailureRetainsEarlierPagesAndOtherSubscriptions(t *testing.T) {
	s := scanner(t, func(_ context.Context, u string, _ int64) ([]byte, error) {
		if u == initial(subA) {
			return withNext(literalPage, "?api-version=2022-12-01&skiptoken=denied"), nil
		}
		if strings.Contains(u, "skiptoken=denied") {
			return nil, errors.New("provider secret payload")
		}
		return []byte(literalPage), nil
	})
	got, err := s.Scan(context.Background(), map[string]string{subA: "Z", subB: "A"})
	if err != nil || len(got.Rows) != 4 || !reflect.DeepEqual(got.Failures, []Failure{{subA, "zone_request_failed"}}) || got.Rows[0].SubscriptionID != subB {
		t.Fatalf("partial: %+v %v", got, err)
	}
	encoded, _ := json.Marshal(got.Failures)
	if strings.Contains(string(encoded), "secret") {
		t.Fatal("provider text leaked")
	}
}

func TestMalformedEmptyAndAmbiguousResponses(t *testing.T) {
	for _, body := range []string{`{}`, `null`, `{"value": null}`, `{"value": {}}`, `{"value":[null]}`, `{"value":[],"value":[]}`, `{"value":[],"Value":[]}`, `{"value":[]} {}`, `{"value":[{"name":"x","availabilityZoneMappings":[null]}]}`, `{"value":[{"name":"x","availabilityZoneMappings":[{"logicalZone":"1"}]}]}`, `{"value":[],"nextLink":false}`, "{\"value\":[],\"x\":\"\xff\"}", `{"value":[],"unknown":` + strings.Repeat("[", 66) + "0" + strings.Repeat("]", 66) + "}"} {
		t.Run(body[:min(len(body), 24)], func(t *testing.T) {
			s := scanner(t, func(context.Context, string, int64) ([]byte, error) { return []byte(body), nil })
			got, e := s.Scan(context.Background(), map[string]string{subA: "A"})
			if e != nil || len(got.Rows) != 0 || len(got.Failures) != 1 || got.Failures[0].Code != "zone_invalid_response" {
				t.Fatalf("invalid response accepted: %+v %v", got, e)
			}
		})
	}
	for _, body := range []string{`{"value":[]}`, `{"value":[{"name":"empty","unknown":{"safe":true}}],"nextLink":null}`} {
		s := scanner(t, func(context.Context, string, int64) ([]byte, error) { return []byte(body), nil })
		got, e := s.Scan(context.Background(), map[string]string{subA: "A"})
		if e != nil || len(got.Rows) != 0 || len(got.Failures) != 0 {
			t.Fatalf("valid empty rejected %+v %v", got, e)
		}
	}
}

func TestUnsafeContinuationBlockedBeforeAuthenticatedGetter(t *testing.T) {
	for _, link := range []string{"https://evil.test/subscriptions/" + subA + "/locations?api-version=2022-12-01", initial(subB), strings.Replace(initial(subA), "2022-12-01", "2020-01-01", 1), initial(subA) + "&api-version=2022-12-01", initial(subA) + "#fragment", strings.Replace(initial(subA), "https://", "http://", 1), strings.Replace(initial(subA), "management.azure.com", "user@management.azure.com", 1), strings.Replace(initial(subA), "/locations", "/%6cocations", 1), "?api-version=2022-12-01&x=\\bad", "?api-version=2022-12-01;x=1", strings.Repeat("x", 8193)} {
		t.Run(link[:min(len(link), 32)], func(t *testing.T) {
			var calls atomic.Int32
			s := scanner(t, func(context.Context, string, int64) ([]byte, error) {
				calls.Add(1)
				return withNext(literalPage, link), nil
			})
			got, e := s.Scan(context.Background(), map[string]string{subA: "A"})
			if e != nil || calls.Load() != 1 || len(got.Rows) != 2 || len(got.Failures) != 1 || got.Failures[0].Code != "zone_unsafe_continuation" {
				t.Fatalf("unsafe continuation invoked getter: calls=%d result=%+v err=%v", calls.Load(), got, e)
			}
		})
	}
}

func TestCyclesAndAdapterLimits(t *testing.T) {
	var calls atomic.Int32
	s := scanner(t, func(context.Context, string, int64) ([]byte, error) {
		calls.Add(1)
		return withNext(literalPage, initial(subA)), nil
	})
	got, _ := s.Scan(context.Background(), map[string]string{subA: "A"})
	if calls.Load() != 1 || got.Failures[0].Code != "zone_continuation_cycle" {
		t.Fatal("cycle not bounded")
	}
	calls.Store(0)
	s = scanner(t, func(context.Context, string, int64) ([]byte, error) {
		n := calls.Add(1)
		return withNext(`{"value":[]}`, fmt.Sprintf("?api-version=2022-12-01&skiptoken=%d", n)), nil
	})
	got, _ = s.Scan(context.Background(), map[string]string{subA: "A"})
	if calls.Load() != 64 || got.Failures[0].Code != "zone_page_limit" {
		t.Fatalf("literal page bound: calls=%d %+v", calls.Load(), got)
	}
	s = scanner(t, func(context.Context, string, int64) ([]byte, error) { return []byte(strings.Repeat("x", 1048577)), nil })
	got, _ = s.Scan(context.Background(), map[string]string{subA: "A"})
	if got.Failures[0].Code != "zone_byte_limit" {
		t.Fatal("literal byte ceiling ignored")
	}
	row := `{"name":"x","availabilityZoneMappings":[{"logicalZone":"1","physicalZone":"x-az1"}]}`
	large := `{"value":[` + strings.TrimSuffix(strings.Repeat(row+",", 2049), ",") + `]}`
	s = scanner(t, func(context.Context, string, int64) ([]byte, error) { return []byte(large), nil })
	got, _ = s.Scan(context.Background(), map[string]string{subA: "A"})
	if len(got.Failures) != 1 || len(got.Rows) != 0 {
		t.Fatal("literal row ceiling ignored")
	}
	// One MiB pages remain individually valid; cumulative bytes are bounded too.
	calls.Store(0)
	s = scanner(t, func(context.Context, string, int64) ([]byte, error) {
		n := calls.Add(1)
		b := withNext(`{"value":[],"padding":"`+strings.Repeat("x", 1048400)+`"}`, fmt.Sprintf("?api-version=2022-12-01&skiptoken=%d", n))
		return b, nil
	})
	got, _ = s.Scan(context.Background(), map[string]string{subA: "A"})
	if calls.Load() != 9 || got.Failures[0].Code != "zone_byte_limit" {
		t.Fatalf("cumulative body bound: %d %+v", calls.Load(), got)
	}
}

func TestAggregateAdmissionAndLaterPageRowLimit(t *testing.T) {
	row := `{"name":"x","availabilityZoneMappings":[{"logicalZone":"1","physicalZone":"x-az1"}]}`
	page := `{"value":[` + strings.TrimSuffix(strings.Repeat(row+",", 2048), ",") + `]}`
	s := scanner(t, func(context.Context, string, int64) ([]byte, error) { return []byte(page), nil })
	scope := map[string]string{}
	for i := 0; i < 33; i++ {
		scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "same display name"
	}
	got, err := s.Scan(context.Background(), scope)
	if err != nil || len(got.Rows) != 65536 || !reflect.DeepEqual(got.Failures, []Failure{{"00000020-1111-4111-8111-111111111111", "zone_row_limit"}}) {
		t.Fatalf("literal aggregate row admission: rows=%d failures=%+v err=%v", len(got.Rows), got.Failures, err)
	}
	s = scanner(t, func(_ context.Context, u string, _ int64) ([]byte, error) {
		if u == initial(subA) {
			return withNext(page, "?api-version=2022-12-01&skiptoken=next"), nil
		}
		return []byte(literalPage), nil
	})
	got, err = s.Scan(context.Background(), map[string]string{subA: "A"})
	if err != nil || len(got.Rows) != 2048 || len(got.Failures) != 1 || got.Failures[0].Code != "zone_row_limit" {
		t.Fatal("later-page row limit discarded valid earlier page or claimed completion")
	}
	got, err = s.Scan(context.Background(), map[string]string{})
	if err != nil || len(got.Rows) != 0 || len(got.Failures) != 0 {
		t.Fatal("empty selected scope rejected")
	}
}

func TestScopeValidationWorkerBoundAndCancellation(t *testing.T) {
	guard := scanner(t, func(context.Context, string, int64) ([]byte, error) {
		t.Error("invalid/cancelled scope invoked getter")
		return nil, nil
	})
	if _, e := guard.Scan(context.Background(), map[string]string{"../escape": "bad"}); e == nil {
		t.Fatal("invalid scope accepted")
	}
	letterID := "abcdefab-abcd-4abc-8abc-abcdefabcdef"
	if _, e := guard.Scan(context.Background(), map[string]string{letterID: "x", strings.ToUpper(letterID): "y"}); e == nil {
		t.Fatal("duplicate normalized subscription accepted")
	}
	tooMany := map[string]string{}
	for i := 0; i < 257; i++ {
		tooMany[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "n"
	}
	if _, e := guard.Scan(context.Background(), tooMany); e == nil {
		t.Fatal("scope ceiling ignored")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, e := guard.Scan(ctx, map[string]string{subA: "A"})
	if !errors.Is(e, context.Canceled) || got.Failures[0].Code != "zone_cancelled" {
		t.Fatal("cancelled scope not surfaced")
	}
	var active, peak atomic.Int32
	started := make(chan struct{}, 5)
	s := scanner(t, func(c context.Context, _ string, _ int64) ([]byte, error) {
		n := active.Add(1)
		defer active.Add(-1)
		for old := peak.Load(); n > old && !peak.CompareAndSwap(old, n); old = peak.Load() {
		}
		started <- struct{}{}
		<-c.Done()
		return nil, c.Err()
	})
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	scope := map[string]string{}
	for i := 0; i < 20; i++ {
		scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "n"
	}
	done := make(chan Result, 1)
	go func() { r, _ := s.Scan(ctx, scope); done <- r }()
	for i := 0; i < 5; i++ {
		select {
		case <-started:
		case <-time.After(3 * time.Second):
			t.Fatal("workers did not start")
		}
	}
	cancel()
	select {
	case got := <-done:
		if peak.Load() != 5 || active.Load() != 0 || len(got.Failures) != 20 {
			t.Fatalf("worker/cancel bound: %d %+v", peak.Load(), got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled workers blocked")
	}
}

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

type credential struct{ calls atomic.Int32 }

func (c *credential) GetToken(_ context.Context, o policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	if !reflect.DeepEqual(o.Scopes, []string{"https://management.usgovcloudapi.net/.default"}) {
		return azcore.AccessToken{}, errors.New("wrong audience")
	}
	return azcore.AccessToken{Token: "synthetic-canary", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestSelectedCloudRealAuthenticatedPipelineBoundary(t *testing.T) {
	cred := &credential{}
	var calls atomic.Int32
	options := azure.DefaultHTTPClientOptions(time.Second)
	options.Scope = "https://management.usgovcloudapi.net/.default"
	options.MaxRetries = -1
	options.Transport = transportFunc(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		if r.Method != "GET" || r.URL.String() != "https://management.usgovcloudapi.net/subscriptions/"+subA+"/locations?api-version=2022-12-01" || r.Header.Get("Authorization") != "Bearer synthetic-canary" {
			t.Error("authenticated request contract")
		}
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(string(withNext(literalPage, initial(subA))))), Request: r}, nil
	})
	s, e := NewScanner("https://management.usgovcloudapi.net/", azure.NewHTTPClient(cred, options))
	if e != nil {
		t.Fatal(e)
	}
	got, e := s.Scan(context.Background(), map[string]string{subA: "A"})
	if e != nil || calls.Load() != 1 || cred.calls.Load() != 1 || len(got.Rows) != 2 || got.Failures[0].Code != "zone_unsafe_continuation" {
		t.Fatalf("cross-cloud continuation boundary: %+v %v calls=%d tokens=%d", got, e, calls.Load(), cred.calls.Load())
	}
	for _, endpoint := range []string{"http://example.test", "https://example.test/path", "https://u@example.test", "https://example.test?x=1", "https://example.test#x", "https://example.test\\x"} {
		if _, e := NewScanner(endpoint, azure.NewHTTPClient(cred, options)); e == nil {
			t.Fatal("unsafe initial endpoint accepted")
		}
	}
}
