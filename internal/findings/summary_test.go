package findings

import (
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestBuildDeduplicatesAndSkipsSLA(t *testing.T) {
	definitions := []assessment.RecommendationDefinition{
		{ID: "rec-1", Recommendation: "Fix widget", Category: "Security", Impact: assessment.ImpactHigh, ResourceType: "Microsoft.Test/widgets"},
		{ID: "sla", Category: assessment.CategorySLA, Impact: assessment.ImpactLow},
	}
	results := []assessment.Finding{
		{RecommendationID: "rec-1", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one", Category: "Security", Impact: assessment.ImpactHigh},
		{RecommendationID: "REC-1", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/ONE", Category: "Security", Impact: assessment.ImpactHigh},
		{RecommendationID: "sla", ResourceID: "resource", Category: assessment.CategorySLA, Impact: assessment.ImpactLow},
	}

	summary := Build(definitions, results, 4)
	if len(summary.Recommendations) != 1 {
		t.Fatalf("got %d recommendations, want 1", len(summary.Recommendations))
	}
	if summary.Recommendations[0].ImpactedResources != 1 {
		t.Fatalf("got %d impacted resources, want 1", summary.Recommendations[0].ImpactedResources)
	}
	if summary.Resources != 4 || summary.ImpactedResources != 1 || summary.RecommendationsFound != 1 {
		t.Fatalf("unexpected summary totals: %+v", summary)
	}
	if summary.ImpactedByImpact[assessment.ImpactHigh] != 1 {
		t.Fatalf("unexpected impact totals: %+v", summary.ImpactedByImpact)
	}
}

func TestBuildIgnoresFindingWithoutDefinition(t *testing.T) {
	summary := Build(nil, []assessment.Finding{{
		RecommendationID: "external",
		Recommendation:   "External recommendation",
		ResourceID:       "resource",
		ResourceType:     "Microsoft.Test/widgets",
		Category:         "Governance",
		Impact:           assessment.ImpactMedium,
	}}, 1)
	if len(summary.Recommendations) != 0 || summary.ImpactedResources != 0 {
		t.Fatalf("unexpected summary for undefined recommendation: %+v", summary)
	}
}

func TestBuildSortsByResourceTypeThenID(t *testing.T) {
	definitions := []assessment.RecommendationDefinition{
		{ID: "z-2", ResourceType: "Microsoft.Z/type", Category: "Security", Impact: assessment.ImpactLow},
		{ID: "a-2", ResourceType: "Microsoft.A/type", Category: "Security", Impact: assessment.ImpactLow},
		{ID: "a-1", ResourceType: "Microsoft.A/type", Category: "Security", Impact: assessment.ImpactLow},
	}
	summary := Build(definitions, nil, 0)
	got := []string{summary.Recommendations[0].ID, summary.Recommendations[1].ID, summary.Recommendations[2].ID}
	want := []string{"a-1", "a-2", "z-2"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("sort order = %v, want %v", got, want)
		}
	}
}
