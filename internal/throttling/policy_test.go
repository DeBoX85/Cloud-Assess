package throttling

import "testing"

func TestKindForURL(t *testing.T) {
	tests := []struct {
		url  string
		want Kind
	}{
		{url: "https://management.azure.com/providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01", want: KindGraph},
		{url: "https://management.azure.com/subscriptions/sub/providers/Microsoft.CostManagement/query", want: KindCost},
		{url: "https://prices.azure.com/api/retail/prices", want: KindRetail},
		{url: "https://management.azure.com/subscriptions", want: KindARM},
	}
	for _, tt := range tests {
		if got := KindForURL(tt.url); got != tt.want {
			t.Fatalf("KindForURL(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}
