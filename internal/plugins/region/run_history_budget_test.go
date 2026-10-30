package region

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestRunHistoryFrozenBudgets(t *testing.T) {
	for _, name := range []string{"meters-at8192", "meters-plus-one8193", "scope-at1000", "scope-plus-one1001", "meter-bytes-at512", "meter-bytes-plus-one513"} {
		t.Run(name, func(t *testing.T) {
			scope := rhScope()
			meters := []HistoricalCostMeter{{"m", 3}}
			reject := false
			switch name {
			case "meters-at8192", "meters-plus-one8193":
				n := 8192
				if name == "meters-plus-one8193" {
					n++
					reject = true
				}
				meters = make([]HistoricalCostMeter, n)
				for i := range meters {
					meters[i] = HistoricalCostMeter{fmt.Sprintf("m%04d", i), 3}
				}
			case "scope-at1000", "scope-plus-one1001":
				n := 1000
				if name == "scope-plus-one1001" {
					n++
					reject = true
				}
				scope = map[string]string{rhA: "a"}
				for i := 0; i < n-1; i++ {
					scope[fmt.Sprintf("%08x-1111-1111-1111-%012x", i, i)] = "s"
				}
			case "meter-bytes-at512", "meter-bytes-plus-one513":
				n := 512
				if name == "meter-bytes-plus-one513" {
					n++
					reject = true
				}
				meters = []HistoricalCostMeter{{strings.Repeat("m", n), 3}}
			}
			got, e := AdmitRunHistory(context.Background(), scope, rhRequest(), rhCollection("USD", "complete", "", meters))
			if reject {
				rhRejected(t, got, e)
				return
			}
			if e != nil || got == nil || !got.Evidence().Admitted || !got.Evidence().History.Complete || len(got.Evidence().History.Meters) != len(meters) {
				t.Fatal("RH_BUDGET_LITERAL", name, got, e)
			}
		})
	}
}
