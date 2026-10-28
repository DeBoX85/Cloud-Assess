from pathlib import Path
import hashlib,json,os,subprocess
root=Path('/workspace/scratch/ee3c2680ffe8/retail-validation-20261020');out=Path(__file__).resolve().parent;base='9fe1261260f55d9796917686a4f0bc10d1cd9490'
def git(*a):return subprocess.check_output(['git',*a],cwd=root)
def owned_names():return sorted(n for n in set(git('diff','--name-only',base).decode().splitlines()+git('ls-files','--others','--exclude-standard').decode().splitlines()) if '__pycache__' not in Path(n).parts)
names=owned_names()
assert all(not n.startswith('internal/rules/upstream/') for n in names)
original={n:(root/n).read_bytes() if (root/n).exists() else None for n in names}
before={}
for n in names:
 p=subprocess.run(['git','show',base+':'+n],cwd=root,stdout=subprocess.PIPE,stderr=subprocess.DEVNULL);before[n]=p.stdout if p.returncode==0 else None
assert len(names)>=20
cmd=['/workspace/scratch/ee3c2680ffe8/tooling/go/bin/go','test','./internal/plugins/region','-count=1']
env=dict(os.environ,GOFLAGS='-mod=readonly',GOTOOLCHAIN='local')
def put(files):
 for n,b in files.items():
  p=root/n
  if b is None:p.unlink(missing_ok=True)
  else:p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(b)
def test(label,args):
 p=subprocess.run(args,cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);(out/(label+'.log')).write_bytes(p.stdout);assert p.returncode==0,(label,p.stdout.decode())
try:
 put(before)
 for n,b in before.items():assert ((root/n).read_bytes() if (root/n).exists() else None)==b
 assert not git('diff',base,'--','internal/plugins/region','scripts/ci/quality-plan.json','.github/workflows/test.yml')
 test('rollback-base-region',cmd)
finally:put(original)
for n,b in original.items():assert ((root/n).read_bytes() if (root/n).exists() else None)==b
assert names==owned_names()
test('restored-retail',cmd+['-run','^TestRetail'])
summary=dict(base=base,files=[dict(path=n,base_sha256=hashlib.sha256(before[n]).hexdigest() if before[n] is not None else None,restored_sha256=hashlib.sha256(b).hexdigest() if b is not None else None) for n,b in original.items()],baseline_region_test='pass',restored_retail_test='pass',exact_bytes_restored=True,scope='local Linux coherent source/runner/workflow/docs rollback; not installed/live rollback')
(out/'rollback-review.json').write_text(json.dumps(summary,indent=2)+'\n');print(json.dumps(dict(files=len(names),exact_bytes_restored=True,baseline_region_test='pass',restored_retail_test='pass')))
