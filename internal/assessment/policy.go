package assessment

// PolicyNonCompliance is one non-compliant Azure Policy state associated with a resource.
type PolicyNonCompliance struct {
	SubscriptionID       string `json:"subscriptionId"`
	SubscriptionName     string `json:"subscriptionName,omitempty"`
	ResourceGroup        string `json:"resourceGroup"`
	ResourceType         string `json:"resourceType"`
	ResourceName         string `json:"resourceName"`
	PolicyDisplayName    string `json:"policyDisplayName"`
	PolicyDescription    string `json:"policyDescription,omitempty"`
	ResourceID           string `json:"resourceId"`
	Timestamp            string `json:"timestamp"`
	PolicyDefinitionName string `json:"policyDefinitionName"`
	PolicyDefinitionID   string `json:"policyDefinitionId"`
	PolicyAssignmentName string `json:"policyAssignmentName"`
	PolicyAssignmentID   string `json:"policyAssignmentId"`
	ComplianceState      string `json:"complianceState"`
}
