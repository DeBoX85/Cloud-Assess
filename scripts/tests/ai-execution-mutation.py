#!/usr/bin/env python3
"""Require public AI guards to reject compiling faults, with restored baselines."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ('public-pre-auth', 'cmd/cloud-assess/command.go', 'if slices.Contains(flags.internalPlugins, aigov.Name) {',
     'if slices.Contains(flags.internalPlugins, aigov.Name) && false {', './cmd/cloud-assess', 'TestAICloudPreflightBeforeFactoriesAndScope'),
    ('discovery-health', 'internal/plugins/aigov/execution.go', 'if discoveryErr != nil || discovered.Health.Status == assessment.StageFailed {',
     'if false {', './internal/plugins/aigov', 'TestExecutionDiscoveryFailureSurvivesEnrichment'),
    ('recorded-tag-snapshot', 'internal/config/filters.go', 'out.resourceScope = maps.Clone(f.resourceScope)',
     'out.resourceScope = nil', './internal/orchestration', 'TestAICoordinatorRecordedTagScopeOwned'),
    ('projected-scope-identity', 'internal/plugins/aigov_projection.go', 'if !ok || row.Cells[0] != name {',
     'if !ok && false || row.Cells[0] != name && false {', './internal/plugins', 'TestAIProjectionRejectsForeignIdentityAndOwnsCells'),
)


def run(directory, package, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$', package],
                          cwd=directory, capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-ai-execution-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, filename, before, after, package, test in MUTATIONS:
        baseline = run(directory, package, test)
        if baseline.returncode:
            raise SystemExit('AI execution baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / filename
        original = path.read_text(encoding='utf-8')
        if original.count(before) != 1:
            raise SystemExit('AI execution anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, package, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout:
            raise SystemExit('AI execution mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, package, test)
        if restored.returncode:
            raise SystemExit('Restored AI execution failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('AI execution mutation rejected:', name)
