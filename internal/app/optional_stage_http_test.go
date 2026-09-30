package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
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

// These synthetic HTTP envelopes test retrieval health and source-shaped row
// projections, not live Azure query execution or default retry sequencing. Only the optional adapter is real;
// healthy scope/inventory/Graph operations use the existing coordinator fixture.
func TestOptionalStageHTTPHealthPersistsReport(t *testing.T) {
	fixtureBytes, err := os.ReadFile(filepath.Join("testdata", "optional-stage-projections.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Row          json.RawMessage  `json:"row"`
		MalformedRow json.RawMessage  `json:"malformedRow"`
		Expected     []map[string]any `json:"expected"`
	}
	if err := json.Unmarshal(fixtureBytes, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, stageName := range []string{stages.Policy, stages.Defender, stages.DefenderRecommendations} {
		for _, tc := range []struct {
			name      string
			status    int
			body      string
			failed    bool
			reason    string
			projected bool
			malformed bool
		}{
			{name: "denied", status: http.StatusForbidden, body: `{"error":{"code":"AuthorizationFailed","message":"synthetic denial"}}`, failed: true, reason: "AuthorizationFailed"},
			{name: "throttled", status: http.StatusTooManyRequests, body: `{"error":{"code":"TooManyRequests","message":"synthetic throttle"}}`, failed: true, reason: "TooManyRequests"},
			{name: "missing-data", status: http.StatusOK, body: `{"count":0}`, failed: true, reason: "missing or null data array"},
			{name: "valid-empty", status: http.StatusOK, body: `{"data":[]}`},
			{name: "non-empty", status: http.StatusOK, projected: true},
			{name: "mixed-malformed", status: http.StatusOK, projected: true, malformed: true},
			{name: "all-malformed", status: http.StatusOK, malformed: true},
		} {
			t.Run(stageName+"/"+tc.name, func(t *testing.T) {
				const sub = "11111111-2222-3333-4444-555555555555"
				payload := tc.body
				wantRecords := 0
				if tc.projected || tc.malformed {
					fixture, ok := fixtures[stageName]
					if !ok || len(fixture.Expected) != 1 {
						t.Fatal("missing projection fixture")
					}
					data := []json.RawMessage{}
					if tc.projected {
						data = append(data, fixture.Row)
						wantRecords = 1
					}
					if tc.malformed {
						data = append(data, fixture.MalformedRow)
					}
					envelope, err := json.Marshal(struct {
						Data []json.RawMessage `json:"data"`
					}{data})
					if err != nil {
						t.Fatal(err)
					}
					payload = string(envelope)
				}
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
					_, _ = io.WriteString(w, payload)
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
				if tc.malformed {
					wantCompleteness, wantStatus = assessment.CompletenessCompleteWithWarnings, assessment.StageCompletedWithWarnings
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
				if len(report.AzurePolicy)+len(report.Defender)+len(report.DefenderRecommendations) != wantRecords {
					t.Fatal("retrieval fabricated optional records")
				}
				if tc.projected {
					var document map[string]json.RawMessage
					if err := json.Unmarshal(encoded, &document); err != nil {
						t.Fatal(err)
					}
					dataset := map[string]string{stages.Policy: "azurePolicy", stages.Defender: "defender", stages.DefenderRecommendations: "defenderRecommendations"}[stageName]
					var actual []map[string]any
					if err := json.Unmarshal(document[dataset], &actual); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(actual, fixtures[stageName].Expected) {
						t.Fatalf("projection mismatch: got %+v want %+v", actual, fixtures[stageName].Expected)
					}
				}
				found := false
				for _, stage := range report.Stages {
					if stage.Name == stageName {
						found = true
						if stage.Status != wantStatus || stage.Records != wantRecords || (stage.Error != nil) != tc.failed {
							t.Fatalf("incorrect persisted stage health: %+v", stage)
						}
						wantWarnings := 0
						if tc.malformed {
							wantWarnings = 1
						}
						if len(stage.Warnings) != wantWarnings {
							t.Fatalf("warnings=%+v want count=%d", stage.Warnings, wantWarnings)
						}
						if tc.malformed {
							code := map[string]string{stages.Policy: "policy_malformed_arg_rows", stages.Defender: "defender_status_malformed_arg_rows", stages.DefenderRecommendations: "defender_recommendations_malformed_arg_rows"}[stageName]
							if stage.Warnings[0].Code != code || !strings.Contains(stage.Warnings[0].Message, "skipped 1 malformed") {
								t.Fatalf("incorrect malformed-row warning: %+v", stage.Warnings)
							}
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
