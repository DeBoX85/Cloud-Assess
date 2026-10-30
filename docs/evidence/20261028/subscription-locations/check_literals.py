"""Independent shape/identity controls; never generate or update expected values."""
import copy,json,sys
from pathlib import Path
D=Path(__file__).resolve().parent

def validate(d):
 if d.get('format')!='cloud-assess-subscription-locations-literals-v1':raise ValueError('format')
 cs=d.get('cases');recipes=d.get('boundary_recipes');ob=d.get('procedural')
 if not isinstance(cs,list) or len(cs)!=35:raise ValueError('case count')
 if len({c['name'] for c in cs})!=35 or any(not c['name'] for c in cs):raise ValueError('case IDs')
 if not isinstance(recipes,list):raise ValueError('recipe identities')
 if len(recipes)!=8 or len(set(recipes))!=8:raise ValueError('recipe count')
 if len(ob)!=8 or len(set(ob))!=8:raise ValueError('obligations')
 for c in cs:
  if c['status'] not in ['complete','partial','unknown'] or c['count_source'] not in ['api','fallback','incomplete']:raise ValueError('status')
  if c['calls']!=len(c['pages']) or c['calls']<1:raise ValueError('calls')
  if c['regions']!=sorted(set(c['regions'])) or set(c['counts'])!=set(c['regions']) or not set(c['mappings'])<=set(c['regions']):raise ValueError('coordinates')
  if c['status']!='complete' and c['count_source']!='incomplete':raise ValueError('incomplete promotion')
  if c['status']=='complete' and c['code']:raise ValueError('success code')
  if c['count_source']=='fallback' and c['mappings']:raise ValueError('invented fallback mapping')
  if any(not isinstance(n,int) or n<0 for n in c['counts'].values()):raise ValueError('counts')
 return d

def main():
 d=json.loads((D/'literals.json').read_text());validate(d)
 mutations=[lambda x:x['cases'].pop(),lambda x:x['cases'].__setitem__(1,copy.deepcopy(x['cases'][0])),lambda x:x['boundary_recipes'].pop(),lambda x:x['boundary_recipes'].__setitem__(1,x['boundary_recipes'][0]),lambda x:x['procedural'].pop(),lambda x:x['cases'][0].__setitem__('status','invented'),lambda x:x['cases'][7].__setitem__('count_source','fallback'),lambda x:x['cases'][0].__setitem__('calls',99)]
 for mutate in mutations:
  x=copy.deepcopy(d);mutate(x)
  try:validate(x)
  except (ValueError,KeyError,TypeError):pass
  else:raise AssertionError('malformed specification accepted')
 z=json.loads((D/'zone-counts-pinned.json').read_text());assert z['source']=='8e4f0577f3615e6c9014c031bcad079f235369cc' and len(z['counts'])==33 and set(z['counts'].values())=={3}
 print(json.dumps({'wire_literals':35,'boundary_recipes':8,'procedural_obligations':8,'malformed_controls':8,'actual_future_API_executed':False}))
if __name__=='__main__':main()
