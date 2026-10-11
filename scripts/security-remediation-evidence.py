#!/usr/bin/env python3
"""Temporary read-only hosted evidence, not acceptance or Git publication."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import urllib.request

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / '.build' / 'security-evidence'
OUT.mkdir(parents=True, exist_ok=True)


def run(args, *, allowed=(0,), capture=True, env=None):
    p = subprocess.run(args, cwd=ROOT, env=env, text=True,
                       stdout=subprocess.PIPE if capture else None,
                       stderr=subprocess.PIPE if capture else None, timeout=600)
    if p.returncode not in allowed:
        raise RuntimeError('unexpected command status: ' + repr(args) + '\n'
                           + str(p.returncode) + '\n'
                           + (p.stdout or '')[-4000:] + (p.stderr or '')[-4000:])
    return p


def emit(name, raw):
    text = raw.decode('utf-8')
    parts = [text[i:i+3000] for i in range(0, len(text), 3000)] or ['']
    sha = hashlib.sha256(raw).hexdigest()
    for i, part in enumerate(parts):
        print('SECURITY_EVIDENCE_JSON ' + json.dumps({
            'name': name, 'sha256': sha, 'index': i, 'count': len(parts),
            'content': part}, ensure_ascii=True), flush=True)


def record(name, args, allowed=(0,), env=None):
    p = run(args, allowed=allowed, env=env)
    raw = (p.stdout or '').encode('utf-8')
    (OUT/name).write_bytes(raw)
    emit(name, raw)
    print('SECURITY_COMMAND_STATUS ' + json.dumps({
        'name': name, 'command': args, 'exit': p.returncode,
        'stderr': p.stderr}), flush=True)
    return p.returncode


def scan(prefix):
    env = dict(os.environ, GOBIN=str(OUT/'tools'))
    run(['go', 'install', 'golang.org/x/vuln/cmd/govulncheck@v1.8.0'], env=env)
    scanner = str(OUT/'tools'/'govulncheck')
    status = record(prefix+'-vuln.json', [scanner, '-json', './...'], (0, 3))
    verbose = record(prefix+'-vuln.txt', [scanner, '-show=verbose', './...'], (0, 3))
    if status != verbose:
        raise RuntimeError('inconsistent diagnostic scanner status')
    return status


phase = sys.argv[1]
print('SECURITY_EXECUTION_IDENTITY ' + json.dumps({
    'phase': phase, 'head': run(['git', 'rev-parse', 'HEAD']).stdout.strip(),
    'tree': run(['git', 'rev-parse', 'HEAD^{tree}']).stdout.strip(),
    'parents': run(['git', 'show', '-s', '--format=%P']).stdout.strip().split(),
    'go': run(['go', 'version']).stdout.strip(),
    'aprl': run(['git', '-C', 'internal/rules/upstream/aprl', 'rev-parse', 'HEAD']).stdout.strip()
}), flush=True)
if phase == 'old':
    if run(['go', 'env', 'GOVERSION']).stdout.strip() != 'go1.26.8':
        raise RuntimeError('unexpected old toolchain')
    record('old-module-graph.json', ['go', 'list', '-m', '-json', 'all'])
    if scan('old') != 3:
        raise RuntimeError('known old reachable scan did not reproduce exit3')
    print('OLD_SECURITY_GATE_REPRODUCED_FAILURE_NOT_ACCEPTANCE', flush=True)
elif phase == 'fixed':
    if run(['go', 'env', 'GOVERSION']).stdout.strip() != 'go1.26.9':
        raise RuntimeError('unexpected fixed toolchain')
    expected = '42d158b4d8f7b61ac0a830567c940a86098fb7aac52e467a5ebec03ef5cc2f8d'
    release_data = json.load(urllib.request.urlopen('https://go.dev/dl/?mode=json&include=all', timeout=60))
    release = next(r for r in release_data if r['version'] == 'go1.26.9')
    entry = next(f for f in release['files'] if f['filename'] == 'go1.26.9.linux-amd64.tar.gz')
    if entry['sha256'] != expected:
        raise RuntimeError('publisher checksum mismatch')
    archive = OUT/'go1.26.9.linux-amd64.tar.gz'
    with urllib.request.urlopen('https://go.dev/dl/'+entry['filename'], timeout=120) as response:
        raw = response.read(90*1024*1024+1)
    if len(raw) != entry['size'] or hashlib.sha256(raw).hexdigest() != expected:
        raise RuntimeError('actual downloaded archive checksum/size mismatch')
    archive.write_bytes(raw)
    emit('go-publisher.json', (json.dumps(entry, indent=2)+'\n').encode())
    run(['go', 'mod', 'edit', '-go=1.26.9'])
    run(['go', 'get', 'golang.org/x/net@v0.60.0'])
    run(['go', 'mod', 'tidy'])
    run(['go', 'mod', 'verify'])
    # Only the isolated proposal checkout is changed. Nothing is pushed.
    for path in ['scripts/dependency-inventory.py', 'scripts/tests/dependency-inventory.py']:
        p = ROOT/path
        p.write_bytes(p.read_bytes().replace(b'1.26.8', b'1.26.9'))
    run(['python3', 'scripts/tests/dependency-inventory.py'], capture=False)
    run(['python3', 'scripts/dependency-inventory.py'], capture=False)
    run(['python3', 'scripts/dependency-inventory.py', '--check'], capture=False)
    record('fixed-module-graph.json', ['go', 'list', '-m', '-json', 'all'])
    for module in ['net', 'crypto', 'sys', 'text', 'term']:
        record('why-'+module+'.txt', ['go', 'mod', 'why', '-m', 'golang.org/x/'+module])
    run(['go', 'build', '-o', str(OUT/'cloud-assess'), './cmd/cloud-assess'])
    record('compiled-cli.txt', ['go', 'version', '-m', str(OUT/'cloud-assess')])
    for target in ['linux', 'windows']:
        env = dict(os.environ, GOOS=target, GOARCH='amd64', CGO_ENABLED='0')
        record(target+'-cli-deps.txt', ['go', 'list', '-buildvcs=false', '-deps', '-f', '{{if .Module}}{{.ImportPath}} {{.Module.Path}} {{.Module.Version}}{{end}}', './cmd/cloud-assess'], env=env)
    status = scan('fixed')
    for path in ['go.mod', 'go.sum', 'docs/dependencies/inventory.json', 'docs/dependencies/NOTICES.md']:
        emit(path, (ROOT/path).read_bytes())
    record('candidate-diff.txt', ['git', 'diff', '--', 'go.mod', 'go.sum', 'scripts/dependency-inventory.py', 'scripts/tests/dependency-inventory.py'])
    if status:
        raise RuntimeError('candidate reachable security scan still failed; evidence retained, no acceptance')
    print('ISOLATED_FIXED_SECURITY_SCAN_SUCCESS_NOT_NATIVE_ACCEPTANCE', flush=True)
else:
    raise RuntimeError('unknown phase')
