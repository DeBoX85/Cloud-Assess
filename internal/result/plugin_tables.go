package result

import (
	"fmt"
	"sort"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const PluginSchemaVersion = "1.1"

// BuildWithPluginTables preserves the ordinary builder and schema when no
// plugin table was explicitly supplied. Otherwise it validates and owns a full
// copy of the versioned table extension.
func BuildWithPluginTables(input Input, tables []assessment.PluginTable) (*AssessmentResult, error) {
	if err := assessment.ValidatePluginTables(tables); err != nil {
		return nil, err
	}
	r := Build(input)
	if len(tables) == 0 {
		return r, nil
	}
	r.SchemaVersion = PluginSchemaVersion
	r.PluginTables = make([]assessment.PluginTable, len(tables))
	for i, table := range tables {
		table.Columns = append([]string(nil), table.Columns...)
		table.Rows = make([]assessment.PluginRow, len(tables[i].Rows))
		for j, row := range tables[i].Rows {
			row.Cells = append([]string(nil), row.Cells...)
			table.Rows[j] = row
		}
		table.Health.Warnings = append([]assessment.AssessmentWarning(nil), table.Health.Warnings...)
		if table.Health.Error != nil {
			e := *table.Health.Error
			table.Health.Error = &e
		}
		r.PluginTables[i] = table
	}
	sort.Slice(r.PluginTables, func(i, j int) bool {
		left, right := r.PluginTables[i], r.PluginTables[j]
		if left.Metadata.Name != right.Metadata.Name {
			return left.Metadata.Name < right.Metadata.Name
		}
		return left.ID < right.ID
	})
	r.Completeness = pluginCompleteness(r.Completeness, r.PluginTables)
	if err := r.ValidatePluginExtension(); err != nil {
		return nil, err
	}
	return r, nil
}

// ValidatePluginExtension rejects inconsistent manually modified library data
// before any renderer mutates an output file. Core-only schema behavior stays
// unchanged; other whole-result invariants belong to their existing contracts.
func (r *AssessmentResult) ValidatePluginExtension() error {
	if len(r.PluginTables) > 0 && r.SchemaVersion != PluginSchemaVersion {
		return fmt.Errorf("plugin table extension requires assessment schema 1.1")
	}
	if len(r.PluginTables) > 0 {
		switch r.Completeness {
		case assessment.CompletenessComplete, assessment.CompletenessCompleteWithWarnings, assessment.CompletenessPartial, assessment.CompletenessFailed:
		default:
			return fmt.Errorf("plugin tables require known assessment completeness")
		}
		if pluginCompleteness(r.Completeness, r.PluginTables) != r.Completeness {
			return fmt.Errorf("assessment completeness conceals plugin table failure or warning")
		}
	}
	return assessment.ValidatePluginTables(r.PluginTables)
}

func pluginCompleteness(current assessment.Completeness, tables []assessment.PluginTable) assessment.Completeness {
	if current == assessment.CompletenessFailed {
		return current
	}
	for _, table := range tables {
		if table.Health.Status == assessment.StageFailed || table.Health.Status == assessment.StageSkipped {
			current = assessment.CompletenessPartial
		}
		if table.Health.Status == assessment.StageCompletedWithWarnings && current == assessment.CompletenessComplete {
			current = assessment.CompletenessCompleteWithWarnings
		}
	}
	return current
}
