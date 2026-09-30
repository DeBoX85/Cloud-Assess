package arg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type fakeTransport struct {
	requests  []Request
	responses []*Response
	errAt     int
}

func (f *fakeTransport) Do(_ context.Context, request Request) (*Response, error) {
	f.requests = append(f.requests, request)
	index := len(f.requests) - 1
	if f.errAt > 0 && len(f.requests) == f.errAt {
		return nil, errors.New("transport failure")
	}
	if index >= len(f.responses) {
		return &Response{}, nil
	}
	return f.responses[index], nil
}

func strptr(value string) *string { return &value }

func TestQueryBatchesAtThreeHundredSubscriptions(t *testing.T) {
	transport := &fakeTransport{responses: []*Response{{}, {}}}
	client := NewClient(transport)
	subscriptions := make(map[string]string, 301)
	for i := 0; i < 301; i++ {
		id := fmt.Sprintf("sub-%03d", i)
		subscriptions[id] = id
	}

	if _, err := client.Query(context.Background(), "resources", subscriptions); err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got, want := len(transport.requests), 2; got != want {
		t.Fatalf("request count = %d, want %d", got, want)
	}
	if len(transport.requests[0].Subscriptions) != 300 || len(transport.requests[1].Subscriptions) != 1 {
		t.Fatalf("unexpected batch sizes: %d and %d", len(transport.requests[0].Subscriptions), len(transport.requests[1].Subscriptions))
	}
	if transport.requests[0].Subscriptions[0] != "sub-000" || transport.requests[1].Subscriptions[0] != "sub-300" {
		t.Fatalf("subscription batching is not deterministic: %#v %#v", transport.requests[0].Subscriptions[:1], transport.requests[1].Subscriptions)
	}
}

func TestQueryFollowsSkipTokenAndAggregatesRows(t *testing.T) {
	transport := &fakeTransport{responses: []*Response{
		{Data: []json.RawMessage{json.RawMessage(`{"page":1}`)}, SkipToken: strptr("next")},
		{Data: []json.RawMessage{json.RawMessage(`{"page":2}`)}},
	}}
	client := NewClient(transport)
	result, err := client.Query(context.Background(), "resources", map[string]string{"sub": "name"})
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(result.Data) != 2 || len(transport.requests) != 2 {
		t.Fatalf("unexpected pagination result: rows=%d requests=%d", len(result.Data), len(transport.requests))
	}
	if transport.requests[0].Options.SkipToken != nil {
		t.Fatal("first request should not contain a skip token")
	}
	if transport.requests[1].Options.SkipToken == nil || *transport.requests[1].Options.SkipToken != "next" {
		t.Fatalf("second request skip token = %#v", transport.requests[1].Options.SkipToken)
	}
}

func TestQueryRejectsRepeatedOrEmptyContinuationToken(t *testing.T) {
	for _, test := range []struct {
		name      string
		responses []*Response
		calls     int
	}{
		{name: "repeated", responses: []*Response{{Data: []json.RawMessage{json.RawMessage(`{"page":1}`)}, SkipToken: strptr("same")}, {Data: []json.RawMessage{json.RawMessage(`{"page":1}`)}, SkipToken: strptr("same")}}, calls: 2},
		{name: "empty", responses: []*Response{{SkipToken: strptr("")}}, calls: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			transport := &fakeTransport{responses: test.responses}
			result, err := NewClient(transport).Query(context.Background(), "resources", map[string]string{"sub": "name"})
			if err == nil || result != nil || len(transport.requests) != test.calls {
				t.Fatalf("result = %#v, error = %v, calls = %d; want explicit failure after %d calls", result, err, len(transport.requests), test.calls)
			}
		})
	}
}

func TestQueryStopsBeforeTransportWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	transport := &fakeTransport{}
	result, err := NewClient(transport).Query(ctx, "resources", map[string]string{"sub": "name"})
	if result != nil || !errors.Is(err, context.Canceled) || len(transport.requests) != 0 {
		t.Fatalf("result = %#v, error = %v, calls = %d; want cancellation before request", result, err, len(transport.requests))
	}
}

func TestQueryUsesReferenceRequestOptions(t *testing.T) {
	transport := &fakeTransport{responses: []*Response{{}}}
	client := NewClient(transport)
	_, err := client.Query(context.Background(), "policyresources", map[string]string{"sub": "name"}, QueryOptions{ManagementGroupScope: true})
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	request := transport.requests[0]
	if request.Options.ResultFormat != ResultFormatObjectArray || request.Options.Top == nil || *request.Options.Top != 5000 {
		t.Fatalf("unexpected options: %+v", request.Options)
	}
	if request.Options.AuthorizationScopeFilter == nil || *request.Options.AuthorizationScopeFilter != "AtScopeAndAbove" {
		t.Fatalf("management group scope filter missing: %+v", request.Options)
	}
}

func TestQueryDoesNotUseManagementGroupScopeByDefault(t *testing.T) {
	transport := &fakeTransport{responses: []*Response{{}}}
	client := NewClient(transport)
	_, err := client.Query(context.Background(), "resources", map[string]string{"sub": "name"})
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if transport.requests[0].Options.AuthorizationScopeFilter != nil {
		t.Fatalf("unexpected authorization scope filter: %+v", transport.requests[0].Options)
	}
}

func TestQueryPropagatesTransportFailure(t *testing.T) {
	transport := &fakeTransport{responses: []*Response{{}}, errAt: 1}
	client := NewClient(transport)
	_, err := client.Query(context.Background(), "resources", map[string]string{"sub": "name"})
	if err == nil || !errors.Is(err, errors.New("transport failure")) && err.Error() != "failed to run resource graph query: transport failure" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestQueryWithNoSubscriptionsMakesNoRequests(t *testing.T) {
	transport := &fakeTransport{}
	client := NewClient(transport)
	result, err := client.Query(context.Background(), "resources", nil)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(transport.requests) != 0 || !reflect.DeepEqual(result.Data, []json.RawMessage{}) {
		t.Fatalf("unexpected empty-scope behavior: requests=%d data=%#v", len(transport.requests), result.Data)
	}
}

func TestQueryLargeSubscriptionSetMakesBoundedDisjointBatches(t *testing.T) {
	const count = 3001
	subscriptions := make(map[string]string, count)
	for i := 0; i < count; i++ {
		subscriptions[fmt.Sprintf("sub-%04d", i)] = "fixture"
	}
	transport := &fakeTransport{}
	result, err := NewClient(transport).Query(context.Background(), "resources", subscriptions)
	if err != nil || result == nil || len(transport.requests) != 11 {
		t.Fatalf("result=%v error=%v request count=%d", result, err, len(transport.requests))
	}
	seen := map[string]bool{}
	for _, request := range transport.requests {
		if len(request.Subscriptions) > 300 || len(request.Subscriptions) == 0 {
			t.Fatalf("unbounded or empty batch: %d", len(request.Subscriptions))
		}
		for _, id := range request.Subscriptions {
			if seen[id] {
				t.Fatalf("subscription queried twice: %s", id)
			}
			seen[id] = true
		}
	}
	if len(seen) != count {
		t.Fatalf("queried %d subscriptions, want %d", len(seen), count)
	}
}

func TestQueryFailureAfterSuccessfulPageDoesNotPublishTruncatedSuccess(t *testing.T) {
	transport := &fakeTransport{responses: []*Response{{Data: []json.RawMessage{json.RawMessage(`{"id":"healthy-page"}`)}, SkipToken: strptr("next")}}, errAt: 2}
	result, err := NewClient(transport).Query(context.Background(), "resources", map[string]string{"sub": "fixture"})
	if err == nil || result != nil || len(transport.requests) != 2 {
		t.Fatalf("result=%v error=%v calls=%d; truncated query must fail", result, err, len(transport.requests))
	}
}
