package config

import (
	"strings"
	"testing"
)

func FuzzExactResourceExclusionPreservesOtherResource(f *testing.F) {
	f.Add("widget")
	f.Add("é")
	f.Fuzz(func(t *testing.T, name string) {
		if name == "" || len(name) > 256 || strings.ContainsAny(name, "/ \t\r\n") {
			return
		}
		prefix := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/"
		id, other := prefix+name, prefix+name+"-other"
		filters := NewFilters()
		filters.Assessment.Exclude.Resources = []string{id}
		filters.RebuildIndexes()
		filters.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
		if !filters.Assessment.IsResourceExcluded(id, nil) || filters.Assessment.IsResourceExcluded(other, nil) {
			t.Fatal("exact exclusion changed an unrelated resource")
		}
		if !filters.Assessment.IsServiceExcluded(id) || filters.Assessment.IsServiceExcluded(other) {
			t.Fatal("inventory and downstream structural scope disagree")
		}
	})
}
