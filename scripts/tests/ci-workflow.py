#!/usr/bin/env python3
"""Independent failures for skipped dependencies, lost checks and weaker budgets."""
import importlib.util
from collections import Counter
import json
import re
from pathlib import Path
import subprocess
import sys
import unittest
from unittest.mock import patch

sys.dont_write_bytecode = True
ROOT = Path(__file__).resolve().parents[2]


def load(name):
    spec = importlib.util.spec_from_file_location(name, ROOT / 'scripts/ci' / (name + '.py'))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


workflow = load('generate_quality_workflow')
gate = load('required_gate')


def require_build_before_consumer(text):
    for platform, consumer in [('linux', 'Test with race detector'), ('windows', 'Test Go packages on Windows')]:
        match = re.search(r'^  ' + platform + r'-native:\n(.*?)(?=^  [a-z][a-z0-9-]*:\n|\Z)', text, re.M | re.S)
        if match is None:
            raise ValueError('missing native job')
        body = match.group(1)
        marker = '      - name: Verify isolated required inventory Build controls\n'
        if body.count(marker) != 1 or body.index(marker) >= body.index('      - name: ' + consumer + '\n'):
            raise ValueError('Build proof missing or hidden behind absent admission consumer')
        block = body.split(marker, 1)[1].split('      - name:', 1)[0]
        if 'scripts/tests/required-inventory-build.py' not in block or 'continue-on-error' in block or 'if:' in block:
            raise ValueError('Build proof weakened')


def require_admission_before_consumer(text):
    for platform, consumer in [('linux', 'Test with race detector'), ('windows', 'Test Go packages on Windows')]:
        body = re.search(r'^  ' + platform + r'-native:\n(.*?)(?=^  [a-z][a-z0-9-]*:\n|\Z)', text, re.M | re.S).group(1)
        marker = '      - name: Verify required evidence Admission controls\n'
        if body.count(marker) != 1 or body.index(marker) >= body.index('      - name: ' + consumer + '\n'):
            raise ValueError('Admission proof missing or late')
        block = body.split(marker, 1)[1].split('      - name:', 1)[0]
        if 'scripts/tests/required-inventory-admission.py' not in block or 'continue-on-error' in block or 'if:' in block:
            raise ValueError('Admission proof weakened')


def require_legacy_admission(text):
    for platform, consumer in [('linux', 'Test with race detector'), ('windows', 'Test Go packages on Windows')]:
        match = re.search(r'^  '+platform+r'-native:\n(.*?)(?=^  [a-z][a-z0-9-]*:\n|\Z)',text,re.M|re.S)
        if match is None: raise ValueError('missing native job')
        body=match.group(1);marker='      - name: Verify legacy mutation assertion admission\n'
        if body.count(marker)!=1 or body.index(marker)>=body.index('      - name: '+consumer+'\n'): raise ValueError('legacy assertion proof missing/late')
        block=body.split(marker,1)[1].split('      - name:',1)[0]
        if 'scripts/tests/mutation-assertions.py' not in block or 'continue-on-error' in block or 'if:' in block: raise ValueError('legacy assertion proof weakened')

def require_tag_collision_controls(text):
    for host,consumer in [('linux','Verify pinned source-data provenance'),('windows','Parse scripts with Windows PowerShell 5.1')]:
        match=re.search(r'^  '+host+r'-native:\n(.*?)(?=^  [a-z]|\Z)',text,re.M|re.S)
        if not match:raise ValueError('missing tag-control native host')
        body=match.group(1);marker='      - name: Verify tag collision rejection controls\n'
        if body.count(marker)!=1 or body.index(marker)>=body.index('      - name: '+consumer+'\n'):raise ValueError('tag controls missing/late')
        block=body.split(marker,1)[1].split('      - name:',1)[0]
        if 'scripts/tests/tag-collision-controls.py' not in block or 'continue-on-error' in block or 'if:' in block:raise ValueError('tag controls weakened')


def require_comparison_json_controls(text):
    for host, consumer in [('linux', 'Verify pinned source-data provenance'), ('windows', 'Parse scripts with Windows PowerShell 5.1')]:
        match = re.search(r'^  '+host+r'-native:\n(.*?)(?=^  [a-z]|\Z)', text, re.M|re.S)
        if not match: raise ValueError('missing comparison JSON native host')
        body = match.group(1); marker = '      - name: Verify comparison JSON admission controls\n'
        if body.count(marker) != 1 or body.index(marker) >= body.index('      - name: '+consumer+'\n'):
            raise ValueError('comparison JSON controls missing/late')
        block = body.split(marker, 1)[1].split('      - name:', 1)[0]
        if 'scripts/tests/comparison-json-controls.py' not in block or 'continue-on-error' in block or 'if:' in block:
            raise ValueError('comparison JSON controls weakened')


def require_coordinator_specification(text):
    for host, consumer in [('linux', 'Verify pinned source-data provenance'), ('windows', 'Parse scripts with Windows PowerShell 5.1')]:
        match = re.search(r'^  '+host+r'-native:\n(.*?)(?=^  [a-z]|\Z)', text, re.M|re.S)
        if not match: raise ValueError('missing coordinator specification native host')
        body = match.group(1); marker = '      - name: Verify coordinator scheduling specification\n'
        if body.count(marker) != 1 or body.index(marker) >= body.index('      - name: '+consumer+'\n'):
            raise ValueError('coordinator specification missing/late')
        block = body.split(marker, 1)[1].split('      - name:', 1)[0]
        if 'scripts/tests/coordinator-scheduling.py' not in block or 'continue-on-error' in block or 'if:' in block:
            raise ValueError('coordinator specification weakened')


def require_subscription_reservation_controls(text):
    for host, consumer in [('linux', 'Verify pinned source-data provenance'), ('windows', 'Parse scripts with Windows PowerShell 5.1')]:
        match = re.search(r'^  '+host+r'-native:\n(.*?)(?=^  [a-z]|\Z)', text, re.M|re.S)
        if not match: raise ValueError('missing subscription reservation controls native host')
        body = match.group(1); marker = '      - name: Verify subscription reservation contract controls\n'
        if body.count(marker) != 1 or body.index(marker) >= body.index('      - name: '+consumer+'\n'):
            raise ValueError('subscription reservation controls missing/late')
        block = body.split(marker, 1)[1].split('      - name:', 1)[0]
        if 'scripts/tests/subscription-reservation-controls.py' not in block or 'continue-on-error' in block or 'if:' in block:
            raise ValueError('subscription reservation controls weakened')


class GateTests(unittest.TestCase):
    def test_coordinator_specification_both_hosts(self):
        require_coordinator_specification(workflow.generate())

    def test_coordinator_specification_cannot_be_removed_or_weakened(self):
        text = workflow.generate(); marker = '      - name: Verify coordinator scheduling specification\n'
        for host in ('linux', 'windows'):
            before, rest = text.split('  '+host+'-native:\n', 1)
            for mutation in ('      - name: Lost coordinator specification\n', marker+'        continue-on-error: true\n', marker+'        if: false\n'):
                with self.assertRaises(ValueError):
                    require_coordinator_specification(before+'  '+host+'-native:\n'+rest.replace(marker, mutation, 1))

    def test_coordinator_specification_cannot_run_late(self):
        text = workflow.generate(); marker = '      - name: Verify coordinator scheduling specification\n'
        start = text.index(marker); end = text.index('      - name:', start+len(marker)); block = text[start:end]
        lost = text[:start]+text[end:]; later = lost.index('      - name: Vet\n')
        with self.assertRaises(ValueError):
            require_coordinator_specification(lost[:later]+block+lost[later:])

    def test_subscription_reservation_controls_both_hosts(self):
        require_subscription_reservation_controls(workflow.generate())

    def test_subscription_reservation_controls_cannot_be_removed_or_weakened(self):
        text = workflow.generate(); marker = '      - name: Verify subscription reservation contract controls\n'
        for host in ('linux', 'windows'):
            before, rest = text.split('  '+host+'-native:\n', 1)
            for mutation in ('      - name: Lost subscription reservation controls\n', marker+'        continue-on-error: true\n', marker+'        if: false\n'):
                with self.assertRaises(ValueError):
                    require_subscription_reservation_controls(before+'  '+host+'-native:\n'+rest.replace(marker, mutation, 1))

    def test_subscription_reservation_controls_cannot_run_late(self):
        text = workflow.generate(); marker = '      - name: Verify subscription reservation contract controls\n'
        start = text.index(marker); end = text.index('      - name:', start+len(marker)); block = text[start:end]
        lost = text[:start]+text[end:]; later = lost.index('      - name: Vet\n')
        with self.assertRaises(ValueError):
            require_subscription_reservation_controls(lost[:later]+block+lost[later:])

    def test_comparison_json_controls_both_hosts(self):
        require_comparison_json_controls(workflow.generate())

    def test_comparison_json_controls_cannot_be_removed_or_weakened(self):
        text = workflow.generate(); marker = '      - name: Verify comparison JSON admission controls\n'
        for host in ('linux', 'windows'):
            before, rest = text.split('  '+host+'-native:\n', 1)
            for mutation in ('      - name: Lost JSON controls\n', marker+'        continue-on-error: true\n', marker+'        if: false\n'):
                with self.assertRaises(ValueError):
                    require_comparison_json_controls(before+'  '+host+'-native:\n'+rest.replace(marker, mutation, 1))

    def test_comparison_json_controls_cannot_run_late(self):
        text = workflow.generate(); marker = '      - name: Verify comparison JSON admission controls\n'
        start = text.index(marker); end = text.index('      - name:', start+len(marker)); block = text[start:end]
        lost = text[:start]+text[end:]; later = lost.index('      - name: Vet\n')
        with self.assertRaises(ValueError):
            require_comparison_json_controls(lost[:later]+block+lost[later:])

    def test_tag_collision_controls_both_hosts(self):
        require_tag_collision_controls(workflow.generate())

    def test_tag_collision_controls_cannot_be_removed_or_weakened(self):
        text=workflow.generate();marker='      - name: Verify tag collision rejection controls\n'
        for bad in [text.replace(marker,'      - name: Lost tag controls\n',1),text.replace(marker,marker+'        continue-on-error: true\n',1)]:
            with self.assertRaises(ValueError):require_tag_collision_controls(bad)

    def test_legacy_assertion_admission_required_both_hosts(self):
        require_legacy_admission(workflow.generate())

    def test_legacy_assertion_admission_cannot_be_removed_or_weakened(self):
        text=workflow.generate();marker='      - name: Verify legacy mutation assertion admission\n'
        for bad in [text.replace(marker,'      - name: Lost assertion proof\n',1),text.replace(marker,marker+'        continue-on-error: true\n',1)]:
            with self.assertRaises(ValueError): require_legacy_admission(bad)


    def test_admission_proof_precedes_consumer(self):
        require_admission_before_consumer(workflow.generate())

    def test_admission_proof_missing_or_weakened(self):
        text = workflow.generate()
        marker = '      - name: Verify required evidence Admission controls\n'
        for bad in [text.replace(marker, '      - name: Missing Admission\n', 1), text.replace(marker, marker+'        continue-on-error: true\n', 1)]:
            with self.assertRaises(ValueError):
                require_admission_before_consumer(bad)


    def test_build_proof_precedes_full_consumer_on_both_hosts(self):
        require_build_before_consumer(workflow.generate())

    def test_lost_or_late_build_proof_rejects(self):
        text = workflow.generate()
        marker = '      - name: Verify isolated required inventory Build controls\n'
        with self.assertRaises(ValueError):
            require_build_before_consumer(text.replace(marker, '      - name: Lost Build proof\n', 1))
        start = text.index(marker)
        end = text.index('      - name:', start + len(marker))
        block = text[start:end]
        lost = text[:start] + text[end:]
        later = lost.index('      - name: Vet\n')
        with self.assertRaises(ValueError):
            require_build_before_consumer(lost[:later] + block + lost[later:])


    def test_actual_cli_success_and_failure_exit(self):
        for platform in ('linux', 'windows'):
            needs = {k: {'result': 'success', 'outputs': {}} for k in gate.worker_ids(platform)}
            for status, expected in [('success', 0), ('skipped', 1), ('failure', 1), ('cancelled', 1)]:
                needs[platform + '-native']['result'] = status
                with patch.dict('os.environ', {'NEEDS_JSON': json.dumps(needs)}):
                    result = subprocess.run([sys.executable, str(ROOT / 'scripts/ci/required_gate.py'), '--platform', platform], capture_output=True, timeout=10)
                self.assertEqual(result.returncode, expected, (platform, status))

    def test_every_required_worker_rejects_non_success(self):
        for platform in ('linux', 'windows'):
            expected = [platform + '-native', platform + '-controls-1', platform + '-controls-2', platform + '-controls-3']
            self.assertEqual(gate.worker_ids(platform), expected)
            for key in expected:
                for status in ('failure', 'cancelled', 'skipped', 'neutral', 'timed_out', '', None):
                    with self.subTest(platform=platform, worker=key, status=status):
                        needs = {k: {'result': 'success'} for k in expected}
                        needs[key]['result'] = status
                        with self.assertRaises(ValueError):
                            gate.require_workers(json.dumps(needs), platform)

    def test_missing_extra_duplicate_and_malformed_dependencies(self):
        good = {k: {'result': 'success'} for k in gate.worker_ids('linux')}
        bad = [[], None, {}, {**good, 'unexpected': {'result': 'success'}}]
        for key in good:
            bad.append({k: value for k, value in good.items() if k != key})
        for value in bad:
            with self.assertRaises(ValueError):
                gate.require_workers(json.dumps(value), 'linux')
        for value in ('{', '{"linux-native":{},"linux-native":{}}'):
            with self.assertRaises(ValueError):
                gate.require_workers(value, 'linux')


class ConservationTests(unittest.TestCase):
    def test_every_original_step_block_runs_on_its_original_platform(self):
        _, before = workflow.split_jobs(workflow.BASELINE.read_text(encoding='utf-8'))
        _, after = workflow.split_jobs(workflow.WORKFLOW.read_text(encoding='utf-8'))
        for platform, old_id, count in [('linux', 'quality', 62), ('windows', 'windows-validation', 47)]:
            _, originals = workflow.split_steps(before[old_id])
            self.assertEqual(len(originals), count)
            expected = Counter(originals.values())
            for name in ('Checkout', 'Set up Go'):
                expected[originals[name]] += 3
            actual = Counter()
            for job_id in [platform + '-native', platform + '-controls-1', platform + '-controls-2', platform + '-controls-3']:
                _, steps = workflow.split_steps(after[job_id])
                self.assertIn('Verify parallel workflow conservation and gate guards', steps)
                del steps['Verify parallel workflow conservation and gate guards']
                plan = json.loads(workflow.PLAN.read_text(encoding='utf-8'))
                for addition in plan['additions'][platform]:
                    if job_id == platform + '-' + addition['worker']:
                        self.assertIn(addition['name'], steps)
                        del steps[addition['name']]
                for block in steps.values():
                    block = block.replace("go-version: '1.26.10'", "go-version: '1.26.9'")
                    block = block.replace('go test -count=1 -covermode=atomic', 'go test -covermode=atomic')
                    block = block.replace('run: go test -count=1 ./...', 'run: go test ./...')
                    actual[block] += 1
            self.assertEqual(actual, expected)

    def test_native_crlf_checkout_preserves_baseline_authority(self):
        raw = workflow.BASELINE.read_bytes().replace(b'\r\n', b'\n')
        with patch.object(Path, 'read_bytes', return_value=raw.replace(b'\n', b'\r\n')):
            workflow.read_baseline()
        with patch.object(Path, 'read_bytes', return_value=raw.replace(b'go test -race', b'go test', 1)):
            with self.assertRaises(ValueError):
                workflow.read_baseline()

    def test_actual_workflow_is_conserved(self):
        workflow.verify(workflow.WORKFLOW.read_text(encoding='utf-8'))

    def test_actual_workflow_rejects_lost_checks_or_weaker_guards(self):
        text = workflow.WORKFLOW.read_text(encoding='utf-8')
        changes = [
            ('go test -race -count=1 ./...', 'go test ./...'),
            ('< 75.0', '< 70.0'),
            ('-fuzztime 500000x', '-fuzztime 50x'),
            ('run: python3 scripts/tests/compute-vm-sku-collector-mutation.py', 'run: true'),
            ('if: ${{ always() }}', 'if: ${{ success() }}'),
            ('linux-native, linux-controls-1, linux-controls-2, linux-controls-3', 'linux-native'),
            ('runs-on: windows-2022', 'runs-on: ubuntu-24.04'),
            ('contents: read', 'contents: write'),
            ("go-version: '1.26.10'", "go-version: '1.26.9'"),
            ('run: go test -count=1 ./...', 'run: go test ./...'),
        ]
        for old, new in changes:
            with self.subTest(guard=old):
                self.assertIn(old, text)
                with self.assertRaises(ValueError):
                    workflow.verify(text.replace(old, new, 1))

    def test_plan_cannot_drop_duplicate_or_move_a_suite(self):
        header, baseline, plan = workflow.read_baseline()
        for operation in ('drop', 'duplicate', 'foreign'):
            changed = json.loads(json.dumps(plan))
            groups = changed['groups']['linux']
            if operation == 'drop':
                groups[0].pop()
            elif operation == 'duplicate':
                groups[0].append(groups[1][0])
            else:
                groups[0][0] = 'Nonexistent or foreign-platform suite'
            with patch.object(workflow, 'read_baseline', return_value=(header, baseline, changed)):
                with self.assertRaises(ValueError):
                    workflow.generate()


class AdditiveInventoryTests(unittest.TestCase):
    def test_process_checks_are_mandatory_on_both_native_workers(self):
        _, jobs = workflow.split_jobs(workflow.WORKFLOW.read_text(encoding='utf-8'))
        for platform in ('linux', 'windows'):
            _, steps = workflow.split_steps(jobs[platform + '-native'])
            for name, script in [('Verify hosted gate contract definitions', 'gate-probe.py'), ('Verify handoff authority consistency', 'handoff.py')]:
                self.assertIn(name, steps)
                self.assertIn('run: ' + ('python3' if platform == 'linux' else 'python') + ' scripts/tests/' + script, steps[name])
                self.assertNotIn('        if:', steps[name])
                self.assertNotIn('continue-on-error', steps[name])

    def test_additive_inventory_rejects_replacement_skip_and_foreign_paths(self):
        header, baseline, plan = workflow.read_baseline()
        bad = [
            {'name': 'Checkout', 'script': 'scripts/tests/gate-probe.py', 'worker': 'native'},
            {'name': 'Extra check', 'script': '../outside.py', 'worker': 'native'},
            {'name': 'Extra check', 'script': 'scripts/tests/gate-probe.py', 'worker': 'other'},
            {'name': 'Extra check', 'script': 'scripts/tests/gate-probe.py', 'worker': 'native', 'if': 'false'},
        ]
        for addition in bad:
            changed = json.loads(json.dumps(plan))
            changed['additions']['linux'].append(addition)
            with patch.object(workflow, 'read_baseline', return_value=(header, baseline, changed)):
                with self.assertRaises(ValueError):
                    workflow.generate()


if __name__ == '__main__':
    unittest.main()
