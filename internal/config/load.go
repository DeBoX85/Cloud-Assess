package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadFilters reads the optional assessment filter YAML file and rebuilds runtime indexes.
// An empty filename returns the default filter configuration.
func LoadFilters(filename string) (*Filters, error) {
	filters := NewFilters()
	if filename == "" {
		return filters, nil
	}

	cleanPath := filepath.Clean(filename)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("read filter file %q: %w", filename, err)
	}
	if err := yaml.Unmarshal(data, filters); err != nil {
		return nil, fmt.Errorf("parse filter YAML %q: %w", filename, err)
	}
	filters.RebuildIndexes()
	if filters.Assessment != nil {
		if err := filters.Assessment.Validate(); err != nil {
			return nil, fmt.Errorf("validate filter file %q: %w", filename, err)
		}
	}
	return filters, nil
}
