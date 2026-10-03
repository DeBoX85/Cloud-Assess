// Copyright (c) Microsoft Corporation.
// Scoring and primary table behavior derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_SELECTION.md for provenance and corrections.
package region

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const (
	Name             = "region-selection"
	MaxComparisons   = 8192
	MaxSubscriptions = 1000
	MaxLabelBytes    = 512
	MaxDetails       = 4096
	MaxCount         = 65536
)

var subscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
var regionID = regexp.MustCompile(`^[a-z][a-z0-9]{0,63}$`)
var columns = []string{"Subscription", "Source Region", "Target Region", "Source Resource Type Count", "Available Resource Types", "Unavailable Resource Types", "Availability %", "Total SKUs Checked", "Available SKUs", "Unavailable SKUs", "Restricted SKUs", "Zone-Restricted SKUs", "Unknown SKUs", "SKU Availability %", "Availability Zones", "Target AZ Mapping", "Avg Latency (ms)", "Avg Cost Difference %", "Recommendation Score", "Score Quality", "Recommendation", "Missing Resource Types", "Unavailable SKUs (detail)", "Restricted SKUs (detail)", "Zone-Restricted SKUs (detail)"}

func Metadata() assessment.PluginMetadata {
	return assessment.PluginMetadata{Name: Name, Version: "0.2.0-beta", Description: "Analyzes optimal Azure region selection based on service availability, network latency, and cost comparison", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}
}

func PendingTable() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "comparisons", Metadata: Metadata(), SheetName: "Region Selection", Description: "Analysis of optimal Azure region selection based on service availability, network latency, and cost factors", Columns: slices.Clone(columns), Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped, Warnings: []assessment.AssessmentWarning{{Code: "plugin_not_run", Message: "requested plugin has not executed"}}}}
}

// Comparison contains selected, decoded input. Project does not collect or infer
// availability, prices or latency. Unknown SKU API checks remain explicit.
type Comparison struct {
	SubscriptionID, SubscriptionName, SourceRegion, TargetRegion          string
	SourceResourceTypeCount, AvailableTypes, UnavailableTypes             int
	AvailabilityPercent, AvgCostDifference, AvgLatencyMs                  float64
	HasCostData, LatencyEstimated                                         bool
	MissingResourceTypes, MissingSKUs, RestrictedSKUs, ZoneRestrictedSKUs []string
	TotalSKUsChecked, AvailableSKUs, UnavailableSKUs, UnknownSKUs         int
	SKUAvailabilityPercent                                                float64
	SourceZoneCount, TargetZoneCount                                      int
	TargetZoneMappings                                                    map[string]string
}

// Project owns its output and bounds work before allocating joined detail cells.
// Inputs must remain stable during the call; no caller-owned maps/slices escape.
func Project(ctx context.Context, subscriptions map[string]string, input []Comparison) (assessment.PluginTable, error) {
	table := PendingTable()
	fail := func(code string) (assessment.PluginTable, error) {
		t := PendingTable()
		t.Health = assessment.StageExecution{Name: Name, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "region comparison input could not be safely projected"}}
		return t, fmt.Errorf("region comparison input could not be safely projected")
	}
	if err := ctx.Err(); err != nil {
		t, _ := fail("region_cancelled")
		return t, err
	}
	if len(input) > MaxComparisons || len(subscriptions) > MaxSubscriptions {
		return fail("region_input_limit")
	}
	textBytes := 0
	text := func(s string) bool {
		textBytes += len(s)
		return len(s) <= MaxLabelBytes && textBytes <= assessment.MaxPluginTextBytes && utf8.ValidString(s) && strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || r == 0xfffe || r == 0xffff }) < 0
	}
	scope := make(map[string]string, len(subscriptions))
	for id, name := range subscriptions {
		if len(id) != 36 {
			return fail("region_scope_invalid")
		}
		key := strings.ToLower(id)
		if !subscriptionID.MatchString(id) || !text(name) {
			return fail("region_scope_invalid")
		}
		if _, exists := scope[key]; exists {
			return fail("region_scope_invalid")
		}
		scope[key] = name
	}
	seen := make(map[string]bool, len(input))
	for _, c := range input {
		if err := ctx.Err(); err != nil {
			t, _ := fail("region_cancelled")
			return t, err
		}
		if len(c.SubscriptionID) != 36 || len(c.SourceRegion) > 64 || len(c.TargetRegion) > 64 {
			return fail("region_input_invalid")
		}
		id := strings.ToLower(c.SubscriptionID)
		name, selected := scope[id]
		if !selected || name != c.SubscriptionName {
			return fail("region_scope_invalid")
		}
		if !regionID.MatchString(c.SourceRegion) || !regionID.MatchString(c.TargetRegion) || !text(c.SubscriptionName) {
			return fail("region_input_invalid")
		}
		key := id + "/" + c.SourceRegion + "/" + c.TargetRegion
		if seen[key] {
			return fail("region_input_invalid")
		}
		seen[key] = true
		for _, count := range []int{c.SourceResourceTypeCount, c.AvailableTypes, c.UnavailableTypes, c.TotalSKUsChecked, c.AvailableSKUs, c.UnavailableSKUs, c.UnknownSKUs, c.SourceZoneCount, c.TargetZoneCount} {
			if count < 0 || count > MaxCount {
				return fail("region_input_invalid")
			}
		}
		if c.AvailableTypes+c.UnavailableTypes != c.SourceResourceTypeCount || c.AvailableSKUs+c.UnavailableSKUs+c.UnknownSKUs+len(c.RestrictedSKUs)+len(c.ZoneRestrictedSKUs) != c.TotalSKUsChecked {
			return fail("region_input_invalid")
		}
		for _, n := range []float64{c.AvailabilityPercent, c.SKUAvailabilityPercent, c.AvgCostDifference, c.AvgLatencyMs} {
			if math.IsNaN(n) || math.IsInf(n, 0) {
				return fail("region_input_invalid")
			}
		}
		if c.AvgLatencyMs < 0 || c.AvgCostDifference < -100 || c.AvailabilityPercent < 0 || c.AvailabilityPercent > 100 || c.SKUAvailabilityPercent < 0 || c.SKUAvailabilityPercent > 100 {
			return fail("region_input_invalid")
		}
		availability := 0.0
		if c.SourceResourceTypeCount > 0 {
			availability = float64(c.AvailableTypes) / float64(c.SourceResourceTypeCount) * 100
		}
		skuAvailability := 0.0
		if c.TotalSKUsChecked > 0 {
			skuAvailability = 100
			if confirmed := c.TotalSKUsChecked - c.UnknownSKUs; confirmed > 0 {
				skuAvailability = float64(c.AvailableSKUs) / float64(confirmed) * 100
			}
		}
		if math.Abs(c.AvailabilityPercent-availability) > 1e-9 || math.Abs(c.SKUAvailabilityPercent-skuAvailability) > 1e-9 {
			return fail("region_input_invalid")
		}
		for _, details := range [][]string{c.MissingResourceTypes, c.MissingSKUs, c.RestrictedSKUs, c.ZoneRestrictedSKUs} {
			if len(details) > MaxDetails {
				return fail("region_input_limit")
			}
			units := max(0, len(details)-1) * 2
			for _, detail := range details {
				if !text(detail) {
					return fail("region_text_limit")
				}
				units += utf16Units(detail)
			}
			if units > assessment.MaxPluginCellUnits {
				return fail("region_text_limit")
			}
		}
		if len(c.TargetZoneMappings) > MaxDetails {
			return fail("region_input_limit")
		}
		units := max(0, len(c.TargetZoneMappings)-1) * 2
		for logical, physical := range c.TargetZoneMappings {
			if !text(logical) || !text(physical) || logical == "" || physical == "" {
				return fail("region_text_limit")
			}
			units += utf16Units(logical) + utf16Units(physical) + 1
		}
		if units > assessment.MaxPluginCellUnits {
			return fail("region_text_limit")
		}
	}
	type scoredRow struct {
		score float64
		row   assessment.PluginRow
	}
	rows := make([]scoredRow, 0, len(input))
	for _, c := range input {
		if err := ctx.Err(); err != nil {
			t, _ := fail("region_cancelled")
			return t, err
		}
		score := calculateScore(c)
		rows = append(rows, scoredRow{score, assessment.PluginRow{SubscriptionID: strings.ToLower(c.SubscriptionID), Cells: projectCells(c, score)}})
	}
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.score != b.score {
			return a.score > b.score
		}
		if a.row.Cells[0] != b.row.Cells[0] {
			return a.row.Cells[0] < b.row.Cells[0]
		}
		if a.row.SubscriptionID != b.row.SubscriptionID {
			return a.row.SubscriptionID < b.row.SubscriptionID
		}
		return slices.Compare(a.row.Cells[1:3], b.row.Cells[1:3]) < 0
	})
	for _, row := range rows {
		table.Rows = append(table.Rows, row.row)
	}
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(table.Rows)}
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		return fail("region_output_limit")
	}
	if err := ctx.Err(); err != nil {
		t, _ := fail("region_cancelled")
		return t, err
	}
	return table, nil
}

func calculateScore(c Comparison) float64 {
	sku := 100.0
	if confirmed := c.TotalSKUsChecked - c.UnknownSKUs; confirmed > 0 {
		sku = (float64(c.AvailableSKUs) + float64(len(c.ZoneRestrictedSKUs))*0.75 + float64(len(c.RestrictedSKUs))*0.5) / float64(confirmed) * 100
		sku = math.Min(sku, 100)
	}
	cost := 100.0
	if c.HasCostData {
		cost = math.Max(0, math.Min(100, 100-c.AvgCostDifference*2))
	}
	latency := 100.0
	if c.AvgLatencyMs > 0 {
		switch {
		case c.AvgLatencyMs < 50:
			latency = 100
		case c.AvgLatencyMs > 200:
			latency = 0
		default:
			latency = 100 - ((c.AvgLatencyMs - 50) / 150 * 100)
		}
	}
	score := c.AvailabilityPercent*0.35 + sku*0.30 + cost*0.15 + latency*0.20
	if c.SourceZoneCount > 0 && c.TargetZoneCount < c.SourceZoneCount {
		score *= 1 - float64(c.SourceZoneCount-c.TargetZoneCount)/float64(c.SourceZoneCount)*0.10
	}
	return score
}

func projectCells(c Comparison, score float64) []string {
	cost, latency, sku := "N/A", "N/A", "N/A"
	if c.HasCostData {
		cost = fmt.Sprintf("%+.2f%%", c.AvgCostDifference)
	}
	if c.AvgLatencyMs > 0 {
		latency = fmt.Sprintf("%.1f", c.AvgLatencyMs)
	}
	if c.TotalSKUsChecked > 0 {
		sku = fmt.Sprintf("%.2f%%", c.SKUAvailabilityPercent)
	}
	zones := fmt.Sprintf("%d → %d", c.SourceZoneCount, c.TargetZoneCount)
	switch {
	case c.SourceZoneCount == 0 && c.TargetZoneCount > 0:
		zones += " ✓"
	case c.SourceZoneCount > 0 && c.TargetZoneCount == 0:
		zones += " ✗"
	case c.SourceZoneCount > c.TargetZoneCount:
		zones += " ⚠"
	}
	keys := make([]string, 0, len(c.TargetZoneMappings))
	for key := range c.TargetZoneMappings {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"→"+c.TargetZoneMappings[key])
	}
	quality := []string{}
	if cost == "N/A" {
		quality = append(quality, "no cost data")
	}
	if latency == "N/A" {
		quality = append(quality, "no latency data")
	} else if c.LatencyEstimated {
		quality = append(quality, "estimated latency")
	}
	qualityText := "Full"
	if len(quality) > 0 {
		qualityText = strings.Join(quality, ", ")
	}
	recommendation := "Not Recommended"
	if score >= 80 {
		recommendation = "Recommended"
	} else if score >= 60 {
		recommendation = "Neutral"
	}
	return []string{c.SubscriptionName, c.SourceRegion, c.TargetRegion, strconv.Itoa(c.SourceResourceTypeCount), strconv.Itoa(c.AvailableTypes), strconv.Itoa(c.UnavailableTypes), fmt.Sprintf("%.2f%%", c.AvailabilityPercent), strconv.Itoa(c.TotalSKUsChecked), strconv.Itoa(c.AvailableSKUs), strconv.Itoa(c.UnavailableSKUs), strconv.Itoa(len(c.RestrictedSKUs)), strconv.Itoa(len(c.ZoneRestrictedSKUs)), strconv.Itoa(c.UnknownSKUs), sku, zones, strings.Join(parts, ", "), latency, cost, fmt.Sprintf("%.2f", score), qualityText, recommendation, strings.Join(c.MissingResourceTypes, "; "), strings.Join(c.MissingSKUs, "; "), strings.Join(c.RestrictedSKUs, "; "), strings.Join(c.ZoneRestrictedSKUs, "; ")}
}

func utf16Units(s string) int {
	units := 0
	for _, r := range s {
		units += utf16.RuneLen(r)
	}
	return units
}
