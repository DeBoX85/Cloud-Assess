package result

import (
	"reflect"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestBuildCreatesCanonicalSummaryAndMetadata(t *testing.T) {
	generated := time.Date(2026, 9, 16, 15, 30, 0, 0, time.FixedZone("CEST", 2*60*60))
	definitions := []assessment.RecommendationDefinition{
		{ID: "rec-2", ResourceType: "Microsoft.Storage/storageAccounts", Recommendation: "Second", Category: "Security", Impact: "High", Source: "APRL"},
		{ID: "rec-1", ResourceType: "Microsoft.Compute/virtualMachines", Recommendation: "First", Category: "Reliability", Impact: "Medium", Source: "APRL"},
	}
	findingRows := []assessment.Finding{
		{RecommendationID: "rec-2", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1", SubscriptionID: "sub", Category: "Security"},
	}

	got := Build(Input{
		GeneratedAt:     generated,
		ScopeID:         "scope-hash",
		Completeness:    assessment.CompletenessComplete,
		Recommendations: definitions,
		Findings:        findingRows,
		Resources: []assessment.Resource{
			{ID: "resource-2"},
			{ID: "resource-1"},
		},
	})

	if got.SchemaVersion != SchemaVersion || got.ScopeID != "scope-hash" || got.Completeness != assessment.CompletenessComplete {
		t.Fatalf("unexpected metadata: %+v", got)
	}
	if got.GeneratedAt.Location() != time.UTC || !got.GeneratedAt.Equal(generated) {
		t.Fatalf("generatedAt = %s, want same instant in UTC", got.GeneratedAt)
	}
	if got.Summary == nil || got.Summary.Resources != 2 || got.Summary.RecommendationsFound != 1 || got.Summary.ImpactedResources != 1 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if len(got.Recommendations) != 2 || got.Recommendations[0].ID != "rec-1" || got.Recommendations[1].ID != "rec-2" {
		t.Fatalf("recommendations are not canonical: %+v", got.Recommendations)
	}
}

func TestBuildDoesNotMutateOrShareCallerOwnedNestedData(t *testing.T) {
	stages := []assessment.StageExecution{{
		Name:     "graph",
		Warnings: []assessment.AssessmentWarning{{Code: "warning", Message: "original"}},
		Error:    &assessment.AssessmentError{Code: "error", Message: "original"},
	}}
	recommendations := []assessment.RecommendationDefinition{{
		ID:        "rec",
		Tags:      []string{"tag-a"},
		LearnMore: []assessment.LearnMoreLink{{Name: "Docs", URL: "https://example.test"}},
	}}
	findings := []assessment.Finding{{RecommendationID: "rec", ResourceID: "resource", Parameters: []string{"p1"}}}
	resources := []assessment.Resource{{ID: "resource", Tags: map[string]string{"env": "prod"}}}

	got := Build(Input{
		Stages:          stages,
		Recommendations: recommendations,
		Findings:        findings,
		Resources:       resources,
	})

	got.Stages[0].Warnings[0].Message = "changed"
	got.Stages[0].Error.Message = "changed"
	got.Recommendations[0].Tags[0] = "changed"
	got.Recommendations[0].LearnMore[0].URL = "changed"
	got.Findings[0].Parameters[0] = "changed"
	got.Resources[0].Tags["env"] = "changed"

	if stages[0].Warnings[0].Message != "original" || stages[0].Error.Message != "original" {
		t.Fatal("result shares stage nested data with caller")
	}
	if recommendations[0].Tags[0] != "tag-a" || recommendations[0].LearnMore[0].URL != "https://example.test" {
		t.Fatal("result shares recommendation nested data with caller")
	}
	if findings[0].Parameters[0] != "p1" || resources[0].Tags["env"] != "prod" {
		t.Fatal("result shares finding/resource nested data with caller")
	}
}

func TestBuildCanonicalOrderingDoesNotModifyInputSlices(t *testing.T) {
	resources := []assessment.Resource{
		{ID: "/subscriptions/b/resourceGroups/z/providers/Microsoft.Storage/storageAccounts/b", SubscriptionID: "b", ResourceGroup: "z", Type: "Microsoft.Storage/storageAccounts", Name: "b"},
		{ID: "/subscriptions/a/resourceGroups/a/providers/Microsoft.Compute/virtualMachines/a", SubscriptionID: "a", ResourceGroup: "a", Type: "Microsoft.Compute/virtualMachines", Name: "a"},
	}
	original := append([]assessment.Resource(nil), resources...)

	got := Build(Input{Resources: resources})
	if len(got.Resources) != 2 || got.Resources[0].SubscriptionID != "a" || got.Resources[1].SubscriptionID != "b" {
		t.Fatalf("unexpected canonical resource order: %+v", got.Resources)
	}
	if !reflect.DeepEqual(resources, original) {
		t.Fatalf("caller resource slice was reordered: got %+v want %+v", resources, original)
	}
}

func TestBuildRetainsAndSortsAuxiliaryDatasets(t *testing.T) {
	got := Build(Input{
		Advisor: []assessment.AdvisorRecommendation{
			{SubscriptionID: "b", ResourceID: "z", RecommendationID: "2"},
			{SubscriptionID: "a", ResourceID: "x", RecommendationID: "1"},
		},
		Defender: []assessment.DefenderPlanStatus{
			{SubscriptionID: "b", Name: "VM"},
			{SubscriptionID: "a", Name: "Storage"},
		},
		AzurePolicy: []assessment.PolicyNonCompliance{
			{SubscriptionID: "b", ResourceID: "z", PolicyDefinitionID: "2"},
			{SubscriptionID: "a", ResourceID: "x", PolicyDefinitionID: "1"},
		},
		ArcSQL: []assessment.ArcSQLRecord{
			{SubscriptionID: "b", SQLInstance: "z"},
			{SubscriptionID: "a", SQLInstance: "x"},
		},
		Costs: []assessment.CostRecord{
			{SubscriptionID: "b", ServiceName: "Storage", Value: "2", Currency: "USD"},
			{SubscriptionID: "a", ServiceName: "Compute", Value: "1", Currency: "USD"},
		},
	})

	if got.Advisor[0].SubscriptionID != "a" || got.Defender[0].SubscriptionID != "a" || got.AzurePolicy[0].SubscriptionID != "a" || got.ArcSQL[0].SubscriptionID != "a" || got.Costs[0].SubscriptionID != "a" {
		t.Fatalf("auxiliary datasets are not canonical: %+v", got)
	}
}
