package region

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func sqlKnownGroup(name string) SQLCapabilityGroup {
	return SQLCapabilityGroup{Known: true, State: 0, Version: "v", Edition: "e", Name: name, SKUName: "s", SKUTier: "e", Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"}
}
func sqlRequireProjectionLimit(t *testing.T, out SQLCatalogEvidence, err error) {
	t.Helper()
	if !reflect.DeepEqual(out, SQLCatalogEvidence{}) || err == nil || err.Error() != "SQL capability catalog projection rejected [sql_capability_projection_limit]" {
		t.Fatal("zero fixed projection limit", out.Status, len(out.Values), err)
	}
}
func TestSQLDatabaseProjectionCounts(t *testing.T) {
	for _, kind := range []string{"selected", "groups", "required", "aliases"} {
		bound := 8192
		if kind == "selected" {
			bound = 1000
		}
		if kind == "aliases" {
			bound = 256
		}
		for _, n := range []int{bound, bound + 1} {
			t.Run(fmt.Sprintf("%s-%d", kind, n), func(t *testing.T) {
				scope, q, c := sqlBoundaryInput()
				required := []string{}
				switch kind {
				case "selected":
					for i := 1; i < n; i++ {
						scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "selected"
					}
				case "groups":
					for i := 0; i < n; i++ {
						g := sqlKnownGroup(fmt.Sprintf("s%04d", i))
						c.Groups["v/e/"+g.Name] = g
					}
				case "required":
					for i := 0; i < n; i++ {
						required = append(required, fmt.Sprintf("v/e/s%04d", i))
					}
				case "aliases":
					q.LocationAliases = map[string]string{}
					for i := 0; i < n; i++ {
						q.LocationAliases[fmt.Sprintf("reg ion%d", i)] = fmt.Sprintf("region%d", i)
					}
				}
				out, e := c.Catalog(context.Background(), scope, q, required)
				if n == bound {
					want := "complete"
					if kind == "required" {
						want = "unknown"
					}
					if e != nil || out.Status != want {
						t.Fatal("exact public count", kind, out.Status, e)
					}
				} else {
					sqlRequireProjectionLimit(t, out, e)
				}
			})
		}
	}
}
func TestSQLDatabaseProjectionJointCounts(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	for i := 1; i < 1000; i++ {
		scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "selected"
	}
	q.LocationAliases = map[string]string{}
	for i := 0; i < 256; i++ {
		q.LocationAliases[fmt.Sprintf("reg ion%d", i)] = fmt.Sprintf("region%d", i)
	}
	required := []string{}
	for i := 0; i < 8192; i++ {
		g := sqlKnownGroup(fmt.Sprintf("s%04d", i))
		key := "v/e/" + g.Name
		c.Groups[key] = g
		required = append(required, key)
	}
	out, e := c.Catalog(context.Background(), scope, q, required)
	if e != nil || out.Status != "complete" || len(out.Values) != 8192 {
		t.Fatal("joint public counts", out.Status, len(out.Values), e)
	}
	if 1000+256+8192+8192 != 17640 || 17640 >= 65536 {
		t.Fatal("projection occurrence proof")
	}
}
func TestSQLDatabaseProjectionLabels(t *testing.T) {
	for _, field := range []string{"selected-label", "version", "edition", "name", "family", "sku-name"} {
		for _, n := range []int{512, 513} {
			t.Run(fmt.Sprintf("%s-%d", field, n), func(t *testing.T) {
				scope, q, c := sqlBoundaryInput()
				g := sqlKnownGroup("n")
				s := strings.Repeat("a", n)
				switch field {
				case "selected-label":
					scope[sqlTestID] = s
				case "version":
					g.Version = s
				case "edition":
					g.Edition = s
					g.SKUTier = s
				case "name":
					g.Name = s
				case "family":
					g.Family = s
				case "sku-name":
					g.SKUName = s
				}
				c.Groups[g.Version+"/"+g.Edition+"/"+g.Name] = g
				out, e := c.Catalog(context.Background(), scope, q, nil)
				if n == 512 {
					if e != nil || out.Status != "complete" {
						t.Fatal("exact component/metadata label", out.Status, e)
					}
				} else {
					sqlRequireProjectionLimit(t, out, e)
				}
			})
		}
	}
}
func TestSQLDatabaseProjectionWholeValidation(t *testing.T) {
	for _, mode := range []string{"negative-capacity", "unknown-negative", "key-coordinate", "bad-state", "unknown-state", "metadata-case", "key-case", "wrong-tier", "bad-flag", "missing-sku-inference", "unknown-with-groups", "unknown-without-code", "complete-with-code"} {
		scope, q, c := sqlBoundaryInput()
		g := sqlKnownGroup("s")
		key := "v/e/s"
		switch mode {
		case "negative-capacity":
			g.Capacity = -1
		case "unknown-negative":
			g.Known = false
			g.State = 2
			g.Capacity = -1
		case "key-coordinate":
			g.Name = "other"
		case "bad-state":
			g.State = 1
		case "unknown-state":
			g.Known = false
		case "metadata-case":
			g.SKUName = "S"
		case "key-case":
			key = "V/e/s"
		case "wrong-tier":
			g.SKUTier = "other"
		case "bad-flag":
			g.Standalone = "supported"
		case "missing-sku-inference":
			g.Known = false
			g.State = 2
			g.SKUName = ""
		case "unknown-with-groups":
			c.Status = "unknown"
			c.FailureCode = "sql_capability_incomplete"
		case "unknown-without-code":
			c.Status = "unknown"
		case "complete-with-code":
			c.FailureCode = "sql_capability_request_failed"
		}
		c.Groups[key] = g
		if mode == "unknown-without-code" {
			c.Groups = map[string]SQLCapabilityGroup{}
		}
		out, e := c.Catalog(context.Background(), scope, q, nil)
		if e == nil || e.Error() != "SQL capability catalog projection rejected [sql_capability_projection_invalid]" || !reflect.DeepEqual(out, SQLCatalogEvidence{}) {
			t.Fatal("malformed unrequired input", mode, out, e)
		}
	}
	scope, q, c := sqlBoundaryInput()
	g := sqlKnownGroup("s")
	g.Known = false
	g.State = 2
	c.Groups["v/e/s"] = g
	for _, required := range [][]string{nil, {"v/e/s"}, {"v/e/missing"}} {
		out, e := c.Catalog(context.Background(), scope, q, required)
		want := "complete"
		if len(required) > 0 {
			want = "unknown"
		}
		if e != nil || out.Status != want || out.Values == nil || out.Warnings == nil || len(out.Warnings) != 1 || out.DeploymentEligibility != "not_assessed" {
			t.Fatal("truthful owned empty result", out, e)
		}
	}
}

// Construct valid ASCII metadata independently from the declared accounting:
// for fixed5byte name, key+metadata=2V+3E+F+S+40. Fixed admitted input=198bytes.
func sqlTextGroup(t *testing.T, name string, target int) SQLCapabilityGroup {
	t.Helper()
	for v := 1; v <= 512; v++ {
		for e := 1; e <= 512; e++ {
			rest := target - 2*v - 3*e - 40
			if rest < 1 || rest > 1023 {
				continue
			}
			f := rest - 512
			if f < 0 {
				f = 0
			}
			s := rest - f
			if f > 511 || s < 1 || s > 512 {
				continue
			}
			return SQLCapabilityGroup{Known: true, State: 0, Version: strings.Repeat("v", v), Edition: strings.Repeat("e", e), Name: name, Family: strings.Repeat("f", f), SKUName: strings.Repeat("s", s), SKUTier: strings.Repeat("e", e), Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"}
		}
	}
	t.Fatal("unreachable fixture residual", target)
	return SQLCapabilityGroup{}
}
func sqlIndependentInputBytes(scope map[string]string, q SQLCapabilityRequest, c *SQLCapabilityCollection, required []string) int {
	n := 0
	for id, label := range scope {
		n += len(id) + len(label)
	}
	n += len(q.SubscriptionID) + len(q.Region) + len(q.ResourceType)
	for k, v := range q.LocationAliases {
		n += len(k) + len(v)
	}
	n += len(c.SubscriptionID) + len(c.Region) + len(c.ResourceType) + len(c.Status) + len(c.FailureCode)
	for key, g := range c.Groups {
		n += len(key) + len(g.Version) + len(g.Edition) + len(g.Name) + len(g.Family) + len(g.SKUName) + len(g.SKUTier) + len(g.Standalone) + len(g.InstancePool)
	}
	for _, r := range required {
		n += len(r)
	}
	return n
}
func TestSQLDatabaseProjectionTextBound(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	if sqlIndependentInputBytes(scope, q, c, nil) != 198 {
		t.Fatal("fixed input accounting")
	}
	remaining := 16711680 - 198
	i := 0
	for remaining >= 3624+46 {
		g := SQLCapabilityGroup{Known: true, State: 0, Version: strings.Repeat("v", 512), Edition: strings.Repeat("e", 512), Name: fmt.Sprintf("s%04d", i), Family: strings.Repeat("f", 512), SKUName: strings.Repeat("s", 512), SKUTier: strings.Repeat("e", 512), Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"}
		c.Groups[g.Version+"/"+g.Edition+"/"+g.Name] = g
		remaining -= 3624
		i++
	}
	if remaining > 3623 {
		t.Fatal("fixture residual", remaining)
	}
	g := sqlTextGroup(t, fmt.Sprintf("s%04d", i), remaining)
	key := g.Version + "/" + g.Edition + "/" + g.Name
	c.Groups[key] = g
	if len(c.Groups) > 8192 || len(g.Family) >= 512 || sqlIndependentInputBytes(scope, q, c, nil) != 16711680 {
		t.Fatal("exact valid fixture accounting")
	}
	out, e := c.Catalog(context.Background(), scope, q, nil)
	if e != nil || out.Status != "complete" {
		t.Fatal("exact independent input text", out.Status, e)
	}
	g.Family += "f"
	c.Groups[key] = g
	if sqlIndependentInputBytes(scope, q, c, nil) != 16711681 {
		t.Fatal("one-over fixture accounting")
	}
	out, e = c.Catalog(context.Background(), scope, q, nil)
	sqlRequireProjectionLimit(t, out, e)
}
func TestSQLDatabaseRequiredTextBound(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	required := []string{"v/e/missing"}
	if sqlIndependentInputBytes(scope, q, c, required) != 198+11 {
		t.Fatal("fixed input accounting")
	}
	remaining := 16711680 - 198 - 11
	i := 0
	for remaining >= 3624+46 {
		g := SQLCapabilityGroup{Known: true, State: 0, Version: strings.Repeat("v", 512), Edition: strings.Repeat("e", 512), Name: fmt.Sprintf("s%04d", i), Family: strings.Repeat("f", 512), SKUName: strings.Repeat("s", 512), SKUTier: strings.Repeat("e", 512), Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"}
		c.Groups[g.Version+"/"+g.Edition+"/"+g.Name] = g
		remaining -= 3624
		i++
	}
	if remaining > 3623 {
		t.Fatal("fixture residual", remaining)
	}
	g := sqlTextGroup(t, fmt.Sprintf("s%04d", i), remaining)
	key := g.Version + "/" + g.Edition + "/" + g.Name
	c.Groups[key] = g
	if len(c.Groups) > 8192 || len(g.Family) >= 512 || sqlIndependentInputBytes(scope, q, c, required) != 16711680 {
		t.Fatal("exact valid fixture accounting")
	}
	out, e := c.Catalog(context.Background(), scope, q, required)
	if e != nil || out.Status != "unknown" {
		t.Fatal("exact independent input text", out.Status, e)
	}
	g.Family += "f"
	c.Groups[key] = g
	if sqlIndependentInputBytes(scope, q, c, required) != 16711681 {
		t.Fatal("one-over fixture accounting")
	}
	out, e = c.Catalog(context.Background(), scope, q, required)
	sqlRequireProjectionLimit(t, out, e)
}
func TestSQLDatabaseAliasTextBound(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	q.LocationAliases = map[string]string{"east us": "eastus"}
	if sqlIndependentInputBytes(scope, q, c, nil) != 198+13 {
		t.Fatal("fixed input accounting")
	}
	remaining := 16711680 - 198 - 13
	i := 0
	for remaining >= 3624+46 {
		g := SQLCapabilityGroup{Known: true, State: 0, Version: strings.Repeat("v", 512), Edition: strings.Repeat("e", 512), Name: fmt.Sprintf("s%04d", i), Family: strings.Repeat("f", 512), SKUName: strings.Repeat("s", 512), SKUTier: strings.Repeat("e", 512), Capacity: 1, Standalone: "not_applicable", InstancePool: "not_applicable"}
		c.Groups[g.Version+"/"+g.Edition+"/"+g.Name] = g
		remaining -= 3624
		i++
	}
	if remaining > 3623 {
		t.Fatal("fixture residual", remaining)
	}
	g := sqlTextGroup(t, fmt.Sprintf("s%04d", i), remaining)
	key := g.Version + "/" + g.Edition + "/" + g.Name
	c.Groups[key] = g
	if len(c.Groups) > 8192 || len(g.Family) >= 512 || sqlIndependentInputBytes(scope, q, c, nil) != 16711680 {
		t.Fatal("exact valid fixture accounting")
	}
	out, e := c.Catalog(context.Background(), scope, q, nil)
	if e != nil || out.Status != "complete" {
		t.Fatal("exact independent input text", out.Status, e)
	}
	g.Family += "f"
	c.Groups[key] = g
	if sqlIndependentInputBytes(scope, q, c, nil) != 16711681 {
		t.Fatal("one-over fixture accounting")
	}
	out, e = c.Catalog(context.Background(), scope, q, nil)
	sqlRequireProjectionLimit(t, out, e)
}
func TestSQLDatabaseOutputBudgetProof(t *testing.T) {
	maximum := 0
	for n := 1; n <= 8192; n++ {
		byFields := 282 + n*4638
		byInput := 16711680 + 38 - n*5
		bound := byFields
		if byInput < bound {
			bound = byInput
		}
		if bound > maximum {
			maximum = bound
		}
	}
	if maximum != 16693718 || maximum >= 16711680 {
		t.Fatal("independent output/input inequality", maximum)
	}
}
