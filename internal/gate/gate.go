package gate

import (
	"fmt"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/findings"
)

type Criterion struct {
	Threshold string
	rank      int
}

type Result struct {
	Recommendations int
	Impacted        int
}

type Failure struct {
	Threshold       string
	Recommendations int
	Impacted        int
}

func (e *Failure) Error() string {
	return fmt.Sprintf("quality gate failed: %d recommendation(s) at or above %s impact with %d impacted resource(s)", e.Recommendations, e.Threshold, e.Impacted)
}

func Parse(value string) (Criterion, error) {
	value = strings.TrimSpace(value)
	rank, ok := assessment.SeverityRank(value)
	if !ok {
		return Criterion{}, fmt.Errorf("invalid fail-on impact %q: supported values are High, Medium, Low", value)
	}
	threshold := assessment.ImpactLow
	switch rank {
	case 3:
		threshold = assessment.ImpactHigh
	case 2:
		threshold = assessment.ImpactMedium
	}
	return Criterion{Threshold: threshold, rank: rank}, nil
}

func Evaluate(summary *findings.Summary, criterion Criterion) Result {
	var result Result
	if summary == nil {
		return result
	}
	for _, recommendation := range summary.Recommendations {
		if recommendation.ImpactedResources == 0 {
			continue
		}
		rank, ok := assessment.SeverityRank(recommendation.Impact)
		if !ok || rank < criterion.rank {
			continue
		}
		result.Recommendations++
		result.Impacted += recommendation.ImpactedResources
	}
	return result
}

func Check(summary *findings.Summary, criterion Criterion) error {
	result := Evaluate(summary, criterion)
	if result.Recommendations == 0 {
		return nil
	}
	return &Failure{Threshold: criterion.Threshold, Recommendations: result.Recommendations, Impacted: result.Impacted}
}
