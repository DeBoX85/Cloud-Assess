package region

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// This helper exercises the existing consumer, without a Storage wire parser.
func storageConsumerCalculation(t *testing.T, input SKUEvidence) *AvailabilityCalculation {
	t.Helper()
	const id = "11111111-1111-1111-1111-111111111111"
	const resourceType = "microsoft.storage/storageaccounts"
	scope := map[string]string{id: "selected"}
	inventory := &InventoryCalculation{Subscriptions: map[string]InventoryCounts{id: {ResourceTypesByRegion: map[string]map[string]int64{"eastus": {resourceType: 3}}, SKUsByTypeAndRegion: map[string]map[string]map[string]int64{resourceType: {"eastus": {"a": 1, "b": 1, "c": 1}}}}}}
	request := AvailabilityRequest{SubscriptionID: id, SourceRegion: "eastus", TargetRegion: "westeurope"}
	evidence := AvailabilityEvidence{SubscriptionID: id, TargetRegion: "westeurope", InventoryComplete: true, ProvidersComplete: true, Locations: map[string]map[string]bool{resourceType: {"westeurope": true}}, SKUEnabled: true, SKUs: map[string]SKUEvidence{resourceType: input}, ZoneCounts: map[string]int{"eastus": 3, "westeurope": 3}}
	got, err := CalculateAvailability(context.Background(), scope, inventory, request, evidence)
	if err != nil || got == nil {
		t.Fatalf("consumer calculation failed: %v", err)
	}
	c := got.Comparison
	if c.SourceResourceTypeCount != 1 || c.AvailableTypes != 1 || c.UnavailableTypes != 0 || c.AvailabilityPercent != 100 || c.SourceZoneCount != 3 || c.TargetZoneCount != 3 || c.TotalSKUsChecked != 3 || len(c.ZoneRestrictedSKUs) != 0 || got.Health.Records != 1 {
		t.Fatalf("non-SKU consumer contract changed: %#v / %#v", c, got.Health)
	}
	return got
}

func TestStorageProjectionConsumerMixedContract(t *testing.T) {
	input := SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{"a": {State: SKUAvailable}, "b": {State: SKURestricted}, "c": {State: SKUUnavailable}}}
	got := storageConsumerCalculation(t, input)
	c := got.Comparison
	if c.AvailableSKUs != 1 || c.UnavailableSKUs != 1 || c.UnknownSKUs != 0 || math.Abs(c.SKUAvailabilityPercent-33.33333333333333) > 1e-9 {
		t.Fatalf("literal mixed counts/percentage mismatch: %#v", c)
	}
	if !reflect.DeepEqual(c.MissingSKUs, []string{"microsoft.storage/storageaccounts:c"}) || !reflect.DeepEqual(c.RestrictedSKUs, []string{"microsoft.storage/storageaccounts:b"}) {
		t.Fatalf("literal mixed details mismatch: %#v", c)
	}
	if got.Health.Status != assessment.StageCompleted || len(got.Health.Warnings) != 0 {
		t.Fatalf("complete Storage evidence changed health: %#v", got.Health)
	}
}

func TestStorageProjectionConsumerEmptyContract(t *testing.T) {
	got := storageConsumerCalculation(t, SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{}})
	c := got.Comparison
	if c.AvailableSKUs != 0 || c.UnavailableSKUs != 3 || c.UnknownSKUs != 0 || c.SKUAvailabilityPercent != 0 {
		t.Fatalf("literal empty counts/percentage mismatch: %#v", c)
	}
	if !reflect.DeepEqual(c.MissingSKUs, []string{"microsoft.storage/storageaccounts:a", "microsoft.storage/storageaccounts:b", "microsoft.storage/storageaccounts:c"}) || len(c.RestrictedSKUs) != 0 {
		t.Fatalf("literal empty details mismatch: %#v", c)
	}
	if got.Health.Status != assessment.StageCompleted || len(got.Health.Warnings) != 0 {
		t.Fatalf("complete-empty Storage evidence changed health: %#v", got.Health)
	}
}

func TestStorageProjectionConsumerUnknownContract(t *testing.T) {
	got := storageConsumerCalculation(t, SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}})
	c := got.Comparison
	// Source-compatible percentage100 must remain paired with unknown counts,
	// details and warning. It is not positive Storage availability evidence.
	if c.AvailableSKUs != 0 || c.UnavailableSKUs != 0 || c.UnknownSKUs != 3 || c.SKUAvailabilityPercent != 100 {
		t.Fatalf("literal unknown counts/percentage mismatch: %#v", c)
	}
	if !reflect.DeepEqual(c.MissingSKUs, []string{"microsoft.storage/storageaccounts:a (unknown)", "microsoft.storage/storageaccounts:b (unknown)", "microsoft.storage/storageaccounts:c (unknown)"}) || len(c.RestrictedSKUs) != 0 {
		t.Fatalf("literal unknown details mismatch: %#v", c)
	}
	if got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Code != "availability_sku_unknown" {
		t.Fatalf("unknown Storage evidence lost its warning: %#v", got.Health)
	}
}
