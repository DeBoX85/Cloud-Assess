package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
)

// This model evaluates the leading ID projection and ordinary column sorting
// over independently constructed rows. Equal keys have different allowed input
// ordering on the second request. It is not a general KQL engine or live evidence.
type tiedInventoryPoster struct {
	calls int
}

var modeledIDProjection = regexp.MustCompile(`(?i)\|\s*project\s+id\s*=\s*tostring\(\s*([a-z][a-z0-9]*)\s*\)\s*,`)

func (p *tiedInventoryPoster) PostStream(_ context.Context, _ string, body io.ReadSeekCloser) (*http.Response, error) {
	var request arg.Request
	if err := json.NewDecoder(body).Decode(&request); err != nil {
		return nil, err
	}
	p.calls++
	if p.calls > 2 || request.Options == nil || request.Options.Top == nil || *request.Options.Top != 1000 {
		return nil, fmt.Errorf("unexpected inventory paging contract")
	}
	offset := 0
	if request.Options.SkipToken != nil {
		var err error
		offset, err = strconv.Atoi(*request.Options.SkipToken)
		if err != nil || offset != 1000 {
			return nil, fmt.Errorf("unexpected inventory continuation")
		}
	}
	project := modeledIDProjection.FindStringSubmatch(request.Query)
	order := strings.Split(request.Query, "| order by ")
	if len(project) != 2 || len(order) != 2 {
		return nil, fmt.Errorf("unsupported fixture projection or ordering")
	}
	type sortColumn struct {
		name      string
		ascending bool
	}
	columns := []sortColumn{}
	for _, text := range strings.Split(order[1], ",") {
		fields := strings.Fields(strings.ToLower(text))
		if len(fields) < 1 || len(fields) > 2 || len(fields) == 2 && fields[1] != "asc" && fields[1] != "desc" {
			return nil, fmt.Errorf("unsupported fixture sort column")
		}
		columns = append(columns, sortColumn{fields[0], len(fields) == 2 && fields[1] == "asc"})
	}
	rows := make([]map[string]string, 0, 1001)
	for index := 0; index < 1001; index++ {
		// Values and identities are derived only from the independent fixture,
		// never from the query's alleged uniqueness or the implementation result.
		raw := map[string]string{
			"id":             fmt.Sprintf("/subscriptions/sub/resourceGroups/shared/providers/Microsoft.Test/widgets/%04d", index),
			"subscriptionid": "sub",
			"resourcegroup":  "shared",
			"type":           "Microsoft.Test/widgets",
			"name":           fmt.Sprintf("%04d", index),
			"kind":           "test",
		}
		id, ok := raw[strings.ToLower(project[1])]
		if !ok {
			return nil, fmt.Errorf("unsupported fixture projection field")
		}
		raw["id"] = id
		rows = append(rows, raw)
	}
	if offset > 0 {
		rows[999], rows[1000] = rows[1000], rows[999]
	}
	for _, col := range columns {
		if _, ok := rows[0][col.name]; !ok {
			return nil, fmt.Errorf("unsupported fixture ordering field")
		}
	}
	slices.SortStableFunc(rows, func(a, b map[string]string) int {
		for _, col := range columns {
			comparison := strings.Compare(a[col.name], b[col.name])
			if comparison != 0 {
				if !col.ascending {
					comparison = -comparison
				}
				return comparison
			}
		}
		return 0
	})
	page := make([]map[string]string, 0, 1000)
	for _, row := range rows[offset:min(offset+1000, len(rows))] {
		page = append(page, map[string]string{
			"id": row["id"], "subscriptionId": row["subscriptionid"],
			"resourceGroup": row["resourcegroup"], "type": row["type"],
			"name": row["name"], "kind": row["kind"],
		})
	}
	payload := map[string]any{"data": page, "count": len(page), "totalRecords": 1001, "resultTruncated": "false"}
	if offset == 0 {
		payload["$skipToken"] = "1000"
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(data)))}, nil
}

func TestInventoryPagingRetainsAllResourcesWithTiedSubscriptionAndGroup(t *testing.T) {
	poster := &tiedInventoryPoster{}
	client := arg.NewClient(arg.NewHTTPTransportWithClient(poster, "https://example.test/graph"))
	inventory, err := DiscoverResources(context.Background(), client, map[string]string{"sub": "Synthetic"}, nil)
	if err != nil || inventory == nil || poster.calls != 2 || len(inventory.Included) != 1001 {
		t.Fatalf("inventory paging failed: error=%v inventory=%v calls=%d", err, inventory, poster.calls)
	}
	seen := make(map[string]bool, 1001)
	for _, resource := range inventory.Included {
		if seen[resource.ID] {
			t.Fatalf("tied sort duplicated resource and concealed missing inventory: %s", resource.ID)
		}
		seen[resource.ID] = true
		if resource.Kind != "test" || resource.SubscriptionID != "sub" || resource.ResourceGroup != "shared" {
			t.Fatal("inventory projection or scope changed")
		}
	}
	for index := 0; index < 1001; index++ {
		id := fmt.Sprintf("/subscriptions/sub/resourceGroups/shared/providers/Microsoft.Test/widgets/%04d", index)
		if !seen[id] {
			t.Fatalf("inventory paging omitted resource: %s", id)
		}
	}
}

func TestInventoryPagingModelUsesActualProjectedAndSortedValues(t *testing.T) {
	const projection = "resources | project id=tostring(id), subscriptionId=tostring(subscriptionId), resourceGroup=tostring(resourceGroup), type=tostring(type), name=tostring(name), kind=tostring(kind)"
	for _, tc := range []struct {
		name, query string
		wantUnique  int
	}{
		{"tied keys", projection + " | order by subscriptionId, resourceGroup", 1000},
		{"unique ascending", projection + " | order by subscriptionId, resourceGroup, id asc", 1001},
		{"unique descending", projection + " | order by subscriptionId, resourceGroup, id desc", 1001},
		{"nonunique projected id", strings.Replace(projection, "id=tostring(id)", "id=tostring(subscriptionId)", 1) + " | order by subscriptionId, resourceGroup, id asc", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			poster := &tiedInventoryPoster{}
			client := arg.NewClient(arg.NewHTTPTransportWithClient(poster, "https://example.test/graph"))
			result, err := client.Query(context.Background(), tc.query, map[string]string{"sub": "Synthetic"})
			if err != nil || result == nil || len(result.Data) != 1001 || poster.calls != 2 {
				t.Fatalf("model query failed: result=%v error=%v calls=%d", result, err, poster.calls)
			}
			unique := map[string]bool{}
			for _, raw := range result.Data {
				var row map[string]string
				if err := json.Unmarshal(raw, &row); err != nil {
					t.Fatal(err)
				}
				unique[row["id"]] = true
			}
			if len(unique) != tc.wantUnique {
				t.Fatalf("model ignored actual projection or ordering: unique=%d want=%d", len(unique), tc.wantUnique)
			}
		})
	}
}
