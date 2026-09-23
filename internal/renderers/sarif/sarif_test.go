package sarif

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

func TestMarshalCanonicalSARIF(t *testing.T) {
	data := testResult()
	encoded, err := Marshal(data, "1.2.3")
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var report logFile
	if err := json.Unmarshal(encoded, &report); err != nil {
		t.Fatalf("unmarshal SARIF: %v", err)
	}
	if report.Version != "2.1.0" || report.Schema != schemaURL {
		t.Fatalf("unexpected SARIF metadata: version=%q schema=%q", report.Version, report.Schema)
	}
	if len(report.Runs) != 1 {
		t.Fatalf("runs = %d, want 1", len(report.Runs))
	}

	run := report.Runs[0]
	brand := branding.Default()
	if run.Tool.Driver.Name != brand.CLIName {
		t.Fatalf("tool name = %q, want %q", run.Tool.Driver.Name, brand.CLIName)
	}
	if run.Tool.Driver.InformationURI != brand.WebsiteURL {
		t.Fatalf("information URI = %q, want %q", run.Tool.Driver.InformationURI, brand.WebsiteURL)
	}
	if run.Tool.Driver.Version != "1.2.3" {
		t.Fatalf("tool version = %q, want 1.2.3", run.Tool.Driver.Version)
	}
	if run.Automation.ID != "cloud-assess/scope-123" {
		t.Fatalf("automation ID = %q", run.Automation.ID)
	}
	if len(run.Tool.Driver.Rules) != 3 {
		t.Fatalf("rules = %d, want 3 impacted non-SLA rules", len(run.Tool.Driver.Rules))
	}
	if len(run.Results) != 3 {
		t.Fatalf("results = %d, want 3 deduplicated non-SLA results", len(run.Results))
	}

	levels := map[string]string{}
	for _, item := range run.Results {
		levels[item.RuleID] = item.Level
		if len(item.PartialFingerprints) != 1 {
			t.Fatalf("fingerprints for %q = %#v", item.RuleID, item.PartialFingerprints)
		}
		if _, ok := item.PartialFingerprints[fingerprintNamespace]; !ok {
			t.Fatalf("missing %q fingerprint for %q", fingerprintNamespace, item.RuleID)
		}
		if _, legacy := item.PartialFingerprints["azqrFinding/v1"]; legacy {
			t.Fatalf("legacy fingerprint namespace leaked for %q", item.RuleID)
		}
	}
	if levels["rec-high"] != "error" || levels["rec-medium"] != "warning" || levels["rec-low"] != "note" {
		t.Fatalf("unexpected SARIF levels: %#v", levels)
	}
	if strings.Contains(string(encoded), "rec-sla") {
		t.Fatal("SLA recommendation leaked into SARIF")
	}
	if strings.Contains(strings.ToLower(string(encoded)), "azqr") {
		t.Fatal("legacy product identity leaked into SARIF")
	}
}

func TestMarshalDeterministic(t *testing.T) {
	data := testResult()
	first, err := Marshal(data, "dev")
	if err != nil {
		t.Fatal(err)
	}
	second, err := Marshal(data, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("SARIF output is not deterministic")
	}
}

func TestMarshalRejectsInvalidInput(t *testing.T) {
	if _, err := Marshal(nil, "dev"); err == nil {
		t.Fatal("Marshal(nil) error = nil, want error")
	}
	if _, err := Marshal(&result.AssessmentResult{}, "dev"); err == nil {
		t.Fatal("Marshal(result without summary) error = nil, want error")
	}
}

func TestWriteFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "assessment.sarif")
	if err := WriteFile(testResult(), filename, "dev"); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	// Go's Windows FileMode reports the read-only attribute, not the file ACL.
	if runtime.GOOS != "windows" {
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("permissions = %o, want 600", got)
		}
	}
	if err := WriteFile(testResult(), "", "dev"); err == nil {
		t.Fatal("WriteFile() with empty filename error = nil, want error")
	}
}

func testResult() *result.AssessmentResult {
	definitions := []assessment.RecommendationDefinition{
		{ID: "rec-high", Recommendation: "High recommendation", Category: "Security", Impact: "High", ResourceType: "Microsoft.Compute/virtualMachines", Source: "APRL"},
		{ID: "rec-medium", Recommendation: "Medium recommendation", Category: "Reliability", Impact: "Medium", ResourceType: "Microsoft.Storage/storageAccounts", Source: "CUSTOM"},
		{ID: "rec-low", Recommendation: "Low recommendation", Category: "Operational Excellence", Impact: "Low", ResourceType: "Microsoft.Network/virtualNetworks", Source: "AOR"},
		{ID: "rec-sla", Recommendation: "SLA information", Category: assessment.CategorySLA, Impact: "Low", ResourceType: "Microsoft.Compute/virtualMachines", Source: "APRL"},
	}
	findings := []assessment.Finding{
		{RecommendationID: "rec-high", Recommendation: "High recommendation", Category: "Security", Impact: "High", ResourceType: "Microsoft.Compute/virtualMachines", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm1", ResourceName: "vm1"},
		{RecommendationID: "REC-HIGH", Recommendation: "High recommendation", Category: "Security", Impact: "High", ResourceType: "Microsoft.Compute/virtualMachines", ResourceID: "/SUBSCRIPTIONS/SUB/RESOURCEGROUPS/RG/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINES/VM1", ResourceName: "vm1"},
		{RecommendationID: "rec-medium", Recommendation: "Medium recommendation", Category: "Reliability", Impact: "Medium", ResourceType: "Microsoft.Storage/storageAccounts", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1", ResourceName: "st1"},
		{RecommendationID: "rec-low", Recommendation: "Low recommendation", Category: "Operational Excellence", Impact: "Low", ResourceType: "Microsoft.Network/virtualNetworks", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet1", ResourceName: "vnet1"},
		{RecommendationID: "rec-sla", Recommendation: "SLA information", Category: assessment.CategorySLA, Impact: "Low", ResourceType: "Microsoft.Compute/virtualMachines", ResourceID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm1", ResourceName: "vm1"},
	}
	return result.Build(result.Input{
		GeneratedAt:     time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC),
		ScopeID:         "scope-123",
		Completeness:    assessment.CompletenessComplete,
		Recommendations: definitions,
		Findings:        findings,
	})
}
