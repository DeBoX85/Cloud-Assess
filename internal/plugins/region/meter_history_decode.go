package region

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

// historyJSON validates all metadata before generic materialisation. Work is
// cumulative across pages; keys use the same Unicode fold as the pinned decoder.
func historyJSON(ctx context.Context, body []byte, work *int) string {
	if !utf8.Valid(body) {
		return "history_decode_invalid"
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var walk func(int) string
	text := func(s string) string {
		if len(s) > 8192 {
			return "history_limit"
		}
		if strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || regionNoncharacter(r) }) >= 0 {
			return "history_decode_invalid"
		}
		return ""
	}
	walk = func(depth int) string {
		if ctx.Err() != nil {
			return "history_decode_invalid"
		}
		if depth > 32 || *work >= 65536 {
			return "history_limit"
		}
		*work = *work + 1
		token, e := d.Token()
		if e != nil {
			return "history_decode_invalid"
		}
		switch v := token.(type) {
		case string:
			return text(v)
		case json.Number:
			n, e := v.Float64()
			if e != nil || math.IsNaN(n) || math.IsInf(n, 0) {
				return "history_decode_invalid"
			}
			return ""
		}
		delim, compound := token.(json.Delim)
		if !compound {
			return ""
		}
		if delim != '{' && delim != '[' {
			return "history_decode_invalid"
		}
		keys := map[string]bool{}
		for d.More() {
			if delim == '{' {
				if len(keys) >= 128 {
					return "history_limit"
				}
				key, e := d.Token()
				s, ok := key.(string)
				if e != nil || !ok {
					return "history_decode_invalid"
				}
				if code := text(s); code != "" {
					return code
				}
				fold := regionFoldKey(s)
				if keys[fold] {
					return "history_decode_invalid"
				}
				keys[fold] = true
			}
			if code := walk(depth + 1); code != "" {
				return code
			}
		}
		end, e := d.Token()
		if e != nil || delim == '{' && end != json.Delim('}') || delim == '[' && end != json.Delim(']') {
			return "history_decode_invalid"
		}
		return ""
	}
	if code := walk(0); code != "" {
		return code
	}
	if _, e := d.Token(); e != io.EOF {
		return "history_decode_invalid"
	}
	return ""
}

func historyField(m map[string]any, name string) (any, bool) {
	for key, value := range m {
		if regionFoldKey(key) == regionFoldKey(name) {
			return value, true
		}
	}
	return nil, false
}

func historyOwnedResource(id, subscription string) bool {
	if len(id) > 4096 || strings.ContainsAny(id, "?#\\") {
		return false
	}
	parts := strings.Split(id, "/")
	if len(parts) < 4 || parts[0] != "" || strings.ToLower(parts[1]) != "subscriptions" || len(parts[2]) != 36 || !subscriptionID.MatchString(parts[2]) || !strings.EqualFold(parts[2], subscription) {
		return false
	}
	for _, p := range parts[3:] {
		if p == "" {
			return false
		}
	}
	return true
}

func historyRawPositive(raw string) bool {
	if strings.HasPrefix(raw, "-") {
		return false
	}
	if i := strings.IndexAny(raw, "eE"); i >= 0 {
		raw = raw[:i]
	}
	return strings.IndexAny(raw, "123456789") >= 0
}

func decodeHistoryPage(ctx context.Context, body []byte, subscription, currency string, work, rowCount *int) ([]HistoricalCostMeter, string, string) {
	bad := func() ([]HistoricalCostMeter, string, string) { return nil, "", "history_decode_invalid" }
	if code := historyJSON(ctx, body, work); code != "" {
		return nil, "", code
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var envelope map[string]any
	if d.Decode(&envelope) != nil || envelope == nil {
		return bad()
	}
	p, _ := historyField(envelope, "properties")
	properties, ok := p.(map[string]any)
	if !ok || properties == nil {
		return bad()
	}
	c, _ := historyField(properties, "columns")
	columns, ok := c.([]any)
	if !ok || columns == nil {
		return bad()
	}
	r, _ := historyField(properties, "rows")
	rows, ok := r.([]any)
	if !ok || rows == nil {
		return bad()
	}
	if len(columns) > 32 || len(rows) > 8192-*rowCount {
		return nil, "", "history_limit"
	}
	next := ""
	if n, present := historyField(properties, "nextLink"); present && n != nil {
		next, ok = n.(string)
		if !ok {
			return bad()
		}
	}
	types := make([]string, len(columns))
	indices := map[string]int{}
	seen := map[string]bool{}
	costIndex := -1
	for i, value := range columns {
		col, ok := value.(map[string]any)
		if !ok || col == nil {
			return bad()
		}
		n, _ := historyField(col, "name")
		name, ok := n.(string)
		if !ok || name == "" {
			return bad()
		}
		v, _ := historyField(col, "type")
		typ, ok := v.(string)
		if !ok || (typ != "String" && typ != "Number") || seen[regionFoldKey(name)] {
			return bad()
		}
		seen[regionFoldKey(name)] = true
		types[i] = typ
		indices[name] = i
		switch name {
		case "PreTaxCost", "Cost", "CostInBillingCurrency":
			if costIndex != -1 || typ != "Number" {
				return bad()
			}
			costIndex = i
		case "ResourceId", "ResourceGuid", "Currency":
			if typ != "String" {
				return bad()
			}
		}
	}
	if len(rows) > 0 {
		if _, ok := indices["ResourceId"]; !ok {
			return bad()
		}
		if _, ok := indices["ResourceGuid"]; !ok {
			return bad()
		}
		if _, ok := indices["Currency"]; !ok || costIndex == -1 {
			return bad()
		}
	}
	out := []HistoricalCostMeter{}
	for _, value := range rows {
		if ctx.Err() != nil {
			return bad()
		}
		row, ok := value.([]any)
		if !ok || row == nil || len(row) != len(columns) {
			return bad()
		}
		for i, cell := range row {
			if types[i] == "String" {
				if _, ok := cell.(string); !ok {
					return bad()
				}
			} else if _, ok := cell.(json.Number); !ok {
				return bad()
			}
		}
		resource := row[indices["ResourceId"]].(string)
		meter := row[indices["ResourceGuid"]].(string)
		if !historyOwnedResource(resource, subscription) || meter == "" || len(meter) > 512 || row[indices["Currency"]].(string) != currency {
			return bad()
		}
		number := row[costIndex].(json.Number)
		cost, e := number.Float64()
		if e != nil || math.IsNaN(cost) || math.IsInf(cost, 0) || math.Abs(cost) > MaxCostValue || cost == 0 && historyRawPositive(number.String()) {
			return bad()
		}
		if cost <= 0 {
			continue
		}
		if cost < MinCostWeight {
			return bad()
		}
		out = append(out, HistoricalCostMeter{MeterID: meter, HistoricalCost: cost})
	}
	*rowCount += len(rows)
	return out, next, ""
}
