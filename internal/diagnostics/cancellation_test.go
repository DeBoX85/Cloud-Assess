package diagnostics

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type terminalDiagnosticsClient struct {
	cancel context.CancelFunc
	plain bool
}

func (p terminalDiagnosticsClient) PostStream(_ context.Context, _ string, body io.ReadSeekCloser) (*http.Response, error) {
	defer body.Close()
	p.cancel()
	if p.plain {
		return nil, errors.New("plain diagnostics transport failure")
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"responses":[{"httpStatusCode":200,"content":{"value":[]}}]}`))}, nil
}

func TestDiagnosticsCallerCancellationCannotCreateMissingSettingsFindings(t *testing.T) {
	for _, plain := range []bool{false, true} {
		t.Run(map[bool]string{false: "valid response", true: "plain error"}[plain], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			resource := assessment.Resource{ID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/one", Type: "Microsoft.Storage/storageAccounts"}
			result, err := NewWithClient(terminalDiagnosticsClient{cancel: cancel, plain: plain}, "https://management.azure.com").Scan(ctx, []assessment.Resource{resource}, nil, nil)
			if !errors.Is(err, context.Canceled) || len(result.Findings) != 0 {
				t.Fatalf("Diagnostics caller cancellation became assessment data: findings=%+v error=%v", result.Findings, err)
			}
		})
	}
}
