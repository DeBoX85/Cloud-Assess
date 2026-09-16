package arcsql

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
)

type fakeGraph struct {
	result *arg.Result
	err    error
	query  string
}

func (f *fakeGraph) Query(_ context.Context, query string, _ map[string]string, _ ...arg.QueryOptions) (*arg.Result, error) {
	f.query = query
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

func TestScanMapsFiltersWarnsAndSorts(t *testing.T) {
	const (
		subA      = "11111111-1111-1111-1111-111111111111"
		subB      = "22222222-2222-2222-2222-222222222222"
		instanceA = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.AzureArcData/sqlServerInstances/a"
		instanceB = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.AzureArcData/sqlServerInstances/b"
		excluded  = "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/rg/providers/Microsoft.AzureArcData/sqlServerInstances/c"
	)
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"subscriptionId":"` + subA + `","status":"Connected","AzureArcServer":"/subscriptions/` + subA + `/resourceGroups/rg/providers/Microsoft.HybridCompute/machines/server-b","SQLInstance":"` + instanceB + `","resourceGroup":"rg","version":"2022","Build":"16.0","patchLevel":"CU1","edition":"Enterprise","vcores":"8","License":"PAYG","DPSStatus":"OK","TELStatus":"__","DefenderStatus":"Protected"}`),
		json.RawMessage(`{"subscriptionId":"` + subA + `","status":"Connected","AzureArcServer":"/subscriptions/` + subA + `/resourceGroups/rg/providers/Microsoft.HybridCompute/machines/server-a","SQLInstance":"` + instanceA + `","resourceGroup":"rg","version":"2019","Build":"15.0","patchLevel":"CU20","edition":"Standard","vcores":"4","License":"SA","DPSStatus":"No Data","TELStatus":"No Data","DefenderStatus":"NotProtected"}`),
		json.RawMessage(`{"subscriptionId":"` + subB + `","status":"Connected","SQLInstance":"` + excluded + `","resourceGroup":"rg","vcores":"2"}`),
		json.RawMessage(`{"subscriptionId":`),
	}}}
	filter := fakeFilter{excludedSubscriptions: map[string]bool{strings.ToLower(subB): true}}

	result, err := NewWithClient(graph).Scan(
		context.Background(),
		map[string]string{subA: "Subscription A", subB: "Subscription B"},
		filter,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("records = %d, want 2: %+v", len(result.Records), result.Records)
	}
	if result.Records[0].SQLInstance != instanceA || result.Records[1].SQLInstance != instanceB {
		t.Fatalf("unexpected deterministic order: %+v", result.Records)
	}
	first := result.Records[0]
	if first.SubscriptionName != "Subscription A" || first.Version != "2019" || first.VCores != "4" || first.License != "SA" {
		t.Fatalf("unexpected Arc SQL mapping: %+v", first)
	}
	if first.DPSStatus != "No Data" || first.TELStatus != "No Data" || first.DefenderStatus != "NotProtected" {
		t.Fatalf("unexpected Arc SQL status mapping: %+v", first)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "arcsql_malformed_arg_rows" {
		t.Fatalf("warnings = %+v", result.Warnings)
	}
	if !strings.Contains(graph.query, "Microsoft.AzureArcData/sqlServerInstances") ||
		!strings.Contains(graph.query, "WindowsAgent.SqlServer") ||
		!strings.Contains(graph.query, "join kind=inner") ||
		!strings.Contains(graph.query, "extend vcores = toint(properties.vCore)") {
		t.Fatalf("Arc SQL query lost source semantics: %s", graph.query)
	}
}

func TestResourceFilterUsesSQLInstance(t *testing.T) {
	const instance = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg/providers/Microsoft.AzureArcData/sqlServerInstances/sql1"
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"subscriptionId":"11111111-1111-1111-1111-111111111111","SQLInstance":"` + instance + `","vcores":"4"}`),
	}}}
	filter := fakeFilter{excludedResources: map[string]bool{strings.ToLower(instance): true}}
	result, err := NewWithClient(graph).Scan(context.Background(), nil, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 0 {
		t.Fatalf("excluded SQL instance produced records: %+v", result.Records)
	}
}

func TestPinnedVCoreTypeContractTreatsNumericVCoreAsMalformed(t *testing.T) {
	graph := &fakeGraph{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"subscriptionId":"sub1","SQLInstance":"/subscriptions/sub1/resourceGroups/rg/providers/Microsoft.AzureArcData/sqlServerInstances/sql1","vcores":4}`),
	}}}

	result, err := NewWithClient(graph).Scan(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 0 {
		t.Fatalf("numeric vcores should preserve pinned string-decoder behavior, got %+v", result.Records)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Code != "arcsql_malformed_arg_rows" {
		t.Fatalf("expected malformed-row warning for numeric vcores, got %+v", result.Warnings)
	}
}

func TestArcSQLQueryFailureIsReturned(t *testing.T) {
	graph := &fakeGraph{err: errors.New("arc sql failed")}
	_, err := NewWithClient(graph).Scan(context.Background(), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "arc sql failed") {
		t.Fatalf("expected query error, got %v", err)
	}
}

func TestArcSQLNilResultIsError(t *testing.T) {
	graph := &fakeGraph{result: nil}
	if _, err := NewWithClient(graph).Scan(context.Background(), nil, nil); err == nil {
		t.Fatal("expected nil Arc SQL ARG result to fail")
	}
}
