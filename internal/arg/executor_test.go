package arg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type fakeQuerier struct {
	mu      sync.Mutex
	calls   []string
	results map[string]*Result
	errors  map[string]error
}

func (f *fakeQuerier) Query(_ context.Context, query string, _ map[string]string, _ ...QueryOptions) (*Result, error) {
	f.mu.Lock()
	f.calls = append(f.calls, query)
	f.mu.Unlock()
	if err := f.errors[query]; err != nil {
		return nil, err
	}
	if result := f.results[query]; result != nil {
		return result, nil
	}
	return &Result{}, nil
}

func findingResult(id, name string) *Result {
	row := fmt.Sprintf(`{"id":%q,"name":%q}`, id, name)
	return &Result{Data: []json.RawMessage{json.RawMessage(row)}}
}

func TestExecuteRecommendationsReturnsDeterministicFindings(t *testing.T) {
	q := &fakeQuerier{results: map[string]*Result{
		"query-b": findingResult("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/b", "b"),
		"query-a": findingResult("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/a", "a"),
	}}
	definitions := []assessment.RecommendationDefinition{
		{ID: "rec-b", Query: "query-b", ResourceType: "Microsoft.Test/widgets"},
		{ID: "rec-a", Query: "query-a", ResourceType: "Microsoft.Test/widgets"},
	}
	findings, warnings, err := ExecuteRecommendations(context.Background(), q, definitions, map[string]string{"sub": "Test"}, nil, 2)
	if err != nil {
		t.Fatalf("ExecuteRecommendations() error = %v", err)
	}
	if len(warnings) != 0 || len(findings) != 2 {
		t.Fatalf("unexpected result: findings=%d warnings=%d", len(findings), len(warnings))
	}
	if findings[0].RecommendationID != "rec-a" || findings[1].RecommendationID != "rec-b" {
		t.Fatalf("findings are not deterministically ordered: %#v", findings)
	}
}

func TestExecuteRecommendationsSkipsUnsupportedLogicalTableWithWarning(t *testing.T) {
	q := &fakeQuerier{errors: map[string]error{
		"bad-table": &azcore.ResponseError{ErrorCode: "DisallowedLogicalTableName", StatusCode: 400},
	}}
	definitions := []assessment.RecommendationDefinition{{ID: "rec", Query: "bad-table", ResourceType: "Microsoft.Test/widgets"}}
	findings, warnings, err := ExecuteRecommendations(context.Background(), q, definitions, nil, nil, 1)
	if err != nil {
		t.Fatalf("unsupported logical table should not fail execution: %v", err)
	}
	if len(findings) != 0 || len(warnings) != 1 || warnings[0].Code != "unsupported_logical_table" {
		t.Fatalf("unexpected skip result: findings=%#v warnings=%#v", findings, warnings)
	}
}

func TestExecuteRecommendationsReturnsOrdinaryQueryFailure(t *testing.T) {
	boom := errors.New("query failed")
	q := &fakeQuerier{errors: map[string]error{"query": boom}}
	definitions := []assessment.RecommendationDefinition{{ID: "rec", Query: "query"}}
	findings, _, err := ExecuteRecommendations(context.Background(), q, definitions, nil, nil, 1)
	if err == nil || !errors.Is(err, boom) {
		t.Fatalf("error = %v, want wrapped query failure", err)
	}
	if findings != nil {
		t.Fatalf("fatal query failure should return nil findings, got %#v", findings)
	}
}

func TestExecuteRecommendationsAppliesResourceFilter(t *testing.T) {
	keep := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/keep"
	drop := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/drop"
	q := &fakeQuerier{results: map[string]*Result{
		"query": {Data: []json.RawMessage{
			json.RawMessage(fmt.Sprintf(`{"id":%q,"name":"keep"}`, keep)),
			json.RawMessage(fmt.Sprintf(`{"id":%q,"name":"drop"}`, drop)),
		}},
	}}
	definitions := []assessment.RecommendationDefinition{{ID: "rec", Query: "query", ResourceType: "Microsoft.Test/widgets"}}
	findings, _, err := ExecuteRecommendations(context.Background(), q, definitions, nil, func(id string) bool { return id == drop }, 1)
	if err != nil {
		t.Fatalf("ExecuteRecommendations() error = %v", err)
	}
	if len(findings) != 1 || findings[0].ResourceID != keep {
		t.Fatalf("resource filter not applied: %#v", findings)
	}
}

func TestExecuteRecommendationsSkipsEmptyQuery(t *testing.T) {
	q := &fakeQuerier{}
	definitions := []assessment.RecommendationDefinition{{ID: "rec", Query: ""}}
	findings, warnings, err := ExecuteRecommendations(context.Background(), q, definitions, nil, nil, 1)
	if err != nil || len(findings) != 0 || len(warnings) != 0 {
		t.Fatalf("unexpected empty-query result: findings=%#v warnings=%#v err=%v", findings, warnings, err)
	}
	if len(q.calls) != 0 {
		t.Fatalf("empty query called querier: %#v", q.calls)
	}
}

func TestExecuteRecommendationsEmptyDefinitions(t *testing.T) {
	q := &fakeQuerier{}
	findings, warnings, err := ExecuteRecommendations(context.Background(), q, nil, nil, nil, 10)
	if err != nil || findings == nil || warnings == nil || len(findings) != 0 || len(warnings) != 0 {
		t.Fatalf("unexpected empty-definition result: findings=%#v warnings=%#v err=%v", findings, warnings, err)
	}
}
