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
	"appi":    {{Key: "appi", Name: "Application Insights", ResourceTypes: []string{"Microsoft.Insights/components", "Microsoft.Insights/activityLogAlerts"}}},
	"arc":     {{Key: "arc", Name: "Azure Arc", ResourceTypes: []string{"Microsoft.AzureArcData/sqlServerInstances"}}},
	"as":      {{Key: "as", Name: "Analysis Services", ResourceTypes: []string{"Microsoft.AnalysisServices/servers"}}},
	"asa":     {{Key: "asa", Name: "Stream Analytics Job", ResourceTypes: []string{"Microsoft.StreamAnalytics/streamingJobs"}}},
	"asp":     {{Key: "asp", Name: "App Service Plan", ResourceTypes: []string{"Microsoft.Web/serverFarms", "Microsoft.Web/sites", "Microsoft.Web/connections", "Microsoft.Web/certificates"}}},
	"avail":   {{Key: "avail", Name: "Availability Set", ResourceTypes: []string{"Microsoft.Compute/availabilitySets"}}},
	"avd":     {{Key: "avd", Name: "Azure Virtual Desktop", ResourceTypes: []string{"Specialized.Workload/AVD"}}},
	"avs":     {{Key: "avs", Name: "Azure VMware Solution", ResourceTypes: []string{"Microsoft.AVS/privateClouds", "Specialized.Workload/AVS"}}},
	"ba":      {{Key: "ba", Name: "Batch Account", ResourceTypes: []string{"Microsoft.Batch/batchAccounts"}}},
	"bastion": {{Key: "bastion", Name: "Bastion Host", ResourceTypes: []string{"Microsoft.Network/bastionHosts"}}},
	"ca":      {{Key: "ca", Name: "Container App", ResourceTypes: []string{"Microsoft.App/containerApps"}}},
	"cae":     {{Key: "cae", Name: "Container Apps Environment", ResourceTypes: []string{"Microsoft.App/managedenvironments"}}},
	"ci":      {{Key: "ci", Name: "Container Instance", ResourceTypes: []string{"Microsoft.ContainerInstance/containerGroups"}}},
	"con":     {{Key: "con", Name: "Connection", ResourceTypes: []string{"Microsoft.Network/connections"}}},
	"cosmos":  {{Key: "cosmos", Name: "Cosmos DB", ResourceTypes: []string{"Microsoft.DocumentDB/databaseAccounts"}}},
	"cr":      {{Key: "cr", Name: "Container Registry", ResourceTypes: []string{"Microsoft.ContainerRegistry/registries"}}},
	"dbw":     {{Key: "dbw", Name: "Databricks Workspace", ResourceTypes: []string{"Microsoft.Databricks/workspaces"}}},
	"ddos":    {{Key: "ddos", Name: "DDoS Protection Plan", ResourceTypes: []string{"Microsoft.Network/ddosProtectionPlans"}}},
	"dec":     {{Key: "dec", Name: "Data Explorer Cluster", ResourceTypes: []string{"Microsoft.Kusto/clusters"}}},
	"disk":    {{Key: "disk", Name: "Disk", ResourceTypes: []string{"Microsoft.Compute/disks"}}},
	"dnsres":  {{Key: "dnsres", Name: "DNS Resolver", ResourceTypes: []string{"Microsoft.Network/dnsResolvers"}}},
	"dnsz":    {{Key: "dnsz", Name: "DNS Zone", ResourceTypes: []string{"Microsoft.Network/dnsZones"}}},
	"domain":  {{Key: "domain", Name: "Domain Services", ResourceTypes: []string{"Microsoft.AAD/domainServices"}}},
	"erc":     {{Key: "erc", Name: "ExpressRoute Circuit", ResourceTypes: []string{"Microsoft.Network/expressRouteCircuits", "Microsoft.Network/ExpressRoutePorts", "Microsoft.Network/expressRouteGateways"}}},
	"evgd":    {{Key: "evgd", Name: "Event Grid Domain", ResourceTypes: []string{"Microsoft.EventGrid/domains"}}},
	"evgt":    {{Key: "evgt", Name: "Event Grid Topic", ResourceTypes: []string{"Microsoft.EventGrid/topics"}}},
	"evh":     {{Key: "evh", Name: "Event Hub", ResourceTypes: []string{"Microsoft.EventHub/namespaces"}}},
	"fabric":  {{Key: "fabric", Name: "Fabric", ResourceTypes: []string{"Microsoft.Fabric/capacities"}}},
	"fdfp":    {{Key: "fdfp", Name: "Front Door Firewall Policy", ResourceTypes: []string{"Microsoft.Network/frontdoorWebApplicationFirewallPolicies"}}},
	"gal":     {{Key: "gal", Name: "Compute Gallery", ResourceTypes: []string{"Microsoft.Compute/galleries"}}},
	"hpc":     {{Key: "hpc", Name: "HPC", ResourceTypes: []string{"Specialized.Workload/HPC"}}},
	"hub":     {{Key: "hub", Name: "Machine Learning Workspace", ResourceTypes: []string{"Microsoft.MachineLearningServices/workspaces", "Microsoft.MachineLearningServices/registries"}}},
	"iot":     {{Key: "iot", Name: "IoT Hub", ResourceTypes: []string{"Microsoft.Devices/IotHubs"}}},
	"it":      {{Key: "it", Name: "Image Template", ResourceTypes: []string{"Microsoft.VirtualMachineImages/imageTemplates"}}},
	"kv":      {{Key: "kv", Name: "Key Vault", ResourceTypes: []string{"Microsoft.KeyVault/vaults"}}},
	"lb":      {{Key: "lb", Name: "Load Balancer", ResourceTypes: []string{"Microsoft.Network/loadBalancers"}}},
	"log":     {{Key: "log", Name: "Log Analytics Workspace", ResourceTypes: []string{"Microsoft.OperationalInsights/workspaces"}}},
	"logic":   {{Key: "logic", Name: "Logic App", ResourceTypes: []string{"Microsoft.Logic/workflows"}}},
	"mysql":   {{Key: "mysql", Name: "MySQL Database", ResourceTypes: []string{"Microsoft.DBforMySQL/servers", "Microsoft.DBforMySQL/flexibleServers"}}},
	"netapp":  {{Key: "netapp", Name: "NetApp Account", ResourceTypes: []string{"Microsoft.NetApp/netAppAccounts"}}},
	"ng":      {{Key: "ng", Name: "NAT Gateway", ResourceTypes: []string{"Microsoft.Network/natGateways"}}},
	"nic":     {{Key: "nic", Name: "Network Interface", ResourceTypes: []string{"Microsoft.Network/networkInterfaces"}}},
	"nsg":     {{Key: "nsg", Name: "Network Security Group", ResourceTypes: []string{"Microsoft.Network/networkSecurityGroups"}}},
	"ntc":     {{Key: "ntc", Name: "Azure Traffic Collector", ResourceTypes: []string{"Microsoft.NetworkFunction/azureTrafficCollectors"}}},
	"nw":      {{Key: "nw", Name: "Network Watcher", ResourceTypes: []string{"Microsoft.Network/networkWatchers"}}},
	"odb":     {{Key: "odb", Name: "Oracle Database", ResourceTypes: []string{"Oracle.Database/cloudExadataInfrastructures", "Oracle.Database/cloudVmClusters"}}},
	"p2svpng": {{Key: "p2svpng", Name: "P2S VPN Gateway", ResourceTypes: []string{"Microsoft.Network/p2sVpnGateways"}}},
	"pdnsz":   {{Key: "pdnsz", Name: "Private DNS Zone", ResourceTypes: []string{"Microsoft.Network/privateDnsZones"}}},
	"pep":     {{Key: "pep", Name: "Private Endpoint", ResourceTypes: []string{"Microsoft.Network/privateEndpoints"}}},
	"pip":     {{Key: "pip", Name: "Public IP Address", ResourceTypes: []string{"Microsoft.Network/publicIPAddresses"}}},
	"psql":    {{Key: "psql", Name: "PostgreSQL Database", ResourceTypes: []string{"Microsoft.DBforPostgreSQL/servers", "Microsoft.DBforPostgreSQL/flexibleServers"}}},
	"redis": {
		{Key: "redis", Name: "Redis Cache", ResourceTypes: []string{"Microsoft.Cache/Redis"}},
		{Key: "redis", Name: "Redis Enterprise", ResourceTypes: []string{"Microsoft.Cache/redisEnterprise"}},
	},
	"resource": {{Key: "resource", Name: "Resource", ResourceTypes: []string{"Microsoft.Resources"}}},
	"rg":       {{Key: "rg", Name: "Resource Group", ResourceTypes: []string{"Microsoft.Resources/resourceGroups"}}},
	"rsv":      {{Key: "rsv", Name: "Recovery Services Vault", ResourceTypes: []string{"Microsoft.RecoveryServices/vaults"}}},
	"rt":       {{Key: "rt", Name: "Route Table", ResourceTypes: []string{"Microsoft.Network/routeTables"}}},
	"sap":      {{Key: "sap", Name: "SAP", ResourceTypes: []string{"Specialized.Workload/SAP"}}},
	"sb":       {{Key: "sb", Name: "Service Bus", ResourceTypes: []string{"Microsoft.ServiceBus/namespaces"}}},
	"sigr":     {{Key: "sigr", Name: "SignalR", ResourceTypes: []string{"Microsoft.SignalRService/SignalR"}}},
	"sql":      {{Key: "sql", Name: "SQL Server", ResourceTypes: []string{"Microsoft.Sql/servers", "Microsoft.Sql/servers/databases", "Microsoft.Sql/servers/elasticPools"}}},
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
