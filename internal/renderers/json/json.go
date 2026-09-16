package json

import (
	stdjson "encoding/json"
	"fmt"
	"os"

	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

// Marshal returns the canonical assessment result as stable, human-readable JSON.
func Marshal(data *result.AssessmentResult) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("assessment result is nil")
	}
	return stdjson.MarshalIndent(data, "", "\t")
}

// String returns the same JSON representation used by file output.
func String(data *result.AssessmentResult) (string, error) {
	encoded, err := Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// WriteFile writes the canonical JSON report to filename.
func WriteFile(data *result.AssessmentResult, filename string) error {
	if filename == "" {
		return fmt.Errorf("JSON output filename is empty")
	}
	encoded, err := Marshal(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, encoded, 0o600); err != nil {
		return fmt.Errorf("write JSON report %q: %w", filename, err)
	}
	return nil
}
