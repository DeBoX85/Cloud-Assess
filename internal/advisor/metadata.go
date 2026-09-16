// Portions of this file reproduce Advisor metadata behavior from Microsoft Azure Quick
// Review (MIT licensed). See NOTICE.md.

package advisor

import (
	"context"
	"encoding/json"
	"fmt"
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
	for nextURL != "" {
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
		nextURL = c.resolveNextLink(page.NextLink)
	}
	return recommendationTypes, nil
}

func (c *MetadataClient) resolveNextLink(nextLink string) string {
	nextLink = strings.TrimSpace(nextLink)
	if nextLink == "" {
		return ""
	}
	if strings.HasPrefix(nextLink, "https://") || strings.HasPrefix(nextLink, "http://") {
		return nextLink
	}
	if strings.HasPrefix(nextLink, "/") {
		return c.endpoint + nextLink
	}
	return c.endpoint + "/" + nextLink
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
