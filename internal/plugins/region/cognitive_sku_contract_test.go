package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type cognitiveWireRequest struct {
	SubscriptionID string            `json:"subscription_id"`
	Region         string            `json:"region"`
	Kind           string            `json:"kind"`
	Aliases        map[string]string `json:"location_aliases"`
}
type cognitiveWireGroup struct {
	Known bool     `json:"known"`
	State int      `json:"state"`
	Zones []string `json:"blocked_zones"`
}
type cognitiveWireCollection struct {
	SubscriptionID string                        `json:"subscription_id"`
	Region         string                        `json:"region"`
	Kind           string                        `json:"kind"`
	ResourceType   string                        `json:"resource_type"`
	Status         string                        `json:"status"`
	Groups         map[string]cognitiveWireGroup `json:"groups"`
	FailureCode    string                        `json:"failure_code"`
}
type cognitiveWireCatalog struct {
	SubscriptionID        string                              `json:"subscription_id"`
	Region                string                              `json:"region"`
	Kind                  string                              `json:"kind"`
	ResourceType          string                              `json:"resource_type"`
	Status                string                              `json:"status"`
	DeploymentEligibility string                              `json:"deployment_eligibility"`
	Values                map[string]vmssContractAvailability `json:"values"`
	Warnings              []assessment.AssessmentWarning      `json:"warnings"`
}
type cognitiveWireCase struct {
	Name  string `json:"name"`
	Input struct {
		Selected map[string]string    `json:"selected_subscriptions"`
		Request  cognitiveWireRequest `json:"request"`
		Required []string             `json:"required_names"`
		Source   struct {
			Kind  string `json:"kind"`
			Pages []struct {
				Status      int    `json:"status"`
				Body        string `json:"body"`
				BodyFile    string `json:"body_file"`
				BodyPointer string `json:"body_pointer"`
				Error       string `json:"error"`
				NilResponse bool   `json:"nil_response"`
			} `json:"pages"`
			Value *cognitiveWireCollection `json:"value"`
		} `json:"source"`
		Context      string `json:"context"`
		TerminalCall int    `json:"terminal_getter_call"`
	} `json:"input"`
	Expected struct {
		Collection *cognitiveWireCollection `json:"collection"`
		Catalog    *cognitiveWireCatalog    `json:"catalog"`
		Error      *string                  `json:"error"`
		Calls      int                      `json:"getter_calls"`
		URLs       []string                 `json:"request_urls"`
		Unchanged  bool                     `json:"input_unchanged"`
		Owned      bool                     `json:"owned_output"`
	} `json:"expected"`
}

func cognitiveFromWire(w *cognitiveWireCollection) *CognitiveSKUCollection {
	if w == nil {
		return nil
	}
	c := &CognitiveSKUCollection{SubscriptionID: w.SubscriptionID, Region: w.Region, Kind: w.Kind, ResourceType: w.ResourceType, Status: w.Status, FailureCode: w.FailureCode, Groups: map[string]CognitiveSKUGroup{}}
	for k, g := range w.Groups {
		c.Groups[k] = CognitiveSKUGroup{Known: g.Known, State: g.State, BlockedZones: append([]string{}, g.Zones...)}
	}
	return c
}
func cognitiveCatalogFromWire(w *cognitiveWireCatalog) CognitiveCatalogEvidence {
	if w == nil {
		return CognitiveCatalogEvidence{}
	}
	c := CognitiveCatalogEvidence{SubscriptionID: w.SubscriptionID, Region: w.Region, Kind: w.Kind, ResourceType: w.ResourceType, Status: w.Status, DeploymentEligibility: w.DeploymentEligibility, Values: map[string]SKUAvailability{}, Warnings: w.Warnings}
	for k, g := range w.Values {
		c.Values[k] = SKUAvailability{State: g.State, BlockedZones: g.BlockedZones}
	}
	return c
}
func cognitiveContractCause(s string) error {
	switch s {
	case "context_canceled":
		return context.Canceled
	case "context_deadline_exceeded":
		return context.DeadlineExceeded
	case "":
		return nil
	default:
		return errors.New("opaque-private-getter-detail")
	}
}
func TestCognitiveSKUFrozenManifest(t *testing.T) {
	var m struct {
		Format  string `json:"format"`
		Purpose string `json:"purpose"`
		Base    string `json:"base_revision"`
		Primary []struct {
			File  string `json:"file"`
			Bytes int    `json:"bytes"`
			SHA   string `json:"sha256"`
			Blob  string `json:"source_blob"`
		} `json:"primary_files"`
		Limitation assessment.AssessmentWarning `json:"limitation"`
		Scoring    struct {
			Integration string `json:"integration"`
			Conversion  bool   `json:"implicit_sku_evidence_conversion"`
		} `json:"scoring_invariant"`
		Cases []cognitiveWireCase `json:"cases"`
	}
	if err := vmssDecodeContract(vmRead(t, "cognitive-sku-target-contract-v1.json"), &m); err != nil {
		t.Fatal(err)
	}
	if m.Format != "cognitive-sku-target-contract-v1" || len(m.Cases) != 143 || m.Scoring.Integration != "not_wired" || m.Scoring.Conversion {
		t.Fatal("frozen manifest identity")
	}
	seen := map[string]bool{}
	collectors, direct := 0, 0
	for _, f := range m.Cases {
		if f.Name == "" || seen[f.Name] {
			t.Fatal("duplicate case")
		}
		seen[f.Name] = true
		if f.Input.Source.Kind == "collector" {
			collectors++
		} else if f.Input.Source.Kind == "collection" {
			direct++
		} else {
			t.Fatal("unknown source")
		}
		t.Run(f.Name, func(t *testing.T) {
			q := CognitiveSKURequest{SubscriptionID: f.Input.Request.SubscriptionID, Region: f.Input.Request.Region, Kind: f.Input.Request.Kind, LocationAliases: f.Input.Request.Aliases}
			ctx := &storageTerminalContext{Context: context.Background(), done: make(chan struct{})}
			if f.Input.Context == "cancelled" {
				ctx.finish(context.Canceled)
			} else if f.Input.Context == "deadline" {
				ctx.finish(context.DeadlineExceeded)
			}
			calls := 0
			urls := []string{}
			getter := storageGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
				calls++
				urls = append(urls, u)
				if limit != 1048576 || calls > len(f.Input.Source.Pages) {
					t.Fatalf("unexpected getter call/limit %d/%d", calls, limit)
				}
				p := f.Input.Source.Pages[calls-1]
				body := []byte(p.Body)
				if p.BodyFile != "" {
					if p.BodyFile != "cognitive-primary-2024-10-01-GetSkus.json" || p.BodyPointer != "/responses/200/body" {
						t.Fatal("unexpected primary pointer")
					}
					var primary struct {
						Responses map[string]struct {
							Body json.RawMessage `json:"body"`
						} `json:"responses"`
					}
					if err := json.Unmarshal(vmRead(t, p.BodyFile), &primary); err != nil {
						t.Fatal(err)
					}
					body = primary.Responses["200"].Body
				}
				if calls == f.Input.TerminalCall {
					if f.Input.Context == "cancelled-during-get" {
						ctx.finish(context.Canceled)
					} else if f.Input.Context == "deadline-during-get" {
						ctx.finish(context.DeadlineExceeded)
					}
				}
				var response *http.Response
				if !p.NilResponse {
					response = &http.Response{StatusCode: p.Status}
				}
				err := cognitiveContractCause(p.Error)
				if err != nil {
					err = fmt.Errorf("private wrapper: %w", err)
				}
				return body, response, err
			})
			before := vmssSnapshot(t, struct {
				Q        CognitiveSKURequest
				Selected map[string]string
				Required []string
				Source   any
			}{q, f.Input.Selected, f.Input.Required, f.Input.Source})
			c := cognitiveFromWire(f.Input.Source.Value)
			var err error
			if f.Input.Source.Kind == "collector" {
				collector, e := NewCognitiveSKUCollector("https://management.azure.com", getter)
				if e != nil {
					t.Fatal(e)
				}
				c, err = collector.Collect(ctx, f.Input.Selected, q)
				if !reflect.DeepEqual(c, cognitiveFromWire(f.Expected.Collection)) {
					t.Fatalf("collection\ngot %#v\nwant %#v", c, cognitiveFromWire(f.Expected.Collection))
				}
			}
			var got CognitiveCatalogEvidence
			if err == nil {
				got, err = c.Catalog(ctx, f.Input.Selected, q, f.Input.Required)
			}
			if f.Expected.Error == nil {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				e := *f.Expected.Error
				cause := cognitiveContractCause(e)
				if strings.HasPrefix(e, "context_") {
					if !errors.Is(err, cause) {
						t.Fatalf("context: %v", err)
					}
				} else if err == nil || err.Error() != map[string]string{"cognitive_sku_request_invalid": "Cognitive SKU request rejected [cognitive_sku_request_invalid]", "cognitive_sku_projection_invalid": "Cognitive SKU catalog projection rejected [cognitive_sku_projection_invalid]"}[e] {
					t.Fatalf("sanitized error: %v", err)
				}
			}
			if !reflect.DeepEqual(got, cognitiveCatalogFromWire(f.Expected.Catalog)) {
				t.Fatalf("catalog\ngot %#v\nwant %#v", got, cognitiveCatalogFromWire(f.Expected.Catalog))
			}
			if calls != f.Expected.Calls || !reflect.DeepEqual(urls, f.Expected.URLs) {
				t.Fatalf("transport calls/URLs %d/%v", calls, urls)
			}
			after := vmssSnapshot(t, struct {
				Q        CognitiveSKURequest
				Selected map[string]string
				Required []string
				Source   any
			}{q, f.Input.Selected, f.Input.Required, f.Input.Source})
			if !f.Expected.Unchanged || !bytes.Equal(before, after) {
				t.Fatal("input mutated")
			}
			if f.Expected.Owned {
				snap := vmssSnapshot(t, c)
				want := vmssSnapshot(t, got)
				for k, v := range got.Values {
					if len(v.BlockedZones) > 0 {
						v.BlockedZones[0] = "999"
					}
					delete(got.Values, k)
				}
				got.Warnings[0].Message = "changed"
				again, e := c.Catalog(ctx, f.Input.Selected, q, f.Input.Required)
				if e != nil || !bytes.Equal(snap, vmssSnapshot(t, c)) || !bytes.Equal(want, vmssSnapshot(t, again)) {
					t.Fatal("output alias")
				}
				for k, g := range c.Groups {
					if len(g.BlockedZones) > 0 {
						g.BlockedZones[0] = "888"
					}
					delete(c.Groups, k)
				}
				if !bytes.Equal(want, vmssSnapshot(t, again)) {
					t.Fatal("input alias")
				}
			}
		})
	}
	if collectors != 108 || direct != 35 {
		t.Fatal("case conservation")
	}
	if _, ok := reflect.TypeOf((*CognitiveSKUCollection)(nil)).MethodByName("Evidence"); ok {
		t.Fatal("implicit scoring conversion")
	}
}
