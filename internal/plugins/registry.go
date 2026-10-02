package plugins

import (
	"fmt"
	"regexp"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
)

const ZoneMapping = "zone-mapping"

// InternalMetadata advertises implemented adapters, not the full source registry.
// Execution/CLI availability remains a separate acceptance boundary.
func InternalMetadata() []assessment.PluginMetadata {
	return []assessment.PluginMetadata{{Name: ZoneMapping, Version: "1.0.0", Description: "Retrieves logical-to-physical availability zone mappings for all Azure regions in each subscription", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}}
}

func ValidateNames(names []string) ([]string, error) {
	for _, name := range names {
		if name != ZoneMapping {
			return nil, fmt.Errorf("internal plugin %q is unavailable", name)
		}
	}
	if len(names) == 0 {
		return nil, nil
	}
	return []string{ZoneMapping}, nil
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
