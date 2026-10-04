// Copyright (c) Microsoft Corporation.
// Availability arithmetic derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_AVAILABILITY_RUNTIME.md for corrections.
package region

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxAvailabilityEntries = 65536

// SKU states match the pinned source. Future states remain unavailable.
const (
	SKUAvailable = iota
	SKURestricted
	SKUUnavailable
	SKUZoneRestricted
)

type AvailabilityRequest struct {
	SubscriptionID string
	SourceRegion string
	TargetRegion string
}

type SKUAvailability struct {
	State int
	BlockedZones []string
}

// Status is complete, unknown or unsupported. Complete empty responses are
// valid; absent declarations do not imply unsupported or successful collection.
type SKUEvidence struct {
	Status string
	Values map[string]SKUAvailability
}

type AvailabilityEvidence struct {
	SubscriptionID string
	TargetRegion string
	InventoryComplete bool
	ProvidersComplete bool
	Locations map[string]map[string]bool
	SKUEnabled bool
	SKUs map[string]SKUEvidence
	ZoneCounts map[string]int
}

type AvailabilityCalculation struct {
	Comparison Comparison
	Health assessment.StageExecution
}

// CalculateAvailability operates on stable, decoded per-run inputs. Health
// describes their declared completeness, not Azure or collection certification.
// No caller-owned slices/maps escape and failures return no partial comparison.
func CalculateAvailability(ctx context.Context, subscriptions map[string]string, inventory *InventoryCalculation, request AvailabilityRequest, evidence AvailabilityEvidence) (*AvailabilityCalculation, error) {
	fail := func(code string) (*AvailabilityCalculation, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("region availability input could not be safely calculated [%s]", code)
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	decodedBytes, entries := 0, 0
	charge := func(n int) bool {
		if ctx.Err() != nil || n > MaxAvailabilityEntries-entries {
			return false
		}
		entries += n
		return true
	}
	text := func(label string) bool { return auxiliaryText(label, &decodedBytes) }
	typeKey := func(key string) bool {
		namespace, name, found := strings.Cut(key, "/")
		return text(key) && found && namespace != "" && name != "" && key == strings.ToLower(key) && !strings.ContainsAny(key, " \t")
	}
	scope, ok := auxiliaryScope(ctx, subscriptions, &decodedBytes)
	if !ok || inventory == nil || len(inventory.Subscriptions) != len(scope) {
		return fail("scope_invalid")
	}
	for id, name := range scope {
		if name == "" {
			return fail("scope_invalid")
		}
		if _, exists := inventory.Subscriptions[id]; !exists {
			return fail("scope_invalid")
		}
	}
	id := strings.ToLower(request.SubscriptionID)
	name, selected := scope[id]
	if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || !selected || evidence.SubscriptionID != id || evidence.TargetRegion != request.TargetRegion {
		return fail("scope_invalid")
	}
	if len(request.SourceRegion) > 64 || len(request.TargetRegion) > 64 || !regionID.MatchString(request.SourceRegion) || !regionID.MatchString(request.TargetRegion) {
		return fail("region_invalid")
	}
	counts := inventory.Subscriptions[id]
	if !charge(len(counts.ResourceTypesByRegion) + len(counts.SKUsByTypeAndRegion) + len(evidence.Locations) + len(evidence.SKUs) + len(evidence.ZoneCounts)) {
		return fail("input_limit")
	}
	for region, types := range counts.ResourceTypesByRegion {
		// Inventory can contain logical/unassigned regions; only requested IDs
		// must be physical-format identifiers. Exact source keys are preserved.
		if !text(region) || !charge(len(types)) {
			return fail("input_limit")
		}
		for key, count := range types {
			if !typeKey(key) || count < 0 || count > MaxCount {
				return fail("input_invalid")
			}
		}
	}
	for key, regions := range counts.SKUsByTypeAndRegion {
		if !typeKey(key) || !charge(len(regions)) {
			return fail("input_limit")
		}
		for region, skus := range regions {
			if !text(region) || !charge(len(skus)) {
				return fail("input_limit")
			}
			for sku, count := range skus {
				if !text(sku) || strings.TrimSpace(sku) == "" || count < 0 || count > MaxCount {
					return fail("input_invalid")
				}
			}
		}
	}
	for key, locations := range evidence.Locations {
		if !typeKey(key) || !charge(len(locations)) {
			return fail("input_limit")
		}
		for location, present := range locations {
			if !text(location) || !present || !regionID.MatchString(location) {
				return fail("input_invalid")
			}
		}
	}
	normalized := map[string]map[string]SKUAvailability{}
	for key, sku := range evidence.SKUs {
		if !typeKey(key) || !charge(len(sku.Values)) || (sku.Status != "complete" && sku.Status != "unknown" && sku.Status != "unsupported") || (sku.Status != "complete" && len(sku.Values) != 0) {
			return fail("evidence_invalid")
		}
		values := map[string]SKUAvailability{}
		for raw, value := range sku.Values {
			if !text(raw) || strings.TrimSpace(raw) == "" || !charge(len(value.BlockedZones)) {
				return fail("input_limit")
			}
			key := strings.ToLower(strings.TrimSpace(raw))
			if _, exists := values[key]; exists {
				return fail("response_collision")
			}
			for _, zone := range value.BlockedZones {
				if !text(zone) || zone == "" {
					return fail("input_invalid")
				}
			}
			values[key] = value
		}
		normalized[key] = values
	}
	for region, count := range evidence.ZoneCounts {
		if !text(region) || !regionID.MatchString(region) || count < 0 || count > MaxCount {
			return fail("input_invalid")
		}
	}
	result := &AvailabilityCalculation{Comparison: Comparison{SubscriptionID: id, SubscriptionName: name, SourceRegion: request.SourceRegion, TargetRegion: request.TargetRegion, MissingResourceTypes: []string{}, MissingSKUs: []string{}, RestrictedSKUs: []string{}, ZoneRestrictedSKUs: []string{}, SourceZoneCount: evidence.ZoneCounts[request.SourceRegion], TargetZoneCount: evidence.ZoneCounts[request.TargetRegion]}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: 1}}
	warn := func(code string) {
		for _, warning := range result.Health.Warnings {
			if warning.Code == code {
				return
			}
		}
		result.Health.Status = assessment.StageCompletedWithWarnings
		result.Health.Warnings = append(result.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: "region availability evidence is incomplete"})
	}
	if !evidence.InventoryComplete {
		warn("availability_inventory_partial")
	}
	if !evidence.ProvidersComplete {
		warn("availability_providers_partial")
	}
	if _, exists := evidence.ZoneCounts[request.SourceRegion]; !exists {
		warn("availability_zones_partial")
	}
	if _, exists := evidence.ZoneCounts[request.TargetRegion]; !exists {
		warn("availability_zones_partial")
	}
	outputBytes, outputEntries := 0, 0
	appendDetail := func(list *[]string, detail string) bool {
		if len(*list) >= MaxDetails || outputEntries >= MaxDetailEntries || !auxiliaryText(detail, &outputBytes) {
			return false
		}
		outputEntries++
		*list = append(*list, detail)
		return true
	}
	c := &result.Comparison
	types, exists := counts.ResourceTypesByRegion[request.SourceRegion]
	if exists {
		c.SourceResourceTypeCount = len(types)
		for key := range types {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			locations := evidence.Locations[key]
			if locations["global"] || locations[request.TargetRegion] {
				c.AvailableTypes++
			} else {
				c.UnavailableTypes++
				if !appendDetail(&c.MissingResourceTypes, key) {
					return fail("output_limit")
				}
			}
		}
		if c.SourceResourceTypeCount > 0 {
			c.AvailabilityPercent = float64(c.AvailableTypes) / float64(c.SourceResourceTypeCount) * 100
		}
		for key, regions := range counts.SKUsByTypeAndRegion {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			skus := regions[request.SourceRegion]
			if len(skus) == 0 {
				continue
			}
			if !evidence.SKUEnabled {
				warn("availability_sku_disabled")
				continue
			}
			skuEvidence, declared := evidence.SKUs[key]
			if !declared {
				return fail("evidence_missing")
			}
			if skuEvidence.Status == "unsupported" {
				warn("availability_sku_unsupported")
				continue
			}
			for raw := range skus {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				c.TotalSKUsChecked++
				identifier := key + ":" + raw
				if skuEvidence.Status == "unknown" {
					c.UnknownSKUs++
					warn("availability_sku_unknown")
					if !appendDetail(&c.MissingSKUs, identifier+" (unknown)") {
						return fail("output_limit")
					}
					continue
				}
				value, found := normalized[key][strings.ToLower(strings.TrimSpace(raw))]
				if !found {
					value.State = SKUUnavailable
				}
				switch value.State {
				case SKUAvailable:
					c.AvailableSKUs++
				case SKURestricted:
					if !appendDetail(&c.RestrictedSKUs, identifier) {
						return fail("output_limit")
					}
				case SKUZoneRestricted:
					if len(value.BlockedZones) > 0 {
						length := len(identifier) + len(" (zones blocked: )") + len(value.BlockedZones) - 1
						for _, zone := range value.BlockedZones {
							length += len(zone)
						}
						if length > MaxLabelBytes {
							return fail("output_limit")
						}
						identifier += " (zones blocked: " + strings.Join(value.BlockedZones, ",") + ")"
					}
					if !appendDetail(&c.ZoneRestrictedSKUs, identifier) {
						return fail("output_limit")
					}
				default:
					c.UnavailableSKUs++
					if !appendDetail(&c.MissingSKUs, identifier) {
						return fail("output_limit")
					}
				}
			}
		}
		if c.TotalSKUsChecked > 0 {
			c.SKUAvailabilityPercent = 100
			if confirmed := c.TotalSKUsChecked - c.UnknownSKUs; confirmed > 0 {
				c.SKUAvailabilityPercent = float64(c.AvailableSKUs) / float64(confirmed) * 100
			}
		}
	}
	for _, list := range [][]string{c.MissingResourceTypes, c.MissingSKUs, c.RestrictedSKUs, c.ZoneRestrictedSKUs} {
		sort.Strings(list)
	}
	sort.Slice(result.Health.Warnings, func(i, j int) bool { return result.Health.Warnings[i].Code < result.Health.Warnings[j].Code })
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return result, nil
}
