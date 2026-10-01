"""Validate fully resolved public profiles and reproduce Go canonical JSON bytes."""
import json
import re
import unicodedata
from urllib.parse import urlsplit

KEYS = ('schemaVersion', 'productName', 'cliName', 'reportTitle', 'reportFilePrefix', 'websiteURL')


def strict_json(data):
    def pairs(items):
        result = {}
        for key, value in items:
            if key in result:
                raise ValueError('duplicate JSON key')
            result[key] = value
        return result
    return json.loads(data.decode('utf-8'), object_pairs_hook=pairs)


def safe_component(value):
    return isinstance(value, str) and re.fullmatch(r'[a-z][a-z0-9_-]{0,63}', value) is not None and re.fullmatch(r'(con|prn|aux|nul|com[0-9]|lpt[0-9])', value, re.I) is None


def text(value, limit):
    return isinstance(value, str) and bool(value) and value.strip() == value and len(value.encode('utf-8')) <= limit and all(unicodedata.category(c) not in ('Cc', 'Cf') and c != '\ufffd' for c in value)


def validate(profile):
    if not isinstance(profile, dict) or set(profile) != set(KEYS) or type(profile['schemaVersion']) is not int or profile['schemaVersion'] != 1:
        raise ValueError('invalid resolved branding profile schema')
    if not text(profile['productName'], 256) or not text(profile['reportTitle'], 256):
        raise ValueError('invalid branding display text')
    if not safe_component(profile['cliName']) or not safe_component(profile['reportFilePrefix']):
        raise ValueError('invalid branding filename component')
    uri = profile['websiteURL']
    if not text(uri, 2048):
        raise ValueError('invalid branding URL')
    parsed = urlsplit(uri)
    if parsed.scheme != 'https' or not uri.startswith('https://') or not parsed.hostname or parsed.username is not None or parsed.password is not None:
        raise ValueError('invalid branding URL')
    if any(c.isspace() or c in '\\<>{}' for c in parsed.netloc) or re.search(r'%(?![0-9A-Fa-f]{2})', uri):
        raise ValueError('invalid branding URL')
    # Evaluate the port property to reject malformed supplied port syntax.
    parsed.port
    return profile


def canonical(profile):
    validate(profile)
    data = json.dumps({k: profile[k] for k in KEYS}, ensure_ascii=False, separators=(',', ':'))
    # encoding/json Marshal's default HTML-safe and Unicode line-separator escapes.
    for char, escape in [('&', '\\u0026'), ('<', '\\u003c'), ('>', '\\u003e'), ('\u2028', '\\u2028'), ('\u2029', '\\u2029')]:
        data = data.replace(char, escape)
    return data.encode('utf-8')


def parse(data):
    if len(data) > 16384:
        raise ValueError('branding profile exceeds size limit')
    profile = validate(strict_json(data))
    if data != canonical(profile):
        raise ValueError('branding profile is not canonical')
    return profile
