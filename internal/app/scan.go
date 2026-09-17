package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/gate"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	csvrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/csv"
	excelrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/excel"
	jsonrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/json"
	sarifrenderer "github.com/DeBoX85/Cloud-Assess/internal/renderers/sarif"
	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

const (
	ExitSuccess       = 0
	ExitExecutionFail = 1
	ExitQualityGate   = 2
	ExitPartial       = 3
)

type AssessmentRunner interface {
	Run(context.Context, orchestration.Request) (*result.AssessmentResult, error)
}

type OutputOptions struct {
	BaseName              string
	XLSX                  bool
	JSON                  bool
	CSV                   bool
	Stdout                bool
	SARIF                 bool
	RedactSubscriptionIDs bool
	Version               string
}

type ScanOptions struct {
	Assessment orchestration.Request
	Outputs    OutputOptions
	FailOn     string
}

type Outcome struct {
	Assessment *result.AssessmentResult
	Files      []string
	ExitCode   int
}

type PartialAssessmentError struct {
	Completeness assessment.Completeness
}

func (e *PartialAssessmentError) Error() string {
	return fmt.Sprintf("assessment completed with incomplete stage data: %s", e.Completeness)
}

type Runner struct {
	assessment AssessmentRunner
	stdout     io.Writer
	now        func() time.Time
}

func NewRunner(assessmentRunner AssessmentRunner) *Runner {
	return &Runner{assessment: assessmentRunner, stdout: os.Stdout, now: time.Now}
}

// Run executes the assessment, renders all requested artifacts, then applies assessment and
// quality-gate exit semantics. Gate/partial failures intentionally happen after rendering.
func (r *Runner) Run(ctx context.Context, options ScanOptions) (Outcome, error) {
	outcome := Outcome{ExitCode: ExitExecutionFail}
	if r == nil || r.assessment == nil {
		return outcome, fmt.Errorf("application assessment runner is not configured")
	}
	if r.now == nil {
		r.now = time.Now
	}

	var criterion gate.Criterion
	gateEnabled := options.FailOn != ""
	if gateEnabled {
		var err error
		criterion, err = gate.Parse(options.FailOn)
		if err != nil {
			return outcome, err
		}
	}

	assessmentResult, executionErr := r.assessment.Run(ctx, options.Assessment)
	outcome.Assessment = assessmentResult
	if assessmentResult == nil {
		if executionErr != nil {
			return outcome, executionErr
		}
		return outcome, fmt.Errorf("assessment returned no result")
	}

	baseName := options.Outputs.BaseName
	if baseName == "" {
		baseName = defaultBaseName(r.now())
	}
	files, err := r.render(assessmentResult, baseName, options.Outputs)
	outcome.Files = files
	if err != nil {
		return outcome, err
	}
	if executionErr != nil {
		return outcome, executionErr
	}

	switch assessmentResult.Completeness {
	case assessment.CompletenessFailed:
		return outcome, fmt.Errorf("assessment completeness is failed")
	case assessment.CompletenessPartial:
		outcome.ExitCode = ExitPartial
		return outcome, &PartialAssessmentError{Completeness: assessmentResult.Completeness}
	}

	if gateEnabled {
		if err := gate.Check(assessmentResult.Summary, criterion); err != nil {
			outcome.ExitCode = ExitQualityGate
			return outcome, err
		}
	}
	outcome.ExitCode = ExitSuccess
	return outcome, nil
}

func (r *Runner) render(data *result.AssessmentResult, baseName string, options OutputOptions) ([]string, error) {
	generated := make([]string, 0, 4)
	tableOptions := tables.Options{RedactSubscriptionIDs: options.RedactSubscriptionIDs}

	if options.XLSX {
		filename := baseName + ".xlsx"
		if err := excelrenderer.WriteFile(data, filename, tableOptions); err != nil {
			return generated, err
		}
		generated = append(generated, filename)
	}
	if options.JSON {
		filename := baseName + ".json"
		if err := jsonrenderer.WriteFile(data, filename); err != nil {
			return generated, err
		}
		generated = append(generated, filename)
	}
	if options.CSV {
		files, err := csvrenderer.Write(data, baseName, tableOptions)
		generated = append(generated, files...)
		if err != nil {
			return generated, err
		}
	}
	if options.SARIF {
		filename := baseName + ".sarif"
		if err := sarifrenderer.WriteFile(data, filename, options.Version); err != nil {
			return generated, err
		}
		generated = append(generated, filename)
	}
	if options.Stdout {
		encoded, err := jsonrenderer.String(data)
		if err != nil {
			return generated, err
		}
		writer := r.stdout
		if writer == nil {
			writer = io.Discard
		}
		if _, err := fmt.Fprintln(writer, encoded); err != nil {
			return generated, fmt.Errorf("write JSON report to stdout: %w", err)
		}
	}
	return generated, nil
}

func defaultBaseName(now time.Time) string {
	brand := branding.Default()
	stamp := now.Format("2006_01_02_T150405")
	return filepath.Clean(fmt.Sprintf("%s_%s", brand.ReportFilePrefix, stamp))
}
