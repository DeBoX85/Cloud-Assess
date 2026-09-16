package json

import (
	stdjson "encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

func sampleResult() *result.AssessmentResult {
	return result.Build(result.Input{
		GeneratedAt:  time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
		ScopeID:      "scope-123",
		Completeness: assessment.CompletenessCompleteWithWarnings,
		Stages: []assessment.StageExecution{
			{
				Name:   "graph",
				Status: assessment.StageCompletedWithWarnings,
				Warnings: []assessment.AssessmentWarning{
					{Code: "sample_warning", Message: "sample warning"},
				},
			},
		},
		Recommendations: []assessment.RecommendationDefinition{
			{
				ID:             "rec-001",
				Recommendation: "Use secure configuration",
				Category:       "Security",
				Impact:         "High",
				ResourceType:   "microsoft.storage/storageaccounts",
				Source:         "CUSTOM",
				Query:          "Resources | where false",
			},
		},
		Findings: []assessment.Finding{
			{
				RecommendationID: "rec-001",
				Category:         "Security",
				Impact:           "High",
				ResourceType:     "microsoft.storage/storageaccounts",
				ResourceID:       "/subscriptions/sub-1/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1",
				SubscriptionID:   "sub-1",
				ResourceGroup:    "rg",
				ResourceName:     "st1",
				Source:           "CUSTOM",
			},
		},
		Resources: []assessment.Resource{
			{
				ID:             "/subscriptions/sub-1/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1",
				SubscriptionID: "sub-1",
				ResourceGroup:  "rg",
				Type:           "microsoft.storage/storageaccounts",
				Name:           "st1",
			},
		},
	})
}

func TestMarshalIncludesCanonicalAssessmentState(t *testing.T) {
	encoded, err := Marshal(sampleResult())
	if err != nil {
		t.Fatal(err)
	}

	var decoded map[string]any
	if err := stdjson.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"schemaVersion", "generatedAt", "scopeId", "completeness", "stages", "summary", "recommendations", "findings", "resources"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("expected top-level key %q", key)
		}
	}
	if got := decoded["schemaVersion"]; got != result.SchemaVersion {
		t.Fatalf("schemaVersion = %v, want %s", got, result.SchemaVersion)
	}
	if got := decoded["completeness"]; got != string(assessment.CompletenessCompleteWithWarnings) {
		t.Fatalf("completeness = %v", got)
	}
	if strings.Contains(string(encoded), "Resources | where false") || strings.Contains(string(encoded), "\"query\"") {
		t.Fatal("internal recommendation query leaked into JSON output")
	}
}

func TestMarshalIsDeterministicForCanonicalResult(t *testing.T) {
	data := sampleResult()
	first, err := Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("JSON output was not deterministic")
	}
}

func TestStringMatchesMarshal(t *testing.T) {
	data := sampleResult()
	encoded, err := Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	text, err := String(data)
	if err != nil {
		t.Fatal(err)
	}
	if text != string(encoded) {
		t.Fatal("stdout JSON representation differs from file representation")
	}
}

func TestWriteFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "assessment.json")
	if err := WriteFile(sampleResult(), filename); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := Marshal(sampleResult())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(expected) {
		t.Fatal("file content differs from canonical JSON encoding")
	}
}

func TestNilAndEmptyFilenameReturnErrors(t *testing.T) {
	if _, err := Marshal(nil); err == nil {
		t.Fatal("expected nil result error")
	}
	if err := WriteFile(sampleResult(), ""); err == nil {
		t.Fatal("expected empty filename error")
	}
}
