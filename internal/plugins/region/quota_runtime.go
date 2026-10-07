// Copyright (c) Microsoft Corporation.
// Quota arithmetic and provider filters derived from MIT-licensed AZQR.
// See NOTICE.md and docs/REGION_QUOTA_RUNTIME.md for corrections and bounds.
package region

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type QuotaRequest struct {
	SubscriptionID, Region, QuotaType string
}

type QuotaUsage struct {
	ResourceName, LocalizedName string
	Current, Limit              int64
	CurrentKnown, LimitKnown    bool
}

// QuotaEvidence.Status is complete, partial, unknown or unsupported.
type QuotaEvidence struct {
	Request QuotaRequest
	Status  string
	Usages  []QuotaUsage
}

type QuotaCalculation struct {
	Rows   []QuotaRow
	Health assessment.StageExecution
	Table  *assessment.PluginTable
}

func quotaFailure(ctx context.Context, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("region quota evidence could not be safely calculated [%s]", code)
}

func quotaProvider(provider string) bool {
	switch provider {
	case "VM", "Network", "SQL", "App Service", "Storage":
		return true
	}
	return false
}

func quotaFiltered(provider, resource string) bool {
	switch provider {
	case "VM":
		return !strings.Contains(resource, "Family")
	case "Network":
		return slices.Contains([]string{"NetworkWatchers", "RouteFilterRulesPerRouteFilter", "RouteFiltersPerExpressRouteBgpPeering", "RoutesPerExpressRouteCircuit", "BgpCommunityFilterRulesPerRouteFilter"}, resource)
	case "SQL":
		return strings.HasSuffix(resource, "PerServer") || strings.HasSuffix(resource, "PerDatabase")
	case "App Service":
		return slices.Contains([]string{"CustomDomains", "HostNameBindings", "SslBindings", "SslConnections", "Certificates"}, resource)
	case "Storage":
		return slices.Contains([]string{"TotalBlobContainers", "TotalBlobs", "TotalContainers", "TotalFileShares", "TotalQueues", "TotalTables"}, resource)
	}
	return false
}

// CalculateQuota performs no collection. Inputs must remain stable during this
// call. It validates all supplied evidence before allocating output rows and
// returns independently owned calculation and table health, even for empty data.
func CalculateQuota(ctx context.Context, subscriptions map[string]string, requests []QuotaRequest, evidence []QuotaEvidence) (*QuotaCalculation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(requests) > MaxAuxRows || len(evidence) > len(requests) {
		return nil, quotaFailure(ctx, "work_limit")
	}
	textBytes := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &textBytes)
	if !ok {
		return nil, quotaFailure(ctx, "scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return nil, quotaFailure(ctx, "scope_invalid")
		}
	}
	outputText := textBytes
	normalize := func(request QuotaRequest) (QuotaRequest, bool) {
		if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || !regionID.MatchString(request.Region) || !quotaProvider(request.QuotaType) {
			return QuotaRequest{}, false
		}
		request.SubscriptionID = strings.ToLower(request.SubscriptionID)
		_, selected := scope[request.SubscriptionID]
		return request, selected
	}
	expected := make(map[QuotaRequest]bool, len(requests))
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key, valid := normalize(request)
		if !valid || expected[key] {
			return nil, quotaFailure(ctx, "request_invalid")
		}
		for _, label := range []string{request.SubscriptionID, request.Region, request.QuotaType} {
			if !auxiliaryText(label, &textBytes) {
				return nil, quotaFailure(ctx, "text_limit")
			}
		}
		expected[key] = true
	}
	observed := make(map[QuotaRequest]int, len(evidence))
	count := 0
	for i, response := range evidence {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key, valid := normalize(response.Request)
		_, duplicate := observed[key]
		if !valid || !expected[key] || duplicate {
			return nil, quotaFailure(ctx, "evidence_invalid")
		}
		if response.Status != "complete" && response.Status != "partial" && response.Status != "unknown" && response.Status != "unsupported" || (response.Status == "unknown" || response.Status == "unsupported") && len(response.Usages) != 0 {
			return nil, quotaFailure(ctx, "status_invalid")
		}
		if len(response.Usages) > MaxAuxRows-count {
			return nil, quotaFailure(ctx, "work_limit")
		}
		count += len(response.Usages)
		seen := make(map[string]bool, len(response.Usages))
		for _, usage := range response.Usages {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if !auxiliaryText(usage.ResourceName, &textBytes) || !auxiliaryText(usage.LocalizedName, &textBytes) {
				return nil, quotaFailure(ctx, "text_limit")
			}
			if usage.Current < 0 || usage.Current > MaxAuxCount || usage.Limit < -MaxAuxCount || usage.Limit > MaxAuxCount || !usage.CurrentKnown && usage.Current != 0 || !usage.LimitKnown && usage.Limit != 0 {
				return nil, quotaFailure(ctx, "count_invalid")
			}
			if usage.ResourceName != "" {
				if seen[usage.ResourceName] {
					return nil, quotaFailure(ctx, "identity_duplicate")
				}
				seen[usage.ResourceName] = true
			}
			if usage.ResourceName != "" && usage.CurrentKnown && usage.LimitKnown && usage.Limit > 0 && !quotaFiltered(key.QuotaType, usage.ResourceName) {
				pct := float64(usage.Limit-usage.Current) / float64(usage.Limit) * 100
				if math.Abs(pct) > float64(MaxAuxCount) {
					return nil, quotaFailure(ctx, "ratio_limit")
				}
				label := usage.LocalizedName
				if label == "" {
					label = usage.ResourceName
				}
				for _, text := range []string{scope[key.SubscriptionID], key.Region, key.QuotaType, usage.ResourceName, label} {
					if !auxiliaryText(text, &outputText) {
						return nil, quotaFailure(ctx, "text_limit")
					}
				}
				outputText += 128
				if outputText > auxTextBudget {
					return nil, quotaFailure(ctx, "text_limit")
				}
			}
		}
		observed[key] = i
	}
	result := &QuotaCalculation{Rows: make([]QuotaRow, 0, count), Health: assessment.StageExecution{Name: Name, Status: assessment.StageCompleted}}
	warn := func(code string) {
		for _, warning := range result.Health.Warnings {
			if warning.Code == code {
				return
			}
		}
		result.Health.Status = assessment.StageCompletedWithWarnings
		result.Health.Warnings = append(result.Health.Warnings, assessment.AssessmentWarning{Code: code, Message: "region quota evidence is incomplete or unusable"})
	}
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key, _ := normalize(request)
		i, exists := observed[key]
		if !exists {
			warn("quota_evidence_missing")
			continue
		}
		response := evidence[i]
		if response.Status != "complete" {
			warn("quota_evidence_" + response.Status)
		}
		for _, usage := range response.Usages {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if usage.ResourceName == "" || !usage.CurrentKnown || !usage.LimitKnown {
				warn("quota_usage_unknown")
				continue
			}
			if quotaFiltered(key.QuotaType, usage.ResourceName) {
				continue
			}
			if usage.Limit <= 0 {
				warn("quota_limit_unusable")
				continue
			}
			available := usage.Limit - usage.Current
			pct := float64(available) / float64(usage.Limit) * 100
			label := usage.LocalizedName
			if label == "" {
				label = usage.ResourceName
			}
			result.Rows = append(result.Rows, QuotaRow{SubscriptionID: key.SubscriptionID, Subscription: scope[key.SubscriptionID], Region: key.Region, QuotaType: key.QuotaType, ResourceName: usage.ResourceName, DisplayName: label, Current: usage.Current, Limit: usage.Limit, Available: available, HeadroomPct: pct, IsNearLimit: available*100 < usage.Limit*15, IsOverLimit: available <= 0})
		}
	}
	result.Health.Records = len(result.Rows)
	table, err := ProjectQuota(ctx, subscriptions, result.Rows)
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
