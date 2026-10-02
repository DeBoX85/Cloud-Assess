// Zone mapping behavior is characterized from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md and docs/ZONE_MAPPING.md.
package zone

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	APIVersion                 = "2022-12-01"
	MaxPageBytes         int64 = 1 << 20
	MaxSubscriptionBytes       = 8 << 20
	MaxPages                   = 64
	MaxSubscriptionRows        = 2048
	MaxSubscriptions           = 256
	MaxRows                    = 65536
	MaxWorkers                 = 5
)

// Getter must honor cancellation, byte bounds and the no-redirect contract.
type Getter interface {
	GetBounded(context.Context, string, int64) ([]byte, error)
}

type Row struct {
	SubscriptionID, SubscriptionName, Location, DisplayName, LogicalZone, PhysicalZone string
}

// Failure is intentionally sanitized: no raw provider body or continuation URL.
type Failure struct{ SubscriptionID, Code string }

// Result retains successful pages, including those before a later-page failure.
// Callers must project Failures into incomplete stage health before exposure.
type Result struct {
	Rows     []Row
	Failures []Failure
}

type Scanner struct {
	origin *url.URL
	getter Getter
}

func NewScanner(endpoint string, getter Getter) (*Scanner, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Fragment != "" || u.RawQuery != "" || u.RawPath != "" || (u.Path != "" && u.Path != "/") || strings.ContainsAny(endpoint, "\\\r\n\t") || getter == nil {
		return nil, fmt.Errorf("invalid zone mapping ARM endpoint or getter")
	}
	u.Path = ""
	return &Scanner{origin: u, getter: getter}, nil
}

var subscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
var errRowLimit = errors.New("zone mapping page row limit exceeded")

type subscription struct{ id, name string }
type outcome struct {
	rows []Row
	code string
}

func (s *Scanner) Scan(ctx context.Context, subscriptions map[string]string) (Result, error) {
	if len(subscriptions) > MaxSubscriptions {
		return Result{}, fmt.Errorf("zone mapping subscription limit exceeded")
	}
	selected := make([]subscription, 0, len(subscriptions))
	seen := make(map[string]bool, len(subscriptions))
	for id, name := range subscriptions {
		if !subscriptionID.MatchString(id) || seen[strings.ToLower(id)] {
			return Result{}, fmt.Errorf("invalid or duplicate zone mapping subscription ID")
		}
		seen[strings.ToLower(id)] = true
		selected = append(selected, subscription{strings.ToLower(id), name})
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].id < selected[j].id })
	outcomes := make([]outcome, len(selected))
	jobs := make(chan int)
	var workers sync.WaitGroup
	for w := 0; w < min(MaxWorkers, len(selected)); w++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for i := range jobs {
				outcomes[i] = s.fetch(ctx, selected[i])
			}
		}()
	}
	scheduled := 0
schedule:
	for i := range selected {
		if ctx.Err() != nil {
			break
		}
		select {
		case jobs <- i:
			scheduled++
		case <-ctx.Done():
			break schedule
		}
	}
	close(jobs)
	workers.Wait()
	for i := scheduled; i < len(outcomes); i++ {
		outcomes[i].code = "zone_cancelled"
	}
	result := Result{Rows: []Row{}, Failures: []Failure{}}
	for i, out := range outcomes {
		if len(result.Rows)+len(out.rows) > MaxRows {
			out.rows = nil
			out.code = "zone_row_limit"
		}
		result.Rows = append(result.Rows, out.rows...)
		if out.code != "" {
			result.Failures = append(result.Failures, Failure{selected[i].id, out.code})
		}
	}
	sort.SliceStable(result.Rows, func(i, j int) bool {
		a, b := result.Rows[i], result.Rows[j]
		left := [6]string{a.SubscriptionName, a.Location, a.LogicalZone, a.SubscriptionID, a.PhysicalZone, a.DisplayName}
		right := [6]string{b.SubscriptionName, b.Location, b.LogicalZone, b.SubscriptionID, b.PhysicalZone, b.DisplayName}
		for k := range left {
			if left[k] != right[k] {
				return left[k] < right[k]
			}
		}
		return false
	})
	return result, ctx.Err()
}

func (s *Scanner) fetch(ctx context.Context, sub subscription) outcome {
	first := *s.origin
	first.Path = "/subscriptions/" + sub.id + "/locations"
	first.RawQuery = "api-version=" + APIVersion
	endpoint := first.String()
	seen := make(map[string]bool)
	var out outcome
	totalBytes := 0
	for page := 0; ; page++ {
		if ctx.Err() != nil {
			out.code = "zone_cancelled"
			return out
		}
		if page >= MaxPages {
			out.code = "zone_page_limit"
			return out
		}
		u, err := s.continuation(&first, endpoint)
		if err != nil {
			out.code = "zone_unsafe_continuation"
			return out
		}
		if seen[u.String()] {
			out.code = "zone_continuation_cycle"
			return out
		}
		seen[u.String()] = true
		body, err := s.getter.GetBounded(ctx, u.String(), MaxPageBytes)
		if err != nil {
			out.code = "zone_request_failed"
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil {
				out.code = "zone_cancelled"
			}
			return out
		}
		totalBytes += len(body)
		if int64(len(body)) > MaxPageBytes || totalBytes > MaxSubscriptionBytes {
			out.code = "zone_byte_limit"
			return out
		}
		rows, next, err := decodePage(body, sub)
		if err != nil {
			out.code = "zone_invalid_response"
			if errors.Is(err, errRowLimit) {
				out.code = "zone_row_limit"
			}
			return out
		}
		if len(out.rows)+len(rows) > MaxSubscriptionRows {
			out.code = "zone_row_limit"
			return out
		}
		out.rows = append(out.rows, rows...)
		if next == "" {
			return out
		}
		endpoint = next
	}
}

// Every URL is checked before invoking the authenticated getter.
func (s *Scanner) continuation(first *url.URL, link string) (*url.URL, error) {
	if len(link) > 8192 || strings.ContainsAny(link, "\\\r\n\t") {
		return nil, fmt.Errorf("unsafe continuation")
	}
	parsed, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("unsafe continuation")
	}
	u := first.ResolveReference(parsed)
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil || u.Scheme != "https" || !strings.EqualFold(u.Host, s.origin.Host) || u.User != nil || u.Fragment != "" || u.RawPath != "" || u.Path != first.Path || len(q["api-version"]) != 1 || q.Get("api-version") != APIVersion {
		return nil, fmt.Errorf("unsafe continuation")
	}
	return u, nil
}

func decodePage(body []byte, sub subscription) ([]Row, string, error) {
	if !utf8.Valid(body) || validateJSON(body) != nil {
		return nil, "", fmt.Errorf("invalid locations JSON")
	}
	var envelope struct {
		Value    json.RawMessage `json:"value"`
		NextLink *string         `json:"nextLink"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, "", fmt.Errorf("invalid locations envelope")
	}
	if len(envelope.Value) == 0 || bytes.Equal(bytes.TrimSpace(envelope.Value), []byte("null")) {
		return nil, "", fmt.Errorf("missing locations array")
	}
	var locations []*struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Mappings    []*struct {
			Logical  string `json:"logicalZone"`
			Physical string `json:"physicalZone"`
		} `json:"availabilityZoneMappings"`
	}
	if err := json.Unmarshal(envelope.Value, &locations); err != nil {
		return nil, "", fmt.Errorf("invalid locations array")
	}
	rows := []Row{}
	for _, loc := range locations {
		if loc == nil {
			return nil, "", fmt.Errorf("null location")
		}
		for _, mapping := range loc.Mappings {
			if mapping == nil || loc.Name == "" || mapping.Logical == "" || mapping.Physical == "" {
				return nil, "", fmt.Errorf("missing zone identity")
			}
			if len(rows) == MaxSubscriptionRows {
				return nil, "", errRowLimit
			}
			rows = append(rows, Row{sub.id, sub.name, loc.Name, loc.DisplayName, mapping.Logical, mapping.Physical})
		}
	}
	next := ""
	if envelope.NextLink != nil {
		next = *envelope.NextLink
	}
	return rows, next, nil
}

// Reject duplicate keys (including case variants understood by encoding/json)
// and excessive depth instead of allowing ambiguous completeness evidence.
func validateJSON(body []byte) error {
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var value func(int) error
	value = func(depth int) error {
		if depth > 64 {
			return fmt.Errorf("JSON depth limit")
		}
		token, err := d.Token()
		if err != nil {
			return err
		}
		delim, container := token.(json.Delim)
		if !container {
			return nil
		}
		if delim != '{' && delim != '[' {
			return fmt.Errorf("invalid JSON container")
		}
		keys := map[string]bool{}
		for d.More() {
			if delim == '{' {
				key, err := d.Token()
				if err != nil {
					return err
				}
				name, ok := key.(string)
				if !ok || keys[strings.ToLower(name)] {
					return fmt.Errorf("duplicate JSON key")
				}
				keys[strings.ToLower(name)] = true
			}
			if err := value(depth + 1); err != nil {
				return err
			}
		}
		_, err = d.Token()
		return err
	}
	if err := value(0); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	return nil
}
