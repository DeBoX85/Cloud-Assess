#!/usr/bin/env python3
"""Check a real built CLI offline; no injected executor or successful Azure scan."""
import argparse
import hashlib
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
import os
from pathlib import Path
import subprocess
import shutil
import tempfile
import threading
import unittest

ROOT = Path(__file__).resolve().parents[2]
BINARY = None
GO = None
EXPECTED_VERSION = "dev"


class BuiltCLITests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.traffic = []

        class Tripwire(BaseHTTPRequestHandler):
            def do_GET(self):
                cls.traffic.append((self.command, self.path))
                self.send_error(503, 'Offline QA tripwire')

            do_POST = do_GET
            do_CONNECT = do_GET

            def log_message(self, *_):
                pass

        cls.server = ThreadingHTTPServer(('127.0.0.1', 0), Tripwire)
        cls.thread = threading.Thread(target=cls.server.serve_forever, daemon=True)
        cls.thread.start()
        cls.endpoint = 'http://127.0.0.1:' + str(cls.server.server_address[1])

    @classmethod
    def tearDownClass(cls):
        cls.server.shutdown()
        cls.server.server_close()
        cls.thread.join(timeout=5)

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix='cloud-assess-built-cli-')
        self.addCleanup(self.temp.cleanup)
        self.directory = Path(self.temp.name) / 'reports with spaces'
        self.directory.mkdir()
        self.executable = self.directory / BINARY.name
        shutil.copy2(BINARY, self.executable)
        self.env = {k: v for k, v in os.environ.items() if not k.upper().startswith('AZURE_')}
        self.env.update({'AZURE_CONFIG_DIR': str(self.directory / 'isolated-azure-config'),
                         'HTTP_PROXY': self.endpoint, 'HTTPS_PROXY': self.endpoint,
                         'http_proxy': self.endpoint, 'https_proxy': self.endpoint,
                         'NO_PROXY': '', 'no_proxy': '', 'MSI_ENDPOINT': self.endpoint,
                         'IDENTITY_ENDPOINT': self.endpoint, 'IDENTITY_HEADER': 'synthetic-qa'})
        self.requests_before = len(self.traffic)

    def tearDown(self):
        self.assertEqual(len(self.traffic), self.requests_before, 'preflight attempted HTTP/authentication through the tripwire')

    def execute(self, args):
        return subprocess.run([str(self.executable)] + args, cwd=self.directory, env=self.env,
                              capture_output=True, text=True, timeout=15)

    def snapshot(self):
        return {str(p.relative_to(self.directory)): 'directory' if p.is_dir() else hashlib.sha256(p.read_bytes()).hexdigest()
                for p in self.directory.rglob('*')}

    def test_help_version_and_documented_output_contract(self):
        cases = [([], 'Cloud Assess'), (['--help'], 'cloud-assess'),
                 (['scan', '--help'], '--redact-subscription-ids'), (['--version'], 'cloud-assess version ' + EXPECTED_VERSION)]
        for args, expected in cases:
            with self.subTest(args=args):
                result = self.execute(args)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertIn(expected, result.stdout)
                self.assertEqual(result.stderr, '')
        help_text = self.execute(['scan', '--help']).stdout
        for flag in ('--json', '--xlsx', '--csv', '--sarif', '--stdout', '--filters', '--fail-on', '--stages'):
            self.assertIn(flag, help_text)
        self.assertIn('SARIF retains stable resource identities', help_text)
        self.assertEqual(list(self.directory.iterdir()), [self.executable], 'help/version created files')

    def test_preflight_rejections_preserve_reports(self):
        malformed = self.directory / 'malformed.yml'
        malformed.write_text('assessment: [unterminated\n')
        invalid_scope = self.directory / 'invalid-scope.yml'
        invalid_scope.write_text('assessment:\n  include:\n    resourceGroups: [not-a-resource-group-id]\n')
        valid = self.directory / 'valid.yml'
        valid.write_text('assessment:\n  include:\n    tags:\n      Environment: dev\n')
        cases = [
            ('unknown-command', ['no-such-command'], 'unknown command'),
            ('unknown-flag', ['scan', '--no-such-flag'], 'unknown flag'),
            ('unexpected-position', ['scan', 'extra'], 'unknown command'),
            ('malformed-boolean', ['scan', '--json=perhaps'], 'invalid argument'),
            ('unknown-stage', ['scan', '--stages', 'nonexistent'], 'unknown stage name'),
            ('mandatory-graph', ['scan', '--stages=-graph'], 'graph stage is mandatory'),
            ('deferred-plugin', ['scan', '--stages', 'plugin'], 'plugin stage is not available'),
            ('malformed-stage-param', ['scan', '--stage-param', 'broken'], 'stage param must be in the form'),
            ('unknown-stage-param', ['scan', '--stage-param', 'nonexistent.key=value'], 'unknown stage:'),
            ('unknown-option', ['scan', '--stage-param', 'plugin.unknown=value'], 'unknown option'),
            ('invalid-impact', ['scan', '--fail-on', 'Critical'], 'invalid fail-on impact'),
            ('missing-filter', ['scan', '--filters', str(self.directory / 'missing.yml')], 'read filter file'),
            ('malformed-filter', ['scan', '--filters', str(malformed)], 'parse filter YAML'),
            ('invalid-filter-scope', ['scan', '--filters', str(invalid_scope)], 'validate filter file'),
            ('valid-filter-deferred-stage', ['scan', '--filters', str(valid), '--stages', 'plugin'], 'plugin stage is not available'),
        ]
        base = self.directory / 'existing report'
        for suffix in ('.json', '.xlsx', '.sarif', '_Recommendations.csv'):
            Path(str(base) + suffix).write_bytes(b'retained report bytes\x00\xff')
        snapshot = self.snapshot()
        for name, args, message in cases:
            with self.subTest(case=name):
                if args[0] == 'scan':
                    args = args + ['--output-name', str(base)]
                    if name != 'malformed-boolean':
                        args += ['--json', '--csv', '--sarif', '--stdout']
                completed = self.execute(args)
                self.assertEqual(completed.returncode, 1, completed.stdout + completed.stderr)
                self.assertIn(message, completed.stderr)
                self.assertEqual(completed.stdout, '', 'rejected scan emitted a report to stdout')
                self.assertEqual(self.snapshot(), snapshot)
                self.assertEqual(len(self.traffic), self.requests_before)
                self.assertFalse((self.directory / 'isolated-azure-config').exists())

    def test_compiled_module_evidence_matches_inventory(self):
        completed = subprocess.run([GO, 'version', '-m', '-json', str(BINARY)], capture_output=True, text=True, timeout=30)
        self.assertEqual(completed.returncode, 0, completed.stderr)
        info = json.loads(completed.stdout)
        manifest = json.loads((ROOT / 'docs/dependencies/inventory.json').read_text())
        self.assertEqual(info['GoVersion'], manifest['goVersion'])
        self.assertEqual(info['Path'], 'github.com/DeBoX85/Cloud-Assess/cmd/cloud-assess')
        self.assertEqual(info['Main']['Path'], 'github.com/DeBoX85/Cloud-Assess')
        settings = {s['Key']: s['Value'] for s in info['Settings']}
        target = settings['GOOS'] + '/' + settings['GOARCH']
        expected_os = 'windows' if os.name == 'nt' else 'linux'
        self.assertEqual(settings['GOOS'], expected_os)
        self.assertEqual(settings['GOARCH'], 'amd64')
        self.assertEqual(settings['CGO_ENABLED'], '0')
        expected = {m['path']: (m['version'], m['sum']) for m in manifest['modules'] if target in m['cliTargets']}
        actual = {}
        for dep in info['Deps']:
            self.assertNotIn('Replace', dep, 'unreviewed compiled replacement')
            self.assertNotIn(dep['Path'], actual, 'duplicate compiled module')
            actual[dep['Path']] = (dep['Version'], dep['Sum'])
        self.assertEqual(actual, expected)
        print('Verified built artifact:', target, 'modules:', len(actual), 'sha256:', hashlib.sha256(BINARY.read_bytes()).hexdigest())


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--binary', required=True)
    parser.add_argument('--expected-version', default='dev')
    parser.add_argument('--go', default=os.environ.get('QA_GO', 'go'))
    args, rest = parser.parse_known_args()
    BINARY = Path(args.binary).resolve(strict=True)
    GO = args.go
    EXPECTED_VERSION = args.expected_version
    unittest.main(argv=[__file__] + rest)
