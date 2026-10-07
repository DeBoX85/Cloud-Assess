#!/usr/bin/env python3
"""Capture unchanged pinned AZQR helpers and synthetic quota/reservation collectors."""
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
CAPTURES = (
    ("aux_capture_test.go.txt", "output/aux_capture_test.go",
     "REGION_AUX_CAPTURE_OUTPUT", "^TestRegionAuxiliaryCapture$",
     "./internal/scanners/plugins/region/output"),
    ("inventory_capture_test.go.txt", "inventory_capture_test.go",
     "REGION_INVENTORY_CAPTURE_OUTPUT", "^TestRegionInventoryCalculationCapture$",
     "./internal/scanners/plugins/region"),
    ("availability_capture_test.go.txt", "availability/availability_capture_test.go",
     "REGION_AVAILABILITY_CAPTURE_OUTPUT", "^TestRegionAvailabilityCalculationCapture$",
     "./internal/scanners/plugins/region/availability"),
    ("latency_capture_test.go.txt", "latency/latency_capture_test.go",
     "REGION_LATENCY_CAPTURE_OUTPUT", "^TestRegionLatencyCapture$",
     "./internal/scanners/plugins/region/latency"),
    ("cost_enrichment_capture_test.go.txt", "cost/cost_enrichment_capture_test.go",
     "REGION_COST_ENRICHMENT_CAPTURE_OUTPUT", "^TestRegionCostEnrichmentCapture$",
     "./internal/scanners/plugins/region/cost"),
    ("quota_capture_test.go.txt", "quota/quota_capture_test.go",
     "REGION_QUOTA_CAPTURE_OUTPUT", "^TestRegionQuotaCapture$",
     "./internal/scanners/plugins/region/quota"),
    ("vm_quota_capture_test.go.txt", "quota/vm_quota_capture_test.go",
     "REGION_VM_QUOTA_CAPTURE_OUTPUT", "^TestRegionVMQuotaCapture$",
     "./internal/scanners/plugins/region/quota"),
    ("reservation_capture_test.go.txt", "crg/reservation_capture_test.go",
     "REGION_RESERVATION_CAPTURE_OUTPUT", "^TestRegionReservationCapture$",
     "./internal/scanners/plugins/region/crg"),
)
FILES = ("source-aux-inputs.json", "source-aux-outputs.json",
         "source-inventory-inputs.json", "source-inventory-outputs.json",
         "source-availability-inputs.json", "source-availability-outputs.json",
         "source-latency-inputs.json", "source-latency-outputs.json",
         "source-latency-data.json",
         "source-cost-enrichment-inputs.json", "source-cost-enrichment-outputs.json",
         "source-quota-inputs.json", "source-quota-outputs.json",
         "source-vm-quota-inputs.json", "source-vm-quota-outputs.json",
         "source-reservation-inputs.json", "source-reservation-outputs.json")


def run(args, directory, timeout=240, capture=False):
    return subprocess.run(args, cwd=directory, check=True, timeout=timeout,
                          text=True, stdout=subprocess.PIPE if capture else None)


def verify_capture_bytes(raw, retained):
    if raw != retained.read_bytes():
        raise SystemExit("Retained source capture byte mismatch: " + retained.name)


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
    output = Path(temporary) / "output"
    expected_untracked = []
    # Pure helpers and injected collector transports; production/SDK/modules unchanged.
    for harness, relative, variable, test, package in CAPTURES:
        target = directory / "internal/scanners/plugins/region" / relative
        shutil.copyfile(ROOT / "internal/plugins/region/testdata" / harness, target)
        expected_untracked.append(str(target.relative_to(directory)).replace("\\", "/"))
        previous = os.environ.get(variable)
        os.environ[variable] = str(output)
        try:
            run(["go", "test", "-mod=readonly", "-count=1", "-timeout=90s", "-v",
                 "-run", test, package], directory, timeout=600)
        finally:
            if previous is None:
                os.environ.pop(variable, None)
            else:
                os.environ[variable] = previous
    run(["git", "diff", "--exit-code"], directory)
    run(["git", "diff", "--exit-code"], directory / "internal/graph/aprl")
    untracked = run(["git", "ls-files", "--others", "--exclude-standard"], directory, capture=True).stdout.splitlines()
    if sorted(untracked) != sorted(expected_untracked):
        raise SystemExit("Unexpected source-copy changes")
    print("REGION_CAPTURE_PROVENANCE " + json.dumps({"source": SOURCE, "tree": SOURCE_TREE, "aprl": APRL}))
    for name in FILES:
        raw = (output / name).read_bytes()
        content = raw.decode("utf-8")
        pieces = [content[i:i + 3000] for i in range(0, len(content), 3000)]
        for index, piece in enumerate(pieces):
            print("REGION_CAPTURE_JSON " + json.dumps({
                "name": name, "sha256": hashlib.sha256(raw).hexdigest(),
                "index": index, "count": len(pieces), "content": piece}, ensure_ascii=True))

    # Emit complete observations first; missing/changed goldens still fail.
    for name in FILES:
        verify_capture_bytes((output / name).read_bytes(),
                             ROOT / "internal/plugins/region/testdata" / name)
