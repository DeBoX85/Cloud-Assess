package assessment

import "strings"

const (
	ImpactHigh   = "High"
	ImpactMedium = "Medium"
	ImpactLow    = "Low"
	CategorySLA  = "SLA"
)

// SeverityRank returns High=3, Medium=2, Low=1.
func SeverityRank(impact string) (int, bool) {
	switch {
	case strings.EqualFold(strings.TrimSpace(impact), ImpactHigh):
		return 3, true
	case strings.EqualFold(strings.TrimSpace(impact), ImpactMedium):
		return 2, true
	case strings.EqualFold(strings.TrimSpace(impact), ImpactLow):
		return 1, true
	default:
		return 0, false
	}
}
