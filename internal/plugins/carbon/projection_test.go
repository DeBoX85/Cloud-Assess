package carbon

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type typeFilter func(string) bool

func (f typeFilter) IsResourceTypeExcluded(v string) bool { return f(v) }
func ptr(v float64) *float64                              { return &v }
func sourceItems() []Item {
	return []Item{{"Microsoft.Compute/virtualMachines", ptr(80), ptr(100), ptr(-20)}, {"Microsoft.Compute/virtualMachines", ptr(20), ptr(0), nil}, {"Microsoft.Storage/storageAccounts", ptr(75), ptr(0), ptr(0)}}
}
func capture(t *testing.T, name string) [][]string {
	t.Helper()
	b, e := os.ReadFile("testdata/source-" + name + ".json")
	if e != nil {
		t.Fatal(e)
	}
	var v []struct {
		Metadata struct {
			Name, Version, Description, Author, License string
			Type                                        int
		}
		Sheet       string `json:"sheet_name"`
		Description string
		Table       [][]string
	}
	if e = json.Unmarshal(b, &v); e != nil || len(v) != 1 {
		t.Fatalf("source capture: %v", e)
	}
	m := Metadata()
	if v[0].Metadata.Name != m.Name || v[0].Metadata.Version != m.Version || v[0].Metadata.Description != m.Description || v[0].Metadata.Author != m.Author || v[0].Metadata.License != m.License || v[0].Metadata.Type != 1 || v[0].Sheet != "Carbon Emissions" || v[0].Description != PendingTable().Description {
		t.Fatal("source complete metadata/sheet")
	}
	return v[0].Table
}
func TestCarbonProjectionActualSourceCapturesAndCalculatedCells(t *testing.T) {
	for _, kind := range []string{"all", "filtered", "batched", "empty"} {
		t.Run(kind, func(t *testing.T) {
			items := sourceItems()
			var filter Filter
			if kind == "filtered" {
				filter = typeFilter(func(v string) bool { return v != "Microsoft.Compute/virtualMachines" })
			}
			if kind == "batched" {
				items = append(items, sourceItems()...)
			}
			if kind == "empty" {
				items = nil
			}
			got, e := Project(context.Background(), "2026-01-01", "2026-03-31", items, filter)
			if e != nil || got.Health.Status != assessment.StageCompleted || assessment.ValidatePluginTables([]assessment.PluginTable{got}) != nil {
				t.Fatalf("healthy canonical: %+v %v", got, e)
			}
			table := [][]string{got.Columns}
			for _, r := range got.Rows {
				if r.SubscriptionID != "" {
					t.Fatal("aggregate fabricated subscription identity")
				}
				table = append(table, r.Cells)
			}
			if !reflect.DeepEqual(table, capture(t, kind)) {
				t.Fatalf("source aggregate/filter/date cells: %+v", table)
			}
		})
	}
	b, e := os.ReadFile("testdata/source-calculations.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Latest, Previous, Change float64
		Row                      []string
	}
	if json.Unmarshal(b, &cases) != nil || len(cases) != 6 {
		t.Fatal("source calculations capture")
	}
	for i, c := range cases {
		got, e := Project(context.Background(), "2026-01-01", "2026-03-31", []Item{{"Microsoft.Compute/virtualMachines", ptr(c.Latest), ptr(c.Previous), ptr(c.Change)}}, nil)
		if e != nil || len(got.Rows) != 1 || !reflect.DeepEqual(got.Rows[0].Cells, c.Row) {
			t.Fatalf("actual source calculation %d: %+v %v", i, got, e)
		}
	}
	// Source hides a denied report as the same header-only table as empty scope.
	// This is characterized evidence, not target HTTP error acceptance.
	if !reflect.DeepEqual(capture(t, "empty"), capture(t, "denied")) {
		t.Fatal("source silent denied behavior changed")
	}
}
func TestCarbonIndependentCaptureHashes(t *testing.T) {
	for name, want := range map[string]string{"all": "bb14c9143d7590f431e629c67ca20fd2e47bd7d153472d41d8583ecd14f7391d", "filtered": "ba2dc5a0b4c9bdf35eed8997fa5789fc05b78cdb89b6865da4ca6d0e7efa4d07", "batched": "6966a5a96ec9b02974c6d6abc847732a041a22d0c8ce419b68d8afb2c87be3f8", "empty": "0152b1ec8081a786e4fbaa7673820a486ffbec411d7057c4cd5c2a4e5724f32f", "denied": "0152b1ec8081a786e4fbaa7673820a486ffbec411d7057c4cd5c2a4e5724f32f", "calculations": "0c9795aa4d6f063bfc8d9e3d3519a8be059c7e43a7a6b211a3c8b2a365fdb872"} {
		b, e := os.ReadFile("testdata/source-" + name + ".json")
		if e != nil || fmt.Sprintf("%x", sha256.Sum256(b)) != want {
			t.Fatalf("independent source capture hash %s", name)
		}
	}
}
func TestCarbonMalformedFiniteNumericBoundsAndDateValidation(t *testing.T) {
	bad := []Item{{"", ptr(1), nil, nil}, {"bad\nlabel", ptr(1), nil, nil}, {strings.Repeat("x", 513), ptr(1), nil, nil}, {"bad\xff", ptr(1), nil, nil}, {"bad\uffff", ptr(1), nil, nil}, {"nil", nil, nil, nil}, {"nan", ptr(math.NaN()), nil, nil}, {"inf", ptr(1), ptr(math.Inf(1)), nil}, {"change", ptr(1), nil, ptr(math.Inf(-1))}}
	good := sourceItems()
	items := append(good, bad...)
	got, e := Project(context.Background(), "2026-01-01", "2026-03-31", items, nil)
	if e != nil || len(got.Rows) != 2 || got.Health.Records != 2 || got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Message != "skipped 9 invalid carbon emission items" {
		t.Fatalf("malformed accounting: %+v %v", got, e)
	}
	for _, dates := range [][2]string{{"bad", "2026-03-31"}, {"2026-01-01", "2026-02-30"}, {"2026-04-01", "2026-03-31"}} {
		v, e := Project(context.Background(), dates[0], dates[1], good, nil)
		if e == nil || len(v.Rows) != 0 || v.Health.Error.Code != "carbon_date_invalid" {
			t.Fatal("invalid date marked healthy")
		}
	}
	overflow := []Item{{"", ptr(1), nil, nil}, {"same", ptr(math.MaxFloat64), nil, nil}, {"same", ptr(math.MaxFloat64), nil, nil}}
	got, e = Project(context.Background(), "2026-01-01", "2026-03-31", overflow, nil)
	if e == nil || len(got.Rows) != 1 || strings.Contains(strings.Join(got.Rows[0].Cells, "|"), "Inf") || got.Health.Error.Code != "carbon_numeric_limit" || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Message != "skipped 1 invalid carbon emission items" {
		t.Fatal("overflow lost prior data or formatted infinity")
	}
	ratio := []Item{{"ratio", ptr(math.MaxFloat64), ptr(math.SmallestNonzeroFloat64), nil}}
	got, e = Project(context.Background(), "2026-01-01", "2026-03-31", ratio, nil)
	if e == nil || len(got.Rows) != 0 || got.Health.Error.Code != "carbon_numeric_limit" {
		t.Fatal("ratio overflow not rejected")
	}
	for _, kind := range []string{"items", "types", "text"} {
		t.Run(kind, func(t *testing.T) {
			var items []Item
			code := "carbon_" + map[string]string{"items": "item", "types": "type", "text": "text"}[kind] + "_limit"
			switch kind {
			case "items":
				items = make([]Item, MaxItems+1)
			case "types":
				for i := 0; i <= MaxResourceTypes; i++ {
					items = append(items, Item{fmt.Sprintf("type-%04d", i), ptr(0), nil, nil})
				}
			case "text":
				for i := 0; i <= MaxTextBytes/MaxLabelBytes; i++ {
					items = append(items, Item{strings.Repeat("a", MaxLabelBytes), ptr(0), nil, nil})
				}
			}
			got, e := Project(context.Background(), "2026-01-01", "2026-03-31", items, nil)
			if e == nil || got.Health.Status != assessment.StageFailed || got.Health.Error.Code != code || assessment.ValidatePluginTables([]assessment.PluginTable{got}) != nil {
				t.Fatalf("bounded truthful failure %s: %+v %v", kind, got, e)
			}
		})
	}
}
func TestCarbonCancellationOwnershipAndConcurrentRuns(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	items := sourceItems()
	got, e := Project(ctx, "2026-01-01", "2026-03-31", items, typeFilter(func(string) bool { cancel(); return false }))
	if !errors.Is(e, context.Canceled) || got.Health.Error.Code != "carbon_cancelled" || len(got.Rows) != 1 || got.Rows[0].Cells[3] != "80.00" {
		t.Fatalf("context prior-row retention: %+v %v", got, e)
	}
	healthy, e := Project(context.Background(), "2026-01-01", "2026-03-31", items, nil)
	if e != nil {
		t.Fatal(e)
	}
	*items[0].Latest = 999
	healthy.Rows[0].Cells[3] = "mutated"
	healthy.Columns[0] = "mutated"
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, e := Project(context.Background(), "2026-01-01", "2026-03-31", sourceItems(), nil)
			if e != nil || v.Rows[0].Cells[3] != "100.00" || v.Columns[0] != "Period From" {
				t.Error("cross-run projection state")
			}
			v.Rows[0].Cells[3] = "mutated returned"
		}()
	}
	wg.Wait()
	a := PendingTable()
	a.Columns[0] = "mutated"
	if PendingTable().Columns[0] != "Period From" {
		t.Fatal("pending headers aliased")
	}
}
