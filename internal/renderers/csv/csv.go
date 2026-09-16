package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

// Write creates one CSV file per applicable report table and returns the generated paths.
// The familiar reference suffixes are preserved, with assessmentStatus added by Cloud Assess.
func Write(data *result.AssessmentResult, baseFilename string, opts tables.Options) ([]string, error) {
	if data == nil {
		return nil, fmt.Errorf("assessment result is nil")
	}
	if baseFilename == "" {
		return nil, fmt.Errorf("CSV output base filename is empty")
	}

	projected, err := tables.Build(data, opts)
	if err != nil {
		return nil, err
	}
	generated := make([]string, 0, len(projected))
	for _, table := range projected {
		if !tables.ShouldRender(data, table) {
			continue
		}
		filename := fmt.Sprintf("%s.%s.csv", baseFilename, table.Key)
		if err := writeTable(filename, table.Rows); err != nil {
			return generated, err
		}
		generated = append(generated, filename)
	}
	return generated, nil
}

func writeTable(filename string, rows [][]string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("create CSV report %q: %w", filename, err)
	}
	writer := csv.NewWriter(file)
	writer.WriteAll(rows)
	writeErr := writer.Error()
	closeErr := file.Close()
	if writeErr != nil {
		return fmt.Errorf("write CSV report %q: %w", filename, writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close CSV report %q: %w", filename, closeErr)
	}
	return nil
}
