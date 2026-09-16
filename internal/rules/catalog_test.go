package rules

import (
	"testing"
	"testing/fstest"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestLoadFilesystemPairsRecommendationWithKQL(t *testing.T) {
	fsys := fstest.MapFS{
		"service/recommendations.yaml": {Data: []byte(`
- description: Fix the widget
  aprlGuid: rec-001
  recommendationTypeId: null
  recommendationControl: Security
  recommendationImpact: High
  recommendationResourceType: Microsoft.Test/widgets
  recommendationMetadataState: Active
  longDescription: Long guidance
  potentialBenefits: Better security
  pgVerified: true
  automationAvailable: true
  tags: [security, test]
  learnMoreLink:
  - name: First
    url: https://example.test/first
  - name: Second
    url: https://example.test/second
`)},
		"service/kql/rec-001.kql": {Data: []byte("resources | where type =~ 'microsoft.test/widgets'")},
	}

	definitions, err := LoadFilesystem(fsys, SourceCustom)
	if err != nil {
		t.Fatalf("LoadFilesystem() error = %v", err)
	}
	if got, want := len(definitions), 1; got != want {
		t.Fatalf("definition count = %d, want %d", got, want)
	}
	definition := definitions[0]
	if definition.ID != "rec-001" || definition.Source != SourceCustom || definition.ValidationMechanism != ValidationARG {
		t.Fatalf("unexpected definition identity: %+v", definition)
	}
	if definition.Query == "" {
		t.Fatal("matching KQL was not attached")
	}
	if definition.AutomationAvailable != "true" {
		t.Fatalf("AutomationAvailable = %q, want true", definition.AutomationAvailable)
	}
	if !definition.PGVerified {
		t.Fatal("PGVerified should be preserved")
	}
	if len(definition.LearnMore) != 2 || definition.LearnMore[0].Name != "First" || definition.LearnMore[1].Name != "Second" {
		t.Fatalf("learn-more ordering not preserved: %#v", definition.LearnMore)
	}
}

func TestCatalogLaterProviderWinsSameResourceAndID(t *testing.T) {
	catalog := NewCatalog()
	catalog.Add(assessment.RecommendationDefinition{ID: "rec", ResourceType: "Microsoft.Test/widgets", Source: SourceAPRL})
	catalog.Add(assessment.RecommendationDefinition{ID: "rec", ResourceType: "Microsoft.Test/widgets", Source: SourceAOR})
	catalog.Add(assessment.RecommendationDefinition{ID: "rec", ResourceType: "Microsoft.Test/widgets", Source: SourceCustom})

	definitions := catalog.ByResourceType("microsoft.test/WIDGETS")
	if len(definitions) != 1 || definitions[0].Source != SourceCustom {
		t.Fatalf("later provider did not win: %#v", definitions)
	}
}

func TestEmbeddedExecutionFiltering(t *testing.T) {
	base := assessment.RecommendationDefinition{ID: "rec", State: "Active", Query: "resources | project id"}
	if !IsExecutableEmbedded(base, nil) {
		t.Fatal("normal embedded rule should execute")
	}

	cases := []assessment.RecommendationDefinition{
		{ID: "rec", State: "disabled", Query: "resources"},
		{ID: "rec", State: "Active", Query: "// cannot-be-validated-with-arg"},
		{ID: "rec", State: "Active", Query: "// under-development"},
		{ID: "rec", State: "Active", Query: "// under development"},
	}
	for _, definition := range cases {
		if IsExecutableEmbedded(definition, nil) {
			t.Fatalf("rule should have been skipped: %+v", definition)
		}
	}

	if IsExecutableEmbedded(base, func(id string) bool { return id == "rec" }) {
		t.Fatal("excluded recommendation should not execute")
	}
}

func TestPluginExecutionOnlyHonorsRecommendationExclusion(t *testing.T) {
	pluginRule := assessment.RecommendationDefinition{
		ID:    "plugin-rec",
		State: "disabled",
		Query: "// cannot-be-validated-with-arg",
	}
	if !IsExecutablePlugin(pluginRule, nil) {
		t.Fatal("plugin rule should preserve reference behavior and remain executable")
	}
	if IsExecutablePlugin(pluginRule, func(id string) bool { return id == "plugin-rec" }) {
		t.Fatal("excluded plugin recommendation should not execute")
	}
}

func TestCatalogOrderingIsDeterministic(t *testing.T) {
	catalog := NewCatalog()
	catalog.Add(assessment.RecommendationDefinition{ID: "b", ResourceType: "Microsoft.B/type"})
	catalog.Add(assessment.RecommendationDefinition{ID: "z", ResourceType: "Microsoft.A/type"})
	catalog.Add(assessment.RecommendationDefinition{ID: "a", ResourceType: "Microsoft.A/type"})

	all := catalog.All()
	if len(all) != 3 || all[0].ID != "a" || all[1].ID != "z" || all[2].ID != "b" {
		t.Fatalf("unexpected deterministic ordering: %#v", all)
	}
}
