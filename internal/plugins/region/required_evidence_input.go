package region

import (
	"context"
	"strings"
)

// The admission traversal charges every supplied row before filtering. A code
// is private control flow, never a provider diagnostic or an input fragment.
type requiredInput struct {
	ctx    context.Context
	budget requiredBudget
	code   string
}

func (v *requiredInput) count(n int) bool {
	if v.ctx.Err() != nil {
		return false
	}
	if n < 0 || n > MaxInventoryEntries-v.budget.entries {
		v.code = "required_evidence_limit"
		return false
	}
	return true
}
func (v *requiredInput) row() bool {
	if !v.count(1) {
		return false
	}
	v.budget.entries++
	return true
}
func (v *requiredInput) text(fields ...string) bool {
	for _, s := range fields {
		if v.ctx.Err() != nil {
			return false
		}
		if !requiredScalar(s) {
			v.code = "required_evidence_invalid"
			return false
		}
		if !v.budget.add(0, s) {
			v.code = "required_evidence_limit"
			return false
		}
	}
	return true
}
func (v *requiredInput) invalid() bool { v.code = "required_evidence_invalid"; return false }
func (v *requiredInput) coordinate(s, kind string) bool {
	if v.ctx.Err() != nil {
		return false
	}
	if _, ok := sqlCapabilityCoordinate(s, kind); !ok {
		return v.invalid()
	}
	if !v.budget.add(0, s) {
		v.code = "required_evidence_limit"
		return false
	}
	return true
}
func (v *requiredInput) requirement(r RequiredInventoryRequirement) bool {
	if !v.row() || !v.text(r.SubscriptionID, r.SourceRegion, r.ResourceType, r.Kind) {
		return false
	}
	if !v.count(len(r.SKUs)) {
		return false
	}
	for _, s := range r.SKUs {
		if !v.row() || !v.text(s) {
			return false
		}
	}
	if !v.count(len(r.Kinds)) {
		return false
	}
	for _, s := range requiredKeys(r.Kinds) {
		if !v.row() || !v.text(s) || !v.count(len(r.Kinds[s])) {
			return false
		}
		for _, k := range r.Kinds[s] {
			if !v.row() || !v.text(k) {
				return false
			}
		}
	}
	if !v.count(len(r.ResourceIDs)) {
		return false
	}
	for _, s := range r.ResourceIDs {
		if !v.row() || !v.text(s) {
			return false
		}
	}
	return true
}
func requiredDiskFailure(s string) bool {
	switch s {
	case "sku_byte_limit", "sku_continuation_cycle", "sku_identity_duplicate", "sku_identity_invalid", "sku_invalid_response", "sku_page_limit", "sku_request_failed", "sku_status_invalid", "sku_text_limit", "sku_unsafe_continuation", "sku_work_limit":
		return true
	}
	return false
}
func requiredStatus(status, failure string, rows int, valid func(string) bool, partial bool) bool {
	return (status == "complete" && failure == "") || (status == "partial" && partial && failure != "" && valid(failure)) || (status == "unknown" && valid(failure) && rows == 0)
}
func (v *requiredInput) zones(state int, known bool, zones []string) bool {
	if state < SKUAvailable || state > SKUZoneRestricted || !known && state != SKUUnavailable || state != SKUZoneRestricted && len(zones) != 0 || state == SKUZoneRestricted && len(zones) == 0 {
		return v.invalid()
	}
	if !v.count(len(zones)) {
		return false
	}
	prior := ""
	for _, z := range zones {
		if !v.row() || !v.text(z) {
			return false
		}
		if !diskZone(z) || z <= prior {
			return v.invalid()
		}
		prior = z
	}
	return true
}
func (v *requiredInput) metadata(sub, region string, q RequiredEvidenceRequest) bool {
	if sub != q.SubscriptionID || region != q.TargetRegion {
		return v.invalid()
	}
	return true
}
func (v *requiredInput) collection(c RequiredCollectionInput, q RequiredEvidenceRequest) bool {
	// Union discrimination precedes any payload traversal.
	n := 0
	for _, has := range []bool{c.Storage != nil, c.Disk != nil, c.VM != nil, c.Cognitive != nil, c.SQL != nil} {
		if has {
			n++
		}
	}
	if n != 1 {
		return v.invalid()
	}
	switch c.ResourceType {
	case "microsoft.storage/storageaccounts":
		p := c.Storage
		if p == nil {
			return v.invalid()
		}
		if len(p.Groups) > MaxStorageSKURows || len(p.UnqualifiedNames) > MaxStorageSKURows {
			v.code = "required_evidence_limit"
			return false
		}
		if !v.row() || !v.text(p.SubscriptionID, p.Region, p.Status, p.FailureCode) || !v.metadata(p.SubscriptionID, p.Region, q) {
			return false
		}
		if !requiredStatus(p.Status, p.FailureCode, len(p.Groups)+len(p.UnqualifiedNames), storageFailureCode, false) {
			return v.invalid()
		}
		if !v.count(len(p.Groups)) {
			return false
		}
		total := 0
		for _, name := range requiredKeys(p.Groups) {
			kinds := p.Groups[name]
			if !v.row() || !v.text(name) {
				return false
			}
			if !diskName(name) || strings.ToLower(name) != name || len(kinds) == 0 {
				return v.invalid()
			}
			if len(kinds) > MaxStorageSKURows-total {
				v.code = "required_evidence_limit"
				return false
			}
			total += len(kinds)
			if !v.count(len(kinds)) {
				return false
			}
			for _, kind := range requiredKeys(kinds) {
				g := kinds[kind]
				if !v.row() || !v.text(kind) {
					return false
				}
				if !storageIdentifier(kind) || strings.ToLower(kind) != kind || g.State < 0 || g.State > SKUUnavailable || !g.Known && g.State != SKUUnavailable {
					return v.invalid()
				}
			}
		}
		if !v.count(len(p.UnqualifiedNames)) {
			return false
		}
		for _, name := range requiredKeys(p.UnqualifiedNames) {
			if !v.row() || !v.text(name) {
				return false
			}
			if !p.UnqualifiedNames[name] || !diskName(name) || strings.ToLower(name) != name {
				return v.invalid()
			}
		}
	case "microsoft.compute/disks":
		p := c.Disk
		if p == nil {
			return v.invalid()
		}
		if len(p.Values) > MaxInventoryEntries {
			v.code = "required_evidence_limit"
			return false
		}
		if !v.row() || !v.text(p.SubscriptionID, p.Region, p.Status, p.FailureCode) || !v.metadata(p.SubscriptionID, p.Region, q) {
			return false
		}
		if !requiredStatus(p.Status, p.FailureCode, len(p.Values), requiredDiskFailure, true) {
			return v.invalid()
		}
		if !v.count(len(p.Values)) {
			return false
		}
		for _, name := range requiredKeys(p.Values) {
			g := p.Values[name]
			if !v.row() || !v.text(name) {
				return false
			}
			if !diskName(name) || strings.ToLower(name) != name {
				return v.invalid()
			}
			if !v.zones(g.State, true, g.BlockedZones) {
				return false
			}
		}
	case vmSKUType, vmSSSKUType:
		p := c.VM
		if p == nil {
			return v.invalid()
		}
		if len(p.Groups) > MaxComputeVMSKUGroups {
			v.code = "required_evidence_limit"
			return false
		}
		if !v.row() || !v.text(p.SubscriptionID, p.Region, p.ResourceType, p.Status, p.FailureCode) || !v.metadata(p.SubscriptionID, p.Region, q) {
			return false
		}
		if p.ResourceType != c.ResourceType || !requiredStatus(p.Status, p.FailureCode, len(p.Groups), vmSKUFailure, true) {
			return v.invalid()
		}
		if !v.count(len(p.Groups)) {
			return false
		}
		for _, name := range requiredKeys(p.Groups) {
			g := p.Groups[name]
			if !v.row() || !v.text(name) {
				return false
			}
			if !diskName(name) || strings.ToLower(name) != name {
				return v.invalid()
			}
			if !v.zones(g.State, g.Known, g.BlockedZones) {
				return false
			}
		}
	case cognitiveSKUType:
		p := c.Cognitive
		if p == nil {
			return v.invalid()
		}
		if len(p.Groups) > MaxCognitiveSKUGroups {
			v.code = "required_evidence_limit"
			return false
		}
		if !v.row() || !v.text(p.SubscriptionID, p.Region, p.Kind, p.ResourceType, p.Status, p.FailureCode) || !v.metadata(p.SubscriptionID, p.Region, q) {
			return false
		}
		if p.ResourceType != c.ResourceType || p.Kind != c.Kind || !storageIdentifier(p.Kind) || !requiredStatus(p.Status, p.FailureCode, len(p.Groups), cognitiveSKUFailure, true) {
			return v.invalid()
		}
		if !v.count(len(p.Groups)) {
			return false
		}
		for _, name := range requiredKeys(p.Groups) {
			g := p.Groups[name]
			if !v.row() || !v.text(name) {
				return false
			}
			if !diskName(name) || strings.ToLower(name) != name {
				return v.invalid()
			}
			if !v.zones(g.State, g.Known, g.BlockedZones) {
				return false
			}
		}
	case sqlDatabaseType, sqlManagedInstanceType:
		p := c.SQL
		if p == nil {
			return v.invalid()
		}
		if len(p.Groups) > MaxSQLCapabilityLeaves {
			v.code = "required_evidence_limit"
			return false
		}
		if !v.row() || !v.text(p.SubscriptionID, p.Region, p.ResourceType, p.Status, p.FailureCode) || !v.metadata(p.SubscriptionID, p.Region, q) {
			return false
		}
		if p.ResourceType != c.ResourceType || !requiredStatus(p.Status, p.FailureCode, len(p.Groups), sqlCapabilityFailure, false) {
			return v.invalid()
		}
		if !v.count(len(p.Groups)) {
			return false
		}
		for _, name := range requiredKeys(p.Groups) {
			g := p.Groups[name]
			if !v.row() || !v.coordinate(name, c.ResourceType) || !v.text(sqlGroupStrings(g)...) {
				return false
			}
			if c.ResourceType == sqlDatabaseType && !sqlValidDatabaseGroup(name, g) || c.ResourceType == sqlManagedInstanceType && !sqlValidManagedInstanceGroup(name, g) {
				return v.invalid()
			}
		}
	default:
		return v.invalid()
	}
	return true
}
