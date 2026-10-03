// Copyright (c) Microsoft Corporation.
// Quota/reservation formatting derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_AUXILIARY.md for provenance and corrections.
package region

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxAuxRows = 8192
const MaxAuxCount int64 = 1000000000000
const auxTextBudget = assessment.MaxPluginTextBytes - (64 << 10)

// QuotaRow retains decoded source values and flags. Arithmetic is not inferred.
type QuotaRow struct {
	SubscriptionID, Subscription, Region, QuotaType, ResourceName string
	Current, Limit, Available                                    int64
	HeadroomPct                                                  float64
	IsNearLimit, IsOverLimit                                      bool
}

// ReservationRow carries the ten source cells with out-of-band scope identity.
type ReservationRow struct {
	SubscriptionID string
	Cells          []string
}

func pendingAuxiliary(quota bool) assessment.PluginTable {
	table := assessment.PluginTable{
		SchemaVersion: assessment.PluginTableSchemaVersion,
		Metadata:      Metadata(),
		Rows:          []assessment.PluginRow{},
		Health:        assessment.StageExecution{Name: Name, Status: assessment.StageSkipped},
	}
	if quota {
		table.ID = "quota"
		table.SheetName = "Quota"
		table.Description = "VM, network, SQL, App Service, storage, and ARM quota usage per subscription and region"
		table.Columns = []string{"Subscription", "Region", "Quota Type", "Resource", "Current", "Limit", "Available", "Headroom %", "Status"}
	} else {
		table.ID = "reservations"
		table.SheetName = "Capacity Reservations"
		table.Description = "Capacity Reservation Group inventory for migration planning (idle/at-capacity/over-allocated flags)"
		table.Columns = []string{"Subscription", "Region", "Resource Group", "CRG Name", "Reservation Name", "SKU", "Reserved", "Allocated", "Available", "Status"}
	}
	return table
}

func auxiliaryFailure(quota bool, code string) (*assessment.PluginTable, error) {
	table := pendingAuxiliary(quota)
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "region auxiliary input could not be safely projected"}}
	return &table, fmt.Errorf("region auxiliary input could not be safely projected")
}

func auxiliaryText(s string, total *int) bool {
	if len(s) > MaxLabelBytes || !utf8.ValidString(s) || strings.IndexFunc(s, func(r rune) bool {
		return unicode.IsControl(r) || r == 0xfffd || r == 0xfffe || r == 0xffff
	}) >= 0 {
		return false
	}
	*total += len(s)
	return *total <= auxTextBudget
}

// auxiliaryScope copies selection and accounts text before output allocation.
func auxiliaryScope(ctx context.Context, subscriptions map[string]string, total *int) (map[string]string, bool) {
	if len(subscriptions) > MaxSubscriptions {
		return nil, false
	}
	scope := make(map[string]string, len(subscriptions))
	for id, name := range subscriptions {
		if ctx.Err() != nil || len(id) != 36 || !subscriptionID.MatchString(id) || !auxiliaryText(name, total) {
			return nil, false
		}
		key := strings.ToLower(id)
		if _, exists := scope[key]; exists {
			return nil, false
		}
		scope[key] = name
	}
	return scope, true
}

func auxiliarySelected(scope map[string]string, id, name string) bool {
	if len(id) != 36 || !subscriptionID.MatchString(id) {
		return false
	}
	selectedName, selected := scope[strings.ToLower(id)]
	return selected && selectedName == name
}

func auxiliaryCancelled(ctx context.Context, quota bool) (*assessment.PluginTable, error) {
	table, _ := auxiliaryFailure(quota, "region_aux_cancelled")
	return table, ctx.Err()
}

func completeAuxiliary(ctx context.Context, table assessment.PluginTable, quota bool) (*assessment.PluginTable, error) {
	if ctx.Err() != nil {
		return auxiliaryCancelled(ctx, quota)
	}
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(table.Rows)}
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		return auxiliaryFailure(quota, "region_aux_output_limit")
	}
	if ctx.Err() != nil {
		return auxiliaryCancelled(ctx, quota)
	}
	return &table, nil
}

// ProjectQuota preserves source received order and absent-empty behavior.
// All input guards complete before allocating rows. Inputs must remain stable.
func ProjectQuota(ctx context.Context, subscriptions map[string]string, input []QuotaRow) (*assessment.PluginTable, error) {
	if ctx.Err() != nil {
		return auxiliaryCancelled(ctx, true)
	}
	if len(input) > MaxAuxRows {
		return auxiliaryFailure(true, "region_aux_input_limit")
	}
	textBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &textBytes)
	if ctx.Err() != nil {
		return auxiliaryCancelled(ctx, true)
	}
	if !ok {
		return auxiliaryFailure(true, "region_aux_scope_invalid")
	}
	seen := make(map[[4]string]bool, len(input))
	for _, row := range input {
		if ctx.Err() != nil {
			return auxiliaryCancelled(ctx, true)
		}
		if !auxiliarySelected(scope, row.SubscriptionID, row.Subscription) {
			return auxiliaryFailure(true, "region_aux_scope_invalid")
		}
		if len(row.Region) > 64 || !regionID.MatchString(row.Region) {
			return auxiliaryFailure(true, "region_aux_input_invalid")
		}
		for _, label := range []string{row.Subscription, row.Region, row.QuotaType, row.ResourceName} {
			if label == "" || !auxiliaryText(label, &textBytes) {
				return auxiliaryFailure(true, "region_aux_text_limit")
			}
		}
		if row.Current < 0 || row.Current > MaxAuxCount || row.Limit < 0 || row.Limit > MaxAuxCount || row.Available < -MaxAuxCount || row.Available > MaxAuxCount || math.IsNaN(row.HeadroomPct) || math.IsInf(row.HeadroomPct, 0) || math.Abs(row.HeadroomPct) > float64(MaxAuxCount) {
			return auxiliaryFailure(true, "region_aux_input_invalid")
		}
		key := [4]string{strings.ToLower(row.SubscriptionID), row.Region, row.QuotaType, row.ResourceName}
		if seen[key] {
			return auxiliaryFailure(true, "region_aux_input_invalid")
		}
		seen[key] = true
		// Worst-case numeric/status output is bounded independently of labels.
		textBytes += 128
		if textBytes > auxTextBudget {
			return auxiliaryFailure(true, "region_aux_text_limit")
		}
	}
	if len(input) == 0 {
		return nil, nil
	}
	table := pendingAuxiliary(true)
	table.Rows = make([]assessment.PluginRow, 0, len(input))
	for _, row := range input {
		if ctx.Err() != nil {
			return auxiliaryCancelled(ctx, true)
		}
		status := "OK"
		if row.IsOverLimit {
			status = "At/Over Limit"
		} else if row.IsNearLimit {
			status = "Near Limit"
		}
		table.Rows = append(table.Rows, assessment.PluginRow{SubscriptionID: strings.ToLower(row.SubscriptionID), Cells: []string{row.Subscription, row.Region, row.QuotaType, row.ResourceName, strconv.FormatInt(row.Current, 10), strconv.FormatInt(row.Limit, 10), strconv.FormatInt(row.Available, 10), fmt.Sprintf("%.1f%%", row.HeadroomPct), status}})
	}
	return completeAuxiliary(ctx, table, true)
}

// ProjectReservations preserves ten source cells and absent-empty behavior.
// Count/status validation bounds decoded input without recalculating arithmetic.
func ProjectReservations(ctx context.Context, subscriptions map[string]string, input []ReservationRow) (*assessment.PluginTable, error) {
	if ctx.Err() != nil {
		return auxiliaryCancelled(ctx, false)
	}
	if len(input) > MaxAuxRows {
		return auxiliaryFailure(false, "region_aux_input_limit")
	}
	textBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &textBytes)
	if ctx.Err() != nil {
		return auxiliaryCancelled(ctx, false)
	}
	if !ok {
		return auxiliaryFailure(false, "region_aux_scope_invalid")
	}
	seen := make(map[[5]string]bool, len(input))
	for _, row := range input {
		if ctx.Err() != nil {
			return auxiliaryCancelled(ctx, false)
		}
		if len(row.Cells) != 10 {
			return auxiliaryFailure(false, "region_aux_input_invalid")
		}
		if !auxiliarySelected(scope, row.SubscriptionID, row.Cells[0]) {
			return auxiliaryFailure(false, "region_aux_scope_invalid")
		}
		if len(row.Cells[1]) > 64 || !regionID.MatchString(row.Cells[1]) {
			return auxiliaryFailure(false, "region_aux_input_invalid")
		}
		for _, cell := range row.Cells {
			if cell == "" || !auxiliaryText(cell, &textBytes) {
				return auxiliaryFailure(false, "region_aux_text_limit")
			}
		}
		for i := 6; i <= 8; i++ {
			count, err := strconv.ParseInt(row.Cells[i], 10, 64)
			if err != nil || count < -MaxAuxCount || count > MaxAuxCount || i < 8 && count < 0 || strconv.FormatInt(count, 10) != row.Cells[i] {
				return auxiliaryFailure(false, "region_aux_input_invalid")
			}
		}
		switch row.Cells[9] {
		case "Idle", "Available", "At-Capacity", "Over-Allocated":
		default:
			return auxiliaryFailure(false, "region_aux_input_invalid")
		}
		key := [5]string{strings.ToLower(row.SubscriptionID), row.Cells[1], row.Cells[2], row.Cells[3], row.Cells[4]}
		if seen[key] {
			return auxiliaryFailure(false, "region_aux_input_invalid")
		}
		seen[key] = true
	}
	if len(input) == 0 {
		return nil, nil
	}
	table := pendingAuxiliary(false)
	table.Rows = make([]assessment.PluginRow, 0, len(input))
	for _, row := range input {
		if ctx.Err() != nil {
			return auxiliaryCancelled(ctx, false)
		}
		table.Rows = append(table.Rows, assessment.PluginRow{SubscriptionID: strings.ToLower(row.SubscriptionID), Cells: slices.Clone(row.Cells)})
	}
	return completeAuxiliary(ctx, table, false)
}
