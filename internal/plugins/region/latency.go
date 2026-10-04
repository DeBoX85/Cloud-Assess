// Copyright (c) Microsoft Corporation.
// Latency arithmetic and pinned data derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_LATENCY.md for provenance and corrections.
package region

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"sort"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const (
	MaxLatencyEntries      = 8192
	MaxLatencyBytes        = 1 << 20
	MaxLatencyMilliseconds = 1000000
	latencySourceBlob      = "64a6abc72e2a54fec286cb27123175d1d6bcda61"
	latencyDataHash        = "45e575040812ee74e006be34623cb253df24e6558f07727cefcef7e55e13326b"
)

//go:embed data/latency.json
var bundledLatencyJSON []byte

type latencyData struct {
	SourceBlob string                        `json:"source_blob"`
	Matrix     map[string]map[string]float64 `json:"matrix"`
	Clusters   map[string]string             `json:"clusters"`
}

type LatencyCalculation struct {
	Comparisons []Comparison
	Health      assessment.StageExecution
}

// EnrichLatency uses the historical pinned P50 snapshot, not live measurements.
// Inputs must remain stable during the call. Output and decoded data are owned
// per call; health describes latency alone and does not replace upstream health.
func EnrichLatency(ctx context.Context, subscriptions map[string]string, input []Comparison) (*LatencyCalculation, error) {
	data, err := decodeLatency(ctx, bundledLatencyJSON)
	if err != nil {
		return nil, err
	}
	return enrichLatency(ctx, subscriptions, input, data)
}

func latencyFailure(ctx context.Context, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("region latency input could not be safely calculated [%s]", code)
}

func decodeLatency(ctx context.Context, raw []byte) (latencyData, error) {
	if err := ctx.Err(); err != nil {
		return latencyData{}, err
	}
	if len(raw) > MaxLatencyBytes {
		return latencyData{}, latencyFailure(ctx, "serialized_limit")
	}
	hash := sha256.Sum256(raw)
	if hex.EncodeToString(hash[:]) != latencyDataHash {
		return latencyData{}, latencyFailure(ctx, "source_invalid")
	}
	var data latencyData
	if err := json.Unmarshal(raw, &data); err != nil {
		return latencyData{}, latencyFailure(ctx, "data_invalid")
	}
	if err := ctx.Err(); err != nil {
		return latencyData{}, err
	}
	return data, nil
}

func enrichLatency(ctx context.Context, subscriptions map[string]string, input []Comparison, data latencyData) (*LatencyCalculation, error) {
	// Reuse the canonical complete comparison admission, including selected
	// identity, detail/UTF16 budgets, duplicates and final projection bounds.
	if _, err := Project(ctx, subscriptions, input); err != nil {
		return nil, latencyFailure(ctx, "comparison_invalid")
	}
	if data.SourceBlob != latencySourceBlob {
		return nil, latencyFailure(ctx, "source_invalid")
	}
	entries, textBytes := 0, 0
	charge := func(n int) bool {
		if n > MaxLatencyEntries-entries {
			return false
		}
		entries += n
		return true
	}
	label := func(s string) bool {
		textBytes += len(s)
		return len(s) <= 64 && textBytes <= MaxLatencyBytes && regionID.MatchString(s)
	}
	if !charge(len(data.Matrix) + len(data.Clusters)) {
		return nil, latencyFailure(ctx, "input_limit")
	}
	for region, cluster := range data.Clusters {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !label(region) || !label(cluster) {
			return nil, latencyFailure(ctx, "label_limit")
		}
	}
	for source, row := range data.Matrix {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !charge(len(row)) {
			return nil, latencyFailure(ctx, "input_limit")
		}
		if !label(source) {
			return nil, latencyFailure(ctx, "label_limit")
		}
		for target, ms := range row {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if !label(target) {
				return nil, latencyFailure(ctx, "label_limit")
			}
			if math.IsNaN(ms) || math.IsInf(ms, 0) || ms < 0 || ms > MaxLatencyMilliseconds {
				return nil, latencyFailure(ctx, "value_invalid")
			}
		}
	}
	// Sum in canonical order. The source's directional membership and every
	// contributing cell (including self cells) are preserved.
	type aggregate struct {
		sum   float64
		count int
	}
	totals := map[string]aggregate{}
	sources := slices.Collect(maps.Keys(data.Matrix))
	sort.Strings(sources)
	for _, source := range sources {
		targets := slices.Collect(maps.Keys(data.Matrix[source]))
		sort.Strings(targets)
		for _, target := range targets {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			a, aKnown := data.Clusters[source]
			b, bKnown := data.Clusters[target]
			if !aKnown || !bKnown {
				continue
			}
			key := a + ":" + b
			value := totals[key]
			value.sum += data.Matrix[source][target]
			value.count++
			if math.IsInf(value.sum, 0) || value.sum > float64(MaxLatencyEntries)*MaxLatencyMilliseconds {
				return nil, latencyFailure(ctx, "sum_invalid")
			}
			totals[key] = value
		}
	}
	result := &LatencyCalculation{Comparisons: make([]Comparison, 0, len(input)), Health: assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: len(input)}}
	warn := func(code string) {
		for _, warning := range result.Health.Warnings {
			if warning.Code == code {
				return
			}
		}
		result.Health.Status = assessment.StageCompletedWithWarnings
		result.Health.Warnings = append(result.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: "region latency evidence is unavailable or estimated"})
	}
	for _, before := range input {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		c := before
		c.MissingResourceTypes = slices.Clone(before.MissingResourceTypes)
		c.MissingSKUs = slices.Clone(before.MissingSKUs)
		c.RestrictedSKUs = slices.Clone(before.RestrictedSKUs)
		c.ZoneRestrictedSKUs = slices.Clone(before.ZoneRestrictedSKUs)
		c.TargetZoneMappings = maps.Clone(before.TargetZoneMappings)
		c.AvgLatencyMs, c.LatencyEstimated = 0, false
		if c.SourceRegion != c.TargetRegion {
			if ms, ok := data.Matrix[c.SourceRegion][c.TargetRegion]; ok {
				c.AvgLatencyMs = ms
			} else if ms, ok := data.Matrix[c.TargetRegion][c.SourceRegion]; ok {
				c.AvgLatencyMs = ms
			} else {
				a, aKnown := data.Clusters[c.SourceRegion]
				b, bKnown := data.Clusters[c.TargetRegion]
				value, exists := totals[a+":"+b]
				if aKnown && bKnown && exists && value.count > 0 && value.sum > 0 {
					c.AvgLatencyMs = value.sum / float64(value.count)
					c.LatencyEstimated = true
					warn("latency_estimated")
				} else {
					warn("latency_unknown")
				}
			}
		}
		result.Comparisons = append(result.Comparisons, c)
	}
	sort.Slice(result.Health.Warnings, func(i, j int) bool { return result.Health.Warnings[i].Code < result.Health.Warnings[j].Code })
	// Changed numeric text must also satisfy the complete output table budget.
	if _, err := Project(ctx, subscriptions, result.Comparisons); err != nil {
		return nil, latencyFailure(ctx, "output_invalid")
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return result, nil
}
