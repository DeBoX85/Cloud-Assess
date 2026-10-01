package arg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type finalPageCancelTransport struct {
	cancel context.CancelFunc
	calls  int
}

func (f *finalPageCancelTransport) Do(context.Context, Request) (*Response, error) {
	f.calls++
	f.cancel()
	return &Response{Data: []json.RawMessage{json.RawMessage(`{"id":"fixture"}`)}}, nil
}
func TestQueryCanceledDuringFinalPageDoesNotPublishSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	transport := &finalPageCancelTransport{cancel: cancel}
	result, err := NewClient(transport).Query(ctx, "resources", map[string]string{"sub": "fixture"})
	if !errors.Is(err, context.Canceled) || result != nil || transport.calls != 1 {
		t.Fatalf("result=%v error=%v calls=%d", result, err, transport.calls)
	}
}

func TestHealthyPaginationTokenStateIsPerBatchAndQuery(t *testing.T) {
	subscriptions := map[string]string{}
	for i := 0; i < 301; i++ {
		subscriptions[fmt.Sprintf("sub-%03d", i)] = "fixture"
	}
	transport := &fakeTransport{}
	// Both batches and both invocations legitimately use the same opaque token.
	for range 4 {
		transport.responses = append(transport.responses,
			&Response{Data: []json.RawMessage{json.RawMessage(`{"page":1}`)}, SkipToken: strptr("same-token")},
			&Response{Data: []json.RawMessage{json.RawMessage(`{"page":2}`)}})
	}
	client := NewClient(transport)
	for run := 0; run < 2; run++ {
		result, err := client.Query(context.Background(), "resources", subscriptions)
		if err != nil || result == nil || len(result.Data) != 4 {
			t.Fatalf("run=%d result=%v error=%v", run, result, err)
		}
	}
	if len(transport.requests) != 8 {
		t.Fatalf("requests=%d", len(transport.requests))
	}
	for i, request := range transport.requests {
		if i%2 == 0 && request.Options.SkipToken != nil {
			t.Fatalf("batch inherited skip token at request %d", i)
		}
		if i%2 == 1 && (request.Options.SkipToken == nil || *request.Options.SkipToken != "same-token") {
			t.Fatalf("missing continuation at request %d", i)
		}
	}
}

type distinctCancelTransport struct {
	cancel context.CancelFunc
	calls  int
}

func (f *distinctCancelTransport) Do(context.Context, Request) (*Response, error) {
	f.calls++
	if f.calls > 6 {
		return nil, fmt.Errorf("fixture runaway tripwire")
	}
	if f.calls == 3 {
		f.cancel()
	}
	token := fmt.Sprintf("opaque-%d", f.calls)
	return &Response{Data: []json.RawMessage{json.RawMessage(`{"id":"fixture"}`)}, SkipToken: &token}, nil
}
func TestDistinctQueryTokensHonorCancellationWithoutPartialResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	transport := &distinctCancelTransport{cancel: cancel}
	result, err := NewClient(transport).Query(ctx, "resources", map[string]string{"sub": "fixture"})
	if !errors.Is(err, context.Canceled) || result != nil || transport.calls != 3 {
		t.Fatalf("result=%v error=%v calls=%d", result, err, transport.calls)
	}
}
