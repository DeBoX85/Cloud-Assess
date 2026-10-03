package aigov

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const fixtureSub = "11111111-1111-4111-8111-111111111111"
const fixtureID = "/subscriptions/" + fixtureSub + "/resourceGroups/fixture-rg/providers/Microsoft.CognitiveServices/accounts/fixture-ai"

func ptr[T any](v T) *T { return &v }
func fixtureAccounts() []Account {
	return []Account{{ID: fixtureID, SubscriptionID: fixtureSub, ResourceGroup: "fixture-rg", Name: "fixture-ai", Kind: "OpenAI", SKU: "S0"}, {ID: strings.ReplaceAll(fixtureID, "fixture-ai", "fixture-empty"), SubscriptionID: fixtureSub, ResourceGroup: "fixture-rg", Name: "fixture-empty"}}
}
func fixtureScope() map[string]string { return map[string]string{fixtureSub: "Synthetic subscription"} }

// Test-only wire extraction uses the literal pinned SDK schema. Expected cells
// come from separately executed source captures, never from this parser or target.
func fixturePoints(t *testing.T) []Point {
	t.Helper()
	data, e := os.ReadFile("testdata/metrics-input.json")
	if e != nil {
		t.Fatal(e)
	}
	var wire struct {
		Values []struct {
			ID     *string `json:"resourceid"`
			Values []struct {
				Series []struct {
					Meta []struct {
						Name *struct {
							Value *string `json:"value"`
						} `json:"name"`
						Value *string `json:"value"`
					} `json:"metadatavalues"`
					Data []struct {
						Timestamp *time.Time `json:"timeStamp"`
						Count     *float64   `json:"count"`
					} `json:"data"`
				} `json:"timeseries"`
			} `json:"value"`
		} `json:"values"`
	}
	if e = json.Unmarshal(data, &wire); e != nil {
		t.Fatal(e)
	}
	var out []Point
	for _, r := range wire.Values {
		if r.ID == nil {
			continue
		}
		for _, m := range r.Values {
			for _, series := range m.Series {
				var dep, model, status *string
				for _, v := range series.Meta {
					if v.Name == nil || v.Name.Value == nil || v.Value == nil {
						continue
					}
					switch strings.ToLower(*v.Name.Value) {
					case "modeldeploymentname":
						dep = v.Value
					case "modelname":
						model = v.Value
					case "statuscode":
						status = v.Value
					}
				}
				for _, v := range series.Data {
					out = append(out, Point{*r.ID, v.Timestamp, v.Count, dep, model, status})
				}
			}
		}
	}
	return out
}
func fixtureDeployments() map[string]DeploymentSet {
	return map[string]DeploymentSet{fixtureID: {Values: []Deployment{{Name: ptr("fixture-deployment"), ModelVersion: ptr("2026-01-01"), ModelFormat: ptr("OpenAI"), Capacity: ptr(int64(42)), Upgrade: ptr("NoAutoUpgrade"), Spillover: ptr("")}, {Name: ptr("ignored-deployment")}}}}
}
func valid(t *testing.T, table assessment.PluginTable) {
	t.Helper()
	if e := assessment.ValidatePluginTables([]assessment.PluginTable{table}); e != nil {
		t.Fatalf("invalid canonical output: %v", e)
	}
}
func TestProjectActualSource(t *testing.T) {
	for _, kind := range []string{"all", "deployment-denied", "empty", "unknown-include-tag"} {
		t.Run(kind, func(t *testing.T) {
			var capture []struct {
				Metadata struct {
					Name, Version, Description, Author, License string
					Type                                        int
					ColumnMetadata                              []struct {
						Name string `json:"name"`
					}
				} `json:"metadata"`
				Sheet       string     `json:"sheet_name"`
				Description string     `json:"description"`
				Table       [][]string `json:"table"`
			}
			data, e := os.ReadFile("testdata/source-" + kind + ".json")
			if e != nil {
				t.Fatal(e)
			}
			if e = json.Unmarshal(data, &capture); e != nil {
				t.Fatal(e)
			}
			accounts, points, deployments := fixtureAccounts(), fixturePoints(t), fixtureDeployments()
			if kind == "deployment-denied" {
				deployments = map[string]DeploymentSet{fixtureID: {Failed: true}}
			}
			if kind == "empty" || kind == "unknown-include-tag" {
				accounts = nil
				points = nil
				deployments = nil
			}
			table, e := Project(context.Background(), fixtureScope(), accounts, points, deployments, nil)
			if e != nil {
				t.Fatal(e)
			}
			valid(t, table)
			source := capture[0]
			meta := table.Metadata
			if meta.Name != source.Metadata.Name || meta.Version != source.Metadata.Version || meta.Description != source.Metadata.Description || meta.Author != source.Metadata.Author || meta.License != source.Metadata.License || meta.Type != "internal" || source.Metadata.Type != 1 || table.SheetName != source.Sheet || table.Description != source.Description || !reflect.DeepEqual(table.Columns, source.Table[0]) {
				t.Fatal("source metadata/columns/sheet mismatch")
			}
			got := [][]string{}
			for _, row := range table.Rows {
				if row.SubscriptionID != fixtureSub {
					t.Fatal("lost privacy correlation")
				}
				got = append(got, row.Cells)
			}
			if !reflect.DeepEqual(got, source.Table[1:]) {
				t.Fatalf("source cells mismatch\ngot %v\nwant %v", got, source.Table[1:])
			}
			if len(points) > 0 && table.Health.Status != assessment.StageCompletedWithWarnings {
				t.Fatal("missing points concealed")
			}
			if kind == "deployment-denied" && (len(table.Health.Warnings) != 2 || table.Health.Warnings[1].Code != "ai_deployment_unavailable") {
				t.Fatal("deployment failure concealed")
			}
		})
	}
}
func TestSourceCaptureHashes(t *testing.T) {
	for name, want := range map[string]string{"source-all.json": "d23f7beca04f4c8479c081b2b6103615875eb3e30277a4bc8beb8ec363681d22", "source-deployment-denied.json": "728d371bf3c81f5fb5887b91c3adb23171d36e637167a9fdcde59be49738e551", "source-empty.json": "b8028ac622173c209adab0495eff64b46d2921781358ca25be56b723c958ab3f", "source-unknown-include-tag.json": "b8028ac622173c209adab0495eff64b46d2921781358ca25be56b723c958ab3f", "metrics-input.json": "b666207abd55744a1e8ac74d5ff4c1a49ce3005f845d7f97614781cac0bb1d48", "deployments-input.json": "ac24ff2a680f1505f861e6768f6989415c3f026b543e3376897657d6c6b9ec7c"} {
		b, e := os.ReadFile("testdata/" + name)
		if e != nil {
			t.Fatal(e)
		}
		if fmt.Sprintf("%x", sha256.Sum256(b)) != want {
			t.Fatalf("source capture changed: %s", name)
		}
	}
}
func point() Point {
	return Point{ResourceID: fixtureID, Timestamp: ptr(time.Date(2026, 9, 30, 12, 5, 0, 0, time.UTC)), Count: ptr(1.5), Deployment: ptr("fixture-deployment"), Model: ptr("fixture-model"), Status: ptr("429")}
}
func TestScopeCorrelationAndInvalidValues(t *testing.T) {
	for name, mutate := range map[string]func(*Point){"foreign": func(p *Point) {
		p.ResourceID = strings.Replace(p.ResourceID, fixtureSub, "22222222-2222-4222-8222-222222222222", 1)
	}, "descendant": func(p *Point) { p.ResourceID += "/deployments/fixture-deployment" }} {
		t.Run(name, func(t *testing.T) {
			p := point()
			mutate(&p)
			table, e := Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{point(), p}, fixtureDeployments(), nil)
			if e == nil || table.Health.Error == nil || table.Health.Error.Code != "ai_metric_scope_invalid" || len(table.Rows) != 1 {
				t.Fatal("unrequested metric correlation accepted or prior rows lost")
			}
			valid(t, table)
		})
	}
	for name, value := range map[string]float64{"nan": math.NaN(), "inf": math.Inf(1), "negative": -1} {
		t.Run(name, func(t *testing.T) {
			p := point()
			p.Count = &value
			table, e := Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{point(), p}, fixtureDeployments(), nil)
			if e != nil || len(table.Rows) != 1 || table.Rows[0].Cells[15] != "2" || table.Health.Status != assessment.StageCompletedWithWarnings {
				t.Fatal("invalid count corrupted valid sum")
			}
			valid(t, table)
		})
	}
	p := point()
	p.ResourceID = strings.ToUpper(p.ResourceID)
	table, e := Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{p}, fixtureDeployments(), nil)
	if e != nil || len(table.Rows) != 1 || table.Rows[0].SubscriptionID != fixtureSub {
		t.Fatal("ARM ID case correlation")
	}
	a := fixtureAccounts()
	a = append(a, a[0])
	if table, e = Project(context.Background(), fixtureScope(), a, nil, nil, nil); e == nil || table.Health.Error.Code != "ai_account_invalid" {
		t.Fatal("duplicate account accepted")
	}
}

type cancelContext struct {
	context.Context
	n int
}

func (c *cancelContext) Err() error {
	c.n--
	if c.n <= 0 {
		return context.Canceled
	}
	return nil
}
func TestProjectionLimitsOverflowAndCancellation(t *testing.T) {
	p := point()
	p.Count = ptr(math.MaxFloat64)
	table, e := Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{p, p}, fixtureDeployments(), nil)
	if e == nil || len(table.Rows) != 1 || table.Health.Error.Code != "ai_numeric_limit" {
		t.Fatal("overflow accepted or valid sum lost")
	}
	valid(t, table)
	table, e = Project(context.Background(), fixtureScope(), fixtureAccounts(), make([]Point, MaxPoints+1), nil, nil)
	if e == nil || table.Health.Error.Code != "ai_input_limit" {
		t.Fatal("input budget bypass")
	}
	valid(t, table)
	points := make([]Point, MaxRows+1)
	for i := range points {
		points[i] = point()
		points[i].Deployment = ptr(fmt.Sprintf("deployment-%d", i))
	}
	table, e = Project(context.Background(), fixtureScope(), fixtureAccounts(), points, fixtureDeployments(), nil)
	if e == nil || len(table.Rows) != MaxRows || table.Health.Error.Code != "ai_row_limit" {
		t.Fatal("row budget bypass or retained rows lost")
	}
	valid(t, table)
	// Scope/account/deployment checks precede each point. This deterministic
	// context cancels after the first aggregate, without timing/sleep assumptions.
	ctx := &cancelContext{Context: context.Background(), n: 8}
	table, e = Project(ctx, fixtureScope(), fixtureAccounts(), []Point{point(), point(), point()}, fixtureDeployments(), nil)
	if !errors.Is(e, context.Canceled) || table.Health.Status != assessment.StageFailed || len(table.Rows) != 1 {
		t.Fatal("cancellation identity lost")
	}
	valid(t, table)
	ctx2, cancel := context.WithCancel(context.Background())
	cancel()
	table, e = Project(ctx2, fixtureScope(), nil, nil, nil, nil)
	if !errors.Is(e, context.Canceled) || table.Health.Error.Code != "ai_cancelled" {
		t.Fatal("pre-cancellation accepted")
	}
	valid(t, table)
}

type excludeFilter struct{}

func (excludeFilter) IsServiceExcluded(string) bool { return true }
func TestDefaultsFiltersOwnershipAndConcurrency(t *testing.T) {
	p := point()
	p.Deployment = ptr("")
	p.Model = nil
	p.Status = nil
	table, e := Project(context.Background(), map[string]string{fixtureSub: ""}, fixtureAccounts(), []Point{p}, nil, nil)
	if e != nil || table.Rows[0].Cells[0] != fixtureSub || table.Rows[0].Cells[5] != "" || table.Rows[0].Cells[6] != "Unknown" || table.Health.Warnings[0].Code != "ai_deployment_unavailable" {
		t.Fatal("source defaults or missing enrichment health")
	}
	valid(t, table)
	table, e = Project(context.Background(), fixtureScope(), fixtureAccounts(), nil, nil, excludeFilter{})
	if e != nil || table.SheetName != "AI Throttling" || len(table.Rows) != 0 {
		t.Fatal("source filtered empty sheet")
	}
	valid(t, table)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a, p, d := fixtureAccounts(), []Point{point()}, fixtureDeployments()
			table, e := Project(context.Background(), fixtureScope(), a, p, d, nil)
			if e != nil {
				t.Error(e)
				return
			}
			a[0].Name = "changed"
			*p[0].Deployment = "changed"
			*d[fixtureID].Values[0].ModelFormat = "changed"
			if table.Rows[0].Cells[2] != "fixture-ai" || table.Rows[0].Cells[5] != "fixture-deployment" || table.Rows[0].Cells[8] != "OpenAI" {
				t.Error("output aliases caller input")
			}
			table.Rows[0].Cells[0] = "changed"
			table.Columns[0] = "changed"
			again, e := Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{point()}, fixtureDeployments(), nil)
			if e != nil || again.Columns[0] != "Subscription" || again.Rows[0].Cells[0] != "Synthetic subscription" {
				t.Error("per-run state leaked")
			}
		}()
	}
	wg.Wait()
}

func TestProjectionOutputTextAndMalformedEnrichment(t *testing.T) {
	a := fixtureAccounts()[:1]
	a[0].Kind = strings.Repeat("k", MaxLabelBytes)
	a[0].SKU = strings.Repeat("s", MaxLabelBytes)
	d := fixtureDeployments()
	v := d[fixtureID]
	v.Values[0].ModelVersion = ptr(strings.Repeat("v", MaxLabelBytes))
	v.Values[0].ModelFormat = ptr(strings.Repeat("f", MaxLabelBytes))
	v.Values[0].Upgrade = ptr(strings.Repeat("u", MaxLabelBytes))
	v.Values[0].Spillover = ptr(strings.Repeat("p", MaxLabelBytes))
	d[fixtureID] = v
	points := make([]Point, MaxRows)
	for i := range points {
		points[i] = point()
		points[i].Model = ptr(fmt.Sprintf("%0512d", i))
	}
	table, e := Project(context.Background(), fixtureScope(), a, points, d, nil)
	if e == nil || table.Health.Error.Code != "ai_output_text_limit" || len(table.Rows) == 0 || len(table.Rows) >= MaxRows {
		t.Fatal("repeated-cell text budget bypass")
	}
	valid(t, table)
	d = fixtureDeployments()
	v = d[fixtureID]
	v.Values[0].Capacity = ptr(int64(-1))
	v.Values = append(v.Values, Deployment{Name: nil})
	d[fixtureID] = v
	table, e = Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{point()}, d, nil)
	if e != nil || table.Health.Status != assessment.StageCompletedWithWarnings || table.Rows[0].Cells[9] != "N/A" {
		t.Fatal("malformed deployment treated as complete metadata")
	}
	valid(t, table)
	p := point()
	p.Deployment = ptr(strings.Repeat("x", MaxLabelBytes+1))
	table, e = Project(context.Background(), fixtureScope(), fixtureAccounts(), []Point{point(), p}, fixtureDeployments(), nil)
	if e != nil || len(table.Rows) != 1 || table.Health.Status != assessment.StageCompletedWithWarnings {
		t.Fatal("label limit bypass")
	}
	valid(t, table)
}
