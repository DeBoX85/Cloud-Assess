from pathlib import Path
import copy,json
out=Path(__file__).resolve().parent
A='aaaaaaaa-1111-4111-8111-111111111111';B='bbbbbbbb-2222-4222-8222-222222222222'
VM='microsoft.compute/virtualmachines';SS='microsoft.compute/virtualmachinescalesets';DISK='microsoft.compute/disks';ST='microsoft.storage/storageaccounts';COG='microsoft.cognitiveservices/accounts';DB='microsoft.sql/servers/databases';MI='microsoft.sql/managedinstances'
def rid(typ,name,sub=A):
 ns,*types=typ.split('/');parts=['','subscriptions',sub,'resourceGroups','rg','providers',ns]
 for t,n in zip(types,['server']*(len(types)-1)+[name]):parts.extend([t,n])
 return '/'.join(parts)
def resource(typ,name,sku,kind='',sub=A,**kw):
 r=dict(id=rid(typ,name,sub),subscriptionId=sub,resourceGroup='rg',name=name,type=typ,location='eastus',skuName=sku,skuTier='',skuFamily='',skuCapacity=0,kind=kind);r.update(kw);return r
def group(typ,ids,skus,kind='',kinds=None,flags=None,sub=A,source='eastus'):
 return dict(subscription_id=sub,source_region=source,resource_type=typ,kind=kind,skus=skus,kinds=kinds or {},resource_ids=ids,flags=flags or [])
def sku(status,values=None):return dict(status=status,values=values or {})
def val(state):return dict(state=state,blocked_zones=[])
def collection(typ,payload,slot,kind='',sub=A,source='eastus',target='westeurope'):
 return dict(subscription_id=sub,source_region=source,target_region=target,resource_type=typ,kind=kind,**{slot:payload})
def basepayload(typ,groups=None,sub=A,region='westeurope',**kw):
 p=dict(SubscriptionID=sub,Region=region,ResourceType=typ,Status='complete',Groups=groups or {},FailureCode='');p.update(kw);return p
v=resource(VM,'vm','Standard_D2s_v5');vgroup=group(VM,[v['id']],['standard_d2s_v5'])
vp=basepayload(VM,{'standard_d2s_v5':dict(Known=True,State=0,BlockedZones=[])})
vcol=collection(VM,vp,'vm')
cases=[]
def case(name,resources=None,groups=None,collections=None,skus=None,catalogs=None,plan='complete',admission='complete',scope='complete',error='',phase='',selected=None,declarations=None,request=None):
 c=dict(name=name,selected=selected if selected is not None else {A:'One'},resources=copy.deepcopy(resources if resources is not None else [v]),declarations=copy.deepcopy(declarations if declarations is not None else [dict(subscription_id=A,status='complete')]),request=copy.deepcopy(request if request is not None else dict(subscription_id=A,source_region='eastus',target_region='westeurope')),collections=copy.deepcopy(collections if collections is not None else [vcol]),want=dict(error=error,error_phase=phase,plan_status=plan,scope_status=scope,requirements=copy.deepcopy(groups if groups is not None else [vgroup]),admission_status=admission,sku_evidence=copy.deepcopy(skus if skus is not None else {VM:sku('complete',{'standard_d2s_v5':val(0)})}),catalogs=catalogs or [],preserve_raw_records=True))
 cases.append(c);return c
case('vm-complete-owned-coordinates',resources=[dict(v,skuTier='Standard',skuFamily='DSv5',skuCapacity=2,tags={'not-a-coordinate':'ignored'},sla='ignored')])
case('vm-no-collection-is-unknown',collections=[],admission='partial',skus={VM:sku('unknown')})
case('vm-complete-empty-is-unavailable',collections=[collection(VM,basepayload(VM),'vm')],skus={VM:sku('complete',{'standard_d2s_v5':val(2)})})
case('vm-partial-prefix-not-eligible',collections=[collection(VM,basepayload(VM,{'standard_d2s_v5':dict(Known=True,State=0,BlockedZones=[])},Status='partial',FailureCode='vm_sku_request_failed'),'vm')],admission='partial',skus={VM:sku('unknown')})
case('scope-missing-is-unknown',declarations=[],plan='partial',scope='unknown',admission='unknown',skus={VM:sku('unknown')})
case('scope-partial-retains-diagnostics',declarations=[dict(subscription_id=A,status='partial')],plan='partial',scope='partial',admission='partial',skus={VM:sku('unknown')})
case('other-scope-partial-does-not-poison-complete',selected={A:'One',B:'Two'},declarations=[dict(subscription_id=A,status='complete'),dict(subscription_id=B,status='partial')],plan='partial')
case('complete-empty-inventory',resources=[],groups=[],collections=[],skus={})
case('empty-selected-plan',resources=[],groups=[],collections=[],skus={},selected={},declarations=[],error='required_evidence_invalid',phase='admission',scope='',admission='')
case('missing-sku-does-not-vanish',resources=[dict(v,skuName='')],groups=[group(VM,[v['id']],[],flags=['missing_sku'])],collections=[],admission='partial',skus={VM:sku('unknown')})
case('unicode-sku-retained-unrepresentable',resources=[dict(v,skuName='Dø')],groups=[group(VM,[v['id']],[],flags=['unrepresentable_sku'])],collections=[],admission='partial',skus={VM:sku('unknown')})
case('sku-alias-raw-retained',resources=[dict(v,skuName=' Standard_D2s_v5 ')])
case('display-source-normalized-raw-retained',resources=[dict(v,location='East US')])
case('global-source-not-guessed',resources=[dict(v,location='global')],groups=[],collections=[],skus={})
for name,res,decl in [
 ('duplicate-resource-id',[v,dict(v)],None),
 ('case-alias-resource-id',[v,dict(v,id=v['id'][:-2]+'VM')],None),
 ('resource-unselected-subscription',[resource(VM,'other','Standard_D2s_v5',sub=B)],None),
 ('id-subscription-mismatch',[dict(v,id=rid(VM,'vm',B))],None),
 ('id-type-mismatch',[dict(v,type=DISK)],None),
 ('id-group-mismatch',[dict(v,resourceGroup='other')],None),
 ('id-empty-path-segment',[dict(v,id=v['id'].replace('/providers/','//providers/'))],None),
 ('negative-capacity',[dict(v,skuCapacity=-1)],None),
 ('oversized-retained-cell',[dict(v,skuTier='x'*513)],None),
 ('duplicate-completeness',None,[dict(subscription_id=A,status='complete'),dict(subscription_id=A.upper(),status='partial')]),
 ('extra-completeness',None,[dict(subscription_id=A,status='complete'),dict(subscription_id=B,status='complete')]),
 ('invalid-completeness',None,[dict(subscription_id=A,status='success')]),
]:case(name,resources=res,declarations=decl,error='required_inventory_invalid',phase='inventory',plan='',scope='',admission='',groups=[],skus={})
for name,cols,req in [
 ('wrong-target-evidence',[dict(vcol,target_region='westus')],None),
 ('wrong-source-evidence',[dict(vcol,source_region='westus')],None),
 ('wrong-selected-scope-evidence',[dict(vcol,subscription_id=B)],None),
 ('duplicate-evidence-binding',[vcol,copy.deepcopy(vcol)],None),
 ('inner-target-binding-mismatch',[collection(VM,basepayload(VM,region='westus'),'vm')],None),
 ('wrong-inner-resource-type',[collection(VM,basepayload(SS),'vm')],None),
 ('wrong-payload-family',[collection(VM,dict(SubscriptionID=A,Region='westeurope',Status='complete',Values={}), 'disk')],None),
 ('nil-union-payload',[dict(subscription_id=A,source_region='eastus',target_region='westeurope',resource_type=VM,kind='')],None),
 ('noncanonical-target',None,dict(subscription_id=A,source_region='eastus',target_region='West Europe')),
]:case(name,collections=cols,request=req,error='required_evidence_invalid',phase='admission',admission='',skus={})
s1=resource(ST,'s1','Standard_LRS','Storage');s2=resource(ST,'s2','Standard_LRS','StorageV2')
sg=group(ST,[s1['id'],s2['id']],['standard_lrs'],kinds={'standard_lrs':['storage','storagev2']})
sp=dict(SubscriptionID=A,Region='westeurope',Status='complete',Groups={'standard_lrs':{'storage':dict(Known=True,State=0),'storagev2':dict(Known=True,State=1)}},UnqualifiedNames={},FailureCode='')
case('storage-all-required-kinds-intersect',resources=[s1,s2],groups=[sg],collections=[collection(ST,sp,'storage')],skus={ST:sku('complete',{'standard_lrs':val(1)})})
case('storage-input-order-independent',resources=[s2,s1],groups=[sg],collections=[collection(ST,sp,'storage')],skus={ST:sku('complete',{'standard_lrs':val(1)})})
missing=resource(ST,'s2','Standard_LRS','')
case('storage-one-missing-kind-not-first-wins',resources=[s1,missing],groups=[group(ST,[s1['id'],missing['id']],['standard_lrs'],kinds={'standard_lrs':['','storage']},flags=['missing_kind'])],collections=[collection(ST,sp,'storage')],skus={ST:sku('unknown')},admission='partial')
c1=resource(COG,'c1','S0','OpenAI');c2=resource(COG,'c2','S0','SpeechServices')
cgroups=[group(COG,[c1['id']],['s0'],kind='openai',flags=['catalog_only']),group(COG,[c2['id']],['s0'],kind='speechservices',flags=['catalog_only'])]
cp1=basepayload(COG,{'s0':dict(Known=True,State=0,BlockedZones=[])},Kind='openai');cp2=basepayload(COG,{'s0':dict(Known=True,State=1,BlockedZones=[])},Kind='speechservices')
def catalog(typ,kind,status,keys=None,states=None):return dict(resource_type=typ,kind=kind,status=status,deployment_eligibility='not_assessed',keys=keys or [],states=states or {})
case('cognitive-same-sku-different-kinds',resources=[c1,c2],groups=cgroups,collections=[collection(COG,cp2,'cognitive','speechservices'),collection(COG,cp1,'cognitive','openai')],skus={COG:sku('unsupported')},catalogs=[catalog(COG,'openai','complete',['s0'],{'s0':0}),catalog(COG,'speechservices','complete',['s0'],{'s0':1})])
case('cognitive-wrong-kind-evidence',resources=[c1],groups=[cgroups[0]],collections=[collection(COG,cp2,'cognitive','speechservices')],error='required_evidence_invalid',phase='admission',admission='',skus={})
case('cognitive-missing-kind-not-guessed-from-tags',resources=[dict(c1,kind='',tags={'kind':'OpenAI'})],groups=[group(COG,[c1['id']],['s0'],flags=['catalog_only','missing_kind'])],collections=[],admission='partial',skus={COG:sku('unsupported')},catalogs=[catalog(COG,'','unknown')])
ss=resource(SS,'ss','Standard_D2s_v5');ssgroup=group(SS,[ss['id']],['standard_d2s_v5'],flags=['catalog_only'])
case('vmss-remains-catalog-only',resources=[ss],groups=[ssgroup],collections=[collection(SS,basepayload(SS,{'standard_d2s_v5':dict(Known=True,State=0,BlockedZones=[])}),'vm')],skus={SS:sku('unsupported')},catalogs=[catalog(SS,'','complete',['standard_d2s_v5'],{'standard_d2s_v5':0})])
for typ,name in [(DB,'sql-database-no-coordinate-inference'),(MI,'sql-mi-no-policy-inference')]:
 row=resource(typ,'sql','12.0/generalpurpose/name',skuTier='GeneralPurpose',skuFamily='Gen5',skuCapacity=4,tags={'edition':'generalpurpose','version':'12.0','instancePool':'enabled'})
 case(name,resources=[row],groups=[group(typ,[row['id']],['12.0/generalpurpose/name'],flags=['catalog_only'])],collections=[collection(typ,basepayload(typ),'sql')],skus={typ:sku('unsupported')},catalogs=[catalog(typ,'','complete')])
d=resource(DISK,'disk','Premium_LRS');dg=group(DISK,[d['id']],['premium_lrs'])
dp=dict(SubscriptionID=A,Region='westeurope',Status='complete',Values={'premium_lrs':dict(State=0,BlockedZones=[])},FailureCode='')
case('disk-validated-owned-evidence',resources=[d],groups=[dg],collections=[collection(DISK,dp,'disk')],skus={DISK:sku('complete',{'premium_lrs':val(0)})})
bad=copy.deepcopy(dp);bad['Values']['premium_lrs']['State']=99
case('disk-forged-complete-state-rejected',resources=[d],groups=[dg],collections=[collection(DISK,bad,'disk')],error='required_evidence_invalid',phase='admission',admission='',skus={})
unknown='microsoft.example/widgets';u=resource(unknown,'u','Opaque')
case('unknown-type-no-adapter-invention',resources=[u],groups=[group(unknown,[u['id']],['opaque'],flags=['unsupported'])],collections=[],skus={unknown:sku('unsupported')})
c3=resource(COG,'c3','S0','openai')
case('cognitive-kind-case-alias-stays-one-owned-group',resources=[c3,c1],groups=[group(COG,[c1['id'],c3['id']],['s0'],kind='openai',flags=['catalog_only'])],collections=[collection(COG,cp1,'cognitive','openai')],skus={COG:sku('unsupported')},catalogs=[catalog(COG,'openai','complete',['s0'],{'s0':0})])
case('invalid-second-family-rejects-prefix',resources=[v,d],groups=[dg,vgroup],collections=[vcol,collection(DISK,bad,'disk')],error='required_evidence_invalid',phase='admission',admission='',skus={})
case('complete-family-survives-other-family-unknown',resources=[v,s1],groups=[vgroup,group(ST,[s1['id']],['standard_lrs'],kinds={'standard_lrs':['storage']})],collections=[vcol],admission='partial',skus={VM:sku('complete',{'standard_d2s_v5':val(0)}),ST:sku('unknown')})
assert len({c['name'] for c in cases})==len(cases)
# Literal desired states above do not execute or inspect future target functions.
(out/'required-inventory-contract.provisional.json').write_text(json.dumps(dict(format='required-inventory-contract-v1',frozen_before_production=False,prepared_against='458d5c5296c321c8e5e961f18eb5422d43c2eb18',cases=cases),indent=2,ensure_ascii=False)+'\n')
print('Prepared',len(cases),'independent provisional scenarios, not frozen or executed')
