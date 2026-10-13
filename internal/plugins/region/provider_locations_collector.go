// Provider request semantics derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_PROVIDER_TARGET_RECONCILIATION_20261012.md.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const MaxProviderLocationWork = 65536

type ProviderLocationRequest struct {
	SubscriptionID string
}

// Locations and all run state are owned by this collection. Status is complete,
// partial or unknown; even an admitted empty page establishes a partial prefix.
type ProviderLocationCollection struct {
	SubscriptionID string
	Status         string
	Locations      map[string]map[string]bool
	FailureCode    string
}

type ProviderLocationCollector struct {
	transport *RESTQuotaCollector
}

func NewProviderLocationCollector(endpoint string, getter RESTQuotaGetter) (*ProviderLocationCollector, error) {
	transport, err := NewRESTQuotaCollector(endpoint, getter)
	if err != nil {
		return nil, fmt.Errorf("invalid provider location origin or getter")
	}
	return &ProviderLocationCollector{transport: transport}, nil
}

var providerIdentifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func providerTypeIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for _, part := range strings.Split(value, "/") {
		if !providerIdentifier.MatchString(part) {
			return false
		}
	}
	return true
}

func providerReturnedID(id, subscription, namespace string) bool {
	parts := strings.Split(id, "/")
	return len(parts) == 5 && parts[0] == "" &&
		len(parts[1]) == len("subscriptions") && strings.EqualFold(parts[1], "subscriptions") &&
		len(parts[2]) == 36 && subscriptionID.MatchString(parts[2]) && strings.EqualFold(parts[2], subscription) &&
		len(parts[3]) == len("providers") && strings.EqualFold(parts[3], "providers") &&
		providerIdentifier.MatchString(parts[4]) && strings.EqualFold(parts[4], namespace)
}

func providerLocation(value string) (string, bool) {
	// Preserve the source's ASCII-space-only display normalization, with
	// explicit ASCII admission before lowercasing to reject Unicode aliases.
	for _, r := range value {
		if r > 127 {
			return "", false
		}
	}
	normalized := strings.ToLower(strings.ReplaceAll(value, " ", ""))
	return normalized, len(normalized) <= 64 && regionID.MatchString(normalized) &&
		(normalized == "global" || costPhysical(normalized))
}

type providerPage struct {
	locations map[string]map[string]bool
	next      string
	work      int
	text      int
}

func providerArray(raw json.RawMessage) ([]json.RawMessage, bool) {
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, false
	}
	var values []json.RawMessage
	if json.Unmarshal(raw, &values) != nil {
		return nil, false
	}
	return values, true
}

func decodeProviderPage(ctx context.Context, body []byte, subscription string, prior map[string]map[string]bool, work, text int) (providerPage, string) {
	page := providerPage{locations: map[string]map[string]bool{}, work: work, text: text}
	if !quotaJSON(ctx, body) {
		return providerPage{}, "provider_invalid_response"
	}
	var envelope struct {
		Value    json.RawMessage
		NextLink json.RawMessage
	}
	if json.Unmarshal(body, &envelope) != nil {
		return providerPage{}, "provider_invalid_response"
	}
	providers, ok := providerArray(envelope.Value)
	if !ok {
		return providerPage{}, "provider_invalid_response"
	}
	if len(envelope.NextLink) != 0 {
		// A present null is not a successful end-of-pagination declaration.
		if bytes.Equal(bytes.TrimSpace(envelope.NextLink), []byte("null")) ||
			json.Unmarshal(envelope.NextLink, &page.next) != nil {
			return providerPage{}, "provider_invalid_response"
		}
	}
	charge := func() bool {
		if page.work >= MaxProviderLocationWork {
			return false
		}
		page.work++
		return true
	}
	label := func(value string) bool { return auxiliaryText(value, &page.text) }
	for _, raw := range providers {
		if ctx.Err() != nil {
			return providerPage{}, "provider_cancelled"
		}
		if !charge() {
			return providerPage{}, "provider_work_limit"
		}
		var provider struct {
			ID            json.RawMessage
			Namespace     string
			ResourceTypes json.RawMessage
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &provider) != nil ||
			!label(provider.Namespace) || !providerIdentifier.MatchString(provider.Namespace) {
			return providerPage{}, "provider_identity_invalid"
		}
		if len(provider.ID) != 0 {
			var id string
			if json.Unmarshal(provider.ID, &id) != nil || !label(id) ||
				!providerReturnedID(id, subscription, provider.Namespace) {
				return providerPage{}, "provider_identity_invalid"
			}
		}
		types, ok := providerArray(provider.ResourceTypes)
		if !ok {
			return providerPage{}, "provider_invalid_response"
		}
		for _, rawType := range types {
			if ctx.Err() != nil {
				return providerPage{}, "provider_cancelled"
			}
			if !charge() {
				return providerPage{}, "provider_work_limit"
			}
			var resourceType struct {
				ResourceType string
				Locations    json.RawMessage
			}
			if bytes.Equal(bytes.TrimSpace(rawType), []byte("null")) || json.Unmarshal(rawType, &resourceType) != nil ||
				!label(resourceType.ResourceType) || !providerTypeIdentifier(resourceType.ResourceType) {
				return providerPage{}, "provider_identity_invalid"
			}
			key := strings.ToLower(provider.Namespace + "/" + resourceType.ResourceType)
			if len(key) > MaxLabelBytes {
				return providerPage{}, "provider_text_limit"
			}
			if _, exists := page.locations[key]; exists {
				return providerPage{}, "provider_type_collision"
			}
			if _, exists := prior[key]; exists {
				return providerPage{}, "provider_type_collision"
			}
			locations, ok := providerArray(resourceType.Locations)
			if !ok {
				return providerPage{}, "provider_invalid_response"
			}
			values := map[string]bool{}
			for _, rawLocation := range locations {
				if ctx.Err() != nil {
					return providerPage{}, "provider_cancelled"
				}
				if !charge() {
					return providerPage{}, "provider_work_limit"
				}
				var value string
				if bytes.Equal(bytes.TrimSpace(rawLocation), []byte("null")) || json.Unmarshal(rawLocation, &value) != nil || !label(value) {
					return providerPage{}, "provider_invalid_response"
				}
				normalized, ok := providerLocation(value)
				if !ok {
					return providerPage{}, "provider_location_invalid"
				}
				values[normalized] = true
			}
			page.locations[key] = values
		}
	}
	return page, ""
}

// Collect validates every destination before credentials and admits entire
// pages atomically. Stable caller inputs are required; no run state is shared.
func (c *ProviderLocationCollector) Collect(ctx context.Context, subscriptions map[string]string, request ProviderLocationRequest) (*ProviderLocationCollection, error) {
	failInput := func(code string) (*ProviderLocationCollection, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("provider location request rejected [%s]", code)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &text)
	if !ok {
		return failInput("scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return failInput("scope_invalid")
		}
	}
	id := strings.ToLower(request.SubscriptionID)
	if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || scope[id] == "" {
		return failInput("scope_invalid")
	}
	first := *c.transport.origin
	first.Path = "/subscriptions/" + id + "/providers"
	first.RawQuery = "api-version=2021-04-01"
	result := &ProviderLocationCollection{SubscriptionID: id, Status: "complete", Locations: map[string]map[string]bool{}}
	admitted, work, total := 0, 0, 0
	seen := map[string]bool{}
	link := first.String()
	fail := func(code string) (*ProviderLocationCollection, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result.Status = "unknown"
		if admitted > 0 {
			result.Status = "partial"
		}
		result.FailureCode = code
		return result, nil
	}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if admitted >= MaxQuotaPages {
			return fail("provider_page_limit")
		}
		u, safe := c.transport.continuation(&first, "2021-04-01", link)
		if !safe {
			return fail("provider_unsafe_continuation")
		}
		if seen[u.String()] {
			return fail("provider_continuation_cycle")
		}
		seen[u.String()] = true
		body, response, err := c.transport.getter.GetBoundedWithResponse(ctx, u.String(), MaxQuotaPageBytes)
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if err != nil || response == nil || response.StatusCode != http.StatusOK {
			return fail("provider_request_failed")
		}
		if len(body) > MaxQuotaPageBytes || len(body) > MaxQuotaTotalBytes-total {
			return fail("provider_byte_limit")
		}
		total += len(body)
		page, code := decodeProviderPage(ctx, body, id, result.Locations, work, text)
		if code != "" {
			return fail(code)
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		for key, locations := range page.locations {
			result.Locations[key] = locations
		}
		work, text = page.work, page.text
		admitted++
		if page.next == "" {
			break
		}
		link = page.next
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
