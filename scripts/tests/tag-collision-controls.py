"""Require compiling alias-scope regressions, restoring every production byte."""
import hashlib,json,subprocess
from pathlib import Path
from mutation_assertions import assertion_witness,run_go_test

ROOT=Path(__file__).resolve().parents[2]
FILE=ROOT/'internal/config/filters.go'
CONTROLS=(
 ('include-map','include = f.Include.Tags','include = nil','TestTagConfigurationAliasesRejected/include/case-conflict'),
 ('exclude-map','exclude = f.Exclude.Tags','exclude = nil','TestTagConfigurationAliasesRejected/exclude/case-conflict'),
 ('duplicate-rejection','if seen[key] {','if false {','TestTagConfigurationAliasesRejected/include/same-value'),
 ('case-fold','key = normalize(key)','key = strings.TrimSpace(key)','TestTagConfigurationAliasesRejected/include/case-conflict'),
 ('trim-fold','key = normalize(key)','key = strings.ToLower(key)','TestTagConfigurationAliasesRejected/include/trim-conflict'),
 ('discovery-fail-closed','func (f *AssessmentFilter) IsResourceExcluded(resourceID string, tags map[string]string) bool {\n\tif f.invalidTagConfig ||','func (f *AssessmentFilter) IsResourceExcluded(resourceID string, tags map[string]string) bool {\n\tif false ||','TestInvalidTagScopeCloneAndRecovery/include'),
 ('downstream-fail-closed','func (f *AssessmentFilter) IsServiceExcluded(resourceID string) bool {\n\tif f.invalidTagConfig ||','func (f *AssessmentFilter) IsServiceExcluded(resourceID string) bool {\n\tif false ||','TestInvalidTagScopeCloneAndRecovery/include'),
 ('clone-invalid-state','out := *f','out := *f\n\tout.invalidTagConfig = false','TestInvalidTagScopeCloneAndRecovery/exclude'),
 ('rebuild-recovery','f.invalidTagConfig = f.validateTagKeys() != nil','f.invalidTagConfig = f.invalidTagConfig || f.validateTagKeys() != nil','TestInvalidTagScopeCloneAndRecovery/include'),
)
def run(test):
 return run_go_test(['go','test','-count=1','-run','^'+test+'$','./internal/config'],cwd=ROOT,capture_output=True,text=True,encoding='utf-8',timeout=180)
original=FILE.read_bytes();tests=(ROOT/'internal/config/tag_collision_test.go').read_bytes()
baseline=run('Test(TagConfigurationAliasesRejected|LoadTagAliasesRejected|InvalidTagScopeCloneAndRecovery|HealthyTagOverlapAndExactValues)');assert baseline.returncode==0,baseline.stdout+baseline.stderr
try:
 for name,before,after,test in CONTROLS:
  source=original.decode().replace('\r\n','\n');assert source.count(before)==1,name
  changed=source.replace(before,after)
  if b'\r\n' in original:changed=changed.replace('\n','\r\n')
  FILE.write_bytes(changed.encode())
  compiled=subprocess.run(['go','test','-mod=readonly','-run','^$','./internal/config'],cwd=ROOT,capture_output=True,text=True,timeout=180);assert compiled.returncode==0,name+compiled.stdout+compiled.stderr
  result=run(test);witness=assertion_witness(result,test);assert witness is not None,name+result.stdout+result.stderr
  FILE.write_bytes(original)
  restored=run(test);assert restored.returncode==0,name+restored.stdout+restored.stderr
  print('TAG_COLLISION_MUTATION_JSON '+json.dumps({'name':name,'compile_exit':0,'witness':witness,'restored_exit':0},separators=(',',':')),flush=True)
  print('Tag collision mutation rejected: '+name,flush=True)
finally:FILE.write_bytes(original)
assert FILE.read_bytes()==original and (ROOT/'internal/config/tag_collision_test.go').read_bytes()==tests
print('TAG_COLLISION_PROOF_JSON '+json.dumps({'controls':len(CONTROLS),'initial_baseline_exit':baseline.returncode,'source_restored_exact':True,'test_source_unchanged':True,'source_sha256':hashlib.sha256(original).hexdigest(),'test_sha256':hashlib.sha256(tests).hexdigest(),'scope':'configuration aliases only; source assertion evidence, not authenticated provenance or live Azure'},separators=(',',':')),flush=True)
