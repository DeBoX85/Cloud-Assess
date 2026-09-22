package equivalence

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/diagnostics"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

const (
	testSubscriptionID = "11111111-2222-3333-4444-555555555555"
	testResourceID     = "/subscriptions/11111111-2222-3333-4444-555555555555/resourceGroups/RG/providers/Microsoft.Storage/storageAccounts/ST1"
)

func TestReferenceAndTargetProjectionNormalizeIntentionalDifferences(t *testing.T) {
	reference, err := LoadReference(bytes.NewReader(referenceFixtureJSON(t)))
	if err != nil {
		t.Fatal(err)
	}
	target, err := LoadTarget(bytes.NewReader(targetFixtureJSON(t, false)))
	if err != nil {
		t.Fatal(err)
	}

	report := Compare(reference, target)
	if !report.Equivalent {
		encoded, _ := json.MarshalIndent(report, "", "  ")
		t.Fatalf("expected semantic equivalence, got:\n%s", encoded)
	}
}

func TestCompareSurfacesChangedSemanticField(t *testing.T) {
	reference, err := LoadReference(bytes.NewReader(referenceFixtureJSON(t)))
	if err != nil {
		t.Fatal(err)
	}
	target, err := LoadTarget(bytes.NewReader(targetFixtureJSON(t, true)))
	if err != nil {
		t.Fatal(err)
	}

	report := Compare(reference, target)
	if report.Equivalent {
		t.Fatal("expected impact change to break equivalence")
	}
	diff := datasetDiff(report, DatasetFindings)
	if len(diff.Changed) != 1 {
		t.Fatalf("changed findings = %#v, want one", diff.Changed)
	}
	if diff.ReferenceCount != 1 || diff.TargetCount != 1 || diff.MissingCount != 0 || diff.ExtraCount != 0 || diff.ChangedCount != 1 {
		t.Fatalf("finding summary counts = %#v", diff)
	}
	if len(diff.Changed[0].Fields) != 1 || diff.Changed[0].Fields[0].Field != "impact" {
		t.Fatalf("finding field deltas = %#v, want impact only", diff.Changed[0].Fields)
	}
}

func TestPartialTargetIsNotComparable(t *testing.T) {
	data := result.Build(result.Input{
		Completeness: assessment.CompletenessPartial,
		Stages: []assessment.StageExecution{{
			Name:   stages.Graph,
			Status: assessment.StageCompleted,
		}},
	})
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	target, err := LoadTarget(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	reference, err := LoadReference(bytes.NewReader([]byte(`{
		"recommendations": [],
		"impacted": [],
		"resourceType": [],
		"inventory": [],
		"outOfScope": []
	}`)))
	if err != nil {
		t.Fatal(err)
	}

	report := Compare(reference, target)
	if report.Equivalent {
		t.Fatal("partial target must not be declared equivalent")
	}
	if len(report.Preconditions) == 0 {
		t.Fatal("expected comparability precondition failure")
	}
}

func TestCoverageMismatchIsExplicit(t *testing.T) {
	reference, err := LoadReference(bytes.NewReader([]byte(`{
		"recommendations": [],
		"impacted": [],
		"resourceType": [],
		"inventory": [],
		"outOfScope": []
	}`)))
	if err != nil {
		t.Fatal(err)
	}
	data := result.Build(result.Input{
		Completeness: assessment.CompletenessComplete,
		Stages: []assessment.StageExecution{
			{Name: stages.Graph, Status: assessment.StageCompleted},
			{Name: stages.Advisor, Status: assessment.StageCompleted},
		},
	})
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	target, err := LoadTarget(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}

	report := Compare(reference, target)
	diff := datasetDiff(report, DatasetAdvisor)
	if !diff.CoverageMismatch || diff.ReferenceEnabled || !diff.TargetEnabled {
		t.Fatalf("advisor coverage diff = %#v", diff)
	}
}

func TestRedactedReportsAreNotComparable(t *testing.T) {
	reference, err := LoadReference(bytes.NewReader([]byte(`{
		"recommendations": [],
		"impacted": [{
			"recommendationId": "rec-1",
			"resourceId": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxx5555555/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st",
			"subscriptionId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxx5555555"
		}],
		"resourceType": [],
		"inventory": [],
		"outOfScope": []
	}`)))
	if err != nil {
		t.Fatal(err)
	}
	if reference.Comparable {
		t.Fatal("redacted reference report must not be comparable")
	}

	targetJSON := targetFixtureJSON(t, false)
	targetJSON = bytes.ReplaceAll(targetJSON, []byte(testSubscriptionID), []byte("xxxxxxxx-xxxx-xxxx-xxxx-xxxxx5555555"))
	target, err := LoadTarget(bytes.NewReader(targetJSON))
	if err != nil {
		t.Fatal(err)
	}
	if target.Comparable {
		t.Fatal("redacted target report must not be comparable")
	}

	report := Compare(reference, target)
	if report.Equivalent || len(report.Preconditions) != 2 {
		t.Fatalf("redacted comparison = %#v", report)
	}
}

func referenceFixtureJSON(t *testing.T) []byte {
	t.Helper()
	payload := map[string]any{
		"recommendations": []map[string]string{{
			"implemented":               "false",
			"numberOfImpactedResources": "1",
			"recommendationSource":      "AZQR",
			"category":                  "MonitoringAndAlerting",
			"recommendation":            "Storage should have diagnostic settings enabled",
			"impact":                    "Low",
			"bestPracticesGuidance":     "",
			"readMore":                  "https://learn.example/diagnostics",
			"recommendationId":          "st-001",
		}},
		"impacted": []map[string]string{{
			"validatedUsing":   "Azure Resource Graph",
			"source":           "AZQR",
			"category":         "MonitoringAndAlerting",
			"impact":           "Low",
			"resourceType":     "Microsoft.Storage/storageAccounts",
			"recommendation":   "Storage should have diagnostic settings enabled",
			"recommendationId": "st-001",
			"subscriptionId":   testSubscriptionID,
			"subscriptionName": "",
			"resourceGroup":    "RG",
			"resourceName":     "ST1",
			"resourceId":       testResourceID,
			"param1":           "",
			"param2":           "",
			"param3":           "",
			"param4":           "",
			"param5":           "",
			"learn":            "https://learn.example/diagnostics",
		}},
		"resourceType": []map[string]string{{
			"subscriptionName":  "Test Subscription",
			"resourceType":      "Microsoft.Storage/storageAccounts",
			"numberOfResources": "1",
		}},
		"inventory": []map[string]string{{
			"subscriptionId": testSubscriptionID,
			"resourceGroup":  "RG",
			"location":       "norwayeast",
			"resourceType":   "Microsoft.Storage/storageAccounts",
			"resourceName":   "ST1",
			"skuName":        "",
			"skuTier":        "",
			"capacity":       "",
			"kind":           "",
			"sla":            "99.9%",
			"resourceId":     testResourceID,
		}},
		"outOfScope": []map[string]string{},
		"costs": []map[string]string{{
			"from":             "2026-08-01",
			"to":               "2026-08-31",
			"subscriptionId":   testSubscriptionID,
			"subscriptionName": "",
			"serviceName":      "Storage",
			"value":            "12.3",
			"currency":         "NOK",
		}},
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func targetFixtureJSON(t *testing.T, changedImpact bool) []byte {
	t.Helper()
	impact := assessment.ImpactLow
	if changedImpact {
		impact = assessment.ImpactMedium
	}
	definition := assessment.RecommendationDefinition{
		ID:             "st-001",
		Recommendation: "Storage should have diagnostic settings enabled",
		Category:       diagnostics.CategoryMonitoringAndAlerting,
		Impact:         assessment.ImpactLow,
		ResourceType:   "microsoft.storage/storageaccounts",
		Source:         diagnostics.Source,
		LearnMore: []assessment.LearnMoreLink{{
			Name: "Diagnostic Settings",
			URL:  "https://learn.example/diagnostics",
		}},
	}
	data := result.Build(result.Input{
		GeneratedAt:  time.Date(2026, 9, 18, 0, 0, 0, 0, time.UTC),
		Completeness: assessment.CompletenessComplete,
		Stages: []assessment.StageExecution{
			{Name: stages.Graph, Status: assessment.StageCompleted},
			{Name: stages.Diagnostics, Status: assessment.StageCompleted},
			{Name: stages.Advisor, Status: assessment.StageSkipped},
			{Name: stages.Defender, Status: assessment.StageSkipped},
			{Name: stages.DefenderRecommendations, Status: assessment.StageSkipped},
			{Name: stages.Policy, Status: assessment.StageSkipped},
			{Name: stages.Arc, Status: assessment.StageSkipped},
			{Name: stages.Cost, Status: assessment.StageCompleted},
		},
		Recommendations: []assessment.RecommendationDefinition{definition},
		Findings: []assessment.Finding{
			{
				RecommendationID:    "st-001",
				Source:              diagnostics.Source,
				ValidationMechanism: diagnostics.ValidationAzureResourceManager,
				Category:            diagnostics.CategoryMonitoringAndAlerting,
				Impact:              impact,
				ResourceType:        "microsoft.storage/storageaccounts",
				Recommendation:      "Storage should have diagnostic settings enabled",
				ResourceID:          testResourceID,
				SubscriptionID:      testSubscriptionID,
				SubscriptionName:    "Test Subscription",
				ResourceGroup:       "RG",
				ResourceName:        "ST1",
				LearnMoreURL:        "https://learn.example/diagnostics",
			},
			{
				RecommendationID: "sla-test",
				Category:         assessment.CategorySLA,
				ResourceID:       testResourceID,
				Parameters:       []string{"99.9%"},
			},
		},
		Resources: []assessment.Resource{{
			ID:             testResourceID,
			SubscriptionID: testSubscriptionID,
			ResourceGroup:  "RG",
			Location:       "norwayeast",
			Type:           "Microsoft.Storage/storageAccounts",
			Name:           "ST1",
		}},
		ResourceTypes: []assessment.ResourceTypeCount{{
			SubscriptionID:   testSubscriptionID,
			SubscriptionName: "Test Subscription",
			ResourceType:     "Microsoft.Storage/storageAccounts",
			Count:            1,
		}},
		Costs: []assessment.CostRecord{{
			From:             time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			To:               time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
			SubscriptionID:   testSubscriptionID,
			SubscriptionName: "Test Subscription",
			ServiceName:      "Storage",
			Value:            "12.300000",
			Currency:         "NOK",
		}},
	})
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func TestFindingSubscriptionNameNormalizationIsDiagnosticsOnly(t *testing.T) {
	if got := normalizeFindingSubscriptionName("st-001", "Development"); got != "" {
		t.Fatalf("diagnostics subscription name = %q, want empty comparison value", got)
	}
	if got := normalizeFindingSubscriptionName("aprl-001", "Development"); got != "Development" {
		t.Fatalf("ordinary finding subscription name = %q, want Development", got)
	}
}

func TestNormalizeSourcePreservesStrictProvenanceBoundaries(t *testing.T) {
	tests := []struct {
		name             string
		recommendationID string
		source           string
		want             string
	}{
		{name: "diagnostics", recommendationID: "st-001", source: "AZQR", want: "DIAGNOSTICS"},
		{name: "legacy custom", recommendationID: "domain-003", source: "AZQR", want: "CUSTOM"},
		{name: "APRL", recommendationID: "st-003", source: "APRL", want: "APRL"},
		{name: "AOR", recommendationID: "orphan-id", source: "AOR", want: "AOR"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := normalizeSource(test.recommendationID, test.source); got != test.want {
				t.Fatalf("normalizeSource(%q, %q) = %q, want %q", test.recommendationID, test.source, got, test.want)
			}
		})
	}
}

func datasetDiff(report Report, name string) DatasetDiff {
	for _, diff := range report.Datasets {
		if diff.Name == name {
			return diff
		}
	}
	return DatasetDiff{Name: name}
}
