package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/servicehealth"
)

type serviceTransport func(context.Context, arg.Request) (*arg.Response, error)

func (f serviceTransport) Do(c context.Context, r arg.Request) (*arg.Response, error) { return f(c, r) }
func literalServiceTable(t *testing.T) assessment.PluginTable {
	t.Helper()
	s := servicehealth.NewWithTransport(serviceTransport(func(context.Context, arg.Request) (*arg.Response, error) {
		return &arg.Response{Data: []json.RawMessage{json.RawMessage(`{"subscriptionId":"11111111-1111-4111-8111-111111111111","targetRegion":"eastus","targetResourceType":"microsoft.compute/virtualmachines","percentageOfTimeWithoutEvents":99.5,"events":2,"affectedResources":3}`)}}, nil
	}))
	v, err := s.Scan(context.Background(), map[string]string{"11111111-1111-4111-8111-111111111111": "A"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func TestServiceHealthBoundaryOwnershipPrivacyAndInvalid(t *testing.T) {
	scope := map[string]string{"11111111-1111-4111-8111-111111111111": "A"}
	start, end := time.Unix(1, 0), time.Unix(2, 0)
	value := literalServiceTable(t)
	got := ServiceHealthTable(value, nil, scope, start, end)
	value.Rows[0].Cells[1] = "caller changed"
	value.Columns[0] = "caller changed"
	if got.Rows[0].Cells[1] != "eastus" || got.Columns[0] != "Subscription ID" || !got.Health.StartedAt.Equal(start) || !got.Health.FinishedAt.Equal(end) {
		t.Fatal("projection ownership/clock")
	}
	for _, kind := range []string{"identity", "metadata", "source metadata", "foreign", "columns", "empty region", "long type", "control region", "nan", "count", "sheet", "empty false success"} {
		t.Run(kind, func(t *testing.T) {
			value := literalServiceTable(t)
			switch kind {
			case "identity":
				value.ID = "other"
			case "metadata":
				value.Metadata.Name = "zone-mapping"
			case "source metadata":
				value.Metadata.Description = "unreviewed description"
			case "foreign":
				value.Rows[0].SubscriptionID = "22222222-2222-4222-8222-222222222222"
				value.Rows[0].Cells[0] = value.Rows[0].SubscriptionID
			case "columns":
				value.Columns[0] = "other"
			case "empty region":
				value.Rows[0].Cells[1] = ""
			case "long type":
				value.Rows[0].Cells[2] = strings.Repeat("x", 513)
			case "control region":
				value.Rows[0].Cells[1] = "east\nus"
			case "nan":
				value.Rows[0].Cells[3] = "NaN%"
			case "count":
				value.Rows[0].Cells[4] = "-1"
			case "sheet":
				value.SheetName = "unreviewed"
			case "empty false success":
				value = servicehealth.PendingTable()
				value.Health.Status = assessment.StageCompleted
				value.Health.Warnings = nil
			}
			got := ServiceHealthTable(value, nil, scope, start, end)
			if got.Health.Status != assessment.StageFailed || got.Health.Error.Code != "service_health_output_invalid" || len(got.Rows) != 0 {
				t.Fatalf("invalid projection: %+v", got)
			}
		})
	}
	value = literalServiceTable(t)
	value.Health.Status = assessment.StageFailed
	value.Health.Error = &assessment.AssessmentError{Code: "provider_secret", Message: "secret-canary https://evil.invalid"}
	value.Health.Warnings = []assessment.AssessmentWarning{{Code: "provider_secret", Message: "secret-canary"}}
	got = ServiceHealthTable(value, errors.New("secret-canary"), scope, start, end)
	b, _ := json.Marshal(got)
	if strings.Contains(string(b), "secret-canary") || strings.Contains(string(b), "provider_secret") || len(got.Rows) != 1 || got.Health.Status != assessment.StageFailed {
		t.Fatalf("privacy/retention: %s", b)
	}
	names, err := ValidateNames([]string{"zone-mapping", "service-health", "service-health"})
	if err != nil || strings.Join(names, ",") != "service-health,zone-mapping" {
		t.Fatal("source registry order/dedup")
	}
}
