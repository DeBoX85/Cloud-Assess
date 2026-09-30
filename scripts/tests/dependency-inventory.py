#!/usr/bin/env python3
import importlib.util
import json
import os
from pathlib import Path
import tempfile
import sys
sys.dont_write_bytecode = True
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('inventory', Path(__file__).resolve().parents[1] / 'dependency-inventory.py')
inventory = importlib.util.module_from_spec(spec)
spec.loader.exec_module(inventory)


class InventoryTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.module = self.root / 'module'
        self.module.mkdir()
        (self.module / 'LICENSE').write_bytes(b'MIT fixture\r\nCopyright test\r\n')
        (self.module / 'NOTICE').write_text('Additional attribution\n')
        self.goroot = self.root / 'goroot'
        self.goroot.mkdir()
        (self.goroot / 'LICENSE').write_text('BSD fixture\n')
        (self.goroot / 'PATENTS').write_text('Patent grant fixture\n')
        for p in ('go.mod', 'go.sum', 'internal/rules/embedded.go', 'internal/skus/skus.go', 'internal/skus/known_skus.yaml', 'LICENSE', 'THIRD_PARTY_LICENSES.md'):
            path = self.root / p
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text('fixture\n')
        self.aprl = self.root / 'internal/rules/upstream/aprl'
        self.aprl.mkdir(parents=True)
        (self.aprl / 'LICENSE').write_text('APRL fixture\n')
        (self.root / 'internal/rules/provenance.go').write_text('\n'.join(k + ' = "' + 'a' * 40 + '"' for k in ('ReferenceRepositoryCommit', 'APRLCommit', 'AORSnapshotTree', 'CustomRulesSnapshotTree')))
        self.calls = []
        self.replace = False
        self.version = 'go1.26.8'

    def command(self, args, cwd=None, env=None):
        self.calls.append((args, cwd, env))
        if args[0] == 'git':
            return '' if 'status' in args else 'a' * 40 + '\n'
        if args[1:3] == ['list', '-m']:
            module = {'Path': 'example.test/module', 'Version': 'v1.0.0'}
            if self.replace:
                module['Replace'] = {'Dir': '/unreviewed'}
            return json.dumps(module)
        if args[1] == 'list':
            return json.dumps({'Module': {'Path': 'example.test/module'}})
        if args[1:3] == ['mod', 'download']:
            return json.dumps({'Path': 'example.test/module', 'Version': 'v1.0.0', 'Dir': str(self.module), 'Sum': 'h1:fixture', 'GoModSum': 'h1:modfixture'})
        if args[-1] == 'GOROOT':
            return str(self.goroot)
        if args[-1] == 'GOVERSION':
            return self.version
        raise AssertionError(args)

    def collect(self):
        with patch.object(inventory, 'ROOT', self.root), patch.object(inventory, 'run', self.command):
            return inventory.collect('go')

    def test_exact_notice_bytes_and_target_accounting(self):
        manifest, notices = self.collect()
        row = json.loads(manifest)['modules'][0]
        self.assertEqual(row['cliTargets'], ['linux/amd64', 'windows/amd64'])
        self.assertIn('MIT fixture\r\nCopyright test\r\n', notices)
        self.assertEqual(row['notices'][0]['sha256'], inventory.digest((self.module / 'LICENSE').read_bytes()))
        self.assertIn('Patent grant fixture', notices)
        download = next(c for c in self.calls if c[0][1:3] == ['mod', 'download'])
        self.assertNotEqual(download[1], self.root)
        self.assertEqual(download[2]['GOWORK'], 'off')

    def test_missing_license_rejected_even_with_notice(self):
        (self.module / 'LICENSE').unlink()
        with self.assertRaisesRegex(ValueError, 'no license evidence'):
            self.collect()

    def test_empty_license_rejected(self):
        (self.module / 'LICENSE').write_text(' ')
        with self.assertRaisesRegex(ValueError, 'empty notice'):
            self.collect()

    def test_symlink_rejected(self):
        (self.module / 'LICENSE').unlink()
        (self.module / 'LICENSE').symlink_to(self.module / 'NOTICE')
        with self.assertRaisesRegex(ValueError, 'symlink notice'):
            self.collect()

    def test_replacement_requires_provenance_policy(self):
        self.replace = True
        with self.assertRaisesRegex(ValueError, 'replaced/unversioned'):
            self.collect()

    def test_wrong_toolchain_rejected(self):
        self.version = 'go1.26.0'
        with self.assertRaisesRegex(ValueError, 'pinned go1.26.8'):
            self.collect()

    def test_failed_collection_does_not_overwrite_outputs(self):
        path = self.root / 'docs/dependencies/inventory.json'
        path.parent.mkdir(parents=True)
        path.write_text('retained evidence')
        with patch.object(inventory, 'ROOT', self.root), patch.object(inventory, 'collect', side_effect=ValueError('missing evidence')), patch('sys.argv', ['inventory']):
            with self.assertRaises(ValueError):
                inventory.main()
        self.assertEqual(path.read_text(), 'retained evidence')

    def test_check_rejects_stale_without_writing(self):
        path = self.root / 'docs/dependencies/inventory.json'
        path.parent.mkdir(parents=True)
        path.write_text('stale')
        with patch.object(inventory, 'ROOT', self.root), patch.object(inventory, 'collect', return_value=('new', 'notices')), patch('sys.argv', ['inventory', '--check']):
            with self.assertRaisesRegex(SystemExit, 'stale inventory'):
                inventory.main()
        self.assertEqual(path.read_text(), 'stale')
        self.assertFalse((path.parent / 'NOTICES.md').exists())

    def test_json_sequence_and_malformed_input(self):
        self.assertEqual(list(inventory.objects(' {"a":1}\n {"b":2} ')), [{'a': 1}, {'b': 2}])
        with self.assertRaises(json.JSONDecodeError):
            list(inventory.objects('{broken'))


if __name__ == '__main__':
    unittest.main()
