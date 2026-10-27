package zone

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// Go's pinned JSON decoder matches struct fields with Unicode simple folding.
// Ambiguous provider pages must fail before any current-page rows are admitted.
func TestZoneRejectsUnicodeFoldDuplicateKeys(t *testing.T) {
	cases := []struct{ name, body string }{
		{"physical zone long s", strings.Replace(literalPage, `"physicalZone":"eastus-az3"`, `"physicalZone":"eastus-az3","phyſicalZone":"forged"`, 1)},
		{"mappings long s", `{"value":[{"name":"eastus","availabilityZoneMappings":[{"logicalZone":"1","physicalZone":"eastus-az1"}],"availabilityZoneMappingſ":[]}]}`},
		{"continuation Kelvin", strings.TrimSuffix(string(withNext(literalPage, "?api-version=2022-12-01&skiptoken=two")), "}") + `,"nextLinK":""}`},
		{"unknown sigma duplicate", strings.TrimSuffix(literalPage, "}") + `,"metadata":{"Σ":1,"ς":2}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			s := scanner(t, func(context.Context, string, int64) ([]byte, error) {
				calls++
				return []byte(tc.body), nil
			})
			got, err := s.Scan(context.Background(), map[string]string{subA: "A"})
			if err != nil || calls != 1 || len(got.Rows) != 0 || !reflect.DeepEqual(got.Failures, []Failure{{subA, "zone_invalid_response"}}) {
				t.Fatalf("Unicode fold duplicate admitted as zone evidence: rows=%+v failures=%+v calls=%d err=%v", got.Rows, got.Failures, calls, err)
			}
		})
	}
}

func TestZoneUniqueUnicodeMetadataRemainsValid(t *testing.T) {
	body := strings.TrimSuffix(literalPage, "}") + `,"metadata":{"Σ":1,"雪":2,"雨":3}}`
	s := scanner(t, func(context.Context, string, int64) ([]byte, error) { return []byte(body), nil })
	got, err := s.Scan(context.Background(), map[string]string{subA: "A"})
	want := []Row{{subA, "A", "eastus", "East US", "1", "eastus-az1"}, {subA, "A", "eastus", "East US", "2", "eastus-az3"}}
	if err != nil || len(got.Failures) != 0 || !reflect.DeepEqual(got.Rows, want) {
		t.Fatalf("unique Unicode metadata changed valid zone evidence: %+v %v", got, err)
	}
}
