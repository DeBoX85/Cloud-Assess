package region

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type inputCase struct {
	Case  string
	Input Comparison
}
type capturedCase struct {
	Case     string
	Score    float64
	Metadata struct {
		Name, Version, Description, Author, License string
		Type                                        int
	}
	Table [][]string
}

func fixtures(t *testing.T) ([]inputCase, []capturedCase) {
	t.Helper()
	var in []inputCase
	var out []capturedCase
	for name, value := range map[string]any{"source-inputs.json": &in, "source-primary.json": &out} {
		b, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(b, value); err != nil {
			t.Fatal(err)
		}
	}
	if len(in) != 42 || len(in) != len(out) {
		t.Fatal("independent source fixture count changed")
	}
	return in, out
}

func TestSourceCaptureHashes(t *testing.T) {
	for name, want := range map[string]string{"source-inputs.json": "121d891edf09aaffc53db7bf6063021250b9d7af0e2277e1b89fafecee1050a6", "source-primary.json": "8ce00e410df46af95630c2240dcfafab7760fd341316e3668c27e55018d8d985"} {
		b, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(b)
		if hex.EncodeToString(sum[:]) != want {
			t.Fatal("independent source bytes changed: " + name)
		}
	}
}

func TestPrimarySourceScoresAndEveryCell(t *testing.T) {
	in, out := fixtures(t)
	for i, c := range in {
		t.Run(c.Case, func(t *testing.T) {
			want := out[i]
			if want.Case != c.Case {
				t.Fatal("source case mismatch")
			}
			if math.Abs(calculateScore(c.Input)-want.Score) > 1e-12 {
				t.Fatalf("score got%.17g want%.17g", calculateScore(c.Input), want.Score)
			}
			table, err := Project(context.Background(), map[string]string{c.Input.SubscriptionID: c.Input.SubscriptionName}, []Comparison{c.Input})
			if err != nil {
				t.Fatal(err)
			}
			m := want.Metadata
			if table.Metadata != (assessment.PluginMetadata{Name: m.Name, Version: m.Version, Description: m.Description, Author: m.Author, License: m.License, Type: "internal"}) || m.Type != 1 {
				t.Fatal("source metadata mismatch")
			}
			if !slices.Equal(table.Columns, want.Table[0]) || len(table.Rows) != 1 || !slices.Equal(table.Rows[0].Cells, want.Table[1]) {
				t.Fatalf("complete captured table mismatch\ngot%#v\nwant%q", table, want.Table)
			}
			if table.Rows[0].SubscriptionID != c.Input.SubscriptionID || table.Health.Status != assessment.StageCompleted || table.Health.Records != 1 || assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
				t.Fatal("canonical scoped health/identity")
			}
		})
	}
}

func TestPrimaryRejectsForeignScope(t *testing.T) {
	in, _ := fixtures(t)
	c := in[0].Input
	table, err := Project(context.Background(), map[string]string{"22222222-2222-2222-2222-222222222222": c.SubscriptionName}, []Comparison{c})
	if err == nil || len(table.Rows) != 0 || table.Health.Status != assessment.StageFailed {
		t.Fatal("foreign comparison was reported")
	}
}

func TestPrimaryRejectsMalformedAndExcessiveInput(t *testing.T) {
	in, _ := fixtures(t)
	mutations := map[string]func(*Comparison){
		"name":               func(c *Comparison) { c.SubscriptionName = "private-foreign" },
		"identity":           func(c *Comparison) { c.SubscriptionID = "not-an-id" },
		"region":             func(c *Comparison) { c.TargetRegion = "../secret" },
		"region-unicode":     func(c *Comparison) { c.SourceRegion = "eastKus" },
		"type-count":         func(c *Comparison) { c.AvailableTypes-- },
		"sku-count":          func(c *Comparison) { c.AvailableSKUs-- },
		"negative-count":     func(c *Comparison) { c.SourceZoneCount = -1 },
		"count-limit":        func(c *Comparison) { c.TargetZoneCount = MaxCount + 1 },
		"type-percent":       func(c *Comparison) { c.AvailabilityPercent = 99 },
		"sku-percent":        func(c *Comparison) { c.SKUAvailabilityPercent = 99 },
		"negative-latency":   func(c *Comparison) { c.AvgLatencyMs = -1 },
		"negative-price":     func(c *Comparison) { c.AvgCostDifference = -100.01 },
		"nonfinite-cost":     func(c *Comparison) { c.AvgCostDifference = math.NaN() },
		"nonfinite-latency":  func(c *Comparison) { c.AvgLatencyMs = math.Inf(1) },
		"nonfinite-types":    func(c *Comparison) { c.AvailabilityPercent = math.Inf(-1) },
		"nonfinite-sku":      func(c *Comparison) { c.SKUAvailabilityPercent = math.NaN() },
		"detail-control":     func(c *Comparison) { c.MissingSKUs = []string{"private\nvalue"} },
		"detail-utf8":        func(c *Comparison) { c.MissingResourceTypes = []string{string([]byte{0xff})} },
		"detail-replacement": func(c *Comparison) { c.MissingSKUs = []string{"bad�"} },
		"detail-label-limit": func(c *Comparison) { c.MissingSKUs = []string{strings.Repeat("s", MaxLabelBytes+1)} },
		"detail-count-limit": func(c *Comparison) { c.MissingSKUs = make([]string, MaxDetails+1) },
		"detail-cell-limit":  func(c *Comparison) { c.MissingSKUs = slices.Repeat([]string{strings.Repeat("s", MaxLabelBytes)}, 128) },
		"mapping-control":    func(c *Comparison) { c.TargetZoneMappings = map[string]string{"1": "private\tzone"} },
		"mapping-empty":      func(c *Comparison) { c.TargetZoneMappings = map[string]string{"": "west"} },
		"mapping-count-limit": func(c *Comparison) {
			c.TargetZoneMappings = map[string]string{}
			for i := 0; i <= MaxDetails; i++ {
				c.TargetZoneMappings[strings.Repeat("s", i+1)] = "west"
			}
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			c := in[0].Input
			mutate(&c)
			table, err := Project(context.Background(), map[string]string{in[0].Input.SubscriptionID: in[0].Input.SubscriptionName}, []Comparison{c})
			if err == nil || table.Health.Status != assessment.StageFailed || len(table.Rows) != 0 || strings.Contains(err.Error(), "private") {
				t.Fatal("malformed input accepted or exposed")
			}
		})
	}
	c := in[0].Input
	scope := map[string]string{c.SubscriptionID: c.SubscriptionName}
	for name, values := range map[string][]Comparison{"duplicate": {c, c}, "comparison-limit": make([]Comparison, MaxComparisons+1)} {
		t.Run(name, func(t *testing.T) {
			if _, err := Project(context.Background(), scope, values); err == nil {
				t.Fatal("excessive/duplicate comparisons accepted")
			}
		})
	}
	alias := strings.ToUpper("abcdefab-1111-1111-1111-111111111111")
	if _, err := Project(context.Background(), map[string]string{alias: "a", strings.ToLower(alias): "b"}, nil); err == nil {
		t.Fatal("duplicate scope alias")
	}
	if _, err := Project(context.Background(), map[string]string{"invalid": "a"}, nil); err == nil {
		t.Fatal("malformed empty scope")
	}
	if _, err := Project(context.Background(), map[string]string{c.SubscriptionID: strings.Repeat("s", MaxLabelBytes+1)}, nil); err == nil {
		t.Fatal("scope text limit")
	}
	overscope := map[string]string{}
	for i := 0; i <= MaxSubscriptions; i++ {
		overscope[strings.Repeat("s", i+1)] = "a"
	}
	if _, err := Project(context.Background(), overscope, nil); err == nil {
		t.Fatal("scope count limit")
	}
}

func TestPrimaryOwnershipOrderingEmptyAndCancellation(t *testing.T) {
	in, _ := fixtures(t)
	c := in[0].Input
	c.RestrictedSKUs = []string{"type:a"}
	c.AvailableSKUs = 3
	c.SKUAvailabilityPercent = 75
	scope := map[string]string{c.SubscriptionID: c.SubscriptionName}
	before, _ := json.Marshal(c)
	table, err := Project(context.Background(), scope, []Comparison{c})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(c)
	if string(before) != string(after) {
		t.Fatal("caller input mutated")
	}
	c.TargetZoneMappings["1"] = "mutated"
	c.RestrictedSKUs[0] = "mutated"
	scope[c.SubscriptionID] = "mutated"
	if strings.Contains(strings.Join(table.Rows[0].Cells, ","), "mutated") {
		t.Fatal("output aliases caller")
	}
	table.Columns[0] = "mutated"
	if PendingTable().Columns[0] != "Subscription" {
		t.Fatal("columns escape")
	}
	empty, err := Project(context.Background(), nil, nil)
	if err != nil || len(empty.Columns) != 25 || len(empty.Rows) != 0 || empty.Health.Status != assessment.StageCompleted {
		t.Fatal("pure empty contract")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	failed, err := Project(ctx, nil, nil)
	if !errors.Is(err, context.Canceled) || failed.Health.Status != assessment.StageFailed || len(failed.Rows) != 0 {
		t.Fatal("cancelled projection")
	}
	// Literal source primary ordering is score descending, then display name.
	a, b, d := in[0].Input, in[0].Input, in[1].Input
	a.SubscriptionID = "22222222-2222-2222-2222-222222222222"
	a.SubscriptionName = "B"
	b.SubscriptionID = "33333333-3333-3333-3333-333333333333"
	b.SubscriptionName = "A"
	d.SubscriptionID = "44444444-4444-4444-4444-444444444444"
	d.SubscriptionName = "A"
	scope = map[string]string{a.SubscriptionID: a.SubscriptionName, b.SubscriptionID: b.SubscriptionName, d.SubscriptionID: d.SubscriptionName}
	ordered, err := Project(context.Background(), scope, []Comparison{d, a, b})
	if err != nil {
		t.Fatal(err)
	}
	if ordered.Rows[0].SubscriptionID != b.SubscriptionID || ordered.Rows[1].SubscriptionID != a.SubscriptionID || ordered.Rows[2].SubscriptionID != d.SubscriptionID {
		t.Fatal("score/name ordering")
	}
	a.SubscriptionName = "A"
	scope[a.SubscriptionID] = "A"
	x, err := Project(context.Background(), scope, []Comparison{b, a})
	if err != nil {
		t.Fatal(err)
	}
	y, err := Project(context.Background(), scope, []Comparison{a, b})
	if err != nil || !reflect.DeepEqual(x, y) || x.Rows[0].SubscriptionID != a.SubscriptionID {
		t.Fatal("deterministic equal-score/name identity tie")
	}
}

func TestPrimaryConcurrentIsolation(t *testing.T) {
	in, _ := fixtures(t)
	var wg sync.WaitGroup
	for _, c := range in {
		wg.Add(1)
		go func(c inputCase) {
			defer wg.Done()
			scope := map[string]string{c.Input.SubscriptionID: c.Input.SubscriptionName}
			first, err := Project(context.Background(), scope, []Comparison{c.Input})
			if err != nil {
				t.Error(err)
				return
			}
			for i := 0; i < 4; i++ {
				again, err := Project(context.Background(), scope, []Comparison{c.Input})
				if err != nil || !reflect.DeepEqual(first, again) {
					t.Error("per-run projection changed")
				}
			}
		}(c)
	}
	wg.Wait()
}

func TestPrimaryWorkLimits(t *testing.T) {
	in, _ := fixtures(t)
	base := in[0].Input
	input := make([]Comparison, MaxComparisons+1)
	for i := range input {
		input[i] = base
		input[i].SourceRegion = fmt.Sprintf("source%d", i)
	}
	scope := map[string]string{base.SubscriptionID: base.SubscriptionName}
	if _, err := Project(context.Background(), scope, input); err == nil {
		t.Fatal("comparison work limit not enforced")
	}
	table, err := Project(context.Background(), scope, input[:MaxComparisons])
	if err != nil || len(table.Rows) != MaxComparisons {
		t.Fatal("exact comparison limit should pass", err)
	}
	base.MissingSKUs = slices.Repeat([]string{strings.Repeat("s", MaxLabelBytes)}, 63)
	if _, err := Project(context.Background(), scope, []Comparison{base}); err != nil {
		t.Fatal("bounded joined detail should pass", err)
	}
	base.MissingSKUs = append(base.MissingSKUs, strings.Repeat("s", MaxLabelBytes))
	if _, err := Project(context.Background(), scope, []Comparison{base}); err == nil {
		t.Fatal("joined cell limit not enforced")
	}
	for i := range input {
		input[i].MissingSKUs = slices.Repeat([]string{strings.Repeat("s", MaxLabelBytes)}, 32)
	}
	failed, err := Project(context.Background(), scope, input[:MaxComparisons])
	if err == nil || failed.Health.Status != assessment.StageFailed || len(failed.Rows) != 0 {
		t.Fatal("aggregate text budget should fail without partial rows")
	}
}

func TestPrimaryAggregateDetailWorkRejectedBeforeProjection(t *testing.T) {
	in, _ := fixtures(t)
	base := in[0].Input
	base.TotalSKUsChecked, base.AvailableSKUs, base.UnavailableSKUs, base.UnknownSKUs = MaxDetails*2, 0, 0, 0
	base.SKUAvailabilityPercent = 0
	base.RestrictedSKUs = make([]string, MaxDetails)
	base.ZoneRestrictedSKUs = make([]string, MaxDetails)
	base.MissingResourceTypes = make([]string, MaxDetails)
	base.MissingSKUs = make([]string, MaxDetails)
	input := make([]Comparison, 5)
	for i := range input {
		input[i] = base
		input[i].SourceRegion = fmt.Sprintf("source%d", i)
	}
	scope := map[string]string{base.SubscriptionID: base.SubscriptionName}
	failed, err := Project(context.Background(), scope, input)
	if err == nil || failed.Health.Error == nil || failed.Health.Error.Code != "region_input_limit" || len(failed.Rows) != 0 {
		t.Fatal("aggregate short/empty detail work must fail before projection")
	}
}

func TestPrimaryJoinedSeparatorsCountBeforeProjection(t *testing.T) {
	in, _ := fixtures(t)
	base := in[0].Input
	base.TotalSKUsChecked, base.AvailableSKUs, base.UnavailableSKUs = 64, 0, 0
	base.SKUAvailabilityPercent = 0
	details := slices.Repeat([]string{strings.Repeat("s", 511)}, 32)
	base.MissingResourceTypes, base.MissingSKUs, base.RestrictedSKUs, base.ZoneRestrictedSKUs = details, details, details, details
	input := make([]Comparison, 256)
	for i := range input {
		input[i] = base
		input[i].SourceRegion = fmt.Sprintf("source%d", i)
	}
	failed, err := Project(context.Background(), map[string]string{base.SubscriptionID: base.SubscriptionName}, input)
	if err == nil || failed.Health.Error == nil || failed.Health.Error.Code != "region_text_limit" || len(failed.Rows) != 0 {
		t.Fatal("joined separator bytes must fail aggregate budget before projection")
	}
}
