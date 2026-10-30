package region

import (
	"fmt"
	"strings"
	"testing"
)

func locPad(body string, n int) string {
	if len(body) > n {
		panic("test recipe too small")
	}
	return body + strings.Repeat(" ", n-len(body))
}
func TestSubscriptionLocationsFrozenBudgets(t *testing.T) {
	for _, name := range locFixtures(t).Boundary {
		t.Run(name, func(t *testing.T) {
			pages := []locPage{}
			status, code, source, calls := "complete", "", "fallback", 0
			switch name {
			case "pages-at64", "pages-plus-one65":
				for i := 0; i < 64; i++ {
					next := ""
					if i < 63 || name == "pages-plus-one65" {
						next = locLink(i + 1)
					}
					pages = append(pages, locPage{Body: locBody("", next), Status: 200})
				}
				calls = 64
				if name == "pages-plus-one65" {
					status, code, source = "partial", "locations_page_limit", "incomplete"
				}
			case "body-at1MiB", "body-plus-one-byte":
				n := 1 << 20
				if name == "body-plus-one-byte" {
					n++
					status, code, source = "unknown", "locations_byte_limit", "incomplete"
				}
				pages = []locPage{{Body: locPad(locBody("", ""), n), Status: 200}}
				calls = 1
			case "aggregate-at8MiB", "aggregate-plus-one-body":
				for i := 0; i < 8; i++ {
					next := ""
					if i < 7 || name == "aggregate-plus-one-body" {
						next = locLink(i + 1)
					}
					pages = append(pages, locPage{Body: locPad(locBody("", next), 1<<20), Status: 200})
				}
				calls = 8
				if name == "aggregate-plus-one-body" {
					pages = append(pages, locPage{Body: locBody("", ""), Status: 200})
					calls = 9
					status, code, source = "partial", "locations_byte_limit", "incomplete"
				}
			case "work-at8192", "work-plus-one8193":
				n := 8192
				if name == "work-plus-one8193" {
					n++
					status, code, source = "unknown", "locations_work_limit", "incomplete"
				}
				rows := make([]string, n)
				for i := range rows {
					rows[i] = fmt.Sprintf(`{"name":"r%d","metadata":{"regionType":"Logical"}}`, i)
				}
				pages = []locPage{{Body: locBody(strings.Join(rows, ","), ""), Status: 200}}
				calls = 1
			default:
				t.Fatal("LOC_UNKNOWN_RECIPE")
			}
			g := &locGetter{pages: pages}
			got, e := locCollect(t, g)
			if e != nil || got == nil {
				t.Fatalf("LOC_BUDGET_RESULT: %v", e)
			}
			v := got.Evidence()
			if v.Status != status || v.FailureCode != code || v.CountSource != source || g.calls() != calls || len(v.Regions) != 0 {
				t.Fatalf("LOC_BUDGET_LITERAL: got=%#v calls=%d want=%s/%s/%s/%d", v, g.calls(), status, code, source, calls)
			}
		})
	}
}
func TestSubscriptionLocationsJSONAdmissionBoundaries(t *testing.T) {
	cases := []struct {
		name, body string
		valid      bool
	}{
		{"depth-at32", `{"value":[],"opaque":` + strings.Repeat("[", 31) + `null` + strings.Repeat("]", 31) + `}`, true},
		{"depth-plus-one33", `{"value":[],"opaque":` + strings.Repeat("[", 32) + `null` + strings.Repeat("]", 32) + `}`, false},
		{"string-at512", `{"value":[],"opaque":"` + strings.Repeat("x", 512) + `"}`, true},
		{"string-plus-one513", `{"value":[],"opaque":"` + strings.Repeat("x", 513) + `"}`, false},
		{"tokens-at65536", `{"value":[],"opaque":[` + strings.Repeat("null,", 65532) + `null]}`, true},
		{"tokens-plus-one65537", `{"value":[],"opaque":[` + strings.Repeat("null,", 65533) + `null]}`, false},
		{"invalid-utf8", `{"value":[],"opaque":"` + string([]byte{255}) + `"}`, false},
		{"escaped-unpaired-surrogate", `{"value":[],"opaque":"\ud800"}`, false},
		{"escaped-control", `{"value":[],"opaque":"\u000a"}`, false},
	}
	for _, n := range []int{127, 128} {
		b := `{"value":[]`
		for i := 0; i < n; i++ {
			b += fmt.Sprintf(`,"k%d":0`, i)
		}
		b += "}"
		cases = append(cases, struct {
			name, body string
			valid      bool
		}{fmt.Sprintf("keys-%d", n+1), b, n == 127})
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := &locGetter{pages: []locPage{{Body: c.body, Status: 200}}}
			r, e := locCollect(t, g)
			if e != nil || r == nil {
				t.Fatal("LOC_JSON_RESULT", e)
			}
			v := r.Evidence()
			wantStatus, wantCode := "complete", ""
			if !c.valid {
				wantStatus, wantCode = "unknown", "locations_invalid_response"
			}
			if v.Status != wantStatus || v.FailureCode != wantCode {
				t.Fatal("LOC_JSON_BOUNDARY", c.name, v)
			}
		})
	}
}
