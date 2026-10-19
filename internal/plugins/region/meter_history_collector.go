// Region meter-history semantics derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_METER_HISTORY_CONTRACT_20261019.md.
package region

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

// MeterHistoryPoster must enforce every-attempt byte/deadline bounds, close
// response bodies, use the selected ARM audience and never follow redirects.
// The existing azure.HTTPClient.PostBounded implements this contract.
type MeterHistoryPoster interface {
	PostBounded(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)
}

type MeterHistoryCollection struct {
	SubscriptionID string
	Currency       string
	Status         string
	FailureCode    string
	Meters         []HistoricalCostMeter
}

// Evidence never promotes a diagnostic page prefix to usable cost weights.
func (c *MeterHistoryCollection) Evidence() CostHistoryEvidence {
	e := CostHistoryEvidence{Meters: []HistoricalCostMeter{}}
	if c != nil && c.Status == "complete" {
		e.Complete = true
		e.Meters = append(e.Meters, c.Meters...)
	}
	return e
}

type MeterHistoryCollector struct {
	origin *url.URL
	poster MeterHistoryPoster
}

func NewMeterHistoryCollector(endpoint string, poster MeterHistoryPoster) (*MeterHistoryCollector, error) {
	u, e := url.Parse(endpoint)
	if e != nil || len(endpoint) > 2048 || u.Scheme != "https" || !validQuotaHost(u.Hostname()) || (u.Port() != "" && u.Port() != "443") || strings.HasSuffix(u.Host, ":") || u.User != nil || u.Opaque != "" || u.ForceQuery || u.RawQuery != "" || u.Fragment != "" || u.RawPath != "" || (u.Path != "" && u.Path != "/") || strings.ContainsAny(endpoint, "%#\\ \r\n\t") || poster == nil {
		return nil, fmt.Errorf("region meter history rejected [history_request_invalid]")
	}
	u.Host = strings.ToLower(u.Hostname())
	u.Path = ""
	return &MeterHistoryCollector{origin: u, poster: poster}, nil
}

func historyTerminal(ctx context.Context, cause error) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if errors.Is(cause, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(cause, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	return nil
}

func historyMeters(totals map[string]float64) []HistoricalCostMeter {
	out := make([]HistoricalCostMeter, 0, len(totals))
	for id, value := range totals {
		out = append(out, HistoricalCostMeter{MeterID: id, HistoricalCost: value})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].MeterID < out[j].MeterID })
	return out
}

func (c *MeterHistoryCollector) continuation(first *url.URL, link string) (*url.URL, bool) {
	if len(link) > 8192 || strings.ContainsAny(link, "#\\ \r\n\t") {
		return nil, false
	}
	ref, e := url.Parse(link)
	if e != nil {
		return nil, false
	}
	u := first.ResolveReference(ref)
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil || u.Scheme != "https" || !validQuotaHost(u.Hostname()) || !strings.EqualFold(u.Hostname(), c.origin.Hostname()) || (u.Port() != "" && u.Port() != "443") || strings.HasSuffix(u.Host, ":") || u.Path != first.Path || u.RawPath != "" || u.Opaque != "" || u.User != nil || u.Fragment != "" || u.ForceQuery || len(q["api-version"]) != 1 || q.Get("api-version") != "2021-10-01" {
		return nil, false
	}
	for key, values := range q {
		if key != "api-version" && key != "$skiptoken" || len(values) != 1 {
			return nil, false
		}
		if key == "$skiptoken" && (values[0] == "" || len(values[0]) > 4096 || !utf8.ValidString(values[0]) || strings.IndexFunc(values[0], func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || regionNoncharacter(r) }) >= 0) {
			return nil, false
		}
	}
	u.Host = c.origin.Host
	u.RawQuery = "api-version=2021-10-01"
	if token := q.Get("$skiptoken"); token != "" {
		u.RawQuery += "&$skiptoken=" + url.QueryEscape(token)
	}
	return u, true
}

// Collect admits entire pages atomically. The caller keeps inputs stable during
// the call; all run state/results are owned and no credential/client is created.
func (c *MeterHistoryCollector) Collect(ctx context.Context, selected map[string]string, q MeterHistoryRequest) (*MeterHistoryCollection, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	if c == nil || c.origin == nil || c.poster == nil {
		return nil, historyRequestError(ctx)
	}
	text := 0
	scope, valid := auxiliaryScope(ctx, selected, &text)
	if !valid {
		return nil, historyRequestError(ctx)
	}
	for _, name := range scope {
		if name == "" {
			return nil, historyRequestError(ctx)
		}
	}
	id := strings.ToLower(q.SubscriptionID)
	if !historyRequestValid(q) || scope[id] == "" {
		return nil, historyRequestError(ctx)
	}
	body, e := MeterHistoryQuery(ctx, q)
	if e != nil {
		return nil, e
	}
	first := *c.origin
	first.Path = "/subscriptions/" + id + "/providers/Microsoft.CostManagement/query"
	first.RawQuery = "api-version=2021-10-01"
	result := &MeterHistoryCollection{SubscriptionID: id, Currency: q.Currency, Status: "unknown", Meters: []HistoricalCostMeter{}}
	totals := map[string]float64{}
	pages, received, work, rowCount := 0, 0, 0, 0
	fail := func(code string) (*MeterHistoryCollection, error) {
		result.FailureCode = code
		result.Meters = historyMeters(totals)
		if pages > 0 {
			result.Status = "partial"
		}
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		return result, nil
	}
	visited := map[string]bool{}
	link := first.String()
	for {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if pages >= 64 {
			return fail("history_limit")
		}
		u, ok := c.continuation(&first, link)
		if !ok || visited[u.String()] {
			return fail("history_continuation_invalid")
		}
		visited[u.String()] = true
		b, response, err := c.poster.PostBounded(ctx, u.String(), azure.NopReadSeekCloser{Reader: bytes.NewReader(body)}, 1<<20)
		if terminal := historyTerminal(ctx, err); terminal != nil {
			return nil, terminal
		}
		if len(b) > 1<<20 || len(b) > 8<<20-received {
			return fail("history_limit")
		}
		received += len(b)
		if response != nil && response.StatusCode != 200 && response.StatusCode != 204 {
			return fail("history_response_invalid")
		}
		if err != nil {
			return fail("history_request_failed")
		}
		if response == nil || response.StatusCode == 204 && len(b) != 0 {
			return fail("history_response_invalid")
		}
		if response.StatusCode == 204 {
			pages++
			break
		}
		rows, next, code := decodeHistoryPage(ctx, b, id, q.Currency, &work, &rowCount)
		if code != "" {
			return fail(code)
		}
		pending := map[string]float64{}
		for _, row := range rows {
			if e := ctx.Err(); e != nil {
				return nil, e
			}
			v, exists := pending[row.MeterID]
			if !exists {
				v = totals[row.MeterID]
			}
			v += row.HistoricalCost
			if v > MaxCostValue {
				return fail("history_decode_invalid")
			}
			pending[row.MeterID] = v
		}
		for key, value := range pending {
			totals[key] = value
		}
		pages++
		if next == "" {
			break
		}
		link = next
	}
	result.Status = "complete"
	result.Meters = historyMeters(totals)
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return result, nil
}
