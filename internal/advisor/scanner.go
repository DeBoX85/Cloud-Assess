// Portions of this file reproduce Advisor scan behavior from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md.

package advisor

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const Query = `
	AdvisorResources
	| where type =~ 'microsoft.advisor/recommendations'
	| where isnotempty(properties.resourceMetadata.resourceId)
	| where isnull(properties.suppressionIds) or array_length(properties.suppressionIds) == 0
	| project SubscriptionId=subscriptionId,
		Category = tostring(properties.category),
		Impact = tostring(properties.impact),
		ImpactedValue = tostring(properties.impactedValue),
		ResourceId = tostring(properties.resourceMetadata.resourceId),
		RecommendationTypeId = tostring(properties.recommendationTypeId)
	| summarize take_any(*) by ResourceId, RecommendationTypeId
	`

type GraphQuerier interface {
	Query(context.Context, string, map[string]string, ...arg.QueryOptions) (*arg.Result, error)
}

type ResourceFilter interface {
	IsSubscriptionExcluded(subscriptionID string) bool
	IsServiceExcluded(resourceID string) bool
}

type Result struct {
	Records  []assessment.AdvisorRecommendation
	Warnings []assessment.AssessmentWarning
}

type Scanner struct {
	graph    GraphQuerier
	metadata RecommendationTypeProvider
}

func New(credential azcore.TokenCredential) *Scanner {
	graphClient := arg.NewClient(arg.NewHTTPTransport(credential))
	httpClient := azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(30*time.Second))
	metadataClient := NewMetadataClient(httpClient, azure.ResourceManagerEndpoint())
	return NewWithClients(graphClient, metadataClient)
}

func NewWithClients(graphClient GraphQuerier, metadataClient RecommendationTypeProvider) *Scanner {
	return &Scanner{graph: graphClient, metadata: metadataClient}
}

func (s *Scanner) Scan(
	ctx context.Context,
	subscriptions map[string]string,
	filter ResourceFilter,
) (Result, error) {
	result := Result{Records: []assessment.AdvisorRecommendation{}, Warnings: []assessment.AssessmentWarning{}}
	if s == nil || s.metadata == nil {
		return result, fmt.Errorf("Advisor metadata provider is not configured")
	}
	if s.graph == nil {
		return result, fmt.Errorf("Advisor Resource Graph client is not configured")
	}

	recommendationTypes, err := s.metadata.RecommendationTypes(ctx)
	if err != nil {
		return result, err
	}

	graphResult, err := s.graph.Query(ctx, Query, subscriptions)
	if err != nil {
		return result, fmt.Errorf("query Advisor recommendations: %w", err)
	}
	if graphResult == nil {
		return result, fmt.Errorf("query Advisor recommendations: Resource Graph returned nil result")
	}

	rows, malformed := arg.DecodeRowsWithStats[advisorRow](graphResult.Data)
	if malformed > 0 {
		result.Warnings = append(result.Warnings, assessment.AssessmentWarning{
			Code:    "advisor_malformed_arg_rows",
			Message: fmt.Sprintf("skipped %d malformed Advisor Resource Graph row(s)", malformed),
		})
	}

	for _, row := range rows {
		if filter != nil && filter.IsSubscriptionExcluded(row.SubscriptionID) {
			continue
		}
		if filter != nil && filter.IsServiceExcluded(row.ResourceID) {
			continue
		}

		result.Records = append(result.Records, assessment.AdvisorRecommendation{
			RecommendationID: row.RecommendationTypeID,
			SubscriptionID:   row.SubscriptionID,
			SubscriptionName: subscriptions[row.SubscriptionID],
			ResourceType:     azure.ResourceTypeFromResourceID(row.ResourceID),
			ResourceName:     row.ImpactedValue,
			ResourceID:       row.ResourceID,
			Category:         row.Category,
			Impact:           row.Impact,
			Description:      recommendationTypes[row.RecommendationTypeID],
		})
	}

	sort.Slice(result.Records, func(i, j int) bool {
		if result.Records[i].SubscriptionID != result.Records[j].SubscriptionID {
			return result.Records[i].SubscriptionID < result.Records[j].SubscriptionID
		}
		if result.Records[i].ResourceID != result.Records[j].ResourceID {
			return result.Records[i].ResourceID < result.Records[j].ResourceID
		}
		return result.Records[i].RecommendationID < result.Records[j].RecommendationID
	})
	return result, nil
}

type advisorRow struct {
	SubscriptionID       string `json:"SubscriptionId"`
	ResourceID           string `json:"ResourceId"`
	ImpactedValue        string `json:"ImpactedValue"`
	Category             string `json:"Category"`
	Impact               string `json:"Impact"`
	RecommendationTypeID string `json:"RecommendationTypeId"`
}
