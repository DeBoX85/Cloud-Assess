package region

import (
	"context"
	"errors"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func requiredAdmissionError(ctx context.Context, observe func(string), code string) (*RequiredEvidenceAdmission, error) {
	if observe != nil {
		observe("before-return")
	}
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, errors.New(code)
}

// AdmitRequiredEvidence performs bounded structural admission only. Caller data
// is not authenticated service provenance. No transport or discovery is used.
func AdmitRequiredEvidence(ctx context.Context, plan *RequiredInventoryPlan, q RequiredEvidenceRequest, collections []RequiredCollectionInput) (*RequiredEvidenceAdmission, error) {
	return admitRequiredEvidenceObserved(ctx, plan, q, collections, nil)
}
func admitRequiredEvidenceObserved(ctx context.Context, plan *RequiredInventoryPlan, q RequiredEvidenceRequest, collections []RequiredCollectionInput, observe func(string)) (*RequiredEvidenceAdmission, error) {
	fail := func(code string) (*RequiredEvidenceAdmission, error) {
		return requiredAdmissionError(ctx, observe, code)
	}
	invalid := func() (*RequiredEvidenceAdmission, error) { return fail("required_evidence_invalid") }
	limit := func() (*RequiredEvidenceAdmission, error) { return fail("required_evidence_limit") }
	if ctx == nil {
		return invalid()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(collections) > MaxInventoryEntries {
		return limit()
	}
	if plan == nil || plan.state == nil {
		return invalid()
	}
	for _, s := range []string{q.SubscriptionID, q.SourceRegion, q.TargetRegion} {
		if !requiredScalar(s) {
			return invalid()
		}
	}
	q.SubscriptionID = strings.ToLower(q.SubscriptionID)
	source, sourceOK := diskPhysical(q.SourceRegion)
	target, targetOK := diskPhysical(q.TargetRegion)
	scope, selected := plan.state.scopes[q.SubscriptionID]
	if len(q.SubscriptionID) != 36 || !subscriptionID.MatchString(q.SubscriptionID) || !selected || !sourceOK || source != q.SourceRegion || !targetOK || target != q.TargetRegion {
		return invalid()
	}
	input := requiredInput{ctx: ctx, code: "required_evidence_invalid"}
	if !input.row() || !input.text(q.SubscriptionID, scope.Name, scope.Status) || !input.row() || !input.text(q.SubscriptionID, q.SourceRegion, q.TargetRegion) {
		return fail(input.code)
	}
	requirements := []RequiredInventoryRequirement{}
	required := map[[4]string]bool{}
	for _, r := range plan.state.requirements {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if r.SubscriptionID != q.SubscriptionID || r.SourceRegion != q.SourceRegion {
			continue
		}
		if !input.requirement(r) {
			return fail(input.code)
		}
		requirements = append(requirements, r)
		required[requiredTuple(r)] = true
	}
	if !input.count(len(collections)) {
		return fail(input.code)
	}
	bindings := make(map[[4]string]RequiredCollectionInput, len(collections))
	for _, raw := range collections {
		if !input.row() || !input.text(raw.SubscriptionID, raw.SourceRegion, raw.TargetRegion, raw.ResourceType, raw.Kind) {
			return fail(input.code)
		}
		if !providerTypeIdentifier(raw.ResourceType) {
			return invalid()
		}
		c := raw
		c.SubscriptionID = strings.ToLower(c.SubscriptionID)
		c.ResourceType = strings.ToLower(c.ResourceType)
		if !storageASCII(c.Kind) {
			return invalid()
		}
		c.Kind = strings.ToLower(c.Kind)
		tuple := [4]string{c.SubscriptionID, c.SourceRegion, c.ResourceType, c.Kind}
		_, duplicate := bindings[tuple]
		if c.SubscriptionID != q.SubscriptionID || c.SourceRegion != q.SourceRegion || c.TargetRegion != q.TargetRegion || !required[tuple] || duplicate || c.ResourceType != cognitiveSKUType && c.Kind != "" {
			return invalid()
		}
		if !input.collection(c, q) {
			return fail(input.code)
		}
		bindings[tuple] = c
	}
	if err := requiredPhase(ctx, observe, "after-input"); err != nil {
		return nil, err
	}
	// Projectors produce bounded temporary values. Reserve the entire owned
	// result and accessor view before adopting them as public/private evidence.
	result := &RequiredEvidenceAdmission{Request: q, ScopeStatus: scope.Status, Status: scope.Status, SKUs: map[string]SKUEvidence{}, Catalogs: []RequiredCatalogBinding{}}
	partial := false
	selectedScope := map[string]string{q.SubscriptionID: scope.Name}
	for _, r := range requirements {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		c, present := bindings[requiredTuple(r)]
		e := SKUEvidence{Status: "unsupported", Values: map[string]SKUAvailability{}}
		badCoordinates := r.MissingSKU || r.UnrepresentableSKU || r.MissingKind || r.UnrepresentableKind
		var err error
		catalogStatus := ""
		switch r.ResourceType {
		case "microsoft.storage/storageaccounts":
			e, err = c.Storage.Evidence(ctx, r.Kinds)
		case vmSKUType:
			e, err = c.VM.Evidence(ctx, r.SKUs)
		case "microsoft.compute/disks":
			e = c.Disk.Evidence()
			if e.Status == "complete" {
				values := map[string]SKUAvailability{}
				for _, key := range r.SKUs {
					if err := ctx.Err(); err != nil {
						return nil, err
					}
					g, found := e.Values[key]
					if !found {
						g = SKUAvailability{State: SKUUnavailable, BlockedZones: []string{}}
					}
					values[key] = g
				}
				e.Values = values
			}
		case cognitiveSKUType:
			var out CognitiveCatalogEvidence
			if !storageIdentifier(r.Kind) {
				out = CognitiveCatalogEvidence{SubscriptionID: q.SubscriptionID, Region: q.TargetRegion, Kind: r.Kind, ResourceType: r.ResourceType, Status: "unknown", DeploymentEligibility: "not_assessed", Values: map[string]SKUAvailability{}, Warnings: []assessment.AssessmentWarning{{Code: "cognitive_sku_catalog_only", Message: cognitiveCatalogWarning}}}
			} else {
				out, err = c.Cognitive.Catalog(ctx, selectedScope, CognitiveSKURequest{SubscriptionID: q.SubscriptionID, Region: q.TargetRegion, Kind: r.Kind}, r.SKUs)
			}
			if badCoordinates {
				out.Status = "unknown"
				out.Values = map[string]SKUAvailability{}
			}
			catalogStatus = out.Status
			result.Catalogs = append(result.Catalogs, RequiredCatalogBinding{SubscriptionID: q.SubscriptionID, SourceRegion: q.SourceRegion, TargetRegion: q.TargetRegion, ResourceType: r.ResourceType, Kind: r.Kind, Status: out.Status, DeploymentEligibility: out.DeploymentEligibility, Cognitive: &out})
		case vmSSSKUType:
			var out VMSSCatalogEvidence
			out, err = c.VM.VMSSCatalog(ctx, selectedScope, ComputeVMSKURequest{SubscriptionID: q.SubscriptionID, Region: q.TargetRegion, ResourceType: r.ResourceType}, r.SKUs)
			if badCoordinates {
				out.Status = "unknown"
				out.Values = map[string]SKUAvailability{}
			}
			catalogStatus = out.Status
			result.Catalogs = append(result.Catalogs, RequiredCatalogBinding{SubscriptionID: q.SubscriptionID, SourceRegion: q.SourceRegion, TargetRegion: q.TargetRegion, ResourceType: r.ResourceType, Status: out.Status, DeploymentEligibility: out.DeploymentEligibility, VMSS: &out})
		case sqlDatabaseType, sqlManagedInstanceType:
			keys := []string{}
			if c.SQL != nil {
				keys = requiredKeys(c.SQL.Groups)
			}
			var out SQLCatalogEvidence
			out, err = c.SQL.Catalog(ctx, selectedScope, SQLCapabilityRequest{SubscriptionID: q.SubscriptionID, Region: q.TargetRegion, ResourceType: r.ResourceType}, keys)
			catalogStatus = out.Status
			result.Catalogs = append(result.Catalogs, RequiredCatalogBinding{SubscriptionID: q.SubscriptionID, SourceRegion: q.SourceRegion, TargetRegion: q.TargetRegion, ResourceType: r.ResourceType, Status: out.Status, DeploymentEligibility: out.DeploymentEligibility, SQL: &out})
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return fail(requiredProjectionCode(err))
		}
		if !r.CatalogOnly && !r.Unsupported && (badCoordinates || scope.Status != "complete") {
			e = SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}}
		}
		result.SKUs[r.ResourceType] = e
		if !r.Unsupported && (!present || badCoordinates || e.Status == "unknown" || catalogStatus == "unknown") {
			partial = true
		}
	}
	if scope.Status == "complete" && partial {
		result.Status = "partial"
	}
	if err := requiredPhase(ctx, observe, "after-derived-work"); err != nil {
		return nil, err
	}
	if !requiredAdmissionOutput(ctx, result) {
		return limit()
	}
	result.request = q
	result.values = requiredCopySKUs(result.SKUs)
	if err := requiredPhase(ctx, observe, "before-return"); err != nil {
		return nil, err
	}
	return result, nil
}

// Exact static local formats only. Unknown errors stay invalid, never pass through.
func requiredProjectionCode(err error) string {
	switch err.Error() {
	case "Storage SKU projection rejected [storage_sku_projection_limit]", "Compute VM SKU projection rejected [vm_sku_projection_limit]", "Cognitive SKU catalog projection rejected [cognitive_sku_projection_limit]", "VMSS catalog projection rejected [vmss_catalog_projection_limit]", "SQL capability catalog projection rejected [sql_capability_projection_limit]":
		return "required_evidence_limit"
	}
	return "required_evidence_invalid"
}
func requiredAdmissionOutput(ctx context.Context, a *RequiredEvidenceAdmission) bool {
	b := requiredBudget{}
	if !b.add(0, a.ScopeStatus, a.Status) {
		return false
	}
	for i := 0; i < 3; i++ {
		if !b.add(1, a.Request.SubscriptionID, a.Request.SourceRegion, a.Request.TargetRegion) {
			return false
		}
		if len(a.SKUs) > MaxInventoryEntries-b.entries {
			return false
		}
		for _, typ := range requiredKeys(a.SKUs) {
			if ctx.Err() != nil {
				return false
			}
			e := a.SKUs[typ]
			if !b.add(1, typ, e.Status) {
				return false
			}
			if len(e.Values) > MaxInventoryEntries-b.entries {
				return false
			}
			for _, key := range requiredKeys(e.Values) {
				g := e.Values[key]
				if !b.add(1, key) {
					return false
				}
				if len(g.BlockedZones) > MaxInventoryEntries-b.entries {
					return false
				}
				for _, z := range g.BlockedZones {
					if ctx.Err() != nil || !b.add(1, z) {
						return false
					}
				}
			}
		}
	}
	if len(a.Catalogs) > MaxInventoryEntries-b.entries {
		return false
	}
	for _, c := range a.Catalogs {
		if ctx.Err() != nil || !b.add(1, c.SubscriptionID, c.SourceRegion, c.TargetRegion, c.ResourceType, c.Kind, c.Status, c.DeploymentEligibility) {
			return false
		}
		var values map[string]SKUAvailability
		var warnings []assessment.AssessmentWarning
		if p := c.Cognitive; p != nil {
			if !b.add(1, p.SubscriptionID, p.Region, p.Kind, p.ResourceType, p.Status, p.DeploymentEligibility) {
				return false
			}
			values = p.Values
			warnings = p.Warnings
		}
		if p := c.VMSS; p != nil {
			if !b.add(1, p.SubscriptionID, p.Region, p.TargetResourceType, p.SourceResourceType, p.Status, p.DeploymentEligibility) {
				return false
			}
			values = p.Values
			warnings = p.Warnings
		}
		if p := c.SQL; p != nil {
			if !b.add(1, p.SubscriptionID, p.Region, p.ResourceType, p.Status, p.DeploymentEligibility) {
				return false
			}
			if len(p.Values) > MaxInventoryEntries-b.entries {
				return false
			}
			for _, key := range requiredKeys(p.Values) {
				if ctx.Err() != nil || !b.add(1, key) || !b.add(0, sqlGroupStrings(p.Values[key])...) {
					return false
				}
			}
			warnings = p.Warnings
		}
		if len(values) > MaxInventoryEntries-b.entries {
			return false
		}
		for _, key := range requiredKeys(values) {
			if ctx.Err() != nil || !b.add(1, key) {
				return false
			}
			if len(values[key].BlockedZones) > MaxInventoryEntries-b.entries {
				return false
			}
			for _, z := range values[key].BlockedZones {
				if ctx.Err() != nil || !b.add(1, z) {
					return false
				}
			}
		}
		if len(warnings) > MaxInventoryEntries-b.entries {
			return false
		}
		for _, w := range warnings {
			if ctx.Err() != nil || !b.add(1, w.Code, w.Message) {
				return false
			}
		}
	}
	return true
}
