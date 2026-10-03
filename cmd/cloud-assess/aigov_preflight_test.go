package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
)

func TestAICloudPreflightBeforeFactoriesAndScope(t *testing.T) {
	keys := []string{azure.EnvAzureAuthorityHost, azure.EnvAzureResourceManagerEndpoint, azure.EnvAzureResourceManagerAudience}
	public := []string{"https://login.microsoftonline.com/", "https://management.azure.com", "https://management.core.windows.net"}
	for _, mode := range []string{"standalone", "mixed", "scanner"} {
		for mask := 0; mask < 12; mask++ {
			t.Run(fmt.Sprintf("%s-%d", mode, mask), func(t *testing.T) {
				cloudName := "public"
				if mask == 8 {
					cloudName = "AzureChina"
				}
				if mask == 9 {
					cloudName = "AzureGovernment"
				}
				if mask == 10 {
					cloudName = "unrecognized"
				}
				t.Setenv(azure.EnvAzureCloud, cloudName)
				for i, key := range keys {
					value := ""
					if mask < 8 && mask&(1<<i) != 0 {
						value = public[i]
					}
					if mask == 11 {
						value = []string{"https://login.invalid", "https://arm.invalid", "https://audience.invalid"}[i]
					}
					t.Setenv(key, value)
				}
				credentialCalls, factoryCalls := 0, 0
				f := scanFlags{internalPlugins: []string{"ai-gov"}, pluginOnly: mode == "standalone"}
				if mode == "scanner" {
					f.scannerKeys = []string{"ca"}
				}
				code, err := executeScanWithFactories(context.Background(), f, func() (azcore.TokenCredential, error) {
					credentialCalls++
					return nil, errors.New("credential factory reached")
				}, func(azcore.TokenCredential) (orchestration.Operations, error) {
					factoryCalls++
					return orchestration.Operations{}, errors.New("scope factory tripwire")
				})
				allowed := mask == 0 || mask == 7
				if code != 1 || err == nil || factoryCalls != 0 || allowed && credentialCalls != 1 || !allowed && (credentialCalls != 0 || !strings.Contains(err.Error(), "cloud configuration required")) {
					t.Fatalf("AI preflight ordering: credentials=%d factories=%d error=%v", credentialCalls, factoryCalls, err)
				}
			})
		}
	}
}
