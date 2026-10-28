"""Closed independent expectation shape, not future collector execution."""
import copy,json
from pathlib import Path
def reject(message):raise ValueError(message)
def unique(pairs):
 out={}
 for k,v in pairs:
  if k in out:reject('duplicate JSON key')
  out[k]=v
 return out
def verify(d):
 if set(d)!={'format','base','scope','cases','budget_recipes','procedural'} or d['format']!='cloud-assess-subscription-reservations-v1' or d['base']!='4627a4d9c2853d9337e8c0f959f241a6d319b416':reject('document identity')
 if len(d['cases'])!=13 or len(d['budget_recipes'])!=8 or len(d['procedural'])!=12 or len(set(d['procedural']))!=12:reject('lost expectations')
 ids=set()
 for c in d['cases']:
  if set(c)!={'id','scenario','status','code','calls','regions','rows'} or c['id'] in ids:reject('case shape/identity')
  ids.add(c['id'])
  if c['status'] not in ['complete','partial','unknown'] or type(c['calls']) is not int or not 1<=c['calls']<=256:reject('case status/calls')
  if not isinstance(c['regions'],list) or len(c['regions'])!=len(set(c['regions'])) or any(not isinstance(r,str) or not r for r in c['regions']):reject('region shape')
  if not isinstance(c['rows'],list):reject('rows shape')
  for row in c['rows']:
   if not isinstance(row,list) or len(row)!=6 or any(not isinstance(s,str) for s in row[:4]) or any(type(n) is not int or n<0 for n in row[4:]) or row[0] not in c['regions']:reject('literal row')
  if c['status']=='complete' and c['code']!='' or c['status']!='complete' and not c['code']:reject('status/code')
 for c in d['budget_recipes']:
  if set(c)!={'id','kind','amount','calls','status','code','rows','allocated'} or c['id'] in ids:reject('recipe shape/identity')
  ids.add(c['id'])
  if c['kind'] not in ['calls','pages','bytes','work'] or any(type(c[k]) is not int or c[k]<0 for k in ['amount','calls','rows','allocated']) or not 1<=c['calls']<=256:reject('recipe cardinality')
  if c['status'] not in ['complete','partial'] or c['status']=='complete' and c['code'] or c['status']=='partial' and not c['code']:reject('recipe status')
 return {'literal_cases':13,'boundary_recipes':8,'procedural_obligations':12,'actual_product_execution':False,'scope':'Independent expectation shape only, not collector/API/timer/memory proof.'}
def main():
 root=Path(__file__).parent;d=json.loads((root/'literals.json').read_text(),object_pairs_hook=unique,parse_constant=lambda s:reject('nonfinite JSON'))
 report=verify(d);bad=[]
 x=copy.deepcopy(d);x['cases'][1]['id']=x['cases'][0]['id'];bad.append(x)
 x=copy.deepcopy(d);x['cases'][0]['calls']=True;bad.append(x)
 x=copy.deepcopy(d);x['cases'][0]['status']='healthy';bad.append(x)
 x=copy.deepcopy(d);x['cases'][0]['rows'][0][4]=True;bad.append(x)
 x=copy.deepcopy(d);del x['cases'][0]['regions'];bad.append(x)
 x=copy.deepcopy(d);x['budget_recipes'][1]['calls']=257;bad.append(x)
 for x in bad:
  try:verify(x)
  except ValueError:pass
  else:raise AssertionError('bad independent expectation admitted')
 for raw in ['{"x":1,"x":2}','{"x":NaN}']:
  try:json.loads(raw,object_pairs_hook=unique,parse_constant=lambda s:reject('nonfinite JSON'))
  except ValueError:pass
  else:raise AssertionError('bad JSON admitted')
 report['malformed_controls']=8
 (root/'literal-shape-verification.json').write_text(json.dumps(report,indent=2)+'\n');print(json.dumps(report))
if __name__=='__main__':main()
