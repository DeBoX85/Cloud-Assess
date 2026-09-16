package azure

import (
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/throttling"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// NewARMClientOptions returns the shared ARM SDK configuration used by production Azure adapters.
func NewARMClientOptions() *arm.ClientOptions {
	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Retry: policy.RetryOptions{
				RetryDelay:    4 * time.Second,
				MaxRetryDelay: 60 * time.Second,
				MaxRetries:    5,
			},
			Cloud:            CloudConfiguration(),
			PerRetryPolicies: []policy.Policy{throttling.NewPolicy()},
		},
	}
}
