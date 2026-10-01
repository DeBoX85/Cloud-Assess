package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
)

func TestMissingExplicitScopePersistsFailedReportWithoutResourceQueries(t *testing.T) {
	const visible = "11111111-2222-3333-4444-555555555555"
	const missing = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	for _, masked := range []bool{false, true} {
		t.Run(map[bool]string{false: "private raw", true: "redacted"}[masked], func(t *testing.T) {
			operations := operationsForEndToEndTest(t, rules.NewCatalog(), assessment.Resource{}, assessment.Finding{})
			operations.DiscoverSubscriptions = func(context.Context, []string, *config.Filters) (map[string]string, error) {
				return map[string]string{visible: "Visible"}, nil
			}
			operations.DiscoverResources = func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error) {
				t.Fatal("resource query ran after missing requested scope")
				return nil, nil
			}
			base := filepath.Join(t.TempDir(), "scope-failure")
			outcome, err := NewRunner(orchestration.NewCoordinator(operations)).Run(context.Background(), ScanOptions{
				Assessment: orchestration.Request{Subscriptions: []string{visible, missing}}, Outputs: OutputOptions{BaseName: base, JSON: true, RedactSubscriptionIDs: masked},
			})
			if err == nil || !strings.Contains(err.Error(), "scope_requested_subscription_unresolved") || outcome.ExitCode != ExitExecutionFail {
				t.Fatalf("scope failure error=%v exit=%d", err, outcome.ExitCode)
			}
			content, err := os.ReadFile(base + ".json")
			if err != nil {
				t.Fatal(err)
			}
			var report result.AssessmentResult
			if err = json.Unmarshal(content, &report); err != nil {
				t.Fatal(err)
			}
			wantMissing := missing
			wantVisible := visible
			if masked {
				wantMissing = redact.SubscriptionID(missing, true)
				wantVisible = redact.SubscriptionID(visible, true)
			}
			if report.Completeness != assessment.CompletenessFailed || report.Scope.Status != "unresolved" || len(report.Scope.UnresolvedSubscriptionIDs) != 1 || report.Scope.UnresolvedSubscriptionIDs[0] != wantMissing || len(report.Scope.ResolvedSubscriptions) != 1 || report.Scope.ResolvedSubscriptions[0].SubscriptionID != wantVisible {
				t.Fatalf("persisted evidence: %#v", report.Scope)
			}
			if masked && (strings.Contains(string(content), missing) || strings.Contains(string(content), visible)) {
				t.Fatal("failed report leaked scope-only IDs")
			}
		})
	}
}
