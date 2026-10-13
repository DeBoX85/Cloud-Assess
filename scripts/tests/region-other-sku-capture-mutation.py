#!/usr/bin/env python3
"""Compiling independent remaining SKU source controls and restored baselines."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ("source-blob", "TestOtherSKUSourceMetadata", "inputs"),
    ("api-version", "TestOtherSKUSourceMetadata", "inputs"),
    ("restriction-values", "TestOtherSKUSourceCells", "outputs"),
    ("zone-order", "TestOtherSKUSourceCells", "outputs"),
    ("global-location", "TestOtherSKUSourceCells", "outputs"),
    ("kind-collapse", "TestOtherSKUSourceCells", "outputs"),
    ("sql-parent-status", "TestOtherSKUSourceCells", "outputs"),
    ("sql-unicode-status", "TestOtherSKUSourceCells", "outputs"),
    ("mi-tier", "TestOtherSKUSourceCells", "outputs"),
    ("later-page-loss", "TestOtherSKUSourceCells", "outputs"),
    ("page-omission", "TestOtherSKUSourceRequests", "outputs"),
    ("method", "TestOtherSKUSourceRequests", "outputs"),
    ("body", "TestOtherSKUSourceRequests", "outputs"),
    ("raw-path", "TestOtherSKUSourceRequests", "outputs"),
    ("sql-next-ignored", "TestOtherSKUSourceRequests", "outputs"),
    ("sql-type-cache", "TestOtherSKUSourceCache", "outputs"),
    ("scope-cache", "TestOtherSKUSourceCache", "outputs"),
)


def run(directory, test):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + test + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-other-sku-control-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    for name, test, kind in MUTATIONS:
        path = directory / ("internal/plugins/region/testdata/source-other-sku-" + kind + ".json")
        original = path.read_bytes()
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit("Other SKU baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "source-blob":
            modified["Blobs"]["sql.go"] = "0" * 40
        elif name == "api-version":
            modified["Cases"]["cognitive-single"]["Pages"][0]["URL"] = modified["Cases"]["cognitive-single"]["Pages"][0]["URL"].replace("2024-10-01", "2023-01-01")
        elif name == "restriction-values":
            modified["Cases"]["storage-location-values-ignored"]["Entries"]["s"]["State"] = 0
        elif name == "zone-order":
            modified["Cases"]["storage-zone-order"]["Entries"]["s"]["BlockedZones"] = ["1", "2"]
        elif name == "global-location":
            modified["Cases"]["cognitive-global-excluded"]["Entries"] = {"s0": {"State": 0, "BlockedZones": None}}
        elif name == "kind-collapse":
            modified["Cases"]["cognitive-kinds-collapse"]["Entries"]["speech_s0"] = {"State": 0, "BlockedZones": None}
        elif name == "sql-parent-status":
            modified["Cases"]["sql-db-parent-status-ignored"]["Entries"]["s"]["State"] = 2
        elif name == "sql-unicode-status":
            modified["Cases"]["sql-db-status-matrix"]["Entries"]["unicodedisabled"]["State"] = 0
        elif name == "mi-tier":
            modified["Cases"]["sql-mi-tier-matrix"]["Entries"]["gp_gen5"] = modified["Cases"]["sql-mi-tier-matrix"]["Entries"].pop("ngpgen5_gen5")
        elif name == "later-page-loss":
            modified["Cases"]["storage-later-http403"]["Entries"] = modified["Cases"]["storage-single"]["Entries"]
        elif name == "page-omission":
            modified["Cases"]["cognitive-paged"]["Requests"] = modified["Cases"]["cognitive-paged"]["Requests"][:1]
        elif name == "method":
            modified["Cases"]["storage-single"]["Requests"][0]["Method"] = "POST"
        elif name == "body":
            modified["Cases"]["sql-db-status-matrix"]["Requests"][0]["EmptyBody"] = False
        elif name == "raw-path":
            modified["Cases"]["sql-mi-direct-display-region"]["Requests"][0]["RawPath"] = ""
        elif name == "sql-next-ignored":
            modified["Cases"]["sql-db-next-ignored"]["Requests"].append(modified["Cases"]["sql-db-next-ignored"]["Requests"][0])
        elif name == "sql-type-cache":
            modified["Cache"]["mi-separate-key"]["Entries"] = modified["Cache"]["db-first"]["Entries"]
        elif name == "scope-cache":
            modified["Cache"]["db-subscription-separate"]["Requests"] = modified["Cache"]["db-subscription-separate"]["Requests"][:2]
        else:
            raise SystemExit("Unknown other SKU control")
        encoded = json.dumps(modified, ensure_ascii=True).encode("utf-8")
        if encoded == original:
            raise SystemExit("Other SKU control made no change")
        try:
            path.write_bytes(encoded)
            outcome = run(directory, test)
        finally:
            path.write_bytes(original)
        combined = outcome.stdout + outcome.stderr
        if outcome.returncode == 0 or "--- FAIL: " + test not in outcome.stdout or "panic:" in combined or "[build failed]" in combined or "undefined:" in combined:
            raise SystemExit("Other SKU mutation survived or failed outside named assertion: " + name + "\n" + combined)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit("Restored other SKU baseline failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Other SKU source mutation rejected:", name)
