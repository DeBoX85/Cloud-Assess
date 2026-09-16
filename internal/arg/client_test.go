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
