// Carbon request behavior is characterized from Microsoft Azure Quick Review
// (MIT licensed) and the 2025-04-01 ARM contract. See docs/CARBON_EMISSIONS.md.
package carbon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const (
	MaxSubscriptions = 3000
	MaxPages         = 64 // Across all batches, not per subscription.
	MaxPageBytes     = 2 << 20
	MaxResponseBytes = 16 << 20 // Successful date and report bodies together.
	MaxPageItems     = 1000
	MaxItemBytes     = 8 << 10
	MaxTokenBytes    = 4096
	MaxDuration      = 5 * time.Minute
)

type BoundedPoster interface {
	PostBounded(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)
}

// Scanner captures immutable destinations. Injected clients/filter callbacks are
// trusted application code, must honor context and must not follow redirects.
type Scanner struct {
	client                       BoundedPoster
	dateEndpoint, reportEndpoint string
}

func NewWithHTTPClient(endpoint string, client BoundedPoster) (*Scanner, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.RawPath != "" || u.Opaque != "" || strings.Contains(endpoint, "#") || (u.Path != "" && u.Path != "/") || client == nil {
		return nil, fmt.Errorf("invalid carbon ARM endpoint or client")
	}
	u.Path = "/providers/Microsoft.Carbon/queryCarbonEmissionDataAvailableDateRange"
	u.RawQuery = "api-version=2025-04-01"
	dateEndpoint := u.String()
	u.Path = "/providers/Microsoft.Carbon/carbonEmissionReports"
	return &Scanner{client: client, dateEndpoint: dateEndpoint, reportEndpoint: u.String()}, nil
}

func New(credential azcore.TokenCredential) (*Scanner, error) {
	return NewWithHTTPClient(azure.ResourceManagerEndpoint(), azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(30*time.Second)))
}

type reportRequest struct {
	ReportType       string   `json:"reportType"`
	SubscriptionList []string `json:"subscriptionList"`
	CarbonScopeList  []string `json:"carbonScopeList"`
	DateRange        struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"dateRange"`
	CategoryType  string `json:"categoryType"`
	OrderBy       string `json:"orderBy"`
	SortDirection string `json:"sortDirection"`
	PageSize      int    `json:"pageSize"`
	SkipToken     string `json:"skipToken,omitempty"`
}

var scopeID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Scan retains valid pages and later batches on service failures. Aggregated
// labels cannot prove subscription membership; access decisions are separate.
// Missing/partial decisions warn; explicit denial fails, even with valid rows.
func (s *Scanner) Scan(ctx context.Context, subscriptions map[string]string, filter Filter) (assessment.PluginTable, error) {
	ctx, cancel := context.WithTimeout(ctx, MaxDuration)
	defer cancel()
	started := time.Now().UTC()
	items := []Item{}
	period := ""
	malformed, unverified, denied := 0, 0, map[string]bool{}
	failureCode := ""
	var failureCause error
	mark := func(code string, cause error) {
		if failureCode == "" {
			failureCode, failureCause = code, cause
		}
		if cause != nil {
			failureCode, failureCause = code, cause
		}
	}
	finish := func() (assessment.PluginTable, error) {
		table := PendingTable()
		table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted}
		var projectionErr error
		if period != "" {
			// No more I/O: project bounded owned successful input even after caller
			// cancellation, then retain cancellation identity in final health/error.
			table, projectionErr = Project(context.WithoutCancel(ctx), period, period, items, filter)
		}
		if projectionErr != nil {
			mark(table.Health.Error.Code, nil)
		}
		if malformed > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "carbon_malformed_items", Message: fmt.Sprintf("skipped %d invalid carbon emission items", malformed)})
		}
		if unverified > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "carbon_access_unverified", Message: fmt.Sprintf("%d carbon response pages have incomplete access decisions", unverified)})
		}
		if len(denied) > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "carbon_access_denied", Message: fmt.Sprintf("carbon access denied for %d selected subscriptions", len(denied))})
			mark("carbon_access_denied", nil)
		}
		if err := ctx.Err(); err != nil {
			mark("carbon_cancelled", err)
		}
		table.Health.StartedAt, table.Health.FinishedAt = started, time.Now().UTC()
		table.Health.Records = len(table.Rows)
		if failureCode != "" {
			table.Health.Status = assessment.StageFailed
			table.Health.Error = &assessment.AssessmentError{Code: failureCode, Message: "carbon emissions returned incomplete data"}
		} else if len(table.Health.Warnings) > 0 {
			table.Health.Status = assessment.StageCompletedWithWarnings
		}
		if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
			table.Rows = []assessment.PluginRow{}
			table.Health.Records = 0
			failureCode = "carbon_output_invalid"
		}
		if failureCode != "" {
			table.Health.Status = assessment.StageFailed
			table.Health.Error = &assessment.AssessmentError{Code: failureCode, Message: "carbon emissions returned incomplete data"}
			if failureCause != nil {
				return table, fmt.Errorf("carbon emissions stopped: %w", failureCause)
			}
			return table, fmt.Errorf("carbon emissions returned incomplete data")
		}
		if len(table.Health.Warnings) > 0 {
			table.Health.Status = assessment.StageCompletedWithWarnings
		}
		return table, nil
	}
	if err := ctx.Err(); err != nil {
		mark("carbon_cancelled", err)
		return finish()
	}
	if len(subscriptions) > MaxSubscriptions {
		mark("carbon_scope_limit", nil)
		return finish()
	}
	ids, selected := []string{}, map[string]bool{}
	for id := range subscriptions {
		normalized := strings.ToLower(id)
		if !scopeID.MatchString(id) || selected[normalized] {
			mark("carbon_scope_invalid", nil)
			return finish()
		}
		selected[normalized] = true
		ids = append(ids, normalized)
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return finish()
	} // No credentials, requests or fabricated period.
	if s == nil || s.client == nil {
		mark("carbon_not_configured", nil)
		return finish()
	}
	responseBytes, pages := 0, 0
	budgetExhausted := false
	post := func(endpoint string, body io.ReadSeekCloser) ([]byte, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		data, response, err := s.client.PostBounded(ctx, endpoint, body, MaxPageBytes)
		if err != nil || response == nil || response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("carbon request failed")
		}
		if len(data) > MaxPageBytes || responseBytes+len(data) > MaxResponseBytes {
			budgetExhausted = true
			return nil, fmt.Errorf("carbon response byte limit")
		}
		responseBytes += len(data)
		return data, nil
	}
	data, err := post(s.dateEndpoint, nil)
	if err != nil {
		mark("carbon_date_request_failed", nil)
		return finish()
	}
	dateFields, err := object(data, []string{"startDate", "endDate"})
	if err != nil {
		mark("carbon_date_invalid", nil)
		return finish()
	}
	start, a := stringField(dateFields, "startDate", 10, true)
	end, b := stringField(dateFields, "endDate", 10, true)
	if _, err = LatestPeriod(start, end); a != nil || b != nil || err != nil {
		mark("carbon_date_invalid", nil)
		return finish()
	}
	period = end
	for offset := 0; offset < len(ids); offset += 100 {
		batch := append([]string(nil), ids[offset:min(offset+100, len(ids))]...)
		batchScope := map[string]bool{}
		for _, id := range batch {
			batchScope[id] = true
		}
		decisions, seenTokens := map[string]string{}, map[string]bool{}
		token := ""
		for {
			if err := ctx.Err(); err != nil {
				mark("carbon_cancelled", err)
				return finish()
			}
			if pages >= MaxPages {
				mark("carbon_page_limit", nil)
				return finish()
			}
			request := reportRequest{ReportType: "ItemDetailsReport", SubscriptionList: batch, CarbonScopeList: []string{"Scope1", "Scope2", "Scope3"}, CategoryType: "ResourceType", OrderBy: "LatestMonthEmissions", SortDirection: "Desc", PageSize: MaxPageItems, SkipToken: token}
			request.DateRange.Start, request.DateRange.End = period, period
			body, _ := json.Marshal(request)
			pages++
			data, err := post(s.reportEndpoint, azure.NopReadSeekCloser{Reader: bytes.NewReader(body)})
			if err != nil {
				mark("carbon_report_request_failed", nil)
				// Total successful-response budget exhaustion stops all subsequent I/O.
				if budgetExhausted {
					failureCode = "carbon_response_limit"
					return finish()
				}
				break
			}
			page, err := decodePage(data, batchScope)
			if err != nil {
				mark("carbon_response_invalid", nil)
				break
			}
			if page.token != "" && seenTokens[page.token] {
				mark("carbon_token_cycle", nil)
				break
			}
			conflict := false
			for id, decision := range page.decisions {
				if old := decisions[id]; old != "" && old != decision {
					conflict = true
				}
			}
			if conflict {
				mark("carbon_access_invalid", nil)
				break
			}
			for id, decision := range page.decisions {
				decisions[id] = decision
				if decision == "Denied" {
					denied[id] = true
				}
			}
			if len(page.decisions) != len(batch) {
				unverified++
			}
			if len(items)+len(page.items) > MaxItems {
				mark("carbon_item_limit", nil)
				return finish()
			}
			malformed += page.malformed
			items = append(items, page.items...)
			if page.token == "" {
				break
			}
			seenTokens[page.token] = true
			token = page.token
		}
	}
	return finish()
}

type responsePage struct {
	items     []Item
	malformed int
	decisions map[string]string
	token     string
}

// object rejects duplicate/case-aliased keys before projecting known fields.
// Unique unknown bounded metadata is permitted for forward compatibility.
func object(raw []byte, known []string) (map[string]json.RawMessage, error) {
	if !utf8.Valid(raw) || !json.Valid(raw) {
		return nil, fmt.Errorf("invalid carbon JSON")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	t, err := d.Token()
	if err != nil || t != json.Delim('{') {
		return nil, fmt.Errorf("invalid carbon object")
	}
	values, seen := map[string]json.RawMessage{}, map[string]bool{}
	for d.More() {
		t, err := d.Token()
		key, ok := t.(string)
		if err != nil || !ok || len(key) > 128 || len(seen) >= 32 || seen[strings.ToLower(key)] || strings.EqualFold(key, "error") {
			return nil, fmt.Errorf("ambiguous carbon object")
		}
		seen[strings.ToLower(key)] = true
		for _, field := range known {
			if strings.EqualFold(key, field) && key != field {
				return nil, fmt.Errorf("aliased carbon field")
			}
		}
		var value json.RawMessage
		if d.Decode(&value) != nil {
			return nil, fmt.Errorf("invalid carbon field")
		}
		values[key] = value
	}
	return values, nil
}

func safeString(value string, limit int) bool {
	return len(value) <= limit && utf8.ValidString(value) && strings.IndexFunc(value, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffe || r == 0xffff || r == 0xfffd }) < 0
}
func stringField(fields map[string]json.RawMessage, key string, limit int, required bool) (string, error) {
	raw := fields[key]
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		if required {
			return "", fmt.Errorf("missing carbon string")
		}
		return "", nil
	}
	var value string
	if json.Unmarshal(raw, &value) != nil || !safeString(value, limit) || required && value == "" {
		return "", fmt.Errorf("invalid carbon string")
	}
	return value, nil
}

func array(raw []byte, limit int, nullable bool) ([]json.RawMessage, error) {
	if nullable && (len(raw) == 0 || bytes.Equal(raw, []byte("null"))) {
		return nil, nil
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	t, err := d.Token()
	if err != nil || t != json.Delim('[') {
		return nil, fmt.Errorf("missing carbon array")
	}
	values := []json.RawMessage{}
	for d.More() {
		if len(values) >= limit {
			return nil, fmt.Errorf("carbon array limit")
		}
		var value json.RawMessage
		if d.Decode(&value) != nil {
			return nil, fmt.Errorf("invalid carbon array")
		}
		values = append(values, value)
	}
	if _, err = d.Token(); err != nil {
		return nil, fmt.Errorf("invalid carbon array")
	}
	return values, nil
}

func decodePage(raw []byte, batch map[string]bool) (responsePage, error) {
	page := responsePage{decisions: map[string]string{}}
	fields, err := object(raw, []string{"value", "skipToken", "subscriptionAccessDecisionList"})
	if err != nil {
		return page, err
	}
	page.token, err = stringField(fields, "skipToken", MaxTokenBytes, false)
	if err != nil {
		return page, err
	}
	decisions, err := array(fields["subscriptionAccessDecisionList"], 100, true)
	if err != nil {
		return page, err
	}
	for _, raw := range decisions {
		if len(raw) > MaxItemBytes {
			return page, fmt.Errorf("carbon decision limit")
		}
		entry, err := object(raw, []string{"subscriptionId", "decision", "denialReason"})
		if err != nil {
			return page, err
		}
		id, a := stringField(entry, "subscriptionId", 36, true)
		decision, b := stringField(entry, "decision", 7, true)
		_, c := stringField(entry, "denialReason", 4096, false)
		id = strings.ToLower(id)
		if a != nil || b != nil || c != nil || !scopeID.MatchString(id) || !batch[id] || page.decisions[id] != "" || (decision != "Allowed" && decision != "Denied") {
			return page, fmt.Errorf("invalid carbon access decisions")
		}
		page.decisions[id] = decision
	}
	rows, err := array(fields["value"], MaxPageItems, false)
	if err != nil {
		return page, err
	}
	for _, row := range rows {
		item, err := decodeItem(row)
		if err != nil {
			page.malformed++
			continue
		}
		page.items = append(page.items, item)
	}
	return page, nil
}

func decodeItem(raw []byte) (Item, error) {
	var item Item
	if len(raw) > MaxItemBytes {
		return item, fmt.Errorf("carbon item limit")
	}
	fields, err := object(raw, []string{"dataType", "categoryType", "itemName", "latestMonthEmissions", "previousMonthEmissions", "monthlyEmissionsChangeValue", "monthOverMonthEmissionsChangeRatio"})
	if err != nil {
		return item, err
	}
	typeName, a := stringField(fields, "dataType", 64, true)
	category, b := stringField(fields, "categoryType", 64, true)
	label, c := stringField(fields, "itemName", MaxLabelBytes, true)
	if a != nil || b != nil || c != nil || typeName != "ItemDetailsData" || category != "ResourceType" || !validLabel(label) {
		return item, fmt.Errorf("invalid carbon item category")
	}
	item.ResourceType = label
	for key, destination := range map[string]**float64{"latestMonthEmissions": &item.Latest, "previousMonthEmissions": &item.Previous, "monthlyEmissionsChangeValue": &item.Change} {
		raw := fields[key]
		if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
			if key == "latestMonthEmissions" {
				return item, fmt.Errorf("missing carbon latest value")
			}
			continue
		}
		var value float64
		if json.Unmarshal(raw, &value) != nil || !finite(value) {
			return item, fmt.Errorf("invalid carbon number")
		}
		*destination = &value
	}
	return item, nil
}
