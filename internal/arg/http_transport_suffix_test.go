package arg

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type argSuffixBody struct {
	io.Reader
	closes int
}

func (b *argSuffixBody) Close() error {
	b.closes++
	return nil
}

type argSuffixReadError struct{ err error }

func (r argSuffixReadError) Read([]byte) (int, error) { return 0, r.err }

func suffixQuery(body *argSuffixBody) (*Result, error) {
	poster := &fakePoster{resp: &http.Response{StatusCode: http.StatusOK, Body: body}}
	return NewClient(NewHTTPTransportWithClient(poster, "https://example.test/graph")).Query(
		context.Background(), "resources", map[string]string{"sub": "name"},
	)
}

func TestHTTPTransportRejectsTrailingData(t *testing.T) {
	prefix := `{"data":[{"id":"one"}],"resultTruncated":"false"}`
	for _, suffix := range []string{
		`{}`, `[]`, `0`, `null`, `true`, "garbage", "//comment", "\xef\xbb\xbf",
		strings.Repeat(" ", 8193) + `{"data":[]}`,
	} {
		for _, split := range []bool{false, true} {
			t.Run(suffix[:min(len(suffix), 20)]+map[bool]string{false: "/buffered", true: "/split"}[split], func(t *testing.T) {
				var reader io.Reader = strings.NewReader(prefix + suffix)
				if split {
					reader = io.MultiReader(strings.NewReader(prefix), strings.NewReader(suffix))
				}
				body := &argSuffixBody{Reader: reader}
				result, err := suffixQuery(body)
				if err == nil || result != nil {
					t.Fatalf("result=%#v error=%v; want failed query, not successful prefix", result, err)
				}
				if body.closes != 1 {
					t.Fatalf("body closes=%d; want one", body.closes)
				}
			})
		}
	}
}

func TestHTTPTransportAcceptsOnlyWhitespaceAfterDocument(t *testing.T) {
	for _, suffix := range []string{"", " \t\r\n", strings.Repeat(" \t\r\n", 4097)} {
		body := &argSuffixBody{Reader: io.MultiReader(
			strings.NewReader(`{"data":[{"id":"one"}],"resultTruncated":"false","future":{"x":1}}`),
			strings.NewReader(suffix),
		)}
		result, err := suffixQuery(body)
		if err != nil || result == nil || len(result.Data) != 1 || string(result.Data[0]) != `{"id":"one"}` {
			t.Fatalf("result=%#v error=%v; want complete original row", result, err)
		}
		if body.closes != 1 {
			t.Fatalf("body closes=%d; want one", body.closes)
		}
	}
}

func TestHTTPTransportPreservesTrailingReadError(t *testing.T) {
	for _, cause := range []error{errors.New("controlled body failure"), context.Canceled, context.DeadlineExceeded} {
		body := &argSuffixBody{Reader: io.MultiReader(
			strings.NewReader(`{"data":[{"id":"one"}]}`), argSuffixReadError{err: cause},
		)}
		result, err := suffixQuery(body)
		if result != nil || !errors.Is(err, cause) {
			t.Fatalf("result=%#v error=%v; want nil and wrapped cause %v", result, err, cause)
		}
		if body.closes != 1 {
			t.Fatalf("body closes=%d; want one", body.closes)
		}
	}
}
