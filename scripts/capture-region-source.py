#!/usr/bin/env python3
"""Capture unchanged pinned AZQR pure region helpers in an isolated checkout."""
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]
SOURCE = "8e4f0577f3615e6c9014c031bcad079f235369cc"
APRL = "60eaddda76541f6adbc1c5ffa686829807e55e29"
SOURCE_TREE = "17d93b20c303f90f7843036be82f0dc32f3260f1"
HARNESS = ROOT / "internal/plugins/region/testdata/aux_capture_test.go.txt"


def run(args, directory, timeout=240, capture=False):
    return subprocess.run(args, cwd=directory, check=True, timeout=timeout,
                          text=True, stdout=subprocess.PIPE if capture else None)


with tempfile.TemporaryDirectory(prefix="cloud-assess-region-source-") as temporary:
    directory = Path(temporary) / "reference"
    directory.mkdir()
    run(["git", "init", "--quiet"], directory)
    run(["git", "remote", "add", "origin", "https://github.com/DeBoX85/azqr.git"], directory)
    run(["git", "fetch", "--quiet", "--depth=1", "origin", SOURCE], directory)
    run(["git", "checkout", "--quiet", "--detach", "FETCH_HEAD"], directory)
    actual = run(["git", "show", "-s", "--format=%H %T"], directory, capture=True).stdout.strip()
    if actual != SOURCE + " " + SOURCE_TREE:
        raise SystemExit("Pinned source identity mismatch")
    run(["git", "submodule", "update", "--init", "--recursive"], directory)
    actual_aprl = run(["git", "rev-parse", "HEAD"], directory / "internal/graph/aprl", capture=True).stdout.strip()
    if actual_aprl != APRL:
        raise SystemExit("Actual APRL mismatch")
    # Pure helper test only; do not alter production or SDK files.
    target = directory / "internal/scanners/plugins/region/output/aux_capture_test.go"
    shutil.copyfile(HARNESS, target)
    output = Path(temporary) / "output"
    previous = os.environ.get("REGION_AUX_CAPTURE_OUTPUT")
    os.environ["REGION_AUX_CAPTURE_OUTPUT"] = str(output)
    try:
        run(["go", "test", "-mod=readonly", "-count=1", "-timeout=90s", "-v",
             "-run", "^TestRegionAuxiliaryCapture$", "./internal/scanners/plugins/region/output"],
            directory, timeout=600)
    finally:
        if previous is None:
            os.environ.pop("REGION_AUX_CAPTURE_OUTPUT", None)
        else:
            os.environ["REGION_AUX_CAPTURE_OUTPUT"] = previous
    run(["git", "diff", "--exit-code"], directory)
    run(["git", "diff", "--exit-code"], directory / "internal/graph/aprl")
    untracked = run(["git", "ls-files", "--others", "--exclude-standard"], directory, capture=True).stdout.splitlines()
    if untracked != ["internal/scanners/plugins/region/output/aux_capture_test.go"]:
        raise SystemExit("Unexpected source-copy changes")
    print("REGION_CAPTURE_PROVENANCE " + json.dumps({"source": SOURCE, "tree": SOURCE_TREE, "aprl": APRL}))
    for name in ("source-aux-inputs.json", "source-aux-outputs.json"):
        raw = (output / name).read_bytes()
        content = raw.decode("utf-8")
        pieces = [content[i:i + 3000] for i in range(0, len(content), 3000)]
        for index, piece in enumerate(pieces):
            print("REGION_CAPTURE_JSON " + json.dumps({
                "name": name, "sha256": hashlib.sha256(raw).hexdigest(),
                "index": index, "count": len(pieces), "content": piece}, ensure_ascii=True))
