package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunReturnsZeroForEquivalentReports(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	target := filepath.Join(dir, "target.json")
	output := filepath.Join(dir, "diff.json")

	writeTestFile(t, reference, `{
		"recommendations": [],
		"impacted": [],
		"resourceType": [],
		"inventory": [],
		"outOfScope": []
	}`)
	writeTestFile(t, target, `{
		"schemaVersion": "1.0",
		"completeness": "complete",
		"stages": [{"name": "graph", "status": "completed"}]
	}`)

	if code := run([]string{"--reference", reference, "--target", target, "--output", output}); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	var report struct {
		Equivalent bool `json:"equivalent"`
	}
	decodeTestJSON(t, output, &report)
	if !report.Equivalent {
		t.Fatal("expected equivalent report")
	}
}

func TestRunReturnsOneForSemanticDifference(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	target := filepath.Join(dir, "target.json")
	output := filepath.Join(dir, "diff.json")

	writeTestFile(t, reference, `{
		"recommendations": [],
		"impacted": [{
			"recommendationId": "rec-1",
			"resourceId": "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st",
			"source": "APRL",
			"category": "Security",
			"impact": "High"
		}],
		"resourceType": [],
		"inventory": [],
		"outOfScope": []
	}`)
	writeTestFile(t, target, `{
		"schemaVersion": "1.0",
		"completeness": "complete",
		"stages": [{"name": "graph", "status": "completed"}]
	}`)

	if code := run([]string{"--reference", reference, "--target", target, "--output", output}); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}

	var report struct {
		Equivalent bool `json:"equivalent"`
	}
	decodeTestJSON(t, output, &report)
	if report.Equivalent {
		t.Fatal("expected semantic difference")
	}
}

func TestRunReturnsTwoForInvalidInvocation(t *testing.T) {
	if code := run(nil); code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func decodeTestJSON(t *testing.T, path string, target any) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(content, target); err != nil {
		t.Fatal(err)
	}
}
