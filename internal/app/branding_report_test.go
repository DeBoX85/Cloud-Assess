package app

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
	"github.com/xuri/excelize/v2"
)

// Test-only export path lets the native paired-profile checker inspect real files.
var brandingEvidenceDir = flag.String("branding-evidence-dir", "", "Synthetic branding report evidence directory (tests only)")

func brandingReportFixture() *result.AssessmentResult {
	const sub = "11111111-2222-3333-4444-555555555555"
	const resourceID = "/subscriptions/" + sub + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store"
	fixed := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	input := result.Input{
		GeneratedAt: fixed, ScopeID: "branding-fixture", Completeness: assessment.CompletenessComplete,
		Resources:               []assessment.Resource{{ID: resourceID, SubscriptionID: sub, ResourceGroup: "rg", Type: "microsoft.storage/storageaccounts", Name: "store", Tags: map[string]string{"Environment": "dev"}}},
		OutOfScope:              []assessment.Resource{{ID: resourceID + "-excluded", SubscriptionID: sub, ResourceGroup: "rg", Type: "microsoft.storage/storageaccounts", Name: "store-excluded"}},
		ResourceTypes:           []assessment.ResourceTypeCount{{SubscriptionID: sub, ResourceType: "microsoft.storage/storageaccounts", Count: 1}},
		Advisor:                 []assessment.AdvisorRecommendation{{RecommendationID: "advisor-1", SubscriptionID: sub, ResourceID: resourceID, Category: "Security", Impact: "High", Description: "Synthetic Advisor"}},
		Defender:                []assessment.DefenderPlanStatus{{SubscriptionID: sub, Name: "StorageAccounts", Tier: "Standard"}},
		DefenderRecommendations: []assessment.DefenderRecommendation{{SubscriptionID: sub, ResourceID: resourceID, RecommendationName: "Synthetic Defender", RecommendationSeverity: "High"}},
		AzurePolicy:             []assessment.PolicyNonCompliance{{SubscriptionID: sub, ResourceID: resourceID, PolicyDefinitionID: "policy-1", PolicyDisplayName: "Synthetic Policy", ComplianceState: "NonCompliant"}},
		ArcSQL:                  []assessment.ArcSQLRecord{{SubscriptionID: sub, SQLInstance: "synthetic-sql", VCores: "4", Version: "2022"}},
		Costs:                   []assessment.CostRecord{{SubscriptionID: sub, From: fixed.AddDate(0, -1, 0), To: fixed, ServiceName: "Storage", Value: "12.34", Currency: "NOK"}},
	}
	for _, name := range []string{stages.Graph, stages.Diagnostics, stages.Advisor, stages.Defender, stages.DefenderRecommendations, stages.Policy, stages.Arc, stages.Cost} {
		input.Stages = append(input.Stages, assessment.StageExecution{Name: name, Status: assessment.StageCompleted, Records: 1, StartedAt: fixed, FinishedAt: fixed.Add(time.Second)})
	}
	for _, source := range []string{"APRL", "AOR", "CUSTOM"} {
		id := "synthetic-" + strings.ToLower(source)
		input.Recommendations = append(input.Recommendations, assessment.RecommendationDefinition{ID: id, Source: source, ResourceType: "microsoft.storage/storageaccounts", Recommendation: "Synthetic " + source + " control", Impact: assessment.ImpactHigh, Category: "Security", LearnMore: []assessment.LearnMoreLink{{URL: "https://learn.microsoft.com/azure/storage/common/storage-introduction"}}})
		input.Findings = append(input.Findings, assessment.Finding{RecommendationID: id, Source: source, ResourceID: resourceID, SubscriptionID: sub, ResourceGroup: "rg", ResourceType: "microsoft.storage/storageaccounts", ResourceName: "store", Recommendation: "Synthetic " + source + " control", Impact: assessment.ImpactHigh, Category: "Security", ValidationMechanism: "Azure Resource Graph", LearnMoreURL: "https://learn.microsoft.com/azure/storage/common/storage-introduction"})
	}
	return result.Build(input)
}

func TestBrandingReportEvidence(t *testing.T) {
	directory := *brandingEvidenceDir
	if directory == "" {
		directory = t.TempDir()
	} else if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	directory, err := filepath.Abs(directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(directory)
	data := brandingReportFixture()
	before, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	fixed := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		name, base string
		redact     bool
	}{
		{"default", "", true}, {"explicit", filepath.Join(directory, "operator override"), true}, {"raw", filepath.Join(directory, "raw"), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			runner := NewRunner(&fakeAssessmentRunner{result: data})
			runner.now = func() time.Time { return fixed }
			var stdout bytes.Buffer
			runner.stdout = &stdout
			outcome, err := runner.Run(context.Background(), ScanOptions{Outputs: OutputOptions{BaseName: test.base, XLSX: true, JSON: true, CSV: true, SARIF: true, Stdout: true, RedactSubscriptionIDs: test.redact, Version: "branding-test"}})
			if err != nil || outcome.ExitCode != ExitSuccess {
				t.Fatalf("run=%d %v", outcome.ExitCode, err)
			}
			base := test.base
			if base == "" {
				base = branding.Default().ReportFilePrefix + "_2026_10_01_T120000"
			}
			if len(outcome.Files) != 15 {
				t.Fatalf("files=%d want 12 CSV+XLSX+JSON+SARIF", len(outcome.Files))
			}
			for _, path := range outcome.Files {
				if !strings.HasPrefix(path, base+".") {
					t.Fatalf("unexpected output path %q", path)
				}
				info, err := os.Stat(path)
				if err != nil {
					t.Fatal(err)
				}
				if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
					t.Fatalf("report mode=%o", info.Mode().Perm())
				}
			}
			encoded, err := os.ReadFile(base + ".json")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(bytes.TrimSpace(stdout.Bytes()), bytes.TrimSpace(encoded)) {
				t.Fatal("stdout and canonical file differ")
			}
			if err := os.WriteFile(test.name+".stdout.json", stdout.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			book, err := excelize.OpenFile(base + ".xlsx")
			if err != nil {
				t.Fatal(err)
			}
			defer book.Close()
			if len(book.GetSheetList()) != 12 {
				t.Fatalf("sheets=%v", book.GetSheetList())
			}
			for _, sheet := range book.GetSheetList() {
				title, err := book.GetCellValue(sheet, "A1")
				if err != nil || title != branding.Default().ReportTitle {
					t.Fatalf("sheet %s title=%q %v", sheet, title, err)
				}
				formula, err := book.GetCellFormula(sheet, "A1")
				if err != nil || formula != "" {
					t.Fatalf("title executed as formula %q %v", formula, err)
				}
			}
		})
	}
	// Separate runners share only immutable profile/result inputs, never output paths.
	var wg sync.WaitGroup
	for i := range 4 {
		wg.Go(func() {
			runner := NewRunner(&fakeAssessmentRunner{result: data})
			var stdout bytes.Buffer
			runner.stdout = &stdout
			base := filepath.Join(directory, fmt.Sprintf("concurrent-%d", i))
			outcome, err := runner.Run(context.Background(), ScanOptions{Outputs: OutputOptions{BaseName: base, XLSX: true, JSON: true, SARIF: true, Stdout: true, RedactSubscriptionIDs: true, Version: "branding-test"}})
			if err != nil || outcome.ExitCode != ExitSuccess {
				t.Errorf("concurrent run %d: %v", i, err)
			}
		})
	}
	wg.Wait()
	after, err := json.Marshal(data)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("rendering mutated assessment input", err)
	}
}
