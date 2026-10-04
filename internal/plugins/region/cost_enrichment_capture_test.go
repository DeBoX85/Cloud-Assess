package region

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

func literalCostEnrichmentComparison(source, target string, percent float64, known bool) availabilityCapturedComparison {
	c := literalLatencyComparison(source, target, 42, true)
	c.AvgCostDifference, c.HasCostData = percent, known
	return c
}

func TestCostEnrichmentCapturedSemantics(t *testing.T) {
	var got map[string][]availabilityCapturedComparison
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-cost-enrichment-outputs.json"), &got); err != nil {
		t.Fatal(err)
	}
	c := literalCostEnrichmentComparison
	want := map[string][]availabilityCapturedComparison{}
	for name, percent := range map[string]float64{"owned-only": 100, "weighted": 62.5, "zero-weight": -12.5, "near-zero-both": 25, "near-zero-source": 100, "exact-threshold": 100, "free-target": -100, "both-free": 0, "equal-prices": 0, "negative-prices": 0, "duplicate-last-one": 50, "duplicate-last-nine": 90} {
		want[name] = []availabilityCapturedComparison{c("eastus", "westeurope", percent, true)}
	}
	for _, name := range []string{"missing-source", "missing-target", "empty-pricing", "empty-meters", "unowned-only", "negative-weight", "nil-shared"} {
		want[name] = []availabilityCapturedComparison{c("eastus", "westeurope", -5, true)}
	}
	want["logical-source"] = []availabilityCapturedComparison{c("global", "westeurope", -5, true)}
	want["logical-target"] = []availabilityCapturedComparison{c("eastus", "europe", -5, true)}
	want["display-names"] = []availabilityCapturedComparison{c("East US", "West Europe", 100, true)}
	want["same-region"] = []availabilityCapturedComparison{c("eastus", "eastus", 0, true)}
	want["nil-without-prior"] = []availabilityCapturedComparison{c("eastus", "westeurope", 0, false)}
	want["empty-results"] = []availabilityCapturedComparison{}
	other := c("eastus", "westeurope", 100, true)
	other.SubscriptionID = "22222222-2222-2222-2222-222222222222"
	want["unbound-subscription"] = []availabilityCapturedComparison{c("eastus", "westeurope", 100, true), other}
	if len(want) != 26 || !reflect.DeepEqual(got, want) {
		t.Fatal("complete independent source cost comparison literals differ")
	}
	var input struct {
		SourceBlob string `json:"source_blob"`

		TypesBlob string `json:"types_blob"`

		Cases map[string]struct {
			Results []availabilityCapturedComparison

			Meters []struct {
				MeterID string

				HistoricalCost float64
			}

			Shared *struct {
				RegionPricing map[string]map[string]float64
			}
		} `json:"cases"`
	}
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-cost-enrichment-inputs.json"), &input); err != nil {
		t.Fatal(err)
	}
	comparisons := 0
	for name, fixture := range input.Cases {
		for _, before := range fixture.Results {
			comparisons++
			known, percent := true, -5.0
			if name == "nil-without-prior" {
				known, percent = false, 0
			}
			if before.HasCostData != known || before.AvgCostDifference != percent || before.AvgLatencyMs != 42 || !before.LatencyEstimated {
				t.Fatal("full preservation/stale-cost premise changed", name)
			}
		}
	}
	duplicate := input.Cases["duplicate-last-one"].Meters
	if input.SourceBlob != "82b357d2c00dd65c08bc5908000311664431c0e6" || input.TypesBlob != "f9b61fb74075e5a2297be1b898067b8cb047f9b3" || len(input.Cases) != 26 || comparisons != 26 || len(duplicate) != 3 || duplicate[0].MeterID != "m1" || duplicate[1].MeterID != "m1" || duplicate[0].HistoricalCost != 9 || duplicate[1].HistoricalCost != 1 || input.Cases["nil-shared"].Shared != nil || input.Cases["negative-weight"].Meters[0].HistoricalCost != -1 || input.Cases["exact-threshold"].Shared.RegionPricing["m1"]["eastus"] != .0001 {
		t.Fatal("complete pinned cost source topology/arithmetic premises changed")
	}
}

func TestCostEnrichmentCapturedHashes(t *testing.T) {
	files := map[string]string{
		"source-cost-enrichment-inputs.json": "112a462e1bb2090afb39322684264189a8d1977b7c7cab938c8b217535ed670d",

		"source-cost-enrichment-outputs.json": "0ae38cf25c54462e3b7314c4f74c31a5bbe0cfac6d22af72554198ffeb447526",
	}
	for name, want := range files {
		raw := availabilityCapturedBytes(t, name)
		hash := sha256.Sum256(raw)
		if hex.EncodeToString(hash[:]) != want || !json.Valid(raw) {
			t.Fatal("pinned cost-enrichment capture bytes changed", name)
		}
	}
}
