"""Require the independent preparation specification; no coordinator execution."""
import hashlib,importlib.util,json,sys
from pathlib import Path
sys.dont_write_bytecode=True
ROOT=Path(__file__).resolve().parents[2]
DIRECTORY=ROOT/'docs/evidence/20261028/coordinator'
manifest=json.loads((DIRECTORY/'integrity.json').read_text())
assert manifest['format']=='cloud-assess-coordinator-preparation-v1'
assert manifest['base']=='08a220a72edd902eb3f7eafc14b8e95f64a46f12'
for row in manifest['files']:
 raw=(ROOT/row['path']).read_bytes().replace(b'\r\n',b'\n')
 assert len(raw)==row['bytes'] and hashlib.sha256(raw).hexdigest()==row['sha256'],row['path']
spec=importlib.util.spec_from_file_location('coordinator_scheduling_reference',DIRECTORY/'check_scheduling.py');module=importlib.util.module_from_spec(spec);spec.loader.exec_module(module)
before=(DIRECTORY/'scheduling-verification.json').read_bytes()
module.main()
assert (DIRECTORY/'scheduling-verification.json').read_bytes().replace(b'\r\n',b'\n')==before.replace(b'\r\n',b'\n')
result=json.loads(before);assert result['literal_cases']==40 and result['negative_controls']==8 and result['actual_product_execution'] is False
print('COORDINATOR_SCHEDULING_SPEC_JSON '+json.dumps({'literal_cases':40,'negative_controls':8,'bound_files':len(manifest['files']),'frozen_source_hashes_verified':True,'actual_product_execution':False,'scope':'preparation reference model, not future coordinator enforcement'},separators=(',',':')))
