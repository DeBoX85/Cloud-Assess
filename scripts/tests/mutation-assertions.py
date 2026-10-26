#!/usr/bin/env python3
"""Real compiled assertion/panic controls plus fail-closed classifier cases."""
import copy,hashlib,json,os,subprocess,sys,tempfile
from pathlib import Path
sys.dont_write_bytecode=True
from mutation_assertions import assertion_witness,run_go_test
from qa_temp_cleanup import cleanup_temporary_directory

def main():
 directory=tempfile.TemporaryDirectory(prefix='cloud-assess-assertion-admission-')
 try:
  root=Path(directory.name);(root/'go.mod').write_text('module assertionfixture\n\ngo 1.26.10\n')
  source='''package assertionfixture
import ("testing";"time";"os")
func contract(t *testing.T) { t.Helper(); t.Fatal("EXPECTED_CONTRACT_VIOLATION") }
func TestExpected(t *testing.T) { contract(t) }
func TestPanic(t *testing.T) { panic("EXPECTED_CONTRACT_VIOLATION") }
func TestAfterAssertion(t *testing.T) { defer func(){panic("EXPECTED_CONTRACT_VIOLATION")}();t.Error("EXPECTED_CONTRACT_VIOLATION") }
func TestHealthy(t *testing.T) { }
func TestHang(t *testing.T) { time.Sleep(time.Second) }
func TestSourceChange(t *testing.T) { b,e:=os.ReadFile("fixture_test.go");if e!=nil{t.Fatal(e)};if e:=os.WriteFile("fixture_test.go",append(b,10),0600);e!=nil{t.Fatal(e)} }
func TestSubtests(t *testing.T) { t.Run("selected",func(t *testing.T){ t.Fatal("SUBTEST_CONTRACT") }) }
''';path=root/'fixture_test.go';path.write_text(source);original=path.read_bytes();go=os.environ.get('QA_GO','go');results={}
  for name in ['TestExpected','TestPanic','TestAfterAssertion','TestHealthy','TestSubtests/selected']:
   result=run_go_test([go,'test','-count=1','-run','^'+name+'$','.'],cwd=root,capture_output=True,text=True,encoding='utf-8',errors='strict',timeout=30);results[name]=result
   print('ASSERTION_ADMISSION_ACTUAL_JSON '+json.dumps({'test':name,'exit':result.returncode,'stdout':result.stdout,'stderr':result.stderr},separators=(',',':')),flush=True)
  good=results['TestExpected'];assert assertion_witness(good,'TestExpected')
  assert assertion_witness(results['TestSubtests/selected'],'TestSubtests/selected')
  assert not assertion_witness(results['TestPanic'],'TestPanic')
  assert not assertion_witness(results['TestAfterAssertion'],'TestAfterAssertion')
  assert not assertion_witness(results['TestHealthy'],'TestHealthy')
  rejected=[]
  for label,change in [
   ('zero-exit',lambda r:setattr(r,'returncode',0)),('wrong-exit',lambda r:setattr(r,'returncode',2)),
   ('stderr',lambda r:setattr(r,'stderr','infrastructure failed')),('unframed',lambda r:setattr(r,'stdout','--- FAIL: TestExpected\n')),
   ('wrong-test',lambda r:setattr(r,'stdout',r.stdout.replace('TestExpected','TestOther'))),
   ('missing-source-provenance',lambda r:setattr(r,'assertion_sites',[])),
   ('wrong-source-line',lambda r:setattr(r,'assertion_sites',[dict(s,line=s['line']+100) for s in r.assertion_sites])),
   ('wrong-source-root',lambda r:setattr(r,'assertion_sites',[dict(s,roots=['TestOther']) for s in r.assertion_sites])),
   ('panic-payload',lambda r:setattr(r,'stdout',r.stdout+json.dumps({'Action':'output','Output':'panic: crash'})+'\n')),
   ('build-error',lambda r:setattr(r,'stdout',r.stdout+json.dumps({'Action':'build-fail'})+'\n')),
   ('timeout',lambda r:setattr(r,'stdout',r.stdout+json.dumps({'Action':'output','Output':'test timed out'})+'\n')),
   ('race',lambda r:setattr(r,'stdout',r.stdout+json.dumps({'Action':'output','Output':'DATA RACE'})+'\n')),
   ('no-package-failure',lambda r:setattr(r,'stdout','\n'.join(l for l in r.stdout.splitlines() if not(json.loads(l).get('Action')=='fail' and not json.loads(l).get('Test')))+'\n')),
   ('duplicate-json-key',lambda r:setattr(r,'stdout',r.stdout+'{"Action":"output","Action":"fail"}\n')),
   ('no-assertion-reason',lambda r:setattr(r,'stdout',r.stdout.replace('EXPECTED_CONTRACT_VIOLATION','').replace('fixture_test.go:', 'not_source:'))),
  ]:
   changed=copy.copy(good);change(changed);assert not assertion_witness(changed,'TestExpected'),label;rejected.append(label)
  assert path.read_bytes()==original
  timed=run_go_test([go,'test','-timeout=50ms','-run','^TestHang$','.'],cwd=root,capture_output=True,text=True,timeout=30);assert not assertion_witness(timed,'TestHang')
  try:run_go_test([go,'test','-run','^TestHang$','.'],cwd=root,capture_output=True,text=True,timeout=0.1)
  except subprocess.TimeoutExpired:pass
  else:raise AssertionError('parent timeout not enforced')
  try:run_go_test([go,'test','-run','^TestSourceChange$','.'],cwd=root,capture_output=True,text=True,timeout=30)
  except RuntimeError as e:assert 'original test source changed' in str(e)
  else:raise AssertionError('changed source accepted')
  finally:path.write_bytes(original)
  # A real compilation error cannot become assertion evidence.
  path.write_text(source+'\nvar broken = DoesNotExist\n');bad=run_go_test([go,'test','-run','^TestExpected$','.'],cwd=root,capture_output=True,text=True,timeout=30);assert not assertion_witness(bad,'TestExpected');path.write_bytes(original)
  restored=run_go_test([go,'test','-run','^TestHealthy$','.'],cwd=root,capture_output=True,text=True,timeout=30);assert restored.returncode==0 and path.read_bytes()==original
  print('MUTATION_ASSERTION_CLASSIFIER_PROOF_JSON '+json.dumps({'actual_assertion_accepted':True,'actual_helper_and_subtest_accepted':True,'actual_panic_rejected':True,'panic_after_assertion_rejected':True,'actual_build_error_rejected':True,'actual_go_timeout_rejected':True,'actual_parent_timeout_raised':True,'actual_source_change_rejected':True,'negative_cases':rejected,'restored_healthy_exit':0,'source_unchanged':True,'go_test_timeout':120,'parent_wall_max':180,'scope':'assertion-site provenance, not semantic completeness or authenticated test provenance'},separators=(',',':')),flush=True)
 finally:cleanup_temporary_directory(directory)
if __name__=='__main__':main()
