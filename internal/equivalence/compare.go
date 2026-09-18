package equivalence

import (
	"sort"
)

func Compare(reference, target Projection) Report {
	report := Report{
		Equivalent:     true,
		ReferenceNotes: append([]string(nil), reference.Notes...),
		TargetNotes:    append([]string(nil), target.Notes...),
	}

	if !reference.Comparable {
		report.Equivalent = false
		report.Preconditions = append(report.Preconditions, "reference projection is not suitable for semantic comparison")
	}
	if !target.Comparable {
		report.Equivalent = false
		report.Preconditions = append(report.Preconditions, "target projection is not suitable for semantic comparison")
	}

	for _, name := range datasetNames(reference, target) {
		left := reference.Datasets[name]
		right := target.Datasets[name]
		diff := compareDataset(name, left, right)
		if diff.CoverageMismatch || len(diff.Missing) > 0 || len(diff.Extra) > 0 || len(diff.Changed) > 0 {
			report.Equivalent = false
		}
		report.Datasets = append(report.Datasets, diff)
	}
	return report
}

func compareDataset(name string, reference, target Dataset) DatasetDiff {
	diff := DatasetDiff{
		Name:             name,
		ReferenceEnabled: reference.Enabled,
		TargetEnabled:    target.Enabled,
		CoverageMismatch: reference.Enabled != target.Enabled,
	}

	keys := make([]string, 0, len(reference.Records)+len(target.Records))
	seen := map[string]struct{}{}
	for key := range reference.Records {
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	for key := range target.Records {
		if _, ok := seen[key]; ok {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		left, leftOK := reference.Records[key]
		right, rightOK := target.Records[key]
		switch {
		case leftOK && !rightOK:
			diff.Missing = append(diff.Missing, cloneRecord(left))
		case !leftOK && rightOK:
			diff.Extra = append(diff.Extra, cloneRecord(right))
		default:
			if changed := compareRecordFields(key, left.Fields, right.Fields); len(changed) > 0 {
				diff.Changed = append(diff.Changed, ChangedRecord{Key: key, Fields: changed})
			}
		}
	}
	return diff
}

func compareRecordFields(key string, reference, target map[string]string) []FieldDiff {
	fields := make([]string, 0, len(reference)+len(target))
	seen := map[string]struct{}{}
	for field := range reference {
		seen[field] = struct{}{}
		fields = append(fields, field)
	}
	for field := range target {
		if _, ok := seen[field]; ok {
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)

	var changed []FieldDiff
	for _, field := range fields {
		left := reference[field]
		right := target[field]
		if left == right {
			continue
		}
		changed = append(changed, FieldDiff{
			Field:     field,
			Reference: left,
			Target:    right,
		})
	}
	return changed
}

func cloneRecord(value Record) Record {
	fields := make(map[string]string, len(value.Fields))
	for key, item := range value.Fields {
		fields[key] = item
	}
	return Record{Key: value.Key, Fields: fields}
}
