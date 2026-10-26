#!/usr/bin/env python3
"""Isolated Build proof. Full native acceptance remains a separate obligation."""
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

def json_stream(text):
 decoder=json.JSONDecoder();rows=[];at=0
 while at<len(text):
  while at<len(text) and text[at].isspace():at+=1
  if at==len(text):break
  value,at=decoder.raw_decode(text,at);rows.append(value)
 return rows

def module_command(command,cwd,env,timeout):
 return subprocess.run(command,cwd=cwd,env=env,capture_output=True,text=True,encoding='utf-8',errors='strict',timeout=timeout)

def bootstrap_modules(root,invoke=module_command):
 # Only the immutable current module graph may fill a cold cache. No go get,
 # @latest query, tidy or dependency-file rewrite is accepted here.
 before={name:(root/name).read_bytes() for name in ['go.mod','go.sum']}
 env=dict(os.environ);env.update(GOTOOLCHAIN='local',GOWORK='off',GOTELEMETRY='off')
 go=os.environ.get('QA_GO','go')
 metadata=invoke([go,'mod','edit','-json'],root,env,30)
 if metadata.returncode:raise RuntimeError('cannot read locked module graph')
 pins={x['Path']:x['Version'] for x in json.loads(metadata.stdout)['Require']}
 downloaded=invoke([go,'mod','download','-json'],root,env,300)
 if downloaded.returncode:raise RuntimeError('locked module cache bootstrap failed\n'+downloaded.stderr)
 rows=json_stream(downloaded.stdout);seen=set()
 for row in rows:
  path=row.get('Path')
  if row.get('Error') or path not in pins or path in seen or row.get('Version')!=pins[path]:raise RuntimeError('download departed from locked module graph')
  seen.add(path)
 if seen!=set(pins):raise RuntimeError('incomplete locked module cache bootstrap')
 verified=invoke([go,'mod','verify'],root,env,90)
 if verified.returncode or verified.stdout.strip()!='all modules verified':raise RuntimeError('module cache verification failed')
 if any((root/name).read_bytes()!=raw for name,raw in before.items()):raise RuntimeError('locked module metadata changed during cache bootstrap')
 return {'locked_modules':len(pins),'modules':sorted(path+'@'+version for path,version in pins.items()),'pin_files':{n:hashlib.sha256(b).hexdigest() for n,b in before.items()},'pin_files_byte_unchanged':True,'module_cache_verified':True,'bootstrap_network':'runner proxy may fetch exact locked versions only','test_phase_proxy':'off'}

def bootstrap_controls():
 directory=tempfile.TemporaryDirectory(prefix='cloud-assess-module-bootstrap-model-')
 try:
  root=Path(directory.name);before={'go.mod':b'module model\n','go.sum':b'locked checksum model\n'}
  def reset():
   for name,raw in before.items():(root/name).write_bytes(raw)
  def fake(mode):
   def invoke(command,cwd,env,timeout):
    if mode=='timeout':raise subprocess.TimeoutExpired(command,timeout)
    if command[1:]==['mod','edit','-json']:return subprocess.CompletedProcess(command,0,json.dumps({'Require':[{'Path':'example.invalid/locked','Version':'v1.2.3'}]}),'')
    if command[1:]==['mod','verify']:return subprocess.CompletedProcess(command,1 if mode=='failed-verification' else 0,'all modules verified\n','')
    if mode=='changed-pin-file':(root/'go.sum').write_bytes(b'changed')
    row={'Path':'example.invalid/locked','Version':'v1.2.4' if mode=='changed-version' else 'v1.2.3'}
    output='{}' if mode=='missing-module' else json.dumps(row)
    if mode=='duplicate-module':output+=json.dumps(row)
    if mode=='malformed-json':output='not JSON'
    return subprocess.CompletedProcess(command,0,output,'')
   return invoke
  reset();bootstrap_modules(root,fake('healthy'))
  rejected=[]
  for mode in ['changed-version','changed-pin-file','failed-verification','missing-module','duplicate-module','malformed-json','timeout']:
   reset()
   try:bootstrap_modules(root,fake(mode))
   except (RuntimeError,ValueError,KeyError,subprocess.TimeoutExpired):rejected.append(mode)
   else:raise RuntimeError('invalid bootstrap accepted: '+mode)
  return rejected
 finally:cleanup_temporary_directory(directory)


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
 bootstrap_rejected=bootstrap_controls()
 directory=tempfile.TemporaryDirectory(prefix='cloud-assess-required-build-')
 try:
  work=Path(directory.name)
  names=subprocess.check_output(['git','ls-files','-z'],cwd=ROOT).decode().split('\0')
  for name in names:
   if not name:continue
   p=ROOT/name
   if p.is_file():dest=work/name;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(p,dest)
  # Explicit Build file-set proof remains separate from the mandatory full
  # consumer and the Admission proof in normal native verification.
  module_proof=bootstrap_modules(work)
  print('REQUIRED_BUILD_MODULES_JSON '+json.dumps(module_proof,separators=(',',':')),flush=True)
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
  print('REQUIRED_BUILD_PROOF_JSON '+json.dumps({'scope':'isolated explicit file-set, not whole native acceptance','frozen_build_scenarios':76,'build_terminal_schedules':8,'build_recipes':6,'compiling_assertion_reversals':len(results),'classifier_negative_controls':rejected,'bootstrap_negative_controls':bootstrap_rejected,'source_sha256':hashlib.sha256(original).hexdigest(),'source_restored_exact':True,'restored_baseline_exit':0,'race':os.name!='nt','admission_implemented':True,'healthy_whole_native_accepted':False},separators=(',',':')),flush=True)
 finally:cleanup_temporary_directory(directory)

if __name__=='__main__':main()
