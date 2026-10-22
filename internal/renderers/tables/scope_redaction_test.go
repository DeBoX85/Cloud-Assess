package tables

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

const scopeRedactionID = "aaaaaaaa-1111-4111-8111-aaaaaaaaaaaa"
const scopeUnrelatedID = "bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"

func scopeRedactionInput(scopeID string) result.Input {
	return result.Input{ScopeID: scopeID, Completeness: assessment.CompletenessComplete,
		Scope: &assessment.ScopeResolution{Selection: "explicit", Status: "resolved", RequestedSubscriptionIDs: []string{scopeRedactionID}, IncludedSubscriptionIDs: []string{scopeRedactionID}, ResolvedSubscriptions: []assessment.ScopeSubscription{{SubscriptionID: scopeRedactionID, SubscriptionName: "selected"}}}}
}
func scopeRedactionProbe() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: "1.0", ID: "scope-probe", Metadata: assessment.PluginMetadata{Name: "scope-probe", Version: "1", Description: "Scope redaction probe", Author: "test", License: "MIT", Type: "internal"}, SheetName: "Scope Probe", Description: "Scope redaction probe", Columns: []string{"Subscription ID"}, Rows: []assessment.PluginRow{{SubscriptionID: scopeRedactionID, Cells: []string{scopeRedactionID}}}, Health: assessment.StageExecution{Name: "scope-probe", Status: assessment.StageCompleted, Records: 1}}
}
func scopeStatusRows(t *testing.T, data *result.AssessmentResult, masked bool) [][]string {
	t.Helper()
	projected, e := Build(data, Options{RedactSubscriptionIDs: masked})
	if e != nil {
		t.Fatal(e)
	}
	for _, table := range projected {
		if table.Key == "assessmentStatus" {
			return table.Rows
		}
	}
	t.Fatal("missing assessment status")
	return nil
}
func TestAssessmentStatusScopeRedaction(t *testing.T) {
	cases := []struct {
		name, scope string
		maskKnown   bool
	}{
		{"path", "/subscriptions/" + scopeRedactionID + "/providers/Microsoft.Sql/locations/eastus", true},
		{"case", "/SUBSCRIPTIONS/" + strings.ToUpper(scopeRedactionID) + "/resourceGroups/test", true},
		{"bare-known", scopeRedactionID, true},
		{"embedded-known", "catalog:" + scopeRedactionID + ":eastus", true},
		{"unrelated-uuid-retained", "tenant:" + scopeUnrelatedID + " subscription " + scopeRedactionID, true},
		{"unrelated-only", scopeUnrelatedID, false},
	}
	for _, withPlugin := range []bool{false, true} {
		for _, c := range cases {
			t.Run(fmt.Sprintf("%s-plugin-%t", c.name, withPlugin), func(t *testing.T) {
				input := scopeRedactionInput(c.scope)
				var data *result.AssessmentResult
				var e error
				if withPlugin {
					data, e = result.BuildWithPluginTables(input, []assessment.PluginTable{scopeRedactionProbe()})
				} else {
					data = result.Build(input)
				}
				if e != nil {
					t.Fatal(e)
				}
				before, e := json.Marshal(data)
				if e != nil {
					t.Fatal(e)
				}
				for _, masked := range []bool{false, true} {
					rows := scopeStatusRows(t, data, masked)
					if len(rows) < 2 {
						t.Fatal("missing status data")
					}
					for _, row := range rows[1:] {
						if len(row) != 14 {
							t.Fatal("status schema changed", row)
						}
						got := row[2]
						if !masked || !c.maskKnown {
							if got != c.scope {
								t.Fatal("unredacted or unrelated scope changed", got, c.scope)
							}
						} else {
							if strings.Contains(strings.ToLower(got), scopeRedactionID) {
								t.Fatal("assessment status scope ID leaks a known subscription when redaction is enabled")
							}
							if strings.Contains(c.scope, scopeUnrelatedID) && !strings.Contains(got, scopeUnrelatedID) {
								t.Fatal("unrelated tenant UUID must be retained")
							}
							if !strings.Contains(got, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxx") {
								t.Fatal("existing subscription mask format")
							}
						}
					}
				}
				after, e := json.Marshal(data)
				if e != nil || !reflect.DeepEqual(before, after) {
					t.Fatal("scope rendering mutated canonical result", e)
				}
			})
		}
	}
}
func TestAssessmentStatusScopeContextWithoutRecords(t *testing.T) {
	input := result.Input{ScopeID: "/SUBSCRIPTIONS/" + strings.ToUpper(scopeRedactionID) + "/locations/eastus", Completeness: assessment.CompletenessComplete}
	data := result.Build(input)
	rows := scopeStatusRows(t, data, true)
	if strings.Contains(strings.ToLower(rows[1][2]), scopeRedactionID) {
		t.Fatal("contextual scope subscription must mask even without discovery records")
	}
}
func TestAssessmentStatusScopeLargeCoreSelection(t *testing.T) {
	input := scopeRedactionInput(scopeRedactionID)
	for i := 0; i < 5000; i++ {
		id := fmt.Sprintf("%08x-3333-4333-8333-cccccccccccc", i)
		input.Scope.RequestedSubscriptionIDs = append(input.Scope.RequestedSubscriptionIDs, id)
		input.Scope.IncludedSubscriptionIDs = append(input.Scope.IncludedSubscriptionIDs, id)
		input.Scope.ResolvedSubscriptions = append(input.Scope.ResolvedSubscriptions, assessment.ScopeSubscription{SubscriptionID: id, SubscriptionName: "selected"})
	}
	data := result.Build(input)
	rows := scopeStatusRows(t, data, true)
	if strings.Contains(strings.ToLower(rows[1][2]), scopeRedactionID) {
		t.Fatal("large core scope still needs masking without adding the plugin4096ID limit")
	}
}
