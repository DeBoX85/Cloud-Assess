package region

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// Required inventory is private coordination evidence, not Azure provenance.
type InventoryScopeDeclaration struct{ SubscriptionID, Status string }
type RequiredInventoryScope struct{ Name, Status string }
type RequiredInventoryResource struct {
	Record                                           assessment.Resource
	SubscriptionID, SourceRegion, ResourceType, Kind string
	Physical                                         bool
}
type RequiredInventoryRequirement struct {
	SubscriptionID, SourceRegion, ResourceType, Kind                 string
	SKUs                                                             []string
	Kinds                                                            map[string][]string
	ResourceIDs                                                      []string
	MissingSKU, MissingKind, UnrepresentableSKU, UnrepresentableKind bool
	CatalogOnly, Unsupported                                         bool
}
type requiredInventoryState struct {
	resources    []RequiredInventoryResource
	requirements []RequiredInventoryRequirement
	scopes       map[string]RequiredInventoryScope
}
type RequiredInventoryPlan struct {
	Resources    []RequiredInventoryResource
	Requirements []RequiredInventoryRequirement
	Scopes       map[string]RequiredInventoryScope
	Status       string
	state        *requiredInventoryState
}
type RequiredEvidenceRequest struct{ SubscriptionID, SourceRegion, TargetRegion string }
type RequiredCollectionInput struct {
	SubscriptionID, SourceRegion, TargetRegion, ResourceType, Kind string
	Storage                                                        *StorageSKUCollection
	Disk                                                           *ComputeDiskSKUCollection
	VM                                                             *ComputeVMSKUCollection
	Cognitive                                                      *CognitiveSKUCollection
	SQL                                                            *SQLCapabilityCollection
}
type RequiredCatalogBinding struct {
	SubscriptionID, SourceRegion, TargetRegion, ResourceType, Kind string
	Status, DeploymentEligibility                                  string
	Cognitive                                                      *CognitiveCatalogEvidence
	SQL                                                            *SQLCatalogEvidence
	VMSS                                                           *VMSSCatalogEvidence
}
type RequiredEvidenceAdmission struct {
	Request             RequiredEvidenceRequest
	ScopeStatus, Status string
	SKUs                map[string]SKUEvidence
	Catalogs            []RequiredCatalogBinding
	request             RequiredEvidenceRequest
	values              map[string]SKUEvidence
}

// Evidence owns every mutable value independently for each call, including nil
// and zero admissions. The admission constructor will establish private state.
func (a *RequiredEvidenceAdmission) Evidence() (RequiredEvidenceRequest, map[string]SKUEvidence) {
	if a == nil {
		return RequiredEvidenceRequest{}, map[string]SKUEvidence{}
	}
	return a.request, requiredCopySKUs(a.values)
}
func requiredCopySKUs(in map[string]SKUEvidence) map[string]SKUEvidence {
	out := make(map[string]SKUEvidence, len(in))
	for typ, e := range in {
		values := make(map[string]SKUAvailability, len(e.Values))
		for k, v := range e.Values {
			v.BlockedZones = append([]string{}, v.BlockedZones...)
			values[k] = v
		}
		out[typ] = SKUEvidence{Status: e.Status, Values: values}
	}
	return out
}

type requiredBudget struct{ entries, text int }

func (b *requiredBudget) add(n int, fields ...string) bool {
	if n < 0 || n > MaxInventoryEntries-b.entries {
		return false
	}
	b.entries += n
	for _, s := range fields {
		if len(s) > auxTextBudget-b.text {
			return false
		}
		b.text += len(s)
	}
	return true
}
func requiredScalar(s string) bool { n := 0; return auxiliaryText(s, &n) }
func requiredKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func requiredResourceFields(r assessment.Resource) []string {
	return []string{r.ID, r.SubscriptionID, r.ResourceGroup, r.Name, r.Type, r.Location, r.Kind, r.SKUName, r.SKUTier, r.SKUFamily}
}
func requiredResourceIdentity(r assessment.Resource, scopes map[string]RequiredInventoryScope) bool {
	if r.ID == "" || r.Name == "" || !providerTypeIdentifier(r.Type) || !subscriptionID.MatchString(r.SubscriptionID) || len(r.SubscriptionID) != 36 || r.SKUCapacity < 0 || int64(r.SKUCapacity) > 2147483647 {
		return false
	}
	if _, ok := scopes[strings.ToLower(r.SubscriptionID)]; !ok {
		return false
	}
	if !inventoryIdentity(r.ID, r.SubscriptionID) {
		return false
	}
	for _, segment := range strings.Split(r.ID[1:], "/") {
		if segment == "" || segment != strings.TrimSpace(segment) {
			return false
		}
	}
	id, err := arm.ParseResourceID(r.ID)
	return err == nil && id.Name != "" && strings.EqualFold(id.SubscriptionID, r.SubscriptionID) && strings.EqualFold(id.ResourceGroupName, r.ResourceGroup) && strings.EqualFold(id.ResourceType.String(), r.Type)
}

// Errors contain fixed codes only. No input, SDK/provider diagnostic or URL is retained.
func requiredBuildError(ctx context.Context, observe func(string), code string) (*RequiredInventoryPlan, error) {
	if observe != nil {
		observe("before-return")
	}
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, errors.New(code)
}

func requiredTwice(b *requiredBudget, n int, fields ...string) bool {
	return b.add(n, fields...) && b.add(n, fields...)
}
func requiredPhase(ctx context.Context, observe func(string), phase string) error {
	if observe != nil {
		observe(phase)
	}
	return ctx.Err()
}
func requiredSupported(typ string) (catalog, unsupported bool) {
	switch typ {
	case "microsoft.compute/virtualmachines", "microsoft.compute/disks", "microsoft.storage/storageaccounts":
		return false, false
	case "microsoft.compute/virtualmachinescalesets", "microsoft.cognitiveservices/accounts", "microsoft.sql/servers/databases", "microsoft.sql/managedinstances":
		return true, false
	default:
		return false, true
	}
}

type requiredGroup struct {
	requirement RequiredInventoryRequirement
	skus        map[string]bool
	kinds       map[string]map[string]bool
}

func requiredTuple(r RequiredInventoryRequirement) [4]string {
	return [4]string{r.SubscriptionID, r.SourceRegion, r.ResourceType, r.Kind}
}
func requiredTupleLess(a, b [4]string) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// BuildRequiredInventory validates stable caller inputs and owns both diagnostic
// and private views. It performs no discovery, network request or scoring.
func BuildRequiredInventory(ctx context.Context, selected map[string]string, resources []assessment.Resource, declarations []InventoryScopeDeclaration) (*RequiredInventoryPlan, error) {
	return buildRequiredInventoryObserved(ctx, selected, resources, declarations, nil)
}
func buildRequiredInventoryObserved(ctx context.Context, selected map[string]string, resources []assessment.Resource, declarations []InventoryScopeDeclaration, observe func(string)) (*RequiredInventoryPlan, error) {
	fail := func(code string) (*RequiredInventoryPlan, error) { return requiredBuildError(ctx, observe, code) }
	invalid := func() (*RequiredInventoryPlan, error) { return fail("required_inventory_invalid") }
	limit := func() (*RequiredInventoryPlan, error) { return fail("required_inventory_limit") }
	if ctx == nil {
		return invalid()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(selected) > MaxSubscriptions || len(declarations) > MaxSubscriptions || len(resources) > MaxInventoryRows {
		return limit()
	}
	input := requiredBudget{}
	scopes := make(map[string]RequiredInventoryScope, len(selected))
	for _, id := range requiredKeys(selected) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		name := selected[id]
		if !input.add(1) {
			return limit()
		}
		for _, v := range []string{id, name} {
			if !requiredScalar(v) {
				return invalid()
			}
			if !input.add(0, v) {
				return limit()
			}
		}
		if len(id) != 36 || !subscriptionID.MatchString(id) || name == "" {
			return invalid()
		}
		key := strings.ToLower(id)
		if _, exists := scopes[key]; exists {
			return invalid()
		}
		scopes[key] = RequiredInventoryScope{Name: name, Status: "unknown"}
	}
	declared := make(map[string]bool, len(declarations))
	for _, d := range declarations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !input.add(1) {
			return limit()
		}
		for _, v := range []string{d.SubscriptionID, d.Status} {
			if !requiredScalar(v) {
				return invalid()
			}
			if !input.add(0, v) {
				return limit()
			}
		}
		if len(d.SubscriptionID) != 36 || !subscriptionID.MatchString(d.SubscriptionID) || (d.Status != "complete" && d.Status != "partial" && d.Status != "unknown") {
			return invalid()
		}
		id := strings.ToLower(d.SubscriptionID)
		s, ok := scopes[id]
		if !ok || declared[id] {
			return invalid()
		}
		declared[id] = true
		s.Status = d.Status
		scopes[id] = s
	}
	seen := make(map[string]bool, len(resources))
	for _, r := range resources {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !input.add(1) {
			return limit()
		}
		for _, v := range requiredResourceFields(r) {
			if !requiredScalar(v) {
				return invalid()
			}
			if !input.add(0, v) {
				return limit()
			}
		}
		if !requiredResourceIdentity(r, scopes) {
			return invalid()
		}
		id := strings.ToLower(r.ID)
		if seen[id] {
			return invalid()
		}
		seen[id] = true
	}
	if err := requiredPhase(ctx, observe, "after-input"); err != nil {
		return nil, err
	}
	status := "complete"
	allUnknown := len(scopes) > 0
	for _, s := range scopes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if s.Status != "complete" {
			status = "partial"
		}
		if s.Status != "unknown" {
			allUnknown = false
		}
	}
	if allUnknown && len(resources) == 0 {
		status = "unknown"
	}
	output := requiredBudget{}
	if !output.add(0, status) {
		return limit()
	}
	for _, id := range requiredKeys(scopes) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		s := scopes[id]
		if !requiredTwice(&output, 1, id, s.Name, s.Status) {
			return limit()
		}
	}
	derived := make([]RequiredInventoryResource, 0, len(resources))
	groups := map[[4]string]*requiredGroup{}
	for _, raw := range resources {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		record := raw
		record.Tags = nil
		record.SLA = ""
		sub, typ, kind := strings.ToLower(raw.SubscriptionID), strings.ToLower(raw.Type), strings.ToLower(raw.Kind)
		source, known := providerLocation(raw.Location)
		physical := known && costPhysical(source)
		if !physical {
			source = ""
		}
		fields := requiredResourceFields(record)
		fields = append(fields, sub, source, typ, kind)
		if !requiredTwice(&output, 1, fields...) {
			return limit()
		}
		derived = append(derived, RequiredInventoryResource{Record: record, SubscriptionID: sub, SourceRegion: source, ResourceType: typ, Kind: kind, Physical: physical})
		if !physical {
			continue
		}
		cognitive, storage := typ == "microsoft.cognitiveservices/accounts", typ == "microsoft.storage/storageaccounts"
		canonicalKind := ""
		missingKind, unrepresentableKind := false, false
		if cognitive || storage {
			missingKind = raw.Kind == ""
			unrepresentableKind = !missingKind && !storageIdentifier(raw.Kind)
			if !missingKind && !unrepresentableKind {
				canonicalKind = kind
			}
		}
		keyKind := ""
		if cognitive {
			keyKind = canonicalKind
		}
		tuple := [4]string{sub, source, typ, keyKind}
		g := groups[tuple]
		if g == nil {
			if !requiredTwice(&output, 1, tuple[:]...) {
				return limit()
			}
			catalog, unsupported := requiredSupported(typ)
			g = &requiredGroup{requirement: RequiredInventoryRequirement{SubscriptionID: sub, SourceRegion: source, ResourceType: typ, Kind: keyKind, SKUs: []string{}, Kinds: map[string][]string{}, ResourceIDs: []string{}, CatalogOnly: catalog, Unsupported: unsupported}, skus: map[string]bool{}, kinds: map[string]map[string]bool{}}
			groups[tuple] = g
		}
		r := &g.requirement
		r.MissingKind = r.MissingKind || missingKind
		r.UnrepresentableKind = r.UnrepresentableKind || unrepresentableKind
		if !requiredTwice(&output, 1, raw.ID) {
			return limit()
		}
		r.ResourceIDs = append(r.ResourceIDs, raw.ID)
		missing := raw.SKUName == ""
		sku := strings.ToLower(strings.TrimSpace(raw.SKUName))
		representable := storageIdentifier(sku)
		// SQL generic SKU aliases remain diagnostics, including slash-bearing labels.
		// This does not establish version/edition/policy or a collected coordinate.
		if typ == "microsoft.sql/servers/databases" || typ == "microsoft.sql/managedinstances" {
			representable = providerTypeIdentifier(sku)
		}
		unrepresentable := !missing && !representable
		r.MissingSKU = r.MissingSKU || missing
		r.UnrepresentableSKU = r.UnrepresentableSKU || unrepresentable
		if missing || unrepresentable {
			continue
		}
		if !g.skus[sku] {
			if !requiredTwice(&output, 1, sku) {
				return limit()
			}
			g.skus[sku] = true
		}
		if storage {
			kinds := g.kinds[sku]
			if kinds == nil {
				if !requiredTwice(&output, 1, sku) {
					return limit()
				}
				kinds = map[string]bool{}
				g.kinds[sku] = kinds
			}
			if !kinds[canonicalKind] {
				if len(kinds) >= 16 {
					return limit()
				}
				if !requiredTwice(&output, 1, canonicalKind) {
					return limit()
				}
				kinds[canonicalKind] = true
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	sort.Slice(derived, func(i, j int) bool {
		return strings.ToLower(derived[i].Record.ID) < strings.ToLower(derived[j].Record.ID)
	})
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	tuples := make([][4]string, 0, len(groups))
	for tuple := range groups {
		tuples = append(tuples, tuple)
	}
	sort.Slice(tuples, func(i, j int) bool { return requiredTupleLess(tuples[i], tuples[j]) })
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	requirements := make([]RequiredInventoryRequirement, 0, len(tuples))
	for _, tuple := range tuples {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		g := groups[tuple]
		r := g.requirement
		r.SKUs = requiredKeys(g.skus)
		for _, sku := range requiredKeys(g.kinds) {
			r.Kinds[sku] = requiredKeys(g.kinds[sku])
		}
		sort.Slice(r.ResourceIDs, func(i, j int) bool { return strings.ToLower(r.ResourceIDs[i]) < strings.ToLower(r.ResourceIDs[j]) })
		requirements = append(requirements, r)
	}
	if err := requiredPhase(ctx, observe, "after-derived-work"); err != nil {
		return nil, err
	}
	privateResources := append([]RequiredInventoryResource{}, derived...)
	privateRequirements := make([]RequiredInventoryRequirement, 0, len(requirements))
	publicRequirements := make([]RequiredInventoryRequirement, 0, len(requirements))
	for _, r := range requirements {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		privateRequirements = append(privateRequirements, requiredCopyRequirement(r))
		publicRequirements = append(publicRequirements, requiredCopyRequirement(r))
	}
	privateScopes := make(map[string]RequiredInventoryScope, len(scopes))
	publicScopes := make(map[string]RequiredInventoryScope, len(scopes))
	for _, id := range requiredKeys(scopes) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		privateScopes[id] = scopes[id]
		publicScopes[id] = scopes[id]
	}
	result := &RequiredInventoryPlan{Resources: append([]RequiredInventoryResource{}, derived...), Requirements: publicRequirements, Scopes: publicScopes, Status: status, state: &requiredInventoryState{resources: privateResources, requirements: privateRequirements, scopes: privateScopes}}
	if err := requiredPhase(ctx, observe, "before-return"); err != nil {
		return nil, err
	}
	return result, nil
}
func requiredCopyRequirement(r RequiredInventoryRequirement) RequiredInventoryRequirement {
	r.SKUs = append([]string{}, r.SKUs...)
	r.ResourceIDs = append([]string{}, r.ResourceIDs...)
	kinds := make(map[string][]string, len(r.Kinds))
	for k, v := range r.Kinds {
		kinds[k] = append([]string{}, v...)
	}
	r.Kinds = kinds
	return r
}
