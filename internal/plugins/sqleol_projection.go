package plugins

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/sqleol"
)

// SQLEOLTable preserves source text formatting; subscription correlation is an
// owned row identity, distinct from the displayed subscription name in cell 0.
func SQLEOLTable(value assessment.PluginTable, scanErr error, subscriptions map[string]string, started, finished time.Time) assessment.PluginTable {
	pending := sqleol.PendingTable()
	invalid := func() assessment.PluginTable {
		t := pending
		t.Health = assessment.StageExecution{Name: sqleol.Name, Status: assessment.StageFailed, StartedAt: started.UTC(), FinishedAt: finished.UTC(), Error: &assessment.AssessmentError{Code: "sql_eol_output_invalid", Message: "sql-eol output could not be safely projected"}}
		return t
	}
	if value.SchemaVersion != pending.SchemaVersion || value.ID != pending.ID || value.Metadata != pending.Metadata || !slices.Equal(value.Columns, pending.Columns) || value.SheetName != pending.SheetName || value.Description != pending.Description || len(value.Health.Warnings) > 4 || assessment.ValidatePluginTables([]assessment.PluginTable{value}) != nil || value.Health.Status == assessment.StageSkipped && scanErr == nil {
		return invalid()
	}
	scope := map[string]bool{}
	for id := range subscriptions {
		scope[strings.ToLower(id)] = true
	}
	for _, row := range value.Rows {
		if !scope[strings.ToLower(row.SubscriptionID)] {
			return invalid()
		}
		for _, cell := range row.Cells {
			if len(cell) > 4096 || strings.IndexFunc(cell, unicode.IsControl) >= 0 {
				return invalid()
			}
		}
	}
	t := value
	t.Columns = append([]string(nil), value.Columns...)
	t.Rows = make([]assessment.PluginRow, len(value.Rows))
	for i, row := range value.Rows {
		t.Rows[i] = assessment.PluginRow{SubscriptionID: row.SubscriptionID, Cells: append([]string(nil), row.Cells...)}
	}
	t.Health = assessment.StageExecution{Name: sqleol.Name, Status: value.Health.Status, Records: len(t.Rows), StartedAt: started.UTC(), FinishedAt: finished.UTC()}
	for _, warning := range value.Health.Warnings {
		code, message := "sql_eol_warning", "sql-eol reported incomplete input data"
		var count int
		if warning.Code == "sql_eol_malformed_rows" {
			code = warning.Code
			if _, err := fmt.Sscanf(warning.Message, "skipped %d invalid sql-eol rows", &count); err == nil && count > 0 && count <= sqleol.MaxRows {
				message = fmt.Sprintf("skipped %d invalid sql-eol rows", count)
			}
		}
		t.Health.Warnings = append(t.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: message})
	}
	if value.Health.Status == assessment.StageFailed || scanErr != nil {
		code := "sql_eol_query_failed"
		if value.Health.Error != nil {
			switch value.Health.Error.Code {
			case "sql_eol_cancelled", "sql_eol_scope_limit", "sql_eol_scope_invalid", "sql_eol_not_configured", "sql_eol_output_invalid":
				code = value.Health.Error.Code
			}
		}
		t.Health.Status = assessment.StageFailed
		t.Health.Error = &assessment.AssessmentError{Code: code, Message: "sql-eol returned incomplete data"}
	}
	return t
}
