"""Specification inventory only; never executes or models the future producer."""
from pathlib import Path
import json,copy,re
D=Path(__file__).resolve().parent
def validate(p):
 assert isinstance(p,dict) and p['format']=='cloud-assess-run-history-binding-plan-v1'
 for k,n in [('literal_cases',35),('known_failure_codes',5),('request_rejections',13),('budget_recipes',6),('cancellation_kinds',2),('relative_integration_literals',2),('expected_named_events',77)]:assert p[k]==n
 assert len(p['case_names'])==len(set(p['case_names']))==35
 assert len(p['root_tests'])==len(set(p['root_tests']))==14
 assert all(isinstance(s,str) and re.fullmatch('[a-z][a-z0-9-]*',s) for s in p['case_names'])
 assert all(isinstance(s,str) and s.startswith('TestRunHistory') for s in p['root_tests'])
plan=json.loads((D/'test-plan.json').read_text());validate(plan)
mutations=[lambda p:p['case_names'].pop(),lambda p:p['case_names'].__setitem__(0,p['case_names'][1]),lambda p:p['root_tests'].pop(),lambda p:p.__setitem__('budget_recipes',5),lambda p:p.__setitem__('expected_named_events',76),lambda p:p.__setitem__('cancellation_kinds',1),lambda p:p.__setitem__('format','unknown'),lambda p:p.__setitem__('known_failure_codes',None)]
for mutate in mutations:
 p=copy.deepcopy(plan);mutate(p)
 try:validate(p)
 except (AssertionError,TypeError,KeyError):continue
 raise AssertionError('malformed plan admitted')
print(json.dumps({'literal_cases':35,'root_tests':14,'future_named_events':77,'negative_controls':8,'actual_future_API_executed':False}))
