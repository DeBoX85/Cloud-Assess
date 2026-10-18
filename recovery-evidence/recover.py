from pathlib import Path,PurePosixPath
import base64,datetime,gzip,hashlib,io,json,os,re,subprocess,sys
workspace=Path(os.environ['GITHUB_WORKSPACE']).resolve()
control=workspace/'control';repo=workspace/'candidate';checkpoint=workspace/'checkpoint';first_repo=workspace/'first-runner'
qa=workspace/'recovered-qa';package=workspace/'package'
SHA='8c47d8105c7c3a6dec02f4a9b22958ab58b4682a';PREVIEW='01547e02e4489184e48c0b421d5c0a71289a0c7f'
CHECKPOINT='9ae98507cc55ec268b7a27b4bdf2a14f2f6e6567'
ALLOWLIST=["39009853499/117329867980.log","39009853499/117329868123.log","39009853499/117329868127.log","39009853499/117329868149.log","39009853499/117329868163.log","39009853499/117329868180.log","39009853499/117329868190.log","39009853499/117329868191.log","39009853499/117331844497.log","39009853499/117333184254.log","39009853499/jobs.json","39009853499/review.json","39009853499/run.json","39162596987/117778845592.log","39162596987/117778845785.log","39162596987/117778845796.log","39162596987/117778845814.log","39162596987/117778845815.log","39162596987/117778845821.log","39162596987/117778845832.log","39162596987/117778845844.log","39162596987/117781203184.log","39162596987/117783463780.log","39162596987/jobs.final.json","39162596987/red-review.json","39162596987/run.final.json","39162596990/117778845585.log","39162596990/jobs.json","39162596990/review.json","39162596990/run.json","39163414427/117781298822.log","39163414427/jobs.json","39163414427/review.json","39163414427/run.json","39163414431/117781299107.log","39163414431/117781299282.log","39163414431/117781299283.log","39163414431/117781299287.log","39163414431/117781299304.log","39163414431/117781299317.log","39163414431/117781299324.log","39163414431/117781299370.log","39163414431/117783784842.log","39163414431/117785451276.log","39163414431/jobs.final.json","39163414431/review.json","39163414431/run.final.json","build-sql-report-evidence.py","build-sql-report-immutable-evidence.py","review-sql-immutable.py","review-sql-report-red.py","review-sql-source.py","review-sql.py","scope-redaction-controls-local.exit.json","scope-redaction-controls-local.log","scope-redaction-core-records-pre-fix.exit.json","scope-redaction-core-records-pre-fix.log","scope-redaction-correction-local.exit.json","scope-redaction-correction-local.log","scope-redaction-final-evidence-inventory.json","scope-redaction-native-before-production.json","scope-redaction-native-review.final.json","scope-redaction-native-review.partial.json","scope-redaction-new-controls.json","scope-redaction-preflight.exit.json","scope-redaction-preflight.log","scope-redaction-publication.json","scope-redaction-rollback.exit.json","scope-redaction-rollback.log","scope-redaction-rollback.py","sql-orchestration-poll-immutable-windows.js","sql-orchestration-poll-windows.js","sql-orchestration-publish-windows.js","sql-orchestration-save-immutable-v2-windows.js","sql-orchestration-save-windows.js","sql-report-renderer-finding-review.json","windows-node-helpers.js","windows-task-launcher.py"]
def digest(b):return hashlib.sha256(b).hexdigest()
def git(where,*args):return subprocess.check_output(['git',*args],cwd=where).decode().strip()
def put(name,raw):
 p=qa/name;assert p.resolve().is_relative_to(qa)
 p.parent.mkdir(parents=True,exist_ok=True)
 with p.open('xb') as f:f.write(raw)
def save(name,data):put(name,(json.dumps(data,indent=2)+'\n').encode())
assert git(repo,'rev-parse','HEAD')==SHA and git(repo,'rev-parse','HEAD^{tree}')=='ef2b61f832fb186593d974a3ff3fb9ba2160e245'
assert git(checkpoint,'rev-parse','HEAD')==CHECKPOINT
assert git(first_repo,'rev-parse','HEAD')=='fd5aa54f451c58e0a565030d8bb0800bf418bf92'
assert not qa.exists() and not package.exists()
qa.mkdir();package.mkdir()
bundle=checkpoint/'docs/evidence/20261018/pr169-sql-report-scope.json.gz.base64'
manifest=checkpoint/'docs/evidence/20261018/pr169-sql-report-scope.manifest.json'
assert bundle.stat().st_size==578141 and digest(bundle.read_bytes())=='0ad1863b030ee9d1a6cb117845caa05ed7dce26dbc19f311eabae97db0c9f711'
verify=[sys.executable,str(repo/'scripts/verify-evidence-bundle.py'),'--bundle',str(bundle),'--manifest',str(manifest)]
result=subprocess.run(verify,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);put('hosted-checkpoint-bounded-verification.log',result.stdout)
assert result.returncode==0,result.stdout.decode()
with gzip.GzipFile(fileobj=io.BytesIO(base64.b64decode(bundle.read_bytes().strip(),validate=True))) as stream:decoded=stream.read(16*1024*1024+1)
assert len(decoded)<=16*1024*1024
payload=json.loads(decoded);inventory=json.loads(manifest.read_bytes())
assert len(payload['members'])==len(ALLOWLIST)==78 and {x['path'] for x in payload['members']}==set(ALLOWLIST)
assert digest(decoded)==inventory['decoded_sha256']
expected={m['path']:m for m in inventory['members']}
for member in payload['members']:
 name=member['path'];assert name in ALLOWLIST and '..' not in PurePosixPath(name).parts
 raw=member['text'].encode();assert len(raw)==expected[name]['bytes'] and digest(raw)==expected[name]['sha256']
 put(name,raw)
print('CHECKPOINT_RECOVERY '+json.dumps({'members':78,'boundedVerifier':'PASS','allowlist':'exact published78names','helpersExecuted':False}),flush=True)
inputs=json.loads((control/'recovery-evidence/inputs.json').read_bytes())
assert inputs['candidate']==SHA and inputs['checkpoint']==CHECKPOINT
for entry in inputs['files']:
 name=entry['path'];assert re.fullmatch(r'[A-Za-z0-9_./-]{1,240}',name) and '..' not in PurePosixPath(name).parts
 raw=(control/'recovery-evidence/qa'/name).read_bytes()
 assert len(raw)==entry['bytes'] and digest(raw)==entry['sha256']
 put(name,raw)
put('hosted-recovery-inputs.json',(control/'recovery-evidence/inputs.json').read_bytes())
for source,target in [('recovery-evidence/recover.py','hosted-recovery-orchestrator.py'),('recovery-evidence/review.py','hosted-recovery-review.py'),('recovery-evidence/rollback.py','hosted-recovery-rollback.py'),('.github/workflows/sql-report-evidence-recovery.yml','hosted-recovery-workflow.yml')]:
 put(target,(control/source).read_bytes())
version=subprocess.check_output(['go','version']).decode().strip();assert version.startswith('go version go1.26.9 ')
provenance={'candidate':SHA,'candidateTree':'ef2b61f832fb186593d974a3ff3fb9ba2160e245','preview':PREVIEW,'checkpoint':CHECKPOINT,'recipe':git(control,'rev-parse','HEAD'),'recoveryRun':os.environ['GITHUB_RUN_ID'],'go':version,'localRuntime':'Disconnected after original Windows rollback/real native/source review; no raw local-only log fabricated. Earlier published scope evidence recovered unchanged; fresh hosted Linux rollback replaces archival dependency on inaccessible local raw log.','azure':'none','recoveredHelpersExecuted':False}
save('hosted-recovery-provenance.json',provenance)
review=subprocess.run([sys.executable,str(control/'recovery-evidence/review.py'),str(qa),str(repo),str(first_repo)],stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
sys.stdout.buffer.write(review.stdout);sys.stdout.flush()
assert review.returncode==0
started=datetime.datetime.now(datetime.timezone.utc).isoformat()
rollback=subprocess.run([sys.executable,str(control/'recovery-evidence/rollback.py'),str(repo),str(qa/'hosted-report-rollback-proof.json')],stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
put('hosted-report-rollback.log',rollback.stdout)
save('hosted-report-rollback.exit.json',{'exit_code':rollback.returncode,'startedAt':started,'completedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'candidate':SHA})
if rollback.returncode!=0:sys.stdout.buffer.write(rollback.stdout);sys.stdout.flush()
assert rollback.returncode==0
assert git(repo,'write-tree')=='ef2b61f832fb186593d974a3ff3fb9ba2160e245' and not git(repo,'status','--porcelain')
print('HOSTED_ROLLBACK '+json.dumps({'files':16,'retainedCorrection':True,'base':'724504b7e7fb454289e02c110abe0ff7fcfb65a4','exactRestoration':True,'staticTablesRegionBothWays':'PASS'}),flush=True)
paths=sorted(p for p in qa.rglob('*') if p.is_file())
assert 1<=len(paths)<=128 and not any(p.is_symlink() for p in paths)
members=[]
for p in paths:
 name=p.relative_to(qa).as_posix();assert re.fullmatch(r'[A-Za-z0-9_./-]{1,240}',name) and '..' not in PurePosixPath(name).parts
 text=p.read_bytes().decode('utf-8')
 assert not re.search(r'gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{30,}|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----',text),name
 members.append({'path':name,'text':text})
raw=json.dumps({'format':'cloud-assess-evidence-v1','members':members},ensure_ascii=False,separators=(',',':')).encode()
assert len(raw)<=16*1024*1024
encoded=base64.b64encode(gzip.compress(raw,mtime=0))+b'\n';assert len(encoded)<=4*1024*1024
mf={'format':'cloud-assess-evidence-v1','classification':'reviewed-synthetic-ci-local','bundle_sha256':digest(encoded),'decoded_sha256':digest(raw),'revisions':{'candidate':SHA,'preview':PREVIEW,'protected':'a6fee427e19feeef28691737ea9bc5026086efa0','checkpoint':CHECKPOINT,'recovery_recipe':provenance['recipe']},'members':[{'path':m['path'],'bytes':len(m['text'].encode()),'sha256':digest(m['text'].encode())} for m in members]}
stem='pr170-sql-report-implementation'
(package/(stem+'.json.gz.base64')).write_bytes(encoded)
(package/(stem+'.manifest.json')).write_bytes((json.dumps(mf,indent=2)+'\n').encode())
r=json.loads((qa/'39165239413/review.json').read_bytes());jobs=r['jobs']
dt=lambda s:datetime.datetime.fromisoformat(s.replace('Z','+00:00'))
workers=[j for j in jobs if j['name'] not in ['quality','windows-validation']]
first=min(dt(j['started_at']) for j in workers);last=max(dt(j['completed_at']) for j in workers)
summed=sum((dt(j['completed_at'])-dt(j['started_at'])).total_seconds() for j in workers)
rows='\n'.join('| '+j['name']+' | '+str(j['id'])+' | '+str(j['bytes'])+' | '+j['sha256']+' |' for j in jobs)
index=f'''# PR170 final internal SQL report evidence

Exact product code {SHA}, treeef2b61f832fb186593d974a3ff3fb9ba2160e245, sole parent verified scope correction724504b7e7fb454289e02c110abe0ff7fcfb65a4. Tested preview{PREVIEW} has identical tree and ordered protecteda6fee427/candidate8c47d810 parents. Draft/open/unmerged. Code16files and metadata/ref were fully read back before QA. Artifact/recovery/admin trees are not this product candidate and inherit no certification.

Own quality39165239413 passed all8workers/bothrequiredgates; all10complete original logs/steps/native runner labels and actual checkout reviewed. Linux551 exact control identities SHA256{r['controls']['linux']['sha256_full_lines_lf']}; Windows546 SHA256{r['controls']['windows']['sha256_full_lines_lf']}. All536/531corrected-base predecessors plus15frozen compiling report controls retained, no missing/extra/duplicates. Baseline455/450independent literal full-line hashes and all prior SQL/scope control script hashes checked. Coverage87.2%, unchanged75%floor. Full tests/race/vet/fuzz/default-custom packaging/branding/isolated installed CLI/docs/module-license/provenance/security passed. Zero reachable/imported vulnerabilities;1module-only advisory OPEN.

Own source39165239406:23captures/315contiguous unique chunks/909865UTF8bytes; all declared hashes and EVERY full file compared with exact candidate Git blobs. Source log{r['source_log']['bytes']}bytes/SHA256{r['source_log']['sha256']}. AZQR8e4f0577/tree17d93b20/APRL60eaddda unchanged. Original report first-runner39160871908 both native missing APIs and independent source39160871913 complete logs/metadata retained; source-only success is not report runtime coverage. All26unchanged report literals execute at final code, old316SQLcases and frozen oracle/budget/first runner preserved. No public registry/scoring/deployment eligibility/schema change.

Observed8worker first{first.isoformat()}, last{last.isoformat()}, span{int((last-first).total_seconds())}seconds; summed parallel worker duration{int(summed)}seconds is not elapsed/CPU/billed time or speedup. Full metadata retains queue/setup/step timestamps.

| Job | ID | UTF8bytes | SHA256 |
| --- | --- | --- | --- |
{rows}

## Recovery and fallback

Local filesystem/Node transport disconnected after original native/source reviews and Windows16file atomic rollback/restoration passed. Raw local-only report logs remain unavailable; none are fabricated or certified from remembered hashes. Previously published scope archive9ae98507 was bounded-verified BEFORE exact78name allowlisted restoration, preserving complete older QA/baseline/pre-fix evidence. No recovered helper was executed. A newly authored hosted read-only reviewer replayed all current quality/control/capture checks against candidate Git blobs. Fresh hosted Linux rollback removed all16report-specific files/changes to exact corrected baseaed4db22, retained all6independent scope production/test/control blobs, restored exact candidateef2b61f8 bytes/tree and passed static/tables/region suites BOTH ways. Full executed command/script/stdout/exit/proof are retained. Broader installed/operator/live/platform rollback remains open.

This explicitly strengthens the original raw-contract fallback: report rollback retains the independently verified privacy correction. Original contract/frozen expectations/all prior heads preserved. Recovery recipe{provenance['recipe']}, hosted run{provenance['recoveryRun']}, Go1.26.9 and exact existing checkout/setup action pins; no Azure/new permissions/action/dependency/SDK/source refresh. Recovery workflow is isolated to its QA branch and changes no product candidate workflow or protected check.

## Bounded retained material

{len(members)}members/{sum(m['bytes'] for m in mf['members'])}UTF8memberbytes/{len(encoded)}bundlebytes; SHA256{digest(encoded)}. Complete original quality/source/preparation logs plus reviewed scope checkpoint and fresh hosted rollback retained. [Bundle]({stem}.json.gz.base64), [manifest]({stem}.manifest.json). Existing scripts/verify-evidence-bundle.py passed on the completed output. It bounds input/decompression/inventory/names/hashes and never extracts or executes; perform it before explicit allowlisted restoration. Limited common-token screening is not privacy certification. UTF8connector-decoded logs preserved, not original transport ZIP hash claims.

Remote artifact publication/full readback and final administrative handover are recorded in live PR170, which owns exact-head closure. Code remains8c47d810; artifact/admin hashes are not code certification. Preserve all prior heads, frozen cases/primaries/captures/pins and unrelated work. No subagents/protected merge/bypass/release. Accepted adapters1/7/source7/7 unchanged; no full parity or independent-person review. Independent backup/retention, human/security/license/Gate004/release/live/scale/freshOS/installed rollback and module advisory remain OPEN.
'''
(package/'PR170_SQL_REPORT_IMPLEMENTATION_EVIDENCE.md').write_bytes(index.encode())
checked=subprocess.run([sys.executable,str(repo/'scripts/verify-evidence-bundle.py'),'--bundle',str(package/(stem+'.json.gz.base64')),'--manifest',str(package/(stem+'.manifest.json'))],stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
sys.stdout.buffer.write(checked.stdout);sys.stdout.flush();assert checked.returncode==0
print('RECOVERY_PACKAGE '+json.dumps({'candidate':SHA,'members':len(members),'member_bytes':sum(m['bytes'] for m in mf['members']),'bundle_bytes':len(encoded),'bundle_sha256':digest(encoded),'recipe':provenance['recipe'],'run':provenance['recoveryRun']},sort_keys=True),flush=True)
for p in sorted(package.iterdir()):
 data=p.read_bytes();text=data.decode();chunks=[text[i:i+8192] for i in range(0,len(text),8192)]
 print('EVIDENCE_FILE_META '+json.dumps({'name':p.name,'bytes':len(data),'sha256':digest(data),'parts':len(chunks)}),flush=True)
 for i,chunk in enumerate(chunks):print('EVIDENCE_FILE_CHUNK '+json.dumps({'name':p.name,'index':i,'count':len(chunks),'content':chunk},ensure_ascii=False,separators=(',',':')),flush=True)
