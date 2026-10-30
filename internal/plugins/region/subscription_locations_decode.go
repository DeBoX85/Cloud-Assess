// Location wire admission follows the independently frozen contract. It does
// not change the older quota/provider JSON or opaque-map admission policies.
package region

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type locationRow struct {
	name     string
	mappings map[string]string
}
type locationPage struct {
	rows       []locationRow
	names      []string
	next       string
	work, text int
}

const (
	locationEnvelope = iota
	locationArray
	locationObject
	locationMetadata
	locationMappingArray
	locationMappingObject
	locationOpaque
)

func locationText(s string, limit int, nonempty bool) bool {
	return (!nonempty || s != "") && len(s) <= limit && utf8.ValidString(s) && strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || regionNoncharacter(r) }) < 0
}

// Exact decoded duplicates apply everywhere. Typed aliases are checked only by
// canonical object readers below; unknown opaque maps remain case-sensitive.
func locationJSON(ctx context.Context, body []byte) bool {
	if !utf8.Valid(body) {
		return false
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	nodes := 0
	var walk func(int, int, int) bool
	walk = func(depth, kind, limit int) bool {
		if ctx.Err() != nil || depth > 32 || nodes >= 65536 {
			return false
		}
		nodes++
		token, e := d.Token()
		if e != nil {
			return false
		}
		if s, ok := token.(string); ok {
			return locationText(s, limit, false)
		}
		delim, compound := token.(json.Delim)
		if !compound {
			return true
		}
		if delim != '{' && delim != '[' {
			return false
		}
		keys := map[string]bool{}
		for d.More() {
			child, nextLimit := locationOpaque, MaxLabelBytes
			if delim == '{' {
				token, e := d.Token()
				key, ok := token.(string)
				if e != nil || !ok || !locationText(key, MaxLabelBytes, false) || len(keys) >= 128 || keys[key] {
					return false
				}
				keys[key] = true
				switch {
				case kind == locationEnvelope && key == "value":
					child = locationArray
				case kind == locationObject && key == "metadata":
					child = locationMetadata
				case kind == locationObject && key == "availabilityZoneMappings":
					child = locationMappingArray
				}
				if kind == locationEnvelope && key == "nextLink" {
					nextLimit = 8192
				}
				if kind == locationObject && key == "id" {
					nextLimit = 2048
				}
			} else {
				if kind == locationArray {
					child = locationObject
				}
				if kind == locationMappingArray {
					child = locationMappingObject
				}
			}
			if !walk(depth+1, child, nextLimit) {
				return false
			}
		}
		end, e := d.Token()
		return e == nil && (delim == '{' && end == json.Delim('}') || delim == '[' && end == json.Delim(']'))
	}
	if !walk(0, locationEnvelope, MaxLabelBytes) {
		return false
	}
	_, e := d.Token()
	return e == io.EOF
}
func locationString(raw json.RawMessage, limit int, nonempty bool) (string, bool) {
	var p *string
	if len(raw) == 0 || json.Unmarshal(raw, &p) != nil || p == nil || !locationText(*p, limit, nonempty) {
		return "", false
	}
	return *p, true
}
func locationReturnedID(id, subscription, name string) bool {
	if !storageASCII(id) {
		return false
	}
	parts := strings.Split(id, "/")
	return len(parts) == 5 && parts[0] == "" && strings.EqualFold(parts[1], "subscriptions") && len(parts[2]) == 36 && subscriptionID.MatchString(parts[2]) && strings.EqualFold(parts[2], subscription) && strings.EqualFold(parts[3], "locations") && strings.EqualFold(parts[4], name)
}
func decodeLocationPage(ctx context.Context, body []byte, subscription string, prior map[string]bool, work, text int) (locationPage, string) {
	invalid := func() (locationPage, string) { return locationPage{}, "locations_invalid_response" }
	if !locationJSON(ctx, body) {
		return invalid()
	}
	envelope, ok := providerObject(body, "value", "nextLink")
	if !ok {
		return invalid()
	}
	entries, ok := providerArray(envelope["value"])
	if !ok {
		return invalid()
	}
	page := locationPage{rows: []locationRow{}, names: []string{}, work: work, text: text}
	if raw, present := envelope["nextLink"]; present {
		var valid bool
		page.next, valid = locationString(raw, 8192, false)
		if !valid {
			return invalid()
		}
	}
	charge := func() bool {
		if page.work >= MaxAuxRows {
			return false
		}
		page.work++
		return true
	}
	label := func(s string) bool {
		if len(s) > auxTextBudget-page.text {
			return false
		}
		page.text += len(s)
		return true
	}
	names := map[string]bool{}
	for _, raw := range entries {
		if ctx.Err() != nil {
			return invalid()
		}
		if !charge() {
			return locationPage{}, "locations_work_limit"
		}
		fields, valid := providerObject(raw, "name", "metadata", "availabilityZoneMappings", "id", "subscriptionId")
		if !valid {
			return invalid()
		}
		name, valid := locationString(fields["name"], MaxLabelBytes, true)
		if !valid || !storageASCII(name) {
			return invalid()
		}
		name = strings.ToLower(name)
		if !regionID.MatchString(name) || !label(name) {
			return invalid()
		}
		if names[name] || prior[name] {
			return locationPage{}, "locations_identity_duplicate"
		}
		names[name] = true
		metadata, valid := providerObject(fields["metadata"], "regionType")
		if !valid {
			return invalid()
		}
		kind, valid := locationString(metadata["regionType"], MaxLabelBytes, true)
		if !valid || !label(kind) {
			return invalid()
		}
		if kind == "Physical" && !costPhysical(name) {
			return invalid()
		}
		if raw, present := fields["id"]; present {
			value, valid := locationString(raw, 2048, true)
			if !valid || !locationReturnedID(value, subscription, name) || !label(value) {
				return invalid()
			}
		}
		if raw, present := fields["subscriptionId"]; present {
			value, valid := locationString(raw, MaxLabelBytes, true)
			if !valid || len(value) != 36 || !subscriptionID.MatchString(value) || !strings.EqualFold(value, subscription) || !label(value) {
				return invalid()
			}
		}
		mappings := map[string]string{}
		if raw, present := fields["availabilityZoneMappings"]; present && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			values, valid := providerArray(raw)
			if !valid {
				return invalid()
			}
			for _, raw := range values {
				if ctx.Err() != nil {
					return invalid()
				}
				if !charge() {
					return locationPage{}, "locations_work_limit"
				}
				fields, valid := providerObject(raw, "logicalZone", "physicalZone")
				if !valid {
					return invalid()
				}
				logical, valid := locationString(fields["logicalZone"], MaxLabelBytes, true)
				if !valid || !label(logical) {
					return invalid()
				}
				physical, valid := locationString(fields["physicalZone"], MaxLabelBytes, true)
				if !valid || !label(physical) {
					return invalid()
				}
				if _, exists := mappings[logical]; exists {
					return locationPage{}, "locations_identity_duplicate"
				}
				mappings[logical] = physical
			}
		}
		page.names = append(page.names, name)
		if kind == "Physical" {
			page.rows = append(page.rows, locationRow{name: name, mappings: mappings})
		}
	}
	return page, ""
}
