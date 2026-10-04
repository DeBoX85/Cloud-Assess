#!/usr/bin/env python3
"""Require compiling named source-oracle assertions on altered synthetic fixtures."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
TEST = "TestLatencyCapturedSemantics"


def run(directory):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + TEST + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-latency-capture-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    for name in ("forward-preference", "directional-absence", "default-average", "dataset-topology"):
        file = "source-latency-data.json" if name == "dataset-topology" else "source-latency-outputs.json"
        path = directory / "internal/plugins/region/testdata" / file
        original = path.read_bytes()
        baseline = run(directory)
        if baseline.returncode:
            raise SystemExit("Latency capture baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "forward-preference":
            modified["synthetic"]["comparisons"][0]["AvgLatencyMs"] = 115
        elif name == "directional-absence":
            modified["one-way-cluster"]["comparisons"][1]["AvgLatencyMs"] = 80
        elif name == "default-average":
            modified["default"]["averages"]["europe:americas"] = 0
        else:
            del modified["matrix"]["eastus"]["westeurope"]
        try:
            path.write_text(json.dumps(modified), encoding="utf-8")
            outcome = run(directory)
        finally:
            path.write_bytes(original)
        if outcome.returncode == 0 or "--- FAIL: " + TEST not in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr or "[build failed]" in outcome.stdout:
            raise SystemExit("Latency capture mutation survived or failed outside named assertions: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory)
        if restored.returncode:
            raise SystemExit("Restored latency capture failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Latency source oracle mutation rejected:", name)
