package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type fakeGetter struct {
	urls []string
}

func (f *fakeGetter) Get(_ context.Context, url string) ([]byte, error) {
	f.urls = append(f.urls, url)
	if strings.Contains(url, "/storageAccounts/bad/") {
		return nil, &azcore.ResponseError{StatusCode: 400, ErrorCode: "BadRequest"}
	}
	if strings.Contains(url, "/storageAccounts/offline/") {
		return nil, errors.New("sensitive resource /subscriptions/private")
	}
	return []byte(`{"value":[]}`), nil
}

func TestProbeUsesOnlyEligibleDiagnosticSettingsGETsAndHidesResourceIDs(t *testing.T) {
	resources := []assessment.Resource{
		{ID: "/subscriptions/private/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/ok", Type: "Microsoft.Storage/storageAccounts"},
		{ID: "/subscriptions/private/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/bad", Type: "Microsoft.Storage/storageAccounts"},
		{ID: "/subscriptions/private/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/offline", Type: "Microsoft.Storage/storageAccounts"},
		{ID: "/subscriptions/private/resourceGroups/rg/providers/Microsoft.Compute/snapshots/ignored", Type: "Microsoft.Compute/snapshots"},
	}
	client := &fakeGetter{}
	result := probe(context.Background(), client, "https://management.azure.com/", eligibleResources(resources), 2)
	if result.EligibleResources != 3 || result.SuccessfulGETs != 1 || len(result.Failures) != 2 || result.OriginalHTTP400Warnings != 2 {
		t.Fatalf("unexpected probe summary: %+v", result)
	}
	for _, url := range client.urls {
		if !strings.HasPrefix(url, "https://management.azure.com/subscriptions/") || !strings.HasSuffix(url, "/providers/microsoft.insights/diagnosticSettings?api-version=2021-05-01-preview") {
			t.Fatalf("unexpected read-only request: %s", url)
		}
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "/subscriptions/") || strings.Contains(string(encoded), "resourceGroups") {
		t.Fatalf("resource IDs leaked in output: %s", encoded)
	}
	if result.Failures[0].HTTPStatus != 400 || result.Failures[0].AzureErrorCode != "BadRequest" {
		t.Fatalf("expected sanitized HTTP 400 result, got %+v", result.Failures[0])
	}
	if result.Failures[1].ErrorClass != "transport_or_authentication_error" {
		t.Fatalf("expected sanitized transport error, got %+v", result.Failures[1])
	}
}

func TestOriginalHTTP400Warnings(t *testing.T) {
	stages := []assessment.StageExecution{{Name: "diagnostics", Warnings: []assessment.AssessmentWarning{
		{Code: "diagnostics_subrequest_non_success", Message: "diagnostic settings batch subrequest returned HTTP 400"},
		{Code: "diagnostics_subrequest_non_success", Message: "diagnostic settings batch subrequest returned HTTP 403"},
		{Code: "diagnostics_subrequest_non_success", Message: "diagnostic settings batch subrequest returned HTTP 400"},
	}}}
	if got := originalHTTP400Warnings(stages); got != 2 {
		t.Fatalf("HTTP 400 warnings = %d, want 2", got)
	}
}
