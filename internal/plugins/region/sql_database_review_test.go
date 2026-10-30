package region

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestSQLDatabaseExactPinnedShapeTypes(t *testing.T) {
	type field struct {
		Type   string `json:"type"`
		Format string `json:"format"`
		Ref    string `json:"$ref"`
		Items  struct {
			Type string `json:"type"`
			Ref  string `json:"$ref"`
		} `json:"items"`
	}
	var schema struct {
		Definitions map[string]struct {
			Properties map[string]field `json:"properties"`
		} `json:"definitions"`
	}
	var common struct {
		Definitions map[string]struct {
			Properties map[string]field `json:"properties"`
		} `json:"definitions"`
	}
	if json.Unmarshal(vmRead(t, "sql-primary-2021-11-01-schema.json"), &schema) != nil || json.Unmarshal(vmRead(t, "sql-primary-2021-11-01-common.json"), &common) != nil {
		t.Fatal("pinned schema decode")
	}
	schema.Definitions["Sku"] = common.Definitions["Sku"]
	for name, definition := range schema.Definitions {
		for key, f := range definition.Properties {
			expected := ""
			if f.Ref != "" {
				parts := strings.Split(f.Ref, "/")
				expected = parts[len(parts)-1]
			} else {
				switch f.Type {
				case "string":
					expected = "string"
				case "boolean":
					expected = "bool"
				case "integer":
					if f.Format != "int32" {
						t.Fatal("unexpected pinned integer format")
					}
					expected = "int32"
				case "number":
					if f.Format != "double" {
						t.Fatal("unexpected pinned number format")
					}
					expected = "double"
				case "array":
					parts := strings.Split(f.Items.Ref, "/")
					expected = "[]" + parts[len(parts)-1]
				default:
					t.Fatal("unsupported pinned field", name, key)
				}
			}
			if sqlCapabilityShapes[name][key] != expected {
				t.Fatal("pinned wire shape type mismatch", name, key, sqlCapabilityShapes[name][key], expected)
			}
		}
	}
}
func TestSQLDatabaseRequiredCoordinateBounds(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	for _, n := range []int{512, 513} {
		required := []string{strings.Repeat("v", n) + "/" + strings.Repeat("e", n) + "/" + strings.Repeat("s", n)}
		out, e := c.Catalog(context.Background(), scope, q, required)
		if n == 512 {
			if e != nil || out.Status != "unknown" {
				t.Fatal("full coordinate component exact", out, e)
			}
		} else if e == nil || e.Error() != "SQL capability catalog projection rejected [sql_capability_projection_invalid]" {
			t.Fatal("invalid oversized coordinate components", out, e)
		}
	}
	out, e := c.Catalog(context.Background(), scope, q, []string{strings.Repeat("x", 2052)})
	sqlRequireProjectionLimit(t, out, e)
}
func TestSQLDatabaseAliasLabelBounds(t *testing.T) {
	scope, q, c := sqlBoundaryInput()
	for _, n := range []int{512, 513} {
		q.LocationAliases = map[string]string{"e" + strings.Repeat(" ", n-6) + "astus": "eastus"}
		out, e := c.Catalog(context.Background(), scope, q, nil)
		if n == 512 {
			if e != nil || out.Status != "complete" {
				t.Fatal("exact alias label", out, e)
			}
		} else {
			sqlRequireProjectionLimit(t, out, e)
		}
	}
}
