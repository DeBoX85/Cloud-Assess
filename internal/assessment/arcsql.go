package assessment

// ArcSQLRecord is one Arc-enabled SQL Server inventory/status record.
type ArcSQLRecord struct {
	SubscriptionID   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName,omitempty"`
	Status           string `json:"status"`
	AzureArcServer   string `json:"azureArcServer"`
	SQLInstance      string `json:"sqlInstance"`
	ResourceGroup    string `json:"resourceGroup"`
	Version          string `json:"version"`
	Build            string `json:"build"`
	PatchLevel       string `json:"patchLevel"`
	Edition          string `json:"edition"`
	VCores           string `json:"vCores"`
	License          string `json:"license"`
	DPSStatus        string `json:"dpsStatus"`
	TELStatus        string `json:"telStatus"`
	DefenderStatus   string `json:"defenderStatus"`
}
