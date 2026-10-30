package region

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

const rhA = "abcdefab-cdef-abcd-efab-cdefabcdefab"
const rhB = "12345678-1234-1234-1234-123456789abc"

func rhRequest() MeterHistoryRequest {
	return MeterHistoryRequest{SubscriptionID: rhA, Months: 1, Now: time.Date(2026, 10, 29, 1, 2, 3, 456, time.FixedZone("offset", 3600)), Currency: "USD"}
}
func rhScope() map[string]string { return map[string]string{rhA: "a", rhB: "b"} }
func rhCollection(currency, status, code string, meters []HistoricalCostMeter) *MeterHistoryCollection {
	return &MeterHistoryCollection{SubscriptionID: rhA, Currency: currency, Status: status, FailureCode: code, Meters: meters}
}

type rhCase struct {
	name, currency, status, code, id string
	meters, want                     []HistoricalCostMeter
	reject, complete                 bool
	observed                         int
}

func rhCases() []rhCase {
	return []rhCase{
		{name: "usd", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 3}}, want: []HistoricalCostMeter{{"m", 3}}, complete: true, observed: 1},
		{name: "nok-kept", currency: "NOK", status: "complete", meters: []HistoricalCostMeter{{"m", 3}}, want: []HistoricalCostMeter{{"m", 3}}, complete: true, observed: 1},
		{name: "syntactic-zzz", currency: "ZZZ", status: "complete", meters: []HistoricalCostMeter{{"m", 3}}, want: []HistoricalCostMeter{{"m", 3}}, complete: true, observed: 1},
		{name: "complete-empty", currency: "USD", status: "complete", want: []HistoricalCostMeter{}, complete: true},
		{name: "partial-prefix", currency: "NOK", status: "partial", code: "history_limit", meters: []HistoricalCostMeter{{"m", 3}}, want: []HistoricalCostMeter{}, observed: 1},
		{name: "partial-empty", currency: "USD", status: "partial", code: "history_request_failed", want: []HistoricalCostMeter{}},
		{name: "unknown-empty", currency: "USD", status: "unknown", code: "history_request_failed", want: []HistoricalCostMeter{}},
		{name: "opaque-case", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"M", 3}, {"m", 7}}, want: []HistoricalCostMeter{{"M", 3}, {"m", 7}}, complete: true, observed: 2},
		{name: "opaque-unicode", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"é", 3}, {"e\u0301", 7}}, want: []HistoricalCostMeter{{"é", 3}, {"e\u0301", 7}}, complete: true, observed: 2},
		{name: "zero-kept", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 0}}, want: []HistoricalCostMeter{{"m", 0}}, complete: true, observed: 1},
		{name: "min-weight", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 1e-12}}, want: []HistoricalCostMeter{{"m", 1e-12}}, complete: true, observed: 1},
		{name: "max-weight", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 1e12}}, want: []HistoricalCostMeter{{"m", 1e12}}, complete: true, observed: 1},
		{name: "lower-currency", currency: "usd", status: "complete", reject: true},
		{name: "zero-status", currency: "USD", reject: true},
		{name: "status-alias", currency: "USD", status: "Complete", reject: true},
		{name: "complete-code", currency: "USD", status: "complete", code: "history_limit", reject: true},
		{name: "partial-no-code", currency: "USD", status: "partial", reject: true},
		{name: "unknown-no-code", currency: "USD", status: "unknown", reject: true},
		{name: "alien-code", currency: "USD", status: "partial", code: "SENSITIVE_provider_secret", reject: true},
		{name: "unknown-prefix", currency: "USD", status: "unknown", code: "history_limit", meters: []HistoricalCostMeter{{"m", 3}}, reject: true},
		{name: "foreign-id", currency: "USD", status: "complete", id: rhB, reject: true},
		{name: "duplicate", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 3}, {"m", 7}}, reject: true},
		{name: "duplicate-equal", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 3}, {"m", 3}}, reject: true},
		{name: "empty-meter", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"", 3}}, reject: true},
		{name: "invalid-utf8", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{string([]byte{255}), 3}}, reject: true},
		{name: "control-meter", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m\n", 3}}, reject: true},
		{name: "replacement-meter", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"\ufffd", 3}}, reject: true},
		{name: "noncharacter-meter", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"\uffff", 3}}, reject: true},
		{name: "negative", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", -1}}, reject: true},
		{name: "nan", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", math.NaN()}}, reject: true},
		{name: "positive-inf", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", math.Inf(1)}}, reject: true},
		{name: "negative-inf", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", math.Inf(-1)}}, reject: true},
		{name: "tiny-positive", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 5e-13}}, reject: true},
		{name: "too-large", currency: "USD", status: "complete", meters: []HistoricalCostMeter{{"m", 1e12 + 1}}, reject: true},
		{name: "partial-invalid-weight", currency: "USD", status: "partial", code: "history_limit", meters: []HistoricalCostMeter{{"m", -1}}, reject: true},
	}
}
func rhRejected(t *testing.T, got *RunHistoryBinding, err error) {
	t.Helper()
	if got != nil || err == nil || err.Error() != "region run history binding rejected [run_history_invalid]" {
		t.Fatal("RH_LITERAL_REJECTION", got, err)
	}
}
func TestRunHistoryFixtureShape(t *testing.T) {
	raw, e := os.ReadFile("../../../docs/evidence/20261029/run-history-binding/test-plan.json")
	if e != nil {
		t.Fatal(e)
	}
	var p struct {
		CaseNames []string `json:"case_names"`
	}
	if json.Unmarshal(raw, &p) != nil {
		t.Fatal("RH_PLAN_JSON")
	}
	names := []string{}
	seen := map[string]bool{}
	for _, c := range rhCases() {
		if seen[c.name] {
			t.Fatal("RH_DUPLICATE_CASE")
		}
		seen[c.name] = true
		names = append(names, c.name)
	}
	if len(names) != 35 || !reflect.DeepEqual(names, p.CaseNames) {
		t.Fatal("RH_PLAN_SHAPE", names)
	}
}
func TestRunHistoryFrozenLiterals(t *testing.T) {
	for _, c := range rhCases() {
		t.Run(c.name, func(t *testing.T) {
			q := rhRequest()
			q.Currency = c.currency
			in := rhCollection(c.currency, c.status, c.code, c.meters)
			if c.id != "" {
				in.SubscriptionID = c.id
			}
			got, e := AdmitRunHistory(context.Background(), rhScope(), q, in)
			if c.reject {
				rhRejected(t, got, e)
				return
			}
			if e != nil || got == nil {
				t.Fatal("RH_LITERAL_RESULT", e)
			}
			v := got.Evidence()
			if !v.Admitted || v.Request.SubscriptionID != rhA || v.Request.Currency != c.currency || v.Request.Months != 1 || !v.Request.Now.Equal(q.Now) || v.Request.Now.Location() != time.UTC || v.Status != c.status || v.FailureCode != c.code || v.ObservedMeters != c.observed || v.History.Complete != c.complete || !reflect.DeepEqual(v.History.Meters, c.want) {
				t.Fatal("RH_LITERAL_VALUES", v)
			}
		})
	}
}
func TestRunHistoryCurrencyBinding(t *testing.T) {
	q := rhRequest()
	q.Currency = "NOK"
	in := rhCollection("USD", "complete", "", []HistoricalCostMeter{{"m", 3}})
	got, e := AdmitRunHistory(context.Background(), rhScope(), q, in)
	rhRejected(t, got, e)
}
func TestRunHistoryFailureCodes(t *testing.T) {
	for _, code := range []string{"history_limit", "history_continuation_invalid", "history_response_invalid", "history_request_failed", "history_decode_invalid"} {
		t.Run(code, func(t *testing.T) {
			got, e := AdmitRunHistory(context.Background(), rhScope(), rhRequest(), rhCollection("USD", "partial", code, nil))
			if e != nil || got == nil || got.Evidence().History.Complete || got.Evidence().FailureCode != code {
				t.Fatal("RH_FAILURE_CODE", got, e)
			}
		})
	}
}
func TestRunHistoryRequestAdmission(t *testing.T) {
	for _, name := range []string{"months-zero", "months-thirteen", "now-zero", "year-one", "currency-empty", "currency-lower", "currency-control", "bad-id", "unselected", "scope-alias", "scope-name-empty", "scope-bad-id", "nil-collection"} {
		t.Run(name, func(t *testing.T) {
			q := rhRequest()
			scope := rhScope()
			in := rhCollection("USD", "complete", "", nil)
			switch name {
			case "months-zero":
				q.Months = 0
			case "months-thirteen":
				q.Months = 13
			case "now-zero":
				q.Now = time.Time{}
			case "year-one":
				q.Now = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
			case "currency-empty":
				q.Currency = ""
			case "currency-lower":
				q.Currency = "usd"
			case "currency-control":
				q.Currency = "U\nD"
			case "bad-id":
				q.SubscriptionID = "SENSITIVE_bad"
			case "unselected":
				delete(scope, rhA)
			case "scope-alias":
				scope[strings.ToUpper(rhA)] = "a"
			case "scope-name-empty":
				scope[rhA] = ""
			case "scope-bad-id":
				scope["bad"] = "SENSITIVE"
			case "nil-collection":
				in = nil
			}
			got, e := AdmitRunHistory(context.Background(), scope, q, in)
			rhRejected(t, got, e)
		})
	}
}
func TestRunHistoryNormalizedRequest(t *testing.T) {
	q := rhRequest()
	q.SubscriptionID = strings.ToUpper(rhA)
	in := rhCollection("USD", "complete", "", nil)
	in.SubscriptionID = strings.ToUpper(rhA)
	got, e := AdmitRunHistory(context.Background(), map[string]string{strings.ToUpper(rhA): "a"}, q, in)
	if e != nil || got == nil || got.Evidence().Request.SubscriptionID != rhA || got.Evidence().Request.Now.Location() != time.UTC {
		t.Fatal("RH_REQUEST_BINDING", got, e)
	}
}
