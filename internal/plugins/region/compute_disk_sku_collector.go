// Compute disk SKU request semantics derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_COMPUTE_DISK_SKU_COLLECTOR_20261013.md.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	MaxComputeDiskSKUWork = 65536
	computeDiskSKUVersion = "2021-07-01"
)

type ComputeDiskSKURequest struct {
	SubscriptionID string
	Region         string
}

// Values retain only admitted pages. Incomplete diagnostic values must never
// be used directly for scoring; Evidence projects them to unknown/empty.
type ComputeDiskSKUCollection struct {
	SubscriptionID string
	Region         string
	Status         string
	Values         map[string]SKUAvailability
	FailureCode    string
}

func (c *ComputeDiskSKUCollection) Evidence() SKUEvidence {
	out := SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}}
	if c == nil || c.Status != "complete" {
		return out
	}
	out.Status = "complete"
	for key, value := range c.Values {
		value.BlockedZones = append([]string{}, value.BlockedZones...)
		out.Values[key] = value
	}
	return out
}

type ComputeDiskSKUCollector struct{ transport *RESTQuotaCollector }

func NewComputeDiskSKUCollector(endpoint string, getter RESTQuotaGetter) (*ComputeDiskSKUCollector, error) {
	transport, err := NewRESTQuotaCollector(endpoint, getter)
	if err != nil {
		return nil, fmt.Errorf("invalid Compute disk SKU origin or getter")
	}
	return &ComputeDiskSKUCollector{transport: transport}, nil
}

func diskPhysical(value string) (string, bool) {
	if len(value) == 0 || len(value) > 64 {
		return "", false
	}
	for _, character := range value {
		if character > 127 {
			return "", false
		}
	}
	value = strings.ToLower(value)
	return value, regionID.MatchString(value) && costPhysical(value)
}

func diskName(value string) bool {
	if value == "" || len(value) > MaxLabelBytes || strings.TrimSpace(value) != value {
		return false
	}
	for _, character := range value {
		if character < 32 || character > 126 {
			return false
		}
	}
	return true
}

func diskZone(value string) bool {
	if value == "" || len(value) > MaxLabelBytes || value[0] < '1' || value[0] > '9' {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func (c *ComputeDiskSKUCollector) continuation(first *url.URL, filter, link string) (*url.URL, bool) {
	if len(link) > 8192 || strings.ContainsAny(link, "#\\ \r\n\t") {
		return nil, false
	}
	ref, err := url.Parse(link)
	if err != nil {
		return nil, false
	}
	target := first.ResolveReference(ref)
	query, err := url.ParseQuery(target.RawQuery)
	if err != nil || target.Scheme != "https" || !strings.EqualFold(target.Host, c.transport.origin.Host) ||
		target.Path != first.Path || target.RawPath != "" || target.Opaque != "" || target.User != nil ||
		target.Fragment != "" || target.ForceQuery || len(query["api-version"]) != 1 || query.Get("api-version") != computeDiskSKUVersion {
		return nil, false
	}
	for key, values := range query {
		if len(values) != 1 {
			return nil, false
		}
		switch key {
		case "api-version":
		case "$filter":
			if values[0] != filter {
				return nil, false
			}
		case "$skiptoken":
			if values[0] == "" || len(values[0]) > 4096 {
				return nil, false
			}
		default:
			return nil, false
		}
	}
	target.Host = c.transport.origin.Host
	target.RawQuery = query.Encode()
	return target, true
}

type diskPageBudget struct {
	ctx        context.Context
	work, text int
}

func (b *diskPageBudget) charge() string {
	if b.ctx.Err() != nil {
		return "sku_cancelled"
	}
	if b.work >= MaxComputeDiskSKUWork {
		return "sku_work_limit"
	}
	b.work++
	return ""
}

func (b *diskPageBudget) label(value string, limit int) string {
	if b.ctx.Err() != nil {
		return "sku_cancelled"
	}
	if len(value) > limit || len(value) > auxTextBudget-b.text {
		return "sku_text_limit"
	}
	b.text += len(value)
	return ""
}

// Charge every object, array element and decoded string (including property
// names, ignored metadata and filtered rows) before projecting any page.
func (b *diskPageBudget) walk(value any, depth int, rootNext bool) string {
	if b.ctx.Err() != nil {
		return "sku_cancelled"
	}
	switch value := value.(type) {
	case map[string]any:
		if code := b.charge(); code != "" {
			return code
		}
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if code := b.label(key, MaxLabelBytes); code != "" {
				return code
			}
			if code := b.walk(value[key], depth+1, depth == 0 && key == "nextLink"); code != "" {
				return code
			}
		}
	case []any:
		for _, entry := range value {
			if code := b.charge(); code != "" {
				return code
			}
			if code := b.walk(entry, depth+1, false); code != "" {
				return code
			}
		}
	case string:
		limit := MaxLabelBytes
		if rootNext {
			limit = 8192
		}
		return b.label(value, limit)
	}
	return ""
}

func diskScalar(raw json.RawMessage, target any) bool {
	return len(raw) != 0 && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) && json.Unmarshal(raw, target) == nil
}

func diskStrings(ctx context.Context, raw json.RawMessage, valid func(string) bool) ([]string, bool) {
	array, ok := providerArray(raw)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(array))
	for _, item := range array {
		var value string
		if ctx.Err() != nil || !diskScalar(item, &value) || valid != nil && !valid(value) {
			return nil, false
		}
		out = append(out, value)
	}
	return out, true
}

func diskOptional(fields map[string]json.RawMessage, key string, valid func(json.RawMessage) bool) bool {
	raw, present := fields[key]
	return !present || valid(raw)
}

func diskOptionalStrings(ctx context.Context, fields map[string]json.RawMessage, keys ...string) bool {
	for _, key := range keys {
		if !diskOptional(fields, key, func(raw json.RawMessage) bool {
			var value string
			return diskScalar(raw, &value)
		}) {
			return false
		}
	}
	return ctx.Err() == nil
}

func diskCapabilities(ctx context.Context, raw json.RawMessage) bool {
	array, ok := providerArray(raw)
	if !ok {
		return false
	}
	for _, item := range array {
		fields, ok := providerObject(item, "name", "value")
		if !ok || !diskOptionalStrings(ctx, fields, "name", "value") {
			return false
		}
	}
	return ctx.Err() == nil
}

func diskMetadata(ctx context.Context, fields map[string]json.RawMessage) bool {
	if !diskOptionalStrings(ctx, fields, "tier", "size", "family", "kind") {
		return false
	}
	for _, key := range []string{"apiVersions"} {
		if !diskOptional(fields, key, func(raw json.RawMessage) bool {
			_, ok := diskStrings(ctx, raw, nil)
			return ok
		}) {
			return false
		}
	}
	if !diskOptional(fields, "capacity", func(raw json.RawMessage) bool {
		capacity, ok := providerObject(raw, "minimum", "maximum", "default", "scaleType")
		if !ok {
			return false
		}
		for _, key := range []string{"minimum", "maximum", "default"} {
			if !diskOptional(capacity, key, func(raw json.RawMessage) bool {
				var value int64
				return diskScalar(raw, &value)
			}) {
				return false
			}
		}
		return diskOptional(capacity, "scaleType", func(raw json.RawMessage) bool {
			var value string
			return diskScalar(raw, &value) && (value == "Automatic" || value == "Manual" || value == "None")
		})
	}) || !diskOptional(fields, "capabilities", func(raw json.RawMessage) bool { return diskCapabilities(ctx, raw) }) {
		return false
	}
	if !diskOptional(fields, "costs", func(raw json.RawMessage) bool {
		array, ok := providerArray(raw)
		if !ok {
			return false
		}
		for _, item := range array {
			cost, ok := providerObject(item, "meterID", "quantity", "extendedUnit")
			if !ok || !diskOptionalStrings(ctx, cost, "meterID", "extendedUnit") ||
				!diskOptional(cost, "quantity", func(raw json.RawMessage) bool {
					var quantity int64
					return diskScalar(raw, &quantity)
				}) {
				return false
			}
		}
		return ctx.Err() == nil
	}) {
		return false
	}
	return diskOptional(fields, "locationInfo", func(raw json.RawMessage) bool {
		array, ok := providerArray(raw)
		if !ok {
			return false
		}
		for _, item := range array {
			location, ok := providerObject(item, "location", "zones", "zoneDetails", "extendedLocations", "type")
			if !ok || !diskOptionalStrings(ctx, location, "type") ||
				!diskOptional(location, "location", func(raw json.RawMessage) bool {
					var value string
					if !diskScalar(raw, &value) {
						return false
					}
					_, ok := diskPhysical(value)
					return ok
				}) {
				return false
			}
			for _, key := range []string{"zones", "extendedLocations"} {
				var validate func(string) bool
				if key == "zones" {
					validate = diskZone
				}
				if !diskOptional(location, key, func(raw json.RawMessage) bool {
					_, ok := diskStrings(ctx, raw, validate)
					return ok
				}) {
					return false
				}
			}
			if !diskOptional(location, "zoneDetails", func(raw json.RawMessage) bool {
				details, ok := providerArray(raw)
				if !ok {
					return false
				}
				for _, item := range details {
					detail, ok := providerObject(item, "Name", "capabilities")
					if !ok || !diskOptional(detail, "Name", func(raw json.RawMessage) bool {
						_, ok := diskStrings(ctx, raw, diskZone)
						return ok
					}) || !diskOptional(detail, "capabilities", func(raw json.RawMessage) bool {
						return diskCapabilities(ctx, raw)
					}) {
						return false
					}
				}
				return ctx.Err() == nil
			}) {
				return false
			}
		}
		return ctx.Err() == nil
	})
}

func diskLocations(ctx context.Context, raw json.RawMessage) ([]string, bool) {
	values, ok := diskStrings(ctx, raw, func(value string) bool {
		_, ok := diskPhysical(value)
		return ok
	})
	for i := range values {
		values[i] = strings.ToLower(values[i])
	}
	return values, ok
}

func diskAvailability(ctx context.Context, row map[string]json.RawMessage, region string) (SKUAvailability, bool) {
	value := SKUAvailability{State: SKUAvailable, BlockedZones: []string{}}
	locations, ok := diskLocations(ctx, row["locations"])
	if !ok {
		return SKUAvailability{}, false
	}
	listed := false
	for _, location := range locations {
		listed = listed || location == region
	}
	restrictions, ok := providerArray(row["restrictions"])
	if !ok {
		return SKUAvailability{}, false
	}
	restricted := false
	blocked := map[string]bool{}
	for _, raw := range restrictions {
		restriction, ok := providerObject(raw, "type", "values", "restrictionInfo", "reasonCode")
		var kind, reason string
		if !ok || ctx.Err() != nil || !diskScalar(restriction["type"], &kind) ||
			!diskScalar(restriction["reasonCode"], &reason) || reason != "QuotaId" && reason != "NotAvailableForSubscription" ||
			kind != "Location" && kind != "Zone" {
			return SKUAvailability{}, false
		}
		var values, scope, zones []string
		if raw, present := restriction["values"]; present {
			if kind == "Location" {
				values, ok = diskLocations(ctx, raw)
			} else {
				values, ok = diskStrings(ctx, raw, nil)
			}
			if !ok {
				return SKUAvailability{}, false
			}
		}
		if raw, present := restriction["restrictionInfo"]; present {
			info, valid := providerObject(raw, "locations", "zones")
			if !valid {
				return SKUAvailability{}, false
			}
			if raw, present := info["locations"]; present {
				scope, ok = diskLocations(ctx, raw)
				if !ok {
					return SKUAvailability{}, false
				}
			}
			if raw, present := info["zones"]; present {
				zones, ok = diskStrings(ctx, raw, diskZone)
				if !ok {
					return SKUAvailability{}, false
				}
			}
		}
		if kind == "Location" {
			if len(values)+len(scope) == 0 {
				return SKUAvailability{}, false
			}
			for _, location := range append(values, scope...) {
				restricted = restricted || location == region
			}
		} else {
			if len(scope) == 0 || len(zones) == 0 {
				return SKUAvailability{}, false
			}
			for _, location := range scope {
				if location == region {
					for _, zone := range zones {
						blocked[zone] = true
					}
				}
			}
		}
	}
	if !listed {
		value.State = SKUUnavailable
	} else if restricted {
		value.State = SKURestricted
	} else if len(blocked) > 0 {
		value.State = SKUZoneRestricted
		for zone := range blocked {
			value.BlockedZones = append(value.BlockedZones, zone)
		}
		sort.Strings(value.BlockedZones)
	}
	return value, ctx.Err() == nil
}

type diskDecodedPage struct {
	values     map[string]SKUAvailability
	next       string
	work, text int
}

func decodeDiskPage(ctx context.Context, body []byte, region string, prior map[string]SKUAvailability, work, text int) (diskDecodedPage, string) {
	page := diskDecodedPage{values: map[string]SKUAvailability{}}
	if !quotaJSON(ctx, body) {
		return diskDecodedPage{}, "sku_invalid_response"
	}
	envelope, valid := providerObject(body, "value", "nextLink")
	if !valid {
		return diskDecodedPage{}, "sku_invalid_response"
	}
	rows, valid := providerArray(envelope["value"])
	if !valid {
		return diskDecodedPage{}, "sku_invalid_response"
	}
	if raw, present := envelope["nextLink"]; present && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) && !diskScalar(raw, &page.next) {
		return diskDecodedPage{}, "sku_invalid_response"
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var document any
	if decoder.Decode(&document) != nil {
		return diskDecodedPage{}, "sku_invalid_response"
	}
	budget := diskPageBudget{ctx: ctx, work: work, text: text}
	if code := budget.walk(document, 0, false); code != "" {
		return diskDecodedPage{}, code
	}
	for _, raw := range rows {
		row, valid := providerObject(raw, "resourceType", "name")
		var kind, name string
		if !valid || ctx.Err() != nil || !diskScalar(row["resourceType"], &kind) || !providerTypeIdentifier(kind) ||
			!diskScalar(row["name"], &name) || !diskName(name) {
			return diskDecodedPage{}, "sku_identity_invalid"
		}
		if strings.ToLower(kind) != "disks" {
			continue
		}
		row, valid = providerObject(raw, "resourceType", "name", "tier", "size", "family", "kind", "capacity",
			"locations", "locationInfo", "apiVersions", "costs", "capabilities", "restrictions")
		if !valid || !diskMetadata(ctx, row) {
			return diskDecodedPage{}, "sku_invalid_response"
		}
		key := strings.ToLower(name)
		if _, exists := page.values[key]; exists {
			return diskDecodedPage{}, "sku_identity_duplicate"
		}
		if _, exists := prior[key]; exists {
			return diskDecodedPage{}, "sku_identity_duplicate"
		}
		value, valid := diskAvailability(ctx, row, region)
		if !valid {
			return diskDecodedPage{}, "sku_invalid_response"
		}
		page.values[key] = value
	}
	page.work, page.text = budget.work, budget.text
	return page, ""
}

// Collect admits entire pages, owns all call state and validates every URL before
// the authenticated getter. Callers must keep scope inputs stable during the call.
func (c *ComputeDiskSKUCollector) Collect(ctx context.Context, subscriptions map[string]string, request ComputeDiskSKURequest) (*ComputeDiskSKUCollection, error) {
	failInput := func(code string) (*ComputeDiskSKUCollection, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Compute disk SKU request rejected [%s]", code)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.transport == nil {
		return failInput("request_invalid")
	}
	text := 0
	scope, valid := auxiliaryScope(ctx, subscriptions, &text)
	if !valid {
		return failInput("scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return failInput("scope_invalid")
		}
	}
	id := strings.ToLower(request.SubscriptionID)
	region, valid := diskPhysical(request.Region)
	if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || scope[id] == "" ||
		!valid || region != request.Region {
		return failInput("request_invalid")
	}
	first := *c.transport.origin
	first.Path = "/subscriptions/" + id + "/providers/Microsoft.Compute/skus"
	filter := "location eq '" + region + "'"
	query := url.Values{}
	query.Set("api-version", computeDiskSKUVersion)
	query.Set("$filter", filter)
	first.RawQuery = query.Encode()
	result := &ComputeDiskSKUCollection{SubscriptionID: id, Region: region, Status: "complete", Values: map[string]SKUAvailability{}}
	admitted, total, work := 0, 0, 0
	seen := map[string]bool{}
	link := first.String()
	fail := func(code string) (*ComputeDiskSKUCollection, error) {
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
			return fail("sku_page_limit")
		}
		target, valid := c.continuation(&first, filter, link)
		if !valid {
			return fail("sku_unsafe_continuation")
		}
		if seen[target.String()] {
			return fail("sku_continuation_cycle")
		}
		seen[target.String()] = true
		body, response, err := c.transport.getter.GetBoundedWithResponse(ctx, target.String(), MaxQuotaPageBytes)
		if cause := ctx.Err(); cause != nil {
			return nil, cause
		}
		if errors.Is(err, context.Canceled) {
			return nil, context.Canceled
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		if response != nil && response.StatusCode != http.StatusOK {
			return fail("sku_status_invalid")
		}
		if err != nil || response == nil {
			return fail("sku_request_failed")
		}
		if len(body) > MaxQuotaPageBytes || len(body) > MaxQuotaTotalBytes-total {
			return fail("sku_byte_limit")
		}
		total += len(body)
		page, code := decodeDiskPage(ctx, body, region, result.Values, work, text)
		if code != "" {
			return fail(code)
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		for key, value := range page.values {
			result.Values[key] = value
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
