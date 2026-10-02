// Service-health behavior is characterized from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md and docs/SERVICE_HEALTH.md.
package servicehealth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const (
	Name             = "service-health"
	MaxSubscriptions = 3000
	MaxPages         = 64
	MaxRows          = 65536
	MaxDataBytes     = 16 << 20
	MaxPageBytes     = 2 << 20
	MaxRowBytes      = 8192
	MaxDuration      = 5 * time.Minute
)

type Filter interface{ IsResourceTypeExcluded(string) bool }
type Scanner struct{ transport arg.Transport }

func Metadata() assessment.PluginMetadata {
	return assessment.PluginMetadata{Name: Name, Version: "0.1.0-beta", Description: "Analyzes Azure service health events to determine the percentage of time resources were unaffected by service issues over the last 90 days.", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}
}

func PendingTable() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "availability", Metadata: Metadata(), SheetName: "Service Health Availability", Description: "Azure service health availability analysis by subscription, region, and resource type", Columns: []string{"Subscription ID", "Target Region", "Target Resource Type", "Percentage Without Events", "Events Count", "Affected Resources"}, Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped, Warnings: []assessment.AssessmentWarning{{Code: "plugin_not_run", Message: "requested plugin has not executed"}}}}
}

// Injected transports are trusted code and must honor context/body/request bounds.
func NewWithTransport(transport arg.Transport) *Scanner { return &Scanner{transport: transport} }

type BoundedPoster interface {
	PostBounded(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)
}

// NewWithHTTPClient validates the selected ARM origin before authenticated POSTs.
func NewWithHTTPClient(endpoint string, client BoundedPoster) (*Scanner, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.RawPath != "" || u.Opaque != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") || client == nil {
		return nil, fmt.Errorf("invalid service-health ARM endpoint or client")
	}
	endpoint = strings.TrimSuffix(endpoint, "/") + "/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01"
	return NewWithTransport(arg.NewHTTPTransportWithClient(boundedPoster{client}, endpoint)), nil
}

func New(credential azcore.TokenCredential) (*Scanner, error) {
	return NewWithHTTPClient(azure.ResourceManagerEndpoint(), azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(30*time.Second)))
}

type boundedPoster struct{ client BoundedPoster }

func (p boundedPoster) PostStream(ctx context.Context, endpoint string, body io.ReadSeekCloser) (*http.Response, error) {
	data, response, err := p.client.PostBounded(ctx, endpoint, body, MaxPageBytes)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("nil bounded response")
	}
	if err := validateEnvelope(data); err != nil {
		return nil, err
	}
	copyResponse := *response
	copyResponse.Body = io.NopCloser(bytes.NewReader(data))
	return &copyResponse, nil
}

// Reject ambiguous top-level accounting before the shared ARG decoder. Unknown
// unique envelope metadata remains permitted; projected rows have their own schema.
func validateEnvelope(data []byte) error {
	if !utf8.Valid(data) || !json.Valid(data) {
		return fmt.Errorf("invalid service-health envelope")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return fmt.Errorf("invalid service-health envelope object")
	}
	seen := map[string]bool{}
	for d.More() {
		token, err := d.Token()
		if err != nil {
			return fmt.Errorf("invalid service-health envelope field")
		}
		key, ok := token.(string)
		if !ok || len(key) > 128 || len(seen) >= 128 || seen[strings.ToLower(key)] {
			return fmt.Errorf("ambiguous service-health envelope")
		}
		seen[strings.ToLower(key)] = true
		for _, expected := range []string{"data", "$skipToken", "resultTruncated"} {
			if strings.EqualFold(key, expected) && key != expected {
				return fmt.Errorf("ambiguous service-health envelope field")
			}
		}
		var value json.RawMessage
		if d.Decode(&value) != nil {
			return fmt.Errorf("invalid service-health envelope value")
		}
	}
	if !seen["data"] {
		return fmt.Errorf("missing service-health envelope data")
	}
	return nil
}

type budgetTransport struct {
	transport    arg.Transport
	pages, bytes int
	data         []json.RawMessage
}

func (b *budgetTransport) Do(ctx context.Context, request arg.Request) (*arg.Response, error) {
	if b.pages >= MaxPages {
		return nil, fmt.Errorf("service-health page limit")
	}
	b.pages++
	r, err := b.transport.Do(ctx, request)
	if err != nil {
		return nil, err
	}
	if r == nil || r.Data == nil {
		return nil, fmt.Errorf("missing service-health data")
	}
	if len(r.Data) > int(arg.MaxRowsPerPage) || len(b.data)+len(r.Data) > MaxRows {
		return nil, fmt.Errorf("service-health row limit")
	}
	pageBytes := 0
	for _, raw := range r.Data {
		pageBytes += len(raw)
	}
	if pageBytes > MaxPageBytes || b.bytes+pageBytes > MaxDataBytes {
		return nil, fmt.Errorf("service-health byte limit")
	}
	b.bytes += pageBytes
	for _, raw := range r.Data {
		b.data = append(b.data, append(json.RawMessage(nil), raw...))
	}
	return r, nil
}

var subscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type row struct {
	SubscriptionID string  `json:"subscriptionId"`
	Region         string  `json:"targetRegion"`
	ResourceType   string  `json:"targetResourceType"`
	Percentage     float64 `json:"percentageOfTimeWithoutEvents"`
	Events         int64   `json:"events"`
	Affected       int64   `json:"affectedResources"`
}

func decode(raw json.RawMessage) (row, error) {
	var r row
	if len(raw) > MaxRowBytes || !utf8.Valid(raw) {
		return r, fmt.Errorf("invalid service-health row")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return r, fmt.Errorf("invalid service-health object")
	}
	seen := map[string]bool{}
	for d.More() {
		keyToken, err := d.Token()
		if err != nil {
			return r, err
		}
		key, ok := keyToken.(string)
		if !ok || seen[key] {
			return r, fmt.Errorf("duplicate service-health field")
		}
		seen[key] = true
		var value json.RawMessage
		if err = d.Decode(&value); err != nil {
			return r, err
		}
		if bytes.Equal(value, []byte("null")) {
			return r, fmt.Errorf("null service-health field")
		}
		switch key {
		case "subscriptionId":
			err = json.Unmarshal(value, &r.SubscriptionID)
		case "targetRegion":
			err = json.Unmarshal(value, &r.Region)
		case "targetResourceType":
			err = json.Unmarshal(value, &r.ResourceType)
		case "percentageOfTimeWithoutEvents":
			err = json.Unmarshal(value, &r.Percentage)
		case "events":
			err = json.Unmarshal(value, &r.Events)
		case "affectedResources":
			err = json.Unmarshal(value, &r.Affected)
		default:
			return r, fmt.Errorf("unknown service-health field")
		}
		if err != nil {
			return r, err
		}
	}
	if _, err = d.Token(); err != nil {
		return r, err
	}
	var extra any
	if d.Decode(&extra) != io.EOF || len(seen) != 6 || !subscriptionID.MatchString(r.SubscriptionID) || !validLabel(r.Region) || !validLabel(r.ResourceType) || math.IsNaN(r.Percentage) || math.IsInf(r.Percentage, 0) || r.Percentage < 0 || r.Percentage > 100 || r.Events < 0 || r.Affected < 0 {
		return r, fmt.Errorf("invalid service-health values")
	}
	return r, nil
}
func validLabel(value string) bool {
	return value != "" && len(value) <= 512 && utf8.ValidString(value) && strings.IndexFunc(value, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffe || r == 0xffff }) < 0
}

// Scan preserves valid completed-page rows on later failure. Provider errors never
// enter table health. All budgets and slices are owned by this execution.
func (s *Scanner) Scan(ctx context.Context, subscriptions map[string]string, filter Filter) (assessment.PluginTable, error) {
	ctx, cancel := context.WithTimeout(ctx, MaxDuration)
	defer cancel()
	table := PendingTable()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, StartedAt: time.Now().UTC()}
	fail := func(code string, err error) (assessment.PluginTable, error) {
		table.Health.Status = assessment.StageFailed
		table.Health.Error = &assessment.AssessmentError{Code: code, Message: "service-health returned incomplete data"}
		table.Health.FinishedAt = time.Now().UTC()
		return table, err
	}
	if err := ctx.Err(); err != nil {
		return fail("service_health_cancelled", err)
	}
	if s == nil || s.transport == nil {
		return fail("service_health_not_configured", fmt.Errorf("service-health transport is not configured"))
	}
	if len(subscriptions) > MaxSubscriptions {
		return fail("service_health_scope_limit", fmt.Errorf("service-health subscription limit"))
	}
	selected := make(map[string]string, len(subscriptions))
	scope := map[string]bool{}
	for id, name := range subscriptions {
		if !subscriptionID.MatchString(id) {
			return fail("service_health_scope_invalid", fmt.Errorf("invalid service-health subscription scope"))
		}
		if scope[strings.ToLower(id)] {
			return fail("service_health_scope_invalid", fmt.Errorf("duplicate service-health subscription scope"))
		}
		scope[strings.ToLower(id)] = true
		selected[id] = name
	}
	budget := &budgetTransport{transport: s.transport, data: []json.RawMessage{}}
	_, queryErr := arg.NewClient(budget).Query(ctx, Query, selected)
	if queryErr != nil {
		if ctx.Err() != nil {
			queryErr = ctx.Err()
		} else {
			queryErr = fmt.Errorf("service-health query failed")
		}
	}
	if budget.pages > 0 && len(budget.data) == 0 && queryErr != nil {
		return fail("service_health_query_failed", queryErr)
	}
	table.SheetName = "Service Issues"
	table.Description = "Azure service health availability analysis showing percentage of time without service health events by subscription, region, and resource type (last 90 days)"
	rows := make([]row, 0, len(budget.data))
	malformed := 0
	for _, raw := range budget.data {
		r, err := decode(raw)
		if err != nil || !scope[strings.ToLower(r.SubscriptionID)] {
			malformed++
			continue
		}
		if filter != nil && filter.IsResourceTypeExcluded(r.ResourceType) {
			continue
		}
		rows = append(rows, r)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.Percentage != b.Percentage {
			return a.Percentage < b.Percentage
		}
		if a.SubscriptionID != b.SubscriptionID {
			return a.SubscriptionID < b.SubscriptionID
		}
		if a.Region != b.Region {
			return a.Region < b.Region
		}
		return a.ResourceType < b.ResourceType
	})
	for _, r := range rows {
		table.Rows = append(table.Rows, assessment.PluginRow{SubscriptionID: r.SubscriptionID, Cells: []string{r.SubscriptionID, r.Region, r.ResourceType, fmt.Sprintf("%.2f%%", r.Percentage), fmt.Sprintf("%d", r.Events), fmt.Sprintf("%d", r.Affected)}})
	}
	table.Health.Records = len(table.Rows)
	if malformed > 0 {
		table.Health.Status = assessment.StageCompletedWithWarnings
		table.Health.Warnings = []assessment.AssessmentWarning{{Code: "service_health_malformed_rows", Message: fmt.Sprintf("skipped %d invalid service-health rows", malformed)}}
	}
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		table.Rows = []assessment.PluginRow{}
		table.Health.Records = 0
		return fail("service_health_output_invalid", fmt.Errorf("invalid service-health output"))
	}
	if queryErr != nil {
		return fail("service_health_query_failed", queryErr)
	}
	if err := ctx.Err(); err != nil {
		return fail("service_health_cancelled", err)
	}
	table.Health.FinishedAt = time.Now().UTC()
	return table, nil
}
