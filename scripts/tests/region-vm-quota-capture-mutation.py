#!/usr/bin/env python3
"""Require compiling full VM SDK source-oracle failures and restored baselines."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
TEST = "TestVMQuotaCapturedSemantics"


def run(directory):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + TEST + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-vm-quota-capture-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    path = directory / "internal/plugins/region/testdata/source-vm-quota-outputs.json"
    original = path.read_bytes()
    for name in ("family-filter", "source-panic", "unsupported-error", "nil-empty", "sdk-request"):
        baseline = run(directory)
        if baseline.returncode:
            raise SystemExit("VM quota capture baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "family-filter":
            modified["selection"]["Entries"][1]["ResourceName"] = "cores"
        elif name == "source-panic":
            modified["missing-current"]["Panicked"] = False
        elif name == "unsupported-error":
            modified["http400"]["Failed"] = False
        elif name == "nil-empty":
            modified["empty"]["Entries"] = []
        else:
            modified["paged"]["Requests"][1] = modified["paged"]["Requests"][1].replace("2024-11-01", "2022-07-01")
        try:
            path.write_text(json.dumps(modified), encoding="utf-8")
            outcome = run(directory)
        finally:
            path.write_bytes(original)
        if outcome.returncode == 0 or "--- FAIL: " + TEST not in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr or "[build failed]" in outcome.stdout:
            raise SystemExit("VM quota capture mutation survived or failed outside named assertion: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory)
        if restored.returncode:
            raise SystemExit("Restored VM quota capture failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("VM quota source oracle mutation rejected:", name)
