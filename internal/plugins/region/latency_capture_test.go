package region

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

type latencyCapturedCase struct {
	Comparisons []availabilityCapturedComparison `json:"comparisons"`
	Averages    map[string]float64               `json:"averages"`
}

type latencyCapturedInput struct {
	Matrix  map[string]map[string]float64
	Results []availabilityCapturedComparison
}

func literalLatencyComparison(source, target string, ms float64, estimated bool) availabilityCapturedComparison {
	return availabilityCapturedComparison{Comparison: Comparison{SubscriptionID: availabilityID, SubscriptionName: "selected", SourceRegion: source, TargetRegion: target, SourceResourceTypeCount: 2, AvailableTypes: 1, UnavailableTypes: 1, AvailabilityPercent: 50, TotalSKUsChecked: 1, UnknownSKUs: 1, SKUAvailabilityPercent: 100, MissingResourceTypes: []string{"microsoft.custom/widgets"}, MissingSKUs: []string{"microsoft.custom/widgets:SKU (unknown)"}, RestrictedSKUs: []string{}, ZoneRestrictedSKUs: []string{}, SourceZoneCount: 3, TargetZoneCount: 2, TargetZoneMappings: map[string]string{"1": "phys1"}, AvgLatencyMs: ms, LatencyEstimated: estimated, HasCostData: true, AvgCostDifference: -5}}
}

func TestLatencyCapturedHashes(t *testing.T) {
	files := map[string]string{
		"source-latency-inputs.json":  "1651ac72d681e39dbf43429ba8c44a9bd4d18f5e0cb801db0af580fa4677a6b7",
		"source-latency-outputs.json": "791c85899fcfa2598ac8de3860194126b232ecd7134f30558219e7bb827f1a5e",
		"source-latency-data.json":    "45e575040812ee74e006be34623cb253df24e6558f07727cefcef7e55e13326b",
	}
	for name, want := range files {
		raw := availabilityCapturedBytes(t, name)
		hash := sha256.Sum256(raw)
		if hex.EncodeToString(hash[:]) != want || !json.Valid(raw) {
			t.Fatal("pinned latency capture bytes changed", name)
		}
	}
}

func TestLatencyCapturedSemantics(t *testing.T) {
	var got map[string]latencyCapturedCase
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-latency-outputs.json"), &got); err != nil {
		t.Fatal(err)
	}
	c := literalLatencyComparison
	want := map[string]latencyCapturedCase{

		"synthetic": {Comparisons: []availabilityCapturedComparison{c("eastus", "westeurope", 85, false), c("westeurope", "eastus", 115, false), c("northeurope", "eastus", 105, false), c("westus", "francesouth", 95, true), c("francesouth", "westus", 115, true), c("northeurope", "francesouth", 30, true), c("westus", "westus2", 999, true), c("eastus", "eastus", 0, false), c("unknown", "westeurope", 1, false), c("westeurope", "unknown", 1, false), c("global", "francesouth", 0, false), c("East US", "West Europe", 85, false)}, Averages: map[string]float64{"americas:americas": 999, "americas:europe": 95, "europe:americas": 115, "europe:europe": 30}},

		"empty": {Comparisons: []availabilityCapturedComparison{c("eastus", "eastus", 0, false), c("eastus", "westeurope", 0, false)}, Averages: map[string]float64{}},

		"empty-results": {Comparisons: []availabilityCapturedComparison{}, Averages: map[string]float64{}},

		"one-way-cluster": {Comparisons: []availabilityCapturedComparison{c("westus", "francesouth", 80, true), c("francesouth", "westus", 0, false)}, Averages: map[string]float64{"americas:europe": 80}},

		"zero": {Comparisons: []availabilityCapturedComparison{c("eastus", "westeurope", 0, false), c("westus", "francesouth", 0, false)}, Averages: map[string]float64{"americas:europe": 0}},

		"negative": {Comparisons: []availabilityCapturedComparison{c("eastus", "westeurope", -2, false), c("westus", "francesouth", 0, false)}, Averages: map[string]float64{"americas:europe": -2}},

		"default": {Comparisons: []availabilityCapturedComparison{c("eastus", "westeurope", 85, false), c("austriaeast", "eastus", 125.54871794871795, true), c("eastus", "austriaeast", 126.02051282051282, true), c("qatarcentral", "westeurope", 121.06666666666666, true), c("unknown", "westeurope", 0, false), c("eastus", "eastus", 0, false), c("australiacentral", "australiacentral2", 3, false)}, Averages: map[string]float64{"americas:americas": 59.67307692307692, "americas:apac": 189.80769230769232, "americas:europe": 126.02051282051282, "americas:mea": 213.64615384615385, "apac:americas": 191.9025641025641, "apac:apac": 86.28061224489795, "apac:europe": 215.8177777777778, "apac:mea": 189.06666666666666, "europe:americas": 125.54871794871795, "europe:apac": 218.59047619047618, "europe:europe": 23.99047619047619, "europe:mea": 121.41333333333333, "mea:americas": 212.72307692307692, "mea:apac": 191.31428571428572, "mea:europe": 121.06666666666666, "mea:mea": 110.1}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatal("complete independent source latency comparison/average literals changed")
	}
	var input map[string]latencyCapturedInput
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-latency-inputs.json"), &input); err != nil {
		t.Fatal(err)
	}
	if len(input) != 7 || len(input["synthetic"].Results) != 12 || len(input["default"].Results) != 7 || len(input["empty-results"].Results) != 0 {
		t.Fatal("source input topology changed")
	}
	branches := 0
	for _, fixture := range input {
		for _, before := range fixture.Results {
			if before.AvgLatencyMs != 42 || !before.LatencyEstimated {
				t.Fatal("overwrite premise changed")
			}
			branches++
		}
	}
	if branches != 27 || input["synthetic"].Matrix["eastus"]["eastus"] != 999 || input["negative"].Matrix["eastus"]["westeurope"] != -2 {
		t.Fatal("source same-region/negative arithmetic premise changed")
	}
	var data struct {
		SourceBlob string                        `json:"source_blob"`
		Matrix     map[string]map[string]float64 `json:"matrix"`
		Clusters   map[string]string             `json:"clusters"`
	}
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-latency-data.json"), &data); err != nil {
		t.Fatal(err)
	}
	cells := 0
	for _, row := range data.Matrix {
		cells += len(row)
	}
	if data.SourceBlob != "64a6abc72e2a54fec286cb27123175d1d6bcda61" || len(data.Matrix) != 49 || cells != 2218 || len(data.Clusters) != 57 || !reflect.DeepEqual(data.Matrix, input["default"].Matrix) || data.Clusters["austriaeast"] != "europe" || data.Clusters["qatarcentral"] != "mea" {
		t.Fatal("complete pinned data topology/provenance changed")
	}
}
