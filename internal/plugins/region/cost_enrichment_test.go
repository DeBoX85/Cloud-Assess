package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type costCapturedInput struct {
	Results []availabilityCapturedComparison

	Meters []HistoricalCostMeter

	Shared *struct {
		RegionPricing map[string]map[string]float64
	}
}

func costRuntimeFixtures(t *testing.T) (map[string]costCapturedInput, map[string][]availabilityCapturedComparison) {
	t.Helper()
	var input struct {
		Cases map[string]costCapturedInput `json:"cases"`
	}
	var output map[string][]availabilityCapturedComparison
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-cost-enrichment-inputs.json"), &input); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-cost-enrichment-outputs.json"), &output); err != nil {
		t.Fatal(err)
	}
	return input.Cases, output
}

func costRuntimeInput() ([]Comparison, *CostEnrichmentEvidence) {
	input := []Comparison{literalCostEnrichmentComparison("eastus", "westeurope", -5, true).Comparison}
	evidence := &CostEnrichmentEvidence{PricingStatus: "complete", History: map[string]CostHistoryEvidence{availabilityID: {Complete: true, Meters: []HistoricalCostMeter{{MeterID: "m1", HistoricalCost: 3}, {MeterID: "m2", HistoricalCost: 1}}}}, RegionPricing: map[string]map[string]float64{"m1": {"eastus": 1, "westeurope": 2}, "m2": {"eastus": 2, "westeurope": 1}}}
	return input, evidence
}

func costHealth(t *testing.T, got *CostCalculation, records int, codes ...string) {
	t.Helper()
	want := assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: records}
	for _, code := range codes {
		want.Status = assessment.StageCompletedWithWarnings
		want.Warnings = append(want.Warnings, assessment.AssessmentWarning{Code: code, Message: "region cost evidence is incomplete or ineligible"})
	}
	if !reflect.DeepEqual(got.Health, want) {
		t.Fatal("complete cost-only health differs", got.Health, want)
	}
}

func TestCostRuntimeCapturedComparisons(t *testing.T) {
	inputs, outputs := costRuntimeFixtures(t)
	checked := 0
	for name, fixture := range inputs {
		switch name {
		case "negative-prices", "negative-weight", "duplicate-last-one", "duplicate-last-nine", "display-names", "unbound-subscription":
			continue // Explicit source-invalid corrections are tested separately.
		}
		t.Run(name, func(t *testing.T) {
			input, want := []Comparison{}, []Comparison{}
			for n, before := range fixture.Results {
				input = append(input, before.Comparison)
				want = append(want, outputs[name][n].Comparison)
			}
			evidence := &CostEnrichmentEvidence{PricingStatus: "complete", History: map[string]CostHistoryEvidence{availabilityID: {Complete: true, Meters: fixture.Meters}}}
			if fixture.Shared == nil {
				evidence.PricingStatus = "unavailable"
			} else {
				evidence.RegionPricing = fixture.Shared.RegionPricing
			}
			codes := []string{}
			switch name {
			case "nil-shared", "nil-without-prior":
				codes = []string{"cost_pricing_unavailable"}
			case "missing-source", "missing-target", "empty-pricing", "unowned-only":
				codes = []string{"cost_no_eligible_meters", "cost_price_missing"}
			case "empty-meters":
				codes = []string{"cost_no_eligible_meters"}
			case "logical-source", "logical-target":
				codes = []string{"cost_region_unsupported"}
			case "near-zero-source":
				codes = []string{"cost_meter_ineligible"}
			}
			// Source-invalid stale state is deliberately cleared, never copied.
			if len(codes) > 0 && name != "near-zero-source" {
				for n := range want {
					want[n].AvgCostDifference, want[n].HasCostData = 0, false
				}
			}
			before, _ := json.Marshal(input)
			got, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
			if err != nil || got == nil || !reflect.DeepEqual(got.Comparisons, want) {
				t.Fatalf("complete cost source/correction differs for %s: %#v %v want %#v", name, got, err, want)
			}
			costHealth(t, got, len(want), codes...)
			after, _ := json.Marshal(input)
			if !bytes.Equal(before, after) {
				t.Fatal("source-style input mutation retained")
			}
			checked++
		})
	}
	if checked != 20 {
		t.Fatal("valid/corrected source case coverage differs", checked)
	}
}

func costReject(t *testing.T, scope map[string]string, input []Comparison, evidence *CostEnrichmentEvidence, code string) {
	t.Helper()
	got, err := EnrichCost(context.Background(), scope, input, evidence)
	if got != nil || err == nil || !strings.Contains(err.Error(), "["+code+"]") {
		t.Fatalf("no-partial %s rejection missing: %#v %v", code, got, err)
	}
}

func TestCostRuntimeAdmission(t *testing.T) {
	input, evidence := costRuntimeInput()
	bad := input[0]
	bad.SubscriptionID = "22222222-2222-2222-2222-222222222222"
	bad.SubscriptionName = ""
	costReject(t, latencyScope(), []Comparison{bad}, nil, "comparison_invalid")
	bad = input[0]
	bad.SourceRegion = "East US"
	costReject(t, latencyScope(), []Comparison{bad}, evidence, "comparison_invalid")
	bad = input[0]
	bad.AvgCostDifference = math.NaN()
	costReject(t, latencyScope(), []Comparison{bad}, evidence, "comparison_invalid")
	costReject(t, latencyScope(), append(input, input[0]), evidence, "comparison_invalid")
	costReject(t, latencyScope(), make([]Comparison, MaxComparisons+1), evidence, "comparison_invalid")
	bad = input[0]
	bad.SubscriptionID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	scope := map[string]string{bad.SubscriptionID: "selected"}
	bad.SubscriptionID = strings.ToUpper(bad.SubscriptionID)
	costReject(t, scope, []Comparison{bad}, nil, "comparison_invalid")
	_, evidence = costRuntimeInput()
	evidence.History["22222222-2222-2222-2222-222222222222"] = evidence.History[availabilityID]
	costReject(t, latencyScope(), input, evidence, "evidence_scope_invalid")
	_, evidence = costRuntimeInput()
	evidence.History[availabilityID] = CostHistoryEvidence{Complete: true, Meters: []HistoricalCostMeter{{MeterID: "m1", HistoricalCost: 9}, {MeterID: "m1", HistoricalCost: 1}}}
	costReject(t, latencyScope(), input, evidence, "duplicate_meter")
	for _, status := range []string{"", "unknown", "Complete"} {
		_, evidence = costRuntimeInput()
		evidence.PricingStatus = status
		costReject(t, latencyScope(), input, evidence, "status_invalid")
	}
	for _, status := range []string{"unavailable", "unsupported"} {
		_, evidence = costRuntimeInput()
		evidence.PricingStatus = status
		costReject(t, latencyScope(), input, evidence, "status_invalid")
	}
	for _, value := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1), MaxCostValue + 1} {
		_, evidence = costRuntimeInput()
		h := evidence.History[availabilityID]
		h.Meters[0].HistoricalCost = value
		costReject(t, latencyScope(), input, evidence, "value_invalid")
		_, evidence = costRuntimeInput()
		evidence.RegionPricing["unowned"] = map[string]float64{"eastus": value}
		costReject(t, latencyScope(), input, evidence, "value_invalid")
	}
	for _, id := range []string{"", "m\n", strings.Repeat("a", MaxLabelBytes+1), string([]byte{0xff})} {
		_, evidence = costRuntimeInput()
		evidence.RegionPricing = map[string]map[string]float64{id: {"eastus": 1}}
		costReject(t, latencyScope(), input, evidence, "text_limit")
		_, evidence = costRuntimeInput()
		evidence.History[availabilityID] = CostHistoryEvidence{Complete: true, Meters: []HistoricalCostMeter{{MeterID: id}}}
		costReject(t, latencyScope(), input, evidence, "text_limit")
	}
	for _, region := range []string{"EastUS", "east us", "", strings.Repeat("a", 65)} {
		_, evidence = costRuntimeInput()
		evidence.RegionPricing = map[string]map[string]float64{"m1": {region: 1}}
		costReject(t, latencyScope(), input, evidence, "region_invalid")
	}
}

func TestCostRuntimeHealthAndSubscriptionWeights(t *testing.T) {
	input, evidence := costRuntimeInput()
	otherID := "22222222-2222-2222-2222-222222222222"
	other := input[0]
	other.SubscriptionID = otherID
	scope := map[string]string{availabilityID: "selected", otherID: "selected"}
	input = append(input, other)
	evidence.History[otherID] = CostHistoryEvidence{Complete: true, Meters: []HistoricalCostMeter{{MeterID: "m1", HistoricalCost: 1}, {MeterID: "m2", HistoricalCost: 3}}}
	got, err := EnrichCost(context.Background(), scope, input, evidence)
	want := []Comparison{literalCostEnrichmentComparison("eastus", "westeurope", 62.5, true).Comparison, literalCostEnrichmentComparison("eastus", "westeurope", -12.5, true).Comparison}
	want[1].SubscriptionID = otherID
	if err != nil || !reflect.DeepEqual(got.Comparisons, want) {
		t.Fatal("distinct selected-subscription weighted full results differ", got, err)
	}
	costHealth(t, got, 2)
	for id, h := range evidence.History {
		slices.Reverse(h.Meters)
		evidence.History[id] = h
	}
	for n := 0; n < 8; n++ {
		again, err := EnrichCost(context.Background(), scope, input, evidence)
		if err != nil || !reflect.DeepEqual(again, got) {
			t.Fatal("received/map order changes weighted result", err)
		}
	}
	delete(evidence.History, otherID)
	got, err = EnrichCost(context.Background(), scope, input, evidence)
	if err != nil || !got.Comparisons[0].HasCostData || got.Comparisons[1].HasCostData || got.Comparisons[1].AvgCostDifference != 0 {
		t.Fatal("missing history reused another subscription", err)
	}
	costHealth(t, got, 2, "cost_history_unavailable")
	input, evidence = costRuntimeInput()
	h := evidence.History[availabilityID]
	h.Complete = false
	evidence.History[availabilityID] = h
	got, err = EnrichCost(context.Background(), latencyScope(), input, evidence)
	if err != nil || got.Comparisons[0].HasCostData || got.Comparisons[0].AvgCostDifference != 0 {
		t.Fatal("partial history treated complete", err)
	}
	costHealth(t, got, 1, "cost_history_partial")
	input, evidence = costRuntimeInput()
	evidence.PricingStatus = "partial"
	delete(evidence.RegionPricing, "m2")
	got, err = EnrichCost(context.Background(), latencyScope(), input, evidence)
	if err != nil || got.Comparisons[0].AvgCostDifference != 100 || !got.Comparisons[0].HasCostData {
		t.Fatal("partial pricing discarded valid owned pair", err)
	}
	costHealth(t, got, 1, "cost_price_missing", "cost_pricing_partial")
	got, err = EnrichCost(context.Background(), latencyScope(), input, nil)
	if err != nil || got.Comparisons[0].HasCostData || got.Comparisons[0].AvgCostDifference != 0 {
		t.Fatal("nil evidence preserved stale cost", err)
	}
	costHealth(t, got, 1, "cost_history_unavailable", "cost_pricing_unavailable")
}

func TestCostRuntimeBudgets(t *testing.T) {
	input, evidence := costRuntimeInput()
	evidence.History = map[string]CostHistoryEvidence{availabilityID: {Complete: true, Meters: []HistoricalCostMeter{}}}
	h := evidence.History[availabilityID]
	for n := 0; n < MaxCostMeters; n++ {
		h.Meters = append(h.Meters, HistoricalCostMeter{MeterID: fmt.Sprintf("m%d", n), HistoricalCost: MaxCostValue})
	}
	evidence.History[availabilityID] = h
	evidence.RegionPricing = nil
	if got, err := EnrichCost(context.Background(), latencyScope(), input, evidence); err != nil || got == nil {
		t.Fatal("exact historical meter/value boundary rejected", err)
	}
	h.Meters = append(h.Meters, HistoricalCostMeter{MeterID: "over"})
	evidence.History[availabilityID] = h
	costReject(t, latencyScope(), input, evidence, "input_limit")
	_, evidence = costRuntimeInput()
	evidence.History = nil
	evidence.RegionPricing = map[string]map[string]float64{}
	for n := 0; n < 1985; n++ {
		prices := map[string]float64{}
		for r := 0; r < 32; r++ {
			prices[fmt.Sprintf("r%d", r)] = MaxCostValue
		}
		evidence.RegionPricing[fmt.Sprintf("m%d", n)] = prices
	}
	prices := map[string]float64{}
	for r := 0; r < 30; r++ {
		prices[fmt.Sprintf("r%d", r)] = 1
	}
	evidence.RegionPricing["last"] = prices
	// 1985*(1+32) + (1+30) =65536; 32 unique regions.
	if got, err := EnrichCost(context.Background(), latencyScope(), nil, evidence); err != nil || got == nil {
		t.Fatal("exact decoded-entry/region/value boundary rejected", err)
	}
	prices["r30"] = 1
	costReject(t, latencyScope(), nil, evidence, "input_limit")
	evidence.RegionPricing = map[string]map[string]float64{"m": {}}
	for r := 0; r < 33; r++ {
		evidence.RegionPricing["m"][fmt.Sprintf("r%d", r)] = 1
	}
	costReject(t, latencyScope(), nil, evidence, "input_limit")
	evidence.RegionPricing = map[string]map[string]float64{}
	for n := 0; n < MaxCostMeters; n++ {
		evidence.RegionPricing[fmt.Sprintf("m%d", n)] = nil
	}
	if _, err := EnrichCost(context.Background(), latencyScope(), nil, evidence); err != nil {
		t.Fatal("exact pricing meter boundary rejected", err)
	}
	evidence.RegionPricing["over"] = nil
	costReject(t, latencyScope(), nil, evidence, "input_limit")
	id, region := strings.Repeat("m", 512), strings.Repeat("r", 64)
	evidence.RegionPricing = map[string]map[string]float64{id: {region: MaxCostValue}}
	if _, err := EnrichCost(context.Background(), latencyScope(), nil, evidence); err != nil {
		t.Fatal("exact admitted label/region boundary rejected", err)
	}
	// Text maximum is independently implied by count/label bounds; no false
	// independently triggerable decoded-text overflow assertion is claimed.
}

func TestCostRuntimeWorkBudget(t *testing.T) {
	_, evidence := costRuntimeInput()
	h := CostHistoryEvidence{Complete: true}
	for n := 0; n < 8192; n++ {
		h.Meters = append(h.Meters, HistoricalCostMeter{MeterID: fmt.Sprintf("m%d", n)})
	}
	evidence.History[availabilityID] = h
	evidence.RegionPricing = nil
	input := []Comparison{}
	for n := 0; n < 128; n++ {
		input = append(input, Comparison{SubscriptionID: availabilityID, SubscriptionName: "selected", SourceRegion: "eastus", TargetRegion: fmt.Sprintf("r%d", n)})
	}
	if got, err := EnrichCost(context.Background(), latencyScope(), input, evidence); err != nil || got == nil {
		t.Fatal("exact1048576 work boundary rejected", err)
	}
	// The smallest possible next comparison with this fully owned meter list
	// exceeds work while staying well within every comparison/text/entry cap.
	input = append(input, Comparison{SubscriptionID: availabilityID, SubscriptionName: "selected", SourceRegion: "eastus", TargetRegion: "over"})
	costReject(t, latencyScope(), input, evidence, "work_limit")
}

func TestCostRuntimeOwnershipAndCancellation(t *testing.T) {
	input, evidence := costRuntimeInput()
	// Exercise all four nonempty detail slices with valid matching SKU counts.
	input[0].RestrictedSKUs = []string{"restricted"}
	input[0].ZoneRestrictedSKUs = []string{"zone-restricted"}
	input[0].TotalSKUsChecked = 3
	input[0].SKUAvailabilityPercent = 0
	evidence.PricingStatus = "partial"
	before, _ := json.Marshal(struct {
		Input []Comparison

		Evidence *CostEnrichmentEvidence
	}{input, evidence})
	baseline, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for n := 0; n < 16; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
			if err != nil || !reflect.DeepEqual(got, baseline) {
				t.Error("concurrent cost state leaked", err)
			}
		}()
	}
	wg.Wait()
	probe := &availabilityCancelContext{Context: context.Background(), at: 1000000}
	if _, err := EnrichCost(probe, latencyScope(), input, evidence); err != nil {
		t.Fatal(err)
	}
	for _, at := range []int{1, probe.calls / 2, probe.calls} {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		got, err := EnrichCost(ctx, latencyScope(), input, evidence)
		if got != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("partial cost returned at cancellation%d: %#v %v", at, got, err)
		}
	}
	baseline.Comparisons[0].MissingResourceTypes[0] = "changed"
	baseline.Comparisons[0].MissingSKUs[0] = "changed"
	baseline.Comparisons[0].RestrictedSKUs[0] = "changed"
	baseline.Comparisons[0].ZoneRestrictedSKUs[0] = "changed"
	baseline.Comparisons[0].TargetZoneMappings["1"] = "changed"
	baseline.Health.Warnings[0].Code = "changed"
	after, _ := json.Marshal(struct {
		Input []Comparison

		Evidence *CostEnrichmentEvidence
	}{input, evidence})
	if !bytes.Equal(before, after) {
		t.Fatal("result aliases input/evidence")
	}
	again, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
	if err != nil || again.Comparisons[0].AvgCostDifference != 62.5 || again.Comparisons[0].TargetZoneMappings["1"] != "phys1" || again.Health.Warnings[0].Code != "cost_pricing_partial" {
		t.Fatal("prior result leaked", err)
	}
	_, different := costRuntimeInput()
	different.RegionPricing["m1"]["westeurope"] = 1
	different.RegionPricing["m2"]["westeurope"] = 2
	zero, err := EnrichCost(context.Background(), latencyScope(), input, different)
	if err != nil || zero.Comparisons[0].AvgCostDifference != 0 || !zero.Comparisons[0].HasCostData {
		t.Fatal("different sequential evidence reused old values", err)
	}
}

func TestCostRuntimeAvailabilityLatencyProjectionIntegration(t *testing.T) {
	scope, inventory, request, availabilityEvidence := availabilityInputs()
	inventory.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.custom/widgets": 1}
	availabilityEvidence.Locations["microsoft.custom/widgets"] = map[string]bool{"westeurope": true}
	availabilityEvidence.InventoryComplete = false
	available, err := CalculateAvailability(context.Background(), scope, inventory, request, availabilityEvidence)
	if err != nil || available.Health.Status != assessment.StageCompletedWithWarnings {
		t.Fatal("availability partial premise missing", err)
	}
	latency, err := EnrichLatency(context.Background(), scope, []Comparison{available.Comparison})
	if err != nil {
		t.Fatal(err)
	}
	_, evidence := costRuntimeInput()
	cost, err := EnrichCost(context.Background(), scope, latency.Comparisons, evidence)
	if err != nil || cost.Health.Status != assessment.StageCompleted || available.Health.Status != assessment.StageCompletedWithWarnings || latency.Health.Status != assessment.StageCompleted {
		t.Fatal("cost repaired upstream completeness or lost integration", err)
	}
	table, err := Project(context.Background(), scope, cost.Comparisons)
	// 100*.35+100*.30+0*.15+(100-35/150*100)*.20=80.3333;
	// one fewer target zone applies0.9666667 =>77.65556.
	want := []string{"selected", "eastus", "westeurope", "1", "1", "0", "100.00%", "0", "0", "0", "0", "0", "0", "N/A", "3 → 2 ⚠", "", "85.0", "+62.50%", "77.66", "Full", "Neutral", "", "", "", ""}
	if err != nil || len(table.Rows) != 1 || !reflect.DeepEqual(table.Rows[0].Cells, want) || table.Rows[0].SubscriptionID != availabilityID {
		t.Fatal("complete independent availability/latency/cost projection differs", table, err)
	}
	if available.Comparison.HasCostData || latency.Comparisons[0].HasCostData {
		t.Fatal("enrichment mutated earlier stage results")
	}
}

func TestCostRuntimeDeterministicSum(t *testing.T) {
	input, evidence := costRuntimeInput()
	evidence.History[availabilityID] = CostHistoryEvidence{Complete: true, Meters: []HistoricalCostMeter{{MeterID: "a", HistoricalCost: 1e12}, {MeterID: "b", HistoricalCost: 1e12}, {MeterID: "c", HistoricalCost: 1}}}
	evidence.RegionPricing = map[string]map[string]float64{"a": {"eastus": 1, "westeurope": 2}, "b": {"eastus": 1, "westeurope": 0}, "c": {"eastus": 1, "westeurope": 1.00001}}
	first, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
	if err != nil || first == nil || math.Abs(first.Comparisons[0].AvgCostDifference-5.000000000030256e-16) > 1e-27 {
		t.Fatal("independent small residual lost", first, err)
	}
	h := evidence.History[availabilityID]
	slices.Reverse(h.Meters)
	evidence.History[availabilityID] = h
	second, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
	if err != nil || !reflect.DeepEqual(first, second) {
		t.Fatal("received-order floating sum differs", first, second, err)
	}
}

func TestCostRuntimeFreeTargetRoundoff(t *testing.T) {
	input, evidence := costRuntimeInput()
	weights := []float64{50472.046742886334, 48492.51122277341, 35678.99645449557, 34607.791901815486, 53847.87957378443, 62348.945279750515, 61245.246478272566, 45814.68000997244}
	h := CostHistoryEvidence{Complete: true}
	evidence.RegionPricing = map[string]map[string]float64{}
	for n, w := range weights {
		id := fmt.Sprintf("m%d", n)
		h.Meters = append(h.Meters, HistoricalCostMeter{MeterID: id, HistoricalCost: w})
		evidence.RegionPricing[id] = map[string]float64{"eastus": 1, "westeurope": 0}
	}
	evidence.History[availabilityID] = h
	got, err := EnrichCost(context.Background(), latencyScope(), input, evidence)
	want := []Comparison{literalCostEnrichmentComparison("eastus", "westeurope", -100, true).Comparison}
	if err != nil || got == nil || !reflect.DeepEqual(got.Comparisons, want) {
		t.Fatal("valid free-target lower-domain rounding rejected", got, err)
	}
	costHealth(t, got, 1)
}
