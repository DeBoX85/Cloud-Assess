package region

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
)

const cognitiveBoundaryID = "abcdefab-1234-4234-8234-123456789abc"
const cognitiveBoundaryType = "microsoft.cognitiveservices/accounts"

func cognitiveBoundaryInput() (map[string]string, CognitiveSKURequest, *CognitiveSKUCollection) {
	return map[string]string{cognitiveBoundaryID: "selected"}, CognitiveSKURequest{SubscriptionID: cognitiveBoundaryID, Region: "eastus", Kind: "openai"}, &CognitiveSKUCollection{SubscriptionID: cognitiveBoundaryID, Region: "eastus", Kind: "openai", ResourceType: cognitiveBoundaryType, Status: "complete", Groups: map[string]CognitiveSKUGroup{}}
}

func cognitiveRequireLimit(t *testing.T, got CognitiveCatalogEvidence, err error) {
	t.Helper()
	if !reflect.DeepEqual(got, CognitiveCatalogEvidence{}) || err == nil || err.Error() != "Cognitive SKU catalog projection rejected [cognitive_sku_projection_limit]" {
		t.Fatal("zero result and fixed limit error required", got.Status, len(got.Values), err)
	}
}

func TestCognitiveCatalogCountBoundaries(t *testing.T) {
	for _, n := range []int{1000, 1001} {
		scope, q, c := cognitiveBoundaryInput()
		for i := 1; i < n; i++ {
			scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "selected"
		}
		got, err := c.Catalog(context.Background(), scope, q, nil)
		if n == 1000 {
			if err != nil || got.Status != "complete" {
				t.Fatal("exact scope", got, err)
			}
		} else {
			cognitiveRequireLimit(t, got, err)
		}
	}
	for _, n := range []int{8192, 8193} {
		scope, q, c := cognitiveBoundaryInput()
		required := make([]string, n)
		for i := range required {
			required[i] = fmt.Sprintf("s%04d", i)
			c.Groups[required[i]] = CognitiveSKUGroup{Known: true, State: 0}
		}
		got, err := c.Catalog(context.Background(), scope, q, nil)
		if n == 8192 {
			if err != nil || got.Status != "complete" {
				t.Fatal("exact groups", got, err)
			}
		} else {
			cognitiveRequireLimit(t, got, err)
		}
		c.Groups = map[string]CognitiveSKUGroup{}
		got, err = c.Catalog(context.Background(), scope, q, required)
		if n == 8192 {
			if err != nil || got.Status != "complete" || len(got.Values) != 8192 || got.Values["s8191"].State != 2 {
				t.Fatal("exact required", got.Status, len(got.Values), err)
			}
		} else {
			cognitiveRequireLimit(t, got, err)
		}
	}
}

func TestCognitiveCatalogLabelBoundaries(t *testing.T) {
	for _, field := range []string{"label", "required", "group", "zone"} {
		for _, n := range []int{512, 513} {
			t.Run(fmt.Sprintf("%s-%d", field, n), func(t *testing.T) {
				scope, q, c := cognitiveBoundaryInput()
				required := []string{}
				text := strings.Repeat("s", n)
				switch field {
				case "label":
					scope[cognitiveBoundaryID] = text
				case "required":
					required = []string{text}
				case "group":
					c.Groups[text] = CognitiveSKUGroup{Known: true, State: 0}
				case "zone":
					c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: []string{strings.Repeat("1", n)}}
					required = []string{"s"}
				}
				got, err := c.Catalog(context.Background(), scope, q, required)
				if n == 512 {
					if err != nil || got.Status != "complete" {
						t.Fatal("exact label", got, err)
					}
				} else {
					cognitiveRequireLimit(t, got, err)
				}
			})
		}
	}
}

func TestCognitiveCatalogWorkBoundary(t *testing.T) {
	// One selected subscription + one group + one required name +65533zones
	// is exactly65536logical occurrences; all names and zones are valid.
	for _, n := range []int{65533, 65534} {
		scope, q, c := cognitiveBoundaryInput()
		zones := make([]string, n)
		for i := range zones {
			zones[i] = fmt.Sprint(100000 + i)
		}
		c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: zones}
		got, err := c.Catalog(context.Background(), scope, q, []string{"s"})
		if n == 65533 {
			if err != nil || len(got.Values["s"].BlockedZones) != 65533 {
				t.Fatal("exact work", got.Status, err)
			}
		} else {
			cognitiveRequireLimit(t, got, err)
		}
	}
}

func TestCognitiveCatalogTextBoundaries(t *testing.T) {
	// Independently counted literal wire-independent text, no production budget:
	// selected36+8; expected36+6+6; collection36+6+6+36+8;
	// group key1 =>185bytes with no required name. Cap is16711680.
	for _, n := range []int{16711680 - 185, 16711681 - 185} {
		scope, q, c := cognitiveBoundaryInput()
		c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: vmssZoneBytes(n)}
		got, err := c.Catalog(context.Background(), scope, q, nil)
		if n == 16711680-185 {
			if err != nil || got.Status != "complete" || len(got.Values) != 0 {
				t.Fatal("exact input text", got.Status, err)
			}
		} else {
			cognitiveRequireLimit(t, got, err)
		}
	}
	// Output binding/status/eligibility/warning/key =248bytes; input with
	// required s =186bytes. This isolates output overrun while input is valid.
	for _, n := range []int{16711680 - 248, 16711681 - 248} {
		scope, q, c := cognitiveBoundaryInput()
		c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: vmssZoneBytes(n)}
		got, err := c.Catalog(context.Background(), scope, q, []string{"s"})
		if n == 16711680-248 {
			if err != nil || got.Status != "complete" || len(got.Values["s"].BlockedZones) == 0 {
				t.Fatal("exact output text", got.Status, err)
			}
		} else {
			cognitiveRequireLimit(t, got, err)
		}
	}
	// Unknown required knowledge makes the entire output unknown, including
	// when another required group's unused output would exceed the output cap.
	scope, q, c := cognitiveBoundaryInput()
	c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: vmssZoneBytes(16711681 - 248)}
	c.Groups["u"] = CognitiveSKUGroup{Known: false, State: 2}
	got, err := c.Catalog(context.Background(), scope, q, []string{"s", "u"})
	if err != nil || got.Status != "unknown" || len(got.Values) != 0 {
		t.Fatal("required uncertainty before unused output limits", got.Status, err)
	}
}

func TestCognitiveCatalogValidationAndSanitization(t *testing.T) {
	for _, alter := range []func(map[string]string, *CognitiveSKURequest, *CognitiveSKUCollection){
		func(s map[string]string, _ *CognitiveSKURequest, _ *CognitiveSKUCollection) {
			s[cognitiveBoundaryID] = "sensitive-canary\n"
		},
		func(_ map[string]string, _ *CognitiveSKURequest, c *CognitiveSKUCollection) {
			c.FailureCode = "sensitive-canary https://example.invalid/token"
		},
		func(_ map[string]string, _ *CognitiveSKURequest, c *CognitiveSKUCollection) {
			c.Groups[" sensitive-canary"] = CognitiveSKUGroup{Known: true, State: 0}
		},
	} {
		scope, q, c := cognitiveBoundaryInput()
		alter(scope, &q, c)
		got, err := c.Catalog(context.Background(), scope, q, nil)
		if !reflect.DeepEqual(got, CognitiveCatalogEvidence{}) || err == nil || err.Error() != "Cognitive SKU catalog projection rejected [cognitive_sku_projection_invalid]" {
			t.Fatal("all unrequired input validated; error sanitized", got, err)
		}
	}
	scope, q, c := cognitiveBoundaryInput()
	c.Groups["standard a1"] = CognitiveSKUGroup{Known: true, State: 0}
	got, err := c.Catalog(context.Background(), scope, q, []string{"STANDARD A1"})
	if err != nil || len(got.Values) != 1 || got.Values["standard a1"].State != 0 {
		t.Fatal("existing internal ASCII-space policy", got, err)
	}
}

func TestCognitiveCatalogScopeIsolation(t *testing.T) {
	for _, mode := range []string{"unselected-nil", "unselected-empty-label"} {
		scope, q, _ := cognitiveBoundaryInput()
		if mode == "unselected-nil" {
			delete(scope, cognitiveBoundaryID)
			scope["bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"] = "other"
		} else {
			scope["bbbbbbbb-2222-4222-8222-bbbbbbbbbbbb"] = ""
		}
		var c *CognitiveSKUCollection
		got, err := c.Catalog(context.Background(), scope, q, nil)
		if !reflect.DeepEqual(got, CognitiveCatalogEvidence{}) || err == nil || err.Error() != "Cognitive SKU catalog projection rejected [cognitive_sku_projection_invalid]" {
			t.Fatal("scope validation independent of collection binding", mode, got, err)
		}
	}
}

func TestCognitiveCatalogFinalReturnContext(t *testing.T) {
	// Fixed late-return schedule, deliberately not recalibrated from the API
	// being tested. Otherwise removing the final check shortens the sweep and
	// silently removes the very cancellation observation it should verify.
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		scope, q, _ := cognitiveBoundaryInput()
		var c *CognitiveSKUCollection
		ctx := &vmssCountingContext{Context: context.Background(), terminalAt: 9, cause: cause}
		got, err := c.Catalog(ctx, scope, q, nil)
		if err != cause || !reflect.DeepEqual(got, CognitiveCatalogEvidence{}) {
			t.Fatal("fixed final-return context priority", got, err, ctx.calls)
		}
	}
}

func TestCognitiveCatalogLateContext(t *testing.T) {
	for _, mode := range []string{"complete", "unknown", "nil", "invalid", "limit"} {
		for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
			scope, q, c := cognitiveBoundaryInput()
			required := []string{"s"}
			c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: []string{"1", "2"}}
			switch mode {
			case "unknown":
				c.Status = "unknown"
				c.Groups = nil
			case "nil":
				c = nil
			case "invalid":
				q.Region = "INVALID"
			case "limit":
				scope[cognitiveBoundaryID] = strings.Repeat("x", 513)
			}
			baseline := &vmssCountingContext{Context: context.Background()}
			_, _ = c.Catalog(baseline, scope, q, required)
			if baseline.calls < 2 {
				t.Fatal("no late-context observation", mode)
			}
			// Sweep every observation, including the last success/error return.
			for at := 1; at <= baseline.calls; at++ {
				ctx := &vmssCountingContext{Context: context.Background(), terminalAt: at, cause: cause}
				got, err := c.Catalog(ctx, scope, q, required)
				if err != cause || !reflect.DeepEqual(got, CognitiveCatalogEvidence{}) {
					t.Fatal("terminal context must win", mode, at, got, err)
				}
			}
		}
	}
}

func TestCognitiveCatalogOwnedConcurrentOutput(t *testing.T) {
	scope, q, c := cognitiveBoundaryInput()
	c.Groups["s"] = CognitiveSKUGroup{Known: true, State: 3, BlockedZones: []string{"1", "2"}}
	required := []string{"s"}
	before := vmssSnapshot(t, []any{scope, q, c, required})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a, err := c.Catalog(context.Background(), scope, q, required)
			if err != nil {
				t.Error(err)
				return
			}
			a.Values["s"].BlockedZones[0] = "99"
			delete(a.Values, "s")
			a.Warnings[0].Code = "changed"
			b, err := c.Catalog(context.Background(), scope, q, required)
			if err != nil || !reflect.DeepEqual(b.Values["s"].BlockedZones, []string{"1", "2"}) || b.Warnings[0].Code != "cognitive_sku_catalog_only" {
				t.Error("shared output", b, err)
			}
		}()
	}
	wg.Wait()
	if !reflect.DeepEqual(before, vmssSnapshot(t, []any{scope, q, c, required})) {
		t.Fatal("output mutation changed inputs")
	}
}
