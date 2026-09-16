package redact

import "strings"

const maskedSubscriptionPrefix = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxx"

func SubscriptionID(subscriptionID string, enabled bool) string {
	if len(subscriptionID) < 36 {
		return ""
	}
	if !enabled {
		return subscriptionID
	}
	return maskedSubscriptionPrefix + subscriptionID[29:]
}

func SubscriptionIDInResourceID(resourceID string, enabled bool) string {
	if len(resourceID) < 51 || !strings.HasPrefix(resourceID, "/subscriptions/") {
		return ""
	}
	if !enabled {
		return resourceID
	}
	return "/subscriptions/" + maskedSubscriptionPrefix + resourceID[44:]
}
