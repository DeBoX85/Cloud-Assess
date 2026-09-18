package json

import (
	stdjson "encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

type Options struct {
	RedactSubscriptionIDs bool
}

// Marshal returns the canonical assessment result as stable, human-readable JSON.
// The legacy-compatible default is unredacted for direct library callers; application
// surfaces should use MarshalWithOptions so CLI redaction preferences are explicit.
func Marshal(data *result.AssessmentResult) ([]byte, error) {
	return MarshalWithOptions(data, Options{})
}

// MarshalWithOptions returns canonical JSON while optionally redacting subscription IDs
// everywhere they occur in the serialized result. This includes nested resource IDs,
// policy identifiers, portal links, warning text, and other string fields containing a
// subscription ID, rather than only the top-level subscriptionId fields.
func MarshalWithOptions(data *result.AssessmentResult, opts Options) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("assessment result is nil")
	}
	encoded, err := stdjson.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}
	if !opts.RedactSubscriptionIDs {
		return encoded, nil
	}
	return redactSubscriptionIDs(encoded, data), nil
}

// String returns the same unredacted JSON representation used by Marshal.
func String(data *result.AssessmentResult) (string, error) {
	return StringWithOptions(data, Options{})
}

// StringWithOptions returns the same JSON representation used by file output.
func StringWithOptions(data *result.AssessmentResult, opts Options) (string, error) {
	encoded, err := MarshalWithOptions(data, opts)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// WriteFile writes the unredacted canonical JSON report to filename.
func WriteFile(data *result.AssessmentResult, filename string) error {
	return WriteFileWithOptions(data, filename, Options{})
}

// WriteFileWithOptions writes the canonical JSON report using the requested redaction behavior.
func WriteFileWithOptions(data *result.AssessmentResult, filename string, opts Options) error {
	if filename == "" {
		return fmt.Errorf("JSON output filename is empty")
	}
	encoded, err := MarshalWithOptions(data, opts)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, encoded, 0o600); err != nil {
		return fmt.Errorf("write JSON report %q: %w", filename, err)
	}
	return nil
}

func redactSubscriptionIDs(encoded []byte, data *result.AssessmentResult) []byte {
	ids := collectSubscriptionIDs(data)
	ids = append(ids, subscriptionIDsFromSerializedJSON(encoded)...)
	if len(ids) == 0 {
		return encoded
	}

	replacements := make(map[string]string, len(ids))
	patterns := make([]string, 0, len(ids))
	for _, id := range ids {
		masked := redact.SubscriptionID(id, true)
		if masked == "" {
			continue
		}
		key := strings.ToLower(id)
		if _, exists := replacements[key]; exists {
			continue
		}
		replacements[key] = masked
		patterns = append(patterns, regexp.QuoteMeta(id))
	}
	if len(patterns) == 0 {
		return encoded
	}
	sort.Strings(patterns)
	matcher, err := regexp.Compile("(?i)(" + strings.Join(patterns, "|") + ")")
	if err != nil {
		// Patterns are generated through regexp.QuoteMeta, so this is defensive only.
		return encoded
	}
	return matcher.ReplaceAllFunc(encoded, func(match []byte) []byte {
		return []byte(replacements[strings.ToLower(string(match))])
	})
}

func collectSubscriptionIDs(data *result.AssessmentResult) []string {
	ids := make([]string, 0)
	add := func(value string) {
		if strings.TrimSpace(value) != "" {
			ids = append(ids, value)
		}
	}
	for _, item := range data.Findings {
		add(item.SubscriptionID)
	}
	for _, item := range data.Resources {
		add(item.SubscriptionID)
	}
	for _, item := range data.OutOfScope {
		add(item.SubscriptionID)
	}
	for _, item := range data.ResourceTypes {
		add(item.SubscriptionID)
	}
	for _, item := range data.Advisor {
		add(item.SubscriptionID)
	}
	for _, item := range data.Defender {
		add(item.SubscriptionID)
	}
	for _, item := range data.DefenderRecommendations {
		add(item.SubscriptionID)
	}
	for _, item := range data.AzurePolicy {
		add(item.SubscriptionID)
	}
	for _, item := range data.ArcSQL {
		add(item.SubscriptionID)
	}
	for _, item := range data.Costs {
		add(item.SubscriptionID)
	}
	return ids
}


const subscriptionGUIDPattern = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`

var (
	subscriptionFieldPattern = regexp.MustCompile(`(?i)"subscriptionId"\s*:\s*"(` + subscriptionGUIDPattern + `)"`)
	subscriptionPathPattern  = regexp.MustCompile(`(?i)/subscriptions/(` + subscriptionGUIDPattern + `)`)
	subscriptionTextPattern  = regexp.MustCompile(`(?i)\bsubscription(?:\s+id)?\s+(` + subscriptionGUIDPattern + `)\b`)
)

func subscriptionIDsFromSerializedJSON(encoded []byte) []string {
	seen := map[string]struct{}{}
	var ids []string
	addMatches := func(pattern *regexp.Regexp) {
		for _, match := range pattern.FindAllSubmatch(encoded, -1) {
			if len(match) < 2 {
				continue
			}
			id := string(match[1])
			key := strings.ToLower(id)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			ids = append(ids, id)
		}
	}
	addMatches(subscriptionFieldPattern)
	addMatches(subscriptionPathPattern)
	addMatches(subscriptionTextPattern)
	return ids
}
