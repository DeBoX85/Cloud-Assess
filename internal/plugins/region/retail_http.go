package region

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

const retailEndpoint = "https://prices.azure.com/api/retail/prices"

var errRetailRequest = errors.New("retail request failed")
var errRetailLimit = errors.New("retail response limit exceeded")

// RetailGetter deliberately cannot receive an ARM audience, credential or header.
// Injected implementations must honour cancellation and the supplied byte bound.
type RetailGetter interface {
	GetBounded(context.Context, string, int64) ([]byte, *http.Response, error)
}

// RetailHTTPClient owns its transport. It never shares an authenticated pipeline.
type RetailHTTPClient struct{ client *http.Client }

func NewRetailHTTPClient() *RetailHTTPClient {
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	transport := &http.Transport{
		Protocols: protocols, DisableKeepAlives: true, DisableCompression: true,
		ResponseHeaderTimeout: 30 * time.Second, TLSHandshakeTimeout: 30 * time.Second,
		MaxResponseHeaderBytes: 1 << 20,
	}
	// HTTP/1 with no reused connections also prevents net/http's transparent
	// idempotent-request retries (Transport.shouldRetryRequest in pinned Go).
	return &RetailHTTPClient{client: &http.Client{
		Transport: transport, Timeout: 30 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}
}

func retailURL(raw string) (*url.URL, bool) {
	if len(raw) > 8192 || strings.ContainsAny(raw, "\\#") || strings.IndexFunc(raw, func(r rune) bool { return unicode.IsControl(r) || unicode.IsSpace(r) }) >= 0 {
		return nil, false
	}
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "https" || u.User != nil || u.Opaque != "" || u.Fragment != "" || u.ForceQuery || u.RawPath != "" || u.Path != "/api/retail/prices" {
		return nil, false
	}
	h := strings.ToLower(u.Host)
	if h != "prices.azure.com" && h != "prices.azure.com:443" {
		return nil, false
	}
	u.Host = "prices.azure.com"
	return u, true
}

func retailCause(ctx context.Context, err error) error {
	if ctx != nil && ctx.Err() != nil {
		return ctx.Err()
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	return nil
}

func (c *RetailHTTPClient) GetBounded(ctx context.Context, raw string, limit int64) ([]byte, *http.Response, error) {
	if ctx == nil || c == nil || c.client == nil || limit < 1 || limit > 1<<20 {
		return nil, nil, errRetailRequest
	}
	if cause := retailCause(ctx, nil); cause != nil {
		return nil, nil, cause
	}
	u, ok := retailURL(raw)
	if !ok {
		return nil, nil, errRetailRequest
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, errRetailRequest
	}
	resp, err := c.client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if cause := retailCause(ctx, err); cause != nil {
		return nil, resp, cause
	}
	if err != nil || resp == nil || resp.Body == nil {
		return nil, resp, errRetailRequest
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if cause := retailCause(ctx, err); cause != nil {
		return nil, resp, cause
	}
	if int64(len(body)) > limit {
		return nil, resp, errRetailLimit
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		// Bounded error bodies are counted by the collector, never disclosed.
		return body, resp, errRetailRequest
	}
	return body, resp, nil
}
