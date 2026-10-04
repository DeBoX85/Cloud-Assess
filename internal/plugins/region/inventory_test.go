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

func inventoryFixture(t *testing.T) ([]assessment.Resource, map[string]auxiliaryCapturedTable) {
	t.Helper()
	var input struct{ Resources []assessment.Resource }
	var outputs map[string]json.RawMessage
	for name, value := range map[string]any{"source-aux-inputs.json": &input, "source-aux-outputs.json": &outputs} {
		raw, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(raw, value); err != nil {
			t.Fatal(err)
		}
	}
	expected := map[string]auxiliaryCapturedTable{}
	for _, name := range []string{"inventory-raw", "inventory-masked", "inventory-empty"} {
		var table auxiliaryCapturedTable
		if err := json.Unmarshal(outputs[name], &table); err != nil {
			t.Fatal(err)
		}
		expected[name] = table
	}
	if len(outputs) != 14 || len(input.Resources) != 2 || len(expected["inventory-raw"].Table) != 3 || len(expected["inventory-masked"].Table) != 3 || len(expected["inventory-empty"].Table) != 1 {
		t.Fatal("independent inventory capture topology changed")
	}
	return input.Resources, expected
}

func inventoryScope() map[string]string {
	return map[string]string{auxSubscription: "Synthetic", otherAuxSubscription: "Synthetic"}
}

func requireInventoryFailure(t *testing.T, table *assessment.PluginTable, err error, code string) {
	t.Helper()
	if err == nil || table == nil || len(table.Rows) != 0 || len(table.Columns) != 10 || table.Health.Status != assessment.StageFailed || table.Health.Error == nil || table.Health.Error.Code != code || strings.Contains(err.Error(), "private") || assessment.ValidatePluginTables([]assessment.PluginTable{*table}) != nil {
		t.Fatalf("unsafe inventory input, guard or partial output: expected=%s err=%v", code, err)
	}
}

func TestInventoryCapturedCellsAndComposition(t *testing.T) {
	resources, expected := inventoryFixture(t)
	for _, masked := range []bool{false, true} {
		name := "inventory-raw"
		if masked {
			name = "inventory-masked"
		}
		got, err := ProjectInventory(context.Background(), inventoryScope(), resources, masked)
		want := expected[name]
		if err != nil || got == nil || got.SheetName != "Region Inventory" || want.SheetName != "Inventory" || got.Description != want.Description || got.Metadata != Metadata() || got.ID != "inventory" || !slices.Equal(got.Columns, want.Table[0]) || len(got.Rows) != 2 || got.Health.Status != assessment.StageCompleted || got.Health.Records != 2 || assessment.ValidatePluginTables([]assessment.PluginTable{*got}) != nil {
			t.Fatal("captured inventory metadata/composition/health changed", err)
		}
		for i, row := range got.Rows {
			if !slices.Equal(row.Cells, want.Table[i+1]) || row.SubscriptionID != resources[i].SubscriptionID {
				t.Fatal("captured raw/masked inventory cell or correlation changed", i)
			}
		}
		got.SheetName = want.SheetName
		if assessment.ValidatePluginTables([]assessment.PluginTable{*got}) == nil {
			t.Fatal("source Inventory collision was silently accepted")
		}
	}
	for _, empty := range [][]assessment.Resource{nil, {}} {
		got, err := ProjectInventory(context.Background(), inventoryScope(), empty, true)
		if err != nil || got == nil || len(got.Rows) != 0 || !slices.Equal(got.Columns, expected["inventory-empty"].Table[0]) || got.Health.Status != assessment.StageCompleted {
			t.Fatal("source header-only inventory branch changed", err)
		}
	}
	// Source received ordering must survive; no sort or core SLA column.
	resources[0], resources[1] = resources[1], resources[0]
	got, err := ProjectInventory(context.Background(), inventoryScope(), resources, false)
	if err != nil || got.Rows[0].Cells[4] != "set" || got.Rows[1].Cells[4] != "=vm" || len(got.Columns) != 10 {
		t.Fatal("source input order/formula text/width changed", err)
	}
}

func TestInventorySourceCapacityAndBoundaries(t *testing.T) {
	resources, _ := inventoryFixture(t)
	for _, test := range []struct {
		sku, kind string
		capacity  int
		want      string
	}{
		{"Standard_D4s_v3", "Microsoft.Compute/virtualMachines", 1, "1"},
		{"Standard_D4s_v3", "Microsoft.Compute/virtualMachines", 0, "4"},
		{"Standard_D4s_v3", "Microsoft.Compute/virtualMachines", -1000000000000, "4"},
		{"Standard_D4s_v3", "MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS", 2, "8"},
		{"Standard_D4s_v3", "Microsoft.Compute/virtualMachineScaleSets", 0, "4"},
		{"standard_d4s_v3", "Microsoft.Compute/virtualMachineScaleSets", 2, "2"},
		{"standard_d4s_v3", "Microsoft.Compute/virtualMachines", 0, ""},
		{"unknown", "Microsoft.Compute/virtualMachines", -1, ""},
		{"unknown", "Microsoft.Compute/virtualMachines", 1000000000000, "1000000000000"},
		{"Standard_D4s_v3", "Microsoft.Compute/virtualMachineScaleSets", 250000000000, "1000000000000"},
	} {
		r := resources[0]
		r.SKUName, r.Type, r.SKUCapacity = test.sku, test.kind, test.capacity
		got, err := ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, false)
		if err != nil || got.Rows[0].Cells[7] != test.want {
			t.Fatal("literal source capacity/fallback/case/product boundary changed", test, err)
		}
	}
	for _, count := range []int{250000000001, 1000000000001, -1000000000001, 9223372036854775807} {
		r := resources[1]
		r.SKUCapacity = count
		got, err := ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, false)
		requireInventoryFailure(t, got, err, "region_inventory_capacity_invalid")
	}
}

func TestInventorySelectedScopeAndIdentity(t *testing.T) {
	resources, _ := inventoryFixture(t)
	for _, scope := range []map[string]string{nil, {auxSubscription: "Synthetic"}, {auxSubscription: ""}, {auxSubscription: "private\x00"}, {"private-not-uuid": "Synthetic"}} {
		got, err := ProjectInventory(context.Background(), scope, resources, true)
		requireInventoryFailure(t, got, err, "region_inventory_scope_invalid")
	}
	got, err := ProjectInventory(context.Background(), map[string]string{auxSubscription: ""}, nil, false)
	requireInventoryFailure(t, got, err, "region_inventory_scope_invalid")
	alias := "ABCDEFAB-1111-1111-1111-111111111111"
	got, err = ProjectInventory(context.Background(), map[string]string{alias: "Synthetic", strings.ToLower(alias): "Synthetic"}, nil, false)
	requireInventoryFailure(t, got, err, "region_inventory_scope_invalid")
	for _, id := range []string{
		"/subscriptions/" + otherAuxSubscription + "/providers/Microsoft.Test/widgets/private",
		"/Subscriptions/" + auxSubscription + "/providers/Microsoft.Test/widgets/private",
		"/subscriptions/" + auxSubscription,
		"/subscriptions/" + auxSubscription + "/",
		"/subscriptions/" + auxSubscription + "x/providers/Microsoft.Test/widgets/private",
		"/subscriptions/" + auxSubscription + "/providers//widgets/private",
		"/subscriptions/" + auxSubscription + "/providers/Microsoft.Test/widgets/private/",
		"/subscriptions/" + auxSubscription + "/providers/Microsoft.Test/widgets/private?token",
		"/subscriptions/" + auxSubscription + "/providers/Microsoft.Test/widgets/private#token",
		"/subscriptions/" + auxSubscription + "/providers/Microsoft.Test/widgets/private\\token",
		"/providers/Microsoft.Test/widgets/private",
	} {
		r := resources[0]
		r.ID = id
		got, err = ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, true)
		requireInventoryFailure(t, got, err, "region_inventory_scope_invalid")
	}
	// Both subscription-level and child/extension suffixes are retained.
	for _, suffix := range []string{"providers/Microsoft.Test/widgets/private", "resourceGroups/rg/providers/Microsoft.Test/widgets/private/children/child/providers/Microsoft.Ext/types/x"} {
		r := resources[0]
		r.ID = "/subscriptions/" + auxSubscription + "/" + suffix
		got, err = ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, false)
		if err != nil || got.Rows[0].Cells[9] != r.ID {
			t.Fatal("valid selected subscription/extension suffix changed", err)
		}
	}
	// Uppercase UUIDs are selected by identity, and retain source displayed case.
	r := resources[0]
	r.SubscriptionID = alias
	r.ID = "/subscriptions/" + strings.ToLower(alias) + "/providers/Microsoft.Test/widgets/private"
	got, err = ProjectInventory(context.Background(), map[string]string{alias: "Synthetic"}, []assessment.Resource{r}, false)
	if err != nil || got.Rows[0].SubscriptionID != strings.ToLower(alias) || got.Rows[0].Cells[0] != alias {
		t.Fatal("case-insensitive identity/source display changed", err)
	}
}

func TestInventoryMalformedAndDuplicates(t *testing.T) {
	resources, _ := inventoryFixture(t)
	for _, invalid := range []string{"private\x00", "private\ufffd", "private\ufffe", "private\uffff", string([]byte{0xff}), strings.Repeat("s", 513)} {
		r := resources[0]
		r.Name = invalid
		got, err := ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, false)
		requireInventoryFailure(t, got, err, "region_inventory_text_limit")
	}
	for _, field := range []string{"id", "type", "name"} {
		r := resources[0]
		switch field {
		case "id":
			r.ID = ""
		case "type":
			r.Type = ""
		case "name":
			r.Name = ""
		}
		got, err := ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, false)
		requireInventoryFailure(t, got, err, "region_inventory_input_invalid")
	}
	duplicate := resources[0]
	duplicate.ID = strings.Replace(duplicate.ID, "/rg/", "/RG/", 1)
	got, err := ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{resources[0], duplicate}, false)
	requireInventoryFailure(t, got, err, "region_inventory_input_invalid")
	r := resources[0]
	r.Name, r.SKUName = strings.Repeat("\U0001f600", 128), strings.Repeat("\u00e9", 256)
	r.ResourceGroup, r.Location, r.SKUTier, r.Kind = "", "", "", ""
	got, err = ProjectInventory(context.Background(), inventoryScope(), []assessment.Resource{r}, true)
	if err != nil || got.Rows[0].Cells[4] != strings.Repeat("\U0001f600", 128) || got.Rows[0].Cells[5] != strings.Repeat("\u00e9", 256) || got.Rows[0].Cells[1] != "" || got.Rows[0].Cells[2] != "" {
		t.Fatal("optional/safe512-byte Unicode labels changed", err)
	}
}

func inventoryMany(t *testing.T, count int, text bool) []assessment.Resource {
	t.Helper()
	fixture, _ := inventoryFixture(t)
	resources := make([]assessment.Resource, count)
	for i := range resources {
		r := fixture[0]
		r.ID = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Test/widgets/%05d", auxSubscription, i)
		if text {
			r.ID += strings.Repeat("i", 512-len(r.ID))
			r.ResourceGroup, r.Location, r.Type, r.Name, r.SKUName, r.SKUTier, r.Kind = strings.Repeat("s", 512), strings.Repeat("s", 512), strings.Repeat("s", 512), strings.Repeat("s", 512), strings.Repeat("s", 512), strings.Repeat("s", 512), strings.Repeat("s", 512)
		}
		resources[i] = r
	}
	return resources
}

func TestInventoryWorkAndTextLimits(t *testing.T) {
	resources := inventoryMany(t, 8192, false)
	got, err := ProjectInventory(context.Background(), inventoryScope(), resources, false)
	if err != nil || len(got.Rows) != 8192 {
		t.Fatal("literal8192-resource boundary rejected", err)
	}
	extra := resources[0]
	extra.ID += "overflow"
	resources = append(resources, extra)
	got, err = ProjectInventory(context.Background(), inventoryScope(), resources, false)
	requireInventoryFailure(t, got, err, "region_inventory_input_limit")
	scope := map[string]string{}
	for i := 0; i < 1000; i++ {
		scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "Synthetic"
	}
	got, err = ProjectInventory(context.Background(), scope, nil, false)
	if err != nil || len(got.Rows) != 0 {
		t.Fatal("literal1000-subscription empty inventory boundary rejected", err)
	}
	scope[auxSubscription] = "Synthetic"
	got, err = ProjectInventory(context.Background(), scope, nil, false)
	requireInventoryFailure(t, got, err, "region_inventory_scope_invalid")
	got, err = ProjectInventory(context.Background(), inventoryScope(), inventoryMany(t, 6000, true), false)
	requireInventoryFailure(t, got, err, "region_inventory_text_limit")
	// Decoded scope text is absent from displayed cells: this case isolates the
	// early decoded budget from the projected/canonical output guard.
	scope = map[string]string{auxSubscription: strings.Repeat("n", 512)}
	for i := 0; i < 999; i++ {
		scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = strings.Repeat("n", 512)
	}
	got, err = ProjectInventory(context.Background(), scope, inventoryMany(t, 4000, true), false)
	requireInventoryFailure(t, got, err, "region_inventory_text_limit")
}

func TestInventoryOwnershipConcurrencyAndCancellation(t *testing.T) {
	resources, _ := inventoryFixture(t)
	scope := inventoryScope()
	original := slices.Clone(resources)
	baseline, err := ProjectInventory(context.Background(), scope, resources, true)
	if err != nil || !reflect.DeepEqual(resources, original) {
		t.Fatal("inventory input mutated", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := ProjectInventory(context.Background(), scope, resources, true)
			if err != nil || !reflect.DeepEqual(got, baseline) {
				t.Error("inventory per-run state/order leaked")
				return
			}
			got.Columns[0] = "owned"
			got.Rows[0].Cells[0] = "owned"
		}()
	}
	wg.Wait()
	for _, at := range []int{1, 4, 6, 7, 8, 9} {
		ctx, cancel := context.WithCancel(context.Background())
		checking := &serviceCancelContext{Context: ctx, cancel: cancel, at: at}
		got, err := ProjectInventory(checking, scope, resources, true)
		cancel()
		requireInventoryFailure(t, got, err, "region_inventory_cancelled")
		if !errors.Is(err, context.Canceled) {
			t.Fatal("inventory context identity changed")
		}
	}
	resources[0].ID, resources[0].Name, scope[auxSubscription] = "changed", "changed", "changed"
	if baseline.Columns[0] != "Subscription Id" || baseline.Rows[0].Cells[4] != "=vm" || baseline.Rows[0].SubscriptionID != auxSubscription {
		t.Fatal("inventory owned output retained input aliases")
	}
}
