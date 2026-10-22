// SQL 2021-11-01 bounded wire shapes derived from MIT-licensed Microsoft REST specifications.
// Pinned source/complete license and all actual26definitions are frozen by the SQL contract.
package region

import (
	"context"
	"encoding/json"
	"math"
	"strings"
)

var sqlCapabilityShapes = map[string]map[string]string{
	"AutoPauseDelayTimeRange":                             {"minValue": "int32", "maxValue": "int32", "stepSize": "int32", "default": "int32", "unit": "string", "doNotPauseValue": "int32"},
	"EditionCapability":                                   {"name": "string", "supportedServiceLevelObjectives": "[]ServiceObjectiveCapability", "zoneRedundant": "bool", "readScale": "ReadScaleCapability", "supportedStorageCapabilities": "[]StorageCapability", "status": "string", "reason": "string"},
	"ElasticPoolEditionCapability":                        {"name": "string", "supportedElasticPoolPerformanceLevels": "[]ElasticPoolPerformanceLevelCapability", "zoneRedundant": "bool", "status": "string", "reason": "string"},
	"ElasticPoolPerDatabaseMaxPerformanceLevelCapability": {"limit": "double", "unit": "string", "supportedPerDatabaseMinPerformanceLevels": "[]ElasticPoolPerDatabaseMinPerformanceLevelCapability", "status": "string", "reason": "string"},
	"ElasticPoolPerDatabaseMinPerformanceLevelCapability": {"limit": "double", "unit": "string", "status": "string", "reason": "string"},
	"ElasticPoolPerformanceLevelCapability":               {"performanceLevel": "PerformanceLevelCapability", "sku": "Sku", "supportedLicenseTypes": "[]LicenseTypeCapability", "maxDatabaseCount": "int32", "includedMaxSize": "MaxSizeCapability", "supportedMaxSizes": "[]MaxSizeRangeCapability", "supportedPerDatabaseMaxSizes": "[]MaxSizeRangeCapability", "supportedPerDatabaseMaxPerformanceLevels": "[]ElasticPoolPerDatabaseMaxPerformanceLevelCapability", "zoneRedundant": "bool", "supportedMaintenanceConfigurations": "[]MaintenanceConfigurationCapability", "status": "string", "reason": "string"},
	"InstancePoolEditionCapability":                       {"name": "string", "supportedFamilies": "[]InstancePoolFamilyCapability", "status": "string", "reason": "string"},
	"InstancePoolFamilyCapability":                        {"name": "string", "supportedLicenseTypes": "[]LicenseTypeCapability", "supportedVcoresValues": "[]InstancePoolVcoresCapability", "status": "string", "reason": "string"},
	"InstancePoolVcoresCapability":                        {"name": "string", "value": "int32", "storageLimit": "MaxSizeCapability", "status": "string", "reason": "string"},
	"LicenseTypeCapability":                               {"name": "string", "status": "string", "reason": "string"},
	"LocationCapabilities":                                {"name": "string", "supportedServerVersions": "[]ServerVersionCapability", "supportedManagedInstanceVersions": "[]ManagedInstanceVersionCapability", "status": "string", "reason": "string"},
	"LogSizeCapability":                                   {"limit": "int32", "unit": "string"},
	"MaintenanceConfigurationCapability":                  {"name": "string", "zoneRedundant": "bool", "status": "string", "reason": "string"},
	"ManagedInstanceEditionCapability":                    {"name": "string", "supportedFamilies": "[]ManagedInstanceFamilyCapability", "supportedStorageCapabilities": "[]StorageCapability", "zoneRedundant": "bool", "status": "string", "reason": "string"},
	"ManagedInstanceFamilyCapability":                     {"name": "string", "sku": "string", "supportedLicenseTypes": "[]LicenseTypeCapability", "supportedVcoresValues": "[]ManagedInstanceVcoresCapability", "status": "string", "reason": "string"},
	"ManagedInstanceMaintenanceConfigurationCapability":   {"name": "string", "status": "string", "reason": "string"},
	"ManagedInstanceVcoresCapability":                     {"name": "string", "value": "int32", "includedMaxSize": "MaxSizeCapability", "supportedStorageSizes": "[]MaxSizeRangeCapability", "instancePoolSupported": "bool", "standaloneSupported": "bool", "supportedMaintenanceConfigurations": "[]ManagedInstanceMaintenanceConfigurationCapability", "status": "string", "reason": "string"},
	"ManagedInstanceVersionCapability":                    {"name": "string", "supportedEditions": "[]ManagedInstanceEditionCapability", "supportedInstancePoolEditions": "[]InstancePoolEditionCapability", "status": "string", "reason": "string"},
	"MaxSizeCapability":                                   {"limit": "int32", "unit": "string"},
	"MaxSizeRangeCapability":                              {"minValue": "MaxSizeCapability", "maxValue": "MaxSizeCapability", "scaleSize": "MaxSizeCapability", "logSize": "LogSizeCapability", "status": "string", "reason": "string"},
	"MinCapacityCapability":                               {"value": "double", "status": "string", "reason": "string"},
	"PerformanceLevelCapability":                          {"value": "double", "unit": "string"},
	"ReadScaleCapability":                                 {"maxNumberOfReplicas": "int32", "status": "string", "reason": "string"},
	"ServerVersionCapability":                             {"name": "string", "supportedEditions": "[]EditionCapability", "supportedElasticPoolEditions": "[]ElasticPoolEditionCapability", "status": "string", "reason": "string"},
	"ServiceObjectiveCapability":                          {"id": "string", "name": "string", "supportedMaxSizes": "[]MaxSizeRangeCapability", "performanceLevel": "PerformanceLevelCapability", "sku": "Sku", "supportedLicenseTypes": "[]LicenseTypeCapability", "includedMaxSize": "MaxSizeCapability", "zoneRedundant": "bool", "supportedAutoPauseDelay": "AutoPauseDelayTimeRange", "supportedMinCapacities": "[]MinCapacityCapability", "computeModel": "string", "supportedMaintenanceConfigurations": "[]MaintenanceConfigurationCapability", "status": "string", "reason": "string"},
	"StorageCapability":                                   {"storageAccountType": "string", "status": "string", "reason": "string"},
	"Sku":                                                 {"name": "string", "tier": "string", "size": "string", "family": "string", "capacity": "int32"},
}

func sqlCapabilityShape(ctx context.Context, raw json.RawMessage, kind string) bool {
	if ctx.Err() != nil {
		return false
	}
	fields, found := sqlCapabilityShapes[kind]
	if !found {
		return false
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	row, ok := providerObject(raw, keys...)
	if !ok {
		return false
	}
	if kind == "Sku" {
		if _, has := row["name"]; !has {
			return false
		}
	}
	for key, value := range row {
		if ctx.Err() != nil {
			return false
		}
		ty, known := fields[key]
		if !known {
			continue
		}
		switch ty {
		case "string":
			var s string
			if !diskScalar(value, &s) {
				return false
			}
		case "bool":
			var b bool
			if !diskScalar(value, &b) {
				return false
			}
		case "int32":
			var n int64
			if !diskScalar(value, &n) || n < -2147483648 || n > 2147483647 {
				return false
			}
		case "double":
			var n float64
			if !diskScalar(value, &n) || math.IsInf(n, 0) || math.IsNaN(n) {
				return false
			}
		default:
			if strings.HasPrefix(ty, "[]") {
				rows, valid := providerArray(value)
				if !valid {
					return false
				}
				for _, item := range rows {
					if !sqlCapabilityShape(ctx, item, strings.TrimPrefix(ty, "[]")) {
						return false
					}
				}
			} else if !sqlCapabilityShape(ctx, value, ty) {
				return false
			}
		}
	}
	return ctx.Err() == nil
}
