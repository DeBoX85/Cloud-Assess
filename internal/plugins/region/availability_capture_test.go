package region

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"testing"
)

type availabilityCapturedComparison struct {
	Comparison
	Score float64
}

type availabilityCapturedOutput struct {
	Cases          map[string]availabilityCapturedComparison `json:"cases"`
	ProviderCalls  map[string]int                            `json:"provider_calls"`
	ProviderLookup []struct {
		ResourceType string
		Region       string
		Available    bool
	} `json:"provider_lookup"`
}

func availabilityCapturedBytes(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestAvailabilityCaptureHashes(t *testing.T) {
	files := map[string]string{
		"source-availability-inputs.json":  "32345bea44a162a1518f9022c93dba3cf5640980c94fb1d3734e172b780cc245",
		"source-availability-outputs.json": "d3e8186fb80ffc9a47191228e5cad222636ffdfea1de3c837747bf896acbaef9",
	}
	for name, want := range files {
		raw := availabilityCapturedBytes(t, name)
		got := sha256.Sum256(raw)
		if hex.EncodeToString(got[:]) != want || !json.Valid(raw) {
			t.Fatal("independent availability capture bytes changed", name)
		}
	}
}

func TestAvailabilityCapturedSemantics(t *testing.T) {
	var got availabilityCapturedOutput
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-availability-outputs.json"), &got); err != nil {
		t.Fatal(err)
	}
	base := func(source, target string) availabilityCapturedComparison {
		return availabilityCapturedComparison{Comparison: Comparison{SourceRegion: source, TargetRegion: target, MissingResourceTypes: []string{}, MissingSKUs: []string{}, RestrictedSKUs: []string{}, ZoneRestrictedSKUs: []string{}, SourceZoneCount: 3, TargetZoneCount: 2}}
	}
	empty := base("eastus", "westeurope")
	caseSource := base("East US", "westeurope")
	caseSource.SourceZoneCount = 0
	mixed := base("eastus", "westeurope")
	mixed.SourceResourceTypeCount, mixed.AvailableTypes, mixed.UnavailableTypes, mixed.AvailabilityPercent = 4, 2, 2, 50
	mixed.MissingResourceTypes = []string{"microsoft.custom/widgets", "microsoft.network/virtualnetworks"}
	noProviders := base("eastus", "westeurope")
	noProviders.SourceResourceTypeCount, noProviders.UnavailableTypes = 4, 4
	noProviders.MissingResourceTypes = []string{"microsoft.compute/virtualmachines", "microsoft.custom/widgets", "microsoft.network/privatednszones", "microsoft.network/virtualnetworks"}
	caseTarget := base("eastus", "West Europe")
	caseTarget.SourceResourceTypeCount, caseTarget.AvailableTypes, caseTarget.UnavailableTypes, caseTarget.AvailabilityPercent, caseTarget.TargetZoneCount = 4, 1, 3, 25, 0
	caseTarget.MissingResourceTypes = []string{"microsoft.compute/virtualmachines", "microsoft.custom/widgets", "microsoft.network/virtualnetworks"}
	disabled := base("eastus", "westeurope")
	disabled.SourceResourceTypeCount, disabled.AvailableTypes, disabled.UnavailableTypes, disabled.AvailabilityPercent = 2, 1, 1, 50
	disabled.MissingResourceTypes = []string{"microsoft.custom/widgets"}
	all := disabled
	all.TotalSKUsChecked, all.AvailableSKUs, all.UnavailableSKUs, all.SKUAvailabilityPercent = 8, 2, 3, 25
	all.MissingSKUs = []string{"microsoft.capture/available:absent", "microsoft.capture/available:future", "microsoft.capture/available:missing"}
	all.RestrictedSKUs = []string{"microsoft.capture/available:restricted"}
	all.ZoneRestrictedSKUs = []string{"microsoft.capture/available:zones (zones blocked: 3,1,3)", "microsoft.capture/available:zones-empty"}
	unknown := base("eastus", "westeurope")
	unknown.SourceResourceTypeCount, unknown.AvailableTypes, unknown.AvailabilityPercent = 1, 1, 100
	unknown.TotalSKUsChecked, unknown.UnknownSKUs, unknown.SKUAvailabilityPercent = 2, 2, 100
	unknown.MissingSKUs = []string{"microsoft.capture/failed:First (unknown)", "microsoft.capture/failed:Second (unknown)"}
	partial := base("eastus", "westeurope")
	partial.SourceResourceTypeCount, partial.AvailableTypes, partial.AvailabilityPercent = 2, 2, 100
	partial.TotalSKUsChecked, partial.AvailableSKUs, partial.UnknownSKUs, partial.SKUAvailabilityPercent = 3, 1, 1, 50
	partial.MissingSKUs = []string{"microsoft.capture/failed:Unknown (unknown)"}
	partial.RestrictedSKUs = []string{"microsoft.capture/available:restricted"}
	want := map[string]availabilityCapturedComparison{
		"cancelled-source-ignores-context": partial,
		"resource-empty-providers":         noProviders,
		"resource-mixed-no-sku":            mixed,
		"sku-all-states":                   all,
		"sku-all-unknown":                  unknown,
		"sku-disabled":                     disabled,
		"sku-mixed-unknown":                partial,
		"source-absent":                    empty,
		"source-key-exact-case":            caseSource,
		"source-present-empty":             empty,
		"target-key-exact-case":            caseTarget,
	}
	if !reflect.DeepEqual(got.Cases, want) {
		t.Fatal("complete literal source availability comparisons changed")
	}
	if !reflect.DeepEqual(got.ProviderCalls, map[string]int{"available": 3, "failed": 3}) {
		t.Fatal("pure provider/cache call topology changed")
	}
	wantTypes := []string{"microsoft.compute/virtualmachines", "microsoft.compute/virtualmachines", "Microsoft.Compute/virtualMachines", "microsoft.network/privatednszones", "microsoft.network/privatednszones", "microsoft.custom/widgets", "malformed", "microsoft.compute/virtualmachines/child"}
	wantRegions := []string{"westeurope", "West Europe", "westeurope", "private-region", "", "westeurope", "westeurope", "westeurope"}
	wantAvailable := []bool{true, false, false, true, true, false, false, false}
	if len(got.ProviderLookup) != 8 {
		t.Fatal("provider lookup oracle topology changed")
	}
	for i, probe := range got.ProviderLookup {
		if probe.ResourceType != wantTypes[i] || probe.Region != wantRegions[i] || probe.Available != wantAvailable[i] {
			t.Fatal("literal exact-key/global/malformed provider lookup changed", i)
		}
	}
	var inputs struct {
		Cases map[string]struct {
			Source     string
			Target     string
			Inventory  InventoryCounts
			SKUEnabled bool
			Cancelled  bool
		} `json:"cases"`
		SKUResponse map[string]struct {
			State        int
			BlockedZones []string
		} `json:"sku_response"`
		SKUError bool `json:"sku_error"`
	}
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-availability-inputs.json"), &inputs); err != nil {
		t.Fatal(err)
	}
	if len(inputs.Cases) != 11 || len(inputs.SKUResponse) != 6 || !inputs.SKUError || !inputs.Cases["cancelled-source-ignores-context"].Cancelled || !inputs.Cases["sku-all-states"].SKUEnabled || inputs.Cases["sku-disabled"].SKUEnabled {
		t.Fatal("literal availability input topology changed")
	}
	counts := inputs.Cases["resource-mixed-no-sku"].Inventory.ResourceTypesByRegion["eastus"]
	raw := inputs.Cases["sku-all-states"].Inventory.SKUsByTypeAndRegion["microsoft.capture/available"]["eastus"]
	if counts["microsoft.compute/virtualmachines"] != 900 || counts["microsoft.network/privatednszones"] != 0 || len(counts) != 4 || len(raw) != 8 || raw["Ready"] != 100 || raw[" ready "] != 1 {
		t.Fatal("source distinct-key/resource-count/raw-SKU input premise changed")
	}
	if inputs.SKUResponse[" Ready "].State != 0 || inputs.SKUResponse["future"].State != 99 || !slices.Equal(inputs.SKUResponse["zones"].BlockedZones, []string{"3", "1", "3"}) {
		t.Fatal("raw response keys/future-state/blocked-zone input changed")
	}
}
