package discovery

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

// AzureScopeClient implements subscription and management-group discovery with Azure SDK clients.
type AzureScopeClient struct {
	subscriptions     *armsubscription.SubscriptionsClient
	managementGroups *armmanagementgroups.ClientFactory
}

var _ SubscriptionLister = (*AzureScopeClient)(nil)
var _ ManagementGroupLister = (*AzureScopeClient)(nil)

// NewAzureScopeClient creates the production scope discovery adapter.
func NewAzureScopeClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*AzureScopeClient, error) {
	if credential == nil {
		return nil, fmt.Errorf("Azure scope discovery credential is nil")
	}
	if options == nil {
		options = &arm.ClientOptions{}
	}

	subscriptions, err := armsubscription.NewSubscriptionsClient(credential, options)
	if err != nil {
		return nil, fmt.Errorf("create Azure subscriptions client: %w", err)
	}
	managementGroups, err := armmanagementgroups.NewClientFactory(credential, options)
	if err != nil {
		return nil, fmt.Errorf("create Azure management groups client factory: %w", err)
	}
	return &AzureScopeClient{
		subscriptions:     subscriptions,
		managementGroups: managementGroups,
	}, nil
}

// ListSubscriptions returns all subscriptions visible to the current Azure identity.
func (c *AzureScopeClient) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	if c == nil || c.subscriptions == nil {
		return nil, fmt.Errorf("Azure subscriptions client is not configured")
	}

	pager := c.subscriptions.NewListPager(nil)
	result := make([]Subscription, 0, 16)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list Azure subscriptions: %w", err)
		}
		for _, item := range page.Value {
			if item == nil || item.SubscriptionID == nil {
				continue
			}
			result = append(result, Subscription{
				ID:          value(item.SubscriptionID),
				DisplayName: value(item.DisplayName),
				State:       subscriptionState(item.State),
			})
		}
	}
	return result, nil
}

// SubscriptionsUnderManagementGroup returns subscriptions directly associated with one management group.
func (c *AzureScopeClient) SubscriptionsUnderManagementGroup(ctx context.Context, groupID string) ([]Subscription, error) {
	if c == nil || c.managementGroups == nil {
		return nil, fmt.Errorf("Azure management groups client is not configured")
	}
	client := c.managementGroups.NewManagementGroupSubscriptionsClient()
	pager := client.NewGetSubscriptionsUnderManagementGroupPager(groupID, nil)
	result := make([]Subscription, 0, 16)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list subscriptions under management group %q: %w", groupID, err)
		}
		for _, item := range page.Value {
			if item == nil || item.Name == nil {
				continue
			}
			entry := Subscription{ID: value(item.Name)}
			if item.Properties != nil {
				entry.DisplayName = value(item.Properties.DisplayName)
				entry.State = value(item.Properties.State)
			}
			result = append(result, entry)
		}
	}
	return result, nil
}

// DescendantManagementGroups returns descendant management-group IDs below one management group.
// The Azure endpoint returns all descendants, not only immediate children. The discovery walker
// deduplicates visited group IDs, so returning every descendant preserves source behavior without
// creating duplicate subscription records.
func (c *AzureScopeClient) DescendantManagementGroups(ctx context.Context, groupID string) ([]string, error) {
	if c == nil || c.managementGroups == nil {
		return nil, fmt.Errorf("Azure management groups client is not configured")
	}
	client := c.managementGroups.NewClient()
	pager := client.NewGetDescendantsPager(groupID, nil)
	result := make([]string, 0, 8)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list descendants for management group %q: %w", groupID, err)
		}
		for _, item := range page.Value {
			if item == nil || item.Type == nil || item.Name == nil {
				continue
			}
			if value(item.Type) == "Microsoft.Management/managementGroups" {
				result = append(result, value(item.Name))
			}
		}
	}
	return result, nil
}

func value[T ~string](pointer *T) string {
	if pointer == nil {
		return ""
	}
	return string(*pointer)
}

func subscriptionState(state *armsubscription.SubscriptionState) string {
	if state == nil {
		return ""
	}
	return string(*state)
}
