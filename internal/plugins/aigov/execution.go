package aigov

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// Execution composes the accepted discovery and request libraries. Each call has
// one enclosing five-minute context; library request/byte budgets stay separate.
type Execution struct {
	discovery *Discovery
	scanner   *Scanner
}

func NewExecution(credential azcore.TokenCredential) (*Execution, error) {
	discovery, err := NewDiscovery(credential)
	if err != nil {
		return nil, err
	}
	scanner, err := New(credential)
	if err != nil {
		return nil, err
	}
	return NewExecutionWithScanners(discovery, scanner)
}

// Explicit scanners are trusted application dependencies with existing bounded
// client contracts. The caller supplies an owned, stable filter snapshot.
func NewExecutionWithScanners(discovery *Discovery, scanner *Scanner) (*Execution, error) {
	if discovery == nil || scanner == nil {
		return nil, fmt.Errorf("AI execution scanners required")
	}
	return &Execution{discovery: discovery, scanner: scanner}, nil
}

func (e *Execution) Scan(ctx context.Context, subscriptions map[string]string, filter Filter) (assessment.PluginTable, error) {
	ctx, cancel := context.WithTimeout(ctx, MaxDuration)
	defer cancel()
	scope := make(map[string]string, len(subscriptions))
	for id, name := range subscriptions {
		scope[id] = name
	}
	table := PendingTable()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted}
	if e == nil || e.discovery == nil || e.scanner == nil {
		table.Health.Status = assessment.StageFailed
		table.Health.Error = &assessment.AssessmentError{Code: "ai_not_configured", Message: "AI execution is not configured"}
		return table, fmt.Errorf("AI execution is not configured")
	}
	discovered, discoveryErr := e.discovery.Discover(ctx, scope, filter)
	var scanErr error
	if len(discovered.Accounts) > 0 {
		table, scanErr = e.scanner.Scan(ctx, scope, discovered.Accounts, nil)
		// The source sheet branch depends on selected discovery, even when no
		// metrics could be received. Failed retrieval must retain those headers.
		table.SheetName = "AI Gov"
	}
	table.Health.Warnings = append(table.Health.Warnings, discovered.Health.Warnings...)
	if discoveryErr != nil || discovered.Health.Status == assessment.StageFailed {
		table.Health.Status = assessment.StageFailed
		table.Health.Error = discovered.Health.Error
	} else if table.Health.Status == assessment.StageCompleted && len(table.Health.Warnings) > 0 {
		table.Health.Status = assessment.StageCompletedWithWarnings
	}
	return table, errors.Join(discoveryErr, scanErr)
}
