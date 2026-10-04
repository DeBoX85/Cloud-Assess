// Copyright (c) Microsoft Corporation.
// Inventory counting derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_INVENTORY_AGGREGATION.md for corrections.
package region

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxInventoryEntries = 65536

type InventoryCounts struct {
	ResourceTypes         map[string]int64
	SKUsByType            map[string]map[string]int64
	LocationCounts        map[string]int64
	ResourceTypesByRegion map[string]map[string]int64
	SKUsByTypeAndRegion   map[string]map[string]map[string]int64
}

type InventoryCalculation struct {
	Subscriptions              map[string]InventoryCounts
	Aggregate                  InventoryCounts
	ContributorSubscriptionIDs []string
	Records                    int
}

func emptyInventoryCounts() InventoryCounts {
	return InventoryCounts{ResourceTypes: map[string]int64{}, SKUsByType: map[string]map[string]int64{}, LocationCounts: map[string]int64{}, ResourceTypesByRegion: map[string]map[string]int64{}, SKUsByTypeAndRegion: map[string]map[string]map[string]int64{}}
}

func normalizeInventoryRegion(region string) string {
	return strings.ToLower(strings.ReplaceAll(region, " ", ""))
}

func inventoryCalculationFailure(code string) (*InventoryCalculation, error) {
	return nil, fmt.Errorf("region inventory calculation input could not be safely aggregated [%s]", code)
}

// incrementInventory charges every newly allocated map key before insertion.
// It counts a resource once, preserving raw SKU keys and normalized locations.
func incrementInventory(counts InventoryCounts, resourceType, location, sku string, insert func(string) bool) bool {
	count := func(values map[string]int64, key string) bool {
		if _, exists := values[key]; !exists && !insert(key) {
			return false
		}
		values[key]++
		return true
	}
	if !count(counts.ResourceTypes, resourceType) || !count(counts.LocationCounts, location) {
		return false
	}
	if counts.ResourceTypesByRegion[location] == nil {
		if !insert(location) {
			return false
		}
		counts.ResourceTypesByRegion[location] = map[string]int64{}
	}
	if !count(counts.ResourceTypesByRegion[location], resourceType) {
		return false
	}
	if sku == "" {
		return true
	}
	if counts.SKUsByType[resourceType] == nil {
		if !insert(resourceType) {
			return false
		}
		counts.SKUsByType[resourceType] = map[string]int64{}
	}
	if !count(counts.SKUsByType[resourceType], sku) {
		return false
	}
	if counts.SKUsByTypeAndRegion[resourceType] == nil {
		if !insert(resourceType) {
			return false
		}
		counts.SKUsByTypeAndRegion[resourceType] = map[string]map[string]int64{}
	}
	if counts.SKUsByTypeAndRegion[resourceType][location] == nil {
		if !insert(location) {
			return false
		}
		counts.SKUsByTypeAndRegion[resourceType][location] = map[string]int64{}
	}
	return count(counts.SKUsByTypeAndRegion[resourceType][location], sku)
}

// CalculateInventory creates only owned per-run counters from selected decoded
// records. Inputs remain stable; no collection/provider completeness is inferred.
func CalculateInventory(ctx context.Context, subscriptions map[string]string, resources []assessment.Resource) (*InventoryCalculation, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	decodedBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &decodedBytes)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if !ok {
		return inventoryCalculationFailure("scope_invalid")
	}
	if len(resources) > MaxInventoryRows {
		return inventoryCalculationFailure("input_limit")
	}
	for _, name := range scope {
		if name == "" {
			return inventoryCalculationFailure("scope_invalid")
		}
	}
	seen := map[string]bool{}
	for _, resource := range resources {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		_, selected := scope[strings.ToLower(resource.SubscriptionID)]
		if len(resource.SubscriptionID) != 36 || !subscriptionID.MatchString(resource.SubscriptionID) || !selected {
			return inventoryCalculationFailure("scope_invalid")
		}
		if resource.ID == "" || resource.Type == "" {
			return inventoryCalculationFailure("input_invalid")
		}
		for _, label := range []string{resource.ID, resource.SubscriptionID, resource.Type, resource.Location, resource.SKUName} {
			if !auxiliaryText(label, &decodedBytes) {
				return inventoryCalculationFailure("text_limit")
			}
		}
		if !inventoryIdentity(resource.ID, resource.SubscriptionID) {
			return inventoryCalculationFailure("scope_invalid")
		}
		if seen[strings.ToLower(resource.ID)] {
			return inventoryCalculationFailure("input_invalid")
		}
		seen[strings.ToLower(resource.ID)] = true
		for _, label := range []string{strings.ToLower(resource.Type), normalizeInventoryRegion(resource.Location)} {
			if !auxiliaryText(label, &decodedBytes) {
				return inventoryCalculationFailure("text_limit")
			}
		}
	}
	entries, outputBytes := 0, 0
	insert := func(key string) bool {
		if ctx.Err() != nil || entries >= MaxInventoryEntries || len(key) > auxTextBudget-outputBytes {
			return false
		}
		entries++
		outputBytes += len(key)
		return true
	}
	result := &InventoryCalculation{Subscriptions: map[string]InventoryCounts{}, Aggregate: emptyInventoryCounts(), ContributorSubscriptionIDs: []string{}, Records: len(resources)}
	for id := range scope {
		if !insert(id) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return inventoryCalculationFailure("output_limit")
		}
		result.Subscriptions[id] = emptyInventoryCounts()
	}
	contributors := map[string]bool{}
	for _, resource := range resources {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		id := strings.ToLower(resource.SubscriptionID)
		counts, selected := result.Subscriptions[id]
		if !selected {
			return inventoryCalculationFailure("output_invalid")
		}
		if !contributors[id] {
			if !insert(id) {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				return inventoryCalculationFailure("output_limit")
			}
			contributors[id] = true
			result.ContributorSubscriptionIDs = append(result.ContributorSubscriptionIDs, id)
		}
		resourceType, location := strings.ToLower(resource.Type), normalizeInventoryRegion(resource.Location)
		if !incrementInventory(counts, resourceType, location, resource.SKUName, insert) || !incrementInventory(result.Aggregate, resourceType, location, resource.SKUName, insert) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return inventoryCalculationFailure("output_limit")
		}
	}
	sort.Strings(result.ContributorSubscriptionIDs)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return result, nil
}
