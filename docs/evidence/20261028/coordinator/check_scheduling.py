"""Independent finite specification checker; imports no product code."""
import copy,json
from pathlib import Path
PHASES=('providers','required','history','retail','quota','reservations')
STOPS=('none','before-plan','after-inventory',*('after-'+p for p in PHASES),'before-publish')
FIELDS={'subscriptions','pairs','requirements_per_pair','cost','meters','quota_requests','reservation_requests','cancel_at','cause','failure_at','partial'}
def unique(pairs):
 d={}
 for k,v in pairs:
  if k in d:raise ValueError('duplicate JSON member')
  d[k]=v
 return d
def reject(message):raise ValueError(message)
def reference(q):
 if not isinstance(q,dict) or set(q)!=FIELDS:reject('input shape')
 for k in ['subscriptions','pairs','requirements_per_pair','meters','quota_requests','reservation_requests']:
  if type(q[k]) is not int or q[k]<0:reject('cardinality type')
 if type(q['cost']) is not bool or q['cancel_at'] not in STOPS or q['cause'] not in ['canceled','deadline_exceeded'] or q['failure_at'] not in ('none',*PHASES):reject('terminal shape')
 if not isinstance(q['partial'],list) or len(set(q['partial']))!=len(q['partial']) or any(p not in PHASES for p in q['partial']):reject('partial shape')
 counts=dict.fromkeys(PHASES,0)
 def terminal(status,publish=False,reason='none'):return {'calls':dict(counts),'status':status,'publish':publish,'reason':reason}
 if q['cancel_at']=='before-plan':return terminal('cancelled',reason=q['cause'])
 active=q['subscriptions']>0 and q['pairs']>0
 planned={'providers':q['subscriptions'] if active else 0,'required':q['pairs']*q['requirements_per_pair'] if active else 0,'history':q['subscriptions'] if active and q['cost'] else 0,'retail':1 if active and q['cost'] else 0,'quota':q['quota_requests'] if active else 0,'reservations':q['reservation_requests'] if active else 0}
 if q['subscriptions']>1000 or q['pairs']>8192 or sum(planned.values())>65536:return terminal('rejected',reason='plan_limit')
 if q['cancel_at']=='after-inventory':return terminal('cancelled',reason=q['cause'])
 for phase in PHASES:
  counts[phase]=planned[phase] if phase!='retail' or q['meters']>0 else 0
  if q['cancel_at']=='after-'+phase:return terminal('cancelled',reason=q['cause'])
  if q['failure_at']==phase:
   if counts[phase]==0:reject('failure requires an invoked dependency')
   return terminal('failed',reason='dependency_failure')
 if q['cancel_at']=='before-publish':return terminal('cancelled',reason=q['cause'])
 return terminal('completed_with_warnings' if any(counts[p]>0 for p in q['partial']) else 'completed',True)
def verify(doc):
 if set(doc)!={'format','scope','limits','cases'} or doc['format']!='cloud-assess-coordinator-scheduling-v1' or not isinstance(doc['scope'],str) or not doc['scope'] or doc['limits']!={'subscriptions':1000,'pairs':8192,'collector_callbacks':65536}:reject('document shape/limits')
 if not isinstance(doc['cases'],list) or not doc['cases']:reject('empty cases')
 ids=set()
 for c in doc['cases']:
  if set(c)!={'id','input','expected'} or not isinstance(c['id'],str) or not c['id'] or c['id'] in ids:reject('case identity')
  ids.add(c['id'])
  e=c['expected']
  if not isinstance(e,dict) or set(e)!={'calls','status','publish','reason'} or set(e['calls'])!=set(PHASES) or type(e['publish']) is not bool:reject('expected shape')
  if any(type(n) is not int or n<0 for n in e['calls'].values()):reject('expected cardinality')
  if reference(c['input'])!=e:reject('literal mismatch: '+c['id'])
 return len(ids)
def main():
 root=Path(__file__).parent;doc=json.loads((root/'scheduling-literals.json').read_text(),object_pairs_hook=unique,parse_constant=lambda s:reject('nonfinite JSON'))
 n=verify(doc);bad=[]
 x=copy.deepcopy(doc);x['cases'].append(x['cases'][0]);bad.append(x)
 x=copy.deepcopy(doc);x['cases'][0]['input']['pairs']=True;bad.append(x)
 x=copy.deepcopy(doc);x['cases'][0]['input']['cancel_at']='surprise';bad.append(x)
 x=copy.deepcopy(doc);x['cases'][0]['expected']['calls']['retail']=2;bad.append(x)
 x=copy.deepcopy(doc);x['limits']['collector_callbacks']+=1;bad.append(x)
 x=copy.deepcopy(doc);del x['cases'][0]['input']['cost'];bad.append(x)
 for x in bad:
  try:verify(x)
  except ValueError:pass
  else:raise AssertionError('malformed specification accepted')
 for raw in ['{"x":1,"x":2}','{"x":NaN}']:
  try:json.loads(raw,object_pairs_hook=unique,parse_constant=lambda s:reject('nonfinite JSON'))
  except ValueError:pass
  else:raise AssertionError('malformed JSON specification accepted')
 report={'literal_cases':n,'negative_controls':len(bad)+2,'actual_product_execution':False,'scope':'Finite boundary/cardinality/terminal reference specification only; not Go API enforcement, adapter completeness, real timer or memory proof.'}
 (root/'scheduling-verification.json').write_text(json.dumps(report,indent=2)+'\n');print(json.dumps(report))
if __name__=='__main__':main()
