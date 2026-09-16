package redact

import "testing"

func TestSubscriptionID(t *testing.T) {
	id := "12345678-1234-1234-1234-123456789012"
	if got := SubscriptionID(id, true); got != "xxxxxxxx-xxxx-xxxx-xxxx-xxxxx6789012" {
		t.Fatalf("masked ID = %q", got)
	}
	if got := SubscriptionID(id, false); got != id {
		t.Fatalf("unmasked ID = %q", got)
	}
	if got := SubscriptionID("short", true); got != "" {
		t.Fatalf("short ID should return empty string, got %q", got)
	}
}

func TestSubscriptionIDInResourceID(t *testing.T) {
	resourceID := "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	want := "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxx6789012/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	if got := SubscriptionIDInResourceID(resourceID, true); got != want {
		t.Fatalf("masked resource ID = %q, want %q", got, want)
	}
	if got := SubscriptionIDInResourceID(resourceID, false); got != resourceID {
		t.Fatalf("unmasked resource ID = %q", got)
	}
	if got := SubscriptionIDInResourceID("/resourceGroups/rg", true); got != "" {
		t.Fatalf("invalid resource ID should return empty string, got %q", got)
	}
}
