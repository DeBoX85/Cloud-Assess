package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func latencyFixtures(t *testing.T) (map[string]latencyCapturedInput, map[string]latencyCapturedCase, latencyData) {
	t.Helper()
	var input map[string]latencyCapturedInput
	var output map[string]latencyCapturedCase
	var data latencyData
	for name, destination := range map[string]any{"source-latency-inputs.json": &input, "source-latency-outputs.json": &output, "source-latency-data.json": &data} {
		if err := json.Unmarshal(availabilityCapturedBytes(t, name), destination); err != nil {
			t.Fatal(err)
		}
	}
	return input, output, data
}

func latencyScope() map[string]string { return map[string]string{availabilityID: "selected"} }

func TestLatencyRuntimeCapturedComparisons(t *testing.T) {
	input, output, data := latencyFixtures(t)
	checked := 0
	for name, fixture := range input {
		if name == "negative" {
			continue // Source retains negative evidence; target deliberately rejects it.
		}
		t.Run(name, func(t *testing.T) {
			comparisons, want := []Comparison{}, []Comparison{}
			for n, before := range fixture.Results {
				if before.SourceRegion == "East US" {
					continue // Canonical identity admission is an explicit correction.
				}
				comparisons = append(comparisons, before.Comparison)
				want = append(want, output[name].Comparisons[n].Comparison)
			}
			before, _ := json.Marshal(comparisons)
			synthetic := data
			synthetic.Matrix = fixture.Matrix
			got, err := enrichLatency(context.Background(), latencyScope(), comparisons, synthetic)
			if name == "default" {
				got, err = EnrichLatency(context.Background(), latencyScope(), comparisons)
			}
			if err != nil || got == nil || !reflect.DeepEqual(got.Comparisons, want) {
				t.Fatalf("complete source latency subset differs: %#v / %v, want %#v", got, err, want)
			}
			codes := []string{}
			for _, warning := range got.Health.Warnings {
				codes = append(codes, warning.Code)
			}
			wantCodes := map[string][]string{"synthetic": {"latency_estimated", "latency_unknown"}, "default": {"latency_estimated", "latency_unknown"}, "one-way-cluster": {"latency_estimated", "latency_unknown"}, "empty": {"latency_unknown"}, "zero": {"latency_unknown"}, "empty-results": {}}[name]
			status := assessment.StageCompleted
			if len(wantCodes) != 0 {
				status = assessment.StageCompletedWithWarnings
			}
			if !reflect.DeepEqual(codes, wantCodes) || got.Health.Status != status || got.Health.Records != len(want) || got.Health.Name != Name {
				t.Fatal("latency-only health missing", got.Health)
			}
			after, _ := json.Marshal(comparisons)
			if !bytes.Equal(before, after) {
				t.Fatal("source-style in-place mutation retained")
			}
			checked += len(want)
		})
	}
	if checked != 24 {
		t.Fatal("valid source-subset coverage changed", checked)
	}
}

func latencyReject(t *testing.T, scope map[string]string, input []Comparison, data latencyData, code string) {
	t.Helper()
	got, err := enrichLatency(context.Background(), scope, input, data)
	if got != nil || err == nil || !strings.Contains(err.Error(), "["+code+"]") {
		t.Fatalf("missing no-partial %s rejection: %#v %v", code, got, err)
	}
}

func TestLatencyRuntimeAdmission(t *testing.T) {
	_, _, data := latencyFixtures(t)
	c := literalLatencyComparison("eastus", "westeurope", 42, true).Comparison
	wrong := c
	wrong.SubscriptionID = "22222222-2222-2222-2222-222222222222"
	// Empty name equals the absent scope lookup: selection alone must reject.
	wrong.SubscriptionName = ""
	latencyReject(t, latencyScope(), []Comparison{wrong}, data, "comparison_invalid")
	wrong = c
	wrong.SubscriptionName = "other"
	latencyReject(t, latencyScope(), []Comparison{wrong}, data, "comparison_invalid")
	wrong = c
	wrong.SourceRegion = "East US"
	latencyReject(t, latencyScope(), []Comparison{wrong}, data, "comparison_invalid")
	latencyReject(t, latencyScope(), []Comparison{c, c}, data, "comparison_invalid")
	latencyReject(t, latencyScope(), make([]Comparison, MaxComparisons+1), data, "comparison_invalid")
	wrong = c
	wrong.MissingSKUs = make([]string, MaxDetails+1)
	latencyReject(t, latencyScope(), []Comparison{wrong}, data, "comparison_invalid")
	wrong = c
	wrong.AvgLatencyMs = math.NaN()
	latencyReject(t, latencyScope(), []Comparison{wrong}, data, "comparison_invalid")
}

func TestLatencyRuntimeDataRejection(t *testing.T) {
	_, _, data := latencyFixtures(t)
	for _, value := range []float64{-2, math.NaN(), math.Inf(1), math.Inf(-1), MaxLatencyMilliseconds + 1} {
		bad := data
		bad.Matrix = map[string]map[string]float64{"eastus": {"westeurope": value}}
		latencyReject(t, latencyScope(), nil, bad, "value_invalid")
	}
	for _, invalid := range []string{"", "EastUS", "east us", strings.Repeat("a", 65), "éast", "a\n"} {
		bad := data
		bad.Matrix = map[string]map[string]float64{invalid: {"westeurope": 1}}
		latencyReject(t, latencyScope(), nil, bad, "label_limit")
		bad.Matrix = map[string]map[string]float64{"eastus": {invalid: 1}}
		latencyReject(t, latencyScope(), nil, bad, "label_limit")
		bad.Matrix = nil
		bad.Clusters = map[string]string{"eastus": invalid}
		latencyReject(t, latencyScope(), nil, bad, "label_limit")
		bad.Clusters = map[string]string{invalid: "europe"}
		latencyReject(t, latencyScope(), nil, bad, "label_limit")
	}
	bad := data
	bad.SourceBlob = "unverified"
	latencyReject(t, latencyScope(), nil, bad, "source_invalid")
	got, err := decodeLatency(context.Background(), bytes.Repeat([]byte{' '}, MaxLatencyBytes+1))
	if err == nil || !strings.Contains(err.Error(), "[serialized_limit]") || got.Matrix != nil {
		t.Fatal("serialized admission missing", err)
	}
	corrupt := append(bytes.Clone(bundledLatencyJSON), ' ')
	got, err = decodeLatency(context.Background(), corrupt)
	if err == nil || !strings.Contains(err.Error(), "[source_invalid]") || got.Matrix != nil {
		t.Fatal("pinned bytes admission missing", err)
	}
	decoded, err := decodeLatency(context.Background(), bundledLatencyJSON)
	if err != nil || !reflect.DeepEqual(decoded, data) || !bytes.Equal(bundledLatencyJSON, availabilityCapturedBytes(t, "source-latency-data.json")) {
		t.Fatal("embedded dataset differs from independently observed complete source", err)
	}
	decoded.Matrix["eastus"]["westeurope"] = 999
	again, err := decodeLatency(context.Background(), bundledLatencyJSON)
	if err != nil || again.Matrix["eastus"]["westeurope"] != 85 {
		t.Fatal("decoded dataset reused mutable state", err)
	}
}

func TestLatencyRuntimeBudgets(t *testing.T) {
	data := latencyData{SourceBlob: latencySourceBlob, Matrix: map[string]map[string]float64{"eastus": {}}, Clusters: map[string]string{}}
	for n := 0; n < MaxLatencyEntries-1; n++ {
		data.Matrix["eastus"][fmt.Sprintf("target%d", n)] = MaxLatencyMilliseconds
	}
	if got, err := enrichLatency(context.Background(), latencyScope(), nil, data); err != nil || got == nil {
		t.Fatal("exact nested entry/value boundary rejected", err)
	}
	data.Matrix["eastus"]["overflow"] = 1
	latencyReject(t, latencyScope(), nil, data, "input_limit")
	data.Matrix = nil
	for n := 0; n < MaxLatencyEntries; n++ {
		data.Clusters[fmt.Sprintf("r%07d", n)+strings.Repeat("a", 56)] = strings.Repeat("b", 64)
	}
	// 8192 owned cluster declarations of two 64-byte labels hit exactly1MiB.
	// Count plus label bounds independently imply this aggregate maximum.
	if _, err := enrichLatency(context.Background(), latencyScope(), nil, data); err != nil {
		t.Fatal("exact decoded-label budget rejected", err)
	}
	data.Clusters["overflow"] = "europe"
	latencyReject(t, latencyScope(), nil, data, "input_limit")
}

func TestLatencyRuntimeOwnershipAndCancellation(t *testing.T) {
	c := literalLatencyComparison("austriaeast", "eastus", 42, true).Comparison
	c.RestrictedSKUs, c.ZoneRestrictedSKUs = []string{"restricted"}, []string{"zone"}
	c.TotalSKUsChecked = 3
	c.SKUAvailabilityPercent = 0
	input := []Comparison{c}
	before, _ := json.Marshal(input)
	baseline, err := EnrichLatency(context.Background(), latencyScope(), input)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for n := 0; n < 16; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := EnrichLatency(context.Background(), latencyScope(), input)
			if err != nil || !reflect.DeepEqual(got, baseline) {
				t.Error("concurrent latency state leaked", err)
			}
		}()
	}
	wg.Wait()
	probe := &availabilityCancelContext{Context: context.Background(), at: 1000000}
	if _, err := EnrichLatency(probe, latencyScope(), input); err != nil {
		t.Fatal(err)
	}
	for _, at := range []int{1, 100, probe.calls} {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		got, err := EnrichLatency(ctx, latencyScope(), input)
		if got != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("latency cancellation returned partial result at%d: %#v %v", at, got, err)
		}
	}
	baseline.Comparisons[0].MissingResourceTypes[0] = "modified"
	baseline.Comparisons[0].MissingSKUs[0] = "modified"
	baseline.Comparisons[0].RestrictedSKUs[0] = "modified"
	baseline.Comparisons[0].ZoneRestrictedSKUs[0] = "modified"
	baseline.Comparisons[0].TargetZoneMappings["1"] = "modified"
	baseline.Health.Warnings[0].Code = "modified"
	after, _ := json.Marshal(input)
	if !bytes.Equal(before, after) {
		t.Fatal("latency result aliases caller-owned proof")
	}
	again, err := EnrichLatency(context.Background(), latencyScope(), input)
	if err != nil || again.Comparisons[0].TargetZoneMappings["1"] != "phys1" || again.Health.Warnings[0].Code != "latency_estimated" {
		t.Fatal("previous result leaked into subsequent latency call", err)
	}
}

func TestLatencyRuntimeAvailabilityProjectionIntegration(t *testing.T) {
	scope, inventory, request, evidence := availabilityInputs()
	inventory.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.custom/widgets": 1}
	evidence.Locations["microsoft.custom/widgets"] = map[string]bool{"westeurope": true}
	available, err := CalculateAvailability(context.Background(), scope, inventory, request, evidence)
	if err != nil {
		t.Fatal(err)
	}
	latency, err := EnrichLatency(context.Background(), scope, []Comparison{available.Comparison})
	if err != nil || latency.Health.Status != assessment.StageCompleted || latency.Comparisons[0].AvgLatencyMs != 85 || available.Comparison.AvgLatencyMs != 0 {
		t.Fatal("availability-to-owned-latency integration failed", err)
	}
	table, err := Project(context.Background(), scope, latency.Comparisons)
	want := []string{"selected", "eastus", "westeurope", "1", "1", "0", "100.00%", "0", "0", "0", "0", "0", "0", "N/A", "3 → 2 ⚠", "", "85.0", "N/A", "92.16", "no cost data", "Recommended", "", "", "", ""}
	if err != nil || len(table.Rows) != 1 || !reflect.DeepEqual(table.Rows[0].Cells, want) || table.Rows[0].SubscriptionID != availabilityID {
		t.Fatal("complete integration projection differs", table, err)
	}
	// Availability incompleteness cannot be repaired by known latency or projection.
	evidence.InventoryComplete = false
	available, err = CalculateAvailability(context.Background(), scope, inventory, request, evidence)
	if err != nil || available.Health.Status != assessment.StageCompletedWithWarnings {
		t.Fatal("availability incompleteness premise missing", err)
	}
	latency, err = EnrichLatency(context.Background(), scope, []Comparison{available.Comparison})
	if err != nil || latency.Health.Status != assessment.StageCompleted || available.Health.Status != assessment.StageCompletedWithWarnings {
		t.Fatal("stage-specific health boundary changed", err)
	}
}
