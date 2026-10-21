package region

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const vmssCatalogWarning = "VMSS catalog evidence does not establish workload deployment eligibility, quota or allocation capacity."

// VMSSCatalogEvidence describes node catalog restrictions only. It is deliberately
// distinct from SKUEvidence: it cannot establish workload placement or eligibility.
type VMSSCatalogEvidence struct {
	SubscriptionID, Region, TargetResourceType, SourceResourceType string
	Status, DeploymentEligibility                                  string
	Values                                                         map[string]SKUAvailability
	Warnings                                                       []assessment.AssessmentWarning
}

// VMSSCatalog projects existing, structurally validated collection data without
// transport or credentials. Inputs must remain stable and ctx must be non-nil.
// A structurally valid caller-built collection is not authenticated provenance.
func (c *ComputeVMSKUCollection) VMSSCatalog(ctx context.Context, selected map[string]string, expected ComputeVMSKURequest, required []string) (VMSSCatalogEvidence, error) {
	fail := func(code string) (VMSSCatalogEvidence, error) {
		if err := ctx.Err(); err != nil {
			return VMSSCatalogEvidence{}, err
		}
		return VMSSCatalogEvidence{}, fmt.Errorf("VMSS catalog projection rejected [%s]", code)
	}
	if ctx.Err() != nil {
		return fail("vmss_catalog_projection_invalid")
	}
	if len(selected) > MaxSubscriptions || len(required) > MaxComputeVMSKUGroups || c != nil && len(c.Groups) > MaxComputeVMSKUGroups {
		return fail("vmss_catalog_projection_limit")
	}
	budget := vmSKUBudget{ctx: ctx}
	// Unlike legacy VM Evidence, this boundary explicitly charges selection and
	// expected binding. Do not change that API's accounting to share this policy.
	for id, label := range selected {
		if !budget.entry() || !budget.label(id, MaxLabelBytes) || !budget.label(label, MaxLabelBytes) {
			return fail("vmss_catalog_projection_limit")
		}
		if label == "" {
			return fail("vmss_catalog_projection_invalid")
		}
	}
	scopeText := 0
	scope, ok := auxiliaryScope(ctx, selected, &scopeText)
	if !ok || len(scope) == 0 {
		return fail("vmss_catalog_projection_invalid")
	}
	for _, s := range []string{expected.SubscriptionID, expected.Region, expected.ResourceType} {
		if !budget.label(s, MaxLabelBytes) {
			return fail("vmss_catalog_projection_limit")
		}
	}
	id := strings.ToLower(expected.SubscriptionID)
	region, validRegion := diskPhysical(expected.Region)
	if len(id) != 36 || !subscriptionID.MatchString(id) || scope[id] == "" || !validRegion || region != expected.Region || expected.ResourceType != vmSSSKUType {
		return fail("vmss_catalog_projection_invalid")
	}
	if c != nil {
		for _, s := range []string{c.SubscriptionID, c.Region, c.ResourceType, c.Status, c.FailureCode} {
			if !budget.label(s, MaxLabelBytes) {
				return fail("vmss_catalog_projection_limit")
			}
		}
		if c.SubscriptionID != id || c.Region != region || c.ResourceType != vmSSSKUType {
			return fail("vmss_catalog_projection_invalid")
		}
		if (c.Status != "complete" && c.Status != "partial" && c.Status != "unknown") || !vmSKUFailure(c.FailureCode) || c.Status == "complete" && c.FailureCode != "" || c.Status == "partial" && c.FailureCode == "" || c.Status == "unknown" && len(c.Groups) != 0 {
			return fail("vmss_catalog_projection_invalid")
		}
		// Validate and charge all groups, including unrequired diagnostic prefixes.
		for name, g := range c.Groups {
			if !budget.entry() || !budget.label(name, MaxLabelBytes) {
				return fail("vmss_catalog_projection_limit")
			}
			if !diskName(name) || strings.ToLower(name) != name || g.State < SKUAvailable || g.State > SKUZoneRestricted || !g.Known && g.State != SKUUnavailable || g.State != SKUZoneRestricted && len(g.BlockedZones) != 0 || g.State == SKUZoneRestricted && len(g.BlockedZones) == 0 {
				return fail("vmss_catalog_projection_invalid")
			}
			prior := ""
			for _, zone := range g.BlockedZones {
				if !budget.entry() || !budget.label(zone, MaxLabelBytes) {
					return fail("vmss_catalog_projection_limit")
				}
				if !diskZone(zone) || zone <= prior {
					return fail("vmss_catalog_projection_invalid")
				}
				prior = zone
			}
		}
	}
	names := make(map[string]bool, len(required))
	for _, s := range required {
		if !budget.entry() || !budget.label(s, MaxLabelBytes) {
			return fail("vmss_catalog_projection_limit")
		}
		name := strings.ToLower(s)
		if !diskName(s) || names[name] {
			return fail("vmss_catalog_projection_invalid")
		}
		names[name] = true
	}
	// Every input has been validated and charged before result allocation.
	newResult := func(status string) VMSSCatalogEvidence {
		return VMSSCatalogEvidence{SubscriptionID: id, Region: region, TargetResourceType: vmSSSKUType, SourceResourceType: vmSKUType, Status: status, DeploymentEligibility: "not_assessed", Values: map[string]SKUAvailability{}, Warnings: []assessment.AssessmentWarning{{Code: "vmss_catalog_only", Message: vmssCatalogWarning}}}
	}
	finish := func(out VMSSCatalogEvidence) (VMSSCatalogEvidence, error) {
		if err := ctx.Err(); err != nil {
			return VMSSCatalogEvidence{}, err
		}
		return out, nil
	}
	if c == nil || c.Status != "complete" {
		return finish(newResult("unknown"))
	}
	// Resolve whole-result uncertainty before allocating a partial output. The
	// answer must not depend on map iteration order or an unused output's size.
	for name := range names {
		if ctx.Err() != nil {
			return fail("vmss_catalog_projection_invalid")
		}
		if g, found := c.Groups[name]; found && !g.Known {
			return finish(newResult("unknown"))
		}
	}
	out := newResult("complete")
	output := vmSKUBudget{ctx: ctx}
	for _, s := range []string{out.SubscriptionID, out.Region, out.TargetResourceType, out.SourceResourceType, out.Status, out.DeploymentEligibility, out.Warnings[0].Code, out.Warnings[0].Message} {
		if !output.label(s, MaxLabelBytes) {
			return fail("vmss_catalog_projection_limit")
		}
	}
	for name := range names {
		if ctx.Err() != nil {
			return fail("vmss_catalog_projection_invalid")
		}
		g, found := c.Groups[name]
		if !found {
			g = ComputeVMSKUGroup{Known: true, State: SKUUnavailable}
		}
		if !output.label(name, MaxLabelBytes) {
			return fail("vmss_catalog_projection_limit")
		}
		zones := make([]string, 0, len(g.BlockedZones))
		for _, zone := range g.BlockedZones {
			if !output.label(zone, MaxLabelBytes) {
				return fail("vmss_catalog_projection_limit")
			}
			zones = append(zones, zone)
		}
		out.Values[name] = SKUAvailability{State: g.State, BlockedZones: zones}
	}
	return finish(out)
}
