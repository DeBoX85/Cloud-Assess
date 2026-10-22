// Region meter-history semantics derived from MIT-licensed Azure Quick Review.
// See NOTICE.md and docs/REGION_METER_HISTORY_CONTRACT_20261019.md.
package region

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type MeterHistoryRequest struct {
	SubscriptionID string
	Months         int
	Now            time.Time
	Currency       string
}

func historyRequestValid(q MeterHistoryRequest) bool {
	if len(q.SubscriptionID) != 36 || !subscriptionID.MatchString(q.SubscriptionID) || q.Months < 1 || q.Months > 12 || q.Now.IsZero() || q.Now.UTC().Year() < 2 || q.Now.UTC().Year() > 9999 || len(q.Currency) != 3 {
		return false
	}
	for i := range q.Currency {
		if q.Currency[i] < 'A' || q.Currency[i] > 'Z' {
			return false
		}
	}
	return true
}

func historyRequestError(ctx context.Context) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	return fmt.Errorf("region meter history rejected [history_request_invalid]")
}

// MeterHistoryQuery owns the read-oriented query bytes and uses an explicit
// clock, preserving the source's second-precision completed UTC months.
func MeterHistoryQuery(ctx context.Context, q MeterHistoryRequest) ([]byte, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	if !historyRequestValid(q) {
		return nil, historyRequestError(ctx)
	}
	now := q.Now.UTC()
	first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	groups := []map[string]string{}
	for _, name := range []string{"ResourceId", "ResourceGuid", "MeterCategory", "MeterSubcategory", "Meter"} {
		groups = append(groups, map[string]string{"type": "Dimension", "name": name})
	}
	body, e := json.Marshal(map[string]any{
		"type": "ActualCost", "timeframe": "Custom",
		"timePeriod": map[string]string{"from": first.AddDate(0, -q.Months, 0).Format(time.RFC3339), "to": first.Add(-time.Second).Format(time.RFC3339)},
		"dataset": map[string]any{
			"granularity": "Monthly", "grouping": groups,
			"aggregation": map[string]any{"PreTaxCost": map[string]string{"name": "PreTaxCost", "function": "Sum"}},
		},
	})
	if e != nil {
		return nil, historyRequestError(ctx)
	}
	if e = ctx.Err(); e != nil {
		return nil, e
	}
	return body, nil
}
