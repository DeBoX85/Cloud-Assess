package discovery

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

type ManagementGroupLister interface {
	SubscriptionsUnderManagementGroup(context.Context, string) ([]Subscription, error)
	DescendantManagementGroups(context.Context, string) ([]string, error)
}

// DiscoverManagementGroupSubscriptions recursively resolves active subscriptions under
// the requested management groups. A visited set prevents duplicate descendant traversal
// while preserving the source implementation's resolved subscription semantics.
func DiscoverManagementGroupSubscriptions(
	ctx context.Context,
	lister ManagementGroupLister,
	groups []string,
	filters *config.Filters,
) (map[string]string, error) {
	if lister == nil {
		return nil, fmt.Errorf("management group lister is not configured")
	}
	result := map[string]string{}
	visited := map[string]struct{}{}

	var visit func(string) error
	visit = func(group string) error {
		normalized := strings.ToLower(strings.TrimSpace(group))
		if normalized == "" {
			return nil
		}
		if _, ok := visited[normalized]; ok {
			return nil
		}
		visited[normalized] = struct{}{}

		subscriptions, err := lister.SubscriptionsUnderManagementGroup(ctx, group)
		if err != nil {
			return fmt.Errorf("list subscriptions under management group %s: %w", group, err)
		}
		for _, subscription := range subscriptions {
			if strings.EqualFold(subscription.State, SubscriptionStateDisabled) || strings.EqualFold(subscription.State, SubscriptionStateDeleted) {
				continue
			}
			if filters != nil && filters.Assessment != nil && filters.Assessment.IsSubscriptionExcluded(subscription.ID) {
				continue
			}
			result[subscription.ID] = subscription.DisplayName
		}

		descendants, err := lister.DescendantManagementGroups(ctx, group)
		if err != nil {
			return fmt.Errorf("list descendants of management group %s: %w", group, err)
		}
		for _, descendant := range descendants {
			if err := visit(descendant); err != nil {
				return err
			}
		}
		return nil
	}

	for _, group := range groups {
		if err := visit(group); err != nil {
			return nil, err
		}
	}
	return result, nil
}
