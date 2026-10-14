package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
)

// This deterministic fixture models the documented service freedom to reorder
// equal sort keys between pages. It does not execute KQL or claim live evidence.
type tiedInventoryPoster struct {
	calls int
}

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
	order := strings.Split(request.Query, "| order by ")
	if len(order) != 2 {
		return nil, fmt.Errorf("inventory sort clause missing")
	}
	unique := false
	for _, column := range strings.Split(order[1], ",") {
		fields := strings.Fields(strings.ToLower(column))
		if len(fields) > 0 && fields[0] == "id" {
			unique = true
		}
	}
	rows := make([]map[string]any, 0, 1000)
	for index := offset; index < min(offset+1000, 1001); index++ {
		identity := index
		if offset == 1000 && !unique {
			// Same subscription/RG keys allow the prior last item to move here,
			// replacing the final distinct item without changing page counts.
			identity = 999
		}
		rows = append(rows, map[string]any{
			"id":             fmt.Sprintf("/subscriptions/sub/resourceGroups/shared/providers/Microsoft.Test/widgets/%04d", identity),
			"subscriptionId": "sub",
			"resourceGroup":  "shared",
			"type":           "Microsoft.Test/widgets",
			"name":           fmt.Sprintf("%04d", identity),
			"kind":           "test",
		})
	}
	payload := map[string]any{"data": rows, "count": len(rows), "totalRecords": 1001, "resultTruncated": "false"}
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
