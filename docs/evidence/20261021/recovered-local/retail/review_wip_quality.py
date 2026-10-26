from pathlib import Path
import collections,hashlib,json,re,subprocess,sys
root=Path(__file__).resolve().parent;repo=Path(sys.argv[1]).resolve();sha,preview=sys.argv[2:4];runid=int(sys.argv[4]);sourceid=int(sys.argv[5]);baseline=root.parent/'39182715043'
def h(b):return hashlib.sha256(b).hexdigest()
def git(*a):return subprocess.check_output(['git',*a],cwd=repo)
def lines(b):return [re.sub(r'^\S+Z ','',s) for s in b.decode('utf-8-sig').splitlines()]
def obj(run,name):return json.loads((root/str(run)/(name+'.json')).read_bytes())
def identities(folder,platform):
 out=[]
 for p in folder.glob('*.log'):
  ss=lines(p.read_bytes())
  if set(ss)&{f'Complete job name: {platform}-{w}' for w in ['native','controls-1','controls-2','controls-3']}:
   out.extend(s for s in ss if 'mutation rejected:' in s or s.startswith('Bootstrap checksum control:'))
 assert len(out)==len(set(out));return sorted(out)
assert git('rev-parse',sha+'^{tree}').strip()==git('rev-parse',preview+'^{tree}').strip()
assert git('show','-s','--format=%P',preview).decode().strip().split()==['a6fee427e19feeef28691737ea9bc5026086efa0',sha]
run,jobs=obj(runid,'run'),obj(runid,'jobs');expected={'quality','windows-validation'}|{p+'-'+w for p in ['linux','windows'] for w in ['native','controls-1','controls-2','controls-3']}
assert run['head_sha']==sha and run['conclusion']=='success' and jobs['total_count']==len(jobs['jobs'])==10
assert {j['name'] for j in jobs['jobs']}==expected
rows=[];native={}
for j in jobs['jobs']:
 assert j['conclusion']=='success' and j.get('runner_id',0)>0
 b=(root/str(runid)/(str(j['id'])+'.log')).read_bytes();ss=lines(b);text='\n'.join(ss)
 assert any('Cleaning up orphan processes' in s for s in ss[-12:]) and f'Complete job name: {j["name"]}' in ss
 assert preview in text and not any(s.startswith('##[error]') for s in ss)
 if j['name'] not in ['quality','windows-validation']:assert any(('windows' if j['name'].startswith('windows-') else 'ubuntu') in label for label in j['labels'])
 for step in j['steps']:
  if step['conclusion']!='success':assert j['name']=='linux-native' and step['name']=='Retain synthetic QA failure evidence' and step['conclusion']=='skipped'
 rows.append({'job':j['id'],'name':j['name'],'bytes':len(b),'sha256':h(b),'started_at':j['started_at'],'completed_at':j['completed_at'],'labels':j['labels'],'steps':j['steps']})
 if j['name'].endswith('-native'):native[j['name']]=text
controls={}
for platform,count,digest in [('linux',580,'db859a2a2cfbcc9db90251884b4d8b48677374ccab4e1eff90a72fa974533ccf'),('windows',575,'2733b1526e4d2000e2e6070dbc11086a39bc0cc4d35acc50da418e445e4142e0')]:
 old=identities(baseline,platform);new=identities(root/str(runid),platform)
 assert old==new and len(new)==count and h(('\n'.join(new)+'\n').encode())==digest
 controls[platform]={'count':count,'sha256_full_lines_lf':digest,'missing':0,'extra':0,'duplicates':0}
for name,text in native.items():
 assert 'Your code is affected by 0 vulnerabilities.' in text
 assert '3 vulnerabilities in packages you import and 1' in text and 'vulnerability in modules you require' in text
coverage=re.findall(r'Total statement coverage: ([0-9]+(?:\.[0-9]+)?)%',native['linux-native']);assert len(coverage)==1 and float(coverage[0])>=75
source=obj(sourceid,'review');assert source['sha']==sha and len(source['captures'])==23 and source['source_bytes']==909865 and source['chunks']==315
from datetime import datetime
elapsed=lambda a,b:(datetime.fromisoformat(b.replace('Z','+00:00'))-datetime.fromisoformat(a.replace('Z','+00:00'))).total_seconds()
worker_timing=[]
for j in jobs['jobs']:
 if j['name'] in ['quality','windows-validation']:continue
 first=j['steps'][0]
 worker_timing.append({'name':j['name'],'job_wall_seconds':elapsed(j['started_at'],j['completed_at']),'creation_to_job_start_seconds':elapsed(run['created_at'],j['started_at']),'runner_setup_step_seconds':elapsed(first['started_at'],first['completed_at'])})
source_jobs=obj(sourceid,'jobs')['jobs'];assert len(source_jobs)==1
source_log=(root/str(sourceid)/(str(source_jobs[0]['id'])+'.log')).read_bytes()
assert any('Cleaning up orphan processes' in x for x in lines(source_log)[-12:])
summary={'candidate' :sha,'preview':preview,'run':runid,'source_run':sourceid,'jobs':rows,'controls':controls,'coverage_percent':float(coverage[0]),'unchanged_floor':75,'native_security':{'reachable':0,'imported_not_called':3,'module_only':1,'residuals':'OPEN; no suppression/waiver'},'source':source,'timing':{'run_started_at':run['run_started_at'],'run_completed_at':run['updated_at'],'first_job_started_at':min(j['started_at'] for j in jobs['jobs']),'last_job_completed_at':max(j['completed_at'] for j in jobs['jobs'])},'worker_timing':worker_timing,'sum_worker_wall_seconds':sum(w['job_wall_seconds'] for w in worker_timing),'timing_note':'Creation-to-start includes scheduling; worker wall sums are not elapsed time or CPU-speedup evidence. Raw job/step timestamps retained.','scope':'exact own head full native/source offline verification, not protected acceptance/independent review/live/release'}
(root/str(runid)/'review.json').write_text(json.dumps(summary,indent=2)+'\n');print(json.dumps({k:summary[k] for k in ['candidate','controls','coverage_percent','native_security','timing']}))
