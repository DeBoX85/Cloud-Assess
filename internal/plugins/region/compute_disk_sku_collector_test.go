package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

const diskTestID = "11111111-1111-4111-8111-111111111111"
const diskFirstURL = "https://management.azure.com/subscriptions/11111111-1111-4111-8111-111111111111/providers/Microsoft.Compute/skus?%24filter=location+eq+%27eastus%27&api-version=2021-07-01"
const diskNextURL = "https://management.azure.com/subscriptions/11111111-1111-4111-8111-111111111111/providers/Microsoft.Compute/skus?%24skiptoken=two&api-version=2021-07-01"
const diskAvailableRow = `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[]}`

type diskReply struct {
	body        string
	status      int
	err         error
	nilResponse bool
	after       func()
}

type diskFakeGetter struct {
	mu      sync.Mutex
	replies []diskReply
	urls    []string
	limits  []int64
}

func (f *diskFakeGetter) GetBoundedWithResponse(_ context.Context, uri string, limit int64) ([]byte, *http.Response, error) {
	f.mu.Lock()
	index := len(f.urls)
	f.urls = append(f.urls, uri)
	f.limits = append(f.limits, limit)
	var reply diskReply
	if index < len(f.replies) {
		reply = f.replies[index]
	} else {
		reply = diskReply{err: errors.New("unexpected request")}
	}
	f.mu.Unlock()
	if reply.after != nil {
		reply.after()
	}
	status := reply.status
	if status == 0 {
		status = 200
	}
	var response *http.Response
	if !reply.nilResponse {
		response = &http.Response{StatusCode: status}
	}
	return []byte(reply.body), response, reply.err
}

func diskEnvelope(row string, next string) string {
	body := `{"value":[` + row + `]`
	if next != "" {
		encoded, _ := json.Marshal(next)
		body += `,"nextLink":` + string(encoded)
	}
	return body + "}"
}

func diskCollect(t *testing.T, replies ...diskReply) (*ComputeDiskSKUCollection, error, *diskFakeGetter) {
	t.Helper()
	getter := &diskFakeGetter{replies: replies}
	collector, err := NewComputeDiskSKUCollector("https://management.azure.com", getter)
	if err != nil {
		t.Fatal(err)
	}
	result, err := collector.Collect(context.Background(), map[string]string{diskTestID: "Selected"}, ComputeDiskSKURequest{diskTestID, "eastus"})
	return result, err, getter
}

func TestComputeDiskSKULiteralStates(t *testing.T) {
	cases := []struct {
		name string
		row  string
		want map[string]SKUAvailability
	}{
		{"empty", "", map[string]SKUAvailability{}},
		{"available", diskAvailableRow, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}},
		{"uppercase-location", `{"resourceType":"DISKS","name":"Premium_LRS","locations":["EASTUS","eastus"],"restrictions":[]}`, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}},
		{"other-region", `{"resourceType":"disks","name":"Premium_LRS","locations":["westus"],"restrictions":[]}`, map[string]SKUAvailability{"premium_lrs": {State: 2, BlockedZones: []string{}}}},
		{"empty-locations", `{"resourceType":"disks","name":"Premium_LRS","locations":[],"restrictions":[]}`, map[string]SKUAvailability{"premium_lrs": {State: 2, BlockedZones: []string{}}}},
		{"values-only-quota", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Location","reasonCode":"QuotaId","values":["eastus"]}]}`, map[string]SKUAvailability{"premium_lrs": {State: 1, BlockedZones: []string{}}}},
		{"info-only-subscription", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Location","reasonCode":"NotAvailableForSubscription","restrictionInfo":{"locations":["eastus"]}}]}`, map[string]SKUAvailability{"premium_lrs": {State: 1, BlockedZones: []string{}}}},
		{"scope-union", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Location","reasonCode":"NotAvailableForSubscription","values":["westus"],"restrictionInfo":{"locations":["eastus"]}}]}`, map[string]SKUAvailability{"premium_lrs": {State: 1, BlockedZones: []string{}}}},
		{"other-restriction", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Location","reasonCode":"QuotaId","values":["westus"]}]}`, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}},
		{"zone-union", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","values":["ignored"],"restrictionInfo":{"locations":["eastus"],"zones":["2","1"]}},{"type":"Zone","reasonCode":"NotAvailableForSubscription","restrictionInfo":{"locations":["eastus"],"zones":["1","3"]}}]}`, map[string]SKUAvailability{"premium_lrs": {State: 3, BlockedZones: []string{"1", "2", "3"}}}},
		{"regional-wins", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["1"]}},{"type":"Location","reasonCode":"QuotaId","values":["eastus"]},{"type":"Zone","reasonCode":"NotAvailableForSubscription","restrictionInfo":{"locations":["eastus"],"zones":["2"]}}]}`, map[string]SKUAvailability{"premium_lrs": {State: 1, BlockedZones: []string{}}}},
		{"absence-wins", `{"resourceType":"disks","name":"Premium_LRS","locations":["westus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["1"]}},{"type":"Location","reasonCode":"QuotaId","values":["eastus"]}]}`, map[string]SKUAvailability{"premium_lrs": {State: 2, BlockedZones: []string{}}}},
		{"other-zone", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["westus"],"zones":["1"]}}]}`, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}},
		{"other-type", `{"resourceType":"virtualMachines","name":"Standard_A0","locations":null,"restrictions":42,"capacity":{"minimum":null}}`, map[string]SKUAvailability{}},
		{"valid-metadata", `{"resourceType":"disks","name":"Premium_LRS","locations":["eastus"],"restrictions":[],"kind":"","family":"Disk","tier":"Standard","size":"Small","apiVersions":[],"capacity":{"minimum":-1,"maximum":9223372036854775807,"default":0,"scaleType":"None"},"costs":[{"meterID":"","quantity":-1,"extendedUnit":""}],"capabilities":[{"name":"","value":""}],"locationInfo":[{"location":"eastus","zones":["1"],"type":"future-edge","extendedLocations":["edge1"],"zoneDetails":[{"Name":["1"],"capabilities":[{"name":"x","value":"y"}]}]}],"unknown":{"a":[null,1,true,"text"]}}`, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			result, err, getter := diskCollect(t, diskReply{body: diskEnvelope(test.row, "")})
			if err != nil || result == nil || result.SubscriptionID != diskTestID || result.Region != "eastus" || result.Status != "complete" || result.FailureCode != "" || !reflect.DeepEqual(result.Values, test.want) || !reflect.DeepEqual(result.Evidence(), SKUEvidence{Status: "complete", Values: test.want}) {
				t.Fatalf("result=%#v error=%v; want complete literal map %#v", result, err, test.want)
			}
			if !reflect.DeepEqual(getter.urls, []string{diskFirstURL}) || !reflect.DeepEqual(getter.limits, []int64{1048576}) {
				t.Fatalf("requests=%v limits=%v", getter.urls, getter.limits)
			}
		})
	}
}

func TestComputeDiskSKURejectsMalformedPages(t *testing.T) {
	cases := []string{
		`null`, `[]`, `{}`, `{"value":null}`, `{"value":{}}`, `{"Value":[]}`,
		`{"value":[null]}`, `{"value":[],"nextLink":null}`, `{"value":[],"NextLink":""}`,
		`{"value":[],"nextLink":2}`, `{"value":[],"value":[]}`, `{"value":[],"VALUE":[]}`,
		`{"value":[]}{}`, `{"value":[]}garbage`, `{"value":[],"x":"\ud800"}`,
		`{"value":[],"x":"\u0000"}`, `{"value":[],"x":"\uffff"}`,
	}
	rows := []string{
		`{}`, `{"resourceType":"disks","name":null}`,
		`{"ResourceType":"disks","name":"x","locations":[],"restrictions":[]}`,
		`{"resourceType":"diſks","name":"x","locations":[],"restrictions":[]}`,
		`{"resourceType":"disks","name":"","locations":[],"restrictions":[]}`,
		`{"resourceType":"disks","name":" x","locations":[],"restrictions":[]}`,
		`{"resourceType":"disks","name":"x ","locations":[],"restrictions":[]}`,
		`{"resourceType":"disks","name":"é","locations":[],"restrictions":[]}`,
		`{"resourceType":"disks","name":"x","locations":null,"restrictions":[]}`,
		`{"resourceType":"disks","name":"x","locations":[]}`,
		`{"resourceType":"disks","name":"x","locations":[],"restrictions":null}`,
		`{"resourceType":"disks","name":"x","locations":["global"],"restrictions":[]}`,
		`{"resourceType":"disks","name":"x","locations":["East US"],"restrictions":[]}`,
		`{"resourceType":"disks","name":"x","locations":["ſouth"],"restrictions":[]}`,
		`{"resourceType":"disks","name":"x","locations":[null],"restrictions":[]}`,
		`{"resourceType":"disks","name":"x","locations":[],"restrictions":[null]}`,
	}
	restrictions := []string{
		`{}`, `{"type":null,"reasonCode":"QuotaId"}`, `{"type":"Location","reasonCode":"Unknown","values":["eastus"]}`,
		`{"type":"Future","reasonCode":"QuotaId","values":["eastus"]}`,
		`{"type":"location","reasonCode":"QuotaId","values":["eastus"]}`,
		`{"type":"Location","reasonCode":"quotaid","values":["eastus"]}`,
		`{"Type":"Location","reasonCode":"QuotaId","values":["eastus"]}`,
		`{"type":"Location","reasonCode":"QuotaId","Values":["eastus"]}`,
		`{"type":"Location","reasonCode":"QuotaId","values":[]}`,
		`{"type":"Location","reasonCode":"QuotaId","values":null}`,
		`{"type":"Location","reasonCode":"QuotaId","values":["global"]}`,
		`{"type":"Location","reasonCode":"QuotaId","values":["eastus"],"restrictionInfo":null}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":[]}}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":[],"zones":["1"]}}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":null}}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["0"]}}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["01"]}}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["١"]}}`,
		`{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":[1]}}`,
	}
	for _, restriction := range restrictions {
		rows = append(rows, `{"resourceType":"disks","name":"x","locations":["eastus"],"restrictions":[`+restriction+"]}")
	}
	metadata := []string{
		`"tier":null`, `"kind":{}`, `"family":[]`, `"size":1`, `"apiVersions":null`, `"apiVersions":[null]`,
		`"capacity":null`, `"capacity":{"minimum":null}`, `"capacity":{"maximum":9223372036854775808}`,
		`"capacity":{"default":1.0}`, `"capacity":{"scaleType":"future"}`, `"capacity":{"Minimum":1}`,
		`"costs":null`, `"costs":[null]`, `"costs":[{"meterID":null}]`, `"costs":[{"quantity":1.1}]`,
		`"capabilities":null`, `"capabilities":[null]`, `"capabilities":[{"value":null}]`,
		`"locationInfo":null`, `"locationInfo":[null]`, `"locationInfo":[{"location":"global"}]`,
		`"locationInfo":[{"zones":["01"]}]`, `"locationInfo":[{"type":null}]`,
		`"locationInfo":[{"extendedLocations":null}]`, `"locationInfo":[{"zoneDetails":[null]}]`,
		`"locationInfo":[{"zoneDetails":[{"name":["1"]}]}]`,
		`"locationInfo":[{"zoneDetails":[{"Name":null}]}]`,
	}
	for _, extra := range metadata {
		rows = append(rows, strings.TrimSuffix(diskAvailableRow, "}")+","+extra+"}")
	}
	for _, row := range rows {
		cases = append(cases, diskEnvelope(row, ""))
	}
	cases = append(cases, "{\"value\":[],\"x\":\""+string([]byte{0xff})+"\"}")
	for index, body := range cases {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			result, err, getter := diskCollect(t, diskReply{body: body})
			if err != nil || result == nil || result.Status != "unknown" || result.FailureCode == "" || len(result.Values) != 0 || !reflect.DeepEqual(result.Evidence(), SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}}) || len(getter.urls) != 1 {
				t.Fatalf("body=%q result=%#v error=%v", body, result, err)
			}
		})
	}
}

func TestComputeDiskSKURequestsAndContinuation(t *testing.T) {
	second := `{"resourceType":"disks","name":"Standard_LRS","locations":["westus"],"restrictions":[]}`
	for _, next := range []string{
		"?api-version=2021-07-01&$skiptoken=two",
		diskNextURL,
		"?$filter=location%20eq%20%27eastus%27&api-version=2021-07-01&$skiptoken=two",
	} {
		result, err, getter := diskCollect(t, diskReply{body: diskEnvelope(diskAvailableRow, next)}, diskReply{body: diskEnvelope(second, "")})
		want := map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}, "standard_lrs": {State: 2, BlockedZones: []string{}}}
		if err != nil || result.Status != "complete" || !reflect.DeepEqual(result.Values, want) || getter.urls[0] != diskFirstURL || len(getter.urls) != 2 {
			t.Fatalf("next=%s result=%#v err=%v requests=%v", next, result, err, getter.urls)
		}
		if !strings.Contains(next, "$filter") && getter.urls[1] != diskNextURL {
			t.Fatalf("second request=%s", getter.urls[1])
		}
	}
	unsafe := []string{
		"https://evil.test/subscriptions/" + diskTestID + "/providers/Microsoft.Compute/skus?api-version=2021-07-01",
		"http://management.azure.com/subscriptions/" + diskTestID + "/providers/Microsoft.Compute/skus?api-version=2021-07-01",
		"https://user@management.azure.com/subscriptions/" + diskTestID + "/providers/Microsoft.Compute/skus?api-version=2021-07-01",
		"/subscriptions/22222222-2222-4222-8222-222222222222/providers/Microsoft.Compute/skus?api-version=2021-07-01",
		"/subscriptions/" + diskTestID + "/providers/Microsoft.Storage/skus?api-version=2021-07-01",
		"/subscriptions/" + diskTestID + "/providers/Microsoft.Compute/%73kus?api-version=2021-07-01",
		"?api-version=2023-03-01", "?api-version=2021-07-01&api-version=2021-07-01",
		"?Api-Version=2021-07-01", "?api-version=2021-07-01&$filter=location+eq+%27westus%27",
		"?api-version=2021-07-01&$filter=location+eq+%27eastus%27&$filter=location+eq+%27eastus%27",
		"?api-version=2021-07-01&includeExtendedLocations=true", "?api-version=2021-07-01&$skiptoken=",
		"?api-version=2021-07-01&$skiptoken=a&$skiptoken=b", "?api-version=2021-07-01#fragment",
		"?api-version=2021-07-01&x=unknown", "?api-version=2021-07-01&$skiptoken=" + strings.Repeat("x", 4097),
		"bad next link", "https:opaque",
	}
	for _, next := range unsafe {
		result, err, getter := diskCollect(t, diskReply{body: diskEnvelope(diskAvailableRow, next)})
		if err != nil || result.Status != "partial" || result.FailureCode != "sku_unsafe_continuation" || len(getter.urls) != 1 || !reflect.DeepEqual(result.Values, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}) || len(result.Evidence().Values) != 0 {
			t.Fatalf("next=%s result=%#v err=%v requests=%v", next, result, err, getter.urls)
		}
	}
	result, err, getter := diskCollect(t, diskReply{body: diskEnvelope("", diskFirstURL)})
	if err != nil || result.Status != "partial" || result.FailureCode != "sku_continuation_cycle" || len(getter.urls) != 1 {
		t.Fatalf("cycle result=%#v err=%v requests=%v", result, err, getter.urls)
	}
	exact := "?api-version=2021-07-01&$skiptoken=" + strings.Repeat("x", 4096)
	result, err, getter = diskCollect(t, diskReply{body: diskEnvelope("", exact)}, diskReply{body: diskEnvelope("", "")})
	if err != nil || result.Status != "complete" || len(getter.urls) != 2 {
		t.Fatalf("exact token result=%#v err=%v calls=%d", result, err, len(getter.urls))
	}
}

func TestComputeDiskSKUAtomicFailureAndSnapshot(t *testing.T) {
	duplicate := `{"resourceType":"disks","name":"premium_lrs","locations":[],"restrictions":[]}`
	for _, second := range []diskReply{
		{body: diskEnvelope(strings.Replace(diskAvailableRow, "Premium_LRS", "Standard_LRS", 1)+",{}", "")},
		{body: diskEnvelope(duplicate, "")},
		{body: `{"value":[null]}`},
		{body: `{"value":[]}`, status: 403},
		{err: errors.New("PRIVATE_PROVIDER_TEXT")},
		{nilResponse: true},
	} {
		result, err, getter := diskCollect(t, diskReply{body: diskEnvelope(diskAvailableRow, "?api-version=2021-07-01&$skiptoken=two")}, second)
		want := map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}
		if err != nil || result.Status != "partial" || !reflect.DeepEqual(result.Values, want) || len(getter.urls) != 2 || strings.Contains(result.FailureCode, "PRIVATE") || !reflect.DeepEqual(result.Evidence(), SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}}) {
			t.Fatalf("result=%#v err=%v requests=%v", result, err, getter.urls)
		}
	}
	result, err, _ := diskCollect(t, diskReply{body: diskEnvelope(diskAvailableRow+","+duplicate, "")})
	if err != nil || result.Status != "unknown" || result.FailureCode != "sku_identity_duplicate" || len(result.Values) != 0 {
		t.Fatalf("same-page duplicate result=%#v err=%v", result, err)
	}
	for _, status := range []int{201, 204, 206, 301, 302, 307, 308, 400, 401, 403, 404, 405, 429, 500, 503} {
		result, err, getter := diskCollect(t, diskReply{body: diskEnvelope("", ""), status: status})
		if err != nil || result.Status != "unknown" || result.FailureCode != "sku_status_invalid" || len(result.Values) != 0 || len(getter.urls) != 1 {
			t.Fatalf("status=%d result=%#v err=%v", status, result, err)
		}
	}
	result, err, _ = diskCollect(t, diskReply{body: diskEnvelope("", "?api-version=2021-07-01&$skiptoken=two")}, diskReply{err: errors.New("controlled")})
	if err != nil || result.Status != "partial" || len(result.Values) != 0 || result.Evidence().Status != "unknown" {
		t.Fatalf("empty prefix result=%#v err=%v", result, err)
	}
	var nilCollection *ComputeDiskSKUCollection
	for _, value := range []*ComputeDiskSKUCollection{nilCollection, {Status: "partial"}, {Status: "unknown"}, {Status: "unsupported"}, {Status: "other"}} {
		if !reflect.DeepEqual(value.Evidence(), SKUEvidence{Status: "unknown", Values: map[string]SKUAvailability{}}) {
			t.Fatal("noncomplete collection manufactured complete evidence")
		}
	}
	owned := &ComputeDiskSKUCollection{Status: "complete", Values: map[string]SKUAvailability{"x": {State: 3, BlockedZones: []string{"1", "2"}}}}
	first := owned.Evidence()
	first.Values["x"].BlockedZones[0] = "changed"
	first.Values["new"] = SKUAvailability{State: 0}
	other := owned.Evidence()
	if !reflect.DeepEqual(other, SKUEvidence{Status: "complete", Values: map[string]SKUAvailability{"x": {State: 3, BlockedZones: []string{"1", "2"}}}}) {
		t.Fatalf("snapshot shares maps/slices: %#v", other)
	}
}

func TestComputeDiskSKURequestAndCancellation(t *testing.T) {
	getter := &diskFakeGetter{}
	for _, origin := range []string{"http://management.azure.com", "https://user@management.azure.com", "https://management.azure.com/x", "https://management.azure.com?x=y", "https://management.azure.com#x", "https://management.azure.com:444", "https://bad..host"} {
		if value, err := NewComputeDiskSKUCollector(origin, getter); err == nil || value != nil {
			t.Fatalf("accepted invalid origin %s", origin)
		}
	}
	if value, err := NewComputeDiskSKUCollector("https://management.azure.com", nil); err == nil || value != nil {
		t.Fatal("accepted nil getter")
	}
	collector, _ := NewComputeDiskSKUCollector("https://management.azure.com", getter)
	cases := []struct {
		scope   map[string]string
		request ComputeDiskSKURequest
	}{
		{nil, ComputeDiskSKURequest{diskTestID, "eastus"}},
		{map[string]string{diskTestID: ""}, ComputeDiskSKURequest{diskTestID, "eastus"}},
		{map[string]string{"bad": "name"}, ComputeDiskSKURequest{diskTestID, "eastus"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{"22222222-2222-4222-8222-222222222222", "eastus"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "EastUS"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "global"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "europe"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "East US"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "ſouth"}},
		{map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, strings.Repeat("x", 65)}},
	}
	for _, test := range cases {
		value, err := collector.Collect(context.Background(), test.scope, test.request)
		if err == nil || value != nil || len(getter.urls) != 0 {
			t.Fatalf("request before input validation: %#v %v %v", value, err, getter.urls)
		}
	}
	aliasID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	value, err := collector.Collect(context.Background(), map[string]string{aliasID: "one", strings.ToUpper(aliasID): "two"}, ComputeDiskSKURequest{aliasID, "eastus"})
	if err == nil || value != nil || len(getter.urls) != 0 {
		t.Fatal("normalized scope collision accepted")
	}
	scope := map[string]string{}
	for i := 0; i < 1001; i++ {
		scope[fmt.Sprintf("%08x-0000-4000-8000-000000000000", i)] = "Selected"
	}
	if value, err := collector.Collect(context.Background(), scope, ComputeDiskSKURequest{diskTestID, "eastus"}); err == nil || value != nil || len(getter.urls) != 0 {
		t.Fatal("scope limit not enforced before request")
	}
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		f := &diskFakeGetter{replies: []diskReply{{err: fmt.Errorf("wrapped: %w", cause)}}}
		c, _ := NewComputeDiskSKUCollector("https://management.azure.com", f)
		value, err := c.Collect(context.Background(), map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "eastus"})
		if value != nil || !errors.Is(err, cause) {
			t.Fatalf("getter context error became evidence: %#v %v", value, err)
		}
	}
	for _, late := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		f := &diskFakeGetter{replies: []diskReply{{body: diskEnvelope(diskAvailableRow, ""), after: cancel}}}
		c, _ := NewComputeDiskSKUCollector("https://management.azure.com", f)
		if !late {
			cancel()
		}
		value, err := c.Collect(ctx, map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "eastus"})
		if value != nil || !errors.Is(err, context.Canceled) || len(f.urls) != map[bool]int{false: 0, true: 1}[late] {
			t.Fatalf("cancel value=%#v err=%v calls=%d", value, err, len(f.urls))
		}
		cancel()
	}
}

func TestComputeDiskSKUBounds(t *testing.T) {
	empty := diskEnvelope("", "")
	for _, n := range []int{1048576, 1048577} {
		body := empty + strings.Repeat(" ", n-len(empty))
		value, err, _ := diskCollect(t, diskReply{body: body})
		want := "complete"
		if n > 1048576 {
			want = "unknown"
		}
		if err != nil || value.Status != want || n > 1048576 && value.FailureCode != "sku_byte_limit" {
			t.Fatalf("page bytes=%d value=%#v err=%v", n, value, err)
		}
	}
	for _, excess := range []bool{false, true} {
		var replies []diskReply
		for i := 0; i < 8; i++ {
			next := ""
			if i < 7 || excess {
				next = "?api-version=2021-07-01&$skiptoken=" + fmt.Sprint(i+1)
			}
			body := diskEnvelope("", next)
			body += strings.Repeat(" ", 1048576-len(body))
			replies = append(replies, diskReply{body: body})
		}
		if excess {
			replies = append(replies, diskReply{body: empty})
		}
		value, err, f := diskCollect(t, replies...)
		want := "complete"
		if excess {
			want = "partial"
		}
		if err != nil || value.Status != want || excess && value.FailureCode != "sku_byte_limit" || len(f.urls) != 8+map[bool]int{false: 0, true: 1}[excess] {
			t.Fatalf("aggregate excess=%v value=%#v err=%v calls=%d", excess, value, err, len(f.urls))
		}
	}
	for _, excess := range []bool{false, true} {
		var replies []diskReply
		for i := 0; i < 64; i++ {
			next := ""
			if i < 63 || excess {
				next = "?api-version=2021-07-01&$skiptoken=" + fmt.Sprint(i+1)
			}
			replies = append(replies, diskReply{body: diskEnvelope("", next)})
		}
		value, err, f := diskCollect(t, replies...)
		want := "complete"
		if excess {
			want = "partial"
		}
		if err != nil || value.Status != want || excess && value.FailureCode != "sku_page_limit" || len(f.urls) != 64 {
			t.Fatalf("page count excess=%v value=%#v err=%v calls=%d", excess, value, err, len(f.urls))
		}
	}
	for _, excess := range []bool{false, true} {
		first := `{"value":[],"scratch":[` + strings.TrimSuffix(strings.Repeat("null,", 32767), ",") + `],"nextLink":"?api-version=2021-07-01&$skiptoken=two"}`
		count := 32767
		if excess {
			count++
		}
		second := `{"value":[],"scratch":[` + strings.TrimSuffix(strings.Repeat("null,", count), ",") + "]}"
		value, err, f := diskCollect(t, diskReply{body: first}, diskReply{body: second})
		want := "complete"
		if excess {
			want = "partial"
		}
		if err != nil || value.Status != want || excess && value.FailureCode != "sku_work_limit" || len(f.urls) != 2 {
			t.Fatalf("raw work excess=%v value=%#v err=%v", excess, value, err)
		}
	}
	for _, n := range []int{512, 513} {
		row := strings.Replace(diskAvailableRow, "Premium_LRS", strings.Repeat("x", n), 1)
		value, err, _ := diskCollect(t, diskReply{body: diskEnvelope(row, "")})
		want := "complete"
		if n > 512 {
			want = "unknown"
		}
		if err != nil || value.Status != want {
			t.Fatalf("name bytes=%d value=%#v err=%v", n, value, err)
		}
	}
	for _, extra := range []int{0, 1} {
		budget := diskPageBudget{ctx: context.Background(), text: auxTextBudget - 3}
		code := budget.walk(map[string]any{"x": strings.Repeat("y", 2+extra)}, 0, false)
		want := ""
		if extra == 1 {
			want = "sku_text_limit"
		}
		if code != want {
			t.Fatalf("decoded text edge extra=%d code=%s", extra, code)
		}
	}
}

type diskCanaryCredential struct {
	calls  atomic.Int32
	mu     sync.Mutex
	scopes []string
}

func (c *diskCanaryCredential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	c.mu.Lock()
	c.scopes = append([]string{}, options.Scopes...)
	c.mu.Unlock()
	return azcore.AccessToken{Token: "OFFLINE_CANARY", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type diskCloseBody struct {
	io.Reader
	closes *atomic.Int32
}

func (b *diskCloseBody) Close() error { b.closes.Add(1); return nil }

type diskCanaryTransport struct {
	t      *testing.T
	calls  atomic.Int32
	closes atomic.Int32
	status int
}

func (f *diskCanaryTransport) Do(request *http.Request) (*http.Response, error) {
	f.calls.Add(1)
	if request.Method != "GET" || request.URL.String() != diskFirstURL || request.Body != nil && request.Body != http.NoBody || request.Header.Get("Authorization") != "Bearer OFFLINE_CANARY" {
		f.t.Errorf("unexpected authenticated request: method=%s uri=%s body=%v", request.Method, request.URL, request.Body)
	}
	status := f.status
	if status == 0 {
		status = 200
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: &diskCloseBody{Reader: strings.NewReader(diskEnvelope(diskAvailableRow, "")), closes: &f.closes}, Request: request}, nil
}

func TestComputeDiskSKUActualBoundedTransportAndIsolation(t *testing.T) {
	for _, status := range []int{200, 403, 307} {
		credential := &diskCanaryCredential{}
		wire := &diskCanaryTransport{t: t, status: status}
		client := azure.NewHTTPClient(credential, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, OperationTimeout: 5 * time.Second, Scope: "https://management.azure.com/.default", Transport: wire})
		c, _ := NewComputeDiskSKUCollector("https://management.azure.com", client)
		value, err := c.Collect(context.Background(), map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "eastus"})
		if err != nil || value == nil || wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 || !reflect.DeepEqual(credential.scopes, []string{"https://management.azure.com/.default"}) {
			t.Fatalf("bounded transport status=%d value=%#v err=%v calls=%d closes=%d token=%d", status, value, err, wire.calls.Load(), wire.closes.Load(), credential.calls.Load())
		}
		if status == 200 && value.Status != "complete" || status != 200 && value.Status != "unknown" {
			t.Fatalf("status=%d value=%#v", status, value)
		}
		_, err = c.Collect(context.Background(), map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "global"})
		if err == nil || wire.calls.Load() != 1 || credential.calls.Load() != 1 {
			t.Fatal("invalid scope authenticated")
		}
	}
	replies := make([]diskReply, 8)
	for i := range replies {
		replies[i] = diskReply{body: diskEnvelope(diskAvailableRow, "")}
	}
	getter := &diskFakeGetter{replies: replies}
	collector, _ := NewComputeDiskSKUCollector("https://management.azure.com", getter)
	results := make(chan *ComputeDiskSKUCollection, 8)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			value, err := collector.Collect(context.Background(), map[string]string{diskTestID: "name"}, ComputeDiskSKURequest{diskTestID, "eastus"})
			if err != nil {
				t.Error(err)
				return
			}
			results <- value
		}()
	}
	group.Wait()
	close(results)
	count := 0
	for value := range results {
		count++
		if value.Status != "complete" || !reflect.DeepEqual(value.Values, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}) {
			t.Fatalf("isolated result=%#v", value)
		}
		value.Values["changed"] = SKUAvailability{State: 1}
	}
	if count != 8 || len(getter.urls) != 8 {
		t.Fatalf("isolated calls/results=%d/%d", len(getter.urls), count)
	}
}

func TestComputeDiskSKUJSONEdgesAndContext(t *testing.T) {
	for _, objects := range []int{31, 32} {
		deep := strings.Repeat(`{"x":`, objects) + "null" + strings.Repeat("}", objects)
		body := `{"value":[],"metadata":` + deep + "}"
		value, err, _ := diskCollect(t, diskReply{body: body})
		want := "complete"
		if objects == 32 {
			want = "unknown"
		}
		if err != nil || value.Status != want {
			t.Fatalf("depth objects=%d value=%#v err=%v", objects, value, err)
		}
	}
	for _, extra := range []int{127, 128} {
		body := `{"value":[]`
		for i := 0; i < extra; i++ {
			body += fmt.Sprintf(",%q:null", fmt.Sprintf("x%03d", i))
		}
		body += "}"
		value, err, _ := diskCollect(t, diskReply{body: body})
		want := "complete"
		if extra == 128 {
			want = "unknown"
		}
		if err != nil || value.Status != want {
			t.Fatalf("object keys=%d value=%#v err=%v", extra+1, value, err)
		}
	}
	for _, count := range []int{65533, 65534} {
		body := `{"value":[],"metadata":[` + strings.TrimSuffix(strings.Repeat("null,", count), ",") + "]}"
		value, err, _ := diskCollect(t, diskReply{body: body})
		want := "complete"
		if count == 65534 {
			want = "unknown"
		}
		if err != nil || value.Status != want {
			t.Fatalf("JSON token edge count=%d value=%#v err=%v", count, value, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	page, code := decodeDiskPage(ctx, []byte(diskEnvelope(diskAvailableRow, "")), "eastus", map[string]SKUAvailability{}, 0, 0)
	if code == "" || page.values != nil {
		t.Fatalf("cancelled decoder admitted page: %#v %s", page, code)
	}
}
