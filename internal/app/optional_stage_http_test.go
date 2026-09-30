package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/defender"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

type optionalStageCredential struct{}

func (optionalStageCredential) GetToken(context.Context, azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "synthetic-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// These synthetic HTTP envelopes test retrieval health, not non-empty Azure row
// projection or default retry sequencing. Only the optional adapter is real;
// healthy scope/inventory/Graph operations use the existing coordinator fixture.
func TestOptionalStageHTTPHealthPersistsReport(t *testing.T) {
	for _, stageName := range []string{stages.Policy, stages.Defender, stages.DefenderRecommendations} {
		for _, tc := range []struct {
			name   string
			status int
			body   string
			failed bool
			reason string
		}{
			{"denied", http.StatusForbidden, `{"error":{"code":"AuthorizationFailed","message":"synthetic denial"}}`, true, "AuthorizationFailed"},
			{"throttled", http.StatusTooManyRequests, `{"error":{"code":"TooManyRequests","message":"synthetic throttle"}}`, true, "TooManyRequests"},
			{"missing-data", http.StatusOK, `{"count":0}`, true, "missing or null data array"},
			{"valid-empty", http.StatusOK, `{"data":[]}`, false, ""},
		} {
			t.Run(stageName+"/"+tc.name, func(t *testing.T) {
				const sub = "11111111-2222-3333-4444-555555555555"
				expectedQuery := policy.Query
				if stageName == stages.Defender {
					expectedQuery = defender.StatusQuery
				}
				if stageName == stages.DefenderRecommendations {
					expectedQuery = defender.RecommendationsQuery
				}
				var calls atomic.Int32
				server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					if r.Method != http.MethodPost || r.URL.Path != "/providers/Microsoft.ResourceGraph/resources" || r.URL.Query().Get("api-version") != arg.ResourceGraphAPIVersion || r.Header.Get("Authorization") != "Bearer synthetic-token" {
						t.Errorf("unexpected request contract: %s %s", r.Method, r.URL)
					}
					var request arg.Request
					if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
						t.Errorf("decode request: %v", err)
					}
					if request.Query != expectedQuery || len(request.Subscriptions) != 1 || request.Subscriptions[0] != sub {
						t.Errorf("unexpected query or scope: %+v", request)
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tc.status)
					_, _ = io.WriteString(w, tc.body)
				}))
				defer server.Close()
				client := azure.NewHTTPClient(optionalStageCredential{}, &azure.HTTPClientOptions{Timeout: 2 * time.Second, MaxRetries: -1, Scope: "https://management.azure.com/.default", Transport: server.Client()})
				graph := arg.NewClient(arg.NewHTTPTransportWithClient(client, server.URL+"/providers/Microsoft.ResourceGraph/resources?api-version="+arg.ResourceGraphAPIVersion))
				resource := assessment.Resource{ID: "/subscriptions/" + sub + "/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/store", SubscriptionID: sub, ResourceGroup: "rg", Type: "Microsoft.Storage/storageAccounts"}
				definition := assessment.RecommendationDefinition{ID: "healthy-fixture", ResourceType: resource.Type, Query: "resources", Source: rules.SourceCustom}
				finding := assessment.Finding{RecommendationID: definition.ID, ResourceID: resource.ID, ResourceType: resource.Type}
				catalog := rules.NewCatalog()
				catalog.Add(definition)
				operations := operationsForEndToEndTest(t, catalog, resource, finding)
				switch stageName {
				case stages.Policy:
					operations.ScanPolicy = func(ctx context.Context, subs map[string]string, filter *config.AssessmentFilter) (policy.Result, error) {
						return policy.NewWithClient(graph).Scan(ctx, subs, filter)
					}
				case stages.Defender:
					operations.ScanDefenderStatus = func(ctx context.Context, subs map[string]string, filter *config.AssessmentFilter) (defender.StatusResult, error) {
						return defender.NewWithClient(graph).ScanStatus(ctx, subs, filter)
					}
				case stages.DefenderRecommendations:
					operations.ScanDefenderRecommendations = func(ctx context.Context, subs map[string]string, filter *config.AssessmentFilter) (defender.RecommendationsResult, error) {
						return defender.NewWithClient(graph).ScanRecommendations(ctx, subs, filter)
					}
				}
				stageConfig := stages.NewDefault()
				for _, name := range []string{stages.Diagnostics, stages.Advisor, stages.Defender} {
					if err := stageConfig.Set(name, false); err != nil {
						t.Fatal(err)
					}
				}
				if err := stageConfig.Set(stageName, true); err != nil {
					t.Fatal(err)
				}
				base := filepath.Join(t.TempDir(), "health")
				outcome, runErr := NewRunner(orchestration.NewCoordinator(operations)).Run(context.Background(), ScanOptions{Assessment: orchestration.Request{Subscriptions: []string{sub}, ScannerKeys: []string{"st"}, Stages: stageConfig}, Outputs: OutputOptions{BaseName: base, JSON: true}})
				wantExit, wantCompleteness, wantStatus := ExitSuccess, assessment.CompletenessComplete, assessment.StageCompleted
				if tc.failed {
					wantExit, wantCompleteness, wantStatus = ExitPartial, assessment.CompletenessPartial, assessment.StageFailed
				}
				if (runErr != nil) != tc.failed || outcome.ExitCode != wantExit {
					t.Fatalf("error=%v exit=%d want failed=%v exit=%d", runErr, outcome.ExitCode, tc.failed, wantExit)
				}
				if calls.Load() != 1 {
					t.Fatalf("HTTP calls=%d want 1 with retries disabled", calls.Load())
				}
				encoded, err := os.ReadFile(base + ".json")
				if err != nil {
					t.Fatalf("report missing: %v", err)
				}
				var report result.AssessmentResult
				if err := json.Unmarshal(encoded, &report); err != nil {
					t.Fatal(err)
				}
				if report.Completeness != wantCompleteness || len(report.Resources) != 1 || len(report.Findings) != 1 || report.Findings[0].RecommendationID != definition.ID {
					t.Fatalf("healthy data or completeness lost: %+v", report)
				}
				if len(report.AzurePolicy)+len(report.Defender)+len(report.DefenderRecommendations) != 0 {
					t.Fatal("retrieval fabricated optional records")
				}
				found := false
				for _, stage := range report.Stages {
					if stage.Name == stageName {
						found = true
						if stage.Status != wantStatus || stage.Records != 0 || (stage.Error != nil) != tc.failed {
							t.Fatalf("incorrect persisted stage health: %+v", stage)
						}
						if tc.failed && (stage.Error.Code != "stage_failed" || !strings.Contains(stage.Error.Message, tc.reason)) {
							t.Fatalf("unexpected failure code: %+v", stage.Error)
						}
					}
				}
				if !found {
					t.Fatal("report omitted optional stage")
				}
			})
		}
	}
}
