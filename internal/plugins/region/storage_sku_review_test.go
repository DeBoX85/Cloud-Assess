package region

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestStorageSKULateUnknownCancellation(t *testing.T) {
	// Cancellation becomes visible after initial admission but before the nil
	// collection's unknown return. Terminal context must still win.
	ctx := &storageCountingContext{Context: context.Background(), limit: 2}
	var c *StorageSKUCollection
	got, err := c.Evidence(ctx, nil)
	if !errors.Is(err, context.Canceled) || len(got.Values) != 0 {
		t.Fatalf("unknown projection lost late cancellation: %#v / %v", got, err)
	}
}

func TestStorageSKURequiredKindIsolation(t *testing.T) {
	c := storageCollectBody(t, `{"value":[{"name":"s","kind":"Storage","resourceType":"storageAccounts","locations":["eastus"],"restrictions":[]},{"name":"s","kind":"BlobStorage","resourceType":"storageAccounts","restrictions":[]}]}`)
	got, err := c.Evidence(context.Background(), map[string][]string{"s": {"Storage"}})
	want := SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{"s": {State: 0, BlockedZones: []string{}}}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("nonrequired incomplete kind poisoned known evidence: %#v / %v", got, err)
	}
	got, err = c.Evidence(context.Background(), map[string][]string{"s": {"Storage", "BlobStorage"}})
	if err != nil || got.Status != "unknown" || len(got.Values) != 0 {
		t.Fatalf("required incomplete kind became complete: %#v / %v", got, err)
	}
}

func TestStorageSKUHardRestrictionDominates(t *testing.T) {
	for _, rows := range []string{
		`{"name":"s","kind":"Storage","resourceType":"storageAccounts","locations":["eastus"],"restrictions":[{"type":"location","values":["eastus"],"reasonCode":"QuotaId"},{"type":"location","values":["eastus"],"reasonCode":"NotAvailableForSubscription"}]}`,
		`{"name":"s","kind":"Storage","resourceType":"storageAccounts","locations":["eastus"],"restrictions":[{"type":"location","values":["eastus"],"reasonCode":"NotAvailableForSubscription"},{"type":"location","values":["eastus"],"reasonCode":"QuotaId"}]}`,
	} {
		c := storageCollectBody(t, `{"value":[`+rows+`]}`)
		got, err := c.Evidence(context.Background(), map[string][]string{"s": {"Storage"}})
		if err != nil || got.Status != "complete" || got.Values["s"].State != 2 {
			t.Fatalf("hard restriction lost order-independent precedence: %#v / %v", got, err)
		}
	}
}
