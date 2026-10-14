package arg

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestExecuteRecommendationsRejectsTerminalContextBeforeWork(t *testing.T) {
	for _, state := range []string{"canceled", "deadline"} {
		for _, empty := range []bool{false, true} {
			t.Run(state+"/"+map[bool]string{false: "definitions", true: "empty"}[empty], func(t *testing.T) {
				var ctx context.Context
				var cancel context.CancelFunc
				want := error(context.Canceled)
				if state == "deadline" {
					ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
					want = context.DeadlineExceeded
				} else {
					ctx, cancel = context.WithCancel(context.Background())
					cancel()
				}
				defer cancel()
				q := &fakeQuerier{}
				definitions := []assessment.RecommendationDefinition{{ID: "one", Query: "query"}}
				if empty {
					definitions = nil
				}
				findings, _, err := ExecuteRecommendations(ctx, q, definitions, nil, nil, 1)
				if !errors.Is(err, want) || findings != nil || len(q.calls) != 0 {
					t.Fatalf("terminal context represented as successful work: findings=%v error=%v calls=%v", findings, err, q.calls)
				}
			})
		}
	}
}

func TestExecuteRecommendationsCancellationDuringFilteringIsNotSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	q := &fakeQuerier{results: map[string]*Result{"query": findingResult(id, "one")}}
	definitions := []assessment.RecommendationDefinition{{ID: "one", Query: "query"}}
	filtered := 0
	findings, _, err := ExecuteRecommendations(ctx, q, definitions, nil, func(string) bool {
		filtered++
		cancel()
		return false
	}, 1)
	if filtered != 1 || !errors.Is(err, context.Canceled) || findings != nil {
		t.Fatalf("cancellation after query represented as success: filtered=%d findings=%v error=%v", filtered, findings, err)
	}
}
