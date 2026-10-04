// Copyright (c) Microsoft Corporation.
// Service availability formatting derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_SERVICE_AVAILABILITY.md for provenance.
package region

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxServiceTargets = 32
const MaxServiceRows = 8192
const MaxServiceEntries = 65536
const serviceTextBudget = assessment.MaxPluginTextBytes - (64 << 10)

// ServiceInventory is an already decoded aggregate. Contributor identities are
// declarations checked against selection, not proof of resource ownership.
// The later collector must construct this aggregate from selected resources.
type ServiceInventory struct {
	SubscriptionIDs []string

	ResourceTypes map[string]int64

	SKUsByType, ResourceTypesByRegion map[string]map[string]int64
}

// ServiceComparison contains only inputs consumed by the source sheet helper.
// RestrictedSKUs is bounded but intentionally has no sheet semantics in source.
type ServiceComparison struct {
	SubscriptionID, SubscriptionName, SourceRegion, TargetRegion string

	MissingResourceTypes, MissingSKUs, RestrictedSKUs, ZoneRestrictedSKUs []string
}

func serviceColumns() []string {
	return []string{"ResourceType", "ResourceCount", "ImplementedRegions", "SKUCount", "SKU", "SKU available", "SKU zone-restricted", "Service available"}
}

func serviceFailure(code string) ([]assessment.PluginTable, error) {
	return []assessment.PluginTable{{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "service-availability", Metadata: Metadata(), SheetName: "Service Availability", Description: "Service and SKU availability could not be safely projected", Columns: serviceColumns(), Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageFailed, Error: &assessment.AssessmentError{Code: code, Message: "region service availability input could not be safely projected"}}}}, fmt.Errorf("region service availability input could not be safely projected")
}

func serviceCancelled(ctx context.Context) ([]assessment.PluginTable, error) {
	tables, _ := serviceFailure("region_service_cancelled")
	return tables, ctx.Err()
}

func serviceSupported(resourceType string) bool {
	// Exact pinned registry.Get keys; no provider construction or FetchSKUs.
	switch strings.ToLower(resourceType) {
	case "microsoft.compute/virtualmachines", "microsoft.compute/virtualmachinescalesets", "microsoft.compute/disks", "microsoft.sql/servers/databases", "microsoft.sql/managedinstances", "microsoft.storage/storageaccounts", "microsoft.cognitiveservices/accounts":
		return true
	default:
		return false
	}
}

type serviceSets struct {
	missingTypes, missingSKUs, zones map[string]bool
}

type servicePreparedRow struct {
	groups [][]string
}

// serviceMeasure includes delimiters before any joined cell allocation.
func serviceMeasure(groups [][]string, total *int) bool {
	for _, group := range groups {
		units := max(0, len(group)-1) * 2
		bytes := units
		for _, s := range group {
			units += utf16Units(s)
			bytes += len(s)
		}
		if units > assessment.MaxPluginCellUnits {
			return false
		}
		*total += bytes
		if *total > serviceTextBudget {
			return false
		}
	}
	return true
}

// ProjectServiceAvailability preserves the source helper's cells and nil empty
// branches. Inputs must remain stable during the call. No input maps or slices
// escape, and no service, credential, score or availability calculation runs.
func ProjectServiceAvailability(ctx context.Context, subscriptions map[string]string, inventory *ServiceInventory, comparisons []ServiceComparison) ([]assessment.PluginTable, error) {
	if ctx.Err() != nil {
		return serviceCancelled(ctx)
	}
	if len(comparisons) > MaxComparisons {
		return serviceFailure("region_service_input_limit")
	}
	decodedBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &decodedBytes)
	if ctx.Err() != nil {
		return serviceCancelled(ctx)
	}
	if !ok {
		return serviceFailure("region_service_scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return serviceFailure("region_service_scope_invalid")
		}
	}
	entries := len(scope)
	entry := func(n int) bool {
		if n > MaxServiceEntries-entries {
			return false
		}
		entries += n
		return true
	}
	text := func(s string) bool { return s != "" && auxiliaryText(s, &decodedBytes) }
	contributors := map[string]bool{}
	if inventory != nil {
		if len(inventory.SubscriptionIDs) > MaxSubscriptions || !entry(len(inventory.SubscriptionIDs)) || !entry(len(inventory.ResourceTypes)) || !entry(len(inventory.SKUsByType)) || !entry(len(inventory.ResourceTypesByRegion)) {
			return serviceFailure("region_service_input_limit")
		}
		for _, id := range inventory.SubscriptionIDs {
			if ctx.Err() != nil {
				return serviceCancelled(ctx)
			}
			key := strings.ToLower(id)
			_, selected := scope[key]
			if len(id) != 36 || !subscriptionID.MatchString(id) || !selected || contributors[key] {
				return serviceFailure("region_service_scope_invalid")
			}
			contributors[key] = true
		}
		if len(inventory.ResourceTypes) > 0 && len(contributors) == 0 {
			return serviceFailure("region_service_scope_invalid")
		}
		validateCounts := func(values map[string]int64) bool {
			for label, count := range values {
				if ctx.Err() != nil || !text(label) || count < 0 || count > MaxAuxCount {
					return false
				}
			}
			return true
		}
		if !validateCounts(inventory.ResourceTypes) {
			if ctx.Err() != nil {
				return serviceCancelled(ctx)
			}
			return serviceFailure("region_service_input_invalid")
		}
		for rt, skus := range inventory.SKUsByType {
			if !entry(len(skus)) {
				return serviceFailure("region_service_input_limit")
			}
			if !text(rt) || !validateCounts(skus) {
				if ctx.Err() != nil {
					return serviceCancelled(ctx)
				}
				return serviceFailure("region_service_input_invalid")
			}
		}
		for location, types := range inventory.ResourceTypesByRegion {
			if !entry(len(types)) {
				return serviceFailure("region_service_input_limit")
			}
			if len(location) > 64 || !regionID.MatchString(location) || !text(location) || !validateCounts(types) {
				if ctx.Err() != nil {
					return serviceCancelled(ctx)
				}
				return serviceFailure("region_service_input_invalid")
			}
		}
	}
	targets := map[string]bool{}
	seen := map[[3]string]bool{}
	if !entry(len(comparisons)) {
		return serviceFailure("region_service_input_limit")
	}
	for _, c := range comparisons {
		if ctx.Err() != nil {
			return serviceCancelled(ctx)
		}
		if !auxiliarySelected(scope, c.SubscriptionID, c.SubscriptionName) || inventory != nil && !contributors[strings.ToLower(c.SubscriptionID)] {
			return serviceFailure("region_service_scope_invalid")
		}
		if len(c.SourceRegion) > 64 || len(c.TargetRegion) > 64 || !regionID.MatchString(c.SourceRegion) || !regionID.MatchString(c.TargetRegion) {
			return serviceFailure("region_service_input_invalid")
		}
		for _, s := range []string{c.SubscriptionName, c.SourceRegion, c.TargetRegion} {
			if !text(s) {
				return serviceFailure("region_service_text_limit")
			}
		}
		key := [3]string{strings.ToLower(c.SubscriptionID), c.SourceRegion, c.TargetRegion}
		if seen[key] {
			return serviceFailure("region_service_input_invalid")
		}
		seen[key] = true
		targets[c.TargetRegion] = true
		if len(targets) > MaxServiceTargets {
			return serviceFailure("region_service_input_limit")
		}
		for _, details := range [][]string{c.MissingResourceTypes, c.MissingSKUs, c.RestrictedSKUs, c.ZoneRestrictedSKUs} {
			if !entry(len(details)) {
				return serviceFailure("region_service_input_limit")
			}
			for _, detail := range details {
				if ctx.Err() != nil {
					return serviceCancelled(ctx)
				}
				if !text(detail) {
					return serviceFailure("region_service_text_limit")
				}
			}
		}
	}
	if ctx.Err() != nil {
		return serviceCancelled(ctx)
	}
	if inventory == nil || len(comparisons) == 0 {
		return nil, nil
	}
	if len(inventory.ResourceTypes) > MaxServiceRows/len(targets) {
		return serviceFailure("region_service_input_limit")
	}
	sortedTargets := make([]string, 0, len(targets))
	for target := range targets {
		sortedTargets = append(sortedTargets, target)
	}
	sort.Strings(sortedTargets)
	sheets := map[string]bool{}
	sets := map[string]serviceSets{}
	for _, target := range sortedTargets {
		name := "Svc Avail " + target
		if len(name) > 31 {
			name = name[:31] // Targets are validated ASCII; source truncates runes.
		}
		if sheets[name] {
			return serviceFailure("region_service_sheet_collision")
		}
		sheets[name] = true
		sets[target] = serviceSets{map[string]bool{}, map[string]bool{}, map[string]bool{}}
	}
	for _, c := range comparisons {
		set := sets[c.TargetRegion]
		for _, rt := range c.MissingResourceTypes {
			set.missingTypes[strings.ToLower(rt)] = true
		}
		for _, sku := range c.MissingSKUs {
			set.missingSKUs[strings.ToLower(sku)] = true
		}
		for _, sku := range c.ZoneRestrictedSKUs {
			// Source splits before lowercasing. Preserve that case-sensitive suffix.
			set.zones[strings.ToLower(strings.SplitN(sku, " (zones", 2)[0])] = true
		}
	}
	impl := map[string][]string{}
	for location, types := range inventory.ResourceTypesByRegion {
		for rt := range types {
			impl[rt] = append(impl[rt], location)
		}
	}
	for _, regions := range impl {
		sort.Strings(regions)
	}
	sortedTypes := make([]string, 0, len(inventory.ResourceTypes))
	skusByType := map[string][]string{}
	for rt := range inventory.ResourceTypes {
		sortedTypes = append(sortedTypes, rt)
		for sku := range inventory.SKUsByType[rt] {
			skusByType[rt] = append(skusByType[rt], sku)
		}
		sort.Strings(skusByType[rt])
	}
	sort.Strings(sortedTypes)
	prepared := make([][]servicePreparedRow, len(sortedTargets))
	outputBytes := 0
	for i, target := range sortedTargets {
		set := sets[target]
		for _, rt := range sortedTypes {
			if ctx.Err() != nil {
				return serviceCancelled(ctx)
			}
			skus := skusByType[rt]
			skuDisplay := skus
			zones := []string{}
			skuStatus := "Available"
			serviceStatus := "Available"
			if set.missingTypes[strings.ToLower(rt)] {
				serviceStatus = "Not available"
			}
			if len(skus) == 0 || !serviceSupported(rt) {
				skuStatus = "N/A"
				zones = []string{"N/A"}
				if len(skus) == 0 {
					skuDisplay = []string{"N/A"}
				}
			} else {
				missing := []string{}
				for _, sku := range skus {
					key := strings.ToLower(rt + ":" + sku)
					if set.missingSKUs[key] {
						missing = append(missing, sku)
					}
					if set.zones[key] {
						zones = append(zones, sku)
					}
				}
				if len(missing) > 0 {
					skuStatus = "Not available"
					skuDisplay = missing
				}
			}
			groups := [][]string{{rt}, {strconv.FormatInt(inventory.ResourceTypes[rt], 10)}, impl[strings.ToLower(rt)], {strconv.Itoa(len(skus))}, skuDisplay, {skuStatus}, zones, {serviceStatus}}
			if !serviceMeasure(groups, &outputBytes) {
				return serviceFailure("region_service_text_limit")
			}
			prepared[i] = append(prepared[i], servicePreparedRow{groups})
		}
	}
	tables := make([]assessment.PluginTable, len(sortedTargets))
	for i, target := range sortedTargets {
		if ctx.Err() != nil {
			return serviceCancelled(ctx)
		}
		sheet := "Svc Avail " + target
		if len(sheet) > 31 {
			sheet = sheet[:31]
		}
		table := assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: fmt.Sprintf("service-availability-%03d", i+1), Metadata: Metadata(), SheetName: sheet, Description: "Service and SKU availability for target region: " + target, Columns: serviceColumns(), Rows: make([]assessment.PluginRow, len(prepared[i]))}
		for j, row := range prepared[i] {
			if ctx.Err() != nil {
				return serviceCancelled(ctx)
			}
			cells := make([]string, len(row.groups))
			for k, group := range row.groups {
				cells[k] = strings.Join(group, ", ")
			}
			table.Rows[j] = assessment.PluginRow{Cells: cells}
		}
		table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(table.Rows)}
		tables[i] = table
	}
	if assessment.ValidatePluginTables(tables) != nil {
		return serviceFailure("region_service_output_limit")
	}
	if ctx.Err() != nil {
		return serviceCancelled(ctx)
	}
	return tables, nil
}
