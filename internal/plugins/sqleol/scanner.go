// SQL EOL behavior is characterized from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md and docs/SQL_EOL.md.
package sqleol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
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
	Name             = "sql-eol"
	MaxSubscriptions = 3000
	MaxPages         = 64
	MaxRows          = 65536
	MaxDataBytes     = 16 << 20
	MaxPageBytes     = 2 << 20
	MaxRowBytes      = 64 << 10
	MaxDuration      = 5 * time.Minute
)

type Filter interface{ IsSubscriptionExcluded(string) bool }
type Scanner struct{ transport arg.Transport }

func Metadata() assessment.PluginMetadata {
	return assessment.PluginMetadata{Name: Name, Version: "0.6.0-beta", Description: "Analyzes SQL Server End-of-Life and Extended Security Update status with host-level ESU billing (once per OSE, per version, at the highest edition), full cost breakdown (VM compute, SQL license, ESU), migration recommendations with conservative GP-only SQL MI cost estimates, and unified SQL MI migration savings and verdict", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}
}

var columns = []string{"Subscription", "Resource Group", "Name", "Location", "Arc Server Name", "Cloud Type", "Service Type", "SQL Version", "Edition", "EOL Status", "ESU Applicable", "ESU Enabled", "ESU Start Date", "ESU End Date", "Migration Target Tier", "Migration Recommendation", "vCores", "Billable Cores", "ESU Monthly Cost/Core", "SQL License Type", "SQL License Cost/Core/Month", "SQL License Monthly Cost", "VM Cost/Core/Month", "Est VM Compute Monthly Cost", "Est ESU Monthly Cost", "ESU Cost Basis", "Patch Ops Monthly Cost", "Current Monthly Cost", "Consolidation Ratio", "Est SQL MI Monthly Cost", "Est SQL MI Monthly Saving", "SQL MI Migration Verdict"}
var fields = []string{"SubscriptionId", "Name", "ResourceGroup", "Subscription", "Location", "ArcServerName", "CloudType", "ServiceType", "SQLVersion", "Edition", "vCores", "BillableCores", "EOLStatus", "ESUApplicable", "ESUEnabled", "MigrationRecommendation", "MigrationTargetTier", "ESUStartDate", "ESUEndDate", "ESUMonthlyCostPerCore", "SQLLicenseType", "SQLLicenseMonthlyCostPerCore", "SQLLicenseMonthlyCost", "VMCostPerCorePerMonth", "EstVMComputeMonthlyCost", "EstESUMonthlyCost", "ESUCostBasis", "PatchOpsMonthlyCost", "CurrentMonthlyCost", "ConsolidationRatio", "EstSQLMIMonthlyCost", "EstSQLMIMonthlySaving", "SQLMIMigrationVerdict"}
var cellFields = []string{"Subscription", "ResourceGroup", "Name", "Location", "ArcServerName", "CloudType", "ServiceType", "SQLVersion", "Edition", "EOLStatus", "ESUApplicable", "ESUEnabled", "ESUStartDate", "ESUEndDate", "MigrationTargetTier", "MigrationRecommendation", "vCores", "BillableCores", "ESUMonthlyCostPerCore", "SQLLicenseType", "SQLLicenseMonthlyCostPerCore", "SQLLicenseMonthlyCost", "VMCostPerCorePerMonth", "EstVMComputeMonthlyCost", "EstESUMonthlyCost", "ESUCostBasis", "PatchOpsMonthlyCost", "CurrentMonthlyCost", "ConsolidationRatio", "EstSQLMIMonthlyCost", "EstSQLMIMonthlySaving", "SQLMIMigrationVerdict"}

func PendingTable() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "lifecycle", Metadata: Metadata(), SheetName: "SQL EOL", Description: "SQL Server End-of-Life and Extended Security Update status with cost analysis", Columns: append([]string(nil), columns...), Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped, Warnings: []assessment.AssessmentWarning{{Code: "plugin_not_run", Message: "requested plugin has not executed"}}}}
}

// Injected transports are trusted code and must honor context/body/request bounds.
func NewWithTransport(transport arg.Transport) *Scanner { return &Scanner{transport: transport} }

type BoundedPoster interface {
	PostBounded(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)
}

// NewWithHTTPClient validates the selected ARM origin before authenticated POSTs.
func NewWithHTTPClient(endpoint string, client BoundedPoster) (*Scanner, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.RawPath != "" || u.Opaque != "" || u.Fragment != "" || strings.Contains(endpoint, "#") || (u.Path != "" && u.Path != "/") || client == nil {
		return nil, fmt.Errorf("invalid sql-eol ARM endpoint or client")
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
		return fmt.Errorf("invalid sql-eol envelope")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return fmt.Errorf("invalid sql-eol envelope object")
	}
	seen := map[string]bool{}
	for d.More() {
		token, err := d.Token()
		if err != nil {
			return fmt.Errorf("invalid sql-eol envelope field")
		}
		key, ok := token.(string)
		if !ok || len(key) > 128 || len(seen) >= 128 || seen[strings.ToLower(key)] {
			return fmt.Errorf("ambiguous sql-eol envelope")
		}
		seen[strings.ToLower(key)] = true
		for _, expected := range []string{"data", "$skipToken", "resultTruncated"} {
			if strings.EqualFold(key, expected) && key != expected {
				return fmt.Errorf("ambiguous sql-eol envelope field")
			}
		}
		var value json.RawMessage
		if d.Decode(&value) != nil {
			return fmt.Errorf("invalid sql-eol envelope value")
		}
	}
	if !seen["data"] {
		return fmt.Errorf("missing sql-eol envelope data")
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
		return nil, fmt.Errorf("sql-eol page limit")
	}
	b.pages++
	r, err := b.transport.Do(ctx, request)
	if err != nil {
		return nil, err
	}
	if r == nil || r.Data == nil {
		return nil, fmt.Errorf("missing sql-eol data")
	}
	if len(r.Data) > int(arg.MaxRowsPerPage) || len(b.data)+len(r.Data) > MaxRows {
		return nil, fmt.Errorf("sql-eol row limit")
	}
	pageBytes := 0
	for _, raw := range r.Data {
		pageBytes += len(raw)
	}
	if pageBytes > MaxPageBytes || b.bytes+pageBytes > MaxDataBytes {
		return nil, fmt.Errorf("sql-eol byte limit")
	}
	b.bytes += pageBytes
	for _, raw := range r.Data {
		b.data = append(b.data, append(json.RawMessage(nil), raw...))
	}
	return r, nil
}

var subscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// decode preserves the source string/null projection without interpreting price text.
func decode(raw json.RawMessage) (assessment.PluginRow, error) {
	var row assessment.PluginRow
	if len(raw) > MaxRowBytes || !utf8.Valid(raw) {
		return row, fmt.Errorf("invalid sql-eol row")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return row, fmt.Errorf("invalid sql-eol object")
	}
	values := map[string]string{}
	for d.More() {
		token, err := d.Token()
		if err != nil {
			return row, fmt.Errorf("invalid sql-eol key")
		}
		key, ok := token.(string)
		if !ok {
			return row, fmt.Errorf("invalid sql-eol key")
		}
		if _, exists := values[key]; exists {
			return row, fmt.Errorf("duplicate sql-eol field")
		}
		known := false
		for _, field := range fields {
			if key == field {
				known = true
				break
			}
		}
		if !known {
			return row, fmt.Errorf("unknown sql-eol field")
		}
		var value string
		if d.Decode(&value) != nil || len(value) > 4096 || !utf8.ValidString(value) || strings.IndexFunc(value, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffe || r == 0xffff }) >= 0 {
			return row, fmt.Errorf("invalid sql-eol value")
		}
		values[key] = value
	}
	if _, err = d.Token(); err != nil {
		return row, err
	}
	var extra any
	if d.Decode(&extra) != io.EOF || len(values) != len(fields) || !subscriptionID.MatchString(values["SubscriptionId"]) {
		return row, fmt.Errorf("invalid sql-eol schema/scope")
	}
	row.SubscriptionID = values["SubscriptionId"]
	for _, key := range cellFields {
		row.Cells = append(row.Cells, values[key])
	}
	return row, nil
}

// Scan owns budgets, scope and table data per call, retains completed pages on
// failure, and preserves received ordering rather than re-sorting price strings.
func (s *Scanner) Scan(ctx context.Context, subscriptions map[string]string, filter Filter) (assessment.PluginTable, error) {
	ctx, cancel := context.WithTimeout(ctx, MaxDuration)
	defer cancel()
	table := PendingTable()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, StartedAt: time.Now().UTC()}
	fail := func(code string, err error) (assessment.PluginTable, error) {
		table.Health.Status = assessment.StageFailed
		table.Health.Error = &assessment.AssessmentError{Code: code, Message: "sql-eol returned incomplete data"}
		table.Health.FinishedAt = time.Now().UTC()
		return table, err
	}
	if err := ctx.Err(); err != nil {
		return fail("sql_eol_cancelled", err)
	}
	if s == nil || s.transport == nil {
		return fail("sql_eol_not_configured", fmt.Errorf("sql-eol transport is not configured"))
	}
	if len(subscriptions) > MaxSubscriptions {
		return fail("sql_eol_scope_limit", fmt.Errorf("sql-eol subscription limit"))
	}
	selected := make(map[string]string, len(subscriptions))
	scope := map[string]bool{}
	for id, name := range subscriptions {
		if !subscriptionID.MatchString(id) || scope[strings.ToLower(id)] {
			return fail("sql_eol_scope_invalid", fmt.Errorf("invalid sql-eol subscription scope"))
		}
		selected[id] = name
		scope[strings.ToLower(id)] = true
	}
	budget := &budgetTransport{transport: s.transport, data: []json.RawMessage{}}
	_, queryErr := arg.NewClient(budget).Query(ctx, Query, selected)
	if queryErr != nil {
		if ctx.Err() != nil {
			queryErr = ctx.Err()
		} else {
			queryErr = fmt.Errorf("sql-eol query failed")
		}
	}
	malformed := 0
	for _, raw := range budget.data {
		row, err := decode(raw)
		if err != nil || !scope[strings.ToLower(row.SubscriptionID)] {
			malformed++
			continue
		}
		if filter != nil && filter.IsSubscriptionExcluded(row.SubscriptionID) {
			continue
		}
		table.Rows = append(table.Rows, row)
	}
	table.Health.Records = len(table.Rows)
	if malformed > 0 {
		table.Health.Status = assessment.StageCompletedWithWarnings
		table.Health.Warnings = []assessment.AssessmentWarning{{Code: "sql_eol_malformed_rows", Message: fmt.Sprintf("skipped %d invalid sql-eol rows", malformed)}}
	}
	if assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
		table.Rows = []assessment.PluginRow{}
		table.Health.Records = 0
		return fail("sql_eol_output_invalid", fmt.Errorf("invalid sql-eol output"))
	}
	if queryErr != nil {
		return fail("sql_eol_query_failed", queryErr)
	}
	if err := ctx.Err(); err != nil {
		return fail("sql_eol_cancelled", err)
	}
	table.Health.FinishedAt = time.Now().UTC()
	return table, nil
}
