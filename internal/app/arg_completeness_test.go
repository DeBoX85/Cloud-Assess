package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/advisor"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

type truncatedGraphPoster struct{ calls int }

func (p *truncatedGraphPoster) PostStream(context.Context, string, io.ReadSeekCloser) (*http.Response, error) {
	p.calls++
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"data":[{"id":"must-not-become-a-record"}],"count":1,"totalRecords":10,"resultTruncated":"true"}`))}, nil
}

func TestTokenlessARGTruncationPersistsFailedInventoryAndGraphStages(t *testing.T) {
	for _, failedStage := range []string{orchestration.StageResourceInventory, stages.Graph} {
		t.Run(failedStage, func(t *testing.T) {
			const sub = "11111111-2222-3333-4444-555555555555"
			resource := assessment.Resource{ID: "/subscriptions/" + sub + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store", SubscriptionID: sub, ResourceGroup: "rg", Type: "Microsoft.Storage/storageAccounts"}
			catalog := rules.NewCatalog()
			catalog.Add(assessment.RecommendationDefinition{ID: "fixture", ResourceType: resource.Type, Query: "resources", Source: rules.SourceCustom})
			operations := operationsForEndToEndTest(t, catalog, resource, assessment.Finding{})
			poster := &truncatedGraphPoster{}
			client := arg.NewClient(arg.NewHTTPTransportWithClient(poster, "https://example.test/graph"))
			graphCalled := false
			operations.ExecuteGraph = func(ctx context.Context, defs []assessment.RecommendationDefinition, subs map[string]string, _ *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
				graphCalled = true
				return arg.ExecuteRecommendations(ctx, client, defs, subs, nil, 1)
			}
			if failedStage == orchestration.StageResourceInventory {
				operations.DiscoverResources = func(ctx context.Context, subs map[string]string, filters *config.Filters) (*discovery.ResourceInventory, error) {
					return discovery.DiscoverResources(ctx, client, subs, filters)
				}
			}
			advisorCalled := false
			operations.ScanAdvisor = func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
				advisorCalled = true
				return advisor.Result{}, nil
			}
			base := filepath.Join(t.TempDir(), "failed-truncation")
			outcome, err := NewRunner(orchestration.NewCoordinator(operations)).Run(context.Background(), ScanOptions{Assessment: orchestration.Request{Subscriptions: []string{sub}, ScannerKeys: []string{"st"}, Stages: stages.NewDefault()}, Outputs: OutputOptions{BaseName: base, JSON: true}})
			if err == nil || !strings.Contains(err.Error(), "truncated without a continuation token") || outcome.ExitCode != ExitExecutionFail || advisorCalled || poster.calls != 1 {
				t.Fatalf("error=%v exit=%d advisor=%v requests=%d", err, outcome.ExitCode, advisorCalled, poster.calls)
			}
			if graphCalled != (failedStage == stages.Graph) {
				t.Fatal("Graph execution boundary incorrect")
			}
			data, err := os.ReadFile(base + ".json")
			if err != nil {
				t.Fatal(err)
			}
			var report result.AssessmentResult
			if err := json.Unmarshal(data, &report); err != nil {
				t.Fatal(err)
			}
			wantResources := 0
			if failedStage == stages.Graph {
				wantResources = 1
			}
			if report.Completeness != assessment.CompletenessFailed || len(report.Resources) != wantResources || len(report.Findings) != 0 {
				t.Fatalf("report lost failure or published incomplete data: %+v", report)
			}
			failed, skipped := false, false
			for _, stage := range report.Stages {
				if stage.Name == failedStage {
					failed = stage.Status == assessment.StageFailed && stage.Error != nil && stage.Error.Code == "stage_failed"
				}
				if stage.Name == stages.Advisor {
					skipped = stage.Status == assessment.StageSkipped
				}
			}
			if !failed || !skipped {
				t.Fatal("failed/skipped stage evidence missing")
			}
		})
	}
}
