#!/usr/bin/env python3
"""Require compiling inventory faults, named assertion failures and restoration."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/inventory.go'
MUTATIONS = (
    ('empty-selected-name', 'if name == "" {', 'if name == "" && false {',
     'TestInventorySelectedScopeAndIdentity', 1),
    ('selected-resource', '|| !selected {', '|| (!selected && false) {',
     'TestInventorySelectedScopeAndIdentity', 1),
    ('embedded-subscription', 'if !inventoryIdentity(resource.ID, resource.SubscriptionID) {',
     'if !inventoryIdentity(resource.ID, resource.SubscriptionID) && false {',
     'TestInventorySelectedScopeAndIdentity', 1),
    ('row-work-limit', 'if len(resources) > MaxInventoryRows {',
     'if len(resources) > MaxInventoryRows && false {', 'TestInventoryWorkAndTextLimits', 1),
    ('decoded-text-limit', 'if !auxiliaryText(label, &decodedBytes) {',
     'if !auxiliaryText(label, &decodedBytes) && false {', 'TestInventoryWorkAndTextLimits', 1),
    ('vmss-product-limit', 'capacity > MaxAuxCount/int64(sku.VCPUs)',
     'capacity > MaxAuxCount/int64(sku.VCPUs) && false', 'TestInventorySourceCapacityAndBoundaries', 1),
    ('masking-resource-id', 'redact.SubscriptionIDInResourceID(resource.ID, mask)',
     'redact.SubscriptionIDInResourceID(resource.ID, false)', 'TestInventoryCapturedCellsAndComposition', 1),
    ('sheet-composition', 'SheetName: "Region Inventory"', 'SheetName: "Inventory"',
     'TestInventoryCapturedCellsAndComposition', 1),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-inventory-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Inventory baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / FILE
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Inventory mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout:
            raise SystemExit('Inventory mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored inventory baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region inventory mutation rejected:', name)
