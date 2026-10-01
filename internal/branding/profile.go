package branding

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MaxProfileBytes bounds both operator input and decoded linker state.
const MaxProfileBytes = 16384

// Profile contains only implemented presentation fields. It is public metadata.
type Profile struct {
	SchemaVersion    int    `json:"schemaVersion"`
	ProductName      string `json:"productName"`
	CLIName          string `json:"cliName"`
	ReportTitle      string `json:"reportTitle"`
	ReportFilePrefix string `json:"reportFilePrefix"`
	WebsiteURL       string `json:"websiteURL"`
}

// embeddedProfile is set only by the supported builder using the Go linker -X.
var embeddedProfile string

var embeddedBrand, embeddedErr = resolveEmbedded(embeddedProfile)

// EmbeddedError rejects malformed linked state before any CLI work.
func EmbeddedError() error { return embeddedErr }

func resolveEmbedded(encoded string) (Branding, error) {
	if encoded == "" {
		return defaults(), nil
	}
	if len(encoded) > base64.RawURLEncoding.EncodedLen(MaxProfileBytes) {
		return Branding{}, fmt.Errorf("invalid embedded branding: size limit")
	}
	data, err := base64.RawURLEncoding.Strict().DecodeString(encoded)
	if err != nil {
		return Branding{}, fmt.Errorf("invalid embedded branding: encoding")
	}
	profile, err := ParseProfile(data)
	if err != nil {
		return Branding{}, fmt.Errorf("invalid embedded branding: %w", err)
	}
	return profile.Branding(), nil
}

// ParseProfile rejects duplicate/unknown keys rather than accepting JSON's last
// value or case-insensitive struct-field matching. Omitted fields use defaults.
func ParseProfile(data []byte) (Profile, error) {
	if len(data) > MaxProfileBytes || !utf8.Valid(data) {
		return Profile{}, fmt.Errorf("branding profile exceeds size limit or is not UTF-8")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return Profile{}, fmt.Errorf("branding profile must be an object")
	}
	brand := defaults()
	p := Profile{1, brand.ProductName, brand.CLIName, brand.ReportTitle, brand.ReportFilePrefix, brand.WebsiteURL}
	seen := map[string]bool{}
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return Profile{}, fmt.Errorf("invalid branding profile key")
		}
		key, ok := keyToken.(string)
		if !ok || seen[key] {
			return Profile{}, fmt.Errorf("duplicate or invalid branding profile key")
		}
		seen[key] = true
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return Profile{}, fmt.Errorf("invalid branding profile value")
		}
		if key == "schemaVersion" {
			if err := json.Unmarshal(raw, &p.SchemaVersion); bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || err != nil || p.SchemaVersion != 1 {
				return Profile{}, fmt.Errorf("branding schemaVersion must be 1")
			}
			continue
		}
		var field *string
		switch key {
		case "productName":
			field = &p.ProductName
		case "cliName":
			field = &p.CLIName
		case "reportTitle":
			field = &p.ReportTitle
		case "reportFilePrefix":
			field = &p.ReportFilePrefix
		case "websiteURL":
			field = &p.WebsiteURL
		default:
			return Profile{}, fmt.Errorf("unsupported branding profile key")
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, field) != nil {
			return Profile{}, fmt.Errorf("branding fields must be strings")
		}
	}
	if token, err = decoder.Token(); err != nil || token != json.Delim('}') {
		return Profile{}, fmt.Errorf("unterminated branding profile")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return Profile{}, fmt.Errorf("trailing branding profile data")
	}
	if !seen["schemaVersion"] {
		return Profile{}, fmt.Errorf("branding schemaVersion is required")
	}
	if err := p.Validate(); err != nil {
		return Profile{}, err
	}
	return p, nil
}

var slug = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
var device = regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])$`)

// SafeComponent restricts presentation-derived filenames on Linux and Windows.
func SafeComponent(value string) bool { return slug.MatchString(value) && !device.MatchString(value) }

func safeText(value string, limit int) bool {
	if value == "" || len(value) > limit || strings.TrimSpace(value) != value || !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == utf8.RuneError {
			return false
		}
	}
	return true
}

func (p Profile) Validate() error {
	if p.SchemaVersion != 1 {
		return fmt.Errorf("branding schemaVersion must be 1")
	}
	if !safeText(p.ProductName, 256) || !safeText(p.ReportTitle, 256) {
		return fmt.Errorf("invalid branding display text")
	}
	if !SafeComponent(p.CLIName) || !SafeComponent(p.ReportFilePrefix) {
		return fmt.Errorf("invalid branding filename component")
	}
	u, err := url.Parse(p.WebsiteURL)
	if !safeText(p.WebsiteURL, 2048) || err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.Opaque != "" {
		return fmt.Errorf("branding websiteURL must be an absolute HTTPS URL without user information")
	}
	return nil
}

// Canonical returns a deterministic, fully resolved profile after validation.
func (p Profile) Canonical() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

func (p Profile) Branding() Branding {
	brand := defaults()
	brand.ProductName, brand.CLIName, brand.ReportTitle = p.ProductName, p.CLIName, p.ReportTitle
	brand.ReportFilePrefix, brand.WebsiteURL = p.ReportFilePrefix, p.WebsiteURL
	return brand
}

// CurrentProfile exports only the five active fields, never placeholders.
func CurrentProfile() Profile {
	b := Default()
	return Profile{1, b.ProductName, b.CLIName, b.ReportTitle, b.ReportFilePrefix, b.WebsiteURL}
}
