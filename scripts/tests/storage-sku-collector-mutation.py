#!/usr/bin/env python3
"""Compile isolated Storage reversals, require named failures and exact restoration."""
from pathlib import Path
import os
import subprocess

ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / 'internal/plugins/region/storage_sku_collector.go'
GO = os.environ.get('QA_GO', 'go')


def run(select):
    return subprocess.run([GO, 'test', './internal/plugins/region', '-run', select, '-count=1', '-v'],
                          cwd=ROOT, capture_output=True, text=True, encoding='utf-8',
                          errors='backslashreplace', timeout=180)


def check(select):
    result = run(select)
    if result.returncode:
        raise SystemExit('Storage restored baseline failed\n' + result.stdout + result.stderr)


original = SOURCE.read_bytes()
text = original.decode('utf-8')
manifest = 'TestStorageSKUFrozenManifest'
controls = (
    ('request-version', 'storageSKUVersion  = "2023-01-01"', 'storageSKUVersion  = "2021-07-01"', manifest),
    ('getter-bound', 'target.String(), MaxStorageSKUBytes)', 'target.String(), 0)', manifest),
    ('pre-auth-scope', 'scope[id] == ""', 'false', manifest),
    ('no-paging', 'envelope["nextLink"]; present {', 'envelope["nextLink"]; present && false {', manifest),
    ('json-admission', 'if !quotaJSON(ctx, body) {', 'if false {', manifest),
    ('display-alias', 'physical, ok := aliases[value]; ok {', 'physical, ok := aliases[value]; ok && false {', manifest),
    ('stage-physical', 'if storageStage(value) {\n\t\treturn "", true\n\t}', 'if storageStage(value) {\n\t\treturn strings.TrimSuffix(value, "(stage)"), true\n\t}', manifest),
    ('regional-union', 'listed: prior.listed || flags.listed', 'listed: flags.listed', manifest),
    ('kind-grouping', 'name, kind = strings.ToLower(name), strings.ToLower(kind)', 'name, kind = strings.ToLower(name), "storage"', manifest),
    ('missing-locations', 'flags.listed = flags.listed || location == region\n\t\t}\n\t} else {\n\t\tflags.unknown = true', 'flags.listed = flags.listed || location == region\n\t\t}\n\t} else {\n\t\tflags.unknown = false', manifest),
    ('restriction-scope', 'if interpretable && known && location == region {', 'if interpretable && known && location != "" {', manifest),
    ('restriction-union', 'restricted: prior.restricted || flags.restricted', 'restricted: flags.restricted', manifest),
    ('unknown-group-gating', 'if found && !group.Known {', 'if found && !group.Known && false {', manifest),
    ('unqualified-gating', 'if missingKind[name] || c.UnqualifiedNames[name] {', 'if missingKind[name] {', manifest),
    ('required-kind', 'group, found := c.Groups[name][kind]', 'group, found := c.Groups[name]["storage"]\n\t\t\t_ = kind', 'TestStorageSKURequiredKindIsolation'),
    ('atomic-admission', 'if code != "" {\n\t\treturn fail(code)\n\t}', 'if code != "" && false {\n\t\treturn fail(code)\n\t}', manifest),
    ('status-error-pair', 'if response != nil && response.StatusCode != http.StatusOK {', 'if err == nil && response != nil && response.StatusCode != http.StatusOK {', 'TestStorageSKUActualHTTPClient'),
    ('context-priority', 'return nil, cause', 'return result, nil', 'TestStorageSKUCancellation'),
    ('body-bytes', 'if len(body) > MaxStorageSKUBytes {', 'if false {', 'TestStorageSKUBodyAndRowsBounds'),
    ('row-limit', 'if len(rows) > MaxStorageSKURows {', 'if len(rows) > MaxStorageSKURows+1 {', 'TestStorageSKUBodyAndRowsBounds'),
    ('raw-work', 'b.work >= MaxStorageSKUWork', 'b.work > MaxStorageSKUWork', 'TestStorageSKUJSONBounds'),
    ('projection-work', 'b.work >= MaxStorageSKUWork', 'b.work > MaxStorageSKUWork', 'TestStorageSKUProjectionBounds'),
    ('projection-text', 'return b.ctx.Err() == nil && auxiliaryText(value, &b.text)', 'ignored := 0\n\treturn b.ctx.Err() == nil && auxiliaryText(value, &ignored)', 'TestStorageSKUProjectionBounds'),
    ('capability-conflict', 'availableValue != "" && availableValue != normalized', 'false && availableValue != normalized', manifest),
    ('late-unknown-context', 'incomplete := func() (SKUEvidence, error) {\n\t\tif err := ctx.Err(); err != nil {\n\t\t\treturn SKUEvidence{}, err\n\t\t}\n\t\treturn unknown, nil\n\t}', 'incomplete := func() (SKUEvidence, error) {\n\t\treturn unknown, nil\n\t}', 'TestStorageSKULateUnknownCancellation'),
)
check('^TestStorageSKU')
try:
    for name, old, replacement, expected in controls:
        if text.count(old) != 1:
            raise SystemExit('Storage production anchor changed: ' + name)
        SOURCE.write_bytes(text.replace(old, replacement, 1).encode('utf-8'))
        compiled = run('^$')
        if compiled.returncode:
            raise SystemExit('Storage mutation did not compile: ' + name + '\n' + compiled.stdout + compiled.stderr)
        failed = run('^' + expected + '$')
        if failed.returncode == 0 or ('--- FAIL: ' + expected) not in failed.stdout or '[build failed]' in failed.stdout + failed.stderr:
            raise SystemExit('Storage mutation missed named assertion: ' + name + '\n' + failed.stdout + failed.stderr)
        SOURCE.write_bytes(original)
        check('^' + expected + '$')
        print('Storage SKU compiling mutation rejected: ' + name + '; restored named baseline passed', flush=True)
finally:
    SOURCE.write_bytes(original)
check('^TestStorageSKU')
if SOURCE.read_bytes() != original:
    raise SystemExit('Storage production source was not restored')
