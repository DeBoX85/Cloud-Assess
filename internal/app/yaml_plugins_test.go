package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

type yamlFixturePoster struct {
	t          *testing.T
	fail       bool
	calls      int
	resourceID string
}

func (p *yamlFixturePoster) PostStream(_ context.Context, endpoint string, body io.ReadSeekCloser) (*http.Response, error) {
	p.calls++
	if endpoint != "https://example.test/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01" {
		p.t.Fatal("wrong Graph endpoint")
	}
	var request arg.Request
	if err := json.NewDecoder(body).Decode(&request); err != nil {
		return nil, err
	}
	if request.Query != "resources | where type =~ 'Microsoft.Storage/storageAccounts' // under-development\n" || len(request.Subscriptions) != 1 || request.Subscriptions[0] != "11111111-2222-3333-4444-555555555555" || request.Options == nil || request.Options.ResultFormat != "objectArray" || request.Options.Top == nil || *request.Options.Top != 1000 {
		p.t.Fatalf("wrong plugin ARG request: %+v", request)
	}
	if p.fail {
		return nil, errors.New("synthetic plugin query denied")
	}
	data, _ := json.Marshal(map[string]any{"data": []map[string]any{{"id": p.resourceID, "subscriptionId": "11111111-2222-3333-4444-555555555555", "resourceGroup": "rg", "name": "store", "type": "Microsoft.Storage/storageAccounts", "location": "norwayeast"}}})
	return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(data))}, nil
}

func TestYAMLGraphToCanonicalReportAndFailureHealth(t *testing.T) {
	// Literal synthetic data traverses file parsing, catalog overlay, actual ARG
	// transport serialization/decoder/executor, coordinator and JSON rendering.
	const subscription = "11111111-2222-3333-4444-555555555555"
	resource := assessment.Resource{ID: "/subscriptions/" + subscription + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store", SubscriptionID: subscription, ResourceGroup: "rg", Name: "store", Type: "Microsoft.Storage/storageAccounts", Location: "norwayeast"}
	root := t.TempDir()
	const yaml = `name: operator-controls
queries:
  - aprlGuid: yaml-control
    description: Operator storage control
    recommendationControl: Security
    recommendationImpact: High
    recommendationResourceType: Microsoft.Storage/storageAccounts
    recommendationMetadataState: disabled
    query: "resources | where type =~ 'Microsoft.Storage/storageAccounts' // under-development\n"
`
	if err := os.WriteFile(filepath.Join(root, "control.yml"), []byte(yaml), 0600); err != nil {
		t.Fatal(err)
	}
	plugins, err := rules.DiscoverYAMLPlugins([]string{root})
	if err != nil {
		t.Fatal(err)
	}
	catalog := rules.NewCatalog()
	catalog.Add(assessment.RecommendationDefinition{ID: "yaml-control", ResourceType: resource.Type, Query: "wrong embedded query", Recommendation: "embedded", Source: rules.SourceAPRL})
	for _, tc := range []struct {
		name                string
		fail, excluded      bool
		wantCalls, wantRows int
		wantCompleteness    assessment.Completeness
		wantExit            int
	}{
		{"healthy external disabled marker still executes", false, false, 1, 1, assessment.CompletenessComplete, 0},
		{"recommendation exclusion", false, true, 0, 0, assessment.CompletenessComplete, 0},
		{"query failure", true, false, 1, 0, assessment.CompletenessFailed, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			poster := &yamlFixturePoster{t: t, fail: tc.fail, resourceID: resource.ID}
			client := arg.NewClient(arg.NewHTTPTransportWithClient(poster, "https://example.test/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01"))
			ops := operationsForEndToEndTest(t, catalog, resource, assessment.Finding{})
			ops.ExecuteGraph = func(ctx context.Context, defs []assessment.RecommendationDefinition, subs map[string]string, filter *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
				return arg.ExecuteRecommendations(ctx, client, defs, subs, filter.IsServiceExcluded, 1)
			}
			configStages := stages.NewDefault()
			if err := configStages.Apply([]string{"-diagnostics", "-advisor", "-defender"}); err != nil {
				t.Fatal(err)
			}
			filters := config.NewFilters()
			if tc.excluded {
				filters.Assessment.Exclude.Recommendations = []string{"yaml-control"}
				filters.RebuildIndexes()
			}
			base := filepath.Join(t.TempDir(), "yaml-result")
			outcome, runErr := NewRunner(orchestration.NewCoordinator(ops)).Run(context.Background(), ScanOptions{Assessment: orchestration.Request{ScannerKeys: []string{"st"}, Filters: filters, Stages: configStages, YAMLRecommendations: rules.PluginDefinitions(plugins)}, Outputs: OutputOptions{BaseName: base, JSON: true}})
			if (runErr != nil) != tc.fail || outcome.ExitCode != tc.wantExit {
				t.Fatalf("outcome: %+v %v", outcome, runErr)
			}
			raw, err := os.ReadFile(base + ".json")
			if err != nil {
				t.Fatal(err)
			}
			var report result.AssessmentResult
			if err := json.Unmarshal(raw, &report); err != nil {
				t.Fatal(err)
			}
			if poster.calls != tc.wantCalls || len(report.Findings) != tc.wantRows || report.Completeness != tc.wantCompleteness {
				t.Fatalf("request/rows/health: calls=%d %+v", poster.calls, report)
			}
			if tc.wantRows == 1 {
				finding := report.Findings[0]
				if finding.RecommendationID != "yaml-control" || finding.ResourceID != resource.ID || finding.Source != "operator-controls" || finding.Recommendation != "Operator storage control" || finding.Category != "Security" || finding.Impact != "High" || finding.SubscriptionName != "Test Subscription" || finding.LearnMoreURL != "" {
					t.Fatalf("plugin projection: %+v", finding)
				}
			}
			if tc.fail {
				healthy := false
				for _, stage := range report.Stages {
					if stage.Name == stages.Graph && stage.Status == assessment.StageFailed && stage.Error != nil {
						healthy = true
					}
				}
				if !healthy || !strings.Contains(string(raw), "synthetic plugin query denied") {
					t.Fatal("failed Graph health not persisted")
				}
			}
			if got := catalog.ByResourceType(resource.Type); len(got) != 1 || got[0].Source != rules.SourceAPRL || got[0].Query != "wrong embedded query" {
				t.Fatal("shared base catalog leaked plugin override")
			}
		})
	}
}
