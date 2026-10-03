#!/usr/bin/env python3
"""Exercise bootstrap preservation and checksum-before-extraction boundaries."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
SCRIPT = ROOT / 'scripts/bootstrap-workspace.sh'


def command(script, args, env=None):
    return subprocess.run(['bash', str(script), *args], env=env,
                          capture_output=True, text=True, timeout=30)


def assert_checksum_rejected(script):
    with tempfile.TemporaryDirectory(prefix='workspace-checksum-') as directory:
        root = Path(directory)
        shims = root / 'shims'
        shims.mkdir()
        curl = shims / 'curl'
        curl.write_text('''#!/usr/bin/env python3
import pathlib, sys
pathlib.Path(sys.argv[sys.argv.index('-o')+1]).write_bytes(b'corrupt archive fixture')
''')
        tar = shims / 'tar'
        tar.write_text('''#!/bin/sh
printf reached > "$BOOTSTRAP_TEST_TAR_MARKER"
exit 99
''')
        curl.chmod(0o755)
        tar.chmod(0o755)
        marker = root / 'extraction-reached'
        workspace = root / 'workspace with spaces'
        env = dict(os.environ, PATH=str(shims)+os.pathsep+os.environ['PATH'],
                   BOOTSTRAP_TEST_TAR_MARKER=str(marker))
        result = command(script, [str(workspace)], env)
        if marker.exists():
            raise AssertionError('checksum fixture reached extraction')
        if result.returncode == 0 or 'FAILED' not in result.stdout:
            raise AssertionError('checksum fixture did not fail at checksum verification')
        if (workspace/'tools/go').exists() or (workspace/'target').exists():
            raise AssertionError('checksum fixture installed tools or cloned repositories')


class BootstrapTests(unittest.TestCase):
    def test_usage(self):
        self.assertEqual(command(SCRIPT, ['--help']).returncode, 0)
        self.assertEqual(command(SCRIPT, []).returncode, 2)

    def test_existing_workspace_preserved(self):
        with tempfile.TemporaryDirectory(prefix='workspace-preserve-') as directory:
            root = Path(directory)
            retained = root/'user-file'
            retained.write_bytes(b'retained user bytes')
            result = command(SCRIPT, [str(root)])
            self.assertNotEqual(result.returncode, 0)
            self.assertEqual(retained.read_bytes(), b'retained user bytes')
            self.assertEqual(sorted(p.name for p in root.iterdir()), ['user-file'])

    def test_checksum_before_extraction(self):
        assert_checksum_rejected(SCRIPT)


def verify_control():
    source = SCRIPT.read_text()
    anchor = '  printf \'%s  %s\\n\' "$expected" "$file" | sha256sum -c -'
    if source.count(anchor) != 1:
        raise AssertionError('checksum control anchor changed')
    with tempfile.TemporaryDirectory(prefix='workspace-control-') as directory:
        mutant = Path(directory)/'bootstrap.sh'
        mutant.write_text(source.replace(anchor, '  : # disabled checksum guard'))
        syntax = subprocess.run(['bash', '-n', str(mutant)], capture_output=True, text=True)
        if syntax.returncode:
            raise AssertionError('checksum control is not valid shell syntax')
        try:
            assert_checksum_rejected(mutant)
        except AssertionError as error:
            if str(error) != 'checksum fixture reached extraction':
                raise
        else:
            raise AssertionError('checksum behavioral control was not detected')
    assert_checksum_rejected(SCRIPT)
    print('Bootstrap checksum control: valid shell mutation rejected before unsafe extraction; restored guard passed')


if __name__ == '__main__':
    suite = unittest.defaultTestLoader.loadTestsFromTestCase(BootstrapTests)
    result = unittest.TextTestRunner().run(suite)
    if not result.wasSuccessful():
        raise SystemExit(1)
    verify_control()
