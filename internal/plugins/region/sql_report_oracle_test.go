// Independent literal SQL report expectations frozen before report APIs exist.
package region

import "github.com/DeBoX85/Cloud-Assess/internal/assessment"

const sqlReportOracleID = "aaaaaaaa-1111-4111-8111-aaaaaaaaaaaa"

type sqlReportOracleCase struct {
	name       string
	selected   map[string]string
	request    SQLCapabilityRequest
	collection *SQLCapabilityCollection
	required   []string
	want       assessment.PluginTable
	wantError  string
}

func sqlReportExpected(family string, cells [][]string, unknown bool, code string) assessment.PluginTable {
	id, sheet, description := "sql-capability-catalog", "SQL Capability Catalog", "SQL capability catalog diagnostics; deployment eligibility is not assessed."
	if family == "db" {
		id, sheet, description = "sql-database-catalog", "SQL Database Catalog", "SQL Database capability catalog diagnostics; deployment eligibility is not assessed."
	}
	if family == "mi" {
		id, sheet, description = "sql-mi-catalog", "SQL MI Catalog", "SQL Managed Instance capability catalog diagnostics; deployment eligibility is not assessed."
	}
	table := assessment.PluginTable{SchemaVersion: "1.0", ID: id, Metadata: assessment.PluginMetadata{Name: "region-selection", Version: "1.0.0", Description: "Bounded SQL capability catalog diagnostics", Author: "Cloud Assess", License: "MIT", Type: "internal"}, SheetName: sheet, Description: description,
		Columns: []string{"Subscription ID", "Region", "Resource type", "Capability coordinate", "Version", "Edition", "Name", "Family", "SKU name", "SKU tier", "Capacity", "Standalone", "Instance pool", "Catalog availability", "Deployment eligibility"}, Rows: []assessment.PluginRow{}}
	for _, row := range cells {
		table.Rows = append(table.Rows, assessment.PluginRow{SubscriptionID: sqlReportOracleID, Cells: row})
	}
	table.Health = assessment.StageExecution{Name: "region-selection", Status: assessment.StageCompletedWithWarnings, Records: len(cells), Warnings: []assessment.AssessmentWarning{{Code: "sql_capability_catalog_only", Message: "SQL capability catalog evidence does not establish deployment eligibility, quota or allocation capacity."}}}
	if unknown {
		table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "sql_capability_catalog_unknown", Message: "Required SQL capability catalog evidence is unknown; no capability rows are reported."})
	}
	if code != "" {
		table.Health = assessment.StageExecution{Name: "region-selection", Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "SQL capability report could not be safely projected"}}
	}
	return table
}

func sqlReportOracleCases() []sqlReportOracleCase {
	scope := map[string]string{sqlReportOracleID: "selected"}
	dbq := SQLCapabilityRequest{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: "microsoft.sql/servers/databases"}
	miq := SQLCapabilityRequest{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: "microsoft.sql/managedinstances"}
	db := func() *SQLCapabilityCollection {
		return &SQLCapabilityCollection{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: "microsoft.sql/servers/databases", Status: "complete", Groups: map[string]SQLCapabilityGroup{
			"12/basic/o1":    {Known: true, State: 0, Version: "12", Edition: "basic", Name: "o1", SKUName: "same_sku", SKUTier: "basic", Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"},
			"12/standard/o2": {Known: true, State: 2, Version: "12", Edition: "standard", Name: "o2", Family: "family", SKUName: "same_sku", SKUTier: "standard", Capacity: 2, Standalone: "not_applicable", InstancePool: "not_applicable"}}}
	}
	mi := func() *SQLCapabilityCollection {
		return &SQLCapabilityCollection{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: "microsoft.sql/managedinstances", Status: "complete", Groups: map[string]SQLCapabilityGroup{
			"12/generalpurpose/gen5/c8":  {Known: true, State: 0, Version: "12", Edition: "generalpurpose", Name: "c8", Family: "gen5", SKUName: "explicit_sku", SKUTier: "generalpurpose", Capacity: 8, Standalone: "unsupported", InstancePool: "unsupported"},
			"12/generalpurpose/gen8/c16": {Known: true, State: 2, Version: "12", Edition: "generalpurpose", Name: "c16", Family: "gen8", SKUName: "explicit_sku", SKUTier: "generalpurpose", Capacity: 16, Standalone: "unknown", InstancePool: "supported"}}}
	}
	dbRows := [][]string{
		{sqlReportOracleID, "eastus", "microsoft.sql/servers/databases", "12/basic/o1", "12", "basic", "o1", "", "same_sku", "basic", "1", "not_applicable", "not_applicable", "available", "not_assessed"},
		{sqlReportOracleID, "eastus", "microsoft.sql/servers/databases", "12/standard/o2", "12", "standard", "o2", "family", "same_sku", "standard", "2", "not_applicable", "not_applicable", "unavailable", "not_assessed"}}
	miRows := [][]string{
		{sqlReportOracleID, "eastus", "microsoft.sql/managedinstances", "12/generalpurpose/gen5/c8", "12", "generalpurpose", "c8", "gen5", "explicit_sku", "generalpurpose", "8", "unsupported", "unsupported", "available", "not_assessed"},
		{sqlReportOracleID, "eastus", "microsoft.sql/managedinstances", "12/generalpurpose/gen8/c16", "12", "generalpurpose", "c16", "gen8", "explicit_sku", "generalpurpose", "16", "unknown", "supported", "unavailable", "not_assessed"}}
	dbUnknown := db()
	dbUnknown.Status = "unknown"
	dbUnknown.FailureCode = "sql_capability_request_failed"
	dbUnknown.Groups = map[string]SQLCapabilityGroup{}
	miUnknown := mi()
	miUnknown.Status = "unknown"
	miUnknown.FailureCode = "sql_capability_incomplete"
	miUnknown.Groups = map[string]SQLCapabilityGroup{}
	dbRequiredUnknown := db()
	g := dbRequiredUnknown.Groups["12/basic/o1"]
	g.Known = false
	g.State = 2
	dbRequiredUnknown.Groups["12/basic/o1"] = g
	miRequiredUnknown := mi()
	g = miRequiredUnknown.Groups["12/generalpurpose/gen5/c8"]
	g.Known = false
	g.State = 2
	miRequiredUnknown.Groups["12/generalpurpose/gen5/c8"] = g
	dbBad := db()
	g = dbBad.Groups["12/standard/o2"]
	g.Capacity = -1
	dbBad.Groups["12/standard/o2"] = g
	miBad := mi()
	g = miBad.Groups["12/generalpurpose/gen8/c16"]
	g.Standalone = "choose_one"
	miBad.Groups["12/generalpurpose/gen8/c16"] = g
	wrongBinding := db()
	wrongBinding.Region = "westus"
	inconsistent := db()
	inconsistent.Status = "unknown"
	inconsistent.FailureCode = "sql_capability_incomplete"
	badKind := dbq
	badKind.ResourceType = "microsoft.sql/servers/elasticpools"
	badAliases := dbq
	badAliases.LocationAliases = map[string]string{"global": "eastus"}
	unselected := dbq
	unselected.SubscriptionID = "bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"
	invalid := sqlReportExpected("", nil, false, "sql_capability_report_invalid")
	invalidError := "SQL capability report rejected [sql_capability_report_invalid]"
	return []sqlReportOracleCase{
		{"db-full-coordinates-same-sku-sorted", scope, dbq, db(), []string{"12/standard/o2", "12/basic/o1"}, sqlReportExpected("db", dbRows, false, ""), ""},
		{"mi-full-coordinates-same-sku-modes-sorted", scope, miq, mi(), []string{"12/generalpurpose/gen8/c16", "12/generalpurpose/gen5/c8"}, sqlReportExpected("mi", miRows, false, ""), ""},
		{"db-known-available", scope, dbq, db(), []string{"12/basic/o1"}, sqlReportExpected("db", dbRows[:1], false, ""), ""},
		{"mi-known-unavailable", scope, miq, mi(), []string{"12/generalpurpose/gen8/c16"}, sqlReportExpected("mi", miRows[1:], false, ""), ""},
		{"db-empty-requirements", scope, dbq, db(), nil, sqlReportExpected("db", nil, false, ""), ""},
		{"mi-empty-requirements", scope, miq, mi(), nil, sqlReportExpected("mi", nil, false, ""), ""},
		{"db-nil-collection", scope, dbq, nil, nil, sqlReportExpected("db", nil, true, ""), ""},
		{"mi-nil-collection", scope, miq, nil, nil, sqlReportExpected("mi", nil, true, ""), ""},
		{"db-operational-unknown", scope, dbq, dbUnknown, []string{"12/basic/o1"}, sqlReportExpected("db", nil, true, ""), ""},
		{"mi-incomplete-unknown", scope, miq, miUnknown, nil, sqlReportExpected("mi", nil, true, ""), ""},
		{"db-missing-required", scope, dbq, db(), []string{"12/basic/missing"}, sqlReportExpected("db", nil, true, ""), ""},
		{"mi-missing-required", scope, miq, mi(), []string{"12/generalpurpose/gen5/missing"}, sqlReportExpected("mi", nil, true, ""), ""},
		{"db-unknown-required", scope, dbq, dbRequiredUnknown, []string{"12/basic/o1"}, sqlReportExpected("db", nil, true, ""), ""},
		{"mi-unknown-required", scope, miq, miRequiredUnknown, []string{"12/generalpurpose/gen5/c8"}, sqlReportExpected("mi", nil, true, ""), ""},
		{"db-unrequired-unknown", scope, dbq, dbRequiredUnknown, []string{"12/standard/o2"}, sqlReportExpected("db", dbRows[1:], false, ""), ""},
		{"mi-unrequired-unknown", scope, miq, miRequiredUnknown, []string{"12/generalpurpose/gen8/c16"}, sqlReportExpected("mi", miRows[1:], false, ""), ""},
		{"db-malformed-unrequired", scope, dbq, dbBad, []string{"12/basic/o1"}, invalid, invalidError},
		{"mi-malformed-unrequired", scope, miq, miBad, []string{"12/generalpurpose/gen5/c8"}, invalid, invalidError},
		{"normalized-required-collision", scope, dbq, db(), []string{"12/basic/o1", "12/BASIC/o1"}, invalid, invalidError},
		{"no-name-only-lookup", scope, miq, mi(), []string{"explicit_sku"}, invalid, invalidError},
		{"collection-binding-invalid", scope, dbq, wrongBinding, nil, invalid, invalidError},
		{"unknown-must-have-empty-groups", scope, dbq, inconsistent, nil, invalid, invalidError},
		{"unsupported-family", scope, badKind, nil, nil, invalid, invalidError},
		{"empty-selection", map[string]string{}, dbq, nil, nil, invalid, invalidError},
		{"invalid-alias", scope, badAliases, nil, nil, invalid, invalidError},
		{"unselected-subscription", scope, unselected, nil, nil, invalid, invalidError},
	}
}
