package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/sqleol"
)

func actualSQLTable(t *testing.T) assessment.PluginTable {
	t.Helper()
	b, err := os.ReadFile("sqleol/testdata/source-input.json")
	if err != nil {
		t.Fatal(err)
	}
	var rows []json.RawMessage
	if json.Unmarshal(b, &rows) != nil {
		t.Fatal("source input")
	}
	s := sqleol.NewWithTransport(serviceTransport(func(context.Context, arg.Request) (*arg.Response, error) { return &arg.Response{Data: rows}, nil }))
	v, err := s.Scan(context.Background(), map[string]string{"11111111-1111-4111-8111-111111111111": "A", "22222222-2222-4222-8222-222222222222": "B"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestSQLProjectionOwnershipSourceAndInvalidOutput(t *testing.T) {
	scope := map[string]string{"11111111-1111-4111-8111-111111111111": "A", "22222222-2222-4222-8222-222222222222": "B"}
	start, end := time.Unix(1, 0), time.Unix(2, 0)
	value := actualSQLTable(t)
	got := SQLEOLTable(value, nil, scope, start, end)
	value.Rows[0].Cells[0] = "mutated"
	value.Columns[0] = "mutated"
	if got.Rows[0].Cells[0] != "B" || got.Rows[0].Cells[30] != "-12.34" || got.Columns[0] != "Subscription" || !got.Health.StartedAt.Equal(start) || !got.Health.FinishedAt.Equal(end) {
		t.Fatal("source string/clock/ownership projection")
	}
	for _, kind := range []string{"metadata description", "identity", "headers", "sheet", "foreign", "missing correlation", "control", "limit", "skipped"} {
		t.Run(kind, func(t *testing.T) {
			value := actualSQLTable(t)
			switch kind {
			case "metadata description":
				value.Metadata.Description = "unreviewed"
			case "identity":
				value.ID = "other"
			case "headers":
				value.Columns[1] = "unreviewed"
			case "sheet":
				value.SheetName = "other"
			case "foreign":
				value.Rows[0].SubscriptionID = "33333333-3333-4333-8333-333333333333"
			case "missing correlation":
				value.Rows[0].SubscriptionID = ""
			case "control":
				value.Rows[0].Cells[2] = "bad\nname"
			case "limit":
				value.Rows[0].Cells[2] = strings.Repeat("x", 4097)
			case "skipped":
				value = sqleol.PendingTable()
			}
			got := SQLEOLTable(value, nil, scope, start, end)
			if got.Health.Status != assessment.StageFailed || got.Health.Error.Code != "sql_eol_output_invalid" || len(got.Rows) != 0 || len(got.Columns) != 32 {
				t.Fatalf("invalid SQL projection: %+v", got)
			}
		})
	}
	value = actualSQLTable(t)
	value.Health.Status = assessment.StageFailed
	value.Health.Error = &assessment.AssessmentError{Code: "provider_secret", Message: "provider-secret-canary"}
	value.Health.Warnings = []assessment.AssessmentWarning{{Code: "provider_secret", Message: "provider-secret-canary"}}
	got = SQLEOLTable(value, errors.New("provider-secret-canary"), scope, start, end)
	b, _ := json.Marshal(got)
	if strings.Contains(string(b), "provider-secret-canary") || strings.Contains(string(b), "provider_secret") || len(got.Rows) != 2 || got.Health.Status != assessment.StageFailed {
		t.Fatalf("SQL error privacy/retention: %s", b)
	}
	names, err := ValidateNames([]string{"zone-mapping", "sql-eol", "service-health", "sql-eol"})
	if err != nil || strings.Join(names, ",") != "service-health,sql-eol,zone-mapping" {
		t.Fatal("three adapter source order/dedup")
	}
}
