from pathlib import Path
import collections,hashlib,json,re,subprocess,sys
root=Path(sys.argv[1]).resolve();repo=Path(sys.argv[2]).resolve();first_repo=Path(sys.argv[3]).resolve()
SHA='8c47d8105c7c3a6dec02f4a9b22958ab58b4682a';PREVIEW='01547e02e4489184e48c0b421d5c0a71289a0c7f'
def digest(b):return hashlib.sha256(b).hexdigest()
def read(run,name):return json.loads((root/str(run)/(name+'.json')).read_bytes())
def clean(b):return [re.sub(r'^\S+Z ','',s) for s in b.decode('utf-8-sig').splitlines()]
def git_bytes(sha,path,where=repo):return subprocess.check_output(['git','show',sha+':'+path],cwd=where)
def save(name,data):
 p=root/name;p.parent.mkdir(parents=True,exist_ok=True)
 with p.open('x',encoding='utf-8',newline='\n') as f:f.write(json.dumps(data,indent=2)+'\n')
def controls(run,platform):
 ids=[]
 for p in (root/str(run)).glob('*.log'):
  lines=clean(p.read_bytes())
  if not any(s==f'Complete job name: {platform}-{w}' for s in lines for w in ['native','controls-1','controls-2','controls-3']):continue
  ids += [s for s in lines if 'mutation rejected:' in s or s.startswith('Bootstrap checksum control:')]
 assert len(ids)==len(set(ids))
 return sorted(ids)
baseline={p:controls(39009853499,p) for p in ['linux','windows']}
assert len(baseline['linux'])==455 and len(baseline['windows'])==450
assert digest(('\n'.join(baseline['linux'])+'\n').encode())=='6c33168bc159995a10934d2891819e2ccc519a0d221e8ff1fe9b227012e3dc3b'
assert digest(('\n'.join(baseline['windows'])+'\n').encode())=='e7894d9a6f3e6b6cccc4b9c9fecfa2007588dc1664401ba759ed9648aeaa61f7'
inventory=json.loads((root/'scope-redaction-new-controls.json').read_bytes())
assert len(inventory['expected_new_full_lines'])==81
for path,h in inventory['scripts'].items():assert digest(git_bytes(SHA,path))==h
plan_raw=git_bytes(SHA,'scripts/tests/sql-report-controls.plan.json')
assert digest(plan_raw)=='fb5be2d44ba2aad6db98beebea79a043c415e394a1946f7753972f11f6176219'
assert digest(git_bytes(SHA,'scripts/tests/sql-report-controls.py'))=='d217822ee4952810fe67cf855c006ba51e2b4682eecb07b22d5e710d965ddd4a'
plan=json.loads(plan_raw);assert len(plan['mutations'])==15
report_ids=['SQL report mutation rejected: '+m['name']+'; compilation succeeded; named assertion failed; exact source restored; baseline passed' for m in plan['mutations']]
assert len(set(report_ids))==15
meta=read(39165239413,'run');jobs=read(39165239413,'jobs')['jobs']
expected={'quality','windows-validation'}|{p+'-'+w for p in ['linux','windows'] for w in ['native','controls-1','controls-2','controls-3']}
assert meta['head_sha']==SHA and meta['conclusion']=='success' and len(jobs)==10
assert {j['name'] for j in jobs}==expected and all(j['conclusion']=='success' for j in jobs)
logs=[]
for j in jobs:
 raw=(root/'39165239413'/(str(j['id'])+'.log')).read_bytes();lines=clean(raw)
 assert PREVIEW in '\n'.join(lines) and not any(s.startswith('##[error]') for s in lines)
 if j['name'] not in ['quality','windows-validation']:assert any(('windows' if j['name'].startswith('windows-') else 'ubuntu') in label for label in j['labels'])
 for step in j['steps']:
  if step['conclusion']!='success':assert j['name']=='linux-native' and step['name']=='Retain synthetic QA failure evidence' and step['conclusion']=='skipped'
 logs.append({'id':j['id'],'name':j['name'],'bytes':len(raw),'sha256':digest(raw),'started_at':j['started_at'],'completed_at':j['completed_at']})
result={}
for platform,count,pre_h in [('linux',551,'1c98bb5df02c5e3dba2a5b7fbc5cbd639b1d725e0dafde06f93e85e6bfd67c0c'),('windows',546,'7901009864c4c9fa37ac0e9e6f11f0b3423104dac85b9491946616fdfdb4f1cd')]:
 ids=controls(39165239413,platform)
 prior=sorted(baseline[platform]+inventory['expected_new_full_lines'])
 assert digest(('\n'.join(prior)+'\n').encode())==pre_h
 assert ids==sorted(prior+report_ids) and len(ids)==count
 result[platform]={'count':count,'missing':0,'extra':0,'duplicates':0,'sha256_full_lines_lf':digest(('\n'.join(ids)+'\n').encode()),'predecessor':len(prior)}
def source(run,sha,where):
 smeta=read(run,'run');sjobs=read(run,'jobs')['jobs']
 assert smeta['head_sha']==sha and smeta['conclusion']=='success' and len(sjobs)==1 and sjobs[0]['conclusion']=='success'
 assert all(s['conclusion']=='success' for s in sjobs[0]['steps'])
 raw=(root/str(run)/(str(sjobs[0]['id'])+'.log')).read_bytes();lines=clean(raw)
 assert not any(s.startswith('##[error]') for s in lines)
 provenance=[json.loads(s.split(' ',1)[1]) for s in lines if s.startswith('REGION_CAPTURE_PROVENANCE ')]
 assert provenance==[{'source':'8e4f0577f3615e6c9014c031bcad079f235369cc','tree':'17d93b20c303f90f7843036be82f0dc32f3260f1','aprl':'60eaddda76541f6adbc1c5ffa686829807e55e29'}]
 chunks=collections.defaultdict(list)
 for s in lines:
  if s.startswith('REGION_CAPTURE_JSON '):
   value=json.loads(s.split(' ',1)[1]);chunks[value['name']].append(value)
 assert len(chunks)==23
 captures=[]
 for name,parts in sorted(chunks.items()):
  assert [p['index'] for p in parts]==list(range(parts[0]['count'])) and all(p['count']==len(parts) for p in parts)
  b=''.join(p['content'] for p in parts).encode()
  assert all(p['sha256']==digest(b) for p in parts)
  assert b==git_bytes(sha,'internal/plugins/region/testdata/'+name,where)
  captures.append({'name':name,'chunks':len(parts),'bytes':len(b),'sha256':digest(b)})
 assert sum(x['chunks'] for x in captures)==315 and sum(x['bytes'] for x in captures)==909865
 summary={'sha':sha,'source_run':run,'captures':captures,'chunks':315,'source_bytes':909865,'source_log':{'bytes':len(raw),'sha256':digest(raw)}}
 save(str(run)+'/review.json',summary);return summary
src=source(39165239406,SHA,repo)
source(39160871913,'fd5aa54f451c58e0a565030d8bb0800bf418bf92',first_repo)
prep=read(39160871908,'run');assert prep['head_sha']=='fd5aa54f451c58e0a565030d8bb0800bf418bf92' and prep['conclusion']=='failure'
prep_jobs=read(39160871908,'jobs')['jobs'];observed=[]
for job_id,name,h in [(117773714360,'windows-native','5af355de42251f0759cac79c06515c57fc9d3341adb79c6ebdd89eb6ad74a33f'),(117773714247,'linux-native','fd852bcb24f05e85bc26d496b22d79a281e440b7c4a1f3a19b0ae58d91516949')]:
 j=next(j for j in prep_jobs if j['id']==job_id);assert j['name']==name and j['conclusion']=='failure'
 b=(root/'39160871908'/(str(job_id)+'.log')).read_bytes();s=b.decode('utf-8-sig')
 assert digest(b)==h and 'undefined: ProjectSQLCapabilityTable' in s and 'undefined: CollectSQLCapabilityTable' in s and 'syntax error' not in s
 observed.append({'job':job_id,'name':name,'bytes':len(b),'sha256':h,'bothAPIsAbsent':True})
save('report-missing-api-review.recovered.json',{'firstRunner':prep['head_sha'],'quality':39160871908,'observed':observed,'beforeProduction':'Previously inspected/published at19:59:34.595UTC; published PR167 snapshot retained. Fresh hosted replay checks complete original native logs, not report runtime coverage.'})
linux=(root/'39165239413/117786690089.log').read_text(encoding='utf-8-sig')
coverage=re.findall(r'Total statement coverage: ([0-9]+(?:\.[0-9]+)?)%',linux);assert coverage==['87.2']
for job_id in [117786690089,117786690197]:
 text=(root/'39165239413'/(str(job_id)+'.log')).read_text(encoding='utf-8-sig')
 assert 'Your code is affected by 0 vulnerabilities.' in text
 assert '0 vulnerabilities in packages you import and 1' in text and 'vulnerability in modules you require' in text
summary={'sha':SHA,'preview':PREVIEW,'run':39165239413,'source_run':39165239406,'complete_quality':True,'controls':result,'jobs':logs,'source_log':src['source_log'],'captures':src['captures'],'coveragePercent':87.2,'floorUnchanged':75,'reachableVulnerabilities':0,'importedVulnerabilities':0,'moduleOnlyAdvisory':'OPEN:1','scope':'Own exact-head complete original CI logs; recovered immutable baseline; fresh hosted deterministic review. No recovered helper was executed.'}
save('39165239413/review.json',summary)
print('QUALITY_REVIEW '+json.dumps({'sha':SHA,'controls':result,'coverage':87.2,'source':{'captures':23,'chunks':315,'bytes':909865}},sort_keys=True),flush=True)
