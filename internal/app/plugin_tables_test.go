package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	csvrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/csv"
	excelrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/excel"
	jsonrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/json"
	sarifrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/sarif"
	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/xuri/excelize/v2"
)

const pluginSubID = "abcdefab-abcd-4abc-8abc-abcdefabcdef"

// These two literal schemas reflect different pinned source plugins. The service
// health fixture checks infrastructure shape, not its unmigrated query/adapter.
func pluginReportFixture(t *testing.T) *result.AssessmentResult {
	t.Helper()
	zone := assessment.PluginTable{SchemaVersion: "1.0", ID: "zones", Metadata: assessment.PluginMetadata{Name: "zone-mapping", Version: "1.0.0", Description: "Retrieves logical-to-physical availability zone mappings for all Azure regions in each subscription", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}, SheetName: "Zone Mapping", Description: "Logical-to-physical availability zone mappings for Azure regions by subscription", Columns: []string{"Subscription", "Location", "Display Name", "Logical Zone", "Physical Zone"}, Rows: []assessment.PluginRow{{SubscriptionID: pluginSubID, Cells: []string{"=1+1", "eastus", "East US", "1", "eastus-az1"}}}, Health: assessment.StageExecution{Name: "zone-mapping", Status: assessment.StageCompleted, Records: 1}}
	health := assessment.PluginTable{SchemaVersion: "1.0", ID: "availability", Metadata: assessment.PluginMetadata{Name: "service-health", Version: "0.1.0-beta", Description: "Analyzes Azure service health events to determine the percentage of time resources were unaffected by service issues over the last 90 days.", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}, SheetName: "Service Health", Description: "Service health fixture", Columns: []string{"Subscription ID", "Target Region", "Target Resource Type", "Percentage Without Events", "Events Count", "Affected Resources"}, Rows: []assessment.PluginRow{{SubscriptionID: pluginSubID, Cells: []string{pluginSubID, "norwayeast", "microsoft.compute/virtualmachines", "99.5", "2", "4"}}}, Health: assessment.StageExecution{Name: "service-health", Status: assessment.StageFailed, Records: 1, Error: &assessment.AssessmentError{Code: "plugin_request_failed", Message: "subscription id " + strings.ToUpper(pluginSubID) + " had incomplete data"}}}
	r, e := result.BuildWithPluginTables(result.Input{Completeness: assessment.CompletenessComplete}, []assessment.PluginTable{zone, health})
	if e != nil {
		t.Fatal(e)
	}
	return r
}

func TestPluginTablesApplicationArtifactsPrivacyAndPartialExit(t *testing.T) {
	for _, mask := range []bool{false, true} {
		t.Run(map[bool]string{false: "raw", true: "masked"}[mask], func(t *testing.T) {
			r := pluginReportFixture(t)
			original, _ := json.Marshal(r)
			var stdout bytes.Buffer
			runner := NewRunner(&fakeAssessmentRunner{result: r})
			runner.stdout = &stdout
			base := filepath.Join(t.TempDir(), "plugin-report")
			out, e := runner.Run(context.Background(), ScanOptions{Outputs: OutputOptions{BaseName: base, XLSX: true, JSON: true, CSV: true, Stdout: true, SARIF: true, RedactSubscriptionIDs: mask, Version: "test"}})
			var partial *PartialAssessmentError
			if out.ExitCode != ExitPartial || !errors.As(e, &partial) {
				t.Fatalf("plugin partial exit: %d %v", out.ExitCode, e)
			}
			onDisk, e := os.ReadFile(base + ".json")
			if e != nil {
				t.Fatal(e)
			}
			if !bytes.Equal(bytes.TrimSpace(onDisk), bytes.TrimSpace(stdout.Bytes())) {
				t.Fatal("JSON/stdout table data differ")
			}
			var decoded result.AssessmentResult
			if e = json.Unmarshal(onDisk, &decoded); e != nil {
				t.Fatal(e)
			}
			if decoded.SchemaVersion != "1.1" || len(decoded.PluginTables) != 2 || decoded.Completeness != assessment.CompletenessPartial {
				t.Fatal("canonical table extension or health missing")
			}
			expectedID := pluginSubID
			if mask {
				expectedID = redact.SubscriptionID(pluginSubID, true)
			}
			if decoded.PluginTables[0].Rows[0].SubscriptionID != expectedID || decoded.PluginTables[0].Rows[0].Cells[0] != expectedID || decoded.PluginTables[1].Rows[0].Cells[0] != "=1+1" {
				t.Fatal("canonical values/privacy")
			}
			var expectedTables []assessment.PluginTable
			expectedBytes, _ := json.Marshal(r.PluginTables)
			if e = json.Unmarshal(expectedBytes, &expectedTables); e != nil {
				t.Fatal(e)
			}
			if mask {
				for i := range expectedTables {
					for j := range expectedTables[i].Rows {
						expectedTables[i].Rows[j].SubscriptionID = expectedID
						for k, cell := range expectedTables[i].Rows[j].Cells {
							expectedTables[i].Rows[j].Cells[k] = strings.ReplaceAll(cell, pluginSubID, expectedID)
						}
					}
					if expectedTables[i].Health.Error != nil {
						expectedTables[i].Health.Error.Message = strings.ReplaceAll(expectedTables[i].Health.Error.Message, strings.ToUpper(pluginSubID), expectedID)
					}
				}
			}
			if !reflect.DeepEqual(decoded.PluginTables, expectedTables) {
				t.Fatal("complete plugin metadata/columns/rows/health JSON comparison differs")
			}
			if mask && strings.Contains(strings.ToLower(string(onDisk)), pluginSubID) {
				t.Fatal("plugin JSON subscription leak")
			}
			book, e := excelize.OpenFile(base + ".xlsx")
			if e != nil {
				t.Fatal(e)
			}
			defer book.Close()
			if !reflect.DeepEqual(book.GetSheetList(), []string{"Assessment Status", "Service Health", "Zone Mapping"}) {
				t.Fatalf("actual sheets: %+v", book.GetSheetList())
			}
			for _, table := range r.PluginTables {
				path := base + "." + table.Key() + ".csv"
				file, e := os.Open(path)
				if e != nil {
					t.Fatal(e)
				}
				records, e := csv.NewReader(file).ReadAll()
				file.Close()
				if e != nil {
					t.Fatal(e)
				}
				want := append([][]string{append([]string(nil), table.Columns...)}, append([]string(nil), table.Rows[0].Cells...))
				if mask {
					for i, row := range want {
						for j, cell := range row {
							want[i][j] = strings.ReplaceAll(cell, pluginSubID, expectedID)
						}
					}
				}
				csvWant := [][]string{want[0], append([]string(nil), want[1]...)}
				if table.ID == "zones" {
					csvWant[1][0] = "'=1+1"
				}
				if !reflect.DeepEqual(records, csvWant) {
					t.Fatalf("all plugin CSV cells: got=%+v want=%+v", records, csvWant)
				}
				for i, cell := range want[0] {
					address, _ := excelize.CoordinatesToCellName(i+1, 4)
					actual, _ := book.GetCellValue(table.SheetName, address)
					if actual != cell {
						t.Fatalf("XLSX header %s %q", address, actual)
					}
				}
				for i, cell := range want[1] {
					address, _ := excelize.CoordinatesToCellName(i+1, 5)
					actual, _ := book.GetCellValue(table.SheetName, address)
					formula, _ := book.GetCellFormula(table.SheetName, address)
					if actual != cell || formula != "" {
						t.Fatalf("XLSX literal cell %s actual=%q formula=%q want=%q", address, actual, formula, cell)
					}
				}
			}
			status, e := os.ReadFile(base + ".assessmentStatus.csv")
			if e != nil {
				t.Fatal(e)
			}
			if !bytes.Contains(status, []byte("plugin:service-health:availability")) || !bytes.Contains(status, []byte("plugin_request_failed")) {
				t.Fatal("plugin failure health absent")
			}
			if mask && strings.Contains(strings.ToLower(string(status)), pluginSubID) {
				t.Fatal("plugin human health subscription leak")
			}
			// Plugin tables do not become invented SARIF rule findings.
			sarif, e := os.ReadFile(base + ".sarif")
			if e != nil {
				t.Fatal(e)
			}
			if bytes.Contains(sarif, []byte("zone-mapping")) || bytes.Contains(sarif, []byte("service-health")) {
				t.Fatal("plugin tables became SARIF findings")
			}
			after, _ := json.Marshal(r)
			if !bytes.Equal(original, after) {
				t.Fatal("renderer mutated canonical plugin data")
			}
		})
	}
}

func TestPluginTextMaskingInferenceAndCallerOwnership(t *testing.T) {
	r := pluginReportFixture(t)
	unknown := "33333333-3333-4333-8333-333333333333"
	r.PluginTables[1].Rows[0].Cells[0] = "subscription id " + strings.ToUpper(pluginSubID) + " /subscriptions/" + unknown + "/resourceGroups/r"
	projected, e := tables.Build(r, tables.Options{RedactSubscriptionIDs: true})
	if e != nil {
		t.Fatal(e)
	}
	for _, table := range projected {
		if table.Key == "plugin_zone-mapping_zones" {
			if strings.Contains(strings.ToLower(table.Rows[1][0]), pluginSubID) || strings.Contains(table.Rows[1][0], unknown) {
				t.Fatal("known or contextual plugin identity leaked")
			}
			table.Rows[0][0] = "caller"
			table.Rows[1][0] = "caller"
		}
	}
	if r.PluginTables[1].Columns[0] != "Subscription" || strings.HasPrefix(r.PluginTables[1].Rows[0].Cells[0], "caller") {
		t.Fatal("projection borrowed cells")
	}
}

func TestInvalidPluginTablesRejectBeforeOutputReplacement(t *testing.T) {
	r := pluginReportFixture(t)
	r.PluginTables[1].SheetName = "../unsafe"
	dir := t.TempDir()
	xlsx := filepath.Join(dir, "existing.xlsx")
	js := filepath.Join(dir, "existing.json")
	csv := filepath.Join(dir, "existing.assessmentStatus.csv")
	for _, p := range []string{xlsx, js, csv} {
		if e := os.WriteFile(p, []byte("preserve"), 0600); e != nil {
			t.Fatal(e)
		}
	}
	if excelrenderer.WriteFile(r, xlsx, tables.Options{}) == nil {
		t.Fatal("invalid Excel table accepted")
	}
	if jsonrenderer.WriteFileWithOptions(r, js, jsonrenderer.Options{}) == nil {
		t.Fatal("invalid JSON table accepted")
	}
	if _, e := csvrenderer.Write(r, filepath.Join(dir, "existing"), tables.Options{}); e == nil {
		t.Fatal("invalid CSV table accepted")
	}
	if _, e := sarifrenderer.Marshal(r, "test"); e == nil {
		t.Fatal("invalid table result accepted for SARIF")
	}
	for _, p := range []string{xlsx, js, csv} {
		b, e := os.ReadFile(p)
		if e != nil || string(b) != "preserve" {
			t.Fatal("validation replaced existing output")
		}
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 3 {
		t.Fatal("validation wrote extra files")
	}
	r = pluginReportFixture(t)
	r.Completeness = assessment.CompletenessComplete
	if _, e := jsonrenderer.Marshal(r); e == nil {
		t.Fatal("false complete plugin report accepted")
	}
}

func TestFailedAndSkippedEmptyPluginTablesStillRender(t *testing.T) {
	for _, status := range []assessment.StageStatus{assessment.StageFailed, assessment.StageSkipped} {
		r := pluginReportFixture(t)
		p := r.PluginTables[0]
		p.Rows = nil
		p.Health.Records = 0
		p.Health.Status = status
		if status == assessment.StageSkipped {
			p.Health.Error = nil
		}
		r, e := result.BuildWithPluginTables(result.Input{Completeness: assessment.CompletenessComplete}, []assessment.PluginTable{p})
		if e != nil {
			t.Fatal(e)
		}
		projected, e := tables.Build(r, tables.Options{})
		if e != nil {
			t.Fatal(e)
		}
		table := projected[len(projected)-1]
		if len(table.Rows) != 1 || !tables.ShouldRender(r, table) || !reflect.DeepEqual(table.Rows[0], p.Columns) {
			t.Fatal("requested empty failed/skipped table disappeared")
		}
	}
}

func TestPluginUnicodeSheetBoundaryActualWorkbook(t *testing.T) {
	r := pluginReportFixture(t)
	name := strings.Repeat("😀", 15) + "x"
	r.PluginTables[1].SheetName = name
	path := filepath.Join(t.TempDir(), "unicode.xlsx")
	if e := excelrenderer.WriteFile(r, path, tables.Options{}); e != nil {
		t.Fatal(e)
	}
	book, e := excelize.OpenFile(path)
	if e != nil {
		t.Fatal(e)
	}
	defer book.Close()
	cell, e := book.GetCellValue(name, "A5")
	if e != nil || cell != "=1+1" {
		t.Fatal("exact 31 UTF16-unit sheet failed actual workbook roundtrip")
	}
}

func TestPluginHumanMaskingBoundFailsClosedBeforeReplacement(t *testing.T) {
	r := pluginReportFixture(t)
	r.Scope = &assessment.ScopeResolution{}
	for i := 0; i < 4097; i++ {
		r.Scope.RequestedSubscriptionIDs = append(r.Scope.RequestedSubscriptionIDs, fmt.Sprintf("%08x-1111-4111-8111-111111111111", i))
	}
	path := filepath.Join(t.TempDir(), "preserve.xlsx")
	if e := os.WriteFile(path, []byte("preserve"), 0600); e != nil {
		t.Fatal(e)
	}
	if e := excelrenderer.WriteFile(r, path, tables.Options{RedactSubscriptionIDs: true}); e == nil {
		t.Fatal("excess masking identities accepted")
	}
	b, e := os.ReadFile(path)
	if e != nil || string(b) != "preserve" {
		t.Fatal("masking failure replaced output")
	}
	if _, e := tables.Build(r, tables.Options{}); e != nil {
		t.Fatal("explicit raw projection unexpectedly used masking limit")
	}
	ordinary := result.Build(result.Input{Scope: r.Scope, Completeness: assessment.CompletenessComplete})
	if _, e := tables.Build(ordinary, tables.Options{RedactSubscriptionIDs: true}); e != nil {
		t.Fatal("new plugin mask bound changed ordinary projection")
	}
}
