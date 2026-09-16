package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

const ResourceInventoryQuery = "resources | project id=tostring(id), subscriptionId=tostring(subscriptionId), resourceGroup=tostring(resourceGroup), location=tostring(location), type=tostring(type), name=tostring(name), tags, skuName=tostring(coalesce(sku.name, properties.sku.name, properties.hardwareProfile.vmSize, properties.tier, sku)), skuTier=tostring(coalesce(sku.tier, properties.sku.tier)), skuFamily=tostring(coalesce(sku.family, properties.sku.family)), skuCapacity=tolong(coalesce(sku.capacity, properties.sku.capacity, 0)), ['kind']=tostring(kind) | order by subscriptionId, resourceGroup"

type ResourceInventory struct {
	Included      []assessment.Resource
	Excluded      []assessment.Resource
	MalformedRows int
}

type resourceQuerier interface {
	Query(context.Context, string, map[string]string, ...arg.QueryOptions) (*arg.Result, error)
}

type resourceRow struct {
	ID             string            `json:"id"`
	SubscriptionID string            `json:"subscriptionId"`
	ResourceGroup  string            `json:"resourceGroup"`
	Location       string            `json:"location"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	SKUName        string            `json:"skuName"`
	SKUTier        string            `json:"skuTier"`
	SKUFamily      string            `json:"skuFamily"`
	SKUCapacity    int               `json:"skuCapacity"`
	Kind           string            `json:"kind"`
	Tags           map[string]string `json:"tags"`
}

// DiscoverResources queries Azure Resource Graph once for inventory and partitions the
// result into in-scope and out-of-scope resources using the assessment filter model.
func DiscoverResources(
	ctx context.Context,
	querier resourceQuerier,
	subscriptions map[string]string,
	filters *config.Filters,
) (*ResourceInventory, error) {
	if querier == nil {
		return nil, fmt.Errorf("resource inventory querier is not configured")
	}

	result, err := querier.Query(ctx, ResourceInventoryQuery, subscriptions)
	if err != nil {
		return nil, fmt.Errorf("query Azure Resource Graph for resource inventory: %w", err)
	}
	inventory := &ResourceInventory{
		Included: make([]assessment.Resource, 0, len(result.Data)),
		Excluded: []assessment.Resource{},
	}

	for _, raw := range result.Data {
		var row resourceRow
		if err := json.Unmarshal(raw, &row); err != nil {
			inventory.MalformedRows++
			continue
		}
		resource := assessment.Resource{
			ID:             row.ID,
			SubscriptionID: row.SubscriptionID,
			ResourceGroup:  row.ResourceGroup,
			Location:       row.Location,
			Type:           row.Type,
			Name:           row.Name,
			SKUName:        row.SKUName,
			SKUTier:        row.SKUTier,
			SKUFamily:      row.SKUFamily,
			SKUCapacity:    row.SKUCapacity,
			Kind:           row.Kind,
			Tags:           row.Tags,
		}

		excluded := false
		if filters != nil && filters.Assessment != nil {
			rgID := azure.ResourceGroupIDFromResourceID(resource.ID)
			excluded = filters.Assessment.IsResourceExcluded(
				resource.ID,
				resource.SubscriptionID,
				rgID,
				resource.Type,
				resource.Tags,
			)
			filters.Assessment.SetResourceScope(resource.ID, !excluded)
		}
		if excluded {
			inventory.Excluded = append(inventory.Excluded, resource)
			continue
		}
		inventory.Included = append(inventory.Included, resource)
	}
	return inventory, nil
}

// CountResourcesByTypeAndSubscription derives deterministic counts from the filtered inventory.
func CountResourcesByTypeAndSubscription(resources []assessment.Resource, subscriptions map[string]string) []assessment.ResourceTypeCount {
	type countKey struct {
		subscriptionID string
		resourceType   string
	}
	counts := map[countKey]int{}
	for _, resource := range resources {
		counts[countKey{subscriptionID: resource.SubscriptionID, resourceType: resource.Type}]++
	}

	out := make([]assessment.ResourceTypeCount, 0, len(counts))
	for key, count := range counts {
		out = append(out, assessment.ResourceTypeCount{
			SubscriptionID:   key.subscriptionID,
			SubscriptionName: subscriptions[key.subscriptionID],
			ResourceType:     key.resourceType,
			Count:            count,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].SubscriptionName != out[j].SubscriptionName {
			return out[i].SubscriptionName < out[j].SubscriptionName
		}
		if out[i].ResourceType != out[j].ResourceType {
			return out[i].ResourceType < out[j].ResourceType
		}
		return out[i].SubscriptionID < out[j].SubscriptionID
	})
	return out
}
