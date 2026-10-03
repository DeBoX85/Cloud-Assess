#!/usr/bin/env python3
"""Require named AI boundary tests to reject compiling faults in an isolated copy."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ("foreign-continuation", "scanner.go", "!strings.EqualFold(u.Host, s.origin.Host)", "false", "TestDeploymentContinuationBeforeAuthentication"),
    ("foreign-metrics", "decoder.go", 'scope[key] == "" || reported[key]', "reported[key]", "TestMetricCoverageCorrelationAndErrors"),
    ("discard-metrics-on-enrichment-limit", "scanner.go", "projectionErr != nil && len(enrichment) != 0", "false", "TestRequestEnrichmentLimitRetainsAllValidMetrics"),
)


def run(directory, test):
    return subprocess.run(
        ["go", "test", "-count=1", "-run", "^" + test + "$", "./internal/plugins/aigov"],
        cwd=directory, capture_output=True, text=True, timeout=120,
    )


with tempfile.TemporaryDirectory(prefix="cloud-assess-ai-mutations-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("assessment", "azure", "throttling", "plugins/aigov"):
        shutil.copytree(ROOT / "internal" / name, directory / "internal" / name)
    for name, filename, before, after, test in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode != 0:
            raise SystemExit("AI mutation baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        path = directory / "internal/plugins/aigov" / filename
        original = path.read_text(encoding="utf-8")
        if original.count(before) != 1:
            raise SystemExit("AI mutation anchor changed: " + name)
        try:
            path.write_text(original.replace(before, after), encoding="utf-8")
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding="utf-8")
        if outcome.returncode == 0 or "--- FAIL: " + test not in outcome.stdout:
            raise SystemExit("AI mutation survived or failed outside named assertions: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode != 0:
            raise SystemExit("Restored AI mutation fixture failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("AI request mutation rejected:", name)
