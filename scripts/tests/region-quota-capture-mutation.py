#!/usr/bin/env python3
"""Require compiled full quota source-oracle assertion failures and restoration."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
TEST = "TestQuotaCapturedSemantics"


def run(directory):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + TEST + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-quota-capture-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    path = directory / "internal/plugins/region/testdata/source-quota-outputs.json"
    original = path.read_bytes()
    for name in ("threshold", "filter", "unsupported-prefix", "denial", "request"):
        baseline = run(directory)
        if baseline.returncode:
            raise SystemExit("Quota capture baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "threshold":
            modified["arithmetic"]["Entries"][0]["IsNearLimit"] = True
        elif name == "filter":
            modified["web-filter"]["Entries"][0]["ResourceName"] = "CustomDomains"
        elif name == "unsupported-prefix":
            modified["later-http404"]["Entries"] = modified["paged"]["Entries"][:1]
        elif name == "denial":
            modified["http403"]["Failed"] = False
        else:
            modified["sql-filter"]["Requests"][0] += "&widerScope=true"
        try:
            path.write_text(json.dumps(modified), encoding="utf-8")
            outcome = run(directory)
        finally:
            path.write_bytes(original)
        if outcome.returncode == 0 or "--- FAIL: " + TEST not in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr or "[build failed]" in outcome.stdout:
            raise SystemExit("Quota capture mutation survived or failed outside named assertion: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory)
        if restored.returncode:
            raise SystemExit("Restored quota capture failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Quota source oracle mutation rejected:", name)
