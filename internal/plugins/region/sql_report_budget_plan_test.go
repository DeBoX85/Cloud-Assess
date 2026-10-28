package region

import (
	"context"
	"fmt"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"strings"
	"testing"
)

// This independent oracle counts the documented existing PluginTable text
// occurrences, never calls the future report code and freezes the8MiB boundary.
func sqlReportOracleText(table assessment.PluginTable) int {
	total := len(table.SheetName) + len(table.Metadata.Version) + len(table.Metadata.Description) + len(table.Metadata.Author) + len(table.Metadata.License) + len(table.Description)
	for _, s := range table.Columns {
		total += len(s)
	}
	for _, r := range table.Rows {
		for _, s := range r.Cells {
			total += len(s)
		}
	}
	for _, w := range table.Health.Warnings {
		total += len(w.Message)
	}
	return total
}
func sqlReportBudgetFixture(t *testing.T, family string, over bool) (map[string]string, SQLCapabilityRequest, *SQLCapabilityCollection, []string) {
	t.Helper()
	kind := "microsoft.sql/servers/databases"
	standalone, pool := "not_applicable", "not_applicable"
	hardware := strings.Repeat("f", 512)
	if family == "mi" {
		kind = "microsoft.sql/managedinstances"
		standalone, pool = "unknown", "unknown"
		hardware = strings.Repeat("f", 193)
	}
	q := SQLCapabilityRequest{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: kind}
	scope := map[string]string{sqlReportOracleID: "selected"}
	c := &SQLCapabilityCollection{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: kind, Status: "complete", Groups: map[string]SQLCapabilityGroup{}}
	required := []string{}
	cells := [][]string{}
	for i := 0; i < 8192; i++ {
		name := fmt.Sprintf("n%04d", i)
		key := "v/e/" + name
		if family == "mi" {
			key = "v/e/" + hardware + "/" + name
		}
		c.Groups[key] = SQLCapabilityGroup{Known: true, State: 0, Version: "v", Edition: "e", Name: name, Family: hardware, SKUName: "s", SKUTier: "e", Capacity: 1, Standalone: standalone, InstancePool: pool}
		required = append(required, key)
		cells = append(cells, []string{sqlReportOracleID, "eastus", kind, key, "v", "e", name, hardware, "s", "e", "1", standalone, pool, "available", "not_assessed"})
	}
	want := sqlReportExpected(family, cells, false, "")
	remaining := 8388608 - sqlReportOracleText(want)
	if remaining <= 0 {
		t.Fatal("independent minimum fixture must be below report cap")
	}
	for i, key := range required {
		extra := remaining
		if extra > 511 {
			extra = 511
		}
		g := c.Groups[key]
		g.SKUName += strings.Repeat("s", extra)
		c.Groups[key] = g
		want.Rows[i].Cells[8] = g.SKUName
		remaining -= extra
	}
	if remaining != 0 || sqlReportOracleText(want) != 8388608 {
		t.Fatal("independent exact output reachability", remaining, sqlReportOracleText(want))
	}
	if over {
		changed := false
		for i, key := range required {
			g := c.Groups[key]
			if len(g.SKUName) < 512 {
				g.SKUName += "s"
				c.Groups[key] = g
				want.Rows[i].Cells[8] = g.SKUName
				changed = true
				break
			}
		}
		if !changed || sqlReportOracleText(want) != 8388609 {
			t.Fatal("independent one-over output reachability")
		}
	}
	// Existing frozen Catalog must admit BOTH fixtures, so inherited guards do not mask the new report boundary.
	catalog, err := c.Catalog(context.Background(), scope, q, required)
	if err != nil || catalog.Status != "complete" || len(catalog.Values) != 8192 {
		t.Fatal("masked report fixture", catalog.Status, len(catalog.Values), err)
	}
	return scope, q, c, required
}
func TestSQLReportBudgetReachability(t *testing.T) {
	for _, family := range []string{"db", "mi"} {
		for _, over := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", family, over), func(t *testing.T) { sqlReportBudgetFixture(t, family, over) })
		}
	}
}
