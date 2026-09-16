// Portions of this file reproduce Cost Management scan behavior from Microsoft Azure Quick
// Review (MIT licensed). See NOTICE.md.

package cost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const (
	apiVersion        = "2021-10-01"
	defaultWorkers    = 2
	requestTimeout    = 120 * time.Second
	typeActualCost    = "ActualCost"
	timeframeCustom   = "Custom"
	functionSum       = "Sum"
	dimensionGrouping = "Dimension"
)

type Poster interface {
	Post(context.Context, string, io.ReadSeekCloser) ([]byte, *http.Response, error)
}

type Result struct {
	Records  []assessment.CostRecord
	Warnings []assessment.AssessmentWarning
}

type Scanner struct {
	client     Poster
	endpoint   string
	maxWorkers int
}

func New(credential azcore.TokenCredential) *Scanner {
	return NewWithClient(
		azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(requestTimeout)),
		azure.ResourceManagerEndpoint(),
	)
}

func NewWithClient(client Poster, endpoint string) *Scanner {
	return &Scanner{
		client:     client,
		endpoint:   strings.TrimSuffix(endpoint, "/"),
		maxWorkers: defaultWorkers,
	}
}

// Scan queries previous-completed-month actual costs for each subscription. The source stage
// uses two workers; Cloud Assess preserves that concurrency ceiling and adds deterministic
// output ordering after collection.
func (s *Scanner) Scan(
	ctx context.Context,
	subscriptions map[string]string,
) (Result, error) {
	return s.scanAt(ctx, subscriptions, time.Now().UTC())
}

func (s *Scanner) scanAt(
	ctx context.Context,
	subscriptions map[string]string,
	now time.Time,
) (Result, error) {
	result := Result{
		Records:  []assessment.CostRecord{},
		Warnings: []assessment.AssessmentWarning{},
	}
	if len(subscriptions) == 0 {
		return result, nil
	}
	if s == nil || s.client == nil {
		return result, fmt.Errorf("Cost Management HTTP client is not configured")
	}
	if s.endpoint == "" {
		return result, fmt.Errorf("Azure Resource Manager endpoint is not configured")
	}

	from, to := PreviousCompletedMonth(now)
	subscriptionIDs := make([]string, 0, len(subscriptions))
	for subscriptionID := range subscriptions {
		subscriptionIDs = append(subscriptionIDs, subscriptionID)
	}
	sort.Strings(subscriptionIDs)

	workerCount := s.maxWorkers
	if workerCount <= 0 {
		workerCount = defaultWorkers
	}
	if workerCount > len(subscriptionIDs) {
		workerCount = len(subscriptionIDs)
	}

	type subscriptionResult struct {
		records  []assessment.CostRecord
		warnings []assessment.AssessmentWarning
		err      error
	}

	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan string)
	results := make(chan subscriptionResult, len(subscriptionIDs))

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for range workerCount {
		go func() {
			defer wg.Done()
			for subscriptionID := range jobs {
				records, warnings, err := s.querySubscription(
					workerCtx,
					subscriptionID,
					subscriptions[subscriptionID],
					from,
					to,
				)
				results <- subscriptionResult{records: records, warnings: warnings, err: err}
				if err != nil {
					cancel()
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, subscriptionID := range subscriptionIDs {
			select {
			case jobs <- subscriptionID:
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
	for subscription := range results {
		result.Records = append(result.Records, subscription.records...)
		result.Warnings = append(result.Warnings, subscription.warnings...)
		if subscription.err != nil && firstErr == nil {
			firstErr = subscription.err
		}
	}
	if firstErr != nil {
		return result, firstErr
	}

	sort.Slice(result.Records, func(i, j int) bool {
		if result.Records[i].SubscriptionID != result.Records[j].SubscriptionID {
			return result.Records[i].SubscriptionID < result.Records[j].SubscriptionID
		}
		if result.Records[i].ServiceName != result.Records[j].ServiceName {
			return result.Records[i].ServiceName < result.Records[j].ServiceName
		}
		if result.Records[i].Currency != result.Records[j].Currency {
			return result.Records[i].Currency < result.Records[j].Currency
		}
		return result.Records[i].Value < result.Records[j].Value
	})
	sort.Slice(result.Warnings, func(i, j int) bool {
		if result.Warnings[i].Code != result.Warnings[j].Code {
			return result.Warnings[i].Code < result.Warnings[j].Code
		}
		return result.Warnings[i].Message < result.Warnings[j].Message
	})
	return result, nil
}

func (s *Scanner) querySubscription(
	ctx context.Context,
	subscriptionID string,
	subscriptionName string,
	from time.Time,
	to time.Time,
) ([]assessment.CostRecord, []assessment.AssessmentWarning, error) {
	request := queryDefinition{
		Type:      typeActualCost,
		Timeframe: timeframeCustom,
		TimePeriod: queryTimePeriod{
			From: from,
			To:   to,
		},
		Dataset: queryDataset{
			Aggregation: map[string]queryAggregation{
				"TotalCost": {
					Name:     "Cost",
					Function: functionSum,
				},
			},
			Grouping: []queryGrouping{{
				Type: dimensionGrouping,
				Name: "ServiceName",
			}},
		},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal cost query for subscription %s: %w", subscriptionID, err)
	}

	url := fmt.Sprintf(
		"%s/subscriptions/%s/providers/Microsoft.CostManagement/query?api-version=%s",
		s.endpoint,
		subscriptionID,
		apiVersion,
	)
	responseBody, _, err := s.client.Post(ctx, url, azure.NopReadSeekCloser{Reader: bytes.NewReader(body)})
	if err != nil {
		if errorCode, ok := azure.IsSkippableResponseError(err); ok {
			return nil, []assessment.AssessmentWarning{{
				Code: "cost_subscription_skipped",
				Message: fmt.Sprintf(
					"skipped Cost Management query for subscription %s because Azure returned %s",
					subscriptionID,
					errorCode,
				),
			}}, nil
		}
		return nil, nil, fmt.Errorf("query costs for subscription %s: %w", subscriptionID, err)
	}
	if len(bytes.TrimSpace(responseBody)) == 0 {
		return []assessment.CostRecord{}, nil, nil
	}

	var response queryResult
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, nil, fmt.Errorf("decode costs for subscription %s: %w", subscriptionID, err)
	}

	records := make([]assessment.CostRecord, 0, len(response.Properties.Rows))
	warnings := []assessment.AssessmentWarning{}
	for _, row := range response.Properties.Rows {
		if len(row) < 3 {
			warnings = append(warnings, assessment.AssessmentWarning{
				Code: "cost_malformed_row",
				Message: fmt.Sprintf(
					"skipped malformed Cost Management row for subscription %s: expected at least 3 values, got %d",
					subscriptionID,
					len(row),
				),
			})
			continue
		}
		records = append(records, assessment.CostRecord{
			From:             from,
			To:               to,
			SubscriptionID:   subscriptionID,
			SubscriptionName: subscriptionName,
			ServiceName:      fmt.Sprintf("%v", row[1]),
			Value:            fmt.Sprintf("%v", row[0]),
			Currency:         fmt.Sprintf("%v", row[2]),
		})
	}
	return records, warnings, nil
}

type queryDefinition struct {
	Type       string          `json:"type"`
	Timeframe  string          `json:"timeframe"`
	TimePeriod queryTimePeriod `json:"timePeriod"`
	Dataset    queryDataset    `json:"dataset"`
}

type queryTimePeriod struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type queryDataset struct {
	Aggregation map[string]queryAggregation `json:"aggregation"`
	Grouping    []queryGrouping             `json:"grouping"`
}

type queryAggregation struct {
	Name     string `json:"name"`
	Function string `json:"function"`
}

type queryGrouping struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type queryResult struct {
	Properties queryProperties `json:"properties"`
}

type queryProperties struct {
	Rows [][]any `json:"rows"`
}
