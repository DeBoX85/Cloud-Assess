// Request hierarchy derived from MIT-licensed AZQR. See NOTICE.md and
// docs/REGION_RESERVATION_COLLECTOR.md for corrections and evidence boundaries.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const MaxReservationCalls = 256

type ReservationCollection struct {
	Evidence    ReservationEvidence
	FailureCode string
}
type ReservationCollector struct{ transport *RESTQuotaCollector }

func NewReservationCollector(endpoint string, getter RESTQuotaGetter) (*ReservationCollector, error) {
	c, e := NewRESTQuotaCollector(endpoint, getter)
	if e != nil {
		return nil, e
	}
	return &ReservationCollector{c}, nil
}

func reservationGroup(id, selected string) (string, bool) {
	sub, _, _, _, ok := reservationIdentity(id + "/capacityReservations/probe")
	return id, ok && sub == selected
}

func reservationVM(id, selected string) bool {
	if len(id) > MaxReservationIDBytes || !strings.HasPrefix(id, "/") || strings.ContainsAny(id, "\\?#%") {
		return false
	}
	p := strings.Split(id[1:], "/")
	if len(p) != 8 && len(p) != 10 {
		return false
	}
	for i, v := range map[int]string{0: "subscriptions", 2: "resourceGroups", 4: "providers", 5: "Microsoft.Compute"} {
		if len(p[i]) != len(v) || !strings.EqualFold(p[i], v) {
			return false
		}
	}
	if !strings.EqualFold(p[1], selected) {
		return false
	}
	if len(p) == 8 {
		if len(p[6]) != len("virtualMachines") || !strings.EqualFold(p[6], "virtualMachines") {
			return false
		}
	} else if len(p[6]) != len("virtualMachineScaleSets") || !strings.EqualFold(p[6], "virtualMachineScaleSets") || len(p[8]) != len("virtualMachines") || !strings.EqualFold(p[8], "virtualMachines") {
		return false
	}
	for _, i := range []int{3, 7, len(p) - 1} {
		n := 0
		if p[i] == "" || p[i] == "." || p[i] == ".." || !auxiliaryText(p[i], &n) {
			return false
		}
	}
	return true
}

type reservationItem struct{ ID, Name, Location string }
type reservationPage struct {
	Value    []reservationItem
	NextLink string
}

func decodeReservationPage(body []byte, group, subscription string, seen map[string]bool) (reservationPage, bool) {
	var raw struct {
		Value    json.RawMessage
		NextLink string
	}
	if json.Unmarshal(body, &raw) != nil || len(raw.Value) == 0 || bytes.Equal(bytes.TrimSpace(raw.Value), []byte("null")) {
		return reservationPage{}, false
	}
	var rows []*reservationItem
	if json.Unmarshal(raw.Value, &rows) != nil {
		return reservationPage{}, false
	}
	out := reservationPage{NextLink: raw.NextLink}
	page := map[string]bool{}
	for _, r := range rows {
		if r == nil {
			return reservationPage{}, false
		}
		if group == "" {
			if _, ok := reservationGroup(r.ID, subscription); !ok {
				return reservationPage{}, false
			}
			parts := strings.Split(r.ID, "/")
			if r.Name != "" && !strings.EqualFold(r.Name, parts[len(parts)-1]) {
				return reservationPage{}, false
			}
			r.Location = strings.ToLower(r.Location)
			if !regionID.MatchString(r.Location) {
				return reservationPage{}, false
			}
		} else {
			id := group + "/capacityReservations/" + r.Name
			sub, _, _, _, ok := reservationIdentity(id)
			if !ok || sub != subscription || r.Name == "" || (r.ID != "" && !strings.EqualFold(r.ID, id)) {
				return reservationPage{}, false
			}
			r.ID = id
		}
		key := strings.ToLower(r.ID)
		if page[key] || seen[key] {
			return reservationPage{}, false
		}
		page[key] = true
		out.Value = append(out.Value, *r)
	}
	for key := range page {
		seen[key] = true
	}
	return out, true
}

func decodeReservation(body []byte, id, region, subscription string) (ReservationUsage, int, bool) {
	var raw struct {
		ID, Name, Location string
		SKU                *struct {
			Name     string
			Capacity *int64
		}
		Properties *struct {
			InstanceView *struct {
				UtilizationInfo *struct{ VirtualMachinesAllocated json.RawMessage }
			}
		}
	}
	if bytes.Equal(bytes.TrimSpace(body), []byte("null")) || json.Unmarshal(body, &raw) != nil || (raw.ID != "" && !strings.EqualFold(raw.ID, id)) || (raw.Name != "" && !strings.EqualFold(raw.Name, id[strings.LastIndex(id, "/")+1:])) || (raw.Location != "" && strings.ToLower(raw.Location) != region) {
		return ReservationUsage{}, 0, false
	}
	r := ReservationUsage{ResourceID: id, Region: region, ResponseName: raw.Name, ResponseRegion: strings.ToLower(raw.Location)}
	if raw.SKU != nil {
		r.SKU = raw.SKU.Name
		if raw.SKU.Capacity != nil {
			r.Reserved = *raw.SKU.Capacity
			r.ReservedKnown = true
			if r.Reserved < 0 || r.Reserved > MaxAuxCount {
				return ReservationUsage{}, 0, false
			}
		}
	}
	count := 0
	if raw.Properties != nil && raw.Properties.InstanceView != nil && raw.Properties.InstanceView.UtilizationInfo != nil {
		b := raw.Properties.InstanceView.UtilizationInfo.VirtualMachinesAllocated
		if len(b) > 0 && !bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
			var refs []*struct{ ID string }
			if json.Unmarshal(b, &refs) != nil || len(refs) > MaxAuxRows {
				return ReservationUsage{}, 0, false
			}
			seen := map[string]bool{}
			for _, ref := range refs {
				if ref == nil || !reservationVM(ref.ID, subscription) || seen[strings.ToLower(ref.ID)] {
					return ReservationUsage{}, 0, false
				}
				seen[strings.ToLower(ref.ID)] = true
			}
			count = len(refs)
			r.Allocated = int64(count)
			r.AllocatedKnown = true
		}
	}
	return r, count, true
}

// Collect owns per-call state. Operational failures are explicit incomplete
// evidence; invalid caller scope and cancellation return nil result and error.
func (c *ReservationCollector) Collect(ctx context.Context, subscriptions map[string]string, request ReservationRequest) (*ReservationCollection, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	text := 0
	scope, ok := auxiliaryScope(ctx, subscriptions, &text)
	if !ok {
		return nil, reservationFailure(ctx, "scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return nil, reservationFailure(ctx, "scope_invalid")
		}
	}
	request.SubscriptionID = strings.ToLower(request.SubscriptionID)
	if !subscriptionID.MatchString(request.SubscriptionID) || len(request.SubscriptionID) != 36 || scope[request.SubscriptionID] == "" || !regionID.MatchString(request.Region) {
		return nil, reservationFailure(ctx, "request_invalid")
	}
	result := &ReservationCollection{Evidence: ReservationEvidence{Request: request, Status: "complete", Reservations: []ReservationUsage{}}}
	accepted := false
	calls, total, work := 0, 0, 0
	stop := false
	fail := func(code string) {
		if result.FailureCode == "" {
			result.FailureCode = code
		}
		result.Evidence.Status = "unknown"
		if accepted {
			result.Evidence.Status = "partial"
		}
	}
	urlFor := func(path string, expanded bool) *url.URL {
		u := *c.transport.origin
		u.Path = path
		u.RawQuery = "api-version=2024-11-01"
		if expanded {
			u.RawQuery = "%24expand=instanceView&api-version=2024-11-01"
		}
		return &u
	}
	get := func(u *url.URL) ([]byte, bool, error) {
		if e := ctx.Err(); e != nil {
			return nil, false, e
		}
		if calls >= MaxReservationCalls {
			fail("reservation_request_limit")
			stop = true
			return nil, false, nil
		}
		calls++
		b, response, e := c.transport.getter.GetBoundedWithResponse(ctx, u.String(), MaxQuotaPageBytes)
		if x := ctx.Err(); x != nil {
			return nil, false, x
		}
		if errors.Is(e, context.Canceled) || errors.Is(e, context.DeadlineExceeded) {
			return nil, false, e
		}
		if e != nil || response == nil || response.StatusCode != http.StatusOK {
			fail("reservation_request_failed")
			return nil, false, nil
		}
		if len(b) > MaxQuotaPageBytes || len(b) > MaxQuotaTotalBytes-total {
			fail("reservation_byte_limit")
			stop = true
			return nil, false, nil
		}
		total += len(b)
		if !quotaJSONWithIDs(ctx, b, true) {
			fail("reservation_invalid_response")
			return nil, false, nil
		}
		return b, true, nil
	}
	consume := func(n int) bool {
		if n > MaxAuxRows-work {
			fail("reservation_work_limit")
			stop = true
			return false
		}
		work += n
		return true
	}
	var list func(string, string, map[string]bool, func(reservationItem) error) error
	list = func(path, group string, seen map[string]bool, visit func(reservationItem) error) error {
		first := urlFor(path, false)
		link := first.String()
		links := map[string]bool{}
		for pages := 0; !stop; pages++ {
			if e := ctx.Err(); e != nil {
				return e
			}
			if pages >= MaxQuotaPages {
				fail("reservation_page_limit")
				return nil
			}
			u, valid := c.transport.continuation(first, "2024-11-01", link)
			if !valid || links[u.String()] {
				fail("reservation_unsafe_continuation")
				return nil
			}
			links[u.String()] = true
			b, valid, e := get(u)
			if e != nil {
				return e
			}
			if !valid {
				return nil
			}
			p, valid := decodeReservationPage(b, group, request.SubscriptionID, seen)
			if !valid {
				fail("reservation_identity_invalid")
				return nil
			}
			if !consume(len(p.Value)) {
				return nil
			}
			accepted = true
			for _, item := range p.Value {
				if e := ctx.Err(); e != nil {
					return e
				}
				if e := visit(item); e != nil {
					return e
				}
				if stop {
					return nil
				}
			}
			if p.NextLink == "" {
				return nil
			}
			link = p.NextLink
		}
		return nil
	}
	groups, reservations := map[string]bool{}, map[string]bool{}
	e := list("/subscriptions/"+request.SubscriptionID+"/providers/Microsoft.Compute/capacityReservationGroups", "", groups, func(group reservationItem) error {
		if group.Location != request.Region {
			return nil
		}
		return list(group.ID+"/capacityReservations", group.ID, reservations, func(summary reservationItem) error {
			b, valid, e := get(urlFor(summary.ID, true))
			if e != nil {
				return e
			}
			if !valid {
				return nil
			}
			r, refs, valid := decodeReservation(b, summary.ID, request.Region, request.SubscriptionID)
			if !valid {
				fail("reservation_identity_invalid")
				return nil
			}
			if !consume(refs) {
				return nil
			}
			result.Evidence.Reservations = append(result.Evidence.Reservations, r)
			return nil
		})
	})
	if e != nil {
		return nil, e
	}
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return result, nil
}
