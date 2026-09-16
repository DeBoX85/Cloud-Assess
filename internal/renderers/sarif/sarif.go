package sarif

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/findings"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

const (
	schemaURL              = "https://json.schemastore.org/sarif-2.1.0.json"
	fingerprintNamespace   = "cloudAssessFinding/v1"
	automationNamespace    = "cloud-assess/"
	logicalLocationKind    = "resource"
)

type logFile struct {
	Version string `json:"version"`
	Schema  string `json:"$schema"`
	Runs    []run  `json:"runs"`
}

type run struct {
	Tool       tool       `json:"tool"`
	Results    []sarifResult `json:"results"`
	Automation automation `json:"automationDetails,omitempty"`
}

type automation struct {
	ID string `json:"id,omitempty"`
}

type tool struct {
	Driver driver `json:"driver"`
}

type driver struct {
	Name           string `json:"name"`
	Version        string `json:"version,omitempty"`
	InformationURI string `json:"informationUri,omitempty"`
	Rules          []rule `json:"rules"`
}

type rule struct {
	ID               string            `json:"id"`
	ShortDescription message           `json:"shortDescription"`
	FullDescription  message           `json:"fullDescription,omitempty"`
	HelpURI          string            `json:"helpUri,omitempty"`
	Properties       map[string]string `json:"properties"`
}

type sarifResult struct {
	RuleID              string            `json:"ruleId"`
	Level               string            `json:"level"`
	Message             message           `json:"message"`
	PartialFingerprints map[string]string `json:"partialFingerprints"`
	Locations           []location        `json:"locations"`
}

type message struct {
	Text string `json:"text"`
}

type location struct {
	LogicalLocations []logicalLocation `json:"logicalLocations"`
}

type logicalLocation struct {
	Name               string `json:"name"`
	FullyQualifiedName string `json:"fullyQualifiedName"`
	Kind               string `json:"kind"`
}

// Marshal renders primary assessment findings as SARIF 2.1.0.
func Marshal(data *result.AssessmentResult, version string) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("assessment result is nil")
	}
	if data.Summary == nil {
		return nil, fmt.Errorf("SARIF rendering requires a findings summary")
	}

	impacted := make(map[string]findings.RecommendationSummary, data.Summary.RecommendationsFound)
	rules := make([]rule, 0, data.Summary.RecommendationsFound)
	for _, recommendation := range data.Summary.Recommendations {
		if recommendation.ImpactedResources == 0 {
			continue
		}
		key := normalize(recommendation.ID)
		if key == "" {
			continue
		}
		impacted[key] = recommendation
		rules = append(rules, buildRule(recommendation))
	}

	type findingKey struct {
		ruleID     string
		resourceID string
	}
	seen := make(map[findingKey]struct{}, data.Summary.ImpactedResources)
	results := make([]sarifResult, 0, data.Summary.ImpactedResources)
	for _, finding := range data.Findings {
		if strings.EqualFold(finding.Category, assessment.CategorySLA) {
			continue
		}
		ruleID := normalize(finding.RecommendationID)
		if _, ok := impacted[ruleID]; !ok {
			continue
		}
		key := findingKey{ruleID: ruleID, resourceID: normalize(finding.ResourceID)}
		if _, duplicate := seen[key]; duplicate {
			continue
		}
		seen[key] = struct{}{}
		results = append(results, buildResult(finding))
	}

	brand := branding.Default()
	output := logFile{
		Version: "2.1.0",
		Schema:  schemaURL,
		Runs: []run{{
			Tool: tool{Driver: driver{
				Name:           brand.CLIName,
				Version:        version,
				InformationURI: brand.WebsiteURL,
				Rules:          rules,
			}},
			Results:    results,
			Automation: automation{ID: automationNamespace + data.ScopeID},
		}},
	}

	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal SARIF report: %w", err)
	}
	return encoded, nil
}

// WriteFile writes the SARIF report with private file permissions.
func WriteFile(data *result.AssessmentResult, filename, version string) error {
	if filename == "" {
		return fmt.Errorf("SARIF output filename is empty")
	}
	encoded, err := Marshal(data, version)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, encoded, 0o600); err != nil {
		return fmt.Errorf("write SARIF report %q: %w", filename, err)
	}
	return nil
}

func buildRule(recommendation findings.RecommendationSummary) rule {
	return rule{
		ID:               recommendation.ID,
		ShortDescription: message{Text: recommendation.Recommendation},
		FullDescription:  message{Text: recommendation.LongDescription},
		HelpURI:          recommendation.LearnURL,
		Properties: map[string]string{
			"category":     recommendation.Category,
			"impact":       recommendation.Impact,
			"resourceType": recommendation.ResourceType,
			"source":       recommendation.Source,
		},
	}
}

func buildResult(finding assessment.Finding) sarifResult {
	return sarifResult{
		RuleID: finding.RecommendationID,
		Level:  level(finding.Impact),
		Message: message{Text: fmt.Sprintf(
			"%s: %s",
			finding.Recommendation,
			finding.ResourceID,
		)},
		PartialFingerprints: map[string]string{
			fingerprintNamespace: fingerprint(finding.RecommendationID, finding.ResourceID),
		},
		Locations: []location{{
			LogicalLocations: []logicalLocation{{
				Name:               finding.ResourceName,
				FullyQualifiedName: finding.ResourceID,
				Kind:               logicalLocationKind,
			}},
		}},
	}
}

func level(impact string) string {
	rank, _ := assessment.SeverityRank(impact)
	switch rank {
	case 3:
		return "error"
	case 2:
		return "warning"
	default:
		return "note"
	}
}

func fingerprint(recommendationID, resourceID string) string {
	sum := sha256.Sum256([]byte(normalize(recommendationID) + "\x00" + normalize(resourceID)))
	return hex.EncodeToString(sum[:16])
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
