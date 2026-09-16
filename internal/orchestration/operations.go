package orchestration

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/advisor"
	"github.com/DeBoX85/Cloud-Assess/internal/arcsql"
	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/cost"
	"github.com/DeBoX85/Cloud-Assess/internal/defender"
	"github.com/DeBoX85/Cloud-Assess/internal/diagnostics"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/policy"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
)

// Operations is the Azure-facing dependency set used by the coordinator. Keeping these
// functions explicit lets orchestration tests exercise the real control flow without Azure.
type Operations struct {
	DiscoverSubscriptions       func(context.Context, []string, *config.Filters) (map[string]string, error)
	DiscoverManagementGroups    func(context.Context, []string, *config.Filters) (map[string]string, error)
	DiscoverResources           func(context.Context, map[string]string, *config.Filters) (*discovery.ResourceInventory, error)
	LoadCatalog                 func() (*rules.Catalog, error)
	ExecuteGraph                func(context.Context, []assessment.RecommendationDefinition, map[string]string, *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error)
	ScanDiagnostics             func(context.Context, []assessment.Resource, *config.AssessmentFilter, map[string]string) (diagnostics.Result, error)
	ScanAdvisor                 func(context.Context, map[string]string, *config.AssessmentFilter) (advisor.Result, error)
	ScanDefenderStatus          func(context.Context, map[string]string, *config.AssessmentFilter) (defender.StatusResult, error)
	ScanDefenderRecommendations func(context.Context, map[string]string, *config.AssessmentFilter) (defender.RecommendationsResult, error)
	ScanPolicy                  func(context.Context, map[string]string, *config.AssessmentFilter) (policy.Result, error)
	ScanArcSQL                  func(context.Context, map[string]string, *config.AssessmentFilter) (arcsql.Result, error)
	ScanCost                    func(context.Context, map[string]string) (cost.Result, error)
}

// NewAzureOperations binds production Azure clients and scanners to the orchestration contract.
func NewAzureOperations(credential azcore.TokenCredential) (Operations, error) {
	scopeClient, err := discovery.NewAzureScopeClient(credential, azure.NewARMClientOptions())
	if err != nil {
		return Operations{}, err
	}
	graphClient := arg.NewClient(arg.NewHTTPTransport(credential))
	diagnosticsScanner := diagnostics.New(credential)
	advisorScanner := advisor.New(credential)
	defenderScanner := defender.New(credential)
	policyScanner := policy.New(credential)
	arcSQLScanner := arcsql.New(credential)
	costScanner := cost.New(credential)

	return Operations{
		DiscoverSubscriptions: func(ctx context.Context, requested []string, filters *config.Filters) (map[string]string, error) {
			return discovery.DiscoverSubscriptions(ctx, scopeClient, requested, filters)
		},
		DiscoverManagementGroups: func(ctx context.Context, groups []string, filters *config.Filters) (map[string]string, error) {
			return discovery.DiscoverManagementGroupSubscriptions(ctx, scopeClient, groups, filters)
		},
		DiscoverResources: func(ctx context.Context, subscriptions map[string]string, filters *config.Filters) (*discovery.ResourceInventory, error) {
			return discovery.DiscoverResources(ctx, graphClient, subscriptions, filters)
		},
		LoadCatalog: func() (*rules.Catalog, error) {
			catalog, _, err := rules.LoadPinnedCatalog()
			return catalog, err
		},
		ExecuteGraph: func(ctx context.Context, definitions []assessment.RecommendationDefinition, subscriptions map[string]string, filter *config.AssessmentFilter) ([]assessment.Finding, []arg.RuleWarning, error) {
			var isExcluded func(string) bool
			if filter != nil {
				isExcluded = filter.IsServiceExcluded
			}
			return arg.ExecuteRecommendations(ctx, graphClient, definitions, subscriptions, isExcluded, arg.DefaultRuleWorkers)
		},
		ScanDiagnostics: func(ctx context.Context, resources []assessment.Resource, filter *config.AssessmentFilter, subscriptions map[string]string) (diagnostics.Result, error) {
			return diagnosticsScanner.Scan(ctx, resources, filter, subscriptions)
		},
		ScanAdvisor: func(ctx context.Context, subscriptions map[string]string, filter *config.AssessmentFilter) (advisor.Result, error) {
			return advisorScanner.Scan(ctx, subscriptions, filter)
		},
		ScanDefenderStatus: func(ctx context.Context, subscriptions map[string]string, filter *config.AssessmentFilter) (defender.StatusResult, error) {
			return defenderScanner.ScanStatus(ctx, subscriptions, filter)
		},
		ScanDefenderRecommendations: func(ctx context.Context, subscriptions map[string]string, filter *config.AssessmentFilter) (defender.RecommendationsResult, error) {
			return defenderScanner.ScanRecommendations(ctx, subscriptions, filter)
		},
		ScanPolicy: func(ctx context.Context, subscriptions map[string]string, filter *config.AssessmentFilter) (policy.Result, error) {
			return policyScanner.Scan(ctx, subscriptions, filter)
		},
		ScanArcSQL: func(ctx context.Context, subscriptions map[string]string, filter *config.AssessmentFilter) (arcsql.Result, error) {
			return arcSQLScanner.Scan(ctx, subscriptions, filter)
		},
		ScanCost: func(ctx context.Context, subscriptions map[string]string) (cost.Result, error) {
			return costScanner.Scan(ctx, subscriptions)
		},
	}, nil
}
