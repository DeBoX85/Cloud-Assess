#!/usr/bin/env python3
"""Require SR-F12 to reject compiling old overwrite behavior and restore bytes."""

from pathlib import Path
import subprocess


ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / "internal/equivalence/projection.go"
PATTERN = "^TestComparison(RejectsConflictingDuplicateEvidence|Preserves.*)$"
ANCHOR = "if existing, exists := dataset.Records[key]; exists {"


def run(pattern):
    completed = subprocess.run(
        ["go", "test", "-count=1", "-timeout=90s", "-run", pattern,
         "./tools/equivalence"],
        cwd=ROOT, capture_output=True, text=True, encoding="utf-8",
        errors="backslashreplace", timeout=180, check=False,
    )
    return completed.returncode, completed.stdout + completed.stderr


def baseline():
    code, output = run(PATTERN)
    if code != 0:
        raise RuntimeError("SR-F12 baseline failed\n" + output)


def main():
    original = SOURCE.read_bytes()
    source = original.decode("utf-8").replace("\r\n", "\n")
    if source.count(ANCHOR) != 1:
        raise RuntimeError("SR-F12 mutation anchor must occur exactly once")
    baseline()
    try:
        # Disabling the whole duplicate branch restores the original last-row
        # overwrite, while its imports remain used and compilation stays valid.
        SOURCE.write_bytes(source.replace(
            ANCHOR, "if existing, exists := dataset.Records[key]; exists && false {",
            1).encode("utf-8"))
        code, output = run("^$")
        if code != 0:
            raise RuntimeError("SR-F12 reversal did not compile\n" + output)
        code, output = run(PATTERN)
        if (code == 0
                or "--- FAIL: TestComparisonRejectsConflictingDuplicateEvidence" not in output
                or "conflicting duplicate evidence" not in output
                or "[build failed]" in output
                or "panic:" in output):
            raise RuntimeError("SR-F12 named compiling assertion absent\n" + output)
    finally:
        SOURCE.write_bytes(original)
    if SOURCE.read_bytes() != original:
        raise RuntimeError("SR-F12 source bytes were not restored")
    baseline()
    print("Section 4 compiling mutation rejected: equivalence-duplicate-conflict; "
          "restored baseline passed")


if __name__ == "__main__":
    main()
