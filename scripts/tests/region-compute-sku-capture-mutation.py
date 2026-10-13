#!/usr/bin/env python3
"""Compiling independent Compute source-cell controls with restored baselines."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ("source-blob", "TestComputeSKUSourceMetadata", "inputs"),
    ("api-version", "TestComputeSKUSourceMetadata", "inputs"),
    ("location-precedence", "TestComputeSKUSourceCells", "outputs"),
    ("zone-order", "TestComputeSKUSourceCells", "outputs"),
    ("lost-page-error", "TestComputeSKUSourceCells", "outputs"),
    ("unicode-alias", "TestComputeSKUSourceCells", "outputs"),
    ("page-omission", "TestComputeSKUSourceRequests", "outputs"),
    ("unsafe-continuation", "TestComputeSKUSourceRequests", "outputs"),
    ("method", "TestComputeSKUSourceRequests", "outputs"),
    ("body", "TestComputeSKUSourceRequests", "outputs"),
    ("subscription-cache", "TestComputeSKUSourceCache", "outputs"),
    ("zone-sharing", "TestComputeSKUSourceCache", "outputs"),
    ("retry", "TestComputeSKUSourceCache", "outputs"),
    ("concurrent-key", "TestComputeSKUSourceCache", "outputs"),
)


def run(directory, test):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + test + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-compute-source-control-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    for name, test, kind in MUTATIONS:
        path = directory / ("internal/plugins/region/testdata/source-compute-sku-" + kind + ".json")
        original = path.read_bytes()
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit("Compute source baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "source-blob":
            modified["Blobs"]["compute.go"] = "0" * 40
        elif name == "api-version":
            modified["Cases"]["vm-single"]["Pages"][0]["URL"] = modified["Cases"]["vm-single"]["Pages"][0]["URL"].replace("2023-03-01", "2021-07-01")
        elif name == "location-precedence":
            modified["Cases"]["disk-location-precedence"]["Entries"]["d"]["State"] = 3
        elif name == "zone-order":
            modified["Cases"]["disk-zone-order-duplicates"]["Entries"]["d"]["BlockedZones"] = ["1", "2"]
        elif name == "lost-page-error":
            modified["Cases"]["disk-later-http403"]["Entries"] = modified["Cases"]["disk-next-foreign"]["Entries"]
        elif name == "unicode-alias":
            modified["Cases"]["disk-unicode-field-alias"]["Entries"]["d"]["State"] = 0
        elif name == "page-omission":
            modified["Cases"]["vm-paged"]["Requests"] = modified["Cases"]["vm-paged"]["Requests"][:1]
        elif name == "unsafe-continuation":
            modified["Cases"]["vm-next-foreign"]["Requests"] = modified["Cases"]["vm-next-foreign"]["Requests"][:1]
        elif name == "method":
            modified["Cases"]["vm-single"]["Requests"][0]["Method"] = "POST"
        elif name == "body":
            modified["Cases"]["vm-single"]["Requests"][0]["EmptyBody"] = False
        elif name == "subscription-cache":
            modified["Cache"]["subscription-separate"]["Entries"] = modified["Cache"]["vmss-separate"]["Entries"]
        elif name == "zone-sharing":
            modified["ZoneAfter"] = modified["ZoneBefore"]
        elif name == "retry":
            modified["Retry"][1]["Failed"] = True
        elif name == "concurrent-key":
            modified["ConcurrentSeparateKeys"]["west"] = modified["ConcurrentSeparateKeys"]["east"]
        else:
            raise SystemExit("Unknown Compute source control")
        encoded = json.dumps(modified, ensure_ascii=True).encode("utf-8")
        if encoded == original:
            raise SystemExit("Compute source control made no change: " + name)
        try:
            path.write_bytes(encoded)
            outcome = run(directory, test)
        finally:
            path.write_bytes(original)
        combined = outcome.stdout + outcome.stderr
        if outcome.returncode == 0 or "--- FAIL: " + test not in outcome.stdout or "panic:" in combined or "[build failed]" in combined or "undefined:" in combined:
            raise SystemExit("Compute source mutation survived or failed outside named assertion: " + name + "\n" + combined)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit("Restored Compute source baseline failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Compute source mutation rejected:", name)
