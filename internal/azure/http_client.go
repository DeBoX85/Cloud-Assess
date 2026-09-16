// Portions of this file reproduce authenticated Azure HTTP behavior from Microsoft
// Azure Quick Review (MIT licensed). See NOTICE.md.

package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/DeBoX85/Cloud-Assess/internal/throttling"
)

// NopReadSeekCloser adapts bytes.Reader for Azure SDK retryable request bodies.
type NopReadSeekCloser struct{ *bytes.Reader }

func (NopReadSeekCloser) Close() error { return nil }

var sharedTransport = &http.Transport{
	MaxIdleConns:        200,
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout:     90 * time.Second,
	ForceAttemptHTTP2:   true,
}

type HTTPClientOptions struct {
	Timeout          time.Duration
	MaxRetries       int32
	OperationTimeout time.Duration
	Scope            string
	Transport        policy.Transporter
}

func DefaultHTTPClientOptions(timeout time.Duration) *HTTPClientOptions {
	return &HTTPClientOptions{
		Timeout:          timeout,
		MaxRetries:       5,
		OperationTimeout: timeout * 10,
		Scope:            ResourceManagerScope(),
	}
}

type HTTPClient struct {
	pipeline runtime.Pipeline
}

func NewHTTPClient(credential azcore.TokenCredential, options *HTTPClientOptions) *HTTPClient {
	if options == nil {
		options = DefaultHTTPClientOptions(30 * time.Second)
	}

	retryOptions := policy.RetryOptions{
		MaxRetries:    options.MaxRetries,
		TryTimeout:    options.Timeout,
		RetryDelay:    4 * time.Second,
		MaxRetryDelay: 60 * time.Second,
	}

	var transport policy.Transporter
	if options.Transport != nil {
		transport = options.Transport
	} else {
		transport = &http.Client{
			Transport: sharedTransport,
			Timeout:   options.Timeout + 5*time.Second,
		}
	}

	clientOptions := &policy.ClientOptions{Retry: retryOptions, Transport: transport}
	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{options.Scope}, nil)
	pipeline := runtime.NewPipeline(
		"cloud-assess-http-client",
		"v1.0.0",
		runtime.PipelineOptions{
			PerRetry: []policy.Policy{authPolicy, throttling.NewPolicy()},
		},
		clientOptions,
	)
	return &HTTPClient{pipeline: pipeline}
}

func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	body, _, err := c.doRequest(ctx, http.MethodGet, url, nil)
	return body, err
}

func (c *HTTPClient) Post(ctx context.Context, url string, body io.ReadSeekCloser) ([]byte, *http.Response, error) {
	return c.doRequest(ctx, http.MethodPost, url, body)
}

// PostStream returns an unread response body. The caller must close it.
func (c *HTTPClient) PostStream(ctx context.Context, url string, body io.ReadSeekCloser) (*http.Response, error) {
	request, err := runtime.NewRequest(ctx, http.MethodPost, url)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		if err := request.SetBody(body, "application/json"); err != nil {
			return nil, fmt.Errorf("set request body: %w", err)
		}
	}

	response, err := c.pipeline.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return response, nil
	}

	responseBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("read error response: %w", readErr)
	}
	_ = response.Body.Close()
	response.Body = io.NopCloser(bytes.NewReader(responseBody))
	return response, runtime.NewResponseError(response)
}

func (c *HTTPClient) doRequest(ctx context.Context, method, url string, body io.ReadSeekCloser) ([]byte, *http.Response, error) {
	request, err := runtime.NewRequest(ctx, method, url)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		if err := request.SetBody(body, "application/json"); err != nil {
			return nil, nil, fmt.Errorf("set request body: %w", err)
		}
	}

	response, err := c.pipeline.Do(request)
	if err != nil {
		return nil, nil, fmt.Errorf("execute request: %w", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, nil, fmt.Errorf("read response body: %w", err)
	}
	_ = response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Body = io.NopCloser(bytes.NewReader(responseBody))
		return responseBody, response, runtime.NewResponseError(response)
	}
	return responseBody, response, nil
}
