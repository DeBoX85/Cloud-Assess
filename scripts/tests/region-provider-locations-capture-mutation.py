#!/usr/bin/env python3
"""Require compiling named provider source-cell assertion failures and restoration."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ("sdk-version", "TestProviderSourceMetadata", "inputs"),
    ("tls-disabled", "TestProviderSourceMetadata", "inputs"),
    ("global-location", "TestProviderSourceCacheSemantics", "outputs"),
    ("overwrite-collision", "TestProviderSourceCacheSemantics", "outputs"),
    ("lost-prefix", "TestProviderSourceCacheSemantics", "outputs"),
    ("null-panic", "TestProviderSourceCacheSemantics", "outputs"),
    ("request-page", "TestProviderSourceRequests", "outputs"),
    ("method", "TestProviderSourceRequests", "outputs"),
    ("caller-fallback", "TestProviderSourceCallerRows", "outputs"),
    ("caller-cancelled", "TestProviderSourceCallerRows", "outputs"),
)


def run(directory, test):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + test + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix="cloud-assess-provider-source-control-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    for name, test, kind in MUTATIONS:
        path = directory / ("internal/plugins/region/testdata/source-provider-locations-" + kind + ".json")
        original = path.read_bytes()
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit("Provider source baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        modified = json.loads(original)
        if name == "sdk-version":
            modified["SDKModule"] = modified["SDKModule"].replace("v1.2.0", "v1.1.1")
        elif name == "tls-disabled":
            modified["TLSVerificationEnabled"] = False
        elif name == "global-location":
            modified["global"]["Cache"]["Data"]["microsoft.compute"]["virtualmachines"] = {"westeurope": {}}
        elif name == "overwrite-collision":
            modified["duplicate-type"]["Cache"]["Data"]["microsoft.compute"]["virtualmachines"]["westeurope"] = {}
        elif name == "lost-prefix":
            modified["later-http403"]["Cache"] = modified["basic"]["Cache"]
        elif name == "null-panic":
            modified["null-provider"]["Panicked"] = False
        elif name == "request-page":
            modified["paged"]["Requests"] = modified["paged"]["Requests"][:1]
        elif name == "method":
            modified["basic"]["Requests"][0] = modified["basic"]["Requests"][0].replace("GET ", "POST ")
        elif name == "caller-fallback":
            modified["later-http403"]["CallerRows"][0]["AvailabilityPercent"] = 100
        elif name == "caller-cancelled":
            modified["cancelled"]["CallerRows"] = []
        else:
            raise SystemExit("Unknown provider mutation")
        try:
            path.write_text(json.dumps(modified, ensure_ascii=True), encoding="utf-8")
            outcome = run(directory, test)
        finally:
            path.write_bytes(original)
        if outcome.returncode == 0 or "--- FAIL: " + test not in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr or "[build failed]" in outcome.stdout:
            raise SystemExit("Provider source mutation survived or failed outside named assertion: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit("Restored provider source baseline failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Provider source mutation rejected:", name)
