#!/usr/bin/env python3
"""Compare real application report artifacts under two immutable linked profiles.

Synthetic assessment input, production application/renderers, no Azure scan.
"""
import base64
import copy
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
import os
from pathlib import Path
import subprocess
import tempfile
import threading
import xml.etree.ElementTree as ET
import zipfile

ROOT = Path(__file__).resolve().parents[2]
GO = os.environ.get('QA_GO', 'go')
SUB = '11111111-2222-3333-4444-555555555555'
DEFAULT = {'schemaVersion': 1, 'productName': 'Cloud Assess', 'cliName': 'cloud-assess',
           'reportTitle': 'Azure Cloud Assessment', 'reportFilePrefix': 'cloud_assessment',
           'websiteURL': 'https://github.com/DeBoX85/Cloud-Assess'}
CUSTOM = {'schemaVersion': 1, 'productName': 'Synthetic Cloud Toolkit', 'cliName': 'synthetic-cloud',
          'reportTitle': '=1+1 & Synthetic "Assessment"', 'reportFilePrefix': 'synthetic_report',
          'websiteURL': 'https://example.test/custom'}
NS = {'s': 'http://schemas.openxmlformats.org/spreadsheetml/2006/main'}


def run(args, env=None):
    result = subprocess.run(args, cwd=ROOT, env=env, capture_output=True, text=True, timeout=180)
    if result.returncode:
        raise AssertionError(result.stdout + result.stderr)
    return result.stdout


def workbook(path, title):
    """Read actual XLSX cells/formulas independently of the Go renderer/tests."""
    with zipfile.ZipFile(path) as archive:
        shared = []
        if 'xl/sharedStrings.xml' in archive.namelist():
            shared = [''.join(x.itertext()) for x in ET.fromstring(archive.read('xl/sharedStrings.xml'))]
        names = [s.attrib['name'] for s in ET.fromstring(archive.read('xl/workbook.xml')).findall('s:sheets/s:sheet', NS)]
        relationships = {r.attrib['Id']: r.attrib['Target'] for r in ET.fromstring(archive.read('xl/_rels/workbook.xml.rels'))}
        sheets = ET.fromstring(archive.read('xl/workbook.xml')).findall('s:sheets/s:sheet', NS)
        output = {}
        for sheet in sheets:
            target = relationships[sheet.attrib['{http://schemas.openxmlformats.org/officeDocument/2006/relationships}id']]
            target = target.lstrip('/') if target.startswith('/') else 'xl/' + target
            xml = ET.fromstring(archive.read(target))
            cells = {}
            for cell in xml.findall('s:sheetData/s:row/s:c', NS):
                value = cell.find('s:v', NS)
                text = '' if value is None else (value.text or '')
                if cell.attrib.get('t') == 's':
                    text = shared[int(text)]
                elif cell.attrib.get('t') == 'inlineStr':
                    text = ''.join(x.text or '' for x in cell.findall('.//s:t', NS))
                formula = cell.find('s:f', NS)
                cells[cell.attrib['r']] = [text, None if formula is None else formula.text]
            if cells.get('A1') != [title, None]:
                raise AssertionError('missing/wrong/formula branding title in ' + sheet.attrib['name'])
            # Only this explicit presentation cell may differ.
            cells['A1'] = ['<presentation-title>', None]
            output[sheet.attrib['name']] = cells
        if len(names) != 12 or any(not any(k.endswith('5') for k in c) for c in output.values()):
            raise AssertionError('fixture did not exercise every non-empty report sheet')
        return output


def artifacts(directory, profile):
    output = {}
    for case, base in [('default', profile['reportFilePrefix'] + '_2026_10_01_T120000'),
                       ('explicit', 'operator override'), ('raw', 'raw')]:
        files = [p for p in directory.glob(base + '.*') if not p.name.endswith('.stdout.json')]
        if len(files) != 15:
            raise AssertionError('unexpected artifact set for ' + case)
        report = json.loads((directory / (base + '.json')).read_text())
        stdout = json.loads((directory / (case + '.stdout.json')).read_text())
        if report != stdout:
            raise AssertionError('stdout changed assessment projection')
        expected_ids = {'synthetic-aprl', 'synthetic-aor', 'synthetic-custom'}
        if {r['source'] for r in report['recommendations']} != {'APRL', 'AOR', 'CUSTOM'} or {r['recommendationId'] for r in report['findings']} != expected_ids:
            raise AssertionError('upstream recommendation identity/provenance changed')
        if case == 'raw' and report['resources'][0]['subscriptionId'] != SUB:
            raise AssertionError('raw identity missing')
        if case != 'raw' and SUB in json.dumps(report):
            raise AssertionError('redacted canonical data leaked raw subscription ID')
        sarif = json.loads((directory / (base + '.sarif')).read_text())
        driver = sarif['runs'][0]['tool']['driver']
        if driver['name'] != profile['cliName'] or driver['informationUri'] != profile['websiteURL']:
            raise AssertionError('SARIF profile metadata mismatch')
        if len(sarif['runs'][0]['results']) != 3 or not all(r.get('partialFingerprints') for r in sarif['runs'][0]['results']):
            raise AssertionError('SARIF fixture lacks stable result identities')
        # Permit exactly these two known presentation fields to differ.
        driver['name'], driver['informationUri'] = '<cli-name>', '<website-url>'
        csv = {p.name[len(base):]: p.read_bytes() for p in files if p.suffix == '.csv'}
        if len(csv) != 12:
            raise AssertionError('CSV fixture coverage incomplete')
        output[case] = {'json': (directory / (base + '.json')).read_bytes(),
                        'stdout': (directory / (case + '.stdout.json')).read_bytes(),
                        'csv': csv, 'sarif': sarif, 'xlsx': workbook(directory / (base + '.xlsx'), profile['reportTitle'])}
    for i in range(4):
        base = 'concurrent-' + str(i)
        if (directory / (base + '.json')).read_bytes() != output['explicit']['json']:
            raise AssertionError('concurrent rendering mixed assessment data')
        if workbook(directory / (base + '.xlsx'), profile['reportTitle']) != output['explicit']['xlsx']:
            raise AssertionError('concurrent rendering mixed worksheet identity/data')
        concurrent = json.loads((directory / (base + '.sarif')).read_text())
        driver = concurrent['runs'][0]['tool']['driver']
        if driver['name'] != profile['cliName'] or driver['informationUri'] != profile['websiteURL']:
            raise AssertionError('concurrent SARIF identity changed')
        driver['name'], driver['informationUri'] = '<cli-name>', '<website-url>'
        if concurrent != output['explicit']['sarif']:
            raise AssertionError('concurrent SARIF record identities changed')
    return output


def compare(default, custom):
    if default != custom:
        raise AssertionError('branding changed canonical data, report rows or stable SARIF identity')


def check_negative_controls(default, custom):
    # These alter independently read evidence after permitted presentation fields
    # are normalized. The same comparison used for acceptance must reject each.
    for field in ('json', 'stdout', 'csv', 'xlsx', 'sarif'):
        changed = copy.deepcopy(custom)
        if field in ('json', 'stdout'):
            changed['explicit'][field] += b'changed assessment data'
        elif field == 'csv':
            key = next(iter(changed['explicit']['csv']))
            changed['explicit']['csv'][key] += b'changed row'
        elif field == 'xlsx':
            changed['explicit']['xlsx']['Recommendations']['A5'][0] = 'changed row'
        else:
            changed['explicit']['sarif']['runs'][0]['automationDetails']['id'] = 'changed-identity'
        try:
            compare(default, changed)
        except AssertionError:
            continue
        raise AssertionError('comparison accepted ' + field + ' mutation')


def main():
    traffic = []
    class Tripwire(BaseHTTPRequestHandler):
        def do_GET(self):
            traffic.append(self.path)
            self.send_error(503)
        do_POST = do_GET
        do_CONNECT = do_GET
        def log_message(self, *_):
            pass
    server = ThreadingHTTPServer(('127.0.0.1', 0), Tripwire)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        with tempfile.TemporaryDirectory(prefix='cloud-assess-brand-reports-') as tmp:
            tmp = Path(tmp)
            env = {k: v for k, v in os.environ.items() if not k.upper().startswith('AZURE_')}
            endpoint = 'http://127.0.0.1:' + str(server.server_address[1])
            env.update({'HTTP_PROXY': endpoint, 'HTTPS_PROXY': endpoint, 'http_proxy': endpoint,
                        'https_proxy': endpoint, 'NO_PROXY': '', 'no_proxy': '', 'AZURE_CONFIG_DIR': str(tmp / 'no-azure'),
                        'MSI_ENDPOINT': endpoint, 'IDENTITY_ENDPOINT': endpoint, 'IDENTITY_HEADER': 'offline-tripwire'})
            paired = []
            for name, profile in [('default', DEFAULT), ('custom', CUSTOM)]:
                binary = tmp / (name + ('.exe' if os.name == 'nt' else ''))
                args = [GO, 'test', '-c', '-o', str(binary)]
                if os.name != 'nt':
                    args += ['-race']
                if name == 'custom':
                    encoded = base64.urlsafe_b64encode(json.dumps(profile, separators=(',', ':')).encode()).decode().rstrip('=')
                    args += ['-ldflags', '-X github.com/DeBoX85/Cloud-Assess/internal/branding.embeddedProfile=' + encoded]
                run(args + ['./internal/app'])
                directory = tmp / (name + ' reports with spaces')
                run([str(binary), '-test.run=^TestBrandingReportEvidence$', '-test.count=1',
                     '-branding-evidence-dir', str(directory)], env)
                paired.append(artifacts(directory, profile))
            compare(paired[0], paired[1])
            check_negative_controls(paired[0], paired[1])
            if traffic or (tmp / 'no-azure').exists():
                raise AssertionError('synthetic report checks attempted HTTP/authentication')
            print('Branding reports: 12 XLSX sheets/12 CSV tables, JSON/stdout, SARIF records/provenance, raw/redacted paths and concurrent isolation passed')
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=5)


if __name__ == '__main__':
    main()
