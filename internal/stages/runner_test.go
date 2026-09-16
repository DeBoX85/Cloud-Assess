package stages

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func testClock() func() time.Time {
	current := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	return func() time.Time {
		value := current
		current = current.Add(time.Second)
		return value
	}
}

func TestRunnerCompleteWithSkippedStage(t *testing.T) {
	runner := newRunnerWithClock(testClock())
	result := runner.Execute(context.Background(), []Task{
		{
			Name:    "graph",
			Enabled: true,
			Run: func(context.Context) (Outcome, error) {
				return Outcome{Records: 7}, nil
			},
		},
		{Name: "cost", Enabled: false},
	})

	if result.Completeness != assessment.CompletenessComplete {
		t.Fatalf("completeness = %q, want complete", result.Completeness)
	}
	if len(result.Stages) != 2 {
		t.Fatalf("stages = %d, want 2", len(result.Stages))
	}
	if result.Stages[0].Status != assessment.StageCompleted || result.Stages[0].Records != 7 {
		t.Fatalf("unexpected graph execution: %+v", result.Stages[0])
	}
	if result.Stages[1].Status != assessment.StageSkipped || len(result.Stages[1].Warnings) != 0 {
		t.Fatalf("disabled stage should be a clean skip: %+v", result.Stages[1])
	}
}

func TestWarningsProduceCompleteWithWarnings(t *testing.T) {
	warning := assessment.AssessmentWarning{Code: "partial_rows", Message: "one malformed row was skipped"}
	result := newRunnerWithClock(testClock()).Execute(context.Background(), []Task{{
		Name:    "advisor",
		Enabled: true,
		Run: func(context.Context) (Outcome, error) {
			return Outcome{Records: 3, Warnings: []assessment.AssessmentWarning{warning}}, nil
		},
	}})

	if result.Completeness != assessment.CompletenessCompleteWithWarnings {
		t.Fatalf("completeness = %q, want complete_with_warnings", result.Completeness)
	}
	stage := result.Stages[0]
	if stage.Status != assessment.StageCompletedWithWarnings || stage.Records != 3 || len(stage.Warnings) != 1 {
		t.Fatalf("unexpected warning stage: %+v", stage)
	}
}

func TestNonCriticalFailureProducesPartialAndContinues(t *testing.T) {
	ranAfterFailure := false
	result := newRunnerWithClock(testClock()).Execute(context.Background(), []Task{
		{
			Name:     "advisor",
			Enabled:  true,
			Critical: false,
			Run: func(context.Context) (Outcome, error) {
				return Outcome{}, errors.New("advisor unavailable")
			},
		},
		{
			Name:    "defender",
			Enabled: true,
			Run: func(context.Context) (Outcome, error) {
				ranAfterFailure = true
				return Outcome{Records: 2}, nil
			},
		},
	})

	if result.Completeness != assessment.CompletenessPartial {
		t.Fatalf("completeness = %q, want partial", result.Completeness)
	}
	if !ranAfterFailure {
		t.Fatal("non-critical failure should not stop later stages")
	}
	if result.Stages[0].Status != assessment.StageFailed || result.Stages[0].Error == nil || result.Stages[0].Error.Code != "stage_failed" {
		t.Fatalf("unexpected failed stage: %+v", result.Stages[0])
	}
	if result.Stages[1].Status != assessment.StageCompleted {
		t.Fatalf("later stage did not complete: %+v", result.Stages[1])
	}
}

func TestCriticalFailureStopsRequestedDownstreamStages(t *testing.T) {
	ran := false
	result := newRunnerWithClock(testClock()).Execute(context.Background(), []Task{
		{
			Name:     "resource-discovery",
			Enabled:  true,
			Critical: true,
			Run: func(context.Context) (Outcome, error) {
				return Outcome{}, errors.New("inventory failed")
			},
		},
		{
			Name:    "graph",
			Enabled: true,
			Run: func(context.Context) (Outcome, error) {
				ran = true
				return Outcome{}, nil
			},
		},
	})

	if result.Completeness != assessment.CompletenessFailed {
		t.Fatalf("completeness = %q, want failed", result.Completeness)
	}
	if ran {
		t.Fatal("requested downstream stage ran after critical failure")
	}
	if result.Stages[0].Status != assessment.StageFailed {
		t.Fatalf("critical stage should fail: %+v", result.Stages[0])
	}
	if result.Stages[1].Status != assessment.StageSkipped || len(result.Stages[1].Warnings) != 1 || result.Stages[1].Warnings[0].Code != "stage_not_run_after_critical_failure" {
		t.Fatalf("downstream stage should expose why it was skipped: %+v", result.Stages[1])
	}
}

func TestContextCancellationStopsEvenNonCriticalStageSequence(t *testing.T) {
	result := newRunnerWithClock(testClock()).Execute(context.Background(), []Task{
		{
			Name:    "advisor",
			Enabled: true,
			Run: func(context.Context) (Outcome, error) {
				return Outcome{}, context.Canceled
			},
		},
		{Name: "defender", Enabled: true, Run: func(context.Context) (Outcome, error) {
			t.Fatal("stage ran after cancellation")
			return Outcome{}, nil
		}},
	})

	if result.Completeness != assessment.CompletenessFailed {
		t.Fatalf("completeness = %q, want failed", result.Completeness)
	}
	if result.Stages[0].Error == nil || result.Stages[0].Error.Code != "assessment_canceled" {
		t.Fatalf("unexpected cancellation error: %+v", result.Stages[0])
	}
	if result.Stages[1].Status != assessment.StageSkipped {
		t.Fatalf("later stage should be skipped: %+v", result.Stages[1])
	}
}

func TestNilNonCriticalStageFunctionProducesPartialAndContinues(t *testing.T) {
	ran := false
	result := newRunnerWithClock(testClock()).Execute(context.Background(), []Task{
		{Name: "advisor", Enabled: true},
		{Name: "defender", Enabled: true, Run: func(context.Context) (Outcome, error) {
			ran = true
			return Outcome{}, nil
		}},
	})

	if result.Completeness != assessment.CompletenessPartial || !ran {
		t.Fatalf("nil non-critical stage should be partial and continue: %+v", result)
	}
	if result.Stages[0].Error == nil || result.Stages[0].Error.Code != "stage_not_configured" {
		t.Fatalf("unexpected missing-stage error: %+v", result.Stages[0])
	}
}

func TestPartialOutranksWarnings(t *testing.T) {
	result := newRunnerWithClock(testClock()).Execute(context.Background(), []Task{
		{Name: "advisor", Enabled: true, Run: func(context.Context) (Outcome, error) {
			return Outcome{}, errors.New("failed")
		}},
		{Name: "defender", Enabled: true, Run: func(context.Context) (Outcome, error) {
			return Outcome{Warnings: []assessment.AssessmentWarning{{Code: "warning", Message: "warning"}}}, nil
		}},
	})
	if result.Completeness != assessment.CompletenessPartial {
		t.Fatalf("completeness = %q, want partial", result.Completeness)
	}
}
