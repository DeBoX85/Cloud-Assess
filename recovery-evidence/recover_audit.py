from pathlib import Path
import os,sys,json,re,hashlib,gzip,base64,subprocess,urllib.request,urllib.error,urllib.parse,copy,importlib.util
ROOT=Path.cwd();REPO=ROOT/'candidate';OUT=ROOT/'package';OUT.mkdir(exist_ok=True)
SHA='1ededa501e663d888a905d6585dafff0d11e9087';TREE='4dce2dbd7e1a9c71643084905423b46606c3b4ce';PREVIEW='5155c76118bb0623aa8540638275f74260c26789';BASE='eb554ee826187bef05126d3c59e81eb4bd68ecf8';PROTECTED='a6fee427e19feeef28691737ea9bc5026086efa0'
API='https://api.github.com/repos/DeBoX85/Cloud-Assess';TOKEN=os.environ['GH_TOKEN'];members={}
def digest(b):return hashlib.sha256(b).hexdigest()
class NoRedirect(urllib.request.HTTPRedirectHandler):
 def redirect_request(self,*args,**kwargs):return None
def fetch(path,log=False):
 assert path.startswith('/') and '..' not in path
 req=urllib.request.Request(API+path,headers={'Authorization':'Bearer '+TOKEN,'Accept':'application/vnd.github+json','User-Agent':'cloud-assess-evidence-recovery'})
 try:
  with urllib.request.build_opener(NoRedirect()).open(req,timeout=45) as response:data=response.read(8*1024*1024+1)
 except urllib.error.HTTPError as e:
  if not log or e.code not in (301,302,303,307,308):raise RuntimeError('GitHub API HTTP '+str(e.code)) from None
  location=e.headers.get('Location','');u=urllib.parse.urlsplit(location)
  assert u.scheme=='https' and not u.username and not u.password and not u.fragment and any(u.hostname.endswith(s) for s in ['.blob.core.windows.net','.amazonaws.com','.actions.githubusercontent.com'])
  # The signed log redirect receives NO GitHub authorization header.
  try:
   with urllib.request.urlopen(urllib.request.Request(location,headers={'User-Agent':'cloud-assess-evidence-recovery'}),timeout=45) as response:data=response.read(8*1024*1024+1)
  except urllib.error.HTTPError as e:raise RuntimeError('Log storage HTTP '+str(e.code)) from None
 assert len(data)<=8*1024*1024
 return data
def obj(path):return json.loads(fetch(path))
def git(*args):return subprocess.check_output(['git',*args],cwd=REPO)
def clean(raw):return [re.sub(r'^\S+Z ','',x) for x in raw.decode('utf-8-sig').splitlines()]
names={'quality','windows-validation'}|{p+'-'+w for p in ['linux','windows'] for w in ['native','controls-1','controls-2','controls-3']}
def verify_metadata(run,jobs,sha):
 assert run['head_sha']==sha and run['status']=='completed' and run['conclusion']=='success' and run['run_attempt']==1
 rows=jobs['jobs'];assert jobs['total_count']==len(rows)==10 and {j['name'] for j in rows}==names and len({j['id'] for j in rows})==10
 for j in rows:
  assert j['run_id']==run['id'] and j['conclusion']=='success' and j['status']=='completed'
  assert all(s['conclusion']=='success' or j['name']=='linux-native' and s['name']=='Retain synthetic QA failure evidence' and s['conclusion']=='skipped' for s in j['steps'])
  if j['name'] not in ['quality','windows-validation']:assert any(('windows' if j['name'].startswith('windows-') else 'ubuntu') in label for label in j['labels'])
def verify_log(raw,name,preview=None):
 ls=clean(raw);assert 'Complete job name: '+name in ls and 'Cleaning up orphan processes' in ls
 assert not any(x.startswith('##[error]') for x in ls)
 if preview:assert preview in '\n'.join(ls)
 return ls
def ids(rows,logs,platform):
 out=[]
 for j in rows:
  if j['name'].startswith(platform+'-'):out += [x for x in clean(logs[j['id']]) if 'mutation rejected:' in x or x.startswith('Bootstrap checksum control:')]
 assert len(out)==len(set(out))
 return sorted(out)
def verify_ids(actual,previous,count,expected_hash):
 assert actual==previous and len(actual)==count and digest(('\n'.join(actual)+'\n').encode())==expected_hash
def rejected(name,fn):
 try:fn()
 except (AssertionError,KeyError,ValueError):print('RECOVERY_GUARD_REJECTED '+name,flush=True);return
 raise AssertionError('recovery guard failed: '+name)
assert git('rev-parse','HEAD').decode().strip()==SHA and git('rev-parse','HEAD^{tree}').decode().strip()==TREE
assert git('show','-s','--format=%P',SHA).decode().strip()==BASE
paths=git('diff','--name-only',BASE,SHA).decode().splitlines();assert len(paths)==16 and all(p.endswith('.md') and (p.startswith('docs/') or p=='README.md') for p in paths)
preview=obj('/git/commits/'+PREVIEW);assert preview['tree']['sha']==TREE and [p['sha'] for p in preview['parents']]==[PROTECTED,SHA]
quality=obj('/actions/runs/39317160104');qjobs=obj('/actions/runs/39317160104/jobs?per_page=100');previous=obj('/actions/runs/39182715043');pjobs=obj('/actions/runs/39182715043/jobs?per_page=100')
verify_metadata(quality,qjobs,SHA);verify_metadata(previous,pjobs,BASE)
bad=copy.deepcopy(quality);bad['head_sha']=BASE;rejected('wrong-head',lambda:verify_metadata(bad,qjobs,SHA))
bad=copy.deepcopy(qjobs);bad['jobs']=bad['jobs'][:-1];rejected('missing-worker',lambda:verify_metadata(quality,bad,SHA))
bad=copy.deepcopy(qjobs);bad['jobs'][0]=bad['jobs'][1];rejected('duplicate-worker',lambda:verify_metadata(quality,bad,SHA))
bad=copy.deepcopy(qjobs);bad['jobs'][0]['conclusion']='cancelled';rejected('cancelled-worker',lambda:verify_metadata(quality,bad,SHA))
logs={};plogs={};reviews=[]
for run,jobs,target,preview_sha in [(quality,qjobs,logs,PREVIEW),(previous,pjobs,plogs,None)]:
 prefix=str(run['id']);members[prefix+'/run.json']=(json.dumps(run,indent=2)+'\n').encode();members[prefix+'/jobs.json']=(json.dumps(jobs,indent=2)+'\n').encode()
 for j in jobs['jobs']:
  raw=fetch('/actions/jobs/'+str(j['id'])+'/logs',True);verify_log(raw,j['name'],preview_sha);target[j['id']]=raw;members[prefix+'/'+str(j['id'])+'.log']=raw
  if run==quality:reviews.append({'id':j['id'],'name':j['name'],'bytes':len(raw),'sha256':digest(raw),'labels':j['labels'],'steps':j['steps'],'started_at':j['started_at'],'completed_at':j['completed_at']})
sample=qjobs['jobs'][0];rejected('truncated-log',lambda:verify_log(logs[sample['id']][:200],sample['name'],PREVIEW))
controls={}
for platform,count,expected_hash in [('linux',580,'db859a2a2cfbcc9db90251884b4d8b48677374ccab4e1eff90a72fa974533ccf'),('windows',575,'2733b1526e4d2000e2e6070dbc11086a39bc0cc4d35acc50da418e445e4142e0')]:
 actual=ids(qjobs['jobs'],logs,platform);old=ids(pjobs['jobs'],plogs,platform);verify_ids(actual,old,count,expected_hash)
 rejected(platform+'-missing-control',lambda:verify_ids(actual[:-1],old,count,expected_hash))
 controls[platform]={'count':count,'sha256_full_lines_lf':expected_hash,'missing':0,'extra':0,'duplicates':0}
native={j['name']:'\n'.join(clean(logs[j['id']])) for j in qjobs['jobs'] if j['name'].endswith('-native')}
for text in native.values():assert 'Your code is affected by 0 vulnerabilities.' in text and '0 vulnerabilities in packages you import and 1' in text and 'vulnerability in modules you require' in text
coverage=re.findall(r'Total statement coverage: ([0-9.]+)%',native['linux-native']);assert coverage==['87.4']
source=obj('/actions/runs/39317160083');sjobs=obj('/actions/runs/39317160083/jobs?per_page=100');assert source['head_sha']==SHA and source['conclusion']=='success' and sjobs['total_count']==len(sjobs['jobs'])==1
sj=sjobs['jobs'][0];assert sj['conclusion']=='success' and all(s['conclusion']=='success' for s in sj['steps']);raw=fetch('/actions/jobs/'+str(sj['id'])+'/logs',True);ls=verify_log(raw,sj['name'])
def reconstruct(lines):
 provenance=[json.loads(s.split(' ',1)[1]) for s in lines if s.startswith('REGION_CAPTURE_PROVENANCE ')];assert provenance==[{'source':'8e4f0577f3615e6c9014c031bcad079f235369cc','tree':'17d93b20c303f90f7843036be82f0dc32f3260f1','aprl':'60eaddda76541f6adbc1c5ffa686829807e55e29'}]
 chunks={}
 for s in lines:
  if s.startswith('REGION_CAPTURE_JSON '):
   p=json.loads(s.split(' ',1)[1]);chunks.setdefault(p['name'],[]).append(p)
 assert len(chunks)==23
 captures=[]
 for name,parts in sorted(chunks.items()):
  assert [p['index'] for p in parts]==list(range(parts[0]['count'])) and all(p['count']==len(parts) for p in parts)
  b=''.join(p['content'] for p in parts).encode();assert all(p['sha256']==digest(b) for p in parts) and b==git('show',SHA+':internal/plugins/region/testdata/'+name)
  captures.append({'name':name,'chunks':len(parts),'bytes':len(b),'sha256':digest(b)})
 assert sum(p['chunks'] for p in captures)==315 and sum(p['bytes'] for p in captures)==909865
 return captures
captures=reconstruct(ls)
# Verify every historical bundle before reading only allowlisted data; never
# execute a helper recovered from evidence. Reconstruct the lost local-only
# stronger consecutive conservation proof from the published original logs.
spec=importlib.util.spec_from_file_location('evidence_verifier',REPO/'scripts/verify-evidence-bundle.py')
verifier=importlib.util.module_from_spec(spec);spec.loader.exec_module(verifier)
checkpoint=ROOT/'checkpoint/docs/evidence/20261019';historical={};total_decoded=0;verified_parts=0
for prefix,limit in [('audit-72h-part',6),('audit-72h-events-part',3)]:
 for number in range(1,limit+1):
  stem=prefix+str(number);bb=verifier.bounded_read(checkpoint/(stem+'.json.gz.base64'),verifier.MAX_BUNDLE);mm=verifier.bounded_read(checkpoint/(stem+'.manifest.json'),verifier.MAX_MANIFEST)
  verifier.verify(bb,mm);decoded=gzip.decompress(base64.b64decode(bb.strip(),validate=True));total_decoded+=len(decoded);assert total_decoded<=128*1024*1024
  payload_data=json.loads(decoded)
  for member in payload_data['members']:
   name=member['path']
   if re.fullmatch(r'pr\d+/(metadata|review)\.json',name) or re.fullmatch(r'pr\d+/\d+/\d+\.log',name):
    assert name not in historical;historical[name]=member['text'].encode('utf-8')
  verified_parts+=1
edges=[];prior=None;prior_pr=None;successful=0
for number in range(148,176):
 item=json.loads(historical['pr'+str(number)+'/review.json'])
 if item['status']!='success':continue
 successful+=1;current={}
 for platform in ['linux','windows']:
  values=[]
  for j in item['jobs']:
   if j['name'].startswith(platform+'-'):
    path='pr'+str(number)+'/'+str(j['run'])+'/'+str(j['id'])+'.log';values += [x for x in clean(historical[path]) if 'mutation rejected:' in x or x.startswith('Bootstrap checksum control:')]
  assert len(values)==len(set(values));ordered=sorted(values);record=item['controls'][platform]
  assert len(ordered)==record['count'] and digest(('\n'.join(ordered)+'\n').encode())==record['sha256_full_lines_lf']
  current[platform]=set(ordered)
  if prior is not None:
   assert prior[platform]<=current[platform];edges.append({'from':prior_pr,'to':number,'platform':platform,'conserved':len(prior[platform]),'added':len(current[platform]-prior[platform]),'missing':0})
  else:assert len(ordered)==(385 if platform=='linux' else 380)
 prior=current;prior_pr=number
assert verified_parts==9 and successful==22 and len(edges)==42 and len(prior['linux'])==580 and len(prior['windows'])==575
historical_proof={'bounded_parts_verified':verified_parts,'successful_heads':successful,'consecutive_platform_edges':42,'edges':edges,'restored_helpers_executed':False}
members['historical-conservation-reconstructed.json']=(json.dumps(historical_proof,indent=2)+'\n').encode()
print('HISTORICAL_CONSERVATION '+json.dumps({'heads':22,'edges':42,'missing':0,'boundedParts':9}),flush=True)

removed=False;bad=[]
for line in ls:
 if not removed and line.startswith('REGION_CAPTURE_JSON '):removed=True;continue
 bad.append(line)
rejected('missing-capture-chunk',lambda:reconstruct(bad))
members['39317160083/run.json']=(json.dumps(source,indent=2)+'\n').encode();members['39317160083/jobs.json']=(json.dumps(sjobs,indent=2)+'\n').encode();members['39317160083/'+str(sj['id'])+'.log']=raw
report={'candidate':SHA,'tree':TREE,'parent':BASE,'preview':PREVIEW,'quality':quality['id'],'source':source['id'],'controls':controls,'coveragePercent':87.4,'unchangedCoverageFloor':75,'moduleOnlyAdvisory':'OPEN:1','jobs':reviews,'captures':captures,'sourceLog':{'bytes':len(raw),'sha256':digest(raw)},'historicalConservation':historical_proof,'recoveryRecipe':os.environ['GITHUB_SHA'],'recoveryRun':os.environ['GITHUB_RUN_ID'],'limits':'Current original hosted evidence reconstructed after local environment409offline; local-only new helper/preflight proof files not fabricated.'}
members['review.json']=(json.dumps(report,indent=2)+'\n').encode()
members['recovery/recover_audit.py']=Path(__file__).read_bytes()
members['recovery/workflow.yml']=(ROOT/'control/.github/workflows/audit-handover-evidence-recovery.yml').read_bytes()
for name in ['CURRENT_STATE.md','RESTART_PROMPT.md','CHAT_TRANSFER_20261019.md','AUDIT_72H_20261019.md','DISCREPANCIES_20261019.md','METER_HISTORY_COMPLETED_20261019.md']:members['candidate/docs/'+name]=git('show',SHA+':docs/'+name)
status='# Executor outage and recovery status\n\nExact handover code remains PR176/'+SHA+'/tree'+TREE+'. Own8workers/bothgates/source passed;580Linux/575Windows identities,87.4%coverage/75%floor,23captures/315chunks/909865bytes verified against exact Git. Local execution disconnected after the code/historical archives were published; fresh hosted read-only reconstruction retained original current/predecessor logs, raw metadata, full controls/source proof and tested negative guards using the unchanged bounded verifier. No local-only unavailable proof is invented. All16candidate Markdown files remain frozen; this storage addendum is the material current executor update, not a coding base. Source/data pins and protected refs unchanged. No critical blocker without a known mitigation found. A new chat must verify its own execution environment or use the tested isolated hosted pattern. Next is task6retail contract; no retail implementation or Azure API request has occurred. Module-only advisory, independent review/backup, operational/live/release obligations remain OPEN.\n'
restart='# Restart after executor interruption\n\nContinue DeBoX85/Cloud-Assess from branchdocs/72h-audit-handover-20261019, exact1ededa501e663d888a905d6585dafff0d11e9087, livePR176 reached through175. Read AGENTS.md, CURRENT_STATE, RESTART_PROMPT, CHAT_TRANSFER_20261019, AUDIT_72H_20261019 and DISCREPANCIES_20261019 there, plus this RECOVERY_STATUS addendum and live176final outcomes. The previous local executor disconnected; do not assume its scratch/tool paths are available or that a local-only log exists. Original code/evidence and all16documents are published. Hosted recovery recipe reconstructs original current logs/control/source/metadata and bounded-verifies without changing the candidate. Never select this archive or recipe branch as coding base. Verify refs/head/tree/parents, full evidence/readback and actual executor before editing. Tasks1-5complete offline; all75windowed runs/444jobs/416available logs audited, with skipped/cancelled metadata explicitly not passes. Then task6retail contract/literals only, researching at90%confidence or lower. No retail implementation yet. Preserve all prior heads/fixtures/pins, privacy fallback169->170 and unrelated work. Continue authorised independent work; stop/document dedicated blocker/wait for Denis only for a critical issue with no clear mitigation. No Azure/fixtures/RBAC/refresh/subagents/protected merge/bypass/release. Use Denis Bogunic identity; keep advisory/human/independent/live/release obligations open.\n'
members['recovery/RECOVERY_STATUS_20261019.md']=status.encode();members['recovery/RESTART_AFTER_EXECUTOR_OUTAGE_20261019.md']=restart.encode()
payload=json.dumps({'format':'cloud-assess-evidence-v1','members':[{'path':p,'text':b.decode('utf-8')} for p,b in sorted(members.items())]},ensure_ascii=True,separators=(',',':')).encode();bundle=base64.b64encode(gzip.compress(payload,mtime=0))+b'\n'
manifest={'format':'cloud-assess-evidence-v1','classification':'reviewed-synthetic-ci-local','bundle_sha256':digest(bundle),'decoded_sha256':digest(payload),'members':[{'path':p,'bytes':len(b),'sha256':digest(b)} for p,b in sorted(members.items())],'revisions':{'candidate':SHA,'preview':PREVIEW,'protected':PROTECTED,'product_predecessor':BASE,'retrospective_archive':'2fa53c346a19798e4ad90124d26275a32d42ae91','recovery_recipe':os.environ['GITHUB_SHA']}}
(OUT/'.gitattributes').write_bytes(b'* -text\n');(OUT/'audit-handover-final.json.gz.base64').write_bytes(bundle);(OUT/'audit-handover-final.manifest.json').write_bytes((json.dumps(manifest,indent=2)+'\n').encode())
subprocess.run([sys.executable,str(REPO/'scripts/tests/evidence-bundle.py')],cwd=REPO,check=True)
subprocess.run([sys.executable,str(REPO/'scripts/verify-evidence-bundle.py'),'--bundle',str(OUT/'audit-handover-final.json.gz.base64'),'--manifest',str(OUT/'audit-handover-final.manifest.json')],check=True)
summary={'candidate':SHA,'members':len(members),'memberBytes':sum(len(b) for b in members.values()),'bundleBytes':len(bundle),'bundleSHA256':digest(bundle),'recoveryRecipe':os.environ['GITHUB_SHA'],'recoveryRun':os.environ['GITHUB_RUN_ID'],'controls':controls,'coverage':87.4}
index='# PR176 exact-head handover QA and executor recovery\n\n'+status+'\nOwn quality39317160104/source39317160083 complete. All raw job/step/runner metadata, original complete logs and predecessor control logs retained; failed/missing-head/control/chunk guards reject under explicit negative probes. Existing bounded-verifier regression suite and output verification passed. The separate retrospective2fa53c34 retains75windowed runs/444jobs/416available logs and32source reconstructions; six skipped/22unassigned-cancelled jobs have no logs and remain nonacceptance. Current recovery is outside that request-time window. Local-only new helper/research/preflight files became inaccessible; none is falsely retained. Product17616files were read back before outage, and no candidate byte was changed for recovery.\n\n[Bundle](audit-handover-final.json.gz.base64), [manifest](audit-handover-final.manifest.json), [current recovery status](RECOVERY_STATUS_20261019.md), [copy-ready restart](RESTART_AFTER_EXECUTOR_OUTAGE_20261019.md). Read restored helpers before executing. Archive storage is not another code certificate, coding base, independent backup or protected/live/release acceptance.\n\n'+json.dumps(summary,indent=2)+'\n'
(OUT/'PR176_AUDIT_HANDOVER_EVIDENCE.md').write_bytes(index.encode());(OUT/'RECOVERY_STATUS_20261019.md').write_bytes(status.encode());(OUT/'RESTART_AFTER_EXECUTOR_OUTAGE_20261019.md').write_bytes(restart.encode())
print('RECOVERY_PACKAGE '+json.dumps(summary,sort_keys=True),flush=True)
for f in sorted(OUT.iterdir()):
 b=f.read_bytes();s=b.decode('utf-8');parts=[s[i:i+8192] for i in range(0,len(s),8192)]
 print('EVIDENCE_FILE_META '+json.dumps({'name':f.name,'bytes':len(b),'sha256':digest(b),'parts':len(parts)}),flush=True)
 for i,part in enumerate(parts):print('EVIDENCE_FILE_CHUNK '+json.dumps({'name':f.name,'index':i,'count':len(parts),'content':part},ensure_ascii=False,separators=(',',':')),flush=True)
