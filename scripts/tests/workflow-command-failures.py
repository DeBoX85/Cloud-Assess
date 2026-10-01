#!/usr/bin/env python3
"""Execute real Windows workflow snippets with synthetic native command failures."""
import argparse
from pathlib import Path
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--powershell', default='pwsh')
    args = parser.parse_args()
    workflow = (ROOT / '.github/workflows/test.yml').read_text(encoding='utf-8').split('  windows-validation:', 1)[1]
    names = ['Validate branding builder and actual native executables',
             'Test Windows candidate package integrity and isolated installation']
    with tempfile.TemporaryDirectory(prefix='cloud-assess-workflow-exit-') as directory:
        directory = Path(directory)
        for index, name in enumerate(names):
            step = workflow.split('      - name: ' + name + '\n', 1)[1].split('\n      - name:', 1)[0]
            snippet = step.split('        run: |\n', 1)[1]
            snippet = '\n'.join(line[10:] for line in snippet.splitlines() if line.startswith('          '))
            path = directory / ('snippet-' + str(index) + '.ps1')
            path.write_text('''param([int]$FailOrdinal)
$ErrorActionPreference = 'Stop'
$script:callCount = 0
function python {
    $script:callCount++
    $global:LASTEXITCODE = if ($script:callCount -eq $FailOrdinal) { 1 } else { 0 }
}
$failed = $false
try {
''' + snippet + '''
} catch { $failed = $true }
if ($failed -ne ($FailOrdinal -gt 0)) { throw 'Workflow masked a native command failure' }
$expected = if ($FailOrdinal -eq 1) { 1 } else { 2 }
if ($script:callCount -ne $expected) { throw 'Workflow continued after a failed native command' }
''', encoding='utf-8')
            for ordinal in (0, 1, 2):
                result = subprocess.run([args.powershell, '-NoProfile', '-NonInteractive', '-File', str(path),
                                         '-FailOrdinal', str(ordinal)], capture_output=True, text=True, timeout=30)
                if result.returncode:
                    raise AssertionError(name + '\n' + result.stdout + result.stderr)
    print('Workflow command failures: both actual Windows blocks reject first/second failures and accept success (6 cases)')


if __name__ == '__main__':
    main()
