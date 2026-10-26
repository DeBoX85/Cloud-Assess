#!/usr/bin/env python3
"""Real compiled assertion/panic controls plus fail-closed classifier cases."""
import ast,copy,hashlib,json,os,subprocess,sys,tempfile
from pathlib import Path
sys.dont_write_bytecode=True
from mutation_assertions import assertion_witness,run_go_test
from qa_temp_cleanup import cleanup_temporary_directory

def verify_driver_source(source):
 tree=ast.parse(source)
 calls=[n for n in ast.walk(tree) if isinstance(n,ast.Call)]
 if sum(isinstance(n.func,ast.Name) and n.func.id=='run_go_test' for n in calls)!=1:raise ValueError('bounded wrapper missing or duplicated')
 if any(isinstance(n.func,ast.Attribute) and ast.unparse(n.func)=='subprocess.run' for n in calls):raise ValueError('unwrapped driver subprocess')
 guards=[]
 for n in ast.walk(tree):
  if isinstance(n,(ast.If,ast.Assert)) and 'intended_failure(' in ast.unparse(n.test):guards.append(n)
 if len(guards)!=1:raise ValueError('strict admission guard missing or duplicated')
 test=guards[0].test
 if isinstance(guards[0],ast.If):
  if not isinstance(test,ast.UnaryOp) or not isinstance(test.op,ast.Not):raise ValueError('weakened rejecting guard')
  test=test.operand
  if not any(isinstance(n,ast.Raise) for n in ast.walk(guards[0])):raise ValueError('guard does not reject')
 if not isinstance(test,ast.Call) or not isinstance(test.func,ast.Name) or test.func.id!='intended_failure' or len(test.args)!=2:raise ValueError('noncanonical admission predicate')
 return tree

def verify_driver_inventory():
 root=Path(__file__).resolve().parents[2];data=json.loads((root/'docs/evidence/20261026/mutation-assurance/driver-inventory.json').read_text());assert len(data['drivers'])==18;names={'MUTATIONS','controls','FILES','FILE','SOURCE','SELECT','manifest'}
 for row in data['drivers']:
  path=root/row['path'];source=path.read_text();tree=verify_driver_source(source);nodes=[n for n in ast.walk(tree) if isinstance(n,ast.Assign) and any(isinstance(a,ast.Name) and a.id in names for a in n.targets)]
  assert hashlib.sha256(ast.dump(ast.Module(body=nodes,type_ignores=[]),include_attributes=False).encode()).hexdigest()==row['control_declarations_sha256'],row['path']
  if 'plan_sha256' in row:assert hashlib.sha256(path.with_suffix('.plan.json').read_bytes()).hexdigest()==row['plan_sha256']
  for bad in [source.replace('run_go_test(', 'subprocess.run(',1),source.replace('intended_failure(', 'missing_admission(',1)]:
   try:verify_driver_source(bad)
   except ValueError:pass
   else:raise AssertionError('driver wiring negative accepted '+row['path'])
 return 18

def main():
 drivers=verify_driver_inventory()
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
  print('MUTATION_ASSERTION_CLASSIFIER_PROOF_JSON '+json.dumps({'drivers_verified':drivers,'driver_wiring_negatives':36,'actual_assertion_accepted':True,'actual_helper_and_subtest_accepted':True,'actual_panic_rejected':True,'panic_after_assertion_rejected':True,'actual_build_error_rejected':True,'actual_go_timeout_rejected':True,'actual_parent_timeout_raised':True,'actual_source_change_rejected':True,'negative_cases':rejected,'restored_healthy_exit':0,'source_unchanged':True,'go_test_timeout':120,'parent_wall_max':180,'scope':'assertion-site provenance, not semantic completeness or authenticated test provenance'},separators=(',',':')),flush=True)
 finally:cleanup_temporary_directory(directory)
if __name__=='__main__':main()
