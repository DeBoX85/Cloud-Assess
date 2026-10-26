package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These literal rows deliberately disagree on impact at the same normalized
// identity. Neither row order nor agreement on the final row proves equivalence.
func TestComparisonRejectsConflictingDuplicateEvidence(t *testing.T) {
	high := `{"resourceId":"/synthetic/private-resource","recommendationId":"rec-1","impact":"High"}`
	low := `{"resourceId":"/SYNTHETIC/PRIVATE-RESOURCE","recommendationId":"REC-1","impact":"Low"}`
	for _, side := range []string{"reference", "target", "both"} {
		for _, rows := range []string{high + "," + low, low + "," + high, high + "," + low + "," + low} {
			t.Run(side+"/"+rows, func(t *testing.T) {
				refRows, targetRows := low, low
				if side != "target" {
					refRows = rows
				}
				if side != "reference" {
					targetRows = rows
				}
				directory := t.TempDir()
				reference, target, output := filepath.Join(directory, "reference.json"), filepath.Join(directory, "target.json"), filepath.Join(directory, "diff.json")
				writeTestFile(t, reference, `{"advisor":[`+refRows+`]}`)
				writeTestFile(t, target, `{"schemaVersion":"1.0","completeness":"complete","advisor":[`+targetRows+`]}`)
				code := run([]string{"--reference", reference, "--target", target, "--output", output})
				if code != 1 {
					t.Fatalf("conflicting duplicate evidence did not fail comparison with exit 1: got %d", code)
				}
				encoded, err := os.ReadFile(output)
				if err != nil {
					t.Fatal(err)
				}
				var report struct {
					Equivalent                  bool
					Preconditions               []string
					ReferenceNotes, TargetNotes []string
				}
				if err := json.Unmarshal(encoded, &report); err != nil {
					t.Fatal(err)
				}
				if report.Equivalent || len(report.Preconditions) == 0 {
					t.Fatal("conflicting duplicate evidence lacks comparability failure")
				}
				if side != "target" && len(report.ReferenceNotes) == 0 {
					t.Fatal("reference conflict lacks diagnostic note")
				}
				if side != "reference" && len(report.TargetNotes) == 0 {
					t.Fatal("target conflict lacks diagnostic note")
				}
				for _, note := range append(report.ReferenceNotes, report.TargetNotes...) {
					if strings.Contains(strings.ToLower(note), "private-resource") {
						t.Fatal("duplicate note exposes record identity")
					}
				}
			})
		}
	}
}

func TestComparisonPreservesNormalizedIdenticalDuplicates(t *testing.T) {
	directory := t.TempDir()
	reference, target, output := filepath.Join(directory, "reference.json"), filepath.Join(directory, "target.json"), filepath.Join(directory, "diff.json")
	writeTestFile(t, reference, `{"advisor":[{"resourceId":"/synthetic/r","recommendationId":"rec-1","impact":"High"},{"resourceId":" /SYNTHETIC/R ","recommendationId":" REC-1 ","impact":" HIGH "}]}`)
	writeTestFile(t, target, `{"schemaVersion":"1.0","completeness":"complete","advisor":[{"resourceId":"/synthetic/r","recommendationId":"rec-1","impact":"high"}]}`)
	if code := run([]string{"--reference", reference, "--target", target, "--output", output}); code != 0 {
		t.Fatalf("normalized identical duplicate rejected: exit %d", code)
	}
}

func TestComparisonPreservesResourceTypeCountAggregation(t *testing.T) {
	directory := t.TempDir()
	reference, target, output := filepath.Join(directory, "reference.json"), filepath.Join(directory, "target.json"), filepath.Join(directory, "diff.json")
	writeTestFile(t, reference, `{"recommendations":[],"impacted":[],"inventory":[],"outOfScope":[],"resourceType":[{"subscriptionName":"Synthetic","resourceType":"Microsoft.Storage/storageAccounts","numberOfResources":"2"},{"subscriptionName":"Synthetic","resourceType":"Microsoft.Storage/storageAccounts","numberOfResources":"3"}]}`)
	writeTestFile(t, target, `{"schemaVersion":"1.0","completeness":"complete","stages":[{"name":"graph","status":"completed"}],"resourceTypes":[{"subscriptionName":"Synthetic","resourceType":"Microsoft.Storage/storageAccounts","count":5}]}`)
	if code := run([]string{"--reference", reference, "--target", target, "--output", output}); code != 0 {
		t.Fatalf("documented count aggregation rejected: exit %d", code)
	}
}
