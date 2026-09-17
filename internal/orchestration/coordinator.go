package orchestration

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/arg"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/discovery"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/scanners"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
)

const (
	StageScopeDiscovery    = "scope"
	StageResourceInventory = "inventory"
)

// Request contains assessment-engine inputs. Rendering and process-exit behavior are
// intentionally handled by the application layer rather than the coordinator.
type Request struct {
	ManagementGroups []string
	Subscriptions    []string
	ResourceGroups   []string
	ScannerKeys      []string
	Filters          *config.Filters
	Stages           *stages.Config
}

// Coordinator composes the independently characterized assessment subsystems into one
// deterministic assessment result.
type Coordinator struct {
	operations Operations
	runner     *stages.Runner
	now        func() time.Time
}

func NewCoordinator(operations Operations) *Coordinator {
	return &Coordinator{
		operations: operations,
		runner:     stages.NewRunner(),
		now:        time.Now,
	}
}

type scanState struct {
	subscriptions           map[string]string
	scopeID                 string
	inventory               *discovery.ResourceInventory
	resourceTypes           []assessment.ResourceTypeCount
	recommendations         []assessment.RecommendationDefinition
	findings                []assessment.Finding
	advisor                 []assessment.AdvisorRecommendation
	defender                []assessment.DefenderPlanStatus
	defenderRecommendations []assessment.DefenderRecommendation
	policy                  []assessment.PolicyNonCompliance
	arcSQL                  []assessment.ArcSQLRecord
	costs                   []assessment.CostRecord
}

// Run executes one assessment. Configuration/scope validation errors are returned before
// execution. Critical execution failures return both the partial canonical result and an error
// so the application layer can decide whether to persist status artifacts.
func (c *Coordinator) Run(ctx context.Context, request Request) (*result.AssessmentResult, error) {
	if c == nil {
		return nil, fmt.Errorf("assessment coordinator is nil")
	}
	prepared, err := prepareRequest(request)
	if err != nil {
		return nil, err
	}
	if err := validateOperations(c.operations); err != nil {
		return nil, err
	}
	if c.runner == nil {
		c.runner = stages.NewRunner()
	}
	if c.now == nil {
		c.now = time.Now
	}

	state := &scanState{}
	tasks := c.tasks(prepared, state)
	run := c.runner.Execute(ctx, tasks)

	inventory := state.inventory
	if inventory == nil {
		inventory = &discovery.ResourceInventory{}
	}
	assessmentResult := result.Build(result.Input{
		GeneratedAt:             c.now().UTC(),
		ScopeID:                 state.scopeID,
		Completeness:            run.Completeness,
		Stages:                  run.Stages,
		Recommendations:         state.recommendations,
		Findings:                state.findings,
		Resources:               inventory.Included,
		OutOfScope:              inventory.Excluded,
		ResourceTypes:           state.resourceTypes,
		Advisor:                 state.advisor,
		Defender:                state.defender,
		DefenderRecommendations: state.defenderRecommendations,
		AzurePolicy:             state.policy,
		ArcSQL:                  state.arcSQL,
		Costs:                   state.costs,
	})

	if run.Completeness == assessment.CompletenessFailed {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return assessmentResult, ctxErr
		}
		for _, execution := range run.Stages {
			if execution.Status == assessment.StageFailed && execution.Error != nil {
				return assessmentResult, fmt.Errorf("assessment failed during %s: %s", execution.Name, execution.Error.Message)
			}
		}
		return assessmentResult, fmt.Errorf("assessment failed")
	}
	return assessmentResult, nil
}

type preparedRequest struct {
	managementGroups []string
	subscriptions    []string
	selectedKeys     []string
	filters          *config.Filters
	stages           *stages.Config
}

func prepareRequest(request Request) (preparedRequest, error) {
	if len(request.ManagementGroups) > 0 && (len(request.Subscriptions) > 0 || len(request.ResourceGroups) > 0) {
		return preparedRequest{}, fmt.Errorf("management group name cannot be used with a subscription ID or resource group name")
	}
	if len(request.Subscriptions) < 1 && len(request.ResourceGroups) > 0 {
		return preparedRequest{}, fmt.Errorf("resource group name can only be used with a subscription ID")
	}
	if len(request.Subscriptions) > 1 && len(request.ResourceGroups) > 0 {
		return preparedRequest{}, fmt.Errorf("resource group name can only be used with one subscription ID")
	}

	filters := cloneFilters(request.Filters)
	filters.Assessment.Include.Subscriptions = append(filters.Assessment.Include.Subscriptions, request.Subscriptions...)
	if len(request.ResourceGroups) > 0 {
		for _, resourceGroup := range request.ResourceGroups {
			filters.Assessment.Include.ResourceGroups = append(filters.Assessment.Include.ResourceGroups,
				fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", request.Subscriptions[0], resourceGroup))
		}
	}
	filters.RebuildIndexes()
	if err := filters.Assessment.Validate(); err != nil {
		return preparedRequest{}, err
	}

	scannerKeys := append([]string(nil), request.ScannerKeys...)
	if len(scannerKeys) == 0 {
		scannerKeys = scanners.Keys()
	}
	selectedKeys := scanners.SelectedKeys(scannerKeys, filters.Assessment.Include.ResourceTypes)
	filters.Assessment.SetAllowedResourceTypes(scanners.ResourceTypes(selectedKeys))

	stageConfig := request.Stages
	if stageConfig == nil {
		stageConfig = stages.NewDefault()
	}
	if err := stageConfig.Validate(); err != nil {
		return preparedRequest{}, err
	}

	return preparedRequest{
		managementGroups: append([]string(nil), request.ManagementGroups...),
		subscriptions:    append([]string(nil), request.Subscriptions...),
		selectedKeys:     selectedKeys,
		filters:          filters,
		stages:           stageConfig,
	}, nil
}

func (c *Coordinator) tasks(request preparedRequest, state *scanState) []stages.Task {
	return []stages.Task{
		{
			Name: StageScopeDiscovery, Enabled: true, Critical: true,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				var subscriptions map[string]string
				var err error
				if len(request.managementGroups) > 0 {
					subscriptions, err = c.operations.DiscoverManagementGroups(ctx, request.managementGroups, request.filters)
				} else {
					subscriptions, err = c.operations.DiscoverSubscriptions(ctx, request.subscriptions, request.filters)
				}
				if err != nil {
					return stages.Outcome{}, err
				}
				state.subscriptions = subscriptions
				state.scopeID = discovery.ScopeID(subscriptions)
				return stages.Outcome{Records: len(subscriptions)}, nil
			},
		},
		{
			Name: StageResourceInventory, Enabled: true, Critical: true,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				inventory, err := c.operations.DiscoverResources(ctx, state.subscriptions, request.filters)
				if err != nil {
					return stages.Outcome{}, err
				}
				state.inventory = inventory
				state.resourceTypes = discovery.CountResourcesByTypeAndSubscription(inventory.Included, state.subscriptions)
				outcome := stages.Outcome{Records: len(inventory.Included)}
				if inventory.MalformedRows > 0 {
					outcome.Warnings = append(outcome.Warnings, assessment.AssessmentWarning{
						Code:    "inventory_malformed_arg_rows",
						Message: fmt.Sprintf("skipped %d malformed resource inventory row(s)", inventory.MalformedRows),
					})
				}
				return outcome, nil
			},
		},
		{
			Name: stages.Graph, Enabled: request.stages.IsEnabled(stages.Graph), Critical: true,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				catalog, err := c.operations.LoadCatalog()
				if err != nil {
					return stages.Outcome{}, err
				}
				phaseOneTypes := scanners.ResourceTypes(request.selectedKeys)
				state.recommendations = executableDefinitions(catalog, phaseOneTypes, request.filters.Assessment)

				phaseTwoTypes := prunedResourceTypes(request.selectedKeys, state.inventory.Included)
				definitions := executableDefinitions(catalog, phaseTwoTypes, request.filters.Assessment)
				findings, ruleWarnings, err := c.operations.ExecuteGraph(ctx, definitions, state.subscriptions, request.filters.Assessment)
				warnings := assessmentWarnings(ruleWarnings)
				if err != nil {
					return stages.Outcome{Warnings: warnings}, err
				}
				state.findings = append(state.findings, findings...)
				return stages.Outcome{Records: len(findings), Warnings: warnings}, nil
			},
		},
		{
			Name: stages.Diagnostics, Enabled: request.stages.IsEnabled(stages.Diagnostics), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanDiagnostics(ctx, state.inventory.Included, request.filters.Assessment, state.subscriptions)
				state.recommendations = append(state.recommendations, value.Recommendations...)
				state.findings = append(state.findings, value.Findings...)
				return stages.Outcome{Records: len(value.Findings), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.Advisor, Enabled: request.stages.IsEnabled(stages.Advisor), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanAdvisor(ctx, state.subscriptions, request.filters.Assessment)
				state.advisor = append(state.advisor, value.Records...)
				return stages.Outcome{Records: len(value.Records), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.Defender, Enabled: request.stages.IsEnabled(stages.Defender), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanDefenderStatus(ctx, state.subscriptions, request.filters.Assessment)
				state.defender = append(state.defender, value.Records...)
				return stages.Outcome{Records: len(value.Records), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.DefenderRecommendations, Enabled: request.stages.IsEnabled(stages.DefenderRecommendations), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanDefenderRecommendations(ctx, state.subscriptions, request.filters.Assessment)
				state.defenderRecommendations = append(state.defenderRecommendations, value.Records...)
				return stages.Outcome{Records: len(value.Records), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.Policy, Enabled: request.stages.IsEnabled(stages.Policy), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanPolicy(ctx, state.subscriptions, request.filters.Assessment)
				state.policy = append(state.policy, value.Records...)
				return stages.Outcome{Records: len(value.Records), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.Arc, Enabled: request.stages.IsEnabled(stages.Arc), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanArcSQL(ctx, state.subscriptions, request.filters.Assessment)
				state.arcSQL = append(state.arcSQL, value.Records...)
				return stages.Outcome{Records: len(value.Records), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.Cost, Enabled: request.stages.IsEnabled(stages.Cost), Critical: false,
			Run: func(ctx context.Context) (stages.Outcome, error) {
				value, err := c.operations.ScanCost(ctx, state.subscriptions)
				state.costs = append(state.costs, value.Records...)
				return stages.Outcome{Records: len(value.Records), Warnings: value.Warnings}, err
			},
		},
		{
			Name: stages.Plugin, Enabled: request.stages.IsEnabled(stages.Plugin), Critical: false,
			Run: nil,
		},
	}
}

func executableDefinitions(catalog *rules.Catalog, resourceTypes []string, filter *config.AssessmentFilter) []assessment.RecommendationDefinition {
	if catalog == nil {
		return nil
	}
	seenTypes := map[string]struct{}{}
	definitions := make([]assessment.RecommendationDefinition, 0)
	for _, resourceType := range resourceTypes {
		normalized := normalize(resourceType)
		if _, exists := seenTypes[normalized]; exists {
			continue
		}
		seenTypes[normalized] = struct{}{}
		for _, definition := range catalog.ByResourceType(resourceType) {
			var excluded func(string) bool
			if filter != nil {
				excluded = filter.IsRecommendationExcluded
			}
			if rules.IsExecutableEmbedded(definition, excluded) {
				definitions = append(definitions, definition)
			}
		}
	}
	sort.Slice(definitions, func(i, j int) bool {
		if normalize(definitions[i].ResourceType) != normalize(definitions[j].ResourceType) {
			return normalize(definitions[i].ResourceType) < normalize(definitions[j].ResourceType)
		}
		return definitions[i].ID < definitions[j].ID
	})
	return definitions
}

func prunedResourceTypes(selectedKeys []string, resources []assessment.Resource) []string {
	deployed := make(map[string]struct{}, len(resources))
	for _, resource := range resources {
		deployed[normalize(resource.Type)] = struct{}{}
	}

	seen := map[string]struct{}{}
	out := make([]string, 0)
	add := func(resourceType string) {
		normalized := normalize(resourceType)
		if _, exists := seen[normalized]; exists {
			return
		}
		seen[normalized] = struct{}{}
		out = append(out, resourceType)
	}

	for _, key := range selectedKeys {
		for _, service := range scanners.ByKey(key) {
			active := false
			for _, resourceType := range service.ResourceTypes {
				if _, exists := deployed[normalize(resourceType)]; exists {
					active = true
					break
				}
			}
			if !active {
				continue
			}
			for _, resourceType := range service.ResourceTypes {
				add(resourceType)
			}
		}
	}
	for _, service := range scanners.ByKey("resource") {
		for _, resourceType := range service.ResourceTypes {
			add(resourceType)
		}
	}
	return out
}

func assessmentWarnings(values []arg.RuleWarning) []assessment.AssessmentWarning {
	out := make([]assessment.AssessmentWarning, 0, len(values))
	for _, warning := range values {
		message := warning.Message
		if warning.RecommendationID != "" {
			message = fmt.Sprintf("recommendation %s (%s): %s", warning.RecommendationID, warning.ResourceType, warning.Message)
		}
		out = append(out, assessment.AssessmentWarning{Code: warning.Code, Message: message})
	}
	return out
}

func validateOperations(operations Operations) error {
	missing := make([]string, 0)
	if operations.DiscoverSubscriptions == nil {
		missing = append(missing, "DiscoverSubscriptions")
	}
	if operations.DiscoverManagementGroups == nil {
		missing = append(missing, "DiscoverManagementGroups")
	}
	if operations.DiscoverResources == nil {
		missing = append(missing, "DiscoverResources")
	}
	if operations.LoadCatalog == nil {
		missing = append(missing, "LoadCatalog")
	}
	if operations.ExecuteGraph == nil {
		missing = append(missing, "ExecuteGraph")
	}
	if operations.ScanDiagnostics == nil {
		missing = append(missing, "ScanDiagnostics")
	}
	if operations.ScanAdvisor == nil {
		missing = append(missing, "ScanAdvisor")
	}
	if operations.ScanDefenderStatus == nil {
		missing = append(missing, "ScanDefenderStatus")
	}
	if operations.ScanDefenderRecommendations == nil {
		missing = append(missing, "ScanDefenderRecommendations")
	}
	if operations.ScanPolicy == nil {
		missing = append(missing, "ScanPolicy")
	}
	if operations.ScanArcSQL == nil {
		missing = append(missing, "ScanArcSQL")
	}
	if operations.ScanCost == nil {
		missing = append(missing, "ScanCost")
	}
	if len(missing) > 0 {
		return fmt.Errorf("orchestration operations are not configured: %s", strings.Join(missing, ", "))
	}
	return nil
}

func cloneFilters(source *config.Filters) *config.Filters {
	if source == nil || source.Assessment == nil {
		return config.NewFilters()
	}
	result := config.NewFilters()
	if source.Assessment.Include != nil {
		result.Assessment.Include.Subscriptions = append([]string(nil), source.Assessment.Include.Subscriptions...)
		result.Assessment.Include.ResourceGroups = append([]string(nil), source.Assessment.Include.ResourceGroups...)
		result.Assessment.Include.ResourceTypes = append([]string(nil), source.Assessment.Include.ResourceTypes...)
		result.Assessment.Include.Tags = cloneStringMap(source.Assessment.Include.Tags)
	}
	if source.Assessment.Exclude != nil {
		result.Assessment.Exclude.Subscriptions = append([]string(nil), source.Assessment.Exclude.Subscriptions...)
		result.Assessment.Exclude.ResourceGroups = append([]string(nil), source.Assessment.Exclude.ResourceGroups...)
		result.Assessment.Exclude.Resources = append([]string(nil), source.Assessment.Exclude.Resources...)
		result.Assessment.Exclude.Recommendations = append([]string(nil), source.Assessment.Exclude.Recommendations...)
		result.Assessment.Exclude.Tags = cloneStringMap(source.Assessment.Exclude.Tags)
	}
	result.RebuildIndexes()
	return result
}

func cloneStringMap(source map[string]string) map[string]string {
	out := make(map[string]string, len(source))
	for key, value := range source {
		out[key] = value
	}
	return out
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
