#!/usr/bin/env python3
"""Require compiling ARG framing regressions to detect selected production faults."""
from pathlib import Path
import os
import subprocess

ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / 'internal/arg/http_transport.go'
GO = os.environ.get('QA_GO', 'go')
SELECT = '^TestHTTPTransport(RejectsTrailingData|AcceptsOnlyWhitespaceAfterDocument|PreservesTrailingReadError)$'


def run(select):
    return subprocess.run([GO, 'test', './internal/arg', '-run', select, '-count=1', '-v'],
                          cwd=ROOT, capture_output=True, text=True, timeout=120)


def baseline():
    result = run(SELECT)
    if result.returncode:
        raise SystemExit('ARG framing baseline failed\n' + result.stdout + result.stderr)


original = SOURCE.read_bytes()
text = original.decode('utf-8')
controls = (
    ('buffered-suffix', 'io.MultiReader(decoder.Buffered(), response.Body)', 'response.Body', 'TestHTTPTransportRejectsTrailingData'),
    ('invalid-suffix', 'return fmt.Errorf("decode ARG response: unexpected trailing data")', 'continue', 'TestHTTPTransportRejectsTrailingData'),
    ('tail-read-error', 'return fmt.Errorf("read ARG response body: %w", err)', 'return nil', 'TestHTTPTransportPreservesTrailingReadError'),
    ('premature-eof', 'if err == io.EOF {', 'if n >= 0 {', 'TestHTTPTransportRejectsTrailingData'),
)
baseline()
try:
    for name, old, new, expected in controls:
        if text.count(old) != 1:
            raise SystemExit('ARG mutation anchor changed: ' + name)
        SOURCE.write_bytes(text.replace(old, new, 1).encode('utf-8'))
        compiled = run('^$')
        if compiled.returncode:
            raise SystemExit('ARG mutation did not compile: ' + name + '\n' + compiled.stdout + compiled.stderr)
        failed = run('^' + expected + '$')
        if failed.returncode == 0 or ('--- FAIL: ' + expected) not in failed.stdout or '[build failed]' in failed.stdout + failed.stderr:
            raise SystemExit('ARG mutation missed named assertion: ' + name + '\n' + failed.stdout + failed.stderr)
        SOURCE.write_bytes(original)
        baseline()
        print('ARG response framing compiling mutation rejected: ' + name + '; restored baseline passed', flush=True)
finally:
    SOURCE.write_bytes(original)
if SOURCE.read_bytes() != original:
    raise SystemExit('ARG framing original source was not restored')
