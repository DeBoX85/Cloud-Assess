#!/usr/bin/env python3
"""Adversarial examples for the additional local orientation checker."""
from copy import deepcopy
import importlib.util
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('handover_check', ROOT / 'scripts/check-handover.py')
check = importlib.util.module_from_spec(spec)
spec.loader.exec_module(check)


class RecoveryTests(unittest.TestCase):
    def setUp(self):
        self.state = check.read_json(ROOT / 'docs/RECOVERY_STATE.json')
        self.inventory = check.read_json(ROOT / 'docs/TASK_CARRY_FORWARD_20261022.json')
        self.followups = check.read_json(ROOT / 'docs/RECOVERY_FOLLOWUPS_20261030.json')

    def test_current_orientation(self):
        check.validate(self.state, self.inventory, self.followups)
        check.check_checkout(self.state, self.inventory)

    def test_preparation_green_cannot_certify_product(self):
        self.state['product_qa']['candidate'] = self.state['roles']['preserved_preparation']['commit']
        with self.assertRaisesRegex(ValueError, 'QA revision'):
            check.validate(self.state, self.inventory, self.followups)

    def test_role_conflation_rejects(self):
        self.state['roles']['verified_product']['commit'] = self.state['roles']['accepted']['commit']
        with self.assertRaisesRegex(ValueError, 'roles conflated'):
            check.validate(self.state, self.inventory, self.followups)

    def test_missing_original_and_provenance_claims_reject(self):
        for field, value in [('original_streams_recovered', True), ('replacement_provenance', '')]:
            bad = deepcopy(self.state)
            bad['evidence'][field] = value
            with self.subTest(field=field), self.assertRaisesRegex(ValueError, 'provenance'):
                check.validate(bad, self.inventory, self.followups)

    def test_missing_retention_steps_reject(self):
        for field in ['remote_readback', 'fresh_restore']:
            bad = deepcopy(self.state)
            bad['evidence'][field] = False
            with self.subTest(field=field), self.assertRaisesRegex(ValueError, 'fresh restore'):
                check.validate(bad, self.inventory, self.followups)

    def test_missing_obligation_and_followup_reject(self):
        bad = deepcopy(self.inventory)
        bad['obligations'].pop()
        with self.assertRaisesRegex(ValueError, 'obligations missing'):
            check.validate(self.state, bad, self.followups)
        with self.assertRaisesRegex(ValueError, 'followups lost'):
            check.validate(self.state, self.inventory, {'items': []})

    def test_coordinator_premature_closure_rejects(self):
        next(x for x in self.inventory['obligations'] if x['id'] == 'OFFLINE-10')['state'] = 'complete'
        with self.assertRaisesRegex(ValueError, 'coordinator task'):
            check.validate(self.state, self.inventory, self.followups)

    def test_premature_caution_exit_and_backup_claim_reject(self):
        self.state['enhanced_verification']['active'] = False
        with self.assertRaisesRegex(ValueError, 'prematurely disabled'):
            check.validate(self.state, self.inventory, self.followups)
        self.state['enhanced_verification']['active'] = True
        self.state['evidence']['independent_backup_verified'] = True
        with self.assertRaisesRegex(ValueError, 'backup is still unverified'):
            check.validate(self.state, self.inventory, self.followups)

    def test_duplicate_json_keys_reject(self):
        with self.assertRaisesRegex(ValueError, 'duplicate JSON key'):
            check.unique_object([('commit', 'one'), ('commit', 'other')])


if __name__ == '__main__':
    unittest.main()
