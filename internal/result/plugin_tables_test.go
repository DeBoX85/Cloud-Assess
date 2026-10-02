package result

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func pluginTableFixture() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: "1.0", ID: "zones", Metadata: assessment.PluginMetadata{Name: "zone-mapping", Version: "1.0.0", Description: "Source description", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}, SheetName: "Zone Mapping", Description: "Logical-to-physical mappings", Columns: []string{"Subscription", "Location", "Display Name", "Logical Zone", "Physical Zone"}, Rows: []assessment.PluginRow{{SubscriptionID: "11111111-1111-4111-8111-111111111111", Cells: []string{"A", "eastus", "East US", "1", "eastus-az1"}}}, Health: assessment.StageExecution{Name: "zone-mapping", Status: assessment.StageCompleted, Records: 1}}
}

func TestPluginCanonicalBuilderOwnsDataAndPreservesDefaults(t *testing.T) {
	input := Input{Completeness: assessment.CompletenessComplete}
	ordinary := Build(input)
	empty, err := BuildWithPluginTables(input, nil)
	if err != nil || !reflect.DeepEqual(ordinary, empty) {
		t.Fatal("empty extension changed ordinary builder")
	}
	b, _ := json.Marshal(empty)
	if strings.Contains(string(b), "pluginTables") || empty.SchemaVersion != "1.0" {
		t.Fatal("default schema changed")
	}
	original := pluginTableFixture()
	original.Health.Status = assessment.StageFailed
	original.Health.Error = &assessment.AssessmentError{Code: "zone_request_failed", Message: "sanitized"}
	original.Health.Warnings = []assessment.AssessmentWarning{{Code: "zone_partial", Message: "earlier rows retained"}}
	tables := []assessment.PluginTable{original}
	first, err := BuildWithPluginTables(input, tables)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildWithPluginTables(input, tables)
	if err != nil {
		t.Fatal(err)
	}
	if first.SchemaVersion != "1.1" || first.Completeness != assessment.CompletenessPartial {
		t.Fatal("partial plugin did not change schema/completeness")
	}
	first.PluginTables[0].Columns[0] = "changed"
	first.PluginTables[0].Rows[0].Cells[0] = "changed"
	first.PluginTables[0].Health.Warnings[0].Message = "changed"
	first.PluginTables[0].Health.Error.Message = "changed"
	if !reflect.DeepEqual(second.PluginTables[0], tables[0]) {
		t.Fatal("result ownership leaked")
	}
	tables[0].Rows[0].Cells[0] = "caller change"
	if second.PluginTables[0].Rows[0].Cells[0] != "A" {
		t.Fatal("caller mutation leaked")
	}
	second.Completeness = assessment.CompletenessComplete
	if second.ValidatePluginExtension() == nil {
		t.Fatal("complete assessment concealed failed table")
	}
}

func TestPluginCanonicalHealthAndOrder(t *testing.T) {
	for _, status := range []assessment.StageStatus{assessment.StageCompleted, assessment.StageCompletedWithWarnings, assessment.StageFailed, assessment.StageSkipped} {
		t.Run(string(status), func(t *testing.T) {
			p := pluginTableFixture()
			p.Health.Status = status
			expected := assessment.CompletenessComplete
			if status == assessment.StageCompletedWithWarnings {
				p.Health.Warnings = []assessment.AssessmentWarning{{Code: "zone_partial", Message: "warning"}}
				expected = assessment.CompletenessCompleteWithWarnings
			}
			if status == assessment.StageFailed {
				p.Health.Error = &assessment.AssessmentError{Code: "zone_request_failed", Message: "safe"}
				expected = assessment.CompletenessPartial
			}
			if status == assessment.StageSkipped {
				p.Rows = nil
				p.Health.Records = 0
				expected = assessment.CompletenessPartial
			}
			got, e := BuildWithPluginTables(Input{Completeness: assessment.CompletenessComplete}, []assessment.PluginTable{p})
			if e != nil || got.Completeness != expected {
				t.Fatalf("health: %+v %v", got, e)
			}
			got, e = BuildWithPluginTables(Input{Completeness: assessment.CompletenessFailed}, []assessment.PluginTable{p})
			if e != nil || got.Completeness != assessment.CompletenessFailed {
				t.Fatal("existing critical failure was concealed")
			}
		})
	}
	zone := pluginTableFixture()
	health := pluginTableFixture()
	health.Metadata.Name = "service-health"
	health.SheetName = "Service Health"
	health.ID = "availability"
	health.Health.Name = "service-health"
	got, e := BuildWithPluginTables(Input{Completeness: assessment.CompletenessComplete}, []assessment.PluginTable{zone, health})
	if e != nil || got.PluginTables[0].Metadata.Name != "service-health" || got.PluginTables[1].Metadata.Name != "zone-mapping" {
		t.Fatal("canonical order")
	}
	short := pluginTableFixture()
	short.Metadata.Name = "zone"
	short.Health.Name = "zone"
	short.SheetName = "Zones"
	got, e = BuildWithPluginTables(Input{Completeness: assessment.CompletenessComplete}, []assessment.PluginTable{zone, short})
	if e != nil || got.PluginTables[0].Metadata.Name != "zone" {
		t.Fatal("table key separator changed canonical plugin-name order")
	}
}

func TestPluginTableValidationRejectsMalformedAndBoundedInputs(t *testing.T) {
	cases := map[string]func(*assessment.PluginTable){
		"Unicode reserved sheet": func(p *assessment.PluginTable) { p.SheetName = "Coſts" },
		"UTF16 sheet ceiling":    func(p *assessment.PluginTable) { p.SheetName = strings.Repeat("😀", 16) },
		"schema":                 func(p *assessment.PluginTable) { p.SchemaVersion = "2.0" }, "id path": func(p *assessment.PluginTable) { p.ID = "../zones" }, "plugin type": func(p *assessment.PluginTable) { p.Metadata.Type = "yaml" }, "empty version": func(p *assessment.PluginTable) { p.Metadata.Version = "" },
		"reserved sheet": func(p *assessment.PluginTable) { p.SheetName = "inventory" }, "invalid sheet": func(p *assessment.PluginTable) { p.SheetName = "a/b" }, "long sheet": func(p *assessment.PluginTable) { p.SheetName = strings.Repeat("x", 32) }, "quoted sheet": func(p *assessment.PluginTable) { p.SheetName = "'sheet" },
		"row width": func(p *assessment.PluginTable) { p.Rows[0].Cells = p.Rows[0].Cells[:4] }, "subscription": func(p *assessment.PluginTable) { p.Rows[0].SubscriptionID = "bad" }, "duplicate columns": func(p *assessment.PluginTable) { p.Columns[1] = "subscription" }, "empty columns": func(p *assessment.PluginTable) { p.Columns = nil }, "column count": func(p *assessment.PluginTable) { p.Columns = make([]string, 65) },
		"UTF8": func(p *assessment.PluginTable) { p.Rows[0].Cells[0] = "\xff" }, "XML control": func(p *assessment.PluginTable) { p.Rows[0].Cells[0] = "\x01" }, "cell ceiling": func(p *assessment.PluginTable) { p.Rows[0].Cells[0] = strings.Repeat("x", 32768) }, "UTF16 ceiling": func(p *assessment.PluginTable) { p.Rows[0].Cells[0] = strings.Repeat("😀", 16384) },
		"unknown health": func(p *assessment.PluginTable) { p.Health.Status = "maybe" }, "row count": func(p *assessment.PluginTable) { p.Health.Records = 2 }, "missing failure": func(p *assessment.PluginTable) { p.Health.Status = assessment.StageFailed }, "missing warning": func(p *assessment.PluginTable) { p.Health.Status = assessment.StageCompletedWithWarnings }, "skipped data": func(p *assessment.PluginTable) { p.Health.Status = assessment.StageSkipped }, "invalid warning": func(p *assessment.PluginTable) {
			p.Health.Status = assessment.StageCompletedWithWarnings
			p.Health.Warnings = []assessment.AssessmentWarning{{Code: "bad code"}}
		},
		"row ceiling": func(p *assessment.PluginTable) {
			p.Rows = make([]assessment.PluginRow, 65537)
			p.Health.Records = len(p.Rows)
		}, "text ceiling": func(p *assessment.PluginTable) {
			p.Rows = nil
			for i := 0; i < 600; i++ {
				p.Rows = append(p.Rows, assessment.PluginRow{Cells: []string{strings.Repeat("x", 32767), "", "", "", ""}})
			}
			p.Health.Records = len(p.Rows)
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			p := pluginTableFixture()
			change(&p)
			if _, e := BuildWithPluginTables(Input{Completeness: assessment.CompletenessComplete}, []assessment.PluginTable{p}); e == nil {
				t.Fatal("invalid plugin table accepted")
			}
		})
	}
	p := pluginTableFixture()
	if assessment.ValidatePluginTables([]assessment.PluginTable{p, p}) == nil {
		t.Fatal("duplicate table admitted")
	}
	q := pluginTableFixture()
	q.ID = "other"
	q.SheetName = "zone mapping"
	if assessment.ValidatePluginTables([]assessment.PluginTable{p, q}) == nil {
		t.Fatal("sheet collision admitted")
	}
	p.SheetName = "Sheet"
	q.SheetName = "ſheet"
	if assessment.ValidatePluginTables([]assessment.PluginTable{p, q}) == nil {
		t.Fatal("cross-plugin Unicode fold sheet collision admitted")
	}
	p = pluginTableFixture()
	p.Columns[0] = "Subscription"
	p.Columns[1] = "ſubscription"
	if assessment.ValidatePluginTables([]assessment.PluginTable{p}) == nil {
		t.Fatal("Unicode fold duplicate column admitted")
	}
	if assessment.ValidatePluginTables(make([]assessment.PluginTable, 65)) == nil {
		t.Fatal("table count ceiling ignored")
	}
	batch := make([]assessment.PluginTable, 5)
	for i := range batch {
		p := pluginTableFixture()
		p.ID = fmt.Sprintf("part-%c", 'a'+i)
		p.SheetName = fmt.Sprintf("Zones %d", i)
		p.Columns = []string{"Value"}
		n := 65536
		if i == 4 {
			n = 1
		}
		p.Rows = make([]assessment.PluginRow, n)
		for j := range p.Rows {
			p.Rows[j].Cells = []string{"x"}
		}
		p.Health.Records = n
		batch[i] = p
	}
	if assessment.ValidatePluginTables(batch) == nil {
		t.Fatal("literal aggregate row ceiling ignored")
	}
	p = pluginTableFixture()
	p.SheetName = strings.Repeat("😀", 15) + "x"
	p.Rows[0].Cells[0] = strings.Repeat("😀", 16383) + "x"
	if assessment.ValidatePluginTables([]assessment.PluginTable{p}) != nil {
		t.Fatal("exact UTF16 ceiling rejected")
	}
}
