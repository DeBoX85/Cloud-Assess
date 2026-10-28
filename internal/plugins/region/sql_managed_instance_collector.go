// Managed Instance capability catalog uses explicit pinned Microsoft SQL wire metadata.
// See the frozen SQL contract and complete MIT notices in NOTICE.md.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const sqlManagedInstanceType = "microsoft.sql/managedinstances"
const MaxSQLManagedInstanceLeaves = 8192

func NewSQLManagedInstanceSKUCollector(endpoint string, getter RESTQuotaGetter) (*SQLCapabilityCollector, error) {
	transport, e := NewRESTQuotaCollector(endpoint, getter)
	if e != nil {
		return nil, fmt.Errorf("invalid SQL capability origin or getter")
	}
	return &SQLCapabilityCollector{transport: transport, family: sqlManagedInstanceType}, nil
}
func sqlManagedInstanceMode(node map[string]json.RawMessage, key string) (string, bool) {
	raw, has := node[key]
	if !has {
		return "unknown", true
	}
	var value bool
	if !diskScalar(raw, &value) {
		return "", false
	}
	if value {
		return "supported", true
	}
	return "unsupported", true
}
func sqlManagedInstanceGroup(root, version, edition, family, leaf map[string]json.RawMessage) (SQLCapabilityGroup, bool) {
	v, _ := sqlCapabilityName(version)
	e, _ := sqlCapabilityName(edition)
	f, _ := sqlCapabilityName(family)
	n, _ := sqlCapabilityName(leaf)
	known, state := sqlCapabilityState(root, version, edition, family, leaf)
	g := SQLCapabilityGroup{Known: known, State: state, Version: v, Edition: e, Family: f, Name: n, SKUTier: e}
	if raw, has := family["sku"]; has {
		var sku string
		if !diskScalar(raw, &sku) || !storageIdentifier(sku) {
			return g, false
		}
		g.SKUName = strings.ToLower(sku)
	} else {
		g.Known = false
	}
	if raw, has := leaf["value"]; has {
		if !diskScalar(raw, &g.Capacity) || g.Capacity <= 0 || g.Capacity > 2147483647 {
			return g, false
		}
	} else {
		g.Known = false
	}
	var ok bool
	g.Standalone, ok = sqlManagedInstanceMode(leaf, "standaloneSupported")
	if !ok {
		return g, false
	}
	g.InstancePool, ok = sqlManagedInstanceMode(leaf, "instancePoolSupported")
	if !ok {
		return g, false
	}
	if !g.Known {
		g.State = SKUUnavailable
	}
	return g, true
}
func decodeSQLManagedInstance(ctx context.Context, body []byte, q SQLCapabilityRequest) (map[string]SQLCapabilityGroup, string) {
	if !quotaJSON(ctx, body) {
		return nil, "invalid_response"
	}
	root, ok := providerObject(body, "name", "status", "reason", "supportedServerVersions", "supportedManagedInstanceVersions", "value", "nextLink")
	if !ok {
		return nil, "invalid_response"
	}
	if _, present := root["nextLink"]; present {
		return nil, "continuation_invalid"
	}
	if _, present := root["value"]; present {
		return nil, "invalid_response"
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var doc any
	if dec.Decode(&doc) != nil {
		return nil, "invalid_response"
	}
	budget := vmSKUBudget{ctx: ctx}
	if code := budget.walk(doc, 0, false); code != "" {
		switch code {
		case "vm_sku_work_limit":
			return nil, "work_limit"
		case "vm_sku_text_limit":
			return nil, "text_limit"
		default:
			return nil, "invalid_response"
		}
	}
	if !sqlCapabilityShape(ctx, body, "LocationCapabilities") {
		return nil, "invalid_response"
	}
	incomplete := false
	var rootName string
	if raw, has := root["name"]; has {
		if !diskScalar(raw, &rootName) {
			return nil, "invalid_response"
		}
		physical, _ := storageLocation(rootName, q.LocationAliases)
		if physical != "" && physical != q.Region {
			return nil, "invalid_response"
		}
		incomplete = physical == ""
	} else {
		incomplete = true
	}

	versions, complete, code := sqlCapabilityChildren(ctx, root, "supportedManagedInstanceVersions")
	if code != "" {
		return nil, code
	}
	incomplete = incomplete || !complete
	groups := map[string]SQLCapabilityGroup{}
	leaves := 0
	for _, version := range versions {
		editions, complete, code := sqlCapabilityChildren(ctx, version, "supportedEditions")
		if code != "" {
			return nil, code
		}
		incomplete = incomplete || !complete
		for _, edition := range editions {
			families, complete, code := sqlCapabilityChildren(ctx, edition, "supportedFamilies")
			if code != "" {
				return nil, code
			}
			incomplete = incomplete || !complete
			for _, family := range families {
				vcores, complete, code := sqlCapabilityChildren(ctx, family, "supportedVcoresValues")
				if code != "" {
					return nil, code
				}
				incomplete = incomplete || !complete
				for _, leaf := range vcores {
					if ctx.Err() != nil {
						return nil, "invalid_response"
					}
					leaves++
					if leaves > MaxSQLManagedInstanceLeaves {
						return nil, "row_limit"
					}
					g, valid := sqlManagedInstanceGroup(root, version, edition, family, leaf)
					if !valid {
						return nil, "invalid_response"
					}
					key := g.Version + "/" + g.Edition + "/" + g.Family + "/" + g.Name
					groups[key] = g
				}
			}
		}
	}
	if incomplete {
		return nil, "incomplete"
	}
	if ctx.Err() != nil {
		return nil, "invalid_response"
	}
	return groups, ""
}
func sqlManagedInstanceCoordinate(value string) (string, bool) {
	if len(value) > 2051 {
		return "", false
	}
	parts := strings.Split(value, "/")
	if len(parts) != 4 {
		return "", false
	}
	for _, p := range parts {
		if !storageIdentifier(p) {
			return "", false
		}
	}
	return strings.ToLower(value), true
}
func sqlCapabilityCoordinate(value, kind string) (string, bool) {
	if kind == sqlDatabaseType {
		return sqlCoordinate(value)
	}
	if kind == sqlManagedInstanceType {
		return sqlManagedInstanceCoordinate(value)
	}
	return "", false
}
func sqlValidManagedInstanceGroup(key string, g SQLCapabilityGroup) bool {
	canonical, ok := sqlManagedInstanceCoordinate(key)
	if !ok || canonical != key || key != g.Version+"/"+g.Edition+"/"+g.Family+"/"+g.Name {
		return false
	}
	for _, s := range []string{g.Version, g.Edition, g.Family, g.Name} {
		if !storageIdentifier(s) || strings.ToLower(s) != s {
			return false
		}
	}
	if g.SKUName != "" && (!storageIdentifier(g.SKUName) || strings.ToLower(g.SKUName) != g.SKUName) {
		return false
	}
	if g.SKUTier != g.Edition || g.State != SKUAvailable && g.State != SKUUnavailable || !g.Known && g.State != SKUUnavailable || g.Capacity < 0 || g.Capacity > 2147483647 {
		return false
	}
	if g.Known && (g.SKUName == "" || g.Capacity <= 0) {
		return false
	}
	for _, mode := range []string{g.Standalone, g.InstancePool} {
		if mode != "supported" && mode != "unsupported" && mode != "unknown" {
			return false
		}
	}
	return true
}
