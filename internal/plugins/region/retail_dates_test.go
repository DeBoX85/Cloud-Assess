package region

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRetailStrictRFC3339(t *testing.T) {
	for _, date := range []string{"2026-01-01T0:00:00Z", "2026-01-01T00:00:00+24:00", "2026-01-01T00:00:00+00:60", "2026-01-01T00:00:00,1Z"} {
		var page map[string]any
		json.Unmarshal([]byte(retailFixture(t)), &page)
		page["Items"].([]any)[0].(map[string]any)["effectiveStartDate"] = date
		body, _ := json.Marshal(page)
		r, e, n := retailRun(t, func(int, string) ([]byte, *http.Response, error) { return body, &http.Response{StatusCode: 200}, nil })
		if e != nil || r == nil || r.Status != "unknown" || r.FailureCode != "retail_decode_invalid" || n != 1 {
			t.Fatal("non-RFC3339 timestamp admitted", date, r, e, n)
		}
	}
}

func TestRetailExactDateOrdering(t *testing.T) {
	for _, v := range []struct {
		dates  [2]string
		prices [2]float64
		want   float64
	}{
		{[2]string{"", "0000-01-01T00:00:00Z"}, [2]float64{2, 10}, 10},
		{[2]string{"2026-01-01T00:00:00.1234567891Z", "2026-01-01T00:00:00.1234567892Z"}, [2]float64{3, 4}, 4},
		{[2]string{"2026-01-01T01:00:00.1234567892+01:00", "2026-01-01T00:00:00.1234567891Z"}, [2]float64{4, 3}, 4},
		{[2]string{"2026-01-01T00:00:00Z", "2026-10-20T00:00:00.0000000001Z"}, [2]float64{2, 9}, 2},
	} {
		for _, reverse := range []bool{false, true} {
			r, e, _ := retailRun(t, func(n int, _ string) ([]byte, *http.Response, error) {
				if n == 1 {
					return []byte(retailFixture(t)), &http.Response{StatusCode: 200}, nil
				}
				items := []any{}
				for i, date := range v.dates {
					var p map[string]any
					json.Unmarshal([]byte(retailFixture(t)), &p)
					row := p["Items"].([]any)[0].(map[string]any)
					row["effectiveStartDate"] = date
					row["retailPrice"] = v.prices[i]
					row["unitPrice"] = v.prices[i]
					items = append(items, row)
				}
				if reverse {
					items[0], items[1] = items[1], items[0]
				}
				body, _ := json.Marshal(map[string]any{"Items": items})
				return body, &http.Response{StatusCode: 200}, nil
			})
			if e != nil || r.Status != "complete" {
				t.Fatal("exact parsed date selection", v, reverse, r, e)
			}
			_, prices := r.Evidence()
			if prices["m1"]["eastus"] != v.want {
				t.Fatal("date precision/missing sentinel", v, reverse, prices)
			}
		}
	}
}
