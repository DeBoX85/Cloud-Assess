package skus

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed known_skus.yaml
var rawYAML []byte

type SKU struct {
	Name   string `yaml:"name"`
	VCPUs  int    `yaml:"vcpus"`
	Family string `yaml:"family"`
}

var skuMap = load()

func load() map[string]SKU {
	var entries []SKU
	if err := yaml.Unmarshal(rawYAML, &entries); err != nil {
		return map[string]SKU{}
	}
	out := make(map[string]SKU, len(entries))
	for _, entry := range entries {
		out[entry.Name] = entry
	}
	return out
}

func Lookup(name string) (SKU, bool) {
	sku, ok := skuMap[name]
	return sku, ok
}

// ComputeCapacity preserves the pinned reference report behavior. VM scale sets report
// instance count multiplied by vCPU count when the SKU is known; otherwise capacity is
// the reported SKU capacity, with a vCPU fallback when no explicit capacity exists.
func ComputeCapacity(skuName string, skuCapacity int, isVMSS bool) string {
	if skuCapacity > 0 {
		if isVMSS {
			if sku, ok := Lookup(skuName); ok && sku.VCPUs > 0 {
				return fmt.Sprint(skuCapacity * sku.VCPUs)
			}
		}
		return fmt.Sprint(skuCapacity)
	}
	if sku, ok := Lookup(skuName); ok && sku.VCPUs > 0 {
		return fmt.Sprint(sku.VCPUs)
	}
	return ""
}
