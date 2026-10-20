package region

import (
	"cmp"
	"context"
	"errors"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RetailPriceItem retains the pinned source's retail fields without adding a
// public report schema. Original/target bindings are recorded separately.
type RetailPriceItem struct {
	CurrencyCode         string  `json:"currencyCode"`
	TierMinimumUnits     float64 `json:"tierMinimumUnits"`
	RetailPrice          float64 `json:"retailPrice"`
	UnitPrice            float64 `json:"unitPrice"`
	ArmRegionName        string  `json:"armRegionName"`
	Location             string  `json:"location"`
	EffectiveStartDate   string  `json:"effectiveStartDate"`
	MeterID              string  `json:"meterId"`
	MeterName            string  `json:"meterName"`
	ProductID            string  `json:"productId"`
	SkuID                string  `json:"skuId"`
	ProductName          string  `json:"productName"`
	SkuName              string  `json:"skuName"`
	ServiceName          string  `json:"serviceName"`
	ServiceID            string  `json:"serviceId"`
	ServiceFamily        string  `json:"serviceFamily"`
	UnitOfMeasure        string  `json:"unitOfMeasure"`
	Type                 string  `json:"type"`
	IsPrimaryMeterRegion bool    `json:"isPrimaryMeterRegion"`
	ArmSkuName           string  `json:"armSkuName"`
}

type RetailPricingRequest struct {
	MeterIDs, Regions []string
	Now               time.Time
}
type RetailMeter struct {
	OriginalMeterID string
	Item            RetailPriceItem
	Conflict        bool
}
type RetailPrice struct {
	OriginalMeterID string
	Item            RetailPriceItem
	Conflict        bool
}
type RetailBinding struct {
	OriginalMeterID string
	Item            RetailPriceItem
	Phase           string
}
type RetailUoMMismatch struct {
	OriginalMeterID, OriginalUnit string
	Item                          RetailPriceItem
}
type RetailCounts struct {
	Requests, RawItems, JSONTokens, Bindings, MatchingWork int
	ReceivedBytes                                          int64
}

type RetailPricingCollection struct {
	Status, FailureCode string
	Meters              []RetailMeter
	Prices              []RetailPrice
	RawItems            []RetailPriceItem
	Bindings            []RetailBinding
	UoMMismatches       []RetailUoMMismatch
	Counts              RetailCounts
	// Eligibility cannot be forged by mutating exported diagnostic fields.
	admission string
	eligible  map[string]map[string]float64
}

func (r *RetailPricingCollection) Evidence() (string, map[string]map[string]float64) {
	out := map[string]map[string]float64{}
	if r == nil {
		return "unavailable", out
	}
	if r.admission != "complete" {
		if r.admission == "partial" {
			return "partial", out
		}
		return "unavailable", out
	}
	for id, regions := range r.eligible {
		out[id] = map[string]float64{}
		for region, price := range regions {
			out[id][region] = price
		}
	}
	return "complete", out
}

type RetailCollector struct{ getter RetailGetter }

func NewRetailCollector(getter RetailGetter) (*RetailCollector, error) {
	if getter == nil {
		return nil, errRetailRequest
	}
	v := reflect.ValueOf(getter)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if v.IsNil() {
			return nil, errRetailRequest
		}
	}
	return &RetailCollector{getter: getter}, nil
}

type retailTuple struct{ MeterName, ProductID, SkuName string }

func retailKey(i RetailPriceItem) retailTuple {
	return retailTuple{i.MeterName, i.ProductID, i.SkuName}
}
func retailIdentity(i RetailPriceItem) RetailPriceItem {
	i.TierMinimumUnits = 0
	i.RetailPrice = 0
	i.UnitPrice = 0
	i.EffectiveStartDate = ""
	return i
}
func retailLess(a, b RetailPriceItem) bool {
	as := [...]string{a.CurrencyCode, a.ArmRegionName, a.Location, a.EffectiveStartDate, a.MeterID, a.MeterName, a.ProductID, a.SkuID, a.ProductName, a.SkuName, a.ServiceName, a.ServiceID, a.ServiceFamily, a.UnitOfMeasure, a.Type, a.ArmSkuName}
	bs := [...]string{b.CurrencyCode, b.ArmRegionName, b.Location, b.EffectiveStartDate, b.MeterID, b.MeterName, b.ProductID, b.SkuID, b.ProductName, b.SkuName, b.ServiceName, b.ServiceID, b.ServiceFamily, b.UnitOfMeasure, b.Type, b.ArmSkuName}
	for i := range as {
		if as[i] != bs[i] {
			return as[i] < bs[i]
		}
	}
	if n := cmp.Or(cmp.Compare(a.TierMinimumUnits, b.TierMinimumUnits), cmp.Compare(a.RetailPrice, b.RetailPrice), cmp.Compare(a.UnitPrice, b.UnitPrice)); n != 0 {
		return n < 0
	}
	return !a.IsPrimaryMeterRegion && b.IsPrimaryMeterRegion
}

func retailDate(i RetailPriceItem) time.Time {
	if i.EffectiveStartDate == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, i.EffectiveStartDate)
	return t.UTC()
}

func retailRequest(q RetailPricingRequest) ([]string, []string, bool) {
	if q.Now.IsZero() || len(q.MeterIDs) > 8192 || len(q.Regions) < 1 || len(q.Regions) > 32 || len(q.MeterIDs)*len(q.Regions) > 65536 {
		return nil, nil, false
	}
	ids := append([]string{}, q.MeterIDs...)
	regions := append([]string{}, q.Regions...)
	seen := map[string]bool{}
	total := 0
	for _, id := range ids {
		total += len(id)
		if id == "" || !retailCell(id) || seen[id] {
			return nil, nil, false
		}
		seen[id] = true
	}
	seen = map[string]bool{}
	for _, r := range regions {
		total += len(r)
		if !retailCell(r) || !regionID.MatchString(r) || strings.ToLower(r) != r || !costPhysical(r) || seen[r] {
			return nil, nil, false
		}
		seen[r] = true
	}
	if total > 4<<20 {
		return nil, nil, false
	}
	sort.Strings(ids)
	sort.Strings(regions)
	return ids, regions, true
}

func retailQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }
func retailInitial(clauses []string, metadata bool) string {
	prefix := "currencyCode eq 'USD' and type eq 'Consumption' and "
	if metadata {
		prefix += "isPrimaryMeterRegion eq true and "
	}
	return retailEndpoint + "?" + url.Values{"$filter": {prefix + "(" + strings.Join(clauses, " or ") + ")"}}.Encode()
}

func retailNext(raw, current, filter string, previous uint64) (string, uint64, bool) {
	if len(raw) > 8192 || strings.ContainsAny(raw, "\\#") || strings.ContainsAny(raw, " \t\r\n") {
		return "", 0, false
	}
	ref, e := url.Parse(raw)
	if e != nil || ref.RawPath != "" || ref.Opaque != "" || ref.User != nil || ref.Fragment != "" || ref.ForceQuery {
		return "", 0, false
	}
	base, _ := url.Parse(current)
	u, ok := retailURL(base.ResolveReference(ref).String())
	if !ok {
		return "", 0, false
	}
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil || len(q) != 2 || len(q["$filter"]) != 1 || q.Get("$filter") != filter || len(q["$skip"]) != 1 {
		return "", 0, false
	}
	skip := q.Get("$skip")
	if skip == "" || strings.IndexFunc(skip, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
		return "", 0, false
	}
	n, e := strconv.ParseUint(skip, 10, 64)
	if e != nil || n == 0 || n <= previous {
		return "", 0, false
	}
	q.Set("$skip", strconv.FormatUint(n, 10))
	u.RawQuery = q.Encode()
	return u.String(), n, true
}

type retailCollectionRun struct {
	ctx               context.Context
	c                 *RetailCollector
	r                 *RetailPricingCollection
	q                 RetailPricingRequest
	ids, regions      []string
	wanted, regionSet map[string]bool
	meters            map[string]RetailMeter
	prices            map[string]map[string]RetailPrice
	admitted          int
}

func (s *retailCollectionRun) bind(id string, item RetailPriceItem, phase string) string {
	if s.r.Counts.Bindings >= 65536 {
		return "retail_limit"
	}
	s.r.Counts.Bindings++
	s.r.Bindings = append(s.r.Bindings, RetailBinding{id, item, phase})
	return ""
}

func (s *retailCollectionRun) metadata(items []RetailPriceItem) string {
	for _, item := range items {
		if s.ctx.Err() != nil {
			return "retail_request_failed"
		}
		if code := s.work(); code != "" {
			return code
		}
		if !s.wanted[item.MeterID] || retailDate(item).After(s.q.Now) {
			continue
		}
		if code := s.bind(item.MeterID, item, "metadata"); code != "" {
			return code
		}
		old, present := s.meters[item.MeterID]
		if !present || item.TierMinimumUnits < old.Item.TierMinimumUnits {
			s.meters[item.MeterID] = RetailMeter{item.MeterID, item, false}
			continue
		}
		if item.TierMinimumUnits != old.Item.TierMinimumUnits {
			continue
		}
		if retailIdentity(item) != retailIdentity(old.Item) {
			old.Conflict = true
		}
		if retailLess(item, old.Item) {
			old.Item = item
		}
		s.meters[item.MeterID] = old
	}
	return ""
}

func (s *retailCollectionRun) work() string {
	if s.r.Counts.MatchingWork >= 1048576 {
		return "retail_limit"
	}
	s.r.Counts.MatchingWork++
	return ""
}

func (s *retailCollectionRun) pricing(items []RetailPriceItem, index map[retailTuple][]string) string {
	for _, item := range items {
		if s.ctx.Err() != nil {
			return "retail_request_failed"
		}
		if code := s.work(); code != "" {
			return code
		}
		if retailDate(item).After(s.q.Now) {
			continue
		}
		for _, id := range index[retailKey(item)] {
			if code := s.work(); code != "" {
				return code
			}
			if code := s.bind(id, item, "pricing"); code != "" {
				return code
			}
			meter := s.meters[id].Item
			if item.TierMinimumUnits != meter.TierMinimumUnits {
				continue
			}
			if item.UnitOfMeasure != meter.UnitOfMeasure {
				s.r.UoMMismatches = append(s.r.UoMMismatches, RetailUoMMismatch{id, meter.UnitOfMeasure, item})
				continue
			}
			if item.RetailPrice <= 0 || !s.regionSet[item.ArmRegionName] {
				continue
			}
			if s.prices[id] == nil {
				s.prices[id] = map[string]RetailPrice{}
			}
			old, present := s.prices[id][item.ArmRegionName]
			date, previous := retailDate(item), retailDate(old.Item)
			if !present || date.After(previous) {
				s.prices[id][item.ArmRegionName] = RetailPrice{id, item, false}
				continue
			}
			if !date.Equal(previous) {
				continue
			}
			if item.RetailPrice != old.Item.RetailPrice || retailIdentity(item) != retailIdentity(old.Item) {
				old.Conflict = true
			}
			if retailLess(item, old.Item) {
				old.Item = item
			}
			s.prices[id][item.ArmRegionName] = old
		}
	}
	return ""
}

// chain validates a whole page before admitting any of its rows. All budgets
// belong to the run, not a page, phase or batch. No retry is performed.
func (s *retailCollectionRun) chain(initial string, metadata bool, index map[retailTuple][]string) (string, error) {
	u, _ := url.Parse(initial)
	filter := u.Query().Get("$filter")
	current := initial
	skip := uint64(0)
	seen := map[string]bool{initial: true}
	for page := 0; page < 64; page++ {
		if cause := retailCause(s.ctx, nil); cause != nil {
			return "", cause
		}
		if s.r.Counts.Requests >= 4096 || s.r.Counts.ReceivedBytes >= 32<<20 {
			return "retail_limit", nil
		}
		s.r.Counts.Requests++
		body, resp, e := s.c.getter.GetBounded(s.ctx, current, 1<<20)
		if cause := retailCause(s.ctx, e); cause != nil {
			return "", cause
		}
		// Account bounded error bodies too, including the single overflow byte
		// retained by our transport's limit detector.
		received := int64(len(body))
		if errors.Is(e, errRetailLimit) && received == 0 {
			received = (1 << 20) + 1
		}
		s.r.Counts.ReceivedBytes += received
		if len(body) > 1<<20 || s.r.Counts.ReceivedBytes > 32<<20 || errors.Is(e, errRetailLimit) {
			return "retail_limit", nil
		}

		if e != nil {
			return "retail_request_failed", nil
		}
		if resp == nil {
			return "retail_response_invalid", nil
		}
		if resp.StatusCode != 200 {
			return "retail_request_failed", nil
		}
		items, next, code := decodeRetailPage(s.ctx, body, &s.r.Counts.JSONTokens, &s.r.Counts.RawItems, metadata)
		if cause := retailCause(s.ctx, nil); cause != nil {
			return "", cause
		}
		if code != "" {
			return code, nil
		}
		s.admitted++
		s.r.RawItems = append(s.r.RawItems, items...)
		if metadata {
			code = s.metadata(items)
		} else {
			code = s.pricing(items, index)
		}
		if cause := retailCause(s.ctx, nil); cause != nil {
			return "", cause
		}
		if code != "" {
			return code, nil
		}
		if next == "" {
			return "", nil
		}
		var ok bool
		current, skip, ok = retailNext(next, current, filter, skip)
		if !ok || seen[current] {
			return "retail_continuation_invalid", nil
		}
		seen[current] = true
	}
	return "retail_limit", nil
}

func (s *retailCollectionRun) batches(clauses []string, max int, metadata bool, index map[retailTuple][]string) (string, error) {
	for len(clauses) > 0 {
		n := min(max, len(clauses))
		for n > 0 && len(retailInitial(clauses[:n], metadata)) > 8192 {
			n--
		}
		if n == 0 {
			return "retail_limit", nil
		}
		code, e := s.chain(retailInitial(clauses[:n], metadata), metadata, index)
		if e != nil || code != "" {
			return code, e
		}
		clauses = clauses[n:]
	}
	return "", nil
}

func (s *retailCollectionRun) finish(code string) (*RetailPricingCollection, error) {
	if cause := retailCause(s.ctx, nil); cause != nil {
		return nil, cause
	}
	r := s.r
	r.Status = "complete"
	if code != "" {
		r.Status = "unknown"
		if s.admitted > 0 {
			r.Status = "partial"
		}
	}
	for _, id := range s.ids {
		meter, ok := s.meters[id]
		if ok {
			r.Meters = append(r.Meters, meter)
		}
		if code == "" && (!ok || meter.Conflict) {
			r.Status = "partial"
		}
		for _, region := range s.regions {
			price, ok := s.prices[id][region]
			if ok {
				r.Prices = append(r.Prices, price)
			}
			if code == "" && (!ok || price.Conflict) {
				r.Status = "partial"
			}
		}
	}
	if code == "" && r.Status != "complete" {
		code = "retail_incomplete"
	}
	r.FailureCode = code
	r.admission = r.Status
	if r.Status == "complete" {
		r.eligible = map[string]map[string]float64{}
		for _, p := range r.Prices {
			if r.eligible[p.OriginalMeterID] == nil {
				r.eligible[p.OriginalMeterID] = map[string]float64{}
			}
			r.eligible[p.OriginalMeterID][p.Item.ArmRegionName] = p.Item.RetailPrice
		}
	}
	// Stable diagnostic order includes ignored/duplicate raw rows and bindings.
	sort.Slice(r.RawItems, func(i, j int) bool { return retailLess(r.RawItems[i], r.RawItems[j]) })
	sort.Slice(r.Bindings, func(i, j int) bool {
		a, b := r.Bindings[i], r.Bindings[j]
		if a.OriginalMeterID != b.OriginalMeterID {
			return a.OriginalMeterID < b.OriginalMeterID
		}
		if a.Phase != b.Phase {
			return a.Phase < b.Phase
		}
		return retailLess(a.Item, b.Item)
	})
	sort.Slice(r.UoMMismatches, func(i, j int) bool {
		a, b := r.UoMMismatches[i], r.UoMMismatches[j]
		if a.OriginalMeterID != b.OriginalMeterID {
			return a.OriginalMeterID < b.OriginalMeterID
		}
		return retailLess(a.Item, b.Item)
	})
	if cause := retailCause(s.ctx, nil); cause != nil {
		return nil, cause
	}
	return r, nil
}

func (c *RetailCollector) Collect(ctx context.Context, q RetailPricingRequest) (*RetailPricingCollection, error) {
	if ctx == nil || c == nil || c.getter == nil {
		return nil, errRetailRequest
	}
	if cause := retailCause(ctx, nil); cause != nil {
		return nil, cause
	}
	ids, regions, ok := retailRequest(q)
	if !ok {
		return nil, errRetailRequest
	}
	s := &retailCollectionRun{ctx: ctx, c: c, q: q, ids: ids, regions: regions, r: &RetailPricingCollection{}, wanted: map[string]bool{}, regionSet: map[string]bool{}, meters: map[string]RetailMeter{}, prices: map[string]map[string]RetailPrice{}}
	clauses := make([]string, 0, len(ids))
	for _, id := range ids {
		s.wanted[id] = true
		clauses = append(clauses, "meterId eq "+retailQuote(id))
	}
	for _, r := range regions {
		s.regionSet[r] = true
	}
	code, e := s.batches(clauses, 15, true, nil)
	if e != nil {
		return nil, e
	}
	if code != "" {
		return s.finish(code)
	}
	index := map[retailTuple][]string{}
	for _, id := range ids {
		if m, ok := s.meters[id]; ok && !m.Conflict {
			index[retailKey(m.Item)] = append(index[retailKey(m.Item)], id)
		}
	}
	clauses = nil
	for k := range index {
		clauses = append(clauses, "(meterName eq "+retailQuote(k.MeterName)+" and productId eq "+retailQuote(k.ProductID)+" and skuName eq "+retailQuote(k.SkuName)+")")
	}
	sort.Strings(clauses)
	code, e = s.batches(clauses, 5, false, index)
	if e != nil {
		return nil, e
	}
	return s.finish(code)
}
