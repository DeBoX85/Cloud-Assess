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

// retailJSON bounds every value, key and delimiter before materialisation.
// Work is cumulative across both phases; keys use the pinned Unicode fold.
func retailJSON(ctx context.Context, body []byte, work *int) string {
	if !utf8.Valid(body) {
		return "retail_decode_invalid"
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var walk func(int) string
	text := func(s string) string {
		if len(s) > 8192 {
			return "retail_limit"
		}
		if strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || regionNoncharacter(r) }) >= 0 {
			return "retail_decode_invalid"
		}
		return ""
	}
	walk = func(depth int) string {
		if ctx.Err() != nil {
			return "retail_decode_invalid"
		}
		if depth > 32 || *work >= 1048576 {
			return "retail_limit"
		}
		*work = *work + 1
		token, e := d.Token()
		if e != nil {
			return "retail_decode_invalid"
		}
		switch v := token.(type) {
		case string:
			return text(v)
		case json.Number:
			n, e := v.Float64()
			if e != nil || math.IsNaN(n) || math.IsInf(n, 0) || n == 0 && retailNonzero(string(v)) {
				return "retail_decode_invalid"
			}
			return ""
		}
		delim, compound := token.(json.Delim)
		if !compound {
			return ""
		}
		if delim != '{' && delim != '[' {
			return "retail_decode_invalid"
		}
		keys := map[string]bool{}
		for d.More() {
			if delim == '{' {
				if len(keys) >= 128 {
					return "retail_limit"
				}
				if *work >= 1048576 {
					return "retail_limit"
				}
				*work++
				key, e := d.Token()
				s, ok := key.(string)
				if e != nil || !ok {
					return "retail_decode_invalid"
				}
				if code := text(s); code != "" {
					return code
				}
				fold := regionFoldKey(s)
				if keys[fold] {
					return "retail_decode_invalid"
				}
				keys[fold] = true
			}
			if code := walk(depth + 1); code != "" {
				return code
			}
		}
		if *work >= 1048576 {
			return "retail_limit"
		}
		*work++
		end, e := d.Token()
		if e != nil || delim == '{' && end != json.Delim('}') || delim == '[' && end != json.Delim(']') {
			return "retail_decode_invalid"
		}
		return ""
	}
	if code := walk(0); code != "" {
		return code
	}
	if _, e := d.Token(); e != io.EOF {
		return "retail_decode_invalid"
	}
	return ""
}

func retailNonzero(raw string) bool {
	if i := strings.IndexAny(raw, "eE"); i >= 0 {
		raw = raw[:i]
	}
	return strings.IndexAny(raw, "123456789") >= 0
}

func retailCell(s string) bool {
	return len(s) <= 512 && utf8.ValidString(s) && strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || regionNoncharacter(r) }) < 0
}

func decodeRetailPage(ctx context.Context, body []byte, work, rawCount *int, metadata bool) ([]RetailPriceItem, string, string) {
	bad := func() ([]RetailPriceItem, string, string) { return nil, "", "retail_decode_invalid" }
	if code := retailJSON(ctx, body, work); code != "" {
		return nil, "", code
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var envelope map[string]any
	if d.Decode(&envelope) != nil || envelope == nil {
		return bad()
	}
	raw, _ := historyField(envelope, "Items")
	rows, ok := raw.([]any)
	if !ok || rows == nil {
		return bad()
	}
	if len(rows) > 65536-*rawCount {
		return nil, "", "retail_limit"
	}
	*rawCount += len(rows)
	next := ""
	if n, exists := historyField(envelope, "NextPageLink"); exists && n != nil {
		next, ok = n.(string)
		if !ok {
			return bad()
		}
	}
	// Envelope display fields are optional; they must never coerce null/numbers.
	for _, name := range []string{"BillingCurrency", "CustomerEntityId", "CustomerEntityType"} {
		if v, exists := historyField(envelope, name); exists {
			s, ok := v.(string)
			if !ok || !retailCell(s) {
				return bad()
			}
		}
	}
	items := make([]RetailPriceItem, 0, len(rows))
	for _, row := range rows {
		if ctx.Err() != nil {
			return bad()
		}
		m, ok := row.(map[string]any)
		if !ok || m == nil {
			return bad()
		}
		var item RetailPriceItem
		fields := []struct {
			name     string
			dest     *string
			required bool
		}{
			{"currencyCode", &item.CurrencyCode, true}, {"armRegionName", &item.ArmRegionName, true},
			{"meterId", &item.MeterID, true}, {"meterName", &item.MeterName, true}, {"productId", &item.ProductID, true},
			{"skuName", &item.SkuName, true}, {"unitOfMeasure", &item.UnitOfMeasure, true}, {"type", &item.Type, true},
			{"location", &item.Location, false}, {"effectiveStartDate", &item.EffectiveStartDate, false},
			{"skuId", &item.SkuID, false}, {"productName", &item.ProductName, false}, {"serviceName", &item.ServiceName, false},
			{"serviceId", &item.ServiceID, false}, {"serviceFamily", &item.ServiceFamily, false}, {"armSkuName", &item.ArmSkuName, false},
		}
		for _, f := range fields {
			v, present := historyField(m, f.name)
			if !present {
				if f.required {
					return bad()
				}
				continue
			}
			s, ok := v.(string)
			if !ok || !retailCell(s) || f.required && s == "" {
				return bad()
			}
			*f.dest = s
		}
		for _, f := range []struct {
			name string
			dest *float64
		}{{"tierMinimumUnits", &item.TierMinimumUnits}, {"retailPrice", &item.RetailPrice}, {"unitPrice", &item.UnitPrice}} {
			v, _ := historyField(m, f.name)
			n, ok := v.(json.Number)
			if !ok {
				return bad()
			}
			x, e := n.Float64()
			if e != nil || x < 0 || x > 1e12 {
				return bad()
			}
			*f.dest = x
		}
		v, _ := historyField(m, "isPrimaryMeterRegion")
		item.IsPrimaryMeterRegion, ok = v.(bool)
		if !ok || metadata && !item.IsPrimaryMeterRegion || item.CurrencyCode != "USD" || item.Type != "Consumption" || item.RetailPrice != item.UnitPrice {
			return bad()
		}
		item.ArmRegionName = strings.ReplaceAll(strings.ToLower(item.ArmRegionName), " ", "")
		if !regionID.MatchString(item.ArmRegionName) || !costPhysical(item.ArmRegionName) {
			return bad()
		}
		if item.EffectiveStartDate != "" {
			if _, valid := retailTimestamp(item.EffectiveStartDate); !valid {
				return bad()
			}
		}
		items = append(items, item)
	}
	return items, next, ""
}
