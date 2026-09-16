package config

import (
	"fmt"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

type Filters struct {
	Assessment *AssessmentFilter `yaml:"assessment" json:"assessment"`
}

type AssessmentFilter struct {
	Include *IncludeFilter `yaml:"include" json:"include"`
	Exclude *ExcludeFilter `yaml:"exclude" json:"exclude"`

	includeSubscriptions   map[string]bool
	includeResourceGroups  map[string]bool
	includeTags            map[string]string
	excludeSubscriptions   map[string]bool
	excludeResourceGroups  map[string]bool
	excludeResources       map[string]bool
	excludeRecommendations map[string]bool
	excludeTags            map[string]string
	resourceScope          map[string]bool

	// allowedResourceTypes is runtime scanner scope, not raw user configuration.
	// The reference implementation builds this set from the selected scanner keys and
	// applies it to both inventory and downstream findings. It must therefore be set by
	// scanner selection before resource/finding filtering begins.
	allowedResourceTypes map[string]bool
}

type IncludeFilter struct {
	Subscriptions  []string          `yaml:"subscriptions,flow" json:"subscriptions"`
	ResourceGroups []string          `yaml:"resourceGroups,flow" json:"resourceGroups"`
	// ResourceTypes preserves the reference configuration name. Despite the name,
	// values are scanner/service keys such as "aks", "ca", and "st", not ARM type strings.
	ResourceTypes []string          `yaml:"resourceTypes,flow" json:"resourceTypes"`
	Tags          map[string]string `yaml:"tags" json:"tags"`
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
	f.RebuildIndexes()
	return f
}

func (f *Filters) RebuildIndexes() {
	if f.Assessment == nil {
		f.Assessment = &AssessmentFilter{}
	}
	f.Assessment.rebuildIndexes()
}

func (f *AssessmentFilter) rebuildIndexes() {
	if f.Include == nil {
		f.Include = &IncludeFilter{Tags: map[string]string{}}
	}
	if f.Exclude == nil {
		f.Exclude = &ExcludeFilter{Tags: map[string]string{}}
	}
	if f.Include.Tags == nil {
		f.Include.Tags = map[string]string{}
	}
	if f.Exclude.Tags == nil {
		f.Exclude.Tags = map[string]string{}
	}
	f.includeSubscriptions = stringSet(f.Include.Subscriptions)
	f.includeResourceGroups = stringSet(f.Include.ResourceGroups)
	f.includeTags = normalizeTags(f.Include.Tags)
	f.excludeSubscriptions = stringSet(f.Exclude.Subscriptions)
	f.excludeResourceGroups = stringSet(f.Exclude.ResourceGroups)
	f.excludeResources = stringSet(f.Exclude.Resources)
	f.excludeRecommendations = stringSet(f.Exclude.Recommendations)
	f.excludeTags = normalizeTags(f.Exclude.Tags)
	f.resourceScope = map[string]bool{}
	// Scanner scope depends on Include.ResourceTypes and the CLI scanner selection,
	// so rebuilding config indexes intentionally invalidates any previous runtime scope.
	f.allowedResourceTypes = nil
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

// SetAllowedResourceTypes installs the ARM resource-type scope derived from the
// selected scanner keys. An empty set intentionally excludes every resource type,
// matching the fail-closed behavior of the reference implementation after scanner loading.
func (f *AssessmentFilter) SetAllowedResourceTypes(resourceTypes []string) {
	f.allowedResourceTypes = stringSet(resourceTypes)
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
	return !f.allowedResourceTypes[normalize(resourceType)]
}

func (f *AssessmentFilter) IsRecommendationExcluded(recommendationID string) bool {
	return f.excludeRecommendations[normalize(recommendationID)]
}

// IsResourceExcluded evaluates structural scope and tags during resource discovery.
// Structural values are derived from the ARM resource ID so inventory and downstream
// findings use the same source-of-truth semantics.
func (f *AssessmentFilter) IsResourceExcluded(resourceID string, tags map[string]string) bool {
	if f.isResourceStructurallyExcluded(resourceID) {
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

// IsServiceExcluded preserves the reference downstream finding filter name and behavior.
// It reapplies structural scanner/subscription/RG/resource scope before tag-scope inheritance.
func (f *AssessmentFilter) IsServiceExcluded(resourceID string) bool {
	if f.isResourceStructurallyExcluded(resourceID) {
		return true
	}
	if len(f.includeTags) == 0 && len(f.excludeTags) == 0 {
		return false
	}
	if included, ok := f.resourceScopeDecision(resourceID); ok {
		return !included
	}
	// Unknown scope is fail-closed for include-tag filters and fail-open for
	// exclude-only tag filters, matching the reference implementation.
	return len(f.includeTags) > 0
}

func (f *AssessmentFilter) isResourceStructurallyExcluded(resourceID string) bool {
	resourceType := azure.ResourceTypeFromResourceID(resourceID)
	if f.IsResourceTypeExcluded(resourceType) {
		return true
	}
	if f.IsSubscriptionExcluded(azure.SubscriptionFromResourceID(resourceID)) {
		return true
	}
	if f.IsResourceGroupExcluded(azure.ResourceGroupIDFromResourceID(resourceID)) {
		return true
	}
	return f.excludeResources[normalize(resourceID)]
}

func (f *AssessmentFilter) resourceScopeDecision(resourceID string) (bool, bool) {
	id := normalize(resourceID)
	for {
		if included, ok := f.resourceScope[id]; ok {
			return included, true
		}
		lastSlash := strings.LastIndexByte(id, '/')
		if lastSlash <= 0 {
			return false, false
		}
		id = id[:lastSlash]
	}
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
