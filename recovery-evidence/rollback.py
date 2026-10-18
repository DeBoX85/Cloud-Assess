from pathlib import Path
import hashlib,json,os,subprocess,sys
os.environ['PYTHONDONTWRITEBYTECODE']='1'
repo=Path(sys.argv[1]).resolve();proof=Path(sys.argv[2]).resolve()
head='8c47d8105c7c3a6dec02f4a9b22958ab58b4682a';base='724504b7e7fb454289e02c110abe0ff7fcfb65a4'
def git(*args):return subprocess.check_output(['git',*args],cwd=repo).decode().strip()
assert git('rev-parse','HEAD')==head and git('rev-parse','HEAD^')==base and not git('status','--porcelain')
files=git('diff','--name-only',base,head).splitlines();assert len(files)==16
allowed={'internal/plugins/region/sql_report.go','internal/plugins/region/sql_report_contract_test.go','internal/plugins/region/sql_report_boundary_test.go','internal/plugins/region/sql_report_renderers_test.go','scripts/tests/sql-report-controls.py','scripts/tests/sql-report-controls.plan.json','scripts/ci/quality-plan.json','.github/workflows/test.yml'}
assert all(f in allowed or f.startswith('docs/') for f in files)
assert not any(f.startswith('internal/renderers/') for f in files)
original={f:(repo/f).read_bytes() for f in files}
retained=['internal/renderers/tables/tables.go','internal/renderers/tables/scope_redaction.go','internal/renderers/tables/scope_redaction_test.go','internal/renderers/tables/scope_redaction_records_test.go','scripts/tests/scope-redaction-controls.py','scripts/tests/scope-redaction-controls.plan.json']
renderer={f:(repo/f).read_bytes() for f in retained}
def test(label):
 for cmd in [[sys.executable,'scripts/dev-preflight.py','--static-only'],['go','test','./internal/renderers/tables','./internal/plugins/region','-count=1']]:
  print('REHEARSAL_COMMAND '+json.dumps({'phase':label,'command':cmd}),flush=True)
  r=subprocess.run(cmd,cwd=repo,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
  sys.stdout.buffer.write(r.stdout);sys.stdout.flush();assert r.returncode==0,(label,cmd,r.returncode)
 for f,raw in renderer.items():assert (repo/f).read_bytes()==raw,('independent scope correction changed',f)
try:
 git('restore','--source='+base,'--staged','--worktree','--',*files)
 assert git('write-tree')=='aed4db229c8a7d2b7f6518d5aa3547e86c30878b'
 test('rollback')
 print('Atomic report+runner rollback: exact corrected base; independent scope fix retained; static/tables/region passed',flush=True)
finally:
 git('restore','--source='+head,'--staged','--worktree','--',*files)
for f,b in original.items():assert (repo/f).read_bytes()==b
assert git('write-tree')=='ef2b61f832fb186593d974a3ff3fb9ba2160e245' and not git('status','--porcelain')
test('restored')
print('Exact16file candidate bytes/tree restored; independent scope fix retained; static/tables/region passed',flush=True)
record={'candidate':head,'candidateTree':git('rev-parse',head+'^{tree}'),'correctedBase':base,'baseTree':'aed4db229c8a7d2b7f6518d5aa3547e86c30878b','files':[{'path':f,'sha256':hashlib.sha256(b).hexdigest()} for f,b in original.items()],'independentScopeFixRetained':{f:hashlib.sha256(b).hexdigest() for f,b in renderer.items()},'rollbackAndRestoreChecks':'PASS','scope':'Fresh isolated hosted Linux atomic report+runner rollback/restoration; original local Windows rehearsal separately observed, raw local log currently unavailable; no installed/live/operator rollback claim'}
with proof.open('x',encoding='utf-8',newline='\n') as f:json.dump(record,f,indent=2);f.write('\n')
