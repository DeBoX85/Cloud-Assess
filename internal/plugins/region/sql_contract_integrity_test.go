package region

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

// Frozen independent future literals, not execution of the absent SQL APIs.
func TestSQLCapabilityContractIntegrity(t *testing.T) {
	for _, f := range []struct {
		name, digest, format   string
		size, count, collector int
	}{
		{"sql-database-capability-contract-v1.json", "67b664515b0407d57a037d401bd830662db49b4b177c2a67140701a34587a5fd", "sql-database-capability-contract-v1", 362752, 153, 124},
		{"sql-managed-instance-capability-contract-v1.json", "4035965123c5c72c93c099c7accd6c7be8d2d5a84a7041a8359f0ec1ecd720a7", "sql-managed-instance-capability-contract-v1", 399973, 163, 134},
	} {
		b := vmRead(t, f.name)
		h := sha256.Sum256(b)
		if len(b) != f.size || hex.EncodeToString(h[:]) != f.digest {
			t.Fatal("frozen SQL oracle bytes changed", f.name)
		}
		var m struct {
			Format string `json:"format"`
			Cases  []struct {
				Name  string `json:"name"`
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
		names := map[string]bool{}
		collector, direct := 0, 0
		for _, c := range m.Cases {
			if c.Name == "" || names[c.Name] {
				t.Fatal("scenario identity", f.name)
			}
			names[c.Name] = true
			switch c.Input.Source.Kind {
			case "collector":
				collector++
			case "collection":
				direct++
			default:
				t.Fatal("source kind")
			}
		}
		if m.Format != f.format || len(m.Cases) != f.count || collector != f.collector || direct != 29 {
			t.Fatal("scenario counts", f.name)
		}
	}
}

func TestSQLCapabilityPrimaryIntegrity(t *testing.T) {
	for _, f := range []struct {
		name, digest, blob string
		size               int
	}{
		{"sql-primary-2021-11-01-schema.json", "8246593e56c6d2708186a3603f9d92ee64ef2cc292931c7db09753471e506040", "d6c7b7fb88fa594932ba62c5523c343ab85108dd", 38062},
		{"sql-primary-2021-11-01-example.json", "f4b992d87b9e6e8bae6a7693092b4028dfa60bc7a79f132a8b658acebc9b57dd", "9815b0a083dd402d2b7b17e29bb9b09df392dfce", 119773},
		{"sql-primary-2021-11-01-common.json", "709781a78cdfb4ac27edcd4bdc7709e4667b3e33eaf4f91ae1d2c21a42b9a1c3", "f7decac418742c0cd8fb0c6b44d1330a5b6cb46d", 5904},
		{"storage-primary-LICENSE.txt", "b37363bf2295aa0acde60ce78a60034a71a0955075fc01aa09cc23eadff44436", "6acfdd593962ad75691ac6cae598a91a405f348b", 1076},
	} {
		b := vmRead(t, f.name)
		h := sha256.Sum256(b)
		blob := sha1.Sum(append([]byte(fmt.Sprintf("blob %d\x00", len(b))), b...))
		if len(b) != f.size || hex.EncodeToString(h[:]) != f.digest || hex.EncodeToString(blob[:]) != f.blob {
			t.Fatal("immutable SQL primary changed", f.name)
		}
	}
}
