// Portions of this file are derived from Microsoft Azure Quick Review (MIT licensed).
// See NOTICE.md for attribution details.

package scanners

import "sort"

// Service describes one assessment service and the Azure resource types it covers.
type Service struct {
	Key           string
	Name          string
	ResourceTypes []string
}

// registry preserves the service/resource-type coverage of the pinned reference implementation.
var registry = map[string][]Service{
	"aa":      {{Key: "aa", Name: "Automation Account", ResourceTypes: []string{"Microsoft.Automation/automationAccounts"}}},
	"adf":     {{Key: "adf", Name: "Data Factory", ResourceTypes: []string{"Microsoft.DataFactory/factories"}}},
	"afd":     {{Key: "afd", Name: "Front Door", ResourceTypes: []string{"Microsoft.Cdn/profiles"}}},
	"afw":     {{Key: "afw", Name: "Azure Firewall", ResourceTypes: []string{"Microsoft.Network/azureFirewalls", "Microsoft.Network/ipGroups"}}},
	"agw":     {{Key: "agw", Name: "Application Gateway", ResourceTypes: []string{"Microsoft.Network/applicationGateways"}}},
	"aif":     {{Key: "aif", Name: "AI Services", ResourceTypes: []string{"Microsoft.CognitiveServices/accounts"}}},
	"aks":     {{Key: "aks", Name: "Azure Kubernetes Service", ResourceTypes: []string{"Microsoft.ContainerService/managedClusters"}}},
	"amg":     {{Key: "amg", Name: "Azure Managed Grafana", ResourceTypes: []string{"Microsoft.Dashboard/grafana"}}},
	"apim":    {{Key: "apim", Name: "API Management", ResourceTypes: []string{"Microsoft.ApiManagement/service"}}},
	"appcs":   {{Key: "appcs", Name: "App Configuration", ResourceTypes: []string{"Microsoft.AppConfiguration/configurationStores"}}},
	"appgw":   {{Key: "appgw", Name: "Application Gateway", ResourceTypes: []string{"Microsoft.Network/applicationGateways"}}},
	"appins":  {{Key: "appins", Name: "Application Insights", ResourceTypes: []string{"Microsoft.Insights/components"}}},
	"appsvc":  {{Key: "appsvc", Name: "App Service", ResourceTypes: []string{"Microsoft.Web/sites", "Microsoft.Web/serverFarms"}}},
	"avd":     {{Key: "avd", Name: "Azure Virtual Desktop", ResourceTypes: []string{"Microsoft.DesktopVirtualization/hostPools", "Microsoft.DesktopVirtualization/workspaces", "Microsoft.DesktopVirtualization/applicationGroups", "Microsoft.DesktopVirtualization/scalingPlans"}}},
	"avs":     {{Key: "avs", Name: "Azure VMware Solution", ResourceTypes: []string{"Microsoft.AVS/privateClouds", "Microsoft.AVS/privateClouds/clusters"}}},
	"bastion": {{Key: "bastion", Name: "Azure Bastion", ResourceTypes: []string{"Microsoft.Network/bastionHosts"}}},
	"ca":      {{Key: "ca", Name: "Container App", ResourceTypes: []string{"Microsoft.App/containerApps"}}},
	"cae":     {{Key: "cae", Name: "Container App Environment", ResourceTypes: []string{"Microsoft.App/managedEnvironments"}}},
	"cache":   {{Key: "cache", Name: "Azure Cache", ResourceTypes: []string{"Microsoft.Cache/redis", "Microsoft.Cache/redisEnterprise"}}},
	"ci":      {{Key: "ci", Name: "Container Instance", ResourceTypes: []string{"Microsoft.ContainerInstance/containerGroups"}}},
	"cog":     {{Key: "cog", Name: "Cognitive Services", ResourceTypes: []string{"Microsoft.CognitiveServices/accounts"}}},
	"cosmos":  {{Key: "cosmos", Name: "Cosmos DB", ResourceTypes: []string{"Microsoft.DocumentDB/databaseAccounts"}}},
	"cr":      {{Key: "cr", Name: "Container Registry", ResourceTypes: []string{"Microsoft.ContainerRegistry/registries"}}},
	"dbw":     {{Key: "dbw", Name: "Databricks Workspace", ResourceTypes: []string{"Microsoft.Databricks/workspaces"}}},
	"disk":    {{Key: "disk", Name: "Managed Disk", ResourceTypes: []string{"Microsoft.Compute/disks"}}},
	"dns":     {{Key: "dns", Name: "DNS Zone", ResourceTypes: []string{"Microsoft.Network/dnsZones", "Microsoft.Network/privateDnsZones"}}},
	"eh":      {{Key: "eh", Name: "Event Hubs", ResourceTypes: []string{"Microsoft.EventHub/namespaces"}}},
	"evgd":    {{Key: "evgd", Name: "Event Grid Domain", ResourceTypes: []string{"Microsoft.EventGrid/domains"}}},
	"evgns":   {{Key: "evgns", Name: "Event Grid Namespace", ResourceTypes: []string{"Microsoft.EventGrid/namespaces"}}},
	"evgst":   {{Key: "evgst", Name: "Event Grid System Topic", ResourceTypes: []string{"Microsoft.EventGrid/systemTopics"}}},
	"evgt":    {{Key: "evgt", Name: "Event Grid Topic", ResourceTypes: []string{"Microsoft.EventGrid/topics"}}},
	"fab":     {{Key: "fab", Name: "Microsoft Fabric", ResourceTypes: []string{"Microsoft.Fabric/capacities"}}},
	"fd":      {{Key: "fd", Name: "Front Door (classic)", ResourceTypes: []string{"Microsoft.Network/frontDoors"}}},
	"func":    {{Key: "func", Name: "Function App", ResourceTypes: []string{"Microsoft.Web/sites"}}},
	"hpc":     {{Key: "hpc", Name: "HPC Cache", ResourceTypes: []string{"Microsoft.StorageCache/caches"}}},
	"ioth":    {{Key: "ioth", Name: "IoT Hub", ResourceTypes: []string{"Microsoft.Devices/IotHubs"}}},
	"kv":      {{Key: "kv", Name: "Key Vault", ResourceTypes: []string{"Microsoft.KeyVault/vaults", "Microsoft.KeyVault/managedHSMs"}}},
	"law":     {{Key: "law", Name: "Log Analytics Workspace", ResourceTypes: []string{"Microsoft.OperationalInsights/workspaces"}}},
	"lb":      {{Key: "lb", Name: "Load Balancer", ResourceTypes: []string{"Microsoft.Network/loadBalancers"}}},
	"logic":   {{Key: "logic", Name: "Logic App", ResourceTypes: []string{"Microsoft.Logic/workflows"}}},
	"mi":      {{Key: "mi", Name: "Managed Identity", ResourceTypes: []string{"Microsoft.ManagedIdentity/userAssignedIdentities"}}},
	"mlw":     {{Key: "mlw", Name: "Machine Learning Workspace", ResourceTypes: []string{"Microsoft.MachineLearningServices/workspaces"}}},
	"mysql":   {{Key: "mysql", Name: "Azure Database for MySQL", ResourceTypes: []string{"Microsoft.DBforMySQL/flexibleServers"}}},
	"nat":     {{Key: "nat", Name: "NAT Gateway", ResourceTypes: []string{"Microsoft.Network/natGateways"}}},
	"nic":     {{Key: "nic", Name: "Network Interface", ResourceTypes: []string{"Microsoft.Network/networkInterfaces"}}},
	"nsg":     {{Key: "nsg", Name: "Network Security Group", ResourceTypes: []string{"Microsoft.Network/networkSecurityGroups"}}},
	"pip":     {{Key: "pip", Name: "Public IP Address", ResourceTypes: []string{"Microsoft.Network/publicIPAddresses"}}},
	"postgres": {{Key: "postgres", Name: "Azure Database for PostgreSQL", ResourceTypes: []string{"Microsoft.DBforPostgreSQL/flexibleServers"}}},
	"privdns": {{Key: "privdns", Name: "Private DNS Zone", ResourceTypes: []string{"Microsoft.Network/privateDnsZones"}}},
	"privend": {{Key: "privend", Name: "Private Endpoint", ResourceTypes: []string{"Microsoft.Network/privateEndpoints"}}},
	"redis": {
		{Key: "redis", Name: "Azure Cache for Redis", ResourceTypes: []string{"Microsoft.Cache/redis"}},
		{Key: "redis", Name: "Azure Managed Redis", ResourceTypes: []string{"Microsoft.Cache/redisEnterprise"}},
	},
	"resource": {{Key: "resource", Name: "Azure Resource", ResourceTypes: []string{"Microsoft.Resources/subscriptions/resourceGroups", "Microsoft.Resources/tags"}}},
	"route":    {{Key: "route", Name: "Route Table", ResourceTypes: []string{"Microsoft.Network/routeTables"}}},
	"sap":      {{Key: "sap", Name: "SAP Virtual Instance", ResourceTypes: []string{"Microsoft.Workloads/sapVirtualInstances"}}},
	"sb":       {{Key: "sb", Name: "Service Bus", ResourceTypes: []string{"Microsoft.ServiceBus/namespaces"}}},
	"sig":      {{Key: "sig", Name: "Compute Gallery", ResourceTypes: []string{"Microsoft.Compute/galleries"}}},
	"sql":      {{Key: "sql", Name: "SQL Database", ResourceTypes: []string{"Microsoft.Sql/servers", "Microsoft.Sql/servers/databases"}}},
	"sqlmi":    {{Key: "sqlmi", Name: "SQL Managed Instance", ResourceTypes: []string{"Microsoft.Sql/managedInstances"}}},
	"srch":     {{Key: "srch", Name: "Search Service", ResourceTypes: []string{"Microsoft.Search/searchServices"}}},
	"st":       {{Key: "st", Name: "Storage Account", ResourceTypes: []string{"Microsoft.Storage/storageAccounts"}}},
	"sub":      {{Key: "sub", Name: "Subscription", ResourceTypes: []string{"Microsoft.Subscription/subscriptions"}}},
	"synw":     {{Key: "synw", Name: "Synapse Workspace", ResourceTypes: []string{"Microsoft.Synapse/workspaces", "Microsoft.Synapse/workspaces/bigDataPools", "Microsoft.Synapse/workspaces/sqlPools"}}},
	"traf":     {{Key: "traf", Name: "Traffic Manager", ResourceTypes: []string{"Microsoft.Network/trafficManagerProfiles"}}},
	"vdpool":   {{Key: "vdpool", Name: "Virtual Desktop Host Pool", ResourceTypes: []string{"Microsoft.DesktopVirtualization/hostPools", "Microsoft.DesktopVirtualization/scalingPlans", "Microsoft.DesktopVirtualization/workspaces"}}},
	"vgw":      {{Key: "vgw", Name: "Virtual Network Gateway", ResourceTypes: []string{"Microsoft.Network/virtualNetworkGateways"}}},
	"vhub":     {{Key: "vhub", Name: "Virtual Hub", ResourceTypes: []string{"Microsoft.Network/virtualHubs"}}},
	"vm":       {{Key: "vm", Name: "Virtual Machine", ResourceTypes: []string{"Microsoft.Compute/virtualMachines"}}},
	"vmss":     {{Key: "vmss", Name: "Virtual Machine Scale Set", ResourceTypes: []string{"Microsoft.Compute/virtualMachineScaleSets"}}},
	"vnet":     {{Key: "vnet", Name: "Virtual Network", ResourceTypes: []string{"Microsoft.Network/virtualNetworks", "Microsoft.Network/virtualNetworks/subnets"}}},
	"vpng":     {{Key: "vpng", Name: "VPN Gateway", ResourceTypes: []string{"Microsoft.Network/vpnGateways"}}},
	"vpns":     {{Key: "vpns", Name: "VPN Site", ResourceTypes: []string{"Microsoft.Network/vpnSites"}}},
	"vrouter":  {{Key: "vrouter", Name: "Virtual Router", ResourceTypes: []string{"Microsoft.Network/virtualRouters"}}},
	"vwan":     {{Key: "vwan", Name: "Virtual WAN", ResourceTypes: []string{"Microsoft.Network/virtualWans"}}},
	"wps":      {{Key: "wps", Name: "Web PubSub", ResourceTypes: []string{"Microsoft.SignalRService/webPubSub"}}},
}

// Keys returns all scanner keys in deterministic order.
func Keys() []string {
	keys := make([]string, 0, len(registry))
	for key := range registry {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// ByKey returns a copy of all service entries for a scanner key.
func ByKey(key string) []Service {
	services := registry[key]
	out := make([]Service, len(services))
	copy(out, services)
	return out
}

// All returns every registered service entry in scanner-key order.
func All() []Service {
	var out []Service
	for _, key := range Keys() {
		out = append(out, ByKey(key)...)
	}
	return out
}

// SelectedKeys reproduces the pinned reference scanner-selection semantics.
// When a normal scan supplies multiple available scanner keys and the filter has
// include.resourceTypes values, those values are interpreted as scanner keys.
// A scanner-specific command (one key) wins over the filter. With no keys, all
// registered scanners are selected.
func SelectedKeys(scannerKeys, configuredResourceTypes []string) []string {
	var candidates []string
	switch {
	case len(scannerKeys) > 1 && len(configuredResourceTypes) > 0:
		candidates = configuredResourceTypes
	case len(scannerKeys) >= 1:
		candidates = scannerKeys
	default:
		return Keys()
	}

	out := make([]string, 0, len(candidates))
	for _, key := range candidates {
		if _, exists := registry[key]; exists {
			out = append(out, key)
		}
	}
	return out
}

// ResourceTypes returns the distinct Azure resource types covered by the selected keys.
func ResourceTypes(keys []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, key := range keys {
		for _, service := range registry[key] {
			for _, resourceType := range service.ResourceTypes {
				if _, ok := seen[resourceType]; ok {
					continue
				}
				seen[resourceType] = struct{}{}
				out = append(out, resourceType)
			}
		}
	}
	return out
}

// AllowedResourceTypes combines scanner selection and resource-type expansion for
// installing runtime scope into the filter model.
func AllowedResourceTypes(scannerKeys, configuredResourceTypes []string) []string {
	return ResourceTypes(SelectedKeys(scannerKeys, configuredResourceTypes))
}
