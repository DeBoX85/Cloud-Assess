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
	filters.Assessment.Include.ResourceTypes = []string{"Microsoft.Test/widgets"}
	filters.Assessment.Include.Tags = map[string]string{"env": "prod", "team": "platform"}
	filters.Assessment.Exclude.Tags = map[string]string{"lifecycle": "retired"}
	filters.RebuildIndexes()

	id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	rg := "/subscriptions/sub/resourceGroups/rg"

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
			got := filters.Assessment.IsResourceExcluded(id, "sub", rg, "microsoft.test/widgets", tt.tags)
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
	filters.RebuildIndexes()

	if filters.Assessment.IsSubscriptionExcluded("sub-a") {
		t.Fatal("explicit include should win over subscription exclusion, matching reference behavior")
	}
	if !filters.Assessment.IsSubscriptionExcluded("sub-b") {
		t.Fatal("subscription outside include set should be excluded")
	}
}

func TestTagScopeAppliesToChildResource(t *testing.T) {
	filters := NewFilters()
	filters.Assessment.Include.Tags = map[string]string{"env": "prod"}
	filters.RebuildIndexes()
	parent := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	filters.Assessment.SetResourceScope(parent, true)
	if filters.Assessment.IsServiceExcluded(parent + "/slots/blue") {
		t.Fatal("child resource should inherit included parent scope")
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
