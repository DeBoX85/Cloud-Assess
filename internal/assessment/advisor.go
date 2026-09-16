package assessment

// AdvisorRecommendation is the normalized auxiliary Advisor dataset retained separately
// from primary assessment findings, matching the reference report architecture.
type AdvisorRecommendation struct {
	RecommendationID string `json:"recommendationId"`
	SubscriptionID   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName,omitempty"`
	ResourceType     string `json:"resourceType"`
	ResourceName     string `json:"resourceName"`
	ResourceID       string `json:"resourceId"`
	Category         string `json:"category"`
	Impact           string `json:"impact"`
	Description      string `json:"description"`
}
