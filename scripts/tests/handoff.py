#!/usr/bin/env python3
"""Keep resume authority explicit without treating historical snapshots as current."""
from pathlib import Path
import json
import re
import unittest

ROOT = Path(__file__).resolve().parents[2]
ENTRYPOINTS = (
    'SESSION_HANDOVER.md', 'HANDOVER_CHECKPOINT_20261016.md',
    'NEW_CHAT_PROMPT_20261016.md', 'DEVELOPMENT_EXECUTION_PLAN.md',
    'QA_PROCESS.md', 'ROADMAP.md', 'TARGET_SPECIFICATION.md',
    'IMPLEMENTATION_PLAN.md', 'ALIGNMENT_REVIEW.md',
    'FEATURE_PARITY_CHARACTERIZATION.md', 'REGION_SELECTION.md',
    'REGION_SKU_TARGET_ANALYSIS_20261013.md', 'REGION_COMPUTE_SKU_SOURCE_20261013.md',
    'REGION_OTHER_SKU_SOURCE_20261013.md', 'REGION_STORAGE_SKU_COLLECTOR_20261016.md',
    'REGION_VM_SKU_COLLECTOR_20261016.md', 'SYSTEMATIC_REVIEW_FINAL_ASSESSMENT_20261016.md',
    'CODE_REVIEW_QA_20261014.md',
)

EVIDENCE_TASK = 'DC-CI-EVIDENCE-20261022'
EVIDENCE_ISSUE = 'https://github.com/DeBoX85/Cloud-Assess/issues/189'
EVIDENCE_ENTRYPOINTS = ('CURRENT_STATE.md', 'ACTION_REGISTER.md', 'CONTINUE.md', 'RESTART_PROMPT.md')
TASK_ID_PATTERN = re.compile(r'\b(?:DC-(?:[A-Z]+-)*\d+(?:-\d+)*|AR-\d{2}|DV-\d{3}|OP-\d{3}|R7-[FOP]\d{2})\b')


def require_task_inventory(inventory, sources):
    legacy = inventory['legacy_ids']
    ids = [row['id'] for row in legacy]
    if len(ids) != len(set(ids)) or inventory['legacy_count'] != len(ids):
        raise ValueError('duplicate or missing legacy task identities')
    obligations = inventory['obligations']
    current = [row['id'] for row in obligations]
    if len(current) != len(set(current)) or inventory['obligation_count'] != len(current):
        raise ValueError('duplicate or missing current obligations')
    required = {'SEC-BASELINE', 'DC-CI-EVIDENCE-20261022', 'AUDIT-REMAINDER', 'R7-F01', 'R7-F03', 'R7-F04', 'ANC-COMPARE', 'ANC-ALTERNATIVE-VM', 'ANC-MCP', 'ANC-TYPES', 'SEC-INDEPENDENT', 'GATE-004'}
    required.update({'R7-F02', 'R7-O01', 'R7-O02', 'R7-O03', 'R7-O04', 'DC-007', 'DC-012', 'DC-013', 'CONTENT-FRESHNESS', 'PROCESS-EVIDENCE', 'PROCESS-STATE', 'PROCESS-TIMING', 'ACCEPT-SECURITY-LICENSE', 'ACCEPT-OPERATIONS'})
    required.update(f'AR-{n:02}' for n in range(1, 9))
    required.update(f'OFFLINE-{n:02}' for n in range(1, 16))
    if not required.issubset(current):
        raise ValueError('required programme/security/acceptance task lost')
    for source in inventory['sources']:
        expected = set(source['ids'])
        if not expected.issubset(ids) or not expected.issubset(TASK_ID_PATTERN.findall(sources[source['path']])):
            raise ValueError('captured source task identity lost')
    if any(row.get('mapped_to') and row['mapped_to'] not in current for row in legacy):
        raise ValueError('task maps to absent obligation')


def require_evidence_task(task, pages):
    if task.get('id') != EVIDENCE_TASK or task.get('issue') != EVIDENCE_ISSUE:
        raise ValueError('missing stable evidence task or issue')
    if task.get('status') not in ('deferred_open', 'in_progress', 'resolved'):
        raise ValueError('invalid evidence task status')
    if set(pages) != set(EVIDENCE_ENTRYPOINTS):
        raise ValueError('missing evidence task entrypoint')
    if any(EVIDENCE_TASK not in text[:6000] or EVIDENCE_ISSUE not in text[:6000] for text in pages.values()):
        raise ValueError('evidence task must remain visible in every resume entrypoint')
    required = ('authorized_destination', 'retention_owner', 'manifest', 'readback', 'restore')
    if task['status'] == 'resolved' and any(not isinstance(task.get('closure', {}).get(key), str) or not task['closure'][key].strip() for key in required):
        raise ValueError('resolved evidence task lacks closure references')


def require_authority(text):
    header = text.replace('\r\n', '\n').split('\n---\n', 1)[0]
    if len(header.splitlines()) > 14 or '[CURRENT_STATE](CURRENT_STATE.md)' not in header:
        raise ValueError('missing compact current-state authority')
    if 'historical snapshots' not in header or 'behavioral requirements remain in force' not in header:
        raise ValueError('historical status and active contracts must be distinguished')
    if any(word in header for word in ('STOP', 'pending at publication', 'chore/', 'PR148')):
        raise ValueError('dated task state must live in CURRENT_STATE, not resume headers')


class HandoffTests(unittest.TestCase):
    def test_captured_task_inventory_is_conserved(self):
        inventory = json.loads((ROOT / 'docs/TASK_CARRY_FORWARD_20261022.json').read_text(encoding='utf-8'))
        sources = {row['path']: (ROOT / row['path']).read_text(encoding='utf-8') for row in inventory['sources']}
        require_task_inventory(inventory, sources)

    def test_dropped_task_inventory_rejects(self):
        inventory = json.loads((ROOT / 'docs/TASK_CARRY_FORWARD_20261022.json').read_text(encoding='utf-8'))
        sources = {row['path']: (ROOT / row['path']).read_text(encoding='utf-8') for row in inventory['sources']}
        for task_id in ('OFFLINE-08', 'ANC-MCP', 'SEC-INDEPENDENT', 'DC-CI-EVIDENCE-20261022'):
            bad = {**inventory, 'obligations': [row for row in inventory['obligations'] if row['id'] != task_id], 'obligation_count': inventory['obligation_count'] - 1}
            with self.subTest(task=task_id), self.assertRaises(ValueError):
                require_task_inventory(bad, sources)
        source = next(row for row in inventory['sources'] if row['ids'])
        task_id = source['ids'][0]
        with self.assertRaises(ValueError):
            require_task_inventory(inventory, {**sources, source['path']: sources[source['path']].replace(task_id, '')})
        with self.assertRaises(ValueError):
            require_task_inventory({**inventory, 'legacy_ids': inventory['legacy_ids'][1:], 'legacy_count': inventory['legacy_count'] - 1}, sources)

    def test_evidence_task_cannot_silently_disappear(self):
        task = json.loads((ROOT / 'docs/CI_EVIDENCE_TASK.json').read_text(encoding='utf-8'))
        pages = {name: (ROOT / 'docs' / name).read_text(encoding='utf-8') for name in EVIDENCE_ENTRYPOINTS}
        require_evidence_task(task, pages)

    def test_lost_or_prematurely_closed_evidence_task_rejects(self):
        task = {'id': EVIDENCE_TASK, 'issue': EVIDENCE_ISSUE, 'status': 'deferred_open'}
        pages = {name: EVIDENCE_TASK + ' ' + EVIDENCE_ISSUE for name in EVIDENCE_ENTRYPOINTS}
        require_evidence_task(task, pages)
        for name in EVIDENCE_ENTRYPOINTS:
            for marker in (EVIDENCE_TASK, EVIDENCE_ISSUE):
                with self.subTest(page=name, marker=marker), self.assertRaises(ValueError):
                    require_evidence_task(task, {**pages, name: pages[name].replace(marker, '')})
        for bad in ({**task, 'status': 'closed'}, {**task, 'status': 'resolved'}, {**task, 'id': 'lost'}, {**task, 'issue': 'wrong'}):
            with self.assertRaises(ValueError):
                require_evidence_task(bad, pages)
        with self.assertRaises(ValueError):
            require_evidence_task(task, {})
        closed = {**task, 'status': 'resolved', 'closure': {key: 'recorded evidence' for key in ('authorized_destination', 'retention_owner', 'manifest', 'readback', 'restore')}}
        require_evidence_task(closed, pages)
        for key in closed['closure']:
            with self.subTest(missing=key), self.assertRaises(ValueError):
                require_evidence_task({**closed, 'closure': {**closed['closure'], key: ''}}, pages)

    def test_every_resume_entrypoint_has_one_stable_authority(self):
        for name in ENTRYPOINTS:
            with self.subTest(document=name):
                require_authority((ROOT / 'docs' / name).read_text(encoding='utf-8'))

    def test_misleading_resume_headers_are_rejected(self):
        good = ('Read [CURRENT_STATE](CURRENT_STATE.md). Earlier status sections are historical snapshots; '
                'behavioral requirements remain in force.\n\n---\n')
        require_authority(good)
        for bad in (good.replace('CURRENT_STATE.md', 'SESSION_HANDOVER.md'),
                    good.replace('historical snapshots', 'current'),
                    good.replace('behavioral requirements remain in force', 'all requirements obsolete'),
                    good.replace('Read ', 'STOP at PR148. Read ')):
            with self.assertRaises(ValueError):
                require_authority(bad)

    def test_agent_and_readme_route_to_current_state(self):
        for name in ('AGENTS.md', 'README.md'):
            self.assertIn('docs/CURRENT_STATE.md', (ROOT / name).read_text(encoding='utf-8')[:6000])


if __name__ == '__main__':
    unittest.main()
