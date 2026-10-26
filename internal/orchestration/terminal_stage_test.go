package orchestration

import (
	"context"
	"errors"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/cost"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

func TestCoordinatorTerminalCallerDuringFinalCostCannotReturnPartialSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	operations := noOpOperations()
	operations.ScanCost = func(context.Context, map[string]string) (cost.Result, error) {
		cancel()
		return cost.Result{}, errors.New("plain transport failure")
	}
	config := graphOnlyStages(t)
	if err := config.Set(stages.Cost, true); err != nil {
		t.Fatal(err)
	}
	result, err := NewCoordinator(operations).Run(ctx, Request{Stages: config})
	if !errors.Is(err, context.Canceled) || result == nil || result.Completeness != assessment.CompletenessFailed {
		t.Fatalf("terminal final Cost became partial successful assessment: result=%+v error=%v", result, err)
	}
	for _, execution := range result.Stages {
		if execution.Name == stages.Cost && (execution.Error == nil || execution.Error.Code != "assessment_canceled") {
			t.Fatalf("Cost context failure lost: %+v", execution)
		}
	}
}
