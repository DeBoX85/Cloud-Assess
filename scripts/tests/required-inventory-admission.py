#!/usr/bin/env python3
"""Strict bounded Admission proofs; full native acceptance remains separate."""
import hashlib,importlib.util,json,os,shutil,subprocess,tempfile,sys
from pathlib import Path
sys.dont_write_bytecode=True
spec=importlib.util.spec_from_file_location('required_build',Path(__file__).with_name('required-inventory-build.py'))
b=importlib.util.module_from_spec(spec);spec.loader.exec_module(b)
ROOT=b.ROOT
b.TESTS += ['required_inventory_contract_test.go','required_evidence_test.go','required_evidence_budget_test.go']
CONTROLS=[
 ('input-text-excess-one','required_inventory.go','len(s) > auxTextBudget-b.text','len(s) > auxTextBudget+1-b.text','TestRequiredAdmissionFrozenInputBudgets/text-first-excess','ADMISSION_INPUT_BUDGET_REJECTION'),
 ('input-ledger-reset','required_evidence.go','if !input.collection(c, q) {','input.budget = requiredBudget{}; if !input.collection(c, q) {','TestRequiredAdmissionFrozenInputBudgets/entries-first-excess','ADMISSION_INPUT_BUDGET_REJECTION'),
 ('output-text-excess-one','required_inventory.go','len(s) > auxTextBudget-b.text','len(s) > auxTextBudget+1-b.text','TestRequiredAdmissionFrozenOutputBudgets/text-first-excess','ADMISSION_OUTPUT_BUDGET_REJECTION'),
 ('output-third-copy-omitted','required_evidence.go','i < 3','i < 2','TestRequiredAdmissionFrozenOutputBudgets/entries-first-excess','ADMISSION_OUTPUT_BUDGET_REJECTION'),
 ('output-exact-entry-rejected','required_inventory.go','n > MaxInventoryEntries-b.entries','n >= MaxInventoryEntries-b.entries','TestRequiredAdmissionFrozenOutputBudgets/entries-exact','ADMISSION_OUTPUT_BUDGET_EXACT'),
 ('output-exact-text-rejected','required_inventory.go','len(s) > auxTextBudget-b.text','len(s) >= auxTextBudget-b.text','TestRequiredAdmissionFrozenOutputBudgets/text-exact','ADMISSION_OUTPUT_BUDGET_EXACT'),
 ('private-score-alias','required_evidence.go','result.values = requiredCopySKUs(result.SKUs)','result.values = result.SKUs','TestRequiredAdmissionForgeryAndOwnership','ADMISSION_PRIVATE_OWNERSHIP'),
 ('public-plan-trusted','required_evidence.go','range plan.state.requirements','range plan.Requirements','TestRequiredAdmissionForgeryAndOwnership','ADMISSION_PLAN_DIAGNOSTIC_OWNERSHIP'),
 ('terminal-escape','required_inventory.go','return ctx.Err()\n}','return nil\n}','TestRequiredAdmissionContext/before-return-cancel-healthy','ADMISSION_LATE_TERMINAL'),
 ('managed-instance-coordinate-collapse','required_evidence_input.go','sqlCapabilityCoordinate(s, kind)','sqlCoordinate(s)','TestRequiredAdmissionManagedInstanceCoordinates','ADMISSION_MI_ACTUAL_COORDINATE'),
 ('disk-terminal-sentinel-admitted','required_evidence_input.go','case "sku_byte_limit",','case "sku_cancelled", "sku_byte_limit",','TestRequiredAdmissionAllPayloadPrevalidation/disk-failure','ADMISSION_PREVALIDATION'),
 ('unused-zone-duplicate-ignored','required_evidence_input.go','z <= prior','z < prior','TestRequiredAdmissionAllPayloadPrevalidation/vm-unused','ADMISSION_PREVALIDATION'),
 ('outer-limit-after-forgery','required_evidence.go','if len(collections) > MaxInventoryEntries {','if len(collections) > MaxInventoryEntries && plan != nil {','TestRequiredAdmissionOuterAndScalar','ADMISSION_OUTER_FIRST'),
]

def main():
 classifier=b.classifier_controls();bootstrap=b.bootstrap_controls()
 directory=tempfile.TemporaryDirectory(prefix='cloud-assess-required-admission-')
 try:
  work=Path(directory.name)
  for name in subprocess.check_output(['git','ls-files','-z'],cwd=ROOT).decode().split('\0'):
   if name and (ROOT/name).is_file():dest=work/name;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(ROOT/name,dest)
  print('REQUIRED_ADMISSION_MODULES_JSON '+json.dumps(b.bootstrap_modules(work),separators=(',',':')),flush=True)
  baseline=b.run(work,race=os.name!='nt')
  if baseline.returncode or baseline.stderr.strip():raise RuntimeError('healthy Admission baseline failed\n'+baseline.stdout+baseline.stderr)
  passed=[e.get('Test','') for e in b.events(baseline.stdout) if e['Action']=='pass']
  for prefix,count in [('TestRequiredInventoryFrozenLiterals/',76),('TestRequiredAdmissionContext/',8),('TestRequiredAdmissionFrozenInputBudgets/',4),('TestRequiredAdmissionFrozenOutputBudgets/',4),('TestRequiredAdmissionAllPayloadPrevalidation/',15)]:
   if len([x for x in passed if x.startswith(prefix)])!=count:raise RuntimeError('lost proof inventory '+prefix)
  for line in baseline.stdout.splitlines():print('REQUIRED_ADMISSION_BASELINE_JSON '+line,flush=True)
  sources={p:(work/'internal/plugins/region'/p).read_bytes() for p in ['required_inventory.go','required_evidence.go','required_evidence_input.go']}
  results=[]
  for name,path,old,new,test,marker in CONTROLS:
   target=work/'internal/plugins/region'/path;text=sources[path].decode().replace('\r\n','\n')
   if text.count(old)!=1:raise RuntimeError('mutation anchor changed '+name)
   target.write_bytes(text.replace(old,new,1).encode())
   compiled=b.run(work,'^$')
   if compiled.returncode or compiled.stderr.strip():raise RuntimeError('mutation did not compile '+name+'\n'+compiled.stdout+compiled.stderr)
   failed=b.run(work,test)
   if not b.intended_failure(failed.returncode,failed.stdout,failed.stderr,test,marker):raise RuntimeError('not intended assertion '+name+'\n'+failed.stdout+failed.stderr)
   result={'name':name,'test':test,'marker':marker,'compile_exit':0,'exit':1,'events':b.events(failed.stdout),'stdout_sha256':hashlib.sha256(failed.stdout.encode()).hexdigest()};results.append(result)
   print('REQUIRED_ADMISSION_MUTATION_JSON '+json.dumps(result,separators=(',',':')),flush=True);print('required inventory Admission mutation rejected: '+name,flush=True)
   target.write_bytes(sources[path])
   if target.read_bytes()!=sources[path]:raise RuntimeError('restore mismatch '+name)
  restored=b.run(work)
  if restored.returncode or restored.stderr.strip():raise RuntimeError('restored baseline failed\n'+restored.stdout+restored.stderr)
  hashes={p:hashlib.sha256(raw).hexdigest() for p,raw in sources.items()}
  for p,raw in sources.items():
   if (ROOT/'internal/plugins/region'/p).read_bytes()!=raw:raise RuntimeError('original source changed')
  print('REQUIRED_ADMISSION_PROOF_JSON '+json.dumps({'frozen_scenarios':76,'terminal_schedules':8,'input_recipes':4,'output_recipes':4,'typed_single_faults':15,'strict_compiling_reversals':len(results),'classifier_negatives':classifier,'bootstrap_negatives':bootstrap,'source_sha256':hashes,'source_restored_exact':True,'restored_baseline_exit':0,'race':os.name!='nt','whole_native_acceptance':False},separators=(',',':')),flush=True)
 finally:b.cleanup_temporary_directory(directory)
if __name__=='__main__':main()
