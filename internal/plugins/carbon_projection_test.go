package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/carbon"
)

func healthyCarbon(t *testing.T) assessment.PluginTable {
	t.Helper()
	a, b, c := 80.0, 100.0, -20.0
	v, e := carbon.Project(context.Background(), "2026-01-01", "2026-03-31", []carbon.Item{{ResourceType: "Microsoft.Compute/virtualMachines", Latest: &a, Previous: &b, Change: &c}}, nil)
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func TestCarbonProjectionOwnershipSourceAndInvalidOutput(t *testing.T) {
	value := healthyCarbon(t)
	projected := CarbonTable(value, nil, map[string]string{"11111111-1111-4111-8111-111111111111": "Dev"}, time.Unix(1, 0), time.Unix(2, 0))
	if !reflect.DeepEqual(projected.Rows, value.Rows) || projected.Health.Status != assessment.StageCompleted || projected.Health.StartedAt != time.Unix(1, 0).UTC() {
		t.Fatal("source table/projection")
	}
	projected.Columns[0] = "mutated"
	projected.Rows[0].Cells[3] = "mutated"
	if value.Columns[0] != "Period From" || value.Rows[0].Cells[3] != "80.00" {
		t.Fatal("caller output aliased")
	}
	for _, mode := range []string{"metadata", "correlation", "headers", "sheet", "period", "date", "unit", "number", "ratio", "duplicate", "empty-scope", "warning-limit"} {
		t.Run(mode, func(t *testing.T) {
			v := healthyCarbon(t)
			scope := map[string]string{"11111111-1111-4111-8111-111111111111": "Dev"}
			switch mode {
			case "metadata":
				v.Metadata.Description = "unreviewed"
			case "correlation":
				v.Rows[0].SubscriptionID = "11111111-1111-4111-8111-111111111111"
			case "headers":
				v.Columns[3] = "unreviewed"
			case "sheet":
				v.SheetName = "Other"
			case "period":
				v.Rows[0].Cells[1] = "2026-01-01"
			case "date":
				v.Rows[0].Cells[0] = "not-a-date"
			case "unit":
				v.Rows[0].Cells[7] = "tons"
			case "number":
				v.Rows[0].Cells[3] = "NaN"
			case "ratio":
				v.Rows[0].Cells[5] = "-20"
			case "duplicate":
				v.Rows = append(v.Rows, v.Rows[0])
				v.Health.Records++
			case "empty-scope":
				scope = nil
			case "warning-limit":
				v.Health.Status = assessment.StageCompletedWithWarnings
				v.Health.Warnings = make([]assessment.AssessmentWarning, 5)
			}
			got := CarbonTable(v, nil, scope, time.Time{}, time.Time{})
			if got.Health.Error == nil || got.Health.Error.Code != "carbon_output_invalid" || len(got.Rows) != 0 {
				t.Fatal("unsafe aggregate contract accepted")
			}
		})
	}
	// Positive source values can round to display zero. Do not treat the
	// formatted cell as an exact positive-value predicate.
	a, b := 0.0002, 0.0001
	v, e := carbon.Project(context.Background(), "2026-01-01", "2026-03-31", []carbon.Item{{ResourceType: "Microsoft.Compute/virtualMachines", Latest: &a, Previous: &b}}, nil)
	if e != nil || v.Rows[0].Cells[4] != "0.00" {
		t.Fatal("tiny positive source fixture")
	}
	got := CarbonTable(v, nil, map[string]string{"scope": "fixture"}, time.Time{}, time.Time{})
	if got.Health.Status != assessment.StageCompleted || got.Rows[0].Cells[4] != "0.00" {
		t.Fatal("valid source rounding rejected")
	}
}

func TestCarbonProjectionPrivacyWarningsFailuresAndRegistryOrder(t *testing.T) {
	for _, mode := range []string{"warning", "failed", "pending"} {
		v := healthyCarbon(t)
		var err error
		v.Health.Warnings = []assessment.AssessmentWarning{{Code: "raw_provider_warning", Message: "private-provider-canary"}}
		v.Health.Status = assessment.StageCompletedWithWarnings
		if mode == "failed" {
			v.Health.Status = assessment.StageFailed
			v.Health.Error = &assessment.AssessmentError{Code: "raw_provider_error", Message: "private-provider-canary"}
			err = errors.New("private-operation-canary")
		}
		if mode == "pending" {
			v = carbon.PendingTable()
			err = errors.New("private-operation-canary")
		}
		got := CarbonTable(v, err, map[string]string{"scope": "fixture"}, time.Time{}, time.Time{})
		b, _ := json.Marshal(got)
		if strings.Contains(string(b), "canary") || assessment.ValidatePluginTables([]assessment.PluginTable{got}) != nil || mode == "failed" && len(got.Rows) != 1 || mode == "pending" && got.Health.Status != assessment.StageFailed {
			t.Fatal("unsafe/false health or lost rows")
		}
	}
	names, e := ValidateNames([]string{"zone-mapping", "carbon-emissions", "sql-eol", "service-health", "carbon-emissions"})
	if e != nil || strings.Join(names, ",") != "carbon-emissions,service-health,sql-eol,zone-mapping" {
		t.Fatal("source four-adapter order/dedup")
	}
	if PendingTable(carbon.Name).ID != "emissions" || PendingZoneTable().Metadata.Name != ZoneMapping {
		t.Fatal("registry metadata alias")
	}
}
