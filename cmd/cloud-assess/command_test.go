package main

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestScanCommandMapsFlags(t *testing.T) {
	var got scanFlags
	root := newRootCommand(func(_ context.Context, flags scanFlags) (int, error) {
		got = flags
		return 0, nil
	})
	root.SetArgs([]string{
		"scan",
		"--subscription-id", "sub-1",
		"--resource-group", "rg-1",
		"--stages", "cost,-advisor",
		"--stage-param", "plugin.target-regions=norwayeast,swedencentral",
		"--json", "--csv", "--sarif", "--stdout",
		"--xlsx=false",
		"--output-name", "assessment",
		"--redact-subscription-ids=false",
		"--filters", "filters.yaml",
		"--fail-on", "Medium",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(got.subscriptions) != 1 || got.subscriptions[0] != "sub-1" {
		t.Fatalf("subscriptions = %#v", got.subscriptions)
	}
	if len(got.resourceGroups) != 1 || got.resourceGroups[0] != "rg-1" {
		t.Fatalf("resource groups = %#v", got.resourceGroups)
	}
	if got.xlsx || !got.json || !got.csv || !got.sarif || !got.stdout {
		t.Fatalf("output flags were not mapped: %#v", got)
	}
	if got.outputName != "assessment" || got.redactSubscriptionIDs || got.filtersFile != "filters.yaml" || got.failOn != "Medium" {
		t.Fatalf("scalar flags were not mapped: %#v", got)
	}
	if len(got.stageNames) != 2 || got.stageNames[0] != "cost" || got.stageNames[1] != "-advisor" {
		t.Fatalf("stage names = %#v", got.stageNames)
	}
	if len(got.stageParams) != 1 {
		t.Fatalf("stage params = %#v", got.stageParams)
	}
}

func TestRootPropagatesApplicationExitCode(t *testing.T) {
	root := newRootCommand(func(context.Context, scanFlags) (int, error) {
		return 2, errors.New("quality gate failed")
	})
	root.SetArgs([]string{"scan"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected command error")
	}
	value, ok := root.Context().Value(exitCodeContextKey{}).(*int)
	if !ok || value == nil || *value != 2 {
		t.Fatalf("stored exit code = %#v, want 2", value)
	}
}

func TestScanDefaultsToExcelAndRedaction(t *testing.T) {
	var got scanFlags
	root := newRootCommand(func(_ context.Context, flags scanFlags) (int, error) {
		got = flags
		return 0, nil
	})
	root.SetArgs([]string{"scan"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !got.xlsx {
		t.Fatal("Excel should be enabled by default")
	}
	if !got.redactSubscriptionIDs {
		t.Fatal("subscription ID redaction should be enabled by default")
	}
}

func TestExecuteScanRejectsDeferredPluginStageBeforeAzureAuthentication(t *testing.T) {
	code, err := executeScan(context.Background(), scanFlags{stageNames: []string{"plugin"}})
	if err == nil {
		t.Fatal("expected plugin availability error")
	}
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(err.Error(), "plugin stage is not available") {
		t.Fatalf("unexpected error: %v", err)
	}
}
