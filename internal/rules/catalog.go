// Portions of this package reproduce rule-loading behavior from Microsoft Azure Quick Review (MIT licensed).
// See NOTICE.md for attribution details.

package rules

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"gopkg.in/yaml.v3"
)

const (
	SourceAPRL        = "APRL"
	SourceAOR         = "AOR"
	SourceCustom      = "CUSTOM"
	SourceDiagnostics = "DIAGNOSTICS"

	ValidationARG = "Azure Resource Graph"
)

type rawRecommendation struct {
	ID                   string        `yaml:"aprlGuid"`
	RecommendationTypeID any           `yaml:"recommendationTypeId"`
	Description          string        `yaml:"description"`
	Category             string        `yaml:"recommendationControl"`
	Impact               string        `yaml:"recommendationImpact"`
	ResourceType         string        `yaml:"recommendationResourceType"`
	MetadataState        string        `yaml:"recommendationMetadataState"`
	LongDescription      string        `yaml:"longDescription"`
	PotentialBenefits    string        `yaml:"potentialBenefits"`
	PGVerified           bool          `yaml:"pgVerified"`
	AutomationAvailable  any           `yaml:"automationAvailable"`
	Tags                 []string      `yaml:"tags,omitempty"`
	LearnMore            []rawLearnMore `yaml:"learnMoreLink,flow"`
}

type rawLearnMore struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Catalog contains normalized recommendations grouped by lower-case resource type and recommendation ID.
type Catalog struct {
	byResourceType map[string]map[string]assessment.RecommendationDefinition
}

func NewCatalog() *Catalog {
	return &Catalog{byResourceType: map[string]map[string]assessment.RecommendationDefinition{}}
}

// Add inserts or replaces a recommendation. Later providers therefore have higher precedence.
func (c *Catalog) Add(definition assessment.RecommendationDefinition) {
	resourceType := normalize(definition.ResourceType)
	if resourceType == "" || strings.TrimSpace(definition.ID) == "" {
		return
	}
	if c.byResourceType[resourceType] == nil {
		c.byResourceType[resourceType] = map[string]assessment.RecommendationDefinition{}
	}
	c.byResourceType[resourceType][definition.ID] = definition
}

func (c *Catalog) AddAll(definitions []assessment.RecommendationDefinition) {
	for _, definition := range definitions {
		c.Add(definition)
	}
}

// ByResourceType returns recommendations in deterministic ID order.
func (c *Catalog) ByResourceType(resourceType string) []assessment.RecommendationDefinition {
	byID := c.byResourceType[normalize(resourceType)]
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]assessment.RecommendationDefinition, 0, len(ids))
	for _, id := range ids {
		out = append(out, byID[id])
	}
	return out
}

// All returns every recommendation in deterministic resource-type then ID order.
func (c *Catalog) All() []assessment.RecommendationDefinition {
	resourceTypes := make([]string, 0, len(c.byResourceType))
	for resourceType := range c.byResourceType {
		resourceTypes = append(resourceTypes, resourceType)
	}
	sort.Strings(resourceTypes)

	var out []assessment.RecommendationDefinition
	for _, resourceType := range resourceTypes {
		out = append(out, c.ByResourceType(resourceType)...)
	}
	return out
}

// LoadFilesystem reads recommendation YAML and matching KQL files from an fs.FS.
// A KQL file named <recommendation-id>.kql is attached to the recommendation with the same ID.
func LoadFilesystem(fsys fs.FS, source string) ([]assessment.RecommendationDefinition, error) {
	queries := map[string]string{}
	if err := fs.WalkDir(fsys, ".", func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.EqualFold(path.Ext(filePath), ".kql") {
			return nil
		}
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return err
		}
		id := strings.TrimSuffix(path.Base(filePath), path.Ext(filePath))
		queries[id] = string(content)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("load %s queries: %w", source, err)
	}

	var definitions []assessment.RecommendationDefinition
	if err := fs.WalkDir(fsys, ".", func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.EqualFold(path.Ext(filePath), ".yaml") {
			return nil
		}
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return err
		}
		var raw []rawRecommendation
		if err := yaml.Unmarshal(content, &raw); err != nil {
			return fmt.Errorf("parse %s: %w", filePath, err)
		}
		for _, recommendation := range raw {
			definition := normalizeRecommendation(recommendation, source)
			if query, ok := queries[definition.ID]; ok {
				definition.Query = query
			}
			definitions = append(definitions, definition)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("load %s recommendations: %w", source, err)
	}

	return definitions, nil
}

// IsExecutableEmbedded reproduces the source filtering applied to embedded recommendation rules.
func IsExecutableEmbedded(definition assessment.RecommendationDefinition, isExcluded func(string) bool) bool {
	if isExcluded != nil && isExcluded(definition.ID) {
		return false
	}
	query := strings.ToLower(definition.Query)
	if strings.Contains(query, "cannot-be-validated-with-arg") ||
		strings.Contains(query, "under-development") ||
		strings.Contains(query, "under development") ||
		strings.EqualFold(definition.State, "disabled") {
		return false
	}
	return true
}

// IsExecutablePlugin reproduces the source behavior for external YAML rules, which only applies recommendation exclusion.
func IsExecutablePlugin(definition assessment.RecommendationDefinition, isExcluded func(string) bool) bool {
	return isExcluded == nil || !isExcluded(definition.ID)
}

func normalizeRecommendation(raw rawRecommendation, source string) assessment.RecommendationDefinition {
	links := make([]assessment.LearnMoreLink, 0, len(raw.LearnMore))
	for _, link := range raw.LearnMore {
		links = append(links, assessment.LearnMoreLink{Name: link.Name, URL: link.URL})
	}
	return assessment.RecommendationDefinition{
		ID:                   raw.ID,
		RecommendationTypeID: scalarString(raw.RecommendationTypeID),
		Recommendation:       raw.Description,
		Category:             raw.Category,
		Impact:               raw.Impact,
		ResourceType:         raw.ResourceType,
		State:                raw.MetadataState,
		LongDescription:      raw.LongDescription,
		PotentialBenefits:    raw.PotentialBenefits,
		PGVerified:           raw.PGVerified,
		AutomationAvailable:  scalarString(raw.AutomationAvailable),
		Tags:                 append([]string(nil), raw.Tags...),
		LearnMore:            links,
		Source:               source,
		ValidationMechanism:  ValidationARG,
	}
}

func scalarString(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
