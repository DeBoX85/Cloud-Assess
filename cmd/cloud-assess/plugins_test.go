package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/branding"
)

func TestPluginRegistryPrecedenceSourcePathAndTerminalSafety(t *testing.T) {
	home := t.TempDir()
	current := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Chdir(current)
	root := filepath.Join(home, "."+branding.Default().CLIName, "plugins")
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("plugins", 0700); err != nil {
		t.Fatal(err)
	}
	yaml := "name: operator-checks\nversion: '1.0'\ndescription: \"safe\\n\\u001b[31m\"\nauthor: Operator\nlicense: MIT\nqueries:\n  - aprlGuid: custom\n    description: Check\n    recommendationResourceType: Microsoft.Storage/storageAccounts\n    query: resources\n"
	path := filepath.Join(root, "plugin.yml")
	if err := os.WriteFile(path, []byte(yaml), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("plugins", "duplicate.yml"), []byte(strings.Replace(yaml, "'1.0'", "'2.0'", 1)), 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"plugins", "list", "--json"}, {"plugins", "info", "operator-checks", "--json"}, {"plugins", "info", "operator-checks"}} {
		var out bytes.Buffer
		root := newRootCommand(func(context.Context, scanFlags) (int, error) { t.Fatal("registry executed scanner"); return 1, nil })
		root.SetOut(&out)
		root.SetErr(io.Discard)
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatal(err)
		}
		if args[len(args)-1] == "--json" {
			var rows []pluginInfo
			if args[1] == "list" {
				if err := json.Unmarshal(out.Bytes(), &rows); err != nil {
					t.Fatal(err)
				}
			} else {
				var row pluginInfo
				if err := json.Unmarshal(out.Bytes(), &row); err != nil {
					t.Fatal(err)
				}
				rows = []pluginInfo{row}
			}
			var selected *pluginInfo
			for i := range rows {
				if rows[i].Name == "operator-checks" {
					selected = &rows[i]
					break
				}
			}
			if selected == nil || selected.Version != "1.0" || selected.SourceFile != path || selected.Recommendations != 1 || len(selected.ResourceTypes) != 1 || selected.ScannerAvailable {
				t.Fatalf("wrong discovery/capability: %#v", rows)
			}
		} else if strings.Contains(out.String(), "\x1b") || !strings.Contains(out.String(), `safe\n\x1b[31m`) {
			t.Fatalf("unsafe terminal output %q", out.String())
		}
	}
	if err := os.WriteFile(filepath.Join("plugins", "broken.yml"), []byte("broken: ["), 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	rootCmd := newRootCommand(nil)
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"plugins", "list", "--json"})
	if err := rootCmd.Execute(); err == nil || out.Len() != 0 {
		t.Fatal("malformed discovery produced partial registry")
	}
}

func TestZoneCLIFlagPreflightAndNoUnusedDiscovery(t *testing.T) {
	for _, args := range [][]string{{"zone-mapping"}, {"scan", "vm", "--plugin", "zone-mapping"}} {
		root := newRootCommand(func(_ context.Context, f scanFlags) (int, error) {
			cfg, err := configureStages(f)
			if err != nil {
				t.Fatal(err)
			}
			if !cfg.IsEnabled("plugin") || f.pluginOnly != (args[0] == "zone-mapping") {
				t.Fatal("mode lost")
			}
			if args[0] == "scan" && len(f.scannerKeys) != 1 {
				t.Fatal("scanner selection lost")
			}
			return 0, nil
		})
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []scanFlags{{internalPlugins: []string{"unknown"}}, {internalPlugins: []string{"zone-mapping"}, stageNames: []string{"-plugin"}}, {pluginOnly: true, internalPlugins: []string{"zone-mapping"}, stageNames: []string{"graph"}}, {pluginOnly: true, internalPlugins: []string{"zone-mapping"}, stageParams: []string{"plugin.target-regions=westus"}}} {
		if _, err := configureStages(f); err == nil {
			t.Fatal("invalid preflight accepted")
		}
	}
	// Help must not read even malformed discovery candidates or authenticate.
	t.Chdir(t.TempDir())
	if err := os.Mkdir("plugins", 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("plugins/broken.yml", []byte("broken: ["), 0600); err != nil {
		t.Fatal(err)
	}
	root := newRootCommand(func(context.Context, scanFlags) (int, error) { t.Fatal("help invoked executor"); return 1, nil })
	root.SetOut(io.Discard)
	root.SetArgs([]string{"zone-mapping", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestZoneOnlySkipsUnusedYAMLBeforeCredentialPreflight(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("AZURE_TOKEN_CREDENTIALS", "invalid-synthetic-selection")
	if err := os.Mkdir("plugins", 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("plugins/broken.yml", []byte("broken: ["), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := executeScan(context.Background(), scanFlags{pluginOnly: true, internalPlugins: []string{"zone-mapping"}})
	if err == nil || !strings.Contains(err.Error(), "invalid AZURE_TOKEN_CREDENTIALS") {
		t.Fatalf("unused discovery reached zone-only path: %v", err)
	}
	_, err = executeScan(context.Background(), scanFlags{internalPlugins: []string{"zone-mapping"}})
	if err == nil || !strings.Contains(err.Error(), "YAML plugin preflight") {
		t.Fatalf("normal YAML preflight bypassed: %v", err)
	}
}
