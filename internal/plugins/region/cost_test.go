package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func costFixture(t *testing.T) (*CostSheetData, *auxiliaryCapturedTable) {
	t.Helper()
	var input struct{ Pricing *CostSheetData }
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
	var expected *auxiliaryCapturedTable
	if err := json.Unmarshal(outputs["cost-full"], &expected); err != nil {
		t.Fatal(err)
	}
	if len(outputs) != 14 || input.Pricing == nil || len(input.Pricing.MeterInputs) != 3 || expected == nil || len(expected.Table) != 4 {
		t.Fatal("independent cost capture topology changed")
	}
	for _, branch := range []string{"cost-nil", "cost-no-meters", "cost-no-prices"} {
		if string(outputs[branch]) != "null" {
			t.Fatal("source nil cost branch changed")
		}
	}
	input.Pricing.SubscriptionIDs = []string{auxSubscription}
	return input.Pricing, expected
}

func requireCostFailure(t *testing.T, table *assessment.PluginTable, err error, code string) {
	t.Helper()
	if err == nil || table == nil || len(table.Rows) != 0 || len(table.Columns) != 5 || table.Health.Status != assessment.StageFailed || table.Health.Error == nil || table.Health.Error.Code != code || strings.Contains(err.Error(), "private") || assessment.ValidatePluginTables([]assessment.PluginTable{*table}) != nil {
		t.Fatalf("unsafe cost input, wrong guard or partial output: expected=%s err=%v", code, err)
	}
}

func TestCostCapturedCellsAndEmpty(t *testing.T) {
	input, expected := costFixture(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	table, err := ProjectCostComparison(context.Background(), scope, input)
	if err != nil || table == nil || table.SheetName != expected.SheetName || table.Description != expected.Description || table.Metadata != Metadata() || !slices.Equal(table.Columns, expected.Table[0]) || len(table.Rows) != 3 || table.Health.Status != assessment.StageCompleted || table.Health.Records != 3 || assessment.ValidatePluginTables([]assessment.PluginTable{*table}) != nil {
		t.Fatal("captured cost metadata/shape/health changed", err)
	}
	for i, row := range table.Rows {
		if !slices.Equal(row.Cells, expected.Table[i+1]) || row.SubscriptionID != "" {
			t.Fatal("captured cost cell or aggregate UUID changed", i)
		}
	}
	for _, empty := range []*CostSheetData{nil, {SubscriptionIDs: []string{auxSubscription}}, {SubscriptionIDs: []string{auxSubscription}, MeterInputs: input.MeterInputs}} {
		table, err = ProjectCostComparison(context.Background(), scope, empty)
		if err != nil || table != nil {
			t.Fatal("source absent cost branch changed")
		}
	}
}

func TestCostMetadataReceivedOrder(t *testing.T) {
	input := &CostSheetData{SubscriptionIDs: []string{auxSubscription}, MeterInputs: []CostMeter{{"meter-b", "same", "p", "sku"}, {"meter-a", "same", "p", "sku"}, {"meter-c", "Same", "p", "sku"}}, RegionPricing: map[string]map[string]float64{"unknown": {"eastus": -1}}, PriceItems: []CostPriceItem{{"same", "p", "sku", "First service", "First product"}, {"same", "p", "sku", "Last service", "Last product"}}}
	original := slices.Clone(input.MeterInputs)
	table, err := ProjectCostComparison(context.Background(), map[string]string{auxSubscription: "Synthetic"}, input)
	if err != nil || !slices.Equal(input.MeterInputs, original) || !slices.Equal(table.Columns, []string{"MeterId", "ServiceName", "MeterName", "ProductName", "SKUName", "eastus-RetailPrice"}) {
		t.Fatal("caller order changed or unknown-meter region omitted", err)
	}
	want := [][]string{{"meter-a", "", "same", "", "sku", ""}, {"meter-b", "First service", "same", "First product", "sku", ""}, {"meter-c", "", "Same", "", "sku", ""}}
	for i, row := range table.Rows {
		if !slices.Equal(row.Cells, want[i]) {
			t.Fatal("first received meter/item or exact case match changed", i, row.Cells)
		}
	}
}

func TestCostPriceSemantics(t *testing.T) {
	input := &CostSheetData{SubscriptionIDs: []string{auxSubscription}, MeterInputs: []CostMeter{{MeterID: "meter"}}, RegionPricing: map[string]map[string]float64{"meter": {"a": 0, "b": -1e12, "c": 0.00000001, "d": 1.23456, "e": 1e12}, "unknown": {"f": 1}}}
	table, err := ProjectCostComparison(context.Background(), map[string]string{auxSubscription: "Synthetic"}, input)
	if err != nil || !slices.Equal(table.Rows[0].Cells, []string{"meter", "", "", "", "", "", "", "0.0000", "1.2346", "1000000000000.0000", ""}) {
		t.Fatal("zero/negative/missing/tiny/positive price semantics changed", err)
	}
}

func TestCostSelectedScopeAndMalformed(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	for _, ids := range [][]string{nil, {otherAuxSubscription}, {auxSubscription, auxSubscription}, {"private-not-an-id"}} {
		input, _ := costFixture(t)
		input.SubscriptionIDs = ids
		table, err := ProjectCostComparison(context.Background(), scope, input)
		requireCostFailure(t, table, err, "region_cost_scope_invalid")
	}
	for _, selection := range []map[string]string{{"bad": "Synthetic"}, {auxSubscription: ""}, {auxSubscription: "private\x00"}, {auxSubscription: strings.Repeat("s", 513)}} {
		table, err := ProjectCostComparison(context.Background(), selection, nil)
		requireCostFailure(t, table, err, "region_cost_scope_invalid")
	}
	alias := "ABCDEFAB-1111-1111-1111-111111111111"
	table, err := ProjectCostComparison(context.Background(), map[string]string{alias: "Synthetic", strings.ToLower(alias): "Synthetic"}, nil)
	requireCostFailure(t, table, err, "region_cost_scope_invalid")
	for _, price := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1e12 + 1, -1e12 - 1} {
		input, _ := costFixture(t)
		input.RegionPricing["ignored"] = map[string]float64{"eastus": price}
		table, err = ProjectCostComparison(context.Background(), scope, input)
		requireCostFailure(t, table, err, "region_cost_input_invalid")
	}
	input, _ := costFixture(t)
	input.MeterInputs = append(input.MeterInputs, input.MeterInputs[0])
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_input_invalid")
	input, _ = costFixture(t)
	input.MeterInputs[0].MeterID = ""
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_input_invalid")
	input, _ = costFixture(t)
	input.RegionPricing["unused"] = map[string]float64{"east\u212aus": 1}
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_input_invalid")
	for _, label := range []string{"private\x00", "private\ufffd", string([]byte{0xff}), strings.Repeat("s", 513)} {
		input, _ = costFixture(t)
		input.PriceItems[0].ServiceName = label
		table, err = ProjectCostComparison(context.Background(), scope, input)
		requireCostFailure(t, table, err, "region_cost_text_limit")
	}
	input, _ = costFixture(t)
	input.SubscriptionIDs = append(input.SubscriptionIDs, otherAuxSubscription)
	table, err = ProjectCostComparison(context.Background(), map[string]string{auxSubscription: "Synthetic", otherAuxSubscription: "Synthetic"}, input)
	if err != nil || table.Rows[0].SubscriptionID != "" {
		t.Fatal("same-name aggregate declaration conflated or invented identity")
	}
}

func TestCostWorkLimits(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	input := &CostSheetData{SubscriptionIDs: []string{auxSubscription}, RegionPricing: map[string]map[string]float64{"unknown": {"eastus": 1}}}
	for i := 0; i < 8192; i++ {
		input.MeterInputs = append(input.MeterInputs, CostMeter{MeterID: fmt.Sprintf("meter%05d", i)})
	}
	table, err := ProjectCostComparison(context.Background(), scope, input)
	if err != nil || len(table.Rows) != 8192 {
		t.Fatal("literal8192-meter boundary rejected", err)
	}
	input.MeterInputs = append(input.MeterInputs, CostMeter{MeterID: "overflow"})
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_input_limit")
	input.MeterInputs = input.MeterInputs[:1]
	input.RegionPricing["unknown"] = map[string]float64{}
	for i := 0; i < 32; i++ {
		input.RegionPricing["unknown"][fmt.Sprintf("r%02d", i)] = 1
	}
	table, err = ProjectCostComparison(context.Background(), scope, input)
	if err != nil || len(table.Columns) != 37 {
		t.Fatal("literal32-region boundary rejected")
	}
	input.RegionPricing["unknown"]["overflow"] = 1
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_input_limit")
	input.RegionPricing = map[string]map[string]float64{"unknown": {"eastus": 1}}
	// selected1+contributor1+meter1+outer1+price1+items65531 =65536.
	input.PriceItems = make([]CostPriceItem, 65531)
	table, err = ProjectCostComparison(context.Background(), scope, input)
	if err != nil || len(table.Rows) != 1 {
		t.Fatal("literal65536-entry boundary rejected", err)
	}
	input.PriceItems = append(input.PriceItems, CostPriceItem{})
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_input_limit")
}

func TestCostTextLimitsAndUnicode(t *testing.T) {
	scope := map[string]string{auxSubscription: "Synthetic"}
	label := strings.Repeat("s", 512)
	input := &CostSheetData{PriceItems: make([]CostPriceItem, 6600)}
	for i := range input.PriceItems {
		input.PriceItems[i] = CostPriceItem{label, label, label, label, label}
	}
	table, err := ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_text_limit")
	// Independent exact decoded16MiB-minus64KiB boundary, including scope.
	// 6528 items*5*512 +scope9 =16711689; remove9 bytes for16711680.
	input.PriceItems = input.PriceItems[:6528]
	input.PriceItems[0].ServiceName = strings.Repeat("s", 503)
	table, err = ProjectCostComparison(context.Background(), scope, input)
	if err != nil || table != nil {
		t.Fatal("literal decoded text budget rejected", err)
	}
	input.PriceItems[0].ServiceName += "s"
	table, err = ProjectCostComparison(context.Background(), scope, input)
	requireCostFailure(t, table, err, "region_cost_text_limit")
	input, _ = costFixture(t)
	input.MeterInputs[0].MeterName = strings.Repeat("\U0001f600", 128)
	input.MeterInputs[0].SKUName = strings.Repeat("\u00e9", 256)
	table, err = ProjectCostComparison(context.Background(), scope, input)
	if err != nil || table.Rows[1].Cells[2] != strings.Repeat("\U0001f600", 128) || table.Rows[1].Cells[4] != strings.Repeat("\u00e9", 256) {
		t.Fatal("safe512-byte supplementary/BMP cost labels changed", err)
	}
}

func TestCostOwnershipConcurrencyAndCancellation(t *testing.T) {
	input, _ := costFixture(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	baseline, err := ProjectCostComparison(context.Background(), scope, input)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := ProjectCostComparison(context.Background(), scope, input)
			if err != nil || !reflect.DeepEqual(got, baseline) {
				t.Error("cost projection leaked state or order")
				return
			}
			got.Columns[0] = "owned"
			got.Rows[0].Cells[0] = "owned"
		}()
	}
	wg.Wait()
	for _, at := range []int{1, 10, 20} {
		ctx, cancel := context.WithCancel(context.Background())
		checking := &serviceCancelContext{Context: ctx, cancel: cancel, at: at}
		got, err := ProjectCostComparison(checking, scope, input)
		cancel()
		requireCostFailure(t, got, err, "region_cost_cancelled")
		if !errors.Is(err, context.Canceled) {
			t.Fatal("cost cancellation identity lost")
		}
	}
	input.MeterInputs[0].MeterID = "changed"
	input.RegionPricing["meter-a"]["eastus"] = 999
	input.PriceItems[0].ServiceName = "changed"
	input.SubscriptionIDs[0] = "changed"
	scope[auxSubscription] = "changed"
	if baseline.Columns[0] != "MeterId" || baseline.Rows[0].Cells[1] != "Virtual Machines" || baseline.Rows[0].Cells[5] != "1.2346" {
		t.Fatal("cost input/output mutation escaped ownership")
	}
}
