#!/usr/bin/env python3
"""Check actual GitHub job outcomes; an arbitrary failed probe is never proof."""
import argparse
import hashlib
import json
from pathlib import Path
import re

VALIDATOR_STEP = 'Require every native worker to succeed'


def verify(run, jobs, case, sha, read_log):
    if case not in ('success', 'failure', 'skipped'):
        raise ValueError('unknown probe case')
    conclusion = 'success' if case == 'success' else 'failure'
    if (run.get('head_sha') != sha or run.get('event') != 'push' or
            run.get('status') != 'completed' or run.get('conclusion') != conclusion or
            run.get('path') != '.github/workflows/gate-contract-probe.yml' or
            not re.fullmatch(r'qa/gate-contract-[A-Za-z0-9-]+-' + case, run.get('head_branch', ''))):
        raise ValueError('wrong revision, branch, workflow or overall outcome')
    expected = {}
    for platform in ('linux', 'windows'):
        for suffix in ('native', 'controls-1', 'controls-2', 'controls-3'):
            name = platform + '-' + suffix
            expected[name] = case if suffix == 'controls-1' and case != 'success' else 'success'
        expected[platform + '-probe-gate'] = conclusion
    rows = jobs.get('jobs', [])
    if len(rows) != 10 or len({j.get('id') for j in rows}) != 10 or {j.get('name') for j in rows} != set(expected):
        raise ValueError('missing, duplicate or unexpected probe jobs')
    fingerprints = []
    for job in rows:
        name = job['name']
        if job.get('run_id') != run.get('id') or job.get('status') != 'completed' or job.get('conclusion') != expected[name]:
            raise ValueError('unexpected job outcome: ' + name)
        if name.endswith('-probe-gate'):
            steps = job.get('steps', [])
            validator = [s for s in steps if s.get('name') == VALIDATOR_STEP]
            if len(validator) != 1 or validator[0].get('conclusion') != conclusion:
                raise ValueError('gate did not execute its actual validator')
            if any(s.get('conclusion') != 'success' for s in steps if s.get('name') != VALIDATOR_STEP):
                raise ValueError('gate failed outside the intended validator')
            raw = read_log(job['id'])
            text = raw.decode('utf-8')
            message = ('Required native gate verified: ' if case == 'success' else
                       'Required native gate rejected incomplete or unsuccessful workers')
            if message not in text:
                raise ValueError('missing actual validator diagnostic')
            fingerprints.append({'job': job['id'], 'name': name, 'bytes': len(raw), 'sha256': hashlib.sha256(raw).hexdigest()})
        else:
            label = 'windows-2022' if name.startswith('windows-') else 'ubuntu-24.04'
            # A job skipped before scheduling has no runner labels or steps.
            if expected[name] != 'skipped' and label not in job.get('labels', []):
                raise ValueError('probe worker ran on the wrong platform')
            if expected[name] != 'skipped':
                step = [s for s in job.get('steps', []) if s.get('name') == 'Controlled worker outcome']
                if len(step) != 1 or step[0].get('conclusion') != expected[name]:
                    raise ValueError('failure was not the controlled worker step')
                if any(s.get('conclusion') != 'success' for s in job.get('steps', []) if s.get('name') != 'Controlled worker outcome'):
                    raise ValueError('worker failed outside its controlled step')
    return {'case': case, 'sha': sha, 'run': run['id'], 'verified_jobs': 10, 'gates': fingerprints}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--run-json', type=Path, required=True)
    parser.add_argument('--jobs-json', type=Path, required=True)
    parser.add_argument('--logs-dir', type=Path, required=True)
    parser.add_argument('--case', choices=('success', 'failure', 'skipped'), required=True)
    parser.add_argument('--sha', required=True)
    args = parser.parse_args()
    report = verify(json.loads(args.run_json.read_text(encoding='utf-8')),
                    json.loads(args.jobs_json.read_text(encoding='utf-8')), args.case, args.sha,
                    lambda job_id: (args.logs_dir / (str(job_id) + '.log')).read_bytes())
    print(json.dumps(report, sort_keys=True))


if __name__ == '__main__':
    main()
