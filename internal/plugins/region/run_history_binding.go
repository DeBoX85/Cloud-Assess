// Owned history admission binds declared scope and currency without new I/O.
// See docs/RUN_HISTORY_BINDING_CONTRACT_20261029.md for limits and provenance.
package region

import (
	"context"
	"errors"
	"math"
	"strings"
)

type RunHistoryEvidence struct {
	Admitted            bool
	Request             MeterHistoryRequest
	Status, FailureCode string
	ObservedMeters      int
	History             CostHistoryEvidence
}

// Exported fields are diagnostics; only the private snapshot admits evidence.
type RunHistoryBinding struct {
	Request             MeterHistoryRequest
	Status, FailureCode string
	Meters              []HistoricalCostMeter
	state               *runHistoryState
}
type runHistoryState struct{ evidence RunHistoryEvidence }

func runHistoryContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}
func runHistoryInvalid(ctx context.Context) error {
	if e := runHistoryContext(ctx); e != nil {
		return e
	}
	return errors.New("region run history binding rejected [run_history_invalid]")
}
func runHistoryFailureCode(code string) bool {
	switch code {
	case "history_limit", "history_continuation_invalid", "history_response_invalid", "history_request_failed", "history_decode_invalid":
		return true
	}
	return false
}

// AdmitRunHistory snapshots stable inputs. It establishes structural admission,
// not authenticated service provenance or proof of the original query window.
func AdmitRunHistory(ctx context.Context, selected map[string]string, q MeterHistoryRequest, in *MeterHistoryCollection) (*RunHistoryBinding, error) {
	if e := runHistoryContext(ctx); e != nil {
		return nil, e
	}
	if ctx == nil {
		return nil, runHistoryInvalid(ctx)
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, selected, &text)
	if !ok {
		return nil, runHistoryInvalid(ctx)
	}
	for _, name := range scope {
		if e := runHistoryContext(ctx); e != nil {
			return nil, e
		}
		if name == "" {
			return nil, runHistoryInvalid(ctx)
		}
	}
	if !historyRequestValid(q) {
		return nil, runHistoryInvalid(ctx)
	}
	q.SubscriptionID = strings.ToLower(q.SubscriptionID)
	q.Now = q.Now.UTC()
	if scope[q.SubscriptionID] == "" || in == nil || len(in.SubscriptionID) != 36 || !subscriptionID.MatchString(in.SubscriptionID) || !strings.EqualFold(in.SubscriptionID, q.SubscriptionID) || in.Currency != q.Currency {
		return nil, runHistoryInvalid(ctx)
	}
	switch in.Status {
	case "complete":
		if in.FailureCode != "" {
			return nil, runHistoryInvalid(ctx)
		}
	case "partial", "unknown":
		if !runHistoryFailureCode(in.FailureCode) || (in.Status == "unknown" && len(in.Meters) != 0) {
			return nil, runHistoryInvalid(ctx)
		}
	default:
		return nil, runHistoryInvalid(ctx)
	}
	if len(in.Meters) > MaxAuxRows {
		return nil, runHistoryInvalid(ctx)
	}
	seen := make(map[string]bool, len(in.Meters))
	owned := make([]HistoricalCostMeter, 0, len(in.Meters))
	diagnostics := make([]HistoricalCostMeter, 0, len(in.Meters))
	for _, meter := range in.Meters {
		if e := runHistoryContext(ctx); e != nil {
			return nil, e
		}
		if meter.MeterID == "" || !auxiliaryText(meter.MeterID, &text) || seen[meter.MeterID] {
			return nil, runHistoryInvalid(ctx)
		}
		weight := meter.HistoricalCost
		if math.IsNaN(weight) || math.IsInf(weight, 0) || weight < 0 || weight > MaxCostValue || (weight > 0 && weight < MinCostWeight) {
			return nil, runHistoryInvalid(ctx)
		}
		seen[meter.MeterID] = true
		owned = append(owned, meter)
		diagnostics = append(diagnostics, meter)
	}
	history := CostHistoryEvidence{Complete: in.Status == "complete", Meters: []HistoricalCostMeter{}}
	if history.Complete {
		history.Meters = owned
	}
	evidence := RunHistoryEvidence{Admitted: true, Request: q, Status: in.Status, FailureCode: in.FailureCode, ObservedMeters: len(in.Meters), History: history}
	result := &RunHistoryBinding{Request: q, Status: in.Status, FailureCode: in.FailureCode, Meters: diagnostics, state: &runHistoryState{evidence: evidence}}
	if e := runHistoryContext(ctx); e != nil {
		return nil, e
	}
	return result, nil
}

// Evidence returns a fresh owned cost projection and never reads diagnostics.
func (r *RunHistoryBinding) Evidence() RunHistoryEvidence {
	if r == nil || r.state == nil {
		return RunHistoryEvidence{Status: "unknown", History: CostHistoryEvidence{Meters: []HistoricalCostMeter{}}}
	}
	out := r.state.evidence
	out.History.Meters = append([]HistoricalCostMeter{}, r.state.evidence.History.Meters...)
	return out
}
