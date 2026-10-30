package region

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const admissionBudgetSub = "aaaaaaaa-1111-4111-8111-111111111111"

func admissionBudgetResource(typ, name, sku, kind string) assessment.Resource {
	return assessment.Resource{ID: "/subscriptions/" + admissionBudgetSub + "/resourceGroups/rg/providers/" + typ + "/" + name, SubscriptionID: admissionBudgetSub, ResourceGroup: "rg", Name: name, Type: typ, Location: "eastus", SKUName: sku, Kind: kind}
}
func admissionBudgetPlan(t *testing.T, rows []assessment.Resource) *RequiredInventoryPlan {
	t.Helper()
	p, e := BuildRequiredInventory(context.Background(), map[string]string{admissionBudgetSub: "One"}, rows, []InventoryScopeDeclaration{{admissionBudgetSub, "complete"}})
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func TestRequiredAdmissionFrozenInputBudgets(t *testing.T) {
	var data struct {
		Cases []struct {
			Name      string
			Kinds     int
			Counts    []int          `json:"group_counts"`
			Width     int            `json:"key_bytes"`
			Overrides map[string]int `json:"key_overrides"`
			Error     string         `json:"want_error"`
			Catalogs  int            `json:"want_catalogs"`
		}
	}
	if e := json.Unmarshal(requiredFixtureRead(t, "../../..", "docs/evidence/20261022/inventory-budget/budget-recipes.json"), &data); e != nil {
		t.Fatal(e)
	}
	if len(data.Cases) != 4 {
		t.Fatal("input recipe inventory")
	}
	for _, c := range data.Cases {
		t.Run(c.Name, func(t *testing.T) {
			rows := []assessment.Resource{}
			cs := []RequiredCollectionInput{}
			ordinal := 0
			for i := 0; i < c.Kinds; i++ {
				kind := fmt.Sprint("c", i)
				rows = append(rows, admissionBudgetResource(cognitiveSKUType, kind, "s0", kind))
				p := &CognitiveSKUCollection{SubscriptionID: admissionBudgetSub, Region: "westeurope", Kind: kind, ResourceType: cognitiveSKUType, Status: "complete", Groups: map[string]CognitiveSKUGroup{}}
				for j := 0; j < c.Counts[i]; j++ {
					width := c.Width
					if n, ok := c.Overrides[strconv.Itoa(ordinal)]; ok {
						width = n
					}
					key := fmt.Sprintf("g%05d", j) + strings.Repeat("x", width-6)
					p.Groups[key] = CognitiveSKUGroup{Known: true, State: SKUAvailable, BlockedZones: []string{}}
					ordinal++
				}
				cs = append(cs, RequiredCollectionInput{SubscriptionID: admissionBudgetSub, SourceRegion: "eastus", TargetRegion: "westeurope", ResourceType: cognitiveSKUType, Kind: kind, Cognitive: p})
			}
			p := admissionBudgetPlan(t, rows)
			a, e := AdmitRequiredEvidence(context.Background(), p, RequiredEvidenceRequest{admissionBudgetSub, "eastus", "westeurope"}, cs)
			if c.Error != "" {
				if a != nil || e == nil || e.Error() != c.Error {
					t.Fatal("ADMISSION_INPUT_BUDGET_REJECTION", a != nil, e)
				}
				return
			}
			if e != nil || a == nil || a.Status != "complete" || len(a.Catalogs) != c.Catalogs || a.SKUs[cognitiveSKUType].Status != "unsupported" || len(a.SKUs[cognitiveSKUType].Values) != 0 {
				t.Fatal("ADMISSION_INPUT_BUDGET_EXACT", e)
			}
			for _, b := range a.Catalogs {
				if b.Status != "complete" || b.DeploymentEligibility != "not_assessed" || len(b.Cognitive.Values) != 1 || b.Cognitive.Values["s0"].State != SKUUnavailable {
					t.Fatal("ADMISSION_INPUT_BUDGET_RESULT")
				}
			}
		})
	}
}
func TestRequiredAdmissionFrozenOutputBudgets(t *testing.T) {
	cases := []struct {
		Name                               string
		VM, Zones, Width, Full, Extra, SQL int
		Family, Error                      string
	}{{"entries-exact", 8191, 13650, 1, 0, 0, 1, "", ""}, {"entries-first-excess", 8191, 13650, 1, 0, 0, 2, "", "required_evidence_limit"}, {"text-exact", 8190, 13650, 400, 546, 4, 1, "ff", ""}, {"text-first-excess", 8190, 13650, 400, 546, 4, 1, "fff", "required_evidence_limit"}}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			rows := []assessment.Resource{}
			vm := &ComputeVMSKUCollection{SubscriptionID: admissionBudgetSub, Region: "westeurope", ResourceType: vmSKUType, Status: "complete", Groups: map[string]ComputeVMSKUGroup{}}
			ordinal := 0
			for i := 0; i < c.VM; i++ {
				key := fmt.Sprintf("g%05d", i)
				rows = append(rows, admissionBudgetResource(vmSKUType, fmt.Sprintf("r%05d", i), key, ""))
				zones := []string{}
				n := 1
				if i < c.Zones-c.VM {
					n = 2
				}
				for j := 0; j < n; j++ {
					width := c.Width
					if ordinal < c.Full {
						width = 512
					} else if ordinal == c.Full {
						width += c.Extra
					}
					zones = append(zones, strings.Repeat(fmt.Sprint(j+1), width))
					ordinal++
				}
				vm.Groups[key] = ComputeVMSKUGroup{Known: true, State: SKUZoneRestricted, BlockedZones: zones}
			}
			sqlRow := admissionBudgetResource(sqlDatabaseType, "db", "x", "")
			sqlRow.ID = "/subscriptions/" + admissionBudgetSub + "/resourceGroups/rg/providers/microsoft.sql/servers/server/databases/db"
			rows = append(rows, sqlRow)
			sql := &SQLCapabilityCollection{SubscriptionID: admissionBudgetSub, Region: "westeurope", ResourceType: sqlDatabaseType, Status: "complete", Groups: map[string]SQLCapabilityGroup{}}
			for _, name := range []string{"a", "b"}[:c.SQL] {
				sql.Groups["v/e/"+name] = SQLCapabilityGroup{Known: true, State: SKUAvailable, Version: "v", Edition: "e", Name: name, Family: c.Family, SKUName: name, SKUTier: "e", Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"}
			}
			p := admissionBudgetPlan(t, rows)
			q := RequiredEvidenceRequest{admissionBudgetSub, "eastus", "westeurope"}
			cs := []RequiredCollectionInput{{SubscriptionID: admissionBudgetSub, SourceRegion: q.SourceRegion, TargetRegion: q.TargetRegion, ResourceType: vmSKUType, VM: vm}, {SubscriptionID: admissionBudgetSub, SourceRegion: q.SourceRegion, TargetRegion: q.TargetRegion, ResourceType: sqlDatabaseType, SQL: sql}}
			a, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
			if c.Error != "" {
				if a != nil || e == nil || e.Error() != c.Error {
					t.Fatal("ADMISSION_OUTPUT_BUDGET_REJECTION", a != nil, e)
				}
				return
			}
			if e != nil || a == nil || a.Status != "complete" || len(a.SKUs[vmSKUType].Values) != c.VM || len(a.Catalogs) != 1 || len(a.Catalogs[0].SQL.Values) != c.SQL {
				t.Fatal("ADMISSION_OUTPUT_BUDGET_EXACT", e)
			}
		})
	}
}
