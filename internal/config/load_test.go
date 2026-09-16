package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFiltersDefaults(t *testing.T) {
	filters, err := LoadFilters("")
	if err != nil {
		t.Fatalf("LoadFilters(\"\") error = %v", err)
	}
	if filters == nil || filters.Assessment == nil || filters.Assessment.Include == nil || filters.Assessment.Exclude == nil {
		t.Fatalf("LoadFilters(\"\") returned incomplete defaults: %#v", filters)
	}
}

func TestLoadFiltersYAML(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "filters.yaml")
	content := []byte(`assessment:
  include:
    subscriptions: ["SUB-1"]
    resourceGroups: ["/subscriptions/SUB-1/resourceGroups/RG-1"]
    resourceTypes: ["vm"]
    tags:
      Environment: prod
  exclude:
    recommendations: ["REC-1"]
`)
	if err := os.WriteFile(filename, content, 0o600); err != nil {
		t.Fatal(err)
	}

	filters, err := LoadFilters(filename)
	if err != nil {
		t.Fatalf("LoadFilters() error = %v", err)
	}
	if filters.Assessment.IsSubscriptionExcluded("sub-1") {
		t.Fatal("included subscription was excluded")
	}
	if !filters.Assessment.IsSubscriptionExcluded("sub-2") {
		t.Fatal("non-included subscription was retained")
	}
	if !filters.Assessment.IsRecommendationExcluded("rec-1") {
		t.Fatal("excluded recommendation was retained")
	}
	if got := filters.Assessment.Include.ResourceTypes; len(got) != 1 || got[0] != "vm" {
		t.Fatalf("resourceTypes = %#v, want [vm]", got)
	}
}

func TestLoadFiltersRejectsInvalidInput(t *testing.T) {
	directory := t.TempDir()
	if _, err := LoadFilters(filepath.Join(directory, "missing.yaml")); err == nil {
		t.Fatal("missing filter file error = nil, want error")
	}

	badYAML := filepath.Join(directory, "bad.yaml")
	if err := os.WriteFile(badYAML, []byte("assessment: ["), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFilters(badYAML); err == nil {
		t.Fatal("invalid YAML error = nil, want error")
	}

	badRG := filepath.Join(directory, "bad-rg.yaml")
	if err := os.WriteFile(badRG, []byte("assessment:\n  include:\n    resourceGroups: [not-an-arm-id]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFilters(badRG); err == nil {
		t.Fatal("invalid resource group error = nil, want error")
	}
}
