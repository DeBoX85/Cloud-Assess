package rules

import (
	"embed"
	"fmt"
	"io/fs"
)

// The normal reference scan loads APRL Azure resources, AOR, and custom rules.
// APRL specialized-workload files remain available in the pinned submodule but are
// intentionally not part of this catalog unless they are explicitly promoted later.
//
//go:embed upstream/aprl/azure-resources
//go:embed upstream/orphan-resources
//go:embed custom/azure-resources
var pinnedRuleFiles embed.FS

type SourceCount struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

// LoadPinnedCatalog loads the exact rule-source snapshots selected for the v1 equivalence baseline.
// Precedence matches the reference: APRL, then AOR, then custom rules.
func LoadPinnedCatalog() (*Catalog, []SourceCount, error) {
	catalog := NewCatalog()
	counts := make([]SourceCount, 0, 3)

	providers := []struct {
		path   string
		source string
	}{
		{path: "upstream/aprl/azure-resources", source: SourceAPRL},
		{path: "upstream/orphan-resources", source: SourceAOR},
		{path: "custom/azure-resources", source: SourceCustom},
	}

	for _, provider := range providers {
		providerFS, err := fs.Sub(pinnedRuleFiles, provider.path)
		if err != nil {
			return nil, nil, fmt.Errorf("open %s rule source: %w", provider.source, err)
		}
		definitions, err := LoadFilesystem(providerFS, provider.source)
		if err != nil {
			return nil, nil, err
		}
		counts = append(counts, SourceCount{Source: provider.source, Count: len(definitions)})
		catalog.AddAll(definitions)
	}

	return catalog, counts, nil
}
