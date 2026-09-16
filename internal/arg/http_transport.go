package arg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const ResourceGraphAPIVersion = "2024-04-01"

type streamPoster interface {
	PostStream(context.Context, string, io.ReadSeekCloser) (*http.Response, error)
}

type HTTPTransport struct {
	client   streamPoster
	endpoint string
}

func NewHTTPTransport(credential azcore.TokenCredential) *HTTPTransport {
	client := azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(120*time.Second))
	return NewHTTPTransportWithClient(client, resourceGraphEndpoint())
}

func NewHTTPTransportWithClient(client streamPoster, endpoint string) *HTTPTransport {
	return &HTTPTransport{client: client, endpoint: endpoint}
}

func resourceGraphEndpoint() string {
	return fmt.Sprintf("%s/providers/Microsoft.ResourceGraph/resources?api-version=%s", azure.ResourceManagerEndpoint(), ResourceGraphAPIVersion)
}

func (t *HTTPTransport) Do(ctx context.Context, request Request) (*Response, error) {
	if t == nil || t.client == nil {
		return nil, fmt.Errorf("ARG HTTP client is not configured")
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal ARG request: %w", err)
	}

	response, err := t.client.PostStream(ctx, t.endpoint, azure.NopReadSeekCloser{Reader: bytes.NewReader(body)})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	queryResponse := &Response{}
	if value := response.Header.Get("x-ms-user-quota-remaining"); value != "" {
		queryResponse.Quota, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("parse ARG quota header: %w", err)
		}
	}
	if value := response.Header.Get("x-ms-user-quota-resets-after"); value != "" {
		queryResponse.RetryAfter, err = parseResetAfter(value)
		if err != nil {
			return nil, fmt.Errorf("parse ARG reset-after header: %w", err)
		}
	}

	if err := json.NewDecoder(response.Body).Decode(queryResponse); err != nil {
		return nil, fmt.Errorf("decode ARG response: %w", err)
	}
	_, _ = io.Copy(io.Discard, response.Body)
	return queryResponse, nil
}

func parseResetAfter(value string) (time.Duration, error) {
	var hours, minutes, seconds int
	if _, err := fmt.Sscanf(value, "%d:%d:%d", &hours, &minutes, &seconds); err != nil {
		return 0, err
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
