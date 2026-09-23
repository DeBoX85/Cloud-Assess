package csv

import (
	encodingcsv "encoding/csv"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

func csvFixture() *result.AssessmentResult {
	return result.Build(result.Input{
		GeneratedAt:  time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
		ScopeID:      "scope-csv",
		Completeness: assessment.CompletenessComplete,
		Stages: []assessment.StageExecution{
			{Name: stages.Graph, Status: assessment.StageCompleted, Records: 1},
			{Name: stages.Advisor, Status: assessment.StageSkipped},
			{Name: stages.Defender, Status: assessment.StageCompleted, Records: 1},
		},
		Recommendations: []assessment.RecommendationDefinition{
			{ID: "rec-1", Recommendation: "A recommendation, with comma", Category: "Security", Impact: "High", ResourceType: "microsoft.storage/storageaccounts", Source: "CUSTOM"},
		},
		Findings: []assessment.Finding{
			{RecommendationID: "rec-1", Recommendation: "A recommendation, with comma", Category: "Security", Impact: "High", ResourceType: "microsoft.storage/storageaccounts", Source: "CUSTOM", ResourceID: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1", SubscriptionID: "11111111-1111-1111-1111-111111111111", ResourceGroup: "rg", ResourceName: "st1"},
		},
		Resources: []assessment.Resource{
			{ID: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1", SubscriptionID: "11111111-1111-1111-1111-111111111111", ResourceGroup: "rg", Type: "microsoft.storage/storageaccounts", Name: "st1"},
		},
		ResourceTypes: []assessment.ResourceTypeCount{{SubscriptionID: "11111111-1111-1111-1111-111111111111", SubscriptionName: "Sub One", ResourceType: "microsoft.storage/storageaccounts", Count: 1}},
		Defender:      []assessment.DefenderPlanStatus{{SubscriptionID: "11111111-1111-1111-1111-111111111111", SubscriptionName: "Sub One", Name: "VirtualMachines", Tier: "Standard"}},
	})
}

func TestWriteCreatesApplicableFilesAndPreservesCSVQuoting(t *testing.T) {
	base := filepath.Join(t.TempDir(), "assessment")
	files, err := Write(csvFixture(), base, tables.Options{RedactSubscriptionIDs: true})
	if err != nil {
		t.Fatal(err)
	}

	present := map[string]bool{}
	for _, filename := range files {
		present[filepath.Base(filename)] = true
	}
	for _, suffix := range []string{"assessment.assessmentStatus.csv", "assessment.recommendations.csv", "assessment.impacted.csv", "assessment.resourceType.csv", "assessment.inventory.csv", "assessment.outofscope.csv", "assessment.defender.csv"} {
		if !present[suffix] {
			t.Fatalf("expected %s in generated files: %v", suffix, files)
		}
	}
	if present["assessment.advisor.csv"] {
		t.Fatal("skipped Advisor stage should not create a CSV")
	}

	file, err := os.Open(base + ".impacted.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows, err := encodingcsv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[1][5] != "A recommendation, with comma" {
		t.Fatalf("unexpected impacted CSV rows: %#v", rows)
	}
	if rows[1][7] == "11111111-1111-1111-1111-111111111111" {
		t.Fatal("subscription ID was not redacted")
	}
}

func TestWriteUsesPrivateFilePermissions(t *testing.T) {
	base := filepath.Join(t.TempDir(), "assessment")
	files, err := Write(csvFixture(), base, tables.Options{})
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(files[0])
	if err != nil {
		t.Fatal(err)
	}
	// Go's Windows FileMode reports the read-only attribute, not the file ACL.
	if runtime.GOOS != "windows" {
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("CSV mode = %o, want 600", got)
		}
	}
}

func TestWriteValidatesInputsAndPropagatesCreateErrors(t *testing.T) {
	if _, err := Write(nil, "x", tables.Options{}); err == nil {
		t.Fatal("expected nil result error")
	}
	if _, err := Write(csvFixture(), "", tables.Options{}); err == nil {
		t.Fatal("expected empty base filename error")
	}
	missingDir := filepath.Join(t.TempDir(), "missing", "assessment")
	if _, err := Write(csvFixture(), missingDir, tables.Options{}); err == nil {
		t.Fatal("expected file creation error")
	}
}
