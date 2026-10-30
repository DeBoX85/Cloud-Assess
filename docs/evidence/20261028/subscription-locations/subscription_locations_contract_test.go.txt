package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
)

const locSub = "abcdefab-cdef-abcd-efab-cdefabcdefab"
const locRoot = "https://management.azure.com/subscriptions/" + locSub + "/locations?api-version=2022-12-01"

type locPage struct {
	Body   string
	Status int
	Error  string
}
type locLiteral struct {
	Name        string
	Pages       []locPage
	Regions     []string
	Counts      map[string]int
	Mappings    map[string]map[string]string
	Status      string
	CountSource string `json:"count_source"`
	Code        string
	Calls       int
}
type locSpec struct {
	Cases      []locLiteral
	Boundary   []string `json:"boundary_recipes"`
	Procedural []string
}

func locFixtures(t *testing.T) locSpec {
	t.Helper()
	b, e := os.ReadFile("../../../docs/evidence/20261028/subscription-locations/literals.json")
	if e != nil {
		t.Fatal(e)
	}
	var s locSpec
	// Metadata is inspected independently; decoding only the named data sections
	// avoids an implementation dependency in this future consumer.
	var raw map[string]json.RawMessage
	if e = json.Unmarshal(b, &raw); e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw["cases"], &s.Cases); e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw["boundary_recipes"], &s.Boundary); e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw["procedural"], &s.Procedural); e != nil {
		t.Fatal(e)
	}
	if len(s.Cases) != 35 || len(s.Boundary) != 8 || len(s.Procedural) != 8 {
		t.Fatal("LOC_FIXTURE_CARDINALITY")
	}
	return s
}
func TestSubscriptionLocationsFixtureShape(t *testing.T) {
	s := locFixtures(t)
	seen := map[string]bool{}
	for _, c := range s.Cases {
		if c.Name == "" || seen[c.Name] || c.Calls < 1 || len(c.Pages) != c.Calls {
			t.Fatal("LOC_FIXTURE_IDENTITY")
		}
		seen[c.Name] = true
	}
}

type locGetter struct {
	mu     sync.Mutex
	pages  []locPage
	urls   []string
	bounds []int64
	hook   func(int)
}

func (g *locGetter) GetBoundedWithResponse(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
	g.mu.Lock()
	i := len(g.urls)
	g.urls = append(g.urls, u)
	g.bounds = append(g.bounds, n)
	var p locPage
	if i < len(g.pages) {
		p = g.pages[i]
	} else {
		p.Error = "unexpected"
	}
	g.mu.Unlock()
	if g.hook != nil {
		g.hook(i)
	}
	switch p.Error {
	case "cancel":
		return nil, nil, fmt.Errorf("SENSITIVE_PROVIDER_TEXT: %w", context.Canceled)
	case "deadline":
		return nil, nil, fmt.Errorf("SENSITIVE_PROVIDER_TEXT: %w", context.DeadlineExceeded)
	case "":
	default:
		return nil, nil, errors.New("SENSITIVE_PROVIDER_TEXT")
	}
	var response *http.Response
	if p.Status != 0 {
		response = &http.Response{StatusCode: p.Status}
	}
	return []byte(strings.ReplaceAll(p.Body, "{SUB}", locSub)), response, nil
}
func (g *locGetter) calls() int { g.mu.Lock(); defer g.mu.Unlock(); return len(g.urls) }
func locCollect(t *testing.T, g *locGetter) (*SubscriptionLocationCollection, error) {
	t.Helper()
	c, e := NewSubscriptionLocationCollector("https://management.azure.com", g)
	if e != nil {
		t.Fatal(e)
	}
	return c.Collect(context.Background(), map[string]string{locSub: "A"}, SubscriptionLocationRequest{SubscriptionID: locSub})
}
func locWire(name string, mappings string) string {
	extra := ""
	if mappings != "" {
		extra = `,"availabilityZoneMappings":` + mappings
	}
	return `{"name":"` + name + `","metadata":{"regionType":"Physical"}` + extra + `}`
}
func locBody(rows string, next string) string {
	tail := ""
	if next != "" {
		b, _ := json.Marshal(next)
		tail = `,"nextLink":` + string(b)
	}
	return `{"value":[` + rows + `]` + tail + `}`
}
func locLink(i int) string { return fmt.Sprintf("?api-version=2022-12-01&$skiptoken=p%d", i) }
func TestSubscriptionLocationsFrozenLiterals(t *testing.T) {
	for _, item := range locFixtures(t).Cases {
		c := item
		t.Run(c.Name, func(t *testing.T) {
			g := &locGetter{pages: c.Pages}
			got, e := locCollect(t, g)
			if e != nil || got == nil {
				t.Fatalf("LOC_LITERAL_RESULT: %v", e)
			}
			v := got.Evidence()
			if v.Status != c.Status || v.FailureCode != c.Code || v.CountSource != c.CountSource || v.SubscriptionID != locSub {
				t.Fatalf("LOC_LITERAL_STATUS: got=%#v want=%#v", v, c)
			}
			if !reflect.DeepEqual(v.Regions, c.Regions) || !reflect.DeepEqual(v.ZoneCounts, c.Counts) || !reflect.DeepEqual(v.ZoneMappings, c.Mappings) {
				t.Fatalf("LOC_LITERAL_VALUES: got=%#v want=%#v", v, c)
			}
			if v.Regions == nil || v.ZoneCounts == nil || v.ZoneMappings == nil || g.calls() != c.Calls {
				t.Fatalf("LOC_LITERAL_CALLS: %#v calls=%d", v, g.calls())
			}
			if g.urls[0] != locRoot {
				t.Fatal("LOC_REQUEST_BINDING", g.urls[0])
			}
			for _, n := range g.bounds {
				if n != 1<<20 {
					t.Fatal("LOC_REQUEST_BOUND")
				}
			}
			if len(g.urls) > 1 && strings.Contains(c.Pages[0].Body, "second") && g.urls[1] != "https://management.azure.com/subscriptions/"+locSub+"/locations?%24skiptoken=second&api-version=2022-12-01" {
				t.Fatal("LOC_CONTINUATION_BINDING", g.urls[1])
			}
			if got.Status != v.Status || got.FailureCode != v.FailureCode || got.CountSource != v.CountSource || got.SubscriptionID != v.SubscriptionID || !reflect.DeepEqual(got.Regions, v.Regions) || !reflect.DeepEqual(got.ZoneCounts, v.ZoneCounts) || !reflect.DeepEqual(got.ZoneMappings, v.ZoneMappings) {
				t.Fatal("LOC_DIAGNOSTIC_PARITY")
			}
		})
	}
}
func TestSubscriptionLocationsAllPinnedFallbacks(t *testing.T) {
	b, e := os.ReadFile("../../../docs/evidence/20261028/subscription-locations/zone-counts-pinned.json")
	if e != nil {
		t.Fatal(e)
	}
	var d struct{ Counts map[string]int }
	if json.Unmarshal(b, &d) != nil || len(d.Counts) != 33 {
		t.Fatal("LOC_PINNED_CARDINALITY")
	}
	for name, n := range d.Counts {
		g := &locGetter{pages: []locPage{{Body: locBody(locWire(name, ""), ""), Status: 200}}}
		v, e := locCollect(t, g)
		if e != nil || v == nil || v.Evidence().CountSource != "fallback" || v.Evidence().ZoneCounts[name] != n || len(v.Evidence().ZoneMappings) != 0 {
			t.Fatal("LOC_PINNED_FALLBACK", name)
		}
	}
}
