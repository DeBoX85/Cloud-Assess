#!/usr/bin/env python3
"""Isolated Build proof while admission is absent. Never whole-native acceptance."""
import hashlib,json,os,re,shutil,subprocess,tempfile,sys
sys.dont_write_bytecode=True
from pathlib import Path
from qa_temp_cleanup import cleanup_temporary_directory

ROOT=Path(__file__).resolve().parents[2]
SOURCE='internal/plugins/region/required_inventory.go'
TESTS=['required_inventory_fixtures_test.go','required_inventory_build_test.go','required_inventory_build_budget_test.go']
CONTROLS=[
 ('text-excess-two','len(s) > auxTextBudget-b.text','len(s) > auxTextBudget+2-b.text','TestRequiredBuildFrozenBudgetRecipes/build-output-text-first-reachable-excess','BUILD_BUDGET_REJECTION'),
 ('entry-excess-two','n > MaxInventoryEntries-b.entries','n > MaxInventoryEntries+2-b.entries','TestRequiredBuildFrozenBudgetRecipes/build-output-entries-first-reachable-excess','BUILD_BUDGET_REJECTION'),
 ('text-exact-rejected','len(s) > auxTextBudget-b.text','len(s) >= auxTextBudget-b.text','TestRequiredBuildFrozenBudgetRecipes/build-output-text-exact','BUILD_BUDGET_EXACT'),
 ('entry-exact-rejected','n > MaxInventoryEntries-b.entries','n >= MaxInventoryEntries-b.entries','TestRequiredBuildFrozenBudgetRecipes/build-output-entries-exact','BUILD_BUDGET_EXACT'),
 ('seventeenth-storage-kind','len(kinds) >= 16','len(kinds) >= 17','TestRequiredBuildStorageKindFirstExcess','BUILD_STORAGE_SEVENTEENTH_KIND'),
 ('final-terminal-escape','return ctx.Err()\n}','return nil\n}','TestRequiredBuildContext/before-return-cancel-healthy','BUILD_LATE_TERMINAL'),
 ('public-private-resource-alias','Resources: append([]RequiredInventoryResource{}, derived...)','Resources: privateResources','TestRequiredBuildOwnership','BUILD_DIAGNOSTIC_OWNERSHIP'),
 ('public-private-requirement-alias','publicRequirements = append(publicRequirements, requiredCopyRequirement(r))','publicRequirements = append(publicRequirements, privateRequirements[len(privateRequirements)-1])','TestRequiredBuildOwnership','BUILD_DIAGNOSTIC_OWNERSHIP'),
 ('full-terminal-type-binding','strings.EqualFold(id.ResourceType.String(), r.Type)','true','TestRequiredBuildFrozenLiterals/id-type-mismatch','inventory rejection:'),
 ('cognitive-kind-collapse','keyKind = kind','keyKind = ""','TestRequiredBuildFrozenLiterals/cognitive-same-sku-different-kinds','requirements mismatch'),
 ('parsed-type-unicode-alias','providerTypeIdentifier(id.ResourceType.String()) && ','','TestRequiredBuildTypeUnicodeAliasRejected','BUILD_TYPE_ASCII_IDENTITY'),
 ('cognitive-unicode-fold','if storageASCII(raw.Kind) {','if true {','TestRequiredBuildUnrepresentableCognitiveKindIsolation','BUILD_COGNITIVE_KIND_ISOLATION'),
]

def events(output):
 rows=[]
 for line in output.splitlines():
  if line.strip():
   v=json.loads(line)
   if not isinstance(v,dict) or not isinstance(v.get('Action'),str):raise ValueError('invalid Go event')
   rows.append(v)
 return rows

def intended_failure(code,stdout,stderr,expected,marker,timed_out=False):
 if timed_out or code!=1 or stderr.strip():return False
 if any(s in stdout for s in ['panic:','fatal error:','[build failed]','undefined:','signal:','test timed out','DATA RACE']):return False
 try:rows=events(stdout)
 except (ValueError,TypeError):return False
 if any(x['Action'] in ['build-fail','skip'] for x in rows):return False
 failed=[x.get('Test') for x in rows if x['Action']=='fail' and x.get('Test')]
 if expected not in failed or any(x not in [expected,expected.split('/')[0]] for x in failed):return False
 return any(x['Action']=='output' and x.get('Test')==expected and marker in x.get('Output','') for x in rows)

def classifier_controls():
 expected='TestExpected/sub';marker='INTENDED_ASSERTION'
 rows=[{'Action':'output','Test':expected,'Output':marker+'\n'},{'Action':'fail','Test':expected}]
 good='\n'.join(json.dumps(x) for x in rows)+'\n'
 assert intended_failure(1,good,'',expected,marker)
 bad=[('zero-exit',(0,good,'',expected,marker)),('compiler',(1,good,'undefined: Missing',expected,marker)),('panic',(1,good+json.dumps({'Action':'output','Output':'panic: '+marker})+'\n','',expected,marker)),('wrong-test',(1,good,'','TestOther',marker)),('wrong-marker',(1,good,'',expected,'OTHER')),('unframed',(1,marker+'\n','',expected,marker)),('no-failure-event',(1,json.dumps(rows[0])+'\n','',expected,marker)),('timeout',(1,good,'',expected,marker,True)),('infrastructure',(1,good,'connection refused',expected,marker)),('race',(1,good+json.dumps({'Action':'output','Output':'DATA RACE'})+'\n','',expected,marker))]
 for name,args in bad:
  if intended_failure(*args):raise ValueError('classifier accepted '+name)
 return [x[0] for x in bad]

def run(root,selector='',race=False):
 region=root/'internal/plugins/region'
 files=sorted(p.name for p in region.glob('*.go') if not p.name.endswith('_test.go'))+TESTS
 cmd=[os.environ.get('QA_GO','go'),'test','-mod=readonly','-count=1','-timeout=45s','-json']
 if race:cmd.append('-race')
 if selector:cmd+=['-run','/'.join('^'+re.escape(s)+'$' for s in selector.split('/'))]
 env=dict(os.environ);env.update(GOTOOLCHAIN='local',GOWORK='off',GOPROXY='off',GOSUMDB='off',GOTELEMETRY='off')
 return subprocess.run(cmd+files,cwd=region,env=env,capture_output=True,text=True,encoding='utf-8',errors='strict',timeout=90)

def main():
 rejected=classifier_controls()
 directory=tempfile.TemporaryDirectory(prefix='cloud-assess-required-build-')
 try:
  work=Path(directory.name)
  names=subprocess.check_output(['git','ls-files','-z'],cwd=ROOT).decode().split('\0')
  for name in names:
   if not name:continue
   p=ROOT/name
   if p.is_file():dest=work/name;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(p,dest)
  # Explicit file-set verification excludes the still-unimplemented admission
  # consumer. The unchanged full consumer remains in every normal native suite.
  baseline=run(work,race=os.name!='nt')
  if baseline.returncode!=0 or baseline.stderr.strip():raise RuntimeError('Build healthy isolated baseline failed\n'+baseline.stdout+baseline.stderr)
  base_events=events(baseline.stdout)
  passed=[x['Test'] for x in base_events if x['Action']=='pass' and x.get('Test')]
  if len([x for x in passed if x.startswith('TestRequiredBuildFrozenLiterals/')])!=76:raise RuntimeError('lost frozen Build scenario')
  if len([x for x in passed if x.startswith('TestRequiredBuildContext/')])!=8:raise RuntimeError('lost terminal schedule')
  if len([x for x in passed if x.startswith('TestRequiredBuildFrozenBudgetRecipes/')])!=6:raise RuntimeError('lost Build recipe')
  for line in baseline.stdout.splitlines():print('REQUIRED_BUILD_BASELINE_JSON '+line,flush=True)
  source=work/SOURCE;original=source.read_bytes();text=original.decode().replace('\r\n','\n');results=[]
  for name,old,replacement,expected,marker in CONTROLS:
   if text.count(old)!=1:raise RuntimeError('production mutation anchor changed: '+name)
   source.write_bytes(text.replace(old,replacement,1).encode())
   compiled=run(work,'^$')
   if compiled.returncode or compiled.stderr.strip():raise RuntimeError('mutation did not compile: '+name+'\n'+compiled.stdout+compiled.stderr)
   failed=run(work,expected)
   if not intended_failure(failed.returncode,failed.stdout,failed.stderr,expected,marker):raise RuntimeError('not the intended assertion: '+name+'\n'+failed.stdout+failed.stderr)
   result={'name':name,'test':expected,'marker':marker,'compile_exit':compiled.returncode,'exit':failed.returncode,'events':events(failed.stdout),'stdout_sha256':hashlib.sha256(failed.stdout.encode()).hexdigest()};results.append(result)
   print('REQUIRED_BUILD_MUTATION_JSON '+json.dumps(result,separators=(',',':')),flush=True)
   print('required inventory Build mutation rejected: '+name,flush=True)
   source.write_bytes(original)
   if source.read_bytes()!=original:raise RuntimeError('disposable production restore failed')
  restored=run(work)
  if restored.returncode or restored.stderr.strip():raise RuntimeError('restored isolated baseline failed\n'+restored.stdout+restored.stderr)
  if (ROOT/SOURCE).read_bytes()!=original:raise RuntimeError('original worktree source changed')
  print('REQUIRED_BUILD_PROOF_JSON '+json.dumps({'scope':'isolated explicit file-set, not whole native acceptance','frozen_build_scenarios':76,'build_terminal_schedules':8,'build_recipes':6,'compiling_assertion_reversals':len(results),'classifier_negative_controls':rejected,'source_sha256':hashlib.sha256(original).hexdigest(),'source_restored_exact':True,'restored_baseline_exit':0,'race':os.name!='nt','admission_implemented':False,'healthy_whole_native_accepted':False},separators=(',',':')),flush=True)
 finally:cleanup_temporary_directory(directory)

if __name__=='__main__':main()
