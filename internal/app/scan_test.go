package app

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

type fakeAssessmentRunner struct {
	result *result.AssessmentResult
	err    error
	calls  int
}

func (f *fakeAssessmentRunner) Run(context.Context, orchestration.Request) (*result.AssessmentResult, error) {
	f.calls++
	return f.result, f.err
}

func TestRunRendersBeforeQualityGateFailure(t *testing.T) {
	assessmentResult := result.Build(result.Input{
		GeneratedAt:  time.Unix(1, 0),
		Completeness: assessment.CompletenessComplete,
		Recommendations: []assessment.RecommendationDefinition{{
			ID: "rec-high", ResourceType: "Microsoft.Storage/storageAccounts", Recommendation: "Fix it", Impact: assessment.ImpactHigh,
		}},
		Findings: []assessment.Finding{{
			RecommendationID: "rec-high", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st", ResourceType: "Microsoft.Storage/storageAccounts", Recommendation: "Fix it", Impact: assessment.ImpactHigh,
		}},
	})
	base := filepath.Join(t.TempDir(), "report")
	runner := NewRunner(&fakeAssessmentRunner{result: assessmentResult})
	outcome, err := runner.Run(context.Background(), ScanOptions{
		Outputs: OutputOptions{BaseName: base, JSON: true},
		FailOn:  "High",
	})
	if err == nil {
		t.Fatal("expected quality gate failure")
	}
	if outcome.ExitCode != ExitQualityGate {
		t.Fatalf("exit code = %d, want %d", outcome.ExitCode, ExitQualityGate)
	}
	if _, statErr := os.Stat(base + ".json"); statErr != nil {
		t.Fatalf("JSON artifact was not written before gate failure: %v", statErr)
	}
}

func TestRunRendersPartialAssessmentThenReturnsExitThree(t *testing.T) {
	assessmentResult := result.Build(result.Input{GeneratedAt: time.Unix(1, 0), Completeness: assessment.CompletenessPartial})
	base := filepath.Join(t.TempDir(), "partial")
	runner := NewRunner(&fakeAssessmentRunner{result: assessmentResult})
	outcome, err := runner.Run(context.Background(), ScanOptions{Outputs: OutputOptions{BaseName: base, JSON: true}})
	if err == nil {
		t.Fatal("expected partial-assessment error")
	}
	if outcome.ExitCode != ExitPartial {
		t.Fatalf("exit code = %d, want %d", outcome.ExitCode, ExitPartial)
	}
	if _, statErr := os.Stat(base + ".json"); statErr != nil {
		t.Fatalf("JSON artifact was not written for partial assessment: %v", statErr)
	}
}

func TestRunPersistsCriticalFailureStatusWhenResultExists(t *testing.T) {
	assessmentResult := result.Build(result.Input{GeneratedAt: time.Unix(1, 0), Completeness: assessment.CompletenessFailed})
	base := filepath.Join(t.TempDir(), "failed")
	runner := NewRunner(&fakeAssessmentRunner{result: assessmentResult, err: errors.New("graph failed")})
	outcome, err := runner.Run(context.Background(), ScanOptions{Outputs: OutputOptions{BaseName: base, JSON: true}})
	if err == nil || err.Error() != "graph failed" {
		t.Fatalf("error = %v, want graph failed", err)
	}
	if outcome.ExitCode != ExitExecutionFail {
		t.Fatalf("exit code = %d, want %d", outcome.ExitCode, ExitExecutionFail)
	}
	if _, statErr := os.Stat(base + ".json"); statErr != nil {
		t.Fatalf("failure-status artifact was not written: %v", statErr)
	}
}

func TestRunRejectsInvalidGateBeforeAssessment(t *testing.T) {
	fake := &fakeAssessmentRunner{result: result.Build(result.Input{Completeness: assessment.CompletenessComplete})}
	runner := NewRunner(fake)
	outcome, err := runner.Run(context.Background(), ScanOptions{FailOn: "Critical"})
	if err == nil {
		t.Fatal("expected invalid gate error")
	}
	if fake.calls != 0 {
		t.Fatal("assessment should not run when fail-on is invalid")
	}
	if outcome.ExitCode != ExitExecutionFail {
		t.Fatalf("exit code = %d, want %d", outcome.ExitCode, ExitExecutionFail)
	}
}

func TestRunStdoutUsesCanonicalJSON(t *testing.T) {
	assessmentResult := result.Build(result.Input{GeneratedAt: time.Unix(1, 0), Completeness: assessment.CompletenessComplete})
	var output bytes.Buffer
	runner := NewRunner(&fakeAssessmentRunner{result: assessmentResult})
	runner.stdout = &output
	outcome, err := runner.Run(context.Background(), ScanOptions{Outputs: OutputOptions{Stdout: true}})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if outcome.ExitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want 0", outcome.ExitCode)
	}
	if !bytes.Contains(output.Bytes(), []byte(`"schemaVersion": "1.0"`)) {
		t.Fatalf("stdout did not contain canonical JSON: %s", output.String())
	}
}

func TestDefaultBaseNameUsesCentralBrandPrefix(t *testing.T) {
	got := defaultBaseName(time.Date(2026, 9, 17, 13, 5, 6, 0, time.FixedZone("CEST", 2*60*60)))
	if got != "cloud_assessment_2026_09_17_T130506" {
		t.Fatalf("default base = %q", got)
	}
}
