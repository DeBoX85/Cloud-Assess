package plugins

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins/aigov"
)

func aiProjectionFixture(t *testing.T) assessment.PluginTable {
	t.Helper()
	sub := "11111111-1111-4111-8111-111111111111"
	id := "/subscriptions/" + sub + "/resourceGroups/fixture-rg/providers/Microsoft.CognitiveServices/accounts/fixture-ai"
	hour := time.Date(2026, 9, 30, 12, 5, 0, 0, time.UTC)
	count := 1.5
	v, err := aigov.Project(context.Background(), map[string]string{sub: "Synthetic subscription"}, []aigov.Account{{ID: id, SubscriptionID: sub, ResourceGroup: "fixture-rg", Name: "fixture-ai", Kind: "OpenAI", SKU: "S0"}}, []aigov.Point{{ResourceID: id, Timestamp: &hour, Count: &count}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestAIProjectionRejectsForeignIdentityAndOwnsCells(t *testing.T) {
	scope := map[string]string{"11111111-1111-4111-8111-111111111111": "Synthetic subscription"}
	for _, mode := range []string{"foreign", "name", "metadata", "headers", "sheet", "width", "negative", "fractional", "capacity", "spillover", "hour", "control", "oversized", "skipped"} {
		t.Run(mode, func(t *testing.T) {
			v := aiProjectionFixture(t)
			switch mode {
			case "foreign":
				v.Rows[0].SubscriptionID = "22222222-2222-4222-8222-222222222222"
			case "name":
				v.Rows[0].Cells[0] = "other subscription"
			case "metadata":
				v.Metadata.Author = "other"
			case "headers":
				v.Columns[15] = "other"
			case "sheet":
				v.SheetName = "AI Throttling"
			case "width":
				v.Rows[0].Cells = v.Rows[0].Cells[:15]
			case "negative":
				v.Rows[0].Cells[15] = "-1"
			case "fractional":
				v.Rows[0].Cells[15] = "1.5"
			case "capacity":
				v.Rows[0].Cells[9] = "-1"
			case "spillover":
				v.Rows[0].Cells[11] = "maybe"
			case "hour":
				v.Rows[0].Cells[13] = "invalid"
			case "control":
				v.Rows[0].Cells[6] = "bad\x1b"
			case "oversized":
				v.Rows[0].Cells[6] = strings.Repeat("a", 513)
			case "skipped":
				v.Rows = nil
				v.Health.Status = assessment.StageSkipped
				v.Health.Records = 0
			}
			got := AIGovernanceTable(v, nil, scope, time.Unix(1, 0), time.Unix(2, 0))
			if len(got.Rows) != 0 || got.Health.Status != assessment.StageFailed || got.Health.Error.Code != "ai_output_invalid" {
				t.Fatal("unsafe AI injected contract admitted")
			}
		})
	}
	v := aiProjectionFixture(t)
	v.Health.Status = assessment.StageFailed
	v.Health.Error = &assessment.AssessmentError{Code: "provider_secret_canary", Message: "provider_secret_canary"}
	v.Health.Warnings = []assessment.AssessmentWarning{{Code: "provider_secret_canary", Message: "provider_secret_canary"}}
	got := AIGovernanceTable(v, errors.New("provider_secret_canary"), scope, time.Unix(1, 0), time.Unix(2, 0))
	v.Rows[0].Cells[2] = "changed"
	v.Columns[0] = "changed"
	v.Health.Warnings[0].Message = "changed"
	if got.Rows[0].Cells[2] != "fixture-ai" || got.Columns[0] != "Subscription" || got.Health.Error.Code != "ai_request_failed" || strings.Contains(got.Health.Warnings[0].Message, "canary") || got.Health.Records != 1 || !got.Health.StartedAt.Equal(time.Unix(1, 0)) {
		t.Fatal("AI owned/sanitized result lost")
	}
}
