package region

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

func TestComputeDiskSKUNullPaginationTerminator(t *testing.T) {
	for _, suffix := range []string{"", `,"nextLink":""`, `,"nextLink":null`} {
		body := strings.TrimSuffix(diskEnvelope(diskAvailableRow, ""), "}") + suffix + "}"
		value, err, getter := diskCollect(t, diskReply{body: body})
		want := map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}
		if err != nil || value == nil || value.Status != "complete" || value.FailureCode != "" || !reflect.DeepEqual(value.Values, want) || len(getter.urls) != 1 {
			t.Fatalf("terminal suffix=%q result=%#v error=%v; want complete literal disk map", suffix, value, err)
		}
	}
	first := diskEnvelope(diskAvailableRow, diskNextURL)
	last := `{"value":[{"resourceType":"disks","name":"Standard_LRS","locations":["eastus"],"restrictions":[]}],"nextLink":null}`
	value, err, getter := diskCollect(t, diskReply{body: first}, diskReply{body: last})
	want := map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}, "standard_lrs": {State: 0, BlockedZones: []string{}}}
	if err != nil || value == nil || value.Status != "complete" || !reflect.DeepEqual(value.Values, want) || !reflect.DeepEqual(value.Evidence(), SKUEvidence{Status: "complete", Values: want}) || len(getter.urls) != 2 {
		t.Fatalf("null final page result=%#v error=%v; want two complete literal disk entries", value, err)
	}
	// Exact primary 2021-07-01 region example body: original rows are VM SKUs,
	// so the disk collector must admit a healthy complete empty disk map.
	const primaryExample = `{"value":[{"resourceType":"virtualMachines","locations":["westus"],"capabilities":[{"name":"MaxResourceVolumeMB","value":"20480"},{"name":"OSVhdSizeMB","value":"1047552"},{"name":"vCPUs","value":"1"},{"name":"HyperVGenerations","value":"V1"},{"name":"MemoryGB","value":"0.75"},{"name":"MaxDataDiskCount","value":"1"},{"name":"LowPriorityCapable","value":"False"},{"name":"PremiumIO","value":"False"},{"name":"vCPUsAvailable","value":"1"},{"name":"ACUs","value":"50"},{"name":"vCPUsPerCore","value":"1"},{"name":"EphemeralOSDiskSupported","value":"False"},{"name":"AcceleratedNetworkingEnabled","value":"False"},{"name":"RdmaEnabled","value":"False"},{"name":"MaxNetworkInterfaces","value":"2"}],"locationInfo":[{"location":"westus","zones":["2","1"],"zoneDetails":[{"Name":["2"],"capabilities":[{"name":"UltraSSDAvailable","value":"True"}]}]}],"name":"Standard_A0","tier":"Standard","size":"A0","family":"standardA0_A7Family"},{"resourceType":"virtualMachines","locations":["westus"],"capabilities":[{"name":"MaxResourceVolumeMB","value":"71680"},{"name":"OSVhdSizeMB","value":"1047552"},{"name":"vCPUs","value":"1"},{"name":"HyperVGenerations","value":"V1"},{"name":"MemoryGB","value":"1.75"},{"name":"MaxDataDiskCount","value":"2"},{"name":"LowPriorityCapable","value":"True"},{"name":"PremiumIO","value":"False"},{"name":"vCPUsAvailable","value":"1"},{"name":"ACUs","value":"100"},{"name":"vCPUsPerCore","value":"1"},{"name":"EphemeralOSDiskSupported","value":"False"},{"name":"AcceleratedNetworkingEnabled","value":"False"},{"name":"RdmaEnabled","value":"False"},{"name":"MaxNetworkInterfaces","value":"2"}],"locationInfo":[{"location":"westus","zones":["1","2","3"]}],"name":"Standard_A1","tier":"Standard","size":"A1","family":"standardA0_A7Family"}],"nextLink":null}`
	value, err, getter = diskCollect(t, diskReply{body: primaryExample})
	if err != nil || value == nil || value.Status != "complete" || len(value.Values) != 0 || value.FailureCode != "" || len(getter.urls) != 1 {
		t.Fatalf("primary null-terminated example result=%#v error=%v; want healthy filtered empty", value, err)
	}
}

func TestComputeDiskSKUHTTPStatusClassification(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusTemporaryRedirect} {
		credential := &diskCanaryCredential{}
		wire := &diskCanaryTransport{t: t, status: status}
		client := azure.NewHTTPClient(credential, &azure.HTTPClientOptions{Timeout: time.Second, MaxRetries: -1, OperationTimeout: 5 * time.Second, Scope: "https://management.azure.com/.default", Transport: wire})
		collector, _ := NewComputeDiskSKUCollector("https://management.azure.com", client)
		value, err := collector.Collect(context.Background(), map[string]string{diskTestID: "Selected"}, ComputeDiskSKURequest{diskTestID, "eastus"})
		if err != nil || value == nil || value.Status != "unknown" || value.FailureCode != "sku_status_invalid" || len(value.Values) != 0 || wire.calls.Load() != 1 || wire.closes.Load() != 1 || credential.calls.Load() != 1 {
			t.Fatalf("actual bounded status=%d result=%#v error=%v calls=%d closes=%d; want sku_status_invalid", status, value, err, wire.calls.Load(), wire.closes.Load())
		}
		value, err, _ = diskCollect(t, diskReply{body: diskEnvelope(diskAvailableRow, diskNextURL)}, diskReply{status: status, err: errors.New("PRIVATE denied status")})
		if err != nil || value == nil || value.Status != "partial" || value.FailureCode != "sku_status_invalid" || !reflect.DeepEqual(value.Values, map[string]SKUAvailability{"premium_lrs": {State: 0, BlockedZones: []string{}}}) || value.Evidence().Status != "unknown" {
			t.Fatalf("later status=%d result=%#v error=%v; want sanitized partial sku_status_invalid", status, value, err)
		}
	}
	value, err, _ := diskCollect(t, diskReply{status: 200, err: errors.New("PRIVATE body read failed")})
	if err != nil || value == nil || value.Status != "unknown" || value.FailureCode != "sku_request_failed" || len(value.Values) != 0 {
		t.Fatalf("successful status with transport failure=%#v error=%v", value, err)
	}
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		value, err, _ := diskCollect(t, diskReply{status: 403, err: cause})
		if value != nil || !errors.Is(err, cause) {
			t.Fatalf("terminal context cause=%v result=%#v error=%v", cause, value, err)
		}
	}
}
