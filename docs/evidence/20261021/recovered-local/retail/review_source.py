from pathlib import Path
import collections,hashlib,json,re,subprocess,sys
ROOT=Path(__file__).resolve().parent
source=int(sys.argv[1]);sha=sys.argv[2];REPO=Path(sys.argv[3]).resolve()
def clean(raw):
 return [re.sub(r'^\S+Z ','',s) for s in raw.decode('utf-8-sig').splitlines()]
smeta=json.loads((ROOT/str(source)/'run.json').read_text());sjobs=json.loads((ROOT/str(source)/'jobs.json').read_text())['jobs']
assert smeta['head_sha']==sha and smeta['conclusion']=='success' and len(sjobs)==1 and sjobs[0]['conclusion']=='success'
assert all(s['conclusion']=='success' for s in sjobs[0]['steps'])
raw=(ROOT/str(source)/(str(sjobs[0]['id'])+'.log')).read_bytes();lines=clean(raw)
assert not any(s.startswith('##[error]') for s in lines)
provenance=[json.loads(s.split(' ',1)[1]) for s in lines if s.startswith('REGION_CAPTURE_PROVENANCE ')]
assert provenance==[{'source':'8e4f0577f3615e6c9014c031bcad079f235369cc','tree':'17d93b20c303f90f7843036be82f0dc32f3260f1','aprl':'60eaddda76541f6adbc1c5ffa686829807e55e29'}]
chunks=collections.defaultdict(list)
for s in lines:
 if s.startswith('REGION_CAPTURE_JSON '):
  v=json.loads(s.split(' ',1)[1]);chunks[v['name']].append(v)
assert len(chunks)==23
captures=[]
for name,parts in sorted(chunks.items()):
 assert [p['index'] for p in parts]==list(range(parts[0]['count'])) and all(p['count']==len(parts) for p in parts)
 b=''.join(p['content'] for p in parts).encode();h=hashlib.sha256(b).hexdigest()
 assert all(p['sha256']==h for p in parts)
 expected_bytes=subprocess.check_output(['git','show',f'{sha}:internal/plugins/region/testdata/{name}'],cwd=REPO)
 assert b==expected_bytes,(sha,name)
 captures.append({'name':name,'chunks':len(parts),'bytes':len(b),'sha256':h})
assert sum(c['chunks'] for c in captures)==315 and sum(c['bytes'] for c in captures)==909865

summary={'sha':sha,'source_run':source,'all_steps_success':True,'captures':captures,'chunks':315,'source_bytes':909865,'source_log':{'bytes':len(raw),'sha256':hashlib.sha256(raw).hexdigest()}}
(ROOT/str(source)/'review.json').write_text(json.dumps(summary,indent=2)+'\n')
print(json.dumps({'captures':len(captures),'chunks':315,'source_bytes':909865,'source_log':summary['source_log']}))
