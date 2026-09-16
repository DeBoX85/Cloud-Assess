// Portions of this package reproduce Azure Resource Graph row handling from
// Microsoft Azure Quick Review (MIT licensed). See NOTICE.md.

package arg

import (
	"bytes"
	"encoding/json"
)

// DecodeRows decodes each Azure Resource Graph row independently.
// Malformed rows are skipped so one bad row does not discard valid results.
func DecodeRows[T any](data []json.RawMessage) []T {
	rows, _ := DecodeRowsWithStats[T](data)
	return rows
}

// DecodeRowsWithStats preserves the reference tolerant row-decoding behavior while
// making skipped-row counts available to the target's explicit warning/completeness model.
func DecodeRowsWithStats[T any](data []json.RawMessage) ([]T, int) {
	out := make([]T, 0, len(data))
	malformed := 0
	for _, raw := range data {
		var row T
		if err := json.Unmarshal(raw, &row); err != nil {
			malformed++
			continue
		}
		out = append(out, row)
	}
	return out, malformed
}

// RawMessageString converts an ARG-projected value to the reference textual representation.
// JSON strings are unquoted; objects, arrays, numbers and booleans retain JSON text; null is empty.
func RawMessageString(value json.RawMessage) string {
	if len(value) == 0 || bytes.Equal(value, []byte("null")) {
		return ""
	}
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		var decoded string
		if err := json.Unmarshal(value, &decoded); err == nil {
			return decoded
		}
	}
	return string(value)
}
