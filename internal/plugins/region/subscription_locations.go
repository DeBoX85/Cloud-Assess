// Subscription location semantics derive from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/SUBSCRIPTION_LOCATION_CONTRACT_20261028.md.
package region

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const subscriptionLocationVersion = "2022-12-01"

type SubscriptionLocationRequest struct{ SubscriptionID string }
type SubscriptionLocationEvidence struct {
	SubscriptionID, Status, FailureCode, CountSource string
	Regions                                          []string
	ZoneCounts                                       map[string]int
	ZoneMappings                                     map[string]map[string]string
}

// Exports are diagnostics; admission is a separately owned private snapshot.
type SubscriptionLocationCollection struct {
	SubscriptionID, Status, FailureCode, CountSource string
	Regions                                          []string
	ZoneCounts                                       map[string]int
	ZoneMappings                                     map[string]map[string]string
	state                                            *subscriptionLocationState
}
type subscriptionLocationState struct{ evidence SubscriptionLocationEvidence }
type SubscriptionLocationCollector struct{ transport *RESTQuotaCollector }

func locationGetterNil(g RESTQuotaGetter) bool {
	if g == nil {
		return true
	}
	v := reflect.ValueOf(g)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}
func NewSubscriptionLocationCollector(endpoint string, g RESTQuotaGetter) (*SubscriptionLocationCollector, error) {
	if locationGetterNil(g) {
		return nil, locationRequestFailure(nil)
	}
	transport, e := NewRESTQuotaCollector(endpoint, g)
	if e != nil {
		return nil, locationRequestFailure(nil)
	}
	return &SubscriptionLocationCollector{transport: transport}, nil
}
func locationRequestFailure(ctx context.Context) error {
	if ctx != nil && ctx.Err() != nil {
		return ctx.Err()
	}
	return errors.New("region subscription locations rejected [locations_request_invalid]")
}
func emptyLocationEvidence() SubscriptionLocationEvidence {
	return SubscriptionLocationEvidence{Status: "unknown", CountSource: "incomplete", Regions: []string{}, ZoneCounts: map[string]int{}, ZoneMappings: map[string]map[string]string{}}
}
func cloneLocationEvidence(ctx context.Context, in SubscriptionLocationEvidence) (SubscriptionLocationEvidence, error) {
	out := SubscriptionLocationEvidence{SubscriptionID: in.SubscriptionID, Status: in.Status, FailureCode: in.FailureCode, CountSource: in.CountSource, Regions: []string{}, ZoneCounts: map[string]int{}, ZoneMappings: map[string]map[string]string{}}
	check := func() error {
		if ctx != nil {
			return ctx.Err()
		}
		return nil
	}
	for _, name := range in.Regions {
		if e := check(); e != nil {
			return SubscriptionLocationEvidence{}, e
		}
		out.Regions = append(out.Regions, name)
	}
	for name, n := range in.ZoneCounts {
		if e := check(); e != nil {
			return SubscriptionLocationEvidence{}, e
		}
		out.ZoneCounts[name] = n
	}
	for name, m := range in.ZoneMappings {
		if e := check(); e != nil {
			return SubscriptionLocationEvidence{}, e
		}
		copyMap := map[string]string{}
		for k, v := range m {
			if e := check(); e != nil {
				return SubscriptionLocationEvidence{}, e
			}
			copyMap[k] = v
		}
		out.ZoneMappings[name] = copyMap
	}
	if e := check(); e != nil {
		return SubscriptionLocationEvidence{}, e
	}
	return out, nil
}
func (c *SubscriptionLocationCollection) Evidence() SubscriptionLocationEvidence {
	if c == nil || c.state == nil {
		return emptyLocationEvidence()
	}
	out, _ := cloneLocationEvidence(nil, c.state.evidence)
	return out
}

// This is the unchanged source compatibility table, not a current registry.
func locationFallbackCount(name string) int {
	switch name {
	case "australiaeast", "brazilsouth", "canadacentral", "centralindia", "centralus", "chinaeast3", "chinanorth3", "eastasia", "eastus", "eastus2", "francecentral", "germanywestcentral", "israelcentral", "italynorth", "japaneast", "koreacentral", "mexicocentral", "newzealandnorth", "northeurope", "norwayeast", "polandcentral", "qatarcentral", "southafricanorth", "southcentralus", "southeastasia", "spaincentral", "swedencentral", "switzerlandnorth", "uaenorth", "uksouth", "westeurope", "westus2", "westus3":
		return 3
	}
	return 0
}
func publishLocations(ctx context.Context, id, status, code string, rows map[string]locationRow) (*SubscriptionLocationCollection, error) {
	evidence := emptyLocationEvidence()
	evidence.SubscriptionID = id
	evidence.Status = status
	evidence.FailureCode = code
	mapped := false
	for name, row := range rows {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		evidence.Regions = append(evidence.Regions, name)
		if len(row.mappings) > 0 {
			mapped = true
		}
	}
	sort.Strings(evidence.Regions)
	if status == "complete" {
		if mapped {
			evidence.CountSource = "api"
		} else {
			evidence.CountSource = "fallback"
		}
	}
	for _, name := range evidence.Regions {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		row := rows[name]
		evidence.ZoneCounts[name] = len(row.mappings)
		if evidence.CountSource == "fallback" {
			evidence.ZoneCounts[name] = locationFallbackCount(name)
		}
		if len(row.mappings) > 0 {
			evidence.ZoneMappings[name] = row.mappings
		}
	}
	diagnostics, e := cloneLocationEvidence(ctx, evidence)
	if e != nil {
		return nil, e
	}
	result := &SubscriptionLocationCollection{SubscriptionID: diagnostics.SubscriptionID, Status: diagnostics.Status, FailureCode: diagnostics.FailureCode, CountSource: diagnostics.CountSource, Regions: diagnostics.Regions, ZoneCounts: diagnostics.ZoneCounts, ZoneMappings: diagnostics.ZoneMappings, state: &subscriptionLocationState{evidence: evidence}}
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return result, nil
}
func (c *SubscriptionLocationCollector) continuation(first *url.URL, link string) (*url.URL, bool) {
	if len(link) > 8192 || strings.ContainsAny(link, "#\\ \r\n\t") {
		return nil, false
	}
	ref, e := url.Parse(link)
	if e != nil {
		return nil, false
	}
	u := first.ResolveReference(ref)
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil || u.Scheme != "https" || !strings.EqualFold(u.Host, first.Host) || u.Path != first.Path || u.RawPath != "" || u.Opaque != "" || u.User != nil || u.Fragment != "" || u.ForceQuery || len(q["api-version"]) != 1 || q.Get("api-version") != subscriptionLocationVersion {
		return nil, false
	}
	for key, values := range q {
		if (key != "api-version" && key != "$skiptoken") || len(values) != 1 {
			return nil, false
		}
		if key == "$skiptoken" && (values[0] == "" || len(values[0]) > 4096 || !utf8.ValidString(values[0]) || strings.IndexFunc(values[0], func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || regionNoncharacter(r) }) >= 0) {
			return nil, false
		}
	}
	u.RawQuery = q.Encode()
	u.Host = first.Host
	return u, true
}

// Collect has only per-call state; transport/Context cooperation remains required.
func (c *SubscriptionLocationCollector) Collect(ctx context.Context, selected map[string]string, q SubscriptionLocationRequest) (*SubscriptionLocationCollection, error) {
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if ctx == nil || c == nil || c.transport == nil || c.transport.origin == nil || locationGetterNil(c.transport.getter) {
		return nil, locationRequestFailure(ctx)
	}
	text := 0
	scope, valid := auxiliaryScope(ctx, selected, &text)
	if !valid {
		return nil, locationRequestFailure(ctx)
	}
	for _, name := range scope {
		if name == "" {
			return nil, locationRequestFailure(ctx)
		}
	}
	id := strings.ToLower(q.SubscriptionID)
	if len(id) != 36 || !subscriptionID.MatchString(id) || scope[id] == "" {
		return nil, locationRequestFailure(ctx)
	}
	first := *c.transport.origin
	first.Path = "/subscriptions/" + id + "/locations"
	first.RawQuery = "api-version=" + subscriptionLocationVersion
	link := first.String()
	seenLinks := map[string]bool{}
	seenNames := map[string]bool{}
	rows := map[string]locationRow{}
	pages, total, work := 0, 0, 0
	fail := func(code string) (*SubscriptionLocationCollection, error) {
		status := "unknown"
		if pages > 0 {
			status = "partial"
		}
		return publishLocations(ctx, id, status, code, rows)
	}
	for {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if pages >= MaxQuotaPages {
			return fail("locations_page_limit")
		}
		u, valid := c.continuation(&first, link)
		if !valid {
			return fail("locations_unsafe_continuation")
		}
		if seenLinks[u.String()] {
			return fail("locations_continuation_cycle")
		}
		seenLinks[u.String()] = true
		body, response, e := c.transport.getter.GetBoundedWithResponse(ctx, u.String(), MaxQuotaPageBytes)
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if errors.Is(e, context.Canceled) {
			return nil, context.Canceled
		}
		if errors.Is(e, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		if e != nil {
			return fail("locations_request_failed")
		}
		if response == nil || response.StatusCode != 200 {
			return fail("locations_status_invalid")
		}
		if len(body) > MaxQuotaPageBytes || len(body) > MaxQuotaTotalBytes-total {
			return fail("locations_byte_limit")
		}
		total += len(body)
		page, code := decodeLocationPage(ctx, body, id, seenNames, work, text)
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if code != "" {
			return fail(code)
		}
		for _, name := range page.names {
			if e := ctx.Err(); e != nil {
				return nil, e
			}
			seenNames[name] = true
		}
		for _, row := range page.rows {
			if e := ctx.Err(); e != nil {
				return nil, e
			}
			rows[row.name] = row
		}
		work, text = page.work, page.text
		pages++
		if page.next == "" {
			break
		}
		link = page.next
	}
	return publishLocations(ctx, id, "complete", "", rows)
}
