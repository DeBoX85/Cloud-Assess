package discovery

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

type fixtureCredential struct{}

func (fixtureCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fixture", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type hierarchyTransport struct {
	requests []string
	failLeaf bool
}

func (f *hierarchyTransport) Do(request *http.Request) (*http.Response, error) {
	if request.Method != http.MethodGet {
		return &http.Response{StatusCode: http.StatusMethodNotAllowed, Body: io.NopCloser(strings.NewReader(`{}`)), Request: request}, nil
	}
	path := request.URL.Path
	f.requests = append(f.requests, path)
	var status = http.StatusOK
	var body string
	switch {
	case strings.HasSuffix(path, "/managementGroups/root/subscriptions"):
		body = `{"value":[{"name":"sub-root","properties":{"displayName":"Root","state":"Active"}}]}`
	case strings.HasSuffix(path, "/managementGroups/root/descendants"):
		// Azure returns all descendants. The child endpoint also returns the leaf.
		body = `{"value":[{"type":"Microsoft.Management/managementGroups","name":"child"},{"type":"Microsoft.Management/managementGroups","name":"leaf"},{"type":"Microsoft.Management/managementGroups/subscriptions","name":"sub-root"}]}`
	case strings.HasSuffix(path, "/managementGroups/child/subscriptions"):
		body = `{"value":[{"name":"sub-child","properties":{"displayName":"Child","state":"Active"}}]}`
	case strings.HasSuffix(path, "/managementGroups/child/descendants"):
		body = `{"value":[{"type":"Microsoft.Management/managementGroups","name":"leaf"}]}`
	case strings.HasSuffix(path, "/managementGroups/leaf/subscriptions"):
		if f.failLeaf {
			status = http.StatusForbidden
			body = `{"error":{"code":"AuthorizationFailed","message":"fixture denial"}}`
		} else {
			body = `{"value":[{"name":"sub-leaf","properties":{"displayName":"Leaf","state":"Active"}},{"name":"sub-disabled","properties":{"state":"Disabled"}}]}`
		}
	case strings.HasSuffix(path, "/managementGroups/leaf/descendants"):
		body = `{"value":[]}`
	default:
		status = http.StatusNotFound
		body = `{"error":{"code":"UnknownFixturePath"}}`
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: request}, nil
}

func TestAzureScopeAdapterResolvesSyntheticNestedHierarchy(t *testing.T) {
	transport := &hierarchyTransport{}
	client, err := NewAzureScopeClient(fixtureCredential{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}})
	if err != nil {
		t.Fatal(err)
	}
	filters := config.NewFilters()
	filters.Assessment.Include.Subscriptions = []string{"sub-root", "sub-leaf"}
	filters.RebuildIndexes()
	resolved, err := DiscoverManagementGroupSubscriptions(context.Background(), client, []string{"root"}, filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 2 || resolved["sub-root"] != "Root" || resolved["sub-leaf"] != "Leaf" {
		t.Fatalf("resolved subscriptions = %#v, want root and leaf only", resolved)
	}
	if len(transport.requests) != 6 {
		t.Fatalf("HTTP requests = %#v, want one subscriptions and descendants request per unique group", transport.requests)
	}
}

func TestAzureScopeAdapterDescendantDenialFailsWithoutPartialScope(t *testing.T) {
	transport := &hierarchyTransport{failLeaf: true}
	client, err := NewAzureScopeClient(fixtureCredential{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := DiscoverManagementGroupSubscriptions(context.Background(), client, []string{"root"}, nil)
	if err == nil || resolved != nil || !strings.Contains(err.Error(), "leaf") {
		t.Fatalf("resolved = %#v, error = %v; want explicit descendant failure without partial scope", resolved, err)
	}
}

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
