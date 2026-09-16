package stages

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// Outcome is the health information returned by one successful stage execution.
type Outcome struct {
	Records  int
	Warnings []assessment.AssessmentWarning
}

// Task is one ordered assessment stage. Critical stages stop the assessment when they fail;
// non-critical stage failures are recorded and execution continues so partial assessments are
// explicit rather than silently losing data.
type Task struct {
	Name     string
	Enabled  bool
	Critical bool
	Run      func(context.Context) (Outcome, error)
}

type RunResult struct {
	Stages       []assessment.StageExecution
	Completeness assessment.Completeness
}

type Runner struct {
	now func() time.Time
}

func NewRunner() *Runner {
	return &Runner{now: time.Now}
}

func newRunnerWithClock(now func() time.Time) *Runner {
	return &Runner{now: now}
}

func (r *Runner) Execute(ctx context.Context, tasks []Task) RunResult {
	result := RunResult{
		Stages:       make([]assessment.StageExecution, 0, len(tasks)),
		Completeness: assessment.CompletenessComplete,
	}
	if r == nil || r.now == nil {
		r = NewRunner()
	}

	stop := false
	for _, task := range tasks {
		if stop {
			result.Stages = append(result.Stages, r.skippedAfterFailure(task.Name))
			continue
		}
		if !task.Enabled {
			result.Stages = append(result.Stages, r.skipped(task.Name))
			continue
		}

		started := r.now().UTC()
		if task.Run == nil {
			finished := r.now().UTC()
			result.Stages = append(result.Stages, assessment.StageExecution{
				Name:       task.Name,
				Status:     assessment.StageFailed,
				StartedAt:  started,
				FinishedAt: finished,
				Error: &assessment.AssessmentError{
					Code:    "stage_not_configured",
					Message: fmt.Sprintf("stage %q has no execution function", task.Name),
				},
			})
			if task.Critical {
				result.Completeness = assessment.CompletenessFailed
				stop = true
			} else {
				result.Completeness = mergeCompleteness(result.Completeness, assessment.CompletenessPartial)
			}
			continue
		}

		outcome, err := task.Run(ctx)
		finished := r.now().UTC()
		execution := assessment.StageExecution{
			Name:       task.Name,
			StartedAt:  started,
			FinishedAt: finished,
			Records:    outcome.Records,
			Warnings:   append([]assessment.AssessmentWarning(nil), outcome.Warnings...),
		}

		if err != nil {
			execution.Status = assessment.StageFailed
			execution.Error = &assessment.AssessmentError{
				Code:    stageErrorCode(err),
				Message: err.Error(),
			}
			result.Stages = append(result.Stages, execution)

			if task.Critical || isContextFailure(err) {
				result.Completeness = assessment.CompletenessFailed
				stop = true
			} else {
				result.Completeness = mergeCompleteness(result.Completeness, assessment.CompletenessPartial)
			}
			continue
		}

		if len(execution.Warnings) > 0 {
			execution.Status = assessment.StageCompletedWithWarnings
			result.Completeness = mergeCompleteness(result.Completeness, assessment.CompletenessCompleteWithWarnings)
		} else {
			execution.Status = assessment.StageCompleted
		}
		result.Stages = append(result.Stages, execution)
	}

	return result
}

func (r *Runner) skipped(name string) assessment.StageExecution {
	now := r.now().UTC()
	return assessment.StageExecution{
		Name:       name,
		Status:     assessment.StageSkipped,
		StartedAt:  now,
		FinishedAt: now,
	}
}

func (r *Runner) skippedAfterFailure(name string) assessment.StageExecution {
	now := r.now().UTC()
	return assessment.StageExecution{
		Name:       name,
		Status:     assessment.StageSkipped,
		StartedAt:  now,
		FinishedAt: now,
		Warnings: []assessment.AssessmentWarning{{
			Code:    "stage_not_run_after_critical_failure",
			Message: "stage was not run because an earlier critical stage failed",
		}},
	}
}

func stageErrorCode(err error) string {
	if errors.Is(err, context.Canceled) {
		return "assessment_canceled"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "assessment_deadline_exceeded"
	}
	return "stage_failed"
}

func isContextFailure(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func mergeCompleteness(current, candidate assessment.Completeness) assessment.Completeness {
	if completenessRank(candidate) > completenessRank(current) {
		return candidate
	}
	return current
}

func completenessRank(value assessment.Completeness) int {
	switch value {
	case assessment.CompletenessComplete:
		return 0
	case assessment.CompletenessCompleteWithWarnings:
		return 1
	case assessment.CompletenessPartial:
		return 2
	case assessment.CompletenessFailed:
		return 3
	default:
		return 3
	}
}
