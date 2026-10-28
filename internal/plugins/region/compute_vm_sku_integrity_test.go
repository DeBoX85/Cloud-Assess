package region

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestComputeVMSKUManifestIdentity(t *testing.T) {
	b := vmRead(t, "compute-vm-target-contract-v1.json")
	h := sha256.Sum256(b)
	if hex.EncodeToString(h[:]) != "ac5eb22cb4582959afe6f9ccaf43a9aaba908231019a7b85b4e257c599292c34" {
		t.Fatal("pre-edit VM scenario oracle bytes changed")
	}
}
