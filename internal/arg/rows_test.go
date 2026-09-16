package arg

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDecodeRowsSkipsMalformedEntries(t *testing.T) {
	type row struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	input := []json.RawMessage{
		json.RawMessage(`{"name":"a","count":1}`),
		json.RawMessage(`{not json}`),
		json.RawMessage(`{"name":"bad","count":"not-an-int"}`),
		json.RawMessage(`{"name":"c","count":3}`),
	}
	want := []row{{Name: "a", Count: 1}, {Name: "c", Count: 3}}
	if got := DecodeRows[row](input); !reflect.DeepEqual(got, want) {
		t.Fatalf("DecodeRows() = %#v, want %#v", got, want)
	}
}

func TestDecodeRowsNilAndEmptyReturnEmptySlice(t *testing.T) {
	type row struct{ Name string }
	if got := DecodeRows[row](nil); len(got) != 0 {
		t.Fatalf("nil input returned %d rows", len(got))
	}
	if got := DecodeRows[row]([]json.RawMessage{}); len(got) != 0 {
		t.Fatalf("empty input returned %d rows", len(got))
	}
}

func TestRawMessageString(t *testing.T) {
	tests := []struct {
		name string
		in   json.RawMessage
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "empty", in: json.RawMessage(``), want: ""},
		{name: "null", in: json.RawMessage(`null`), want: ""},
		{name: "string", in: json.RawMessage(`"hello"`), want: "hello"},
		{name: "escaped string", in: json.RawMessage(`"a \"b\" c"`), want: `a "b" c`},
		{name: "number", in: json.RawMessage(`42`), want: "42"},
		{name: "bool", in: json.RawMessage(`true`), want: "true"},
		{name: "object", in: json.RawMessage(`{"k":"v"}`), want: `{"k":"v"}`},
		{name: "array", in: json.RawMessage(`["a","b"]`), want: `["a","b"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RawMessageString(tt.in); got != tt.want {
				t.Fatalf("RawMessageString(%q) = %q, want %q", string(tt.in), got, tt.want)
			}
		})
	}
}
