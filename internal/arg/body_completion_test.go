package arg

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

type completionReadBody struct {
	reader *strings.Reader

	terminal error

	closed int
}

func (b *completionReadBody) Read(p []byte) (int, error) {
	if b.reader.Len() > 0 {
		return b.reader.Read(p)
	}
	return 0, b.terminal
}

func (b *completionReadBody) Close() error {
	b.closed++
	return nil
}

func TestQueryRejectsTerminalResponseReadFailure(t *testing.T) {
	sentinel := errors.New("synthetic terminal response failure")
	for _, payload := range []string{`{"data":[]}`, `{"data":[{"id":"one"}]}`} {
		t.Run(payload, func(t *testing.T) {
			for _, terminal := range []error{sentinel, io.ErrUnexpectedEOF, context.Canceled, context.DeadlineExceeded} {
				t.Run(terminal.Error(), func(t *testing.T) {
					body := &completionReadBody{
						reader: strings.NewReader(payload),

						terminal: terminal,
					}
					poster := &fakePoster{resp: &http.Response{
						StatusCode: http.StatusOK,
						Body:       body,
					}}
					result, err := NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(
						context.Background(), "resources", map[string]string{"sub": "name"},
					)
					if result != nil || !errors.Is(err, terminal) {
						t.Fatalf("terminal response failure reported as success: result=%#v error=%v; want nil result and %v", result, err, terminal)
					}
					if body.closed != 1 {
						t.Fatalf("body closed %d times; want exactly one", body.closed)
					}
				})
			}
		})
	}
}

func TestQueryAcceptsCompletedResponseAndClosesBody(t *testing.T) {
	for _, payload := range []string{`{"data":[]}`, `{"data":[{"id":"one"}]}`} {
		t.Run(payload, func(t *testing.T) {
			body := &completionReadBody{
				reader: strings.NewReader(payload),

				terminal: io.EOF,
			}
			poster := &fakePoster{resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}}
			result, err := NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(
				context.Background(), "resources", map[string]string{"sub": "name"},
			)
			want := 0
			if payload == `{"data":[{"id":"one"}]}` {
				want = 1
			}
			if err != nil || result == nil || len(result.Data) != want || body.closed != 1 {
				t.Fatalf("completed response result=%#v error=%v closed=%d; want %d rows and one close", result, err, body.closed, want)
			}
		})
	}
}

type completionTransportFunc func(*http.Request) (*http.Response, error)

func (f completionTransportFunc) Do(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestSharedHTTPPipelineRejectsTerminalReadFailure(t *testing.T) {
	sentinel := errors.New("synthetic terminal response failure")
	for _, payload := range []string{`{"data":[]}`, `{"data":[{"id":"one"}]}`} {
		t.Run(payload, func(t *testing.T) {
			for _, terminal := range []error{sentinel, io.ErrUnexpectedEOF, context.Canceled, context.DeadlineExceeded} {
				t.Run(terminal.Error(), func(t *testing.T) {
					body := &completionReadBody{
						reader: strings.NewReader(payload),

						terminal: terminal,
					}
					calls := 0
					client := azure.NewHTTPClient(argTestCredential{}, &azure.HTTPClientOptions{
						MaxRetries:       -1,
						OperationTimeout: time.Second,
						Scope:            "https://management.azure.com/.default",
						Transport: completionTransportFunc(func(request *http.Request) (*http.Response, error) {
							calls++
							if request.Method != http.MethodPost || request.URL.Path != "/graph" || request.Header.Get("Authorization") != "Bearer test-token" {
								t.Fatalf("unexpected authenticated request: method=%s path=%s", request.Method, request.URL.Path)
							}
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     make(http.Header),
								Body:       body,
								Request:    request,
							}, nil
						}),
					})
					result, err := NewClient(NewHTTPTransportWithClient(client, "https://management.azure.com/graph")).Query(
						context.Background(), "resources", map[string]string{"sub": "name"},
					)
					if result != nil || !errors.Is(err, terminal) || calls != 1 || body.closed != 1 {
						t.Fatalf("shared pipeline result=%#v error=%v calls=%d closed=%d; want nil result, %v, one call and one close", result, err, calls, body.closed, terminal)
					}
				})
			}
		})
	}
}
