"""Mandatory native missing-API/fixture-shape observations, not mutation evidence."""
import hashlib,json,subprocess,sys,tempfile
from pathlib import Path
sys.dont_write_bytecode=True
ROOT=Path(__file__).resolve().parents[2]
D=ROOT/'docs/evidence/20261028/subscription-locations'
status=subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)
pins={p:(ROOT/p).read_bytes() for p in ['go.mod','go.sum']}
manifest=json.loads((D/'integrity.json').read_text());assert len(manifest['files'])==9
for row in manifest['files']:
 raw=(ROOT/row['path']).read_bytes().replace(b'\r\n',b'\n');assert len(raw)==row['bytes'] and hashlib.sha256(raw).hexdigest()==row['sha256']
subprocess.run([sys.executable,str(D/'check_literals.py')],cwd=ROOT,check=True)
stub='''package region
import "context"
type SubscriptionLocationRequest struct{SubscriptionID string}
type SubscriptionLocationEvidence struct{SubscriptionID,Status,FailureCode,CountSource string;Regions []string;ZoneCounts map[string]int;ZoneMappings map[string]map[string]string}
type SubscriptionLocationCollection struct{SubscriptionID,Status,FailureCode,CountSource string;Regions []string;ZoneCounts map[string]int;ZoneMappings map[string]map[string]string}
type SubscriptionLocationCollector struct{}
func NewSubscriptionLocationCollector(string,RESTQuotaGetter)(*SubscriptionLocationCollector,error){panic("LOCATION_SHAPE_ONLY_API_MUST_NOT_EXECUTE")}
func(*SubscriptionLocationCollector)Collect(context.Context,map[string]string,SubscriptionLocationRequest)(*SubscriptionLocationCollection,error){panic("LOCATION_SHAPE_ONLY_API_MUST_NOT_EXECUTE")}
func(*SubscriptionLocationCollection)Evidence()SubscriptionLocationEvidence{panic("LOCATION_SHAPE_ONLY_API_MUST_NOT_EXECUTE")}
'''
with tempfile.TemporaryDirectory(prefix='cloud-assess-location-preproduction-') as tmp:
 T=Path(tmp);replace={str(ROOT/'internal/plugins/region'/p.name.removesuffix('.txt')):str(p) for p in sorted(D.glob('*_test.go.txt'))};assert len(replace)==3
 overlay=T/'overlay.json';overlay.write_text(json.dumps({'Replace':replace}))
 command=['go','test','-mod=readonly','-overlay',str(overlay),'-count=1','-run','^$','./internal/plugins/region']
 missing=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)
 body=missing.stdout+missing.stderr;symbols=['NewSubscriptionLocationCollector','SubscriptionLocationCollection','SubscriptionLocationRequest']
 assert missing.returncode==1 and all('undefined: '+s in body for s in symbols),body
 assert not any(s in body for s in ['panic:','test timed out','DATA RACE','fatal error:']),body
 shape=T/'shape.go';shape.write_text(stub);replace[str(ROOT/'internal/plugins/region/zz_location_shape.go')]=str(shape);overlay.write_text(json.dumps({'Replace':replace}))
 command=['go','test','-mod=readonly','-overlay',str(overlay),'-count=1','-json','-run','^TestSubscriptionLocationsFixtureShape$','./internal/plugins/region']
 shape_result=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)
 assert shape_result.returncode==0,shape_result.stdout+shape_result.stderr
 events=[json.loads(line) for line in shape_result.stdout.splitlines() if line.strip()];assert not any(e['Action'] in ['fail','skip','build-fail'] for e in events)
 assert [e['Test'] for e in events if e['Action']=='pass' and e.get('Test')]==['TestSubscriptionLocationsFixtureShape']
assert all((ROOT/p).read_bytes()==raw for p,raw in pins.items())
assert subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)==status
print('SUBSCRIPTION_LOCATION_PREPRODUCTION_JSON '+json.dumps({'bound_files':9,'wire_literals':35,'boundary_recipes':8,'procedural_obligations':8,'missing_api_compile_exit':1,'missing_symbols':symbols,'missing_compile_stdout_sha256':hashlib.sha256(missing.stdout.encode()).hexdigest(),'missing_compile_stderr_sha256':hashlib.sha256(missing.stderr.encode()).hexdigest(),'fixture_shape_exit':0,'fixture_passed_tests':['TestSubscriptionLocationsFixtureShape'],'shape_stub_methods_panic_if_executed':True,'actual_future_API_executed':False,'compiler_failure_is_runtime_mutation_credit':False,'module_pins_unchanged':True,'checkout_unchanged':True},separators=(',',':')),flush=True)
