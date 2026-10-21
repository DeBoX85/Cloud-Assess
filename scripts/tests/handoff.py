#!/usr/bin/env python3
"""Keep resume authority explicit without treating historical snapshots as current."""
from pathlib import Path
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


def require_authority(text):
    header = text.replace('\r\n', '\n').split('\n---\n', 1)[0]
    if len(header.splitlines()) > 14 or '[CURRENT_STATE](CURRENT_STATE.md)' not in header:
        raise ValueError('missing compact current-state authority')
    if 'historical snapshots' not in header or 'behavioral requirements remain in force' not in header:
        raise ValueError('historical status and active contracts must be distinguished')
    if any(word in header for word in ('STOP', 'pending at publication', 'chore/', 'PR148')):
        raise ValueError('dated task state must live in CURRENT_STATE, not resume headers')


class HandoffTests(unittest.TestCase):
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
