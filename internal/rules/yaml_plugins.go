// Source schema/selection characterized from Microsoft AZQR (MIT licensed).
// See NOTICE.md and docs/YAML_GRAPH_PLUGINS.md.
package rules

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"gopkg.in/yaml.v3"
)

const (
	MaxPluginFileBytes  = 1 << 20
	MaxPluginTotalBytes = 16 << 20
	MaxPluginFiles      = 256
	MaxPluginQueries    = 4096
	MaxPluginEntries    = 8192
	MaxPluginDepth      = 16
)

// YAMLPlugin is local configuration, not an executable internal table plugin.
type YAMLPlugin struct {
	SourceFile  string
	Name        string
	Version     string
	Description string
	Author      string
	License     string
	Definitions []assessment.RecommendationDefinition
}

type yamlPluginConfig struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	License     string            `yaml:"license"`
	Queries     []yamlPluginQuery `yaml:"queries"`
}

type yamlPluginQuery struct {
	Description          string         `yaml:"description"`
	ID                   string         `yaml:"aprlGuid"`
	RecommendationTypeID *string        `yaml:"recommendationTypeId"`
	Category             string         `yaml:"recommendationControl"`
	Impact               string         `yaml:"recommendationImpact"`
	ResourceType         string         `yaml:"recommendationResourceType"`
	State                string         `yaml:"recommendationMetadataState"`
	LongDescription      string         `yaml:"longDescription"`
	PotentialBenefits    string         `yaml:"potentialBenefits"`
	PGVerified           bool           `yaml:"pgVerified"`
	AutomationAvailable  bool           `yaml:"automationAvailable"`
	Tags                 []string       `yaml:"tags"`
	LearnMore            []rawLearnMore `yaml:"learnMoreLink"`
	Query                string         `yaml:"query"`
	QueryFile            string         `yaml:"queryFile"`
}

func reservedPluginName(name string) bool {
	switch strings.ToLower(name) {
	case "aprl", "aor", "custom", "diagnostics", "zone-mapping", "service-health", "sql-eol", "carbon-emissions", "ai-gov", "region-selection":
		return true
	}
	return false
}

type pluginBudget struct{ bytes, files, queries, entries int }

// DiscoverYAMLPlugins preserves first-name discovery precedence, then sorts names
// for Graph override registration. Missing roots are skipped; candidate errors fail.
func DiscoverYAMLPlugins(directories []string) ([]YAMLPlugin, error) {
	var plugins []YAMLPlugin
	seen := map[string]bool{}
	budget := &pluginBudget{}
	for _, directory := range directories {
		root, err := os.OpenRoot(directory)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("open plugin directory: %w", err)
		}
		err = walkPluginDirectory(root, ".", 0, budget, func(relative string) error {
			budget.files++
			if budget.files > MaxPluginFiles {
				return fmt.Errorf("plugin file count exceeds %d", MaxPluginFiles)
			}
			plugin, err := loadYAMLPlugin(root, relative, budget)
			if err != nil {
				return fmt.Errorf("plugin %q: %w", relative, err)
			}
			if !seen[plugin.Name] {
				plugin.SourceFile = filepath.Join(directory, relative)
				seen[plugin.Name] = true
				plugins = append(plugins, plugin)
			}
			return nil
		})
		closeErr := root.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}
	}
	sort.Slice(plugins, func(i, j int) bool { return plugins[i].Name < plugins[j].Name })
	return plugins, nil
}

func walkPluginDirectory(root *os.Root, directory string, depth int, budget *pluginBudget, candidate func(string) error) error {
	if depth > MaxPluginDepth {
		return fmt.Errorf("plugin directory depth exceeds %d", MaxPluginDepth)
	}
	file, err := root.Open(directory)
	if err != nil {
		return fmt.Errorf("open plugin directory: %w", err)
	}
	var entries []fs.DirEntry
	for {
		batch, readErr := file.ReadDir(128)
		budget.entries += len(batch)
		if budget.entries > MaxPluginEntries {
			file.Close()
			return fmt.Errorf("plugin entries exceed %d", MaxPluginEntries)
		}
		entries = append(entries, batch...)
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			file.Close()
			return fmt.Errorf("read plugin directory: %w", readErr)
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		relative := filepath.Join(directory, entry.Name())
		if entry.IsDir() {
			if err := walkPluginDirectory(root, relative, depth+1, budget, candidate); err != nil {
				return err
			}
		} else if ext := filepath.Ext(entry.Name()); ext == ".yaml" || ext == ".yml" {
			if err := candidate(relative); err != nil {
				return err
			}
		}
	}
	return nil
}

func readPluginFile(root *os.Root, relative string, budget *pluginBudget) ([]byte, error) {
	info, err := root.Lstat(relative)
	if err != nil {
		return nil, fmt.Errorf("stat plugin input: %w", err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("plugin input must be a regular file without a final symlink")
	}
	if info.Size() > MaxPluginFileBytes {
		return nil, fmt.Errorf("plugin input exceeds %d bytes", MaxPluginFileBytes)
	}
	file, err := root.Open(relative)
	if err != nil {
		return nil, fmt.Errorf("open confined plugin input: %w", err)
	}
	actual, err := file.Stat()
	if err != nil || !actual.Mode().IsRegular() {
		file.Close()
		return nil, fmt.Errorf("opened plugin input is not a regular file")
	}
	data, readErr := io.ReadAll(io.LimitReader(file, MaxPluginFileBytes+1))
	closeErr := file.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read plugin input: %w", readErr)
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if len(data) > MaxPluginFileBytes {
		return nil, fmt.Errorf("plugin input exceeds %d bytes", MaxPluginFileBytes)
	}
	budget.bytes += len(data)
	if budget.bytes > MaxPluginTotalBytes {
		return nil, fmt.Errorf("plugin inputs exceed %d total bytes", MaxPluginTotalBytes)
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("plugin input must be UTF-8")
	}
	return data, nil
}

func loadYAMLPlugin(root *os.Root, relative string, budget *pluginBudget) (YAMLPlugin, error) {
	data, err := readPluginFile(root, relative, budget)
	if err != nil {
		return YAMLPlugin{}, err
	}
	// Inspect syntax separately so aliases/merge keys are rejected, not expanded.
	var document yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&document); err != nil {
		return YAMLPlugin{}, fmt.Errorf("invalid YAML document")
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return YAMLPlugin{}, fmt.Errorf("exactly one YAML document is required")
	}
	if !unambiguousYAML(&document, 0) {
		return YAMLPlugin{}, fmt.Errorf("YAML aliases, merge keys or excessive nesting are not supported")
	}
	var config yamlPluginConfig
	strict := yaml.NewDecoder(bytes.NewReader(data))
	strict.KnownFields(true)
	if err := strict.Decode(&config); err != nil {
		return YAMLPlugin{}, fmt.Errorf("invalid plugin schema, unknown/duplicate fields or field types")
	}
	if config.Name == "" || len(config.Name) > 128 || config.Name != strings.TrimSpace(config.Name) || strings.IndexFunc(config.Name, unicode.IsControl) >= 0 || strings.ContainsAny(config.Name, `/\`) || reservedPluginName(config.Name) {
		return YAMLPlugin{}, fmt.Errorf("plugin name is invalid or reserved")
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}
	if len(config.Queries) == 0 {
		return YAMLPlugin{}, fmt.Errorf("plugin needs at least one query")
	}
	budget.queries += len(config.Queries)
	if budget.queries > MaxPluginQueries {
		return YAMLPlugin{}, fmt.Errorf("plugin queries exceed %d", MaxPluginQueries)
	}
	plugin := YAMLPlugin{Name: config.Name, Version: config.Version, Description: config.Description, Author: config.Author, License: config.License}
	queryRoot, err := root.OpenRoot(filepath.Dir(relative))
	if err != nil {
		return YAMLPlugin{}, fmt.Errorf("open query directory: %w", err)
	}
	defer queryRoot.Close()
	for index, query := range config.Queries {
		if strings.TrimSpace(query.ID) == "" || strings.TrimSpace(query.Description) == "" || strings.TrimSpace(query.ResourceType) == "" || query.ResourceType != strings.TrimSpace(query.ResourceType) {
			return YAMLPlugin{}, fmt.Errorf("query %d requires aprlGuid, description and recommendationResourceType", index+1)
		}
		if query.QueryFile != "" {
			if !portableQueryPath(query.QueryFile) {
				return YAMLPlugin{}, fmt.Errorf("query %d has a non-confined queryFile", index+1)
			}
			bytes, err := readPluginFile(queryRoot, filepath.FromSlash(query.QueryFile), budget)
			if err != nil {
				return YAMLPlugin{}, fmt.Errorf("query %d: %w", index+1, err)
			}
			query.Query = string(bytes)
		}
		if strings.TrimSpace(query.Query) == "" {
			return YAMLPlugin{}, fmt.Errorf("query %d requires nonempty query text", index+1)
		}
		definition := normalizeRecommendation(rawRecommendation{
			ID: query.ID, Description: query.Description, Category: query.Category, Impact: query.Impact,
			ResourceType: query.ResourceType, MetadataState: query.State, LongDescription: query.LongDescription,
			PotentialBenefits: query.PotentialBenefits, PGVerified: query.PGVerified, AutomationAvailable: query.AutomationAvailable,
			Tags: query.Tags, LearnMore: query.LearnMore,
		}, config.Name)
		// Pinned YAML conversion accepts recommendationTypeId but does not project it.
		definition.Query = query.Query
		plugin.Definitions = append(plugin.Definitions, definition)
	}
	return plugin, nil
}

func unambiguousYAML(node *yaml.Node, depth int) bool {
	if depth > 32 || node.Kind == yaml.AliasNode || node.Tag == "!!merge" {
		return false
	}
	for _, child := range node.Content {
		if !unambiguousYAML(child, depth+1) {
			return false
		}
	}
	return true
}

func portableQueryPath(value string) bool {
	if len(value) > 4096 || len(strings.Split(value, "/")) > MaxPluginDepth || strings.ContainsAny(value, `\:<>"|?*`) || !fs.ValidPath(value) || !filepath.IsLocal(filepath.FromSlash(value)) || path.IsAbs(value) {
		return false
	}
	for _, character := range value {
		if character < 32 {
			return false
		}
	}
	for _, part := range strings.Split(value, "/") {
		base := strings.TrimRight(strings.ToUpper(strings.SplitN(part, ".", 2)[0]), " ")
		if base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" || base == "CONIN$" || base == "CONOUT$" || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") {
			return false
		}
		if len(base) >= 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && ((len(base[3:]) == 1 && strings.Contains("123456789", base[3:])) || base[3:] == "¹" || base[3:] == "²" || base[3:] == "³") {
			return false
		}
	}
	return true
}

// PluginDefinitions returns independent definitions in sorted plugin-name order.
func PluginDefinitions(plugins []YAMLPlugin) []assessment.RecommendationDefinition {
	ordered := append([]YAMLPlugin(nil), plugins...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Name < ordered[j].Name })
	var definitions []assessment.RecommendationDefinition
	for _, plugin := range ordered {
		for _, definition := range plugin.Definitions {
			definitions = append(definitions, cloneDefinition(definition))
		}
	}
	return definitions
}

func cloneDefinition(definition assessment.RecommendationDefinition) assessment.RecommendationDefinition {
	definition.Tags = append([]string(nil), definition.Tags...)
	definition.LearnMore = append([]assessment.LearnMoreLink(nil), definition.LearnMore...)
	return definition
}

// CopyDefinitions separates mutable metadata slices from caller-owned inputs.
func CopyDefinitions(definitions []assessment.RecommendationDefinition) []assessment.RecommendationDefinition {
	copied := make([]assessment.RecommendationDefinition, len(definitions))
	for index, definition := range definitions {
		copied[index] = cloneDefinition(definition)
	}
	return copied
}
