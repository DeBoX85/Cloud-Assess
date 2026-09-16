package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// NewCredential creates the default Azure SDK credential chain for the selected cloud.
// Credential selection can be constrained through Azure SDK environment settings such
// as AZURE_TOKEN_CREDENTIALS. Errors are returned to the orchestration boundary.
func NewCredential() (azcore.TokenCredential, error) {
	options := &azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{Cloud: CloudConfiguration()},
	}
	credential, err := azidentity.NewDefaultAzureCredential(options)
	if err != nil {
		return nil, fmt.Errorf("create Azure credential: %w", err)
	}
	return credential, nil
}
