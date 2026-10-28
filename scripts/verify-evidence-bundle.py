#!/usr/bin/env python3
"""Verify a bounded synthetic-evidence bundle without extraction or execution."""
import argparse
import base64
import gzip
import hashlib
import io
import json
from pathlib import Path, PurePosixPath
import re
import sys
import zlib

FORMAT = 'cloud-assess-evidence-v1'
MAX_BUNDLE = 4 * 1024 * 1024
MAX_DECODED = 16 * 1024 * 1024
MAX_MANIFEST = 256 * 1024
MAX_MEMBERS = 128


def digest(data):
    return hashlib.sha256(data).hexdigest()


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError('duplicate JSON key')
        result[key] = value
    return result


def parse(data):
    return json.loads(data.decode('utf-8'), object_pairs_hook=unique_object)


def member_name(name):
    if not isinstance(name, str) or not re.fullmatch(r'[A-Za-z0-9_./-]{1,240}', name):
        raise ValueError('invalid member name')
    path = PurePosixPath(name)
    if not path.parts or path.is_absolute() or '..' in path.parts or str(path) != name or name.endswith('/'):
        raise ValueError('unsafe member name')
    return name


def verify(bundle, manifest_bytes):
    if not bundle or len(bundle) > MAX_BUNDLE or not manifest_bytes or len(manifest_bytes) > MAX_MANIFEST:
        raise ValueError('evidence input exceeds bounds or is empty')
    manifest = parse(manifest_bytes)
    if not isinstance(manifest, dict) or set(manifest) != {'format', 'classification', 'bundle_sha256', 'decoded_sha256', 'members', 'revisions'}:
        raise ValueError('invalid manifest shape')
    if manifest['format'] != FORMAT or manifest['classification'] != 'reviewed-synthetic-ci-local':
        raise ValueError('unsupported evidence contract')
    if digest(bundle) != manifest['bundle_sha256']:
        raise ValueError('bundle hash mismatch')
    compressed = base64.b64decode(bundle.strip(), validate=True)
    with gzip.GzipFile(fileobj=io.BytesIO(compressed)) as stream:
        decoded = stream.read(MAX_DECODED + 1)
    if len(decoded) > MAX_DECODED or digest(decoded) != manifest['decoded_sha256']:
        raise ValueError('decoded evidence exceeds bounds or hash mismatch')
    payload = parse(decoded)
    if not isinstance(payload, dict) or set(payload) != {'format', 'members'} or payload['format'] != FORMAT:
        raise ValueError('invalid payload shape')
    listed, actual = manifest['members'], payload['members']
    if not isinstance(listed, list) or not isinstance(actual, list) or not 1 <= len(listed) <= MAX_MEMBERS or len(listed) != len(actual):
        raise ValueError('invalid member inventory')
    revisions = manifest['revisions']
    if not isinstance(revisions, dict) or not revisions or any(not isinstance(k, str) or not isinstance(v, str) or not re.fullmatch(r'[0-9a-f]{40}', v) for k, v in revisions.items()):
        raise ValueError('invalid revision evidence')
    expected = {}
    for entry in listed:
        if not isinstance(entry, dict) or set(entry) != {'path', 'bytes', 'sha256'}:
            raise ValueError('invalid inventory entry')
        name = member_name(entry['path'])
        if name in expected or type(entry['bytes']) is not int or entry['bytes'] < 0 or not isinstance(entry['sha256'], str) or not re.fullmatch(r'[0-9a-f]{64}', entry['sha256']):
            raise ValueError('duplicate or invalid inventory entry')
        expected[name] = entry
    seen = set()
    for entry in actual:
        if not isinstance(entry, dict) or set(entry) != {'path', 'text'} or not isinstance(entry['text'], str):
            raise ValueError('invalid payload member')
        name = member_name(entry['path'])
        if name in seen or name not in expected:
            raise ValueError('duplicate or unexpected member')
        seen.add(name)
        raw = entry['text'].encode('utf-8')
        if len(raw) != expected[name]['bytes'] or digest(raw) != expected[name]['sha256']:
            raise ValueError('member length or hash mismatch')
    if seen != set(expected):
        raise ValueError('missing member')
    return {'members': len(seen), 'bytes': sum(e['bytes'] for e in listed), 'revisions': revisions}


def bounded_read(path, limit):
    path = Path(path)
    if path.is_symlink() or not path.is_file():
        raise ValueError('evidence input must be a regular non-symlink file')
    with path.open('rb') as stream:
        data = stream.read(limit + 1)
    if len(data) > limit:
        raise ValueError('evidence input exceeds bounds')
    return data


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--bundle', required=True)
    parser.add_argument('--manifest', required=True)
    args = parser.parse_args()
    try:
        result = verify(bounded_read(args.bundle, MAX_BUNDLE), bounded_read(args.manifest, MAX_MANIFEST))
    except (ValueError, TypeError, KeyError, OSError, EOFError, UnicodeError, RecursionError, zlib.error) as error:
        print('Evidence verification failed: ' + str(error), file=sys.stderr)
        return 1
    print('Evidence integrity verified: ' + json.dumps(result, sort_keys=True))
    print('Integrity only; no publisher authentication, privacy certification or QA acceptance is established.')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
