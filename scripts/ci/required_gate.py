#!/usr/bin/env python3
"""Keep the two protected check names strict across isolated worker jobs."""
import argparse
import json
import os


def worker_ids(platform):
    if platform not in ('linux', 'windows'):
        raise ValueError('unknown native platform')
    return [platform + '-native'] + [platform + '-controls-' + str(i) for i in (1, 2, 3)]


def unique_object(pairs):
    value = {}
    for key, item in pairs:
        if key in value:
            raise ValueError('duplicate dependency identity')
        value[key] = item
    return value


def require_workers(payload, platform):
    needs = json.loads(payload, object_pairs_hook=unique_object)
    expected = worker_ids(platform)
    if not isinstance(needs, dict) or set(needs) != set(expected):
        raise ValueError('missing or unexpected worker dependency')
    for key in expected:
        if not isinstance(needs[key], dict) or needs[key].get('result') != 'success':
            raise ValueError('required worker did not succeed: ' + key)
    return expected


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--platform', choices=('linux', 'windows'), required=True)
    args = parser.parse_args()
    try:
        workers = require_workers(os.environ.get('NEEDS_JSON', ''), args.platform)
    except (ValueError, TypeError):
        raise SystemExit('Required native gate rejected incomplete or unsuccessful workers')
    print('Required native gate verified: ' + ', '.join(workers))


if __name__ == '__main__':
    main()
