package plugins

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/zone"
)

func TestZoneProjectionSourceMetadataAndOwnership(t *testing.T) {
	wantColumns := []string{"Subscription", "Location", "Display Name", "Logical Zone", "Physical Zone"}
	source := zone.Result{Rows: []zone.Row{{SubscriptionID: "11111111-1111-4111-8111-111111111111", SubscriptionName: "Dev", Location: "westus", DisplayName: "West US", LogicalZone: "1", PhysicalZone: "westus-az1"}}}
	table := ZoneTable(source, nil, time.Unix(1, 0), time.Unix(2, 0))
	if !reflect.DeepEqual(table.Columns, wantColumns) || table.Description != "Logical-to-physical availability zone mappings for Azure regions by subscription" || table.Metadata.Description != "Retrieves logical-to-physical availability zone mappings for all Azure regions in each subscription" || table.Metadata.Author != "Azure Quick Review Team" || table.Metadata.Version != "1.0.0" || table.Metadata.License != "MIT" || table.SheetName != "Zone Mapping" {
		t.Fatalf("source metadata differs: %#v", table)
	}
	if table.Health.Status != assessment.StageCompleted || table.Health.Records != 1 || !reflect.DeepEqual(table.Rows[0].Cells, []string{"Dev", "westus", "West US", "1", "westus-az1"}) {
		t.Fatal("source rows/health differ")
	}
	source.Rows[0].SubscriptionName = "mutated"
	table.Columns[0] = "mutated"
	table.Rows[0].Cells[0] = "mutated"
	meta := InternalMetadata()
	meta[0].Name = "mutated"
	if PendingZoneTable().Columns[0] != "Subscription" || InternalMetadata()[0].Name != ZoneMapping {
		t.Fatal("registry state leaked")
	}
}

func TestZoneProjectionFailurePrivacyAndInvalidOutput(t *testing.T) {
	for _, value := range []zone.Result{
		{Failures: []zone.Failure{{SubscriptionID: "https://secret.example/canary", Code: "raw secret canary"}}},
		{Rows: []zone.Row{{SubscriptionID: "bad", SubscriptionName: "canary"}}},
		{Failures: make([]zone.Failure, zone.MaxSubscriptions+1)},
	} {
		table := ZoneTable(value, errors.New("secret canary provider body"), time.Time{}, time.Time{})
		if table.Health.Status != assessment.StageFailed || assessment.ValidatePluginTables([]assessment.PluginTable{table}) != nil {
			t.Fatalf("invalid failure health: %#v", table)
		}
		encoded, _ := json.Marshal(table)
		if strings.Contains(string(encoded), "canary") || strings.Contains(string(encoded), "secret.example") {
			t.Fatal("provider error leaked")
		}
	}
	value := zone.Result{Rows: []zone.Row{{SubscriptionID: "11111111-1111-4111-8111-111111111111", SubscriptionName: "Dev"}}, Failures: []zone.Failure{{SubscriptionID: "22222222-2222-4222-8222-222222222222", Code: "zone_request_failed"}}}
	table := ZoneTable(value, nil, time.Time{}, time.Time{})
	if len(table.Rows) != 1 || table.Health.Records != 1 || table.Health.Error.Code != "zone_incomplete" || !strings.Contains(table.Health.Warnings[0].Message, "22222222-2222-4222-8222-222222222222") {
		t.Fatal("healthy rows or failure correlation lost")
	}
}
