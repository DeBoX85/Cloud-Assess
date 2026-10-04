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

func TestInventoryCalculationCaptureHashesAndTopology(t *testing.T) {
	files := map[string]string{
		"source-inventory-inputs.json":  "cdaa6b25e3250e6d701999279f96d5ea225c927e1249ee43c32cfe0b5072b17c",
		"source-inventory-outputs.json": "756647e2981f8cd0ac4f0d5091ba6b4a74345ad6266de1bc15f8c833afeb2006",
	}
	decoded := map[string]map[string]json.RawMessage{}
	for file, want := range files {
		raw, err := os.ReadFile("testdata/" + file)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(raw)
		if hex.EncodeToString(digest[:]) != want {
			t.Fatal("independent captured inventory bytes changed", file)
		}
		var data map[string]json.RawMessage
		if err := json.Unmarshal(raw, &data); err != nil {
			t.Fatal(err)
		}
		decoded[file] = data
	}
	inputs, outputs := decoded["source-inventory-inputs.json"], decoded["source-inventory-outputs.json"]
	var resources []json.RawMessage
	var selectors, regions []string
	for _, item := range []struct {
		name  string
		value any
	}{{"resources", &resources}, {"selectors", &selectors}, {"regions", &regions}} {
		if err := json.Unmarshal(inputs[item.name], item.value); err != nil {
			t.Fatal(err)
		}
	}
	if len(inputs) != 3 || len(resources) != 10 || len(selectors) != 5 || len(regions) != 30 || selectors[2] != "abcdefab-1111-1111-1111-111111111111" || selectors[3] != "ABCDEFAB-1111-1111-1111-111111111111" {
		t.Fatal("literal source inventory input topology changed")
	}
	names := make([]string, 0, len(outputs))
	for name := range outputs {
		names = append(names, name)
	}
	slices.Sort(names)
	wantNames := []string{"build-empty", "build-nil", "build-nil-resource-panics", "build-selected-0", "build-selected-1", "build-selected-2", "build-selected-3", "build-selected-4", "merge-empty", "merge-nil-source-panics", "merge-repeated", "merge-selected", "normalization"}
	if !slices.Equal(names, wantNames) {
		t.Fatal("literal source inventory output branches changed")
	}
	for _, name := range []string{"build-empty", "build-nil", "build-selected-0", "build-selected-1", "build-selected-2", "build-selected-3", "build-selected-4", "merge-empty", "merge-repeated", "merge-selected"} {
		var maps map[string]json.RawMessage
		if err := json.Unmarshal(outputs[name], &maps); err != nil {
			t.Fatal(err)
		}
		keys := make([]string, 0, len(maps))
		for key, raw := range maps {
			keys = append(keys, key)
			if len(raw) == 0 || raw[0] != '{' {
				t.Fatal("source inventory map became nil or non-object", name, key)
			}
			if name == "build-empty" || name == "build-nil" || name == "build-selected-2" || name == "build-selected-4" || name == "merge-empty" {
				if string(raw) != "{}" {
					t.Fatal("source empty/case-sensitive inventory branch changed", name, key)
				}
			}
		}
		slices.Sort(keys)
		if !slices.Equal(keys, []string{"LocationCounts", "ResourceTypes", "ResourceTypesByRegion", "SKUsByType", "SKUsByTypeAndRegion"}) {
			t.Fatal("source five-map shape changed", name)
		}
	}
	for _, name := range []string{"build-nil-resource-panics", "merge-nil-source-panics"} {
		if string(outputs[name]) != "true" {
			t.Fatal("observed source nil panic branch changed", name)
		}
	}
	var one map[string]json.RawMessage
	if err := json.Unmarshal(outputs["build-selected-0"], &one); err != nil {
		t.Fatal(err)
	}
	var counts map[string]int
	if err := json.Unmarshal(one["LocationCounts"], &counts); err != nil || !reflect.DeepEqual(counts, map[string]int{"": 1, "east\tus": 1, "eastus": 3, "europe": 1, "westeurope": 1}) {
		t.Fatal("source logical/empty/tab/space location counts changed", err)
	}
	// Decode independently into fresh state; the prior map must not retain keys.
	counts = nil
	if err := json.Unmarshal(one["ResourceTypes"], &counts); err != nil || !reflect.DeepEqual(counts, map[string]int{"microsoft.compute/virtualmachines": 4, "microsoft.compute/virtualmachinescalesets": 1, "microsoft.custom/widgets": 1, "microsoft.network/virtualnetworks": 1}) {
		t.Fatal("source resource rather than capacity counts changed", err)
	}
	var normalized []map[string]any
	if err := json.Unmarshal(outputs["normalization"], &normalized); err != nil || len(normalized) != 30 {
		t.Fatal("source normalization topology changed", err)
	}
	wantNormalized := []string{"eastus", "eastus", "east\tus", "east\nus", "east-us", "eastus", "east\u00a0us", "eastkus", "", "unassigned", "global", "global", "europe", "unitedstates", "asia", "asiapacific", "australia", "brazil", "canada", "france", "germany", "india", "japan", "korea", "norway", "southafrica", "switzerland", "uae", "uk", "private-region"}
	for i, record := range normalized {
		if record["region"] != regions[i] || record["normalized"] != wantNormalized[i] || record["physical"] != (i < 8 || i == 29) {
			t.Fatal("source ASCII-only spaces/Unicode/meta/unknown predicate changed", i)
		}
	}
}
