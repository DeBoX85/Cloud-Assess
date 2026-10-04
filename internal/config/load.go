package config

import (
	"bytes"
	"fmt"
	"io"
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
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(filters); err != nil && err != io.EOF {
		return nil, fmt.Errorf("parse filter YAML %q: %w", filename, err)
	}
	if filters.Assessment == nil {
		return nil, fmt.Errorf("validate filter file %q: assessment must not be null", filename)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("parse filter YAML %q: expected one YAML document", filename)
	}
	filters.RebuildIndexes()
	if filters.Assessment != nil {
		if err := filters.Assessment.Validate(); err != nil {
			return nil, fmt.Errorf("validate filter file %q: %w", filename, err)
		}
	}
	return filters, nil
}
