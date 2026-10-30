package region

import (
	"reflect"
	"testing"
)

// Expected values are literals from the already-frozen ownership contract,
// never snapshots returned by the API being checked.
func TestSubscriptionLocationsAccessorLiteralSnapshot(t *testing.T) {
	g := &locGetter{pages: []locPage{{Body: locBody(locWire("eastus", `[{"logicalZone":"1","physicalZone":"eastus-az2"}]`), ""), Status: 200}}}
	got, err := locCollect(t, g)
	if err != nil || got == nil {
		t.Fatal("LOC_ACCESSOR_SETUP", err)
	}
	first := got.Evidence()
	first.Regions[0] = "mutated"
	first.ZoneCounts["eastus"] = 999
	first.ZoneMappings["eastus"]["1"] = "mutated"
	fresh := got.Evidence()
	if fresh.SubscriptionID != locSub || fresh.Status != "complete" || fresh.FailureCode != "" || fresh.CountSource != "api" || !reflect.DeepEqual(fresh.Regions, []string{"eastus"}) || !reflect.DeepEqual(fresh.ZoneCounts, map[string]int{"eastus": 1}) || !reflect.DeepEqual(fresh.ZoneMappings, map[string]map[string]string{"eastus": {"1": "eastus-az2"}}) {
		t.Fatal("LOC_ACCESSOR_SNAPSHOT", fresh)
	}
}
