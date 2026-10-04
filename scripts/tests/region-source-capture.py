#!/usr/bin/env python3
"""Exercise the real retained-byte guard without running source/network checkout."""
import ast
from pathlib import Path
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
tree = ast.parse((ROOT / 'scripts/capture-region-source.py').read_text(encoding='utf-8'))
functions = [node for node in tree.body
             if isinstance(node, ast.FunctionDef) and node.name == 'verify_capture_bytes']
if len(functions) != 1:
    raise SystemExit('Capture byte guard function missing or ambiguous')
module = ast.Module(body=functions, type_ignores=[])
namespace = {}
exec(compile(module, '<actual source byte guard>', 'exec'), namespace)
guard = namespace['verify_capture_bytes']


class CaptureBytesTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix='cloud-assess-capture-bytes-')
        self.addCleanup(self.temporary.cleanup)
        self.retained = Path(self.temporary.name) / 'synthetic.json'
        self.raw = '{"label":"\u00a0\u212a\U0001f600"}\n'.encode('utf-8')
        self.retained.write_bytes(self.raw)

    def test_exact_bytes(self):
        self.assertIsNone(guard(self.raw, self.retained))

    def test_changed_bytes(self):
        with self.assertRaises(SystemExit) as error:
            guard(b'{"label":"private"}\n', self.retained)
        self.assertEqual(str(error.exception), 'Retained source capture byte mismatch: synthetic.json')
        self.assertNotIn('private', str(error.exception))

    def test_equal_json_different_bytes(self):
        with self.assertRaises(SystemExit):
            guard(self.raw.rstrip(b'\n'), self.retained)

    def test_missing_capture(self):
        self.retained.unlink()
        with self.assertRaises(FileNotFoundError):
            guard(self.raw, self.retained)


suite = unittest.defaultTestLoader.loadTestsFromTestCase(CaptureBytesTests)
baseline = unittest.TextTestRunner(verbosity=2).run(suite)
if not baseline.wasSuccessful():
    raise SystemExit(1)

# A syntactically valid disabled guard must fail the named mismatch assertion.
checks = [node for node in ast.walk(functions[0]) if isinstance(node, ast.If)]
if len(checks) != 1:
    raise SystemExit('Capture guard mutation anchor changed')
original = checks[0].test
try:
    checks[0].test = ast.Constant(value=False)
    ast.fix_missing_locations(module)
    mutated = {}
    exec(compile(module, '<compiling disabled source byte guard>', 'exec'), mutated)
    guard = mutated['verify_capture_bytes']
    result = unittest.TestResult()
    CaptureBytesTests('test_changed_bytes').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Capture byte guard survived or failed outside named assertion')
finally:
    checks[0].test = original
    restored = {}
    exec(compile(module, '<restored actual source byte guard>', 'exec'), restored)
    guard = restored['verify_capture_bytes']

restored_result = unittest.TestResult()
unittest.defaultTestLoader.loadTestsFromTestCase(CaptureBytesTests).run(restored_result)
if not restored_result.wasSuccessful():
    raise SystemExit('Restored capture byte guard baseline failed')
print('Region source capture compiling mutation rejected: retained-byte-guard; restored four assertions passed')
