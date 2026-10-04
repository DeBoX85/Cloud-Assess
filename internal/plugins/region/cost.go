// Copyright (c) Microsoft Corporation.
// CostComparison behavior derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_COST_COMPARISON.md for provenance.
package region

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxCostMeters = 8192
const MaxCostRegions = 32
const MaxCostEntries = 65536

// CostMeter contains only fields consumed by the source sheet helper.
type CostMeter struct {
	MeterID, MeterName, ProductID, SKUName string
}

type CostPriceItem struct {
	MeterName, ProductID, SKUName, ServiceName, ProductName string
}

// CostSheetData is already decoded per-run data. SubscriptionIDs declares
// aggregate contributors; the later collector must establish actual ownership.
type CostSheetData struct {
	SubscriptionIDs []string

	MeterInputs []CostMeter

	RegionPricing map[string]map[string]float64

	PriceItems []CostPriceItem
}

func pendingCostSheet() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "cost-comparison", Metadata: Metadata(), SheetName: "CostComparison", Description: "Retail price comparison for resource meters across Azure regions", Columns: []string{"MeterId", "ServiceName", "MeterName", "ProductName", "SKUName"}, Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped}}
}

func costFailure(code string) (*assessment.PluginTable, error) {
	table := pendingCostSheet()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "region cost comparison input could not be safely projected"}}
	return &table, fmt.Errorf("region cost comparison input could not be safely projected")
}

func costCancelled(ctx context.Context) (*assessment.PluginTable, error) {
	table, _ := costFailure("region_cost_cancelled")
	return table, ctx.Err()
}

func costPriceCell(price float64) string {
	if price > 0 {
		return fmt.Sprintf("%.4f", price)
	}
	return ""
}

// ProjectCostComparison owns output without pricing I/O or inference. Inputs
// must remain stable during projection. All guards precede report-row allocation.
func ProjectCostComparison(ctx context.Context, subscriptions map[string]string, input *CostSheetData) (*assessment.PluginTable, error) {
	if ctx.Err() != nil {
		return costCancelled(ctx)
	}
	decodedBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &decodedBytes)
	if ctx.Err() != nil {
		return costCancelled(ctx)
	}
	if !ok {
		return costFailure("region_cost_scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return costFailure("region_cost_scope_invalid")
		}
	}
	if input == nil {
		return nil, nil
	}
	if len(input.SubscriptionIDs) > MaxSubscriptions || len(input.MeterInputs) > MaxCostMeters {
		return costFailure("region_cost_input_limit")
	}
	entries := len(scope)
	entry := func(n int) bool {
		if n > MaxCostEntries-entries {
			return false
		}
		entries += n
		return true
	}
	if !entry(len(input.SubscriptionIDs)) || !entry(len(input.MeterInputs)) || !entry(len(input.RegionPricing)) || !entry(len(input.PriceItems)) {
		return costFailure("region_cost_input_limit")
	}
	contributors := map[string]bool{}
	for _, id := range input.SubscriptionIDs {
		if ctx.Err() != nil {
			return costCancelled(ctx)
		}
		key := strings.ToLower(id)
		_, selected := scope[key]
		if len(id) != 36 || !subscriptionID.MatchString(id) || !selected || contributors[key] {
			return costFailure("region_cost_scope_invalid")
		}
		contributors[key] = true
	}
	if len(input.MeterInputs) > 0 && len(contributors) == 0 {
		return costFailure("region_cost_scope_invalid")
	}
	seenMeters := map[string]bool{}
	firstMeter := map[[3]string]string{}
	for _, meter := range input.MeterInputs {
		if ctx.Err() != nil {
			return costCancelled(ctx)
		}
		if meter.MeterID == "" || seenMeters[meter.MeterID] {
			return costFailure("region_cost_input_invalid")
		}
		for _, s := range []string{meter.MeterID, meter.MeterName, meter.ProductID, meter.SKUName} {
			if !auxiliaryText(s, &decodedBytes) {
				return costFailure("region_cost_text_limit")
			}
		}
		seenMeters[meter.MeterID] = true
		key := [3]string{meter.MeterName, meter.ProductID, meter.SKUName}
		if _, exists := firstMeter[key]; !exists {
			firstMeter[key] = meter.MeterID
		}
	}
	regions := map[string]bool{}
	for meterID, pricing := range input.RegionPricing {
		if ctx.Err() != nil {
			return costCancelled(ctx)
		}
		if !entry(len(pricing)) {
			return costFailure("region_cost_input_limit")
		}
		if meterID == "" || !auxiliaryText(meterID, &decodedBytes) {
			return costFailure("region_cost_text_limit")
		}
		for region, price := range pricing {
			if ctx.Err() != nil {
				return costCancelled(ctx)
			}
			if len(region) > 64 || !regionID.MatchString(region) || !auxiliaryText(region, &decodedBytes) || math.IsNaN(price) || math.IsInf(price, 0) || math.Abs(price) > float64(MaxAuxCount) {
				return costFailure("region_cost_input_invalid")
			}
			regions[region] = true
			if len(regions) > MaxCostRegions {
				return costFailure("region_cost_input_limit")
			}
		}
	}
	metadata := map[string][2]string{}
	for _, item := range input.PriceItems {
		if ctx.Err() != nil {
			return costCancelled(ctx)
		}
		for _, s := range []string{item.MeterName, item.ProductID, item.SKUName, item.ServiceName, item.ProductName} {
			if !auxiliaryText(s, &decodedBytes) {
				return costFailure("region_cost_text_limit")
			}
		}
		if id, match := firstMeter[[3]string{item.MeterName, item.ProductID, item.SKUName}]; match {
			if _, exists := metadata[id]; !exists {
				metadata[id] = [2]string{item.ServiceName, item.ProductName}
			}
		}
	}
	if ctx.Err() != nil {
		return costCancelled(ctx)
	}
	if len(input.MeterInputs) == 0 || len(regions) == 0 {
		return nil, nil
	}
	sortedRegions := make([]string, 0, len(regions))
	for region := range regions {
		sortedRegions = append(sortedRegions, region)
	}
	sort.Strings(sortedRegions)
	sortedMeters := make([]CostMeter, len(input.MeterInputs))
	copy(sortedMeters, input.MeterInputs)
	sort.Slice(sortedMeters, func(i, j int) bool { return sortedMeters[i].MeterID < sortedMeters[j].MeterID })
	projectedBytes := 0
	charge := func(s string) bool {
		projectedBytes += len(s)
		return projectedBytes <= auxTextBudget
	}
	// Charge exact displayed text across every row, including repeated metadata.
	// Labels are already safe512-byte cells; numeric cells have bounded length.
	for _, meter := range sortedMeters {
		if ctx.Err() != nil {
			return costCancelled(ctx)
		}
		meta := metadata[meter.MeterID]
		for _, s := range []string{meter.MeterID, meta[0], meter.MeterName, meta[1], meter.SKUName} {
			if !charge(s) {
				return costFailure("region_cost_text_limit")
			}
		}
		for _, region := range sortedRegions {
			if !charge(costPriceCell(input.RegionPricing[meter.MeterID][region])) {
				return costFailure("region_cost_text_limit")
			}
		}
	}
	table := pendingCostSheet()
	for _, region := range sortedRegions {
		table.Columns = append(table.Columns, region+"-RetailPrice")
	}
	table.Rows = make([]assessment.PluginRow, 0, len(sortedMeters))
	for _, meter := range sortedMeters {
		if ctx.Err() != nil {
			return costCancelled(ctx)
		}
		meta := metadata[meter.MeterID]
		cells := []string{meter.MeterID, meta[0], meter.MeterName, meta[1], meter.SKUName}
		for _, region := range sortedRegions {
			cells = append(cells, costPriceCell(input.RegionPricing[meter.MeterID][region]))
		}
		table.Rows = append(table.Rows, assessment.PluginRow{Cells: cells})
	}
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(table.Rows)}
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		return costFailure("region_cost_output_limit")
	}
	if ctx.Err() != nil {
		return costCancelled(ctx)
	}
	return &table, nil
}
