package azure

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func TestScopeARMOptionsOriginAndCallerOwnership(t *testing.T) {
	for _, endpoint := range []string{"https://management.azure.com", "https://management.usgovcloudapi.net", "https://management.chinacloudapi.cn", "https://arm.example.invalid:8443"} {
		t.Run(endpoint, func(t *testing.T) {
			originalPolicy := armOriginPolicy{}
			backing := make([]policy.Policy, 1, 4)
			backing[0] = originalPolicy
			input := &arm.ClientOptions{ClientOptions: policy.ClientOptions{Cloud: cloud.Configuration{Services: map[cloud.ServiceName]cloud.ServiceConfiguration{cloud.ResourceManager: {Endpoint: endpoint, Audience: endpoint}}}, PerCallPolicies: backing}}
			output, err := ScopeARMOptions(input)
			if err != nil {
				t.Fatal(err)
			}
			if input.DisableRPRegistration || len(input.PerCallPolicies) != 1 || !output.DisableRPRegistration || len(output.PerCallPolicies) != 2 {
				t.Fatal("caller options changed or hardening missing")
			}
			transport, ok := output.Transport.(*http.Client)
			if !ok || transport.CheckRedirect(nil, nil) != http.ErrUseLastResponse {
				t.Fatal("default transport permits redirects")
			}
			req, err := runtime.NewRequest(context.Background(), http.MethodGet, endpoint+"/subscriptions")
			if err != nil {
				t.Fatal(err)
			}
			// An allowed origin reaches Next (no policies in this synthetic request).
			_, err = output.PerCallPolicies[0].Do(req)
			if err == nil || err.Error() != "no more policies" {
				t.Fatalf("valid origin blocked: %v", err)
			}
		})
	}
	for _, endpoint := range []string{"http://arm.invalid", "https://user@arm.invalid", "https://arm.invalid#fragment", "https://arm.invalid?query=1", "\\bad"} {
		_, err := ScopeARMOptions(&arm.ClientOptions{ClientOptions: policy.ClientOptions{Cloud: cloud.Configuration{Services: map[cloud.ServiceName]cloud.ServiceConfiguration{cloud.ResourceManager: {Endpoint: endpoint}}}}})
		if err == nil {
			t.Fatalf("unsafe endpoint accepted: %q", endpoint)
		}
	}
}
