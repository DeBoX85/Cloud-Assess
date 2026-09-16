package azure

import "testing"

const testResourceID = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myRG/providers/Microsoft.Compute/virtualMachines/myVM"

func TestResourceIDHelpers(t *testing.T) {
	if got := SubscriptionFromResourceID(testResourceID); got != "12345678-1234-1234-1234-123456789012" {
		t.Fatalf("subscription = %q", got)
	}
	if got := ResourceGroupFromResourceID(testResourceID); got != "myRG" {
		t.Fatalf("resource group = %q", got)
	}
	if got := ResourceGroupIDFromResourceID(testResourceID); got != "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myRG" {
		t.Fatalf("resource group ID = %q", got)
	}
	if got := ResourceTypeFromResourceID(testResourceID); got != "Microsoft.Compute/virtualMachines" {
		t.Fatalf("resource type = %q", got)
	}
	if got := ResourceNameFromResourceID(testResourceID); got != "myVM" {
		t.Fatalf("resource name = %q", got)
	}
}

func TestMalformedResourceIDReturnsEmptyComponents(t *testing.T) {
	for name, fn := range map[string]func(string) string{
		"subscription": SubscriptionFromResourceID,
		"resourceGroup": ResourceGroupFromResourceID,
		"resourceGroupID": ResourceGroupIDFromResourceID,
		"resourceType": ResourceTypeFromResourceID,
		"resourceName": ResourceNameFromResourceID,
	} {
		if got := fn(""); got != "" {
			t.Fatalf("%s empty input = %q", name, got)
		}
	}
}
