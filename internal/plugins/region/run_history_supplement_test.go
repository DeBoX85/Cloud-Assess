package region

import (
	"context"
	"strings"
	"testing"
)

func TestRunHistoryReturnedCurrencySpelling(t *testing.T) {
	got, err := AdmitRunHistory(context.Background(), rhScope(), rhRequest(), rhCollection("usd", "complete", "", []HistoricalCostMeter{{"m", 3}}))
	if got != nil || err == nil || err.Error() != "region run history binding rejected [run_history_invalid]" {
		t.Fatal("RH_CURRENCY_EXACT", got, err)
	}
}

func TestRunHistoryUTF8ByteBoundary(t *testing.T) {
	for _, tc := range []struct {
		name    string
		repeats int
		reject  bool
	}{
		{"at512bytes", 256, false}, {"above514bytes", 257, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id := strings.Repeat("é", tc.repeats)
			got, err := AdmitRunHistory(context.Background(), rhScope(), rhRequest(), rhCollection("USD", "complete", "", []HistoricalCostMeter{{id, 3}}))
			if tc.reject {
				if got != nil || err == nil || err.Error() != "region run history binding rejected [run_history_invalid]" {
					t.Fatal("RH_UTF8_BYTES", got, err)
				}
				return
			}
			if err != nil || got == nil {
				t.Fatal("RH_UTF8_BYTES", got, err)
			}
			e := got.Evidence()
			if !e.Admitted || !e.History.Complete || len(e.History.Meters) != 1 || e.History.Meters[0].MeterID != id || e.History.Meters[0].HistoricalCost != 3 {
				t.Fatal("RH_UTF8_BYTES", e)
			}
		})
	}
}

func TestRunHistoryOtherScopeAdmission(t *testing.T) {
	scope := rhScope()
	scope[rhB] = ""
	got, err := AdmitRunHistory(context.Background(), scope, rhRequest(), rhCollection("USD", "complete", "", []HistoricalCostMeter{{"m", 3}}))
	if got != nil || err == nil || err.Error() != "region run history binding rejected [run_history_invalid]" {
		t.Fatal("RH_SCOPE_ALL_NAMES", got, err)
	}
}
