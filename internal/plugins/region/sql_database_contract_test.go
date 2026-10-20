package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type sqlWireRequest struct {
	SubscriptionID string            `json:"subscription_id"`
	Region         string            `json:"region"`
	ResourceType   string            `json:"resource_type"`
	Aliases        map[string]string `json:"location_aliases"`
}
type sqlWireGroup struct {
	Known        bool   `json:"known"`
	State        int    `json:"state"`
	Version      string `json:"version"`
	Edition      string `json:"edition"`
	Name         string `json:"name"`
	Family       string `json:"family"`
	SKUName      string `json:"sku_name"`
	SKUTier      string `json:"sku_tier"`
	Capacity     int64  `json:"capacity"`
	Standalone   string `json:"standalone"`
	InstancePool string `json:"instance_pool"`
}
type sqlWireCollection struct {
	SubscriptionID string                  `json:"subscription_id"`
	Region         string                  `json:"region"`
	ResourceType   string                  `json:"resource_type"`
	Status         string                  `json:"status"`
	Groups         map[string]sqlWireGroup `json:"groups"`
	FailureCode    string                  `json:"failure_code"`
}
type sqlWireCatalog struct {
	SubscriptionID        string                         `json:"subscription_id"`
	Region                string                         `json:"region"`
	ResourceType          string                         `json:"resource_type"`
	Status                string                         `json:"status"`
	DeploymentEligibility string                         `json:"deployment_eligibility"`
	Values                map[string]sqlWireGroup        `json:"values"`
	Warnings              []assessment.AssessmentWarning `json:"warnings"`
}
type sqlWireCase struct {
	Name  string `json:"name"`
	Input struct {
		Selected map[string]string `json:"selected_subscriptions"`
		Request  sqlWireRequest    `json:"request"`
		Required []string          `json:"required_coordinates"`
		Source   struct {
			Kind  string `json:"kind"`
			Pages []struct {
				Status      int    `json:"status"`
				Body        string `json:"body"`
				BodyFile    string `json:"body_file"`
				BodyPointer string `json:"body_pointer"`
				NilResponse bool   `json:"nil_response"`
				Error       string `json:"error"`
			} `json:"pages"`
			Value *sqlWireCollection `json:"value"`
		} `json:"source"`
		Context      string `json:"context"`
		TerminalCall int    `json:"terminal_getter_call"`
	} `json:"input"`
	Expected struct {
		Collection *sqlWireCollection `json:"collection"`
		Catalog    *sqlWireCatalog    `json:"catalog"`
		Error      *string            `json:"error"`
		Calls      int                `json:"getter_calls"`
		URLs       []string           `json:"request_urls"`
		Unchanged  bool               `json:"input_unchanged"`
		Owned      bool               `json:"owned_output"`
	} `json:"expected"`
}
type sqlWireManifest struct {
	Format       string `json:"format"`
	Purpose      string `json:"purpose"`
	Base         string `json:"base_revision"`
	ResourceType string `json:"resource_type"`
	Primary      []struct {
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
	Cases []sqlWireCase `json:"cases"`
}

func sqlGroupFromWire(w sqlWireGroup) SQLCapabilityGroup {
	return SQLCapabilityGroup{Known: w.Known, State: w.State, Version: w.Version, Edition: w.Edition, Name: w.Name, Family: w.Family, SKUName: w.SKUName, SKUTier: w.SKUTier, Capacity: w.Capacity, Standalone: w.Standalone, InstancePool: w.InstancePool}
}
func sqlCollectionFromWire(w *sqlWireCollection) *SQLCapabilityCollection {
	if w == nil {
		return nil
	}
	c := &SQLCapabilityCollection{SubscriptionID: w.SubscriptionID, Region: w.Region, ResourceType: w.ResourceType, Status: w.Status, FailureCode: w.FailureCode, Groups: map[string]SQLCapabilityGroup{}}
	for key, g := range w.Groups {
		c.Groups[key] = sqlGroupFromWire(g)
	}
	return c
}
func sqlCatalogFromWire(w *sqlWireCatalog) SQLCatalogEvidence {
	if w == nil {
		return SQLCatalogEvidence{}
	}
	c := SQLCatalogEvidence{SubscriptionID: w.SubscriptionID, Region: w.Region, ResourceType: w.ResourceType, Status: w.Status, DeploymentEligibility: w.DeploymentEligibility, Values: map[string]SQLCapabilityGroup{}, Warnings: append([]assessment.AssessmentWarning{}, w.Warnings...)}
	for key, g := range w.Values {
		c.Values[key] = sqlGroupFromWire(g)
	}
	return c
}
func sqlContractCause(s string) error {
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
func sqlRunFrozenManifest(t *testing.T, file, format, resourceType string, count, collectors int, constructor func(string, RESTQuotaGetter) (*SQLCapabilityCollector, error)) {
	t.Helper()
	var m sqlWireManifest
	if err := vmssDecodeContract(vmRead(t, file), &m); err != nil {
		t.Fatal(err)
	}
	if m.Format != format || m.ResourceType != resourceType || len(m.Cases) != count || m.Scoring.Integration != "not_wired" || m.Scoring.Conversion || m.Limitation.Code != "sql_capability_catalog_only" {
		t.Fatal("frozen manifest identity")
	}
	seen := map[string]bool{}
	nc, nd := 0, 0
	for _, f := range m.Cases {
		if f.Name == "" || seen[f.Name] {
			t.Fatal("case identity")
		}
		seen[f.Name] = true
		switch f.Input.Source.Kind {
		case "collector":
			nc++
		case "collection":
			nd++
		default:
			t.Fatal("unknown source")
		}
		t.Run(f.Name, func(t *testing.T) {
			q := SQLCapabilityRequest{SubscriptionID: f.Input.Request.SubscriptionID, Region: f.Input.Request.Region, ResourceType: f.Input.Request.ResourceType, LocationAliases: f.Input.Request.Aliases}
			ctx := &storageTerminalContext{Context: context.Background(), done: make(chan struct{})}
			switch f.Input.Context {
			case "cancelled":
				ctx.finish(context.Canceled)
			case "deadline":
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
					if p.BodyFile != "sql-primary-2021-11-01-example.json" || p.BodyPointer != "/responses/200/body" {
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
					switch f.Input.Context {
					case "cancelled-during-get":
						ctx.finish(context.Canceled)
					case "deadline-during-get":
						ctx.finish(context.DeadlineExceeded)
					}
				}
				var response *http.Response
				if !p.NilResponse {
					response = &http.Response{StatusCode: p.Status}
				}
				err := sqlContractCause(p.Error)
				if err != nil {
					err = fmt.Errorf("private wrapper: %w", err)
				}
				return body, response, err
			})
			before := vmssSnapshot(t, f.Input)
			c := sqlCollectionFromWire(f.Input.Source.Value)
			var err error
			if f.Input.Source.Kind == "collector" {
				collector, e := constructor("https://management.azure.com", getter)
				if e != nil {
					t.Fatal(e)
				}
				c, err = collector.Collect(ctx, f.Input.Selected, q)
				if !reflect.DeepEqual(c, sqlCollectionFromWire(f.Expected.Collection)) {
					t.Fatalf("collection got %#v want %#v", c, sqlCollectionFromWire(f.Expected.Collection))
				}
			}
			collectionBefore := vmssSnapshot(t, c)
			var got SQLCatalogEvidence
			if err == nil {
				got, err = c.Catalog(ctx, f.Input.Selected, q, f.Input.Required)
			}
			if f.Expected.Error == nil {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				e := *f.Expected.Error
				if strings.HasPrefix(e, "context_") {
					if !errors.Is(err, sqlContractCause(e)) {
						t.Fatalf("context: %v", err)
					}
				} else {
					want := map[string]string{"sql_capability_request_invalid": "SQL capability request rejected [sql_capability_request_invalid]", "sql_capability_projection_invalid": "SQL capability catalog projection rejected [sql_capability_projection_invalid]"}[e]
					if want == "" || err == nil || err.Error() != want {
						t.Fatalf("sanitized error: %v want %q", err, want)
					}
				}
			}
			if !reflect.DeepEqual(got, sqlCatalogFromWire(f.Expected.Catalog)) {
				t.Fatalf("catalog got %#v want %#v", got, sqlCatalogFromWire(f.Expected.Catalog))
			}
			if calls != f.Expected.Calls || !reflect.DeepEqual(urls, f.Expected.URLs) {
				t.Fatalf("transport calls/URLs %d/%v", calls, urls)
			}
			if !f.Expected.Unchanged || !bytes.Equal(before, vmssSnapshot(t, f.Input)) || !bytes.Equal(collectionBefore, vmssSnapshot(t, c)) {
				t.Fatal("input mutated")
			}
			if f.Expected.Owned {
				want := vmssSnapshot(t, got)
				for key, g := range got.Values {
					g.SKUName = "changed"
					got.Values[key] = g
					delete(got.Values, key)
				}
				got.Warnings[0].Message = "changed"
				again, e := c.Catalog(ctx, f.Input.Selected, q, f.Input.Required)
				if e != nil || !bytes.Equal(collectionBefore, vmssSnapshot(t, c)) || !bytes.Equal(want, vmssSnapshot(t, again)) {
					t.Fatal("output alias")
				}
				for key, g := range c.Groups {
					g.Name = "changed"
					c.Groups[key] = g
					delete(c.Groups, key)
				}
				if !bytes.Equal(want, vmssSnapshot(t, again)) {
					t.Fatal("input alias")
				}
			}
		})
	}
	if nc != collectors || nd != 29 {
		t.Fatal("case conservation")
	}
	if _, ok := reflect.TypeOf((*SQLCapabilityCollection)(nil)).MethodByName("Evidence"); ok {
		t.Fatal("implicit scoring conversion")
	}
}
func TestSQLDatabaseFrozenManifest(t *testing.T) {
	sqlRunFrozenManifest(t, "sql-database-capability-contract-v1.json", "sql-database-capability-contract-v1", "microsoft.sql/servers/databases", 153, 124, NewSQLDatabaseSKUCollector)
}
func TestSQLFrozenRunnerStrictDecode(t *testing.T) {
	b := vmRead(t, "sql-database-capability-contract-v1.json")
	for _, bad := range [][]byte{append(append([]byte{}, b...), []byte(" {}")...), bytes.Replace(b, []byte(`"known": true`), []byte(`"unknown_field":0,"known": true`), 1)} {
		var m sqlWireManifest
		if vmssDecodeContract(bad, &m) == nil {
			t.Fatal("runner accepted trailing/unknown input")
		}
	}
}
