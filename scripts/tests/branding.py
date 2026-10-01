#!/usr/bin/env python3
"""Exercise the supported branding builder and linked CLI offline on native CI."""
import hashlib
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
import os
from pathlib import Path
import subprocess
import tempfile
import threading

ROOT = Path(__file__).resolve().parents[2]
GO = os.environ.get('QA_GO', 'go')
PROFILE = {'schemaVersion': 1, 'productName': 'Example Cloud $(touch unexpected-file)',
           'cliName': 'example-cloud', 'reportTitle': 'Example "Cloud" Assessment',
           'reportFilePrefix': 'example_report', 'websiteURL': 'https://example.test/tool'}


def execute(args, cwd=ROOT, env=None, timeout=120):
    return subprocess.run(args, cwd=cwd, env=env, capture_output=True, text=True, timeout=timeout)


def success(result):
    if result.returncode:
        raise AssertionError(result.stdout + result.stderr)
    return result.stdout


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
        with tempfile.TemporaryDirectory(prefix='cloud-assess-branding-') as tmp:
            tmp = Path(tmp)
            profile = tmp / 'profile with spaces.json'
            profile.write_text(json.dumps(PROFILE), encoding='utf-8')
            builder = [GO, 'run', './tools/brand-build']
            default = Path(success(execute(builder + ['--output', str(tmp / 'default')])).strip()).resolve()
            custom = Path(success(execute(builder + ['--profile', str(profile), '--output', str(tmp / 'custom with spaces')])).strip()).resolve()
            if custom.name != 'example-cloud' + ('.exe' if os.name == 'nt' else ''):
                raise AssertionError('custom executable basename mismatch')
            env = {k: v for k, v in os.environ.items() if not k.upper().startswith('AZURE_')}
            endpoint = 'http://127.0.0.1:' + str(server.server_address[1])
            env.update({'HTTP_PROXY': endpoint, 'HTTPS_PROXY': endpoint, 'http_proxy': endpoint,
                        'https_proxy': endpoint, 'NO_PROXY': '', 'no_proxy': '', 'AZURE_CONFIG_DIR': str(tmp / 'no-azure'),
                        'MSI_ENDPOINT': endpoint, 'IDENTITY_ENDPOINT': endpoint, 'IDENTITY_HEADER': 'offline-tripwire'})
            default_profile = json.loads(success(execute([str(default), 'branding'], tmp, env)))
            if default_profile['productName'] != 'Cloud Assess' or default_profile['cliName'] != 'cloud-assess':
                raise AssertionError('default build identity changed')
            if json.loads(success(execute([str(custom), 'branding'], tmp, env))) != PROFILE:
                raise AssertionError('linked canonical profile mismatch')
            for binary, name, product in [(default, 'cloud-assess', 'Cloud Assess'), (custom, 'example-cloud', PROFILE['productName'])]:
                if success(execute([str(binary), '--version'], tmp, env)).strip() != name + ' version dev':
                    raise AssertionError('version identity mismatch')
                help_text = success(execute([str(binary), '--help'], tmp, env))
                if name not in help_text or product not in help_text:
                    raise AssertionError('help identity mismatch')
                rejected = execute([str(binary), 'scan', '--assessment-timeout=-1s'], tmp, env)
                if rejected.returncode != 1 or rejected.stdout or 'assessment timeout cannot be negative' not in rejected.stderr:
                    raise AssertionError('custom preflight contract changed')
            # Existing executable remains byte-identical after a rejected rebuild.
            digest = hashlib.sha256(custom.read_bytes()).hexdigest()
            rejected = execute(builder + ['--profile', str(profile), '--output', str(custom.parent)])
            if rejected.returncode == 0 or hashlib.sha256(custom.read_bytes()).hexdigest() != digest:
                raise AssertionError('builder replaced existing executable')
            # Invalid profiles cannot create the destination or run a build.
            for index, data in enumerate(['{"schemaVersion":null}', '{"schemaVersion":1,"cliName":"../escape"}',
                                          '{"schemaVersion":1,"cliName":"safe","cliName":"other"}',
                                          '{"schemaVersion":1,"companyName":"unused"}']):
                profile.write_text(data, encoding='utf-8')
                output = tmp / ('invalid-' + str(index))
                rejected = execute(builder + ['--profile', str(profile), '--output', str(output)])
                if rejected.returncode == 0 or output.exists():
                    raise AssertionError('invalid profile reached build/output stage')
            malformed = tmp / ('malformed' + ('.exe' if os.name == 'nt' else ''))
            build_env = dict(os.environ, CGO_ENABLED='0')
            success(execute([GO, 'build', '-trimpath', '-buildvcs=true', '-ldflags',
                             '-X github.com/DeBoX85/Cloud-Assess/internal/branding.embeddedProfile=e30',
                             '-o', str(malformed), './cmd/cloud-assess'], env=build_env))
            retained = tmp / 'existing.json'
            retained.write_bytes(b'retained report bytes')
            for args in [['--help'], ['branding'], ['scan', '--output-name', str(tmp / 'existing'), '--json']]:
                rejected = execute([str(malformed)] + args, tmp, env)
                if rejected.returncode != 1 or rejected.stdout or 'invalid embedded branding' not in rejected.stderr:
                    raise AssertionError('malformed embedded state did not fail closed')
                if retained.read_bytes() != b'retained report bytes':
                    raise AssertionError('malformed embedded state replaced output')
            if traffic or (tmp / 'no-azure').exists() or (tmp / 'unexpected-file').exists() or (ROOT / 'unexpected-file').exists():
                raise AssertionError('branding/preflight attempted authentication or shell execution')
            print('Branding builder/CLI: default and custom native binaries, canonical identity, preservation and malformed-state tripwire checks passed')
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=5)


if __name__ == '__main__':
    main()
