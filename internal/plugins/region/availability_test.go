package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const availabilityID = "11111111-1111-1111-1111-111111111111"

func availabilityInputs() (map[string]string, *InventoryCalculation, AvailabilityRequest, AvailabilityEvidence) {
	return map[string]string{availabilityID: "selected"}, &InventoryCalculation{Subscriptions: map[string]InventoryCounts{availabilityID: emptyInventoryCounts()}}, AvailabilityRequest{SubscriptionID: availabilityID, SourceRegion: "eastus", TargetRegion: "westeurope"}, AvailabilityEvidence{SubscriptionID: availabilityID, TargetRegion: "westeurope", InventoryComplete: true, ProvidersComplete: true, Locations: map[string]map[string]bool{}, SKUEnabled: true, SKUs: map[string]SKUEvidence{}, ZoneCounts: map[string]int{"eastus": 3, "westeurope": 2}}
}

func availabilityReject(t *testing.T, scope map[string]string, inventory *InventoryCalculation, request AvailabilityRequest, evidence AvailabilityEvidence, code string) {
	t.Helper()
	got, err := CalculateAvailability(context.Background(), scope, inventory, request, evidence)
	if got != nil || err == nil || !strings.Contains(err.Error(), "["+code+"]") {
		t.Fatalf("expected safe %s rejection, got %#v / %v", code, got, err)
	}
}

func TestAvailabilityRuntimeCapturedComparisons(t *testing.T) {
	var input struct {
		Cases map[string]struct {
			Source string
			Target string
			Inventory InventoryCounts
			Providers struct { Data map[string]map[string]map[string]struct{} }
			Zones map[string]int
			SKUEnabled bool
			Cancelled bool
		}
		SKUResponse map[string]SKUAvailability `json:"sku_response"`
	}
	var output availabilityCapturedOutput
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-availability-inputs.json"), &input); err != nil { t.Fatal(err) }
	if err := json.Unmarshal(availabilityCapturedBytes(t, "source-availability-outputs.json"), &output); err != nil { t.Fatal(err) }
	checked := 0
	for key, fixture := range input.Cases {
		if key == "source-key-exact-case" || key == "target-key-exact-case" || fixture.Cancelled { continue }
		t.Run(key, func(t *testing.T) {
			scope, inventory, request, evidence := availabilityInputs()
			inventory.Subscriptions[availabilityID] = fixture.Inventory
			request.SourceRegion, request.TargetRegion = fixture.Source, fixture.Target
			evidence.ZoneCounts, evidence.SKUEnabled = fixture.Zones, fixture.SKUEnabled
			for namespace, types := range fixture.Providers.Data {
				for name, locations := range types {
					entry := map[string]bool{}
					for region := range locations { entry[region] = true }
					evidence.Locations[namespace+"/"+name] = entry
				}
			}
			evidence.SKUs["microsoft.capture/available"] = SKUEvidence{Status: "complete", Values: input.SKUResponse}
			evidence.SKUs["microsoft.capture/failed"] = SKUEvidence{Status: "unknown"}
			evidence.SKUs["microsoft.custom/widgets"] = SKUEvidence{Status: "unsupported"}
			got, err := CalculateAvailability(context.Background(), scope, inventory, request, evidence)
			if err != nil { t.Fatal(err) }
			want := output.Cases[key].Comparison
			want.SubscriptionID, want.SubscriptionName = availabilityID, "selected"
			if !reflect.DeepEqual(got.Comparison, want) { t.Fatalf("entire source comparison mismatch: %#v want %#v", got.Comparison, want) }
			partial := got.Comparison.UnknownSKUs > 0 || key == "sku-all-states" || key == "sku-disabled" || key == "resource-mixed-no-sku" || key == "resource-empty-providers"
			wantStatus := assessment.StageCompleted
			if partial { wantStatus = assessment.StageCompletedWithWarnings }
			if got.Health.Status != wantStatus || got.Health.Records != 1 { t.Fatalf("source corrections missing: %#v", got.Health) }
		})
		checked++
	}
	if checked != 8 { t.Fatal("source subset topology changed", checked) }
}

func TestAvailabilityRuntimeScopeAndEvidence(t *testing.T) {
	tests := []struct { name string; change func(map[string]string, *InventoryCalculation, *AvailabilityRequest, *AvailabilityEvidence); code string }{
		{"wrong-selected", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { r.SubscriptionID = "22222222-2222-2222-2222-222222222222" }, "scope_invalid"},
		{"wrong-evidence-subscription", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { e.SubscriptionID = "22222222-2222-2222-2222-222222222222" }, "scope_invalid"},
		{"wrong-evidence-region", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { e.TargetRegion = "eastus" }, "scope_invalid"},
		{"extra-inventory", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { i.Subscriptions["extra"] = emptyInventoryCounts() }, "scope_invalid"},
		{"missing-inventory", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { delete(i.Subscriptions, availabilityID) }, "scope_invalid"},
		{"empty-name", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { s[availabilityID] = "" }, "scope_invalid"},
		{"region-case", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { r.SourceRegion = "East US" }, "region_invalid"},
		{"negative-count", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": -1} }, "input_invalid"},
		{"invalid-location", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { e.Locations["microsoft.compute/vms"] = map[string]bool{"west europe": true} }, "input_invalid"},
		{"invalid-key", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { e.Locations["malformed"] = map[string]bool{} }, "input_limit"},
		{"normalization-collision", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{"Ready": {}, " ready ": {State: SKURestricted}}} }, "response_collision"},
		{"ambiguous-unknown", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{"Ready": {}}} }, "evidence_invalid"},
		{"missing-sku", func(s map[string]string, i *InventoryCalculation, r *AvailabilityRequest, e *AvailabilityEvidence) { i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": 1}; i.Subscriptions[availabilityID].SKUsByTypeAndRegion["microsoft.compute/vms"] = map[string]map[string]int64{"eastus": {"Ready": 1}} }, "evidence_missing"},
	}
	for _, test := range tests { t.Run(test.name, func(t *testing.T) { s, i, r, e := availabilityInputs(); test.change(s, i, &r, &e); availabilityReject(t, s, i, r, e, test.code) }) }
}

func TestAvailabilityRuntimeUnknownAndCompleteness(t *testing.T) {
	s, i, r, e := availabilityInputs()
	i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": 1}
	i.Subscriptions[availabilityID].SKUsByTypeAndRegion["microsoft.compute/vms"] = map[string]map[string]int64{"eastus": {"First": 1, "Second": 900}}
	e.Locations["microsoft.compute/vms"] = map[string]bool{"westeurope": true}
	e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "unknown"}
	got, err := CalculateAvailability(context.Background(), s, i, r, e)
	if err != nil { t.Fatal(err) }
	want := Comparison{SubscriptionID: availabilityID, SubscriptionName: "selected", SourceRegion: "eastus", TargetRegion: "westeurope", SourceResourceTypeCount: 1, AvailableTypes: 1, AvailabilityPercent: 100, MissingResourceTypes: []string{}, MissingSKUs: []string{"microsoft.compute/vms:First (unknown)", "microsoft.compute/vms:Second (unknown)"}, RestrictedSKUs: []string{}, ZoneRestrictedSKUs: []string{}, TotalSKUsChecked: 2, UnknownSKUs: 2, SKUAvailabilityPercent: 100, SourceZoneCount: 3, TargetZoneCount: 2}
	if !reflect.DeepEqual(got.Comparison, want) || got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Code != "availability_sku_unknown" { t.Fatalf("unknown evidence treated as complete: %#v", got) }
	e.InventoryComplete, e.ProvidersComplete, e.SKUEnabled = false, false, false
	e.ZoneCounts = nil
	got, err = CalculateAvailability(context.Background(), s, i, r, e)
	if err != nil { t.Fatal(err) }
	wantCodes := []string{"availability_inventory_partial", "availability_providers_partial", "availability_sku_disabled", "availability_zones_partial"}
	codes := []string{}
	for _, warning := range got.Health.Warnings { codes = append(codes, warning.Code) }
	if got.Health.Status != assessment.StageCompletedWithWarnings || !reflect.DeepEqual(codes, wantCodes) { t.Fatal("incomplete evidence lost", got.Health) }
}

func TestAvailabilityRuntimeWorkAndText(t *testing.T) {
	s, i, r, e := availabilityInputs()
	for n := 0; n < MaxAvailabilityEntries-2; n++ { e.Locations[fmt.Sprintf("microsoft.work/type%d", n)] = map[string]bool{} }
	if _, err := CalculateAvailability(context.Background(), s, i, r, e); err != nil { t.Fatal("exact input entry boundary rejected", err) }
	e.Locations["microsoft.work/overflow"] = map[string]bool{}
	availabilityReject(t, s, i, r, e, "input_limit")
	s, i, r, e = availabilityInputs()
	e.Locations["microsoft."+strings.Repeat("a", MaxLabelBytes)] = map[string]bool{}
	availabilityReject(t, s, i, r, e, "input_limit")
	s, i, r, e = availabilityInputs()
	// Valid 512-byte SKU labels exceed the decoded budget before output, even
	// though every confirmed available SKU produces no detail text.
	i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": 1}
	values := map[string]SKUAvailability{}
	for n := 0; n < 33000; n++ { raw := fmt.Sprintf("%06d", n)+strings.Repeat("a", 506); values[raw] = SKUAvailability{State: SKUAvailable} }
	e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "complete", Values: values}
	availabilityReject(t, s, i, r, e, "input_limit")
}

func TestAvailabilityRuntimeOutputBounds(t *testing.T) {
	s, i, r, e := availabilityInputs()
	types := map[string]int64{}
	for n := 0; n < MaxDetails; n++ { types[fmt.Sprintf("microsoft.work/type%d", n)] = 0 }
	i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = types
	if got, err := CalculateAvailability(context.Background(), s, i, r, e); err != nil || len(got.Comparison.MissingResourceTypes) != MaxDetails { t.Fatal("exact detail boundary rejected", err) }
	types["microsoft.work/overflow"] = 0
	availabilityReject(t, s, i, r, e, "output_limit")
	s, i, r, e = availabilityInputs()
	i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": 1}
	i.Subscriptions[availabilityID].SKUsByTypeAndRegion["microsoft.compute/vms"] = map[string]map[string]int64{"eastus": {"Ready": 1}}
	// 21 type + colon + 5 SKU + 18 suffix + 467 zone =512 bytes.
	e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{"ready": {State: SKUZoneRestricted, BlockedZones: []string{strings.Repeat("z", 467)}}}}
	got, err := CalculateAvailability(context.Background(), s, i, r, e)
	if err != nil || len(got.Comparison.ZoneRestrictedSKUs[0]) != MaxLabelBytes { t.Fatal("exact joined detail boundary rejected", err) }
	e.SKUs["microsoft.compute/vms"].Values["ready"] = SKUAvailability{State: SKUZoneRestricted, BlockedZones: []string{strings.Repeat("z", 468)}}
	availabilityReject(t, s, i, r, e, "output_limit")
	s, i, r, e = availabilityInputs()
	i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": 1}
	i.Subscriptions[availabilityID].SKUsByTypeAndRegion["microsoft.compute/vms"] = map[string]map[string]int64{"eastus": {strings.Repeat("r", MaxLabelBytes): 1}}
	e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "unknown"}
	availabilityReject(t, s, i, r, e, "output_limit")
}

type availabilityCancelContext struct { context.Context; calls int; at int }
func (c *availabilityCancelContext) Err() error { c.calls++; if c.calls >= c.at { return context.Canceled }; return nil }

func TestAvailabilityRuntimeOwnershipAndCancellation(t *testing.T) {
	s, i, r, e := availabilityInputs()
	i.Subscriptions[availabilityID].ResourceTypesByRegion["eastus"] = map[string]int64{"microsoft.compute/vms": 1}
	i.Subscriptions[availabilityID].SKUsByTypeAndRegion["microsoft.compute/vms"] = map[string]map[string]int64{"eastus": {"Ready": 1}}
	e.SKUs["microsoft.compute/vms"] = SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{"ready": {State: SKUZoneRestricted, BlockedZones: []string{"3", "1", "3"}}}}
	before, _ := json.Marshal([]any{s, i, r, e})
	baseline, err := CalculateAvailability(context.Background(), s, i, r, e)
	if err != nil { t.Fatal(err) }
	var wg sync.WaitGroup
	for n := 0; n < 16; n++ { wg.Add(1); go func() { defer wg.Done(); got, err := CalculateAvailability(context.Background(), s, i, r, e); if err != nil || !reflect.DeepEqual(got, baseline) { t.Error("concurrent result leaked", err) } }() }
	wg.Wait()
	probe := &availabilityCancelContext{Context: context.Background(), at: 1000000}
	if _, err := CalculateAvailability(probe, s, i, r, e); err != nil { t.Fatal(err) }
	for _, at := range []int{1, 4, probe.calls} {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		got, err := CalculateAvailability(ctx, s, i, r, e)
		if got != nil || !errors.Is(err, context.Canceled) { t.Fatalf("cancelled work returned comparison at%d: %#v %v", at, got, err) }
	}
	baseline.Comparison.ZoneRestrictedSKUs[0] = "modified"
	baseline.Health.Warnings = append(baseline.Health.Warnings, assessment.AssessmentWarning{Code: "modified"})
	after, _ := json.Marshal([]any{s, i, r, e})
	if string(before) != string(after) { t.Fatal("input ownership violated") }
	got, err := CalculateAvailability(context.Background(), s, i, r, e)
	if err != nil || got.Comparison.ZoneRestrictedSKUs[0] != "microsoft.compute/vms:Ready (zones blocked: 3,1,3)" { t.Fatal("result aliases caller or previous run", err) }
}
