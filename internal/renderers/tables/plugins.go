package tables

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/redact"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
)

var pluginContextSubscriptionID = regexp.MustCompile(`(?i)(?:/subscriptions/|\bsubscription(?:\s+id)?\s+)([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
var pluginMaskSubscriptionID = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// Plugin values are ordinary text, including identity-bearing links. Mask all
// known subscription IDs rather than assuming one dedicated column per plugin.
const maxPluginMaskIDs = 4096

func pluginMasker(data *result.AssessmentResult, enabled bool) (func(string) string, error) {
	if !enabled || len(data.PluginTables) == 0 {
		return func(s string) string { return s }, nil
	}
	ids := map[string]string{}
	overflow := false
	add := func(id string) {
		if pluginMaskSubscriptionID.MatchString(id) {
			key := strings.ToLower(id)
			if _, exists := ids[key]; !exists && len(ids) == maxPluginMaskIDs {
				overflow = true
				return
			}
			ids[key] = redact.SubscriptionID(id, true)
		}
	}
	infer := func(text string) {
		for _, m := range pluginContextSubscriptionID.FindAllStringSubmatch(text, -1) {
			add(m[1])
		}
	}
	for _, table := range data.PluginTables {
		for _, text := range []string{table.Description, table.Metadata.Description, table.Metadata.Author, table.Metadata.License} {
			infer(text)
		}
		for _, col := range table.Columns {
			infer(col)
		}
		for _, warning := range table.Health.Warnings {
			infer(warning.Message)
		}
		if table.Health.Error != nil {
			infer(table.Health.Error.Message)
		}
		for _, row := range table.Rows {
			add(row.SubscriptionID)
			for _, cell := range row.Cells {
				infer(cell)
			}
		}
	}
	if scope := data.Scope; scope != nil {
		for _, group := range [][]string{scope.RequestedSubscriptionIDs, scope.IncludedSubscriptionIDs, scope.ExcludedSubscriptionIDs, scope.UnresolvedSubscriptionIDs} {
			for _, id := range group {
				add(id)
			}
		}
		for _, sub := range scope.ResolvedSubscriptions {
			add(sub.SubscriptionID)
		}
	}
	// Mixed assessments may be constructed without a scope resolution by trusted
	// library callers. Include identities from their canonical core records too.
	for _, r := range data.Resources {
		add(r.SubscriptionID)
	}
	for _, r := range data.Findings {
		add(r.SubscriptionID)
	}
	for _, r := range data.OutOfScope {
		add(r.SubscriptionID)
	}
	for _, r := range data.Advisor {
		add(r.SubscriptionID)
	}
	for _, r := range data.Defender {
		add(r.SubscriptionID)
	}
	for _, r := range data.DefenderRecommendations {
		add(r.SubscriptionID)
	}
	for _, r := range data.AzurePolicy {
		add(r.SubscriptionID)
	}
	for _, r := range data.ArcSQL {
		add(r.SubscriptionID)
	}
	for _, r := range data.Costs {
		add(r.SubscriptionID)
	}
	for _, r := range data.ResourceTypes {
		add(r.SubscriptionID)
	}
	if overflow {
		return nil, fmt.Errorf("plugin masking identity limit exceeded")
	}
	if len(ids) == 0 {
		return func(s string) string { return s }, nil
	}
	patterns := make([]string, 0, len(ids))
	for id := range ids {
		patterns = append(patterns, regexp.QuoteMeta(id))
	}
	sort.Strings(patterns)
	pattern, err := regexp.Compile("(?i)(" + strings.Join(patterns, "|") + ")")
	if err != nil {
		return nil, fmt.Errorf("plugin masking pattern could not be constructed")
	}
	return func(s string) string {
		return pattern.ReplaceAllStringFunc(s, func(id string) string { return ids[strings.ToLower(id)] })
	}, nil
}

func pluginTables(data *result.AssessmentResult, mask func(string) string) []Table {
	projected := make([]Table, 0, len(data.PluginTables))
	for _, table := range data.PluginTables {
		rows := make([][]string, 1, len(table.Rows)+1)
		rows[0] = make([]string, len(table.Columns))
		for i, col := range table.Columns {
			rows[0][i] = mask(col)
		}
		for _, row := range table.Rows {
			cells := make([]string, len(row.Cells))
			for i, cell := range row.Cells {
				cells[i] = mask(cell)
			}
			rows = append(rows, cells)
		}
		// Explicit requested tables render even if empty/failed/skipped. Their
		// health remains visible on Assessment Status, unlike absent tables.
		projected = append(projected, Table{Key: table.Key(), SheetName: mask(table.SheetName), Rows: rows})
	}
	return projected
}

func pluginStatusRows(data *result.AssessmentResult, mask func(string) string) [][]string {
	rows := make([][]string, 0, len(data.PluginTables))
	for _, table := range data.PluginTables {
		h := table.Health
		codes := make([]string, len(h.Warnings))
		for i, w := range h.Warnings {
			codes[i] = w.Code
		}
		errorCode, errorMessage := "", ""
		if h.Error != nil {
			errorCode, errorMessage = h.Error.Code, mask(h.Error.Message)
		}
		rows = append(rows, []string{string(data.Completeness), data.SchemaVersion, mask(data.ScopeID), formatTime(data.GeneratedAt), "plugin:" + table.Metadata.Name + ":" + table.ID, string(h.Status), strconv.Itoa(h.Records), strconv.Itoa(len(h.Warnings)), strings.Join(codes, ","), errorCode, errorMessage, formatTime(h.StartedAt), formatTime(h.FinishedAt), h.FinishedAt.Sub(h.StartedAt).String()})
	}
	return rows
}
