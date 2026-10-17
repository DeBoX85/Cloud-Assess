package region

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
)

const vmssBoundaryID = "abcdefab-1234-4234-8234-123456789abc"
const vmssBoundaryType = "microsoft.compute/virtualmachinescalesets"

func vmssBoundaryInput() (map[string]string, ComputeVMSKURequest, *ComputeVMSKUCollection) {
	return map[string]string{vmssBoundaryID: "selected"}, ComputeVMSKURequest{SubscriptionID: vmssBoundaryID, Region: "eastus", ResourceType: vmssBoundaryType}, &ComputeVMSKUCollection{SubscriptionID: vmssBoundaryID, Region: "eastus", ResourceType: vmssBoundaryType, Status: "complete", Groups: map[string]ComputeVMSKUGroup{}}
}

func vmssRequireLimit(t *testing.T, got VMSSCatalogEvidence, err error) {
	t.Helper()
	if !reflect.DeepEqual(got, VMSSCatalogEvidence{}) || err == nil || err.Error() != "VMSS catalog projection rejected [vmss_catalog_projection_limit]" {
		t.Fatal("zero result and fixed limit error required", got, err)
	}
}

func TestVMSSCatalogCountBoundaries(t *testing.T) {
	for _, n := range []int{1000, 1001} {
		scope, q, c := vmssBoundaryInput()
		for i := 1; i < n; i++ {
			scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "selected"
		}
		got, err := c.VMSSCatalog(context.Background(), scope, q, nil)
		if n == 1000 {
			if err != nil || got.Status != "complete" {
				t.Fatal("exact scope", got, err)
			}
		} else {
			vmssRequireLimit(t, got, err)
		}
	}
	for _, n := range []int{8192, 8193} {
		scope, q, c := vmssBoundaryInput()
		required := make([]string, n)
		for i := range required {
			required[i] = fmt.Sprintf("s%04d", i)
			c.Groups[required[i]] = ComputeVMSKUGroup{Known: true, State: 0}
		}
		got, err := c.VMSSCatalog(context.Background(), scope, q, nil)
		if n == 8192 {
			if err != nil || got.Status != "complete" {
				t.Fatal("exact groups", got, err)
			}
		} else {
			vmssRequireLimit(t, got, err)
		}
		c.Groups = map[string]ComputeVMSKUGroup{}
		got, err = c.VMSSCatalog(context.Background(), scope, q, required)
		if n == 8192 {
			if err != nil || got.Status != "complete" || len(got.Values) != 8192 || got.Values["s8191"].State != 2 {
				t.Fatal("exact required", got.Status, len(got.Values), err)
			}
		} else {
			vmssRequireLimit(t, got, err)
		}
	}
}

func TestVMSSCatalogLabelBoundaries(t *testing.T) {
	for _, field := range []string{"label", "required", "group", "zone"} {
		for _, n := range []int{512, 513} {
			t.Run(fmt.Sprintf("%s-%d", field, n), func(t *testing.T) {
				scope, q, c := vmssBoundaryInput()
				required := []string{}
				text := strings.Repeat("s", n)
				switch field {
				case "label":
					scope[vmssBoundaryID] = text
				case "required":
					required = []string{text}
				case "group":
					c.Groups[text] = ComputeVMSKUGroup{Known: true, State: 0}
				case "zone":
					c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: []string{strings.Repeat("1", n)}}
					required = []string{"s"}
				}
				got, err := c.VMSSCatalog(context.Background(), scope, q, required)
				if n == 512 {
					if err != nil || got.Status != "complete" {
						t.Fatal("exact label", got, err)
					}
				} else {
					vmssRequireLimit(t, got, err)
				}
			})
		}
	}
}

func TestVMSSCatalogWorkBoundary(t *testing.T) {
	// One selected subscription + one group + one required name +65533zones
	// is exactly65536logical occurrences; all names and zones are valid.
	for _, n := range []int{65533, 65534} {
		scope, q, c := vmssBoundaryInput()
		zones := make([]string, n)
		for i := range zones {
			zones[i] = fmt.Sprint(100000 + i)
		}
		c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: zones}
		got, err := c.VMSSCatalog(context.Background(), scope, q, []string{"s"})
		if n == 65533 {
			if err != nil || len(got.Values["s"].BlockedZones) != 65533 {
				t.Fatal("exact work", got.Status, err)
			}
		} else {
			vmssRequireLimit(t, got, err)
		}
	}
}

func vmssZoneBytes(n int) []string {
	// Unique positive decimals, each <=512bytes and strictly lexically sorted.
	zones := make([]string, 0, n/512+1)
	for i := 0; n >= 512; i++ {
		zones = append(zones, fmt.Sprint(10000000+i)+strings.Repeat("1", 504))
		n -= 512
	}
	if n > 0 {
		zones = append(zones, strings.Repeat("2", n))
	}
	sort.Strings(zones)
	return zones
}

func TestVMSSCatalogTextBoundaries(t *testing.T) {
	// Independently counted literal wire-independent text, no production budget:
	// selected36+8; expected36+6+41; collection36+6+41+8;
	// group key1 =>219bytes with no required name. Cap is16711680.
	for _, n := range []int{16711680 - 219, 16711681 - 219} {
		scope, q, c := vmssBoundaryInput()
		c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: vmssZoneBytes(n)}
		got, err := c.VMSSCatalog(context.Background(), scope, q, nil)
		if n == 16711680-219 {
			if err != nil || got.Status != "complete" || len(got.Values) != 0 {
				t.Fatal("exact input text", got.Status, err)
			}
		} else {
			vmssRequireLimit(t, got, err)
		}
	}
	// Output binding/status/eligibility/warning/key =257bytes; input with
	// required s =220bytes. This isolates output overrun while input is valid.
	for _, n := range []int{16711680 - 257, 16711681 - 257} {
		scope, q, c := vmssBoundaryInput()
		c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: vmssZoneBytes(n)}
		got, err := c.VMSSCatalog(context.Background(), scope, q, []string{"s"})
		if n == 16711680-257 {
			if err != nil || got.Status != "complete" || len(got.Values["s"].BlockedZones) == 0 {
				t.Fatal("exact output text", got.Status, err)
			}
		} else {
			vmssRequireLimit(t, got, err)
		}
	}
	// Unknown required knowledge makes the entire output unknown, including
	// when another required group's unused output would exceed the output cap.
	scope, q, c := vmssBoundaryInput()
	c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: vmssZoneBytes(16711681 - 257)}
	c.Groups["u"] = ComputeVMSKUGroup{Known: false, State: 2}
	got, err := c.VMSSCatalog(context.Background(), scope, q, []string{"s", "u"})
	if err != nil || got.Status != "unknown" || len(got.Values) != 0 {
		t.Fatal("required uncertainty before unused output limits", got.Status, err)
	}
}

func TestVMSSCatalogValidationAndSanitization(t *testing.T) {
	for _, alter := range []func(map[string]string, *ComputeVMSKURequest, *ComputeVMSKUCollection){
		func(s map[string]string, _ *ComputeVMSKURequest, _ *ComputeVMSKUCollection) {
			s[vmssBoundaryID] = "sensitive-canary\n"
		},
		func(_ map[string]string, _ *ComputeVMSKURequest, c *ComputeVMSKUCollection) {
			c.FailureCode = "sensitive-canary https://example.invalid/token"
		},
		func(_ map[string]string, _ *ComputeVMSKURequest, c *ComputeVMSKUCollection) {
			c.Groups[" sensitive-canary"] = ComputeVMSKUGroup{Known: true, State: 0}
		},
	} {
		scope, q, c := vmssBoundaryInput()
		alter(scope, &q, c)
		got, err := c.VMSSCatalog(context.Background(), scope, q, nil)
		if !reflect.DeepEqual(got, VMSSCatalogEvidence{}) || err == nil || err.Error() != "VMSS catalog projection rejected [vmss_catalog_projection_invalid]" {
			t.Fatal("all unrequired input validated; error sanitized", got, err)
		}
	}
	scope, q, c := vmssBoundaryInput()
	c.Groups["standard a1"] = ComputeVMSKUGroup{Known: true, State: 0}
	got, err := c.VMSSCatalog(context.Background(), scope, q, []string{"STANDARD A1"})
	if err != nil || len(got.Values) != 1 || got.Values["standard a1"].State != 0 {
		t.Fatal("existing internal ASCII-space policy", got, err)
	}
}

func TestVMSSCatalogScopeIsolation(t *testing.T) {
	for _, mode := range []string{"unselected-nil", "unselected-empty-label"} {
		scope, q, _ := vmssBoundaryInput()
		if mode == "unselected-nil" {
			delete(scope, vmssBoundaryID)
			scope["bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"] = "other"
		} else {
			scope["bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"] = ""
		}
		var c *ComputeVMSKUCollection
		got, err := c.VMSSCatalog(context.Background(), scope, q, nil)
		if !reflect.DeepEqual(got, VMSSCatalogEvidence{}) || err == nil || err.Error() != "VMSS catalog projection rejected [vmss_catalog_projection_invalid]" {
			t.Fatal("scope validation independent of collection binding", mode, got, err)
		}
	}
}

type vmssCountingContext struct {
	context.Context
	calls, terminalAt int
	cause             error
}

func (c *vmssCountingContext) Err() error {
	c.calls++
	if c.terminalAt > 0 && c.calls >= c.terminalAt {
		return c.cause
	}
	return nil
}

func TestVMSSCatalogFinalReturnContext(t *testing.T) {
	// Fixed late-return schedule, deliberately not recalibrated from the API
	// being tested. Otherwise removing the final check shortens the sweep and
	// silently removes the very cancellation observation it should verify.
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		scope, q, _ := vmssBoundaryInput()
		var c *ComputeVMSKUCollection
		ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 9, cause: cause}
		got, err := c.VMSSCatalog(ctx, scope, q, nil)
		if err != cause || !reflect.DeepEqual(got, VMSSCatalogEvidence{}) {
			t.Fatal("fixed final-return context priority", got, err, ctx.calls)
		}
	}
}

func TestVMSSCatalogLateContext(t *testing.T) {
	for _, mode := range []string{"complete", "unknown", "nil", "invalid", "limit"} {
		for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
			scope, q, c := vmssBoundaryInput()
			required := []string{"s"}
			c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: []string{"1", "2"}}
			switch mode {
			case "unknown":
				c.Status = "unknown"
				c.Groups = nil
			case "nil":
				c = nil
			case "invalid":
				q.Region = "INVALID"
			case "limit":
				scope[vmssBoundaryID] = strings.Repeat("x", 513)
			}
			baseline := &vmssCountingContext{Context: context.Background()}
			_, _ = c.VMSSCatalog(baseline, scope, q, required)
			if baseline.calls < 2 {
				t.Fatal("no late-context observation", mode)
			}
			// Sweep every observation, including the last success/error return.
			for at := 1; at <= baseline.calls; at++ {
				ctx := &vmssCountingContext{Context: context.Background(), terminalAt: at, cause: cause}
				got, err := c.VMSSCatalog(ctx, scope, q, required)
				if err != cause || !reflect.DeepEqual(got, VMSSCatalogEvidence{}) {
					t.Fatal("terminal context must win", mode, at, got, err)
				}
			}
		}
	}
}

func TestVMSSCatalogOwnedConcurrentOutput(t *testing.T) {
	scope, q, c := vmssBoundaryInput()
	c.Groups["s"] = ComputeVMSKUGroup{Known: true, State: 3, BlockedZones: []string{"1", "2"}}
	required := []string{"s"}
	before := vmssSnapshot(t, []any{scope, q, c, required})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a, err := c.VMSSCatalog(context.Background(), scope, q, required)
			if err != nil {
				t.Error(err)
				return
			}
			a.Values["s"].BlockedZones[0] = "99"
			delete(a.Values, "s")
			a.Warnings[0].Code = "changed"
			b, err := c.VMSSCatalog(context.Background(), scope, q, required)
			if err != nil || !reflect.DeepEqual(b.Values["s"].BlockedZones, []string{"1", "2"}) || b.Warnings[0].Code != "vmss_catalog_only" {
				t.Error("shared output", b, err)
			}
		}()
	}
	wg.Wait()
	if !reflect.DeepEqual(before, vmssSnapshot(t, []any{scope, q, c, required})) {
		t.Fatal("output mutation changed inputs")
	}
}
