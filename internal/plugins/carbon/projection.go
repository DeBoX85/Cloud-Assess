// Carbon report behavior is characterized from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md and docs/CARBON_EMISSIONS.md.
package carbon

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const (
	Name             = "carbon-emissions"
	MaxItems         = 65536
	MaxResourceTypes = 4096
	MaxTextBytes     = 16 << 20
	MaxLabelBytes    = 512
)

var columns = []string{"Period From", "Period To", "Resource Type", "Latest Month Emissions", "Previous Month Emissions", "Month-over-Month Change Ratio", "Monthly Change Value", "Unit"}

func Metadata() assessment.PluginMetadata {
	return assessment.PluginMetadata{Name: Name, Version: "1.0.0", Description: "Analyzes carbon emissions by Azure resource type", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}
}

func PendingTable() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "emissions", Metadata: Metadata(), SheetName: "Carbon Emissions", Description: "Analysis of carbon emissions by Azure resource type for the previous month", Columns: append([]string(nil), columns...), Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped, Warnings: []assessment.AssessmentWarning{{Code: "plugin_not_run", Message: "requested plugin has not executed"}}}}
}

// Item is decoded resource-type data, not a per-subscription row. Optional
// numeric pointers preserve missing versus zero input without inventing totals.
// HTTP shape/access/pagination validation belongs to the bounded request adapter.
type Item struct {
	ResourceType string
	Latest       *float64
	Previous     *float64
	Change       *float64
}
type Filter interface{ IsResourceTypeExcluded(string) bool }
type total struct{ latest, previous, change float64 }

// LatestPeriod preserves the source's service-selected latest available date.
// Reject reversed date ranges rather than trusting contradictory service data.
func LatestPeriod(start, end string) (time.Time, error) {
	a, err := time.Parse("2006-01-02", start)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid carbon available date range")
	}
	b, err := time.Parse("2006-01-02", end)
	if err != nil || a.After(b) {
		return time.Time{}, fmt.Errorf("invalid carbon available date range")
	}
	return b, nil
}

// Project aggregates source values across already selected subscription batches.
// It makes no HTTP requests and cannot establish scope/access completeness. Rows
// are sorted by exact resource-type text because source map iteration is unstable.
func Project(ctx context.Context, start, end string, items []Item, filter Filter) (assessment.PluginTable, error) {
	table := PendingTable()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted}
	sums := map[string]total{}
	malformed, textBytes := 0, 0
	finish := func() {
		keys := make([]string, 0, len(sums))
		for k := range sums {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		table.Rows = make([]assessment.PluginRow, 0, len(keys))
		for _, k := range keys {
			table.Rows = append(table.Rows, assessment.PluginRow{Cells: emissionRow(end, k, sums[k])})
		}
		table.Health.Records = len(table.Rows)
		if malformed > 0 {
			table.Health.Warnings = []assessment.AssessmentWarning{{Code: "carbon_malformed_items", Message: fmt.Sprintf("skipped %d invalid carbon emission items", malformed)}}
		}
	}
	fail := func(code string, cause error) (assessment.PluginTable, error) {
		finish()
		table.Health.Status = assessment.StageFailed
		table.Health.Error = &assessment.AssessmentError{Code: code, Message: "carbon emissions projection returned incomplete data"}
		if cause != nil {
			return table, fmt.Errorf("carbon emissions projection stopped: %w", cause)
		}
		return table, fmt.Errorf("carbon emissions projection returned incomplete data")
	}
	if _, err := LatestPeriod(start, end); err != nil {
		return fail("carbon_date_invalid", nil)
	}
	if len(items) > MaxItems {
		return fail("carbon_item_limit", nil)
	}
	for _, item := range items {
		if err := ctx.Err(); err != nil {
			return fail("carbon_cancelled", err)
		}
		textBytes += len(item.ResourceType)
		if textBytes > MaxTextBytes {
			return fail("carbon_text_limit", nil)
		}
		if !validLabel(item.ResourceType) || item.Latest == nil || !finite(*item.Latest) || item.Previous != nil && !finite(*item.Previous) || item.Change != nil && !finite(*item.Change) {
			malformed++
			continue
		}
		if filter != nil && filter.IsResourceTypeExcluded(item.ResourceType) {
			continue
		}
		old, exists := sums[item.ResourceType]
		if !exists && len(sums) >= MaxResourceTypes {
			return fail("carbon_type_limit", nil)
		}
		next := old
		next.latest += *item.Latest
		if item.Previous != nil {
			next.previous += *item.Previous
		}
		if item.Change != nil {
			next.change += *item.Change
		}
		if !finite(next.latest) || !finite(next.previous) || !finite(next.change) || next.previous != 0 && !finite(((next.latest-next.previous)/next.previous)*100) {
			return fail("carbon_numeric_limit", nil)
		}
		sums[item.ResourceType] = next
	}
	if err := ctx.Err(); err != nil {
		return fail("carbon_cancelled", err)
	}
	finish()
	if malformed > 0 {
		table.Health.Status = assessment.StageCompletedWithWarnings
	}
	return table, nil
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func validLabel(v string) bool {
	return v != "" && len(v) <= MaxLabelBytes && utf8.ValidString(v) && strings.IndexFunc(v, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffe || r == 0xffff }) < 0
}
func emissionRow(date, kind string, agg total) []string {
	row := []string{date, date, kind, fmt.Sprintf("%.2f", agg.latest), "", "", "", "kgCO2e"}
	if agg.previous > 0 {
		row[4] = fmt.Sprintf("%.2f", agg.previous)
	}
	if agg.previous != 0 {
		row[5] = fmt.Sprintf("%.2f%%", ((agg.latest-agg.previous)/agg.previous)*100)
	}
	if agg.change != 0 {
		row[6] = fmt.Sprintf("%.2f", agg.change)
	}
	return row
}
