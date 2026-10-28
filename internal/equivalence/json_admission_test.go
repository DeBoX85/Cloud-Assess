package equivalence

import (
	"strings"
	"testing"
)

// Literal expectations are frozen before the admission implementation. They
// exercise the public loaders, not a proposed helper or its internal layout.
func TestJSONAdmissionRejectsAmbiguousReference(t *testing.T) {
	cases := map[string]string{
		"section":         `{"advisor":[{"impact":"High"}],"advisor":[]}`,
		"escaped-section": `{"advisor":[],"ad\u0076isor":[]}`,
		"row":             `{"advisor":[{"impact":"High","impact":"Low"}]}`,
		"same-value":      `{"advisor":[],"advisor":[]}`,
		"unknown-object":  `{"advisor":[],"extension":{"private-key":1,"private-key":2}}`,
		"unknown-array":   `{"advisor":[],"extension":[{"nested":{"x":1,"x":2}}]}`,
		"empty-name":      `{"advisor":[],"extension":{"":1,"":2}}`,
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := LoadReference(strings.NewReader(input)); err == nil {
				t.Fatal("JSON_ADMISSION_REFERENCE: ambiguous member names accepted")
			} else if strings.Contains(err.Error(), "private-key") {
				t.Fatal("JSON_ADMISSION_PRIVACY: object member name exposed")
			}
		})
	}
}

func TestJSONAdmissionRejectsAmbiguousTarget(t *testing.T) {
	cases := map[string]string{
		"completeness":         `"completeness":"partial","completeness":"complete"`,
		"reverse-completeness": `"completeness":"complete","completeness":"partial"`,
		"same-value":           `"completeness":"complete","completeness":"complete"`,
		"escaped":              `"completeness":"partial","\u0063ompleteness":"complete"`,
		"case-alias":           `"completeness":"partial","COMPLETENESS":"complete"`,
		"unicode-alias":        `"completeness":"complete","\u017fchemaVersion":"1.0"`,
		"array-field":          `"completeness":"complete","advisor":[{"impact":"High","IMPACT":"Low"}]`,
		"pointer-field":        `"completeness":"complete","scope":{"status":"resolved","STATUS":"unknown"}`,
		"nested-warning":       `"completeness":"complete","stages":[{"warnings":[{"code":"a","CODE":"b"}]}]`,
		"plugin-metadata":      `"completeness":"complete","pluginTables":[{"metadata":{"author":"a","AUTHOR":"b"}}]`,
		"opaque-exact":         `"completeness":"complete","resources":[{"tags":{"private-key":"a","private-key":"b"}}]`,
		"unknown-exact":        `"completeness":"complete","extension":{"x":1,"x":2}`,
	}
	for name, fields := range cases {
		t.Run(name, func(t *testing.T) {
			input := `{"schemaVersion":"1.0",` + fields + `}`
			if _, err := LoadTarget(strings.NewReader(input)); err == nil {
				t.Fatal("JSON_ADMISSION_TARGET: duplicate or typed-field alias accepted")
			} else if strings.Contains(err.Error(), "private-key") {
				t.Fatal("JSON_ADMISSION_PRIVACY: object member name exposed")
			}
		})
	}
}

func TestJSONAdmissionPreservesOpaqueMapsAndSingleAliases(t *testing.T) {
	for name, input := range map[string]string{
		"canonical":            `{"schemaVersion":"1.0","completeness":"complete"}`,
		"single-case-alias":    `{"SCHEMAVERSION":"1.0","COMPLETENESS":"complete"}`,
		"single-unicode-alias": `{"\u017fchemaVersion":"1.0","completeness":"complete"}`,
		"opaque-tags":          `{"schemaVersion":"1.0","completeness":"complete","resources":[{"tags":{"Scope":"a","scope":"b"}}]}`,
		"opaque-summary":       `{"schemaVersion":"1.0","completeness":"complete","summary":{"impactedByImpact":{"High":1,"high":2}}}`,
		"unknown-case":         `{"schemaVersion":"1.0","completeness":"complete","Extension":1,"extension":2,"other":{"X":1,"x":2}}`,
		"unknown-number":       `{"schemaVersion":"1.0","completeness":"complete","extension":1e10000}`,
		"separate-objects":     `{"schemaVersion":"1.0","completeness":"complete","extension":[{"x":1},{"x":2}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			p, err := LoadTarget(strings.NewReader(input))
			if err != nil || !p.Comparable {
				t.Fatalf("JSON_ADMISSION_COMPATIBILITY: healthy input rejected: %v", err)
			}
		})
	}
	ref, err := LoadReference(strings.NewReader(`{"advisor":[{"resourceId":"r","recommendationId":"a","impact":"High","Impact":"Low"}],"Advisor":[]}`))
	if err != nil {
		t.Fatal("JSON_ADMISSION_OPAQUE: case-sensitive reference map rejected")
	}
	target, err := LoadTarget(strings.NewReader(`{"schemaVersion":"1.0","completeness":"complete","advisor":[{"resourceId":"r","recommendationId":"a","impact":"High"}]}`))
	if err != nil || !Compare(ref, target).Equivalent {
		t.Fatal("JSON_ADMISSION_OPAQUE: distinct reference map spelling changed selected evidence")
	}
}

func TestJSONAdmissionPreservesExistingSyntaxPolicy(t *testing.T) {
	for _, input := range []string{`{"advisor":[]} {}`, `{"advisor":[],}`, `{"advisor":[`} {
		if _, err := LoadReference(strings.NewReader(input)); err == nil {
			t.Fatal("JSON_ADMISSION_SYNTAX: malformed or trailing JSON accepted")
		}
	}
	// F03 introduces no arbitrary 32/64-container or byte cap. The pinned
	// decoder's existing syntax/depth rules and downstream type checks remain.
	input := `{"advisor":[],"extension":` + strings.Repeat(`[`, 128) + `0` + strings.Repeat(`]`, 128) + `}`
	if _, err := LoadReference(strings.NewReader(input)); err != nil {
		t.Fatalf("JSON_ADMISSION_DEPTH_COMPATIBILITY: existing valid depth rejected: %v", err)
	}
}
