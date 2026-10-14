#!/usr/bin/env python3
"""Verify Section 2 tests reject compiling broken production variants."""

from pathlib import Path
import subprocess


ROOT = Path(__file__).resolve().parents[2]
CASES = [
    (
        "runner-caller-terminal",
        "internal/stages/runner.go",
        "./internal/stages",
        "^TestRunnerCallerTerminalContextWinsOverPlainTaskError$",
        "\t\t\tif contextErr := ctx.Err(); contextErr != nil {\n"
        "\t\t\t\terr = contextErr\n\t\t\t}",
        "\t\t\tif err == nil {\n\t\t\t\terr = ctx.Err()\n\t\t\t}",
        "caller terminal context concealed by task error",
    ),
    (
        "scanner-nested-ownership",
        "internal/scanners/registry.go",
        "./internal/scanners",
        "^TestRegistryAccessorsOwnNestedResourceTypes$",
        "\tfor i := range out {\n"
        "\t\tout[i].ResourceTypes = append([]string(nil), "
        "services[i].ResourceTypes...)\n\t}\n",
        "",
        "caller mutated global scanner resource-type scope",
    ),
    (
        "cost-caller-terminal",
        "internal/cost/scanner.go",
        "./internal/cost",
        "^TestCostCallerCancellationCannotBecomeSuccessfulDataset$",
        "\tif err := ctx.Err(); err != nil {\n"
        "\t\treturn result, err\n\t}\n\tif firstErr != nil {",
        "\tif firstErr != nil {",
        "Cost caller cancellation lost",
    ),
    (
        "diagnostics-caller-terminal",
        "internal/diagnostics/scanner.go",
        "./internal/diagnostics",
        "^TestDiagnosticsCallerCancellationCannotCreate"
        "MissingSettingsFindings$",
        "\tif err := ctx.Err(); err != nil {\n"
        "\t\treturn enabled, warnings, err\n\t}\n"
        "\tif firstErr != nil {",
        "\tif firstErr != nil {",
        "Diagnostics caller cancellation became assessment data",
    ),
    (
        "diagnostics-single-json",
        "internal/diagnostics/scanner.go",
        "./internal/diagnostics",
        "^TestDiagnosticsRejectsTrailingBatchResponseContent$",
        "\tvar trailing any\n"
        "\tif err := decoder.Decode(&trailing); err != io.EOF {\n"
        "\t\tif err != nil {\n"
        "\t\t\treturn nil, nil, fmt.Errorf("
        "\"read diagnostic settings batch trailing content: %w\", err)\n"
        "\t\t}\n"
        "\t\treturn nil, nil, fmt.Errorf("
        "\"diagnostic settings batch response has trailing content\")\n"
        "\t}\n",
        "",
        "trailing batch response content became assessment findings",
    ),
    (
        "diagnostics-batch-scope",
        "internal/diagnostics/scanner.go",
        "./internal/diagnostics",
        "^TestDiagnosticsRejectsSettingsOutsideRequestedBatch$",
        "\t\t\tif !requested[resourceID] {\n"
        "\t\t\t\treturn nil, warnings, fmt.Errorf("
        "\"diagnostic settings response contains a resource "
        "outside the requested batch\")\n\t\t\t}\n",
        "",
        "foreign diagnostic setting accepted for requested batch",
    ),
    (
        "cost-completed-status",
        "internal/cost/scanner.go",
        "./internal/cost",
        "^TestCostHTTPRequiresDocumentedQueryCompletionStatus$",
        "\tif response.StatusCode != http.StatusOK && "
        "response.StatusCode != http.StatusNoContent {\n"
        "\t\treturn nil, nil, fmt.Errorf("
        "\"Cost Management response status %d is not a completed query\", "
        "response.StatusCode)\n\t}\n",
        "",
        "undocumented Cost status admitted as completed query",
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
            f"Section 2 compiling mutation rejected: {name}; "
            "restored baseline passed"
        )


if __name__ == "__main__":
    main()
