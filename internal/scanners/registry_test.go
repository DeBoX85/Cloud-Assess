package scanners

import (
	"reflect"
	"testing"
)

func TestRegistryMatchesReferenceCoverage(t *testing.T) {
	keys := Keys()
	if got, want := len(keys), 87; got != want {
		t.Fatalf("scanner key count = %d, want %d", got, want)
	}

	for _, required := range []string{"aks", "apim", "avd", "avs", "hpc", "kv", "redis", "resource", "sap", "sql", "st", "vm", "vnet"} {
		if len(ByKey(required)) == 0 {
			t.Fatalf("required scanner key %q missing", required)
		}
	}
}

func TestRegistryPreservesSpecializedMappings(t *testing.T) {
	tests := []struct {
		key  string
		want []string
	}{
		{key: "avd", want: []string{"Specialized.Workload/AVD"}},
		{key: "avs", want: []string{"Microsoft.AVS/privateClouds", "Specialized.Workload/AVS"}},
		{key: "hpc", want: []string{"Specialized.Workload/HPC"}},
		{key: "sap", want: []string{"Specialized.Workload/SAP"}},
		{key: "vnet", want: []string{"Microsoft.Network/virtualNetworks", "Microsoft.Network/virtualNetworks/subnets"}},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			services := ByKey(tt.key)
			if len(services) != 1 {
				t.Fatalf("ByKey(%q) returned %d services, want 1", tt.key, len(services))
			}
			if !reflect.DeepEqual(services[0].ResourceTypes, tt.want) {
				t.Fatalf("resource types = %#v, want %#v", services[0].ResourceTypes, tt.want)
			}
		})
	}
}

func TestRedisPreservesTwoScannerEntries(t *testing.T) {
	services := ByKey("redis")
	if got, want := len(services), 2; got != want {
		t.Fatalf("redis service count = %d, want %d", got, want)
	}
	if services[0].Name != "Redis Cache" || services[1].Name != "Redis Enterprise" {
		t.Fatalf("unexpected redis services: %#v", services)
	}
}

func TestKeysAreSorted(t *testing.T) {
	keys := Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Fatalf("keys are not sorted at %q > %q", keys[i-1], keys[i])
		}
	}
}

func TestByKeyReturnsCopy(t *testing.T) {
	first := ByKey("st")
	first[0].Name = "changed"
	second := ByKey("st")
	if second[0].Name != "Storage Account" {
		t.Fatal("ByKey exposed mutable registry slice")
	}
}

func TestSelectedKeysUsesFilterResourceTypesForNormalScan(t *testing.T) {
	available := Keys()
	selected := SelectedKeys(available, []string{"aks", "ca", "does-not-exist", "st"})
	want := []string{"aks", "ca", "st"}
	if !reflect.DeepEqual(selected, want) {
		t.Fatalf("SelectedKeys() = %#v, want %#v", selected, want)
	}
}

func TestSelectedKeysScannerSpecificCommandWinsOverFilter(t *testing.T) {
	selected := SelectedKeys([]string{"vm"}, []string{"aks", "st"})
	if !reflect.DeepEqual(selected, []string{"vm"}) {
		t.Fatalf("SelectedKeys() = %#v, want vm", selected)
	}
}

func TestSelectedKeysDefaultsToAllScanners(t *testing.T) {
	selected := SelectedKeys(nil, nil)
	if !reflect.DeepEqual(selected, Keys()) {
		t.Fatal("empty scanner selection should default to all scanners")
	}
}

func TestAllowedResourceTypesExpandsSelectedScannerKeys(t *testing.T) {
	got := AllowedResourceTypes(Keys(), []string{"aks", "vnet"})
	want := []string{
		"Microsoft.ContainerService/managedClusters",
		"Microsoft.Network/virtualNetworks",
		"Microsoft.Network/virtualNetworks/subnets",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AllowedResourceTypes() = %#v, want %#v", got, want)
	}
}
