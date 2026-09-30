# Development candidate packages

The project can prepare development ZIP candidates for the native Linux and Windows amd64 configurations exercised by CI. This is packaging infrastructure, not a supported-platform declaration or a release decision. Packages remain labeled `development-candidate-not-release-approved` until the separate release process is completed.

## Build and package

Use a clean committed checkout, pinned APRL initialization, Go 1.26.8 and Python 3. On Linux:

```sh
mkdir -p .build
CGO_ENABLED=0 go build -trimpath -buildvcs=true -o .build/cloud-assess ./cmd/cloud-assess
python3 scripts/package-candidate.py --binary .build/cloud-assess --version dev
python3 scripts/tests/package-candidate.py --binary .build/cloud-assess
```

Windows PowerShell:

```powershell
New-Item -ItemType Directory -Path .build -Force | Out-Null
$env:CGO_ENABLED = '0'
go build -trimpath -buildvcs=true -o .build/cloud-assess.exe ./cmd/cloud-assess
if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
python scripts/package-candidate.py --binary .build/cloud-assess.exe --version dev
if ($LASTEXITCODE -ne 0) { throw 'Package failed' }
python scripts/tests/package-candidate.py --binary .build/cloud-assess.exe
if ($LASTEXITCODE -ne 0) { throw 'Package validation failed' }
```

`--go` selects an explicit toolchain binary and `--output` changes the staging/output directory. Existing ZIP/sidecar paths are never overwritten; exclusive same-filesystem hard links also reject concurrent publication. The output filesystem must support these links. The current integration case deliberately uses the independently expected version `dev`. Future version-stamped candidates must supply a matching known CLI version and extend integration expectations before being treated as validated releases.

The packager rejects dirty source trees, unstamped/mismatched/dirty binaries, non-CGO-disabled targets, missing trimpath, unexpected toolchains/modules/checksums, stale committed inventory inputs and inconsistent version labels. ZIP bytes use fixed source-commit timestamps, sorted members and fixed file modes. Repeated packaging of the **same binary and committed inputs within one runner** must be byte-equal; this does not prove separately compiled or cross-host artifacts are reproducible.

## Contents and validation

Each ZIP carries an executable, LICENSE, NOTICE.md, THIRD_PARTY_LICENSES.md, DEPENDENCY_NOTICES.md, dependency-inventory.json, BUILD_INFO.json, PACKAGE_MANIFEST.json and [INSTALL.md](PACKAGE_INSTALL.md). Source files and timestamps are read from the captured commit SHA rather than repeatedly resolving a movable HEAD. Source files are read as committed Git blobs to preserve exact notice bytes on both hosts. The package manifest contains the actual build's source commit/tree, pinned library references and payload hashes; the sidecar hashes the ZIP. This is unsigned diagnostic provenance and integrity evidence, not publisher authentication or a formal SBOM/legal approval.

Thirteen fixtures cover archive/payload corruption, missing notices, unsafe paths, case-insensitive duplicates, symlinks, excessive permissions, existing destination retention, unsupported targets and build-provenance mismatches. One required native integration case builds ZIPs twice, verifies equality, validates all members before creating an extraction directory, checks preserved notices/executable bytes and reruns the existing offline built-CLI checks through the extracted executable. It also verifies existing-output and simulated concurrent-publication rejection, retaining the competing archive and producing no false checksum sidecar. An injected checksum-publication failure after the ZIP write leaves a detectable orphan: extraction and overwrite retry are rejected and the existing ZIP is preserved. A local run without `--binary` explicitly skips this integration and is not equivalent to CI acceptance.

Archives/checksums are created locally under `.build/packages` by the manual packager. Required CI tests use temporary directories and do not upload or publish release artifacts. ZIP and sidecar publication is not a multi-file atomic transaction; never consume a package without its complete verified checksum/manifest evidence.

## Remaining release work

The runner already contains development tools and this is isolated-directory execution, not installation on a fresh operating system. Go/Python are not needed by the executable for help/version, but packaging/QA scripts require them. Additional architectures, real OS installers, PATH/service integration, independently verified source rebuilds, signing, trustworthy publisher provenance, complete package-file licensing review, Azure authentication/scans from installed binaries and operations/security approval remain separate requirements. Preserve old directories for rollback; no release is automatically published or marked approved here.

Build flag behavior follows the [official Go command documentation](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies). Archive member validation precedes writes rather than relying on implicit extraction behavior; see [Python ZIP documentation](https://docs.python.org/3/library/zipfile.html).
