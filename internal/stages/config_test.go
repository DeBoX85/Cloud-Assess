package stages

import "testing"

func TestDefaults(t *testing.T) {
	c := NewDefault()
	for _, name := range []string{Graph, Diagnostics, Advisor, Defender} {
		if !c.IsEnabled(name) {
			t.Fatalf("expected %s enabled by default", name)
		}
	}
	for _, name := range []string{DefenderRecommendations, Arc, Policy, Cost, Plugin} {
		if c.IsEnabled(name) {
			t.Fatalf("expected %s disabled by default", name)
		}
	}
}

func TestApplyEnableDisableAndCommaSeparatedValues(t *testing.T) {
	c := NewDefault()
	if err := c.Apply([]string{"cost,policy", "-advisor"}); err != nil {
		t.Fatal(err)
	}
	if !c.IsEnabled(Cost) || !c.IsEnabled(Policy) || c.IsEnabled(Advisor) {
		t.Fatalf("unexpected stage state: %v", c.EnabledStages())
	}
}

func TestUnknownStageReturnsError(t *testing.T) {
	c := NewDefault()
	if err := c.Apply([]string{"does-not-exist"}); err == nil {
		t.Fatal("expected invalid stage error")
	}
}

func TestGraphIsMandatory(t *testing.T) {
	c := NewDefault()
	if err := c.Apply([]string{"-graph"}); err != nil {
		t.Fatal(err)
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected graph validation error")
	}
}
