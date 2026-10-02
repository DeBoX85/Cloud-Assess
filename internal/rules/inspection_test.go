package rules

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestInspectionSelectionAndPrecedence(t *testing.T) {
	catalog := NewCatalog()
	definition := func(id, resourceType, query, state string) assessment.RecommendationDefinition {
		return assessment.RecommendationDefinition{ID: id, ResourceType: resourceType, Query: query, State: state, Recommendation: resourceType}
	}
	catalog.AddAll([]assessment.RecommendationDefinition{
		definition("z-empty-query", "Type/A", "", ""), // Source inspection includes metadata-only rules.
		definition("collision", "Type/A", "resources", ""),
		definition("cross-type", "Type/A", "resources", ""),
		definition("cross-type", "Type/B", "resources", ""),
		definition("unsupported", "Type/C", "resources", ""),
		definition("disabled", "Type/A", "resources", "Disabled"),
		definition("manual", "Type/A", "// cannot-be-validated-with-arg", ""),
		definition("development", "Type/A", "// under-development", ""),
		definition("development-space", "Type/A", "// under development", ""),
	})
	diagnostic := definition("collision", "type/a", "", "disabled")
	diagnostic.Recommendation = "diagnostic wins"
	diagnostic.LearnMore = []assessment.LearnMoreLink{{URL: "https://example.test/first"}, {URL: "https://example.test/second"}}
	rows := Inspect(catalog, []string{"TYPE/A", "Type/B"}, []assessment.RecommendationDefinition{diagnostic})
	want := []InspectionRow{
		{RecommendationID: "collision", ResourceType: "type/a", Recommendation: "diagnostic wins", LearnMoreURL: "https://example.test/first"},
		{RecommendationID: "cross-type", ResourceType: "Type/B", Recommendation: "Type/B"},
		{RecommendationID: "z-empty-query", ResourceType: "Type/A", Recommendation: "Type/A"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("inspection selection = %#v, want %#v", rows, want)
	}
	if catalog.ResourceTypeCount() != 3 {
		t.Fatal("heading count must include catalog types outside scanner selection")
	}
	rows[0].Recommendation = "changed"
	if got := Inspect(catalog, []string{"Type/A", "Type/B"}, []assessment.RecommendationDefinition{diagnostic}); !reflect.DeepEqual(got, want) {
		t.Fatal("inspection exposed mutable output state")
	}
}

func TestEmptyInspectionIsJSONArray(t *testing.T) {
	rows := Inspect(NewCatalog(), nil, nil)
	data, err := json.Marshal(rows)
	if err != nil || string(data) != "[]" {
		t.Fatalf("empty projection = %s, %v", data, err)
	}
}
