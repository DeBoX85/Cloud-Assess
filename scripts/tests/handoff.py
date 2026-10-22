#!/usr/bin/env python3
"""Keep resume authority explicit without treating historical snapshots as current."""
from pathlib import Path
import json
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
