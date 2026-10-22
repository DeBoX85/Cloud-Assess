package region

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

// Integrity only: future collector expectations are not runtime coverage.
func TestRetailPricingFrozenIntegrity(t *testing.T) {
	for _, v := range []struct {
		path string
		size int
		hash string
	}{
		{"testdata/retail-pricing-contract-v1.json", 92775, "73b9005094216e579156459a77a4ee9d0cfa18c46182d0adcfb7fbafe0b36d3a"},
		{"testdata/retail-primary-cost.go.txt", 32982, "c2a80218cf061e94b49f8491f7970a8b668272f06fd60e7d754969f21442120b"},
		{"testdata/retail-primary-types.go.txt", 11560, "153dd427c31523e13f29d601c042415cceaad27643ccf95e538ac28db3fdfc3b"},
		{"testdata/retail-primary-http.go.txt", 7848, "acf65aff604df4c96c7b726423053eb73520e032e78c9aa72ab9c8500016e688"},
	} {
		b, e := os.ReadFile(v.path)
		if e != nil || len(b) != v.size || fmt.Sprintf("%x", sha256.Sum256(b)) != v.hash {
			t.Fatalf("frozen retail input changed: %s", v.path)
		}
	}
}
