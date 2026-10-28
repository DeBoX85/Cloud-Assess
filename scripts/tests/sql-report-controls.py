#!/usr/bin/env python3
"""SQL report compiling behavioral reversals; own isolated checkout only."""
from mutation_assertions import intended_failure, run_go_test
import hashlib,json,subprocess,sys
from pathlib import Path
root=Path(__file__).resolve().parents[2]
plan=json.loads(Path(__file__).with_name('sql-report-controls.plan.json').read_text(encoding='utf-8'))
assert plan['format']=='sql-report-controls-v1' and len(plan['mutations'])==15
original={m['file']:(root/m['file']).read_bytes() for m in plan['mutations']}
assert len({m['name'] for m in plan['mutations']})==15

def test(name):
 return run_go_test(['go','test','./internal/plugins/region','-run','^'+name+'$','-count=1'],cwd=root,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,encoding='utf-8',errors='strict')
baseline=test('TestSQL.*');assert baseline.returncode==0,('initial healthy baseline failed',baseline.stdout)
try:
 for m in plan['mutations']:
  path=root/m['file'];before=original[m['file']];source=before.decode('utf-8').replace('\r\n','\n')
  assert source.count(m['old'])==m['count'],('control anchor',m['name'])
  try:
   path.write_bytes(source.replace(m['old'],m['replace']).encode('utf-8'))
   compiled=test('');assert compiled.returncode==0,('uncompilable reversal',m['name'],compiled.stdout)
   result=test(m['test'])
   assert intended_failure(result,m['test']),('not original named assertion',m['name'],result.stdout)
  finally:
   path.write_bytes(before)
  assert path.read_bytes()==before,('source restoration',m['name'])
  baseline=test(m['test']);assert baseline.returncode==0,('restored baseline failed',m['name'],baseline.stdout)
  print('SQL report mutation rejected: '+m['name']+'; compilation succeeded; named assertion failed; exact source restored; baseline passed',flush=True)
 baseline=test('TestSQL.*');assert baseline.returncode==0,('complete restored SQL baseline',baseline.stdout)
 for name,raw in original.items():assert (root/name).read_bytes()==raw,('final byte restoration',name)
 print('SQL report controls: all15 compiling reversals, exact original byte restoration and complete passing SQL baseline',flush=True)
finally:
 for name,raw in original.items():(root/name).write_bytes(raw)
