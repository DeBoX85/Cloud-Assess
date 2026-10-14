package cost

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

type terminalCostPoster struct {
	cancel context.CancelFunc
	plain bool
}

func (p terminalCostPoster) Post(_ context.Context, _ string, body io.ReadSeekCloser) ([]byte, *http.Response, error) {
	defer body.Close()
	p.cancel()
	if p.plain {
		return nil, nil, errors.New("plain cost transport failure")
	}
	return []byte(`{"properties":{"rows":[[1,"Storage","USD"]]}}`), &http.Response{StatusCode: http.StatusOK}, nil
}

func TestCostCallerCancellationCannotBecomeSuccessfulDataset(t *testing.T) {
	for _, plain := range []bool{false, true} {
		t.Run(map[bool]string{false: "valid response", true: "plain error"}[plain], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			result, err := NewWithClient(terminalCostPoster{cancel: cancel, plain: plain}, "https://management.azure.com").scanAt(ctx, map[string]string{"sub": "Synthetic"}, time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC))
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Cost caller cancellation lost: result=%+v error=%v", result, err)
			}
		})
	}
}
