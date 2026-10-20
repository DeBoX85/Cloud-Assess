package region

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

// These are independent future expectations. This integrity check does not
// exercise the absent query builder or collector and is not runtime coverage.
func TestMeterHistoryFrozenOracleIntegrity(t *testing.T) {
	b, err := os.ReadFile("testdata/meter-history-contract-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 70949 || fmt.Sprintf("%x", sha256.Sum256(b)) != "fc87a540dc0893ad44b1b54f1397973aa0d0de9db212f806f691e17fe521ebc8" {
		t.Fatal("frozen independent meter-history oracle changed")
	}
}
