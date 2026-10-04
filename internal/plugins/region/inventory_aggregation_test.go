package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func aggregationFixture(t *testing.T) ([]assessment.Resource, []string, map[string]json.RawMessage) {
	t.Helper()
	var input struct {
		Resources []assessment.Resource
		Selectors []string
	}
	var output map[string]json.RawMessage
	for name, value := range map[string]any{"source-inventory-inputs.json": &input, "source-inventory-outputs.json": &output} {
		raw, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(raw, value); err != nil {
			t.Fatal(err)
		}
	}
	if len(input.Resources) != 10 || len(input.Selectors) != 5 || len(output) != 13 {
		t.Fatal("independent aggregation fixture topology changed")
	}
	return input.Resources, input.Selectors, output
}

func requireAggregationFailure(t *testing.T, got *InventoryCalculation, err error, code string) {
	t.Helper()
	want := "region inventory calculation input could not be safely aggregated [" + code + "]"
	if got != nil || err == nil || err.Error() != want || strings.Contains(err.Error(), "private") {
		t.Fatalf("wrong aggregation guard or partial result: expected=%s err=%v", code, err)
	}
}

func TestAggregationCapturedMapsAndIdentityCorrection(t *testing.T) {
	resources, selectors, captured := aggregationFixture(t)
	for _, index := range []int{1, 3} {
		input := []assessment.Resource{}
		for _, resource := range resources {
			if resource.SubscriptionID == selectors[index] {
				input = append(input, resource)
			}
		}
		var want InventoryCounts
		if err := json.Unmarshal(captured[fmt.Sprintf("build-selected-%d", index)], &want); err != nil {
			t.Fatal(err)
		}
		id := strings.ToLower(selectors[index])
		got, err := CalculateInventory(context.Background(), map[string]string{id: "Synthetic"}, input)
		if err != nil || got == nil || !reflect.DeepEqual(got.Subscriptions[id], want) || !reflect.DeepEqual(got.Aggregate, want) || !slices.Equal(got.ContributorSubscriptionIDs, []string{id}) || got.Records != len(input) {
			t.Fatal("complete actual five-map oracle or canonical case correction changed", index, err)
		}
		if index == 3 {
			var sourceLower InventoryCounts
			if err := json.Unmarshal(captured["build-selected-2"], &sourceLower); err != nil || len(sourceLower.ResourceTypes) != 0 || len(got.Aggregate.ResourceTypes) != 1 {
				t.Fatal("source exact UUID filtering discrepancy was hidden", err)
			}
		}
	}
	var empty InventoryCounts
	if err := json.Unmarshal(captured["build-empty"], &empty); err != nil {
		t.Fatal(err)
	}
	got, err := CalculateInventory(context.Background(), inventoryScope(), nil)
	if err != nil || got == nil || !reflect.DeepEqual(got.Aggregate, empty) || len(got.Subscriptions) != 2 || len(got.ContributorSubscriptionIDs) != 0 || got.Records != 0 {
		t.Fatal("actual empty aggregate/no-contributor shape changed", err)
	}
	for _, counts := range got.Subscriptions {
		if !reflect.DeepEqual(counts, empty) {
			t.Fatal("selected empty five maps changed")
		}
	}
	// The full captured source fixture contains a control-bearing tab location.
	// Target deliberately rejects it, without pretending that row was normalized.
	got, err = CalculateInventory(context.Background(), map[string]string{selectors[0]: "Synthetic", selectors[1]: "Synthetic", selectors[2]: "Synthetic"}, resources)
	requireAggregationFailure(t, got, err, "text_limit")
}

func TestAggregationSourceNormalization(t *testing.T) {
	_, _, output := aggregationFixture(t)
	var expected []map[string]any
	if err := json.Unmarshal(output["normalization"], &expected); err != nil || len(expected) != 30 {
		t.Fatal("source normalization oracle changed", err)
	}
	for i, record := range expected {
		if normalizeInventoryRegion(record["region"].(string)) != record["normalized"].(string) {
			t.Fatal("actual source ASCII spaces/Unicode/case normalization changed", i)
		}
	}
}

func TestAggregationLiteralCountsAndIgnoredFields(t *testing.T) {
	resources, selectors, _ := aggregationFixture(t)
	valid := []assessment.Resource{}
	for _, resource := range resources {
		if resource.SubscriptionID == selectors[0] && resource.Location != "East\tUS" {
			// Only consumed fields belong to this arithmetic contract.
			resource.Name, resource.ResourceGroup, resource.Kind, resource.SKUTier, resource.SKUFamily, resource.SLA = "private\x00", "private\x00", "private\x00", "private\x00", "private\x00", "private\x00"
			resource.SKUCapacity = 9223372036854775807
			resource.Tags = map[string]string{"private\x00": "private\x00"}
			valid = append(valid, resource)
		}
	}
	want := InventoryCounts{
		ResourceTypes:         map[string]int64{"microsoft.compute/virtualmachines": 3, "microsoft.compute/virtualmachinescalesets": 1, "microsoft.custom/widgets": 1, "microsoft.network/virtualnetworks": 1},
		SKUsByType:            map[string]map[string]int64{"microsoft.compute/virtualmachines": {"Standard_D4s_v3": 2, " Standard_D4s_v3 ": 1}, "microsoft.compute/virtualmachinescalesets": {"Standard_D4s_v3": 1}, "microsoft.custom/widgets": {"=Widget": 1}},
		LocationCounts:        map[string]int64{"": 1, "eastus": 3, "europe": 1, "westeurope": 1},
		ResourceTypesByRegion: map[string]map[string]int64{"": {"microsoft.custom/widgets": 1}, "eastus": {"microsoft.compute/virtualmachines": 2, "microsoft.network/virtualnetworks": 1}, "europe": {"microsoft.compute/virtualmachines": 1}, "westeurope": {"microsoft.compute/virtualmachinescalesets": 1}},
		SKUsByTypeAndRegion:   map[string]map[string]map[string]int64{"microsoft.compute/virtualmachines": {"eastus": {"Standard_D4s_v3": 2}, "europe": {" Standard_D4s_v3 ": 1}}, "microsoft.compute/virtualmachinescalesets": {"westeurope": {"Standard_D4s_v3": 1}}, "microsoft.custom/widgets": {"": {"=Widget": 1}}},
	}
	got, err := CalculateInventory(context.Background(), map[string]string{selectors[0]: "Synthetic"}, valid)
	if err != nil || got.Records != 6 || !reflect.DeepEqual(got.Aggregate, want) || !reflect.DeepEqual(got.Subscriptions[selectors[0]], want) {
		t.Fatal("literal raw SKU/empty/logical/type/resource-not-capacity maps changed", err)
	}
	// The excluded tab row is explicit; no source oracle bytes were rewritten.
	if len(valid) != 6 {
		t.Fatal("literal valid resource subset changed")
	}
}

func TestAggregationSelectedScopeAndMalformed(t *testing.T) {
	resources, _, _ := aggregationFixture(t)
	input := []assessment.Resource{resources[0]}
	for _, scope := range []map[string]string{nil, {otherAuxSubscription: "Synthetic"}, {auxSubscription: ""}, {auxSubscription: "private\x00"}, {"bad": "Synthetic"}} {
		got, err := CalculateInventory(context.Background(), scope, input)
		requireAggregationFailure(t, got, err, "scope_invalid")
	}
	got, err := CalculateInventory(context.Background(), map[string]string{auxSubscription: ""}, nil)
	requireAggregationFailure(t, got, err, "scope_invalid")
	alias := "ABCDEFAB-1111-1111-1111-111111111111"
	got, err = CalculateInventory(context.Background(), map[string]string{alias: "Synthetic", strings.ToLower(alias): "Synthetic"}, nil)
	requireAggregationFailure(t, got, err, "scope_invalid")
	r := resources[0]
	r.ID = strings.Replace(r.ID, auxSubscription, otherAuxSubscription, 1)
	got, err = CalculateInventory(context.Background(), inventoryScope(), []assessment.Resource{r})
	requireAggregationFailure(t, got, err, "scope_invalid")
	r = resources[0]
	duplicate := r
	duplicate.ID = strings.Replace(r.ID, "/rg/", "/RG/", 1)
	got, err = CalculateInventory(context.Background(), inventoryScope(), []assessment.Resource{r, duplicate})
	requireAggregationFailure(t, got, err, "input_invalid")
	for _, invalid := range []string{"private\x00", "private\ufffd", "private\ufffe", "private\uffff", string([]byte{0xff}), strings.Repeat("s", 513)} {
		r = resources[0]
		r.SKUName = invalid
		got, err = CalculateInventory(context.Background(), inventoryScope(), []assessment.Resource{r})
		requireAggregationFailure(t, got, err, "text_limit")
	}
	for _, field := range []string{"id", "type"} {
		r = resources[0]
		if field == "id" {
			r.ID = ""
		} else {
			r.Type = ""
		}
		got, err = CalculateInventory(context.Background(), inventoryScope(), []assessment.Resource{r})
		requireAggregationFailure(t, got, err, "input_invalid")
	}
	// Unicode casing can expand a safe raw key from512 to768 UTF8 bytes.
	r = resources[0]
	r.Type = strings.Repeat("\u023a", 256)
	if len(r.Type) != 512 || len(strings.ToLower(r.Type)) != 768 {
		t.Fatal("literal Unicode casing expansion premise changed")
	}
	got, err = CalculateInventory(context.Background(), inventoryScope(), []assessment.Resource{r})
	requireAggregationFailure(t, got, err, "text_limit")
}

func aggregationMany(t *testing.T, count int, mode string) []assessment.Resource {
	t.Helper()
	fixture, _, _ := aggregationFixture(t)
	resources := make([]assessment.Resource, count)
	for i := range resources {
		r := fixture[0]
		r.ID = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Test/widgets/r%05d", auxSubscription, i)
		if mode == "entries" || mode == "projected" {
			r.Type, r.Location, r.SKUName = fmt.Sprintf("Microsoft.Test/t%05d", i), fmt.Sprintf("r%05d", i), fmt.Sprintf("sku%05d", i)
		}
		if mode == "projected" {
			r.Type += strings.Repeat("t", 512-len(r.Type))
			r.Location += strings.Repeat("l", 512-len(r.Location))
			r.SKUName += strings.Repeat("s", 512-len(r.SKUName))
		}
		if mode == "decoded" {
			r.ID += strings.Repeat("i", 512-len(r.ID))
			r.Type, r.Location, r.SKUName = strings.Repeat("s", 512), strings.Repeat("s", 512), strings.Repeat("s", 512)
		}
		resources[i] = r
	}
	return resources
}

func TestAggregationWorkAndEntryLimits(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	resources := aggregationMany(t, 8192, "common")
	got, err := CalculateInventory(context.Background(), scope, resources)
	if err != nil || got.Records != 8192 || got.Aggregate.ResourceTypes["microsoft.compute/virtualmachines"] != 8192 {
		t.Fatal("literal8192 resource boundary rejected", err)
	}
	extra := resources[0]
	extra.ID += "overflow"
	got, err = CalculateInventory(context.Background(), scope, append(resources, extra))
	requireAggregationFailure(t, got, err, "input_limit")
	selection := map[string]string{}
	for i := 0; i < 1000; i++ {
		selection[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "Synthetic"
	}
	got, err = CalculateInventory(context.Background(), selection, nil)
	if err != nil || len(got.Subscriptions) != 1000 {
		t.Fatal("literal1000 selected empty inventories rejected", err)
	}
	selection[auxSubscription] = "Synthetic"
	got, err = CalculateInventory(context.Background(), selection, nil)
	requireAggregationFailure(t, got, err, "scope_invalid")
	// 3640 distinct type/location/SKU records:9 keys per five-map inventory,
	// doubled per-sub/aggregate. Three selected IDs+two contributors+4 no-SKU
	// per-sub keys+5 first-SKU per-sub keys+2 new aggregate SKU keys =65536.
	resources = aggregationMany(t, 3640, "entries")
	one := resources[0]
	one.SubscriptionID = otherAuxSubscription
	one.ID = "/subscriptions/" + otherAuxSubscription + "/providers/Microsoft.Test/widgets/no-sku"
	one.SKUName = ""
	two := one
	two.ID = "/subscriptions/" + otherAuxSubscription + "/providers/Microsoft.Test/widgets/new-sku"
	two.SKUName = "new-sku"
	resources = append(resources, one, two)
	selection = map[string]string{auxSubscription: "Synthetic", otherAuxSubscription: "Synthetic", "33333333-3333-3333-3333-333333333333": "Synthetic"}
	got, err = CalculateInventory(context.Background(), selection, resources)
	if err != nil || got.Records != 3642 || len(got.ContributorSubscriptionIDs) != 2 {
		t.Fatal("independent literal65536 output-entry boundary rejected", err)
	}
	extra = resources[0]
	extra.ID += "overflow"
	extra.SKUName = "overflow-sku"
	got, err = CalculateInventory(context.Background(), selection, append(resources, extra))
	requireAggregationFailure(t, got, err, "output_limit")
}

func TestAggregationDecodedAndProjectedText(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	// Repeated type/region/SKU means output keys are tiny, but decoded raw and
	// normalized strings consume5400*(6*512+36)=16783200 bytes before scope.
	got, err := CalculateInventory(context.Background(), scope, aggregationMany(t, 5400, "decoded"))
	requireAggregationFailure(t, got, err, "text_limit")
	// Distinct safe512-byte keys produce18*512 output bytes/resource.
	// 1813*9216+72=16708680 fits;1814*9216+72=16717896 exceeds16711680.
	got, err = CalculateInventory(context.Background(), scope, aggregationMany(t, 1813, "projected"))
	if err != nil || got.Records != 1813 {
		t.Fatal("literal projected key-text boundary rejected", err)
	}
	got, err = CalculateInventory(context.Background(), scope, aggregationMany(t, 1814, "projected"))
	requireAggregationFailure(t, got, err, "output_limit")
}

func TestAggregationOwnershipConcurrencyAndCancellation(t *testing.T) {
	resources, selectors, _ := aggregationFixture(t)
	input := []assessment.Resource{}
	for _, resource := range resources {
		if resource.SubscriptionID == selectors[1] || resource.SubscriptionID == selectors[3] {
			input = append(input, resource)
		}
	}
	scope := map[string]string{selectors[1]: "Synthetic", selectors[2]: "Synthetic"}
	original := slices.Clone(input)
	baseline, err := CalculateInventory(context.Background(), scope, input)
	if err != nil || !reflect.DeepEqual(input, original) || !slices.Equal(baseline.ContributorSubscriptionIDs, []string{selectors[1], selectors[2]}) || baseline.Aggregate.ResourceTypes["microsoft.compute/virtualmachines"] != 3 {
		t.Fatal("owned aggregate/actual contributor ordering/input changed", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := CalculateInventory(context.Background(), scope, input)
			if err != nil || !reflect.DeepEqual(got, baseline) {
				t.Error("aggregation per-run maps or counts leaked")
				return
			}
			counts := got.Subscriptions[selectors[1]]
			counts.ResourceTypes["microsoft.compute/virtualmachines"] = 999
			counts.SKUsByType["microsoft.compute/virtualmachines"]["Standard_F4s_v2"] = 999
			counts.SKUsByTypeAndRegion["microsoft.compute/virtualmachines"]["eastus"]["Standard_F4s_v2"] = 999
			got.ContributorSubscriptionIDs[0] = "owned"
			if got.Aggregate.ResourceTypes["microsoft.compute/virtualmachines"] != 3 || got.Aggregate.SKUsByType["microsoft.compute/virtualmachines"]["Standard_F4s_v2"] != 1 || got.Aggregate.SKUsByTypeAndRegion["microsoft.compute/virtualmachines"]["eastus"]["Standard_F4s_v2"] != 1 {
				t.Error("aggregate aliases a subscription counter")
			}
		}()
	}
	wg.Wait()
	many := aggregationMany(t, 100, "common")
	for _, at := range []int{1, 4, 8, 20, 80, 160} {
		ctx, cancel := context.WithCancel(context.Background())
		checking := &serviceCancelContext{Context: ctx, cancel: cancel, at: at}
		got, err := CalculateInventory(checking, map[string]string{auxSubscription: "Synthetic"}, many)
		cancel()
		if got != nil || !errors.Is(err, context.Canceled) {
			t.Fatal("aggregation cancellation returned partial result or lost identity", at, err)
		}
	}
	input[0].Type, scope[selectors[1]] = "changed", "changed"
	if baseline.Aggregate.ResourceTypes["microsoft.compute/virtualmachines"] != 3 || baseline.ContributorSubscriptionIDs[0] != selectors[1] {
		t.Fatal("aggregation output retained mutable input aliases")
	}
}
