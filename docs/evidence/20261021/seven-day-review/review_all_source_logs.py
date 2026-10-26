from pathlib import Path
import ast
import collections
import hashlib
import json
import re
import subprocess

ROOT = Path('/workspace/scratch/81f11bb42cee')
REPO = ROOT / 'audit-code'
RUNS = {r['id']: r for r in json.loads((ROOT/'audit-runs-initial.json').read_text())}
JOBS = {r['run_id']: r['jobs'] for r in json.loads((ROOT/'audit-jobs.json').read_text())}
LOGS = {r['job_id']: r for r in json.loads((ROOT/'audit-log-analysis.json').read_text())}
PROVENANCE = {'source':'8e4f0577f3615e6c9014c031bcad079f235369cc','tree':'17d93b20c303f90f7843036be82f0dc32f3260f1','aprl':'60eaddda76541f6adbc1c5ffa686829807e55e29'}

def unique(pairs):
    out = {}
    for k,v in pairs:
        if k in out:
            raise ValueError('duplicate JSON key')
        out[k] = v
    return out

def parse(s):
    return json.loads(s, object_pairs_hook=unique)

def verify(raw, sha):
    lines = [re.sub(r'^\S+Z ', '', s) for s in raw.decode('utf-8-sig').splitlines()]
    if any(s.startswith('##[error]') for s in lines):
        raise ValueError('successful source log contains error marker')
    if not any('Cleaning up orphan processes' in s for s in lines[-15:]):
        raise ValueError('missing source footer')
    provenance = [parse(s.split(' ',1)[1]) for s in lines if s.startswith('REGION_CAPTURE_PROVENANCE ')]
    if provenance != [PROVENANCE]:
        raise ValueError('source provenance mismatch')
    chunks = collections.defaultdict(list)
    for s in lines:
        if s.startswith('REGION_CAPTURE_JSON '):
            v = parse(s.split(' ',1)[1])
            chunks[v['name']].append(v)
    capture_script = subprocess.check_output(['git','show',sha+':scripts/capture-region-source.py'],cwd=REPO)
    nodes = [n for n in ast.parse(capture_script).body if isinstance(n, ast.Assign) and any(isinstance(t, ast.Name) and t.id == 'FILES' for t in n.targets)]
    if len(nodes) != 1:
        raise ValueError('ambiguous exact-head capture manifest')
    expected_names = ast.literal_eval(nodes[0].value)
    if len(expected_names) != len(set(expected_names)) or len(expected_names) != 23 or set(chunks) != set(expected_names):
        raise ValueError('capture inventory mismatch')
    captures = []
    for name, parts in sorted(chunks.items()):
        if not re.fullmatch(r'source-[a-z-]+\.json', name):
            raise ValueError('unsafe source member')
        if [p['index'] for p in parts] != list(range(parts[0]['count'])) or any(p['count'] != len(parts) for p in parts):
            raise ValueError('noncontiguous or inconsistent source chunks')
        raw_member = ''.join(p['content'] for p in parts).encode()
        digest = hashlib.sha256(raw_member).hexdigest()
        if any(p['sha256'] != digest for p in parts):
            raise ValueError('capture hash mismatch')
        expected = subprocess.check_output(['git','show',sha+':internal/plugins/region/testdata/'+name],cwd=REPO)
        if raw_member != expected:
            raise ValueError('capture differs from exact candidate Git')
        captures.append({'name':name,'chunks':len(parts),'bytes':len(raw_member),'sha256':digest})
    if sum(v['chunks'] for v in captures) != 315 or sum(v['bytes'] for v in captures) != 909865:
        raise ValueError('capture totals mismatch')
    return captures

out = []
for run in RUNS.values():
    if run['name'] != 'Pinned region source characterization':
        continue
    row = {'run':run['id'],'sha':run['head_sha'],'observed_conclusion':run['conclusion']}
    try:
        jobs = JOBS[run['id']]
        if run['conclusion'] != 'success':
            row['disposition'] = 'not certified; preserve original cancelled/failed status'
        else:
            if len(jobs) != 1 or jobs[0]['conclusion'] != 'success' or any(s['conclusion'] != 'success' for s in jobs[0]['steps']):
                raise ValueError('source success has incomplete job/step inventory')
            meta = LOGS[jobs[0]['id']]
            raw = Path(meta['log']['path']).read_bytes()
            row['captures'] = verify(raw,run['head_sha'])
            row['log'] = meta['log']
            row['disposition'] = 'full source reconstruction matches exact candidate Git'
    except Exception as error:
        row['error'] = str(error)
        row['disposition'] = 'UNRESOLVED'
    out.append(row)
(ROOT/'audit-source-review.json').write_text(json.dumps(out,indent=2)+'\n')
print(json.dumps({'source_runs':len(out),'reconstructed':sum('captures' in r for r in out),'unresolved':[{k:r[k] for k in ('run','error')} for r in out if 'error' in r]}))
