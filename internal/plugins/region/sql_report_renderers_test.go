package region

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	csvreport "github.com/DeBoX85/Cloud-Assess/internal/renderers/csv"
	excelreport "github.com/DeBoX85/Cloud-Assess/internal/renderers/excel"
	jsonreport "github.com/DeBoX85/Cloud-Assess/internal/renderers/json"
	sarifreport "github.com/DeBoX85/Cloud-Assess/internal/renderers/sarif"
	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/xuri/excelize/v2"
)

func sqlReportCoreInput() result.Input {
	id := "/subscriptions/" + sqlReportOracleID + "/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/vm"
	return result.Input{GeneratedAt: time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC), ScopeID: "/subscriptions/" + sqlReportOracleID + "/providers/Microsoft.Sql/locations/eastus",
		Scope:           &assessment.ScopeResolution{Selection: "explicit", Status: "resolved", IncludedSubscriptionIDs: []string{sqlReportOracleID}, ResolvedSubscriptions: []assessment.ScopeSubscription{{SubscriptionID: sqlReportOracleID, SubscriptionName: "selected"}}},
		Completeness:    assessment.CompletenessComplete,
		Recommendations: []assessment.RecommendationDefinition{{ID: "core-rule", Recommendation: "Core fixture", Category: "Security", Impact: "High", ResourceType: "microsoft.compute/virtualmachines", Source: "CUSTOM"}},
		Findings:        []assessment.Finding{{RecommendationID: "core-rule", Recommendation: "Core fixture", Category: "Security", Impact: "High", ResourceType: "microsoft.compute/virtualmachines", Source: "CUSTOM", ResourceID: id, SubscriptionID: sqlReportOracleID, ResourceGroup: "test", ResourceName: "vm"}}}
}
func sqlReportRendererFixture(t *testing.T, mode string) *result.AssessmentResult {
	t.Helper()
	plugin := []assessment.PluginTable{}
	for _, family := range []string{"db", "mi"} {
		scope, q, key, body := sqlReportTestInput(family)
		calls := 0
		ctorFamily := family
		if mode == "failed" {
			if family == "db" {
				ctorFamily = "mi"
			} else {
				ctorFamily = "db"
			}
		}
		col := sqlReportTestCollector(t, ctorFamily, storageGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			calls++
			if mode == "unknown" {
				return []byte("PROVIDER_CANARY_PRIVATE_BODY"), &http.Response{StatusCode: 503}, fmt.Errorf("PROVIDER_CANARY_PRIVATE_ERROR")
			}
			return []byte(body), &http.Response{StatusCode: 200}, nil
		}))
		table, e := CollectSQLCapabilityTable(context.Background(), col, scope, q, []string{key})
		if mode == "failed" {
			sqlReportRequireFailure(t, table, e, "sql_capability_report_invalid")
			if calls != 0 {
				t.Fatal("failed family report made GET")
			}
		} else if e != nil || calls != 1 {
			t.Fatal("actual report fixture", mode, e, calls)
		}
		if mode == "known" && len(table.Rows) != 1 || mode != "known" && len(table.Rows) != 0 {
			t.Fatal("fixture row state", mode, table)
		}
		if mode == "unknown" && (len(table.Health.Warnings) != 2 || table.Health.Warnings[1].Code != "sql_capability_catalog_unknown") {
			t.Fatal("unknown table health", table)
		}
		plugin = append(plugin, table)
	}
	input := sqlReportCoreInput()
	core := result.Build(input)
	data, e := result.BuildWithPluginTables(input, plugin)
	if e != nil {
		t.Fatal(e)
	}
	wantCompleteness := assessment.CompletenessCompleteWithWarnings
	if mode == "failed" {
		wantCompleteness = assessment.CompletenessPartial
	}
	if data.SchemaVersion != "1.1" || data.Completeness != wantCompleteness || !reflect.DeepEqual(data.Findings, core.Findings) || !reflect.DeepEqual(data.Summary, core.Summary) || !reflect.DeepEqual(data.Recommendations, core.Recommendations) {
		t.Fatal("diagnostics preserve core findings/scoring and expose health", data.Completeness)
	}
	// SARIF is deliberately primary-findings-only, not a plugin health export.
	before, e := sarifreport.Marshal(core, "report-test")
	if e != nil {
		t.Fatal(e)
	}
	after, e := sarifreport.Marshal(data, "report-test")
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("no synthetic catalog SARIF findings/rules or schema change", e)
	}
	var parsed struct {
		Runs []struct{ Results []json.RawMessage }
	}
	if e = json.Unmarshal(after, &parsed); e != nil || len(parsed.Runs) != 1 || len(parsed.Runs[0].Results) != 1 {
		t.Fatal("nonempty core SARIF retained", e)
	}
	for i := range plugin {
		plugin[i].Columns[0] = "mutated"
		if len(plugin[i].Rows) > 0 {
			plugin[i].Rows[0].Cells[0] = "mutated"
		}
	}
	if data.PluginTables[0].Columns[0] != "Subscription ID" {
		t.Fatal("result owns report tables")
	}
	return data
}
func sqlReportCheckMask(t *testing.T, text string, masked bool) {
	t.Helper()
	has := strings.Contains(strings.ToLower(text), sqlReportOracleID)
	if has == masked {
		t.Fatal("actual renderer subscription redaction", masked, has)
	}
	if strings.Contains(text, "PROVIDER_CANARY") {
		t.Fatal("raw provider diagnostic escaped")
	}
}
func TestSQLReportActualRenderers(t *testing.T) {
	for _, mode := range []string{"known", "unknown", "failed"} {
		for _, masked := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", mode, masked), func(t *testing.T) {
				data := sqlReportRendererFixture(t, mode)
				snapshot, e := json.Marshal(data)
				if e != nil {
					t.Fatal(e)
				}
				encoded, e := jsonreport.MarshalWithOptions(data, jsonreport.Options{RedactSubscriptionIDs: masked})
				if e != nil {
					t.Fatal(e)
				}
				sqlReportCheckMask(t, string(encoded), masked)
				text, e := jsonreport.StringWithOptions(data, jsonreport.Options{RedactSubscriptionIDs: masked})
				if e != nil || text != string(encoded) {
					t.Fatal("shared stdout JSON path", e)
				}
				var parsed map[string]any
				if e = json.Unmarshal(encoded, &parsed); e != nil || parsed["schemaVersion"] != "1.1" || len(parsed["pluginTables"].([]any)) != 2 {
					t.Fatal("canonical JSON plugin extension", e)
				}
				dir := t.TempDir()
				jsonPath := filepath.Join(dir, "assessment.json")
				if e = jsonreport.WriteFileWithOptions(data, jsonPath, jsonreport.Options{RedactSubscriptionIDs: masked}); e != nil {
					t.Fatal(e)
				}
				actual, e := os.ReadFile(jsonPath)
				if e != nil || !bytes.Equal(actual, encoded) {
					t.Fatal("actual JSON file readback", e)
				}
				opts := tables.Options{RedactSubscriptionIDs: masked}
				projected, e := tables.Build(data, opts)
				if e != nil {
					t.Fatal(e)
				}
				var status tables.Table
				found := map[string]tables.Table{}
				for _, table := range projected {
					if table.Key == "assessmentStatus" {
						status = table
					}
					if strings.HasPrefix(table.Key, "plugin_region-selection_sql-") {
						found[table.SheetName] = table
						if len(table.Rows) == 0 || len(table.Rows[0]) != 15 || table.Rows[0][3] != "Capability coordinate" || table.Rows[0][14] != "Deployment eligibility" {
							t.Fatal("diagnostic columns", table)
						}
						if mode == "known" {
							if len(table.Rows) != 2 || table.Rows[1][14] != "not_assessed" {
								t.Fatal("catalog-only rendered rows", table)
							}
							wantKey := "12/basic/o1"
							if table.SheetName == "SQL MI Catalog" {
								wantKey = "12/generalpurpose/gen5/c8"
							}
							if table.Rows[1][3] != wantKey {
								t.Fatal("full coordinate retained", table)
							}
						} else if len(table.Rows) != 1 {
							t.Fatal("unknown/failure has no capability prefix", table)
						}
						sqlReportCheckMask(t, strings.Join(table.Rows[len(table.Rows)-1], "|")+strings.Join(status.Rows[len(status.Rows)-1], "|"), masked)
					}
				}
				if len(found) != 2 {
					t.Fatal("both SQL families render separately", len(found))
				}
				statusCount := 0
				for _, row := range status.Rows {
					if len(row) > 8 && strings.HasPrefix(row[4], "plugin:region-selection:sql-") {
						statusCount++
						if !strings.Contains(row[2], "locations/eastus") {
							t.Fatal("empty diagnostic retains bound result region", row)
						}
						sqlReportCheckMask(t, strings.Join(row, "|"), masked)
						if mode == "failed" {
							if row[5] != "failed" || row[9] != "sql_capability_report_invalid" {
								t.Fatal("failed health visible", row)
							}
						} else if row[5] != "completed_with_warnings" || !strings.Contains(row[8], "sql_capability_catalog_only") || mode == "unknown" && !strings.Contains(row[8], "sql_capability_catalog_unknown") {
							t.Fatal("catalog limitation/unknown health visible", row)
						}
					}
				}
				if statusCount != 2 {
					t.Fatal("both status records present", statusCount)
				}
				paths, e := csvreport.Write(data, filepath.Join(dir, "assessment"), opts)
				if e != nil {
					t.Fatal(e)
				}
				for _, path := range paths {
					raw, e := os.ReadFile(path)
					if e != nil {
						t.Fatal(e)
					}
					if masked || strings.Contains(filepath.Base(path), "assessmentStatus") || mode == "known" && strings.Contains(filepath.Base(path), "plugin_region-selection_sql-") {
						sqlReportCheckMask(t, string(raw), masked)
					} else if strings.Contains(string(raw), "PROVIDER_CANARY") {
						t.Fatal("raw provider diagnostic escaped")
					}
					if strings.Contains(filepath.Base(path), "plugin_region-selection_sql-") {
						rows, e := csv.NewReader(bytes.NewReader(raw)).ReadAll()
						if e != nil {
							t.Fatal(e)
						}
						matched := false
						for _, table := range found {
							if reflect.DeepEqual(rows, table.Rows) {
								matched = true
							}
						}
						if !matched {
							t.Fatal("actual CSV diagnostic readback", path)
						}
					}
				}
				xlsx := filepath.Join(dir, "assessment.xlsx")
				if e = excelreport.WriteFile(data, xlsx, opts); e != nil {
					t.Fatal(e)
				}
				book, e := excelize.OpenFile(xlsx)
				if e != nil {
					t.Fatal(e)
				}
				for sheet, table := range found {
					rows, e := book.GetRows(sheet)
					header := -1
					for i, row := range rows {
						if reflect.DeepEqual(row, table.Rows[0]) {
							header = i
							break
						}
					}
					if e != nil || header < 2 || len(rows[0]) == 0 || rows[0][0] == "" || len(rows[1]) == 0 || rows[1][0] != sheet || !reflect.DeepEqual(rows[header:], table.Rows) {
						book.Close()
						t.Fatal("actual XLSX diagnostic readback", sheet, e, rows, table.Rows)
					}
				}
				rows, e := book.GetRows("Assessment Status")
				if e != nil {
					book.Close()
					t.Fatal(e)
				}
				joined := ""
				for _, row := range rows {
					joined += strings.Join(row, "|") + "\n"
				}
				sqlReportCheckMask(t, joined, masked)
				if e = book.Close(); e != nil {
					t.Fatal(e)
				}
				after, e := json.Marshal(data)
				if e != nil || !bytes.Equal(after, snapshot) {
					t.Fatal("renderers mutated canonical report", e)
				}
			})
		}
	}
}

func TestSQLReportScopeRedaction(t *testing.T) {
	data := sqlReportRendererFixture(t, "known")
	projected, e := tables.Build(data, tables.Options{RedactSubscriptionIDs: true})
	if e != nil {
		t.Fatal(e)
	}
	for _, table := range projected {
		if table.Key == "assessmentStatus" {
			for _, row := range table.Rows {
				for _, cell := range row {
					if strings.Contains(strings.ToLower(cell), sqlReportOracleID) {
						t.Error("shared assessment status scope ID leaks a subscription under redaction")
					}
				}
			}
		}
	}
	dir := t.TempDir()
	paths, e := csvreport.Write(data, filepath.Join(dir, "scope"), tables.Options{RedactSubscriptionIDs: true})
	if e != nil {
		t.Fatal(e)
	}
	for _, path := range paths {
		if strings.Contains(filepath.Base(path), "assessmentStatus") {
			raw, e := os.ReadFile(path)
			if e != nil {
				t.Fatal(e)
			}
			if strings.Contains(strings.ToLower(string(raw)), sqlReportOracleID) {
				t.Error("actual CSV assessment status scope ID leaks a subscription under redaction")
			}
		}
	}
	output := filepath.Join(dir, "scope.xlsx")
	if e = excelreport.WriteFile(data, output, tables.Options{RedactSubscriptionIDs: true}); e != nil {
		t.Fatal(e)
	}
	book, e := excelize.OpenFile(output)
	if e != nil {
		t.Fatal(e)
	}
	rows, e := book.GetRows("Assessment Status")
	if e != nil {
		book.Close()
		t.Fatal(e)
	}
	for _, row := range rows {
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), sqlReportOracleID) {
				t.Error("actual XLSX assessment status scope ID leaks a subscription under redaction")
			}
		}
	}
	if e = book.Close(); e != nil {
		t.Fatal(e)
	}
}
