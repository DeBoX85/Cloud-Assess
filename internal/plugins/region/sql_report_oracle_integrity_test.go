package region

import (
	"crypto/sha256"
	"fmt"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"os"
	"testing"
)

func TestSQLReportOracleIntegrity(t *testing.T) {
	for file, want := range map[string]string{
		"sql_report_oracle_test.go":      "ab90d0ea7142f825961117e171a8d56c86cbe6a5d89249c810ef981a9a849cdb",
		"sql_report_budget_plan_test.go": "b9c738407e48412dd15097108c67918ed48007909379f39555a027facd6af95e",
	} {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(raw)) != want {
			t.Fatal("frozen SQL report oracle drift", file)
		}
	}
	cases := sqlReportOracleCases()
	if len(cases) != 26 {
		t.Fatal("frozen report case count", len(cases))
	}
	seen := map[string]bool{}
	for _, c := range cases {
		if c.name == "" || seen[c.name] {
			t.Fatal("frozen report case identity", c.name)
		}
		seen[c.name] = true
		if err := assessment.ValidatePluginTables([]assessment.PluginTable{c.want}); err != nil {
			t.Fatal(c.name, err)
		}
	}
}
