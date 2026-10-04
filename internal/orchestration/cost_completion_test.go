package orchestration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/advisor"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/cost"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

type incompleteCostPoster struct{}

func (incompleteCostPoster) Post(_ context.Context, _ string, body io.ReadSeekCloser) ([]byte, *http.Response, error) {
	defer body.Close()
	if _, err := io.ReadAll(body); err != nil {
		return nil, nil, err
	}
	return []byte(`{"properties":{"rows":[[12.34,"Storage","USD"]],"nextLink":"https://management.azure.com/subscriptions/sub/providers/Microsoft.CostManagement/query?api-version=2021-10-01&$skiptoken=next"}}`), &http.Response{StatusCode: http.StatusOK}, nil
}

func TestUnconsumedCostContinuationMarksAssessmentPartial(t *testing.T) {
	operations := noOpOperations()
	operations.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
		return map[string]string{"sub": "Synthetic"}, nil
	}
	operations.ScanCost = cost.NewWithClient(incompleteCostPoster{}, "https://management.azure.com").Scan
	operations.ScanAdvisor = func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error) {
		return advisor.Result{Records: []assessment.AdvisorRecommendation{{RecommendationID: "retained-advisor"}}}, nil
	}
	stageConfig := graphOnlyStages(t)
	if err := stageConfig.Set(stages.Cost, true); err != nil {
		t.Fatal(err)
	}
	if err := stageConfig.Set(stages.Advisor, true); err != nil {
		t.Fatal(err)
	}
	got, err := NewCoordinator(operations).Run(context.Background(), Request{Subscriptions: []string{"sub"}, Stages: stageConfig})
	if err != nil || got == nil || got.Completeness != assessment.CompletenessPartial || len(got.Costs) != 0 || len(got.Advisor) != 1 || got.Advisor[0].RecommendationID != "retained-advisor" {
		t.Fatalf("incomplete cost assessment=%#v error=%v; want partial, no accepted costs and retained Advisor", got, err)
	}
	found := false
	for _, stage := range got.Stages {
		if stage.Name == stages.Cost {
			found = true
			if stage.Status != assessment.StageFailed || stage.Error == nil || !strings.Contains(stage.Error.Message, "continuation") || stage.Records != 0 {
				t.Fatalf("Cost health=%#v; want explicit failed continuation stage", stage)
			}
		}
		if stage.Name == stages.Graph && stage.Status != assessment.StageCompleted {
			t.Fatalf("healthy Graph stage changed: %#v", stage)
		}
	}
	if !found {
		t.Fatal("Cost stage absent")
	}
}
