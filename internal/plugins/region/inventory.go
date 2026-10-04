// Copyright (c) Microsoft Corporation.
// Inventory formatting derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_INVENTORY.md for provenance and corrections.
package region

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/skus"
)

const MaxInventoryRows = 8192

func pendingInventory() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "inventory", Metadata: Metadata(), SheetName: "Region Inventory", Description: "Resource inventory collected during region selection analysis", Columns: []string{"Subscription Id", "Resource Group", "Location", "Resource Type", "Resource Name", "Sku Name", "Sku Tier", "Capacity", "Kind", "Resource Id"}, Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped}}
}

func inventoryFailure(code string) (*assessment.PluginTable, error) {
	table := pendingInventory()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "region inventory input could not be safely projected"}}
	return &table, fmt.Errorf("region inventory input could not be safely projected")
}

func inventoryCancelled(ctx context.Context) (*assessment.PluginTable, error) {
	table, _ := inventoryFailure("region_inventory_cancelled")
	return table, ctx.Err()
}

// inventoryIdentity validates correlation and the exact source masking prefix.
// It does not claim full ARM provider syntax or derive displayed resource fields.
func inventoryIdentity(id, subscription string) bool {
	if len(id) < 53 || !strings.HasPrefix(id, "/subscriptions/") || id[51] != '/' || !strings.EqualFold(id[15:51], subscription) || strings.ContainsAny(id, "?#\\") {
		return false
	}
	for _, segment := range strings.Split(id[52:], "/") {
		if segment == "" {
			return false
		}
	}
	return true
}

func inventoryCells(resource assessment.Resource, mask bool) []string {
	return []string{redact.SubscriptionID(resource.SubscriptionID, mask), resource.ResourceGroup, resource.Location, resource.Type, resource.Name, resource.SKUName, resource.SKUTier, skus.ComputeCapacity(resource.SKUName, resource.SKUCapacity, strings.EqualFold(resource.Type, "microsoft.compute/virtualmachinescalesets")), resource.Kind, redact.SubscriptionIDInResourceID(resource.ID, mask)}
}

// ProjectInventory preserves source cells/order with an explicit sheet-name
// correction. Inputs remain stable during the call; output owns mutable state.
func ProjectInventory(ctx context.Context, subscriptions map[string]string, resources []assessment.Resource, mask bool) (*assessment.PluginTable, error) {
	if ctx.Err() != nil {
		return inventoryCancelled(ctx)
	}
	decodedBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &decodedBytes)
	if ctx.Err() != nil {
		return inventoryCancelled(ctx)
	}
	if !ok {
		return inventoryFailure("region_inventory_scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return inventoryFailure("region_inventory_scope_invalid")
		}
	}
	if len(resources) > MaxInventoryRows {
		return inventoryFailure("region_inventory_input_limit")
	}
	seen := map[string]bool{}
	projectedBytes := 0
	for _, resource := range resources {
		if ctx.Err() != nil {
			return inventoryCancelled(ctx)
		}
		key := strings.ToLower(resource.SubscriptionID)
		_, selected := scope[key]
		if len(resource.SubscriptionID) != 36 || !subscriptionID.MatchString(resource.SubscriptionID) || !selected {
			return inventoryFailure("region_inventory_scope_invalid")
		}
		if resource.ID == "" || resource.Type == "" || resource.Name == "" {
			return inventoryFailure("region_inventory_input_invalid")
		}
		for _, label := range []string{resource.ID, resource.SubscriptionID, resource.ResourceGroup, resource.Location, resource.Type, resource.Name, resource.SKUName, resource.SKUTier, resource.Kind} {
			if !auxiliaryText(label, &decodedBytes) {
				return inventoryFailure("region_inventory_text_limit")
			}
		}
		if !inventoryIdentity(resource.ID, resource.SubscriptionID) {
			return inventoryFailure("region_inventory_scope_invalid")
		}
		idKey := strings.ToLower(resource.ID)
		if seen[idKey] {
			return inventoryFailure("region_inventory_input_invalid")
		}
		seen[idKey] = true
		capacity := int64(resource.SKUCapacity)
		if capacity < -MaxAuxCount || capacity > MaxAuxCount {
			return inventoryFailure("region_inventory_capacity_invalid")
		}
		if capacity > 0 && strings.EqualFold(resource.Type, "microsoft.compute/virtualmachinescalesets") {
			if sku, known := skus.Lookup(resource.SKUName); known && sku.VCPUs > 0 && capacity > MaxAuxCount/int64(sku.VCPUs) {
				return inventoryFailure("region_inventory_capacity_invalid")
			}
		}
		// Compute only small transient cells after checked integer arithmetic.
		// Account every displayed byte before allocating any report rows.
		for _, cell := range inventoryCells(resource, mask) {
			projectedBytes += len(cell)
			if projectedBytes > auxTextBudget {
				return inventoryFailure("region_inventory_text_limit")
			}
		}
	}
	table := pendingInventory()
	table.Rows = make([]assessment.PluginRow, 0, len(resources))
	for _, resource := range resources {
		if ctx.Err() != nil {
			return inventoryCancelled(ctx)
		}
		table.Rows = append(table.Rows, assessment.PluginRow{SubscriptionID: strings.ToLower(resource.SubscriptionID), Cells: inventoryCells(resource, mask)})
	}
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(table.Rows)}
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		return inventoryFailure("region_inventory_output_limit")
	}
	if ctx.Err() != nil {
		return inventoryCancelled(ctx)
	}
	return &table, nil
}
