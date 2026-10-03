// Metrics/deployment requests reproduce the pinned public-cloud source contract
// with explicit incomplete health and guarded continuations. See request plan.
package aigov

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
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const (
	MaxRequests      = 512
	MaxPages         = 64
	MaxResponseBytes = 16 << 20 // Successful bodies across all metrics/deployment calls.
	MaxWorkers       = 5
	MaxDuration      = 5 * time.Minute
)

var regionName = regexp.MustCompile(`^[a-z][a-z0-9]{0,63}$`)
var errBudget = errors.New("AI request budget exceeded")

type BoundedPoster interface {
	PostBounded(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)
}
type BoundedGetter interface {
	GetBoundedWithResponse(context.Context, string, int64) ([]byte, *http.Response, error)
}
type LocatedAccount struct {
	Account
	Region string
}
type Scanner struct {
	origin      *url.URL
	metrics     BoundedPoster
	deployments BoundedGetter
	now         func() time.Time
}

func originURL(endpoint string) (*url.URL, error) {
	u, e := url.Parse(endpoint)
	if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Opaque != "" || u.RawPath != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.ContainsAny(endpoint, "#\\\r\n\t") || (u.Path != "" && u.Path != "/") {
		return nil, fmt.Errorf("invalid AI ARM endpoint")
	}
	u.Path = ""
	return u, nil
}

// Explicit clients are trusted application code, must honor context/byte limits
// and never follow redirects. Production clients capture separate audiences.
func NewWithClients(armEndpoint string, metrics BoundedPoster, deployments BoundedGetter, now func() time.Time) (*Scanner, error) {
	origin, e := originURL(armEndpoint)
	if e != nil || metrics == nil || deployments == nil || now == nil {
		return nil, fmt.Errorf("invalid AI endpoint, clients or clock")
	}
	return &Scanner{origin: origin, metrics: metrics, deployments: deployments, now: now}, nil
}
func New(credential azcore.TokenCredential) (*Scanner, error) {
	configuration := azure.CloudConfiguration()
	public := cloud.AzurePublic
	arm := azure.ResourceManagerEndpoint()
	scope := azure.ResourceManagerScope()
	if credential == nil || arm != strings.TrimRight(public.Services[cloud.ResourceManager].Endpoint, "/") || scope != strings.TrimRight(public.Services[cloud.ResourceManager].Audience, "/")+"/.default" || strings.TrimRight(configuration.ActiveDirectoryAuthorityHost, "/") != strings.TrimRight(public.ActiveDirectoryAuthorityHost, "/") {
		return nil, fmt.Errorf("AI governance metrics public-cloud configuration required")
	}
	armOptions := azure.DefaultHTTPClientOptions(30 * time.Second)
	armOptions.Scope = scope
	metricOptions := azure.DefaultHTTPClientOptions(30 * time.Second)
	metricOptions.Scope = "https://metrics.monitor.azure.com/.default"
	return NewWithClients(arm, azure.NewHTTPClient(credential, metricOptions), azure.NewHTTPClient(credential, armOptions), time.Now)
}

type requestBudget struct {
	requests atomic.Int64
	bytes    atomic.Int64
	entries  atomic.Int64
}

func (b *requestBudget) before(ctx context.Context) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if b.bytes.Load() > MaxResponseBytes || b.requests.Add(1) > MaxRequests {
		return errBudget
	}
	return nil
}
func (b *requestBudget) after(data []byte) error {
	if len(data) > MaxPageBytes || b.bytes.Add(int64(len(data))) > MaxResponseBytes {
		return errBudget
	}
	return nil
}
func safeCause(ctx context.Context, e error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if errors.Is(e, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(e, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	return nil
}
func (s *Scanner) Scan(ctx context.Context, subscriptions map[string]string, located []LocatedAccount, filter Filter) (assessment.PluginTable, error) {
	ctx, cancel := context.WithTimeout(ctx, MaxDuration)
	defer cancel()
	budget := &requestBudget{}
	accounts := make([]Account, 0, len(located))
	for _, a := range located {
		accounts = append(accounts, a.Account)
	}
	preflight, e := Project(ctx, subscriptions, accounts, nil, nil, nil)
	if e != nil {
		return preflight, e
	}
	selected := []Account{}
	points := []Point{}
	enrichment := map[string]DeploymentSet{}
	malformed, metricFailures, deploymentFailures := 0, 0, 0
	failureCode := ""
	var failureCause error
	mark := func(code string, e error) {
		if failureCode == "" {
			failureCode = code
		}
		if cause := safeCause(ctx, e); cause != nil {
			failureCode = "ai_cancelled"
			failureCause = cause
		}
	}
	finish := func() (assessment.PluginTable, error) {
		table, projectionErr := Project(context.WithoutCancel(ctx), subscriptions, selected, points, enrichment, nil)
		if projectionErr != nil && len(enrichment) != 0 {
			// Enrichment is secondary to received metrics. Its aggregate limits can
			// fail before Project reaches any points; retry without metadata so
			// previously accepted metrics survive with source-default cells. Keep
			// the original failed health even if the metric-only projection succeeds.
			code := table.Health.Error.Code
			table, _ = Project(context.WithoutCancel(ctx), subscriptions, selected, points, nil, nil)
			mark(code, projectionErr)
		}
		if malformed > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "ai_malformed_response", Message: fmt.Sprintf("skipped %d invalid AI response entries", malformed)})
		}
		if metricFailures > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "ai_metrics_incomplete", Message: fmt.Sprintf("metrics coverage incomplete for %d account responses", metricFailures)})
		}
		if deploymentFailures > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "ai_deployment_requests_failed", Message: fmt.Sprintf("deployment retrieval incomplete for %d accounts", deploymentFailures)})
		}
		if projectionErr != nil && table.Health.Error != nil {
			mark(table.Health.Error.Code, projectionErr)
		}
		if failureCode != "" {
			table.Health.Status = assessment.StageFailed
			table.Health.Error = &assessment.AssessmentError{Code: failureCode, Message: "AI governance requests returned incomplete data"}
		} else if len(table.Health.Warnings) > 0 {
			table.Health.Status = assessment.StageCompletedWithWarnings
		}
		if e := assessment.ValidatePluginTables([]assessment.PluginTable{table}); e != nil {
			table.Rows = []assessment.PluginRow{}
			table.Health.Records = 0
			table.Health.Status = assessment.StageFailed
			table.Health.Error = &assessment.AssessmentError{Code: "ai_output_invalid", Message: "AI governance output contract invalid"}
			return table, fmt.Errorf("AI governance output contract invalid")
		}
		if failureCode != "" {
			if failureCause != nil {
				return table, fmt.Errorf("AI governance requests stopped: %w", failureCause)
			}
			return table, fmt.Errorf("AI governance requests returned incomplete data")
		}
		return table, nil
	}
	groups := map[string][]LocatedAccount{}
	for _, a := range located {
		if e := ctx.Err(); e != nil {
			mark("ai_cancelled", e)
			return finish()
		}
		region := strings.ToLower(a.Region)
		if !regionName.MatchString(region) {
			mark("ai_region_invalid", nil)
			return finish()
		}
		a.Region = region
		if filter != nil && filter.IsServiceExcluded(a.ID) {
			continue
		}
		if e := ctx.Err(); e != nil {
			mark("ai_cancelled", e)
			return finish()
		}
		selected = append(selected, a.Account)
		key := strings.ToLower(a.SubscriptionID) + "|" + region
		groups[key] = append(groups[key], a)
	}
	if len(selected) == 0 {
		return finish()
	}
	end := s.now().UTC().Truncate(time.Second)
	if end.Year() < 1 || end.Year() > 9999 || end.Add(-167*time.Hour).Year() < 1 {
		mark("ai_clock_invalid", nil)
		return finish()
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		group := groups[k]
		for start := 0; start < len(group); start += MaxBatchAccounts {
			if e := ctx.Err(); e != nil {
				mark("ai_cancelled", e)
				return finish()
			}
			batch := group[start:min(start+MaxBatchAccounts, len(group))]
			ids := make([]string, len(batch))
			for i, a := range batch {
				ids[i] = a.ID
			}
			endpoint := url.URL{Scheme: "https", Host: batch[0].Region + ".metrics.monitor.azure.com", Path: "/subscriptions/" + strings.ToLower(batch[0].SubscriptionID) + "/metrics:getBatch"}
			q := url.Values{"api-version": {"2024-02-01"}, "metricnamespace": {"Microsoft.CognitiveServices/accounts"}, "metricnames": {"AzureOpenAIRequests"}, "aggregation": {"Count"}, "interval": {"PT1H"}, "filter": {"StatusCode eq '*' and ModelDeploymentName eq '*' and ModelName eq '*'"}, "starttime": {end.Add(-167 * time.Hour).Format(time.RFC3339)}, "endtime": {end.Format(time.RFC3339)}}
			endpoint.RawQuery = q.Encode()
			body, _ := json.Marshal(struct {
				IDs []string `json:"resourceids"`
			}{ids})
			if e := budget.before(ctx); e != nil {
				mark("ai_request_limit", e)
				return finish()
			}
			raw, response, e := s.metrics.PostBounded(ctx, endpoint.String(), azure.NopReadSeekCloser{Reader: bytes.NewReader(body)}, MaxPageBytes)
			if e != nil || response == nil || response.StatusCode != http.StatusOK {
				mark("ai_metrics_request_failed", e)
				metricFailures += len(batch)
				if ctx.Err() != nil {
					return finish()
				}
				continue
			}
			if e = budget.after(raw); e != nil {
				mark("ai_response_limit", nil)
				return finish()
			}
			page, decodeErr := DecodeMetrics(ctx, raw, ids)
			malformed += page.Malformed
			metricFailures += len(page.FailedIDs) + len(page.Unreported)
			if decodeErr != nil {
				mark("ai_metrics_decode_failed", decodeErr)
			}
			if len(page.FailedIDs)+len(page.Unreported) > 0 {
				mark("ai_metrics_incomplete", nil)
			}
			if len(points)+len(page.Points) > MaxPoints {
				mark("ai_point_limit", nil)
				return finish()
			}
			points = append(points, page.Points...)
			if ctx.Err() != nil {
				mark("ai_cancelled", ctx.Err())
				return finish()
			}
			haveMetrics := map[string]bool{}
			for _, p := range page.Points {
				if p.Timestamp != nil && p.Count != nil {
					haveMetrics[strings.ToLower(p.ResourceID)] = true
				}
			}
			outcomes := make([]deploymentOutcome, len(batch))
			jobs := make(chan int)
			var wg sync.WaitGroup
			for w := 0; w < min(MaxWorkers, len(batch)); w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := range jobs {
						if haveMetrics[strings.ToLower(batch[i].ID)] {
							outcomes[i] = s.fetchDeployments(ctx, batch[i], budget)
						}
					}
				}()
			}
			for i := range batch {
				jobs <- i
			}
			close(jobs)
			wg.Wait()
			for i, o := range outcomes {
				if !haveMetrics[strings.ToLower(batch[i].ID)] {
					continue
				}
				enrichment[batch[i].ID] = DeploymentSet{Values: o.values, Failed: o.code != ""}
				malformed += o.malformed
				if o.code != "" {
					deploymentFailures++
					if o.code == "ai_request_limit" || o.code == "ai_response_limit" || o.code == "ai_deployment_limit" {
						mark(o.code, o.cause)
					}
					if cause := safeCause(ctx, o.cause); cause != nil {
						mark("ai_cancelled", cause)
					}
				}
			}
			if ctx.Err() != nil {
				mark("ai_cancelled", ctx.Err())
				return finish()
			}
		}
	}
	return finish()
}

type deploymentOutcome struct {
	values    []Deployment
	malformed int
	code      string
	cause     error
}

func (s *Scanner) continuation(raw, path string) (string, error) {
	if len(raw) > MaxNextLinkBytes || strings.ContainsAny(raw, "#\\\r\n\t") {
		return "", fmt.Errorf("invalid AI deployment continuation")
	}
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "https" || !strings.EqualFold(u.Host, s.origin.Host) || u.User != nil || u.Opaque != "" || u.RawPath != "" || u.Fragment != "" || u.ForceQuery || !strings.EqualFold(u.Path, path) {
		return "", fmt.Errorf("unsafe AI deployment continuation")
	}
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil || len(q) > 8 || q.Get("api-version") != "2025-06-01" {
		return "", fmt.Errorf("invalid AI deployment continuation query")
	}
	seen := map[string]bool{}
	for k, v := range q {
		if len(v) != 1 || !safe(k, 64) || k == "" || !safe(v[0], MaxNextLinkBytes) || seen[strings.ToLower(k)] || strings.EqualFold(k, "api-version") && k != "api-version" {
			return "", fmt.Errorf("ambiguous AI deployment continuation query")
		}
		seen[strings.ToLower(k)] = true
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
func (s *Scanner) fetchDeployments(ctx context.Context, a LocatedAccount, budget *requestBudget) deploymentOutcome {
	out := deploymentOutcome{values: []Deployment{}}
	first := *s.origin
	first.Path = a.ID + "/deployments"
	first.RawQuery = "api-version=2025-06-01"
	endpoint := first.String()
	seen := map[string]bool{}
	for pageIndex := 0; ; pageIndex++ {
		if pageIndex >= MaxPages {
			out.code = "ai_deployment_page_limit"
			return out
		}
		next, e := s.continuation(endpoint, first.Path)
		if e != nil || seen[next] {
			out.code = "ai_deployment_continuation_invalid"
			return out
		}
		seen[next] = true
		if e = budget.before(ctx); e != nil {
			out.code = "ai_request_limit"
			out.cause = e
			return out
		}
		raw, response, e := s.deployments.GetBoundedWithResponse(ctx, next, MaxPageBytes)
		if e != nil || response == nil || response.StatusCode != http.StatusOK {
			out.code = "ai_deployment_request_failed"
			out.cause = e
			return out
		}
		if e = budget.after(raw); e != nil {
			out.code = "ai_response_limit"
			return out
		}
		page, e := DecodeDeployments(ctx, raw)
		if budget.entries.Add(int64(len(page.Values))) > MaxDeployments {
			out.code = "ai_deployment_limit"
			return out
		}
		out.malformed += page.Malformed
		if len(out.values)+len(page.Values) > MaxDeployments {
			out.code = "ai_deployment_limit"
			return out
		}
		out.values = append(out.values, page.Values...)
		if e != nil {
			out.code = "ai_deployment_decode_failed"
			out.cause = e
			return out
		}
		if page.NextLink == "" {
			return out
		}
		endpoint = page.NextLink
	}
}
