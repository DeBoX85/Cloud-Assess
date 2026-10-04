package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
)

func TestAuditInvalidFilterSchemaStopsBeforeCredentials(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{"legacy-root", "azqr:\n  include:\n    subscriptions: [expected-subscription]\n"},
		{"unknown-include-field", "assessment:\n  include:\n    subscriptionIds: [expected-subscription]\n"},
		{"legacy-resource-exclusion", "assessment:\n  exclude:\n    services: [/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/retained]\n"},
		{"null-assessment", "assessment: null\n"},
		{"second-document", "assessment: {}\n---\nassessment:\n  include:\n    subscriptions: [expected-subscription]\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "scope-filter.yml")
			if err := os.WriteFile(path, []byte(test.yaml), 0600); err != nil {
				t.Fatal(err)
			}
			credentials, operations := 0, 0
			code, err := executeScanWithFactories(context.Background(), scanFlags{
				filtersFile: path, pluginOnly: true, internalPlugins: []string{"zone-mapping"},
			}, func() (azcore.TokenCredential, error) {
				credentials++
				return nil, errors.New("credential factory tripwire")
			}, func(azcore.TokenCredential) (orchestration.Operations, error) {
				operations++
				return orchestration.Operations{}, errors.New("operations factory tripwire")
			})
			if code != 1 || err == nil || credentials != 0 || operations != 0 || !strings.Contains(err.Error(), "filter") {
				t.Fatalf("invalid filter reached credential boundary: code=%d credentials=%d operations=%d error=%v", code, credentials, operations, err)
			}
		})
	}
}

func TestAuditHealthyFilterSchemaReachesExpectedBoundary(t *testing.T) {
	for _, content := range []string{"", "assessment: {}\n", "assessment:\n  include:\n    subscriptions: [expected-subscription]\n  exclude:\n    resources: [/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/excluded]\n"} {
		path := filepath.Join(t.TempDir(), "valid-filter.yml")
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		calls := 0
		code, err := executeScanWithFactories(context.Background(), scanFlags{
			filtersFile: path, pluginOnly: true, internalPlugins: []string{"zone-mapping"},
		}, func() (azcore.TokenCredential, error) {
			calls++
			return nil, errors.New("healthy credential boundary")
		}, func(azcore.TokenCredential) (orchestration.Operations, error) {
			t.Fatal("operation factory reached")
			return orchestration.Operations{}, nil
		})
		if code != 1 || calls != 1 || err == nil || err.Error() != "healthy credential boundary" {
			t.Fatalf("healthy filter rejected before expected boundary: code=%d calls=%d error=%v", code, calls, err)
		}
	}
}
