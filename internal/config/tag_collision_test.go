package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTagConfigurationAliasesRejected(t *testing.T) {
	for _, side := range []string{"include", "exclude"} {
		for _, fixture := range []struct{ name, left, right, value string }{
			{"case-conflict", "env", "ENV", "dev"},
			{"trim-conflict", "env", " env ", "dev"},
			{"same-value", "env", "ENV", "prod"},
			{"unicode-case", "ÅREA", "årea", "dev"},
			{"unicode-trim", "env", "\u2003env\u2003", "dev"},
		} {
			t.Run(side+"/"+fixture.name, func(t *testing.T) {
				f := NewFilters()
				tags := map[string]string{fixture.left: "prod", fixture.right: fixture.value}
				if side == "include" {
					f.Assessment.Include.Tags = tags
				} else {
					f.Assessment.Exclude.Tags = tags
				}
				for i := 0; i < 16; i++ {
					f.RebuildIndexes()
					err := f.Assessment.Validate()
					if err == nil || !strings.Contains(err.Error(), side+".tags") {
						t.Fatal("TAG_ALIAS_VALIDATION: ambiguous configuration accepted or wrong map identified")
					}
					if len(tags) != 2 {
						t.Fatal("TAG_ALIAS_OWNERSHIP: raw map changed")
					}
				}
			})
		}
	}
}

func TestLoadTagAliasesRejected(t *testing.T) {
	for _, side := range []string{"include", "exclude"} {
		t.Run(side, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "filter.yaml")
			body := "assessment:\n  " + side + ":\n    tags:\n      env: private-value-one\n      ENV: private-value-two\n"
			if err := os.WriteFile(p, []byte(body), 0600); err != nil {
				t.Fatal(err)
			}
			f, err := LoadFilters(p)
			if err == nil || f != nil {
				t.Fatal("TAG_ALIAS_LOAD: ambiguous YAML accepted")
			}
			if !strings.Contains(err.Error(), side+".tags") || strings.Contains(err.Error(), "private-value") {
				t.Fatal("TAG_ALIAS_ERROR: wrong map or tag value in diagnostic")
			}
		})
	}
}

func TestInvalidTagScopeCloneAndRecovery(t *testing.T) {
	const id = "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	for _, side := range []string{"include", "exclude"} {
		t.Run(side, func(t *testing.T) {
			f := NewFilters()
			if side == "include" {
				f.Assessment.Include.Tags = map[string]string{"env": "prod", "ENV": "dev"}
			} else {
				f.Assessment.Exclude.Tags = map[string]string{"env": "prod", "ENV": "dev"}
			}
			f.RebuildIndexes()
			f.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
			f.Assessment.SetResourceScope(id, true)
			for _, current := range []*AssessmentFilter{f.Assessment, f.Assessment.Clone()} {
				if !current.IsResourceExcluded(id, map[string]string{"env": "prod"}) || !current.IsServiceExcluded(id) || !current.IsServiceExcluded(id+"/slots/child") {
					t.Fatal("TAG_ALIAS_SCOPE: invalid configuration widened discovery or recorded downstream scope")
				}
			}
			f.Assessment.Include.Tags = nil
			f.Assessment.Exclude.Tags = nil
			f.RebuildIndexes()
			f.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
			if f.Assessment.Validate() != nil || f.Assessment.IsResourceExcluded(id, nil) || f.Assessment.IsServiceExcluded(id) {
				t.Fatal("TAG_ALIAS_RECOVERY: corrected configuration remained invalid")
			}
		})
	}
}

func TestHealthyTagOverlapAndExactValues(t *testing.T) {
	f := NewFilters()
	f.Assessment.Include.Tags = map[string]string{" ENV ": "prod"}
	f.Assessment.Exclude.Tags = map[string]string{"env": "dev"}
	f.RebuildIndexes()
	f.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
	if err := f.Assessment.Validate(); err != nil {
		t.Fatal(err)
	}
	const id = "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	for _, x := range []struct {
		value    string
		excluded bool
	}{{"prod", false}, {"Prod", true}, {"dev", true}} {
		if f.Assessment.IsResourceExcluded(id, map[string]string{"Env": x.value}) != x.excluded {
			t.Fatal("TAG_HEALTHY: exact value or cross-map behavior changed")
		}
	}
	f.Assessment.Exclude.Tags = map[string]string{"env": "prod"}
	f.RebuildIndexes()
	f.Assessment.SetAllowedResourceTypes([]string{"Microsoft.Test/widgets"})
	if !f.Assessment.IsResourceExcluded(id, map[string]string{"env": "prod"}) {
		t.Fatal("TAG_HEALTHY: exclude precedence lost")
	}
}
