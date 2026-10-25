package region

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func sqlReportInputSnapshot(t *testing.T, c sqlReportOracleCase) []byte {
	t.Helper()
	b, e := json.Marshal(map[string]any{"selected": c.selected, "request": c.request, "collection": c.collection, "required": c.required})
	if e != nil {
		t.Fatal(e)
	}
	return b
}

// This runner is published with BOTH report APIs absent. Literal expectations,
// every metadata/health/row field, errors, input preservation and output ownership
// are checked without deriving expectations from report implementation.
func TestSQLReportFrozenOracle(t *testing.T) {
	for _, c := range sqlReportOracleCases() {
		t.Run(c.name, func(t *testing.T) {
			before := sqlReportInputSnapshot(t, c)
			got, e := ProjectSQLCapabilityTable(context.Background(), c.selected, c.request, c.collection, c.required)
			message := ""
			if e != nil {
				message = e.Error()
			}
			if message != c.wantError || !reflect.DeepEqual(got, c.want) {
				t.Fatalf("frozen report full output/error: got=%+v error=%q want=%+v error=%q", got, message, c.want, c.wantError)
			}
			if !bytes.Equal(before, sqlReportInputSnapshot(t, c)) {
				t.Fatal("report mutated input")
			}
			got.Columns[0] = "mutated"
			if len(got.Rows) > 0 {
				got.Rows[0].Cells[0] = "mutated"
			}
			if len(got.Health.Warnings) > 0 {
				got.Health.Warnings[0].Message = "mutated"
			}
			if got.Health.Error != nil {
				got.Health.Error.Message = "mutated"
			}
			again, e := ProjectSQLCapabilityTable(context.Background(), c.selected, c.request, c.collection, c.required)
			message = ""
			if e != nil {
				message = e.Error()
			}
			if message != c.wantError || !reflect.DeepEqual(again, c.want) || !bytes.Equal(before, sqlReportInputSnapshot(t, c)) {
				t.Fatal("report output ownership/repeat", again, e)
			}
		})
	}
}
func TestSQLReportCollectFrozenSmoke(t *testing.T) {
	oracle := sqlReportOracleCases()
	for _, test := range []struct {
		name, kind, key, body string
		constructor           func(string, RESTQuotaGetter) (*SQLCapabilityCollector, error)
		want                  sqlReportOracleCase
	}{
		{"db", "microsoft.sql/servers/databases", "12/basic/o1", "{\"name\":\"eastus\",\"status\":\"Available\",\"supportedServerVersions\":[{\"name\":\"12\",\"status\":\"Available\",\"supportedEditions\":[{\"name\":\"basic\",\"status\":\"Available\",\"supportedServiceLevelObjectives\":[{\"name\":\"o1\",\"status\":\"Available\",\"sku\":{\"name\":\"same_sku\",\"tier\":\"basic\",\"capacity\":1}}]}]}]}", NewSQLDatabaseSKUCollector, oracle[2]},
		{"mi", "microsoft.sql/managedinstances", "12/generalpurpose/gen5/c8", "{\"name\":\"eastus\",\"status\":\"Available\",\"supportedManagedInstanceVersions\":[{\"name\":\"12\",\"status\":\"Available\",\"supportedEditions\":[{\"name\":\"GeneralPurpose\",\"status\":\"Available\",\"supportedFamilies\":[{\"name\":\"Gen5\",\"status\":\"Available\",\"sku\":\"explicit_sku\",\"supportedVcoresValues\":[{\"name\":\"c8\",\"value\":8,\"status\":\"Available\",\"standaloneSupported\":false,\"instancePoolSupported\":false}]}]}]}]}", NewSQLManagedInstanceSKUCollector, oracle[1]},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			selected := map[string]string{sqlReportOracleID: "selected"}
			q := SQLCapabilityRequest{SubscriptionID: sqlReportOracleID, Region: "eastus", ResourceType: test.kind}
			collector, e := test.constructor("https://management.azure.com", storageGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
				calls++
				if u != "https://management.azure.com/subscriptions/"+sqlReportOracleID+"/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01" || limit != 1048576 {
					t.Fatal("exact one bounded report GET", u, limit)
				}
				return []byte(test.body), &http.Response{StatusCode: http.StatusOK}, nil
			}))
			if e != nil {
				t.Fatal(e)
			}
			want := test.want.want
			if test.name == "mi" {
				want = sqlReportExpected("mi", [][]string{oracle[1].want.Rows[0].Cells}, false, "")
			}
			got, e := CollectSQLCapabilityTable(context.Background(), collector, selected, q, []string{test.key})
			if e != nil || calls != 1 || !reflect.DeepEqual(got, want) {
				t.Fatal("real collector->Catalog->report", got, e, calls, want)
			}
		})
	}
}
