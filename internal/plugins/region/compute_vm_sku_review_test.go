package region

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestComputeVMSKUPageCollisionAndMetadata(t *testing.T) {
	calls := 0
	c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		body := vmReadyBody
		if calls == 1 {
			body = strings.TrimSuffix(body, "}") + `,"nextLink":"` + vmNext(1) + `"}`
		}
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	r, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest())
	if e != nil || calls != 2 || r.Status != "partial" || r.FailureCode != "vm_sku_collision" || len(r.Groups) != 1 || r.Groups["s"].State != 0 {
		t.Fatal("cross-page collision lost prefix", r, e, calls)
	}
	evidence, e := r.Evidence(context.Background(), []string{"s"})
	if e != nil || evidence.Status != "unknown" || len(evidence.Values) != 0 {
		t.Fatal("partial collision scored", evidence, e)
	}
	valid := `{"value":[{"resourceType":"virtualMachines","name":"s","locations":["eastus"],"restrictions":[],"capacity":{"minimum":0,"maximum":10,"default":1,"scaleType":"Future"},"costs":[{"meterID":"m","quantity":-1,"extendedUnit":"unit"}],"apiVersions":["x"],"locationInfo":[{"location":"eastus","zones":["2","1"],"zoneDetails":[{"Name":["1"],"capabilities":[{"name":"x","value":"y"}]}],"extendedLocations":["Seattle"],"type":"Future"}]}]}`
	if got := vmCollect(t, valid); got.Status != "complete" || !got.Groups["s"].Known || got.Groups["s"].State != 0 {
		t.Fatal("valid optional metadata", got)
	}
	for _, body := range []string{strings.Replace(valid, `"minimum":0`, `"minimum":-1`, 1), strings.Replace(valid, `"quantity":-1`, `"quantity":1.5`, 1), strings.Replace(valid, `"Name":["1"]`, `"name":["1"]`, 1), strings.Replace(valid, `"extendedLocations":["Seattle"]`, `"extendedLocations":null`, 1), strings.Replace(valid, `"zones":["2","1"]`, `"zones":["01"]`, 1), strings.Replace(valid, `"apiVersions":["x"]`, `"apiVersions":{}`, 1)} {
		vmFault(t, body, "vm_sku_invalid_response")
	}
}

func TestComputeVMSKUAccountingAndWorkCancellation(t *testing.T) {
	budget := vmSKUBudget{ctx: context.Background(), text: 16711679}
	if !budget.label("x", 512) || budget.label("y", 512) {
		t.Fatal("private decoded-text exact/excess guard")
	}
	body := `{"value":[],"ignored":[` + strings.TrimSuffix(strings.Repeat(`{"x":"value"},`, 4000), ",") + `]}`
	c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		return []byte(body), &http.Response{StatusCode: 200}, nil
	}))
	ctx := &storageCountingContext{Context: context.Background(), limit: 1000}
	if r, e := c.Collect(ctx, map[string]string{vmTestID: "selected"}, vmRequest()); r != nil || !errors.Is(e, context.Canceled) {
		t.Fatal("decode cancellation", r, e)
	}
	required := []string{}
	for i := 0; i < 8192; i++ {
		required = append(required, fmt.Sprint(i))
	}
	ctx = &storageCountingContext{Context: context.Background(), limit: 1000}
	if r, e := vmEmptyCollection().Evidence(ctx, required); len(r.Values) != 0 || !errors.Is(e, context.Canceled) {
		t.Fatal("projection cancellation", r, e)
	}
}

func TestComputeVMSKUGetterTerminalAndPreauth(t *testing.T) {
	if c, e := NewComputeVMSKUCollector("https://management.azure.com", nil); c != nil || e == nil {
		t.Fatal("nil getter", c, e)
	}
	calls := 0
	c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
		calls++
		return []byte(vmReadyBody), &http.Response{StatusCode: 200}, nil
	}))
	for _, n := range []int{512, 513} {
		before := calls
		r, e := c.Collect(context.Background(), map[string]string{vmTestID: strings.Repeat("x", n)}, vmRequest())
		if n == 512 && (e != nil || r == nil || calls != before+1) || n == 513 && (e == nil || r != nil || calls != before) {
			t.Fatal("preauth display-name bound", n, r, e, calls)
		}
	}
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		c, _ := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(context.Context, string, int64) ([]byte, *http.Response, error) {
			return []byte(vmReadyBody), &http.Response{StatusCode: 403}, fmt.Errorf("canary: %w", cause)
		}))
		if r, e := c.Collect(context.Background(), map[string]string{vmTestID: "selected"}, vmRequest()); r != nil || !errors.Is(e, cause) {
			t.Fatal("wrapped terminal getter", r, e)
		}
	}
}

func TestComputeVMSKUProjectionConsumer(t *testing.T) {
	body := `{"value":[{"resourceType":"virtualMachines","name":"s","locations":["eastus"],"restrictions":[]},{"resourceType":"virtualMachines","name":"z","locations":["eastus"],"restrictions":[{"type":"Zone","reasonCode":"QuotaId","restrictionInfo":{"locations":["eastus"],"zones":["1"]}}]}]}`
	for _, kind := range []string{vmTestType, "microsoft.compute/virtualmachinescalesets"} {
		collection := vmCollect(t, body)
		collection.ResourceType = kind
		evidence, e := collection.Evidence(context.Background(), []string{"s", "z"})
		if e != nil {
			t.Fatal(e)
		}
		inventory := &InventoryCalculation{Subscriptions: map[string]InventoryCounts{vmTestID: {ResourceTypesByRegion: map[string]map[string]int64{"westus": {kind: 2}}, SKUsByTypeAndRegion: map[string]map[string]map[string]int64{kind: {"westus": {"s": 1, "z": 1}}}}}}
		got, e := CalculateAvailability(context.Background(), map[string]string{vmTestID: "selected"}, inventory, AvailabilityRequest{SubscriptionID: vmTestID, SourceRegion: "westus", TargetRegion: "eastus"}, AvailabilityEvidence{SubscriptionID: vmTestID, TargetRegion: "eastus", InventoryComplete: true, ProvidersComplete: true, Locations: map[string]map[string]bool{kind: {"eastus": true}}, SKUEnabled: true, SKUs: map[string]SKUEvidence{kind: evidence}, ZoneCounts: map[string]int{"westus": 3, "eastus": 3}})
		if e != nil || got == nil {
			t.Fatal(got, e)
		}
		x := got.Comparison
		if kind == vmTestType {
			if x.TotalSKUsChecked != 2 || x.AvailableSKUs != 1 || x.UnknownSKUs != 0 || x.SKUAvailabilityPercent != 50 || len(x.ZoneRestrictedSKUs) != 1 || x.ZoneRestrictedSKUs[0] != "microsoft.compute/virtualmachines:z (zones blocked: 1)" || got.Health.Status != assessment.StageCompleted {
				t.Fatal("literal VM arithmetic", x, got.Health)
			}
		} else if x.TotalSKUsChecked != 0 || x.UnknownSKUs != 0 || x.AvailableSKUs != 0 || x.SKUAvailabilityPercent != 0 || got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Code != "availability_sku_unsupported" {
			t.Fatal("VMSS unsupported catalog must not be scored", x, got.Health)
		}
	}
}
