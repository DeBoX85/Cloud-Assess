package plugins

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/carbon"
)

// CarbonTable owns aggregate cells, validates the source contract and does not
// invent per-subscription identity. Raw operation/provider health stays private.
func CarbonTable(value assessment.PluginTable, scanErr error, subscriptions map[string]string, started, finished time.Time) assessment.PluginTable {
	pending := carbon.PendingTable()
	invalid := func() assessment.PluginTable {
		t := pending
		t.Health = assessment.StageExecution{Name: carbon.Name, Status: assessment.StageFailed, StartedAt: started.UTC(), FinishedAt: finished.UTC(), Error: &assessment.AssessmentError{Code: "carbon_output_invalid", Message: "carbon emissions output could not be safely projected"}}
		return t
	}
	if value.SchemaVersion != pending.SchemaVersion || value.ID != pending.ID || value.Metadata != pending.Metadata || !slices.Equal(value.Columns, pending.Columns) || value.SheetName != pending.SheetName || value.Description != pending.Description || len(value.Health.Warnings) > 4 || len(value.Rows) > carbon.MaxResourceTypes || assessment.ValidatePluginTables([]assessment.PluginTable{value}) != nil || value.Health.Status == assessment.StageSkipped && scanErr == nil || len(subscriptions) == 0 && len(value.Rows) > 0 {
		return invalid()
	}
	period, lastLabel := "", ""
	for _, row := range value.Rows {
		if row.SubscriptionID != "" {
			return invalid()
		}
		cells := row.Cells
		if _, err := time.Parse("2006-01-02", cells[0]); err != nil || cells[0] != cells[1] || period != "" && period != cells[0] || cells[7] != "kgCO2e" || cells[2] == "" || len(cells[2]) > carbon.MaxLabelBytes || strings.IndexFunc(cells[2], unicode.IsControl) >= 0 || lastLabel != "" && lastLabel >= cells[2] {
			return invalid()
		}
		period, lastLabel = cells[0], cells[2]
		for i := 3; i <= 6; i++ {
			if i != 3 && cells[i] == "" {
				continue
			}
			text, pattern := cells[i], "%.2f"
			if i == 5 {
				if !strings.HasSuffix(text, "%") {
					return invalid()
				}
				text = strings.TrimSuffix(text, "%")
				pattern = "%.2f%%"
			}
			n, err := strconv.ParseFloat(text, 64)
			if err != nil || math.IsNaN(n) || math.IsInf(n, 0) || fmt.Sprintf(pattern, n) != cells[i] || i == 4 && (n < 0 || math.Signbit(n)) {
				return invalid()
			}
		}
	}
	t := value
	t.Columns = append([]string(nil), value.Columns...)
	t.Rows = make([]assessment.PluginRow, len(value.Rows))
	for i, row := range value.Rows {
		t.Rows[i] = assessment.PluginRow{Cells: append([]string(nil), row.Cells...)}
	}
	t.Health = assessment.StageExecution{Name: carbon.Name, Status: value.Health.Status, Records: len(t.Rows), StartedAt: started.UTC(), FinishedAt: finished.UTC()}
	for _, warning := range value.Health.Warnings {
		code, message := "carbon_warning", "carbon emissions reported incomplete input data"
		pattern, maximum := "", 0
		switch warning.Code {
		case "carbon_malformed_items":
			pattern, maximum = "skipped %d invalid carbon emission items", carbon.MaxItems
		case "carbon_access_unverified":
			pattern, maximum = "%d carbon response pages have incomplete access decisions", carbon.MaxPages
		case "carbon_access_denied":
			pattern, maximum = "carbon access denied for %d selected subscriptions", carbon.MaxSubscriptions
		}
		if pattern != "" {
			code = warning.Code
			var count int
			if _, err := fmt.Sscanf(warning.Message, pattern, &count); err == nil && count > 0 && count <= maximum {
				message = fmt.Sprintf(pattern, count)
			}
		}
		t.Health.Warnings = append(t.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: message})
	}
	if value.Health.Status == assessment.StageFailed || scanErr != nil {
		code := "carbon_request_failed"
		if value.Health.Error != nil {
			switch value.Health.Error.Code {
			case "carbon_cancelled", "carbon_scope_limit", "carbon_scope_invalid", "carbon_not_configured", "carbon_date_request_failed", "carbon_date_invalid", "carbon_report_request_failed", "carbon_response_limit", "carbon_response_invalid", "carbon_access_invalid", "carbon_access_denied", "carbon_token_cycle", "carbon_page_limit", "carbon_item_limit", "carbon_type_limit", "carbon_text_limit", "carbon_numeric_limit", "carbon_output_invalid":
				code = value.Health.Error.Code
			}
		}
		t.Health.Status = assessment.StageFailed
		t.Health.Error = &assessment.AssessmentError{Code: code, Message: "carbon emissions returned incomplete data"}
	}
	return t
}
