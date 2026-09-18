package equivalence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/diagnostics"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/skus"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

var diagnosticRecommendationIDs = func() map[string]struct{} {
	ids := map[string]struct{}{}
	for _, definition := range diagnostics.Recommendations() {
		ids[lower(definition.ID)] = struct{}{}
	}
	return ids
}()

func LoadReference(reader io.Reader) (Projection, error) {
	if reader == nil {
		return Projection{}, fmt.Errorf("reference reader is nil")
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return Projection{}, fmt.Errorf("read reference JSON: %w", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(content, &raw); err != nil {
		return Projection{}, fmt.Errorf("decode reference JSON: %w", err)
	}

	projection := newProjection()
	if containsRedactedSubscriptionID(content) {
		projection.Comparable = false
		projection.Notes = append(projection.Notes, "reference report contains redacted subscription IDs; rerun with --mask=false")
	}
	sections := map[string]string{
		"recommendations":         DatasetRecommendations,
		"impacted":                DatasetFindings,
		"resourceType":            DatasetResourceTypes,
		"inventory":               DatasetInventory,
		"outOfScope":              DatasetOutOfScope,
		"advisor":                 DatasetAdvisor,
		"defender":                DatasetDefender,
		"defenderRecommendations": DatasetDefenderRecommendations,
		"azurePolicy":             DatasetAzurePolicy,
		"arcSQL":                  DatasetArcSQL,
		"costs":                   DatasetCosts,
	}

	for section, datasetName := range sections {
		payload, exists := raw[section]
		dataset := projection.Datasets[datasetName]
		dataset.Enabled = exists
		projection.Datasets[datasetName] = dataset
		if !exists {
			continue
		}
		var rows []map[string]string
		if string(payload) != "null" {
			if err := json.Unmarshal(payload, &rows); err != nil {
				return Projection{}, fmt.Errorf("decode reference section %s: %w", section, err)
			}
		}
		if err := projectReferenceRows(&projection, datasetName, rows); err != nil {
			return Projection{}, err
		}
	}

	if _, ok := raw["externalPlugins"]; ok {
		projection.Notes = append(projection.Notes, "reference external plugin output is intentionally not compared by the core equivalence harness")
	}
	return projection, nil
}

func LoadTarget(reader io.Reader) (Projection, error) {
	if reader == nil {
		return Projection{}, fmt.Errorf("target reader is nil")
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return Projection{}, fmt.Errorf("read target JSON: %w", err)
	}
	var data result.AssessmentResult
	if err := json.Unmarshal(content, &data); err != nil {
		return Projection{}, fmt.Errorf("decode target JSON: %w", err)
	}
	if strings.TrimSpace(data.SchemaVersion) == "" {
		return Projection{}, fmt.Errorf("target JSON is missing schemaVersion")
	}

	projection := newProjection()
	if containsRedactedSubscriptionID(content) {
		projection.Comparable = false
		projection.Notes = append(projection.Notes, "target report contains redacted subscription IDs; rerun with --redact-subscription-ids=false")
	}
	switch data.Completeness {
	case assessment.CompletenessFailed, assessment.CompletenessPartial:
		projection.Comparable = false
		projection.Notes = append(projection.Notes, "target assessment completeness is "+string(data.Completeness))
	case assessment.CompletenessCompleteWithWarnings:
		projection.Notes = append(projection.Notes, "target assessment completed with warnings; review Assessment Status alongside equivalence deltas")
	}

	graphFallback := len(data.Recommendations) > 0 || len(data.Findings) > 0 || len(data.Resources) > 0 || len(data.OutOfScope) > 0 || len(data.ResourceTypes) > 0
	graphEnabled := targetStageEnabled(data.Stages, stages.Graph, graphFallback)
	setTargetEnabled(&projection, graphEnabled,
		DatasetRecommendations,
		DatasetFindings,
		DatasetResourceTypes,
		DatasetInventory,
		DatasetOutOfScope,
	)
	setTargetEnabled(&projection, targetStageEnabled(data.Stages, stages.Advisor, len(data.Advisor) > 0), DatasetAdvisor)
	setTargetEnabled(&projection, targetStageEnabled(data.Stages, stages.Defender, len(data.Defender) > 0), DatasetDefender)
	setTargetEnabled(&projection, targetStageEnabled(data.Stages, stages.DefenderRecommendations, len(data.DefenderRecommendations) > 0), DatasetDefenderRecommendations)
	setTargetEnabled(&projection, targetStageEnabled(data.Stages, stages.Policy, len(data.AzurePolicy) > 0), DatasetAzurePolicy)
	setTargetEnabled(&projection, targetStageEnabled(data.Stages, stages.Arc, len(data.ArcSQL) > 0), DatasetArcSQL)
	setTargetEnabled(&projection, targetStageEnabled(data.Stages, stages.Cost, len(data.Costs) > 0), DatasetCosts)

	projectTargetRecommendations(&projection, &data)
	projectTargetFindings(&projection, &data)
	projectTargetResources(&projection, DatasetInventory, data.Resources, data.Findings)
	projectTargetResources(&projection, DatasetOutOfScope, data.OutOfScope, data.Findings)
	projectTargetResourceTypes(&projection, &data)
	projectTargetAdvisor(&projection, &data)
	projectTargetDefender(&projection, &data)
	projectTargetDefenderRecommendations(&projection, &data)
	projectTargetPolicy(&projection, &data)
	projectTargetArcSQL(&projection, &data)
	projectTargetCosts(&projection, &data)
	return projection, nil
}

func setTargetEnabled(projection *Projection, enabled bool, names ...string) {
	for _, name := range names {
		dataset := projection.Datasets[name]
		dataset.Enabled = enabled
		projection.Datasets[name] = dataset
	}
}

func targetStageEnabled(executions []assessment.StageExecution, name string, fallback bool) bool {
	for _, execution := range executions {
		if strings.EqualFold(strings.TrimSpace(execution.Name), name) {
			return execution.Status != assessment.StageSkipped
		}
	}
	return fallback
}

func projectReferenceRows(projection *Projection, name string, rows []map[string]string) error {
	switch name {
	case DatasetRecommendations:
		for _, row := range rows {
			id := lower(get(row, "recommendationId"))
			put(projection, name, id, map[string]string{
				"source":            normalizeSource(id, get(row, "recommendationSource")),
				"category":          lower(get(row, "category")),
				"impact":            lower(get(row, "impact")),
				"recommendation":    text(get(row, "recommendation")),
				"longDescription":   text(get(row, "bestPracticesGuidance")),
				"learn":             strings.TrimSpace(get(row, "readMore")),
				"impactedResources": integer(get(row, "numberOfImpactedResources")),
				"implemented":       lower(get(row, "implemented")),
			})
		}
	case DatasetFindings:
		for _, row := range rows {
			id := lower(get(row, "recommendationId"))
			resourceID := lower(get(row, "resourceId"))
			put(projection, name, recordKey(id, resourceID), map[string]string{
				"recommendationId": id,
				"resourceId":       resourceID,
				"source":           normalizeSource(id, get(row, "source")),
				"category":         lower(get(row, "category")),
				"impact":           lower(get(row, "impact")),
				"resourceType":     lower(get(row, "resourceType")),
				"recommendation":   text(get(row, "recommendation")),
				"subscriptionId":   lower(get(row, "subscriptionId")),
				"subscriptionName": text(get(row, "subscriptionName")),
				"resourceGroup":    lower(get(row, "resourceGroup")),
				"resourceName":     lower(get(row, "resourceName")),
				"param1":           text(get(row, "param1")),
				"param2":           text(get(row, "param2")),
				"param3":           text(get(row, "param3")),
				"param4":           text(get(row, "param4")),
				"param5":           text(get(row, "param5")),
				"learn":            strings.TrimSpace(get(row, "learn")),
			})
		}
	case DatasetInventory, DatasetOutOfScope:
		for _, row := range rows {
			resourceID := lower(get(row, "resourceId"))
			put(projection, name, resourceID, referenceResourceFields(row))
		}
	case DatasetResourceTypes:
		for _, row := range rows {
			key := recordKey(lower(get(row, "subscriptionName")), lower(get(row, "resourceType")))
			addCountRecord(projection, name, key, "count", integer(get(row, "numberOfResources")), map[string]string{
				"subscriptionName": text(get(row, "subscriptionName")),
				"resourceType":     lower(get(row, "resourceType")),
			})
		}
	case DatasetAdvisor:
		for _, row := range rows {
			resourceID := lower(get(row, "resourceId"))
			recID := lower(get(row, "recommendationId"))
			put(projection, name, recordKey(resourceID, recID), map[string]string{
				"subscriptionId":   lower(get(row, "subscriptionId")),
				"subscriptionName": text(get(row, "subscriptionName")),
				"resourceType":     lower(get(row, "resourceType")),
				"resourceName":     lower(get(row, "resourceName")),
				"category":         lower(get(row, "category")),
				"impact":           lower(get(row, "impact")),
				"description":      text(get(row, "description")),
				"resourceId":       resourceID,
				"recommendationId": recID,
			})
		}
	case DatasetDefender:
		for _, row := range rows {
			subID := lower(get(row, "subscriptionId"))
			nameValue := lower(get(row, "name"))
			put(projection, name, recordKey(subID, nameValue), map[string]string{
				"subscriptionId":   subID,
				"subscriptionName": text(get(row, "subscriptionName")),
				"name":             nameValue,
				"tier":             lower(get(row, "tier")),
			})
		}
	case DatasetDefenderRecommendations:
		for _, row := range rows {
			resourceID := lower(get(row, "resourceId"))
			category := lower(get(row, "category"))
			recName := text(get(row, "recommendationName"))
			put(projection, name, recordKey(resourceID, category, lower(recName)), map[string]string{
				"subscriptionId":         lower(get(row, "subscriptionId")),
				"subscriptionName":       text(get(row, "subscriptionName")),
				"resourceGroup":          lower(get(row, "resourceGroup")),
				"resourceType":           lower(get(row, "resourceType")),
				"resourceName":           lower(get(row, "resourceName")),
				"category":               category,
				"recommendationSeverity": lower(get(row, "recommendationSeverity")),
				"recommendationName":     recName,
				"actionDescription":      text(get(row, "actionDescription")),
				"remediationDescription": text(get(row, "remediationDescription")),
				"azurePortalLink":        strings.TrimSpace(get(row, "azPortalLink")),
				"resourceId":             resourceID,
			})
		}
	case DatasetAzurePolicy:
		for _, row := range rows {
			resourceID := lower(get(row, "resourceId"))
			definitionID := lower(get(row, "policyDefinitionId"))
			put(projection, name, recordKey(resourceID, definitionID), map[string]string{
				"subscriptionId":       lower(get(row, "subscriptionId")),
				"subscriptionName":     text(get(row, "subscriptionName")),
				"resourceGroup":        lower(get(row, "resourceGroup")),
				"resourceType":         lower(get(row, "resourceType")),
				"resourceName":         lower(get(row, "resourceName")),
				"policyDisplayName":    text(get(row, "policyDisplayName")),
				"policyDescription":    text(get(row, "policyDescription")),
				"resourceId":           resourceID,
				"policyDefinitionName": lower(get(row, "policyDefinitionName")),
				"policyDefinitionId":   definitionID,
				"policyAssignmentName": lower(get(row, "policyAssignmentName")),
				"policyAssignmentId":   lower(get(row, "policyAssignmentId")),
				"complianceState":      lower(get(row, "complianceState")),
			})
		}
	case DatasetArcSQL:
		for _, row := range rows {
			subID := lower(get(row, "subscriptionId"))
			instance := lower(get(row, "sqlInstance"))
			put(projection, name, recordKey(subID, instance), map[string]string{
				"subscriptionId":   subID,
				"subscriptionName": text(get(row, "subscriptionName")),
				"azureArcServer":   lower(get(row, "azureArcServer")),
				"sqlInstance":      instance,
				"resourceGroup":    lower(get(row, "resourceGroup")),
				"version":          text(get(row, "version")),
				"build":            text(get(row, "build")),
				"patchLevel":       text(get(row, "patchLevel")),
				"edition":          text(get(row, "edition")),
				"vCores":           number(get(row, "vCores")),
				"license":          lower(get(row, "license")),
				"dpsStatus":        lower(get(row, "dpsStatus")),
				"telStatus":        lower(get(row, "telStatus")),
				"defenderStatus":   lower(get(row, "defenderStatus")),
			})
		}
	case DatasetCosts:
		for _, row := range rows {
			subID := lower(get(row, "subscriptionId"))
			service := lower(get(row, "serviceName"))
			currency := strings.ToUpper(strings.TrimSpace(get(row, "currency")))
			from := strings.TrimSpace(get(row, "from"))
			to := strings.TrimSpace(get(row, "to"))
			put(projection, name, recordKey(subID, service, currency, from, to), map[string]string{
				"subscriptionId": subID,
				"serviceName":    service,
				"value":          number(get(row, "value")),
				"currency":       currency,
				"from":           from,
				"to":             to,
			})
		}
	default:
		return fmt.Errorf("unsupported reference dataset %s", name)
	}
	return nil
}

func referenceResourceFields(row map[string]string) map[string]string {
	return map[string]string{
		"subscriptionId": lower(get(row, "subscriptionId")),
		"resourceGroup":  lower(get(row, "resourceGroup")),
		"location":       lower(get(row, "location")),
		"resourceType":   lower(get(row, "resourceType")),
		"resourceName":   lower(get(row, "resourceName")),
		"skuName":        lower(get(row, "skuName")),
		"skuTier":        lower(get(row, "skuTier")),
		"capacity":       text(get(row, "capacity")),
		"kind":           lower(get(row, "kind")),
		"sla":            text(get(row, "sla")),
		"resourceId":     lower(get(row, "resourceId")),
	}
}

func projectTargetRecommendations(projection *Projection, data *result.AssessmentResult) {
	if data.Summary == nil {
		return
	}
	deployed := map[string]bool{"microsoft.resources": true}
	for _, item := range data.ResourceTypes {
		deployed[lower(item.ResourceType)] = true
	}
	for _, item := range data.Summary.Recommendations {
		id := lower(item.ID)
		implemented := "false"
		if !deployed[lower(item.ResourceType)] {
			implemented = "n/a"
		} else if item.ImpactedResources == 0 {
			implemented = "true"
		}
		put(projection, DatasetRecommendations, id, map[string]string{
			"source":            normalizeSource(id, item.Source),
			"category":          lower(item.Category),
			"impact":            lower(item.Impact),
			"recommendation":    text(item.Recommendation),
			"longDescription":   text(item.LongDescription),
			"learn":             strings.TrimSpace(item.LearnURL),
			"impactedResources": strconv.Itoa(item.ImpactedResources),
			"implemented":       implemented,
		})
	}
}

func projectTargetFindings(projection *Projection, data *result.AssessmentResult) {
	for _, finding := range data.Findings {
		if strings.EqualFold(finding.Category, assessment.CategorySLA) {
			continue
		}
		id := lower(finding.RecommendationID)
		resourceID := lower(finding.ResourceID)
		params := [5]string{}
		for i := range params {
			if i < len(finding.Parameters) {
				params[i] = text(finding.Parameters[i])
			}
		}
		put(projection, DatasetFindings, recordKey(id, resourceID), map[string]string{
			"recommendationId": id,
			"resourceId":       resourceID,
			"source":           normalizeSource(id, finding.Source),
			"category":         lower(finding.Category),
			"impact":           lower(finding.Impact),
			"resourceType":     lower(finding.ResourceType),
			"recommendation":   text(finding.Recommendation),
			"subscriptionId":   lower(finding.SubscriptionID),
			"subscriptionName": text(finding.SubscriptionName),
			"resourceGroup":    lower(finding.ResourceGroup),
			"resourceName":     lower(finding.ResourceName),
			"param1":           params[0],
			"param2":           params[1],
			"param3":           params[2],
			"param4":           params[3],
			"param5":           params[4],
			"learn":            strings.TrimSpace(finding.LearnMoreURL),
		})
	}
}

func projectTargetResources(projection *Projection, datasetName string, resources []assessment.Resource, findings []assessment.Finding) {
	slaByID := map[string]string{}
	for _, finding := range findings {
		if !strings.EqualFold(finding.Category, assessment.CategorySLA) || len(finding.Parameters) == 0 {
			continue
		}
		slaByID[lower(finding.ResourceID)] = text(finding.Parameters[0])
	}
	for _, resource := range resources {
		resourceID := lower(resource.ID)
		sla := text(resource.SLA)
		if inherited, ok := slaByID[resourceID]; ok {
			sla = inherited
		}
		put(projection, datasetName, resourceID, map[string]string{
			"subscriptionId": lower(resource.SubscriptionID),
			"resourceGroup":  lower(resource.ResourceGroup),
			"location":       lower(resource.Location),
			"resourceType":   lower(resource.Type),
			"resourceName":   lower(resource.Name),
			"skuName":        lower(resource.SKUName),
			"skuTier":        lower(resource.SKUTier),
			"capacity":       text(skus.ComputeCapacity(resource.SKUName, resource.SKUCapacity, strings.EqualFold(resource.Type, "microsoft.compute/virtualmachinescalesets"))),
			"kind":           lower(resource.Kind),
			"sla":            sla,
			"resourceId":     resourceID,
		})
	}
}

func projectTargetResourceTypes(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.ResourceTypes {
		key := recordKey(lower(item.SubscriptionName), lower(item.ResourceType))
		addCountRecord(projection, DatasetResourceTypes, key, "count", strconv.Itoa(item.Count), map[string]string{
			"subscriptionName": text(item.SubscriptionName),
			"resourceType":     lower(item.ResourceType),
		})
	}
}

func projectTargetAdvisor(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.Advisor {
		resourceID := lower(item.ResourceID)
		recID := lower(item.RecommendationID)
		put(projection, DatasetAdvisor, recordKey(resourceID, recID), map[string]string{
			"subscriptionId":   lower(item.SubscriptionID),
			"subscriptionName": text(item.SubscriptionName),
			"resourceType":     lower(item.ResourceType),
			"resourceName":     lower(item.ResourceName),
			"category":         lower(item.Category),
			"impact":           lower(item.Impact),
			"description":      text(item.Description),
			"resourceId":       resourceID,
			"recommendationId": recID,
		})
	}
}

func projectTargetDefender(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.Defender {
		subID := lower(item.SubscriptionID)
		nameValue := lower(item.Name)
		put(projection, DatasetDefender, recordKey(subID, nameValue), map[string]string{
			"subscriptionId":   subID,
			"subscriptionName": text(item.SubscriptionName),
			"name":             nameValue,
			"tier":             lower(item.Tier),
		})
	}
}

func projectTargetDefenderRecommendations(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.DefenderRecommendations {
		resourceID := lower(item.ResourceID)
		category := lower(item.Category)
		recName := text(item.RecommendationName)
		put(projection, DatasetDefenderRecommendations, recordKey(resourceID, category, lower(recName)), map[string]string{
			"subscriptionId":         lower(item.SubscriptionID),
			"subscriptionName":       text(item.SubscriptionName),
			"resourceGroup":          lower(item.ResourceGroup),
			"resourceType":           lower(item.ResourceType),
			"resourceName":           lower(item.ResourceName),
			"category":               category,
			"recommendationSeverity": lower(item.RecommendationSeverity),
			"recommendationName":     recName,
			"actionDescription":      text(item.ActionDescription),
			"remediationDescription": text(item.RemediationDescription),
			"azurePortalLink":        strings.TrimSpace(item.AzurePortalLink),
			"resourceId":             resourceID,
		})
	}
}

func projectTargetPolicy(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.AzurePolicy {
		resourceID := lower(item.ResourceID)
		definitionID := lower(item.PolicyDefinitionID)
		put(projection, DatasetAzurePolicy, recordKey(resourceID, definitionID), map[string]string{
			"subscriptionId":       lower(item.SubscriptionID),
			"subscriptionName":     text(item.SubscriptionName),
			"resourceGroup":        lower(item.ResourceGroup),
			"resourceType":         lower(item.ResourceType),
			"resourceName":         lower(item.ResourceName),
			"policyDisplayName":    text(item.PolicyDisplayName),
			"policyDescription":    text(item.PolicyDescription),
			"resourceId":           resourceID,
			"policyDefinitionName": lower(item.PolicyDefinitionName),
			"policyDefinitionId":   definitionID,
			"policyAssignmentName": lower(item.PolicyAssignmentName),
			"policyAssignmentId":   lower(item.PolicyAssignmentID),
			"complianceState":      lower(item.ComplianceState),
		})
	}
}

func projectTargetArcSQL(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.ArcSQL {
		subID := lower(item.SubscriptionID)
		instance := lower(azure.ResourceNameFromResourceID(item.SQLInstance))
		put(projection, DatasetArcSQL, recordKey(subID, instance), map[string]string{
			"subscriptionId":   subID,
			"subscriptionName": text(item.SubscriptionName),
			"azureArcServer":   lower(azure.ResourceNameFromResourceID(item.AzureArcServer)),
			"sqlInstance":      instance,
			"resourceGroup":    lower(item.ResourceGroup),
			"version":          text(item.Version),
			"build":            text(item.Build),
			"patchLevel":       text(item.PatchLevel),
			"edition":          text(item.Edition),
			"vCores":           number(item.VCores),
			"license":          lower(item.License),
			"dpsStatus":        lower(item.DPSStatus),
			"telStatus":        lower(item.TELStatus),
			"defenderStatus":   lower(item.DefenderStatus),
		})
	}
}

func projectTargetCosts(projection *Projection, data *result.AssessmentResult) {
	for _, item := range data.Costs {
		subID := lower(item.SubscriptionID)
		service := lower(item.ServiceName)
		currency := strings.ToUpper(strings.TrimSpace(item.Currency))
		from := item.From.Format("2006-01-02")
		to := item.To.Format("2006-01-02")
		put(projection, DatasetCosts, recordKey(subID, service, currency, from, to), map[string]string{
			"subscriptionId": subID,
			"serviceName":    service,
			"value":          number(item.Value),
			"currency":       currency,
			"from":           from,
			"to":             to,
		})
	}
}

func put(projection *Projection, datasetName, key string, fields map[string]string) {
	dataset := projection.Datasets[datasetName]
	if dataset.Records == nil {
		dataset.Records = map[string]Record{}
	}
	dataset.Records[key] = Record{Key: key, Fields: fields}
	projection.Datasets[datasetName] = dataset
}

func addCountRecord(projection *Projection, datasetName, key, countField, countValue string, fields map[string]string) {
	dataset := projection.Datasets[datasetName]
	if existing, ok := dataset.Records[key]; ok {
		current, _ := strconv.Atoi(existing.Fields[countField])
		add, _ := strconv.Atoi(countValue)
		existing.Fields[countField] = strconv.Itoa(current + add)
		dataset.Records[key] = existing
		projection.Datasets[datasetName] = dataset
		return
	}
	fields[countField] = countValue
	put(projection, datasetName, key, fields)
}

func get(row map[string]string, key string) string {
	return row[key]
}

func recordKey(parts ...string) string {
	return strings.Join(parts, "|")
}

func text(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func lower(value string) string {
	return strings.ToLower(text(value))
}

func integer(value string) string {
	value = strings.TrimSpace(value)
	if numberValue, err := strconv.Atoi(value); err == nil {
		return strconv.Itoa(numberValue)
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return strconv.FormatInt(int64(floatValue), 10)
	}
	return text(value)
}

func number(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if parsed, err := strconv.ParseFloat(value, 64); err == nil {
		return strconv.FormatFloat(parsed, 'g', -1, 64)
	}
	return text(value)
}

func normalizeSource(recommendationID, source string) string {
	normalized := strings.ToUpper(text(source))
	if normalized == "AZQR" {
		if _, ok := diagnosticRecommendationIDs[lower(recommendationID)]; ok {
			return diagnostics.Source
		}
		return rules.SourceCustom
	}
	return normalized
}


var redactedSubscriptionMarker = []byte("xxxxxxxx-xxxx-xxxx-xxxx-xxxxx")

func containsRedactedSubscriptionID(content []byte) bool {
	return bytes.Contains(bytes.ToLower(content), redactedSubscriptionMarker)
}
