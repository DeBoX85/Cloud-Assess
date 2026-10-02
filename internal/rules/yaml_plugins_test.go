package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const pluginFixture = `name: tenant-checks
description: Operator query library
author: Test Engineer
license: MIT
queries:
  - aprlGuid: sample-id
    recommendationTypeId: accepted-but-not-projected
    description: Enable the control
    recommendationControl: Security
    recommendationImpact: High
    recommendationResourceType: Microsoft.Storage/storageAccounts
    recommendationMetadataState: disabled
    longDescription: Details
    potentialBenefits: Protection
    pgVerified: true
    automationAvailable: false
    tags: [test]
    learnMoreLink:
      - name: Guide
        url: https://example.test/guide
    query: "resources | where false // under-development\n"
`

func writePluginTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	target := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestYAMLPluginLiteralSourceMapping(t *testing.T) {
	root := t.TempDir()
	writePluginTestFile(t, root, "plugin.yml", pluginFixture)
	plugins, err := DiscoverYAMLPlugins([]string{filepath.Join(root, "missing"), root})
	if err != nil {
		t.Fatal(err)
	}
	expected := []YAMLPlugin{{SourceFile: filepath.Join(root, "plugin.yml"), Name: "tenant-checks", Version: "1.0.0", Description: "Operator query library", Author: "Test Engineer", License: "MIT", Definitions: []assessment.RecommendationDefinition{{
		ID: "sample-id", Recommendation: "Enable the control", Category: "Security", Impact: "High", ResourceType: "Microsoft.Storage/storageAccounts", State: "disabled", LongDescription: "Details", PotentialBenefits: "Protection", PGVerified: true, AutomationAvailable: "false", Tags: []string{"test"}, LearnMore: []assessment.LearnMoreLink{{Name: "Guide", URL: "https://example.test/guide"}}, Query: "resources | where false // under-development\n", Source: "tenant-checks", ValidationMechanism: ValidationARG,
	}}}}
	if !reflect.DeepEqual(plugins, expected) {
		t.Fatalf("source mapping mismatch:\ngot %+v\nwant %+v", plugins, expected)
	}
	definitions := PluginDefinitions(plugins)
	definitions[0].Tags[0] = "changed"
	definitions[0].LearnMore[0].URL = "changed"
	if plugins[0].Definitions[0].Tags[0] != "test" || plugins[0].Definitions[0].LearnMore[0].URL != "https://example.test/guide" {
		t.Fatal("definition slices alias plugin input")
	}
}

func TestYAMLQueryFileOverridesInlineAndRetainsBytes(t *testing.T) {
	root := t.TempDir()
	writePluginTestFile(t, root, "nested/plugin.yaml", pluginFixture+"    queryFile: kql/query.kql\n")
	query := "resources\r\n| where name == 'abc'\r\n"
	writePluginTestFile(t, root, "nested/kql/query.kql", query)
	plugins, err := DiscoverYAMLPlugins([]string{root})
	if err != nil {
		t.Fatal(err)
	}
	if plugins[0].Definitions[0].Query != query {
		t.Fatal("queryFile did not preserve/override query text")
	}
}

func TestYAMLDiscoveryAndOverridePrecedence(t *testing.T) {
	home, current := t.TempDir(), t.TempDir()
	writePluginTestFile(t, home, "z.yaml", strings.ReplaceAll(pluginFixture, "Enable the control", "home"))
	writePluginTestFile(t, current, "a.yaml", strings.ReplaceAll(pluginFixture, "Enable the control", "current"))
	writePluginTestFile(t, current, "z.yml", strings.ReplaceAll(strings.ReplaceAll(pluginFixture, "tenant-checks", "z-last"), "Enable the control", "last sorted name"))
	writePluginTestFile(t, current, "ignored.YAML", "not a plugin")
	plugins, err := DiscoverYAMLPlugins([]string{home, current})
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 2 || plugins[0].Name != "tenant-checks" || plugins[0].Definitions[0].Recommendation != "home" {
		t.Fatalf("first-name precedence: %+v", plugins)
	}
	base := NewCatalog()
	base.Add(assessment.RecommendationDefinition{ID: "sample-id", ResourceType: "Microsoft.Storage/storageAccounts", Recommendation: "embedded", Tags: []string{"original"}})
	overlaid := WithPluginDefinitions(base, PluginDefinitions(plugins))
	got := overlaid.ByResourceType("Microsoft.Storage/storageAccounts")
	if len(got) != 1 || got[0].Recommendation != "last sorted name" || !overlaid.IsPlugin(got[0]) {
		t.Fatalf("sorted override mismatch: %+v", got)
	}
	if base.ByResourceType("Microsoft.Storage/storageAccounts")[0].Recommendation != "embedded" {
		t.Fatal("base catalog mutated")
	}
	plugins[1].Definitions[0].Tags[0] = "changed"
	if got[0].Tags[0] != "test" {
		t.Fatal("catalog retained mutable input")
	}
	overlaid.Add(assessment.RecommendationDefinition{ID: "sample-id", ResourceType: "Microsoft.Storage/storageAccounts"})
	if overlaid.IsPlugin(overlaid.ByResourceType("Microsoft.Storage/storageAccounts")[0]) {
		t.Fatal("ordinary add retained stale origin")
	}
	// Duplicate same-type query IDs within one file are source last-query-wins.
	writePluginTestFile(t, home, "z.yaml", pluginFixture+"  - aprlGuid: sample-id\n    description: later query\n    recommendationResourceType: Microsoft.Storage/storageAccounts\n    query: resources\n")
	plugins, err = DiscoverYAMLPlugins([]string{home})
	if err != nil {
		t.Fatal(err)
	}
	got = WithPluginDefinitions(nil, PluginDefinitions(plugins)).ByResourceType("Microsoft.Storage/storageAccounts")
	if len(got) != 1 || got[0].Recommendation != "later query" {
		t.Fatal("within-file override mismatch")
	}
}

func TestYAMLInvalidInputsAndPortablePaths(t *testing.T) {
	tests := map[string]string{
		"unknown field":          pluginFixture + "unexpected: secret-payload\n",
		"duplicate field":        pluginFixture + "name: duplicate\n",
		"duplicate nested field": pluginFixture + "    query: duplicate\n",
		"extra document":         pluginFixture + "---\nname: second\n",
		"alias":                  strings.Replace(pluginFixture, "description: Enable the control", "description: &a Enable the control", 1) + "    longDescription: *a\n",
		"merge":                  pluginFixture + "    <<: {description: other}\n",
		"invalid bool":           strings.Replace(pluginFixture, "pgVerified: true", "pgVerified: not-a-bool-secret", 1),
		"blank ID":               strings.Replace(pluginFixture, "aprlGuid: sample-id", "aprlGuid: '   '", 1),
		"blank description":      strings.Replace(pluginFixture, "description: Enable the control", "description: '   '", 1),
		"blank type":             strings.Replace(pluginFixture, "recommendationResourceType: Microsoft.Storage/storageAccounts", "recommendationResourceType: ''", 1),
		"blank query":            strings.Replace(pluginFixture, `query: "resources | where false // under-development\n"`, `query: " "`, 1),
		"no queries":             "name: example\nqueries: []\n",
		"invalid name":           strings.Replace(pluginFixture, "name: tenant-checks", "name: '../secret'", 1),
		"unknown nested link":    pluginFixture + "    learnMoreLink: [{url: https://example.test, extra: secret}]\n",
	}
	for _, name := range []string{"APRL", "aor", "CUSTOM", "DIAGNOSTICS", "zone-mapping", "service-health", "sql-eol", "carbon-emissions", "ai-gov", "region-selection"} {
		tests["reserved "+name] = strings.Replace(pluginFixture, "name: tenant-checks", "name: "+name, 1)
	}
	for _, value := range []string{"../outside.kql", "/absolute.kql", `C:\outside.kql`, `\\host\share\x.kql`, "queries/../../outside.kql", "NUL", "nested/COM1.txt", "x:stream", "queries\\x.kql", "queries/../x.kql", "queries/x. ", "query.", "NUL .txt", "COM¹.txt", "LPT²", "que?ry.kql"} {
		tests["query path "+value] = pluginFixture + fmt.Sprintf("    queryFile: '%s'\n", value)
	}
	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writePluginTestFile(t, root, "plugin.yaml", content)
			if _, err := DiscoverYAMLPlugins([]string{root}); err == nil {
				t.Fatal("invalid candidate was accepted")
			} else if strings.Contains(err.Error(), "secret-payload") || strings.Contains(err.Error(), "not-a-bool-secret") {
				t.Fatal("error echoed YAML input")
			}
		})
	}
	root := t.TempDir()
	writePluginTestFile(t, root, "plugin.yaml", pluginFixture+"    queryFile: missing.kql\n")
	if _, err := DiscoverYAMLPlugins([]string{root}); err == nil {
		t.Fatal("missing queryFile used inline fallback")
	}
	if _, err := DiscoverYAMLPlugins([]string{filepath.Join(root, "plugin.yaml")}); err == nil {
		t.Fatal("non-directory root accepted")
	}
}

func TestYAMLInputBounds(t *testing.T) {
	root := t.TempDir()
	writePluginTestFile(t, root, "plugin.yaml", pluginFixture)
	handle, err := os.OpenRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	for _, tc := range []struct {
		name   string
		budget pluginBudget
		want   string
	}{
		{"bytes", pluginBudget{bytes: MaxPluginTotalBytes}, "total bytes"},
		{"queries", pluginBudget{queries: MaxPluginQueries}, "queries exceed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loadYAMLPlugin(handle, "plugin.yaml", &tc.budget)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("limit: %v", err)
			}
		})
	}
	budget := pluginBudget{entries: MaxPluginEntries}
	if err := walkPluginDirectory(handle, ".", 0, &budget, func(string) error { return nil }); err == nil {
		t.Fatal("entry bound not enforced")
	}
	if err := walkPluginDirectory(handle, ".", MaxPluginDepth+1, &pluginBudget{}, func(string) error { return nil }); err == nil {
		t.Fatal("depth bound not enforced")
	}
	writePluginTestFile(t, root, "large.kql", strings.Repeat("x", 1048576))
	if _, err := readPluginFile(handle, "large.kql", &pluginBudget{}); err != nil {
		t.Fatal(err)
	}
	writePluginTestFile(t, root, "large.kql", strings.Repeat("x", 1048577))
	if _, err := readPluginFile(handle, "large.kql", &pluginBudget{}); err == nil {
		t.Fatal("file byte bound not enforced")
	}
	for index := 0; index < MaxPluginFiles; index++ {
		writePluginTestFile(t, root, fmt.Sprintf("%03d.yaml", index), pluginFixture)
	}
	if _, err := DiscoverYAMLPlugins([]string{root}); err == nil || !strings.Contains(err.Error(), "file count exceeds") {
		t.Fatalf("candidate file bound: %v", err)
	}
}

func TestYAMLQuerySymlinkAndIntermediateEscape(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	writePluginTestFile(t, outside, "query.kql", "resources")
	if err := os.Symlink(filepath.Join(outside, "query.kql"), filepath.Join(root, "query.kql")); err != nil {
		t.Fatal(err)
	}
	writePluginTestFile(t, root, "plugin.yaml", pluginFixture+"    queryFile: query.kql\n")
	if _, err := DiscoverYAMLPlugins([]string{root}); err == nil {
		t.Fatal("final query symlink accepted")
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	writePluginTestFile(t, root, "plugin.yaml", pluginFixture+"    queryFile: escape/query.kql\n")
	if _, err := DiscoverYAMLPlugins([]string{root}); err == nil {
		t.Fatal("intermediate symlink escape accepted")
	}
}

func TestYAMLConversionMatchesPinnedLoaderCapture(t *testing.T) {
	root := t.TempDir()
	writePluginTestFile(t, root, "plugin.yaml", pluginFixture)
	plugins, err := DiscoverYAMLPlugins([]string{root})
	if err != nil {
		t.Fatal(err)
	}
	capture, err := os.ReadFile("testdata/yaml-source-conversion.json")
	if err != nil {
		t.Fatal(err)
	}
	var expected map[string]any
	if err := json.Unmarshal(capture, &expected); err != nil {
		t.Fatal(err)
	}
	plugin := plugins[0]
	var rows []map[string]any
	for _, d := range plugin.Definitions {
		links := make([]map[string]string, 0, len(d.LearnMore))
		for _, link := range d.LearnMore {
			links = append(links, map[string]string{"Name": link.Name, "Url": link.URL})
		}
		rows = append(rows, map[string]any{"RecommendationID": d.ID, "Recommendation": d.Recommendation, "Category": d.Category, "Impact": d.Impact, "ResourceType": d.ResourceType, "MetadataState": d.State, "LongDescription": d.LongDescription, "PotentialBenefits": d.PotentialBenefits, "PgVerified": d.PGVerified, "AutomationAvailable": d.AutomationAvailable, "Tags": d.Tags, "GraphQuery": d.Query, "LearnMoreLink": links, "Source": d.Source})
	}
	actual := map[string]any{"Metadata": map[string]any{"Name": plugin.Name, "Version": plugin.Version, "Description": plugin.Description, "Author": plugin.Author, "License": plugin.License, "Type": 0, "CommandPath": "<fixture-path>", "ColumnMetadata": nil}, "Recommendations": rows}
	// Normalize only JSON number representation, not source values or field names.
	encoded, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}
	var actualJSON map[string]any
	if err := json.Unmarshal(encoded, &actualJSON); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actualJSON, expected) {
		t.Fatalf("pinned loader capture mismatch:\n%v\n%v", actualJSON, expected)
	}
}

func TestPortableQueryPathsAgreeWithDocumentedPolicy(t *testing.T) {
	for _, value := range []string{"query.kql", "nested/query file.kql", "unicode/æ-query.kql", "COM12.kql"} {
		if !portableQueryPath(value) {
			t.Fatalf("valid portable path rejected: %q", value)
		}
	}
	for _, value := range []string{"COM¹.kql", "NUL .txt", "../x.kql", "/x.kql", "query?x.kql", "query\x00.kql", "query\n.kql", strings.Repeat("x", 4097), strings.Repeat("a/", 16) + "query.kql"} {
		if portableQueryPath(value) {
			t.Fatal("invalid portable path accepted")
		}
	}
}

func TestPluginOverlayPreservesSameIDInOtherType(t *testing.T) {
	base := NewCatalog()
	base.Add(assessment.RecommendationDefinition{ID: "same", ResourceType: "Microsoft.Compute/virtualMachines", Source: SourceAPRL})
	base.Add(assessment.RecommendationDefinition{ID: "same", ResourceType: "Microsoft.Storage/storageAccounts", Source: SourceAPRL})
	input := []assessment.RecommendationDefinition{{ID: "same", ResourceType: "microsoft.storage/storageAccounts", Source: "operator", Tags: []string{"owned"}}}
	overlaid := WithPluginDefinitions(base, input)
	vm := overlaid.ByResourceType("Microsoft.Compute/virtualMachines")[0]
	storage := overlaid.ByResourceType("Microsoft.Storage/storageAccounts")[0]
	if overlaid.IsPlugin(vm) || !overlaid.IsPlugin(storage) || vm.Source != SourceAPRL || storage.Source != "operator" {
		t.Fatal("overlay changed same ID in other type")
	}
	copied := CopyDefinitions(input)
	input[0].Tags[0] = "changed"
	if copied[0].Tags[0] != "owned" || overlaid.ByResourceType("Microsoft.Storage/storageAccounts")[0].Tags[0] != "owned" {
		t.Fatal("mutable caller slice retained")
	}
}

func TestYAMLPluginKeepsUnicodeAndSpaceLabels(t *testing.T) {
	root := t.TempDir()
	writePluginTestFile(t, root, "plugin.yaml", strings.Replace(pluginFixture, "name: tenant-checks", "name: Kundens ærlige kontroller", 1))
	plugins, err := DiscoverYAMLPlugins([]string{root})
	if err != nil {
		t.Fatal(err)
	}
	if plugins[0].Name != "Kundens ærlige kontroller" || plugins[0].Definitions[0].Source != "Kundens ærlige kontroller" {
		t.Fatal("configured label changed")
	}
}
