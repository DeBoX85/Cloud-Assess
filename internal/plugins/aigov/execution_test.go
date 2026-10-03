package aigov

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type executionPost func(context.Context, string, io.ReadSeekCloser, int64) ([]byte, *http.Response, error)

func (f executionPost) PostBounded(c context.Context, u string, b io.ReadSeekCloser, n int64) ([]byte, *http.Response, error) {
	return f(c, u, b, n)
}

func TestExecutionDiscoveryFailureSurvivesEnrichment(t *testing.T) {
	for _, failure := range []string{"denied", "malformed", "warnings"} {
		t.Run(failure, func(t *testing.T) {
			var calls atomic.Int64
			d, err := NewDiscoveryWithClient("https://arm.test", executionPost(func(ctx context.Context, _ string, _ io.ReadSeekCloser, _ int64) ([]byte, *http.Response, error) {
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) > MaxDuration || time.Until(deadline) <= 0 {
					t.Error("public execution enclosing deadline absent")
				}
				if calls.Add(1) == 2 {
					if failure == "denied" {
						return nil, nil, errors.New("private provider canary")
					}
					return []byte(`{"count":0}`), nil, nil
				}
				row := `{"id":"` + fixtureID + `","subscriptionId":"` + fixtureSub + `","resourceGroup":"fixture-rg","name":"fixture-ai","location":"westeurope","type":"Microsoft.CognitiveServices/accounts","kind":"OpenAI","sku_name":"S0"}`
				if failure == "warnings" {
					return []byte(`{"count":2,"totalRecords":2,"resultTruncated":"false","data":[` + row + `,{}]}`), nil, nil
				}
				return []byte(`{"count":1,"totalRecords":2,"resultTruncated":"false","$skipToken":"opaque","data":[` + row + `]}`), nil, nil
			}))
			if err != nil {
				t.Fatal(err)
			}
			var h *requestHarness
			h = harness(t, func(r *http.Request) (*http.Response, error) {
				return response(r, 200, metricWire(fixtureID, oneSeries), &h.closed), nil
			}, func(r *http.Request) (*http.Response, error) { return response(r, 200, `{"value":[]}`, &h.closed), nil })
			e, err := NewExecutionWithScanners(d, h.scanner)
			if err != nil {
				t.Fatal(err)
			}
			v, err := e.Scan(context.Background(), fixtureScope(), nil)
			if len(v.Rows) != 1 || v.Rows[0].Cells[15] != "2" || v.SheetName != "AI Gov" || h.metrics.Load() != 1 || h.deployments.Load() != 1 {
				t.Fatal("discovery prefix failed to produce retained enriched source rows")
			}
			if failure == "warnings" {
				if err != nil || v.Health.Status != assessment.StageCompletedWithWarnings || len(v.Health.Warnings) != 1 || v.Health.Warnings[0].Code != "ai_discovery_invalid_rows" {
					t.Fatal("discovery warning overwritten by successful enrichment")
				}
			} else if err == nil || v.Health.Status != assessment.StageFailed || v.Health.Error == nil || !strings.HasPrefix(v.Health.Error.Code, "ai_discovery_") {
				t.Fatal("failed discovery became successful enrichment")
			}
			valid(t, v)
		})
	}
}

func TestExecutionEarlierDeadlineEmptyCancellationAndConcurrentIsolation(t *testing.T) {
	var calls atomic.Int64
	d, _ := NewDiscoveryWithClient("https://arm.test", executionPost(func(ctx context.Context, _ string, _ io.ReadSeekCloser, _ int64) ([]byte, *http.Response, error) {
		calls.Add(1)
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > time.Minute {
			t.Error("earlier enclosing deadline lost")
		}
		return []byte(`{"count":0,"totalRecords":0,"resultTruncated":"false","data":[]}`), nil, nil
	}))
	h := harness(t, func(*http.Request) (*http.Response, error) {
		t.Error("empty authenticated metrics")
		return nil, errors.New("tripwire")
	}, func(*http.Request) (*http.Response, error) {
		t.Error("empty authenticated deployment")
		return nil, errors.New("tripwire")
	})
	e, _ := NewExecutionWithScanners(d, h.scanner)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := e.Scan(ctx, fixtureScope(), nil)
			if err != nil || v.SheetName != "AI Throttling" || len(v.Rows) != 0 || v.Health.Status != assessment.StageCompleted {
				t.Error("empty concurrent public result")
			}
			v.Columns[0] = "mutated"
		}()
	}
	wg.Wait()
	v, err := e.Scan(ctx, nil, nil)
	if err != nil || calls.Load() != 8 || v.Columns[0] != "Subscription" {
		t.Fatal("empty scope or result ownership")
	}
	cancel()
	v, err = e.Scan(ctx, fixtureScope(), nil)
	if !errors.Is(err, context.Canceled) || v.Health.Status != assessment.StageFailed || calls.Load() != 8 {
		t.Fatal("cancelled public execution authenticated or lost identity")
	}
}
