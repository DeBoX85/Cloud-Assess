#!/usr/bin/env bash
# Create a fresh Linux amd64 development workspace without changing global tools.
set -euo pipefail

if [[ $# != 1 || $1 == --help ]]; then
  printf 'Usage: bash scripts/bootstrap-workspace.sh NEW_DIRECTORY\n'
  printf 'Requires Linux amd64, bash, git with HTTPS, curl, tar, sha256sum, python3, rg and gcc.\n'
  [[ $# == 1 && $1 == --help ]] && exit 0
  exit 2
fi
[[ $(uname -s) == Linux && $(uname -m) == x86_64 ]] || {
  printf 'This bootstrap supports Linux amd64 only.\n' >&2
  exit 1
}
for tool in git curl tar sha256sum python3 rg gcc; do
  command -v "$tool" >/dev/null || { printf 'Missing prerequisite: %s\n' "$tool" >&2; exit 1; }
done
[[ -x $(git --exec-path)/git-remote-https ]] || {
  printf 'Git HTTPS helper is missing.\n' >&2
  exit 1
}

# mkdir refuses existing directories, files and symlinks. Never reset/reuse a tree.
mkdir -- "$1"
workspace=$(cd -- "$1" && pwd -P)
trap 'printf "Bootstrap failed. Partial task-owned workspace retained: %s\n" "$workspace" >&2' ERR
mkdir -p "$workspace/tools/downloads" "$workspace/tools/bin" "$workspace/cache" "$workspace/logs"

go_version=1.26.8
go_sha=d0f743b33e8d8945e6b1f432edd15785c70507121d6e2a723b21285eddf8b57b
pwsh_version=7.6.6
pwsh_sha=ddbc4a2d113bbd46d283cfedcbcd117a70caefd7673f41f2b4e0000badf103bc
source_pin=8e4f0577f3615e6c9014c031bcad079f235369cc
aprl_pin=60eaddda76541f6adbc1c5ffa686829807e55e29

download() {
  local url=$1 file=$2 expected=$3
  curl --fail --location --silent --show-error --retry 2 --max-time 300 "$url" -o "$file"
  printf '%s  %s\n' "$expected" "$file" | sha256sum -c -
}
download "https://go.dev/dl/go${go_version}.linux-amd64.tar.gz" \
  "$workspace/tools/downloads/go.tar.gz" "$go_sha"
tar --no-same-owner -xzf "$workspace/tools/downloads/go.tar.gz" -C "$workspace/tools"
download "https://github.com/PowerShell/PowerShell/releases/download/v${pwsh_version}/powershell-${pwsh_version}-linux-x64.tar.gz" \
  "$workspace/tools/downloads/powershell.tar.gz" "$pwsh_sha"
mkdir "$workspace/tools/powershell"
tar --no-same-owner -xzf "$workspace/tools/downloads/powershell.tar.gz" -C "$workspace/tools/powershell"
chmod u+x "$workspace/tools/powershell/pwsh"

cat > "$workspace/env.sh" <<'ENV'
#!/usr/bin/env bash
# Source per shell. Tool/cache settings are local; proxy and credential settings stay intact.
cloud_assess_workspace="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
export PATH="$cloud_assess_workspace/tools/go/bin:$cloud_assess_workspace/tools/powershell:$cloud_assess_workspace/tools/bin:$PATH"
export GOPATH="$cloud_assess_workspace/cache/gopath"
export GOMODCACHE="$cloud_assess_workspace/cache/modules"
export GOCACHE="$cloud_assess_workspace/cache/build"
export GOTOOLCHAIN=local
export GOMAXPROCS=4
export GOFLAGS='-mod=readonly -p=4'
ENV
# shellcheck source=/dev/null
source "$workspace/env.sh"
[[ $(go version) == "go version go${go_version} linux/amd64" ]]
[[ $(pwsh -NoLogo -NoProfile -NonInteractive -Command '$PSVersionTable.PSVersion.ToString()') == "$pwsh_version" ]]

git clone --single-branch --branch bootstrap/core-v1 \
  https://github.com/DeBoX85/Cloud-Assess.git "$workspace/target"
git -C "$workspace/target" config user.name 'Denis Bogunic'
git -C "$workspace/target" config user.email '134433647+DeBoX85@users.noreply.github.com'
[[ $(awk '$1 == "go" { print $2 }' "$workspace/target/go.mod") == "$go_version" ]] || {
  printf 'Repository Go version changed. Review/update this recipe before continuing.\n' >&2
  exit 1
}
git -C "$workspace/target" submodule update --init --recursive
git clone --no-checkout https://github.com/DeBoX85/azqr.git "$workspace/azqr-reference"
git -C "$workspace/azqr-reference" checkout --detach "$source_pin"
git -C "$workspace/azqr-reference" submodule update --init --recursive
[[ $(git -C "$workspace/target/internal/rules/upstream/aprl" rev-parse HEAD) == "$aprl_pin" ]]
[[ $(git -C "$workspace/azqr-reference/internal/graph/aprl" rev-parse HEAD) == "$aprl_pin" ]]

for repository in target azqr-reference; do
  (
    cd -- "$workspace/$repository"
    go mod download
    go mod verify
    [[ -z $(git status --porcelain) ]]
  ) > "$workspace/logs/${repository}-modules.log" 2>&1
done
(
  cd -- "$workspace/target"
  GOBIN="$workspace/tools/bin" go install golang.org/x/vuln/cmd/govulncheck@v1.8.0
) > "$workspace/logs/vulnerability-tool-install.log" 2>&1

python3 - "$workspace" "$go_version" "$go_sha" "$pwsh_version" "$pwsh_sha" <<'PY'
import json, pathlib, subprocess, sys
root = pathlib.Path(sys.argv[1])
def git(path, *args):
    return subprocess.check_output(['git', '-C', str(path), *args], text=True).strip()
pins = {
    'internal/rules/upstream/aprl': '60eaddda76541f6adbc1c5ffa686829807e55e29',
    'internal/rules/upstream/orphan-resources': 'a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5',
    'internal/rules/custom': '674b9b3dcb443ce6dc445b48e1db47e4a0ca7082',
}
for path, expected in pins.items():
    if git(root/'target', 'rev-parse', 'HEAD:'+path) != expected:
        raise SystemExit('Target provenance changed: '+path)
if git(root/'target', 'hash-object', 'internal/skus/known_skus.yaml') != 'a2d97a60ec445ce4023a1a5a9b6c4dca76e31f5a':
    raise SystemExit('SKU provenance changed')
data = {'schemaVersion': 1, 'platform': 'linux/amd64',
        'go': {'version': sys.argv[2], 'archiveSHA256': sys.argv[3]},
        'powershell': {'version': sys.argv[4], 'archiveSHA256': sys.argv[5]},
        'govulncheck': 'v1.8.0', 'repositories': {}, 'targetPins': pins,
        'qaExecuted': False}
for name in ['target', 'azqr-reference']:
    path = root/name
    if git(path, 'status', '--porcelain'):
        raise SystemExit('Unexpected changed files: '+name)
    data['repositories'][name] = {'head': git(path, 'rev-parse', 'HEAD'),
                                  'tree': git(path, 'rev-parse', 'HEAD^{tree}'),
                                  'parents': git(path, 'show', '-s', '--format=%P').split()}
(root/'workspace-manifest.json').write_text(json.dumps(data, indent=2)+'\n')
PY
trap - ERR
printf 'Workspace created and tools/modules/pins verified: %s\n' "$workspace"
printf 'Activate with: source %q\n' "$workspace/env.sh"
printf 'Read the live handover/open PRs and run offline QA before editing. No Azure scan was run.\n'
