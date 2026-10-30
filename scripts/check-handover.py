#!/usr/bin/env python3
"""Check local recovery orientation. Never fetch, publish, or certify QA."""
import argparse
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError('duplicate JSON key: ' + key)
        result[key] = value
    return result


def read_json(path):
    raw = path.read_bytes()
    if len(raw) > 1024 * 1024:
        raise ValueError('orientation record exceeds 1 MiB')
    return json.loads(raw.decode('utf-8'), object_pairs_hook=unique_object)


def require(condition, message):
    if not condition:
        raise ValueError(message)


def validate(state, inventory, followups):
    require(state['schema'] == 'cloud-assess-recovery-state-v1', 'unsupported state schema')
    roles = state['roles']
    for name in ['accepted', 'verified_product', 'preserved_wip', 'preserved_preparation', 'security_review']:
        require(re.fullmatch('[0-9a-f]{40}', roles[name]['commit']) is not None,
                'invalid role commit: ' + name)
    product = roles['verified_product']
    require(len({roles[n]['commit'] for n in ['accepted', 'verified_product', 'preserved_wip', 'preserved_preparation', 'security_review']}) == 5,
            'distinct recovery revision roles conflated')
    require(roles['handover_proposal']['parent'] == product['commit'], 'handover parent differs from verified product')
    require(roles['handover_proposal']['branch'] != product['branch'], 'handover and verified product branch conflated')
    require(product['merged'] is False, 'snapshot cannot silently claim protected acceptance')
    qa = state['product_qa']
    require(qa['candidate'] == product['commit'], 'QA revision differs from verified product')
    require(qa['conclusions'] == ['success', 'success'] and qa['complete_logs_reviewed'] == 11,
            'completed product QA requires both workflows and eleven reviewed complete logs')
    require(type(qa['quality_run']) is int and type(qa['source_run']) is int and
            qa['quality_run'] > 0 and qa['source_run'] > 0 and qa['quality_run'] != qa['source_run'],
            'missing or conflated workflow identities')
    require(0 <= qa['coverage_floor'] <= qa['coverage_percent'] <= 100, 'coverage below declared floor or invalid')
    evidence = state['evidence']
    require(evidence['remote_readback'] is True and evidence['fresh_restore'] is True and
            evidence['executed_members'] == 0, 'retained evidence requires readback and non-executing fresh restore')
    for key in ['capture_commit', 'receipt_commit']:
        require(re.fullmatch('[0-9a-f]{40}', evidence[key]) is not None, 'invalid evidence commit')
    require(evidence['members'] > 0 and evidence['bytes'] > 0, 'empty retained evidence')
    require(evidence['original_streams_recovered'] is False and evidence['lost_original_streams'] == 134 and
            evidence['replacement_provenance'].strip(), 'lost originals/replacement provenance changed without reconciliation')
    require(evidence['independent_backup_verified'] is False, 'independent backup is still unverified in this recovery snapshot')
    ids = [row['id'] for row in inventory['obligations']]
    legacy = [row['id'] for row in inventory['legacy_ids']]
    require(len(ids) == len(set(ids)) == inventory['obligation_count'] == state['work']['obligations'],
            'current obligations missing or duplicated')
    require(len(legacy) == len(set(legacy)) == inventory['legacy_count'] == state['work']['legacy_ids'],
            'legacy identities missing or duplicated')
    task = next(row for row in inventory['obligations'] if row['id'] == state['work']['active_task'])
    require(task['state'] == state['work']['state'] == 'implementation_in_progress', 'coordinator task silently closed')
    require(all(next(row for row in inventory['obligations'] if row['id'] == 'OFFLINE-' + str(n))['state'] == 'not_started'
                for n in range(11, 16)), 'future task status differs from recovery snapshot')
    require(set(state['work']['followups']) == {row['id'] for row in followups['items']}, 'recovered followups lost')
    caution = state['enhanced_verification']
    require(caution['active'] is True and caution['minimum_completed_material_slices'] >= 3 and
            0 <= caution['completed_since_handover'] < caution['minimum_completed_material_slices'] and
            caution['exit'].strip(), 'enhanced verification prematurely disabled or untracked')
    require(len(state['uncertainty']) >= 4 and state['work']['next_action'].strip(), 'uncertainty or next action missing')


def git(*args):
    return subprocess.check_output(['git', '-C', str(ROOT), *args], stderr=subprocess.PIPE,
                                   timeout=20).decode('utf-8').strip()


def check_checkout(state, inventory):
    product = state['roles']['verified_product']
    require(git('show', '-s', '--format=%T', product['commit']) == product['tree'], 'verified product Git tree differs')
    git('merge-base', '--is-ancestor', product['commit'], 'HEAD')
    # Read tracked pins from both the verified object and current working files.
    for name in ['go.mod', 'go.sum', '.gitmodules']:
        raw = subprocess.check_output(['git', '-C', str(ROOT), 'show', product['commit'] + ':' + name], timeout=20)
        require((ROOT / name).read_bytes() == raw, 'tracked pin bytes changed: ' + name)
    require('go ' + state['pins']['go'] in (ROOT / 'go.mod').read_text(encoding='utf-8'), 'Go pin differs')
    reference = (ROOT / 'docs/REFERENCE_POLICY.md').read_text(encoding='utf-8')
    require(state['pins']['azqr'] in reference and state['pins']['aprl'] in reference, 'reference pins differ')
    links = git('ls-tree', '-r', 'HEAD')
    require('160000 commit ' + state['pins']['aprl'] in links, 'APRL gitlink differs')
    frozen = json.loads(git('show', product['commit'] + ':docs/TASK_CARRY_FORWARD_20261022.json'))
    for key in ['legacy_ids', 'sources', 'legacy_count', 'obligation_count']:
        require(inventory[key] == frozen[key], 'frozen obligation input changed: ' + key)
    require([x['id'] for x in inventory['obligations']] == [x['id'] for x in frozen['obligations']],
            'ordered obligation identities changed')
    for path in state['authority']:
        require((ROOT / path).is_file(), 'authority file missing: ' + path)
    for path in ['CONTINUE.md', 'CURRENT_STATE.md', 'RESTART_PROMPT.md']:
        text = (ROOT / 'docs' / path).read_text(encoding='utf-8')
        require('RECOVERY_STATE.json' in text and 'RECOVERY_HANDOVER_20261030.md' in text,
                'current recovery routing missing: ' + path)
        require('DC-CI-EVIDENCE-20261022' in text[:6000] and
                'https://github.com/DeBoX85/Cloud-Assess/issues/189' in text[:6000],
                'evidence duty missing from compact entrypoint: ' + path)
    # Check new local Markdown links; historical ledgers retain their own dated paths.
    for path in ['docs/CONTINUE.md', 'docs/CURRENT_STATE.md', 'docs/RESTART_PROMPT.md',
                 'docs/RECOVERY_HANDOVER_20261030.md', 'docs/HANDOVER_PROCESS_REVIEW_20261030.md']:
        page = ROOT / path
        for target in re.findall(r'\]\(([^)]+)\)', page.read_text(encoding='utf-8')):
            if '://' not in target and not target.startswith('#'):
                require((page.parent / target.split('#', 1)[0]).is_file(), 'broken local orientation link: ' + target)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.parse_args()
    try:
        state = read_json(ROOT / 'docs/RECOVERY_STATE.json')
        inventory = read_json(ROOT / 'docs/TASK_CARRY_FORWARD_20261022.json')
        validate(state, inventory, read_json(ROOT / 'docs/RECOVERY_FOLLOWUPS_20261030.json'))
        check_checkout(state, inventory)
    except (ValueError, KeyError, TypeError, StopIteration, OSError, subprocess.SubprocessError) as error:
        print('Local handover consistency FAILED: ' + str(error), file=sys.stderr)
        return 1
    print(json.dumps({'local_handover_consistency': 'pass', 'verified_product': state['product_qa']['candidate'],
                      'remote_state_checked': False, 'semantic_accuracy_certified': False,
                      'new_head_QA_inherited': False, 'uncertainties_remain': state['uncertainty']}))
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
