// Cognitive account catalog behavior derived from MIT-licensed Azure Quick Review
// and Microsoft REST specifications. See NOTICE.md and the frozen Cognitive contract.
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

const cognitiveSKUVersion = "2024-10-01"
const cognitiveSKUType = "microsoft.cognitiveservices/accounts"
const MaxCognitiveSKURows = 8192
const MaxCognitiveSKUGroups = 8192
const MaxCognitiveSKUAliases = 256

type CognitiveSKURequest struct {
	SubscriptionID, Region, Kind string
	LocationAliases              map[string]string
}
type CognitiveSKUGroup struct {
	Known        bool
	State        int
	BlockedZones []string
}
type CognitiveSKUCollection struct {
	SubscriptionID, Region, Kind, ResourceType, Status string
	Groups                                             map[string]CognitiveSKUGroup
	FailureCode                                        string
}
type CognitiveSKUCollector struct{ transport *RESTQuotaCollector }

func NewCognitiveSKUCollector(endpoint string, getter RESTQuotaGetter) (*CognitiveSKUCollector, error) {
	t, e := NewRESTQuotaCollector(endpoint, getter)
	if e != nil {
		return nil, fmt.Errorf("invalid Cognitive SKU origin or getter")
	}
	return &CognitiveSKUCollector{transport: t}, nil
}

// Alias provenance is deliberately narrower than general display-name guessing.
func cognitiveAliases(ctx context.Context, aliases map[string]string) (map[string]string, bool) {
	if len(aliases) > MaxCognitiveSKUAliases {
		return nil, false
	}
	out := make(map[string]string, len(aliases))
	for key, value := range aliases {
		physical, ok := diskPhysical(value)
		if ctx.Err() != nil || !diskName(key) || strings.ToLower(key) != key || !ok || physical != value || strings.ReplaceAll(key, " ", "") != value || storageStage(key) || !costPhysical(key) {
			return nil, false
		}
		if p, physicalKey := diskPhysical(key); physicalKey && p != value {
			return nil, false
		}
		out[key] = value
	}
	return out, true
}
func cognitiveBinding(scope map[string]string, q CognitiveSKURequest) (string, string, bool) {
	if len(q.SubscriptionID) != 36 || len(q.Kind) > MaxLabelBytes {
		return "", "", false
	}
	id := strings.ToLower(q.SubscriptionID)
	kind := strings.ToLower(q.Kind)
	region, ok := diskPhysical(q.Region)
	return id, kind, len(id) == 36 && subscriptionID.MatchString(id) && scope[id] != "" && ok && region == q.Region && storageIdentifier(q.Kind)
}

type cognitiveFlags struct {
	listed, restricted, unknown bool
	zones, scopes               map[string]bool
}

func (f cognitiveFlags) group() CognitiveSKUGroup {
	g := CognitiveSKUGroup{Known: !f.unknown, State: SKUUnavailable, BlockedZones: []string{}}
	if f.unknown || !f.listed {
		return g
	}
	g.State = SKUAvailable
	if f.restricted {
		g.State = SKURestricted
	} else if len(f.zones) > 0 {
		g.State = SKUZoneRestricted
		for z := range f.zones {
			g.BlockedZones = append(g.BlockedZones, z)
		}
		sort.Strings(g.BlockedZones)
	}
	return g
}
func cognitiveRow(ctx context.Context, row map[string]json.RawMessage, region string, aliases map[string]string) (cognitiveFlags, bool) {
	f := cognitiveFlags{zones: map[string]bool{}, scopes: map[string]bool{}}
	if !diskOptionalStrings(ctx, row, "tier") {
		return f, false
	}
	if raw, has := row["locations"]; has {
		locations, ok := diskStrings(ctx, raw, nil)
		if !ok {
			return f, false
		}
		if len(locations) == 0 {
			f.unknown = true
			f.scopes["@unknown"] = true
		}
		for _, loc := range locations {
			physical, known := storageLocation(loc, aliases)
			scope := physical
			if scope == "" {
				scope = "@" + strings.ToLower(loc)
			}
			f.scopes[scope] = true
			f.unknown = f.unknown || !known
			f.listed = f.listed || physical == region
		}
	} else {
		f.unknown = true
		f.scopes["@unknown"] = true
	}
	if raw, has := row["restrictions"]; has {
		restrictions, ok := providerArray(raw)
		if !ok {
			return f, false
		}
		for _, raw := range restrictions {
			r, ok := providerObject(raw, "type", "reasonCode", "values", "restrictionInfo")
			if !ok || !diskOptionalStrings(ctx, r, "type", "reasonCode") {
				return f, false
			}
			var kind, reason string
			if v, has := r["type"]; has && !diskScalar(v, &kind) {
				return f, false
			}
			if v, has := r["reasonCode"]; has && !diskScalar(v, &reason) {
				return f, false
			}
			var values, locations, zones []string
			if v, has := r["values"]; has {
				values, ok = diskStrings(ctx, v, nil)
				if !ok {
					return f, false
				}
			}
			if v, has := r["restrictionInfo"]; has {
				info, valid := providerObject(v, "locations", "zones")
				if !valid {
					return f, false
				}
				if v, has := info["locations"]; has {
					locations, ok = diskStrings(ctx, v, nil)
					if !ok {
						return f, false
					}
				}
				if v, has := info["zones"]; has {
					zones, ok = diskStrings(ctx, v, diskZone)
					if !ok {
						return f, false
					}
				}
			}
			// All supplied scopes are validated even for unknown/foreign restrictions.
			physicalValues := []string{}
			physicalLocations := []string{}
			for _, loc := range values {
				p, known := storageLocation(loc, aliases)
				f.unknown = f.unknown || !known
				physicalValues = append(physicalValues, p)
			}
			for _, loc := range locations {
				p, known := storageLocation(loc, aliases)
				f.unknown = f.unknown || !known
				physicalLocations = append(physicalLocations, p)
			}
			if (reason != "QuotaId" && reason != "NotAvailableForSubscription") || (kind != "Location" && kind != "Zone") {
				f.unknown = true
				continue
			}
			if kind == "Location" {
				if len(values)+len(locations) == 0 {
					f.unknown = true
				}
				for _, loc := range append(physicalValues, physicalLocations...) {
					f.restricted = f.restricted || loc == region
				}
			} else {
				if len(locations) == 0 || len(zones) == 0 {
					f.unknown = true
				}
				for _, loc := range physicalLocations {
					if loc == region {
						for _, z := range zones {
							f.zones[z] = true
						}
					}
				}
			}
		}
	} else {
		f.unknown = true
	}
	return f, ctx.Err() == nil
}

type cognitivePage struct {
	flags            map[string]cognitiveFlags
	next             string
	rows, work, text int
}

func decodeCognitivePage(ctx context.Context, body []byte, q CognitiveSKURequest, prior map[string]cognitiveFlags, rows, work, text int) (cognitivePage, string) {
	bad := func(code string) (cognitivePage, string) { return cognitivePage{}, "cognitive_sku_" + code }
	if !quotaJSON(ctx, body) {
		return bad("invalid_response")
	}
	root, ok := providerObject(body, "value", "nextLink")
	if !ok {
		return bad("invalid_response")
	}
	rawRows, ok := providerArray(root["value"])
	if !ok {
		return bad("invalid_response")
	}
	page := cognitivePage{flags: map[string]cognitiveFlags{}, rows: rows + len(rawRows)}
	if page.rows > MaxCognitiveSKURows {
		return bad("row_limit")
	}
	if raw, has := root["nextLink"]; has && !diskScalar(raw, &page.next) {
		return bad("invalid_response")
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var doc any
	if dec.Decode(&doc) != nil {
		return bad("invalid_response")
	}
	budget := vmSKUBudget{ctx: ctx, work: work, text: text}
	if code := budget.walk(doc, 0, false); code != "" {
		switch code {
		case "vm_sku_work_limit":
			return bad("work_limit")
		case "vm_sku_text_limit":
			return bad("text_limit")
		default:
			return bad("invalid_response")
		}
	}
	// Shadow prior state. A rejected page can never mutate admitted diagnostics.
	for name, f := range prior {
		if ctx.Err() != nil {
			return bad("invalid_response")
		}
		copyFlags := f
		copyFlags.zones = map[string]bool{}
		copyFlags.scopes = map[string]bool{}
		for z := range f.zones {
			copyFlags.zones[z] = true
		}
		for loc := range f.scopes {
			copyFlags.scopes[loc] = true
		}
		page.flags[name] = copyFlags
	}
	for _, raw := range rawRows {
		row, ok := providerObject(raw, "resourceType", "name", "kind", "tier", "locations", "restrictions")
		var resourceType, name, kind string
		if !ok || ctx.Err() != nil || !diskScalar(row["resourceType"], &resourceType) || !storageIdentifier(resourceType) || !diskScalar(row["kind"], &kind) || !storageIdentifier(kind) || !diskScalar(row["name"], &name) || !diskName(name) {
			return bad("invalid_response")
		}
		f, ok := cognitiveRow(ctx, row, q.Region, q.LocationAliases)
		if !ok {
			return bad("invalid_response")
		}
		if !strings.EqualFold(resourceType, "accounts") || !strings.EqualFold(kind, q.Kind) {
			continue
		}
		name = strings.ToLower(name)
		old, exists := page.flags[name]
		if !exists {
			if len(page.flags) >= MaxCognitiveSKUGroups {
				return bad("row_limit")
			}
			page.flags[name] = f
			continue
		}
		for loc := range f.scopes {
			if old.scopes[loc] {
				return bad("collision")
			}
		}
		for loc := range f.scopes {
			old.scopes[loc] = true
		}
		for z := range f.zones {
			old.zones[z] = true
		}
		old.listed = old.listed || f.listed
		old.restricted = old.restricted || f.restricted
		old.unknown = old.unknown || f.unknown
		page.flags[name] = old
	}
	page.work, page.text = budget.work, budget.text
	return page, ""
}
func (c *CognitiveSKUCollector) Collect(ctx context.Context, selected map[string]string, q CognitiveSKURequest) (*CognitiveSKUCollection, error) {
	reject := func() (*CognitiveSKUCollection, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		return nil, fmt.Errorf("Cognitive SKU request rejected [cognitive_sku_request_invalid]")
	}
	if ctx.Err() != nil || c == nil || c.transport == nil {
		return reject()
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, selected, &text)
	if !ok {
		return reject()
	}
	for _, label := range scope {
		if label == "" {
			return reject()
		}
	}
	id, kind, ok := cognitiveBinding(scope, q)
	if !ok {
		return reject()
	}
	aliases, ok := cognitiveAliases(ctx, q.LocationAliases)
	if !ok {
		return reject()
	}
	q.LocationAliases = aliases
	q.Kind = kind
	first := *c.transport.origin
	first.Path = "/subscriptions/" + id + "/providers/Microsoft.CognitiveServices/skus"
	first.RawQuery = url.Values{"api-version": {cognitiveSKUVersion}}.Encode()
	result := &CognitiveSKUCollection{SubscriptionID: id, Region: q.Region, Kind: kind, ResourceType: cognitiveSKUType, Status: "complete", Groups: map[string]CognitiveSKUGroup{}}
	pages, total, rows, work, decodedText := 0, 0, 0, 0, 0
	flags := map[string]cognitiveFlags{}
	seen := map[string]bool{}
	link := first.String()
	fail := func(code string) (*CognitiveSKUCollection, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		result.Status = "unknown"
		if pages > 0 {
			result.Status = "partial"
		}
		result.FailureCode = "cognitive_sku_" + code
		return result, nil
	}
	for {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if pages >= MaxQuotaPages {
			return fail("page_limit")
		}
		target, ok := c.transport.continuation(&first, cognitiveSKUVersion, link)
		if !ok {
			return fail("continuation_invalid")
		}
		if seen[target.String()] {
			return fail("cycle")
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
			return fail("status_invalid")
		}
		if e != nil || response == nil {
			return fail("request_failed")
		}
		if len(body) > MaxQuotaPageBytes {
			return fail("byte_limit")
		}
		if len(body) > MaxQuotaTotalBytes-total {
			return fail("aggregate_limit")
		}
		page, code := decodeCognitivePage(ctx, body, q, flags, rows, work, decodedText)
		if code != "" {
			return fail(strings.TrimPrefix(code, "cognitive_sku_"))
		}
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		flags = page.flags
		result.Groups = make(map[string]CognitiveSKUGroup, len(flags))
		for name, f := range flags {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			result.Groups[name] = f.group()
		}
		rows, work, decodedText = page.rows, page.work, page.text
		total += len(body)
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
func cognitiveSKUFailure(code string) bool {
	switch code {
	case "", "cognitive_sku_request_failed", "cognitive_sku_status_invalid", "cognitive_sku_byte_limit", "cognitive_sku_aggregate_limit", "cognitive_sku_invalid_response", "cognitive_sku_collision", "cognitive_sku_row_limit", "cognitive_sku_work_limit", "cognitive_sku_text_limit", "cognitive_sku_page_limit", "cognitive_sku_continuation_invalid", "cognitive_sku_cycle":
		return true
	}
	return false
}
