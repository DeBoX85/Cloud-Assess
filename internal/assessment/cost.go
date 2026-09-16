package assessment

import "time"

// CostRecord is one previous-month actual-cost aggregation for an Azure service.
type CostRecord struct {
	From             time.Time `json:"from"`
	To               time.Time `json:"to"`
	SubscriptionID   string    `json:"subscriptionId"`
	SubscriptionName string    `json:"subscriptionName,omitempty"`
	ServiceName      string    `json:"serviceName"`
	Value            string    `json:"value"`
	Currency         string    `json:"currency"`
}
