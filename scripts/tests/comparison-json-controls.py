"""Require compiling F03 assertion failures and exact restoration on each host."""
import hashlib,json,re,subprocess,sys
from pathlib import Path

ROOT=Path(__file__).resolve().parents[2]
CHECKOUT_STATUS=subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)
sys.dont_write_bytecode=True
from mutation_assertions import assertion_witness,run_go_test

DIRECTORY=ROOT/'internal/equivalence'
CONTROLS=(
 ('reference-entry','projection.go','validateComparisonJSON(content, nil)','validateComparisonJSON([]byte("{}"), nil)','TestJSONAdmissionRejectsAmbiguousReference/section'),
 ('target-entry','projection.go','validateComparisonJSON(content, reflect.TypeFor[result.AssessmentResult]())','validateComparisonJSON([]byte("{}"), reflect.TypeFor[result.AssessmentResult]())','TestJSONAdmissionRejectsAmbiguousTarget/completeness'),
 ('decoded-names','json_admission.go','if frame.names[name] {','if false {','TestJSONAdmissionRejectsAmbiguousReference/escaped-section'),
 ('nested-names','json_admission.go','if frame.names[name] {','if frame.names[name] && len(stack) == 1 {','TestJSONAdmissionRejectsAmbiguousReference/row'),
 ('typed-field-alias','json_admission.go','if index >= 0 && frame.fields[index] {','if false {','TestJSONAdmissionRejectsAmbiguousTarget/case-alias'),
 ('typed-array','json_admission.go','valueType = comparisonJSONArrayElement(frame.valueType)','valueType = nil','TestJSONAdmissionRejectsAmbiguousTarget/array-field'),
 ('pointer-field','json_admission.go','for t != nil && t.Kind() == reflect.Pointer {','for false {','TestJSONAdmissionRejectsAmbiguousTarget/pointer-field'),
 ('opaque-fold','json_admission.go','if frame.names[name] {','name = strings.ToLower(name)\n\t\t\t\tif frame.names[name] {','TestJSONAdmissionPreservesOpaqueMapsAndSingleAliases/opaque-tags'),
 ('numeric-scan','json_admission.go','decoder.UseNumber()','// decoder.UseNumber()','TestJSONAdmissionPreservesOpaqueMapsAndSingleAliases/unknown-number'),
)

def run(test,package='./internal/equivalence',pattern=False):
 selector='^'+test+'$' if pattern else '/'.join('^'+re.escape(p)+'$' for p in test.split('/'))
 return run_go_test(['go','test','-count=1','-timeout=120s','-run',selector,package],cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)

originals={n:(DIRECTORY/n).read_bytes() for n in ['projection.go','json_admission.go']}
test_paths=['internal/equivalence/json_admission_test.go','tools/equivalence/json_admission_test.go']
tests={n:(ROOT/n).read_bytes() for n in test_paths}
baseline=run('TestJSONAdmission.*',pattern=True);assert baseline.returncode==0,baseline.stdout+baseline.stderr
cli_baseline=run('TestComparisonJSONAdmissionFailsBeforeOutput','./tools/equivalence');assert cli_baseline.returncode==0,cli_baseline.stdout+cli_baseline.stderr
baseline_tests=[]
for outcome,expected_roots in [(baseline,{'TestJSONAdmissionRejectsAmbiguousReference','TestJSONAdmissionRejectsAmbiguousTarget','TestJSONAdmissionPreservesOpaqueMapsAndSingleAliases','TestJSONAdmissionPreservesExistingSyntaxPolicy'}),(cli_baseline,{'TestComparisonJSONAdmissionFailsBeforeOutput'})]:
 events=[json.loads(line) for line in outcome.stdout.splitlines() if line.strip()]
 assert not any(e['Action'] in ['skip','fail','build-fail'] for e in events)
 passed=[e['Test'] for e in events if e['Action']=='pass' and e.get('Test')]
 assert {t for t in passed if '/' not in t}==expected_roots,'missing actual healthy baseline tests'
 baseline_tests.extend(passed)
try:
 for name,filename,before,after,test in CONTROLS:
  raw=originals[filename];source=raw.decode().replace('\r\n','\n');assert source.count(before)==1,name
  changed=source.replace(before,after)
  if b'\r\n' in raw:changed=changed.replace('\n','\r\n')
  (DIRECTORY/filename).write_bytes(changed.encode())
  compiled=subprocess.run(['go','test','-mod=readonly','-run','^$','./internal/equivalence'],cwd=ROOT,capture_output=True,text=True,timeout=180);assert compiled.returncode==0,name+compiled.stdout+compiled.stderr
  result=run(test);witness=assertion_witness(result,test);assert witness is not None,name+result.stdout+result.stderr
  for n,b in originals.items():(DIRECTORY/n).write_bytes(b)
  restored=run(test);assert restored.returncode==0,name+restored.stdout+restored.stderr
  print('COMPARISON_JSON_MUTATION_JSON '+json.dumps({'name':name,'compile_exit':0,'witness':witness,'restored_exit':0},separators=(',',':')),flush=True)
  print('Comparison JSON mutation rejected: '+name,flush=True)
finally:
 for n,b in originals.items():(DIRECTORY/n).write_bytes(b)

assert all((DIRECTORY/n).read_bytes()==b for n,b in originals.items())
assert all((ROOT/n).read_bytes()==b for n,b in tests.items())
assert subprocess.check_output(['git','status','--porcelain','--untracked-files=all'],cwd=ROOT)==CHECKOUT_STATUS,'comparison JSON driver changed checkout hygiene'
print('COMPARISON_JSON_PROOF_JSON '+json.dumps({'controls':len(CONTROLS),'initial_baseline_exit':0,'cli_baseline_exit':0,'baseline_passed_tests':baseline_tests,'source_restored_exact':True,'test_source_unchanged':True,'checkout_hygiene_unchanged':True,'source_sha256':{n:hashlib.sha256(b).hexdigest() for n,b in originals.items()},'test_sha256':{n:hashlib.sha256(b).hexdigest() for n,b in tests.items()},'scope':'comparison JSON admission only; bounded source assertion evidence, not live Azure or authenticated provenance'},separators=(',',':')),flush=True)
