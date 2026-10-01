// Portions of this file reproduce Advisor metadata behavior from Microsoft Azure Quick
// Review (MIT licensed). See NOTICE.md.

package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const metadataAPIVersion = "2020-01-01"

type httpGetter interface {
	Get(context.Context, string) ([]byte, error)
}

type RecommendationTypeProvider interface {
	RecommendationTypes(context.Context) (map[string]string, error)
}

type MetadataClient struct {
	client   httpGetter
	endpoint string
}

func NewMetadataClient(client httpGetter, endpoint string) *MetadataClient {
	return &MetadataClient{client: client, endpoint: strings.TrimSuffix(endpoint, "/")}
}

// RecommendationTypes reproduces the source SDK's paged ARM request to
// /providers/Microsoft.Advisor/metadata?api-version=2020-01-01 and extracts the
// recommendationType metadata entity into recommendation ID -> display name.
func (c *MetadataClient) RecommendationTypes(ctx context.Context) (map[string]string, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("Advisor metadata HTTP client is not configured")
	}
	if c.endpoint == "" {
		return nil, fmt.Errorf("Azure Resource Manager endpoint is not configured")
	}

	nextURL := c.endpoint + "/providers/Microsoft.Advisor/metadata?api-version=" + metadataAPIVersion
	recommendationTypes := map[string]string{}
	seen := map[string]bool{}
	for nextURL != "" {
		var err error
		nextURL, err = c.resolveNextLink(nextURL)
		if err != nil {
			return nil, err
		}
		if seen[nextURL] {
			return nil, fmt.Errorf("Advisor metadata pagination repeated a continuation URL")
		}
		seen[nextURL] = true
		body, err := c.client.Get(ctx, nextURL)
		if err != nil {
			return nil, fmt.Errorf("list Advisor recommendation metadata: %w", err)
		}

		var page metadataListResult
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("decode Advisor recommendation metadata: %w", err)
		}
		for _, entity := range page.Value {
			if !strings.EqualFold(entity.Name, "recommendationType") {
				continue
			}
			for _, value := range entity.Properties.SupportedValues {
				if strings.TrimSpace(value.ID) == "" {
					continue
				}
				recommendationTypes[value.ID] = value.DisplayName
			}
		}
		nextURL = strings.TrimSpace(page.NextLink)
	}
	return recommendationTypes, nil
}

// Pagination may not move the ARM credential to another origin. Error messages
// deliberately omit the supplied URL, which can contain credentials or signed query data.
func (c *MetadataClient) resolveNextLink(nextLink string) (string, error) {
	base, err := url.Parse(c.endpoint)
	if err != nil || base.Scheme != "https" || base.Host == "" || base.User != nil || base.RawQuery != "" || base.Fragment != "" {
		return "", fmt.Errorf("Advisor metadata requires a valid HTTPS ARM endpoint")
	}
	nextLink = strings.TrimSpace(nextLink)
	link, err := url.Parse(nextLink)
	if err != nil || link.User != nil || link.Fragment != "" || strings.HasPrefix(nextLink, "//") || strings.Contains(nextLink, `\`) {
		return "", fmt.Errorf("Advisor metadata continuation URL is invalid")
	}
	if !link.IsAbs() {
		if link.Host != "" {
			return "", fmt.Errorf("Advisor metadata continuation URL is invalid")
		}
		link, err = url.Parse(c.endpoint + "/" + strings.TrimPrefix(nextLink, "/"))
		if err != nil {
			return "", fmt.Errorf("Advisor metadata continuation URL is invalid")
		}
	}
	if !strings.EqualFold(link.Scheme, base.Scheme) || !strings.EqualFold(link.Host, base.Host) {
		return "", fmt.Errorf("Advisor metadata continuation URL must remain on the configured ARM origin")
	}
	link.Scheme = strings.ToLower(link.Scheme)
	link.Host = strings.ToLower(link.Host)
	return link.String(), nil
}

type metadataListResult struct {
	Value    []metadataEntity `json:"value"`
	NextLink string           `json:"nextLink"`
}

type metadataEntity struct {
	Name       string                   `json:"name"`
	Properties metadataEntityProperties `json:"properties"`
}

type metadataEntityProperties struct {
	SupportedValues []metadataSupportedValue `json:"supportedValues"`
}

type metadataSupportedValue struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// Compile-time assertion that the shared Azure client supplies the required GET contract.
var _ httpGetter = (*azure.HTTPClient)(nil)
