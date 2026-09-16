// Portions of this file reproduce diagnostic-settings recommendation metadata from
// Microsoft Azure Quick Review (MIT licensed). See NOTICE.md.

package diagnostics

import (
	"sort"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const (
	Source                         = "DIAGNOSTICS"
	ValidationAzureResourceManager = "Azure Resource Manager"
	CategoryMonitoringAndAlerting  = "MonitoringAndAlerting"
)

type recommendationMetadata struct {
	ID        string
	Text      string
	LearnMore string
}

// supportedResourceTypes mirrors the reference implementation's diagnostic-settings
// support table. A nil value means the resource type is queried for diagnostic settings
// but intentionally has no dedicated recommendation definition.
var supportedResourceTypes = map[string]*recommendationMetadata{
	"microsoft.datafactory/factories": {
		ID:        "adf-001",
		Text:      "Azure Data Factory should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/data-factory/monitor-configure-diagnostics",
	},
	"microsoft.cdn/profiles": {
		ID:        "afd-001",
		Text:      "Azure FrontDoor should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/frontdoor/standard-premium/how-to-logs",
	},
	"microsoft.network/azurefirewalls": {
		ID:        "afw-001",
		Text:      "Azure Firewall should have diagnostic settings enabled",
		LearnMore: "https://docs.microsoft.com/en-us/azure/firewall/logs-and-metrics",
	},
	"microsoft.network/applicationgateways": {
		ID:        "agw-005",
		Text:      "Application Gateway: Monitor and Log the configurations and traffic",
		LearnMore: "https://learn.microsoft.com/en-us/azure/application-gateway/application-gateway-diagnostics#diagnostic-logging",
	},
	"microsoft.cognitiveservices/accounts": {
		ID:        "aif-001",
		Text:      "Service should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/event-hubs/monitor-event-hubs#collection-and-routing",
	},
	"microsoft.containerservice/managedclusters": {
		ID:        "aks-001",
		Text:      "AKS Cluster should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/aks/monitor-aks#collect-resource-logs",
	},
	"microsoft.apimanagement/service": {
		ID:        "apim-001",
		Text:      "APIM should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/api-management/api-management-howto-use-azure-monitor#resource-logs",
	},
	"microsoft.appconfiguration/configurationstores": {
		ID:        "appcs-001",
		Text:      "AppConfiguration should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/azure-app-configuration/monitor-app-configuration?tabs=portal",
	},
	"microsoft.analysisservices/servers": {
		ID:        "as-001",
		Text:      "Azure Analysis Service should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/analysis-services/analysis-services-logging",
	},
	"microsoft.web/serverfarms": {
		ID:   "asp-001",
		Text: "Plan should have diagnostic settings enabled",
	},
	"microsoft.web/sites": {
		ID:        "app-001",
		Text:      "App should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/app-service/troubleshoot-diagnostic-logs#send-logs-to-azure-monitor",
	},
	"microsoft.app/managedenvironments": {
		ID:        "cae-001",
		Text:      "Container Apps Environment should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/container-apps/log-options#diagnostic-settings",
	},
	"microsoft.documentdb/databaseaccounts": {
		ID:        "cosmos-001",
		Text:      "CosmosDB should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/cosmos-db/monitor-resource-logs",
	},
	"microsoft.containerregistry/registries": {
		ID:        "cr-001",
		Text:      "ContainerRegistry should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/container-registry/monitor-service",
	},
	"microsoft.databricks/workspaces": {
		ID:        "dbw-001",
		Text:      "Azure Databricks should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/databricks/administration-guide/account-settings/audit-log-delivery",
	},
	"microsoft.kusto/clusters": {
		ID:        "dec-001",
		Text:      "Azure Data Explorer should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/data-explorer/using-diagnostic-logs",
	},
	"microsoft.eventgrid/domains": {
		ID:        "evgd-001",
		Text:      "Event Grid Domain should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/event-grid/diagnostic-logs",
	},
	"microsoft.eventhub/namespaces": {
		ID:        "evh-001",
		Text:      "Event Hub Namespace should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/event-hubs/monitor-event-hubs#collection-and-routing",
	},
	"microsoft.machinelearningservices/workspaces": {
		ID:        "hub-006",
		Text:      "Service should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/event-hubs/monitor-event-hubs#collection-and-routing",
	},
	"microsoft.keyvault/vaults": {
		ID:        "kv-001",
		Text:      "Key Vault should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/key-vault/general/monitor-key-vault",
	},
	"microsoft.network/loadbalancers": {
		ID:        "lb-001",
		Text:      "Load Balancer should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/load-balancer/monitor-load-balancer#creating-a-diagnostic-setting",
	},
	"microsoft.logic/workflows": {
		ID:        "logic-001",
		Text:      "Logic App should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/logic-apps/monitor-workflows-collect-diagnostic-data",
	},
	"microsoft.dbformariadb/servers": {
		ID:   "maria-001",
		Text: "MariaDB should have diagnostic settings enabled",
	},
	"microsoft.dbformysql/servers": {
		ID:        "mysql-001",
		Text:      "Azure Database for MySQL - Single Server should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/mysql/single-server/concepts-monitoring#server-logs",
	},
	"microsoft.dbformysql/flexibleservers": {
		ID:        "mysqlf-001",
		Text:      "Azure Database for MySQL - Flexible Server should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/mysql/flexible-server/tutorial-query-performance-insights#set-up-diagnostics",
	},
	"microsoft.network/natgateways": {
		ID:        "ng-001",
		Text:      "NAT Gateway should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/nat-gateway/nat-metrics",
	},
	"microsoft.network/networksecuritygroups": {
		ID:        "nsg-001",
		Text:      "NSG should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/virtual-network/virtual-network-nsg-manage-log",
	},
	"microsoft.dbforpostgresql/servers": {
		ID:        "psql-001",
		Text:      "PostgreSQL should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/postgresql/single-server/concepts-server-logs#resource-logs",
	},
	"microsoft.dbforpostgresql/flexibleservers": {
		ID:        "psqlf-001",
		Text:      "PostgreSQL should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/howto-configure-and-access-logs",
	},
	"microsoft.cache/redis": {
		ID:        "redis-001",
		Text:      "Redis should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/cache-monitor-diagnostic-settings",
	},
	"microsoft.servicebus/namespaces": {
		ID:        "sb-001",
		Text:      "Service Bus should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/service-bus-messaging/monitor-service-bus#collection-and-routing",
	},
	"microsoft.signalrservice/signalr": {
		ID:        "sigr-001",
		Text:      "SignalR should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/azure-signalr/signalr-howto-diagnostic-logs",
	},
	"microsoft.sql/servers/databases": {
		ID:   "sqldb-001",
		Text: "SQL Database should have diagnostic settings enabled",
	},
	"microsoft.search/searchservices": {
		ID:        "srch-006",
		Text:      "Azure AI Search should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/search/search-monitor-enable-logging",
	},
	"microsoft.storage/storageaccounts": {
		ID:        "st-001",
		Text:      "Storage should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/storage/blobs/monitor-blob-storage",
	},
	"microsoft.synapse/workspaces": {
		ID:        "synw-001",
		Text:      "Azure Synapse Workspace should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/data-factory/monitor-configure-diagnostics",
	},
	"microsoft.network/trafficmanagerprofiles": {
		ID:        "traf-001",
		Text:      "Traffic Manager should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/traffic-manager/traffic-manager-diagnostic-logs",
	},
	"microsoft.network/virtualnetworkgateways": {
		ID:        "vgw-001",
		Text:      "Virtual Network Gateway should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/vpn-gateway/monitor-vpn-gateway",
	},
	"microsoft.network/virtualnetworks": {
		ID:        "vnet-001",
		Text:      "Virtual Network should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/virtual-network/monitor-virtual-network#collection-and-routing",
	},
	"microsoft.network/virtualwans": {
		ID:        "vwa-001",
		Text:      "Virtual WAN should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/virtual-wan/monitor-virtual-wan",
	},
	"microsoft.signalrservice/webpubsub": {
		ID:        "wps-001",
		Text:      "Web Pub Sub should have diagnostic settings enabled",
		LearnMore: "https://learn.microsoft.com/en-us/azure/azure-web-pubsub/howto-troubleshoot-resource-logs",
	},

	// Legacy/system resource types that the reference implementation queries for
	// diagnostic settings but does not associate with a dedicated recommendation.
	"microsoft.network/networkinterfaces":                       nil,
	"microsoft.network/routetables":                             nil,
	"microsoft.recoveryservices/vaults":                         nil,
	"specialized.workload/avd":                                  nil,
	"microsoft.compute/virtualmachines":                         nil,
	"microsoft.network/virtualnetworks/subnets":                 nil,
	"specialized.workload/hpc":                                  nil,
	"microsoft.automation/automationaccounts":                   nil,
	"microsoft.dashboard/grafana":                               nil,
	"microsoft.virtualmachineimages/imagetemplates":             nil,
	"microsoft.devices/iothubs":                                 nil,
	"microsoft.compute/disks":                                   nil,
	"microsoft.network/connections":                             nil,
	"microsoft.app/containerapps":                               nil,
	"microsoft.network/frontdoorwebapplicationfirewallpolicies": nil,
	"microsoft.batch/batchaccounts":                             nil,
	"microsoft.network/publicipaddresses":                       nil,
	"microsoft.sql/servers":                                     nil,
	"microsoft.sql/servers/elasticpools":                        nil,
	"microsoft.operationalinsights/workspaces":                  nil,
	"microsoft.insights/components":                             nil,
	"microsoft.compute/virtualmachinescalesets":                 nil,
	"microsoft.network/privateendpoints":                        nil,
	"microsoft.containerinstance/containergroups":               nil,
	"microsoft.resources/resourcegroups":                        nil,
	"microsoft.network/ipgroups":                                nil,
	"microsoft.compute/galleries":                               nil,
	"microsoft.network/privatednszones":                         nil,
	"microsoft.network/networkwatchers":                         nil,
	"microsoft.compute/availabilitysets":                        nil,
	"microsoft.web/connections":                                 nil,
	"microsoft.web/certificates":                                nil,
	"specialized.workload/sap":                                  nil,
}

func Supports(resourceType string) bool {
	_, ok := supportedResourceTypes[normalize(resourceType)]
	return ok
}

func RecommendationFor(resourceType string) (assessment.RecommendationDefinition, bool) {
	metadata, ok := supportedResourceTypes[normalize(resourceType)]
	if !ok || metadata == nil {
		return assessment.RecommendationDefinition{}, false
	}

	definition := assessment.RecommendationDefinition{
		ID:                  metadata.ID,
		Recommendation:      metadata.Text,
		Category:            CategoryMonitoringAndAlerting,
		Impact:              assessment.ImpactLow,
		ResourceType:        normalize(resourceType),
		Source:              Source,
		ValidationMechanism: ValidationAzureResourceManager,
	}
	if metadata.LearnMore != "" {
		definition.LearnMore = []assessment.LearnMoreLink{{
			Name: "Diagnostic Settings",
			URL:  metadata.LearnMore,
		}}
	}
	return definition, true
}

func Recommendations() []assessment.RecommendationDefinition {
	definitions := make([]assessment.RecommendationDefinition, 0, len(supportedResourceTypes))
	for resourceType := range supportedResourceTypes {
		if definition, ok := RecommendationFor(resourceType); ok {
			definitions = append(definitions, definition)
		}
	}
	sort.Slice(definitions, func(i, j int) bool {
		if definitions[i].ResourceType != definitions[j].ResourceType {
			return definitions[i].ResourceType < definitions[j].ResourceType
		}
		return definitions[i].ID < definitions[j].ID
	})
	return definitions
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
