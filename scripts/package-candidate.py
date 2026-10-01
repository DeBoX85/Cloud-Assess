#!/usr/bin/env python3
"""Build development candidate ZIPs; never publish or approve a release."""
import argparse
import datetime
import hashlib
import importlib.util
import json
import os
from pathlib import Path, PurePosixPath
import re
import stat
import subprocess
import sys
import tempfile
import zipfile

ROOT = Path(__file__).resolve().parents[1]
COMMON = {'LICENSE', 'NOTICE.md', 'THIRD_PARTY_LICENSES.md', 'DEPENDENCY_NOTICES.md',
          'dependency-inventory.json', 'BUILD_INFO.json', 'INSTALL.md', 'BRANDING_PROFILE.json'}

sys.dont_write_bytecode = True
profile_spec = importlib.util.spec_from_file_location('branding_profile', ROOT / 'scripts/branding-profile.py')
profile_module = importlib.util.module_from_spec(profile_spec)
profile_spec.loader.exec_module(profile_module)


def installation(profile, version, target):
    name = profile['cliName']
    executable = name + ('.exe' if target.startswith('windows/') else '')
    root = name + '-' + version + '-' + target.replace('/', '-')
    return (source('docs/PACKAGE_INSTALL.md').decode('utf-8').replace('cloud-assess', name).encode('utf-8') +
            ('\n## This candidate\n\nExecutable: `' + executable + '`. Archive: `' + root + '.zip`.\n'
             'Inspect the embedded profile with `' + executable + ' branding`.\n'
             'From the extracted directory: `./' + executable + ' --version` (Linux) or '
             '`.\\' + executable + ' --version` (PowerShell).\n'
             'BRANDING_PROFILE.json records the canonical public profile; the manifest records its SHA-256.\n').encode())


def sha(data):
    return hashlib.sha256(data).hexdigest()


def command(args):
    r = subprocess.run(args, cwd=ROOT, capture_output=True, timeout=120)
    if r.returncode:
        raise ValueError(f'command failed: {args[0]} {args[1]}\n{r.stderr.decode("utf-8", errors="replace")[-3000:]}')
    return r.stdout


def source(path, revision='HEAD'):
    # Read committed bytes, avoiding checkout line-ending transformations.
    return command(['git', 'show', revision + ':' + path])


def validate_info(info, inventory, revision, target):
    settings = {s['Key']: s['Value'] for s in info['Settings']}
    if info['GoVersion'] != inventory['goVersion'] or info['Path'] != 'github.com/DeBoX85/Cloud-Assess/cmd/cloud-assess':
        raise ValueError('unexpected executable identity/toolchain')
    if info['Main']['Path'] != 'github.com/DeBoX85/Cloud-Assess':
        raise ValueError('unexpected main module')
    if settings.get('vcs.revision') != revision or settings.get('vcs.modified') != 'false':
        raise ValueError('binary must have matching clean VCS revision')
    if settings.get('-trimpath') != 'true':
        raise ValueError('candidate binary must use trimpath')
    if settings.get('CGO_ENABLED') != '0' or settings.get('GOOS') + '/' + settings.get('GOARCH') != target:
        raise ValueError('unexpected compiled target/CGO configuration')
    expected = {m['path']: (m['version'], m['sum']) for m in inventory['modules'] if target in m['cliTargets']}
    actual = {}
    for dep in info['Deps']:
        if 'Replace' in dep or dep['Path'] in actual:
            raise ValueError('unreviewed replacement/duplicate compiled dependency')
        actual[dep['Path']] = (dep['Version'], dep['Sum'])
    if actual != expected or not expected:
        raise ValueError('compiled dependencies do not match inventory')


def package(binary, version, output, go):
    if not re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9._-]{0,63}', version):
        raise ValueError('invalid candidate version label')
    if command(['git', 'status', '--porcelain', '--untracked-files=normal']).strip():
        raise ValueError('candidate packaging requires a clean source checkout')
    revision = command(['git', 'rev-parse', 'HEAD']).decode().strip()
    tree = command(['git', 'rev-parse', revision + '^{tree}']).decode().strip()
    target = ('windows' if os.name == 'nt' else 'linux') + '/amd64'
    inventory_bytes = source('docs/dependencies/inventory.json', revision)
    inventory = json.loads(inventory_bytes)
    if target not in inventory['targets']:
        raise ValueError('target has no reviewed inventory')
    for path, expected in inventory['inputs'].items():
        if sha(source(path, revision)) != expected:
            raise ValueError('stale committed inventory input: ' + path)
    info = json.loads(command([go, 'version', '-m', '-json', str(binary)]))
    validate_info(info, inventory, revision, target)
    profile_bytes = command([str(binary), 'branding']).rstrip(b'\r\n')
    profile = profile_module.parse(profile_bytes)
    if command([str(binary), '--version']).decode().strip() != profile['cliName'] + ' version ' + version:
        raise ValueError('candidate version does not match executable')
    executable = profile['cliName'] + ('.exe' if os.name == 'nt' else '')
    payload = {executable: binary.read_bytes(), 'LICENSE': source('LICENSE', revision), 'NOTICE.md': source('NOTICE.md', revision),
               'THIRD_PARTY_LICENSES.md': source('THIRD_PARTY_LICENSES.md', revision),
               'DEPENDENCY_NOTICES.md': source('docs/dependencies/NOTICES.md', revision),
               'dependency-inventory.json': inventory_bytes,
               'BUILD_INFO.json': (json.dumps(info, indent=2) + '\n').encode(),
               'INSTALL.md': installation(profile, version, target),
               'BRANDING_PROFILE.json': profile_bytes}
    manifest = {'schemaVersion': 2, 'status': 'development-candidate-not-release-approved',
                'version': version, 'target': target, 'sourceCommit': revision, 'sourceTree': tree,
                'referencePins': inventory['referencePins'], 'brandingSHA256': sha(profile_bytes),
                'files': {p: sha(data) for p, data in sorted(payload.items())}}
    payload['PACKAGE_MANIFEST.json'] = (json.dumps(manifest, indent=2) + '\n').encode()
    root = profile['cliName'] + '-' + version + '-' + target.replace('/', '-')
    archive = output / (root + '.zip')
    checksum = output / (root + '.zip.sha256')
    if archive.exists() or checksum.exists():
        raise ValueError('refusing to overwrite an existing candidate')
    timestamp = int(command(['git', 'show', '-s', '--format=%ct', revision]).strip())
    dt = datetime.datetime.fromtimestamp(timestamp, datetime.timezone.utc)
    date_time = (max(1980, dt.year), dt.month, dt.day, dt.hour, dt.minute, dt.second)
    output.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(prefix='candidate-stage-', dir=output) as staging:
        staged = Path(staging) / archive.name
        with zipfile.ZipFile(staged, 'w', compression=zipfile.ZIP_DEFLATED, compresslevel=9) as z:
            for name, data in sorted(payload.items()):
                item = zipfile.ZipInfo(root + '/' + name, date_time=date_time)
                item.create_system = 3
                item.compress_type = zipfile.ZIP_DEFLATED
                item.external_attr = (stat.S_IFREG | (0o755 if name == executable else 0o644)) << 16
                z.writestr(item, data, compresslevel=9)
        checksum_text = sha(staged.read_bytes()) + '  ' + archive.name + '\n'
        staged_checksum = Path(staging) / checksum.name
        staged_checksum.write_bytes(checksum_text.encode())
        # Two files are published, not a crash-atomic multi-file transaction.
        # Exclusive hard links also reject a concurrent publisher without overwrite.
        os.link(staged, archive)
        os.link(staged_checksum, checksum)
    return archive, checksum


def verify_extract(archive, checksum, destination):
    if destination.exists():
        raise ValueError('extraction requires a new destination')
    expected_line = checksum.read_text().strip()
    if archive.stat().st_size > 256 * 1024 * 1024:
        raise ValueError('archive exceeds size limit')
    if expected_line != sha(archive.read_bytes()) + '  ' + archive.name:
        raise ValueError('archive checksum mismatch')
    with zipfile.ZipFile(archive) as z:
        entries = z.infolist()
        if len(entries) != len(COMMON) + 2 or sum(e.file_size for e in entries) > 256 * 1024 * 1024:
            raise ValueError('unexpected archive size/member count')
        seen = set()
        roots = set()
        payload = {}
        modes = {}
        for e in entries:
            path = PurePosixPath(e.filename)
            parts = path.parts
            if e.filename != '/'.join(parts) or path.is_absolute() or len(parts) != 2 or any(p in ('.', '..') for p in parts) or '\\' in e.filename or ':' in e.filename:
                raise ValueError('unsafe archive path')
            if e.filename.casefold() in seen or e.is_dir() or stat.S_IFMT(e.external_attr >> 16) != stat.S_IFREG:
                raise ValueError('duplicate/nonregular archive member')
            seen.add(e.filename.casefold())
            roots.add(parts[0])
            payload[parts[1]] = z.read(e)
            modes[parts[1]] = (e.external_attr >> 16) & 0o777
        if len(roots) != 1 or 'PACKAGE_MANIFEST.json' not in payload:
            raise ValueError('invalid candidate layout')
        manifest = profile_module.strict_json(payload['PACKAGE_MANIFEST.json'])
        if not isinstance(manifest, dict):
            raise ValueError('unsupported candidate manifest')
        target = manifest.get('target')
        if target not in ('linux/amd64', 'windows/amd64') or type(manifest.get('schemaVersion')) is not int or manifest.get('schemaVersion') != 2 or manifest.get('status') != 'development-candidate-not-release-approved':
            raise ValueError('unsupported candidate manifest')
        version = manifest.get('version', '')
        if not isinstance(version, str) or not re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9._-]{0,63}', version):
            raise ValueError('invalid manifest version')
        if 'BRANDING_PROFILE.json' not in payload:
            raise ValueError('missing branding profile')
        profile = profile_module.parse(payload['BRANDING_PROFILE.json'])
        if manifest.get('brandingSHA256') != sha(payload['BRANDING_PROFILE.json']):
            raise ValueError('branding profile checksum mismatch')
        expected_root = profile['cliName'] + '-' + version + '-' + target.replace('/', '-')
        if roots != {expected_root}:
            raise ValueError('candidate root mismatch')
        executable = profile['cliName'] + ('.exe' if target.startswith('windows/') else '')
        if not isinstance(manifest.get('files'), dict):
            raise ValueError('invalid payload checksum map')
        if set(payload) != COMMON | {executable, 'PACKAGE_MANIFEST.json'} or set(manifest['files']) != COMMON | {executable}:
            raise ValueError('missing/unexpected candidate payload')
        for name in payload:
            if modes[name] != (0o755 if name == executable else 0o644):
                raise ValueError('unexpected payload permissions')
        for name, expected in manifest['files'].items():
            if sha(payload[name]) != expected:
                raise ValueError('payload checksum mismatch: ' + name)
        # Validate everything before creating output paths or applying executable mode.
        destination.mkdir(parents=True)
        for name, data in payload.items():
            path = destination / name
            path.write_bytes(data)
            if os.name != 'nt':
                path.chmod(modes[name])
    return destination / executable


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--binary', required=True, type=Path)
    parser.add_argument('--version', default='dev')
    parser.add_argument('--output', default='.build/packages', type=Path)
    parser.add_argument('--go', default=os.environ.get('QA_GO', 'go'))
    args = parser.parse_args()
    archive, checksum = package(args.binary.resolve(strict=True), args.version, args.output.resolve(), args.go)
    print(archive)
    print(checksum)


if __name__ == '__main__':
    main()
