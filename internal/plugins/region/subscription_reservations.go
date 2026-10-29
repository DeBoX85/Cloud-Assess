// Subscription traversal derives from the MIT-licensed AZQR CRG contract.
// See NOTICE.md and docs/SUBSCRIPTION_RESERVATION_CONTRACT_20261028.md.
package region

import (
	"context"
	"errors"
)

type SubscriptionReservationRequest struct{ SubscriptionID string }

// Status covers the entire subscription traversal, including empty prefixes.
// Regional values remain raw scoped evidence, not a capacity/eligibility claim.
type SubscriptionReservationEvidence struct {
	SubscriptionID, Status, FailureCode string
	Requests                            []ReservationRequest
	Regional                            []ReservationEvidence
}

// Exported fields are diagnostics. Evidence uses a separate owned snapshot.
type SubscriptionReservationCollection struct {
	SubscriptionID, Status, FailureCode string
	Regions                             []string
	Reservations                        []ReservationUsage
	state                               *subscriptionReservationState
}

type subscriptionReservationState struct {
	evidence SubscriptionReservationEvidence
}

type subscriptionReservationTraversal struct {
	regions []string
	seen    map[string]bool
}

func (s *subscriptionReservationTraversal) observe(region string) {
	if !s.seen[region] {
		s.seen[region] = true
		s.regions = append(s.regions, region)
	}
}

func subscriptionReservationCopy(in SubscriptionReservationEvidence) SubscriptionReservationEvidence {
	out := SubscriptionReservationEvidence{SubscriptionID: in.SubscriptionID, Status: in.Status, FailureCode: in.FailureCode, Requests: append([]ReservationRequest{}, in.Requests...), Regional: make([]ReservationEvidence, len(in.Regional))}
	for i, regional := range in.Regional {
		out.Regional[i] = ReservationEvidence{Request: regional.Request, Status: regional.Status, Reservations: append([]ReservationUsage{}, regional.Reservations...)}
	}
	return out
}

// Evidence never promotes exported diagnostic mutations or caller-built values.
// Every mutable slice belongs to this invocation, including empty collections.
func (c *SubscriptionReservationCollection) Evidence() SubscriptionReservationEvidence {
	if c == nil || c.state == nil {
		return SubscriptionReservationEvidence{Status: "unknown", Requests: []ReservationRequest{}, Regional: []ReservationEvidence{}}
	}
	return subscriptionReservationCopy(c.state.evidence)
}

// CollectSubscription traverses all admitted groups once with the existing
// global request/body/work limits. There is no context/cache retained on c.
func (c *ReservationCollector) CollectSubscription(ctx context.Context, selected map[string]string, request SubscriptionReservationRequest) (*SubscriptionReservationCollection, error) {
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if ctx == nil || c == nil || c.transport == nil || c.transport.origin == nil || c.transport.getter == nil {
		return nil, errors.New("region reservation evidence could not be safely calculated [request_invalid]")
	}
	scan := &subscriptionReservationTraversal{regions: []string{}, seen: map[string]bool{}}
	collected, err := c.collectReservationTraversal(ctx, selected, ReservationRequest{SubscriptionID: request.SubscriptionID}, scan)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// Strip any surrounding dependency text from context errors. The
		// regional method deliberately retains its existing error contract.
		if errors.Is(err, context.Canceled) {
			return nil, context.Canceled
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	evidence := SubscriptionReservationEvidence{SubscriptionID: collected.Evidence.Request.SubscriptionID, Status: collected.Evidence.Status, FailureCode: collected.FailureCode, Requests: []ReservationRequest{}, Regional: []ReservationEvidence{}}
	indices := make(map[string]int, len(scan.regions))
	for _, region := range scan.regions {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		q := ReservationRequest{SubscriptionID: evidence.SubscriptionID, Region: region}
		indices[region] = len(evidence.Regional)
		evidence.Requests = append(evidence.Requests, q)
		evidence.Regional = append(evidence.Regional, ReservationEvidence{Request: q, Status: evidence.Status, Reservations: []ReservationUsage{}})
	}
	for _, row := range collected.Evidence.Reservations {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		index, known := indices[row.Region]
		if !known {
			return nil, reservationFailure(ctx, "reservation_identity_invalid")
		}
		evidence.Regional[index].Reservations = append(evidence.Regional[index].Reservations, row)
	}
	result := &SubscriptionReservationCollection{SubscriptionID: evidence.SubscriptionID, Status: evidence.Status, FailureCode: evidence.FailureCode, Regions: append([]string{}, scan.regions...), Reservations: append([]ReservationUsage{}, collected.Evidence.Reservations...), state: &subscriptionReservationState{evidence: evidence}}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return result, nil
}
