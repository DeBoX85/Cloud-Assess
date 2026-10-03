// AI governance behavior is characterized from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md and docs/AI_GOVERNANCE.md.
package aigov

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

const (
	Name           = "ai-gov"
	MaxAccounts    = 4096
	MaxPoints      = 65536
	MaxRows        = 8192
	MaxDeployments = 65536
	MaxTextBytes   = 16 << 20
	MaxLabelBytes  = 512
)

var columns = []string{"Subscription", "Resource Group", "Account Name", "Kind", "SKU", "Deployment Name", "Model Name", "Model Version", "Model Format", "SKU Capacity", "Version Upgrade Option", "Spillover Enabled", "Spillover Deployment", "Hour", "Status Code", "Request Count"}
var accountID = regexp.MustCompile(`(?i)^/subscriptions/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/resourceGroups/([^/]+)/providers/Microsoft\.CognitiveServices/accounts/([^/]+)$`)
var subscriptionID = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func Metadata() assessment.PluginMetadata {
	return assessment.PluginMetadata{Name: Name, Version: "1.0.0", Description: "Checks AI Governance", Author: "Azure Quick Review Team", License: "MIT", Type: "internal"}
}
func PendingTable() assessment.PluginTable {
	return assessment.PluginTable{SchemaVersion: assessment.PluginTableSchemaVersion, ID: "requests", Metadata: Metadata(), SheetName: "AI Throttling", Description: "Analysis of AI/Cognitive Services accounts by hour, model, and status code", Columns: append([]string(nil), columns...), Rows: []assessment.PluginRow{}, Health: assessment.StageExecution{Name: Name, Status: assessment.StageSkipped, Warnings: []assessment.AssessmentWarning{{Code: "plugin_not_run", Message: "requested plugin has not executed"}}}}
}

// Account is an already decoded discovery row. HTTP/discovery and tag-scope
// construction are separate obligations; this projection never queries Azure.
type Account struct{ ID, SubscriptionID, ResourceGroup, Name, Kind, SKU string }

// Point preserves missing fields and explicit empty dimensions independently.
// The decoder must use last matching case-insensitive dimension-name values,
// as the source does. Missing timestamp/count is accounted as incomplete input.
type Point struct {
	ResourceID                string
	Timestamp                 *time.Time
	Count                     *float64
	Deployment, Model, Status *string
}

// Deployment distinguishes absent metadata from present empty/zero values.
type Deployment struct {
	Name, ModelVersion, ModelFormat, Upgrade, Spillover *string
	Capacity                                            *int64
}

// Failed must be set by the request adapter for failed/incomplete enrichment.
// It retains metric rows with source defaults and exposes warning health.
type DeploymentSet struct {
	Values []Deployment
	Failed bool
}
type Filter interface{ IsServiceExcluded(string) bool }
type key struct{ resource, hour, deployment, model, status string }
type info struct{ version, format, capacity, upgrade, spilloverEnabled, spillover string }

func defaults() info { return info{"N/A", "N/A", "N/A", "N/A", "No", "N/A"} }
func safe(s string, limit int) bool {
	return len(s) <= limit && utf8.ValidString(s) && strings.IndexFunc(s, func(r rune) bool { return unicode.IsControl(r) || r == 0xfffd || r == 0xfffe || r == 0xffff }) < 0
}
func dimension(s *string) string {
	if s == nil {
		return "Unknown"
	}
	return *s
}
func deploymentInfo(d Deployment) (info, bool) {
	v := defaults()
	for _, s := range []*string{d.Name, d.ModelVersion, d.ModelFormat, d.Upgrade, d.Spillover} {
		if s != nil && !safe(*s, MaxLabelBytes) {
			return v, false
		}
	}
	if d.ModelVersion != nil {
		v.version = *d.ModelVersion
	}
	if d.ModelFormat != nil {
		v.format = *d.ModelFormat
	}
	if d.Upgrade != nil {
		v.upgrade = *d.Upgrade
	}
	if d.Capacity != nil {
		if *d.Capacity < 0 {
			return v, false
		}
		v.capacity = fmt.Sprintf("%d", *d.Capacity)
	}
	if d.Spillover != nil {
		v.spilloverEnabled = "Yes"
		v.spillover = *d.Spillover
	}
	return v, true
}

// Project preserves source accumulation order, float formatting, exact deployment
// matching and timestamp-zone display. Only final row order is normalized. Fatal
// limits/cancellation retain previously accepted sums with failed health. Inputs
// are caller-owned and must not be concurrently mutated; outputs share no slices.
func Project(ctx context.Context, subscriptions map[string]string, accounts []Account, points []Point, deployments map[string]DeploymentSet, filter Filter) (assessment.PluginTable, error) {
	table := PendingTable()
	table.Health = assessment.StageExecution{Name: Name, Status: assessment.StageCompleted}
	selected := map[string]Account{}
	seenAccounts := map[string]bool{}
	names := map[string]string{}
	enrichment := map[string]map[string]info{}
	sums := map[key]float64{}
	malformed, unavailable, textBytes, deploymentCount, rowText := 0, 0, 0, 0, 0
	build := func(k key, count float64) assessment.PluginRow {
		a := selected[k.resource]
		i := defaults()
		if v, ok := enrichment[k.resource][k.deployment]; ok {
			i = v
		}
		kind, sku := a.Kind, a.SKU
		if kind == "" {
			kind = "Unknown"
		}
		if sku == "" {
			sku = "Unknown"
		}
		return assessment.PluginRow{SubscriptionID: strings.ToLower(a.SubscriptionID), Cells: []string{names[strings.ToLower(a.SubscriptionID)], a.ResourceGroup, a.Name, kind, sku, k.deployment, k.model, i.version, i.format, i.capacity, i.upgrade, i.spilloverEnabled, i.spillover, k.hour, k.status, fmt.Sprintf("%.0f", count)}}
	}
	finish := func() {
		table.Rows = make([]assessment.PluginRow, 0, len(sums))
		for k, v := range sums {
			table.Rows = append(table.Rows, build(k, v))
		}
		sort.Slice(table.Rows, func(i, j int) bool {
			a, b := table.Rows[i], table.Rows[j]
			if c := slices.Compare(a.Cells, b.Cells); c != 0 {
				return c < 0
			}
			return a.SubscriptionID < b.SubscriptionID
		})
		table.Health.Records = len(table.Rows)
		table.Health.Warnings = nil
		if malformed > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "ai_malformed_input", Message: fmt.Sprintf("skipped %d invalid AI governance input entries", malformed)})
		}
		if unavailable > 0 {
			table.Health.Warnings = append(table.Health.Warnings, assessment.AssessmentWarning{Code: "ai_deployment_unavailable", Message: fmt.Sprintf("deployment enrichment incomplete for %d accounts", unavailable)})
		}
	}
	fail := func(code string, cause error) (assessment.PluginTable, error) {
		finish()
		table.Health.Status = assessment.StageFailed
		table.Health.Error = &assessment.AssessmentError{Code: code, Message: "AI governance projection returned incomplete data"}
		if cause != nil {
			return table, fmt.Errorf("AI governance projection stopped: %w", cause)
		}
		return table, fmt.Errorf("AI governance projection returned incomplete data")
	}
	addText := func(values ...string) bool {
		for _, v := range values {
			textBytes += len(v)
		}
		return textBytes <= MaxTextBytes
	}
	if len(accounts) > MaxAccounts || len(subscriptions) > MaxAccounts || len(points) > MaxPoints || len(deployments) > MaxAccounts {
		return fail("ai_input_limit", nil)
	}
	for id, name := range subscriptions {
		if err := ctx.Err(); err != nil {
			return fail("ai_cancelled", err)
		}
		normalized := strings.ToLower(id)
		if !subscriptionID.MatchString(id) || !safe(name, MaxLabelBytes) {
			return fail("ai_scope_invalid", nil)
		}
		if _, ok := names[normalized]; ok {
			return fail("ai_scope_invalid", nil)
		}
		if name == "" {
			name = id
		}
		if !addText(id, name) {
			return fail("ai_text_limit", nil)
		}
		names[normalized] = name
	}
	for _, a := range accounts {
		if err := ctx.Err(); err != nil {
			return fail("ai_cancelled", err)
		}
		m := accountID.FindStringSubmatch(a.ID)
		normalized := strings.ToLower(a.ID)
		if len(m) != 4 || !strings.EqualFold(m[1], a.SubscriptionID) || !strings.EqualFold(m[2], a.ResourceGroup) || !strings.EqualFold(m[3], a.Name) || !safe(a.ID, 2048) || !safe(a.ResourceGroup, MaxLabelBytes) || !safe(a.Name, MaxLabelBytes) || !safe(a.Kind, MaxLabelBytes) || !safe(a.SKU, MaxLabelBytes) {
			return fail("ai_account_invalid", nil)
		}
		if _, ok := names[strings.ToLower(a.SubscriptionID)]; !ok {
			return fail("ai_scope_invalid", nil)
		}
		if seenAccounts[normalized] {
			return fail("ai_account_invalid", nil)
		}
		seenAccounts[normalized] = true
		if !addText(a.ID, a.SubscriptionID, a.ResourceGroup, a.Name, a.Kind, a.SKU) {
			return fail("ai_text_limit", nil)
		}
		if filter != nil && filter.IsServiceExcluded(a.ID) {
			continue
		}
		selected[normalized] = a
	}
	if len(selected) > 0 {
		table.SheetName = "AI Gov"
	}
	seenDeployments := map[string]bool{}
	for id, set := range deployments {
		if err := ctx.Err(); err != nil {
			return fail("ai_cancelled", err)
		}
		normalized := strings.ToLower(id)
		if _, ok := selected[normalized]; !ok || seenDeployments[normalized] {
			return fail("ai_deployment_scope_invalid", nil)
		}
		seenDeployments[normalized] = true
		deploymentCount += len(set.Values)
		if deploymentCount > MaxDeployments {
			return fail("ai_deployment_limit", nil)
		}
		if set.Failed {
			unavailable++
		}
		if !addText(id) {
			return fail("ai_text_limit", nil)
		}
		enrichment[normalized] = map[string]info{}
		for _, d := range set.Values {
			if err := ctx.Err(); err != nil {
				return fail("ai_cancelled", err)
			}
			for _, v := range []*string{d.Name, d.ModelVersion, d.ModelFormat, d.Upgrade, d.Spillover} {
				if v != nil && !addText(*v) {
					return fail("ai_text_limit", nil)
				}
			}
			v, ok := deploymentInfo(d)
			if !ok || d.Name == nil {
				malformed++
				continue
			}
			enrichment[normalized][*d.Name] = v
		}
	}
	for _, p := range points {
		if err := ctx.Err(); err != nil {
			return fail("ai_cancelled", err)
		}
		normalized := strings.ToLower(p.ResourceID)
		if _, ok := selected[normalized]; !ok {
			return fail("ai_metric_scope_invalid", nil)
		}
		dep, model, status := dimension(p.Deployment), dimension(p.Model), dimension(p.Status)
		if !addText(p.ResourceID, dep, model, status) {
			return fail("ai_text_limit", nil)
		}
		if p.Timestamp == nil || p.Count == nil || p.Timestamp.Year() < 1 || p.Timestamp.Year() > 9999 || !safe(dep, MaxLabelBytes) || !safe(model, MaxLabelBytes) || !safe(status, MaxLabelBytes) || math.IsNaN(*p.Count) || math.IsInf(*p.Count, 0) || *p.Count < 0 {
			malformed++
			continue
		}
		if !seenDeployments[normalized] {
			unavailable++
			seenDeployments[normalized] = true
		}
		k := key{normalized, p.Timestamp.Format("2006-01-02 15:00"), dep, model, status}
		previous, exists := sums[k]
		next := previous + *p.Count
		if math.IsInf(next, 0) || math.IsNaN(next) {
			return fail("ai_numeric_limit", nil)
		}
		if !exists {
			if len(sums) >= MaxRows {
				return fail("ai_row_limit", nil)
			}
			// Reserve 512 bytes for changing count text; any finite float rendered with
			// %.0f fits. Include repeated account/enrichment cells, not just input labels.
			needed := 512
			for _, cell := range build(k, 0).Cells[:15] {
				needed += len(cell)
			}
			if rowText+needed > assessment.MaxPluginTextBytes-4096 {
				return fail("ai_output_text_limit", nil)
			}
			rowText += needed
		}
		sums[k] = next
	}
	if err := ctx.Err(); err != nil {
		return fail("ai_cancelled", err)
	}
	finish()
	if len(table.Health.Warnings) > 0 {
		table.Health.Status = assessment.StageCompletedWithWarnings
	}
	return table, nil
}
