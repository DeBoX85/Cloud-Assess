package skus

import "testing"

func TestComputeCapacity(t *testing.T) {
	known := "Standard_D2s_v5"
	if sku, ok := Lookup(known); !ok || sku.VCPUs <= 0 {
		t.Fatalf("expected pinned SKU %s with vCPU data", known)
	}

	tests := []struct {
		name       string
		sku        string
		capacity   int
		vmss       bool
		want       string
		useKnownVC bool
	}{
		{name: "explicit non vmss capacity", sku: known, capacity: 3, want: "3"},
		{name: "unknown sku explicit capacity", sku: "does-not-exist", capacity: 4, vmss: true, want: "4"},
		{name: "unknown sku no capacity", sku: "does-not-exist", want: ""},
		{name: "known sku no capacity uses vcpu", sku: known, useKnownVC: true},
		{name: "vmss multiplies instances by vcpu", sku: known, capacity: 3, vmss: true, useKnownVC: true},
	}

	knownSKU, _ := Lookup(known)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			if tt.useKnownVC {
				if tt.vmss && tt.capacity > 0 {
					want = itoa(tt.capacity * knownSKU.VCPUs)
				} else {
					want = itoa(knownSKU.VCPUs)
				}
			}
			if got := ComputeCapacity(tt.sku, tt.capacity, tt.vmss); got != want {
				t.Fatalf("ComputeCapacity() = %q, want %q", got, want)
			}
		})
	}
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	var digits [20]byte
	i := len(digits)
	for value > 0 {
		i--
		digits[i] = byte('0' + value%10)
		value /= 10
	}
	if negative {
		i--
		digits[i] = '-'
	}
	return string(digits[i:])
}
