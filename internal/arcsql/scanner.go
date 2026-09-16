// Portions of this file reproduce Arc-enabled SQL scan behavior from Microsoft Azure Quick
// Review (MIT licensed). See NOTICE.md.

package arcsql

import (
	"context"
	"fmt"
	"sort"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const Query = `
	resources
	| where type =~ "Microsoft.AzureArcData/sqlServerInstances"
	| extend SQLInstance = id, AzureArcServer = tolower(tostring(properties.containerResourceId))
	| extend version = tostring(properties.version)
	| extend edition = tostring(properties.edition)
	| extend Build = tostring(properties.currentVersion)
	| extend DefenderStatus = tostring(properties.azureDefenderStatus)
	| extend patchLevel = tostring(properties.patchLevel)
	| extend vcores = toint(properties.vCore)
	| join kind=inner (resources
	| where type == 'microsoft.hybridcompute/machines/extensions' 
	| where properties.type == "WindowsAgent.SqlServer"
	| order by ['id'] asc
	| extend License = case(properties.settings.LicenseType == "Paid","SA",properties.settings.LicenseType == "PAYG","PAYG","unset")
	| extend Serverid = tolower(tostring(split(id,'/extensions/WindowsAgent.SqlServer')[0]))
	| parse properties with * 'uploadStatus : ' DPSStatus ';' *
	| parse properties with * 'telemetryUploadStatus : ' TELStatusRaw ';' *
	| extend DPSStatus = iff(DPSStatus == "", "No Data",DPSStatus)
	| extend TELStatuslogs = (parse_json(replace('.\"','\"',TELStatusRaw))).logs
	| extend TELStatus = iff(TELStatuslogs.status == "OK","__",iff(TELStatuslogs.message == "","No Data",TELStatuslogs.message))
	) on $left.AzureArcServer == $right.Serverid
	| join kind=inner (resources
	| where type == "microsoft.hybridcompute/machines"
	| extend status = tostring(properties.status)
	| project id = tolower(id),status) on $left.AzureArcServer == $right.id
	| project subscriptionId,status,AzureArcServer,SQLInstance,resourceGroup,version,Build,patchLevel,edition,vcores,License,DPSStatus,TELStatus,DefenderStatus
	`

type GraphQuerier interface {
	Query(context.Context, string, map[string]string, ...arg.QueryOptions) (*arg.Result, error)
}

type Filter interface {
	IsSubscriptionExcluded(subscriptionID string) bool
	IsServiceExcluded(resourceID string) bool
}

type Result struct {
	Records  []assessment.ArcSQLRecord
	Warnings []assessment.AssessmentWarning
}

type Scanner struct {
	graph GraphQuerier
}

func New(credential azcore.TokenCredential) *Scanner {
	return NewWithClient(arg.NewClient(arg.NewHTTPTransport(credential)))
}

func NewWithClient(graphClient GraphQuerier) *Scanner {
	return &Scanner{graph: graphClient}
}

func (s *Scanner) Scan(
	ctx context.Context,
	subscriptions map[string]string,
	filter Filter,
) (Result, error) {
	result := Result{
		Records:  []assessment.ArcSQLRecord{},
		Warnings: []assessment.AssessmentWarning{},
	}
	if s == nil || s.graph == nil {
		return result, fmt.Errorf("Arc SQL Resource Graph client is not configured")
	}

	graphResult, err := s.graph.Query(ctx, Query, subscriptions)
	if err != nil {
		return result, fmt.Errorf("query Arc SQL resources: %w", err)
	}
	if graphResult == nil {
		return result, fmt.Errorf("query Arc SQL resources: Resource Graph returned nil result")
	}

	rows, malformed := arg.DecodeRowsWithStats[arcSQLRow](graphResult.Data)
	if malformed > 0 {
		result.Warnings = append(result.Warnings, assessment.AssessmentWarning{
			Code:    "arcsql_malformed_arg_rows",
			Message: fmt.Sprintf("skipped %d malformed Arc SQL Resource Graph row(s)", malformed),
		})
	}

	for _, row := range rows {
		if filter != nil && filter.IsSubscriptionExcluded(row.SubscriptionID) {
			continue
		}
		if filter != nil && filter.IsServiceExcluded(row.SQLInstance) {
			continue
		}

		result.Records = append(result.Records, assessment.ArcSQLRecord{
			SubscriptionID:   row.SubscriptionID,
			SubscriptionName: subscriptions[row.SubscriptionID],
			Status:           row.Status,
			AzureArcServer:   row.AzureArcServer,
			SQLInstance:      row.SQLInstance,
			ResourceGroup:    row.ResourceGroup,
			Version:          row.Version,
			Build:            row.Build,
			PatchLevel:       row.PatchLevel,
			Edition:          row.Edition,
			VCores:           row.VCores,
			License:          row.License,
			DPSStatus:        row.DPSStatus,
			TELStatus:        row.TELStatus,
			DefenderStatus:   row.DefenderStatus,
		})
	}

	sort.Slice(result.Records, func(i, j int) bool {
		if result.Records[i].SubscriptionID != result.Records[j].SubscriptionID {
			return result.Records[i].SubscriptionID < result.Records[j].SubscriptionID
		}
		return result.Records[i].SQLInstance < result.Records[j].SQLInstance
	})
	return result, nil
}

type arcSQLRow struct {
	SubscriptionID string `json:"subscriptionId"`
	Status         string `json:"status"`
	AzureArcServer string `json:"AzureArcServer"`
	SQLInstance    string `json:"SQLInstance"`
	ResourceGroup  string `json:"resourceGroup"`
	Version        string `json:"version"`
	Build          string `json:"Build"`
	PatchLevel     string `json:"patchLevel"`
	Edition        string `json:"edition"`
	VCores         string `json:"vcores"`
	License        string `json:"License"`
	DPSStatus      string `json:"DPSStatus"`
	TELStatus      string `json:"TELStatus"`
	DefenderStatus string `json:"DefenderStatus"`
}
