package scanners

import "testing"

func TestRegistryAccessorsOwnNestedResourceTypes(t *testing.T) {
	for _, accessor := range []string{"ByKey", "All"} {
		t.Run(accessor, func(t *testing.T) {
			original := registry["st"][0].ResourceTypes[0]
			defer func() { registry["st"][0].ResourceTypes[0] = original }()
			services := ByKey("st")
			if accessor == "All" {
				services = All()
			}
			for i := range services {
				if services[i].Key == "st" {
					services[i].ResourceTypes[0] = "Microsoft.Test/changed"
				}
			}
			fresh := ByKey("st")
			if len(fresh) != 1 || len(fresh[0].ResourceTypes) != 1 || fresh[0].ResourceTypes[0] != "Microsoft.Storage/storageAccounts" {
				t.Fatalf("caller mutated global scanner resource-type scope through %s: %+v", accessor, fresh)
			}
		})
	}
}
