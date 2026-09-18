package main

import (
	"context"
	"fmt"

	"github.com/DeBoX85/Cloud-Assess/internal/app"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/config"
	"github.com/DeBoX85/Cloud-Assess/internal/gate"
	"github.com/DeBoX85/Cloud-Assess/internal/orchestration"
	"github.com/DeBoX85/Cloud-Assess/internal/stages"
	"github.com/spf13/cobra"
)

type exitCodeContextKey struct{}

type scanFlags struct {
	managementGroups      []string
	subscriptions         []string
	resourceGroups        []string
	stageNames            []string
	stageParams           []string
	outputName            string
	filtersFile           string
	failOn                string
	xlsx                  bool
	json                  bool
	csv                   bool
	stdout                bool
	sarif                 bool
	redactSubscriptionIDs bool
}

type scanExecutor func(context.Context, scanFlags) (int, error)

func newRootCommand(executor scanExecutor) *cobra.Command {
	brand := branding.Default()
	exitCode := 0
	ctx := context.WithValue(context.Background(), exitCodeContextKey{}, &exitCode)
	root := &cobra.Command{
		Use:           brand.CLIName,
		Short:         brand.ReportTitle,
		Long:          fmt.Sprintf("%s is an Azure cloud assessment toolkit.", brand.ProductName),
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetContext(ctx)
	root.AddCommand(newScanCommand(executor, &exitCode))
	return root
}

func newScanCommand(executor scanExecutor, exitCode *int) *cobra.Command {
	flags := scanFlags{}
	command := &cobra.Command{
		Use:   "scan",
		Short: "Assess Azure resources",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			code, err := executor(command.Context(), flags)
			*exitCode = code
			return err
		},
	}

	command.Flags().StringSliceVar(&flags.managementGroups, "management-group-id", nil, "Azure Management Group ID")
	command.Flags().StringSliceVarP(&flags.subscriptions, "subscription-id", "s", nil, "Azure Subscription ID")
	command.Flags().StringSliceVarP(&flags.resourceGroups, "resource-group", "g", nil, "Azure Resource Group (use with one --subscription-id)")
	command.Flags().StringSliceVar(&flags.stageNames, "stages", nil, "Control assessment stages (advisor, defender, defender-recommendations, arc, policy, cost, diagnostics)")
	command.Flags().StringArrayVar(&flags.stageParams, "stage-param", nil, "Stage option in the form stage.key=value (repeatable)")
	command.Flags().BoolVar(&flags.xlsx, "xlsx", true, "Create Excel report")
	command.Flags().BoolVar(&flags.json, "json", false, "Create canonical JSON report")
	command.Flags().BoolVar(&flags.csv, "csv", false, "Create CSV report files")
	command.Flags().BoolVar(&flags.stdout, "stdout", false, "Write canonical JSON to stdout")
	command.Flags().BoolVar(&flags.sarif, "sarif", false, "Create SARIF 2.1.0 report")
	command.Flags().StringVarP(&flags.outputName, "output-name", "o", "", "Output base filename without extension")
	command.Flags().BoolVar(
		&flags.redactSubscriptionIDs,
		"redact-subscription-ids",
		true,
		"Redact subscription IDs in XLSX, CSV, JSON and stdout output; SARIF retains stable resource identities",
	)
	command.Flags().StringVarP(&flags.filtersFile, "filters", "e", "", "Assessment filters file (YAML)")
	command.Flags().StringVar(&flags.failOn, "fail-on", "", "Return exit code 2 when findings meet or exceed this impact (High, Medium, Low)")
	return command
}

func executeScan(ctx context.Context, flags scanFlags) (int, error) {
	filters, err := config.LoadFilters(flags.filtersFile)
	if err != nil {
		return app.ExitExecutionFail, err
	}
	stageConfig := stages.NewDefault()
	if err := stageConfig.Apply(flags.stageNames); err != nil {
		return app.ExitExecutionFail, err
	}
	if err := stageConfig.ApplyParams(flags.stageParams); err != nil {
		return app.ExitExecutionFail, fmt.Errorf("apply stage parameters: %w", err)
	}
	if err := stageConfig.Validate(); err != nil {
		return app.ExitExecutionFail, err
	}
	if stageConfig.IsEnabled(stages.Plugin) {
		return app.ExitExecutionFail, fmt.Errorf("plugin stage is not available in the current core-v1 build")
	}
	if flags.failOn != "" {
		if _, err := gate.Parse(flags.failOn); err != nil {
			return app.ExitExecutionFail, err
		}
	}

	credential, err := azure.NewCredential()
	if err != nil {
		return app.ExitExecutionFail, err
	}
	operations, err := orchestration.NewAzureOperations(credential)
	if err != nil {
		return app.ExitExecutionFail, err
	}
	coordinator := orchestration.NewCoordinator(operations)
	runner := app.NewRunner(coordinator)
	outcome, runErr := runner.Run(ctx, app.ScanOptions{
		Assessment: orchestration.Request{
			ManagementGroups: flags.managementGroups,
			Subscriptions:    flags.subscriptions,
			ResourceGroups:   flags.resourceGroups,
			Filters:          filters,
			Stages:           stageConfig,
		},
		Outputs: app.OutputOptions{
			BaseName:              flags.outputName,
			XLSX:                  flags.xlsx,
			JSON:                  flags.json,
			CSV:                   flags.csv,
			Stdout:                flags.stdout,
			SARIF:                 flags.sarif,
			RedactSubscriptionIDs: flags.redactSubscriptionIDs,
			Version:               version,
		},
		FailOn: flags.failOn,
	})
	return outcome.ExitCode, runErr
}
