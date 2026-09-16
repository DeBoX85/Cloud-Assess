package discovery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/config"
)

const (
	SubscriptionStateDisabled = "Disabled"
	SubscriptionStateDeleted  = "Deleted"
)

type Subscription struct {
	ID          string
	DisplayName string
	State       string
}

type SubscriptionLister interface {
	ListSubscriptions(context.Context) ([]Subscription, error)
}

// DiscoverSubscriptions resolves either all accessible active subscriptions or the
// explicitly requested subset, then applies subscription filters.
func DiscoverSubscriptions(
	ctx context.Context,
	lister SubscriptionLister,
	requested []string,
	filters *config.Filters,
) (map[string]string, error) {
	if lister == nil {
		return nil, fmt.Errorf("subscription lister is not configured")
	}
	available, err := lister.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list Azure subscriptions: %w", err)
	}

	requestedSet := make(map[string]struct{}, len(requested))
	for _, id := range requested {
		requestedSet[strings.ToLower(strings.TrimSpace(id))] = struct{}{}
	}
	result := map[string]string{}
	for _, subscription := range available {
		if strings.EqualFold(subscription.State, SubscriptionStateDisabled) || strings.EqualFold(subscription.State, SubscriptionStateDeleted) {
			continue
		}
		if len(requestedSet) > 0 {
			if _, ok := requestedSet[strings.ToLower(subscription.ID)]; !ok {
				continue
			}
		}
		if filters != nil && filters.Assessment != nil && filters.Assessment.IsSubscriptionExcluded(subscription.ID) {
			continue
		}
		result[subscription.ID] = subscription.DisplayName
	}
	return result, nil
}

// ScopeID creates the deterministic identifier used to correlate reports for a resolved scope.
func ScopeID(subscriptions map[string]string) string {
	ids := make([]string, 0, len(subscriptions))
	for id := range subscriptions {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	hash := sha256.Sum256([]byte(strings.Join(ids, ",")))
	return hex.EncodeToString(hash[:])
}
