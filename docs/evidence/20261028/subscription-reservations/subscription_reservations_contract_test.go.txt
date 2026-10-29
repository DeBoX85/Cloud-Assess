package region

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

const bulkSubID = "abcdefab-cdef-abcd-efab-cdefabcdefab"
const bulkForeignID = "22222222-2222-2222-2222-222222222222"
const bulkOrigin = "https://management.azure.com"
const bulkRoot = "/subscriptions/" + bulkSubID
const bulkGroupsPath = bulkRoot + "/providers/Microsoft.Compute/capacityReservationGroups"

type bulkLiteral struct {
	ID, Scenario, Status, Code string
	Calls                      int
	Regions                    []string
	Rows                       [][]any
}
type bulkRecipe struct {
	ID, Kind, Status, Code         string
	Amount, Calls, Rows, Allocated int
}
type bulkLiterals struct {
	Format, Base, Scope string
	Cases               []bulkLiteral
	Recipes             []bulkRecipe `json:"budget_recipes"`
	Procedural          []string
}

func bulkContract(t *testing.T) bulkLiterals {
	t.Helper()
	raw, e := os.ReadFile("../../../docs/evidence/20261028/subscription-reservations/literals.json")
	if e != nil {
		t.Fatal(e)
	}
	var d bulkLiterals
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if e = decoder.Decode(&d); e != nil {
		t.Fatal(e)
	}
	if d.Format != "cloud-assess-subscription-reservations-v1" || d.Base != "4627a4d9c2853d9337e8c0f959f241a6d319b416" || len(d.Cases) != 13 || len(d.Recipes) != 8 || len(d.Procedural) != 12 {
		t.Fatal("BULK_LITERAL_INVENTORY", d)
	}
	return d
}

type bulkWireGroup struct{ rg, name, region string }

func bulkID(g bulkWireGroup) string {
	return bulkRoot + "/resourceGroups/" + g.rg + "/providers/Microsoft.Compute/capacityReservationGroups/" + g.name
}
func bulkThree() []bulkWireGroup {
	return []bulkWireGroup{{"rg1", "east", "eastus"}, {"rg2", "west", "westeurope"}, {"rg3", "north", "northeurope"}}
}
func bulkPage(groups []bulkWireGroup, next string) string {
	rows := []map[string]string{}
	for _, g := range groups {
		rows = append(rows, map[string]string{"id": bulkID(g), "name": g.name, "location": g.region})
	}
	raw, _ := json.Marshal(map[string]any{"value": rows, "nextLink": next})
	return string(raw)
}
func bulkSummary(names ...string) string {
	rows := []map[string]string{}
	for _, name := range names {
		rows = append(rows, map[string]string{"name": name})
	}
	raw, _ := json.Marshal(map[string]any{"value": rows})
	return string(raw)
}
func bulkDetail(g bulkWireGroup, name, region string, allocated int) string {
	refs := []map[string]string{}
	for i := 0; i < allocated; i++ {
		refs = append(refs, map[string]string{"id": fmt.Sprintf("%s/resourceGroups/g/providers/Microsoft.Compute/virtualMachines/%d", bulkRoot, i)})
	}
	raw, _ := json.Marshal(map[string]any{"id": bulkID(g) + "/capacityReservations/" + name, "name": name, "location": region, "sku": map[string]any{"name": "S", "capacity": 4}, "properties": map[string]any{"instanceView": map[string]any{"utilizationInfo": map[string]any{"virtualMachinesAllocated": refs}}}})
	return string(raw)
}
func bulkHTTP(body string, status int) ([]byte, *http.Response, error) {
	return []byte(body), &http.Response{StatusCode: status}, nil
}
func bulkScenarioGetter(t *testing.T, scenario string) func(context.Context, string, int64) ([]byte, *http.Response, error) {
	t.Helper()
	groups := bulkThree()
	switch scenario {
	case "empty", "unsafe-next":
		groups = nil
	case "three":
	case "same-region":
		groups = groups[:2]
		groups[1].region = "eastus"
	case "child-failure":
		groups = groups[:2]
	default:
		groups = groups[:1]
	}
	return func(_ context.Context, raw string, limit int64) ([]byte, *http.Response, error) {
		u, e := url.Parse(raw)
		if e != nil || u.Scheme != "https" || u.Host != "management.azure.com" || u.Query().Get("api-version") != "2024-11-01" || limit != 1048576 {
			t.Fatalf("BULK_REQUEST_BINDING: %s %d", raw, limit)
		}
		if u.Path == bulkGroupsPath {
			if u.Query().Get("$skiptoken") == "bad" {
				return bulkHTTP(strings.ReplaceAll(bulkPage(groups, ""), bulkSubID, bulkForeignID), 200)
			}
			switch scenario {
			case "initial-failure":
				return bulkHTTP("", 500)
			case "foreign-group":
				return bulkHTTP(strings.ReplaceAll(bulkPage(groups, ""), bulkSubID, bulkForeignID), 200)
			case "duplicate-group":
				return bulkHTTP(bulkPage([]bulkWireGroup{groups[0], groups[0]}, ""), 200)
			case "unsafe-next":
				return bulkHTTP(bulkPage(nil, "https://foreign.invalid/next"), 200)
			case "later-invalid":
				return bulkHTTP(bulkPage(groups, bulkOrigin+bulkGroupsPath+"?api-version=2024-11-01&%24skiptoken=bad"), 200)
			default:
				return bulkHTTP(bulkPage(groups, ""), 200)
			}
		}
		for _, g := range groups {
			if u.Path == bulkID(g)+"/capacityReservations" {
				if u.Query().Get("$expand") != "" {
					t.Fatal("BULK_LIST_EXPAND")
				}
				switch scenario {
				case "empty-group":
					return bulkHTTP(bulkSummary(), 200)
				case "child-failure":
					if g.name == "east" {
						return bulkHTTP("", 500)
					}
				case "get-failure":
					return bulkHTTP(bulkSummary("r", "second"), 200)
				case "duplicate-summary":
					return bulkHTTP(bulkSummary("r", "r"), 200)
				}
				return bulkHTTP(bulkSummary("r"), 200)
			}
			for _, name := range []string{"r", "second"} {
				if u.Path == bulkID(g)+"/capacityReservations/"+name {
					if u.RawQuery != "%24expand=instanceView&api-version=2024-11-01" {
						t.Fatalf("BULK_DETAIL_EXPAND: %s", raw)
					}
					if scenario == "get-failure" && name == "r" {
						return bulkHTTP("", 500)
					}
					region := g.region
					if scenario == "wrong-region" {
						region = "westeurope"
					}
					return bulkHTTP(bulkDetail(g, name, region, 0), 200)
				}
			}
		}
		t.Fatalf("BULK_UNEXPECTED_REQUEST: %s", raw)
		return nil, nil, nil
	}
}
func bulkCollector(t *testing.T, getter quotaGetterFunc) *ReservationCollector {
	t.Helper()
	c, e := NewReservationCollector(bulkOrigin, getter)
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func bulkCollect(t *testing.T, scenario string) (*SubscriptionReservationCollection, int) {
	t.Helper()
	calls := 0
	get := bulkScenarioGetter(t, scenario)
	c := bulkCollector(t, func(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
		calls++
		return get(ctx, u, n)
	})
	result, e := c.CollectSubscription(context.Background(), map[string]string{bulkSubID: "Selected"}, SubscriptionReservationRequest{SubscriptionID: bulkSubID})
	if e != nil {
		t.Fatal("BULK_UNEXPECTED_ERROR", e)
	}
	return result, calls
}
func bulkRows(t *testing.T, input []ReservationUsage) [][]any {
	t.Helper()
	rows := [][]any{}
	for _, r := range input {
		p := strings.Split(strings.TrimPrefix(r.ResourceID, "/"), "/")
		if len(p) != 10 || p[1] != bulkSubID || !r.ReservedKnown || !r.AllocatedKnown || r.SKU != "S" {
			t.Fatalf("BULK_ROW_BINDING: %#v", r)
		}
		rows = append(rows, []any{r.Region, p[3], p[7], p[9], float64(r.Reserved), float64(r.Allocated)})
	}
	return rows
}
func bulkEvidenceCheck(t *testing.T, c *SubscriptionReservationCollection, status, code string, regions []string, rows [][]any) {
	t.Helper()
	if c == nil || c.SubscriptionID != bulkSubID || c.Status != status || c.FailureCode != code || !reflect.DeepEqual(c.Regions, regions) || !reflect.DeepEqual(bulkRows(t, c.Reservations), rows) {
		t.Fatalf("BULK_LITERAL_DIAGNOSTICS: %#v", c)
	}
	e := c.Evidence()
	if e.SubscriptionID != bulkSubID || e.Status != status || e.FailureCode != code || len(e.Requests) != len(regions) || len(e.Regional) != len(regions) {
		t.Fatalf("BULK_LITERAL_ENVELOPE: %#v", e)
	}
	for i, region := range regions {
		wantRequest := ReservationRequest{SubscriptionID: bulkSubID, Region: region}
		wantRows := [][]any{}
		for _, row := range rows {
			if row[0] == region {
				wantRows = append(wantRows, row)
			}
		}
		if e.Requests[i] != wantRequest || e.Regional[i].Request != wantRequest || e.Regional[i].Status != status || !reflect.DeepEqual(bulkRows(t, e.Regional[i].Reservations), wantRows) {
			t.Fatalf("BULK_REGIONAL_BINDING: %#v", e)
		}
	}
	if len(regions) == 0 && (e.Requests == nil || e.Regional == nil) {
		t.Fatal("BULK_EMPTY_ENVELOPE")
	}
}
func TestSubscriptionReservationsFrozenLiterals(t *testing.T) {
	for _, tc := range bulkContract(t).Cases {
		t.Run(tc.ID, func(t *testing.T) {
			c, calls := bulkCollect(t, tc.Scenario)
			if calls != tc.Calls {
				t.Fatalf("BULK_LITERAL_CALLS: %d want %d", calls, tc.Calls)
			}
			bulkEvidenceCheck(t, c, tc.Status, tc.Code, tc.Regions, tc.Rows)
		})
	}
}
