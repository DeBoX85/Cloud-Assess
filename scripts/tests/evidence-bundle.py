#!/usr/bin/env python3
"""Independent integrity, admission and actual CLI failure cases for evidence."""
import base64
import copy
import gzip
import hashlib
import importlib.util
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

sys.dont_write_bytecode = True
ROOT = Path(__file__).resolve().parents[2]
SCRIPT = ROOT / 'scripts/verify-evidence-bundle.py'
spec = importlib.util.spec_from_file_location('evidence_verifier', SCRIPT)
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def fixture(payload=None):
    if payload is None:
        payload = {'format': 'cloud-assess-evidence-v1', 'members': [
            {'path': 'run/job.log', 'text': 'completed\n'},
            {'path': 'review.json', 'text': '{"verified":false}\n'}]}
    decoded = json.dumps(payload).encode()
    bundle = base64.b64encode(gzip.compress(decoded, mtime=0)) + b'\n'
    manifest = {'format': 'cloud-assess-evidence-v1', 'classification': 'reviewed-synthetic-ci-local',
                'bundle_sha256': sha(bundle), 'decoded_sha256': sha(decoded),
                'revisions': {'candidate': '1' * 40}, 'members': [
                    {'path': 'run/job.log', 'bytes': 10, 'sha256': sha(b'completed\n')},
                    {'path': 'review.json', 'bytes': 19, 'sha256': sha(b'{"verified":false}\n')}]}
    return bundle, manifest


class EvidenceTests(unittest.TestCase):
    def assert_rejects(self, bundle, manifest):
        raw_manifest = manifest if isinstance(manifest, bytes) else json.dumps(manifest).encode()
        with self.assertRaises(ValueError):
            module.verify(bundle, raw_manifest)
        with tempfile.TemporaryDirectory(prefix='ca-evidence-negative-') as directory:
            path = Path(directory)
            (path / 'bundle').write_bytes(bundle)
            (path / 'manifest').write_bytes(raw_manifest)
            result = subprocess.run([sys.executable, str(SCRIPT), '--bundle', str(path/'bundle'), '--manifest', str(path/'manifest')], capture_output=True, timeout=10)
            self.assertEqual(result.returncode, 1, result.stderr)
            self.assertIn(b'Evidence verification failed', result.stderr)
            self.assertEqual({p.name for p in path.iterdir()}, {'bundle', 'manifest'})

    def test_valid_independent_fixture(self):
        bundle, manifest = fixture()
        self.assertEqual(module.verify(bundle, json.dumps(manifest).encode())['members'], 2)

    def test_corrupt_bundle_and_manifest(self):
        bundle, manifest = fixture()
        for candidate in (bundle[:-3] + b'xx\n', b''):
            self.assert_rejects(candidate, manifest)
        for key, value in [('decoded_sha256', '0' * 64), ('classification', 'private-azure'), ('format', 'unknown')]:
            bad = copy.deepcopy(manifest); bad[key] = value
            self.assert_rejects(bundle, bad)

    def test_member_tampering_and_inventory(self):
        original = {'format': 'cloud-assess-evidence-v1', 'members': [
            {'path': 'run/job.log', 'text': 'completed\n'}, {'path': 'review.json', 'text': '{"verified":false}\n'}]}
        variants = []
        bad = copy.deepcopy(original); bad['members'][0]['text'] = 'corrupted\n'; variants.append(bad)
        bad = copy.deepcopy(original); bad['members'].pop(); variants.append(bad)
        bad = copy.deepcopy(original); bad['members'].append({'path': 'extra', 'text': ''}); variants.append(bad)
        bad = copy.deepcopy(original); bad['members'][1] = bad['members'][0]; variants.append(bad)
        for bad in variants:
            bundle, manifest = fixture(bad)
            self.assert_rejects(bundle, manifest)
        bundle, manifest = fixture()
        for change in ('length', 'hash', 'duplicate', 'boolean-length', 'revision'):
            bad = copy.deepcopy(manifest)
            if change == 'length': bad['members'][0]['bytes'] = 11
            if change == 'hash': bad['members'][0]['sha256'] = '0' * 64
            if change == 'duplicate': bad['members'][1] = bad['members'][0]
            if change == 'boolean-length': bad['members'][0]['bytes'] = True
            if change == 'revision': bad['revisions']['candidate'] = 'dev'
            self.assert_rejects(bundle, bad)

    def test_unsafe_names_and_duplicate_json(self):
        for name in ('../x', '/tmp/x', 'a/../b', 'C:/x', 'a\\b', './x', 'a//b', 'a/', '.', ''):
            payload = {'format': 'cloud-assess-evidence-v1', 'members': [
                {'path': name, 'text': 'completed\n'}, {'path': 'review.json', 'text': '{"verified":false}\n'}]}
            bundle, manifest = fixture(payload)
            self.assert_rejects(bundle, manifest)
        bundle, manifest = fixture()
        raw = json.dumps(manifest).encode().replace(b'{', b'{"format":"duplicate",', 1)
        self.assert_rejects(bundle, raw)

    def test_expansion_bound(self):
        decoded = b'x' * (16 * 1024 * 1024 + 1)
        bundle = base64.b64encode(gzip.compress(decoded, mtime=0)) + b'\n'
        _, manifest = fixture(); manifest['bundle_sha256'] = sha(bundle); manifest['decoded_sha256'] = sha(decoded)
        self.assert_rejects(bundle, manifest)

    def test_actual_cli_success_failure_and_no_extraction(self):
        with tempfile.TemporaryDirectory(prefix='ca-evidence-test-') as directory:
            path = Path(directory); bundle, manifest = fixture()
            bp, mp = path/'bundle', path/'manifest'
            bp.write_bytes(bundle); mp.write_text(json.dumps(manifest))
            command = [sys.executable, str(SCRIPT), '--bundle', str(bp), '--manifest', str(mp)]
            result = subprocess.run(command, capture_output=True, timeout=10)
            self.assertEqual(result.returncode, 0, result.stderr)
            bp.write_bytes(b'bad archive')
            result = subprocess.run(command, capture_output=True, timeout=10)
            self.assertEqual(result.returncode, 1)
            self.assertIn(b'Evidence verification failed', result.stderr)
            self.assertEqual({p.name for p in path.iterdir()}, {'bundle', 'manifest'})

    def test_bundle_checkout_preserves_bytes(self):
        path = 'docs/evidence/20261017/pr150-pr151.json.gz.base64'
        result = subprocess.run(['git', '-c', 'core.autocrlf=true', 'check-attr', 'text', '--', path], cwd=ROOT, capture_output=True, text=True, timeout=10)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), path + ': text: unset')

    def test_committed_real_evidence(self):
        base = ROOT / 'docs/evidence/20261017'
        result = subprocess.run([sys.executable, str(SCRIPT), '--bundle', str(base/'pr150-pr151.json.gz.base64'), '--manifest', str(base/'manifest.json')], capture_output=True, timeout=10)
        self.assertEqual(result.returncode, 0, result.stderr)


if __name__ == '__main__':
    unittest.main()
