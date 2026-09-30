package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/app"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

// The child is the test executable using the production dispatcher and application
// renderer with synthetic assessment data. It does not authenticate to Azure.
func TestCLIProcessExitPreservesReport(t *testing.T) {
	for _, tc := range []struct {
		name         string
		exit         int
		completeness assessment.Completeness
	}{
		{"success", 0, assessment.CompletenessComplete},
		{"failure", 1, assessment.CompletenessFailed},
		{"severity", 2, assessment.CompletenessComplete},
		{"partial", 3, assessment.CompletenessPartial},
	} {
		t.Run(tc.name, func(t *testing.T) {
			base := filepath.Join(t.TempDir(), "report")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestCLIProcessHelper$")
			cmd.Env = append(os.Environ(), "CLOUD_ASSESS_TEST_PROCESS_CASE="+tc.name, "CLOUD_ASSESS_TEST_PROCESS_OUTPUT="+base)
			output, err := cmd.CombinedOutput()
			if ctx.Err() != nil {
				t.Fatalf("child timed out: %v\n%s", ctx.Err(), output)
			}
			exit := 0
			if err != nil {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("start child: %v", err)
				}
				exit = exitErr.ExitCode()
			}
			if exit != tc.exit {
				t.Fatalf("process exit = %d, want %d\n%s", exit, tc.exit, output)
			}
			data, err := os.ReadFile(base + ".json")
			if err != nil {
				t.Fatalf("report missing after child exit: %v", err)
			}
			var report result.AssessmentResult
			if err := json.Unmarshal(data, &report); err != nil {
				t.Fatal(err)
			}
			if report.Completeness != tc.completeness || len(report.Resources) != 1 {
				t.Fatalf("unexpected persisted report: completeness=%s resources=%d", report.Completeness, len(report.Resources))
			}
			wantFindings := 1
			if tc.name == "failure" {
				wantFindings = 0
			}
			if len(report.Findings) != wantFindings {
				t.Fatalf("persisted findings = %d, want %d", len(report.Findings), wantFindings)
			}
		})
	}
}

type processAssessmentFixture struct {
	report *result.AssessmentResult
	err    error
}

func (f processAssessmentFixture) Run(context.Context, orchestration.Request) (*result.AssessmentResult, error) {
	return f.report, f.err
}

func TestCLIProcessHelper(t *testing.T) {
	name := os.Getenv("CLOUD_ASSESS_TEST_PROCESS_CASE")
	if name == "" {
		return
	}
	input := result.Input{
		GeneratedAt:  time.Unix(1, 0),
		Completeness: assessment.CompletenessComplete,
		Resources:    []assessment.Resource{{ID: "/subscriptions/fixture/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/example", SubscriptionID: "fixture"}},
		Recommendations: []assessment.RecommendationDefinition{{
			ID: "fixture-high", Recommendation: "Fixture recommendation", Impact: assessment.ImpactHigh,
		}},
		Findings: []assessment.Finding{{
			RecommendationID: "fixture-high", Recommendation: "Fixture recommendation", Impact: assessment.ImpactHigh,
		}},
	}
	var executionErr error
	args := []string{"scan", "--json", "--xlsx=false", "--output-name", os.Getenv("CLOUD_ASSESS_TEST_PROCESS_OUTPUT")}
	switch name {
	case "success":
	case "failure":
		input.Completeness = assessment.CompletenessFailed
		input.Findings = nil
		executionErr = errors.New("fixture graph failure")
	case "severity":
		args = append(args, "--fail-on", "High")
	case "partial":
		input.Completeness = assessment.CompletenessPartial
	default:
		t.Fatalf("unknown child case %q", name)
	}
	runner := app.NewRunner(processAssessmentFixture{report: result.Build(input), err: executionErr})
	os.Exit(runWithExecutor(args, func(ctx context.Context, flags scanFlags) (int, error) {
		outcome, err := runner.Run(ctx, app.ScanOptions{
			Outputs: app.OutputOptions{BaseName: flags.outputName, JSON: flags.json, XLSX: flags.xlsx},
			FailOn:  flags.failOn,
		})
		return outcome.ExitCode, err
	}))
}
