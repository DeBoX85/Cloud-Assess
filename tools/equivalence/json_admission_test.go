package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComparisonJSONAdmissionFailsBeforeOutput(t *testing.T) {
	for _, side := range []string{"reference", "target"} {
		t.Run(side, func(t *testing.T) {
			dir := t.TempDir()
			ref, target, output := filepath.Join(dir, "reference.json"), filepath.Join(dir, "target.json"), filepath.Join(dir, "output.json")
			referenceJSON := `{"advisor":[]}`
			targetJSON := `{"schemaVersion":"1.0","completeness":"complete","stages":[{"name":"advisor","status":"completed"}]}`
			if side == "reference" {
				referenceJSON = `{"advisor":[{"impact":"High"}],"advisor":[]}`
			} else {
				targetJSON = `{"schemaVersion":"1.0","completeness":"partial","COMPLETENESS":"complete","stages":[{"name":"advisor","status":"completed"}]}`
			}
			writeTestFile(t, ref, referenceJSON)
			writeTestFile(t, target, targetJSON)
			writeTestFile(t, output, "original output must remain")
			if code := run([]string{"--reference", ref, "--target", target, "--output", output}); code != 2 {
				t.Fatalf("JSON_ADMISSION_CLI: ambiguous input exit = %d, want 2", code)
			}
			b, err := os.ReadFile(output)
			if err != nil || string(b) != "original output must remain" {
				t.Fatal("JSON_ADMISSION_OUTPUT: rejected comparison replaced existing output")
			}
		})
	}
}
