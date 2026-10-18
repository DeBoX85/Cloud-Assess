package stages

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestRunnerCallerTerminalContextWinsOverPlainTaskError(t *testing.T) {
	for _, mode := range []string{"canceled", "deadline", "ordinary"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			if mode == "deadline" {
				cancel()
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			}
			defer cancel()
			calls := 0
			run := NewRunner().Execute(ctx, []Task{
				{Name: "optional-final", Enabled: true, Run: func(taskContext context.Context) (Outcome, error) {
					calls++
					if mode == "canceled" {
						cancel()
					}
					if mode != "ordinary" {
						<-taskContext.Done()
					}
					return Outcome{Records: 2, Warnings: []assessment.AssessmentWarning{{Code: "retained", Message: "retained evidence"}}}, errors.New("plain transport failure")
				}},
				{Name: "disabled", Enabled: false},
			})
			wantCompleteness := assessment.CompletenessFailed
			wantCode := "assessment_canceled"
			if mode == "deadline" {
				wantCode = "assessment_deadline_exceeded"
			}
			if mode == "ordinary" {
				wantCompleteness = assessment.CompletenessPartial
				wantCode = "stage_failed"
			}
			if calls != 1 || run.Completeness != wantCompleteness || run.Stages[0].Error == nil || run.Stages[0].Error.Code != wantCode {
				t.Fatalf("caller terminal context concealed by task error: mode=%s calls=%d result=%+v stage=%+v", mode, calls, run, run.Stages[0])
			}
			if run.Stages[0].Records != 2 || len(run.Stages[0].Warnings) != 1 || run.Stages[1].Status != assessment.StageSkipped {
				t.Fatalf("stage evidence or disabled-stage status lost: %+v", run)
			}
		})
	}
}
