package region

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

// This freezes independent future cases. It deliberately does not execute the
// absent CognitiveSKUCollector/Catalog APIs or certify model deployment support.
func TestCognitiveSKUContractIntegrity(t *testing.T) {
	b := vmRead(t, "cognitive-sku-target-contract-v1.json")
	h := sha256.Sum256(b)
	if len(b) != 268084 || hex.EncodeToString(h[:]) != "c12604b2a6ac8b7fed6eb7c66397cf181bba8c0505b512bd5e844972f6aaa81b" {
		t.Fatal("frozen Cognitive future oracle bytes changed")
	}
	var m struct {
		Format string `json:"format"`
		Cases  []struct {
			Input struct {
				Source struct {
					Kind string `json:"kind"`
				} `json:"source"`
			} `json:"input"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	collector, direct := 0, 0
	for _, c := range m.Cases {
		if c.Input.Source.Kind == "collector" {
			collector++
		} else if c.Input.Source.Kind == "collection" {
			direct++
		} else {
			t.Fatal("unknown future source kind")
		}
	}
	if m.Format != "cognitive-sku-target-contract-v1" || len(m.Cases) != 143 || collector != 108 || direct != 35 {
		t.Fatal("future scenario identity/count changed")
	}
}

func TestCognitiveSKUPrimaryIntegrity(t *testing.T) {
	for _, f := range []struct {
		name, sha256, blob string
		bytes              int
	}{
		{"cognitive-primary-2024-10-01-schema.json", "e6efa93f54d62b025c4afe953b7491ea9c50f497646a6b2545f87e82d03c306a", "f0919ebcbf91755f10591b931bd80138df13f547", 239116},
		{"cognitive-primary-2024-10-01-GetSkus.json", "1c859ac36d080942f817958db9085603cc00b3ce42bbbab5d7a236cf53fc33d6", "57c0132a67d4958fb40cf20914537b034a10f39f", 55486},
	} {
		b := vmRead(t, f.name)
		h := sha256.Sum256(b)
		blob := sha1.Sum(append([]byte(fmt.Sprintf("blob %d\x00", len(b))), b...))
		if len(b) != f.bytes || hex.EncodeToString(h[:]) != f.sha256 || hex.EncodeToString(blob[:]) != f.blob {
			t.Fatal("immutable primary identity changed", f.name)
		}
	}
}
