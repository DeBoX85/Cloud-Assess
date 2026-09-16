// Portions of this file reproduce Azure Policy scan behavior from Microsoft Azure Quick
// Review (MIT licensed). See NOTICE.md.

package policy

import (
	"context"
	"fmt"
	"sort"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const Query = `
	PolicyResources
	| where type == 'microsoft.policyinsights/policystates'
	| extend 
		resourceId = tostring(properties.resourceId),
		subscriptionId = tostring(properties.subscriptionId),
		policyAssignmentId = tostring(properties.policyAssignmentId),
		policyAssignmentName = tostring(properties.policyAssignmentName),
		policyDefinitionId = tostring(properties.policyDefinitionId),
		policyDefinitionName = tostring(properties.policyDefinitionName),
		timestamp = todatetime(properties.timestamp),
		complianceState = tostring(properties.complianceState)
	| where complianceState == 'NonCompliant'
	| join kind=leftouter (
		PolicyResources
		| where type == 'microsoft.authorization/policydefinitions'
		| extend policyDefinitionId = tolower(id)
		| project policyDefinitionId, policyDescription = tostring(properties.description), policyDefinitionDisplayName = properties.displayName
	) on policyDefinitionId
	| join kind=leftouter (
		ResourceContainers
		| where type == 'microsoft.resources/subscriptions'
		| project subscriptionId = tolower(subscriptionId), subscriptionName = name
	) on subscriptionId
	| project subscriptionId, subscriptionName, resourceId, policyAssignmentId, policyAssignmentName, policyDefinitionId, policyDefinitionName, timestamp, policyDefinitionDisplayName, policyDescription, complianceState
	`

type GraphQuerier interface {
	Query(context.Context, string, map[string]string, ...arg.QueryOptions) (*arg.Result, error)
}

type Filter interface {
	IsSubscriptionExcluded(subscriptionID string) bool
	IsServiceExcluded(resourceID string) bool
}

type Result struct {
	Records  []assessment.PolicyNonCompliance
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

func (s *Scanner) Scan(
	ctx context.Context,
	subscriptions map[string]string,
	filter Filter,
) (Result, error) {
	result := Result{
		Records:  []assessment.PolicyNonCompliance{},
		Warnings: []assessment.AssessmentWarning{},
	}
	if s == nil || s.graph == nil {
		return result, fmt.Errorf("Azure Policy Resource Graph client is not configured")
	}

	graphResult, err := s.graph.Query(ctx, Query, subscriptions, arg.QueryOptions{ManagementGroupScope: true})
	if err != nil {
		return result, fmt.Errorf("query Azure Policy non-compliance: %w", err)
	}
	if graphResult == nil {
		return result, fmt.Errorf("query Azure Policy non-compliance: Resource Graph returned nil result")
	}

	rows, malformed := arg.DecodeRowsWithStats[policyRow](graphResult.Data)
	if malformed > 0 {
		result.Warnings = append(result.Warnings, assessment.AssessmentWarning{
			Code:    "policy_malformed_arg_rows",
			Message: fmt.Sprintf("skipped %d malformed Azure Policy Resource Graph row(s)", malformed),
		})
	}

	type policyKey struct {
		resourceID         string
		policyDefinitionID string
	}
	seen := make(map[policyKey]struct{}, len(rows))

	for _, row := range rows {
		if filter != nil && filter.IsSubscriptionExcluded(row.SubscriptionID) {
			continue
		}
		if filter != nil && filter.IsServiceExcluded(row.ResourceID) {
			continue
		}

		key := policyKey{
			resourceID:         row.ResourceID,
			policyDefinitionID: row.PolicyDefinitionID,
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		result.Records = append(result.Records, assessment.PolicyNonCompliance{
			SubscriptionID:       row.SubscriptionID,
			SubscriptionName:     row.SubscriptionName,
			ResourceGroup:        azure.ResourceGroupFromResourceID(row.ResourceID),
			ResourceType:         azure.ResourceTypeFromResourceID(row.ResourceID),
			ResourceName:         azure.ResourceNameFromResourceID(row.ResourceID),
			PolicyDisplayName:    row.PolicyDefinitionDisplay,
			PolicyDescription:    row.PolicyDescription,
			ResourceID:           row.ResourceID,
			Timestamp:            row.Timestamp,
			PolicyDefinitionName: row.PolicyDefinitionName,
			PolicyDefinitionID:   row.PolicyDefinitionID,
			PolicyAssignmentName: row.PolicyAssignmentName,
			PolicyAssignmentID:   row.PolicyAssignmentID,
			ComplianceState:      row.ComplianceState,
		})
	}

	sort.Slice(result.Records, func(i, j int) bool {
		if result.Records[i].SubscriptionID != result.Records[j].SubscriptionID {
			return result.Records[i].SubscriptionID < result.Records[j].SubscriptionID
		}
		if result.Records[i].ResourceID != result.Records[j].ResourceID {
			return result.Records[i].ResourceID < result.Records[j].ResourceID
		}
		return result.Records[i].PolicyDefinitionID < result.Records[j].PolicyDefinitionID
	})
	return result, nil
}

type policyRow struct {
	SubscriptionID          string `json:"subscriptionId"`
	SubscriptionName        string `json:"subscriptionName"`
	ResourceID              string `json:"resourceId"`
	PolicyDefinitionDisplay string `json:"policyDefinitionDisplayName"`
	PolicyDescription       string `json:"policyDescription"`
	Timestamp               string `json:"timestamp"`
	PolicyDefinitionName    string `json:"policyDefinitionName"`
	PolicyDefinitionID      string `json:"policyDefinitionId"`
	PolicyAssignmentName    string `json:"policyAssignmentName"`
	PolicyAssignmentID      string `json:"policyAssignmentId"`
	ComplianceState         string `json:"complianceState"`
}
