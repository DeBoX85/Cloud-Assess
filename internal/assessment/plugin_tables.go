package assessment

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	PluginTableSchemaVersion = "1.0"
	MaxPluginTables          = 64
	MaxPluginColumns         = 64
	MaxPluginTableRows       = 65536
	MaxPluginRows            = 262144
	MaxPluginTextBytes       = 16 << 20
	MaxPluginCellUnits       = 32767
)

type PluginMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	License     string `json:"license"`
	Type        string `json:"type"`
}

type PluginRow struct {
	SubscriptionID string   `json:"subscriptionId,omitempty"`
	Cells          []string `json:"cells"`
}

type PluginTable struct {
	SchemaVersion string         `json:"schemaVersion"`
	ID            string         `json:"id"`
	Metadata      PluginMetadata `json:"metadata"`
	SheetName     string         `json:"sheetName"`
	Description   string         `json:"description"`
	Columns       []string       `json:"columns"`
	Rows          []PluginRow    `json:"rows"`
	Health        StageExecution `json:"health"`
}

func (t PluginTable) Key() string { return "plugin_" + t.Metadata.Name + "_" + t.ID }

var pluginLabel = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
var pluginSubscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
var pluginHealthCode = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,127}$`)

// ValidatePluginTables checks the whole set before rendering any output. Errors
// identify the guard without echoing provider cells, labels or health messages.
func ValidatePluginTables(tables []PluginTable) error {
	if len(tables) > MaxPluginTables {
		return fmt.Errorf("plugin table count limit exceeded")
	}
	keys := map[string]bool{}
	sheets := []string{}
	for _, name := range []string{"Assessment Status", "Recommendations", "ImpactedResources", "ResourceTypes", "Inventory", "Advisor", "Azure Policy", "Arc SQL", "DefenderRecommendations", "Defender", "OutOfScope", "Costs"} {
		sheets = append(sheets, name)
	}
	textBytes, rows := 0, 0
	text := func(s string, maxUnits int) bool {
		if !utf8.ValidString(s) {
			return false
		}
		units := 0
		for _, r := range s {
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' || r == 0xfffe || r == 0xffff {
				return false
			}
			units++
			if r > 0xffff {
				units++
			}
			if units > maxUnits {
				return false
			}
		}
		textBytes += len(s)
		return textBytes <= MaxPluginTextBytes
	}
	for _, table := range tables {
		if table.SchemaVersion != PluginTableSchemaVersion || !pluginLabel.MatchString(table.ID) || !pluginLabel.MatchString(table.Metadata.Name) || table.Metadata.Type != "internal" || keys[table.Key()] {
			return fmt.Errorf("invalid or duplicate plugin table identity")
		}
		keys[table.Key()] = true
		if strings.TrimSpace(table.SheetName) == "" || utf8.RuneCountInString(table.SheetName) > 31 || strings.ContainsAny(table.SheetName, `:/\?*[]`) || strings.HasPrefix(table.SheetName, "'") || strings.HasSuffix(table.SheetName, "'") || strings.IndexFunc(table.SheetName, unicode.IsControl) >= 0 || slices.ContainsFunc(sheets, func(name string) bool { return strings.EqualFold(name, table.SheetName) }) {
			return fmt.Errorf("invalid or colliding plugin sheet name")
		}
		sheets = append(sheets, table.SheetName)
		if strings.TrimSpace(table.Metadata.Version) == "" || !text(table.SheetName, 31) || !text(table.Metadata.Version, 128) || !text(table.Metadata.Description, 4096) || !text(table.Metadata.Author, 512) || !text(table.Metadata.License, 128) || !text(table.Description, 4096) {
			return fmt.Errorf("invalid or excessive plugin metadata text")
		}
		if len(table.Columns) == 0 || len(table.Columns) > MaxPluginColumns {
			return fmt.Errorf("invalid plugin column count")
		}
		columns := []string{}
		for _, col := range table.Columns {
			if strings.TrimSpace(col) == "" || slices.ContainsFunc(columns, func(name string) bool { return strings.EqualFold(name, col) }) || !text(col, 128) || strings.IndexFunc(col, unicode.IsControl) >= 0 {
				return fmt.Errorf("invalid or duplicate plugin column")
			}
			columns = append(columns, col)
		}
		rows += len(table.Rows)
		if len(table.Rows) > MaxPluginTableRows || rows > MaxPluginRows {
			return fmt.Errorf("plugin row limit exceeded")
		}
		for _, row := range table.Rows {
			if len(row.Cells) != len(table.Columns) || row.SubscriptionID != "" && !pluginSubscriptionID.MatchString(row.SubscriptionID) {
				return fmt.Errorf("invalid plugin row width or subscription identity")
			}
			for _, cell := range row.Cells {
				if !text(cell, MaxPluginCellUnits) {
					return fmt.Errorf("invalid or excessive plugin cell text")
				}
			}
		}
		h := table.Health
		if h.Name != table.Metadata.Name || h.Records != len(table.Rows) || len(h.Warnings) > 256 {
			return fmt.Errorf("inconsistent plugin table health")
		}
		for _, warning := range h.Warnings {
			if !pluginHealthCode.MatchString(warning.Code) || !text(warning.Message, 4096) {
				return fmt.Errorf("invalid plugin warning")
			}
		}
		if h.Error != nil && (!pluginHealthCode.MatchString(h.Error.Code) || strings.TrimSpace(h.Error.Message) == "" || !text(h.Error.Message, 4096)) {
			return fmt.Errorf("invalid plugin error")
		}
		switch h.Status {
		case StageCompleted:
			if h.Error != nil || len(h.Warnings) != 0 {
				return fmt.Errorf("completed plugin table has failure information")
			}
		case StageCompletedWithWarnings:
			if h.Error != nil || len(h.Warnings) == 0 {
				return fmt.Errorf("warning plugin table has inconsistent health")
			}
		case StageFailed:
			if h.Error == nil {
				return fmt.Errorf("failed plugin table lacks error")
			}
		case StageSkipped:
			if h.Error != nil || len(table.Rows) != 0 {
				return fmt.Errorf("skipped plugin table contains data or error")
			}
		default:
			return fmt.Errorf("unknown plugin table status")
		}
	}
	return nil
}
