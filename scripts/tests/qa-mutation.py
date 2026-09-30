"""Confirm selected comparator tests reject faults in an isolated standard-library module."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILES = ("compare.go", "types.go", "compare_guard_test.go")
MUTATIONS = (
    ("forced-pass", "compare.go", "\treturn report", "\treport.Equivalent = true\n\treturn report"),
    ("ignored-field", "compare.go", "if left == right {", "if true {"),
    ("ignored-coverage", "compare.go", "reference.Enabled != target.Enabled,", "false,"),
    ("ignored-precondition", "compare.go", "if !target.Comparable {", "if false {"),
)


def run(directory):
    return subprocess.run(
        ["go", "test", "-count=1", "-run", "^TestCompareGuards$", "."],
        cwd=directory, capture_output=True, text=True, timeout=60,
    )


with tempfile.TemporaryDirectory(prefix="cloud-assess-mutations-") as temporary:
    directory = Path(temporary)
    (directory / "go.mod").write_text("module comparatorfixture\n\ngo 1.26.8\n")
    for name in FILES:
        shutil.copyfile(ROOT / "internal/equivalence" / name, directory / name)
    baseline = run(directory)
    if baseline.returncode != 0:
        raise SystemExit("Mutation baseline failed:\n" + baseline.stdout + baseline.stderr)
    for name, filename, before, after in MUTATIONS:
        path = directory / filename
        original = path.read_text()
        if original.count(before) != 1:
            raise SystemExit("Mutation anchor changed: " + name)
        path.write_text(original.replace(before, after))
        outcome = run(directory)
        path.write_text(original)
        if outcome.returncode == 0 or "--- FAIL: TestCompareGuards" not in outcome.stdout:
            raise SystemExit("Mutation survived or failed outside assertions: " + name + "\n" + outcome.stdout + outcome.stderr)
        print("Comparator mutation rejected:", name)
