package region

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestRequiredBuildFrozenBudgetRecipes(t *testing.T) {
	var data struct {
		Cases []struct {
			Name   string `json:"name"`
			Rows   int    `json:"rows"`
			Scopes int    `json:"scopes"`
			Short  int    `json:"short_names"`
			Next   int    `json:"next_name"`
			Error  string `json:"want_error"`
		} `json:"cases"`
	}
	raw := requiredFixtureRead(t, "../../..", "docs/evidence/20261026/inventory-output/budget-recipes.json")
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, c := range data.Cases {
		if !strings.HasPrefix(c.Name, "build-") {
			continue
		}
		count++
		t.Run(c.Name, func(t *testing.T) {
			const sub = "aaaaaaaa-1111-4111-8111-111111111111"
			selected := map[string]string{sub: "One"}
			declarations := []InventoryScopeDeclaration{{sub, "complete"}}
			rows := make([]assessment.Resource, 0, c.Rows)
			if strings.Contains(c.Name, "entries") {
				for i := 1; i < c.Scopes; i++ {
					id := strings.Repeat(string(rune('a'+i)), 8) + "-1111-4111-8111-111111111111"
					selected[id] = "One"
					declarations = append(declarations, InventoryScopeDeclaration{id, "complete"})
				}
				for i := 0; i < c.Rows; i++ {
					name := fmt.Sprintf("r%05d", i)
					rows = append(rows, assessment.Resource{ID: "/subscriptions/" + sub + "/resourceGroups/rg/providers/microsoft.compute/virtualmachines/" + name, SubscriptionID: sub, ResourceGroup: "rg", Name: name, Type: "microsoft.compute/virtualmachines", Location: name, SKUName: fmt.Sprintf("g%05d", i)})
				}
			} else {
				for i := 0; i < c.Rows; i++ {
					width := 512
					if i < c.Short {
						width = 1
					} else if i == c.Short {
						width = c.Next
					}
					rows = append(rows, assessment.Resource{ID: fmt.Sprintf("/subscriptions/%s/resourceGroups/rg/providers/microsoft.example/widgets/r%04d", sub, i), SubscriptionID: sub, ResourceGroup: "rg", Name: strings.Repeat("n", width), Type: "microsoft.example/widgets", Location: "global", Kind: strings.Repeat("k", 512), SKUName: strings.Repeat("s", 512), SKUTier: strings.Repeat("t", 512), SKUFamily: strings.Repeat("f", 512)})
				}
			}
			p, err := BuildRequiredInventory(context.Background(), selected, rows, declarations)
			if c.Error != "" {
				if p != nil || err == nil || err.Error() != c.Error {
					t.Fatal("BUILD_BUDGET_REJECTION", p != nil, err)
				}
				return
			}
			if err != nil || p == nil || len(p.Resources) != c.Rows || p.Status != "complete" {
				t.Fatal("BUILD_BUDGET_EXACT", err)
			}
			if strings.Contains(c.Name, "entries") && (len(p.Requirements) != 8191 || len(p.Scopes) != 4) {
				t.Fatal("BUILD_ENTRIES_EXACT_SHAPE")
			}
		})
	}
	if count != 6 {
		t.Fatal("frozen Build recipe inventory changed")
	}
}
func TestRequiredBuildOuterAndScalarLimits(t *testing.T) {
	selected, rows, d := buildInput(t)
	many := make([]assessment.Resource, 8193)
	p, e := BuildRequiredInventory(context.Background(), selected, many, d)
	if p != nil || e == nil || e.Error() != "required_inventory_limit" {
		t.Fatal("BUILD_OUTER_LIMIT")
	}
	declarations := make([]InventoryScopeDeclaration, 1001)
	p, e = BuildRequiredInventory(context.Background(), selected, rows, declarations)
	if p != nil || e == nil || e.Error() != "required_inventory_limit" {
		t.Fatal("BUILD_DECLARATION_OUTER_LIMIT")
	}
	scope := map[string]string{}
	for i := 0; i < 1001; i++ {
		scope[fmt.Sprint(i)] = ""
	}
	p, e = BuildRequiredInventory(context.Background(), scope, rows, d)
	if p != nil || e == nil || e.Error() != "required_inventory_limit" {
		t.Fatal("BUILD_SCOPE_OUTER_LIMIT")
	}
	rows[0].Name = strings.Repeat("x", 513)
	p, e = BuildRequiredInventory(context.Background(), selected, rows, d)
	if p != nil || e == nil || e.Error() != "required_inventory_invalid" {
		t.Fatal("BUILD_SCALAR_INVALID")
	}
}
func TestRequiredBuildStorageKindFirstExcess(t *testing.T) {
	const sub = "aaaaaaaa-1111-4111-8111-111111111111"
	selected := map[string]string{sub: "One"}
	d := []InventoryScopeDeclaration{{sub, "complete"}}
	rows := []assessment.Resource{}
	for i := 0; i < 17; i++ {
		kind := ""
		if i > 0 {
			kind = fmt.Sprintf("k%02d", i)
		}
		rows = append(rows, assessment.Resource{ID: fmt.Sprintf("/subscriptions/%s/resourceGroups/rg/providers/microsoft.storage/storageaccounts/s%02d", sub, i), SubscriptionID: sub, ResourceGroup: "rg", Name: fmt.Sprintf("s%02d", i), Type: "microsoft.storage/storageaccounts", Location: "eastus", SKUName: "Standard_LRS", Kind: kind})
		if i == 15 {
			p, e := BuildRequiredInventory(context.Background(), selected, rows, d)
			if e != nil || p == nil || !p.Requirements[0].MissingKind || len(p.Requirements[0].Kinds["standard_lrs"]) != 16 {
				t.Fatal("BUILD_STORAGE_SIXTEEN_KINDS", e)
			}
		}
	}
	p, e := BuildRequiredInventory(context.Background(), selected, rows, d)
	if p != nil || e == nil || e.Error() != "required_inventory_limit" {
		t.Fatal("BUILD_STORAGE_SEVENTEENTH_KIND", p != nil, e)
	}
}
