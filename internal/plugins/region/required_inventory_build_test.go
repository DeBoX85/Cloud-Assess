package region

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

func TestRequiredBuildFrozenLiterals(t *testing.T) {
	for _, c := range requiredFixtureLoad(t) {
		t.Run(c.Name, func(t *testing.T) {
			plan := requiredFixtureBuild(t, c)
			declarations := make([]InventoryScopeDeclaration, len(c.Declarations))
			for i, d := range c.Declarations {
				declarations[i] = InventoryScopeDeclaration{d.SubscriptionID, d.Status}
			}
			core, err := buildRequiredInventoryObserved(context.Background(), c.Selected, c.Resources, declarations, nil)
			if c.Want.ErrorPhase == "inventory" {
				if core != nil || err == nil || err.Error() != c.Want.Error {
					t.Fatal("BUILD_CORE_REJECTION")
				}
				return
			}
			if err != nil || !reflect.DeepEqual(core, plan) {
				t.Fatal("BUILD_CORE_EQUIVALENCE")
			}
		})
	}
}
func buildInput(t *testing.T) (map[string]string, []assessment.Resource, []InventoryScopeDeclaration) {
	t.Helper()
	c := requiredFixtureLoad(t)[0]
	d := make([]InventoryScopeDeclaration, len(c.Declarations))
	for i, x := range c.Declarations {
		d[i] = InventoryScopeDeclaration{x.SubscriptionID, x.Status}
	}
	return c.Selected, c.Resources, d
}
func TestRequiredBuildOwnership(t *testing.T) {
	selected, rows, declarations := buildInput(t)
	p, err := BuildRequiredInventory(context.Background(), selected, rows, declarations)
	if err != nil {
		t.Fatal(err)
	}
	original := p.state.resources[0].Record
	id := p.state.resources[0].SubscriptionID
	rows[0].Name = "changed"
	rows[0].Tags["not-a-coordinate"] = "changed"
	selected[id] = "changed"
	declarations[0].Status = "unknown"
	if !reflect.DeepEqual(p.Resources[0].Record, original) {
		t.Fatal("BUILD_INPUT_OWNERSHIP")
	}
	p.Resources[0].Record.Name = "diagnostic"
	p.Requirements[0].SKUs[0] = "diagnostic"
	p.Requirements[0].ResourceIDs[0] = "diagnostic"
	p.Scopes[id] = RequiredInventoryScope{Name: "diagnostic", Status: "unknown"}
	p.Status = "unknown"
	if p.state.resources[0].Record.Name != "vm" || p.state.requirements[0].SKUs[0] != "standard_d2s_v5" || p.state.requirements[0].ResourceIDs[0] != original.ID || p.state.scopes[id].Status != "complete" {
		t.Fatal("BUILD_DIAGNOSTIC_OWNERSHIP")
	}
}
func TestRequiredBuildStorageOwnership(t *testing.T) {
	cases := requiredFixtureLoad(t)
	for _, c := range cases {
		if c.Name != "storage-all-required-kinds-intersect" {
			continue
		}
		p := requiredFixtureBuild(t, c)
		key := p.Requirements[0].SKUs[0]
		want := append([]string{}, p.state.requirements[0].Kinds[key]...)
		p.Requirements[0].Kinds[key][0] = "changed"
		p.Requirements[0].Kinds["changed"] = []string{"changed"}
		if !reflect.DeepEqual(p.state.requirements[0].Kinds[key], want) || len(p.state.requirements[0].Kinds) != 1 {
			t.Fatal("BUILD_STORAGE_OWNERSHIP")
		}
		return
	}
	t.Fatal("missing literal")
}

type buildSignal struct {
	mu       sync.Mutex
	done     chan struct{}
	err      error
	deadline bool
}

func (s *buildSignal) Deadline() (time.Time, bool) { return time.Unix(1, 0), s.deadline }
func (s *buildSignal) Done() <-chan struct{}       { return s.done }
func (s *buildSignal) Err() error                  { s.mu.Lock(); defer s.mu.Unlock(); return s.err }
func (s *buildSignal) Value(any) any               { return nil }
func (s *buildSignal) trip(e error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err == nil {
		s.err = e
		close(s.done)
	}
}
func TestRequiredBuildContext(t *testing.T) {
	selected, rows, d := buildInput(t)
	p, err := BuildRequiredInventory(nil, selected, rows, d)
	if p != nil || err == nil || err.Error() != "required_inventory_invalid" {
		t.Fatal("BUILD_NIL_CONTEXT")
	}
	canceled, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("must-not-leak"))
	expired, stop := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer stop()
	for _, ctx := range []context.Context{canceled, expired} {
		p, err = BuildRequiredInventory(ctx, selected, []assessment.Resource{{}}, d)
		if p != nil || err != ctx.Err() {
			t.Fatal("BUILD_ENTRY_TERMINAL")
		}
	}
	for _, phase := range []string{"after-input", "after-derived-work", "before-return"} {
		for _, deadline := range []bool{false, true} {
			for _, invalid := range []bool{false, true} {
				if invalid && phase != "before-return" {
					continue
				}
				t.Run(phase+map[bool]string{false: "-cancel", true: "-deadline"}[deadline]+map[bool]string{false: "-healthy", true: "-invalid"}[invalid], func(t *testing.T) {
					ctx := &buildSignal{done: make(chan struct{}), deadline: deadline}
					cause := context.Canceled
					if deadline {
						cause = context.DeadlineExceeded
					}
					hits := 0
					reached := []string{}
					observer := func(name string) {
						reached = append(reached, name)
						if name == phase {
							hits++
							ctx.trip(cause)
						}
					}
					input := rows
					if invalid {
						input = []assessment.Resource{{}}
					}
					p, err := buildRequiredInventoryObserved(ctx, selected, input, d, observer)
					if hits != 1 || p != nil || !errors.Is(err, cause) {
						t.Fatal("BUILD_LATE_TERMINAL", hits, p, err)
					}
					allowed := []string{"after-input", "after-derived-work", "before-return"}
					if invalid {
						allowed = []string{"before-return"}
					}
					at := -1
					for _, name := range reached {
						next := -1
						for i, v := range allowed {
							if v == name {
								next = i
							}
						}
						if next <= at {
							t.Fatal("BUILD_PHASE_ORDER")
						}
						at = next
					}
				})
			}
		}
	}
}
func TestRequiredBuildRepeatConcurrent(t *testing.T) {
	selected, rows, d := buildInput(t)
	plans := make([]*RequiredInventoryPlan, 8)
	var wg sync.WaitGroup
	for i := range plans {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p, e := BuildRequiredInventory(context.Background(), selected, rows, d)
			if e != nil {
				t.Error(e)
			}
			plans[i] = p
		}(i)
	}
	wg.Wait()
	for _, p := range plans {
		if p == nil || !reflect.DeepEqual(p, plans[0]) {
			t.Fatal("BUILD_REPEAT_CONCURRENCY")
		}
	}
	plans[0].Requirements[0].SKUs[0] = "changed"
	if plans[1].Requirements[0].SKUs[0] != "standard_d2s_v5" {
		t.Fatal("BUILD_CROSS_RESULT_ALIAS")
	}
}
func TestRequiredZeroAdmissionAccessor(t *testing.T) {
	for _, a := range []*RequiredEvidenceAdmission{nil, {}} {
		r, v := a.Evidence()
		if r != (RequiredEvidenceRequest{}) || v == nil || len(v) != 0 {
			t.Fatal("ZERO_ADMISSION_ACCESSOR")
		}
		v["changed"] = SKUEvidence{}
		_, again := a.Evidence()
		if len(again) != 0 {
			t.Fatal("ZERO_ACCESSOR_ALIAS")
		}
	}
}

func TestRequiredBuildTypeUnicodeAliasRejected(t *testing.T) {
	selected, rows, d := buildInput(t)
	rows[0].ID = "/subscriptions/" + rows[0].SubscriptionID + "/resourceGroups/rg/providers/microsoft.Kompute/virtualmachines/vm"
	rows[0].Type = "microsoft.kompute/virtualmachines"
	p, e := BuildRequiredInventory(context.Background(), selected, rows, d)
	if p != nil || e == nil || e.Error() != "required_inventory_invalid" {
		t.Fatal("BUILD_TYPE_ASCII_IDENTITY", p, e)
	}
}

func TestRequiredBuildUnrepresentableCognitiveKindIsolation(t *testing.T) {
	selected, rows, d := buildInput(t)
	base := rows[0]
	base.ID = "/subscriptions/" + base.SubscriptionID + "/resourceGroups/rg/providers/microsoft.cognitiveservices/accounts/c1"
	base.Type = "microsoft.cognitiveservices/accounts"
	base.Kind = "k"
	base.SKUName = "S0"
	other := base
	other.ID = "/subscriptions/" + base.SubscriptionID + "/resourceGroups/rg/providers/microsoft.cognitiveservices/accounts/c2"
	other.Kind = "K"
	p, e := BuildRequiredInventory(context.Background(), selected, []assessment.Resource{base, other}, d)
	if e != nil || p == nil || len(p.Requirements) != 2 {
		t.Fatal("BUILD_COGNITIVE_KIND_ISOLATION", e)
	}
	if p.Requirements[0].Kind != "k" || p.Requirements[0].UnrepresentableKind || p.Requirements[1].Kind != "K" || !p.Requirements[1].UnrepresentableKind || p.Requirements[1].MissingKind || p.Resources[1].Kind != "K" {
		t.Fatal("BUILD_COGNITIVE_KIND_ISOLATION")
	}
}

func TestRequiredBuildObserverHealthyOrder(t *testing.T) {
	selected, rows, d := buildInput(t)
	phases := []string{}
	p, e := buildRequiredInventoryObserved(context.Background(), selected, rows, d, func(phase string) { phases = append(phases, phase) })
	if e != nil || p == nil || !reflect.DeepEqual(phases, []string{"after-input", "after-derived-work", "before-return"}) {
		t.Fatal("BUILD_HEALTHY_PHASE_ORDER", phases, e)
	}
	public, e := BuildRequiredInventory(context.Background(), selected, rows, d)
	if e != nil || !reflect.DeepEqual(p, public) {
		t.Fatal("BUILD_OBSERVER_OUTPUT_EQUIVALENCE")
	}
	phases = []string{}
	p, e = buildRequiredInventoryObserved(context.Background(), selected, []assessment.Resource{{}}, d, func(phase string) { phases = append(phases, phase) })
	if p != nil || e == nil || e.Error() != "required_inventory_invalid" || !reflect.DeepEqual(phases, []string{"before-return"}) {
		t.Fatal("BUILD_ERROR_PHASE_ORDER", phases, e)
	}
}
