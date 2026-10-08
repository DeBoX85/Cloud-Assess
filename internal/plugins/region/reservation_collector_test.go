package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const crGroup = "/subscriptions/" + quotaTestID + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g"
const crResource = crGroup + "/capacityReservations/r"
const crGroups = `{"value":[{"id":"` + crGroup + `","name":"g","location":"EASTUS"}]}`
const crSummaries = `{"value":[{"name":"r"}]}`

func crRequest() ReservationRequest { return ReservationRequest{quotaTestID, "eastus"} }
func crDetail(n int) string {
	refs := []map[string]string{}
	for i := 0; i < n; i++ {
		refs = append(refs, map[string]string{"id": fmt.Sprintf("/subscriptions/%s/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/v%d", quotaTestID, i)})
	}
	b, _ := json.Marshal(refs)
	return `{"id":"` + crResource + `","name":"r","location":"eastus","sku":{"name":"Standard_D2s_v5","capacity":4},"properties":{"instanceView":{"utilizationInfo":{"virtualMachinesAllocated":` + string(b) + `}}}}`
}
func crCollector(t *testing.T, g quotaGetterFunc) *ReservationCollector {
	t.Helper()
	c, e := NewReservationCollector("https://management.azure.com", g)
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func crURLs() []string {
	return []string{"https://management.azure.com/subscriptions/" + quotaTestID + "/providers/Microsoft.Compute/capacityReservationGroups?api-version=2024-11-01", "https://management.azure.com" + crGroup + "/capacityReservations?api-version=2024-11-01", "https://management.azure.com" + crResource + "?%24expand=instanceView&api-version=2024-11-01"}
}
func crBodies(t *testing.T, bodies []string) (*ReservationCollection, int) {
	t.Helper()
	calls := 0
	urls := crURLs()
	c := crCollector(t, func(_ context.Context, u string, n int64) ([]byte, *http.Response, error) {
		if calls >= len(bodies) || calls >= len(urls) || u != urls[calls] || n != 1048576 {
			t.Fatalf("reservation exact request %d: %s/%d", calls, u, n)
		}
		b := bodies[calls]
		calls++
		return quotaPage(b)
	})
	r, e := c.Collect(context.Background(), quotaScope(), crRequest())
	if e != nil {
		t.Fatal(e)
	}
	return r, calls
}

// Separately adapted synthetic fixtures align source group/Get regions and ARM
// VM references. Original source bodies intentionally differ and are unchanged.
func TestReservationCollectorLiteral(t *testing.T) {
	for _, f := range []struct {
		allocated         int
		available, status string
	}{{0, "4", "Idle"}, {1, "3", "Available"}, {4, "0", "At-Capacity"}, {5, "-1", "Over-Allocated"}} {
		r, calls := crBodies(t, []string{crGroups, crSummaries, crDetail(f.allocated)})
		want := &ReservationCollection{Evidence: ReservationEvidence{Request: crRequest(), Status: "complete", Reservations: []ReservationUsage{{ResourceID: crResource, Region: "eastus", ResponseName: "r", ResponseRegion: "eastus", SKU: "Standard_D2s_v5", Reserved: 4, Allocated: int64(f.allocated), ReservedKnown: true, AllocatedKnown: true}}}}
		if calls != 3 || !reflect.DeepEqual(r, want) {
			t.Fatalf("full reservation evidence: %#v", r)
		}
		calc, e := CalculateReservations(context.Background(), quotaScope(), []ReservationRequest{crRequest()}, []ReservationEvidence{r.Evidence})
		cells := []string{"Selected", "eastus", "rg", "g", "r", "Standard_D2s_v5", "4", fmt.Sprint(f.allocated), f.available, f.status}
		if e != nil || calc.Health.Status != assessment.StageCompleted || calc.Health.Records != 1 || calc.Table == nil || !reflect.DeepEqual(calc.Table.Rows[0].Cells, cells) || !reflect.DeepEqual(calc.Table.Health, calc.Health) {
			t.Fatalf("literal reservation calculation: %#v %v", calc, e)
		}
	}
	for _, f := range []struct {
		body                string
		reserved, allocated bool
	}{{`{"sku":{"name":"Standard_D2s_v5","capacity":0},"properties":{"instanceView":{"utilizationInfo":{"virtualMachinesAllocated":[]}}}}`, true, true}, {`{"sku":{"name":"Standard_D2s_v5","capacity":4}}`, true, false}, {`{"sku":{"name":"Standard_D2s_v5"},"properties":{"instanceView":{"utilizationInfo":{"virtualMachinesAllocated":[]}}}}`, false, true}} {
		r, _ := crBodies(t, []string{crGroups, crSummaries, f.body})
		row := r.Evidence.Reservations[0]
		if r.Evidence.Status != "complete" || row.ReservedKnown != f.reserved || row.AllocatedKnown != f.allocated {
			t.Fatalf("reservation presence: %#v", r)
		}
		calc, e := CalculateReservations(context.Background(), quotaScope(), []ReservationRequest{crRequest()}, []ReservationEvidence{r.Evidence})
		if e != nil {
			t.Fatal(e)
		}
		if !f.reserved || !f.allocated {
			if calc.Table != nil || calc.Health.Status != assessment.StageCompletedWithWarnings || len(calc.Health.Warnings) != 1 || calc.Health.Warnings[0].Code != "reservation_usage_unknown" {
				t.Fatalf("unknown not Idle: %#v", calc)
			}
		} else if calc.Reservations[0].Status != "Idle" {
			t.Fatal("known zero not Idle")
		}
	}
}

func TestReservationCollectorIdentity(t *testing.T) {
	for _, stage := range []string{"summary", "get"} {
		summary, detail := crSummaries, crDetail(1)
		forged := strings.Replace(crResource, "/subscriptions/", "/ſubscriptions/", 1)
		if stage == "summary" {
			summary = `{"value":[{"name":"r","id":"` + forged + `"}]}`
		} else {
			detail = strings.Replace(detail, crResource, forged, 1)
		}
		r, calls := crBodies(t, []string{crGroups, summary, detail})
		wantCalls := 3
		if stage == "summary" {
			wantCalls = 2
		}
		if calls != wantCalls || r.Evidence.Status != "partial" || len(r.Evidence.Reservations) != 0 || r.FailureCode != "reservation_identity_invalid" {
			t.Fatalf("Unicode fixed ARM segment accepted at %s: %#v calls%d", stage, r, calls)
		}
	}
	for _, f := range []struct {
		group, summary, detail string
		calls                  int
	}{{strings.Replace(crGroups, "EASTUS", "global", 1), crSummaries, crDetail(1), 1}, {strings.Replace(crGroups, "EASTUS", "eastKus", 1), crSummaries, crDetail(1), 1}, {`{"value":[null]}`, crSummaries, crDetail(1), 1}, {`{"value":[{"id":"/bad","location":"eastus"}]}`, crSummaries, crDetail(1), 1}, {strings.Replace(crGroups, quotaTestID, "11111111-1111-1111-1111-111111111111", 1), crSummaries, crDetail(1), 1}, {strings.Replace(crGroups, `"name":"g"`, `"name":"other"`, 1), crSummaries, crDetail(1), 1}, {crGroups, `{"value":[null]}`, crDetail(1), 2}, {crGroups, `{"value":[{"name":"r/evil"}]}`, crDetail(1), 2}, {crGroups, `{"value":[{"name":"r","id":"/bad"}]}`, crDetail(1), 2}, {crGroups, `{"value":[{"name":"r"},{"name":"R"}]}`, crDetail(1), 2}, {crGroups, crSummaries, strings.Replace(crDetail(1), `"location":"eastus"`, `"location":"westus"`, 1), 3}, {crGroups, crSummaries, strings.Replace(crDetail(1), `"name":"r"`, `"name":"wrong"`, 1), 3}, {crGroups, crSummaries, strings.Replace(crDetail(1), crResource, crResource+"other", 1), 3}, {crGroups, crSummaries, strings.Replace(crDetail(1), `"capacity":4`, `"capacity":-1`, 1), 3}, {crGroups, crSummaries, strings.Replace(crDetail(1), `"capacity":4`, `"capacity":1.5`, 1), 3}, {crGroups, crSummaries, strings.Replace(crDetail(1), `"virtualMachinesAllocated":[{`, `"virtualMachinesAllocated":[null,{`, 1), 3}, {crGroups, crSummaries, strings.Replace(crDetail(1), "/virtualMachines/v0", "/disks/v0", 1), 3}, {crGroups, crSummaries, `null`, 3}} {
		r, calls := crBodies(t, []string{f.group, f.summary, f.detail})
		status := "unknown"
		if f.calls > 1 {
			status = "partial"
		}
		if calls != f.calls || r.Evidence.Status != status || len(r.Evidence.Reservations) != 0 || r.FailureCode == "" {
			t.Fatalf("identity before downstream/prefix: %#v calls%d", r, calls)
		}
	}
	g := quotaGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		t.Fatal("invalid scope reached auth")
		return nil, nil, nil
	})
	c := crCollector(t, g)
	for _, req := range []ReservationRequest{{quotaTestID, "global"}, {"bad", "eastus"}, {quotaTestID, "EastUS"}, {"11111111-1111-1111-1111-111111111111", "eastus"}} {
		if r, e := c.Collect(context.Background(), quotaScope(), req); r != nil || e == nil {
			t.Fatalf("invalid selected request: %#v %v", r, e)
		}
	}
	if !quotaJSONWithIDs(context.Background(), []byte(`{"id":"`+strings.Repeat("a", 2048)+`"}`), true) || quotaJSONWithIDs(context.Background(), []byte(`{"id":"`+strings.Repeat("a", 2049)+`"}`), true) || quotaJSON(context.Background(), []byte(`{"id":"`+strings.Repeat("a", 513)+`"}`)) {
		t.Fatal("separate identity/quota string caps")
	}
}

func TestReservationCollectorUnicodeDuplicates(t *testing.T) {
	groupSigma := strings.Replace(crGroup, "/resourceGroups/rg/", "/resourceGroups/Σ/", 1)
	groupFinal := strings.Replace(crGroup, "/resourceGroups/rg/", "/resourceGroups/ς/", 1)
	for _, f := range []struct {
		group, summary, detail string
		calls                  int
	}{
		{`{"value":[{"id":"` + groupSigma + `","location":"eastus"},{"id":"` + groupFinal + `","location":"eastus"}]}`, crSummaries, crDetail(1), 1},
		{crGroups, `{"value":[{"name":"Σ"},{"name":"ς"}]}`, crDetail(1), 2},
		{crGroups, crSummaries, strings.ReplaceAll(strings.Replace(crDetail(2), "/virtualMachines/v0", "/virtualMachines/Σ", 1), "/virtualMachines/v1", "/virtualMachines/ς"), 3},
	} {
		r, calls := crBodies(t, []string{f.group, f.summary, f.detail})
		want := "partial"
		if f.calls == 1 {
			want = "unknown"
		}
		if calls != f.calls || r.Evidence.Status != want || r.FailureCode != "reservation_identity_invalid" || len(r.Evidence.Reservations) != 0 {
			t.Fatalf("Unicode duplicate admitted: %#v calls%d", r, calls)
		}
	}

}

func TestReservationRuntimeUnicodeDuplicate(t *testing.T) {
	groupSigma := strings.Replace(crGroup, "/resourceGroups/rg/", "/resourceGroups/Σ/", 1)
	groupFinal := strings.Replace(crGroup, "/resourceGroups/rg/", "/resourceGroups/ς/", 1)
	rows := []ReservationUsage{{ResourceID: groupSigma + "/capacityReservations/r", Region: "eastus", ReservedKnown: true, AllocatedKnown: true}, {ResourceID: groupFinal + "/capacityReservations/r", Region: "eastus", ReservedKnown: true, AllocatedKnown: true}}
	if r, e := CalculateReservations(context.Background(), quotaScope(), []ReservationRequest{crRequest()}, []ReservationEvidence{{Request: crRequest(), Status: "complete", Reservations: rows}}); r != nil || e == nil || !strings.Contains(e.Error(), "identity_duplicate") {
		t.Fatalf("calculator Unicode duplicate admitted: %#v %v", r, e)
	}
}

func TestReservationCollectorUnicodeJSONFields(t *testing.T) {
	for _, body := range []string{
		`{"sku":{"name":"SKU","capacity":4},"ſku":{"name":"SKU","capacity":9}}`,
		`{"properties":{"instanceView":{"utilizationInfo":{"virtualMachinesAllocated":[]}}},"propertieſ":{}}`,
	} {
		r, calls := crBodies(t, []string{crGroups, crSummaries, body})
		if calls != 3 || r.Evidence.Status != "partial" || r.FailureCode != "reservation_invalid_response" || len(r.Evidence.Reservations) != 0 {
			t.Fatalf("Unicode JSON field alias admitted: %#v calls%d", r, calls)
		}
		if quotaJSON(context.Background(), []byte(body)) || quotaJSONWithIDs(context.Background(), []byte(body), true) {
			t.Fatal("shared strict JSON accepted fold-duplicate keys")
		}
	}
}

func TestReservationCollectorFailurePagination(t *testing.T) {
	for _, status := range []int{400, 403, 404, 405, 429, 500, 302, 204} {
		for failAt := 0; failAt < 3; failAt++ {
			calls := 0
			bodies := []string{crGroups, crSummaries, crDetail(0)}
			c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
				i := calls
				calls++
				if i == failAt {
					return []byte(bodies[i]), &http.Response{StatusCode: status}, nil
				}
				return quotaPage(bodies[i])
			})
			r, e := c.Collect(context.Background(), quotaScope(), crRequest())
			want := "unknown"
			if failAt > 0 {
				want = "partial"
			}
			if e != nil || calls != failAt+1 || r.Evidence.Status != want || r.FailureCode != "reservation_request_failed" || len(r.Evidence.Reservations) != 0 {
				t.Fatalf("reservation denied not healthy empty: %#v %v", r, e)
			}
		}
	}
	c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(crGroups), &http.Response{StatusCode: 200}, errors.New("secret body token URL")
	})
	if r, e := c.Collect(context.Background(), quotaScope(), crRequest()); e != nil || r.Evidence.Status != "unknown" || r.FailureCode != "reservation_request_failed" || len(r.Evidence.Reservations) != 0 {
		t.Fatalf("getter error with healthy body: %#v %v", r, e)
	}
	start := crURLs()[0]
	for _, next := range []string{start, start + "&other=x", strings.Replace(start, "management.azure.com", "evil.example", 1), strings.Replace(start, "2024-11-01", "2023-01-01", 1), strings.Replace(start, quotaTestID, "11111111-1111-1111-1111-111111111111", 1)} {
		calls := 0
		c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			b, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": next})
			return quotaPage(string(b))
		})
		r, e := c.Collect(context.Background(), quotaScope(), crRequest())
		if e != nil || calls != 1 || r.Evidence.Status != "partial" || r.FailureCode != "reservation_unsafe_continuation" {
			t.Fatalf("unsafe reservation continuation reached auth: %#v %v calls%d", r, e, calls)
		}
	}
	calls := 0
	c = crCollector(t, func(_ context.Context, u string, _ int64) ([]byte, *http.Response, error) {
		i := calls
		calls++
		if i == 0 {
			return quotaPage(strings.TrimSuffix(crGroups, "}") + `,"nextLink":"?api-version=2024-11-01&$skiptoken=two"}`)
		}
		if i == 1 {
			return quotaPage(crSummaries)
		}
		if i == 2 {
			return quotaPage(crDetail(1))
		}
		if u != start[:strings.Index(start, "?")]+"?%24skiptoken=two&api-version=2024-11-01" {
			t.Fatalf("exact canonical reservation pagination: %s", u)
		}
		return nil, &http.Response{StatusCode: 403}, nil
	})
	r, e := c.Collect(context.Background(), quotaScope(), crRequest())
	if e != nil || calls != 4 || r.Evidence.Status != "partial" || len(r.Evidence.Reservations) != 1 || r.Evidence.Reservations[0].Allocated != 1 {
		t.Fatalf("later group prefix retained: %#v %v", r, e)
	}
}

// Operational failures in an inner chain must not suppress independent peers.
func TestReservationCollectorPeerRecovery(t *testing.T) {
	for _, failed := range []string{"list", "get", "later-list"} {
		t.Run(failed, func(t *testing.T) {
			otherGroup := strings.Replace(crGroup, "/capacityReservationGroups/g", "/capacityReservationGroups/h", 1)
			groups := `{"value":[{"id":"` + crGroup + `","location":"eastus"},{"id":"` + otherGroup + `","location":"eastus"}]}`
			urls := crURLs()
			secondList := "https://management.azure.com" + otherGroup + "/capacityReservations?api-version=2024-11-01"
			secondGet := "https://management.azure.com" + otherGroup + "/capacityReservations/r?%24expand=instanceView&api-version=2024-11-01"
			expected := []string{urls[0], urls[1]}
			if failed != "list" {
				expected = append(expected, urls[2])
			}
			later := strings.Split(urls[1], "?")[0] + "?%24skiptoken=two&api-version=2024-11-01"
			if failed == "later-list" {
				expected = append(expected, later)
			}
			expected = append(expected, secondList, secondGet)
			calls := 0
			c := crCollector(t, func(_ context.Context, u string, _ int64) ([]byte, *http.Response, error) {
				if calls >= len(expected) || u != expected[calls] {
					t.Fatalf("peer request%d: %s", calls, u)
				}
				calls++
				switch u {
				case urls[0]:
					return quotaPage(groups)
				case urls[1]:
					if failed == "list" {
						return nil, &http.Response{StatusCode: 403}, nil
					}
					if failed == "later-list" {
						return quotaPage(`{"value":[{"name":"r"}],"nextLink":"?api-version=2024-11-01&$skiptoken=two"}`)
					}
					return quotaPage(crSummaries)
				case urls[2]:
					if failed == "get" {
						return nil, &http.Response{StatusCode: 403}, nil
					}
					return quotaPage(crDetail(1))
				case later:
					return nil, &http.Response{StatusCode: 403}, nil
				case secondList:
					return quotaPage(crSummaries)
				case secondGet:
					return quotaPage(strings.ReplaceAll(crDetail(2), crGroup, otherGroup))
				}
				t.Fatal("unreachable request")
				return nil, nil, nil
			})
			r, e := c.Collect(context.Background(), quotaScope(), crRequest())
			rows := 1
			if failed == "later-list" {
				rows = 2
			}
			if e != nil || calls != len(expected) || r.Evidence.Status != "partial" || r.FailureCode != "reservation_request_failed" || len(r.Evidence.Reservations) != rows {
				t.Fatalf("peer recovery: %#v %v calls%d", r, e, calls)
			}
			last := r.Evidence.Reservations[rows-1]
			if last.ResourceID != otherGroup+"/capacityReservations/r" || last.Allocated != 2 {
				t.Fatalf("independent peer lost: %#v", last)
			}
		})
	}
}

func TestReservationCollectorCancellationOwnership(t *testing.T) {
	c := crCollector(t, func(_ context.Context, u string, _ int64) ([]byte, *http.Response, error) {
		if u == crURLs()[0] {
			return quotaPage(crGroups)
		}
		if u == crURLs()[1] {
			return quotaPage(crSummaries)
		}
		return quotaPage(crDetail(1))
	})
	probe := &availabilityCancelContext{Context: context.Background(), at: 100000}
	if _, e := c.Collect(probe, quotaScope(), crRequest()); e != nil {
		t.Fatal(e)
	}
	for at := 1; at <= probe.calls; at++ {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		if r, e := c.Collect(ctx, quotaScope(), crRequest()); r != nil || !errors.Is(e, context.Canceled) {
			t.Fatalf("reservation cancellation%d: %#v %v", at, r, e)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := c.Collect(context.Background(), quotaScope(), crRequest())
			if e != nil || r.Evidence.Reservations[0].Allocated != 1 {
				t.Errorf("owned concurrent result: %#v %v", r, e)
				return
			}
			r.Evidence.Reservations[0].SKU = "changed"
		}()
	}
	wg.Wait()
	c = crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return nil, nil, context.DeadlineExceeded
	})
	if r, e := c.Collect(context.Background(), quotaScope(), crRequest()); r != nil || !errors.Is(e, context.DeadlineExceeded) {
		t.Fatalf("reservation getter deadline: %#v %v", r, e)
	}
}

func TestReservationCollectorRetainedSource(t *testing.T) {
	var in struct {
		Cases map[string]struct {
			Pages []struct {
				URL    string
				Status int
				Body   string
			}
		}
	}
	b, e := os.ReadFile("testdata/source-reservation-inputs.json")
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(b, &in); e != nil {
		t.Fatal(e)
	}
	const id = "11111111-1111-1111-1111-111111111111"
	for _, name := range []string{"empty", "group-http403", "null-group", "missing-utilization", "identity-conflict"} {
		t.Run(name, func(t *testing.T) {
			f := in.Cases[name]
			calls := 0
			c := crCollector(t, func(_ context.Context, u string, n int64) ([]byte, *http.Response, error) {
				if calls >= len(f.Pages) || u != f.Pages[calls].URL || n != 1048576 {
					t.Fatalf("unchanged source request: %s", u)
				}
				p := f.Pages[calls]
				calls++
				return []byte(p.Body), &http.Response{StatusCode: p.Status}, nil
			})
			r, e := c.Collect(context.Background(), map[string]string{id: "Selected"}, ReservationRequest{id, "eastus"})
			if e != nil {
				t.Fatal(e)
			}
			switch name {
			case "empty":
				if r.Evidence.Status != "complete" || calls != 1 || len(r.Evidence.Reservations) != 0 {
					t.Fatal("actual source complete empty")
				}
			case "group-http403", "null-group":
				if r.Evidence.Status != "unknown" || calls != 1 {
					t.Fatalf("actual malformed/denied source: %#v", r)
				}
			case "missing-utilization":
				if calls != 3 || r.Evidence.Status != "complete" || len(r.Evidence.Reservations) != 1 || r.Evidence.Reservations[0].AllocatedKnown {
					t.Fatalf("source unknown-utilization correction: %#v", r)
				}
			case "identity-conflict":
				if r.Evidence.Status != "partial" || len(r.Evidence.Reservations) != 0 {
					t.Fatalf("source conflict correction: %#v", r)
				}
			}
		})
	}
}

func TestReservationCollectorAuthenticatedClient(t *testing.T) {
	closed := 0
	calls := 0
	urls := crURLs()
	bodies := []string{crGroups, crSummaries, crDetail(1)}
	client := azure.NewHTTPClient(quotaCredential{t, "https://management.usgovcloudapi.net/.default"}, &azure.HTTPClientOptions{MaxRetries: -1, Scope: "https://management.usgovcloudapi.net/.default", Transport: quotaTransport(func(r *http.Request) (*http.Response, error) {
		if calls >= 3 || r.URL.String() != strings.Replace(urls[calls], "management.azure.com", "management.usgovcloudapi.net", 1) || r.Method != "GET" || r.Body != nil || r.Header.Get("Authorization") != "Bearer synthetic-quota-token" {
			t.Fatalf("reservation sovereign authenticated boundary: %#v", r)
		}
		b := bodies[calls]
		calls++
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: crClosedBody{strings.NewReader(b), &closed}, Request: r}, nil
	})})
	c, e := NewReservationCollector("https://management.usgovcloudapi.net", client)
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.Collect(context.Background(), quotaScope(), crRequest())
	if e != nil || r.Evidence.Status != "complete" || calls != 3 || closed != 3 || len(r.Evidence.Reservations) != 1 {
		t.Fatalf("actual client/body ownership: %#v %v calls%d closed%d", r, e, calls, closed)
	}
}

type crClosedBody struct {
	*strings.Reader
	closed *int
}

func TestReservationCollectorBudgets(t *testing.T) {
	for _, pages := range []int{64, 65} {
		calls := 0
		c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			next := ""
			if calls < pages {
				next = fmt.Sprintf("?api-version=2024-11-01&$skiptoken=%d", calls)
			}
			b, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": next})
			return quotaPage(string(b))
		})
		r, e := c.Collect(context.Background(), quotaScope(), crRequest())
		status, code := "complete", ""
		if pages == 65 {
			status, code = "partial", "reservation_page_limit"
		}
		if e != nil || calls != 64 || r.Evidence.Status != status || r.FailureCode != code {
			t.Fatalf("exact list-page ceiling%d: %#v %v calls%d", pages, r, e, calls)
		}
	}
	for _, size := range []int{1048576, 1048577} {
		body := `{"value":[]}` + strings.Repeat(" ", size-len(`{"value":[]}`))
		calls := 0
		c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) { calls++; return quotaPage(body) })
		r, e := c.Collect(context.Background(), quotaScope(), crRequest())
		want := "complete"
		if size > 1048576 {
			want = "unknown"
		}
		if e != nil || calls != 1 || r.Evidence.Status != want || (size > 1048576 && r.FailureCode != "reservation_byte_limit") {
			t.Fatalf("exact page-byte ceiling%d: %#v %v", size, r, e)
		}
	}
	for _, pages := range []int{8, 9} {
		calls := 0
		c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			next := ""
			if calls < pages {
				next = fmt.Sprintf("?api-version=2024-11-01&$skiptoken=%d", calls)
			}
			b, _ := json.Marshal(map[string]any{"value": []any{}, "nextLink": next})
			size := 1048576
			if calls == 9 {
				size = len(b)
			}
			return quotaPage(string(b) + strings.Repeat(" ", size-len(b)))
		})
		r, e := c.Collect(context.Background(), quotaScope(), crRequest())
		want, code := "complete", ""
		if pages == 9 {
			want, code = "partial", "reservation_byte_limit"
		}
		if e != nil || calls != pages || r.Evidence.Status != want || r.FailureCode != code {
			t.Fatalf("aggregate exact8MiB ceiling: %#v %v calls%d", r, e, calls)
		}
	}
	for _, count := range []int{254, 255} {
		summaries := []map[string]string{}
		for i := 0; i < count; i++ {
			summaries = append(summaries, map[string]string{"name": fmt.Sprintf("r%d", i)})
		}
		b, _ := json.Marshal(map[string]any{"value": summaries})
		calls := 0
		c := crCollector(t, func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			if calls == 1 {
				return quotaPage(crGroups)
			}
			if calls == 2 {
				return quotaPage(string(b))
			}
			return quotaPage(`{"sku":{"name":"SKU","capacity":1},"properties":{"instanceView":{"utilizationInfo":{"virtualMachinesAllocated":[]}}}}`)
		})
		r, e := c.Collect(context.Background(), quotaScope(), crRequest())
		want, code := "complete", ""
		if count == 255 {
			want, code = "partial", "reservation_request_limit"
		}
		if e != nil || calls != 256 || len(r.Evidence.Reservations) != 254 || r.Evidence.Status != want || r.FailureCode != code {
			t.Fatalf("exact256 calls/prefix: %#v %v calls%d", r, e, calls)
		}
	}
	for _, count := range []int{4095, 4096} {
		bodies := []string{crGroups, `{"value":[{"name":"r"},{"name":"a"}]}`, crDetail(4094), strings.ReplaceAll(crDetail(count), `"r"`, `"a"`)}
		bodies[3] = strings.ReplaceAll(bodies[3], crResource, crGroup+"/capacityReservations/a")
		urls := append(crURLs(), "https://management.azure.com"+crGroup+"/capacityReservations/a?%24expand=instanceView&api-version=2024-11-01")
		calls := 0
		c := crCollector(t, func(_ context.Context, u string, n int64) ([]byte, *http.Response, error) {
			if calls >= len(bodies) || u != urls[calls] || n != MaxQuotaPageBytes || len(bodies[calls]) > MaxQuotaPageBytes {
				t.Fatalf("work boundary request/body %d: %s", calls, u)
			}
			body := bodies[calls]
			calls++
			return quotaPage(body)
		})
		r, e := c.Collect(context.Background(), quotaScope(), crRequest())
		want, code, rows := "complete", "", 2
		if count == 4096 {
			want, code, rows = "partial", "reservation_work_limit", 1
		}
		if e != nil || calls != 4 || r.Evidence.Status != want || r.FailureCode != code || len(r.Evidence.Reservations) != rows || r.Evidence.Reservations[0].Allocated != 4094 {
			t.Fatalf("raw group+summaries+references exact8192/prefix: %#v %v calls%d", r, e, calls)
		}
		if rows == 2 && r.Evidence.Reservations[1].Allocated != 4095 {
			t.Fatal("second reservation count lost")
		}
	}
}

func (b crClosedBody) Close() error { *b.closed++; return nil }
