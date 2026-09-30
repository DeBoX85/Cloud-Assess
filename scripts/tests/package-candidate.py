#!/usr/bin/env python3
import argparse
import copy
import importlib.util
import json
import os
from pathlib import Path
import stat
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch
import zipfile

sys.dont_write_bytecode = True
ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('package_candidate', ROOT / 'scripts/package-candidate.py')
pkg = importlib.util.module_from_spec(spec)
spec.loader.exec_module(pkg)
BINARY = None
GO = None


class PackageIntegrityTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.directory = Path(self.temp.name)
        self.root = 'cloud-assess-dev-linux-amd64'
        self.payload = {p: b'fixture' for p in pkg.COMMON | {'cloud-assess'}}
        self.manifest = {'schemaVersion': 1, 'status': 'development-candidate-not-release-approved',
                         'version': 'dev', 'target': 'linux/amd64',
                         'files': {p: pkg.sha(data) for p, data in self.payload.items()}}

    def archive(self, mutate=None):
        payload = dict(self.payload)
        payload['PACKAGE_MANIFEST.json'] = json.dumps(self.manifest).encode()
        rows = []
        for name, data in sorted(payload.items()):
            item = zipfile.ZipInfo(self.root + '/' + name)
            item.create_system = 3
            item.external_attr = (stat.S_IFREG | (0o755 if name == 'cloud-assess' else 0o644)) << 16
            rows.append((item, data))
        if mutate:
            mutate(rows)
        path = self.directory / 'candidate.zip'
        with zipfile.ZipFile(path, 'w') as z:
            for item, data in rows:
                z.writestr(item, data)
        checksum = self.directory / 'candidate.zip.sha256'
        checksum.write_text(pkg.sha(path.read_bytes()) + '  candidate.zip\n')
        return path, checksum

    def reject(self, mutate, message):
        archive, checksum = self.archive(mutate)
        destination = self.directory / 'extract'
        with self.assertRaisesRegex(ValueError, message):
            pkg.verify_extract(archive, checksum, destination)
        self.assertFalse(destination.exists(), 'invalid package wrote extraction files')

    def test_valid_payload_and_executable_permission(self):
        archive, checksum = self.archive()
        executable = pkg.verify_extract(archive, checksum, self.directory / 'clean extraction')
        self.assertEqual(executable.read_bytes(), b'fixture')
        self.assertEqual((executable.parent / 'LICENSE').read_bytes(), b'fixture')
        if os.name != 'nt':
            self.assertEqual(executable.stat().st_mode & 0o777, 0o755)

    def test_bad_archive_checksum_before_writes(self):
        archive, checksum = self.archive()
        checksum.write_text('0' * 64 + '  candidate.zip\n')
        with self.assertRaisesRegex(ValueError, 'archive checksum'):
            pkg.verify_extract(archive, checksum, self.directory / 'extract')
        self.assertFalse((self.directory / 'extract').exists())

    def test_payload_corruption_even_with_updated_archive_checksum(self):
        def mutate(rows):
            index = next(i for i, r in enumerate(rows) if r[0].filename.endswith('/LICENSE'))
            rows[index] = (rows[index][0], b'corrupt license')
        self.reject(mutate, 'payload checksum')

    def test_missing_notice_rejected(self):
        self.reject(lambda rows: rows.pop(next(i for i, r in enumerate(rows) if r[0].filename.endswith('/DEPENDENCY_NOTICES.md'))), 'member count')

    def test_parent_traversal_rejected(self):
        self.reject(lambda rows: setattr(rows[0][0], 'filename', '../escape'), 'unsafe archive path')

    def test_absolute_path_rejected(self):
        self.reject(lambda rows: setattr(rows[0][0], 'filename', '/escape'), 'unsafe archive path')

    def test_duplicate_casefold_member_rejected(self):
        self.reject(lambda rows: setattr(rows[1][0], 'filename', rows[0][0].filename.upper()), 'duplicate')

    def test_symlink_rejected(self):
        self.reject(lambda rows: setattr(rows[0][0], 'external_attr', (stat.S_IFLNK | 0o777) << 16), 'nonregular')

    def test_excess_permissions_rejected(self):
        self.reject(lambda rows: setattr(rows[0][0], 'external_attr', (stat.S_IFREG | 0o777) << 16), 'permissions')

    def test_existing_destination_retained(self):
        archive, checksum = self.archive()
        destination = self.directory / 'extract'
        destination.mkdir()
        (destination / 'existing').write_bytes(b'retain')
        with self.assertRaisesRegex(ValueError, 'new destination'):
            pkg.verify_extract(archive, checksum, destination)
        self.assertEqual((destination / 'existing').read_bytes(), b'retain')

    def test_invalid_manifest_target(self):
        self.manifest['target'] = 'linux/arm64'
        self.reject(None, 'unsupported candidate manifest')

    def test_compiled_provenance_and_dependency_rejections(self):
        inventory = {'goVersion': 'go1.26.8', 'modules': [{'path': 'example/module', 'version': 'v1', 'sum': 'h1:fixture', 'cliTargets': ['linux/amd64']}]}
        info = {'GoVersion': 'go1.26.8', 'Path': 'github.com/DeBoX85/Cloud-Assess/cmd/cloud-assess',
                'Main': {'Path': 'github.com/DeBoX85/Cloud-Assess'}, 'Deps': [{'Path': 'example/module', 'Version': 'v1', 'Sum': 'h1:fixture'}],
                'Settings': [{'Key': k, 'Value': v} for k, v in {'vcs.revision': 'a' * 40, 'vcs.modified': 'false', 'CGO_ENABLED': '0', 'GOOS': 'linux', 'GOARCH': 'amd64', '-trimpath': 'true'}.items()]}
        pkg.validate_info(info, inventory, 'a' * 40, 'linux/amd64')
        for key, value in [('vcs.revision', 'b' * 40), ('vcs.modified', 'true'), ('CGO_ENABLED', '1'), ('-trimpath', 'false')]:
            with self.subTest(setting=key):
                bad = copy.deepcopy(info)
                next(s for s in bad['Settings'] if s['Key'] == key)['Value'] = value
                with self.assertRaises(ValueError):
                    pkg.validate_info(bad, inventory, 'a' * 40, 'linux/amd64')
        bad = copy.deepcopy(info)
        bad['Deps'][0]['Sum'] = 'h1:wrong'
        with self.assertRaisesRegex(ValueError, 'dependencies'):
            pkg.validate_info(bad, inventory, 'a' * 40, 'linux/amd64')


class ActualCandidateTests(unittest.TestCase):
    def test_repeat_archive_install_and_offline_cli(self):
        if BINARY is None:
            self.skipTest('provide --binary for actual candidate integration')
        with tempfile.TemporaryDirectory(prefix='cloud-assess-package-qa-') as temp:
            root = Path(temp)
            archive, checksum = pkg.package(BINARY, 'dev', root / 'first', GO)
            other, _ = pkg.package(BINARY, 'dev', root / 'second', GO)
            self.assertEqual(archive.read_bytes(), other.read_bytes(), 'same inputs produced different ZIPs')
            executable = pkg.verify_extract(archive, checksum, root / 'isolated install with spaces')
            self.assertEqual(executable.read_bytes(), BINARY.read_bytes())
            self.assertEqual((executable.parent / 'DEPENDENCY_NOTICES.md').read_bytes(), pkg.source('docs/dependencies/NOTICES.md'))
            completed = subprocess.run([sys.executable, str(ROOT / 'scripts/tests/built-cli.py'), '--binary', str(executable), '--go', GO], capture_output=True, text=True, timeout=120)
            self.assertEqual(completed.returncode, 0, completed.stdout + completed.stderr)
            with self.assertRaisesRegex(ValueError, 'version does not match'):
                pkg.package(BINARY, 'deliberately-wrong', root / 'wrong-version', GO)
            self.assertFalse((root / 'wrong-version').exists())
            with self.assertRaisesRegex(ValueError, 'overwrite'):
                pkg.package(BINARY, 'dev', root / 'first', GO)
            self.assertEqual(archive.read_bytes(), other.read_bytes())
            original_link = os.link
            def race(source, target):
                if str(target).endswith('.zip'):
                    Path(target).write_bytes(b'concurrent publisher bytes')
                return original_link(source, target)
            with patch.object(pkg.os, 'link', side_effect=race):
                with self.assertRaises(FileExistsError):
                    pkg.package(BINARY, 'dev', root / 'race', GO)
            raced = next((root / 'race').glob('*.zip'))
            self.assertEqual(raced.read_bytes(), b'concurrent publisher bytes')
            self.assertFalse(list((root / 'race').glob('*.sha256')))
            print('Candidate ZIP verified:', archive.name, 'sha256:', pkg.sha(archive.read_bytes()))


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--binary')
    parser.add_argument('--go', default=os.environ.get('QA_GO', 'go'))
    args, rest = parser.parse_known_args()
    BINARY = Path(args.binary).resolve(strict=True) if args.binary else None
    GO = args.go
    unittest.main(argv=[__file__] + rest)
