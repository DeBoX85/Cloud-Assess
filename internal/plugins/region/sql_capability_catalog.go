package region

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const sqlCatalogWarning = "SQL capability catalog evidence does not establish deployment eligibility, quota or allocation capacity."

type SQLCatalogEvidence struct {
	SubscriptionID, Region, ResourceType, Status, DeploymentEligibility string
	Values                                                              map[string]SQLCapabilityGroup
	Warnings                                                            []assessment.AssessmentWarning
}

func sqlCoordinate(value string) (string, bool) {
	if len(value) > MaxSQLCoordinateBytes {
		return "", false
	}
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return "", false
	}
	for _, s := range parts {
		if !storageIdentifier(s) {
			return "", false
		}
	}
	return strings.ToLower(value), true
}
func sqlGroupStrings(g SQLCapabilityGroup) []string {
	return []string{g.Version, g.Edition, g.Name, g.Family, g.SKUName, g.SKUTier, g.Standalone, g.InstancePool}
}
func sqlValidDatabaseGroup(key string, g SQLCapabilityGroup) bool {
	normalized, valid := sqlCoordinate(key)
	if !valid || normalized != key || key != g.Version+"/"+g.Edition+"/"+g.Name {
		return false
	}
	for _, s := range []string{g.Version, g.Edition, g.Name} {
		if !storageIdentifier(s) || strings.ToLower(s) != s {
			return false
		}
	}
	for _, s := range []string{g.Family, g.SKUName, g.SKUTier} {
		if s != "" && (!storageIdentifier(s) || strings.ToLower(s) != s) {
			return false
		}
	}
	if g.State != SKUAvailable && g.State != SKUUnavailable || !g.Known && g.State != SKUUnavailable || g.Capacity < 0 || g.Capacity > 2147483647 {
		return false
	}
	if g.SKUTier != "" && g.SKUTier != g.Edition || g.Standalone != "not_applicable" || g.InstancePool != "not_applicable" {
		return false
	}
	if g.Known && (g.SKUName == "" || g.SKUTier == "" || g.Capacity <= 0) {
		return false
	}
	// A wholly missing SKU cannot contribute tier/capacity/family by inference.
	if g.SKUName == "" && (g.SKUTier != "" || g.Capacity != 0 || g.Family != "") {
		return false
	}
	return true
}

// Catalog performs only structural projection on stable inputs with a nonnil ctx.
// Caller-built evidence is not authenticated service provenance or deployment eligibility.
func (c *SQLCapabilityCollection) Catalog(ctx context.Context, selected map[string]string, expected SQLCapabilityRequest, required []string) (SQLCatalogEvidence, error) {
	fail := func(code string) (SQLCatalogEvidence, error) {
		if err := ctx.Err(); err != nil {
			return SQLCatalogEvidence{}, err
		}
		return SQLCatalogEvidence{}, fmt.Errorf("SQL capability catalog projection rejected [sql_capability_projection_%s]", code)
	}
	if ctx.Err() != nil {
		return fail("invalid")
	}
	if len(selected) > MaxSubscriptions || len(expected.LocationAliases) > MaxSQLCapabilityAliases || len(required) > MaxSQLCapabilityLeaves || c != nil && len(c.Groups) > MaxSQLCapabilityLeaves {
		return fail("limit")
	}
	budget := vmSKUBudget{ctx: ctx}
	for id, label := range selected {
		if !budget.entry() || !budget.label(id, MaxLabelBytes) || !budget.label(label, MaxLabelBytes) {
			return fail("limit")
		}
		if label == "" {
			return fail("invalid")
		}
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, selected, &text)
	if !ok || len(scope) == 0 {
		return fail("invalid")
	}
	for _, s := range []string{expected.SubscriptionID, expected.Region, expected.ResourceType} {
		if !budget.label(s, MaxLabelBytes) {
			return fail("limit")
		}
	}
	id, kind, ok := sqlCapabilityBinding(scope, expected)
	if !ok {
		return fail("invalid")
	}
	for key, value := range expected.LocationAliases {
		if !budget.entry() || !budget.label(key, MaxLabelBytes) || !budget.label(value, MaxLabelBytes) {
			return fail("limit")
		}
	}
	if _, valid := cognitiveAliases(ctx, expected.LocationAliases); !valid {
		return fail("invalid")
	}
	if c != nil {
		for _, s := range []string{c.SubscriptionID, c.Region, c.ResourceType, c.Status, c.FailureCode} {
			if !budget.label(s, MaxLabelBytes) {
				return fail("limit")
			}
		}
		if c.SubscriptionID != id || c.Region != expected.Region || c.ResourceType != kind {
			return fail("invalid")
		}
		if c.Status != "complete" && c.Status != "unknown" || c.Status == "complete" && c.FailureCode != "" || c.Status == "unknown" && (!sqlCapabilityFailure(c.FailureCode) || len(c.Groups) != 0) {
			return fail("invalid")
		}
		for key, g := range c.Groups {
			if !budget.entry() || !budget.label(key, MaxSQLCoordinateBytes) {
				return fail("limit")
			}
			for _, s := range sqlGroupStrings(g) {
				if !budget.label(s, MaxLabelBytes) {
					return fail("limit")
				}
			}
			if (kind == sqlDatabaseType && !sqlValidDatabaseGroup(key, g)) || (kind == sqlManagedInstanceType && !sqlValidManagedInstanceGroup(key, g)) {
				return fail("invalid")
			}
		}
	}
	names := make(map[string]bool, len(required))
	for _, s := range required {
		if !budget.entry() || !budget.label(s, MaxSQLCoordinateBytes) {
			return fail("limit")
		}
		key, valid := sqlCapabilityCoordinate(s, kind)
		if !valid || names[key] {
			return fail("invalid")
		}
		names[key] = true
	}
	newResult := func(status string) SQLCatalogEvidence {
		return SQLCatalogEvidence{SubscriptionID: id, Region: expected.Region, ResourceType: kind, Status: status, DeploymentEligibility: "not_assessed", Values: map[string]SQLCapabilityGroup{}, Warnings: []assessment.AssessmentWarning{{Code: "sql_capability_catalog_only", Message: sqlCatalogWarning}}}
	}
	finish := func(out SQLCatalogEvidence) (SQLCatalogEvidence, error) {
		if err := ctx.Err(); err != nil {
			return SQLCatalogEvidence{}, err
		}
		return out, nil
	}
	if c == nil || c.Status != "complete" {
		return finish(newResult("unknown"))
	}
	for key := range names {
		if ctx.Err() != nil {
			return fail("invalid")
		}
		g, found := c.Groups[key]
		if !found || !g.Known {
			return finish(newResult("unknown"))
		}
	}
	output := vmSKUBudget{ctx: ctx}
	out := newResult("complete")
	for _, s := range []string{out.SubscriptionID, out.Region, out.ResourceType, out.Status, out.DeploymentEligibility, out.Warnings[0].Code, out.Warnings[0].Message} {
		if !output.label(s, MaxLabelBytes) {
			return fail("limit")
		}
	}
	for key := range names {
		if ctx.Err() != nil {
			return fail("invalid")
		}
		g := c.Groups[key]
		if !output.label(key, MaxSQLCoordinateBytes) {
			return fail("limit")
		}
		for _, s := range sqlGroupStrings(g) {
			if !output.label(s, MaxLabelBytes) {
				return fail("limit")
			}
		}
		out.Values[key] = g
	}
	return finish(out)
}
