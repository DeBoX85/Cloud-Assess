package arg

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

type propertyRow struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// This table characterizes JSON type tolerance separately from row-level
// composition. Null/empty objects remain accepted zero rows for source parity;
// acceptance is not proof that an Azure business record is semantically valid.
func TestDecodeRowsShapeAccounting(t *testing.T) {
	for _, tc := range []struct {
		name, raw string
		want      []propertyRow
		rejected  int
	}{
		{"object", `{"name":"ok","count":2}`, []propertyRow{{"ok", 2}}, 0},
		{"unknown-field", `{"name":"ok","count":2,"extra":true}`, []propertyRow{{"ok", 2}}, 0},
		{"null", `null`, []propertyRow{{}}, 0},
		{"empty-object", `{}`, []propertyRow{{}}, 0},
		{"null-fields", `{"name":null,"count":null}`, []propertyRow{{}}, 0},
		{"wrong-string-type", `{"name":2}`, []propertyRow{}, 1},
		{"wrong-number-type", `{"count":"2"}`, []propertyRow{}, 1},
		{"fractional-number", `{"count":2.5}`, []propertyRow{}, 1},
		{"array", `[]`, []propertyRow{}, 1},
		{"scalar", `true`, []propertyRow{}, 1},
		{"truncated", `{"name":`, []propertyRow{}, 1},
		{"trailing-document", `{} {}`, []propertyRow{}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, rejected := DecodeRowsWithStats[propertyRow]([]json.RawMessage{json.RawMessage(tc.raw)})
			if !reflect.DeepEqual(got, tc.want) || rejected != tc.rejected {
				t.Fatalf("rows=%+v rejected=%d want %+v/%d", got, rejected, tc.want, tc.rejected)
			}
		})
	}
}

func FuzzDecodeRowsIsolationAndAccounting(f *testing.F) {
	for _, seed := range []string{`{"name":"middle","count":3}`, `{"name":2}`, `null`, `[]`, `{}`, `{"name":`, "", `{} {}`, `{"name":"\u0000"}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 8192 {
			t.Skip()
		}
		original := bytes.Clone(raw)
		input := []json.RawMessage{json.RawMessage(`{"name":"before","count":1}`), raw, json.RawMessage(`{"name":"after","count":2}`)}
		rows, rejected := DecodeRowsWithStats[propertyRow](input)
		if len(rows)+rejected != len(input) || rejected < 0 || rejected > 1 {
			t.Fatalf("row accounting: rows=%d rejected=%d", len(rows), rejected)
		}
		if rows[0] != (propertyRow{"before", 1}) || rows[len(rows)-1] != (propertyRow{"after", 2}) {
			t.Fatalf("one row contaminated healthy neighbors: %+v", rows)
		}
		prefix, prefixRejected := DecodeRowsWithStats[propertyRow](input[:2])
		suffix, suffixRejected := DecodeRowsWithStats[propertyRow](input[2:])
		combined := append(prefix, suffix...)
		if !reflect.DeepEqual(rows, combined) || rejected != prefixRejected+suffixRejected {
			t.Fatalf("decoding depends on partition: whole=%+v split=%+v", rows, combined)
		}
		if !bytes.Equal(raw, original) {
			t.Fatal("decoder mutated caller-owned row bytes")
		}
	})
}
