// Copyright (c) Microsoft Corporation.
// Reservation arithmetic/status derived from MIT-licensed AZQR.
// See NOTICE.md and docs/REGION_RESERVATION_RUNTIME.md for corrections.
package region

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const MaxReservationIDBytes = 2048

type ReservationRequest struct{ SubscriptionID, Region string }

type ReservationUsage struct {
	ResourceID, Region, ResponseName, ResponseRegion, SKU string
	Reserved, Allocated                                   int64
	ReservedKnown, AllocatedKnown                         bool
}

// Status is complete, partial, unknown or unsupported.
type ReservationEvidence struct {
	Request      ReservationRequest
	Status       string
	Reservations []ReservationUsage
}

type ReservationValue struct {
	ResourceID, SubscriptionID, Subscription, Region string
	ResourceGroup, GroupName, ReservationName, SKU   string
	Reserved, Allocated, Available                   int64
	Status                                           string
}

type ReservationCalculation struct {
	Reservations []ReservationValue
	Health       assessment.StageExecution
	Table        *assessment.PluginTable
}

func reservationFailure(ctx context.Context, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("region reservation evidence could not be safely calculated [%s]", code)
}

// Canonicalize each Unicode simple-fold orbit, matching strings.EqualFold.
// Lowercasing alone does not join sigma and final-sigma resource labels.
func reservationFoldKey(id string) string {
	return strings.Map(func(r rune) rune {
		minimum := r
		for next := unicode.SimpleFold(r); next != r; next = unicode.SimpleFold(next) {
			if next < minimum {
				minimum = next
			}
		}
		return minimum
	}, id)
}

// This parses structural identity, not all Azure service naming constraints.
func reservationIdentity(id string) (subscription, rg, group, name string, ok bool) {
	if len(id) > MaxReservationIDBytes || !strings.HasPrefix(id, "/") || strings.ContainsAny(id, "\\?#%") {
		return
	}
	parts := strings.Split(id[1:], "/")
	if len(parts) != 10 {
		return
	}
	for index, fixed := range map[int]string{0: "subscriptions", 2: "resourceGroups", 4: "providers", 5: "Microsoft.Compute", 6: "capacityReservationGroups", 8: "capacityReservations"} {
		if len(parts[index]) != len(fixed) || !strings.EqualFold(parts[index], fixed) {
			return
		}
	}
	if len(parts[1]) != 36 || !subscriptionID.MatchString(parts[1]) {
		return
	}
	for _, index := range []int{3, 7, 9} {
		budget := 0
		if parts[index] == "" || parts[index] == "." || parts[index] == ".." || !auxiliaryText(parts[index], &budget) {
			return
		}
	}
	return strings.ToLower(parts[1]), parts[3], parts[7], parts[9], true
}

// CalculateReservations performs no collection. Inputs must remain stable.
// All raw admission and independent output budgets precede output allocation.
func CalculateReservations(ctx context.Context, subscriptions map[string]string, requests []ReservationRequest, evidence []ReservationEvidence) (*ReservationCalculation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(requests) > MaxAuxRows || len(evidence) > len(requests) {
		return nil, reservationFailure(ctx, "work_limit")
	}
	textBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &textBytes)
	if !ok {
		return nil, reservationFailure(ctx, "scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return nil, reservationFailure(ctx, "scope_invalid")
		}
	}
	outputText := textBytes
	normalize := func(request ReservationRequest) (ReservationRequest, bool) {
		if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || !regionID.MatchString(request.Region) {
			return ReservationRequest{}, false
		}
		request.SubscriptionID = strings.ToLower(request.SubscriptionID)
		_, selected := scope[request.SubscriptionID]
		return request, selected
	}
	expected := make(map[ReservationRequest]bool, len(requests))
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key, valid := normalize(request)
		if !valid || expected[key] {
			return nil, reservationFailure(ctx, "request_invalid")
		}
		for _, label := range []string{request.SubscriptionID, request.Region} {
			if !auxiliaryText(label, &textBytes) {
				return nil, reservationFailure(ctx, "text_limit")
			}
		}
		expected[key] = true
	}
	observed := make(map[ReservationRequest]int, len(evidence))
	seen := make(map[string]bool)
	count := 0
	for i, response := range evidence {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key, valid := normalize(response.Request)
		_, duplicate := observed[key]
		if !valid || !expected[key] || duplicate {
			return nil, reservationFailure(ctx, "evidence_invalid")
		}
		if response.Status != "complete" && response.Status != "partial" && response.Status != "unknown" && response.Status != "unsupported" || (response.Status == "unknown" || response.Status == "unsupported") && len(response.Reservations) != 0 {
			return nil, reservationFailure(ctx, "status_invalid")
		}
		if len(response.Reservations) > MaxAuxRows-count {
			return nil, reservationFailure(ctx, "work_limit")
		}
		count += len(response.Reservations)
		for _, usage := range response.Reservations {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			sub, rg, group, name, valid := reservationIdentity(usage.ResourceID)
			if !valid || sub != key.SubscriptionID || usage.Region != key.Region || !regionID.MatchString(usage.Region) || usage.ResponseName != "" && !strings.EqualFold(usage.ResponseName, name) || usage.ResponseRegion != "" && (usage.ResponseRegion != usage.Region || !regionID.MatchString(usage.ResponseRegion)) {
				return nil, reservationFailure(ctx, "identity_invalid")
			}
			canonical := reservationFoldKey(usage.ResourceID)
			if seen[canonical] {
				return nil, reservationFailure(ctx, "identity_duplicate")
			}
			seen[canonical] = true
			textBytes += len(usage.ResourceID)
			if textBytes > auxTextBudget {
				return nil, reservationFailure(ctx, "text_limit")
			}
			for _, label := range []string{usage.Region, usage.ResponseName, usage.ResponseRegion, usage.SKU} {
				if !auxiliaryText(label, &textBytes) {
					return nil, reservationFailure(ctx, "text_limit")
				}
			}
			if usage.Reserved < 0 || usage.Reserved > MaxAuxCount || usage.Allocated < 0 || usage.Allocated > MaxAuxCount || !usage.ReservedKnown && usage.Reserved != 0 || !usage.AllocatedKnown && usage.Allocated != 0 {
				return nil, reservationFailure(ctx, "count_invalid")
			}
			if usage.ReservedKnown && usage.AllocatedKnown && usage.SKU != "" {
				if usage.ResponseName != "" {
					name = usage.ResponseName
				}
				for _, label := range []string{scope[sub], usage.Region, rg, group, name, usage.SKU} {
					if !auxiliaryText(label, &outputText) {
						return nil, reservationFailure(ctx, "output_limit")
					}
				}
				outputText += 128
				if outputText > auxTextBudget {
					return nil, reservationFailure(ctx, "output_limit")
				}
			}
		}
		observed[key] = i
	}
	result := &ReservationCalculation{Reservations: make([]ReservationValue, 0, count), Health: assessment.StageExecution{Name: Name, Status: assessment.StageCompleted}}
	warn := func(code string) {
		for _, warning := range result.Health.Warnings {
			if warning.Code == code {
				return
			}
		}
		result.Health.Status = assessment.StageCompletedWithWarnings
		result.Health.Warnings = append(result.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: "region reservation evidence is incomplete or unusable"})
	}
	rows := make([]ReservationRow, 0, count)
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key, _ := normalize(request)
		i, exists := observed[key]
		if !exists {
			warn("reservation_evidence_missing")
			continue
		}
		response := evidence[i]
		if response.Status != "complete" {
			warn("reservation_evidence_" + response.Status)
		}
		for _, usage := range response.Reservations {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if !usage.ReservedKnown || !usage.AllocatedKnown {
				warn("reservation_usage_unknown")
				continue
			}
			if usage.SKU == "" {
				warn("reservation_sku_unknown")
				continue
			}
			sub, rg, group, name, _ := reservationIdentity(usage.ResourceID)
			if usage.ResponseName != "" {
				name = usage.ResponseName
			}
			available := usage.Reserved - usage.Allocated
			status := "Available"
			switch {
			case usage.Allocated == 0:
				status = "Idle"
			case available < 0:
				status = "Over-Allocated"
			case available == 0:
				status = "At-Capacity"
			}
			value := ReservationValue{ResourceID: strings.ToLower(usage.ResourceID), SubscriptionID: sub, Subscription: scope[sub], Region: usage.Region, ResourceGroup: rg, GroupName: group, ReservationName: name, SKU: usage.SKU, Reserved: usage.Reserved, Allocated: usage.Allocated, Available: available, Status: status}
			result.Reservations = append(result.Reservations, value)
			rows = append(rows, ReservationRow{SubscriptionID: sub, Cells: []string{value.Subscription, value.Region, rg, group, name, value.SKU, strconv.FormatInt(value.Reserved, 10), strconv.FormatInt(value.Allocated, 10), strconv.FormatInt(value.Available, 10), value.Status}})
		}
	}
	result.Health.Records = len(result.Reservations)
	table, err := ProjectReservations(ctx, subscriptions, rows)
	if err != nil {
		return nil, err
	}
	if table != nil {
		table.Health = result.Health
		table.Health.Warnings = slices.Clone(result.Health.Warnings)
		result.Table = table
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
