"""Native preparation observations only; no compiler-failure runtime credit."""
from pathlib import Path
import subprocess,json,hashlib,tempfile,sys
sys.dont_write_bytecode=True
ROOT=Path(__file__).resolve().parents[2];D=ROOT/'docs/evidence/20261029/run-history-binding'
status=subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT);pins={p:(ROOT/p).read_bytes() for p in ['go.mod','go.sum']}
shape_only='--shape-only' in sys.argv
if not shape_only:
 manifest=json.loads((D/'integrity.json').read_text());assert len(manifest['files'])==8
 for row in manifest['files']:
  b=(ROOT/row['path']).read_bytes().replace(b'\r\n',b'\n');assert len(b)==row['bytes'] and hashlib.sha256(b).hexdigest()==row['sha256']
subprocess.run([sys.executable,str(D/'check_literals.py')],cwd=ROOT,check=True)
stub='''package region
import "context"
type RunHistoryEvidence struct{Admitted bool;Request MeterHistoryRequest;Status,FailureCode string;ObservedMeters int;History CostHistoryEvidence}
type RunHistoryBinding struct{Request MeterHistoryRequest;Status,FailureCode string;Meters []HistoricalCostMeter}
func AdmitRunHistory(context.Context,map[string]string,MeterHistoryRequest,*MeterHistoryCollection)(*RunHistoryBinding,error){panic("RH_SHAPE_API_MUST_NOT_RUN")}
func(*RunHistoryBinding)Evidence()RunHistoryEvidence{panic("RH_SHAPE_API_MUST_NOT_RUN")}
'''
with tempfile.TemporaryDirectory(prefix='cloud-assess-history-preproduction-') as tmp:
 T=Path(tmp);replace={str(ROOT/'internal/plugins/region'/p.name.removesuffix('.txt')):str(p) for p in sorted(D.glob('*_test.go.txt'))};assert len(replace)==3
 overlay=T/'overlay.json';overlay.write_text(json.dumps({'Replace':replace}));missing=None
 if not shape_only:
  command=['go','test','-mod=readonly','-gcflags=github.com/DeBoX85/Cloud-Assess/internal/plugins/region=-e','-overlay',str(overlay),'-count=1','-run','^$','./internal/plugins/region']
  missing=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180);body=missing.stdout+missing.stderr
  symbols=['AdmitRunHistory','RunHistoryBinding','RunHistoryEvidence'];assert missing.returncode==1 and all('undefined: '+s in body for s in symbols),body
  assert not any(s in body for s in ['panic:','test timed out','DATA RACE','fatal error:']),body
 shape=T/'shape.go';shape.write_text(stub);replace[str(ROOT/'internal/plugins/region/zz_run_history_shape.go')]=str(shape);overlay.write_text(json.dumps({'Replace':replace}))
 command=['go','test','-mod=readonly','-overlay',str(overlay),'-count=1','-json','-run','^TestRunHistoryFixtureShape$','./internal/plugins/region']
 r=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180);assert r.returncode==0,r.stdout+r.stderr
 ev=[json.loads(l) for l in r.stdout.splitlines() if l.strip()];assert not any(e['Action'] in ['fail','skip','build-fail'] for e in ev);assert [e['Test'] for e in ev if e['Action']=='pass' and e.get('Test')]==['TestRunHistoryFixtureShape']
assert all((ROOT/p).read_bytes()==b for p,b in pins.items());assert subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)==status
proof={'bound_files':None if shape_only else 8,'literal_cases':35,'future_named_events':77,'root_tests':14,'metadata_negatives':8,'fixture_shape_exit':0,'fixture_passed_tests':['TestRunHistoryFixtureShape'],'shape_stub_methods_panic_if_invoked':True,'actual_future_API_executed':False,'compiler_failure_is_runtime_mutation_credit':False,'pins_and_checkout_unchanged':True,'shape_stub_sha256':hashlib.sha256(stub.encode()).hexdigest()}
if missing is not None:proof.update({'missing_api_compile_exit':1,'missing_symbols':symbols,'compiler_stdout_sha256':hashlib.sha256(missing.stdout.encode()).hexdigest(),'compiler_stderr_sha256':hashlib.sha256(missing.stderr.encode()).hexdigest()})
print('RUN_HISTORY_PREPRODUCTION_JSON '+json.dumps(proof,separators=(',',':')),flush=True)
