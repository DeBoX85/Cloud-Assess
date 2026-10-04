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

func serviceFixture(t *testing.T) (*ServiceInventory, []ServiceComparison, []auxiliaryCapturedTable) {
	t.Helper()
	var input struct {
		Inventory   *ServiceInventory
		Comparisons []ServiceComparison
	}
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
	var tables []auxiliaryCapturedTable
	if err := json.Unmarshal(outputs["service-full"], &tables); err != nil {
		t.Fatal(err)
	}
	if input.Inventory == nil || len(input.Comparisons) != 2 || len(tables) != 2 || len(outputs) != 14 || string(outputs["service-nil-inventory"]) != "null" || string(outputs["service-no-results"]) != "null" {
		t.Fatal("independent service fixture topology changed")
	}
	input.Inventory.SubscriptionIDs = []string{auxSubscription}
	for i := range input.Comparisons {
		input.Comparisons[i].SubscriptionID = auxSubscription
		input.Comparisons[i].SubscriptionName = "Synthetic"
		input.Comparisons[i].SourceRegion = "eastus"
	}
	return input.Inventory, input.Comparisons, tables
}

func requireServiceFailure(t *testing.T, tables []assessment.PluginTable, err error, code string) {
	t.Helper()
	if err == nil || len(tables) != 1 || len(tables[0].Rows) != 0 || len(tables[0].Columns) != 8 || tables[0].Health.Status != assessment.StageFailed || tables[0].Health.Error == nil || tables[0].Health.Error.Code != code || strings.Contains(err.Error(), "private") || assessment.ValidatePluginTables(tables) != nil {
		t.Fatalf("unsafe input accepted, wrong guard, leaked or retained partial rows: code=%s tables=%#v err=%v", code, tables, err)
	}
}

func TestServiceCapturedEveryCellAndEmpty(t *testing.T) {
	inventory, comparisons, expected := serviceFixture(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	tables, err := ProjectServiceAvailability(context.Background(), scope, inventory, comparisons)
	if err != nil || len(tables) != len(expected) || assessment.ValidatePluginTables(tables) != nil {
		t.Fatal("canonical captured service projection failed", err)
	}
	for i, table := range tables {
		want := expected[i]
		if table.SheetName != want.SheetName || table.Description != want.Description || table.Metadata != Metadata() || !slices.Equal(table.Columns, want.Table[0]) || len(table.Rows) != len(want.Table)-1 || table.Health.Status != assessment.StageCompleted || table.Health.Records != len(table.Rows) {
			t.Fatal("source metadata/shape mismatch", i)
		}
		for j, row := range table.Rows {
			if !slices.Equal(row.Cells, want.Table[j+1]) || row.SubscriptionID != "" {
				t.Fatal("source cell mismatch or fabricated aggregate identity", i, j)
			}
		}
	}
	for _, pair := range []struct {
		inventory   *ServiceInventory
		comparisons []ServiceComparison
	}{{nil, comparisons}, {inventory, nil}, {nil, nil}} {
		tables, err = ProjectServiceAvailability(context.Background(), scope, pair.inventory, pair.comparisons)
		if err != nil || tables != nil {
			t.Fatal("source absent-empty helper branch changed")
		}
	}
	empty := &ServiceInventory{SubscriptionIDs: []string{auxSubscription}}
	tables, err = ProjectServiceAvailability(context.Background(), scope, empty, comparisons)
	if err != nil || len(tables) != 2 || len(tables[0].Rows) != 0 || tables[0].Health.Status != assessment.StageCompleted {
		t.Fatal("non-nil empty inventory must retain header-only target sheets")
	}
}

func TestServiceRegistryAndDetailSemantics(t *testing.T) {
	keys := []string{"microsoft.compute/virtualmachines", "microsoft.compute/virtualmachinescalesets", "microsoft.compute/disks", "microsoft.sql/servers/databases", "microsoft.sql/managedinstances", "microsoft.storage/storageaccounts", "microsoft.cognitiveservices/accounts"}
	inventory := &ServiceInventory{SubscriptionIDs: []string{auxSubscription}, ResourceTypes: map[string]int64{}, SKUsByType: map[string]map[string]int64{}, ResourceTypesByRegion: map[string]map[string]int64{"westus": {}, "eastus": {}}}
	c := ServiceComparison{SubscriptionID: auxSubscription, SubscriptionName: "Synthetic", SourceRegion: "westus", TargetRegion: "eastus"}
	for _, rt := range keys {
		inventory.ResourceTypes[rt] = 2
		inventory.SKUsByType[rt] = map[string]int64{"Z": 0, "A": 2}
		inventory.ResourceTypesByRegion["eastus"][rt] = 2
		inventory.ResourceTypesByRegion["westus"][rt] = 0
		c.MissingSKUs = append(c.MissingSKUs, strings.ToUpper(rt+":Z"))
		c.ZoneRestrictedSKUs = append(c.ZoneRestrictedSKUs, strings.ToUpper(rt+":A")+" (zones blocked: 2)")
	}
	// Literal expected registry support, not the production lookup as an oracle.
	inventory.ResourceTypes["microsoft.custom/widgets"] = 1
	inventory.SKUsByType["microsoft.custom/widgets"] = map[string]int64{"=Widget": 1}
	c.MissingSKUs = append(c.MissingSKUs, "microsoft.custom/widgets:=Widget")
	c.ZoneRestrictedSKUs = append(c.ZoneRestrictedSKUs, "microsoft.custom/widgets:=Widget")
	tables, err := ProjectServiceAvailability(context.Background(), map[string]string{auxSubscription: "Synthetic"}, inventory, []ServiceComparison{c})
	if err != nil || len(tables[0].Rows) != 8 {
		t.Fatal(err)
	}
	for _, row := range tables[0].Rows {
		if row.Cells[0] == "microsoft.custom/widgets" {
			if !slices.Equal(row.Cells[4:7], []string{"=Widget", "N/A", "N/A"}) {
				t.Fatal("unsupported SKU provider fabricated a decision")
			}
		} else if !slices.Equal(row.Cells[1:], []string{"2", "eastus, westus", "2", "Z", "Not available", "A", "Available"}) {
			t.Fatal("pinned registry, case, zero-count membership or overlap semantics changed", row.Cells)
		}
	}
	// Source ignores plain restrictions and an unknown suffix for membership.
	rt := "Microsoft.Compute/virtualMachines"
	inventory = &ServiceInventory{SubscriptionIDs: []string{auxSubscription}, ResourceTypes: map[string]int64{rt: 1}, SKUsByType: map[string]map[string]int64{rt: {"A": 1}}, ResourceTypesByRegion: map[string]map[string]int64{"eastus": {rt: 1}}}
	c.MissingSKUs = []string{rt + ":A (unknown)"}
	c.RestrictedSKUs = []string{rt + ":A"}
	c.ZoneRestrictedSKUs = []string{rt + ":A (Zones blocked: 2)"}
	c.MissingResourceTypes = []string{strings.ToUpper(rt)}
	tables, err = ProjectServiceAvailability(context.Background(), map[string]string{auxSubscription: "Synthetic"}, inventory, []ServiceComparison{c})
	if err != nil || !slices.Equal(tables[0].Rows[0].Cells, []string{rt, "1", "", "1", "A", "Available", "", "Not available"}) {
		t.Fatal("source lower-key implementation lookup, unknown/plain restriction or case-sensitive zone suffix changed")
	}
}

func TestServiceSelectedScope(t *testing.T) {
	inventory, comparisons, _ := serviceFixture(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	foreign := comparisons[0]
	foreign.SubscriptionID = otherAuxSubscription
	tables, err := ProjectServiceAvailability(context.Background(), scope, inventory, []ServiceComparison{foreign})
	requireServiceFailure(t, tables, err, "region_service_scope_invalid")
	foreign = comparisons[0]
	foreign.SubscriptionName = "private-different"
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, []ServiceComparison{foreign})
	requireServiceFailure(t, tables, err, "region_service_scope_invalid")
	for _, ids := range [][]string{{otherAuxSubscription}, {auxSubscription, auxSubscription}, nil, {"private-not-an-id"}} {
		copy := *inventory
		copy.SubscriptionIDs = ids
		tables, err = ProjectServiceAvailability(context.Background(), scope, &copy, comparisons)
		requireServiceFailure(t, tables, err, "region_service_scope_invalid")
	}
	for _, invalid := range []map[string]string{{"bad": "Synthetic"}, {auxSubscription: "private\x00"}, {auxSubscription: strings.Repeat("s", MaxLabelBytes+1)}} {
		tables, err = ProjectServiceAvailability(context.Background(), invalid, nil, nil)
		requireServiceFailure(t, tables, err, "region_service_scope_invalid")
	}
	alias := "ABCDEFAB-1111-1111-1111-111111111111"
	tables, err = ProjectServiceAvailability(context.Background(), map[string]string{alias: "Synthetic", strings.ToLower(alias): "Synthetic"}, nil, nil)
	requireServiceFailure(t, tables, err, "region_service_scope_invalid")
	// Equal names across distinct contributor UUIDs remain permitted.
	inventory.SubscriptionIDs = append(inventory.SubscriptionIDs, otherAuxSubscription)
	comparisons[1].SubscriptionID = otherAuxSubscription
	scope[otherAuxSubscription] = "Synthetic"
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, comparisons)
	if err != nil || len(tables) != 2 || tables[0].Rows[0].SubscriptionID != "" {
		t.Fatal("aggregate scope fabricated or conflated contributor identity")
	}
}

func TestServiceMalformedAndCollision(t *testing.T) {
	for _, mutate := range []func(*ServiceInventory, *[]ServiceComparison){
		func(i *ServiceInventory, c *[]ServiceComparison) { (*c)[0].TargetRegion = "eastKus" },
		func(i *ServiceInventory, c *[]ServiceComparison) { (*c)[0].SourceRegion = "" },
		func(i *ServiceInventory, c *[]ServiceComparison) { *c = append(*c, (*c)[0]) },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypes["private"] = -1 },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypes["private"] = MaxAuxCount + 1 },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.SKUsByType["private"] = map[string]int64{"A": -1} },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypesByRegion["private\n"] = nil },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypes["private\x00"] = 1 },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypes[string([]byte{0xff})] = 1 },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypes["private?"] = 1 },
		func(i *ServiceInventory, c *[]ServiceComparison) { i.ResourceTypes[strings.Repeat("s", 513)] = 1 },
	} {
		inventory, comparisons, _ := serviceFixture(t)
		mutate(inventory, &comparisons)
		tables, err := ProjectServiceAvailability(context.Background(), map[string]string{auxSubscription: "Synthetic"}, inventory, comparisons)
		requireServiceFailure(t, tables, err, "region_service_input_invalid")
	}
	inventory, comparisons, _ := serviceFixture(t)
	comparisons[0].MissingSKUs = []string{"private\x00"}
	tables, err := ProjectServiceAvailability(context.Background(), map[string]string{auxSubscription: "Synthetic"}, inventory, comparisons)
	requireServiceFailure(t, tables, err, "region_service_text_limit")
	comparisons[0].MissingSKUs = nil
	comparisons[0].TargetRegion = strings.Repeat("a", 21) + "b"
	comparisons[1].TargetRegion = strings.Repeat("a", 21) + "c"
	tables, err = ProjectServiceAvailability(context.Background(), map[string]string{auxSubscription: "Synthetic"}, inventory, comparisons)
	requireServiceFailure(t, tables, err, "region_service_sheet_collision")
	tables, err = ProjectServiceAvailability(context.Background(), map[string]string{auxSubscription: "Synthetic"}, inventory, comparisons[:1])
	if err != nil || tables[0].SheetName != "Svc Avail "+strings.Repeat("a", 21) {
		t.Fatal("noncolliding source sheet truncation changed")
	}
}

func serviceTargets(n int) []ServiceComparison {
	comparisons := make([]ServiceComparison, n)
	for i := range comparisons {
		comparisons[i] = ServiceComparison{SubscriptionID: auxSubscription, SubscriptionName: "Synthetic", SourceRegion: "eastus", TargetRegion: fmt.Sprintf("region%02d", i)}
	}
	return comparisons
}

func TestServiceWorkBounds(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	inventory := &ServiceInventory{SubscriptionIDs: []string{auxSubscription}, ResourceTypes: map[string]int64{}}
	for i := 0; i < 8192; i++ {
		inventory.ResourceTypes[fmt.Sprintf("microsoft.custom/type%05d", i)] = 1
	}
	tables, err := ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(1))
	if err != nil || len(tables[0].Rows) != 8192 {
		t.Fatal("literal8192 output row boundary rejected", err)
	}
	inventory.ResourceTypes["microsoft.custom/overflow"] = 1
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(1))
	requireServiceFailure(t, tables, err, "region_service_input_limit")
	delete(inventory.ResourceTypes, "microsoft.custom/overflow")
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(2))
	requireServiceFailure(t, tables, err, "region_service_input_limit")
	inventory.ResourceTypes = nil
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(32))
	if err != nil || len(tables) != 32 {
		t.Fatal("literal32 target boundary rejected")
	}
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(33))
	requireServiceFailure(t, tables, err, "region_service_input_limit")
	// 1 selected +1 contributor +1 comparison +65533 ignored details =65536.
	c := serviceTargets(1)
	c[0].RestrictedSKUs = make([]string, 65533)
	for i := range c[0].RestrictedSKUs {
		c[0].RestrictedSKUs[i] = "x"
	}
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, c)
	if err != nil || len(tables) != 1 {
		t.Fatal("literal65536 global entry boundary rejected", err)
	}
	c[0].RestrictedSKUs = append(c[0].RestrictedSKUs, "x")
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, c)
	requireServiceFailure(t, tables, err, "region_service_input_limit")
	tooMany := make([]ServiceComparison, 8193)
	tables, err = ProjectServiceAvailability(context.Background(), scope, nil, tooMany)
	requireServiceFailure(t, tables, err, "region_service_input_limit")
}

func serviceSKUMap(n, width int) map[string]int64 {
	values := map[string]int64{}
	for i := 0; i < n; i++ {
		prefix := fmt.Sprintf("%03d-", i)
		values[prefix+strings.Repeat("s", width-len(prefix))] = 1
	}
	return values
}

func TestServiceJoinedAndReplicatedTextBounds(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	rt := "microsoft.custom/widgets"
	inventory := &ServiceInventory{SubscriptionIDs: []string{auxSubscription}, ResourceTypes: map[string]int64{rt: 1}, SKUsByType: map[string]map[string]int64{rt: serviceSKUMap(64, 510)}}
	// 64*510 +63*2 +1 =32767 UTF-16 units. One more byte exceeds the cell.
	delete(inventory.SKUsByType[rt], "000-"+strings.Repeat("s", 506))
	inventory.SKUsByType[rt]["000-"+strings.Repeat("s", 507)] = 1
	tables, err := ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(1))
	if err != nil || len(tables[0].Rows[0].Cells[4]) != 32767 {
		t.Fatal("literal32767 joined cell boundary rejected", err)
	}
	delete(inventory.SKUsByType[rt], "000-"+strings.Repeat("s", 507))
	inventory.SKUsByType[rt]["000-"+strings.Repeat("s", 508)] = 1
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(1))
	requireServiceFailure(t, tables, err, "region_service_text_limit")
	// Unicode UTF-16 units are counted exactly, not as twice every rune.
	inventory.SKUsByType[rt] = map[string]int64{strings.Repeat("??", 128): 1}
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(1))
	if err != nil || tables[0].Rows[0].Cells[4] != strings.Repeat("??", 128) {
		t.Fatal("bounded Unicode cell rejected")
	}
	// Small decoded inventory replicated into32 sheets exceeds16MiB output.
	inventory.ResourceTypes = map[string]int64{}
	inventory.SKUsByType = map[string]map[string]int64{}
	for i := 0; i < 256; i++ {
		typeName := fmt.Sprintf("microsoft.custom/type%03d", i)
		inventory.ResourceTypes[typeName] = 1
		inventory.SKUsByType[typeName] = serviceSKUMap(16, 64)
	}
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(32))
	if err != nil || len(tables) != 32 || len(tables[0].Rows) != 256 {
		t.Fatal("bounded replicated inventory rejected", err)
	}
	for typeName := range inventory.ResourceTypes {
		inventory.SKUsByType[typeName] = serviceSKUMap(32, 64)
	}
	tables, err = ProjectServiceAvailability(context.Background(), scope, inventory, serviceTargets(32))
	requireServiceFailure(t, tables, err, "region_service_text_limit")
}

type serviceCancelContext struct {
	context.Context

	cancel context.CancelFunc

	calls int

	at int
}

func (c *serviceCancelContext) Err() error {
	c.calls++
	if c.calls >= c.at {
		c.cancel()
	}
	return c.Context.Err()
}

func TestServiceOwnershipConcurrencyAndCancellation(t *testing.T) {
	inventory, comparisons, _ := serviceFixture(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	baseline, err := ProjectServiceAvailability(context.Background(), scope, inventory, comparisons)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := ProjectServiceAvailability(context.Background(), scope, inventory, comparisons)
			if err != nil || !reflect.DeepEqual(got, baseline) {
				t.Error("service projection nondeterministic/shared state")
				return
			}
			got[0].Columns[0] = "owned"
			got[0].Rows[0].Cells[0] = "owned"
		}()
	}
	wg.Wait()
	for _, at := range []int{1, 10, 25, 50} {
		ctx, cancel := context.WithCancel(context.Background())
		checking := &serviceCancelContext{Context: ctx, cancel: cancel, at: at}
		got, err := ProjectServiceAvailability(checking, scope, inventory, comparisons)
		cancel()
		if checking.calls >= at {
			requireServiceFailure(t, got, err, "region_service_cancelled")
			if !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation identity lost")
			}
		} else if err != nil {
			t.Fatal(err)
		}
	}
	inventory.ResourceTypes["microsoft.compute/virtualmachines"] = 999
	inventory.SKUsByType["microsoft.compute/virtualmachines"]["changed"] = 1
	comparisons[0].MissingSKUs[0] = "changed"
	scope[auxSubscription] = "changed"
	if baseline[0].Columns[0] != "ResourceType" || baseline[0].Rows[0].Cells[1] != "4" || baseline[1].Rows[0].Cells[4] != "Standard_D4s_v3" {
		t.Fatal("caller input or returned output mutation escaped ownership")
	}
}
