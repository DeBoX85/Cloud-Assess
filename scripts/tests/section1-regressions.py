#!/usr/bin/env python3
"""Verify Section 1 regressions reject compiling broken production variants."""

from pathlib import Path
import subprocess


ROOT = Path(__file__).resolve().parents[2]
CASES = [
    (
        "arg-completed-status",
        ROOT / "internal/arg/http_transport.go",
        "./internal/arg",
        "^TestARGHTTPTransportRequiresDocumentedCompletedQueryStatus$",
        "\tif response.StatusCode != http.StatusOK {\n"
        "\t\treturn nil, fmt.Errorf("
        "\"ARG response status %d is not a completed query\", "
        "response.StatusCode)\n"
        "\t}\n",
        "",
        "undocumented ARG status admitted as complete query",
    ),
    (
        "parent-terminal-priority",
        ROOT / "internal/arg/executor.go",
        "./internal/arg",
        "^TestExecuteRecommendationsParentContextWinsOverPlainQueryError$",
        "\tif err := parentContext.Err(); err != nil {\n"
        "\t\treturn nil, warnings, fmt.Errorf("
        "\"recommendation execution interrupted: %w\", err)\n"
        "\t}\n",
        "",
        "parent terminal priority lost",
    ),
    (
        "inventory-tied-order",
        ROOT / "internal/discovery/resources.go",
        "./internal/discovery",
        "^TestInventoryPagingRetainsAllResourcesWithTiedSubscriptionAndGroup$",
        " | order by subscriptionId, resourceGroup, id asc",
        " | order by subscriptionId, resourceGroup",
        "tied sort duplicated resource",
    ),
    (
        "inventory-nonunique-projection",
        ROOT / "internal/discovery/resources.go",
        "./internal/discovery",
        "^TestInventoryPagingRetainsAllResourcesWithTiedSubscriptionAndGroup$",
        "id=tostring(id)",
        "id=tostring(subscriptionId)",
        "tied sort duplicated resource",
    ),
]


def run_go(package, pattern):
    """Run matched Go tests with bounded execution and explicit diagnostics."""
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
    """Reject broken variants and restore exact production bytes each time."""
    for case in CASES:
        name, path, package, pattern, before, after, expected = case
        original = path.read_bytes()
        source = original.decode("utf-8").replace("\r\n", "\n")
        if source.count(before) != 1:
            raise RuntimeError(
                f"{name}: mutation anchor must occur exactly once"
            )
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
                    f"{name}: expected compiling regression not observed\n"
                    f"{output}"
                )
        finally:
            path.write_bytes(original)
        if path.read_bytes() != original:
            raise RuntimeError(f"{name}: production bytes were not restored")
        code, output = run_go(package, pattern)
        if code != 0:
            raise RuntimeError(f"{name}: restored baseline failed\n{output}")
        print(
            f"Section 1 compiling mutation rejected: {name}; "
            "restored baseline passed"
        )


if __name__ == "__main__":
    main()
