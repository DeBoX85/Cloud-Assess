package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

type fakeResourceQuerier struct {
	query  string
	result *arg.Result
	err    error
}

func (f *fakeResourceQuerier) Query(_ context.Context, query string, _ map[string]string, _ ...arg.QueryOptions) (*arg.Result, error) {
	f.query = query
	return f.result, f.err
}

func TestDiscoverResourcesUsesReferenceProjectionAndPartitionsScope(t *testing.T) {
	includedID := "/subscriptions/sub/resourceGroups/keep/providers/Microsoft.Test/widgets/one"
	excludedID := "/subscriptions/sub/resourceGroups/drop/providers/Microsoft.Test/widgets/two"
	q := &fakeResourceQuerier{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"id":"` + includedID + `","subscriptionId":"sub","resourceGroup":"keep","location":"westeurope","type":"Microsoft.Test/widgets","name":"one","tags":{"env":"prod"},"skuName":"S1","skuTier":"Standard","skuFamily":"S","skuCapacity":2,"kind":"test"}`),
		json.RawMessage(`{"id":"` + excludedID + `","subscriptionId":"sub","resourceGroup":"drop","location":"westeurope","type":"Microsoft.Test/widgets","name":"two","tags":{"env":"prod"}}`),
	}}}
	filters := config.NewFilters()
	filters.Assessment.Exclude.ResourceGroups = []string{"/subscriptions/sub/resourceGroups/drop"}
	filters.RebuildIndexes()

	inventory, err := DiscoverResources(context.Background(), q, map[string]string{"sub": "Test"}, filters)
	if err != nil {
		t.Fatalf("DiscoverResources() error = %v", err)
	}
	if q.query != ResourceInventoryQuery {
		t.Fatalf("inventory query changed:\n%s", q.query)
	}
	if len(inventory.Included) != 1 || inventory.Included[0].ID != includedID {
		t.Fatalf("unexpected included resources: %#v", inventory.Included)
	}
	if len(inventory.Excluded) != 1 || inventory.Excluded[0].ID != excludedID {
		t.Fatalf("unexpected excluded resources: %#v", inventory.Excluded)
	}
	resource := inventory.Included[0]
	if resource.SKUName != "S1" || resource.SKUTier != "Standard" || resource.SKUFamily != "S" || resource.SKUCapacity != 2 || resource.Kind != "test" {
		t.Fatalf("SKU/kind projection lost data: %+v", resource)
	}
	if filters.Assessment.IsServiceExcluded(includedID + "/slots/blue") {
		t.Fatal("included parent scope should be inherited by child resources")
	}
}

func TestDiscoverResourcesAppliesTagAndResourceTypeFilters(t *testing.T) {
	q := &fakeResourceQuerier{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"id":"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one","subscriptionId":"sub","resourceGroup":"rg","type":"Microsoft.Test/widgets","name":"one","tags":{"ENV":"prod"}}`),
		json.RawMessage(`{"id":"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Other/things/two","subscriptionId":"sub","resourceGroup":"rg","type":"Microsoft.Other/things","name":"two","tags":{"env":"prod"}}`),
		json.RawMessage(`{"id":"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/three","subscriptionId":"sub","resourceGroup":"rg","type":"Microsoft.Test/widgets","name":"three","tags":{"env":"dev"}}`),
	}}}
	filters := config.NewFilters()
	filters.Assessment.Include.ResourceTypes = []string{"Microsoft.Test/widgets"}
	filters.Assessment.Include.Tags = map[string]string{"env": "prod"}
	filters.RebuildIndexes()

	inventory, err := DiscoverResources(context.Background(), q, map[string]string{"sub": "Test"}, filters)
	if err != nil {
		t.Fatalf("DiscoverResources() error = %v", err)
	}
	if len(inventory.Included) != 1 || inventory.Included[0].Name != "one" || len(inventory.Excluded) != 2 {
		t.Fatalf("unexpected filter partition: included=%#v excluded=%#v", inventory.Included, inventory.Excluded)
	}
}

func TestDiscoverResourcesSkipsMalformedRowsAndRecordsWarningSignal(t *testing.T) {
	q := &fakeResourceQuerier{result: &arg.Result{Data: []json.RawMessage{
		json.RawMessage(`{"id":"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one","subscriptionId":"sub","resourceGroup":"rg","type":"Microsoft.Test/widgets","name":"one"}`),
		json.RawMessage(`{"skuCapacity":"not-an-int"}`),
		json.RawMessage(`{not-json}`),
	}}}
	inventory, err := DiscoverResources(context.Background(), q, nil, nil)
	if err != nil {
		t.Fatalf("DiscoverResources() error = %v", err)
	}
	if len(inventory.Included) != 1 || inventory.MalformedRows != 2 {
		t.Fatalf("unexpected malformed-row handling: %+v", inventory)
	}
}

func TestDiscoverResourcesReturnsQueryFailure(t *testing.T) {
	boom := errors.New("ARG unavailable")
	q := &fakeResourceQuerier{err: boom}
	inventory, err := DiscoverResources(context.Background(), q, nil, nil)
	if err == nil || !errors.Is(err, boom) || inventory != nil {
		t.Fatalf("unexpected failure result: inventory=%#v err=%v", inventory, err)
	}
}

func TestCountResourcesByTypeAndSubscriptionIsDeterministic(t *testing.T) {
	resources := []struct {
		sub      string
		typeName string
	}{
		{sub: "sub-b", typeName: "Microsoft.Storage/storageAccounts"},
		{sub: "sub-a", typeName: "Microsoft.Compute/virtualMachines"},
		{sub: "sub-a", typeName: "Microsoft.Compute/virtualMachines"},
	}
	converted := make([]assessment.Resource, 0, len(resources))
	for _, resource := range resources {
		converted = append(converted, assessment.Resource{SubscriptionID: resource.sub, Type: resource.typeName})
	}
	counts := CountResourcesByTypeAndSubscription(converted, map[string]string{"sub-a": "Alpha", "sub-b": "Beta"})
	if len(counts) != 2 {
		t.Fatalf("count rows = %d, want 2", len(counts))
	}
	if counts[0].SubscriptionName != "Alpha" || counts[0].Count != 2 || counts[0].ResourceType != "Microsoft.Compute/virtualMachines" {
		t.Fatalf("unexpected first count: %+v", counts[0])
	}
	if counts[1].SubscriptionName != "Beta" || counts[1].Count != 1 {
		t.Fatalf("unexpected second count: %+v", counts[1])
	}
}
