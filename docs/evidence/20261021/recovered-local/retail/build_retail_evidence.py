from pathlib import Path
import base64,gzip,hashlib,json,subprocess
root=Path(__file__).resolve().parent;repo=root.parent.parent/'retail-validation-20261020';out=root/'publish';out.mkdir(exist_ok=True)
sha='458d5c5296c321c8e5e961f18eb5422d43c2eb18';preview='dd0b28416a493fba17fd7de9809cebb23d1c4725'
def h(b):return hashlib.sha256(b).hexdigest()
def git(*a):return subprocess.check_output(['git',*a],cwd=repo)
q=json.loads((root/'39450901879/review.json').read_bytes());assert q['candidate']==sha and len(q['jobs'])==10
assert q['controls']['linux']['count']==622 and q['controls']['windows']['count']==617
assert q['source']['sha']==sha and q['source']['source_bytes']==909865
assert json.loads((root/'depth-local-git-review.json').read_bytes())['candidate']==sha
roll=json.loads((root/'rollback-review.json').read_bytes());assert len(roll['files'])==25 and roll['exact_bytes_restored']
assert 'all42 compiling reversals' in (root/'depth-controls-local.log').read_text()
assert (root/'depth-full-preflight.log').read_text().rstrip().endswith('Developer preflight passed. Full final-revision native/source gates remain required.')
ci={};proof={}
for rid in [39446889947,39446889918,39449387918,39449387761,39450901879,39450901886]:
 folder=root/str(rid);assert json.loads((folder/'run.json').read_bytes())['conclusion']=='success'
 for p in sorted(folder.iterdir()):
  if p.is_file():ci['runs/'+str(rid)+'/'+p.name]=p.read_bytes()
# Preserve the complete predecessor control logs, not just a count or digest.
base=root.parent/'39182715043';jobs=json.loads((base/'jobs.json').read_bytes());jobs=jobs['jobs'] if isinstance(jobs,dict) else jobs
for j in jobs:
 if j['name'] in {'quality','windows-validation'}:continue
 key='id' if 'id' in j else 'job';p=base/(str(j[key])+'.log');ci['baseline175/'+p.name]=p.read_bytes()
for name in ['jobs.json','review.json']:
 p=base/name
 if p.exists():ci['baseline175/'+name]=p.read_bytes()
# First consumer's actual absent-API compilation is preparation, not a product pass.
first=root.parent/'39350624037'
for p in sorted(first.iterdir()):
 if p.is_file() and (p.suffix=='.log' or p.name in {'run.json','jobs.json','missing-api-review.json'}):ci['first-consumer178/'+p.name]=p.read_bytes()
for p in sorted(root.rglob('*')):
 if not p.is_file() or '__pycache__' in p.parts or 'publish' in p.parts or 'remote' in p.parts or p.name=='archive-summary.json':continue
 rel=p.relative_to(root)
 if rel.parts[0].isdigit():continue
 proof['implementation/'+rel.as_posix()]=p.read_bytes()
for p in ['retail-boundary.prepared.go.txt','retail-boundary-preproduction.json']:
 proof['preproduction/'+p]=(root.parent/p).read_bytes()
paths=['go.mod','go.sum','docs/REGION_RETAIL_CONTRACT_20261020.md','docs/REGION_RETAIL_IMPLEMENTATION_20261020.md','docs/CURRENT_STATE.md','docs/RESTART_PROMPT.md','scripts/tests/retail-controls.py','scripts/tests/retail-controls.plan.json']
paths+=git('ls-tree','-r','--name-only',sha,'internal/plugins/region').decode().splitlines()
paths=[p for p in paths if not p.startswith('internal/plugins/region/') or '/retail' in p]
for p in paths:proof['candidate/'+p]=git('show',sha+':'+p)
for p in ['internal/plugins/region/retail_decode.go','internal/plugins/region/retail_pricing.go']:
 proof['preserved181/'+p]=git('show','44114c6357caed81d5153b69be5a231e7eabcbec:'+p)
proof['preserved182/internal/plugins/region/retail_decode.go']=git('show','aa304cf71282d28d1586a72f0ce232ce2c1b2d3e:internal/plugins/region/retail_decode.go')
revisions=dict(candidate=sha,preview=preview,preserved182='aa304cf71282d28d1586a72f0ce232ce2c1b2d3e',wip181='44114c6357caed81d5153b69be5a231e7eabcbec',verified_security180='9fe1261260f55d9796917686a4f0bc10d1cd9490',first_consumer178='9e15db3df8cd34c35842394e26acaee0ac5d01ef',protected='a6fee427e19feeef28691737ea9bc5026086efa0')
summary={}
for label,members in [('retail-ci',ci),('retail-proof',proof)]:
 assert len(members)<=128,(label,len(members))
 assert all(b.decode('utf8').encode('utf8')==b for b in members.values())
 payload={'format':'cloud-assess-evidence-v1','members':[dict(path=p,text=b.decode('utf8')) for p,b in sorted(members.items())]}
 raw=json.dumps(payload,ensure_ascii=True,separators=(',',':')).encode();assert len(raw)<=16<<20,(label,len(raw))
 bundle=base64.b64encode(gzip.compress(raw,mtime=0))+b'\n';assert len(bundle)<=4<<20
 manifest=dict(format='cloud-assess-evidence-v1',classification='reviewed-synthetic-ci-local',bundle_sha256=h(bundle),decoded_sha256=h(raw),members=[dict(path=p,bytes=len(b),sha256=h(b)) for p,b in sorted(members.items())],revisions=revisions)
 (out/(label+'.json.gz.base64')).write_bytes(bundle);(out/(label+'.manifest.json')).write_text(json.dumps(manifest,indent=2)+'\n')
 summary[label]=dict(members=len(members),member_bytes=sum(len(b) for b in members.values()),decoded_payload_bytes=len(raw),bundle_bytes=len(bundle),bundle_sha256=h(bundle))
(out/'.gitattributes').write_text('* -text\n')
(root/'archive-summary.json').write_text(json.dumps(summary,indent=2)+'\n');print(json.dumps(summary))
