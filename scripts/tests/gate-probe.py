#!/usr/bin/env python3
"""Check genuine validator rejection rather than unrelated infrastructure failures."""
import importlib.util
from pathlib import Path
import unittest
import copy

ROOT = Path(__file__).resolve().parents[2]


def load(name):
    spec = importlib.util.spec_from_file_location(name, ROOT / 'scripts/ci' / (name + '.py'))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


probe = load('generate_gate_probe')
verifier = load('verify_gate_probe')


def fixture(case):
    result = 'success' if case == 'success' else 'failure'
    run = {'id': 10, 'head_sha': 'abc', 'event': 'push', 'status': 'completed', 'conclusion': result,
           'path': '.github/workflows/gate-contract-probe.yml', 'head_branch': 'qa/gate-contract-test-' + case}
    jobs = []
    for platform in ('linux', 'windows'):
        for suffix in ('native', 'controls-1', 'controls-2', 'controls-3', 'probe-gate'):
            name = platform + '-' + suffix
            status = result if suffix == 'probe-gate' else case if suffix == 'controls-1' else 'success'
            jobs.append({'id': len(jobs) + 1, 'run_id': 10, 'name': name, 'status': 'completed', 'conclusion': status,
                         'labels': [] if status == 'skipped' else ['windows-2022' if platform == 'windows' else 'ubuntu-24.04'],
                         'steps': [] if status == 'skipped' else [{'name': verifier.VALIDATOR_STEP if suffix == 'probe-gate' else 'Controlled worker outcome', 'conclusion': status}]})
    log = b'Required native gate verified: ' if case == 'success' else b'Required native gate rejected incomplete or unsuccessful workers'
    return run, {'jobs': jobs}, lambda _: log


class ProbeTests(unittest.TestCase):
    def test_generated_wiring_and_branch_isolation(self):
        text = probe.OUTPUT.read_text(encoding='utf-8')
        probe.verify(text)
        self.assertNotIn('pull_request:', text)
        self.assertNotIn('continue-on-error', text)
        for original in ('quality', 'windows-validation'):
            self.assertIn(probe.production_jobs()[original].rstrip(), text)
        for old, new in [('always()', 'success()'), ('contents: read', 'contents: write'), ('-skipped', '-ignored')]:
            with self.assertRaises(ValueError):
                probe.verify(text.replace(old, new, 1))

    def test_three_literal_expected_outcomes(self):
        for case in ('success', 'failure', 'skipped'):
            run, jobs, logs = fixture(case)
            self.assertEqual(verifier.verify(run, jobs, case, 'abc', logs)['verified_jobs'], 10)

    def test_rejects_wrong_head_masked_failure_missing_job_and_infrastructure_failure(self):
        run, jobs, logs = fixture('failure')
        bad_jobs = []
        missing = copy.deepcopy(jobs); missing['jobs'].pop(); bad_jobs.append(missing)
        masked = copy.deepcopy(jobs); masked['jobs'][4]['conclusion'] = 'success'; bad_jobs.append(masked)
        unrelated = copy.deepcopy(jobs); unrelated['jobs'][0]['conclusion'] = 'failure'; bad_jobs.append(unrelated)
        infrastructure = copy.deepcopy(jobs); infrastructure['jobs'][4]['steps'].append({'name': 'Checkout gate validator', 'conclusion': 'failure'}); bad_jobs.append(infrastructure)
        for bad in bad_jobs:
            with self.assertRaises(ValueError):
                verifier.verify(run, bad, 'failure', 'abc', logs)
        with self.assertRaises(ValueError):
            verifier.verify(run, jobs, 'failure', 'wrong-head', logs)
        with self.assertRaises(ValueError):
            verifier.verify(run, jobs, 'failure', 'abc', lambda _: b'network failed')


if __name__ == '__main__':
    unittest.main()
