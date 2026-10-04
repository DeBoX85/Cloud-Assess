package main

import (
	"path/filepath"
	"testing"
)

func TestAuditComparisonRejectsUnprovenInputs(t *testing.T) {
	tests := []struct {
		name      string
		reference string
		target    string
	}{
		{
			name:      "missing-completeness",
			reference: `{"inventory":[],"outOfScope":[],"resourceType":[],"recommendations":[],"impacted":[]}`,
			target:    `{"schemaVersion":"1.0","stages":[{"name":"graph","status":"completed"}]}`,
		},
		{
			name:      "unknown-completeness",
			reference: `{"inventory":[],"outOfScope":[],"resourceType":[],"recommendations":[],"impacted":[]}`,
			target:    `{"schemaVersion":"1.0","completeness":"unknown","stages":[{"name":"graph","status":"completed"}]}`,
		},
		{
			name:      "unsupported-schema",
			reference: `{"inventory":[],"outOfScope":[],"resourceType":[],"recommendations":[],"impacted":[]}`,
			target:    `{"schemaVersion":"future-schema","completeness":"complete","stages":[{"name":"graph","status":"completed"}]}`,
		},
		{
			name:      "empty-reference-object",
			reference: `{}`,
			target:    `{"schemaVersion":"1.0","completeness":"complete","stages":[{"name":"graph","status":"skipped"}]}`,
		},
		{
			name:      "null-reference",
			reference: `null`,
			target:    `{"schemaVersion":"1.0","completeness":"complete","stages":[{"name":"graph","status":"skipped"}]}`,
		},
		{
			name:      "unknown-reference-sections",
			reference: `{"unrecognized":[]}`,
			target:    `{"schemaVersion":"1.0","completeness":"complete","stages":[{"name":"graph","status":"skipped"}]}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			reference := filepath.Join(directory, "reference.json")
			target := filepath.Join(directory, "target.json")
			output := filepath.Join(directory, "diff.json")
			writeTestFile(t, reference, test.reference)
			writeTestFile(t, target, test.target)
			if code := run([]string{"--reference", reference, "--target", target, "--output", output}); code == 0 {
				t.Fatal("invalid or unproven comparison input was declared equivalent")
			}
		})
	}
}

func TestAuditComparisonPreservesHealthyEmptyEvidence(t *testing.T) {
	for _, completeness := range []string{"complete", "complete_with_warnings"} {
		t.Run(completeness, func(t *testing.T) {
			directory := t.TempDir()
			reference := filepath.Join(directory, "reference.json")
			target := filepath.Join(directory, "target.json")
			output := filepath.Join(directory, "diff.json")
			writeTestFile(t, reference, `{"advisor":[]}`)
			writeTestFile(t, target, `{"schemaVersion":"1.0","completeness":"`+completeness+`","stages":[{"name":"advisor","status":"completed"}]}`)
			if code := run([]string{"--reference", reference, "--target", target, "--output", output}); code != 0 {
				t.Fatalf("healthy empty evidence rejected with exit code %d", code)
			}
		})
	}
}
