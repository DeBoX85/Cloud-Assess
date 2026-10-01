package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/advisor"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

func TestAssessmentBudgetPersistsDeadlineFailureAndHealthyData(t *testing.T) {
	resource := assessment.Resource{ID: "/subscriptions/fixture/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store", SubscriptionID: "fixture", Type: "Microsoft.Storage/storageAccounts"}
	definition := assessment.RecommendationDefinition{ID: "storage-test", ResourceType: resource.Type, Query: "resources", Source: rules.SourceCustom}
	finding := assessment.Finding{RecommendationID: definition.ID, ResourceID: resource.ID, ResourceType: resource.Type}
	catalog := rules.NewCatalog()
	catalog.Add(definition)
	operations := operationsForEndToEndTest(t, catalog, resource, finding)
	observedDeadline := false
	operations.ScanAdvisor = func(ctx context.Context, _ map[string]string, _ *config.AssessmentFilter) (advisor.Result, error) {
		deadline, ok := ctx.Deadline()
		observedDeadline = ok && time.Until(deadline) < 2*time.Second
		<-ctx.Done()
		return advisor.Result{}, ctx.Err()
	}
	parent, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	base := filepath.Join(t.TempDir(), "budget")
	outcome, err := NewRunner(orchestration.NewCoordinator(operations)).Run(parent, ScanOptions{AssessmentTimeout: time.Second, Assessment: orchestration.Request{Subscriptions: []string{"fixture"}, ScannerKeys: []string{"st"}, Stages: stages.NewDefault()}, Outputs: OutputOptions{BaseName: base, JSON: true}})
	if err == nil || outcome.ExitCode != ExitExecutionFail || !observedDeadline || parent.Err() != nil {
		t.Fatalf("error=%v exit=%d deadline=%v parent=%v", err, outcome.ExitCode, observedDeadline, parent.Err())
	}
	raw, err := os.ReadFile(base + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var report result.AssessmentResult
	if err = json.Unmarshal(raw, &report); err != nil {
		t.Fatal(err)
	}
	if report.Completeness != assessment.CompletenessFailed || len(report.Resources) != 1 || len(report.Findings) != 1 {
		t.Fatalf("lost healthy data or failure: %+v", report)
	}
	for _, stage := range report.Stages {
		if stage.Name == stages.Advisor && (stage.Error == nil || stage.Error.Code != "assessment_deadline_exceeded") {
			t.Fatalf("deadline status missing: %+v", stage)
		}
	}
}

type contextObserver struct{ observed context.Context }

func (f *contextObserver) Run(ctx context.Context, _ orchestration.Request) (*result.AssessmentResult, error) {
	f.observed = ctx
	return result.Build(result.Input{Completeness: assessment.CompletenessComplete}), nil
}
func TestAssessmentBudgetZeroEarlierParentAndNegative(t *testing.T) {
	for _, budget := range []time.Duration{0, time.Hour} {
		parent, cancel := context.WithTimeout(context.Background(), time.Minute)
		observer := &contextObserver{}
		_, err := NewRunner(observer).Run(parent, ScanOptions{AssessmentTimeout: budget})
		if err != nil {
			t.Fatal(err)
		}
		want, _ := parent.Deadline()
		got, ok := observer.observed.Deadline()
		if !ok || !got.Equal(want) {
			t.Fatal("earlier parent deadline replaced")
		}
		if parent.Err() != nil {
			t.Fatal("caller context canceled")
		}
		cancel()
	}
	observer := &contextObserver{}
	_, err := NewRunner(observer).Run(context.Background(), ScanOptions{AssessmentTimeout: -time.Second})
	if err == nil || observer.observed != nil {
		t.Fatal("negative timeout started assessment")
	}
	observer = &contextObserver{}
	_, err = NewRunner(observer).Run(context.Background(), ScanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := observer.observed.Deadline(); ok {
		t.Fatal("default added a deadline")
	}
}
