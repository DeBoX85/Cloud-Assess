package plugins

import (
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/carbon"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/servicehealth"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/sqleol"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
)

const ZoneMapping = "zone-mapping"

// InternalMetadata advertises implemented adapters, not the full source registry.
// Execution/CLI availability remains a separate acceptance boundary.
func InternalMetadata() []assessment.PluginMetadata {
	return []assessment.PluginMetadata{{Name: ZoneMapping, Version: "1.0.0", Description: "Retrieves logical-to-physical availability zone mappings for all Azure regions in each subscription", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}, servicehealth.Metadata(), sqleol.Metadata(), carbon.Metadata()}
}

func ValidateNames(names []string) ([]string, error) {
	selected := map[string]bool{}
	for _, name := range names {
		if name != ZoneMapping && name != servicehealth.Name && name != sqleol.Name && name != carbon.Name {
			return nil, fmt.Errorf("internal plugin %q is unavailable", name)
		}
		selected[name] = true
	}
	result := []string{}
	for _, name := range []string{carbon.Name, servicehealth.Name, sqleol.Name, ZoneMapping} {
		if selected[name] {
			result = append(result, name)
		}
	}
	return result, nil
}

func PendingTable(name string) assessment.PluginTable {
	if name == carbon.Name {
		return carbon.PendingTable()
	}
	if name == servicehealth.Name {
		return servicehealth.PendingTable()
	}
	if name == sqleol.Name {
		return sqleol.PendingTable()
	}
	return PendingZoneTable()
}

// ServiceHealthTable validates the adapter identity and owns all mutable data.
// Raw injected/provider health text never crosses the orchestration boundary.
func ServiceHealthTable(value assessment.PluginTable, scanErr error, subscriptions map[string]string, started, finished time.Time) assessment.PluginTable {
	pending := servicehealth.PendingTable()
	invalid := func() assessment.PluginTable {
		t := pending
		t.Health = assessment.StageExecution{Name: servicehealth.Name, Status: assessment.StageFailed, StartedAt: started.UTC(), FinishedAt: finished.UTC(), Error: &assessment.AssessmentError{Code: "service_health_output_invalid", Message: "service-health output could not be safely projected"}}
		return t
	}
	if value.SchemaVersion != pending.SchemaVersion || value.ID != pending.ID || value.Metadata != pending.Metadata || !slices.Equal(value.Columns, pending.Columns) || len(value.Health.Warnings) > 4 || assessment.ValidatePluginTables([]assessment.PluginTable{value}) != nil {
		return invalid()
	}
	if value.SheetName == pending.SheetName {
		if value.Description != pending.Description || len(value.Rows) != 0 {
			return invalid()
		}
	} else if value.SheetName != "Service Issues" || value.Description != "Azure service health availability analysis showing percentage of time without service health events by subscription, region, and resource type (last 90 days)" {
		return invalid()
	}
	if value.Health.Name != servicehealth.Name || value.Health.Status == assessment.StageSkipped && scanErr == nil || value.SheetName == pending.SheetName && (value.Health.Status == assessment.StageCompleted || value.Health.Status == assessment.StageCompletedWithWarnings) {
		return invalid()
	}
	scope := map[string]bool{}
	for id := range subscriptions {
		scope[strings.ToLower(id)] = true
	}
	for _, r := range value.Rows {
		if len(r.Cells) != 6 || r.SubscriptionID != r.Cells[0] || !scope[strings.ToLower(r.SubscriptionID)] {
			return invalid()
		}
		for _, label := range r.Cells[1:3] {
			if label == "" || len(label) > 512 || strings.IndexFunc(label, unicode.IsControl) >= 0 {
				return invalid()
			}
		}
		percentage, err := strconv.ParseFloat(strings.TrimSuffix(r.Cells[3], "%"), 64)
		if err != nil || math.IsNaN(percentage) || math.IsInf(percentage, 0) || percentage < 0 || percentage > 100 || fmt.Sprintf("%.2f%%", percentage) != r.Cells[3] {
			return invalid()
		}
		for _, cell := range r.Cells[4:] {
			n, err := strconv.ParseInt(cell, 10, 64)
			if err != nil || n < 0 || strconv.FormatInt(n, 10) != cell {
				return invalid()
			}
		}
	}
	t := value
	t.Columns = append([]string(nil), value.Columns...)
	t.Rows = make([]assessment.PluginRow, len(value.Rows))
	for i, r := range value.Rows {
		t.Rows[i] = assessment.PluginRow{SubscriptionID: r.SubscriptionID, Cells: append([]string(nil), r.Cells...)}
	}
	t.Health = assessment.StageExecution{Name: servicehealth.Name, Status: value.Health.Status, Records: len(t.Rows), StartedAt: started.UTC(), FinishedAt: finished.UTC()}
	for _, warning := range value.Health.Warnings {
		code, message := "service_health_warning", "service-health reported incomplete input data"
		var count int
		if warning.Code == "service_health_malformed_rows" {
			code = warning.Code
			if _, err := fmt.Sscanf(warning.Message, "skipped %d invalid service-health rows", &count); err == nil && count > 0 && count <= servicehealth.MaxRows {
				message = fmt.Sprintf("skipped %d invalid service-health rows", count)
			}
		}
		t.Health.Warnings = append(t.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: message})
	}
	if value.Health.Status == assessment.StageFailed || scanErr != nil {
		code := "service_health_query_failed"
		if value.Health.Error != nil {
			switch value.Health.Error.Code {
			case "service_health_cancelled", "service_health_scope_limit", "service_health_scope_invalid", "service_health_not_configured", "service_health_output_invalid":
				code = value.Health.Error.Code
			}
		}
		t.Health.Status = assessment.StageFailed
		t.Health.Error = &assessment.AssessmentError{Code: code, Message: "service-health returned incomplete data"}
	}
	return t
}

func PendingZoneTable() assessment.PluginTable {
	m := InternalMetadata()[0]
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "zones", Metadata: m, SheetName: "Zone Mapping", Description: "Logical-to-physical availability zone mappings for Azure regions by subscription", Columns: []string{"Subscription", "Location", "Display Name", "Logical Zone", "Physical Zone"}, Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: ZoneMapping, Status: assessment.StageSkipped, Warnings: []assessment.AssessmentWarning{{Code: "plugin_not_run", Message: "requested plugin has not executed"}}}}
}

var subscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ZoneTable owns row/cell/health slices and never copies provider errors/URLs.
// A malformed trusted operation result becomes a visible failed empty table.
func ZoneTable(value zone.Result, scanErr error, started, finished time.Time) assessment.PluginTable {
	t := PendingZoneTable()
	t.Health = assessment.StageExecution{Name: ZoneMapping, Status: assessment.StageCompleted, StartedAt: started.UTC(), FinishedAt: finished.UTC()}
	if len(value.Rows) > zone.MaxRows || len(value.Failures) > zone.MaxSubscriptions {
		return invalidZoneTable(started, finished)
	}
	for _, r := range value.Rows {
		t.Rows = append(t.Rows, assessment.PluginRow{SubscriptionID: r.SubscriptionID, Cells: []string{r.SubscriptionName, r.Location, r.DisplayName, r.LogicalZone, r.PhysicalZone}})
	}
	t.Health.Records = len(t.Rows)
	if scanErr != nil || len(value.Failures) > 0 {
		t.Health.Status = assessment.StageFailed
		t.Health.Error = &assessment.AssessmentError{Code: "zone_incomplete", Message: "zone mapping returned incomplete data"}
	}
	for _, f := range value.Failures {
		code := f.Code
		switch code {
		case "zone_cancelled", "zone_page_limit", "zone_unsafe_continuation", "zone_continuation_cycle", "zone_byte_limit", "zone_invalid_response", "zone_row_limit":
		default:
			code = "zone_request_failed"
		}
		message := "zone mapping subscription retrieval failed"
		if subscriptionID.MatchString(f.SubscriptionID) {
			message = "subscription id " + f.SubscriptionID + " zone retrieval incomplete"
		}
		t.Health.Warnings = append(t.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: message})
	}
	if assessment.ValidatePluginTables([]assessment.PluginTable{t}) != nil {
		return invalidZoneTable(started, finished)
	}
	return t
}

func invalidZoneTable(started, finished time.Time) assessment.PluginTable {
	t := PendingZoneTable()
	t.Health = assessment.StageExecution{Name: ZoneMapping, Status: assessment.StageFailed, StartedAt: started.UTC(), FinishedAt: finished.UTC(), Error: &assessment.AssessmentError{Code: "zone_output_invalid", Message: "zone mapping output could not be safely projected"}}
	return t
}
