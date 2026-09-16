package assessment

import "time"

type StageStatus string

const (
	StageCompleted             StageStatus = "completed"
	StageCompletedWithWarnings StageStatus = "completed_with_warnings"
	StageSkipped               StageStatus = "skipped"
	StageFailed                StageStatus = "failed"
)

type Completeness string

const (
	CompletenessComplete             Completeness = "complete"
	CompletenessCompleteWithWarnings Completeness = "complete_with_warnings"
	CompletenessPartial              Completeness = "partial"
	CompletenessFailed               Completeness = "failed"
)

type AssessmentWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AssessmentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StageExecution struct {
	Name       string              `json:"name"`
	Status     StageStatus         `json:"status"`
	StartedAt  time.Time           `json:"startedAt"`
	FinishedAt time.Time           `json:"finishedAt"`
	Records    int                 `json:"records"`
	Warnings   []AssessmentWarning `json:"warnings,omitempty"`
	Error      *AssessmentError    `json:"error,omitempty"`
}

type LearnMoreLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type RecommendationDefinition struct {
	ID                  string          `json:"id"`
	Recommendation      string          `json:"recommendation"`
	Category            string          `json:"category"`
	Impact              string          `json:"impact"`
	ResourceType        string          `json:"resourceType"`
	State               string          `json:"state,omitempty"`
	LongDescription     string          `json:"longDescription,omitempty"`
	PotentialBenefits   string          `json:"potentialBenefits,omitempty"`
	AutomationAvailable string          `json:"automationAvailable,omitempty"`
	Tags                []string        `json:"tags,omitempty"`
	Query               string          `json:"-"`
	LearnMore           []LearnMoreLink `json:"learnMore,omitempty"`
	Source               string          `json:"source"`
}

type Finding struct {
	RecommendationID    string            `json:"recommendationId"`
	Source              string            `json:"source"`
	Category            string            `json:"category"`
	Impact              string            `json:"impact"`
	ResourceType        string            `json:"resourceType"`
	Recommendation      string            `json:"recommendation"`
	LongDescription     string            `json:"longDescription,omitempty"`
	PotentialBenefits   string            `json:"potentialBenefits,omitempty"`
	ResourceID          string            `json:"resourceId"`
	SubscriptionID      string            `json:"subscriptionId"`
	SubscriptionName    string            `json:"subscriptionName,omitempty"`
	ResourceGroup       string            `json:"resourceGroup"`
	ResourceName        string            `json:"resourceName"`
	Tags                map[string]string `json:"tags,omitempty"`
	Parameters          []string          `json:"parameters,omitempty"`
	LearnMoreURL        string            `json:"learnMoreUrl,omitempty"`
	AutomationAvailable string            `json:"automationAvailable,omitempty"`
}
