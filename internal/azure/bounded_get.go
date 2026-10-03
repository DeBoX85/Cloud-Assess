package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// GetBounded closes the response and bounds both successful and error bodies.
// The caller must validate the destination before calling this authenticated
// client. Injected transports are trusted and must not follow redirects.
func (c *HTTPClient) GetBounded(ctx context.Context, endpoint string, maxBytes int64) ([]byte, error) {
	body, _, err := c.GetBoundedWithResponse(ctx, endpoint, maxBytes)
	return body, err
}

// GetBoundedWithResponse also exposes status/headers for APIs with exact success
// semantics. The transport body is closed; returned bytes are owned. Callers
// still validate destinations before authentication and do not read response.Body.
func (c *HTTPClient) GetBoundedWithResponse(ctx context.Context, endpoint string, maxBytes int64) ([]byte, *http.Response, error) {
	return c.requestBounded(ctx, http.MethodGet, endpoint, nil, maxBytes)
}

// PostBounded bounds every attempt before authentication/retry policies read it,
// closes the response, and returns its headers with owned successful body bytes.
// Callers must validate the destination and read-oriented request contract.
func (c *HTTPClient) PostBounded(ctx context.Context, endpoint string, body io.ReadSeekCloser, maxBytes int64) ([]byte, *http.Response, error) {
	return c.requestBounded(ctx, http.MethodPost, endpoint, body, maxBytes)
}

func (c *HTTPClient) requestBounded(ctx context.Context, method, endpoint string, requestBody io.ReadSeekCloser, maxBytes int64) ([]byte, *http.Response, error) {
	if maxBytes <= 0 || maxBytes > 64<<20 {
		return nil, nil, fmt.Errorf("invalid %s response byte limit", method)
	}
	ctx, cancel := c.operationContext(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, boundedRequestLimitKey{}, maxBytes)
	request, err := runtime.NewRequest(ctx, method, endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("create bounded %s request: %w", method, err)
	}
	if requestBody != nil {
		if err := request.SetBody(requestBody, "application/json"); err != nil {
			return nil, nil, fmt.Errorf("set bounded request body: %w", err)
		}
	}
	response, err := c.pipeline.Do(request)
	if err != nil {
		return nil, nil, fmt.Errorf("execute bounded %s request: %w", method, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return nil, nil, fmt.Errorf("read bounded %s response: %w", method, err)
	}
	if int64(len(body)) > maxBytes {
		return nil, nil, fmt.Errorf("%s response exceeds byte limit", method)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		// The original transport body is still closed by the defer above.
		copyResponse := *response
		copyResponse.Body = io.NopCloser(bytes.NewReader(body))
		return nil, response, runtime.NewResponseError(&copyResponse)
	}
	return body, response, nil
}

type boundedRequestLimitKey struct{}

// Apply the bound at the transport response, before SDK authentication/retry
// policies can consume a challenge or failed-attempt body. Ordinary calls have
// no private context limit and retain their existing behavior.
type boundedTransport struct{ next policy.Transporter }

func (t boundedTransport) Do(request *http.Request) (*http.Response, error) {
	response, err := t.next.Do(request)
	if err == nil && response != nil && response.Body != nil {
		if limit, ok := request.Context().Value(boundedRequestLimitKey{}).(int64); ok {
			response.Body = &boundedResponseBody{ReadCloser: response.Body, remaining: limit + 1}
		}
	}
	return response, err
}

type boundedResponseBody struct {
	io.ReadCloser
	remaining int64
}

func (b *boundedResponseBody) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if b.remaining == 0 {
		return 0, fmt.Errorf("response exceeds byte limit")
	}
	if int64(len(p)) > b.remaining {
		p = p[:b.remaining]
	}
	n, err := b.ReadCloser.Read(p)
	b.remaining -= int64(n)
	return n, err
}
