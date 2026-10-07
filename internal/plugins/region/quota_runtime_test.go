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

const quotaTestID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

func quotaScope() map[string]string { return map[string]string{quotaTestID: "Selected"} }
func quotaRequest(provider string) QuotaRequest {
	return QuotaRequest{SubscriptionID: quotaTestID, Region: "eastus", QuotaType: provider}
}
func quotaUsage(name string, current, limit int64) QuotaUsage {
	return QuotaUsage{ResourceName: name, Current: current, Limit: limit, CurrentKnown: true, LimitKnown: true}
}
func calculateTestQuota(usages []QuotaUsage) (*QuotaCalculation, error) {
	request := quotaRequest("Network")
	return CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "complete", Usages: usages}})
}

func TestQuotaRuntimeArithmeticAndDisplay(t *testing.T) {
	input := []QuotaUsage{quotaUsage("exact", 85, 100), quotaUsage("below", 86, 100), quotaUsage("at", 100, 100), quotaUsage("over", 110, 100)}
	input[0].LocalizedName, input[1].LocalizedName = "Shared label", "Shared label"
	got, err := calculateTestQuota(input)
	if err != nil {
		t.Fatal(err)
	}
	want := []QuotaRow{
		{SubscriptionID: quotaTestID, Subscription: "Selected", Region: "eastus", QuotaType: "Network", ResourceName: "exact", DisplayName: "Shared label", Current: 85, Limit: 100, Available: 15, HeadroomPct: 15},
		{SubscriptionID: quotaTestID, Subscription: "Selected", Region: "eastus", QuotaType: "Network", ResourceName: "below", DisplayName: "Shared label", Current: 86, Limit: 100, Available: 14, HeadroomPct: 14.000000000000002, IsNearLimit: true},
		{SubscriptionID: quotaTestID, Subscription: "Selected", Region: "eastus", QuotaType: "Network", ResourceName: "at", DisplayName: "at", Current: 100, Limit: 100, Available: 0, IsNearLimit: true, IsOverLimit: true},
		{SubscriptionID: quotaTestID, Subscription: "Selected", Region: "eastus", QuotaType: "Network", ResourceName: "over", DisplayName: "over", Current: 110, Limit: 100, Available: -10, HeadroomPct: -10, IsNearLimit: true, IsOverLimit: true},
	}
	if !reflect.DeepEqual(got.Rows, want) {
		t.Fatalf("quota arithmetic rows: %#v", got.Rows)
	}
	wantCells := [][]string{
		{"Selected", "eastus", "Network", "Shared label", "85", "100", "15", "15.0%", "OK"},
		{"Selected", "eastus", "Network", "Shared label", "86", "100", "14", "14.0%", "Near Limit"},
		{"Selected", "eastus", "Network", "at", "100", "100", "0", "0.0%", "At/Over Limit"},
		{"Selected", "eastus", "Network", "over", "110", "100", "-10", "-10.0%", "At/Over Limit"},
	}
	if got.Table == nil || got.Health.Status != assessment.StageCompleted || got.Health.Records != 4 || !reflect.DeepEqual(got.Table.Health, got.Health) {
		t.Fatalf("quota healthy projection: %#v", got)
	}
	for i, row := range got.Table.Rows {
		if row.SubscriptionID != quotaTestID || !reflect.DeepEqual(row.Cells, wantCells[i]) {
			t.Fatalf("quota table row %d: %#v", i, row)
		}
	}
	old := want[0]
	old.DisplayName = ""
	table, err := ProjectQuota(context.Background(), quotaScope(), []QuotaRow{old})
	if err != nil || table.Rows[0].Cells[3] != "exact" {
		t.Fatalf("legacy display compatibility: %v %#v", err, table)
	}
}

func TestQuotaRuntimeRetainedSourcePositiveRows(t *testing.T) {
	type sourceRow struct {
		ResourceName, LocalizedName    string
		CurrentValue, Limit, Available int64
		HeadroomPct                    float64
		IsNearLimit, IsAtOrOverLimit   bool
	}
	for _, fixture := range []struct{ file, provider string }{{"source-quota-outputs.json", "Network"}, {"source-vm-quota-outputs.json", "VM"}} {
		raw, err := os.ReadFile("testdata/" + fixture.file)
		if err != nil {
			t.Fatal(err)
		}
		var results map[string]struct{ Entries []sourceRow }
		if err := json.Unmarshal(raw, &results); err != nil {
			t.Fatal(err)
		}
		matched := 0
		for scenario, result := range results {
			provider := fixture.provider
			switch scenario {
			case "sql-filter", "names":
				provider = "SQL"
			case "storage-filter":
				provider = "Storage"
			case "web-filter":
				provider = "App Service"
			}
			for _, source := range result.Entries {
				// Source unsafe negative usage and absent identity are deliberate target
				// corrections, tested separately, not equivalence claims.
				if source.CurrentValue < 0 || source.ResourceName == "" {
					continue
				}
				request := quotaRequest(provider)
				usage := quotaUsage(source.ResourceName, source.CurrentValue, source.Limit)
				usage.LocalizedName = source.LocalizedName
				got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "complete", Usages: []QuotaUsage{usage}}})
				if err != nil || len(got.Rows) != 1 {
					t.Fatalf("source %s/%s: %v %#v", fixture.file, scenario, err, got)
				}
				row := got.Rows[0]
				if row.ResourceName != source.ResourceName || row.Current != source.CurrentValue || row.Limit != source.Limit || row.Available != source.Available || row.HeadroomPct != source.HeadroomPct || row.IsNearLimit != source.IsNearLimit || row.IsOverLimit != source.IsAtOrOverLimit {
					t.Fatalf("retained source full calculation mismatch %s/%s: %#v vs %#v", fixture.file, scenario, row, source)
				}
				label := source.LocalizedName
				if label == "" {
					label = source.ResourceName
				}
				if row.DisplayName != label || got.Table.Rows[0].Cells[3] != label {
					t.Fatal("retained source display mismatch")
				}
				matched++
			}
		}
		if matched == 0 {
			t.Fatal("no retained source positive rows compared")
		}
	}
}

func TestQuotaRuntimeHealthAndOwnership(t *testing.T) {
	for _, status := range []string{"complete", "partial", "unknown", "unsupported", "missing"} {
		t.Run(status, func(t *testing.T) {
			request := quotaRequest("Network")
			var evidence []QuotaEvidence
			if status != "missing" {
				evidence = []QuotaEvidence{{Request: request, Status: status}}
			}
			got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, evidence)
			if err != nil || got.Table != nil || len(got.Rows) != 0 || got.Health.Records != 0 {
				t.Fatalf("empty evidence: %v %#v", err, got)
			}
			if status == "complete" {
				if got.Health.Status != assessment.StageCompleted || len(got.Health.Warnings) != 0 {
					t.Fatal("complete empty is not healthy")
				}
				return
			}
			if got.Health.Status != assessment.StageCompletedWithWarnings || len(got.Health.Warnings) != 1 || got.Health.Warnings[0].Code != "quota_evidence_"+status {
				t.Fatalf("missing health: %#v", got.Health)
			}
		})
	}
	request := quotaRequest("Network")
	input := []QuotaEvidence{{Request: request, Status: "partial", Usages: []QuotaUsage{quotaUsage("PublicIPAddresses", 9, 10)}}}
	got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, input)
	if err != nil || got.Health.Status != assessment.StageCompletedWithWarnings || got.Health.Records != 1 || !reflect.DeepEqual(got.Table.Health, got.Health) {
		t.Fatalf("partial projection health: %v %#v", err, got)
	}
	input[0].Usages[0].ResourceName = "mutated"
	got.Health.Warnings[0].Code = "mutated"
	got.Rows[0].DisplayName = "mutated"
	if got.Table.Health.Warnings[0].Code != "quota_evidence_partial" || got.Table.Rows[0].Cells[3] != "PublicIPAddresses" {
		t.Fatal("quota ownership alias")
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repeated, err := calculateTestQuota([]QuotaUsage{quotaUsage("PublicIPAddresses", 9, 10)})
			if err != nil || repeated.Rows[0].ResourceName != "PublicIPAddresses" || repeated.Health.Status != assessment.StageCompleted {
				t.Errorf("quota per-call state: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestQuotaRuntimeProviderFilters(t *testing.T) {
	cases := []struct {
		provider           string
		excluded, included []string
	}{
		{"VM", []string{"cores", "standardfamily"}, []string{"standardDSFamily", "NotFamilySuffix"}},
		{"Network", []string{"NetworkWatchers", "RouteFilterRulesPerRouteFilter", "RouteFiltersPerExpressRouteBgpPeering", "RoutesPerExpressRouteCircuit", "BgpCommunityFilterRulesPerRouteFilter"}, []string{"networkwatchers", "PublicIPAddresses"}},
		{"SQL", []string{"dtusPerServer", "unitsPerDatabase"}, []string{"dtusPerserver", "Servers"}},
		{"App Service", []string{"CustomDomains", "HostNameBindings", "SslBindings", "SslConnections", "Certificates"}, []string{"customdomains", "Workers"}},
		{"Storage", []string{"TotalBlobContainers", "TotalBlobs", "TotalContainers", "TotalFileShares", "TotalQueues", "TotalTables"}, []string{"totalblobs", "StorageAccounts"}},
	}
	for _, tc := range cases {
		t.Run(tc.provider, func(t *testing.T) {
			request := quotaRequest(tc.provider)
			var usages []QuotaUsage
			for _, name := range append(append([]string{}, tc.excluded...), tc.included...) {
				usages = append(usages, quotaUsage(name, 1, 10))
			}
			got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "complete", Usages: usages}})
			if err != nil || len(got.Rows) != len(tc.included) {
				t.Fatalf("source filter: %v %#v", err, got)
			}
			for i, row := range got.Rows {
				if row.ResourceName != tc.included[i] {
					t.Fatal("source filter/order mismatch")
				}
			}
		})
	}
}

func TestQuotaRuntimeAdmissionCorrectionsAndBounds(t *testing.T) {
	for _, usage := range []QuotaUsage{
		quotaUsage("x", -1, 100), quotaUsage("x", MaxAuxCount+1, 100), quotaUsage("x", 1, MaxAuxCount+1), quotaUsage("x", 1, -MaxAuxCount-1),
		quotaUsage("x", MaxAuxCount, 1), quotaUsage(strings.Repeat("x", MaxLabelBytes+1), 1, 10), quotaUsage("bad\n", 1, 10), quotaUsage(string([]byte{255}), 1, 10),
		{ResourceName: "x", Current: 1, Limit: 10, LimitKnown: true}, {ResourceName: "x", Limit: 1, CurrentKnown: true},
		quotaUsage("NetworkWatchers", -1, 100),
	} {
		got, err := calculateTestQuota([]QuotaUsage{usage})
		if err == nil || got != nil {
			t.Fatalf("unsafe admission: %#v %v", got, err)
		}
	}
	for _, usage := range []QuotaUsage{quotaUsage("", 1, 10), {ResourceName: "x", Limit: 10, LimitKnown: true}, quotaUsage("x", 1, 0), quotaUsage("x", 1, -1)} {
		got, err := calculateTestQuota([]QuotaUsage{usage})
		if err != nil || got.Table != nil || len(got.Rows) != 0 || got.Health.Status != assessment.StageCompletedWithWarnings {
			t.Fatalf("unknown/unusable correction: %v %#v", err, got)
		}
	}
	for _, usage := range []QuotaUsage{quotaUsage(strings.Repeat("x", MaxLabelBytes), 0, MaxAuxCount), quotaUsage("x", 10000000001, 1), quotaUsage("x", MaxAuxCount, MaxAuxCount)} {
		got, err := calculateTestQuota([]QuotaUsage{usage})
		if err != nil || len(got.Rows) != 1 {
			t.Fatalf("valid bound rejected: %v", err)
		}
	}
	got, err := calculateTestQuota([]QuotaUsage{quotaUsage("x", 10000000002, 1)})
	if err == nil || got != nil {
		t.Fatal("ratio bound bypass")
	}
	got, err = calculateTestQuota([]QuotaUsage{quotaUsage("x", 1, 10), quotaUsage("x", 2, 10)})
	if err == nil || got != nil {
		t.Fatal("duplicate raw quota identity accepted")
	}
	usages := make([]QuotaUsage, MaxAuxRows)
	got, err = calculateTestQuota([]QuotaUsage{quotaUsage("NetworkWatchers", 1, 10), quotaUsage("NetworkWatchers", 2, 10)})
	if err == nil || got != nil {
		t.Fatal("filtered duplicate raw identity accepted")
	}
	for i := range usages {
		usages[i] = quotaUsage(fmt.Sprintf("cores%d", i), 1, 10)
	}
	request := quotaRequest("VM")
	got, err = CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "complete", Usages: usages}})
	if err != nil || len(got.Rows) != 0 || got.Health.Status != assessment.StageCompleted {
		t.Fatalf("exact filtered work limit: %v", err)
	}
	usages = append(usages, quotaUsage("cores-extra", 1, 10))
	got, err = CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "complete", Usages: usages}})
	if err == nil || got != nil {
		t.Fatal("filtered input work limit bypass")
	}
	usages = usages[:1]
	usages[0].LocalizedName = strings.Repeat("x", MaxLabelBytes+1)
	got, err = CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "complete", Usages: usages}})
	if err == nil || got != nil {
		t.Fatal("filtered text limit bypass")
	}
}

func TestQuotaRuntimeScopeRequestsAndCancellation(t *testing.T) {
	request := quotaRequest("Network")
	valid := QuotaEvidence{Request: request, Status: "complete", Usages: []QuotaUsage{quotaUsage("x", 1, 10)}}
	for _, tc := range []struct {
		scope    map[string]string
		requests []QuotaRequest
		evidence []QuotaEvidence
	}{
		{map[string]string{}, []QuotaRequest{request}, []QuotaEvidence{valid}},
		{map[string]string{}, []QuotaRequest{request}, nil},
		{map[string]string{quotaTestID: ""}, []QuotaRequest{request}, []QuotaEvidence{valid}},
		{map[string]string{quotaTestID: "a", strings.ToUpper(quotaTestID): "a"}, []QuotaRequest{request}, []QuotaEvidence{valid}},
		{quotaScope(), []QuotaRequest{request, request}, nil},
		{quotaScope(), []QuotaRequest{request}, []QuotaEvidence{valid, valid}},
		{quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: quotaRequest("VM"), Status: "complete"}}},
		{quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "invalid"}}},
		{quotaScope(), []QuotaRequest{request}, []QuotaEvidence{{Request: request, Status: "unknown", Usages: valid.Usages}}},
	} {
		got, err := CalculateQuota(context.Background(), tc.scope, tc.requests, tc.evidence)
		if err == nil || got != nil {
			t.Fatal("quota request/scope guard bypass")
		}
	}
	for _, bad := range []QuotaRequest{{SubscriptionID: quotaTestID, Region: "EastUS", QuotaType: "Network"}, {SubscriptionID: quotaTestID, Region: "eastus", QuotaType: "AppService"}, {SubscriptionID: "bad", Region: "eastus", QuotaType: "Network"}} {
		got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{bad}, nil)
		if err == nil || got != nil {
			t.Fatal("quota canonical request guard bypass")
		}
	}
	upper := request
	upper.SubscriptionID = strings.ToUpper(quotaTestID)
	got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{upper}, []QuotaEvidence{valid})
	if err != nil || got.Rows[0].SubscriptionID != quotaTestID {
		t.Fatalf("canonical UUID: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err = CalculateQuota(ctx, quotaScope(), []QuotaRequest{request}, []QuotaEvidence{valid})
	if got != nil || !errors.Is(err, context.Canceled) {
		t.Fatal("quota cancellation bypass")
	}
	probe := &availabilityCancelContext{Context: context.Background(), at: 10000}
	if _, err := CalculateQuota(probe, quotaScope(), []QuotaRequest{request}, []QuotaEvidence{valid}); err != nil {
		t.Fatal(err)
	}
	for at := 1; at <= probe.calls; at++ {
		ctx := &availabilityCancelContext{Context: context.Background(), at: at}
		got, err := CalculateQuota(ctx, quotaScope(), []QuotaRequest{request}, []QuotaEvidence{valid})
		if got != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("quota deterministic cancellation at %d: %v %#v", at, err, got)
		}
	}
}

func TestQuotaRuntimeDeclaredOrderAndAggregateWork(t *testing.T) {
	network, sql := quotaRequest("Network"), quotaRequest("SQL")
	evidence := []QuotaEvidence{
		{Request: sql, Status: "complete", Usages: []QuotaUsage{quotaUsage("Servers", 1, 10)}},
		{Request: network, Status: "partial", Usages: []QuotaUsage{quotaUsage("PublicIPAddresses", 1, 10), {ResourceName: "missing1"}, {ResourceName: "missing2"}}},
	}
	got, err := CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{network, sql}, evidence)
	if err != nil || len(got.Rows) != 2 || got.Rows[0].ResourceName != "PublicIPAddresses" || got.Rows[1].ResourceName != "Servers" {
		t.Fatalf("declared query order: %v %#v", err, got)
	}
	wantHealth := assessment.StageExecution{Name: Name, Status: assessment.StageCompletedWithWarnings, Records: 2, Warnings: []assessment.AssessmentWarning{
		{Code: "quota_evidence_partial", Message: "region quota evidence is incomplete or unusable"},
		{Code: "quota_usage_unknown", Message: "region quota evidence is incomplete or unusable"},
	}}
	if !reflect.DeepEqual(got.Health, wantHealth) || !reflect.DeepEqual(got.Table.Health, wantHealth) {
		t.Fatalf("deduplicated ordered health: %#v", got.Health)
	}
	requests := make([]QuotaRequest, MaxAuxRows+1)
	got, err = CalculateQuota(context.Background(), quotaScope(), requests, nil)
	if err == nil || got != nil {
		t.Fatal("request admission limit bypass")
	}
	for i := range evidence {
		evidence[i].Usages = make([]QuotaUsage, MaxAuxRows/2+1)
		for j := range evidence[i].Usages {
			evidence[i].Usages[j] = quotaUsage(fmt.Sprintf("quota%d", j), 1, 10)
		}
	}
	got, err = CalculateQuota(context.Background(), quotaScope(), []QuotaRequest{network, sql}, evidence)
	if err == nil || got != nil {
		t.Fatal("aggregate usage limit bypass")
	}
}
