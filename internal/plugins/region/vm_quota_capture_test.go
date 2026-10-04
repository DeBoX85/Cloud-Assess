package region

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

type vmQuotaCapturedResult struct {
	Entries          []quotaCapturedUsage
	Risk             []string
	Failed, Panicked bool
	Requests         []string
}

func TestVMQuotaCapturedSemantics(t *testing.T) {
	var got map[string]vmQuotaCapturedResult
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-vm-quota-outputs.json"), &got); err != nil {
		t.Fatal(err)
	}
	const root = "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Compute/locations/eastus/usages?"
	const start = root + "api-version=2024-11-01"
	const next = root + "%24skiptoken=synthetic-second&api-version=2024-11-01"
	want := map[string]vmQuotaCapturedResult{
		"arithmetic": {Entries: []quotaCapturedUsage{
			{"ExactFamily", "ExactFamily", 85, 100, 15, 15, false, false},
			{"BelowFamily", "BelowFamily", 86, 100, 14, 14.000000000000002, true, false},
			{"AtFamily", "AtFamily", 100, 100, 0, 0, true, true},
			{"OverFamily", "OverFamily", 110, 100, -10, -10, true, true},
			{"NegativeFamily", "NegativeFamily", -5, 100, 105, 105, false, false},
		}, Risk: []string{"BelowFamily (86/100, 86% used)", "AtFamily (100/100, 100% used)", "OverFamily (110/100, 110% used)"}, Requests: []string{start}},
		"selection":       {Entries: []quotaCapturedUsage{{"standardDSFamily", "standardDSFamily", 10, 100, 90, 90, false, false}, {"NotFamilySuffix", "NotFamilySuffix", 10, 100, 90, 90, false, false}}, Requests: []string{start}},
		"localized":       {Entries: []quotaCapturedUsage{{"NoLabelFamily", "", 9, 10, 1, 10, true, false}, {"LabelFamily", "VM Label", 9, 10, 1, 10, true, false}}, Risk: []string{"NoLabelFamily (9/10, 90% used)", "VM Label (9/10, 90% used)"}, Requests: []string{start}},
		"paged":           {Entries: []quotaCapturedUsage{{"FirstFamily", "FirstFamily", 1, 10, 9, 90, false, false}, {"SecondFamily", "SecondFamily", 9, 10, 1, 10, true, false}}, Risk: []string{"SecondFamily (9/10, 90% used)"}, Requests: []string{start, next}},
		"empty":           {Requests: []string{start}},
		"missing-skipped": {Requests: []string{start}},
		"missing-current": {Panicked: true, Requests: []string{start}},
		"null-item":       {Panicked: true, Requests: []string{start}},
		"cancelled":       {Failed: true, Requests: []string{}},
	}
	for _, name := range []string{"http400", "http403", "http404", "http405", "malformed"} {
		want[name] = vmQuotaCapturedResult{Failed: true, Requests: []string{start}}
	}
	for _, name := range []string{"later-http400", "later-http403", "later-http404", "later-http405", "later-malformed"} {
		want[name] = vmQuotaCapturedResult{Failed: true, Requests: []string{start, next}}
	}
	if len(want) != 19 || !reflect.DeepEqual(got, want) {
		t.Fatal("complete independent pinned VM quota request/row/risk/panic/failure literals differ")
	}
	var input struct {
		SourceBlob, SDKModule string
		DisableRPRegistration bool
		MaxRetries            int
		Cases                 map[string]json.RawMessage
	}
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-vm-quota-inputs.json"), &input); err != nil {
		t.Fatal(err)
	}
	if input.SourceBlob != "9d98d3d2c221b925fd5b4bbde4ddf78476146491" || input.SDKModule != "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6 v6.4.0" || !input.DisableRPRegistration || input.MaxRetries != -1 || len(input.Cases) != 19 {
		t.Fatal("pinned SDK source/fixture policy premises differ")
	}
}

func TestVMQuotaCapturedHashes(t *testing.T) {
	for name, want := range map[string]string{
		"source-vm-quota-inputs.json":  "5b5319746110e4e86cc6753ec5b20a3306d6290d03b0e829682d2511324d3a1c",
		"source-vm-quota-outputs.json": "b6311db11c7a62c74818ef694c2991cd1ab46d3a8c0bf1337ed5fcd7cb480c32",
	} {
		raw := availabilityCapturedBytes(t, name)
		hash := sha256.Sum256(raw)
		if hex.EncodeToString(hash[:]) != want || !json.Valid(raw) {
			t.Fatal("pinned VM quota capture bytes changed", name)
		}
	}
}
