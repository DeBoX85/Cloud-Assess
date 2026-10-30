"""Require strict compiling reversals for the frozen subscription collector contract."""
import hashlib,json,re,subprocess,sys
from pathlib import Path
sys.dont_write_bytecode=True
from mutation_assertions import assertion_witness,run_go_test
ROOT=Path(__file__).resolve().parents[2]
STATUS=subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)
COLLECTOR='internal/plugins/region/reservation_collector.go'
SUBSCRIPTION='internal/plugins/region/subscription_reservations.go'
CONTROLS=(
 ('outside-pair',COLLECTOR,'if whole == nil && group.Location != request.Region {','if (whole == nil && group.Location != request.Region) || (whole != nil && group.Location == "northeurope") {','TestSubscriptionReservationsFrozenLiterals/all-three-regions','BULK_LITERAL_CALLS'),
 ('whole-status',SUBSCRIPTION,'Status: collected.Evidence.Status','Status: "complete"','TestSubscriptionReservationsFrozenLiterals/get-failure-peer-recovers','BULK_LITERAL_DIAGNOSTICS'),
 ('unadmitted-default',SUBSCRIPTION,'return SubscriptionReservationEvidence{Status: "unknown", Requests:','return SubscriptionReservationEvidence{Status: "complete", Requests:','TestSubscriptionReservationsAdmissionOwnership','BULK_FABRICATED_ADMISSION'),
 ('accessor-ownership',SUBSCRIPTION,'return subscriptionReservationCopy(c.state.evidence)','return c.state.evidence','TestSubscriptionReservationsAdmissionOwnership','BULK_RUN_STATE_LEAK'),
 ('diagnostic-binding',SUBSCRIPTION,'return subscriptionReservationCopy(c.state.evidence)','return subscriptionReservationCopy(SubscriptionReservationEvidence{SubscriptionID: c.SubscriptionID, Status: c.Status, FailureCode: c.FailureCode, Requests: c.state.evidence.Requests, Regional: c.state.evidence.Regional})','TestSubscriptionReservationsAdmissionOwnership','BULK_DIAGNOSTIC_ALIAS'),
 ('group-region',COLLECTOR,'region = group.Location','region = "eastus"','TestSubscriptionReservationsFrozenLiterals/all-three-regions','BULK_LITERAL_DIAGNOSTICS'),
 ('observed-empty',COLLECTOR,'whole.observe(region)','if false { whole.observe(region) }','TestSubscriptionReservationsFrozenLiterals/observed-empty-group','BULK_LITERAL_DIAGNOSTICS'),
 ('wrapped-context',SUBSCRIPTION,'if errors.Is(err, context.Canceled) {\n\t\t\treturn nil, context.Canceled\n\t\t}','if errors.Is(err, context.Canceled) {\n\t\t\treturn nil, err\n\t\t}','TestSubscriptionReservationsCancellation','BULK_CONTEXT_SANITIZATION'),
 ('counter-reset',COLLECTOR,'whole.observe(region)','whole.observe(region)\n\t\t\tcalls = 0','TestSubscriptionReservationsFrozenBudgets/calls-plus-one','BULK_BUDGET_LITERAL'),
 ('unknown-empty-promotion',SUBSCRIPTION,'return subscriptionReservationCopy(c.state.evidence)','out := subscriptionReservationCopy(c.state.evidence)\n\tif len(out.Requests) == 0 { out.Status = "complete" }\n\treturn out','TestSubscriptionReservationsFrozenLiterals/initial-operational-failure','BULK_LITERAL_ENVELOPE'),
)
sources={name:(ROOT/name).read_bytes() for name in [COLLECTOR,SUBSCRIPTION]}
directory=ROOT/'docs/evidence/20261028/subscription-reservations'
manifest=json.loads((directory/'integrity.json').read_text());assert len(manifest['files'])==9
for row in manifest['files']:
 raw=(ROOT/row['path']).read_bytes().replace(b'\r\n',b'\n');assert len(raw)==row['bytes'] and hashlib.sha256(raw).hexdigest()==row['sha256']
tests={}
for p in directory.glob('subscription_reservations_*_test.go.txt'):
 actual=ROOT/'internal/plugins/region'/p.name.removesuffix('.txt');raw=actual.read_bytes();assert raw.replace(b'\r\n',b'\n')==p.read_bytes().replace(b'\r\n',b'\n');tests[str(actual.relative_to(ROOT))]=raw
assert len(tests)==3
pins={p:(ROOT/p).read_bytes() for p in ['go.mod','go.sum']}
race=sys.platform!='win32'
def run(test):
 selector='/'.join('^'+re.escape(part)+'$' for part in test.split('/'))
 if test=='ALL':selector='^TestSubscriptionReservations'
 command=['go','test','-count=1','-run',selector]
 if race:command.append('-race')
 command.append('./internal/plugins/region')
 return run_go_test(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)
baseline=run('ALL');assert baseline.returncode==0,baseline.stdout+baseline.stderr
events=[json.loads(line) for line in baseline.stdout.splitlines() if line.strip()]
assert not any(e['Action'] in ['fail','skip','build-fail'] for e in events)
passed=[e['Test'] for e in events if e['Action']=='pass' and e.get('Test')];assert len(passed)==len(set(passed))==50
try:
 for name,path,before,after,test,marker in CONTROLS:
  original=sources[path];text=original.decode().replace('\r\n','\n');assert text.count(before)==1,(name,text.count(before))
  changed=text.replace(before,after)
  if b'\r\n' in original:changed=changed.replace('\n','\r\n')
  (ROOT/path).write_bytes(changed.encode())
  command=['go','test','-mod=readonly','-run','^$']
  if race:command.append('-race')
  command.append('./internal/plugins/region')
  compiled=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180);assert compiled.returncode==0,name+compiled.stdout+compiled.stderr
  result=run(test);witness=assertion_witness(result,test);assert witness is not None and marker in witness['diagnostic'],name+result.stdout+result.stderr
  (ROOT/path).write_bytes(original)
  restored=run(test);assert restored.returncode==0,name+restored.stdout+restored.stderr
  print('SUBSCRIPTION_RESERVATION_MUTATION_JSON '+json.dumps({'name':name,'compile_exit':0,'witness':witness,'restored_exit':0},separators=(',',':')),flush=True)
  print('Subscription reservation mutation rejected: '+name,flush=True)
finally:
 for path,raw in sources.items():(ROOT/path).write_bytes(raw)
assert all((ROOT/p).read_bytes()==raw for p,raw in sources.items())
assert all((ROOT/p).read_bytes()==raw for p,raw in tests.items())
assert all((ROOT/p).read_bytes()==raw for p,raw in pins.items())
assert subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)==STATUS,'subscription controls changed checkout hygiene'
print('SUBSCRIPTION_RESERVATION_PROOF_JSON '+json.dumps({'controls':len(CONTROLS),'baseline_exit':0,'baseline_passed_tests':passed,'race':race,'frozen_bound_files':9,'literal_cases':13,'budget_recipes':8,'procedural_obligations':12,'source_sha256':{p:hashlib.sha256(raw).hexdigest() for p,raw in sources.items()},'test_sha256':{p:hashlib.sha256(raw).hexdigest() for p,raw in tests.items()},'source_restored_exact':True,'frozen_tests_unchanged':True,'module_pins_unchanged':True,'checkout_hygiene_unchanged':True,'scope':'actual private collector contract only; no coordinator/live/credential/deployment eligibility acceptance'},separators=(',',':')),flush=True)
