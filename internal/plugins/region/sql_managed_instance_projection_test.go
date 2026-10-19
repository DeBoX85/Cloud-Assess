package region

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func sqlMIKnownGroup(name string) SQLCapabilityGroup {
	return SQLCapabilityGroup{Known: true, State: 0, Version: "v", Edition: "e", Name: name, Family: "f", SKUName: "s", SKUTier: "e", Capacity: 1, Standalone: "unknown", InstancePool: "unknown"}
}
func sqlMIRequireProjectionLimit(t *testing.T, out SQLCatalogEvidence, err error) {
	t.Helper()
	if !reflect.DeepEqual(out, SQLCatalogEvidence{}) || err == nil || err.Error() != "SQL capability catalog projection rejected [sql_capability_projection_limit]" {
		t.Fatal("zero fixed projection limit", out.Status, len(out.Values), err)
	}
}
func TestSQLMIProjectionCounts(t *testing.T) {
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
				scope, q, c := sqlMIInput()
				required := []string{}
				switch kind {
				case "selected":
					for i := 1; i < n; i++ {
						scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "selected"
					}
				case "groups":
					for i := 0; i < n; i++ {
						g := sqlMIKnownGroup(fmt.Sprintf("s%04d", i))
						c.Groups["v/e/f/"+g.Name] = g
					}
				case "required":
					for i := 0; i < n; i++ {
						required = append(required, fmt.Sprintf("v/e/f/s%04d", i))
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
					sqlMIRequireProjectionLimit(t, out, e)
				}
			})
		}
	}
}
func TestSQLMIProjectionJointCounts(t *testing.T) {
	scope, q, c := sqlMIInput()
	for i := 1; i < 1000; i++ {
		scope[fmt.Sprintf("%08x-1111-4111-8111-111111111111", i)] = "selected"
	}
	q.LocationAliases = map[string]string{}
	for i := 0; i < 256; i++ {
		q.LocationAliases[fmt.Sprintf("reg ion%d", i)] = fmt.Sprintf("region%d", i)
	}
	required := []string{}
	for i := 0; i < 8192; i++ {
		g := sqlMIKnownGroup(fmt.Sprintf("s%04d", i))
		key := "v/e/f/" + g.Name
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
func TestSQLMIProjectionLabels(t *testing.T) {
	for _, field := range []string{"selected-label", "version", "edition", "name", "family", "sku-name"} {
		for _, n := range []int{512, 513} {
			t.Run(fmt.Sprintf("%s-%d", field, n), func(t *testing.T) {
				scope, q, c := sqlMIInput()
				g := sqlMIKnownGroup("n")
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
				c.Groups[g.Version+"/"+g.Edition+"/"+g.Family+"/"+g.Name] = g
				out, e := c.Catalog(context.Background(), scope, q, nil)
				if n == 512 {
					if e != nil || out.Status != "complete" {
						t.Fatal("exact component/metadata label", out.Status, e)
					}
				} else {
					sqlMIRequireProjectionLimit(t, out, e)
				}
			})
		}
	}
}
func TestSQLMIProjectionWholeValidation(t *testing.T) {
	for _, mode := range []string{"negative-capacity", "known-zero-capacity", "known-empty-sku", "unknown-negative", "key-coordinate", "bad-state", "unknown-state", "metadata-case", "key-case", "wrong-tier", "bad-flag", "missing-family-coordinate", "unknown-with-groups", "unknown-without-code", "complete-with-code"} {
		scope, q, c := sqlMIInput()
		g := sqlMIKnownGroup("s")
		key := "v/e/f/s"
		switch mode {
		case "negative-capacity":
			g.Capacity = -1
		case "known-zero-capacity":
			g.Capacity = 0
		case "known-empty-sku":
			g.SKUName = ""
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
			key = "V/e/f/s"
		case "wrong-tier":
			g.SKUTier = "other"
		case "bad-flag":
			g.Standalone = "not_applicable"
		case "missing-family-coordinate":
			g.Known = false
			g.State = 2
			g.Family = ""
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
	scope, q, c := sqlMIInput()
	g := sqlMIKnownGroup("s")
	g.Known = false
	g.State = 2
	c.Groups["v/e/f/s"] = g
	for _, required := range [][]string{nil, {"v/e/f/s"}, {"v/e/f/missing"}} {
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

// For fixed5byte name and unknown modes, key+metadata=2V+3E+2F+S+27.
func sqlMITextGroup(t *testing.T, name string, target int) SQLCapabilityGroup {
	t.Helper()
	for v := 1; v <= 512; v++ {
		for e := 1; e <= 512; e++ {
			rest := target - 2*v - 3*e - 27
			if rest < 3 || rest > 1535 {
				continue
			}
			f := 1
			if rest > 513 {
				f = (rest - 511 + 1) / 2
			}
			s := rest - 2*f
			if f < 1 || f > 512 || s < 1 || s > 511 {
				continue
			}
			return SQLCapabilityGroup{Known: true, State: 0, Version: strings.Repeat("v", v), Edition: strings.Repeat("e", e), Name: name, Family: strings.Repeat("f", f), SKUName: strings.Repeat("s", s), SKUTier: strings.Repeat("e", e), Capacity: 1, Standalone: "unknown", InstancePool: "unknown"}
		}
	}
	t.Fatal("unreachable MI residual", target)
	return SQLCapabilityGroup{}
}
func sqlMITextFixture(t *testing.T, requiredMode, aliasMode bool) (map[string]string, SQLCapabilityRequest, *SQLCapabilityCollection, []string, string) {
	t.Helper()
	scope, q, c := sqlMIInput()
	required := []string{}
	overhead := 196
	if requiredMode {
		required = []string{"v/e/f/missing"}
		overhead += 13
	}
	if aliasMode {
		q.LocationAliases = map[string]string{"east us": "eastus"}
		overhead += 13
	}
	if sqlIndependentInputBytes(scope, q, c, required) != overhead {
		t.Fatal("fixed MI input overhead")
	}
	remaining := 16711680 - overhead
	i := 0
	for remaining >= 4123+35 {
		g := SQLCapabilityGroup{Known: true, State: 0, Version: strings.Repeat("v", 512), Edition: strings.Repeat("e", 512), Family: strings.Repeat("f", 512), Name: fmt.Sprintf("s%04d", i), SKUName: strings.Repeat("s", 512), SKUTier: strings.Repeat("e", 512), Capacity: 1, Standalone: "unknown", InstancePool: "unknown"}
		c.Groups[g.Version+"/"+g.Edition+"/"+g.Family+"/"+g.Name] = g
		remaining -= 4123
		i++
	}
	g := sqlMITextGroup(t, fmt.Sprintf("s%04d", i), remaining)
	key := g.Version + "/" + g.Edition + "/" + g.Family + "/" + g.Name
	c.Groups[key] = g
	if len(c.Groups) > 8192 || len(g.SKUName) >= 512 || sqlIndependentInputBytes(scope, q, c, required) != 16711680 {
		t.Fatal("exact valid MI input fixture")
	}
	return scope, q, c, required, key
}
func sqlMIRunTextBoundary(t *testing.T, requiredMode, aliasMode bool) {
	scope, q, c, required, key := sqlMITextFixture(t, requiredMode, aliasMode)
	out, e := c.Catalog(context.Background(), scope, q, required)
	want := "complete"
	if requiredMode {
		want = "unknown"
	}
	if e != nil || out.Status != want {
		t.Fatal("exact independent MI input text", out.Status, e)
	}
	g := c.Groups[key]
	g.SKUName += "s"
	c.Groups[key] = g
	if sqlIndependentInputBytes(scope, q, c, required) != 16711681 {
		t.Fatal("MI one-over accounting")
	}
	out, e = c.Catalog(context.Background(), scope, q, required)
	sqlMIRequireProjectionLimit(t, out, e)
}
func TestSQLMIProjectionTextBound(t *testing.T) { sqlMIRunTextBoundary(t, false, false) }
func TestSQLMIRequiredTextBound(t *testing.T)   { sqlMIRunTextBoundary(t, true, false) }
func TestSQLMIAliasTextBound(t *testing.T)      { sqlMIRunTextBoundary(t, false, true) }
func TestSQLMIOutputBudgetProof(t *testing.T) {
	maximum := 0
	for n := 1; n <= 8192; n++ {
		a := 281 + n*5145
		b := 16711680 + 39 - n*7
		if b < a {
			a = b
		}
		if a > maximum {
			maximum = a
		}
	}
	if maximum != 16689011 || maximum >= 16711680 {
		t.Fatal("MI output/input inequality", maximum)
	}
}
func TestSQLMIFullCoordinateBounds(t *testing.T) {
	scope, q, c := sqlMIInput()
	coordinate := strings.Repeat("v", 512) + "/" + strings.Repeat("e", 512) + "/" + strings.Repeat("f", 512) + "/" + strings.Repeat("n", 512)
	if len(coordinate) != 2051 {
		t.Fatal("MI full coordinate fixture")
	}
	out, e := c.Catalog(context.Background(), scope, q, []string{coordinate})
	if e != nil || out.Status != "unknown" {
		t.Fatal("exact MI2051coordinate", out, e)
	}
	out, e = c.Catalog(context.Background(), scope, q, []string{coordinate + "n"})
	sqlMIRequireProjectionLimit(t, out, e)
	out, e = c.Catalog(context.Background(), scope, q, []string{"v/e/" + strings.Repeat("f", 513) + "/n"})
	if e == nil || e.Error() != "SQL capability catalog projection rejected [sql_capability_projection_invalid]" {
		t.Fatal("MI component guard", out, e)
	}
}

func TestSQLMICoordinateDomains(t *testing.T) {
	scope, q, _ := sqlMIInput()
	var c *SQLCapabilityCollection
	for _, coordinate := range []string{"gp_gen5", "v/e/n"} {
		out, e := c.Catalog(context.Background(), scope, q, []string{coordinate})
		if e == nil || e.Error() != "SQL capability catalog projection rejected [sql_capability_projection_invalid]" || !reflect.DeepEqual(out, SQLCatalogEvidence{}) {
			t.Fatal("MI requires four explicit coordinates", coordinate, out, e)
		}
	}
}
