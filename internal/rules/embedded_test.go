package rules

import "testing"

func TestLoadPinnedCatalogLoadsAllSources(t *testing.T) {
	catalog, counts, err := LoadPinnedCatalog()
	if err != nil {
		t.Fatalf("LoadPinnedCatalog() error = %v", err)
	}
	if catalog == nil {
		t.Fatal("catalog is nil")
	}
	if got, want := len(counts), 3; got != want {
		t.Fatalf("source count entries = %d, want %d", got, want)
	}

	for _, sourceCount := range counts {
		if sourceCount.Count == 0 {
			t.Fatalf("source %s loaded zero recommendations", sourceCount.Source)
		}
	}
}

func TestPinnedCatalogContainsKnownCustomRecommendation(t *testing.T) {
	catalog, _, err := LoadPinnedCatalog()
	if err != nil {
		t.Fatalf("LoadPinnedCatalog() error = %v", err)
	}

	definitions := catalog.ByResourceType("Microsoft.AAD/domainServices")
	for _, definition := range definitions {
		if definition.ID != "domain-003" {
			continue
		}
		if definition.Source != SourceCustom {
			t.Fatalf("domain-003 source = %q, want %q", definition.Source, SourceCustom)
		}
		if definition.Query == "" {
			t.Fatal("domain-003 should have its matching KQL attached")
		}
		if definition.Category != "SLA" || definition.Impact != "High" {
			t.Fatalf("unexpected domain-003 metadata: %+v", definition)
		}
		return
	}
	t.Fatal("domain-003 not found in pinned custom rule corpus")
}

func TestPinnedProvenanceMatchesImportedSnapshots(t *testing.T) {
	provenance := PinnedProvenance()
	if len(provenance) != 3 {
		t.Fatalf("provenance entries = %d, want 3", len(provenance))
	}
	if provenance[0].Source != SourceAPRL || provenance[0].Revision != APRLCommit {
		t.Fatalf("unexpected APRL provenance: %+v", provenance[0])
	}
	if provenance[1].Source != SourceAOR || provenance[1].Revision != AORSnapshotTree {
		t.Fatalf("unexpected AOR provenance: %+v", provenance[1])
	}
	if provenance[2].Source != SourceCustom || provenance[2].Revision != CustomRulesSnapshotTree {
		t.Fatalf("unexpected custom provenance: %+v", provenance[2])
	}
}
