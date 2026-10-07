package region

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

type capturedReservationEntry struct {
	SubscriptionID, SubscriptionName, ResourceGroup, CRGName, ReservationName, Location, SKU string
	Reserved, Allocated, Available                                                           int
	Status                                                                                   string
}
type capturedReservationResult struct {
	Entries          []capturedReservationEntry
	Failed, Panicked bool
	Requests         []string
}

// Complete literal expectations are independently specified from reviewed source
// and SDK request contracts, not generated from the captured outputs.
func TestReservationCapturedSemantics(t *testing.T) {
	raw, err := os.ReadFile("testdata/source-reservation-outputs.json")
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]capturedReservationResult
	if err = json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	const sub = "11111111-1111-1111-1111-111111111111"
	const root = "https://management.azure.com/subscriptions/" + sub
	const groups = root + "/providers/Microsoft.Compute/capacityReservationGroups?api-version=2024-11-01"
	const list = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g/capacityReservations?api-version=2024-11-01"
	const other = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/other/capacityReservations?api-version=2024-11-01"
	const get = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g/capacityReservations/r?%24expand=instanceView&api-version=2024-11-01"
	const bad = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g/capacityReservations/bad?%24expand=instanceView&api-version=2024-11-01"
	const otherGet = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/other/capacityReservations/r?%24expand=instanceView&api-version=2024-11-01"
	const secondGet = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g/capacityReservations/second?%24expand=instanceView&api-version=2024-11-01"
	const nextGroup = root + "/providers/Microsoft.Compute/capacityReservationGroups?%24skiptoken=second&api-version=2024-11-01"
	const nextList = root + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g/capacityReservations?%24skiptoken=second&api-version=2024-11-01"
	entry := func(group, name, location, sku string, reserved, allocated, available int, status string) capturedReservationEntry {
		return capturedReservationEntry{sub, "Synthetic Subscription", "rg", group, name, location, sku, reserved, allocated, available, status}
	}
	row := func(name string, reserved, allocated, available int, status string) capturedReservationEntry {
		return entry("g", name, "eastus", "Standard_D2s_v5", reserved, allocated, available, status)
	}
	result := func(e []capturedReservationEntry, failed, panicked bool, requests ...string) capturedReservationResult {
		return capturedReservationResult{e, failed, panicked, requests}
	}
	one := func(e capturedReservationEntry) capturedReservationResult {
		return result([]capturedReservationEntry{e}, false, false, groups, list, get)
	}
	available := row("r", 4, 1, 3, "Available")
	want := map[string]capturedReservationResult{
		"idle": one(row("r", 4, 0, 4, "Idle")), "available": one(available), "at": one(row("r", 4, 4, 0, "At-Capacity")), "over": one(row("r", 4, 5, -1, "Over-Allocated")),
		"zero": one(row("r", 0, 0, 0, "Idle")), "negative": one(row("r", -1, 0, -1, "Idle")),
		"missing-utilization": one(row("r", 4, 0, 4, "Idle")), "missing-capacity": one(row("r", 0, 1, -1, "Over-Allocated")),
		"missing-get-fields": one(entry("g", "", "eastus", "", 0, 0, 0, "Idle")), "identity-conflict": one(row("returned-name", 4, 1, 3, "Available")),
		"missing-group-location": one(entry("g", "r", "", "Standard_D2s_v5", 4, 1, 3, "Available")),
		"empty":                  result(nil, false, false, groups), "skipped-groups": result(nil, false, false, groups), "skipped-summary": result(nil, false, false, groups, list),
		"malformed-group-id": result(nil, false, false, groups), "null-group": result(nil, false, true, groups), "null-summary": result(nil, false, true, groups, list),
		"cancelled":   {Entries: nil, Failed: true, Panicked: false, Requests: []string{}},
		"cancel-list": result(nil, false, false, groups, list), "cancel-get": result(nil, false, false, groups, list, get),
		"group-http403": result(nil, true, false, groups), "group-http404": result(nil, true, false, groups), "group-malformed": result(nil, true, false, groups),
		"list-http403":   result([]capturedReservationEntry{entry("other", "r", "eastus", "Standard_D2s_v5", 4, 1, 3, "Available")}, false, false, groups, list, other, otherGet),
		"list-http404":   result([]capturedReservationEntry{entry("other", "r", "eastus", "Standard_D2s_v5", 4, 1, 3, "Available")}, false, false, groups, list, other, otherGet),
		"get-http403":    result([]capturedReservationEntry{available}, false, false, groups, list, bad, get),
		"get-http404":    result([]capturedReservationEntry{available}, false, false, groups, list, bad, get),
		"list-malformed": result(nil, false, false, groups, list), "get-malformed": result(nil, false, false, groups, list, get),
		"group-paged":         result([]capturedReservationEntry{available}, false, false, groups, list, get, nextGroup, other),
		"later-group-http403": result(nil, true, false, groups, list, get, nextGroup), "later-group-malformed": result(nil, true, false, groups, list, get, nextGroup),
		"list-paged":           result([]capturedReservationEntry{available, row("second", 4, 4, 0, "At-Capacity")}, false, false, groups, list, get, nextList, secondGet, other),
		"later-list-http403":   result([]capturedReservationEntry{available}, false, false, groups, list, get, nextList, other),
		"later-list-malformed": result([]capturedReservationEntry{available}, false, false, groups, list, get, nextList, other),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("complete reservation source records differ\ngot: %#v\nwant: %#v", got, want)
	}
}

func TestReservationCaptureHashesAndPremises(t *testing.T) {
	for name, hash := range map[string]string{
		"source-reservation-inputs.json":  "3e8e9870bb4a4d4a91490ac400dccfc28e7070ca4b8a2d22434176d1c61b3c09",
		"source-reservation-outputs.json": "4393a4ae45d0b5bacd116aa48bc44ed5698a4cde5b991b2e6e6ab6429974c2d5",
	} {
		raw, e := os.ReadFile("testdata/" + name)
		if e != nil {
			t.Fatal(e)
		}
		if fmt.Sprintf("%x", sha256.Sum256(raw)) != hash {
			t.Fatal("reservation capture byte provenance", name)
		}
	}
	raw, e := os.ReadFile("testdata/source-reservation-inputs.json")
	if e != nil {
		t.Fatal(e)
	}
	var p struct {
		SourceBlob, SDKModule string
		DisableRPRegistration bool
		MaxRetries            int
		Cases                 map[string]json.RawMessage
	}
	if e = json.Unmarshal(raw, &p); e != nil {
		t.Fatal(e)
	}
	if p.SourceBlob != "1ffb6e37cc1c21f6d8b5c9ce483f793387e0e6bc" || p.SDKModule != "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6 v6.4.0" || !p.DisableRPRegistration || p.MaxRetries != -1 || len(p.Cases) != 35 {
		t.Fatal("reservation fixture premises")
	}
}
