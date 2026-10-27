"""Bounded Go-JSON assertion admission. Does not modify product or test bytes."""
import hashlib,json,os,re,subprocess
from pathlib import Path

TOOL=Path(__file__).resolve().parents[2]/'tools/mutationassertions/main.go'
_CACHE={}
FORBIDDEN=('panic:','fatal error:','[build failed]','[setup failed]','undefined:','syntax error:','signal:','test timed out','DATA RACE','runtime: out of memory')

def _unique(pairs):
 result={}
 for k,v in pairs:
  if k in result:raise ValueError('duplicate Go event key')
  result[k]=v
 return result

def _catalog(go,directory,cwd):
 files={p:p.read_bytes() for p in directory.glob('*_test.go')}
 if not files:raise RuntimeError('no original test source for assertion admission')
 key=(str(directory),tuple((str(p),hashlib.sha256(b).hexdigest()) for p,b in sorted(files.items())))
 if key not in _CACHE:
  result=subprocess.run([go,'run',str(TOOL),'--',str(directory)],cwd=cwd,capture_output=True,text=True,encoding='utf-8',errors='strict',timeout=90)
  if result.returncode or result.stderr.strip():raise RuntimeError('cannot inventory original assertion source\n'+result.stdout+result.stderr)
  rows=json.loads(result.stdout,object_pairs_hook=_unique)
  for row in rows:
   path=directory/row['file']
   if path not in files or hashlib.sha256(files[path]).hexdigest()!=row['sha256'] or row['line']<1:raise RuntimeError('assertion source inventory mismatch')
  _CACHE[key]=rows
 return _CACHE[key],files

def run_go_test(args,**kwargs):
 if len(args)<2 or args[1]!='test':raise ValueError('only explicit go test commands supported')
 cwd=Path(kwargs.get('cwd','.')).resolve();packages=[a for a in args[2:] if a=='.' or a.startswith('./')]
 if len(packages)!=1 or '...' in packages[0]:raise ValueError('one explicit mutation-test package required')
 directory=(cwd/packages[0]).resolve();sites,originals=_catalog(args[0],directory,cwd)
 command=list(args)
 if '-json' not in command:command.insert(2,'-json')
 if not any(a=='-timeout' or a.startswith('-timeout=') for a in command):command.insert(2,'-timeout=120s')
 if not any(a=='-mod' or a.startswith('-mod=') for a in command):command.insert(2,'-mod=readonly')
 kwargs['timeout']=min(kwargs.get('timeout',180),180)
 try:result=subprocess.run(command,**kwargs)
 finally:
  for p,raw in originals.items():
   if p.read_bytes()!=raw:raise RuntimeError('original test source changed during mutation invocation')
 result.assertion_sites=sites
 return result

def assertion_witness(result,expected):
 stdout=result.stdout or '';stderr=result.stderr or ''
 if not isinstance(stdout,str) or not isinstance(stderr,str) or result.returncode!=1 or stderr.strip() or any(s in stdout+stderr for s in FORBIDDEN):return None
 try:
  rows=[json.loads(line,object_pairs_hook=_unique) for line in stdout.splitlines() if line.strip()]
  if not rows or any(not isinstance(r,dict) or not isinstance(r.get('Action'),str) for r in rows):return None
 except (ValueError,TypeError):return None
 if any(r['Action'] in ('build-fail','skip') for r in rows):return None
 failed=[r.get('Test') for r in rows if r['Action']=='fail' and r.get('Test')]
 if not failed or not any(t==expected or t.startswith(expected+'/') for t in failed):return None
 if any(not(t==expected or t.startswith(expected+'/') or expected.startswith(t+'/')) for t in failed):return None
 if not any(r['Action']=='fail' and not r.get('Test') for r in rows):return None
 root=expected.split('/')[0];sites=getattr(result,'assertion_sites',[])
 for row in rows:
  test=row.get('Test','')
  if row['Action']!='output' or test not in failed or not(test==expected or test.startswith(expected+'/')):continue
  for line in row.get('Output','').splitlines():
   match=re.match(r'^\s*([A-Za-z0-9_.-]+_test\.go):(\d+):\s+(.+)$',line)
   if not match:continue
   site=next((s for s in sites if s['file']==match[1] and s['line']==int(match[2]) and root in s['roots']),None)
   if site:
    return {'expected':expected,'failed_test':test,'assertion_site':site,'diagnostic':line[:1024],'diagnostic_sha256':hashlib.sha256(line.encode()).hexdigest(),'stdout_sha256':hashlib.sha256(stdout.encode()).hexdigest(),'exit':1,'raw_child_output_retained':False}
 return None

def intended_failure(result,expected):
 witness=assertion_witness(result,expected)
 if witness is None:return False
 print('MUTATION_ASSERTION_WITNESS_JSON '+json.dumps(witness,separators=(',',':')),flush=True)
 return True
