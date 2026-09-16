package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

func env(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestCloudConfigurationDefaultsToPublic(t *testing.T) {
	configuration := cloudConfigurationFrom(env(nil))
	if got, want := configuration.Services[cloud.ResourceManager].Endpoint, cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint; got != want {
		t.Fatalf("ARM endpoint = %q, want %q", got, want)
	}
}

func TestCloudConfigurationPredefinedAliases(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected cloud.Configuration
	}{
		{name: "government", value: "AzureGovernment", expected: cloud.AzureGovernment},
		{name: "government alias", value: "usgovernment", expected: cloud.AzureGovernment},
		{name: "china", value: "China", expected: cloud.AzureChina},
		{name: "public", value: " public ", expected: cloud.AzurePublic},
		{name: "unknown falls back public", value: "something-else", expected: cloud.AzurePublic},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration := cloudConfigurationFrom(env(map[string]string{EnvAzureCloud: tt.value}))
			got := configuration.Services[cloud.ResourceManager].Endpoint
			want := tt.expected.Services[cloud.ResourceManager].Endpoint
			if got != want {
				t.Fatalf("ARM endpoint = %q, want %q", got, want)
			}
		})
	}
}

func TestCustomCloudTakesPrecedence(t *testing.T) {
	configuration := cloudConfigurationFrom(env(map[string]string{
		EnvAzureCloud:                   "AzureChina",
		EnvAzureAuthorityHost:           "https://login.example.test/",
		EnvAzureResourceManagerEndpoint: "https://management.example.test/",
		EnvAzureResourceManagerAudience: "https://audience.example.test/",
	}))
	if configuration.ActiveDirectoryAuthorityHost != "https://login.example.test/" {
		t.Fatalf("authority host = %q", configuration.ActiveDirectoryAuthorityHost)
	}
	service := configuration.Services[cloud.ResourceManager]
	if service.Endpoint != "https://management.example.test/" || service.Audience != "https://audience.example.test/" {
		t.Fatalf("unexpected custom resource manager service: %+v", service)
	}
}

func TestIncompleteCustomCloudFallsBackToNamedCloud(t *testing.T) {
	configuration := cloudConfigurationFrom(env(map[string]string{
		EnvAzureCloud:                   "AzureGovernment",
		EnvAzureResourceManagerEndpoint: "https://management.example.test/",
	}))
	if got, want := configuration.Services[cloud.ResourceManager].Endpoint, cloud.AzureGovernment.Services[cloud.ResourceManager].Endpoint; got != want {
		t.Fatalf("ARM endpoint = %q, want %q", got, want)
	}
}
