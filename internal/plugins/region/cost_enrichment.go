// Copyright (c) Microsoft Corporation.
// Weighted arithmetic derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_COST_RUNTIME.md for source and corrections.
package region

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxCostWork = 1048576
const MaxCostValue = 1000000000000

type HistoricalCostMeter struct {
	MeterID string

	HistoricalCost float64
}

type CostHistoryEvidence struct {
	Complete bool

	Meters []HistoricalCostMeter
}

// Evidence declares completeness; the future collector must establish actual
// resource ownership, currency, unit/tier, clock and pagination correctness.
type CostEnrichmentEvidence struct {
	History map[string]CostHistoryEvidence

	PricingStatus string

	RegionPricing map[string]map[string]float64
}

type CostCalculation struct {
	Comparisons []Comparison

	Health assessment.StageExecution
}

func costEnrichmentFailure(ctx context.Context, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("region cost input could not be safely calculated [%s]", code)
}

// costPhysical retains the pinned IsPhysicalRegion exclusions, not a registry
// of currently available Azure locations. Format admission happens separately.
func costPhysical(region string) bool {
	switch region {
	case "", "unassigned", "global", "europe", "unitedstates", "asia", "asiapacific", "australia", "brazil", "canada", "france", "germany", "india", "japan", "korea", "norway", "southafrica", "switzerland", "uae", "uk":
		return false
	}
	return true
}

// EnrichCost owns results and health per call. Inputs must remain stable during
// the call. No evidence or caller-owned detail maps/slices escape. Health is
// cost-only; it cannot replace availability, latency or collector health.
func EnrichCost(ctx context.Context, subscriptions map[string]string, input []Comparison, evidence *CostEnrichmentEvidence) (*CostCalculation, error) {
	fail := func(code string) (*CostCalculation, error) { return nil, costEnrichmentFailure(ctx, code) }
	if _, err := Project(ctx, subscriptions, input); err != nil {
		return fail("comparison_invalid")
	}
	for _, c := range input {
		if c.SubscriptionID != strings.ToLower(c.SubscriptionID) {
			return fail("comparison_invalid")
		}
	}
	decodedBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &decodedBytes)
	if !ok {
		return fail("scope_invalid")
	}
	if evidence == nil {
		evidence = &CostEnrichmentEvidence{PricingStatus: "unavailable"}
	}
	switch evidence.PricingStatus {
	case "complete", "partial":
	case "unavailable", "unsupported":
		if len(evidence.RegionPricing) != 0 {
			return fail("status_invalid")
		}
	default:
		return fail("status_invalid")
	}
	entries, meters := 0, 0
	charge := func(n int) bool {
		if n > MaxCostEntries-entries {
			return false
		}
		entries += n
		return true
	}
	text := func(s string) bool { return auxiliaryText(s, &decodedBytes) }
	validValue := func(n float64) bool { return !math.IsNaN(n) && !math.IsInf(n, 0) && n >= 0 && n <= MaxCostValue }
	if len(evidence.History) > MaxSubscriptions || len(evidence.RegionPricing) > MaxCostMeters || !charge(len(evidence.History)) || !charge(len(evidence.RegionPricing)) {
		return fail("input_limit")
	}
	// Validate all evidence, including unconsumed prices, before allocating
	// owned indexes or results. Receipt order cannot silently pick duplicates.
	for id, history := range evidence.History {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if _, selected := scope[id]; !selected || id != strings.ToLower(id) || len(id) != 36 || !subscriptionID.MatchString(id) {
			return fail("evidence_scope_invalid")
		}
		if len(history.Meters) > MaxCostMeters-meters || !charge(len(history.Meters)) {
			return fail("input_limit")
		}
		meters += len(history.Meters)
		seen := map[string]bool{}
		for _, meter := range history.Meters {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			if meter.MeterID == "" || !text(meter.MeterID) {
				return fail("text_limit")
			}
			if seen[meter.MeterID] {
				return fail("duplicate_meter")
			}
			seen[meter.MeterID] = true
			if !validValue(meter.HistoricalCost) {
				return fail("value_invalid")
			}
		}
	}
	regions := map[string]bool{}
	for id, prices := range evidence.RegionPricing {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if !charge(len(prices)) {
			return fail("input_limit")
		}
		if id == "" || !text(id) {
			return fail("text_limit")
		}
		for region, price := range prices {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			if len(region) > 64 || !regionID.MatchString(region) || !text(region) {
				return fail("region_invalid")
			}
			if !validValue(price) {
				return fail("value_invalid")
			}
			regions[region] = true
			if len(regions) > MaxCostRegions {
				return fail("input_limit")
			}
		}
	}
	work := 0
	for _, c := range input {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		n := len(evidence.History[c.SubscriptionID].Meters)
		if n > MaxCostWork-work {
			return fail("work_limit")
		}
		work += n
	}
	ordered := map[string][]HistoricalCostMeter{}
	for id, history := range evidence.History {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		ordered[id] = slices.Clone(history.Meters)
		sort.Slice(ordered[id], func(i, j int) bool { return ordered[id][i].MeterID < ordered[id][j].MeterID })
	}
	result := &CostCalculation{Comparisons: make([]Comparison, 0, len(input)), Health: assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(input)}}
	warn := func(code string) {
		for _, warning := range result.Health.Warnings {
			if warning.Code == code {
				return
			}
		}
		result.Health.Status = assessment.StageCompletedWithWarnings
		result.Health.Warnings = append(result.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: "region cost evidence is incomplete or ineligible"})
	}
	for _, before := range input {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		c := before
		c.MissingResourceTypes = slices.Clone(before.MissingResourceTypes)
		c.MissingSKUs = slices.Clone(before.MissingSKUs)
		c.RestrictedSKUs = slices.Clone(before.RestrictedSKUs)
		c.ZoneRestrictedSKUs = slices.Clone(before.ZoneRestrictedSKUs)
		c.TargetZoneMappings = maps.Clone(before.TargetZoneMappings)
		c.AvgCostDifference, c.HasCostData = 0, false
		history, declared := evidence.History[c.SubscriptionID]
		eligible := true
		if !declared {
			warn("cost_history_unavailable")
			eligible = false
		} else if !history.Complete {
			warn("cost_history_partial")
			eligible = false
		}
		switch evidence.PricingStatus {
		case "partial":
			warn("cost_pricing_partial")
		case "unavailable", "unsupported":
			warn("cost_pricing_" + evidence.PricingStatus)
			eligible = false
		}
		if !costPhysical(c.SourceRegion) || !costPhysical(c.TargetRegion) {
			warn("cost_region_unsupported")
			eligible = false
		}
		weighted, total := 0.0, 0.0
		if eligible {
			for _, meter := range ordered[c.SubscriptionID] {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				weight := meter.HistoricalCost
				if weight == 0 {
					weight = 1
				}
				prices := evidence.RegionPricing[meter.MeterID]
				source, hasSource := prices[c.SourceRegion]
				target, hasTarget := prices[c.TargetRegion]
				if !hasSource || !hasTarget {
					warn("cost_price_missing")
					continue
				}
				diff := 0.0
				if source < 0.0001 && target < 0.0001 {
					// Free pairs still contribute denominator weight.
				} else if source < 0.0001 {
					warn("cost_meter_ineligible")
					continue
				} else {
					diff = (target - source) / source * 100
				}
				weighted += diff * weight
				total += weight
				if math.IsNaN(weighted) || math.IsInf(weighted, 0) || math.IsNaN(total) || math.IsInf(total, 0) {
					return fail("sum_invalid")
				}
			}
			if total > 0 {
				// Nonnegative prices/weights prove the exact lower domain;
				// IEEE754 accumulation can round an all-free mean below it.
				c.AvgCostDifference = math.Max(-100, weighted/total)
				c.HasCostData = true
			} else {
				warn("cost_no_eligible_meters")
			}
		}
		result.Comparisons = append(result.Comparisons, c)
	}
	sort.Slice(result.Health.Warnings, func(i, j int) bool { return result.Health.Warnings[i].Code < result.Health.Warnings[j].Code })
	if _, err := Project(ctx, subscriptions, result.Comparisons); err != nil {
		return fail("output_invalid")
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return result, nil
}
