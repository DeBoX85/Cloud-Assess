"""Require compiling, assertion-specific reversals of the frozen location contract."""
import hashlib,json,re,subprocess,sys
from pathlib import Path
sys.dont_write_bytecode=True
from mutation_assertions import assertion_witness,run_go_test
ROOT=Path(__file__).resolve().parents[2]
STATUS=subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)
MAIN='internal/plugins/region/subscription_locations.go'
DECODE='internal/plugins/region/subscription_locations_decode.go'
# Each entry is one deliberate fault; opaque folding requires two coordinated edits.
CONTROLS=(
 ('api-version',MAIN,[( 'const subscriptionLocationVersion = "2022-12-01"','const subscriptionLocationVersion = "2024-12-01"')],'TestSubscriptionLocationsFrozenLiterals/fallback-sorted','LOC_REQUEST_BINDING'),
 ('physical-only',DECODE,[( 'if kind == "Physical" {','if true {')],'TestSubscriptionLocationsFrozenLiterals/logical-mapping-does-not-disable-fallback','LOC_LITERAL_STATUS'),
 ('no-peer-fallback',MAIN,[( 'if evidence.CountSource == "fallback" {','if evidence.CountSource == "fallback" || len(row.mappings) == 0 {')],'TestSubscriptionLocationsFrozenLiterals/api-mixed-no-peer-fallback','LOC_LITERAL_VALUES'),
 ('page-limit',MAIN,[( 'if pages >= MaxQuotaPages {','if pages > MaxQuotaPages {')],'TestSubscriptionLocationsFrozenBudgets/pages-plus-one65','LOC_BUDGET_LITERAL'),
 ('work-limit',DECODE,[( 'if page.work >= MaxAuxRows {','if page.work > MaxAuxRows {')],'TestSubscriptionLocationsFrozenBudgets/work-plus-one8193','LOC_BUDGET_LITERAL'),
 ('opaque-case',DECODE,[( '|| keys[key] {','|| keys[regionFoldKey(key)] {'),('keys[key] = true','keys[regionFoldKey(key)] = true')],'TestSubscriptionLocationsFrozenLiterals/opaque-metadata-case-keys','LOC_LITERAL_STATUS'),
 ('typed-alias',DECODE,[( 'providerObject(raw, "name", "metadata", "availabilityZoneMappings", "id", "subscriptionId")','providerObject(raw, "name", "metadata", "id", "subscriptionId")')],'TestSubscriptionLocationsFrozenLiterals/unicode-typed-alias','LOC_LITERAL_STATUS'),
 ('typed-nil-getter',MAIN,[( 'return v.IsNil()','return false')],'TestSubscriptionLocationsConstructorAdmission','LOC_NIL_GETTER'),
 ('getter-context',MAIN,[( 'c.transport.getter.GetBoundedWithResponse(ctx, u.String(), MaxQuotaPageBytes)','c.transport.getter.GetBoundedWithResponse(context.Background(), u.String(), MaxQuotaPageBytes)')],'TestSubscriptionLocationsContextPropagation','LOC_CONTEXT_PROPAGATION'),
 ('canonical-context',MAIN,[( 'return nil, context.Canceled','return nil, e')],'TestSubscriptionLocationsCanonicalContext/cancel','LOC_CONTEXT_SANITIZATION'),
 ('diagnostic-ownership',MAIN,[( 'cloneLocationEvidence(nil, c.state.evidence)','cloneLocationEvidence(nil, SubscriptionLocationEvidence{SubscriptionID:c.SubscriptionID, Status:c.Status, FailureCode:c.FailureCode, CountSource:c.CountSource, Regions:c.Regions, ZoneCounts:c.ZoneCounts, ZoneMappings:c.ZoneMappings})')],'TestSubscriptionLocationsOwnership','LOC_DIAGNOSTIC_ALIAS'),
 ('accessor-ownership',MAIN,[( 'out, _ := cloneLocationEvidence(nil, c.state.evidence)','out := c.state.evidence')],'TestSubscriptionLocationsAccessorLiteralSnapshot','LOC_ACCESSOR_SNAPSHOT'),
 ('body-limit',MAIN,[( 'if len(body) > MaxQuotaPageBytes ||','if len(body) > MaxQuotaPageBytes+1 ||')],'TestSubscriptionLocationsFrozenBudgets/body-plus-one-byte','LOC_BUDGET_LITERAL'),
 ('pinned-fallback',MAIN,[( '"norwayeast", ','')],'TestSubscriptionLocationsAllPinnedFallbacks','LOC_PINNED_FALLBACK'),
)
sources={p:(ROOT/p).read_bytes() for p in [MAIN,DECODE]}
tests={}
for directory,count in [('subscription-locations',9),('subscription-locations-supplement',2)]:
 D=ROOT/'docs/evidence/20261028'/directory
 manifest=json.loads((D/'integrity.json').read_text());assert len(manifest['files'])==count
 for row in manifest['files']:
  raw=(ROOT/row['path']).read_bytes().replace(b'\r\n',b'\n');assert len(raw)==row['bytes'] and hashlib.sha256(raw).hexdigest()==row['sha256']
 for p in D.glob('*_test.go.txt'):
  actual=ROOT/'internal/plugins/region'/p.name.removesuffix('.txt');raw=actual.read_bytes()
  assert raw.replace(b'\r\n',b'\n')==p.read_bytes().replace(b'\r\n',b'\n');tests[str(actual.relative_to(ROOT))]=raw
assert len(tests)==4
pins={p:(ROOT/p).read_bytes() for p in ['go.mod','go.sum']}
race=sys.platform!='win32'
def run(test):
 selector='^TestSubscriptionLocations' if test=='ALL' else '/'.join('^'+re.escape(p)+'$' for p in test.split('/'))
 command=['go','test','-count=1','-run',selector]
 if race:command.append('-race')
 command.append('./internal/plugins/region')
 return run_go_test(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)
baseline=run('ALL');assert baseline.returncode==0,baseline.stdout+baseline.stderr
events=[json.loads(line) for line in baseline.stdout.splitlines() if line.strip()]
assert not any(e['Action'] in ['fail','skip','build-fail'] for e in events)
passed=[e['Test'] for e in events if e['Action']=='pass' and e.get('Test')];assert len(passed)==len(set(passed))==98
try:
 for name,path,edits,test,marker in CONTROLS:
  original=sources[path];changed=original.decode().replace('\r\n','\n')
  for before,after in edits:
   assert changed.count(before)==1,(name,changed.count(before));changed=changed.replace(before,after)
  if b'\r\n' in original:changed=changed.replace('\n','\r\n')
  (ROOT/path).write_bytes(changed.encode())
  command=['go','test','-mod=readonly','-count=1','-run','^$']
  if race:command.append('-race')
  command.append('./internal/plugins/region')
  compiled=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180);assert compiled.returncode==0,name+compiled.stdout+compiled.stderr
  result=run(test);witness=assertion_witness(result,test)
  assert witness is not None and marker in witness['diagnostic'],name+result.stdout+result.stderr
  (ROOT/path).write_bytes(original)
  restored=run(test);assert restored.returncode==0,name+restored.stdout+restored.stderr
  print('SUBSCRIPTION_LOCATION_MUTATION_JSON '+json.dumps({'name':name,'compile_exit':0,'witness':witness,'restored_exit':0},separators=(',',':')),flush=True)
  print('Subscription location mutation rejected: '+name,flush=True)
finally:
 for p,raw in sources.items():(ROOT/p).write_bytes(raw)
assert all((ROOT/p).read_bytes()==raw for p,raw in sources.items())
assert all((ROOT/p).read_bytes()==raw for p,raw in tests.items())
assert all((ROOT/p).read_bytes()==raw for p,raw in pins.items())
assert subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)==STATUS
print('SUBSCRIPTION_LOCATION_PROOF_JSON '+json.dumps({'controls':len(CONTROLS),'baseline_exit':0,'baseline_passed_tests':passed,'race':race,'original_bound_files':9,'supplemental_bound_files':2,'source_sha256':{p:hashlib.sha256(b).hexdigest() for p,b in sources.items()},'test_sha256':{p:hashlib.sha256(b).hexdigest() for p,b in tests.items()},'source_restored_exact':True,'frozen_tests_unchanged':True,'module_pins_unchanged':True,'checkout_hygiene_unchanged':True,'scope':'private subscription locations only; no coordinator/live/deployment eligibility/heap/timer acceptance'},separators=(',',':')),flush=True)
