package region

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

// Freeze the independent pre-edit oracle; this does not execute the future API.
func TestVMSSCatalogContractIntegrity(t *testing.T) {
	data, err := os.ReadFile("testdata/vmss-catalog-target-contract-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	got := sha256.Sum256(data)
	if hex.EncodeToString(got[:]) != "e09f3814ebc647265ea2aa7e8614e855ee24a96e4eefd4fa4897d7d739e881c0" {
		t.Fatal("frozen VMSS catalog oracle bytes changed")
	}
}
