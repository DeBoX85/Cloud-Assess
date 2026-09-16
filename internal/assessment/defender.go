package assessment

// DefenderPlanStatus is the subscription-level Microsoft Defender for Cloud pricing/tier record.
type DefenderPlanStatus struct {
	SubscriptionID   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName,omitempty"`
	Name             string `json:"name"`
	Tier             string `json:"tier"`
}

// DefenderRecommendation is an unhealthy Microsoft Defender for Cloud assessment record.
type DefenderRecommendation struct {
	SubscriptionID         string `json:"subscriptionId"`
	SubscriptionName       string `json:"subscriptionName,omitempty"`
	ResourceGroup          string `json:"resourceGroup"`
	ResourceType           string `json:"resourceType"`
	ResourceName           string `json:"resourceName"`
	Category               string `json:"category"`
	RecommendationSeverity string `json:"recommendationSeverity"`
	RecommendationName     string `json:"recommendationName"`
	ActionDescription      string `json:"actionDescription,omitempty"`
	RemediationDescription string `json:"remediationDescription,omitempty"`
	AzurePortalLink        string `json:"azurePortalLink,omitempty"`
	ResourceID             string `json:"resourceId"`
}
