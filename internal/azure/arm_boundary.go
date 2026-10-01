package azure

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// ScopeARMOptions copies caller options and constrains SDK scope requests before
// authentication. Caller-supplied transports remain trusted code and must not
// follow redirects internally or mutate destinations after pipeline validation.
func ScopeARMOptions(options *arm.ClientOptions) (*arm.ClientOptions, error) {
	result := arm.ClientOptions{}
	if options != nil {
		result = *options
	}
	result.DisableRPRegistration = true
	configuration := result.Cloud
	endpoint := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if service, ok := configuration.Services[cloud.ResourceManager]; ok {
		endpoint = service.Endpoint
	}
	origin, err := url.Parse(endpoint)
	if err != nil || origin.Scheme != "https" || origin.Host == "" || origin.User != nil || origin.Fragment != "" || origin.RawQuery != "" || strings.Contains(endpoint, `\`) {
		return nil, fmt.Errorf("invalid HTTPS ARM scope endpoint")
	}
	policies := make([]policy.Policy, 0, len(result.PerCallPolicies)+1)
	policies = append(policies, armOriginPolicy{origin: origin})
	result.PerCallPolicies = append(policies, result.PerCallPolicies...)
	if result.Transport == nil {
		result.Transport = &http.Client{Transport: http.DefaultTransport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	return &result, nil
}

type armOriginPolicy struct{ origin *url.URL }

func (p armOriginPolicy) Do(request *policy.Request) (*http.Response, error) {
	target := request.Raw().URL
	if target == nil || target.Scheme != "https" || !strings.EqualFold(target.Host, p.origin.Host) || target.User != nil || target.Fragment != "" || strings.Contains(target.String(), `\`) {
		return nil, fmt.Errorf("ARM scope request rejected: unsafe destination")
	}
	return request.Next()
}
