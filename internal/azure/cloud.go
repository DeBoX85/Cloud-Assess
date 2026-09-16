package azure

import (
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

const (
	EnvAzureCloud                   = "AZURE_CLOUD"
	EnvAzureAuthorityHost           = "AZURE_AUTHORITY_HOST"
	EnvAzureResourceManagerEndpoint = "AZURE_RESOURCE_MANAGER_ENDPOINT"
	EnvAzureResourceManagerAudience = "AZURE_RESOURCE_MANAGER_AUDIENCE"
)

// CloudConfiguration returns the Azure cloud configuration selected by environment variables.
func CloudConfiguration() cloud.Configuration {
	return cloudConfigurationFrom(os.Getenv)
}

func cloudConfigurationFrom(getenv func(string) string) cloud.Configuration {
	authorityHost := getenv(EnvAzureAuthorityHost)
	armEndpoint := getenv(EnvAzureResourceManagerEndpoint)
	if authorityHost != "" && armEndpoint != "" {
		service := cloud.ServiceConfiguration{Endpoint: armEndpoint}
		if audience := getenv(EnvAzureResourceManagerAudience); audience != "" {
			service.Audience = audience
		}
		return cloud.Configuration{
			ActiveDirectoryAuthorityHost: authorityHost,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: service,
			},
		}
	}

	switch strings.ToLower(strings.TrimSpace(getenv(EnvAzureCloud))) {
	case "azuregovernment", "azureusgovernment", "usgovernment":
		return cloud.AzureGovernment
	case "azurechina", "china":
		return cloud.AzureChina
	case "azurepublic", "public", "":
		return cloud.AzurePublic
	default:
		return cloud.AzurePublic
	}
}

// ResourceManagerEndpoint returns the selected ARM endpoint without a trailing slash.
func ResourceManagerEndpoint() string {
	configuration := CloudConfiguration()
	if service, ok := configuration.Services[cloud.ResourceManager]; ok && service.Endpoint != "" {
		return strings.TrimSuffix(service.Endpoint, "/")
	}
	return strings.TrimSuffix(cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint, "/")
}

// ResourceManagerScope returns the OAuth scope used by ARM-backed APIs.
func ResourceManagerScope() string {
	return ResourceManagerEndpoint() + "/.default"
}
