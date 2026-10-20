// Storage request behavior derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_STORAGE_SKU_CONTRACT_20261016.md for corrections.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	MaxStorageSKUBytes = 1 << 20
	MaxStorageSKURows  = 8192
	MaxStorageSKUWork  = 65536
	storageSKUVersion  = "2023-01-01"
)

type StorageSKURequest struct {
	SubscriptionID  string
	Region          string
	LocationAliases map[string]string
}

type StorageSKUGroup struct {
	Known bool
	State int
}

// Groups are diagnostic owned values for Region, qualified by name and kind.
// Only Evidence may project them to the existing name-only scoring interface.
type StorageSKUCollection struct {
	SubscriptionID   string
	Region           string
	Status           string
	Groups           map[string]map[string]StorageSKUGroup
	UnqualifiedNames map[string]bool
	FailureCode      string
}

type StorageSKUCollector struct{ transport *RESTQuotaCollector }

func NewStorageSKUCollector(endpoint string, getter RESTQuotaGetter) (*StorageSKUCollector, error) {
	transport, err := NewRESTQuotaCollector(endpoint, getter)
	if err != nil {
		return nil, fmt.Errorf("invalid Storage SKU origin or getter")
	}
	return &StorageSKUCollector{transport: transport}, nil
}

func storageASCII(value string) bool {
	for _, r := range value {
		if r > 127 {
			return false
		}
	}
	return true
}

func storageIdentifier(value string) bool {
	return len(value) <= MaxLabelBytes && providerIdentifier.MatchString(value)
}

func storageStage(value string) bool {
	if !strings.HasSuffix(value, "(stage)") {
		return false
	}
	_, ok := diskPhysical(strings.TrimSuffix(value, "(stage)"))
	return ok
}

// known distinguishes an explicit nonphysical stage from an unresolved label.
func storageLocation(value string, aliases map[string]string) (location string, known bool) {
	if !storageASCII(value) {
		return "", false
	}
	value = strings.ToLower(value)
	if physical, ok := diskPhysical(value); ok {
		return physical, true
	}
	if storageStage(value) {
		return "", true
	}
	if physical, ok := aliases[value]; ok {
		return physical, true
	}
	return "", false
}

type storageBudget struct {
	ctx        context.Context
	work, text int
}

func (b *storageBudget) entry() bool {
	if b.ctx.Err() != nil || b.work >= MaxStorageSKUWork {
		return false
	}
	b.work++
	return true
}

func (b *storageBudget) label(value string) bool {
	return b.ctx.Err() == nil && auxiliaryText(value, &b.text)
}

// Charge every object and array element, including filtered/ignored metadata.
func (b *storageBudget) walk(value any) string {
	if b.ctx.Err() != nil {
		return "storage_sku_work_limit"
	}
	switch value := value.(type) {
	case map[string]any:
		if !b.entry() {
			return "storage_sku_work_limit"
		}
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if !b.label(key) {
				return "storage_sku_text_limit"
			}
			if code := b.walk(value[key]); code != "" {
				return code
			}
		}
	case []any:
		for _, element := range value {
			if !b.entry() {
				return "storage_sku_work_limit"
			}
			if code := b.walk(element); code != "" {
				return code
			}
		}
	case string:
		if !b.label(value) {
			return "storage_sku_text_limit"
		}
	}
	return ""
}

type storageFlags struct {
	listed, hard, restricted, unknown bool
}

// Optional availability metadata can be valid but insufficient, never global.
func storageRowFlags(ctx context.Context, row map[string]json.RawMessage, region string, aliases map[string]string) (storageFlags, bool) {
	flags := storageFlags{}
	if !diskOptionalStrings(ctx, row, "tier") {
		return flags, false
	}
	if raw, present := row["locations"]; present {
		locations, ok := diskStrings(ctx, raw, nil)
		if !ok {
			return flags, false
		}
		for _, value := range locations {
			location, known := storageLocation(value, aliases)
			flags.unknown = flags.unknown || !known
			flags.listed = flags.listed || location == region
		}
	} else {
		flags.unknown = true
	}
	if raw, present := row["restrictions"]; present {
		restrictions, ok := providerArray(raw)
		if !ok {
			return flags, false
		}
		for _, rawRestriction := range restrictions {
			if ctx.Err() != nil {
				return flags, false
			}
			restriction, valid := providerObject(rawRestriction, "type", "values", "reasonCode")
			if !valid || !diskOptionalStrings(ctx, restriction, "type", "reasonCode") {
				return flags, false
			}
			var kind, reason string
			if raw, present := restriction["type"]; present {
				if !diskScalar(raw, &kind) {
					return flags, false
				}
			}
			if raw, present := restriction["reasonCode"]; present {
				if !diskScalar(raw, &reason) {
					return flags, false
				}
			}
			var values []string
			if raw, present := restriction["values"]; present {
				values, valid = diskStrings(ctx, raw, nil)
				if !valid {
					return flags, false
				}
			}
			interpretable := storageASCII(kind) && strings.ToLower(kind) == "location" && storageASCII(reason) &&
				(strings.ToLower(reason) == "quotaid" || strings.ToLower(reason) == "notavailableforsubscription") && len(values) > 0
			flags.unknown = flags.unknown || !interpretable
			for _, value := range values {
				location, known := storageLocation(value, aliases)
				flags.unknown = flags.unknown || !known
				if interpretable && known && location == region {
					flags.hard = flags.hard || strings.ToLower(reason) == "quotaid"
					flags.restricted = flags.restricted || strings.ToLower(reason) == "notavailableforsubscription"
				}
			}
		}
	} else {
		flags.unknown = true
	}
	if raw, present := row["capabilities"]; present {
		capabilities, ok := providerArray(raw)
		if !ok {
			return flags, false
		}
		availableValue := ""
		for _, rawCapability := range capabilities {
			capability, valid := providerObject(rawCapability, "name", "value")
			if !valid || !diskOptionalStrings(ctx, capability, "name", "value") {
				return flags, false
			}
			var name, value string
			if raw, present := capability["name"]; present {
				if !diskScalar(raw, &name) {
					return flags, false
				}
			}
			if raw, present := capability["value"]; present {
				if !diskScalar(raw, &value) {
					return flags, false
				}
			}
			if name == "" || value == "" {
				flags.unknown = true
			}
			if storageASCII(name) && strings.ToLower(name) == "available" {
				normalized := strings.ToLower(value)
				if !storageASCII(value) || normalized != "true" && normalized != "false" {
					flags.unknown = true
				} else {
					flags.unknown = flags.unknown || availableValue != "" && availableValue != normalized
					availableValue = normalized
					flags.hard = flags.hard || normalized == "false"
				}
			}
		}
	}
	return flags, ctx.Err() == nil
}

func storageDecode(ctx context.Context, body []byte, region string, aliases map[string]string) (map[string]map[string]StorageSKUGroup, map[string]bool, string) {
	if !quotaJSON(ctx, body) {
		return nil, nil, "storage_sku_invalid_response"
	}
	envelope, valid := providerObject(body, "value", "nextLink")
	if !valid {
		return nil, nil, "storage_sku_invalid_response"
	}
	rows, valid := providerArray(envelope["value"])
	if !valid {
		return nil, nil, "storage_sku_invalid_response"
	}
	if _, present := envelope["nextLink"]; present {
		return nil, nil, "storage_sku_unexpected_paging"
	}
	if len(rows) > MaxStorageSKURows {
		return nil, nil, "storage_sku_row_limit"
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var document any
	if decoder.Decode(&document) != nil {
		return nil, nil, "storage_sku_invalid_response"
	}
	budget := storageBudget{ctx: ctx}
	if code := budget.walk(document); code != "" {
		return nil, nil, code
	}
	groups := map[string]map[string]storageFlags{}
	unqualified := map[string]bool{}
	groupCount := 0
	for _, raw := range rows {
		row, valid := providerObject(raw, "name", "resourceType", "kind", "tier", "locations", "restrictions", "capabilities")
		var name, resourceType, kind string
		if ctx.Err() != nil || !valid || !diskScalar(row["name"], &name) || !diskName(name) {
			return nil, nil, "storage_sku_invalid_response"
		}
		if raw, present := row["resourceType"]; present && (!diskScalar(raw, &resourceType) || !storageIdentifier(resourceType)) {
			return nil, nil, "storage_sku_invalid_response"
		}
		if raw, present := row["kind"]; present && (!diskScalar(raw, &kind) || !storageIdentifier(kind)) {
			return nil, nil, "storage_sku_invalid_response"
		}
		flags, valid := storageRowFlags(ctx, row, region, aliases)
		if !valid {
			return nil, nil, "storage_sku_invalid_response"
		}
		if resourceType != "" && strings.ToLower(resourceType) != "storageaccounts" {
			continue
		}
		name, kind = strings.ToLower(name), strings.ToLower(kind)
		if resourceType == "" || kind == "" {
			unqualified[name] = true
			continue
		}
		if groups[name] == nil {
			groups[name] = map[string]storageFlags{}
		}
		prior, found := groups[name][kind]
		if !found {
			if groupCount >= MaxStorageSKURows {
				return nil, nil, "storage_sku_row_limit"
			}
			groupCount++
		}
		groups[name][kind] = storageFlags{listed: prior.listed || flags.listed, hard: prior.hard || flags.hard, restricted: prior.restricted || flags.restricted, unknown: prior.unknown || flags.unknown}
	}
	out := map[string]map[string]StorageSKUGroup{}
	for name, kinds := range groups {
		out[name] = map[string]StorageSKUGroup{}
		for kind, flags := range kinds {
			if ctx.Err() != nil {
				return nil, nil, "storage_sku_work_limit"
			}
			value := StorageSKUGroup{Known: !flags.unknown, State: SKUUnavailable}
			if value.Known && flags.listed && !flags.hard {
				value.State = SKUAvailable
				if flags.restricted {
					value.State = SKURestricted
				}
			}
			out[name][kind] = value
		}
	}
	return out, unqualified, ""
}

func (c *StorageSKUCollector) Collect(ctx context.Context, subscriptions map[string]string, request StorageSKURequest) (*StorageSKUCollection, error) {
	reject := func() (*StorageSKUCollection, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Storage SKU request rejected [storage_sku_request_invalid]")
	}
	if ctx.Err() != nil || c == nil || c.transport == nil {
		return reject()
	}
	text := 0
	scope, valid := auxiliaryScope(ctx, subscriptions, &text)
	if !valid {
		return reject()
	}
	for _, name := range scope {
		if name == "" {
			return reject()
		}
	}
	id := strings.ToLower(request.SubscriptionID)
	region, valid := diskPhysical(request.Region)
	if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || scope[id] == "" || !valid || region != request.Region || len(request.LocationAliases) > 256 {
		return reject()
	}
	aliases := map[string]string{}
	for key, value := range request.LocationAliases {
		physical, valid := diskPhysical(value)
		canonical, isCanonical := diskPhysical(key)
		if ctx.Err() != nil || !diskName(key) || strings.ToLower(key) != key || storageStage(key) || !valid || physical != value || isCanonical && canonical != value || !auxiliaryText(key, &text) || !auxiliaryText(value, &text) {
			return reject()
		}
		aliases[key] = value
	}
	target := *c.transport.origin
	target.Path = "/subscriptions/" + id + "/providers/Microsoft.Storage/skus"
	target.RawQuery = "api-version=" + storageSKUVersion
	result := &StorageSKUCollection{SubscriptionID: id, Region: region, Status: "unknown", Groups: map[string]map[string]StorageSKUGroup{}, UnqualifiedNames: map[string]bool{}}
	fail := func(code string) (*StorageSKUCollection, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result.FailureCode = code
		return result, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	body, response, err := c.transport.getter.GetBoundedWithResponse(ctx, target.String(), MaxStorageSKUBytes)
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
		return fail("storage_sku_status_invalid")
	}
	if err != nil || response == nil {
		return fail("storage_sku_request_failed")
	}
	if len(body) > MaxStorageSKUBytes {
		return fail("storage_sku_byte_limit")
	}
	groups, unqualified, code := storageDecode(ctx, body, region, aliases)
	if code != "" {
		return fail(code)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	result.Status, result.Groups, result.UnqualifiedNames = "complete", groups, unqualified
	return result, nil
}

func storageFailureCode(code string) bool {
	switch code {
	case "", "storage_sku_status_invalid", "storage_sku_request_failed", "storage_sku_byte_limit", "storage_sku_invalid_response", "storage_sku_unexpected_paging", "storage_sku_row_limit", "storage_sku_work_limit", "storage_sku_text_limit":
		return true
	}
	return false
}

// Evidence validates stable caller inputs, then owns the complete projection.
// Unknown entries cannot be represented individually by the current consumer.
func (c *StorageSKUCollection) Evidence(ctx context.Context, requiredKinds map[string][]string) (SKUEvidence, error) {
	unknown := SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}}
	incomplete := func() (SKUEvidence, error) {
		if err := ctx.Err(); err != nil {
			return SKUEvidence{}, err
		}
		return unknown, nil
	}
	fail := func(code string) (SKUEvidence, error) {
		if err := ctx.Err(); err != nil {
			return SKUEvidence{}, err
		}
		return SKUEvidence{}, fmt.Errorf("Storage SKU projection rejected [%s]", code)
	}
	if ctx.Err() != nil {
		return fail("storage_sku_projection_invalid")
	}
	budget := storageBudget{ctx: ctx}
	if c != nil {
		for _, label := range []string{c.SubscriptionID, c.Region, c.Status, c.FailureCode} {
			if !budget.label(label) {
				return fail("storage_sku_projection_limit")
			}
		}
		region, physical := diskPhysical(c.Region)
		if len(c.SubscriptionID) != 36 || !subscriptionID.MatchString(c.SubscriptionID) || strings.ToLower(c.SubscriptionID) != c.SubscriptionID || !physical || region != c.Region ||
			(c.Status != "complete" && c.Status != "unknown") || !storageFailureCode(c.FailureCode) || c.Status == "complete" && c.FailureCode != "" ||
			c.Status == "unknown" && (len(c.Groups) != 0 || len(c.UnqualifiedNames) != 0) {
			return fail("storage_sku_projection_invalid")
		}
		if len(c.Groups) > MaxStorageSKURows || len(c.UnqualifiedNames) > MaxStorageSKURows {
			return fail("storage_sku_projection_limit")
		}
		groups := 0
		for name, kinds := range c.Groups {
			if !budget.entry() || !budget.label(name) {
				return fail("storage_sku_projection_limit")
			}
			if !diskName(name) || strings.ToLower(name) != name || len(kinds) == 0 {
				return fail("storage_sku_projection_invalid")
			}
			for kind, value := range kinds {
				groups++
				if groups > MaxStorageSKURows || !budget.entry() || !budget.label(kind) {
					return fail("storage_sku_projection_limit")
				}
				if !storageIdentifier(kind) || strings.ToLower(kind) != kind || value.State < SKUAvailable || value.State > SKUUnavailable || !value.Known && value.State != SKUUnavailable {
					return fail("storage_sku_projection_invalid")
				}
			}
		}
		for name, present := range c.UnqualifiedNames {
			if !budget.entry() || !budget.label(name) {
				return fail("storage_sku_projection_limit")
			}
			if !present || !diskName(name) || strings.ToLower(name) != name {
				return fail("storage_sku_projection_invalid")
			}
		}
	}
	if len(requiredKinds) > MaxStorageSKURows {
		return fail("storage_sku_projection_limit")
	}
	required := map[string]map[string]bool{}
	missingKind := map[string]bool{}
	for raw, kinds := range requiredKinds {
		if !budget.entry() || !budget.label(raw) || len(kinds) > 16 {
			return fail("storage_sku_projection_limit")
		}
		name := strings.ToLower(raw)
		if !diskName(raw) || required[name] != nil {
			return fail("storage_sku_projection_invalid")
		}
		required[name] = map[string]bool{}
		missingKind[name] = len(kinds) == 0
		for _, rawKind := range kinds {
			if !budget.entry() || !budget.label(rawKind) {
				return fail("storage_sku_projection_limit")
			}
			if rawKind == "" {
				missingKind[name] = true
				continue
			}
			if !storageIdentifier(rawKind) {
				return fail("storage_sku_projection_invalid")
			}
			required[name][strings.ToLower(rawKind)] = true
		}
	}
	if c == nil || c.Status != "complete" {
		return incomplete()
	}
	out := SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{}}
	outputText := 0
	for name, kinds := range required {
		if ctx.Err() != nil {
			return fail("storage_sku_projection_invalid")
		}
		if missingKind[name] || c.UnqualifiedNames[name] {
			return incomplete()
		}
		state := SKUAvailable
		for kind := range kinds {
			if ctx.Err() != nil {
				return fail("storage_sku_projection_invalid")
			}
			group, found := c.Groups[name][kind]
			if found && !group.Known {
				return incomplete()
			}
			if !found || group.State == SKUUnavailable {
				state = SKUUnavailable
			} else if group.State == SKURestricted && state != SKUUnavailable {
				state = SKURestricted
			}
		}
		if !auxiliaryText(name, &outputText) {
			return fail("storage_sku_projection_limit")
		}
		out.Values[name] = SKUAvailability{State: state, BlockedZones: []string{}}
	}
	if ctx.Err() != nil {
		return fail("storage_sku_projection_invalid")
	}
	return out, nil
}
