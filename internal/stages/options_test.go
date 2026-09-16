package stages

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePluginTargetRegions(t *testing.T) {
	got, err := ParseParams([]string{`plugin.target-regions="eastus,westeurope"`})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]any{
		Plugin: {"target-regions": "eastus,westeurope"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parsed params = %#v, want %#v", got, want)
	}
}

func TestApplyParamsStoresAndMergesOptions(t *testing.T) {
	config := NewDefault()
	if err := config.ApplyParams([]string{"plugin.target-regions=eastus"}); err != nil {
		t.Fatal(err)
	}
	if err := config.SetOptions(Plugin, map[string]any{"future-option": true}); err != nil {
		t.Fatal(err)
	}
	got := config.Options(Plugin)
	if got["target-regions"] != "eastus" || got["future-option"] != true {
		t.Fatalf("unexpected stored options: %#v", got)
	}

	got["target-regions"] = "mutated"
	if config.Options(Plugin)["target-regions"] != "eastus" {
		t.Fatal("Options should return a defensive copy")
	}
}

func TestDuplicateStageParamUsesLastValue(t *testing.T) {
	got, err := ParseParams([]string{"plugin.target-regions=eastus", "plugin.target-regions=norwayeast"})
	if err != nil {
		t.Fatal(err)
	}
	if got[Plugin]["target-regions"] != "norwayeast" {
		t.Fatalf("last value did not win: %#v", got)
	}
}

func TestMalformedStageParamReturnsError(t *testing.T) {
	for _, value := range []string{
		"plugin.target-regions",
		"plugin=region",
		".target-regions=eastus",
		"plugin.=eastus",
	} {
		if _, err := ParseParams([]string{value}); err == nil || !strings.Contains(err.Error(), "stage param must be in the form") {
			t.Fatalf("%q: expected format error, got %v", value, err)
		}
	}
}

func TestUnknownStageAndOptionReturnErrors(t *testing.T) {
	if _, err := ParseParams([]string{"advisor.target-regions=eastus"}); err == nil || !strings.Contains(err.Error(), "unknown stage") {
		t.Fatalf("expected unknown-stage error, got %v", err)
	}
	if _, err := ParseParams([]string{"plugin.unknown=value"}); err == nil || !strings.Contains(err.Error(), "unknown option") {
		t.Fatalf("expected unknown-option error, got %v", err)
	}
}

func TestEmptyStageParamsAreIgnored(t *testing.T) {
	got, err := ParseParams([]string{"", "   "})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("empty params should produce no options: %#v", got)
	}
}
