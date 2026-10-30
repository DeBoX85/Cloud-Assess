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

func TestExecuteRecommendationsCancellationStopsRemainingDefinitions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	id := "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Test/widgets/one"
	q := &fakeQuerier{results: map[string]*Result{"query": findingResult(id, "one")}}
	definitions := []assessment.RecommendationDefinition{{ID: "one", Query: "query"}, {ID: "two", Query: "query"}}
	filtered := 0
	findings, _, err := ExecuteRecommendations(ctx, q, definitions, nil, func(string) bool {
		filtered++
		cancel()
		return false
	}, 1)
	if filtered != 1 || len(q.calls) != 1 || !errors.Is(err, context.Canceled) || findings != nil {
		t.Fatalf("remaining definitions started after cancellation: filtered=%d calls=%v findings=%v error=%v", filtered, q.calls, findings, err)
	}
}

type terminalExecutorQuerier func(context.Context) (*Result, error)

func (q terminalExecutorQuerier) Query(ctx context.Context, _ string, _ map[string]string, _ ...QueryOptions) (*Result, error) {
	return q(ctx)
}

func TestExecuteRecommendationsParentContextWinsOverPlainQueryError(t *testing.T) {
	for _, mode := range []string{"canceled", "deadline", "ordinary"} {
		t.Run(mode, func(t *testing.T) {
			plain := errors.New("plain transport failure")
			ctx, cancel := context.WithCancel(context.Background())
			if mode == "deadline" {
				cancel()
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			}
			defer cancel()
			calls := 0
			q := terminalExecutorQuerier(func(queryContext context.Context) (*Result, error) {
				calls++
				if mode == "canceled" {
					cancel()
				}
				if mode != "ordinary" {
					<-queryContext.Done()
				}
				return nil, plain
			})
			definitions := []assessment.RecommendationDefinition{{ID: "one", Query: "query"}}
			findings, _, err := ExecuteRecommendations(ctx, q, definitions, nil, nil, 1)
			want := error(plain)
			if mode == "canceled" {
				want = context.Canceled
			}
			if mode == "deadline" {
				want = context.DeadlineExceeded
			}
			if calls != 1 || findings != nil || !errors.Is(err, want) {
				t.Fatalf("parent terminal priority lost: mode=%s calls=%d findings=%v error=%v want=%v", mode, calls, findings, err, want)
			}
			if mode == "ordinary" && errors.Is(err, context.Canceled) {
				t.Fatal("executor sibling cancellation replaced ordinary initiating error")
			}
		})
	}
}
