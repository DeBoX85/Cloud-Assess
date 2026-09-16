package tables

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/skus"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

type Options struct {
	RedactSubscriptionIDs bool
}

type Table struct {
	Key       string
	SheetName string
	Stage     string
	Rows      [][]string
}

// Build returns report tables in the intended human-report order. Assessment Status is
// intentionally first so incomplete assessments cannot be mistaken for clean reports.
func Build(data *result.AssessmentResult, opts Options) ([]Table, error) {
	if data == nil {
		return nil, fmt.Errorf("assessment result is nil")
	}
	return []Table{
		{Key: "assessmentStatus", SheetName: "Assessment Status", Rows: assessmentStatus(data)},
		{Key: "recommendations", SheetName: "Recommendations", Stage: stages.Graph, Rows: recommendations(data)},
		{Key: "impacted", SheetName: "ImpactedResources", Stage: stages.Graph, Rows: impacted(data, opts)},
		{Key: "resourceType", SheetName: "ResourceTypes", Stage: stages.Graph, Rows: resourceTypes(data)},
		{Key: "inventory", SheetName: "Inventory", Stage: stages.Graph, Rows: resources(data, data.Resources, opts)},
		{Key: "advisor", SheetName: "Advisor", Stage: stages.Advisor, Rows: advisor(data, opts)},
		{Key: "azurePolicy", SheetName: "Azure Policy", Stage: stages.Policy, Rows: policy(data, opts)},
		{Key: "arcSQL", SheetName: "Arc SQL", Stage: stages.Arc, Rows: arcSQL(data, opts)},
		{Key: "defenderRecommendations", SheetName: "DefenderRecommendations", Stage: stages.DefenderRecommendations, Rows: defenderRecommendations(data, opts)},
		{Key: "defender", SheetName: "Defender", Stage: stages.Defender, Rows: defender(data, opts)},
		{Key: "outofscope", SheetName: "OutOfScope", Stage: stages.Graph, Rows: resources(data, data.OutOfScope, opts)},
		{Key: "costs", SheetName: "Costs", Stage: stages.Cost, Rows: costs(data, opts)},
	}, nil
}

// ShouldRender preserves stage-gated report behavior while allowing a failed requested stage
// to retain an empty table whose reason is explained on Assessment Status.
func ShouldRender(data *result.AssessmentResult, table Table) bool {
	if table.Stage == "" {
		return true
	}
	for _, execution := range data.Stages {
		if !strings.EqualFold(execution.Name, table.Stage) {
			continue
		}
		if execution.Status != assessment.StageSkipped {
			return true
		}
		for _, warning := range execution.Warnings {
			if warning.Code == "stage_not_run_after_critical_failure" {
				return true
			}
		}
		return false
	}
	return len(table.Rows) > 1
}

func assessmentStatus(data *result.AssessmentResult) [][]string {
	headers := []string{"Assessment Completeness", "Schema Version", "Scope Id", "Generated At", "Stage", "Status", "Records", "Warning Count", "Warning Codes", "Error Code", "Error Message", "Started At", "Finished At", "Duration"}
	rows := make([][]string, 1, len(data.Stages)+1)
	rows[0] = headers
	if len(data.Stages) == 0 {
		rows = append(rows, []string{string(data.Completeness), data.SchemaVersion, data.ScopeID, formatTime(data.GeneratedAt), "", "", "0", "0", "", "", "", "", "", ""})
		return rows
	}
	for _, execution := range data.Stages {
		warningCodes := make([]string, 0, len(execution.Warnings))
		for _, warning := range execution.Warnings {
			warningCodes = append(warningCodes, warning.Code)
		}
		sort.Strings(warningCodes)
		errorCode, errorMessage := "", ""
		if execution.Error != nil {
			errorCode = execution.Error.Code
			errorMessage = execution.Error.Message
		}
		duration := ""
		if !execution.StartedAt.IsZero() && !execution.FinishedAt.IsZero() {
			duration = execution.FinishedAt.Sub(execution.StartedAt).String()
		}
		rows = append(rows, []string{
			string(data.Completeness), data.SchemaVersion, data.ScopeID, formatTime(data.GeneratedAt),
			execution.Name, string(execution.Status), fmt.Sprint(execution.Records), fmt.Sprint(len(execution.Warnings)), strings.Join(warningCodes, ","),
			errorCode, errorMessage, formatTime(execution.StartedAt), formatTime(execution.FinishedAt), duration,
		})
	}
	return rows
}

func recommendations(data *result.AssessmentResult) [][]string {
	headers := []string{"Implemented", "Number of Impacted Resources", "Azure Service / Well-Architected", "Recommendation Source", "Azure Service Category / Well-Architected Area", "Azure Service / Well-Architected Topic", "Category", "Recommendation", "Impact", "Best Practices Guidance", "Read More", "Recommendation Id"}
	rows := make([][]string, 1, len(data.Summary.Recommendations)+1)
	rows[0] = headers
	deployed := map[string]bool{"microsoft.resources": true}
	for _, item := range data.ResourceTypes {
		deployed[normalize(item.ResourceType)] = true
	}
	for _, recommendation := range data.Summary.Recommendations {
		implemented := "false"
		if !deployed[normalize(recommendation.ResourceType)] {
			implemented = "N/A"
		} else if recommendation.ImpactedResources == 0 {
			implemented = "true"
		}
		parts := strings.Split(recommendation.ResourceType, "/")
		provider, service := "", ""
		if len(parts) > 0 {
			provider = parts[0]
		}
		if len(parts) > 1 {
			service = parts[1]
		}
		rows = append(rows, []string{
			implemented, fmt.Sprint(recommendation.ImpactedResources), "Azure Service", recommendation.Source,
			provider, service, recommendation.Category, recommendation.Recommendation, recommendation.Impact,
			recommendation.LongDescription, recommendation.LearnURL, recommendation.ID,
		})
	}
	return rows
}

func impacted(data *result.AssessmentResult, opts Options) [][]string {
	headers := []string{"Validated Using", "Source", "Category", "Impact", "Resource Type", "Recommendation", "Recommendation Id", "Subscription Id", "Subscription Name", "Resource Group", "Resource Name", "Resource Id", "Param1", "Param2", "Param3", "Param4", "Param5", "Learn"}
	rows := make([][]string, 1, len(data.Findings)+1)
	rows[0] = headers
	seen := map[string]struct{}{}
	for _, finding := range data.Findings {
		if strings.EqualFold(finding.Category, assessment.CategorySLA) {
			continue
		}
		key := normalize(finding.RecommendationID) + "\x00" + normalize(finding.ResourceID)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		validatedUsing := finding.ValidationMechanism
		if validatedUsing == "" {
			validatedUsing = "Azure Resource Graph"
		}
		params := [5]string{}
		copy(params[:], finding.Parameters)
		rows = append(rows, []string{
			validatedUsing, finding.Source, finding.Category, finding.Impact, finding.ResourceType, finding.Recommendation, finding.RecommendationID,
			redact.SubscriptionID(finding.SubscriptionID, opts.RedactSubscriptionIDs), finding.SubscriptionName, finding.ResourceGroup, finding.ResourceName,
			redact.SubscriptionIDInResourceID(finding.ResourceID, opts.RedactSubscriptionIDs), params[0], params[1], params[2], params[3], params[4], finding.LearnMoreURL,
		})
	}
	return rows
}

func resourceTypes(data *result.AssessmentResult) [][]string {
	rows := make([][]string, 1, len(data.ResourceTypes)+1)
	rows[0] = []string{"Subscription Name", "Resource Type", "Number of Resources"}
	for _, item := range data.ResourceTypes {
		rows = append(rows, []string{item.SubscriptionName, item.ResourceType, fmt.Sprint(item.Count)})
	}
	return rows
}

func resources(data *result.AssessmentResult, values []assessment.Resource, opts Options) [][]string {
	rows := make([][]string, 1, len(values)+1)
	rows[0] = []string{"Subscription Id", "Resource Group", "Location", "Resource Type", "Resource Name", "Sku Name", "Sku Tier", "Capacity", "Kind", "SLA", "Resource Id"}
	sla := map[string]string{}
	for _, finding := range data.Findings {
		if !strings.EqualFold(finding.Category, assessment.CategorySLA) || len(finding.Parameters) == 0 {
			continue
		}
		sla[normalize(finding.ResourceID)] = finding.Parameters[0]
	}
	for _, resource := range values {
		resourceSLA := resource.SLA
		if value, ok := sla[normalize(resource.ID)]; ok {
			resourceSLA = value
		}
		capacity := skus.ComputeCapacity(resource.SKUName, resource.SKUCapacity, strings.EqualFold(resource.Type, "microsoft.compute/virtualmachinescalesets"))
		rows = append(rows, []string{
			redact.SubscriptionID(resource.SubscriptionID, opts.RedactSubscriptionIDs), resource.ResourceGroup, resource.Location, resource.Type, resource.Name,
			resource.SKUName, resource.SKUTier, capacity, resource.Kind, resourceSLA, redact.SubscriptionIDInResourceID(resource.ID, opts.RedactSubscriptionIDs),
		})
	}
	return rows
}

func advisor(data *result.AssessmentResult, opts Options) [][]string {
	rows := make([][]string, 1, len(data.Advisor)+1)
	rows[0] = []string{"Subscription Id", "Subscription Name", "Resource Type", "Resource Name", "Category", "Impact", "Description", "Resource Id", "Recommendation Id"}
	for _, item := range data.Advisor {
		rows = append(rows, []string{
			redact.SubscriptionID(item.SubscriptionID, opts.RedactSubscriptionIDs), item.SubscriptionName, item.ResourceType, item.ResourceName,
			item.Category, item.Impact, item.Description, redact.SubscriptionIDInResourceID(item.ResourceID, opts.RedactSubscriptionIDs), item.RecommendationID,
		})
	}
	return rows
}

func defender(data *result.AssessmentResult, opts Options) [][]string {
	rows := make([][]string, 1, len(data.Defender)+1)
	rows[0] = []string{"Subscription Id", "Subscription Name", "Name", "Tier"}
	for _, item := range data.Defender {
		rows = append(rows, []string{redact.SubscriptionID(item.SubscriptionID, opts.RedactSubscriptionIDs), item.SubscriptionName, item.Name, item.Tier})
	}
	return rows
}

func defenderRecommendations(data *result.AssessmentResult, opts Options) [][]string {
	rows := make([][]string, 1, len(data.DefenderRecommendations)+1)
	rows[0] = []string{"Subscription Id", "Subscription Name", "Resource Group", "Resource Type", "Resource Name", "Category", "Recommendation Severity", "Recommendation Name", "Action Description", "Remediation Description", "AzPortal Link", "Resource Id"}
	for _, item := range data.DefenderRecommendations {
		rows = append(rows, []string{
			redact.SubscriptionID(item.SubscriptionID, opts.RedactSubscriptionIDs), item.SubscriptionName, item.ResourceGroup, item.ResourceType, item.ResourceName,
			item.Category, item.RecommendationSeverity, item.RecommendationName, item.ActionDescription, item.RemediationDescription, item.AzurePortalLink,
			redact.SubscriptionIDInResourceID(item.ResourceID, opts.RedactSubscriptionIDs),
		})
	}
	return rows
}

func policy(data *result.AssessmentResult, opts Options) [][]string {
	rows := make([][]string, 1, len(data.AzurePolicy)+1)
	rows[0] = []string{"Subscription Id", "Subscription Name", "Resource Group", "Resource Type", "Resource Name", "Policy Display Name", "Policy Description", "Resource Id", "Time Stamp", "Policy Definition Name", "Policy Definition Id", "Policy Assignment Name", "Policy Assignment Id", "Compliance State"}
	for _, item := range data.AzurePolicy {
		rows = append(rows, []string{
			redact.SubscriptionID(item.SubscriptionID, opts.RedactSubscriptionIDs), item.SubscriptionName, item.ResourceGroup, item.ResourceType, item.ResourceName,
			item.PolicyDisplayName, item.PolicyDescription, redact.SubscriptionIDInResourceID(item.ResourceID, opts.RedactSubscriptionIDs), item.Timestamp,
			item.PolicyDefinitionName, item.PolicyDefinitionID, item.PolicyAssignmentName, item.PolicyAssignmentID, item.ComplianceState,
		})
	}
	return rows
}

func arcSQL(data *result.AssessmentResult, opts Options) [][]string {
	rows := make([][]string, 1, len(data.ArcSQL)+1)
	rows[0] = []string{"Subscription Id", "Subscription Name", "Azure Arc Server", "SQL Instance", "Resource Group", "Version", "Build", "Patch Level", "Edition", "VCores", "License", "DPS Status", "TEL Status", "Defender Status"}
	for _, item := range data.ArcSQL {
		rows = append(rows, []string{
			redact.SubscriptionID(item.SubscriptionID, opts.RedactSubscriptionIDs), item.SubscriptionName,
			azure.ResourceNameFromResourceID(item.AzureArcServer), azure.ResourceNameFromResourceID(item.SQLInstance), item.ResourceGroup,
			item.Version, item.Build, item.PatchLevel, item.Edition, item.VCores, item.License, item.DPSStatus, item.TELStatus, item.DefenderStatus,
		})
	}
	return rows
}

func costs(data *result.AssessmentResult, opts Options) [][]string {
	rows := make([][]string, 1, len(data.Costs)+1)
	rows[0] = []string{"From", "To", "Subscription Id", "Subscription Name", "Service Name", "Value", "Currency"}
	for _, item := range data.Costs {
		rows = append(rows, []string{
			item.From.Format("2006-01-02"), item.To.Format("2006-01-02"), redact.SubscriptionID(item.SubscriptionID, opts.RedactSubscriptionIDs),
			item.SubscriptionName, item.ServiceName, item.Value, item.Currency,
		})
	}
	return rows
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
