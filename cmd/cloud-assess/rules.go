package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/diagnostics"
	"github.com/DeBoX85/Cloud-Assess/internal/rules"
	"github.com/DeBoX85/Cloud-Assess/internal/scanners"
	"github.com/spf13/cobra"
)

func newRulesCommand() *cobra.Command {
	jsonOutput := false
	command := &cobra.Command{
		Use: "rules", Short: "Print supported embedded and Diagnostics recommendations without Azure access",
		Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			catalog, _, err := rules.LoadPinnedCatalog()
			if err != nil {
				return err
			}
			var resourceTypes []string
			for _, service := range scanners.All() {
				resourceTypes = append(resourceTypes, service.ResourceTypes...)
			}
			rows := rules.Inspect(catalog, resourceTypes, diagnostics.Recommendations())
			return writeRules(command.OutOrStdout(), rows, catalog.ResourceTypeCount(), jsonOutput)
		},
	}
	command.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Print the rules list as a JSON array")
	return command
}

func writeRules(writer io.Writer, rows []rules.InspectionRow, resourceTypes int, jsonOutput bool) error {
	if jsonOutput {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "\t")
		return encoder.Encode(rows)
	}
	var output strings.Builder
	fmt.Fprintf(&output, "## Recommendations List\n\nTotal Supported Azure Resource Types: %d\n\n", resourceTypes)
	output.WriteString("|  | Id | Resource Type | Category | Impact | Recommendation | Learn |\n")
	output.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	for i, row := range rows {
		// Show the guidance URI as escaped text. This preserves arbitrary pinned
		// URLs without letting parentheses or HTML alter the Markdown structure.
		fmt.Fprintf(&output, "| %d | %s | %s | %s | %s | %s | %s |\n", i+1,
			markdownCell(row.RecommendationID), markdownCell(row.ResourceType),
			markdownCell(row.Category), markdownCell(row.Impact),
			markdownCell(row.Recommendation), markdownCell(row.LearnMoreURL))
	}
	_, err := io.WriteString(writer, output.String())
	return err
}

func markdownCell(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "|", "&#124;",
		"\r\n", "<br>", "\r", "<br>", "\n", "<br>").Replace(value)
}
