"""Require compiling, assertion-specific reversals of frozen history consumers."""
import hashlib,json,re,subprocess,sys,os
from pathlib import Path
sys.dont_write_bytecode=True
from mutation_assertions import assertion_witness,run_go_test
ROOT=Path(__file__).resolve().parents[2]
STATUS=subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)
MAIN='internal/plugins/region/run_history_binding.go'
AUX='internal/plugins/region/auxiliary.go'
# Each entry is one semantic fault; shared auxiliary text helper is restored exactly.
CONTROLS=(
 ('request-admission',MAIN,[('if !historyRequestValid(q) {','if !historyRequestValid(q) && q.Months != 13 {')],'TestRunHistoryRequestAdmission/months-thirteen','RH_LITERAL_REJECTION'),
 ('scope-alias',MAIN,[('scope, ok := auxiliaryScope(ctx, selected, &text)','folded := map[string]string{}; for id, name := range selected { folded[strings.ToLower(id)] = name }; scope, ok := auxiliaryScope(ctx, folded, &text)')],'TestRunHistoryRequestAdmission/scope-alias','RH_LITERAL_REJECTION'),
 ('all-scope-names',MAIN,[('if name == "" {','if name == "" && false {')],'TestRunHistoryOtherScopeAdmission','RH_SCOPE_ALL_NAMES'),
 ('returned-currency-exact',MAIN,[('in.Currency != q.Currency','!strings.EqualFold(in.Currency, q.Currency)')],'TestRunHistoryReturnedCurrencySpelling','RH_CURRENCY_EXACT'),
 ('currency-binding',MAIN,[('in.Currency != q.Currency','(in.Currency != q.Currency && false)')],'TestRunHistoryCurrencyBinding','RH_LITERAL_REJECTION'),
 ('complete-code',MAIN,[('if in.FailureCode != "" {','if in.FailureCode != "" && false {')],'TestRunHistoryFrozenLiterals/complete-code','RH_LITERAL_REJECTION'),
 ('unknown-prefix',MAIN,[('(in.Status == "unknown" && len(in.Meters) != 0)','(in.Status == "unknown" && len(in.Meters) != 0 && false)')],'TestRunHistoryFrozenLiterals/unknown-prefix','RH_LITERAL_REJECTION'),
 ('partial-eligibility',MAIN,[('Complete: in.Status == "complete",','Complete: in.Status == "complete" || in.Status == "partial",')],'TestRunHistoryFrozenLiterals/partial-prefix','RH_LITERAL_VALUES'),
 ('exact-duplicates',MAIN,[('|| seen[meter.MeterID] {','|| (seen[meter.MeterID] && false) {')],'TestRunHistoryFrozenLiterals/duplicate-equal','RH_LITERAL_REJECTION'),
 ('opaque-case',MAIN,[('|| seen[meter.MeterID] {','|| seen[strings.ToLower(meter.MeterID)] {'),('seen[meter.MeterID] = true','seen[strings.ToLower(meter.MeterID)] = true')],'TestRunHistoryFrozenLiterals/opaque-case','RH_LITERAL_RESULT'),
 ('meter-limit',MAIN,[('if len(in.Meters) > MaxAuxRows {','if len(in.Meters) > MaxAuxRows+1 {')],'TestRunHistoryFrozenBudgets/meters-plus-one8193','RH_LITERAL_REJECTION'),
 ('utf8-byte-limit',AUX,[('if len(s) > MaxLabelBytes ||','if utf8.RuneCountInString(s) > MaxLabelBytes ||')],'TestRunHistoryUTF8ByteBoundary/above514bytes','RH_UTF8_BYTES'),
 ('nonnegative-weight',MAIN,[('weight < 0 ||','(weight < 0 && false) ||')],'TestRunHistoryFrozenLiterals/negative','RH_LITERAL_REJECTION'),
 ('finite-weight',MAIN,[('math.IsNaN(weight) ||','(math.IsNaN(weight) && false) ||')],'TestRunHistoryFrozenLiterals/nan','RH_LITERAL_REJECTION'),
 ('minimum-positive-weight',MAIN,[('(weight > 0 && weight < MinCostWeight)','(weight > 0 && weight < MinCostWeight && false)')],'TestRunHistoryFrozenLiterals/tiny-positive','RH_LITERAL_REJECTION'),
 ('input-ownership',MAIN,[('history.Meters = owned','history.Meters = in.Meters')],'TestRunHistoryInputOwnership','RH_INPUT_OWNERSHIP'),
 ('diagnostic-ownership',MAIN,[('out := r.state.evidence','out := RunHistoryEvidence{Admitted:true, Request:r.Request, Status:r.Status, FailureCode:r.FailureCode, ObservedMeters:len(r.Meters), History:CostHistoryEvidence{Complete:r.Status=="complete", Meters:r.Meters}}')],'TestRunHistoryOwnership','RH_DIAGNOSTIC_OWNERSHIP'),
 ('accessor-ownership',MAIN,[('out.History.Meters = append([]HistoricalCostMeter{}, r.state.evidence.History.Meters...)','out.History.Meters = r.state.evidence.History.Meters')],'TestRunHistoryOwnership','RH_ACCESSOR_OWNERSHIP'),
 ('forgery-admission',MAIN,[('return RunHistoryEvidence{Status: "unknown",','return RunHistoryEvidence{Admitted:true, Status: "unknown",')],'TestRunHistoryForgery','RH_FORGED_ADMISSION'),
 ('context-observation',MAIN,[('return ctx.Err()','_ = ctx.Err(); return nil')],'TestRunHistoryCancellation/cancel','RH_CONTEXT_BOUNDARY'),
 ('UTC-request',MAIN,[('q.Now = q.Now.UTC()','q.Now = q.Now.Add(0)')],'TestRunHistoryNormalizedRequest','RH_REQUEST_BINDING'),
 ('sanitized-error',MAIN,[('errors.New("region run history binding rejected [run_history_invalid]")','errors.New("SENSITIVE invalid input")')],'TestRunHistoryCurrencyBinding','RH_LITERAL_REJECTION'),
)
sources={p:(ROOT/p).read_bytes() for p in [MAIN,AUX]}
tests={}
for directory,count in [('run-history-binding',8),('run-history-supplement',2)]:
 D=ROOT/'docs/evidence/20261029'/directory
 manifest=json.loads((D/'integrity.json').read_text());assert len(manifest['files'])==count
 for row in manifest['files']:
  raw=(ROOT/row['path']).read_bytes().replace(b'\r\n',b'\n');assert len(raw)==row['bytes'] and hashlib.sha256(raw).hexdigest()==row['sha256']
 for p in D.glob('*_test.go.txt'):
  actual=ROOT/'internal/plugins/region'/p.name.removesuffix('.txt');raw=actual.read_bytes()
  assert raw.replace(b'\r\n',b'\n')==p.read_bytes().replace(b'\r\n',b'\n');tests[str(actual.relative_to(ROOT))]=raw
assert len(tests)==4
pins={p:(ROOT/p).read_bytes() for p in ['go.mod','go.sum']}
race=sys.platform!='win32'
PRIVATE=None
if '--private-log-dir' in sys.argv:
 PRIVATE=Path(sys.argv[sys.argv.index('--private-log-dir')+1]).resolve();assert not PRIVATE.is_relative_to(ROOT);PRIVATE.mkdir(parents=True,exist_ok=True)
def retain(label,result):
 if PRIVATE is not None:
  for stream in ['stdout','stderr']:(PRIVATE/(label+'.'+stream)).write_text(getattr(result,stream) or '',encoding='utf-8',newline='')
def run(test):
 selector='^TestRunHistory' if test=='ALL' else '/'.join('^'+re.escape(p)+'$' for p in test.split('/'))
 command=['go','test','-count=1','-run',selector]
 if race:command.append('-race')
 command.append('./internal/plugins/region')
 return run_go_test(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)
baseline=run('ALL');retain('baseline',baseline);assert baseline.returncode==0,baseline.stdout+baseline.stderr
events=[json.loads(line) for line in baseline.stdout.splitlines() if line.strip()]
assert not any(e['Action'] in ['fail','skip','build-fail'] for e in events)
passed=[e['Test'] for e in events if e['Action']=='pass' and e.get('Test')];assert len(passed)==len(set(passed))==82
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
  compiled=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180);retain(name+'-compile',compiled);assert compiled.returncode==0,name+compiled.stdout+compiled.stderr
  result=run(test);retain(name+'-fault',result);witness=assertion_witness(result,test)
  assert witness is not None and marker in witness['diagnostic'],name+result.stdout+result.stderr
  (ROOT/path).write_bytes(original)
  restored=run(test);retain(name+'-restored',restored);assert restored.returncode==0,name+restored.stdout+restored.stderr
  print('RUN_HISTORY_MUTATION_JSON '+json.dumps({'name':name,'compile_exit':0,'witness':witness,'restored_exit':0},separators=(',',':')),flush=True)
  print('Run history mutation rejected: '+name,flush=True)
finally:
 for p,raw in sources.items():(ROOT/p).write_bytes(raw)
assert all((ROOT/p).read_bytes()==raw for p,raw in sources.items())
assert all((ROOT/p).read_bytes()==raw for p,raw in tests.items())
assert all((ROOT/p).read_bytes()==raw for p,raw in pins.items())
assert subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)==STATUS
print('RUN_HISTORY_PROOF_JSON '+json.dumps({'controls':len(CONTROLS),'baseline_exit':0,'baseline_passed_tests':passed,'race':race,'original_bound_files':8,'supplemental_bound_files':2,'source_sha256':{p:hashlib.sha256(b).hexdigest() for p,b in sources.items()},'test_sha256':{p:hashlib.sha256(b).hexdigest() for p,b in tests.items()},'source_restored_exact':True,'frozen_tests_unchanged':True,'module_pins_unchanged':True,'checkout_hygiene_unchanged':True,'local_complete_child_streams_retained':PRIVATE is not None,'scope':'pure structural history admission only; no coordinator/live/service or query-window authentication/heap/timer acceptance'},separators=(',',':')),flush=True)
