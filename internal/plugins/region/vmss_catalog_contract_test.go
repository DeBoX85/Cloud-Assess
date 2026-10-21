package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// Wire types deliberately name every underscored field. Decoding into production
// structs would silently drop blocked_zones and other independent expectations.
type vmssContractAvailability struct {
	State        int      `json:"state"`
	BlockedZones []string `json:"blocked_zones"`
}
type vmssContractCatalog struct {
	SubscriptionID        string                              `json:"subscription_id"`
	Region                string                              `json:"region"`
	TargetResourceType    string                              `json:"target_resource_type"`
	SourceResourceType    string                              `json:"source_resource_type"`
	Status                string                              `json:"status"`
	DeploymentEligibility string                              `json:"deployment_eligibility"`
	Values                map[string]vmssContractAvailability `json:"values"`
	Warnings              []assessment.AssessmentWarning      `json:"warnings"`
}
type vmssContractCollection struct {
	SubscriptionID string `json:"subscription_id"`
	Region         string `json:"region"`
	ResourceType   string `json:"resource_type"`
	Status         string `json:"status"`
	Groups         map[string]struct {
		Known        bool     `json:"known"`
		State        int      `json:"state"`
		BlockedZones []string `json:"blocked_zones"`
	} `json:"groups"`
	FailureCode string `json:"failure_code"`
}
type vmssContractCase struct {
	Name  string `json:"name"`
	Input struct {
		Selected map[string]string `json:"selected_subscriptions"`
		Request  struct {
			SubscriptionID string `json:"subscription_id"`
			Region         string `json:"region"`
			ResourceType   string `json:"resource_type"`
		} `json:"request"`
		Required []string `json:"required_names"`
		Source   struct {
			Kind  string `json:"kind"`
			Pages []struct {
				Status int    `json:"status"`
				Body   string `json:"body"`
			} `json:"pages"`
			Value *vmssContractCollection `json:"value"`
		} `json:"source"`
		Context string `json:"context"`
	} `json:"input"`
	Expected struct {
		Catalog        *vmssContractCatalog `json:"catalog"`
		Error          *string              `json:"error"`
		Calls          int                  `json:"getter_calls"`
		URLs           []string             `json:"request_urls"`
		InputUnchanged bool                 `json:"input_unchanged"`
		OwnedOutput    bool                 `json:"owned_output"`
	} `json:"expected"`
}
type vmssContractManifest struct {
	Format       string                       `json:"format"`
	Purpose      string                       `json:"purpose"`
	BaseRevision string                       `json:"base_revision"`
	SourceType   string                       `json:"catalog_source_resource_type"`
	TargetType   string                       `json:"target_resource_type"`
	Limitation   assessment.AssessmentWarning `json:"limitation"`
	Scoring      struct {
		API      string                              `json:"api"`
		Status   string                              `json:"valid_complete_vmss_status"`
		Values   map[string]vmssContractAvailability `json:"valid_complete_vmss_values"`
		Consumer struct {
			Total     int                    `json:"total_skus_checked"`
			Unknown   int                    `json:"unknown_skus"`
			Available int                    `json:"available_skus"`
			Percent   float64                `json:"sku_availability_percent"`
			Health    assessment.StageStatus `json:"health"`
			Warning   string                 `json:"warning_code"`
		} `json:"consumer"`
	} `json:"scoring_invariant"`
	Cases []vmssContractCase `json:"cases"`
}

func vmssDecodeContract(b []byte, dst any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err := d.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("contract has trailing JSON")
	}
	return nil
}

func vmssLoadContract(t *testing.T) vmssContractManifest {
	t.Helper()
	var m vmssContractManifest
	if err := vmssDecodeContract(vmRead(t, "vmss-catalog-target-contract-v1.json"), &m); err != nil {
		t.Fatal(err)
	}
	if m.Format != "vmss-catalog-target-contract-v1" || len(m.Cases) != 61 {
		t.Fatal("frozen contract identity/count")
	}
	return m
}

func vmssSnapshot(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func vmssExpectedCatalog(w *vmssContractCatalog) VMSSCatalogEvidence {
	if w == nil {
		return VMSSCatalogEvidence{}
	}
	v := make(map[string]SKUAvailability, len(w.Values))
	for k, x := range w.Values {
		v[k] = SKUAvailability{State: x.State, BlockedZones: x.BlockedZones}
	}
	return VMSSCatalogEvidence{SubscriptionID: w.SubscriptionID, Region: w.Region, TargetResourceType: w.TargetResourceType, SourceResourceType: w.SourceResourceType, Status: w.Status, DeploymentEligibility: w.DeploymentEligibility, Values: v, Warnings: w.Warnings}
}

func TestVMSSCatalogFrozenManifest(t *testing.T) {
	m := vmssLoadContract(t)
	seen := map[string]bool{}
	collectors, direct := 0, 0
	for _, f := range m.Cases {
		if f.Name == "" || seen[f.Name] {
			t.Fatal("empty/duplicate case")
		}
		seen[f.Name] = true
		if f.Input.Source.Kind == "collector" {
			collectors++
		} else if f.Input.Source.Kind == "collection" {
			direct++
		} else {
			t.Fatal("unknown source kind")
		}
		t.Run(f.Name, func(t *testing.T) {
			q := ComputeVMSKURequest{SubscriptionID: f.Input.Request.SubscriptionID, Region: f.Input.Request.Region, ResourceType: f.Input.Request.ResourceType}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			switch f.Input.Context {
			case "active":
			case "cancelled":
				cancel()
			case "deadline":
				var end context.CancelFunc
				ctx, end = context.WithDeadline(context.Background(), time.Unix(1, 0))
				defer end()
			default:
				t.Fatal("unknown context")
			}
			var c *ComputeVMSKUCollection
			calls := 0
			if f.Input.Source.Kind == "collector" {
				collector, err := NewComputeVMSKUCollector("https://management.azure.com", vmGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
					if calls >= len(f.Input.Source.Pages) || calls >= len(f.Expected.URLs) {
						t.Fatal("unexpected getter call", calls)
					}
					if limit != 1048576 || u != f.Expected.URLs[calls] {
						t.Fatalf("literal URL/bound mismatch %q/%d", u, limit)
					}
					p := f.Input.Source.Pages[calls]
					calls++
					return []byte(p.Body), &http.Response{StatusCode: p.Status}, nil
				}))
				if err != nil {
					t.Fatal(err)
				}
				c, err = collector.Collect(ctx, f.Input.Selected, q)
				if err != nil {
					t.Fatal("collector preparation", err)
				}
			} else if w := f.Input.Source.Value; w != nil {
				c = &ComputeVMSKUCollection{SubscriptionID: w.SubscriptionID, Region: w.Region, ResourceType: w.ResourceType, Status: w.Status, FailureCode: w.FailureCode, Groups: make(map[string]ComputeVMSKUGroup, len(w.Groups))}
				for k, x := range w.Groups {
					c.Groups[k] = ComputeVMSKUGroup{Known: x.Known, State: x.State, BlockedZones: x.BlockedZones}
				}
			}
			inputs := []any{f.Input.Selected, q, f.Input.Required, c}
			before := vmssSnapshot(t, inputs)
			got, err := c.VMSSCatalog(ctx, f.Input.Selected, q, f.Input.Required)
			if !f.Expected.InputUnchanged || !bytes.Equal(before, vmssSnapshot(t, inputs)) {
				t.Fatal("projection mutated input")
			}
			if calls != f.Expected.Calls || calls != len(f.Expected.URLs) {
				t.Fatal("getter count", calls, f.Expected.Calls)
			}
			want := vmssExpectedCatalog(f.Expected.Catalog)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("literal full catalog mismatch: %#v want %#v", got, want)
			}
			if f.Expected.Error == nil {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				switch *f.Expected.Error {
				case "context_canceled":
					if err != context.Canceled {
						t.Fatal("canonical cancellation", err)
					}
				case "context_deadline_exceeded":
					if err != context.DeadlineExceeded {
						t.Fatal("canonical deadline", err)
					}
				case "vmss_catalog_projection_invalid":
					if err == nil || err.Error() != "VMSS catalog projection rejected [vmss_catalog_projection_invalid]" {
						t.Fatal("sanitized invalid error", err)
					}
				default:
					t.Fatal("unknown expected error")
				}
			}
			if f.Expected.OwnedOutput {
				if c == nil || len(got.Values) == 0 {
					t.Fatal("ownership case missing data")
				}
				pristine := vmssSnapshot(t, got)
				for k, g := range c.Groups {
					if len(g.BlockedZones) > 0 {
						g.BlockedZones[0] = "999"
					}
					c.Groups[k] = ComputeVMSKUGroup{Known: false, State: 2}
				}
				for k := range f.Input.Selected {
					f.Input.Selected[k] = "changed"
				}
				for i := range f.Input.Required {
					f.Input.Required[i] = "changed"
				}
				if !bytes.Equal(pristine, vmssSnapshot(t, got)) {
					t.Fatal("caller mutation changed output")
				}
				for k, v := range got.Values {
					if len(v.BlockedZones) > 0 {
						v.BlockedZones[0] = "888"
					}
					delete(got.Values, k)
				}
				got.Warnings[0].Message = "changed"
				// A fresh call from a fresh wire reconstruction must keep its own warning/map.
				for _, other := range m.Cases {
					if other.Name == "direct-known-catalog" {
						w := other.Input.Source.Value
						r := &ComputeVMSKUCollection{SubscriptionID: w.SubscriptionID, Region: w.Region, ResourceType: w.ResourceType, Status: w.Status, Groups: map[string]ComputeVMSKUGroup{"standard_a1": {Known: true, State: 0}}}
						next, e := r.VMSSCatalog(context.Background(), other.Input.Selected, ComputeVMSKURequest{SubscriptionID: other.Input.Request.SubscriptionID, Region: other.Input.Request.Region, ResourceType: other.Input.Request.ResourceType}, other.Input.Required)
						if e != nil || !reflect.DeepEqual(next, vmssExpectedCatalog(other.Expected.Catalog)) {
							t.Fatal("shared output storage", next, e)
						}
					}
				}
			}
		})
	}
	if collectors != 25 || direct != 36 {
		t.Fatal("source coverage changed", collectors, direct)
	}
}

func TestVMSSCatalogScoringInvariant(t *testing.T) {
	m := vmssLoadContract(t)
	id := "abcdefab-1234-4234-8234-123456789abc"
	kind := "microsoft.compute/virtualmachinescalesets"
	c := &ComputeVMSKUCollection{SubscriptionID: id, Region: "eastus", ResourceType: kind, Status: "complete", Groups: map[string]ComputeVMSKUGroup{"standard_a1": {Known: true, State: 0}}}
	diagnostic, err := c.VMSSCatalog(context.Background(), map[string]string{id: "selected"}, ComputeVMSKURequest{SubscriptionID: id, Region: "eastus", ResourceType: kind}, []string{"standard_a1"})
	if err != nil || diagnostic.Status != "complete" {
		t.Fatal("diagnostic preparation", diagnostic, err)
	}
	evidence, err := c.Evidence(context.Background(), []string{"standard_a1"})
	if err != nil || evidence.Status != m.Scoring.Status || len(evidence.Values) != len(m.Scoring.Values) || m.Scoring.API != "ComputeVMSKUCollection.Evidence" {
		t.Fatal("unsupported scoring projection", evidence, err)
	}
	inventory := &InventoryCalculation{Subscriptions: map[string]InventoryCounts{id: {ResourceTypesByRegion: map[string]map[string]int64{"westus": {kind: 1}}, SKUsByTypeAndRegion: map[string]map[string]map[string]int64{kind: {"westus": {"standard_a1": 1}}}}}}
	got, err := CalculateAvailability(context.Background(), map[string]string{id: "selected"}, inventory, AvailabilityRequest{SubscriptionID: id, SourceRegion: "westus", TargetRegion: "eastus"}, AvailabilityEvidence{SubscriptionID: id, TargetRegion: "eastus", InventoryComplete: true, ProvidersComplete: true, Locations: map[string]map[string]bool{kind: {"eastus": true}}, SKUEnabled: true, SKUs: map[string]SKUEvidence{kind: evidence}, ZoneCounts: map[string]int{"westus": 3, "eastus": 3}})
	if err != nil || got == nil {
		t.Fatal(got, err)
	}
	x, w := got.Comparison, m.Scoring.Consumer
	if x.TotalSKUsChecked != w.Total || x.UnknownSKUs != w.Unknown || x.AvailableSKUs != w.Available || x.SKUAvailabilityPercent != w.Percent || got.Health.Status != w.Health || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Code != w.Warning {
		t.Fatal("VMSS diagnostic entered score", x, got.Health)
	}
}

func TestVMSSCatalogContractStrictDecode(t *testing.T) {
	b := vmRead(t, "vmss-catalog-target-contract-v1.json")
	for _, bad := range [][]byte{bytes.Replace(b, []byte(`"blocked_zones"`), []byte(`"blocked_zone_typo"`), 1), append(append([]byte{}, b...), []byte(` {}`)...)} {
		var m vmssContractManifest
		if vmssDecodeContract(bad, &m) == nil {
			t.Fatal("malformed/unknown contract field admitted")
		}
	}
}
