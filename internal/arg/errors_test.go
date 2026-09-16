package arg

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestUnsupportedLogicalTableGenericMessages(t *testing.T) {
	cases := []struct {
		message string
		want    bool
	}{
		{message: "DisallowedLogicalTableName", want: true},
		{message: "Logical table Foo is invalid, unsupported or disallowed", want: true},
		{message: "Table Foo is invalid, unsupported or disallowed", want: false},
		{message: "permission denied", want: false},
	}
	for _, tt := range cases {
		if got := IsUnsupportedLogicalTableError(errors.New(tt.message)); got != tt.want {
			t.Fatalf("IsUnsupportedLogicalTableError(%q) = %v, want %v", tt.message, got, tt.want)
		}
	}
	if IsUnsupportedLogicalTableError(nil) {
		t.Fatal("nil error should not be classified as unsupported logical table")
	}
}

func TestUnsupportedLogicalTableResponseErrorCode(t *testing.T) {
	err := &azcore.ResponseError{ErrorCode: "DisallowedLogicalTableName", StatusCode: 400}
	if !IsUnsupportedLogicalTableError(err) {
		t.Fatal("DisallowedLogicalTableName response error should be skipped")
	}
}

func TestUnsupportedLogicalTableNestedDetailAndBodyRestoration(t *testing.T) {
	u, _ := url.Parse("https://example.test/graph")
	body := `{"error":{"code":"BadRequest","message":"query failed","details":[{"code":"DisallowedLogicalTableName","message":"not available"}]}}`
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Request:    &http.Request{Method: http.MethodPost, URL: u},
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	err := &azcore.ResponseError{ErrorCode: "BadRequest", StatusCode: 400, RawResponse: response}
	if !IsUnsupportedLogicalTableError(err) {
		t.Fatal("nested DisallowedLogicalTableName should be skipped")
	}
	restored, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		t.Fatalf("read restored body: %v", readErr)
	}
	if string(restored) != body {
		t.Fatalf("response body was not restored: %q", string(restored))
	}
}
