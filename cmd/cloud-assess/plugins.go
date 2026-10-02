package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/plugins"
	"github.com/spf13/cobra"
)

type pluginInfo struct {
	assessment.PluginMetadata
	ScannerAvailable bool     `json:"scannerAvailable"`
	SourceFile       string   `json:"sourceFile,omitempty"`
	Recommendations  int      `json:"recommendations"`
	ResourceTypes    []string `json:"resourceTypes,omitempty"`
}

func registeredPluginInfo() ([]pluginInfo, error) {
	yamlPlugins, err := discoverYAMLPlugins()
	if err != nil {
		return nil, fmt.Errorf("YAML plugin preflight: %w", err)
	}
	rows := []pluginInfo{}
	for _, m := range plugins.InternalMetadata() {
		rows = append(rows, pluginInfo{PluginMetadata: m, ScannerAvailable: true})
	}
	for _, p := range yamlPlugins {
		types := map[string]bool{}
		for _, d := range p.Definitions {
			types[d.ResourceType] = true
		}
		row := pluginInfo{PluginMetadata: assessment.PluginMetadata{Name: p.Name, Version: p.Version, Description: p.Description, Author: p.Author, License: p.License, Type: "yaml"}, SourceFile: p.SourceFile, Recommendations: len(p.Definitions)}
		for typ := range types {
			row.ResourceTypes = append(row.ResourceTypes, typ)
		}
		sort.Strings(row.ResourceTypes)
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows, nil
}

// Escape controls, quotes and path backslashes while retaining printable Unicode.
func terminalText(value string) string { q := strconv.Quote(value); return q[1 : len(q)-1] }

func newPluginsCommand() *cobra.Command {
	root := &cobra.Command{Use: "plugins", Short: "Inspect implemented internal and discovered YAML plugins offline", Args: cobra.NoArgs}
	listJSON, infoJSON := false, false
	list := &cobra.Command{Use: "list", Short: "List available plugin metadata without Azure access", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		rows, err := registeredPluginInfo()
		if err != nil {
			return err
		}
		if listJSON {
			return writePluginJSON(cmd, rows)
		}
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "NAME\tVERSION\tTYPE\tDESCRIPTION"); err != nil {
			return err
		}
		for _, p := range rows {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", terminalText(p.Name), terminalText(p.Version), terminalText(p.Type), terminalText(p.Description)); err != nil {
				return err
			}
		}
		return w.Flush()
	}}
	info := &cobra.Command{Use: "info <plugin-name>", Short: "Inspect plugin metadata and capabilities without Azure access", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		rows, err := registeredPluginInfo()
		if err != nil {
			return err
		}
		for _, p := range rows {
			if p.Name != args[0] {
				continue
			}
			if infoJSON {
				return writePluginJSON(cmd, p)
			}
			out := fmt.Sprintf("Plugin: %s\nVersion: %s\nDescription: %s\nType: %s\n", terminalText(p.Name), terminalText(p.Version), terminalText(p.Description), terminalText(p.Type))
			if p.Author != "" {
				out += "Author: " + terminalText(p.Author) + "\n"
			}
			if p.License != "" {
				out += "License: " + terminalText(p.License) + "\n"
			}
			if p.SourceFile != "" {
				out += "Command Path: " + terminalText(p.SourceFile) + "\n"
			}
			if p.Type == "yaml" {
				out += fmt.Sprintf("Graph recommendations: %d\nResource Types:", p.Recommendations)
				for _, typ := range p.ResourceTypes {
					out += " " + terminalText(typ)
				}
				out += "\n"
			} else {
				out += "Internal scanner: available\n"
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), out)
			return err
		}
		return fmt.Errorf("internal or YAML plugin %q is unavailable", args[0])
	}}
	list.Flags().BoolVarP(&listJSON, "json", "j", false, "Print registry metadata as JSON")
	info.Flags().BoolVarP(&infoJSON, "json", "j", false, "Print plugin metadata as JSON")
	root.AddCommand(list, info)
	return root
}

func writePluginJSON(cmd *cobra.Command, value any) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
