// Portions of this file reproduce Defender scan behavior from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md.

package defender

import (
	"context"
	"fmt"
	"sort"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const StatusQuery = `
	SecurityResources
	| join kind=inner (
		resourcecontainers
		| where type == 'microsoft.resources/subscriptions'
		| project subscriptionId, subscriptionName = name)
	on subscriptionId
	| where type == 'microsoft.security/pricings'
	| project SubscriptionId = subscriptionId, SubscriptionName = subscriptionName, Name = name, Tier = properties.pricingTier
	`

const RecommendationsQuery = `
	securityresources
	| where type == 'microsoft.security/assessments'
	| where properties.status.code == 'Unhealthy'
	| mvexpand Category = properties.metadata.categories
	| extend
		AssessmentId = id,
		AssessmentKey = name,
		ResourceId = tostring(properties.resourceDetails.Id),
		ResourceIdsplit = split(properties.resourceDetails.Id, '/'),
		RecommendationName = tostring(properties.displayName),
		RecommendationState = tostring(properties.status.code),
		ActionDescription = tostring(properties.metadata.description),
		RemediationDescription = tostring(properties.metadata.remediationDescription),
		RecommendationSeverity = tostring(properties.metadata.severity),
		PolicyDefinitionId = properties.metadata.policyDefinitionId,
		AssessmentType = properties.metadata.assessmentType,
		Threats = properties.metadata.threats,
		UserImpact = properties.metadata.userImpact,
		AzPortalLink = tostring(properties.links.azurePortal),
		CategoryString = tostring(Category)
	| extend
		ResourceSubId = tostring(ResourceIdsplit[2]),
		ResourceGroupName = tostring(ResourceIdsplit[4]),
		ResourceType = tostring(ResourceIdsplit[6]),
		ResourceName = tostring(ResourceIdsplit[8])
	| project SubscriptionId=subscriptionId, ResourceGroupName, ResourceType,
		ResourceName, Category=CategoryString, RecommendationSeverity, RecommendationName, ActionDescription,
		RemediationDescription, AzPortalLink, ResourceId
	| distinct SubscriptionId, ResourceGroupName, ResourceType, ResourceName, Category, RecommendationSeverity, RecommendationName, ActionDescription, RemediationDescription, AzPortalLink, ResourceId
	`

type GraphQuerier interface {
	Query(context.Context, string, map[string]string, ...arg.QueryOptions) (*arg.Result, error)
}

type Filter interface {
	IsSubscriptionExcluded(subscriptionID string) bool
	IsServiceExcluded(resourceID string) bool
}

type StatusResult struct {
	Records  []assessment.DefenderPlanStatus
	Warnings []assessment.AssessmentWarning
}

type RecommendationsResult struct {
	Records  []assessment.DefenderRecommendation
	Warnings []assessment.AssessmentWarning
}

type Scanner struct {
	graph GraphQuerier
}

func New(credential azcore.TokenCredential) *Scanner {
	return NewWithClient(arg.NewClient(arg.NewHTTPTransport(credential)))
}

func NewWithClient(graphClient GraphQuerier) *Scanner {
	return &Scanner{graph: graphClient}
}

func (s *Scanner) ScanStatus(
	ctx context.Context,
	subscriptions map[string]string,
	filter Filter,
) (StatusResult, error) {
	result := StatusResult{
		Records:  []assessment.DefenderPlanStatus{},
		Warnings: []assessment.AssessmentWarning{},
	}
	if s == nil || s.graph == nil {
		return result, fmt.Errorf("Defender Resource Graph client is not configured")
	}

	graphResult, err := s.graph.Query(ctx, StatusQuery, subscriptions)
	if err != nil {
		return result, fmt.Errorf("query Defender status: %w", err)
	}
	if graphResult == nil {
		return result, fmt.Errorf("query Defender status: Resource Graph returned nil result")
	}

	rows, malformed := arg.DecodeRowsWithStats[statusRow](graphResult.Data)
	if malformed > 0 {
		result.Warnings = append(result.Warnings, assessment.AssessmentWarning{
			Code:    "defender_status_malformed_arg_rows",
			Message: fmt.Sprintf("skipped %d malformed Defender status Resource Graph row(s)", malformed),
		})
	}

	for _, row := range rows {
		if filter != nil && filter.IsSubscriptionExcluded(row.SubscriptionID) {
			continue
		}
		result.Records = append(result.Records, assessment.DefenderPlanStatus{
			SubscriptionID:   row.SubscriptionID,
			SubscriptionName: row.SubscriptionName,
			Name:             row.Name,
			Tier:             row.Tier,
		})
	}

	sort.Slice(result.Records, func(i, j int) bool {
		if result.Records[i].SubscriptionID != result.Records[j].SubscriptionID {
			return result.Records[i].SubscriptionID < result.Records[j].SubscriptionID
		}
		if result.Records[i].Name != result.Records[j].Name {
			return result.Records[i].Name < result.Records[j].Name
		}
		return result.Records[i].Tier < result.Records[j].Tier
	})
	return result, nil
}

func (s *Scanner) ScanRecommendations(
	ctx context.Context,
	subscriptions map[string]string,
	filter Filter,
) (RecommendationsResult, error) {
	result := RecommendationsResult{
		Records:  []assessment.DefenderRecommendation{},
		Warnings: []assessment.AssessmentWarning{},
	}
	if s == nil || s.graph == nil {
		return result, fmt.Errorf("Defender Resource Graph client is not configured")
	}

	graphResult, err := s.graph.Query(ctx, RecommendationsQuery, subscriptions)
	if err != nil {
		return result, fmt.Errorf("query Defender recommendations: %w", err)
	}
	if graphResult == nil {
		return result, fmt.Errorf("query Defender recommendations: Resource Graph returned nil result")
	}

	rows, malformed := arg.DecodeRowsWithStats[recommendationRow](graphResult.Data)
	if malformed > 0 {
		result.Warnings = append(result.Warnings, assessment.AssessmentWarning{
			Code:    "defender_recommendations_malformed_arg_rows",
			Message: fmt.Sprintf("skipped %d malformed Defender recommendation Resource Graph row(s)", malformed),
		})
	}

	type recommendationKey struct {
		resourceID         string
		category           string
		recommendationName string
	}
	seen := make(map[recommendationKey]struct{}, len(rows))

	for _, row := range rows {
		if filter != nil && filter.IsServiceExcluded(row.ResourceID) {
			continue
		}

		key := recommendationKey{
			resourceID:         row.ResourceID,
			category:           row.Category,
			recommendationName: row.RecommendationName,
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		result.Records = append(result.Records, assessment.DefenderRecommendation{
			SubscriptionID:         row.SubscriptionID,
			SubscriptionName:       subscriptions[row.SubscriptionID],
			ResourceGroup:          row.ResourceGroupName,
			ResourceType:           row.ResourceType,
			ResourceName:           row.ResourceName,
			Category:               row.Category,
			RecommendationSeverity: row.RecommendationSeverity,
			RecommendationName:     row.RecommendationName,
			ActionDescription:      row.ActionDescription,
			RemediationDescription: row.RemediationDescription,
			AzurePortalLink:        fmt.Sprintf("https://%s", row.AzurePortalLink),
			ResourceID:             row.ResourceID,
		})
	}

	sort.Slice(result.Records, func(i, j int) bool {
		if result.Records[i].SubscriptionID != result.Records[j].SubscriptionID {
			return result.Records[i].SubscriptionID < result.Records[j].SubscriptionID
		}
		if result.Records[i].ResourceID != result.Records[j].ResourceID {
			return result.Records[i].ResourceID < result.Records[j].ResourceID
		}
		if result.Records[i].Category != result.Records[j].Category {
			return result.Records[i].Category < result.Records[j].Category
		}
		return result.Records[i].RecommendationName < result.Records[j].RecommendationName
	})
	return result, nil
}

type statusRow struct {
	SubscriptionID   string `json:"SubscriptionId"`
	SubscriptionName string `json:"SubscriptionName"`
	Name             string `json:"Name"`
	Tier             string `json:"Tier"`
}

type recommendationRow struct {
	SubscriptionID         string `json:"SubscriptionId"`
	ResourceGroupName      string `json:"ResourceGroupName"`
	ResourceType           string `json:"ResourceType"`
	ResourceName           string `json:"ResourceName"`
	Category               string `json:"Category"`
	RecommendationSeverity string `json:"RecommendationSeverity"`
	RecommendationName     string `json:"RecommendationName"`
	ActionDescription      string `json:"ActionDescription"`
	RemediationDescription string `json:"RemediationDescription"`
	AzurePortalLink        string `json:"AzPortalLink"`
	ResourceID             string `json:"ResourceId"`
}
