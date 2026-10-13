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
manifest_nodes = [n for n in tree.body if isinstance(n, ast.Assign) and any(isinstance(t, ast.Name) and t.id == 'FILES' for t in n.targets)]
if len(manifest_nodes) != 1:
    raise SystemExit('Capture file manifest missing or ambiguous')
capture_files = ast.literal_eval(manifest_nodes[0].value)


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


    def test_required_manifest(self):
        self.assertEqual(capture_files, ('source-aux-inputs.json','source-aux-outputs.json','source-inventory-inputs.json','source-inventory-outputs.json','source-availability-inputs.json','source-availability-outputs.json','source-latency-inputs.json','source-latency-outputs.json','source-latency-data.json','source-cost-enrichment-inputs.json','source-cost-enrichment-outputs.json','source-quota-inputs.json','source-quota-outputs.json','source-vm-quota-inputs.json','source-vm-quota-outputs.json','source-reservation-inputs.json','source-reservation-outputs.json','source-provider-locations-inputs.json','source-provider-locations-outputs.json','source-compute-sku-inputs.json','source-compute-sku-outputs.json','source-other-sku-inputs.json','source-other-sku-outputs.json'))

    def test_all_files_unconditionally_guarded(self):
        loops = [n for n in ast.walk(tree) if isinstance(n, ast.For) and any(isinstance(c, ast.Call) and isinstance(c.func, ast.Name) and c.func.id == 'verify_capture_bytes' for c in ast.walk(n))]
        self.assertEqual(len(loops), 1)
        self.assertIsInstance(loops[0].iter, ast.Name)
        self.assertEqual(loops[0].iter.id, 'FILES')
        self.assertEqual(len(loops[0].body), 1)
        self.assertIsInstance(loops[0].body[0], ast.Expr)
        self.assertIsInstance(loops[0].body[0].value, ast.Call)
        self.assertEqual(loops[0].body[0].value.func.id, 'verify_capture_bytes')

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
print('Region source capture compiling mutation rejected: retained-byte-guard; restored six assertions passed')

# A compiling omission of the provider records must fail the manifest assertion.
original_files = capture_files
try:
    capture_files = tuple(name for name in capture_files if not name.startswith('source-provider-locations-'))
    result = unittest.TestResult()
    CaptureBytesTests('test_required_manifest').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Provider records omission did not fail named manifest assertion')
finally:
    capture_files = original_files

# A compiling guarded subset must fail the actual-loop coverage assertion.
loops = [n for n in ast.walk(tree) if isinstance(n, ast.For) and any(isinstance(c, ast.Call) and isinstance(c.func, ast.Name) and c.func.id == 'verify_capture_bytes' for c in ast.walk(n))]
original_iterator = loops[0].iter
try:
    loops[0].iter = ast.Tuple(elts=[ast.Constant(value=name) for name in original_files if not name.startswith('source-provider-locations-')], ctx=ast.Load())
    ast.fix_missing_locations(tree)
    compile(tree, '<compiling source provider guard subset>', 'exec')
    result = unittest.TestResult()
    CaptureBytesTests('test_all_files_unconditionally_guarded').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Provider guard subset did not fail named coverage assertion')
finally:
    loops[0].iter = original_iterator

restored_result = unittest.TestResult()
unittest.defaultTestLoader.loadTestsFromTestCase(CaptureBytesTests).run(restored_result)
if not restored_result.wasSuccessful():
    raise SystemExit('Restored complete capture manifest/loop baseline failed')
print('Region source capture compiling mutation rejected: provider-manifest-omission; restored six assertions passed')
print('Region source capture compiling mutation rejected: provider-guard-subset; restored six assertions passed')


# Preserve the provider omission control above and independently omit Compute records.
try:
    capture_files = tuple(name for name in original_files if not name.startswith('source-compute-sku-'))
    result = unittest.TestResult()
    CaptureBytesTests('test_required_manifest').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Compute SKU records omission did not fail named manifest assertion')
finally:
    capture_files = original_files

# A different compiling subset expression explicitly excludes only Compute records.
try:
    loops[0].iter = ast.Tuple(elts=[ast.Constant(value=name) for name in original_files
                                  if not name.startswith('source-compute-sku-')], ctx=ast.Load())
    ast.fix_missing_locations(tree)
    compile(tree, '<compiling source Compute guard subset>', 'exec')
    result = unittest.TestResult()
    CaptureBytesTests('test_all_files_unconditionally_guarded').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Compute SKU guard subset did not fail named coverage assertion')
finally:
    loops[0].iter = original_iterator

restored_result = unittest.TestResult()
unittest.defaultTestLoader.loadTestsFromTestCase(CaptureBytesTests).run(restored_result)
if not restored_result.wasSuccessful():
    raise SystemExit('Restored complete Compute/provider capture manifest/loop baseline failed')
print('Region source capture compiling mutation rejected: compute-sku-manifest-omission; restored six assertions passed')
print('Region source capture compiling mutation rejected: compute-sku-guard-subset; restored six assertions passed')

# Keep both earlier family controls and add independent remaining-SKU exclusions.
try:
    capture_files = tuple(name for name in original_files if not name.startswith('source-other-sku-'))
    result = unittest.TestResult()
    CaptureBytesTests('test_required_manifest').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Other SKU records omission did not fail named manifest assertion')
finally:
    capture_files = original_files

try:
    loops[0].iter = ast.Tuple(elts=[ast.Constant(value=name) for name in original_files
                                  if not name.startswith('source-other-sku-')], ctx=ast.Load())
    ast.fix_missing_locations(tree)
    compile(tree, '<compiling remaining source SKU guard subset>', 'exec')
    result = unittest.TestResult()
    CaptureBytesTests('test_all_files_unconditionally_guarded').run(result)
    if len(result.failures) != 1 or result.errors or result.testsRun != 1:
        raise SystemExit('Other SKU guard subset did not fail named coverage assertion')
finally:
    loops[0].iter = original_iterator

restored_result = unittest.TestResult()
unittest.defaultTestLoader.loadTestsFromTestCase(CaptureBytesTests).run(restored_result)
if not restored_result.wasSuccessful():
    raise SystemExit('Restored complete remaining SKU capture manifest/loop baseline failed')
print('Region source capture compiling mutation rejected: other-sku-manifest-omission; restored six assertions passed')
print('Region source capture compiling mutation rejected: other-sku-guard-subset; restored six assertions passed')
