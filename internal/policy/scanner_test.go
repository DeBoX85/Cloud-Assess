package policy

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
)

type fakeGraph struct {
	result        *arg.Result
	err           error
	query         string
	subscriptions map[string]string
	opts          []arg.QueryOptions
}

func (f *fakeGraph) Query(_ context.Context, query string, subscriptions map[string]string, opts ...arg.QueryOptions) (*arg.Result, error) {
	f.query = query
	f.subscriptions = subscriptions
	f.opts = append([]arg.QueryOptions(nil), opts...)
	if f.err != nil {
		return nil, f.err
	}
	return f.result, nil
}

type fakeFilter struct {
	excludedSubscriptions map[string]bool
	excludedResources     map[string]bool
}

func (f fakeFilter) IsSubscriptionExcluded(subscriptionID string) bool {
	return f.excludedSubscriptions[strings.ToLower(subscriptionID)]
}

func (f fakeFilter) IsServiceExcluded(resourceID string) bool {
	return f.excludedResources[strings.ToLower(resourceID)]
}

func TestScanMapsFiltersDeduplicatesWarnsAndUsesManagementGroupScope(t *testing.T) {
	const (
		subA     = "11111111-1111-1111-1111-111111111111"
		subB     = "22222222-2222-2222-2222-222222222222"
		resource = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/acct1"
		excluded = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/skipme"
		otherSub = "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/rg2/providers/Microsoft.Storage/storageAccounts/acct2"
	)
	row := json.RawMessage(`{"subscriptionId":"` + subA + `","subscriptionName":"Subscription A","resourceId":"` + resource + `","policyDefinitionDisplayName":"Audit storage","policyDescription":"desc","timestamp":"2026-09-16T10:00:00Z","policyDefinitionName":"def","policyDefinitionId":"def-1","policyAssignmentName":"assign","policyAssignmentId":"assign-1","complianceState":"NonCompliant"}`)
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		row,
		row,
		json.RawMessage(`{"subscriptionId":"` + subA + `","subscriptionName":"Subscription A","resourceId":"` + excluded + `","policyDefinitionDisplayName":"Audit vault","policyDescription":"desc","timestamp":"2026-09-16T10:00:00Z","policyDefinitionName":"def-vault","policyDefinitionId":"def-2","policyAssignmentName":"assign","policyAssignmentId":"assign-1","complianceState":"NonCompliant"}`),
		json.RawMessage(`{"subscriptionId":"` + subB + `","subscriptionName":"Subscription B","resourceId":"` + otherSub + `","policyDefinitionDisplayName":"Audit storage","policyDescription":"desc","timestamp":"2026-09-16T10:00:00Z","policyDefinitionName":"def-other","policyDefinitionId":"def-3","policyAssignmentName":"assign","policyAssignmentId":"assign-1","complianceState":"NonCompliant"}`),
		json.RawMessage(`{"subscriptionId":`),
	}}}
	filter := fakeFilter{
		excludedSubscriptions: map[string]bool{strings.ToLower(subB): true},
		excludedResources:     map[string]bool{strings.ToLower(excluded): true},
	}

	result, err := NewWithClient(graph).Scan(
		context.Background(),
		map[string]string{subA: "Subscription A", subB: "Subscription B"},
		filter,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("records = %d, want 1: %+v", len(result.Records), result.Records)
	}
	record := result.Records[0]
	if record.SubscriptionID != subA || record.SubscriptionName != "Subscription A" {
		t.Fatalf("unexpected subscription mapping: %+v", record)
	}
	if record.ResourceGroup != "rg1" || record.ResourceType != "Microsoft.Storage/storageAccounts" || record.ResourceName != "acct1" {
		t.Fatalf("unexpected resource ID mapping: %+v", record)
	}
	if record.PolicyDisplayName != "Audit storage" || record.PolicyDefinitionID != "def-1" || record.PolicyAssignmentID != "assign-1" {
		t.Fatalf("unexpected policy mapping: %+v", record)
	}
	if record.ComplianceState != "NonCompliant" || record.Timestamp != "2026-09-16T10:00:00Z" {
		t.Fatalf("unexpected compliance mapping: %+v", record)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "policy_malformed_arg_rows" {
		t.Fatalf("warnings = %+v", result.Warnings)
	}
	if len(graph.opts) != 1 || !graph.opts[0].ManagementGroupScope {
		t.Fatalf("Policy query must use management-group-aware ARG scope: %+v", graph.opts)
	}
	if !strings.Contains(graph.query, "complianceState == 'NonCompliant'") ||
		!strings.Contains(graph.query, "microsoft.authorization/policydefinitions") ||
		!strings.Contains(graph.query, "microsoft.resources/subscriptions") {
		t.Fatalf("Policy query lost source semantics: %s", graph.query)
	}
}

func TestDedupKeyIncludesPolicyDefinitionID(t *testing.T) {
	const resourceID = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/st1"
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"subscriptionId":"11111111-1111-1111-1111-111111111111","subscriptionName":"Sub","resourceId":"` + resourceID + `","policyDefinitionId":"def-1","complianceState":"NonCompliant"}`),
		json.RawMessage(`{"subscriptionId":"11111111-1111-1111-1111-111111111111","subscriptionName":"Sub","resourceId":"` + resourceID + `","policyDefinitionId":"def-2","complianceState":"NonCompliant"}`),
	}}}

	result, err := NewWithClient(graph).Scan(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records = %d, want 2 distinct policy definitions", len(result.Records))
	}
}

func TestPolicyOutputIsDeterministic(t *testing.T) {
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"subscriptionId":"b","resourceId":"/subscriptions/b/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/z","policyDefinitionId":"def-2","complianceState":"NonCompliant"}`),
		json.RawMessage(`{"subscriptionId":"a","resourceId":"/subscriptions/a/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/y","policyDefinitionId":"def-3","complianceState":"NonCompliant"}`),
		json.RawMessage(`{"subscriptionId":"a","resourceId":"/subscriptions/a/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/x","policyDefinitionId":"def-1","complianceState":"NonCompliant"}`),
	}}}

	result, err := NewWithClient(graph).Scan(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 3 {
		t.Fatalf("records = %d, want 3", len(result.Records))
	}
	if result.Records[0].SubscriptionID != "a" || !strings.HasSuffix(result.Records[0].ResourceID, "/x") || result.Records[2].SubscriptionID != "b" {
		t.Fatalf("unexpected deterministic order: %+v", result.Records)
	}
}

func TestPolicyQueryFailureIsReturned(t *testing.T) {
	graph := &fakeGraph{err: errors.New("policy failed")}
	_, err := NewWithClient(graph).Scan(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "policy failed") {
		t.Fatalf("expected query error, got %v", err)
	}
}

func TestPolicyNilResultIsError(t *testing.T) {
	graph := &fakeGraph{result: nil}
	if _, err := NewWithClient(graph).Scan(context.Background(), nil, nil); err == nil {
		t.Fatal("expected nil Azure Policy ARG result to fail")
	}
}
