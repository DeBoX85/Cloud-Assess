// Compute SKU behavior derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_VM_VMSS_SKU_CONTRACT_20261016.md for corrections.
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

const vmSKUVersion = "2021-07-01"
const vmSKUType = "microsoft.compute/virtualmachines"
const vmSSSKUType = "microsoft.compute/virtualmachinescalesets"
const MaxComputeVMSKUWork = 65536
const MaxComputeVMSKUGroups = 8192

type ComputeVMSKURequest struct{ SubscriptionID, Region, ResourceType string }
type ComputeVMSKUGroup struct {
	Known        bool
	State        int
	BlockedZones []string
}
type ComputeVMSKUCollection struct {
	SubscriptionID, Region, ResourceType, Status string
	Groups                                       map[string]ComputeVMSKUGroup
	FailureCode                                  string
}
type ComputeVMSKUCollector struct{ transport *RESTQuotaCollector }

func NewComputeVMSKUCollector(endpoint string, getter RESTQuotaGetter) (*ComputeVMSKUCollector, error) {
	t, e := NewRESTQuotaCollector(endpoint, getter)
	if e != nil {
		return nil, fmt.Errorf("invalid Compute VM SKU origin or getter")
	}
	return &ComputeVMSKUCollector{transport: t}, nil
}

type vmSKUBudget struct {
	ctx        context.Context
	work, text int
}

func (b *vmSKUBudget) entry() bool {
	if b.ctx.Err() != nil || b.work >= MaxComputeVMSKUWork {
		return false
	}
	b.work++
	return true
}
func (b *vmSKUBudget) label(s string, limit int) bool {
	if b.ctx.Err() != nil || len(s) > limit || len(s) > auxTextBudget-b.text {
		return false
	}
	b.text += len(s)
	return true
}
func (b *vmSKUBudget) walk(v any, depth int, next bool) string {
	if b.ctx.Err() != nil {
		return "vm_sku_work_limit"
	}
	switch v := v.(type) {
	case map[string]any:
		if !b.entry() {
			return "vm_sku_work_limit"
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if !b.label(k, MaxLabelBytes) {
				return "vm_sku_text_limit"
			}
			if code := b.walk(v[k], depth+1, depth == 0 && k == "nextLink"); code != "" {
				return code
			}
		}
	case []any:
		for _, v := range v {
			if !b.entry() {
				return "vm_sku_work_limit"
			}
			if code := b.walk(v, depth+1, false); code != "" {
				return code
			}
		}
	case string:
		limit := MaxLabelBytes
		if next {
			limit = 8192
		}
		if !b.label(v, limit) {
			return "vm_sku_text_limit"
		}
	}
	return ""
}

// Reuse the unchanged strict metadata reader except for the deliberately opaque
// scaleType and nonnegative capacity policy declared in this family's contract.
func vmSKUMetadata(ctx context.Context, row map[string]json.RawMessage) bool {
	rest := row
	if raw, has := row["capacity"]; has {
		capacity, ok := providerObject(raw, "minimum", "maximum", "default", "scaleType")
		if !ok || !diskOptionalStrings(ctx, capacity, "scaleType") {
			return false
		}
		for _, key := range []string{"minimum", "maximum", "default"} {
			if !diskOptional(capacity, key, func(raw json.RawMessage) bool { var n int64; return diskScalar(raw, &n) && n >= 0 }) {
				return false
			}
		}
		rest = make(map[string]json.RawMessage, len(row))
		for k, v := range row {
			if k != "capacity" {
				rest[k] = v
			}
		}
	}
	return diskMetadata(ctx, rest)
}

func vmSKUAvailability(ctx context.Context, row map[string]json.RawMessage, region string) (ComputeVMSKUGroup, bool) {
	out := ComputeVMSKUGroup{Known: true, State: SKUUnavailable, BlockedZones: []string{}}
	listed, hard, restricted := false, false, false
	blocked := map[string]bool{}
	if raw, has := row["locations"]; has {
		locations, ok := diskLocations(ctx, raw)
		if !ok {
			return out, false
		}
		for _, loc := range locations {
			listed = listed || loc == region
		}
	} else {
		out.Known = false
	}
	if raw, has := row["restrictions"]; has {
		restrictions, ok := providerArray(raw)
		if !ok {
			return out, false
		}
		for _, raw := range restrictions {
			r, ok := providerObject(raw, "type", "reasonCode", "values", "restrictionInfo")
			if !ok || ctx.Err() != nil || !diskOptionalStrings(ctx, r, "type", "reasonCode") {
				return out, false
			}
			var kind, reason string
			if v, has := r["type"]; has && !diskScalar(v, &kind) {
				return out, false
			}
			if v, has := r["reasonCode"]; has && !diskScalar(v, &reason) {
				return out, false
			}
			var values, locations, zones []string
			if v, has := r["values"]; has {
				values, ok = diskLocations(ctx, v)
				if !ok {
					return out, false
				}
			}
			if v, has := r["restrictionInfo"]; has {
				info, valid := providerObject(v, "locations", "zones")
				if !valid {
					return out, false
				}
				if v, has := info["locations"]; has {
					locations, ok = diskLocations(ctx, v)
					if !ok {
						return out, false
					}
				}
				if v, has := info["zones"]; has {
					zones, ok = diskStrings(ctx, v, diskZone)
					if !ok {
						return out, false
					}
				}
			}
			if (reason != "QuotaId" && reason != "NotAvailableForSubscription") || (kind != "Location" && kind != "Zone") {
				out.Known = false
				continue
			}
			if kind == "Location" {
				if len(values)+len(locations) == 0 {
					out.Known = false
				}
				for _, loc := range append(values, locations...) {
					restricted = restricted || loc == region
				}
			} else {
				if len(locations) == 0 || len(zones) == 0 {
					out.Known = false
				}
				for _, loc := range locations {
					if loc == region {
						for _, z := range zones {
							blocked[z] = true
						}
					}
				}
			}
		}
	} else {
		out.Known = false
	}
	if raw, has := row["capabilities"]; has {
		caps, ok := providerArray(raw)
		if !ok {
			return out, false
		}
		prior := ""
		for _, raw := range caps {
			cap, ok := providerObject(raw, "name", "value")
			if !ok || !diskOptionalStrings(ctx, cap, "name", "value") {
				return out, false
			}
			var name, value string
			if v, has := cap["name"]; has {
				if !diskScalar(v, &name) {
					return out, false
				}
			} else {
				out.Known = false
			}
			if v, has := cap["value"]; has {
				if !diskScalar(v, &value) {
					return out, false
				}
			} else {
				out.Known = false
			}
			if storageASCII(name) && strings.EqualFold(name, "available") {
				value = strings.ToLower(value)
				if !storageASCII(value) || (value != "true" && value != "false") || prior != "" && prior != value {
					out.Known = false
				}
				prior = value
				hard = hard || value == "false"
			}
		}
	}
	if !out.Known {
		return out, ctx.Err() == nil
	}
	if !listed || hard {
		return out, ctx.Err() == nil
	}
	out.State = SKUAvailable
	if restricted {
		out.State = SKURestricted
	} else if len(blocked) > 0 {
		out.State = SKUZoneRestricted
		for z := range blocked {
			out.BlockedZones = append(out.BlockedZones, z)
		}
		sort.Strings(out.BlockedZones)
	}
	return out, ctx.Err() == nil
}

type vmSKUPage struct {
	groups     map[string]ComputeVMSKUGroup
	next       string
	work, text int
}

func decodeVMSKUPage(ctx context.Context, body []byte, region string, prior map[string]ComputeVMSKUGroup, work, text int) (vmSKUPage, string) {
	bad := func(code string) (vmSKUPage, string) { return vmSKUPage{}, code }
	if !quotaJSON(ctx, body) {
		return bad("vm_sku_invalid_response")
	}
	root, ok := providerObject(body, "value", "nextLink")
	if !ok {
		return bad("vm_sku_invalid_response")
	}
	rows, ok := providerArray(root["value"])
	if !ok {
		return bad("vm_sku_invalid_response")
	}
	page := vmSKUPage{groups: map[string]ComputeVMSKUGroup{}}
	if raw, has := root["nextLink"]; has && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) && !diskScalar(raw, &page.next) {
		return bad("vm_sku_invalid_response")
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var doc any
	if dec.Decode(&doc) != nil {
		return bad("vm_sku_invalid_response")
	}
	budget := vmSKUBudget{ctx: ctx, work: work, text: text}
	if code := budget.walk(doc, 0, false); code != "" {
		return bad(code)
	}
	for _, raw := range rows {
		row, ok := providerObject(raw, "resourceType", "name", "tier", "size", "family", "kind", "capacity", "locations", "locationInfo", "apiVersions", "costs", "capabilities", "restrictions")
		var kind, name string
		if !ok || ctx.Err() != nil || !diskScalar(row["resourceType"], &kind) || !diskName(kind) || !diskScalar(row["name"], &name) || !diskName(name) || !vmSKUMetadata(ctx, row) {
			return bad("vm_sku_invalid_response")
		}
		value, ok := vmSKUAvailability(ctx, row, region)
		if !ok {
			return bad("vm_sku_invalid_response")
		}
		if !strings.EqualFold(kind, "virtualMachines") {
			continue
		}
		name = strings.ToLower(name)
		if _, exists := prior[name]; exists {
			return bad("vm_sku_collision")
		}
		if _, exists := page.groups[name]; exists {
			return bad("vm_sku_collision")
		}
		if len(prior)+len(page.groups) >= MaxComputeVMSKUGroups {
			return bad("vm_sku_work_limit")
		}
		page.groups[name] = value
	}
	page.work, page.text = budget.work, budget.text
	return page, ""
}

func (c *ComputeVMSKUCollector) continuation(first *url.URL, filter, link string) (*url.URL, bool) {
	// The established disk URL guard is read-only; this family additionally requires
	// the unchanged location filter on every continuation, including the first URL.
	shared := ComputeDiskSKUCollector{transport: c.transport}
	u, ok := shared.continuation(first, filter, link)
	if !ok {
		return nil, false
	}
	q := u.Query()
	if len(q["$filter"]) != 1 || q.Get("$filter") != filter || q.Get("api-version") != vmSKUVersion {
		return nil, false
	}
	return u, true
}

func (c *ComputeVMSKUCollector) Collect(ctx context.Context, subscriptions map[string]string, q ComputeVMSKURequest) (*ComputeVMSKUCollection, error) {
	reject := func() (*ComputeVMSKUCollection, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		return nil, fmt.Errorf("Compute VM SKU request rejected [vm_sku_request_invalid]")
	}
	if ctx.Err() != nil || c == nil || c.transport == nil {
		return reject()
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &text)
	if !ok {
		return reject()
	}
	for _, v := range scope {
		if v == "" {
			return reject()
		}
	}
	id := strings.ToLower(q.SubscriptionID)
	region, ok := diskPhysical(q.Region)
	if len(q.SubscriptionID) != 36 || !subscriptionID.MatchString(q.SubscriptionID) || scope[id] == "" || !ok || region != q.Region || (q.ResourceType != vmSKUType && q.ResourceType != vmSSSKUType) {
		return reject()
	}
	first := *c.transport.origin
	first.Path = "/subscriptions/" + id + "/providers/Microsoft.Compute/skus"
	filter := "location eq '" + region + "'"
	query := url.Values{}
	query.Set("api-version", vmSKUVersion)
	query.Set("$filter", filter)
	first.RawQuery = query.Encode()
	result := &ComputeVMSKUCollection{SubscriptionID: id, Region: region, ResourceType: q.ResourceType, Status: "complete", Groups: map[string]ComputeVMSKUGroup{}}
	pages, total, work, decodedText := 0, 0, 0, 0
	seen := map[string]bool{}
	link := first.String()
	fail := func(code string) (*ComputeVMSKUCollection, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		result.Status = "unknown"
		if pages > 0 {
			result.Status = "partial"
		}
		result.FailureCode = code
		return result, nil
	}
	for {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if pages >= MaxQuotaPages {
			return fail("vm_sku_page_limit")
		}
		target, ok := c.continuation(&first, filter, link)
		if !ok {
			return fail("vm_sku_continuation_invalid")
		}
		if seen[target.String()] {
			return fail("vm_sku_cycle")
		}
		seen[target.String()] = true
		body, response, e := c.transport.getter.GetBoundedWithResponse(ctx, target.String(), MaxQuotaPageBytes)
		if cause := ctx.Err(); cause != nil {
			return nil, cause
		}
		if errors.Is(e, context.Canceled) {
			return nil, context.Canceled
		}
		if errors.Is(e, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		if response != nil && response.StatusCode != http.StatusOK {
			return fail("vm_sku_status_invalid")
		}
		if e != nil || response == nil {
			return fail("vm_sku_request_failed")
		}
		if len(body) > MaxQuotaPageBytes {
			return fail("vm_sku_byte_limit")
		}
		if len(body) > MaxQuotaTotalBytes-total {
			return fail("vm_sku_aggregate_limit")
		}
		page, code := decodeVMSKUPage(ctx, body, region, result.Groups, work, decodedText)
		if code != "" {
			return fail(code)
		}
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		for k, v := range page.groups {
			result.Groups[k] = v
		}
		total += len(body)
		work, decodedText = page.work, page.text
		pages++
		if page.next == "" {
			break
		}
		link = page.next
	}
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return result, nil
}

func vmSKUFailure(code string) bool {
	switch code {
	case "", "vm_sku_request_failed", "vm_sku_status_invalid", "vm_sku_byte_limit", "vm_sku_aggregate_limit", "vm_sku_invalid_response", "vm_sku_collision", "vm_sku_work_limit", "vm_sku_text_limit", "vm_sku_page_limit", "vm_sku_continuation_invalid", "vm_sku_cycle":
		return true
	}
	return false
}

func (c *ComputeVMSKUCollection) Evidence(ctx context.Context, required []string) (SKUEvidence, error) {
	empty := func(status string) (SKUEvidence, error) {
		if e := ctx.Err(); e != nil {
			return SKUEvidence{}, e
		}
		return SKUEvidence{Status: status, Values: map[string]SKUAvailability{}}, nil
	}
	fail := func(code string) (SKUEvidence, error) {
		if e := ctx.Err(); e != nil {
			return SKUEvidence{}, e
		}
		return SKUEvidence{}, fmt.Errorf("Compute VM SKU projection rejected [%s]", code)
	}
	if ctx.Err() != nil {
		return fail("vm_sku_projection_invalid")
	}
	budget := vmSKUBudget{ctx: ctx}
	if c != nil {
		for _, s := range []string{c.SubscriptionID, c.Region, c.ResourceType, c.Status} {
			if !budget.label(s, MaxLabelBytes) {
				return fail("vm_sku_projection_limit")
			}
		}
		region, ok := diskPhysical(c.Region)
		if len(c.SubscriptionID) != 36 || !subscriptionID.MatchString(c.SubscriptionID) || strings.ToLower(c.SubscriptionID) != c.SubscriptionID || !ok || region != c.Region || (c.ResourceType != vmSKUType && c.ResourceType != vmSSSKUType) || (c.Status != "complete" && c.Status != "partial" && c.Status != "unknown") || !vmSKUFailure(c.FailureCode) || c.Status == "complete" && c.FailureCode != "" || c.Status == "unknown" && len(c.Groups) != 0 || c.Status == "partial" && c.FailureCode == "" {
			return fail("vm_sku_projection_invalid")
		}
		// FailureCode is a validated fixed enum; its bytes still charge the input budget.
		budget.text += len(c.FailureCode)
		if len(c.Groups) > MaxComputeVMSKUGroups {
			return fail("vm_sku_projection_limit")
		}
		for name, g := range c.Groups {
			if !budget.entry() || !budget.label(name, MaxLabelBytes) {
				return fail("vm_sku_projection_limit")
			}
			if !diskName(name) || strings.ToLower(name) != name || g.State < SKUAvailable || g.State > SKUZoneRestricted || !g.Known && g.State != SKUUnavailable || g.State != SKUZoneRestricted && len(g.BlockedZones) != 0 || g.State == SKUZoneRestricted && len(g.BlockedZones) == 0 {
				return fail("vm_sku_projection_invalid")
			}
			prior := ""
			for _, z := range g.BlockedZones {
				if !budget.entry() || !budget.label(z, MaxLabelBytes) {
					return fail("vm_sku_projection_limit")
				}
				if !diskZone(z) || z <= prior {
					return fail("vm_sku_projection_invalid")
				}
				prior = z
			}
		}
	}
	if len(required) > MaxComputeVMSKUGroups {
		return fail("vm_sku_projection_limit")
	}
	names := map[string]bool{}
	for _, s := range required {
		if !budget.entry() || !budget.label(s, MaxLabelBytes) {
			return fail("vm_sku_projection_limit")
		}
		name := strings.ToLower(s)
		if !diskName(s) || names[name] {
			return fail("vm_sku_projection_invalid")
		}
		names[name] = true
	}
	if c == nil || c.Status != "complete" {
		return empty("unknown")
	}
	if c.ResourceType == vmSSSKUType {
		return empty("unsupported")
	}
	out := SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{}}
	output := 0
	for name := range names {
		if e := ctx.Err(); e != nil {
			return SKUEvidence{}, e
		}
		g, found := c.Groups[name]
		if found && !g.Known {
			return empty("unknown")
		}
		if !found {
			g = ComputeVMSKUGroup{Known: true, State: SKUUnavailable}
		}
		if !auxiliaryText(name, &output) {
			return fail("vm_sku_projection_limit")
		}
		zones := make([]string, 0, len(g.BlockedZones))
		for _, z := range g.BlockedZones {
			if e := ctx.Err(); e != nil {
				return SKUEvidence{}, e
			}
			if !auxiliaryText(z, &output) {
				return fail("vm_sku_projection_limit")
			}
			zones = append(zones, z)
		}
		out.Values[name] = SKUAvailability{State: g.State, BlockedZones: zones}
	}
	if e := ctx.Err(); e != nil {
		return SKUEvidence{}, e
	}
	return out, nil
}
