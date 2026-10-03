// Discovery reproduces the pinned AZQR query with bounded, explicit coverage.
// See docs/AI_GOVERNANCE_EXECUTION.md and NOTICE.md.
package aigov

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const (
	DiscoveryQuery = `resources
| where type =~ "Microsoft.CognitiveServices/accounts"
| where isempty(kind) or kind contains "openai" or kind contains "aiservices"
| project id, subscriptionId, resourceGroup, location, type, name, sku.name, sku.tier, kind
| order by subscriptionId, resourceGroup`
	MaxDiscoveryRequests = 64
	MaxDiscoveryPages    = 32
	MaxDiscoveryBatch    = 300
	MaxDiscoveryPageRows = 1000
)

type Discovery struct {
	origin *url.URL
	client BoundedPoster
}

// DiscoveryResult is independently owned. Failed discovery can retain accounts;
// successful enrichment must not replace this coverage health with success.
type DiscoveryResult struct {
	Accounts []LocatedAccount
	Health   assessment.StageExecution
}

// An injected client/filter is trusted cooperative application code. The client
// must enforce byte/context limits and must never follow redirects.
func NewDiscoveryWithClient(endpoint string, client BoundedPoster) (*Discovery, error) {
	origin, err := originURL(endpoint)
	if err != nil || client == nil {
		return nil, fmt.Errorf("invalid AI discovery endpoint or client")
	}
	return &Discovery{origin: origin, client: client}, nil
}

func NewDiscovery(credential azcore.TokenCredential) (*Discovery, error) {
	return newDiscoveryWithTransport(credential, nil)
}

func newDiscoveryWithTransport(credential azcore.TokenCredential, transport policy.Transporter) (*Discovery, error) {
	arm, scope, err := aiARMConfiguration(credential)
	if err != nil {
		return nil, err
	}
	options := azure.DefaultHTTPClientOptions(30 * time.Second)
	options.Scope, options.Transport = scope, transport
	return NewDiscoveryWithClient(arm, azure.NewHTTPClient(credential, options))
}

type discoveryOptions struct {
	Top          int    `json:"$top"`
	ResultFormat string `json:"resultFormat"`
	SkipToken    string `json:"$skipToken,omitempty"`
}
type discoveryRequest struct {
	Subscriptions []string         `json:"subscriptions"`
	Query         string           `json:"query"`
	Options       discoveryOptions `json:"options"`
}
type discoveryPage struct {
	rows  []json.RawMessage
	total int64
	token string
}

func discoveryCount(raw json.RawMessage) (int64, error) {
	var n int64
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) || json.Unmarshal(raw, &n) != nil || n < 0 {
		return 0, fmt.Errorf("invalid AI discovery count")
	}
	return n, nil
}

// Validate all envelope/continuation fields before admitting current-page rows.
func decodeDiscoveryPage(ctx context.Context, raw []byte) (discoveryPage, error) {
	var page discoveryPage
	if err := validateJSON(ctx, raw); err != nil {
		return page, err
	}
	object, err := wireObject(bytes.TrimSpace(raw), "count", "totalRecords", "resultTruncated", "data", "$skipToken")
	if err != nil {
		return page, err
	}
	count, err := discoveryCount(object["count"])
	if err != nil {
		return page, err
	}
	page.total, err = discoveryCount(object["totalRecords"])
	if err != nil || count > page.total {
		return page, fmt.Errorf("invalid AI discovery total")
	}
	flag, err := wireString(object["resultTruncated"], 5)
	if err != nil || flag == nil || (*flag != "true" && *flag != "false") {
		return page, fmt.Errorf("invalid AI discovery truncation")
	}
	page.rows, err = wireArray(object["data"], MaxDiscoveryPageRows, false)
	if err != nil || int64(len(page.rows)) != count {
		return page, fmt.Errorf("invalid AI discovery array or count")
	}
	for _, row := range page.rows {
		// Nonobjects and aliased known columns are ambiguous page data. Wrong
		// typed projected values are isolated invalid rows instead.
		if _, err := wireObject(row, discoveryFields...); err != nil {
			return page, err
		}
	}
	if token, present := object["$skipToken"]; present {
		value, err := wireString(token, MaxNextLinkBytes)
		if err != nil || value == nil || *value == "" {
			return page, fmt.Errorf("invalid AI discovery continuation")
		}
		page.token = *value
	}
	if *flag == "true" && page.token == "" {
		return page, fmt.Errorf("incomplete AI discovery continuation")
	}
	return page, nil
}

var discoveryFields = []string{"id", "subscriptionId", "resourceGroup", "location", "type", "name", "sku_name", "sku_tier", "kind"}

func decodeDiscoveryAccount(raw json.RawMessage, batch map[string]bool) (LocatedAccount, int, bool) {
	var account LocatedAccount
	object, err := wireObject(raw, discoveryFields...)
	if err != nil {
		return account, 0, false
	}
	values := make([]string, len(discoveryFields))
	text := 0
	valid := true
	for i, field := range discoveryFields {
		limit := MaxLabelBytes
		if field == "id" {
			limit = 2048
		}
		rawValue := object[field]
		if len(rawValue) == 0 || bytes.Equal(rawValue, []byte("null")) {
			continue
		}
		var value string
		if json.Unmarshal(rawValue, &value) != nil {
			valid = false
			continue
		}
		// Charge decoded projected text even when this row will be rejected.
		text += len(value)
		if !safe(value, limit) {
			valid = false
		} else {
			values[i] = value
		}
	}
	if !valid {
		return account, text, false
	}
	account = LocatedAccount{Account: Account{ID: values[0], SubscriptionID: values[1], ResourceGroup: values[2], Name: values[5], SKU: values[6], Kind: values[8]}, Region: strings.ToLower(values[3])}
	match := accountID.FindStringSubmatch(account.ID)
	parts := strings.Split(account.ID, "/")
	// ARM structural names and regional DNS keys permit ASCII casing, not
	// Unicode case-fold aliases (for example long-s or the Kelvin sign).
	if len(parts) != 9 || strings.ToLower(parts[1]) != "subscriptions" || strings.ToLower(parts[3]) != "resourcegroups" || strings.ToLower(parts[5]) != "providers" || strings.ToLower(parts[6]) != "microsoft.cognitiveservices" || strings.ToLower(parts[7]) != "accounts" || strings.ToLower(values[4]) != "microsoft.cognitiveservices/accounts" || len(values[3]) > 64 || strings.IndexFunc(values[3], func(r rune) bool { return r > 127 }) >= 0 {
		return account, text, false
	}
	if len(match) != 4 || !batch[strings.ToLower(account.SubscriptionID)] || !strings.EqualFold(match[1], account.SubscriptionID) || !strings.EqualFold(match[2], account.ResourceGroup) || !strings.EqualFold(match[3], account.Name) || !regionName.MatchString(account.Region) || strings.ContainsAny(account.ID, "%?#\\") || account.ResourceGroup == "." || account.ResourceGroup == ".." || account.Name == "." || account.Name == ".." {
		return account, text, false
	}
	kind := strings.ToLower(account.Kind)
	if kind != "" && !strings.Contains(kind, "openai") && !strings.Contains(kind, "aiservices") {
		return account, text, false
	}
	return account, text, true
}

func (d *Discovery) Discover(ctx context.Context, subscriptions map[string]string, filter Filter) (DiscoveryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, MaxDuration)
	defer cancel()
	out := DiscoveryResult{Accounts: []LocatedAccount{}, Health: assessment.StageExecution{Name: "ai-discovery", Status: assessment.StageCompleted}}
	invalid := 0
	finish := func(code string, cause error) (DiscoveryResult, error) {
		out.Health.Records = len(out.Accounts)
		if invalid > 0 {
			out.Health.Warnings = []assessment.AssessmentWarning{{Code: "ai_discovery_invalid_rows", Message: fmt.Sprintf("skipped %d invalid or duplicate AI account rows", invalid)}}
			out.Health.Status = assessment.StageCompletedWithWarnings
		}
		if code != "" {
			out.Health.Status = assessment.StageFailed
			out.Health.Error = &assessment.AssessmentError{Code: code, Message: "AI account discovery returned incomplete coverage"}
			if cancellation := safeCause(ctx, cause); cancellation != nil {
				out.Health.Error.Code = "ai_discovery_cancelled"
				return out, fmt.Errorf("AI account discovery stopped: %w", cancellation)
			}
			return out, fmt.Errorf("AI account discovery returned incomplete coverage")
		}
		return out, nil
	}
	if len(subscriptions) > MaxAccounts {
		return finish("ai_discovery_scope_limit", nil)
	}
	scope := make([]string, 0, len(subscriptions))
	ownership := map[string]bool{}
	textBytes := 0
	for id, name := range subscriptions {
		if err := ctx.Err(); err != nil {
			return finish("ai_discovery_cancelled", err)
		}
		key := strings.ToLower(id)
		if !subscriptionID.MatchString(id) || !safe(name, MaxLabelBytes) || ownership[key] {
			return finish("ai_discovery_scope_invalid", nil)
		}
		ownership[key] = true
		scope = append(scope, key)
		textBytes += len(key) + len(name)
	}
	sort.Strings(scope)
	requests, received, responseBytes := 0, 0, 0
	seenIDs := map[string]bool{}
	for start := 0; start < len(scope); start += MaxDiscoveryBatch {
		end := min(start+MaxDiscoveryBatch, len(scope))
		batch := map[string]bool{}
		for _, id := range scope[start:end] {
			batch[id] = true
		}
		request := discoveryRequest{Subscriptions: append([]string(nil), scope[start:end]...), Query: DiscoveryQuery, Options: discoveryOptions{Top: MaxDiscoveryPageRows, ResultFormat: "objectArray"}}
		seenTokens := map[string]bool{}
		var total int64 = -1
		batchRows := int64(0)
		for pageNumber := 0; ; pageNumber++ {
			if err := ctx.Err(); err != nil {
				return finish("ai_discovery_cancelled", err)
			}
			if pageNumber >= MaxDiscoveryPages || requests >= MaxDiscoveryRequests {
				return finish("ai_discovery_request_limit", nil)
			}
			requests++
			body, _ := json.Marshal(request)
			raw, _, err := d.client.PostBounded(ctx, d.origin.String()+"/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01", azure.NopReadSeekCloser{Reader: bytes.NewReader(body)}, MaxPageBytes)
			if err != nil {
				return finish("ai_discovery_request_failed", err)
			}
			responseBytes += len(raw)
			if len(raw) > MaxPageBytes || responseBytes > MaxResponseBytes {
				return finish("ai_discovery_body_limit", nil)
			}
			page, err := decodeDiscoveryPage(ctx, raw)
			if err != nil {
				return finish("ai_discovery_page_invalid", err)
			}
			if (total >= 0 && total != page.total) || (page.token != "" && seenTokens[page.token]) || batchRows+int64(len(page.rows)) > page.total || (page.token == "" && batchRows+int64(len(page.rows)) != page.total) {
				return finish("ai_discovery_coverage_invalid", nil)
			}
			if received+len(page.rows) > MaxAccounts {
				return finish("ai_discovery_row_limit", nil)
			}
			textBytes += len(page.token)
			if textBytes > MaxTextBytes {
				return finish("ai_discovery_text_limit", nil)
			}
			total = page.total
			batchRows += int64(len(page.rows))
			received += len(page.rows)
			for _, row := range page.rows {
				if err := ctx.Err(); err != nil {
					return finish("ai_discovery_cancelled", err)
				}
				account, text, valid := decodeDiscoveryAccount(row, batch)
				textBytes += text
				if textBytes > MaxTextBytes {
					return finish("ai_discovery_text_limit", nil)
				}
				key := strings.ToLower(account.ID)
				if !valid || seenIDs[key] {
					invalid++
					continue
				}
				seenIDs[key] = true
				excluded := filter != nil && filter.IsServiceExcluded(account.ID)
				if err := ctx.Err(); err != nil {
					return finish("ai_discovery_cancelled", err)
				}
				if !excluded {
					out.Accounts = append(out.Accounts, account)
				}
			}
			if page.token == "" {
				break
			}
			seenTokens[page.token] = true
			request.Options.SkipToken = page.token
		}
	}
	if err := ctx.Err(); err != nil {
		return finish("ai_discovery_cancelled", err)
	}
	return finish("", nil)
}
