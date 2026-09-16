// Portions of this file reproduce diagnostic-settings scan behavior from Microsoft
// Azure Quick Review (MIT licensed). See NOTICE.md.

package diagnostics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const (
	defaultBatchSize  = 20
	defaultMaxWorkers = 30
	batchAPIVersion   = "2020-06-01"
	diagnosticsAPI    = "2021-05-01-preview"
)

type ResourceFilter interface {
	IsServiceExcluded(resourceID string) bool
}

type BatchClient interface {
	PostStream(context.Context, string, io.ReadSeekCloser) (*http.Response, error)
}

type Result struct {
	Recommendations []assessment.RecommendationDefinition
	Findings        []assessment.Finding
	Warnings        []assessment.AssessmentWarning
}

type Scanner struct {
	client     BatchClient
	endpoint   string
	batchSize  int
	maxWorkers int
}

func New(credential azcore.TokenCredential) *Scanner {
	return NewWithClient(
		azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(60*time.Second)),
		azure.ResourceManagerEndpoint(),
	)
}

func NewWithClient(client BatchClient, endpoint string) *Scanner {
	return &Scanner{
		client:     client,
		endpoint:   strings.TrimSuffix(endpoint, "/"),
		batchSize:  defaultBatchSize,
		maxWorkers: defaultMaxWorkers,
	}
}

// Scan reproduces the reference diagnostics decision: resources with at least one
// diagnostic setting are compliant; supported resources without a setting receive the
// dedicated diagnostic recommendation when one exists.
func (s *Scanner) Scan(
	ctx context.Context,
	resources []assessment.Resource,
	filter ResourceFilter,
	subscriptions map[string]string,
) (Result, error) {
	result := Result{Recommendations: Recommendations()}

	enabled, warnings, err := s.ListResourcesWithDiagnosticSettings(ctx, resources)
	result.Warnings = append(result.Warnings, warnings...)
	if err != nil {
		return result, err
	}

	for _, resource := range resources {
		resourceID := normalize(resource.ID)
		if enabled[resourceID] {
			continue
		}
		if filter != nil && filter.IsServiceExcluded(resource.ID) {
			continue
		}

		definition, ok := RecommendationFor(resource.Type)
		if !ok {
			continue
		}

		learnURL := ""
		if len(definition.LearnMore) > 0 {
			learnURL = definition.LearnMore[0].URL
		}
		result.Findings = append(result.Findings, assessment.Finding{
			RecommendationID:    definition.ID,
			Source:              definition.Source,
			ValidationMechanism: definition.ValidationMechanism,
			Category:            definition.Category,
			Impact:              definition.Impact,
			ResourceType:        resource.Type,
			Recommendation:      definition.Recommendation,
			LongDescription:     definition.LongDescription,
			PotentialBenefits:   definition.PotentialBenefits,
			ResourceID:          resource.ID,
			SubscriptionID:      resource.SubscriptionID,
			SubscriptionName:    subscriptions[resource.SubscriptionID],
			ResourceGroup:       resource.ResourceGroup,
			ResourceName:        resource.Name,
			LearnMoreURL:        learnURL,
			AutomationAvailable: definition.AutomationAvailable,
		})
	}

	return result, nil
}

// ListResourcesWithDiagnosticSettings returns lower-case resource IDs for resources that
// have at least one diagnostic setting. Only resource types in the reference support table
// are sent to ARM. Batch size and worker limits intentionally match the reference.
func (s *Scanner) ListResourcesWithDiagnosticSettings(
	ctx context.Context,
	resources []assessment.Resource,
) (map[string]bool, []assessment.AssessmentWarning, error) {
	enabled := map[string]bool{}
	warnings := []assessment.AssessmentWarning{}

	resourceIDs := make([]string, 0, len(resources))
	for _, resource := range resources {
		if Supports(resource.Type) {
			resourceIDs = append(resourceIDs, resource.ID)
		}
	}
	if len(resourceIDs) == 0 {
		return enabled, warnings, nil
	}
	if s == nil || s.client == nil {
		return nil, warnings, fmt.Errorf("diagnostics batch client is not configured")
	}
	if s.endpoint == "" {
		return nil, warnings, fmt.Errorf("Azure Resource Manager endpoint is not configured")
	}
	if len(resourceIDs) > 5000 {
		warnings = append(warnings, assessment.AssessmentWarning{
			Code:    "diagnostics_large_scope",
			Message: fmt.Sprintf("%d resources support diagnostic settings; the diagnostics scan may take longer than usual", len(resourceIDs)),
		})
	}

	batchSize := s.batchSize
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	batches := make([][]string, 0, (len(resourceIDs)+batchSize-1)/batchSize)
	for i := 0; i < len(resourceIDs); i += batchSize {
		end := i + batchSize
		if end > len(resourceIDs) {
			end = len(resourceIDs)
		}
		batch := append([]string(nil), resourceIDs[i:end]...)
		batches = append(batches, batch)
	}

	workerCount := s.maxWorkers
	if workerCount <= 0 {
		workerCount = defaultMaxWorkers
	}
	if workerCount > len(batches) {
		workerCount = len(batches)
	}

	type batchResult struct {
		enabled  map[string]bool
		warnings []assessment.AssessmentWarning
		err      error
	}

	jobs := make(chan []string)
	results := make(chan batchResult, len(batches))
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for range workerCount {
		go func() {
			defer wg.Done()
			for ids := range jobs {
				batchEnabled, batchWarnings, err := s.scanBatch(workerCtx, ids)
				results <- batchResult{enabled: batchEnabled, warnings: batchWarnings, err: err}
				if err != nil {
					cancel()
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, batch := range batches {
			select {
			case jobs <- batch:
			case <-workerCtx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var firstErr error
	for batch := range results {
		warnings = append(warnings, batch.warnings...)
		if batch.err != nil {
			if firstErr == nil {
				firstErr = batch.err
			}
			continue
		}
		for resourceID := range batch.enabled {
			enabled[resourceID] = true
		}
	}
	if firstErr != nil {
		return enabled, warnings, firstErr
	}
	return enabled, warnings, nil
}

func (s *Scanner) scanBatch(
	ctx context.Context,
	resourceIDs []string,
) (map[string]bool, []assessment.AssessmentWarning, error) {
	request := armBatchRequest{Requests: make([]armBatchRequestItem, 0, len(resourceIDs))}
	for _, resourceID := range resourceIDs {
		request.Requests = append(request.Requests, armBatchRequestItem{
			HTTPMethod:  http.MethodGet,
			RelativeURL: resourceID + "/providers/microsoft.insights/diagnosticSettings?api-version=" + diagnosticsAPI,
		})
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal diagnostics batch request: %w", err)
	}

	url := s.endpoint + "/batch?api-version=" + batchAPIVersion
	response, err := s.client.PostStream(ctx, url, azure.NopReadSeekCloser{Reader: bytes.NewReader(body)})
	if err != nil {
		return nil, nil, fmt.Errorf("query diagnostic settings batch: %w", err)
	}
	defer response.Body.Close()

	var batch armBatchResponse
	if err := json.NewDecoder(response.Body).Decode(&batch); err != nil {
		return nil, nil, fmt.Errorf("decode diagnostic settings batch response: %w", err)
	}

	enabled := map[string]bool{}
	warnings := []assessment.AssessmentWarning{}
	for _, item := range batch.Responses {
		if item.HTTPStatusCode != http.StatusOK {
			// The reference implementation leaves this resource absent from the enabled map,
			// which means it may still receive a missing-diagnostics finding. Preserve that
			// output behavior, but make the uncertainty explicit for stage completeness.
			warnings = append(warnings, assessment.AssessmentWarning{
				Code:    "diagnostics_subrequest_non_success",
				Message: fmt.Sprintf("diagnostic settings batch subrequest returned HTTP %d", item.HTTPStatusCode),
			})
			continue
		}

		var payload diagnosticSettingsPayload
		if err := json.Unmarshal(item.Content, &payload); err != nil {
			return nil, warnings, fmt.Errorf("decode diagnostic settings subresponse: %w", err)
		}
		for _, setting := range payload.Value {
			resourceID, ok := resourceIDFromDiagnosticSetting(setting.ID)
			if !ok {
				warnings = append(warnings, assessment.AssessmentWarning{
					Code:    "diagnostics_malformed_setting_id",
					Message: "diagnostic settings response contained a setting without a parseable resource ID",
				})
				continue
			}
			enabled[resourceID] = true
		}
	}
	return enabled, warnings, nil
}

type armBatchRequest struct {
	Requests []armBatchRequestItem `json:"requests"`
}

type armBatchRequestItem struct {
	HTTPMethod  string `json:"httpMethod"`
	RelativeURL string `json:"relativeUrl"`
}

type armBatchResponse struct {
	Responses []armBatchResponseItem `json:"responses"`
}

type armBatchResponseItem struct {
	HTTPStatusCode int             `json:"httpStatusCode"`
	Content        json.RawMessage `json:"content"`
}

type diagnosticSettingsPayload struct {
	Value []diagnosticSetting `json:"value"`
}

type diagnosticSetting struct {
	ID string `json:"id"`
}

func resourceIDFromDiagnosticSetting(settingID string) (string, bool) {
	id := normalize(settingID)
	const marker = "/providers/microsoft.insights/diagnosticsettings/"
	index := strings.Index(id, marker)
	if index <= 0 {
		return "", false
	}
	return id[:index], true
}
