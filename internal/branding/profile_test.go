package branding

import (
	"encoding/base64"
	"strings"
	"sync"
	"testing"
)

func TestProfileDefaultsAndCanonical(t *testing.T) {
	p, err := ParseProfile([]byte(`{"schemaVersion":1}`))
	if err != nil {
		t.Fatal(err)
	}
	if p.Branding() != defaults() {
		t.Fatal("defaults changed")
	}
	data, err := p.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	again, err := ParseProfile(data)
	if err != nil || again != p {
		t.Fatalf("round trip %v", err)
	}
	data2, _ := again.Canonical()
	if string(data) != string(data2) {
		t.Fatal("not canonical")
	}
}

func TestProfileRejectsInvalidInput(t *testing.T) {
	cases := []string{
		`null`, `[]`, `{}`, `{"schemaVersion":null}`, `{"schemaVersion":1.0}`, `{"schemaVersion":2}`,
		`{"schemaVersion":1,"schemaVersion":1}`, `{"schemaVersion":1,"cliName":"safe","cliName":"other"}`,
		`{"schemaVersion":1,"CLIName":"name"}`, `{"schemaVersion":1,"companyName":"Unused"}`,
		`{"schemaVersion":1,"cliName":null}`, `{"schemaVersion":1,"cliName":1}`,
		`{"schemaVersion":1} {}`, `{"schemaVersion":1,"productName":"\nSecret"}`,
		`{"schemaVersion":1,"productName":"\ud800"}`, `{"schemaVersion":1,"productName":"a\u202eb"}`,
		`{"schemaVersion":1,"productName":" Leading"}`, `{"schemaVersion":1,"reportTitle":""}`,
		`{"schemaVersion":1,"websiteURL":"http://example.test"}`, `{"schemaVersion":1,"websiteURL":"https://user:password@example.test"}`,
		`{"schemaVersion":1,"websiteURL":"https://"}`, `{"schemaVersion":1,"websiteURL":"https://example.test\n"}`,
		"{\"schemaVersion\":1,\"productName\":\"\xff\"}", strings.Repeat(" ", MaxProfileBytes+1),
	}
	for _, input := range cases {
		if _, err := ParseProfile([]byte(input)); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
	for _, name := range []string{"../x", "..", ".", "a/b", `a\b`, "C:drive", "nul", "con", "com1", "lpt0", "trailing.", "trailing ", "-flag", "Upper", strings.Repeat("a", 65)} {
		for _, key := range []string{"cliName", "reportFilePrefix"} {
			input := `{"schemaVersion":1,"` + key + `":"` + strings.ReplaceAll(name, `\`, `\\`) + `"}`
			if _, err := ParseProfile([]byte(input)); err == nil {
				t.Errorf("accepted %s=%q", key, name)
			}
		}
	}
}

func TestEmbeddedProfileAndIsolation(t *testing.T) {
	p, err := ParseProfile([]byte(`{"schemaVersion":1,"productName":"Example Cloud Toolkit","cliName":"example-cloud","reportTitle":"Example Assessment","reportFilePrefix":"example_report","websiteURL":"https://example.test/tool"}`))
	if err != nil {
		t.Fatal(err)
	}
	data, _ := p.Canonical()
	b, err := resolveEmbedded(base64.RawURLEncoding.EncodeToString(data))
	if err != nil || b != p.Branding() {
		t.Fatal(err, b)
	}
	for _, input := range []string{"!", "e30", base64.RawURLEncoding.EncodeToString([]byte(`{"schemaVersion":1,"cliName":"../bad"}`)), strings.Repeat("a", MaxProfileBytes*2)} {
		if _, err := resolveEmbedded(input); err == nil {
			t.Error("accepted malformed embedded state")
		}
	}
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			copy := Default()
			copy.ProductName = "changed"
			if Default() != defaults() {
				t.Error("mutable profile leaked")
			}
		})
	}
	wg.Wait()
}
