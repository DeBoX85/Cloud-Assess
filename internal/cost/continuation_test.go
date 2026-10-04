package cost

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCostRejectsUnconsumedContinuation(t *testing.T) {
	for _, rows := range []string{`[]`, `[[12.34,"Storage","USD"]]`} {
		t.Run(rows, func(t *testing.T) {
			poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
				return []byte(`{"properties":{"columns":[{"name":"Cost","type":"Number"},{"name":"ServiceName","type":"String"},{"name":"Currency","type":"String"}],"nextLink":"https://management.azure.com/subscriptions/sub/providers/Microsoft.CostManagement/query?api-version=2021-10-01&$skiptoken=next","rows":` + rows + `}}`), &http.Response{StatusCode: http.StatusOK}, nil
			}}
			from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
			records, warnings, err := NewWithClient(poster, "https://management.azure.com").querySubscription(
				context.Background(), "sub", "Synthetic", from, from.AddDate(0, 1, 0),
			)
			if err == nil || !strings.Contains(err.Error(), "continuation") || len(records) != 0 || len(warnings) != 0 {
				t.Fatalf("unconsumed continuation returned records=%#v warnings=%#v error=%v; want explicit continuation failure and no accepted costs", records, warnings, err)
			}
			if len(poster.calls) != 1 {
				t.Fatalf("posted %d requests; fail-closed guard must not follow provider URLs", len(poster.calls))
			}
		})
	}
}

func TestCostAcceptsCompletePageWithNullOrEmptyContinuation(t *testing.T) {
	for _, link := range []string{`null`, `""`} {
		t.Run(link, func(t *testing.T) {
			poster := &fakePoster{handler: func(string, []byte) ([]byte, *http.Response, error) {
				return []byte(`{"properties":{"nextLink":` + link + `,"rows":[[12.34,"Storage","USD"]]}}`), &http.Response{StatusCode: http.StatusOK}, nil
			}}
			from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
			records, warnings, err := NewWithClient(poster, "https://management.azure.com").querySubscription(
				context.Background(), "sub", "Synthetic", from, from.AddDate(0, 1, 0),
			)
			if err != nil || len(records) != 1 || records[0].Value != "12.34" || records[0].ServiceName != "Storage" || records[0].Currency != "USD" || len(warnings) != 0 || len(poster.calls) != 1 {
				t.Fatalf("complete page records=%#v warnings=%#v error=%v calls=%d", records, warnings, err, len(poster.calls))
			}
		})
	}
}
