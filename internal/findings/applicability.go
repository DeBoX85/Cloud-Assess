package findings

import "strings"

type Applicability string

const (
	NotApplicable Applicability = "not_applicable"
	Compliant     Applicability = "compliant"
	NonCompliant  Applicability = "noncompliant"
)

// ApplicabilityFor reproduces the reference report semantics without exposing
// the legacy string values N/A/true/false as the internal domain model.
func ApplicabilityFor(resourceType string, impactedResources int, deployedResourceTypes map[string]bool) Applicability {
	if strings.EqualFold(resourceType, "Microsoft.Resources") {
		return stateFromImpact(impactedResources)
	}
	if !deployedResourceTypes[strings.ToLower(resourceType)] {
		return NotApplicable
	}
	return stateFromImpact(impactedResources)
}

func stateFromImpact(impactedResources int) Applicability {
	if impactedResources == 0 {
		return Compliant
	}
	return NonCompliant
}
