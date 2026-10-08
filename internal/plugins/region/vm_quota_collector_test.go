package region

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

// The unchanged SDK input bodies and requests are provenance, not target oracles.
// Full expected raw evidence below is authored independently, including corrections.
func TestVMQuotaCollectorRetainedEvidence(t *testing.T) {
	var inputs struct {
		Cases map[string]struct {
			Pages []struct {
				Status int
				Body   string
			}
			Cancelled bool
		}
	}
	var source map[string]struct{ Requests []string }
	for path, out := range map[string]any{"testdata/source-vm-quota-inputs.json": &inputs, "testdata/source-vm-quota-outputs.json": &source} {
		b, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		if e = json.Unmarshal(b, out); e != nil {
			t.Fatal(e)
		}
	}
	const id = "11111111-1111-1111-1111-111111111111"
	scope := map[string]string{id: "Selected"}
	req := QuotaRequest{id, "eastus", "VM"}
	first := QuotaUsage{"FirstFamily", "FirstFamily", 1, 10, true, true}
	wants := map[string]struct {
		status, code string
		rows         []QuotaUsage
	}{
		"selection":       {"complete", "", []QuotaUsage{{"standardDSFamily", "standardDSFamily", 10, 100, true, true}, {"standarddsfamily", "standarddsfamily", 10, 100, true, true}, {"cores", "cores", 10, 100, true, true}, {"NotFamilySuffix", "NotFamilySuffix", 10, 100, true, true}}},
		"localized":       {"complete", "", []QuotaUsage{{"NoLabelFamily", "", 9, 10, true, true}, {"LabelFamily", "VM Label", 9, 10, true, true}}},
		"paged":           {"complete", "", []QuotaUsage{first, {"SecondFamily", "SecondFamily", 9, 10, true, true}}},
		"empty":           {"complete", "", []QuotaUsage{}},
		"missing-current": {"complete", "", []QuotaUsage{{ResourceName: "MissingCurrentFamily", Limit: 10, LimitKnown: true}}},
		"missing-skipped": {"complete", "", []QuotaUsage{{Current: 1, Limit: 10, CurrentKnown: true, LimitKnown: true}, {Current: 1, Limit: 10, CurrentKnown: true, LimitKnown: true}, {Current: 1, Limit: 10, CurrentKnown: true, LimitKnown: true}, {ResourceName: "NoLimitFamily", Current: 1, CurrentKnown: true}, {ResourceName: "ZeroFamily", LimitKnown: true}, {ResourceName: "NegativeLimitFamily", Limit: -1, LimitKnown: true}, {ResourceName: "cores", Limit: 10, LimitKnown: true}}},
		"null-item":       {"unknown", "quota_invalid_response", []QuotaUsage{}},
		"arithmetic":      {"unknown", "quota_count_invalid", []QuotaUsage{}},
		"malformed":       {"unknown", "quota_invalid_response", []QuotaUsage{}},
	}
	for _, code := range []string{"400", "403", "404", "405"} {
		wants["http"+code] = struct {
			status, code string
			rows         []QuotaUsage
		}{"unknown", "quota_status_invalid", []QuotaUsage{}}
		wants["later-http"+code] = struct {
			status, code string
			rows         []QuotaUsage
		}{"partial", "quota_status_invalid", []QuotaUsage{first}}
	}
	wants["later-malformed"] = struct {
		status, code string
		rows         []QuotaUsage
	}{"partial", "quota_invalid_response", []QuotaUsage{first}}
	if len(inputs.Cases) != 19 || len(wants) != 18 || len(source) != 19 {
		t.Fatal("retained scenario inventory drift")
	}
	for name, fixture := range inputs.Cases {
		t.Run(name, func(t *testing.T) {
			calls := []string{}
			c := quotaCollector(t, func(_ context.Context, u string, n int64) ([]byte, *http.Response, error) {
				i := len(calls)
				if i >= len(source[name].Requests) || u != source[name].Requests[i] || n != 1048576 {
					t.Fatalf("VM exact retained request: %q/%d", u, n)
				}
				calls = append(calls, u)
				p := fixture.Pages[i]
				return []byte(p.Body), &http.Response{StatusCode: p.Status}, nil
			})
			ctx := context.Background()
			if fixture.Cancelled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			r, e := c.Collect(ctx, scope, req)
			if !reflect.DeepEqual(calls, source[name].Requests) {
				t.Fatalf("VM request sequence: %#v", calls)
			}
			if fixture.Cancelled {
				if r != nil || !errors.Is(e, context.Canceled) {
					t.Fatalf("VM cancelled: %#v %v", r, e)
				}
				return
			}
			w, ok := wants[name]
			if !ok {
				t.Fatal("missing independent oracle")
			}
			want := &RESTQuotaCollection{Evidence: QuotaEvidence{Request: req, Status: w.status, Usages: w.rows}, FailureCode: w.code}
			if e != nil || !reflect.DeepEqual(r, want) {
				t.Fatalf("VM full evidence %s: got%#v want%#v error%v", name, r, want, e)
			}
			calc, e := CalculateQuota(context.Background(), scope, []QuotaRequest{req}, []QuotaEvidence{r.Evidence})
			if e != nil {
				t.Fatal(e)
			}
			cells := [][]string{}
			if calc.Table != nil {
				for _, row := range calc.Table.Rows {
					cells = append(cells, row.Cells)
				}
			}
			expected := [][]string{}
			warnings := []string{}
			switch name {
			case "selection":
				expected = [][]string{{"Selected", "eastus", "VM", "standardDSFamily", "10", "100", "90", "90.0%", "OK"}, {"Selected", "eastus", "VM", "NotFamilySuffix", "10", "100", "90", "90.0%", "OK"}}
			case "localized":
				expected = [][]string{{"Selected", "eastus", "VM", "NoLabelFamily", "9", "10", "1", "10.0%", "Near Limit"}, {"Selected", "eastus", "VM", "VM Label", "9", "10", "1", "10.0%", "Near Limit"}}
			case "paged":
				expected = [][]string{{"Selected", "eastus", "VM", "FirstFamily", "1", "10", "9", "90.0%", "OK"}, {"Selected", "eastus", "VM", "SecondFamily", "9", "10", "1", "10.0%", "Near Limit"}}
			case "missing-current":
				warnings = []string{"quota_usage_unknown"}
			case "missing-skipped":
				warnings = []string{"quota_usage_unknown"}
			}
			if w.status == "partial" {
				expected = [][]string{{"Selected", "eastus", "VM", "FirstFamily", "1", "10", "9", "90.0%", "OK"}}
				warnings = []string{"quota_evidence_partial"}
			}
			if w.status == "unknown" {
				warnings = []string{"quota_evidence_unknown"}
			}
			gotWarnings := []string{}
			for _, warning := range calc.Health.Warnings {
				gotWarnings = append(gotWarnings, warning.Code)
			}
			status := assessment.StageCompleted
			if len(warnings) > 0 {
				status = assessment.StageCompletedWithWarnings
			}
			if !reflect.DeepEqual(cells, expected) || !reflect.DeepEqual(gotWarnings, warnings) || calc.Health.Status != status || calc.Health.Records != len(expected) || (len(expected) == 0 && calc.Table != nil) || (len(expected) > 0 && (calc.Table == nil || !reflect.DeepEqual(calc.Table.Health, calc.Health))) {
				t.Fatalf("VM literal calculation %s: cells%#v warnings%#v health%#v", name, cells, gotWarnings, calc.Health)
			}
		})
	}
}

func TestVMQuotaCollectorBoundary(t *testing.T) {
	path := "/subscriptions/" + quotaTestID + "/providers/Microsoft.Compute/locations/eastus/usages"
	for _, next := range []string{"https://evil.example" + path + "?api-version=2024-11-01", path + "?api-version=2023-01-01", strings.Replace(path, "Microsoft.Compute", "Microsoft.Network", 1) + "?api-version=2024-11-01", strings.Replace(path, quotaTestID, "11111111-1111-1111-1111-111111111111", 1) + "?api-version=2024-11-01"} {
		calls := 0
		b, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": next})
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			return quotaPage(string(b))
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("VM"))
		if e != nil || calls != 1 || r.Evidence.Status != "partial" || r.FailureCode != "quota_unsafe_continuation" || len(r.Evidence.Usages) != 0 {
			t.Fatalf("VM unsafe authenticated continuation: %#v %v calls%d", r, e, calls)
		}
	}
	for _, status := range []int{400, 404, 405} {
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return nil, &http.Response{StatusCode: status}, errors.New("secret token URL body")
		})
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("VM"))
		if e != nil || r.Evidence.Status != "unknown" || r.FailureCode != "quota_request_failed" || len(r.Evidence.Usages) != 0 {
			t.Fatalf("VM errored status: %#v %v", r, e)
		}
	}
	for _, body := range []string{`{"value":[{"name":{"value":"ZeroFamily"},"currentValue":0,"limit":10}]}`, `{"value":[{"name":{"value":"ZeroFamily"},"limit":10}]}`} {
		c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { return quotaPage(body) })
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("VM"))
		if e != nil || r.Evidence.Status != "complete" || len(r.Evidence.Usages) != 1 || r.Evidence.Usages[0].CurrentKnown != strings.Contains(body, "currentValue") {
			t.Fatalf("VM zero/presence: %#v %v", r, e)
		}
	}
	c := quotaCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return nil, nil, context.DeadlineExceeded
	})
	if r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("VM")); r != nil || !errors.Is(e, context.DeadlineExceeded) {
		t.Fatalf("VM deadline: %#v %v", r, e)
	}
}

func TestVMQuotaCollectorAuthenticatedClient(t *testing.T) {
	for _, status := range []int{200, 404, 405, 403, 302} {
		closed := false
		calls := 0
		client := azure.NewHTTPClient(quotaCredential{t, "https://management.usgovcloudapi.net/.default"}, &azure.HTTPClientOptions{MaxRetries: -1, Scope: "https://management.usgovcloudapi.net/.default", Transport: quotaTransport(func(r *http.Request) (*http.Response, error) {
			calls++
			if r.Method != "GET" || r.Body != nil || r.URL.String() != "https://management.usgovcloudapi.net/subscriptions/"+quotaTestID+"/providers/Microsoft.Compute/locations/eastus/usages?api-version=2024-11-01" || r.Header.Get("Authorization") != "Bearer synthetic-quota-token" {
				t.Fatalf("VM authenticated request: %#v", r)
			}
			return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: quotaClosedBody{strings.NewReader(`{"value":[]}`), &closed}, Request: r}, nil
		})})
		c, e := NewRESTQuotaCollector("https://management.usgovcloudapi.net", client)
		if e != nil {
			t.Fatal(e)
		}
		r, e := c.Collect(context.Background(), quotaScope(), quotaRequest("VM"))
		want := "unknown"
		if status == 200 {
			want = "complete"
		}
		if e != nil || r.Evidence.Status != want || calls != 1 || !closed {
			t.Fatalf("VM status/ownership%d: %#v %v calls%d closed%v", status, r, e, calls, closed)
		}
	}
}
