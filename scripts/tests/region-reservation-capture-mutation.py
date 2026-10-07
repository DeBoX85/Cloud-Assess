#!/usr/bin/env python3
"""Require compiling complete reservation source-oracle failures and restoration."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
TEST = "TestReservationCapturedSemantics"


def run(directory):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + TEST + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-reservation-capture-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    path = directory / "internal/plugins/region/testdata/source-reservation-outputs.json"
    original = path.read_bytes()
    for name in ("status", "source-panic", "partial-success", "inner-cancellation",
                 "identity", "location", "sdk-request"):
        baseline = run(directory)
        if baseline.returncode:
            raise SystemExit("Reservation source baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "status":
            modified["negative"]["Entries"][0]["Status"] = "Over-Allocated"
        elif name == "source-panic":
            modified["null-group"]["Panicked"] = False
        elif name == "partial-success":
            modified["later-list-http403"]["Entries"] = None
        elif name == "inner-cancellation":
            modified["cancel-get"]["Failed"] = True
        elif name == "identity":
            modified["identity-conflict"]["Entries"][0]["ReservationName"] = "r"
        elif name == "location":
            modified["available"]["Entries"][0]["Location"] = "westus"
        else:
            modified["available"]["Requests"][2] = modified["available"]["Requests"][2].replace("%24expand=instanceView&", "")
        try:
            path.write_text(json.dumps(modified), encoding="utf-8")
            outcome = run(directory)
        finally:
            path.write_bytes(original)
        log = outcome.stdout + outcome.stderr
        if outcome.returncode == 0 or "--- FAIL: " + TEST not in outcome.stdout or "panic:" in log or "[build failed]" in log:
            raise SystemExit("Reservation source mutation survived or failed outside named assertion: " + name + "\n" + log)
        restored = run(directory)
        if restored.returncode:
            raise SystemExit("Restored reservation source oracle failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Reservation source oracle mutation rejected:", name)
