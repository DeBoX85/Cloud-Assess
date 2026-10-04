#!/usr/bin/env python3
"""Require compiling named source-oracle assertions on altered synthetic fixtures."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
TEST = "TestCostEnrichmentCapturedSemantics"


def run(directory):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + TEST + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-cost-enrichment-capture-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    path = directory / "internal/plugins/region/testdata/source-cost-enrichment-outputs.json"
    original = path.read_bytes()
    for name in ("weighted", "free-denominator", "zero-weight", "duplicate-history", "stale-state", "unbound-identity"):
        baseline = run(directory)
        if baseline.returncode:
            raise SystemExit("Cost enrichment capture baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "weighted":
            modified["weighted"][0]["AvgCostDifference"] = 75
        elif name == "free-denominator":
            modified["near-zero-both"][0]["AvgCostDifference"] = 100
        elif name == "zero-weight":
            modified["zero-weight"][0]["AvgCostDifference"] = -50
        elif name == "duplicate-history":
            modified["duplicate-last-one"][0]["AvgCostDifference"] = 90
        elif name == "stale-state":
            modified["nil-shared"][0]["HasCostData"] = False
        else:
            modified["unbound-subscription"][1]["SubscriptionID"] = "11111111-1111-1111-1111-111111111111"
        try:
            path.write_text(json.dumps(modified), encoding="utf-8")
            outcome = run(directory)
        finally:
            path.write_bytes(original)
        if outcome.returncode == 0 or "--- FAIL: " + TEST not in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr or "[build failed]" in outcome.stdout:
            raise SystemExit("Cost enrichment capture mutation survived or failed outside named assertions: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory)
        if restored.returncode:
            raise SystemExit("Restored cost enrichment capture failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Cost enrichment source oracle mutation rejected:", name)
