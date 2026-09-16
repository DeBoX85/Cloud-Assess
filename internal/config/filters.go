package config

import (
	"fmt"
	"strings"
)

type Filters struct {
	Assessment *AssessmentFilter `yaml:"assessment" json:"assessment"`
}

type AssessmentFilter struct {
	Include *IncludeFilter `yaml:"include" json:"include"`
	Exclude *ExcludeFilter `yaml:"exclude" json:"exclude"`

	includeSubscriptions   map[string]bool
	includeResourceGroups  map[string]bool
	includeResourceTypes   map[string]bool
	includeTags            map[string]string
	excludeSubscriptions   map[string]bool
	excludeResourceGroups  map[string]bool
	excludeResources       map[string]bool
	excludeRecommendations map[string]bool
	excludeTags            map[string]string
	resourceScope          map[string]bool
}

type IncludeFilter struct {
	Subscriptions  []string          `yaml:"subscriptions,flow" json:"subscriptions"`
	ResourceGroups []string          `yaml:"resourceGroups,flow" json:"resourceGroups"`
	ResourceTypes  []string          `yaml:"resourceTypes,flow" json:"resourceTypes"`
	Tags           map[string]string `yaml:"tags" json:"tags"`
}

type ExcludeFilter struct {
	Subscriptions   []string          `yaml:"subscriptions,flow" json:"subscriptions"`
	ResourceGroups  []string          `yaml:"resourceGroups,flow" json:"resourceGroups"`
	Resources       []string          `yaml:"resources,flow" json:"resources"`
	Recommendations []string          `yaml:"recommendations,flow" json:"recommendations"`
	Tags            map[string]string `yaml:"tags" json:"tags"`
}

func NewFilters() *Filters {
	f := &Filters{Assessment: &AssessmentFilter{
		Include: &IncludeFilter{Tags: map[string]string{}},
		Exclude: &ExcludeFilter{Tags: map[string]string{}},
	}}
	f.Assessment.initialize()
	return f
}

func (f *AssessmentFilter) initialize() {
	if f.Include == nil {
		f.Include = &IncludeFilter{Tags: map[string]string{}}
	}
	if f.Exclude == nil {
		f.Exclude = &ExcludeFilter{Tags: map[string]string{}}
	}
	f.includeSubscriptions = stringSet(f.Include.Subscriptions)
	f.includeResourceGroups = stringSet(f.Include.ResourceGroups)
	f.includeResourceTypes = stringSet(f.Include.ResourceTypes)
	f.includeTags = normalizeTags(f.Include.Tags)
	f.excludeSubscriptions = stringSet(f.Exclude.Subscriptions)
	f.excludeResourceGroups = stringSet(f.Exclude.ResourceGroups)
	f.excludeResources = stringSet(f.Exclude.Resources)
	f.excludeRecommendations = stringSet(f.Exclude.Recommendations)
	f.excludeTags = normalizeTags(f.Exclude.Tags)
	f.resourceScope = map[string]bool{}
}

func (f *AssessmentFilter) Validate() error {
	for _, id := range append(append([]string{}, f.Include.ResourceGroups...), f.Exclude.ResourceGroups...) {
		if err := ValidateResourceGroupID(id); err != nil {
			return err
		}
	}
	return nil
}

func ValidateResourceGroupID(resourceGroupID string) error {
	parts := strings.Split(resourceGroupID, "/")
	if len(parts) != 5 || parts[0] != "" || !strings.EqualFold(parts[1], "subscriptions") || !strings.EqualFold(parts[3], "resourceGroups") {
		return fmt.Errorf("resource group ID %q has incorrect format; expected /subscriptions/{subscription-id}/resourceGroups/{resource-group-name}", resourceGroupID)
	}
	if parts[2] == "" {
		return fmt.Errorf("resource group ID %q has empty subscription ID", resourceGroupID)
	}
	if parts[4] == "" {
		return fmt.Errorf("resource group ID %q has empty resource group name", resourceGroupID)
	}
	return nil
}

func (f *AssessmentFilter) IsSubscriptionExcluded(subscriptionID string) bool {
	id := normalize(subscriptionID)
	if f.includeSubscriptions[id] {
		return false
	}
	if len(f.includeSubscriptions) > 0 {
		return true
	}
	return f.excludeSubscriptions[id]
}

func (f *AssessmentFilter) IsResourceGroupExcluded(resourceGroupID string) bool {
	id := normalize(resourceGroupID)
	if f.includeResourceGroups[id] {
		return false
	}
	if len(f.includeResourceGroups) > 0 {
		return true
	}
	return f.excludeResourceGroups[id]
}

func (f *AssessmentFilter) IsResourceTypeExcluded(resourceType string) bool {
	if len(f.includeResourceTypes) == 0 {
		return false
	}
	return !f.includeResourceTypes[normalize(resourceType)]
}

func (f *AssessmentFilter) IsRecommendationExcluded(recommendationID string) bool {
	return f.excludeRecommendations[normalize(recommendationID)]
}

func (f *AssessmentFilter) IsResourceExcluded(resourceID, subscriptionID, resourceGroup, resourceType string, tags map[string]string) bool {
	if f.IsSubscriptionExcluded(subscriptionID) || f.IsResourceGroupExcluded(resourceGroup) || f.IsResourceTypeExcluded(resourceType) || f.excludeResources[normalize(resourceID)] {
		return true
	}
	normalizedTags := normalizeTags(tags)
	for key, value := range f.excludeTags {
		if actual, ok := normalizedTags[key]; ok && actual == value {
			return true
		}
	}
	for key, value := range f.includeTags {
		if actual, ok := normalizedTags[key]; !ok || actual != value {
			return true
		}
	}
	return false
}

func (f *AssessmentFilter) SetResourceScope(resourceID string, included bool) {
	f.resourceScope[normalize(resourceID)] = included
}

func (f *AssessmentFilter) IsServiceExcluded(resourceID string) bool {
	id := normalize(resourceID)
	if f.excludeResources[id] {
		return true
	}
	if len(f.includeTags) == 0 && len(f.excludeTags) == 0 {
		return false
	}
	for {
		if included, ok := f.resourceScope[id]; ok {
			return !included
		}
		lastSlash := strings.LastIndexByte(id, '/')
		if lastSlash <= 0 {
			break
		}
		id = id[:lastSlash]
	}
	return len(f.includeTags) > 0
}

func stringSet(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[normalize(value)] = true
	}
	return out
}

func normalizeTags(tags map[string]string) map[string]string {
	out := make(map[string]string, len(tags))
	for key, value := range tags {
		out[normalize(key)] = value
	}
	return out
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
