#!/usr/bin/env python3
"""Check authored local links and operator CLI/PowerShell examples without scans."""
import argparse
import os
from pathlib import Path
import re
import subprocess
import tempfile
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[2]


def check_links():
    files = [ROOT / 'README.md'] + sorted((ROOT / 'docs').glob('*.md'))
    count = 0
    for path in files:
        for target in re.findall(r'\]\(([^)]+)\)', path.read_text(encoding='utf-8')):
            # Scope is authored inline Markdown destinations, not raw third-party notices.
            parsed = urlsplit(target)
            if parsed.scheme or parsed.netloc or not parsed.path:
                continue
            destination = path.parent / unquote(parsed.path)
            if not destination.exists():
                raise ValueError(f'missing local link: {path.relative_to(ROOT)} -> {target}')
            count += 1
    print(f'Authored Markdown local-file links: {count} checked')


def command(args, env=None):
    result = subprocess.run(args, cwd=ROOT, env=env, capture_output=True, text=True, timeout=30)
    if result.returncode:
        raise ValueError(f'offline documentation command failed: {args[0]}\n{result.stderr[-3000:]}')
    return result.stdout


def check_operator_examples(binary, powershell):
    runbook = (ROOT / 'docs/OPERATIONS.md').read_text(encoding='utf-8')
    help_text = command([str(binary), '--help']) + command([str(binary), 'scan', '--help'])
    cloud_text = re.sub(r'`az\s+[^`]+`', '', runbook)
    flags = sorted(set(re.findall(r'(?<![\w-])--[a-z][a-z0-9-]*', cloud_text)))
    for flag in flags:
        if not re.search(re.escape(flag) + r'(?=[\s=,]|$)', help_text):
            raise ValueError('runbook flag absent from actual CLI help: ' + flag)
    snippets = re.findall(r'```powershell\s*\n(.*?)\n```', runbook, re.DOTALL)
    if not snippets:
        raise ValueError('operator runbook has no PowerShell example to validate')
    with tempfile.TemporaryDirectory(prefix='cloud-assess-doc-qa-') as directory:
        directory = Path(directory)
        parser = directory / 'parse-only.ps1'
        parser.write_text('''param([string]$SnippetPath)
$tokens = $null
$errors = $null
[System.Management.Automation.Language.Parser]::ParseFile($SnippetPath, [ref]$tokens, [ref]$errors) | Out-Null
if ($errors.Count -gt 0) { $errors | ForEach-Object { Write-Error $_.Message }; exit 1 }
''', encoding='utf-8')
        env = dict(os.environ)
        # A PowerShell 7 parent can leave incompatible module paths for Windows 5.1.
        env.pop('PSModulePath', None)
        for index, snippet in enumerate(snippets):
            path = directory / f'example-{index}.ps1'
            path.write_text(snippet, encoding='utf-8')
            command([powershell, '-NoProfile', '-NonInteractive', '-File', str(parser), '-SnippetPath', str(path)], env)
    print(f'Operator examples: {len(flags)} flags and {len(snippets)} PowerShell snippet(s) checked; no scan executed')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--binary', type=Path)
    parser.add_argument('--powershell')
    args = parser.parse_args()
    if bool(args.binary) != bool(args.powershell):
        parser.error('--binary and --powershell must be supplied together')
    check_links()
    if args.binary:
        check_operator_examples(args.binary.resolve(), args.powershell)
