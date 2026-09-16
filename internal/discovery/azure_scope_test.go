package discovery

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

func TestNewAzureScopeClientRejectsNilCredential(t *testing.T) {
	client, err := NewAzureScopeClient(nil, nil)
	if err == nil {
		t.Fatal("NewAzureScopeClient(nil) error = nil, want error")
	}
	if client != nil {
		t.Fatalf("NewAzureScopeClient(nil) client = %#v, want nil", client)
	}
}

func TestValueAndSubscriptionState(t *testing.T) {
	text := "value"
	if got := value(&text); got != text {
		t.Fatalf("value() = %q, want %q", got, text)
	}
	if got := value[string](nil); got != "" {
		t.Fatalf("value(nil) = %q, want empty", got)
	}

	state := armsubscription.SubscriptionStateEnabled
	if got := subscriptionState(&state); got != string(state) {
		t.Fatalf("subscriptionState() = %q, want %q", got, state)
	}
	if got := subscriptionState(nil); got != "" {
		t.Fatalf("subscriptionState(nil) = %q, want empty", got)
	}
}
