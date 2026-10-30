package region

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestSubscriptionReservationsFrozenBudgets(t *testing.T) {
	for _, tc := range bulkContract(t).Recipes {
		t.Run(tc.ID, func(t *testing.T) {
			calls := 0
			g := bulkThree()[0]
			c := bulkCollector(t, func(_ context.Context, raw string, limit int64) ([]byte, *http.Response, error) {
				calls++
				u, e := url.Parse(raw)
				if e != nil || u.Host != "management.azure.com" || u.Query().Get("api-version") != "2024-11-01" || limit != 1048576 {
					t.Fatalf("BULK_BUDGET_REQUEST: %s %d", raw, limit)
				}
				if tc.Kind == "pages" || tc.Kind == "bytes" {
					if u.Path != bulkGroupsPath {
						t.Fatal("BULK_BUDGET_GROUP_SCOPE", raw)
					}
					if tc.Kind == "bytes" && calls == 9 {
						return bulkHTTP(" ", 200)
					} // Exactly one byte beyond eight fully accounted MiB.
					next := ""
					if calls < tc.Amount {
						next = fmt.Sprintf("%s%s?api-version=2024-11-01&%%24skiptoken=%d", bulkOrigin, bulkGroupsPath, calls+1)
					}
					body := bulkPage(nil, next)
					if tc.Kind == "bytes" {
						body += strings.Repeat(" ", 1048576-len(body))
					}
					return bulkHTTP(body, 200)
				}
				if u.Path == bulkGroupsPath {
					return bulkHTTP(bulkPage([]bulkWireGroup{g}, ""), 200)
				}
				if u.Path == bulkID(g)+"/capacityReservations" {
					names := []string{"r"}
					if tc.Kind == "calls" {
						names = nil
						for i := 0; i < tc.Amount; i++ {
							names = append(names, fmt.Sprintf("r%d", i))
						}
					}
					return bulkHTTP(bulkSummary(names...), 200)
				}
				prefix := bulkID(g) + "/capacityReservations/"
				if strings.HasPrefix(u.Path, prefix) && u.RawQuery == "%24expand=instanceView&api-version=2024-11-01" {
					allocated := 0
					if tc.Kind == "work" {
						allocated = tc.Amount
					}
					body := bulkDetail(g, strings.TrimPrefix(u.Path, prefix), g.region, allocated)
					if len(body) > 1048576 {
						t.Fatalf("BULK_RECIPE_BODY_OVERSHADOWED: %d", len(body))
					}
					return bulkHTTP(body, 200)
				}
				t.Fatal("BULK_BUDGET_UNEXPECTED_REQUEST", raw)
				return nil, nil, nil
			})
			result, e := c.CollectSubscription(context.Background(), map[string]string{bulkSubID: "Selected"}, SubscriptionReservationRequest{SubscriptionID: bulkSubID})
			if e != nil || result == nil {
				t.Fatalf("BULK_BUDGET_EXECUTION: %#v %v", result, e)
			}
			if calls != tc.Calls || result.Status != tc.Status || result.FailureCode != tc.Code || len(result.Reservations) != tc.Rows {
				t.Fatalf("BULK_BUDGET_LITERAL: %s calls=%d result=%#v", tc.ID, calls, result)
			}
			if tc.Rows == 1 && int(result.Reservations[0].Allocated) != tc.Allocated {
				t.Fatal("BULK_BUDGET_ALLOCATED", result.Reservations)
			}
			evidence := result.Evidence()
			if evidence.Status != tc.Status || evidence.FailureCode != tc.Code {
				t.Fatal("BULK_BUDGET_ADMISSION", evidence)
			}
			for _, regional := range evidence.Regional {
				if regional.Status != tc.Status {
					t.Fatal("BULK_BUDGET_PREFIX_PROMOTION", evidence)
				}
			}
		})
	}
}
