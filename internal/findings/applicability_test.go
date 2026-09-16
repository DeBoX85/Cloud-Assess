package findings

import "testing"

func TestApplicabilityFor(t *testing.T) {
	deployed := map[string]bool{"microsoft.test/widgets": true}
	tests := []struct {
		name     string
		typeName string
		impacted int
		want     Applicability
	}{
		{"absent type", "Microsoft.Other/things", 0, NotApplicable},
		{"deployed compliant", "Microsoft.Test/widgets", 0, Compliant},
		{"deployed noncompliant", "Microsoft.Test/widgets", 2, NonCompliant},
		{"microsoft resources always applicable", "Microsoft.Resources", 0, Compliant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplicabilityFor(tt.typeName, tt.impacted, deployed); got != tt.want {
				t.Fatalf("ApplicabilityFor() = %q, want %q", got, tt.want)
			}
		})
	}
}
