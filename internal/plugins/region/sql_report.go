// SQL report adapters preserve catalog-only evidence and existing report schemas.
package region

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sort"
	"strconv"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const sqlReportTextCap = 8 << 20

func sqlReportColumns() []string {
	return []string{"Subscription ID", "Region", "Resource type", "Capability coordinate", "Version", "Edition", "Name", "Family", "SKU name", "SKU tier", "Capacity", "Standalone", "Instance pool", "Catalog availability", "Deployment eligibility"}
}
func sqlReportBase(kind string) assessment.PluginTable {
	id, sheet, description := "sql-capability-catalog", "SQL Capability Catalog", "SQL capability catalog diagnostics; deployment eligibility is not assessed."
	switch kind {
	case sqlDatabaseType:
		id, sheet, description = "sql-database-catalog", "SQL Database Catalog", "SQL Database capability catalog diagnostics; deployment eligibility is not assessed."
	case sqlManagedInstanceType:
		id, sheet, description = "sql-mi-catalog", "SQL MI Catalog", "SQL Managed Instance capability catalog diagnostics; deployment eligibility is not assessed."
	}
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: id, Metadata: assessment.PluginMetadata{Name: Name, Version: "1.0.0", Description: "Bounded SQL capability catalog diagnostics", Author: "Cloud Assess", License: "MIT", Type: "internal"}, SheetName: sheet, Description: description, Columns: sqlReportColumns(), Rows: []assessment.PluginRow{}}
}
func sqlReportFailure(ctx context.Context, kind, code string, cause error) (assessment.PluginTable, error) {
	var terminal error
	if err := ctx.Err(); err != nil {
		terminal = err
	}
	if terminal == nil {
		if errors.Is(cause, context.Canceled) {
			terminal = context.Canceled
		} else if errors.Is(cause, context.DeadlineExceeded) {
			terminal = context.DeadlineExceeded
		}
	}
	if terminal != nil {
		code = "sql_capability_report_cancelled"
	}
	table := sqlReportBase(kind)
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "SQL capability report could not be safely projected"}}
	if terminal != nil {
		return table, terminal
	}
	return table, fmt.Errorf("SQL capability report rejected [%s]", code)
}
func sqlReportProjectionCode(err error) string {
	if err != nil && err.Error() == "SQL capability catalog projection rejected [sql_capability_projection_limit]" {
		return "sql_capability_report_limit"
	}
	return "sql_capability_report_invalid"
}
func sqlReportAvailability(g SQLCapabilityGroup) string {
	if g.State == SKUAvailable {
		return "available"
	}
	return "unavailable"
}
func sqlReportHeaderBytes(table assessment.PluginTable) int {
	n := len(table.SheetName) + len(table.Metadata.Version) + len(table.Metadata.Description) + len(table.Metadata.Author) + len(table.Metadata.License) + len(table.Description)
	for _, s := range table.Columns {
		n += len(s)
	}
	for _, w := range table.Health.Warnings {
		n += len(w.Message)
	}
	return n
}
func sqlReportRowBytes(c SQLCatalogEvidence, key string, g SQLCapabilityGroup) int {
	digits := 1
	for n := g.Capacity; n >= 10; n /= 10 {
		digits++
	}
	return len(c.SubscriptionID) + len(c.Region) + len(c.ResourceType) + len(key) + len(g.Version) + len(g.Edition) + len(g.Name) + len(g.Family) + len(g.SKUName) + len(g.SKUTier) + digits + len(g.Standalone) + len(g.InstancePool) + len(sqlReportAvailability(g)) + len("not_assessed")
}
func sqlReportFinish(ctx context.Context, kind string, table assessment.PluginTable) (assessment.PluginTable, error) {
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		return sqlReportFailure(ctx, kind, "sql_capability_report_invalid", nil)
	}
	if err := ctx.Err(); err != nil {
		return sqlReportFailure(ctx, kind, "sql_capability_report_cancelled", err)
	}
	return table, nil
}

// ProjectSQLCapabilityTable structurally validates stable caller data through
// Catalog. It does not authenticate collection provenance or deployment eligibility.
func ProjectSQLCapabilityTable(ctx context.Context, selected map[string]string, q SQLCapabilityRequest, collection *SQLCapabilityCollection, required []string) (assessment.PluginTable, error) {
	catalog, err := collection.Catalog(ctx, selected, q, required)
	if err != nil {
		return sqlReportFailure(ctx, "", sqlReportProjectionCode(err), err)
	}
	kind := catalog.ResourceType
	table := sqlReportBase(kind)
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompletedWithWarnings, Warnings: []assessment.AssessmentWarning{{Code: "sql_capability_catalog_only", Message: sqlCatalogWarning}}}
	if catalog.Status == "unknown" {
		table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "sql_capability_catalog_unknown", Message: "Required SQL capability catalog evidence is unknown; no capability rows are reported."})
	}
	keys := make([]string, 0, len(catalog.Values))
	for key := range catalog.Values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	total := sqlReportHeaderBytes(table)
	if total > sqlReportTextCap {
		return sqlReportFailure(ctx, kind, "sql_capability_report_limit", nil)
	}
	// Admission precedes every nonempty row/cell allocation, so no rejected prefix escapes.
	for _, key := range keys {
		if err := ctx.Err(); err != nil {
			return sqlReportFailure(ctx, kind, "sql_capability_report_cancelled", err)
		}
		total += sqlReportRowBytes(catalog, key, catalog.Values[key])
		if total > sqlReportTextCap {
			return sqlReportFailure(ctx, kind, "sql_capability_report_limit", nil)
		}
	}
	table.Rows = make([]assessment.PluginRow, len(keys))
	for i, key := range keys {
		if err := ctx.Err(); err != nil {
			return sqlReportFailure(ctx, kind, "sql_capability_report_cancelled", err)
		}
		g := catalog.Values[key]
		table.Rows[i] = assessment.PluginRow{SubscriptionID: catalog.SubscriptionID, Cells: []string{catalog.SubscriptionID, catalog.Region, catalog.ResourceType, key, g.Version, g.Edition, g.Name, g.Family, g.SKUName, g.SKUTier, strconv.FormatInt(g.Capacity, 10), g.Standalone, g.InstancePool, sqlReportAvailability(g), "not_assessed"}}
	}
	table.Health.Records = len(table.Rows)
	return sqlReportFinish(ctx, kind, table)
}

// CollectSQLCapabilityTable connects only an injected existing bounded collector
// to diagnostic projection. It performs no discovery, credential lookup or writes.
func CollectSQLCapabilityTable(ctx context.Context, collector *SQLCapabilityCollector, selected map[string]string, q SQLCapabilityRequest, required []string) (assessment.PluginTable, error) {
	if err := ctx.Err(); err != nil {
		return sqlReportFailure(ctx, "", "sql_capability_report_cancelled", err)
	}
	var absent *SQLCapabilityCollection
	preflight, err := absent.Catalog(ctx, selected, q, required)
	if err != nil {
		return sqlReportFailure(ctx, "", sqlReportProjectionCode(err), err)
	}
	selected = maps.Clone(selected)
	q.LocationAliases = maps.Clone(q.LocationAliases)
	required = append([]string(nil), required...)
	collection, err := collector.Collect(ctx, selected, q)
	if err != nil {
		return sqlReportFailure(ctx, preflight.ResourceType, "sql_capability_report_invalid", err)
	}
	return ProjectSQLCapabilityTable(ctx, selected, q, collection, required)
}
