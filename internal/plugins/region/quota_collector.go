// REST request contracts derived from MIT-licensed AZQR. See NOTICE.md and
// docs/REGION_QUOTA_COLLECTOR.md for deliberate corrections and admission bounds.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MaxQuotaPageBytes  = 1 << 20
	MaxQuotaTotalBytes = 8 << 20
	MaxQuotaPages      = 64
)

// RESTQuotaGetter must honor cancellation, bound all attempt bodies, use the
// selected ARM audience, close the response and never follow redirects.
type RESTQuotaGetter interface {
	GetBoundedWithResponse(context.Context, string, int64) ([]byte, *http.Response, error)
}

type RESTQuotaCollection struct {
	Evidence    QuotaEvidence
	FailureCode string
}

type RESTQuotaCollector struct {
	origin *url.URL
	getter RESTQuotaGetter
}

var quotaHost = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9.-]*[a-zA-Z0-9])?$`)

func validQuotaHost(host string) bool {
	if len(host) > 253 || !quotaHost.MatchString(host) {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
	}
	return true
}

func NewRESTQuotaCollector(endpoint string, getter RESTQuotaGetter) (*RESTQuotaCollector, error) {
	u, err := url.Parse(endpoint)
	if err != nil || len(endpoint) > 2048 || u.Scheme != "https" || !validQuotaHost(u.Hostname()) || (u.Port() != "" && u.Port() != "443") || u.User != nil || u.Opaque != "" || u.ForceQuery || u.RawQuery != "" || u.Fragment != "" || u.RawPath != "" || (u.Path != "" && u.Path != "/") || strings.ContainsAny(endpoint, "%#\\ \r\n\t") || getter == nil {
		return nil, fmt.Errorf("invalid REST quota origin or getter")
	}
	u.Path = ""
	return &RESTQuotaCollector{origin: u, getter: getter}, nil
}

func restQuotaProvider(kind string) (string, string) {
	switch kind {
	case "Network":
		return "Microsoft.Network", "2022-07-01"
	case "SQL":
		return "Microsoft.Sql", "2021-11-01"
	case "App Service":
		return "Microsoft.Web", "2023-01-01"
	case "Storage":
		return "Microsoft.Storage", "2023-01-01"
	}
	return "", ""
}

// Collect preserves only fully admitted pages. Operational failure is explicit
// incomplete evidence; invalid request and cancellation are errors with no result.
// Inputs must remain stable for the call. The collector owns no mutable run state.
func (c *RESTQuotaCollector) Collect(ctx context.Context, subscriptions map[string]string, request QuotaRequest) (*RESTQuotaCollection, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	text := 0
	scope, valid := auxiliaryScope(ctx, subscriptions, &text)
	if !valid {
		return nil, quotaFailure(ctx, "scope_invalid")
	}
	for _, name := range scope {
		if name == "" {
			return nil, quotaFailure(ctx, "scope_invalid")
		}
	}
	request.SubscriptionID = strings.ToLower(request.SubscriptionID)
	provider, version := restQuotaProvider(request.QuotaType)
	if len(request.SubscriptionID) != 36 || !subscriptionID.MatchString(request.SubscriptionID) || scope[request.SubscriptionID] == "" || !regionID.MatchString(request.Region) || provider == "" {
		return nil, quotaFailure(ctx, "request_invalid")
	}
	first := *c.origin
	first.Path = "/subscriptions/" + request.SubscriptionID + "/providers/" + provider + "/locations/" + request.Region + "/usages"
	first.RawQuery = "api-version=" + version
	result := &RESTQuotaCollection{Evidence: QuotaEvidence{Request: request, Status: "complete", Usages: []QuotaUsage{}}}
	pages, total := 0, 0
	seenLinks := map[string]bool{}
	seenNames := map[string]bool{}
	link := first.String()
	fail := func(code string, status int) (*RESTQuotaCollection, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result.FailureCode = code
		result.Evidence.Status = "unknown"
		if pages > 0 {
			result.Evidence.Status = "partial"
		} else if status == 404 || status == 405 {
			result.Evidence.Status = "unsupported"
		}
		return result, nil
	}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if pages >= MaxQuotaPages {
			return fail("quota_page_limit", 0)
		}
		u, ok := c.continuation(&first, version, link)
		if !ok {
			return fail("quota_unsafe_continuation", 0)
		}
		if seenLinks[u.String()] {
			return fail("quota_continuation_cycle", 0)
		}
		seenLinks[u.String()] = true
		body, response, err := c.getter.GetBoundedWithResponse(ctx, u.String(), MaxQuotaPageBytes)
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if err != nil {
			status := 0
			if response != nil {
				status = response.StatusCode
			}
			return fail("quota_request_failed", status)
		}
		if response == nil || response.StatusCode != http.StatusOK {
			status := 0
			if response != nil {
				status = response.StatusCode
			}
			return fail("quota_status_invalid", status)
		}
		if len(body) > MaxQuotaPageBytes || len(body) > MaxQuotaTotalBytes-total {
			return fail("quota_byte_limit", 0)
		}
		total += len(body)
		rows, next, code := decodeQuotaPage(ctx, body)
		if code != "" {
			return fail(code, 0)
		}
		if len(rows) > MaxAuxRows-len(result.Evidence.Usages) {
			return fail("quota_row_limit", 0)
		}
		pageNames := map[string]bool{}
		for _, row := range rows {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if row.ResourceName != "" && (seenNames[row.ResourceName] || pageNames[row.ResourceName]) {
				return fail("quota_identity_duplicate", 0)
			}
			pageNames[row.ResourceName] = true
			if !auxiliaryText(row.ResourceName, &text) || !auxiliaryText(row.LocalizedName, &text) {
				return fail("quota_text_limit", 0)
			}
		}
		for name := range pageNames {
			seenNames[name] = true
		}
		result.Evidence.Usages = append(result.Evidence.Usages, rows...)
		pages++
		if next == "" {
			break
		}
		link = next
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RESTQuotaCollector) continuation(first *url.URL, version, link string) (*url.URL, bool) {
	if len(link) > 8192 || strings.ContainsAny(link, "#\\ \r\n\t") {
		return nil, false
	}
	ref, err := url.Parse(link)
	if err != nil {
		return nil, false
	}
	u := first.ResolveReference(ref)
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil || u.Scheme != "https" || !strings.EqualFold(u.Host, c.origin.Host) || u.Path != first.Path || u.RawPath != "" || u.Opaque != "" || u.User != nil || u.Fragment != "" || u.ForceQuery || len(q["api-version"]) != 1 || q.Get("api-version") != version {
		return nil, false
	}
	for key, values := range q {
		if key != "api-version" && key != "$skiptoken" || len(values) != 1 || (key == "$skiptoken" && (values[0] == "" || len(values[0]) > 4096)) {
			return nil, false
		}
	}
	u.RawQuery = q.Encode()
	return u, true
}

func quotaJSON(ctx context.Context, body []byte) bool {
	if !utf8.Valid(body) {
		return false
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	tokens := 0
	var walk func(int, bool) bool
	walk = func(depth int, large bool) bool {
		if ctx.Err() != nil || depth > 32 || tokens >= 65536 {
			return false
		}
		tokens++
		t, err := d.Token()
		if err != nil {
			return false
		}
		if s, ok := t.(string); ok {
			limit := MaxLabelBytes
			if large {
				limit = 8192
			}
			return len(s) <= limit && utf8.ValidString(s) && strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || r == 0xfffe || r == 0xffff }) < 0
		}
		delim, compound := t.(json.Delim)
		if !compound {
			return true
		}
		if delim != '{' && delim != '[' {
			return false
		}
		keys := map[string]bool{}
		for d.More() {
			longString := false
			if delim == '{' {
				k, err := d.Token()
				s, ok := k.(string)
				n := 0
				if err != nil || !ok || !auxiliaryText(s, &n) || len(keys) >= 128 || keys[strings.ToLower(s)] {
					return false
				}
				keys[strings.ToLower(s)] = true
				longString = depth == 0 && strings.EqualFold(s, "nextLink")
			}
			if !walk(depth+1, longString) {
				return false
			}
		}
		end, err := d.Token()
		return err == nil && (delim == '{' && end == json.Delim('}') || delim == '[' && end == json.Delim(']'))
	}
	if !walk(0, false) {
		return false
	}
	_, err := d.Token()
	return err == io.EOF
}

func decodeQuotaPage(ctx context.Context, body []byte) ([]QuotaUsage, string, string) {
	if !quotaJSON(ctx, body) {
		return nil, "", "quota_invalid_response"
	}
	var envelope struct {
		Value    json.RawMessage `json:"value"`
		NextLink *string         `json:"nextLink"`
	}
	if json.Unmarshal(body, &envelope) != nil || len(envelope.Value) == 0 || bytes.Equal(bytes.TrimSpace(envelope.Value), []byte("null")) {
		return nil, "", "quota_invalid_response"
	}
	var items []json.RawMessage
	if json.Unmarshal(envelope.Value, &items) != nil {
		return nil, "", "quota_invalid_response"
	}
	if len(items) > MaxAuxRows {
		return nil, "", "quota_row_limit"
	}
	rows := make([]QuotaUsage, 0, len(items))
	for _, raw := range items {
		if ctx.Err() != nil {
			return nil, "", "quota_cancelled"
		}
		var item struct {
			Name    json.RawMessage `json:"name"`
			Current *int64          `json:"currentValue"`
			Limit   *int64          `json:"limit"`
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &item) != nil {
			return nil, "", "quota_invalid_response"
		}
		row := QuotaUsage{}
		if len(item.Name) != 0 && !bytes.Equal(bytes.TrimSpace(item.Name), []byte("null")) {
			if json.Unmarshal(item.Name, &row.ResourceName) == nil {
				row.LocalizedName = row.ResourceName
			} else {
				var name struct {
					Value     string `json:"value"`
					Localized string `json:"localizedValue"`
				}
				if json.Unmarshal(item.Name, &name) != nil {
					return nil, "", "quota_invalid_response"
				}
				row.ResourceName = name.Value
				row.LocalizedName = name.Localized
			}
		}
		if item.Current != nil {
			row.Current = *item.Current
			row.CurrentKnown = true
		}
		if item.Limit != nil {
			row.Limit = *item.Limit
			row.LimitKnown = true
		}
		if row.Current < 0 || row.Current > MaxAuxCount || row.Limit < -MaxAuxCount || row.Limit > MaxAuxCount {
			return nil, "", "quota_count_invalid"
		}
		rows = append(rows, row)
	}
	next := ""
	if envelope.NextLink != nil {
		next = *envelope.NextLink
	}
	return rows, next, ""
}
