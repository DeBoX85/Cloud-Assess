package arg

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

const (
	MaxSubscriptionsPerRequest = 300
	MaxRowsPerPage             = int32(5000)
	ResultFormatObjectArray    = "objectArray"
	ManagementGroupScopeFilter = "AtScopeAndAbove"
)

type QueryOptions struct {
	ManagementGroupScope bool
}

type RequestOptions struct {
	ResultFormat             string  `json:"resultFormat,omitempty"`
	Top                      *int32  `json:"$top,omitempty"`
	SkipToken                *string `json:"$skipToken,omitempty"`
	AuthorizationScopeFilter *string `json:"authorizationScopeFilter,omitempty"`
}

type Request struct {
	Subscriptions []string        `json:"subscriptions"`
	Query         string          `json:"query"`
	Options       *RequestOptions `json:"options"`
}

type Response struct {
	Data       []json.RawMessage `json:"data"`
	SkipToken  *string           `json:"$skipToken,omitempty"`
	Quota      int               `json:"-"`
	RetryAfter time.Duration     `json:"-"`
}

type Result struct {
	Data []json.RawMessage
}

// Transport performs one Resource Graph request. Authentication, retry and HTTP details
// belong in the transport implementation rather than the batching/pagination client.
type Transport interface {
	Do(context.Context, Request) (*Response, error)
}

type Client struct {
	transport Transport
}

func NewClient(transport Transport) *Client {
	return &Client{transport: transport}
}

// Query executes one KQL query over all supplied subscriptions, preserving the reference
// batching and pagination semantics. Subscription IDs are sorted before batching to make
// request composition deterministic.
func (c *Client) Query(
	ctx context.Context,
	query string,
	subscriptions map[string]string,
	opts ...QueryOptions,
) (*Result, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("ARG transport is not configured")
	}

	result := &Result{Data: make([]json.RawMessage, 0, MaxRowsPerPage)}
	subscriptionIDs := make([]string, 0, len(subscriptions))
	for subscriptionID := range subscriptions {
		subscriptionIDs = append(subscriptionIDs, subscriptionID)
	}
	sort.Strings(subscriptionIDs)

	managementGroupScope := len(opts) > 0 && opts[0].ManagementGroupScope
	for start := 0; start < len(subscriptionIDs); start += MaxSubscriptionsPerRequest {
		end := min(start+MaxSubscriptionsPerRequest, len(subscriptionIDs))
		batch := append([]string(nil), subscriptionIDs[start:end]...)

		var skipToken *string
		for {
			top := MaxRowsPerPage
			requestOptions := &RequestOptions{
				ResultFormat: ResultFormatObjectArray,
				Top:          &top,
				SkipToken:    skipToken,
			}
			if managementGroupScope {
				value := ManagementGroupScopeFilter
				requestOptions.AuthorizationScopeFilter = &value
			}

			response, err := c.transport.Do(ctx, Request{
				Subscriptions: batch,
				Query:         query,
				Options:       requestOptions,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run resource graph query: %w", err)
			}
			if response == nil {
				return nil, fmt.Errorf("failed to run resource graph query: transport returned nil response")
			}

			result.Data = append(result.Data, response.Data...)
			if response.SkipToken == nil {
				break
			}
			token := *response.SkipToken
			skipToken = &token
		}
	}

	return result, nil
}
