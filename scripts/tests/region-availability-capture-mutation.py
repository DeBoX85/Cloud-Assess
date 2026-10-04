#!/usr/bin/env python3
"""Require compiling named source-oracle assertions on altered synthetic fixtures."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = "internal/plugins/region/testdata/source-availability-outputs.json"
TEST = "TestAvailabilityCapturedSemantics"


def run(directory):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + TEST + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-availability-capture-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    path = directory / FILE
    original = path.read_bytes()
    for name in ("unique-types", "all-unknown-percent", "blocked-zone-order", "global-provider"):
        baseline = run(directory)
        if baseline.returncode:
            raise SystemExit("Availability capture baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "unique-types":
            modified["cases"]["resource-mixed-no-sku"]["SourceResourceTypeCount"] = 909
        elif name == "all-unknown-percent":
            modified["cases"]["sku-all-unknown"]["SKUAvailabilityPercent"] = 0
        elif name == "blocked-zone-order":
            modified["cases"]["sku-all-states"]["ZoneRestrictedSKUs"][0] = "microsoft.capture/available:zones (zones blocked: 1,3)"
        else:
            modified["provider_lookup"][3]["Available"] = False
        try:
            path.write_text(json.dumps(modified), encoding="utf-8")
            outcome = run(directory)
        finally:
            path.write_bytes(original)
        if outcome.returncode == 0 or "--- FAIL: " + TEST not in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr or "[build failed]" in outcome.stdout:
            raise SystemExit("Availability capture mutation survived or failed outside named assertions: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory)
        if restored.returncode:
            raise SystemExit("Restored availability capture failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Availability source oracle mutation rejected:", name)
