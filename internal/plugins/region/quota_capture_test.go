package region

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

type quotaCapturedUsage struct {
	ResourceName, LocalizedName    string
	CurrentValue, Limit, Available int
	HeadroomPct                    float64
	IsNearLimit, IsAtOrOverLimit   bool
}
type quotaCapturedResult struct {
	Entries  []quotaCapturedUsage
	Risk     []string
	Failed   bool
	Requests []string
}

func TestQuotaCapturedSemantics(t *testing.T) {
	var got map[string]quotaCapturedResult
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-quota-outputs.json"), &got); err != nil {
		t.Fatal(err)
	}
	const root = "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/"
	web := root + "Microsoft.Web/locations/eastus/usages?api-version=2023-01-01"
	network := root + "Microsoft.Network/locations/eastus/usages?api-version=2022-07-01"
	storage := root + "Microsoft.Storage/locations/eastus/usages?api-version=2023-01-01"
	sql := root + "Microsoft.Sql/locations/eastus/usages?api-version=2021-11-01"
	want := map[string]quotaCapturedResult{
		"arithmetic": {Entries: []quotaCapturedUsage{
			{"exact15", "exact15", 85, 100, 15, 15, false, false},
			{"below15", "below15", 86, 100, 14, 14.000000000000002, true, false},
			{"at", "at", 100, 100, 0, 0, true, true},
			{"over", "over", 110, 100, -10, -10, true, true},
			{"negative", "negative", -5, 100, 105, 105, false, false},
		}, Risk: []string{"below15 (86/100, 86% used)", "at (100/100, 100% used)", "over (110/100, 110% used)"}, Requests: []string{web}},
		"names": {Entries: []quotaCapturedUsage{
			{"Servers", "SQL Servers", 9, 10, 1, 10, true, false},
			{"ElasticPools", "ElasticPools", 1, 10, 9, 90, false, false},
			{"", "", 1, 10, 9, 90, false, false},
		}, Risk: []string{"SQL Servers (9/10, 90% used)"}, Requests: []string{sql}},
		"paged": {Entries: []quotaCapturedUsage{
			{"first", "first", 1, 10, 9, 90, false, false},
			{"second", "second", 9, 10, 1, 10, true, false},
		}, Risk: []string{"second (9/10, 90% used)"}, Requests: []string{web, web + "&$skiptoken=synthetic-second"}},
		"empty":          {Entries: []quotaCapturedUsage{}, Requests: []string{web}},
		"cancelled":      {Failed: true, Requests: []string{}},
		"web-filter":     {Entries: []quotaCapturedUsage{{"customdomains", "customdomains", 1, 10, 9, 90, false, false}, {"Workers", "Workers", 1, 10, 9, 90, false, false}}, Requests: []string{web}},
		"network-filter": {Entries: []quotaCapturedUsage{{"networkwatchers", "networkwatchers", 1, 10, 9, 90, false, false}, {"PublicIPAddresses", "PublicIPAddresses", 1, 10, 9, 90, false, false}}, Requests: []string{network}},
		"storage-filter": {Entries: []quotaCapturedUsage{{"totalblobs", "totalblobs", 1, 10, 9, 90, false, false}, {"StorageAccounts", "StorageAccounts", 1, 10, 9, 90, false, false}}, Requests: []string{storage}},
		"sql-filter":     {Entries: []quotaCapturedUsage{{"dtusPerserver", "dtusPerserver", 1, 10, 9, 90, false, false}, {"Servers", "Servers", 1, 10, 9, 90, false, false}}, Requests: []string{sql}},
	}
	for _, name := range []string{"http400", "http404", "http405", "http403", "malformed", "invalid-name"} {
		want[name] = quotaCapturedResult{Failed: name == "http403" || name == "malformed" || name == "invalid-name", Requests: []string{web}}
	}
	for _, name := range []string{"later-http400", "later-http404", "later-http405", "later-http403", "later-malformed"} {
		want[name] = quotaCapturedResult{Failed: name == "later-http403" || name == "later-malformed", Requests: []string{web, web + "&$skiptoken=synthetic-second"}}
	}
	if len(want) != 20 || !reflect.DeepEqual(got, want) {
		t.Fatal("complete independent pinned quota request/row/risk/failure literals differ")
	}
}

func TestQuotaCapturedHashes(t *testing.T) {
	for name, want := range map[string]string{
		"source-quota-inputs.json":  "804a2e6f24536be69b290be9f9d7bfc0015646ee0fc99132bcc646893150918c",
		"source-quota-outputs.json": "0796a14301763d7a67c9d47fef0d8827b10aca8844932ee047eddd71dff1d0ed",
	} {
		raw := availabilityCapturedBytes(t, name)
		hash := sha256.Sum256(raw)
		if hex.EncodeToString(hash[:]) != want || !json.Valid(raw) {
			t.Fatal("pinned quota capture bytes changed", name)
		}
	}
}
