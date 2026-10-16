#!/usr/bin/env python3
"""Verify Section 3 tests reject compiling broken production variants."""

from pathlib import Path
import subprocess


ROOT = Path(__file__).resolve().parents[2]
CASES = [
    (
        "sarif-published-descriptor",
        "internal/renderers/sarif/sarif.go",
        "./internal/renderers/sarif",
        "^TestRuleReferencesUsePublishedDescriptorIdentity$",
        "\t\tRuleID: publishedRuleID,",
        "\t\tRuleID: finding.RecommendationID,",
        "SARIF result does not reference its published descriptor",
    ),
]


def run_go(package, pattern):
    """Return a bounded matched-test result with explicit diagnostics."""
    completed = subprocess.run(
        ["go", "test", "-count=1", "-timeout=90s", "-run", pattern, package],
        cwd=ROOT,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="backslashreplace",
        timeout=180,
        check=False,
    )
    return completed.returncode, completed.stdout + completed.stderr


def main():
    """Check baseline, reject compiling mutation and restore exact bytes."""
    for case in CASES:
        name, relative, package, pattern, before, after, expected = case
        path = ROOT / relative
        original = path.read_bytes()
        source = original.decode("utf-8").replace("\r\n", "\n")
        if source.count(before) != 1:
            raise RuntimeError(f"{name}: mutation anchor must occur once")
        code, output = run_go(package, pattern)
        if code != 0:
            raise RuntimeError(f"{name}: healthy baseline failed\n{output}")
        try:
            path.write_bytes(source.replace(before, after, 1).encode("utf-8"))
            code, output = run_go(package, pattern)
            if (
                code == 0
                or expected not in output
                or "[build failed]" in output
            ):
                raise RuntimeError(
                    f"{name}: compiling regression not observed\n{output}"
                )
        finally:
            path.write_bytes(original)
        if path.read_bytes() != original:
            raise RuntimeError(f"{name}: production bytes not restored")
        code, output = run_go(package, pattern)
        if code != 0:
            raise RuntimeError(f"{name}: restored baseline failed\n{output}")
        print(
            f"Section 3 compiling mutation rejected: {name}; "
            "restored baseline passed"
        )


if __name__ == "__main__":
    main()
