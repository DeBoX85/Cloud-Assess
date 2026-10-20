// SQL capability catalog behavior uses MIT-licensed Azure Quick Review and Microsoft
// REST specifications. See NOTICE.md and SQL_CAPABILITY_CONTRACT_20261017.md.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const sqlCapabilityVersion = "2021-11-01"
const sqlDatabaseType = "microsoft.sql/servers/databases"
const MaxSQLCapabilityLeaves = 8192
const MaxSQLCapabilityAliases = 256
const MaxSQLCoordinateBytes = 2051

type SQLCapabilityRequest struct {
	SubscriptionID, Region, ResourceType string
	LocationAliases                      map[string]string
}
type SQLCapabilityGroup struct {
	Known                                            bool
	State                                            int
	Version, Edition, Name, Family, SKUName, SKUTier string
	Capacity                                         int64
	Standalone, InstancePool                         string
}
type SQLCapabilityCollection struct {
	SubscriptionID, Region, ResourceType, Status string
	Groups                                       map[string]SQLCapabilityGroup
	FailureCode                                  string
}
type SQLCapabilityCollector struct {
	transport *RESTQuotaCollector
	family    string
}

func NewSQLDatabaseSKUCollector(endpoint string, getter RESTQuotaGetter) (*SQLCapabilityCollector, error) {
	transport, err := NewRESTQuotaCollector(endpoint, getter)
	if err != nil {
		return nil, fmt.Errorf("invalid SQL capability origin or getter")
	}
	return &SQLCapabilityCollector{transport: transport, family: sqlDatabaseType}, nil
}
func sqlCapabilityBinding(scope map[string]string, q SQLCapabilityRequest) (string, string, bool) {
	if len(q.SubscriptionID) != 36 || len(q.ResourceType) > MaxLabelBytes {
		return "", "", false
	}
	id := strings.ToLower(q.SubscriptionID)
	kind := strings.ToLower(q.ResourceType)
	region, ok := diskPhysical(q.Region)
	return id, kind, subscriptionID.MatchString(id) && scope[id] != "" && ok && region == q.Region && (kind == sqlDatabaseType || kind == sqlManagedInstanceType)
}
func sqlCapabilityFailure(code string) bool {
	switch code {
	case "sql_capability_request_failed", "sql_capability_status_invalid", "sql_capability_byte_limit", "sql_capability_invalid_response", "sql_capability_row_limit", "sql_capability_work_limit", "sql_capability_text_limit", "sql_capability_continuation_invalid", "sql_capability_collision", "sql_capability_incomplete":
		return true
	}
	return false
}
func sqlCapabilityState(nodes ...map[string]json.RawMessage) (bool, int) {
	known, disabled := true, false
	for _, node := range nodes {
		var status string
		if raw, has := node["status"]; !has || !diskScalar(raw, &status) {
			known = false
		} else {
			switch status {
			case "Available", "Default":
			case "Disabled":
				disabled = true
			default:
				known = false
			}
		}
	}
	if !known {
		return false, SKUUnavailable
	}
	if disabled {
		return true, SKUUnavailable
	}
	return true, SKUAvailable
}
func sqlCapabilityName(node map[string]json.RawMessage) (string, bool) {
	var s string
	if !diskScalar(node["name"], &s) || !storageIdentifier(s) {
		return "", false
	}
	return strings.ToLower(s), true
}

// Selected sibling identities are unique even when duplicate parents have disjoint descendants.
func sqlCapabilityChildren(ctx context.Context, node map[string]json.RawMessage, key string) ([]map[string]json.RawMessage, bool, string) {
	raw, has := node[key]
	if !has {
		return nil, false, ""
	}
	rows, ok := providerArray(raw)
	if !ok {
		return nil, false, "invalid_response"
	}
	out := make([]map[string]json.RawMessage, 0, len(rows))
	seen := map[string]bool{}
	for _, raw := range rows {
		if ctx.Err() != nil {
			return nil, false, "invalid_response"
		}
		var row map[string]json.RawMessage
		if json.Unmarshal(raw, &row) != nil || row == nil {
			return nil, false, "invalid_response"
		}
		name, valid := sqlCapabilityName(row)
		if !valid {
			return nil, false, "invalid_response"
		}
		if seen[name] {
			return nil, false, "collision"
		}
		seen[name] = true
		out = append(out, row)
	}
	return out, true, ""
}
func sqlDatabaseGroup(root, version, edition, leaf map[string]json.RawMessage) (SQLCapabilityGroup, bool) {
	v, _ := sqlCapabilityName(version)
	e, _ := sqlCapabilityName(edition)
	n, _ := sqlCapabilityName(leaf)
	known, state := sqlCapabilityState(root, version, edition, leaf)
	g := SQLCapabilityGroup{Known: known, State: state, Version: v, Edition: e, Name: n, Standalone: "not_applicable", InstancePool: "not_applicable"}
	raw, has := leaf["sku"]
	if !has {
		g.Known = false
		g.State = SKUUnavailable
		return g, true
	}
	var sku map[string]json.RawMessage
	if json.Unmarshal(raw, &sku) != nil || sku == nil {
		return g, false
	}
	var name string
	if !diskScalar(sku["name"], &name) || !storageIdentifier(name) {
		return g, false
	}
	g.SKUName = strings.ToLower(name)
	if raw, has := sku["tier"]; has {
		var tier string
		if !diskScalar(raw, &tier) || !storageIdentifier(tier) || !strings.EqualFold(tier, e) {
			return g, false
		}
		g.SKUTier = strings.ToLower(tier)
	} else {
		g.Known = false
	}
	if raw, has := sku["capacity"]; has {
		if !diskScalar(raw, &g.Capacity) || g.Capacity <= 0 || g.Capacity > 2147483647 {
			return g, false
		}
	} else {
		g.Known = false
	}
	if raw, has := sku["family"]; has {
		var family string
		if !diskScalar(raw, &family) || family != "" && !storageIdentifier(family) {
			return g, false
		}
		g.Family = strings.ToLower(family)
	}
	if !g.Known {
		g.State = SKUUnavailable
	}
	return g, true
}
func decodeSQLDatabase(ctx context.Context, body []byte, q SQLCapabilityRequest) (map[string]SQLCapabilityGroup, string) {
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
	versions, complete, code := sqlCapabilityChildren(ctx, root, "supportedServerVersions")
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
			objectives, complete, code := sqlCapabilityChildren(ctx, edition, "supportedServiceLevelObjectives")
			if code != "" {
				return nil, code
			}
			incomplete = incomplete || !complete
			for _, leaf := range objectives {
				if ctx.Err() != nil {
					return nil, "invalid_response"
				}
				leaves++
				if leaves > MaxSQLCapabilityLeaves {
					return nil, "row_limit"
				}
				g, valid := sqlDatabaseGroup(root, version, edition, leaf)
				if !valid {
					return nil, "invalid_response"
				}
				key := g.Version + "/" + g.Edition + "/" + g.Name
				groups[key] = g
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
func (c *SQLCapabilityCollector) Collect(ctx context.Context, selected map[string]string, q SQLCapabilityRequest) (*SQLCapabilityCollection, error) {
	reject := func() (*SQLCapabilityCollection, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		return nil, fmt.Errorf("SQL capability request rejected [sql_capability_request_invalid]")
	}
	if ctx.Err() != nil || c == nil || c.transport == nil {
		return reject()
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, selected, &text)
	if !ok || len(scope) == 0 {
		return reject()
	}
	for _, label := range scope {
		if label == "" {
			return reject()
		}
	}
	id, kind, ok := sqlCapabilityBinding(scope, q)
	if !ok || kind != c.family {
		return reject()
	}
	aliases, ok := cognitiveAliases(ctx, q.LocationAliases)
	if !ok {
		return reject()
	}
	q.LocationAliases = aliases
	target := *c.transport.origin
	target.Path = "/subscriptions/" + id + "/providers/Microsoft.Sql/locations/" + q.Region + "/capabilities"
	target.RawQuery = url.Values{"api-version": {sqlCapabilityVersion}}.Encode()
	result := &SQLCapabilityCollection{SubscriptionID: id, Region: q.Region, ResourceType: kind, Status: "complete", Groups: map[string]SQLCapabilityGroup{}}
	fail := func(code string) (*SQLCapabilityCollection, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		result.Status = "unknown"
		result.FailureCode = "sql_capability_" + code
		result.Groups = map[string]SQLCapabilityGroup{}
		return result, nil
	}
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	body, response, err := c.transport.getter.GetBoundedWithResponse(ctx, target.String(), MaxQuotaPageBytes)
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	if errors.Is(err, context.Canceled) {
		return nil, context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, context.DeadlineExceeded
	}
	if response != nil && response.StatusCode != http.StatusOK {
		return fail("status_invalid")
	}
	if err != nil || response == nil {
		return fail("request_failed")
	}
	if len(body) > MaxQuotaPageBytes {
		return fail("byte_limit")
	}
	var groups map[string]SQLCapabilityGroup
	var code string
	if kind == sqlManagedInstanceType {
		groups, code = decodeSQLManagedInstance(ctx, body, q)
	} else {
		groups, code = decodeSQLDatabase(ctx, body, q)
	}
	if code != "" {
		return fail(code)
	}
	result.Groups = groups
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return result, nil
}
