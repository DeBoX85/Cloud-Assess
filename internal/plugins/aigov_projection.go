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
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/aigov"
)

// AIGovernanceTable validates the source table contract and owns cells/health.
// Provider and injected error text never crosses the report boundary.
func AIGovernanceTable(value assessment.PluginTable, scanErr error, subscriptions map[string]string, started, finished time.Time) assessment.PluginTable {
	pending := aigov.PendingTable()
	invalid := func() assessment.PluginTable {
		t := pending
		t.Health = assessment.StageExecution{Name: aigov.Name, Status: assessment.StageFailed, StartedAt: started.UTC(), FinishedAt: finished.UTC(), Error: &assessment.AssessmentError{Code: "ai_output_invalid", Message: "AI governance output could not be safely projected"}}
		return t
	}
	if value.SchemaVersion != pending.SchemaVersion || value.ID != pending.ID || value.Metadata != pending.Metadata || !slices.Equal(value.Columns, pending.Columns) || (value.SheetName != "AI Gov" && value.SheetName != pending.SheetName) || value.Description != pending.Description || len(value.Rows) > aigov.MaxRows || len(value.Health.Warnings) > 8 || assessment.ValidatePluginTables([]assessment.PluginTable{value}) != nil || value.Health.Status == assessment.StageSkipped && scanErr == nil || len(value.Rows) > 0 && value.SheetName != "AI Gov" {
		return invalid()
	}
	scope := map[string]string{}
	for id, name := range subscriptions {
		if name == "" {
			name = id
		}
		scope[strings.ToLower(id)] = name
	}
	for _, row := range value.Rows {
		name, ok := scope[strings.ToLower(row.SubscriptionID)]
		if !ok || row.Cells[0] != name {
			return invalid()
		}
		for _, cell := range row.Cells {
			if len(cell) > aigov.MaxLabelBytes || strings.IndexFunc(cell, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd }) >= 0 {
				return invalid()
			}
		}
		for _, i := range []int{1, 2} {
			if row.Cells[i] == "" || row.Cells[i] == "." || row.Cells[i] == ".." || strings.ContainsAny(row.Cells[i], "/%?#\\") {
				return invalid()
			}
		}
		if _, err := time.Parse("2006-01-02 15:04", row.Cells[13]); err != nil {
			return invalid()
		}
		if row.Cells[11] != "Yes" && row.Cells[11] != "No" {
			return invalid()
		}
		if row.Cells[9] != "N/A" {
			n, err := strconv.ParseInt(row.Cells[9], 10, 64)
			if err != nil || n < 0 || strconv.FormatInt(n, 10) != row.Cells[9] {
				return invalid()
			}
		}
		n, err := strconv.ParseFloat(row.Cells[15], 64)
		if err != nil || n < 0 || math.Signbit(n) || math.IsNaN(n) || math.IsInf(n, 0) || fmt.Sprintf("%.0f", n) != row.Cells[15] {
			return invalid()
		}
	}
	t := value
	t.Columns = slices.Clone(value.Columns)
	t.Rows = make([]assessment.PluginRow, len(value.Rows))
	for i, row := range value.Rows {
		t.Rows[i] = assessment.PluginRow{SubscriptionID: strings.ToLower(row.SubscriptionID), Cells: slices.Clone(row.Cells)}
	}
	t.Health = assessment.StageExecution{Name: aigov.Name, Status: value.Health.Status, Records: len(t.Rows), StartedAt: started.UTC(), FinishedAt: finished.UTC()}
	for _, warning := range value.Health.Warnings {
		code, message := "ai_warning", "AI governance reported incomplete input data"
		pattern, maximum := "", 0
		switch warning.Code {
		case "ai_discovery_invalid_rows":
			pattern, maximum = "skipped %d invalid or duplicate AI account rows", aigov.MaxAccounts
		case "ai_malformed_input":
			pattern, maximum = "skipped %d invalid AI governance input entries", aigov.MaxPoints+aigov.MaxDeployments
		case "ai_deployment_unavailable":
			pattern, maximum = "deployment enrichment incomplete for %d accounts", aigov.MaxAccounts
		case "ai_malformed_response":
			pattern, maximum = "skipped %d invalid AI response entries", aigov.MaxPoints+aigov.MaxDeployments
		case "ai_metrics_incomplete":
			pattern, maximum = "metrics coverage incomplete for %d account responses", aigov.MaxAccounts
		case "ai_deployment_requests_failed":
			pattern, maximum = "deployment retrieval incomplete for %d accounts", aigov.MaxAccounts
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
		code := "ai_request_failed"
		if value.Health.Error != nil {
			switch value.Health.Error.Code {
			case "ai_cancelled", "ai_scope_invalid", "ai_input_limit", "ai_account_invalid", "ai_text_limit", "ai_deployment_scope_invalid", "ai_deployment_limit", "ai_row_limit", "ai_numeric_limit", "ai_not_configured", "ai_metrics_request_failed", "ai_metrics_response_invalid", "ai_deployment_request_failed", "ai_request_limit", "ai_response_limit", "ai_output_invalid", "ai_discovery_cancelled", "ai_discovery_scope_limit", "ai_discovery_scope_invalid", "ai_discovery_request_limit", "ai_discovery_request_failed", "ai_discovery_body_limit", "ai_discovery_page_invalid", "ai_discovery_coverage_invalid", "ai_discovery_row_limit", "ai_discovery_text_limit":
				code = value.Health.Error.Code
			}
		}
		t.Health.Status = assessment.StageFailed
		t.Health.Error = &assessment.AssessmentError{Code: code, Message: "AI governance returned incomplete data"}
	}
	return t
}
