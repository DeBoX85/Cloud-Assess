package arg

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const DefaultRuleWorkers = 10

type Querier interface {
	Query(context.Context, string, map[string]string, ...QueryOptions) (*Result, error)
}

type RuleWarning struct {
	RecommendationID string `json:"recommendationId"`
	ResourceType     string `json:"resourceType"`
	Code             string `json:"code"`
	Message          string `json:"message"`
}

type ruleOutcome struct {
	findings []assessment.Finding
	warning  *RuleWarning
	err      error
}

// ExecuteRecommendations evaluates recommendation queries with bounded concurrency.
// Definitions are expected to have already passed catalog/filter applicability checks.
func ExecuteRecommendations(
	ctx context.Context,
	querier Querier,
	definitions []assessment.RecommendationDefinition,
	subscriptions map[string]string,
	isResourceExcluded func(string) bool,
	workers int,
) ([]assessment.Finding, []RuleWarning, error) {
	if querier == nil {
		return nil, nil, fmt.Errorf("ARG querier is not configured")
	}
	if workers <= 0 {
		workers = DefaultRuleWorkers
	}
	if workers > len(definitions) && len(definitions) > 0 {
		workers = len(definitions)
	}
	if len(definitions) == 0 {
		return []assessment.Finding{}, []RuleWarning{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan assessment.RecommendationDefinition)
	outcomes := make(chan ruleOutcome, len(definitions))
	var wait sync.WaitGroup

	worker := func() {
		defer wait.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case definition, ok := <-jobs:
				if !ok {
					return
				}
				if definition.Query == "" {
					outcomes <- ruleOutcome{}
					continue
				}

				result, err := querier.Query(ctx, definition.Query, subscriptions)
				if err != nil {
					if IsUnsupportedLogicalTableError(err) {
						outcomes <- ruleOutcome{warning: &RuleWarning{
							RecommendationID: definition.ID,
							ResourceType:     definition.ResourceType,
							Code:             "unsupported_logical_table",
							Message:          err.Error(),
						}}
						continue
					}
					outcomes <- ruleOutcome{err: fmt.Errorf("recommendation %s query failed: %w", definition.ID, err)}
					cancel()
					return
				}
				if result == nil {
					outcomes <- ruleOutcome{err: fmt.Errorf("recommendation %s query failed: nil ARG result", definition.ID)}
					cancel()
					return
				}

				rows, malformed := DecodeRowsWithStats[FindingRow](result.Data)
				findings := FindingsFromRows(definition, rows, subscriptions)
				if isResourceExcluded != nil {
					filtered := findings[:0]
					for _, finding := range findings {
						if !isResourceExcluded(finding.ResourceID) {
							filtered = append(filtered, finding)
						}
					}
					findings = filtered
				}

				outcome := ruleOutcome{findings: findings}
				if malformed > 0 {
					outcome.warning = &RuleWarning{
						RecommendationID: definition.ID,
						ResourceType:     definition.ResourceType,
						Code:             "malformed_arg_rows",
						Message:          fmt.Sprintf("skipped %d/%d malformed ARG row(s)", malformed, len(result.Data)),
					}
				}
				outcomes <- outcome
			}
		}
	}

	wait.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		defer close(jobs)
		for _, definition := range definitions {
			select {
			case <-ctx.Done():
				return
			case jobs <- definition:
			}
		}
	}()

	go func() {
		wait.Wait()
		close(outcomes)
	}()

	var findings []assessment.Finding
	var warnings []RuleWarning
	var firstErr error
	for outcome := range outcomes {
		if outcome.err != nil && firstErr == nil {
			firstErr = outcome.err
		}
		if outcome.warning != nil {
			warnings = append(warnings, *outcome.warning)
		}
		findings = append(findings, outcome.findings...)
	}
	if firstErr != nil {
		return nil, warnings, firstErr
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].RecommendationID != findings[j].RecommendationID {
			return findings[i].RecommendationID < findings[j].RecommendationID
		}
		return findings[i].ResourceID < findings[j].ResourceID
	})
	sort.Slice(warnings, func(i, j int) bool {
		if warnings[i].RecommendationID != warnings[j].RecommendationID {
			return warnings[i].RecommendationID < warnings[j].RecommendationID
		}
		return warnings[i].ResourceType < warnings[j].ResourceType
	})
	return findings, warnings, nil
}
