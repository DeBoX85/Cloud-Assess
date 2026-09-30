package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/reportfile"
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
	if err := reportfile.Write(filename, func(file *os.File) error {
		writer := csv.NewWriter(file)
		for _, row := range rows {
			safe := make([]string, len(row))
			for i, value := range row {
				safe[i] = spreadsheetText(value)
			}
			if err := writer.Write(safe); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	}); err != nil {
		return fmt.Errorf("write CSV report %q: %w", filename, err)
	}
	return nil
}

// CSV quoting does not prevent spreadsheet formula interpretation. Prefix risky
// cells as text in this human-oriented export; canonical JSON remains unchanged.
func spreadsheetText(value string) string {
	trimmed := strings.TrimLeftFunc(value, unicode.IsSpace)
	if strings.HasPrefix(value, "\t") || strings.HasPrefix(value, "\r") || strings.HasPrefix(value, "\n") || strings.HasPrefix(trimmed, "=") || strings.HasPrefix(trimmed, "+") || strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "@") || strings.HasPrefix(trimmed, "＝") || strings.HasPrefix(trimmed, "＋") || strings.HasPrefix(trimmed, "－") || strings.HasPrefix(trimmed, "＠") {
		return "'" + value
	}
	return value
}
