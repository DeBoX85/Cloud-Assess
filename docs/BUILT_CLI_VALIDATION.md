# Built CLI offline validation

Required Linux and Windows jobs run `scripts/tests/built-cli.py` against their actual `.build/cloud-assess` or `.build/cloud-assess.exe`. The program is copied into a temporary directory containing spaces and executed directly, without `go run`, an injected executor or a product-side fixture switch.

## What is checked

- Four root/help/scan-help/version cases verify product identity, output/redaction flags, success exits, empty stderr and absence of new files.
- Seventeen rejection cases verify exit **1**, the precise expected error, empty stdout, and unchanged hashes for pre-existing JSON/XLSX/CSV/SARIF sentinels and the rest of the temporary directory. The invalid-input matrix includes CLI parsing, malformed/negative assessment timeouts, stage configuration, mandatory Graph, unavailable plugins, stage parameters, impact and filter loading/validation. A valid scalar-tag filter is paired with an unavailable plugin to exercise successful YAML loading without Azure execution.
- A loopback proxy/identity-endpoint tripwire records HTTP attempts while Azure environment credentials are removed and Azure CLI config is isolated. These preflight cases must produce zero observed requests. This is focused boundary evidence, not a general network sandbox or proof that every adapter is read-only.
- `go version -m -json` inspects the actual executable: Go 1.26.8, main/package identity, host OS/amd64, CGO disabled and exact compiled module versions/checksums against [the inventory](dependencies/inventory.json). Current membership is 24 modules on Linux and 26 on Windows. A binary SHA-256 is printed for run diagnostics; it is not signed provenance.

## Local execution

Use Go 1.26.8, Python 3 and the pinned APRL checkout. Build with CGO disabled, then run:

```sh
CGO_ENABLED=0 go build -o .build/cloud-assess ./cmd/cloud-assess
python3 scripts/tests/built-cli.py --binary .build/cloud-assess
```

On Windows PowerShell:

```powershell
$env:CGO_ENABLED = '0'
go build -o .build/cloud-assess.exe ./cmd/cloud-assess
if ($LASTEXITCODE -ne 0) { throw 'CLI build failed' }
python scripts/tests/built-cli.py --binary .build/cloud-assess.exe
if ($LASTEXITCODE -ne 0) { throw 'Built CLI checks failed' }
```

Use `--expected-version` to provide an independently known version when inspecting a version-stamped candidate; CI currently expects `dev`. Use `--go /absolute/path/to/go` for metadata inspection when Go is outside PATH. An unsupported target, modified dependency set or unreviewed replacement must fail instead of being silently skipped. CI's ordinary build keeps VCS stamping; the local checkpoint used `-buildvcs=false` because its environment failed VCS status acquisition, documented in the failure register. That local workaround is not release provenance.

## What remains separate

These cases do not execute successful Azure scans or generate complete reports with the standalone CLI. Existing `cmd/cloud-assess/process_test.go` uses the Go test executable with synthetic assessment injection to verify exits 0/1/2/3 and persisted JSON. Existing application/coordinator tests independently verify real renderer formats and redaction with synthetic data. Keep these evidence types distinct; do not add a hidden production fixture mode to make an offline full scan appear real.

Authenticating/scanning through the built CLI, installed-package behavior, additional architectures/build flags, signed/checksummed distributions and security/operations approval remain release work. The Gate 004 matrix stays open. The build-information inspection uses the [official Go command's version semantics](https://pkg.go.dev/cmd/go#hdr-Print_Go_version).


AI public candidate adds a twelfth compiled-CLI test: source metadata/standalone inherited help and six partial cloud triples across standalone/mixed/aif-scanner paths reject before credential construction, preserving reports/directory hashes and empty stdout. China/Government/unknown cloud rejects; complete public authority/endpoint/audience reaches the deliberately invalid credential selection. The managed HTTP/auth tripwire records no traffic. Region-selection remains the actual unavailable-plugin oracle; all five implemented names are listed alphabetically. Native Windows acceptance and live Azure execution stay distinct.
