package gate

import (
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/findings"
)

func TestParseThresholds(t *testing.T) {
	for _, value := range []string{"High", "medium", " LOW "} {
		if _, err := Parse(value); err != nil {
			t.Fatalf("Parse(%q) returned %v", value, err)
		}
	}
	if _, err := Parse("critical"); err == nil {
		t.Fatal("expected invalid severity error")
	}
}

func TestEvaluateAtOrAboveThreshold(t *testing.T) {
	summary := &findings.Summary{Recommendations: []findings.RecommendationSummary{
		{ID: "high", Impact: assessment.ImpactHigh, ImpactedResources: 2},
		{ID: "medium", Impact: assessment.ImpactMedium, ImpactedResources: 3},
		{ID: "low", Impact: assessment.ImpactLow, ImpactedResources: 4},
	}}
	criterion, _ := Parse("Medium")
	result := Evaluate(summary, criterion)
	if result.Recommendations != 2 || result.Impacted != 5 {
		t.Fatalf("Evaluate() = %+v, want 2 recommendations / 5 impacted", result)
	}
}

func TestCheckReturnsTypedFailure(t *testing.T) {
	summary := &findings.Summary{Recommendations: []findings.RecommendationSummary{{Impact: assessment.ImpactHigh, ImpactedResources: 1}}}
	criterion, _ := Parse("High")
	err := Check(summary, criterion)
	failure, ok := err.(*Failure)
	if !ok {
		t.Fatalf("expected *Failure, got %T (%v)", err, err)
	}
	if failure.Threshold != assessment.ImpactHigh || failure.Recommendations != 1 || failure.Impacted != 1 {
		t.Fatalf("unexpected failure: %+v", failure)
	}
}

func TestNilSummaryPasses(t *testing.T) {
	criterion, _ := Parse("Low")
	if err := Check(nil, criterion); err != nil {
		t.Fatalf("nil summary should not fail gate: %v", err)
	}
}
