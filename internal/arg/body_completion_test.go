package arg

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
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
