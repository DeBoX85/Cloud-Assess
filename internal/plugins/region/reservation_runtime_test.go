package region

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const reservationTestID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

func reservationScope() map[string]string { return map[string]string{reservationTestID: "Selected"} }
func reservationRequest() ReservationRequest {
	return ReservationRequest{SubscriptionID: reservationTestID, Region: "eastus"}
}
func reservationUsage(name string, reserved, allocated int64) ReservationUsage {
	return ReservationUsage{ResourceID: "/subscriptions/" + reservationTestID + "/resourceGroups/rg/providers/Microsoft.Compute/capacityReservationGroups/g/capacityReservations/" + name, Region: "eastus", SKU: "Standard_D2s_v5", Reserved: reserved, Allocated: allocated, ReservedKnown: true, AllocatedKnown: true}
}
func calculateTestReservations(usages []ReservationUsage) (*ReservationCalculation, error) {
	r := reservationRequest()
	return CalculateReservations(context.Background(), reservationScope(), []ReservationRequest{r}, []ReservationEvidence{{Request: r, Status: "complete", Reservations: usages}})
}

func TestReservationRuntimeLiteralStatusAndCells(t *testing.T) {
	input := []ReservationUsage{reservationUsage("idle", 4, 0), reservationUsage("available", 4, 1), reservationUsage("at", 4, 4), reservationUsage("over", 4, 5), reservationUsage("zero", 0, 0), reservationUsage("zero-over", 0, 1)}
	got, err := calculateTestReservations(input)
	if err != nil {
		t.Fatal(err)
	}
	wantCells := [][]string{
		{"Selected", "eastus", "rg", "g", "idle", "Standard_D2s_v5", "4", "0", "4", "Idle"},
		{"Selected", "eastus", "rg", "g", "available", "Standard_D2s_v5", "4", "1", "3", "Available"},
		{"Selected", "eastus", "rg", "g", "at", "Standard_D2s_v5", "4", "4", "0", "At-Capacity"},
		{"Selected", "eastus", "rg", "g", "over", "Standard_D2s_v5", "4", "5", "-1", "Over-Allocated"},
		{"Selected", "eastus", "rg", "g", "zero", "Standard_D2s_v5", "0", "0", "0", "Idle"},
		{"Selected", "eastus", "rg", "g", "zero-over", "Standard_D2s_v5", "0", "1", "-1", "Over-Allocated"},
	}
	wantCounts := [][3]int64{{4, 0, 4}, {4, 1, 3}, {4, 4, 0}, {4, 5, -1}, {0, 0, 0}, {0, 1, -1}}
	if len(got.Reservations) != 6 || got.Table == nil {
		t.Fatalf("reservation literal output: %#v", got)
	}
	wantHealth := assessment.StageExecution{Name: Name, Status: assessment.StageCompleted, Records: 6}
	if !reflect.DeepEqual(got.Health, wantHealth) || !reflect.DeepEqual(got.Table.Health, wantHealth) {
		t.Fatalf("healthy reservation health: %#v", got.Health)
	}
	for i, record := range got.Reservations {
		want := ReservationValue{ResourceID: strings.ToLower(input[i].ResourceID), SubscriptionID: reservationTestID, Subscription: "Selected", Region: "eastus", ResourceGroup: "rg", GroupName: "g", ReservationName: wantCells[i][4], SKU: "Standard_D2s_v5", Reserved: wantCounts[i][0], Allocated: wantCounts[i][1], Available: wantCounts[i][2], Status: wantCells[i][9]}
		if !reflect.DeepEqual(record, want) || got.Table.Rows[i].SubscriptionID != reservationTestID || !reflect.DeepEqual(got.Table.Rows[i].Cells, wantCells[i]) {
			t.Fatalf("reservation full literal %d: %#v %#v", i, record, got.Table.Rows[i])
		}
	}
	if !reflect.DeepEqual(got.Table.Columns, []string{"Subscription", "Region", "Resource Group", "CRG Name", "Reservation Name", "SKU", "Reserved", "Allocated", "Available", "Status"}) {
		t.Fatal("reservation columns changed")
	}
}

func TestReservationRuntimeRetainedKnownSourceCases(t *testing.T) {
	type sourceEntry struct {
		SubscriptionID, SubscriptionName, ResourceGroup, CRGName, ReservationName, Location, SKU string
		Reserved, Allocated, Available                                                           int64
		Status                                                                                   string
	}
	var results map[string]struct{ Entries []sourceEntry }
	raw, err := os.ReadFile("testdata/source-reservation-outputs.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &results); err != nil {
		t.Fatal(err)
	}
	// These are decoded source count/status comparisons with synthetic admission
	// IDs, not execution of a collector or captured Get response identities.
	for _, key := range []string{"idle", "available", "at", "over", "zero"} {
		result, ok := results[key]
		if !ok || len(result.Entries) != 1 {
			t.Fatalf("retained source case absent: %s", key)
		}
		e := result.Entries[0]
		request := ReservationRequest{SubscriptionID: e.SubscriptionID, Region: e.Location}
		usage := ReservationUsage{ResourceID: "/subscriptions/" + e.SubscriptionID + "/resourceGroups/" + e.ResourceGroup + "/providers/Microsoft.Compute/capacityReservationGroups/" + e.CRGName + "/capacityReservations/" + e.ReservationName, Region: e.Location, ResponseName: e.ReservationName, SKU: e.SKU, Reserved: e.Reserved, Allocated: e.Allocated, ReservedKnown: true, AllocatedKnown: true}
		got, err := CalculateReservations(context.Background(), map[string]string{e.SubscriptionID: e.SubscriptionName}, []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "complete", Reservations: []ReservationUsage{usage}}})
		if err != nil || len(got.Reservations) != 1 {
			t.Fatalf("source %s: %v %#v", key, err, got)
		}
		want := ReservationValue{ResourceID: strings.ToLower(usage.ResourceID), SubscriptionID: e.SubscriptionID, Subscription: e.SubscriptionName, Region: e.Location, ResourceGroup: e.ResourceGroup, GroupName: e.CRGName, ReservationName: e.ReservationName, SKU: e.SKU, Reserved: e.Reserved, Allocated: e.Allocated, Available: e.Available, Status: e.Status}
		if !reflect.DeepEqual(got.Reservations[0], want) {
			t.Fatalf("source full calculated fields %s: %#v vs %#v", key, got.Reservations[0], want)
		}
	}
}

func TestReservationRuntimeHealthAndOwnership(t *testing.T) {
	request := reservationRequest()
	for _, status := range []string{"complete", "partial", "unknown", "unsupported", "missing"} {
		t.Run(status, func(t *testing.T) {
			var evidence []ReservationEvidence
			if status != "missing" {
				evidence = []ReservationEvidence{{Request: request, Status: status}}
			}
			got, err := CalculateReservations(context.Background(), reservationScope(), []ReservationRequest{request}, evidence)
			if err != nil || got.Table != nil || len(got.Reservations) != 0 || got.Health.Records != 0 {
				t.Fatalf("empty evidence: %v %#v", err, got)
			}
			if status == "complete" {
				if got.Health.Status != assessment.StageCompleted || len(got.Health.Warnings) != 0 {
					t.Fatal("complete empty health")
				}
				return
			}
			if got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Code != "reservation_evidence_"+status {
				t.Fatalf("incomplete evidence health: %#v", got.Health)
			}
		})
	}
	unknown := reservationUsage("unknown", 4, 0)
	unknown.AllocatedKnown = false
	missingSKU := reservationUsage("sku", 4, 0)
	missingSKU.SKU = ""
	input := []ReservationEvidence{{Request: request, Status: "partial", Reservations: []ReservationUsage{reservationUsage("healthy", 4, 1), unknown, missingSKU}}}
	got, err := CalculateReservations(context.Background(), reservationScope(), []ReservationRequest{request}, input)
	wantHealth := assessment.StageExecution{Name: Name, Status: assessment.StageCompletedWithWarnings, Records: 1, Warnings: []assessment.AssessmentWarning{
		{Code: "reservation_evidence_partial", Message: "region reservation evidence is incomplete or unusable"},
		{Code: "reservation_usage_unknown", Message: "region reservation evidence is incomplete or unusable"},
		{Code: "reservation_sku_unknown", Message: "region reservation evidence is incomplete or unusable"},
	}}
	if err != nil || got.Table == nil || len(got.Reservations) != 1 || !reflect.DeepEqual(got.Health, wantHealth) || !reflect.DeepEqual(got.Table.Health, wantHealth) {
		t.Fatalf("partial source corrections: %v %#v", err, got)
	}
	input[0].Reservations[0].SKU = "mutated"
	got.Reservations[0].SKU = "mutated"
	got.Health.Warnings[0].Code = "mutated"
	if got.Table.Rows[0].Cells[5] != "Standard_D2s_v5" || got.Table.Health.Warnings[0].Code != "reservation_evidence_partial" {
		t.Fatal("reservation ownership alias")
	}
	for _, usage := range []ReservationUsage{unknown, missingSKU} {
		again, err := calculateTestReservations([]ReservationUsage{usage})
		if err != nil || again.Table != nil || again.Health.Status != assessment.StageCompletedWithWarnings {
			t.Fatalf("unknown must not be Idle: %v %#v", err, again)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := calculateTestReservations([]ReservationUsage{reservationUsage("fresh", 4, 0)})
			if err != nil || r.Reservations[0].SKU != "Standard_D2s_v5" || r.Health.Status != assessment.StageCompleted {
				t.Errorf("reservation per-call isolation: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestReservationRuntimeIdentityAdmission(t *testing.T) {
	valid := reservationUsage("r", 4, 0)
	for _, id := range []string{"", valid.ResourceID + "/", valid.ResourceID + "/extra", strings.Replace(valid.ResourceID, "Microsoft.Compute", "Microsoft.Other", 1), strings.Replace(valid.ResourceID, "subscriptions", "subſcriptions", 1), strings.Replace(valid.ResourceID, reservationTestID, "11111111-1111-1111-1111-111111111111", 1), strings.Replace(valid.ResourceID, "/g/", "/../", 1), valid.ResourceID + "?token=canary", valid.ResourceID + "#x", strings.Replace(valid.ResourceID, "/r", "/r%2fx", 1), strings.Replace(valid.ResourceID, "resourceGroups/rg", "resourceGroups/", 1), strings.Replace(valid.ResourceID, "/g/", "/bad\\name/", 1)} {
		u := valid
		u.ResourceID = id
		got, err := calculateTestReservations([]ReservationUsage{u})
		if err == nil || got != nil {
			t.Fatalf("reservation structural identity bypass: %#v", got)
		}
		if strings.Contains(err.Error(), "canary") {
			t.Fatal("reservation error echoes evidence")
		}
	}
	for _, change := range []func(*ReservationUsage){func(u *ReservationUsage) { u.Region = "westus" }, func(u *ReservationUsage) { u.Region = "EastUS" }, func(u *ReservationUsage) { u.ResponseName = "foreign" }, func(u *ReservationUsage) { u.ResponseRegion = "westus" }, func(u *ReservationUsage) { u.ResponseRegion = "EastUS" }, func(u *ReservationUsage) { u.Reserved = -1 }, func(u *ReservationUsage) { u.Allocated = -1 }, func(u *ReservationUsage) { u.Reserved = MaxAuxCount + 1 }, func(u *ReservationUsage) { u.Allocated = MaxAuxCount + 1 }, func(u *ReservationUsage) { u.ReservedKnown = false }, func(u *ReservationUsage) { u.AllocatedKnown = false; u.Allocated = 1 }, func(u *ReservationUsage) { u.SKU = string([]byte{255}) }, func(u *ReservationUsage) { u.SKU = "bad\n" }, func(u *ReservationUsage) { u.SKU = strings.Repeat("x", MaxLabelBytes+1) }} {
		u := valid
		change(&u)
		got, err := calculateTestReservations([]ReservationUsage{u})
		if err == nil || got != nil {
			t.Fatalf("reservation unsafe evidence admitted: %#v", got)
		}
	}
	upper := valid
	upper.ResourceID = strings.ToUpper(valid.ResourceID)
	upper.ResponseName = "R"
	got, err := calculateTestReservations([]ReservationUsage{upper})
	if err != nil || got.Reservations[0].ResourceID != strings.ToLower(valid.ResourceID) || got.Reservations[0].ReservationName != "R" {
		t.Fatalf("reservation canonical identity/display: %v %#v", err, got)
	}
	got, err = calculateTestReservations([]ReservationUsage{valid, upper})
	if err == nil || got != nil {
		t.Fatal("case-insensitive duplicate reservation ID accepted")
	}
	unknown := valid
	unknown.Reserved = 0
	unknown.ReservedKnown = false
	foreignUnknown := unknown
	foreignUnknown.ResourceID = strings.Replace(valid.ResourceID, reservationTestID, "11111111-1111-1111-1111-111111111111", 1)
	got, err = calculateTestReservations([]ReservationUsage{foreignUnknown})
	if err == nil || got != nil {
		t.Fatal("foreign identity in unknown evidence accepted")
	}
	negativeUnknown := unknown
	negativeUnknown.ReservedKnown = true
	negativeUnknown.Reserved = -1
	negativeUnknown.AllocatedKnown = false
	got, err = calculateTestReservations([]ReservationUsage{negativeUnknown})
	if err == nil || got != nil {
		t.Fatal("negative count in unknown evidence accepted")
	}
	longUnknown := unknown
	longUnknown.SKU = strings.Repeat("x", MaxLabelBytes+1)
	got, err = calculateTestReservations([]ReservationUsage{longUnknown})
	if err == nil || got != nil {
		t.Fatal("oversized label in unknown evidence accepted")
	}
	got, err = calculateTestReservations([]ReservationUsage{unknown, unknown})
	if err == nil || got != nil {
		t.Fatal("unknown reservation identity duplicate accepted")
	}
	for _, counts := range [][2]int64{{MaxAuxCount, 0}, {MaxAuxCount, MaxAuxCount}, {0, MaxAuxCount}} {
		u := valid
		u.Reserved, u.Allocated = counts[0], counts[1]
		got, err := calculateTestReservations([]ReservationUsage{u})
		if err != nil || len(got.Reservations) != 1 {
			t.Fatalf("reservation count bound rejected: %v", err)
		}
	}
}

func TestReservationRuntimeRequestsAndCancellation(t *testing.T) {
	request := reservationRequest()
	valid := ReservationEvidence{Request: request, Status: "complete", Reservations: []ReservationUsage{reservationUsage("r", 4, 0)}}
	for _, tc := range []struct {
		scope    map[string]string
		requests []ReservationRequest
		evidence []ReservationEvidence
	}{
		{map[string]string{}, []ReservationRequest{request}, nil},
		{map[string]string{reservationTestID: ""}, []ReservationRequest{request}, nil},
		{map[string]string{reservationTestID: "a", strings.ToUpper(reservationTestID): "a"}, nil, nil},
		{reservationScope(), []ReservationRequest{request, request}, nil},
		{reservationScope(), []ReservationRequest{request}, []ReservationEvidence{valid, valid}},
		{reservationScope(), []ReservationRequest{request}, []ReservationEvidence{{Request: ReservationRequest{SubscriptionID: reservationTestID, Region: "westus"}, Status: "complete"}}},
		{reservationScope(), []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "invalid"}}},
		{reservationScope(), []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "unknown", Reservations: valid.Reservations}}},
	} {
		got, err := CalculateReservations(context.Background(), tc.scope, tc.requests, tc.evidence)
		if err == nil || got != nil {
			t.Fatal("reservation scope/request evidence guard bypass")
		}
	}
	for _, r := range []ReservationRequest{{SubscriptionID: "bad", Region: "eastus"}, {SubscriptionID: reservationTestID, Region: "EastUS"}} {
		got, err := CalculateReservations(context.Background(), reservationScope(), []ReservationRequest{r}, nil)
		if err == nil || got != nil {
			t.Fatal("reservation canonical request guard bypass")
		}
	}
	probe := &availabilityCancelContext{Context: context.Background(), at: 10000}
	if _, err := CalculateReservations(probe, reservationScope(), []ReservationRequest{request}, []ReservationEvidence{valid}); err != nil {
		t.Fatal(err)
	}
	for at := 1; at <= probe.calls; at++ {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		got, err := CalculateReservations(ctx, reservationScope(), []ReservationRequest{request}, []ReservationEvidence{valid})
		if got != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("reservation cancellation at %d: %v %#v", at, err, got)
		}
	}
	quota, err := calculateTestQuota([]QuotaUsage{quotaUsage("PublicIPAddresses", 9, 10)})
	if err != nil || quota.Health.Status != assessment.StageCompleted {
		t.Fatal("quota coexistence after reservation failures")
	}
}

func TestReservationRuntimeWorkAndDeclaredOrder(t *testing.T) {
	request := reservationRequest()
	west := request
	west.Region = "westus"
	westUsage := reservationUsage("west", 4, 0)
	westUsage.Region = "westus"
	evidence := []ReservationEvidence{{Request: west, Status: "complete", Reservations: []ReservationUsage{westUsage}}, {Request: request, Status: "complete", Reservations: []ReservationUsage{reservationUsage("east", 4, 1)}}}
	got, err := CalculateReservations(context.Background(), reservationScope(), []ReservationRequest{request, west}, evidence)
	if err != nil || len(got.Reservations) != 2 || got.Reservations[0].ReservationName != "east" || got.Reservations[1].ReservationName != "west" {
		t.Fatalf("reservation declared order: %v %#v", err, got)
	}
	usages := make([]ReservationUsage, MaxAuxRows)
	for i := range usages {
		usages[i] = reservationUsage(fmt.Sprintf("r%d", i), 0, 0)
		usages[i].ReservedKnown = false
	}
	got, err = calculateTestReservations(usages)
	if err != nil || len(got.Reservations) != 0 || got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 {
		t.Fatalf("exact skipped work bound: %v %#v", err, got)
	}
	usages = append(usages, reservationUsage("extra", 4, 0))
	got, err = calculateTestReservations(usages)
	if err == nil || got != nil {
		t.Fatal("skipped reservation work budget bypass")
	}
	requests := make([]ReservationRequest, MaxAuxRows+1)
	got, err = CalculateReservations(context.Background(), reservationScope(), requests, nil)
	if err == nil || got != nil {
		t.Fatal("reservation request work limit bypass")
	}
	for i := range evidence {
		evidence[i].Reservations = make([]ReservationUsage, MaxAuxRows/2+1)
		for j := range evidence[i].Reservations {
			evidence[i].Reservations[j] = reservationUsage(fmt.Sprintf("r%d-%d", i, j), 0, 0)
			evidence[i].Reservations[j].Region = evidence[i].Request.Region
			evidence[i].Reservations[j].ReservedKnown = false
		}
	}
	got, err = CalculateReservations(context.Background(), reservationScope(), []ReservationRequest{request, west}, evidence)
	if err == nil || got != nil {
		t.Fatal("reservation aggregate raw work budget bypass")
	}
}

func TestReservationRuntimeRawAndOutputTextBudgets(t *testing.T) {
	// Distinct maximum-length ID names and unknown counts exercise raw admission
	// independently of output. Exact per-label maximum is valid.
	request := reservationRequest()
	scope := map[string]string{reservationTestID: strings.Repeat("s", MaxLabelBytes)}
	long := func(i int) ReservationUsage {
		name := fmt.Sprintf("%0508d", i) + "tail"
		u := reservationUsage(name, 0, 0)
		u.ResourceID = strings.Replace(u.ResourceID, "resourceGroups/rg", "resourceGroups/"+strings.Repeat("g", MaxLabelBytes), 1)
		u.ResourceID = strings.Replace(u.ResourceID, "capacityReservationGroups/g/", "capacityReservationGroups/"+strings.Repeat("c", MaxLabelBytes)+"/", 1)
		u.SKU = strings.Repeat("k", MaxLabelBytes)
		return u
	}
	u := long(0)
	got, err := CalculateReservations(context.Background(), scope, []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "complete", Reservations: []ReservationUsage{u}}})
	if err != nil || len(got.Reservations) != 1 {
		t.Fatalf("reservation label maximum: %v", err)
	}
	usages := make([]ReservationUsage, MaxAuxRows)
	for i := range usages {
		usages[i] = long(i)
		usages[i].ReservedKnown = false
		usages[i].ResponseName = fmt.Sprintf("%0508d", i) + "tail"
		usages[i].ResponseRegion = "eastus"
	}
	got, err = CalculateReservations(context.Background(), scope, []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "complete", Reservations: usages}})
	if err == nil || got != nil || !strings.Contains(err.Error(), "text_limit") {
		t.Fatalf("skipped raw text budget bypass: %v %#v", err, got)
	}
	// With absent optional response labels, raw input is smaller than output.
	// At this size raw fits, but repeated selected display names exhaust output.
	usages = usages[:6500]
	for i := range usages {
		usages[i].ReservedKnown = true
		usages[i].ResponseName = ""
		usages[i].ResponseRegion = ""
	}
	got, err = CalculateReservations(context.Background(), scope, []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "complete", Reservations: usages}})
	if err == nil || got != nil || !strings.Contains(err.Error(), "output_limit") {
		t.Fatalf("independent projected text budget bypass: %v %#v", err, got)
	}
	// Independently account the pre-edit schema's byte contract, without calling
	// target admission helpers. Tune only SKU lengths, preserving all identities.
	budget := assessment.MaxPluginTextBytes - (64 << 10)
	for _, output := range []bool{false, true} {
		fixtures := make([]ReservationUsage, 6500)
		total := len(scope[reservationTestID])
		if !output {
			total += len(request.SubscriptionID) + len(request.Region)
		}
		for i := range fixtures {
			fixtures[i] = long(i)
			if output {
				total += 5*MaxLabelBytes + len("eastus") + 128
			} else {
				fixtures[i].ReservedKnown = false
				fixtures[i].ResponseName = fmt.Sprintf("%0508d", i) + "tail"
				fixtures[i].ResponseRegion = "eastus"
				total += len(fixtures[i].ResourceID) + len("eastus") + MaxLabelBytes + len("eastus") + MaxLabelBytes
			}
		}
		excess := total - budget
		if excess <= 0 {
			t.Fatal("exact text fixture does not reach the declared ceiling")
		}
		adjusted := -1
		floor := 0
		if output {
			floor = 1
		}
		for i := range fixtures {
			remove := min(excess, len(fixtures[i].SKU)-floor)
			fixtures[i].SKU = fixtures[i].SKU[:len(fixtures[i].SKU)-remove]
			excess -= remove
			if remove > 0 {
				adjusted = i
			}
			if excess == 0 {
				break
			}
		}
		if excess != 0 || adjusted < 0 {
			t.Fatal("cannot construct independent exact text fixture")
		}
		got, err = CalculateReservations(context.Background(), scope, []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "complete", Reservations: fixtures}})
		if err != nil {
			t.Fatalf("exact reservation text ceiling (output=%v): %v", output, err)
		}
		if output && len(got.Reservations) != 6500 || !output && (got.Table != nil || got.Health.Status != assessment.StageCompletedWithWarnings) {
			t.Fatal("exact ceiling output/unknown semantics")
		}
		fixtures[adjusted].SKU += "x"
		got, err = CalculateReservations(context.Background(), scope, []ReservationRequest{request}, []ReservationEvidence{{Request: request, Status: "complete", Reservations: fixtures}})
		code := "text_limit"
		if output {
			code = "output_limit"
		}
		if err == nil || got != nil || !strings.Contains(err.Error(), code) {
			t.Fatalf("one-byte ceiling overflow (output=%v): %v", output, err)
		}
	}
}
