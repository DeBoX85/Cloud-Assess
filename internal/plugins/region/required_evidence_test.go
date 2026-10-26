package region

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func admissionInput(t *testing.T) (*RequiredInventoryPlan, RequiredEvidenceRequest, []RequiredCollectionInput) {
	t.Helper()
	c := requiredFixtureLoad(t)[0]
	p := requiredFixtureBuild(t, c)
	q := RequiredEvidenceRequest{c.Request.SubscriptionID, c.Request.SourceRegion, c.Request.TargetRegion}
	cs := make([]RequiredCollectionInput, len(c.Collections))
	for i, x := range c.Collections {
		cs[i] = RequiredCollectionInput{x.SubscriptionID, x.SourceRegion, x.TargetRegion, x.ResourceType, x.Kind, x.Storage, x.Disk, x.VM, x.Cognitive, x.SQL}
	}
	cs[0].VM.Groups["standard_d2s_v5"] = ComputeVMSKUGroup{Known: true, State: SKUZoneRestricted, BlockedZones: []string{"1", "2"}}
	return p, q, cs
}
func TestRequiredAdmissionCoreEquivalence(t *testing.T) {
	count := 0
	for _, c := range requiredFixtureLoad(t) {
		if c.Want.ErrorPhase == "inventory" {
			continue
		}
		count++
		t.Run(c.Name, func(t *testing.T) {
			p := requiredFixtureBuild(t, c)
			q := RequiredEvidenceRequest{c.Request.SubscriptionID, c.Request.SourceRegion, c.Request.TargetRegion}
			cs := []RequiredCollectionInput{}
			for _, x := range c.Collections {
				cs = append(cs, RequiredCollectionInput{x.SubscriptionID, x.SourceRegion, x.TargetRegion, x.ResourceType, x.Kind, x.Storage, x.Disk, x.VM, x.Cognitive, x.SQL})
			}
			a, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
			b, f := admitRequiredEvidenceObserved(context.Background(), p, q, cs, nil)
			if !reflect.DeepEqual(a, b) || (e == nil) != (f == nil) || e != nil && e.Error() != f.Error() {
				t.Fatal("ADMISSION_CORE_EQUIVALENCE")
			}
		})
	}
	if count == 0 {
		t.Fatal("no admission cases")
	}
}
func TestRequiredAdmissionForgeryAndOwnership(t *testing.T) {
	p, q, cs := admissionInput(t)
	for _, forged := range []*RequiredInventoryPlan{nil, {}, {Resources: p.Resources, Requirements: p.Requirements, Scopes: p.Scopes, Status: p.Status}} {
		a, e := AdmitRequiredEvidence(context.Background(), forged, q, cs)
		if a != nil || e == nil || e.Error() != "required_evidence_invalid" {
			t.Fatal("ADMISSION_FORGED_PLAN")
		}
	}
	a, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
	if e != nil {
		t.Fatal(e)
	}
	wantQ, want := a.Evidence()
	// Mutate every exported plan container, including nested coordinate views.
	p.Resources[0].Record.Name = "changed"
	p.Requirements[0].SKUs[0] = "changed"
	p.Requirements[0].ResourceIDs[0] = "changed"
	p.Requirements[0].Kinds["changed"] = []string{"changed"}
	p.Scopes[q.SubscriptionID] = RequiredInventoryScope{"changed", "unknown"}
	p.Status = "unknown"
	b, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
	if e != nil {
		t.Fatal(e)
	}
	gotQ, got := b.Evidence()
	if gotQ != wantQ || !reflect.DeepEqual(got, want) {
		t.Fatal("ADMISSION_PLAN_DIAGNOSTIC_OWNERSHIP")
	}
	cs[0].VM.SubscriptionID = "changed"
	g := cs[0].VM.Groups["standard_d2s_v5"]
	g.BlockedZones[0] = "changed"
	cs[0].VM.Groups["standard_d2s_v5"] = g
	delete(cs[0].VM.Groups, "standard_d2s_v5")
	a.Request.SubscriptionID = "changed"
	a.SKUs[vmSKUType].Values["standard_d2s_v5"].BlockedZones[0] = "changed"
	delete(a.SKUs, vmSKUType)
	gotQ, got = a.Evidence()
	if gotQ != wantQ || !reflect.DeepEqual(got, want) {
		t.Fatal("ADMISSION_PRIVATE_OWNERSHIP")
	}
	got[vmSKUType].Values["standard_d2s_v5"].BlockedZones[0] = "changed"
	delete(got, vmSKUType)
	againQ, again := a.Evidence()
	if againQ != wantQ || !reflect.DeepEqual(again, want) {
		t.Fatal("ADMISSION_ACCESSOR_OWNERSHIP")
	}
}
func TestRequiredAdmissionConstructorInputOwnership(t *testing.T) {
	selected, rows, ds := buildInput(t)
	p, e := BuildRequiredInventory(context.Background(), selected, rows, ds)
	if e != nil {
		t.Fatal(e)
	}
	_, q, cs := admissionInput(t)
	a, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
	if e != nil {
		t.Fatal(e)
	}
	rows[0].SKUName = "changed"
	rows[0].Name = "changed"
	rows[0].Tags["not-a-coordinate"] = "changed"
	selected[q.SubscriptionID] = "changed"
	ds[0].Status = "unknown"
	b, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
	if e != nil || !reflect.DeepEqual(a, b) {
		t.Fatal("ADMISSION_CONSTRUCTOR_INPUT_OWNERSHIP")
	}
}
func TestRequiredAdmissionConcurrent(t *testing.T) {
	p, q, cs := admissionInput(t)
	out := make([]*RequiredEvidenceAdmission, 8)
	var wg sync.WaitGroup
	for i := range out {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a, e := AdmitRequiredEvidence(context.Background(), p, q, cs)
			if e != nil {
				t.Error(e)
			}
			out[i] = a
		}(i)
	}
	wg.Wait()
	for _, a := range out {
		if a == nil || !reflect.DeepEqual(a, out[0]) {
			t.Fatal("ADMISSION_CONCURRENT_EQUALITY")
		}
	}
	out[0].SKUs[vmSKUType].Values["standard_d2s_v5"].BlockedZones[0] = "changed"
	if out[1].SKUs[vmSKUType].Values["standard_d2s_v5"].BlockedZones[0] == "changed" {
		t.Fatal("ADMISSION_CROSS_RESULT_ALIAS")
	}
	requests := make([]RequiredEvidenceRequest, 8)
	values := make([]map[string]SKUEvidence, 8)
	for i := range values {
		wg.Add(1)
		go func(i int) { defer wg.Done(); requests[i], values[i] = out[0].Evidence() }(i)
	}
	wg.Wait()
	for i := range values {
		if requests[i] != q || !reflect.DeepEqual(values[i], values[0]) {
			t.Fatal("ADMISSION_CONCURRENT_ACCESSOR")
		}
	}
	values[0][vmSKUType].Values["standard_d2s_v5"].BlockedZones[0] = "changed"
	if values[1][vmSKUType].Values["standard_d2s_v5"].BlockedZones[0] == "changed" {
		t.Fatal("ADMISSION_CONCURRENT_ACCESSOR_ALIAS")
	}
}
func TestRequiredAdmissionContext(t *testing.T) {
	p, q, cs := admissionInput(t)
	a, e := AdmitRequiredEvidence(nil, p, q, cs)
	if a != nil || e == nil || e.Error() != "required_evidence_invalid" {
		t.Fatal("ADMISSION_NIL_CONTEXT")
	}
	canceled, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("private-canary"))
	expired, stop := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer stop()
	for _, ctx := range []context.Context{canceled, expired} {
		for _, rows := range [][]RequiredCollectionInput{{{}}, make([]RequiredCollectionInput, 65537)} {
			a, e := AdmitRequiredEvidence(ctx, nil, q, rows)
			if a != nil || e != ctx.Err() {
				t.Fatal("ADMISSION_ENTRY_TERMINAL")
			}
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
					seen := []string{}
					observe := func(name string) {
						seen = append(seen, name)
						if name == phase {
							hits++
							ctx.trip(cause)
						}
					}
					rows := cs
					if invalid {
						rows = append(append([]RequiredCollectionInput{}, cs...), RequiredCollectionInput{})
					}
					a, e := admitRequiredEvidenceObserved(ctx, p, q, rows, observe)
					if hits != 1 || a != nil || e != cause {
						t.Fatal("ADMISSION_LATE_TERMINAL", hits, e)
					}
					want := []string{"after-input", "after-derived-work", "before-return"}
					if invalid {
						want = []string{"before-return"}
					}
					if !reflect.DeepEqual(seen, want[:len(seen)]) {
						t.Fatal("ADMISSION_PHASE_ORDER")
					}
				})
			}
		}
	}
	for _, invalid := range []bool{false, true} {
		seen := []string{}
		rows := cs
		want := []string{"after-input", "after-derived-work", "before-return"}
		if invalid {
			rows = append(append([]RequiredCollectionInput{}, cs...), RequiredCollectionInput{})
			want = []string{"before-return"}
		}
		a, e := admitRequiredEvidenceObserved(context.Background(), p, q, rows, func(s string) { seen = append(seen, s) })
		b, f := AdmitRequiredEvidence(context.Background(), p, q, rows)
		if !reflect.DeepEqual(seen, want) || !reflect.DeepEqual(a, b) || (e == nil) != (f == nil) {
			t.Fatal("ADMISSION_OBSERVER_EQUIVALENCE")
		}
	}
}
