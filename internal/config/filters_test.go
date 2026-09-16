package config

import "testing"

func TestValidateResourceGroupID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		ok   bool
	}{
		{"valid", "/subscriptions/123/resourceGroups/test-rg", true},
		{"name only", "test-rg", false},
		{"missing leading slash", "subscriptions/123/resourceGroups/test-rg", false},
		{"empty subscription", "/subscriptions//resourceGroups/test-rg", false},
		{"empty resource group", "/subscriptions/123/resourceGroups/", false},
		{"extra segment", "/subscriptions/123/resourceGroups/test-rg/providers", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceGroupID(tt.id)
			if tt.ok && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestTagFiltersPreserveReferenceSemantics(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Include.Tags = map[string]string{"env": "prod", "team": "platform"}
	filters.Assessment.Exclude.Tags = map[string]string{"lifecycle": "retired"}
	filters.RebuildIndexes()
	filters.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})

	id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	tests := []struct {
		name string
		tags map[string]string
		want bool
	}{
		{"all include tags match", map[string]string{"ENV": "prod", "Team": "platform"}, false},
		{"include missing", map[string]string{"env": "prod"}, true},
		{"tag values exact", map[string]string{"env": "Prod", "team": "platform"}, true},
		{"exclude wins", map[string]string{"env": "prod", "team": "platform", "lifecycle": "retired"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filters.Assessment.IsResourceExcluded(id, tt.tags)
			if got != tt.want {
				t.Fatalf("IsResourceExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncludeAndExcludePrecedence(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Include.Subscriptions = []string{"SUB-A"}
	filters.Assessment.Exclude.Subscriptions = []string{"sub-a", "sub-b"}
	filters.Assessment.Include.ResourceGroups = []string{"/subscriptions/sub-a/resourceGroups/keep"}
	filters.Assessment.Exclude.ResourceGroups = []string{"/subscriptions/sub-a/resourceGroups/keep"}
	filters.RebuildIndexes()

	if filters.Assessment.IsSubscriptionExcluded("sub-a") {
		t.Fatal("explicit include should win over subscription exclusion, matching reference behavior")
	}
	if !filters.Assessment.IsSubscriptionExcluded("sub-b") {
		t.Fatal("subscription outside include set should be excluded")
	}
	if filters.Assessment.IsResourceGroupExcluded("/subscriptions/sub-a/resourceGroups/keep") {
		t.Fatal("explicit resource-group include should win over matching exclusion")
	}
	if !filters.Assessment.IsResourceGroupExcluded("/subscriptions/sub-a/resourceGroups/other") {
		t.Fatal("resource group outside include set should be excluded")
	}
}

func TestStructuralScopeIsReappliedToDownstreamFindings(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Include.Subscriptions = []string{"sub-a"}
	filters.Assessment.Include.ResourceGroups = []string{"/subscriptions/sub-a/resourceGroups/keep"}
	excludedID := "/subscriptions/sub-a/resourceGroups/keep/providers/Microsoft.Test/widgets/excluded"
	filters.Assessment.Exclude.Resources = []string{excludedID}
	filters.RebuildIndexes()
	filters.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})

	inScope := "/subscriptions/sub-a/resourceGroups/keep/providers/Microsoft.Test/widgets/one"
	if filters.Assessment.IsServiceExcluded(inScope) {
		t.Fatal("in-scope resource should remain included")
	}
	if !filters.Assessment.IsServiceExcluded("/subscriptions/sub-a/resourceGroups/keep/providers/Microsoft.Other/things/one") {
		t.Fatal("resource type outside scanner scope should be excluded")
	}
	if !filters.Assessment.IsServiceExcluded("/subscriptions/sub-b/resourceGroups/keep/providers/Microsoft.Test/widgets/one") {
		t.Fatal("resource outside included subscription should be excluded")
	}
	if !filters.Assessment.IsServiceExcluded("/subscriptions/sub-a/resourceGroups/other/providers/Microsoft.Test/widgets/one") {
		t.Fatal("resource outside included resource group should be excluded")
	}
	if !filters.Assessment.IsServiceExcluded(excludedID) {
		t.Fatal("explicit resource exclusion should apply downstream")
	}
}

func TestScannerScopeMustBeInstalledBeforeResourceFiltering(t *testing.T) {
	filters := NewFilters()
	id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	if !filters.Assessment.IsResourceExcluded(id, nil) {
		t.Fatal("unconfigured scanner scope should fail closed, matching reference LoadFilters assumptions")
	}
	filters.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
	if filters.Assessment.IsResourceExcluded(id, nil) {
		t.Fatal("allowed resource type should be included after scanner scope is installed")
	}
}

func TestTagScopeAppliesToChildResource(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Include.Tags = map[string]string{"env": "prod"}
	filters.RebuildIndexes()
	filters.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
	parent := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	filters.Assessment.SetResourceScope(parent, true)
	if filters.Assessment.IsServiceExcluded(parent + "/slots/blue") {
		t.Fatal("child resource should inherit included parent scope")
	}
}

func TestExcludeOnlyTagScopeKeepsUnknownDownstreamResource(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Exclude.Tags = map[string]string{"lifecycle": "retired"}
	filters.RebuildIndexes()
	filters.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
	id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/unknown"
	if filters.Assessment.IsServiceExcluded(id) {
		t.Fatal("unknown tag scope should remain included for exclude-only tag filters")
	}
}

func TestRecommendationExclusionIsCaseInsensitive(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Exclude.Recommendations = []string{"REC-001"}
	filters.RebuildIndexes()
	if !filters.Assessment.IsRecommendationExcluded("rec-001") {
		t.Fatal("recommendation exclusion should be case-insensitive")
	}
}
