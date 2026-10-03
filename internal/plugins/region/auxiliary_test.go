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

const auxSubscription = "11111111-1111-1111-1111-111111111111"
const otherAuxSubscription = "22222222-2222-2222-2222-222222222222"

type auxiliaryCapturedTable struct {
	SheetName   string     `json:"sheet_name"`
	Description string     `json:"description"`
	Table       [][]string `json:"table"`
}

func auxiliaryFixtures(t *testing.T) ([]QuotaRow, []ReservationRow, map[string]*auxiliaryCapturedTable) {
	t.Helper()
	var inputs struct {
		Quotas       []QuotaRow
		Reservations [][]string
	}
	var outputs map[string]*auxiliaryCapturedTable
	for name, target := range map[string]any{"source-aux-inputs.json": &inputs, "source-aux-outputs.json": &outputs} {
		raw, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(raw, target); err != nil {
			t.Fatal(err)
		}
	}
	if len(inputs.Quotas) != 3 || len(inputs.Reservations) != 4 || len(outputs) != 14 {
		t.Fatal("independent auxiliary fixture shape changed")
	}
	for i := range inputs.Quotas {
		inputs.Quotas[i].SubscriptionID = auxSubscription
	}
	reservations := make([]ReservationRow, len(inputs.Reservations))
	for i, cells := range inputs.Reservations {
		reservations[i] = ReservationRow{SubscriptionID: auxSubscription, Cells: cells}
	}
	return inputs.Quotas, reservations, outputs
}

func TestAuxiliaryCapturedBytesAndEveryCell(t *testing.T) {
	for name, want := range map[string]string{
		"source-aux-inputs.json":  "f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28",
		"source-aux-outputs.json": "efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06",
	} {
		raw, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(raw)
		if hex.EncodeToString(sum[:]) != want {
			t.Fatal("source auxiliary bytes changed")
		}
	}
	quotas, reservations, outputs := auxiliaryFixtures(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	quota, err := ProjectQuota(context.Background(), scope, quotas)
	if err != nil {
		t.Fatal(err)
	}
	reservation, err := ProjectReservations(context.Background(), scope, reservations)
	if err != nil {
		t.Fatal(err)
	}
	for name, table := range map[string]*assessment.PluginTable{"quota-full": quota, "reservations-full": reservation} {
		want := outputs[name]
		if table == nil || table.SheetName != want.SheetName || table.Description != want.Description || table.Metadata != Metadata() || !slices.Equal(table.Columns, want.Table[0]) || len(table.Rows) != len(want.Table)-1 {
			t.Fatal("captured metadata/shape mismatch:", name)
		}
		for i, row := range table.Rows {
			if !slices.Equal(row.Cells, want.Table[i+1]) || row.SubscriptionID != auxSubscription {
				t.Fatal("captured cell/correlation mismatch:", name, i)
			}
		}
		if table.Health.Status != assessment.StageCompleted || table.Health.Records != len(table.Rows) || assessment.ValidatePluginTables([]assessment.PluginTable{*table}) != nil {
			t.Fatal("canonical auxiliary health invalid")
		}
	}
	quota, err = ProjectQuota(context.Background(), scope, nil)
	if err != nil || quota != nil || outputs["quota-empty"] != nil {
		t.Fatal("source absent empty quota branch")
	}
	reservation, err = ProjectReservations(context.Background(), scope, nil)
	if err != nil || reservation != nil || outputs["reservations-empty"] != nil {
		t.Fatal("source absent empty reservation branch")
	}
	// Flags, not the displayed rounded percentage, determine source status.
	quotas[0].HeadroomPct = 14.999
	quotas[0].IsNearLimit = false
	quotas[0].IsOverLimit = true
	quota, err = ProjectQuota(context.Background(), scope, quotas[:1])
	if err != nil || quota.Rows[0].Cells[7] != "15.0%" || quota.Rows[0].Cells[8] != "At/Over Limit" {
		t.Fatal("flags recalculated from rounded cells")
	}
}

func requireAuxiliaryFailure(t *testing.T, table *assessment.PluginTable, err error, width int) {
	t.Helper()
	if err == nil || table == nil || len(table.Rows) != 0 || len(table.Columns) != width || table.Health.Status != assessment.StageFailed || table.Health.Error == nil || strings.Contains(err.Error(), "private") || assessment.ValidatePluginTables([]assessment.PluginTable{*table}) != nil {
		t.Fatal("unsafe input accepted, leaked or retained partial rows")
	}
}

func TestAuxiliarySelectedScope(t *testing.T) {
	quotas, reservations, _ := auxiliaryFixtures(t)
	foreign := map[string]string{otherAuxSubscription: "Synthetic"}
	table, err := ProjectQuota(context.Background(), foreign, quotas)
	requireAuxiliaryFailure(t, table, err, 9)
	table, err = ProjectReservations(context.Background(), foreign, reservations)
	requireAuxiliaryFailure(t, table, err, 10)
	scope := map[string]string{auxSubscription: "private-different-name"}
	table, err = ProjectQuota(context.Background(), scope, quotas)
	requireAuxiliaryFailure(t, table, err, 9)
	table, err = ProjectReservations(context.Background(), scope, reservations)
	requireAuxiliaryFailure(t, table, err, 10)
	alias := "ABCDEFAB-1111-1111-1111-111111111111"
	for _, invalid := range []map[string]string{
		{"bad": "Synthetic"},
		{alias: "Synthetic", strings.ToLower(alias): "Synthetic"},
		{auxSubscription: strings.Repeat("s", MaxLabelBytes+1)},
	} {
		table, err = ProjectQuota(context.Background(), invalid, nil)
		requireAuxiliaryFailure(t, table, err, 9)
		table, err = ProjectReservations(context.Background(), invalid, nil)
		requireAuxiliaryFailure(t, table, err, 10)
	}
	// UUID, not the repeated display name, correlates rows.
	quotas[1].SubscriptionID = otherAuxSubscription
	reservations[1].SubscriptionID = otherAuxSubscription
	scope = map[string]string{auxSubscription: "Synthetic", otherAuxSubscription: "Synthetic"}
	table, err = ProjectQuota(context.Background(), scope, quotas)
	if err != nil || table.Rows[1].SubscriptionID != otherAuxSubscription {
		t.Fatal("same display name conflated quota identity")
	}
	table, err = ProjectReservations(context.Background(), scope, reservations)
	if err != nil || table.Rows[1].SubscriptionID != otherAuxSubscription {
		t.Fatal("same display name conflated reservation identity")
	}
}

func TestAuxiliaryMalformedAndDuplicates(t *testing.T) {
	quotas, reservations, _ := auxiliaryFixtures(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	for name, change := range map[string]func(*QuotaRow){
		"id":       func(r *QuotaRow) { r.SubscriptionID = "private-not-an-id" },
		"region":   func(r *QuotaRow) { r.Region = "eastKus" },
		"empty":    func(r *QuotaRow) { r.ResourceName = "" },
		"control":  func(r *QuotaRow) { r.ResourceName = "private\nlabel" },
		"utf8":     func(r *QuotaRow) { r.ResourceName = string([]byte{0xff}) },
		"replace":  func(r *QuotaRow) { r.ResourceName = "bad�" },
		"label":    func(r *QuotaRow) { r.ResourceName = strings.Repeat("s", MaxLabelBytes+1) },
		"negative": func(r *QuotaRow) { r.Current = -1 },
		"count":    func(r *QuotaRow) { r.Limit = MaxAuxCount + 1 },
		"available": func(r *QuotaRow) {
			r.Available = -MaxAuxCount - 1
		},
		"nan":   func(r *QuotaRow) { r.HeadroomPct = math.NaN() },
		"inf":   func(r *QuotaRow) { r.HeadroomPct = math.Inf(1) },
		"large": func(r *QuotaRow) { r.HeadroomPct = float64(MaxAuxCount) + 1 },
	} {
		t.Run("quota-"+name, func(t *testing.T) {
			row := quotas[0]
			change(&row)
			table, err := ProjectQuota(context.Background(), scope, []QuotaRow{quotas[1], row})
			requireAuxiliaryFailure(t, table, err, 9)
		})
	}
	for name, change := range map[string]func(*ReservationRow){
		"id":      func(r *ReservationRow) { r.SubscriptionID = "private-not-an-id" },
		"width":   func(r *ReservationRow) { r.Cells = r.Cells[:9] },
		"region":  func(r *ReservationRow) { r.Cells[1] = "../private" },
		"control": func(r *ReservationRow) { r.Cells[2] = "private\tlabel" },
		"utf8":    func(r *ReservationRow) { r.Cells[2] = string([]byte{0xff}) },
		"empty":   func(r *ReservationRow) { r.Cells[3] = "" },
		"label":   func(r *ReservationRow) { r.Cells[4] = strings.Repeat("s", MaxLabelBytes+1) },
		"count":   func(r *ReservationRow) { r.Cells[6] = "1000000000001" },
		"sign":    func(r *ReservationRow) { r.Cells[6] = "-1" },
		"leading": func(r *ReservationRow) { r.Cells[6] = "010" },
		"overflow": func(r *ReservationRow) {
			r.Cells[8] = "999999999999999999999999999999"
		},
		"status": func(r *ReservationRow) { r.Cells[9] = "private-unknown" },
	} {
		t.Run("reservation-"+name, func(t *testing.T) {
			row := ReservationRow{SubscriptionID: auxSubscription, Cells: slices.Clone(reservations[0].Cells)}
			change(&row)
			table, err := ProjectReservations(context.Background(), scope, []ReservationRow{reservations[1], row})
			requireAuxiliaryFailure(t, table, err, 10)
		})
	}
	table, err := ProjectQuota(context.Background(), scope, []QuotaRow{quotas[0], quotas[0]})
	requireAuxiliaryFailure(t, table, err, 9)
	table, err = ProjectReservations(context.Background(), scope, []ReservationRow{reservations[0], reservations[0]})
	requireAuxiliaryFailure(t, table, err, 10)
}

func TestAuxiliaryWorkLimits(t *testing.T) {
	quotas, reservations, _ := auxiliaryFixtures(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	inputQ := make([]QuotaRow, MaxAuxRows+1)
	inputR := make([]ReservationRow, MaxAuxRows+1)
	for i := range inputQ {
		inputQ[i] = quotas[0]
		inputQ[i].ResourceName = fmt.Sprintf("resource%d", i)
		inputR[i] = ReservationRow{SubscriptionID: auxSubscription, Cells: slices.Clone(reservations[0].Cells)}
		inputR[i].Cells[4] = fmt.Sprintf("reservation%d", i)
	}
	table, err := ProjectQuota(context.Background(), scope, inputQ)
	requireAuxiliaryFailure(t, table, err, 9)
	table, err = ProjectReservations(context.Background(), scope, inputR)
	requireAuxiliaryFailure(t, table, err, 10)
	table, err = ProjectQuota(context.Background(), scope, inputQ[:MaxAuxRows])
	if err != nil || len(table.Rows) != MaxAuxRows {
		t.Fatal("exact quota row bound rejected", err)
	}
	table, err = ProjectReservations(context.Background(), scope, inputR[:MaxAuxRows])
	if err != nil || len(table.Rows) != MaxAuxRows {
		t.Fatal("exact reservation row bound rejected", err)
	}
	for i := range inputR[:MaxAuxRows] {
		for _, column := range []int{0, 2, 3, 4, 5} {
			inputR[i].Cells[column] = strings.Repeat("s", MaxLabelBytes)
		}
		inputR[i].Cells[4] = fmt.Sprintf("%08d", i) + strings.Repeat("s", MaxLabelBytes-8)
	}
	scope[auxSubscription] = strings.Repeat("s", MaxLabelBytes)
	table, err = ProjectReservations(context.Background(), scope, inputR[:MaxAuxRows])
	requireAuxiliaryFailure(t, table, err, 10)
	if table.Health.Error.Code != "region_aux_text_limit" {
		t.Fatal("aggregate output text must fail before projection")
	}
	quotas[0].Current, quotas[0].Limit, quotas[0].Available = MaxAuxCount, MaxAuxCount, -MaxAuxCount
	quotas[0].HeadroomPct = -float64(MaxAuxCount)
	table, err = ProjectQuota(context.Background(), map[string]string{auxSubscription: "Synthetic"}, quotas[:1])
	if err != nil || table.Rows[0].Cells[4] != "1000000000000" || table.Rows[0].Cells[6] != "-1000000000000" || table.Rows[0].Cells[7] != "-1000000000000.0%" {
		t.Fatal("exact numeric bound should preserve source formatting")
	}
	scope = map[string]string{}
	for i := 0; i <= MaxSubscriptions; i++ {
		scope[fmt.Sprintf("%08x-1111-1111-1111-111111111111", i)] = "Synthetic"
	}
	table, err = ProjectQuota(context.Background(), scope, nil)
	requireAuxiliaryFailure(t, table, err, 9)
	table, err = ProjectReservations(context.Background(), scope, nil)
	requireAuxiliaryFailure(t, table, err, 10)
}

// stepContext deterministically cancels during synchronous work, without timing.
type stepContext struct {
	context.Context
	remaining int
}

func (c *stepContext) Err() error {
	c.remaining--
	if c.remaining <= 0 {
		return context.Canceled
	}
	return nil
}

func TestAuxiliaryOwnershipIsolationAndCancellation(t *testing.T) {
	quotas, reservations, _ := auxiliaryFixtures(t)
	scope := map[string]string{auxSubscription: "Synthetic"}
	q, err := ProjectQuota(context.Background(), scope, quotas)
	if err != nil {
		t.Fatal(err)
	}
	r, err := ProjectReservations(context.Background(), scope, reservations)
	if err != nil {
		t.Fatal(err)
	}
	beforeQ, _ := json.Marshal(quotas)
	beforeR, _ := json.Marshal(reservations)
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a, e := ProjectQuota(context.Background(), scope, quotas)
			b, f := ProjectReservations(context.Background(), scope, reservations)
			if e != nil || f != nil || !reflect.DeepEqual(a, q) || !reflect.DeepEqual(b, r) {
				t.Error("per-run auxiliary isolation")
			}
		}()
	}
	wg.Wait()
	afterQ, _ := json.Marshal(quotas)
	afterR, _ := json.Marshal(reservations)
	if string(beforeQ) != string(afterQ) || string(beforeR) != string(afterR) {
		t.Fatal("input mutated")
	}
	for _, stop := range []int{1, 3, 5, 8} {
		a, e := ProjectQuota(&stepContext{Context: context.Background(), remaining: stop}, scope, quotas)
		requireAuxiliaryFailure(t, a, e, 9)
		if !errors.Is(e, context.Canceled) {
			t.Fatal("quota cancellation lost")
		}
		b, f := ProjectReservations(&stepContext{Context: context.Background(), remaining: stop}, scope, reservations)
		requireAuxiliaryFailure(t, b, f, 10)
		if !errors.Is(f, context.Canceled) {
			t.Fatal("reservation cancellation lost")
		}
	}
	reservations[0].Cells[0] = "mutated"
	quotas[0].Subscription = "mutated"
	scope[auxSubscription] = "mutated"
	if r.Rows[0].Cells[0] != "Synthetic" || q.Rows[0].Cells[0] != "Synthetic" {
		t.Fatal("output aliases input")
	}
	r.Columns[0] = "mutated"
	r.Rows[0].Cells[1] = "mutated"
	if pendingAuxiliary(false).Columns[0] != "Subscription" {
		t.Fatal("auxiliary columns escaped")
	}
}
