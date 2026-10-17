package arg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ARG response status %d is not a completed query", response.StatusCode)
	}

	queryResponse := &Response{}
	if value := headerValue(response.Header, "x-ms-user-quota-remaining"); value != "" {
		queryResponse.Quota, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("parse ARG quota header: %w", err)
		}
	}
	if value := headerValue(response.Header, "x-ms-user-quota-resets-after"); value != "" {
		queryResponse.RetryAfter, err = parseResetAfter(value)
		if err != nil {
			return nil, fmt.Errorf("parse ARG reset-after header: %w", err)
		}
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(queryResponse); err != nil {
		return nil, fmt.Errorf("decode ARG response: %w", err)
	}
	if queryResponse.Data == nil {
		return nil, fmt.Errorf("decode ARG response: missing or null data array")
	}
	if err := readARGResponseEnd(io.MultiReader(decoder.Buffered(), response.Body)); err != nil {
		return nil, err
	}
	return queryResponse, nil
}

// Reject any non-whitespace after the one response document, including bytes
// already buffered by json.Decoder. Stream the suffix with fixed memory instead
// of decoding and allocating another arbitrary JSON value.
func readARGResponseEnd(reader io.Reader) error {
	var buffer [1024]byte
	for {
		n, err := reader.Read(buffer[:])
		for _, value := range buffer[:n] {
			switch value {
			case ' ', '\t', '\r', '\n':
			default:
				return fmt.Errorf("decode ARG response: unexpected trailing data")
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read ARG response body: %w", err)
		}
	}
}

func headerValue(headers http.Header, name string) string {
	for key, values := range headers {
		if strings.EqualFold(key, name) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func parseResetAfter(value string) (time.Duration, error) {
	var hours, minutes, seconds int
	if _, err := fmt.Sscanf(value, "%d:%d:%d", &hours, &minutes, &seconds); err != nil {
		return 0, err
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
