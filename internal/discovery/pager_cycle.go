package discovery

import (
	"fmt"
	"net/url"
	"strings"
)

// Each pager owns its continuation set; independent listings must not interfere.
// Query/path values are opaque. Normalize only URL scheme/host case.
func checkScopeContinuation(next *string, seen map[string]struct{}) error {
	if next == nil || *next == "" {
		return nil
	}
	value := *next
	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("scope pager returned an invalid continuation")
	}
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	value = parsed.String()
	if _, exists := seen[value]; exists {
		return fmt.Errorf("scope pager repeated a continuation")
	}
	seen[value] = struct{}{}
	return nil
}
