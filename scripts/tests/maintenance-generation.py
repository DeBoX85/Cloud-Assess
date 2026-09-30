"""Run workflow generation recipes with real Git/Go in disposable repositories.

The pinned reference archive is downloaded once (or supplied locally); the curl
fixture replays it. APRL cloning uses the checked-out pinned submodule as a local
Git mirror, retaining the production URL in .gitmodules. No GitHub refs are pushed.
"""
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import tarfile
import textwrap
import unittest

ROOT = Path(__file__).resolve().parents[2]
CORE = 'bootstrap/core-v1'
REFERENCE = '8e4f0577f3615e6c9014c031bcad079f235369cc'
APRL = '60eaddda76541f6adbc1c5ffa686829807e55e29'
APRL_URL = 'https://github.com/Azure/Azure-Proactive-Resiliency-Library-v2.git'
GO = os.environ.get('QA_GO', shutil.which('go'))


def recipe(workflow, name):
    """Extract the actual named literal Bash block, failing if its shape changes."""
    lines = (ROOT / '.github/workflows' / workflow).read_text().splitlines()
    start = lines.index('      - name: '+name)
    end = next((i for i in range(start+1, len(lines)) if lines[i].startswith('      - ')), len(lines))
    run = next(i for i in range(start+1, end) if lines[i].startswith('        run:'))
    if lines[run] != '        run: |':
        return lines[run].removeprefix('        run: ')
    block = []
    for line in lines[run+1:end]:
        if line and not line.startswith('          '):
            break
        block.append(line)
    return textwrap.dedent('\n'.join(block))


class GenerationTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.shared = tempfile.TemporaryDirectory()
        cls.archive = Path(cls.shared.name) / 'reference.tar.gz'
        supplied = os.environ.get('QA_REFERENCE_ARCHIVE')
        if supplied:
            shutil.copyfile(supplied, cls.archive)
        else:
            subprocess.run(['curl', '--fail', '--location', '--retry', '3',
                            'https://codeload.github.com/DeBoX85/azqr/tar.gz/'+REFERENCE,
                            '--output', str(cls.archive)], check=True, timeout=120)
        cls.aprl_source = ROOT / 'internal/rules/upstream/aprl'
        pin = subprocess.check_output(['git', '-C', str(cls.aprl_source), 'rev-parse', 'HEAD'], text=True).strip()
        if pin != APRL:
            raise AssertionError('APRL source mirror must be initialized at the exact pin')
        if not GO:
            raise AssertionError('Go toolchain required')

    @classmethod
    def tearDownClass(cls):
        cls.shared.cleanup()

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.home = Path(self.temp.name)
        self.repo = self.home / 'repo'
        self.remote = self.home / 'remote.git'
        self.summary = self.home / 'summary'
        self.env = os.environ.copy()
        self.env.update(GIT_CONFIG_GLOBAL=os.devnull, GIT_CONFIG_NOSYSTEM='1',
                        GIT_TERMINAL_PROMPT='0', GITHUB_RUN_ID='23456',
                        GITHUB_REPOSITORY='fixture/synthetic', GITHUB_STEP_SUMMARY=str(self.summary),
                        GIT_AUTHOR_NAME='fixture', GIT_AUTHOR_EMAIL='fixture@example.invalid',
                        GIT_COMMITTER_NAME='fixture', GIT_COMMITTER_EMAIL='fixture@example.invalid',
                        GIT_CONFIG_COUNT='2', GIT_CONFIG_KEY_0='protocol.file.allow', GIT_CONFIG_VALUE_0='always',
                        GIT_CONFIG_KEY_1='url.'+str(self.aprl_source)+'.insteadOf', GIT_CONFIG_VALUE_1=APRL_URL,
                        QA_REPLAY_ARCHIVE=str(self.archive))
        self.command('git', 'clone', '--depth', '1', '--no-checkout', ROOT.as_uri(), str(self.repo), cwd=self.home)
        self.git('checkout', '-B', CORE, 'HEAD')
        self.git('remote', 'remove', 'origin')
        self.command('git', 'init', '--bare', str(self.remote), cwd=self.home)
        # The synthetic server must accept the depth-one fixture baseline.
        # This setting applies only to the disposable local bare remote.
        self.git('--git-dir='+str(self.remote), 'config', 'receive.shallowUpdate', 'true')
        self.git('remote', 'add', 'origin', str(self.remote))
        self.git('push', 'origin', CORE)
        self.base = self.git('rev-parse', 'HEAD').stdout.strip()
        shim = self.home / 'bin'
        shim.mkdir()
        curl = shim / 'curl'
        curl.write_text('#!/usr/bin/env python3\nimport os, shutil, sys\n'
                        'assert "https://codeload.github.com/DeBoX85/azqr/tar.gz/'+REFERENCE+'" in sys.argv\n'
                        'shutil.copyfile(os.environ["QA_REPLAY_ARCHIVE"], sys.argv[sys.argv.index("--output")+1])\n')
        curl.chmod(0o755)
        self.env['PATH'] = str(shim)+os.pathsep+self.env['PATH']

    def command(self, *args, cwd=None, check=True):
        try:
            return subprocess.run(args, cwd=cwd or self.repo, env=self.env,
                                  text=True, capture_output=True, check=check, timeout=120)
        except subprocess.CalledProcessError as error:
            error.add_note('Fixture command stderr:\n'+(error.stderr or '')[-4000:])
            raise

    def git(self, *args, **kwargs):
        return self.command('git', *args, **kwargs)

    def generate_rules(self, check=True):
        return self.command('bash', '-c', recipe('import-rule-baseline.yml', 'Import pinned reference rule snapshots'), check=check)

    def publish(self, kind):
        return self.command('bash', str(ROOT / 'scripts/publish-maintenance.sh'), kind)

    def assert_rules_match_pins(self):
        tree = self.git('write-tree').stdout.strip()
        for path, expected in [('internal/rules/custom', '674b9b3dcb443ce6dc445b48e1db47e4a0ca7082'),
                               ('internal/rules/upstream/orphan-resources', 'a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5'),
                               ('internal/rules/upstream/aprl', APRL)]:
            self.assertEqual(self.git('rev-parse', tree+':'+path).stdout.strip(), expected)

    def assert_core_unchanged(self):
        self.assertEqual(self.git('rev-parse', CORE).stdout.strip(), self.base)
        self.assertEqual(self.git('--git-dir='+str(self.remote), 'rev-parse', CORE).stdout.strip(), self.base)

    def test_pinned_rule_regeneration_is_no_change(self):
        self.generate_rules()
        self.git('add', '.gitmodules', 'internal/rules/custom', 'internal/rules/upstream')
        self.assert_rules_match_pins()
        self.assertEqual(self.git('diff', '--cached', '--name-only').stdout, '')
        self.publish('rule-import')
        self.assertFalse(self.summary.exists())
        self.assertEqual(self.git('rev-parse', 'HEAD').stdout.strip(), self.base)
        self.assert_core_unchanged()

    def test_corrupted_rule_regenerates_scoped_proposal(self):
        target = sorted((self.repo / 'internal/rules/custom/azure-resources').rglob('*.yaml'))[0]
        relative = str(target.relative_to(self.repo))
        target.write_text('synthetic corrupted rule\n')
        self.git('add', relative)
        self.git('commit', '-m', 'synthetic corruption baseline')
        self.git('push', 'origin', CORE)
        self.base = self.git('rev-parse', 'HEAD').stdout.strip()
        self.generate_rules()
        self.publish('rule-import')
        self.assert_rules_match_pins()
        self.assertEqual(self.git('diff-tree', '--no-commit-id', '--name-only', '-r', 'HEAD').stdout.splitlines(), [relative])
        self.assert_core_unchanged()

    def test_bad_archive_fails_before_library_mutation(self):
        bad = self.home / 'bad.tar.gz'
        bad.write_text('not an archive')
        self.env['QA_REPLAY_ARCHIVE'] = str(bad)
        self.assertNotEqual(self.generate_rules(check=False).returncode, 0)
        self.assertEqual(self.git('status', '--porcelain').stdout, '')
        self.assertFalse(self.summary.exists())
        self.assert_core_unchanged()

    def test_missing_source_layout_fails_before_library_mutation(self):
        incomplete = self.home / 'incomplete.tar.gz'
        source = self.home / 'source'
        source.mkdir()
        (source / 'marker').write_text('synthetic incomplete snapshot\n')
        with tarfile.open(incomplete, 'w:gz') as archive:
            archive.add(source, arcname='fixture')
        self.env['QA_REPLAY_ARCHIVE'] = str(incomplete)
        self.assertNotEqual(self.generate_rules(check=False).returncode, 0)
        self.assertEqual(self.git('status', '--porcelain').stdout, '')
        self.assertFalse(self.summary.exists())
        self.assert_core_unchanged()

    def test_wrong_generated_snapshot_is_rejected_before_publication(self):
        wrong = self.home / 'wrong.tar.gz'
        source = self.home / 'wrong-source'
        for name in ['internal/graph/azqr/azure-resources', 'internal/graph/azure-orphan-resources']:
            path = source / name
            path.mkdir(parents=True)
            (path / 'synthetic.yaml').write_text('wrong snapshot\n')
        with tarfile.open(wrong, 'w:gz') as archive:
            archive.add(source, arcname='fixture')
        self.env['QA_REPLAY_ARCHIVE'] = str(wrong)
        self.assertNotEqual(self.generate_rules(check=False).returncode, 0)
        self.assertEqual(self.git('rev-parse', 'HEAD').stdout.strip(), self.base)
        self.assertFalse(self.summary.exists())
        self.assert_core_unchanged()

    def test_actual_tidy_removes_unused_requirement_and_then_stabilizes(self):
        # Replace only fixture files; execute the real workflow tidy command.
        dependency = self.home / 'unused'
        dependency.mkdir()
        (dependency / 'go.mod').write_text('module fixture.invalid/unused\n\ngo 1.26.0\n')
        (dependency / 'unused.go').write_text('package unused\n')
        fixture = self.home / 'module'
        self.command('git', 'init', '-b', CORE, str(fixture), cwd=self.home)
        self.repo = fixture
        (fixture / 'go.mod').write_text('module fixture.invalid/main\n\ngo 1.26.0\n\nrequire fixture.invalid/unused v0.0.0\n\nreplace fixture.invalid/unused => '+str(dependency)+'\n')
        (fixture / 'go.sum').write_text('')
        (fixture / 'main.go').write_text('package main\nfunc main() {}\n')
        self.git('add', '.')
        self.git('commit', '-m', 'synthetic untidy module')
        self.remote = self.home / 'module-remote.git'
        self.command('git', 'init', '--bare', str(self.remote), cwd=self.home)
        self.git('remote', 'add', 'origin', str(self.remote))
        self.git('push', 'origin', CORE)
        self.base = self.git('rev-parse', 'HEAD').stdout.strip()
        self.env['PATH'] = str(Path(GO).parent)+os.pathsep+self.env['PATH']
        self.env.update(GOPROXY='off', GOSUMDB='off', GOTOOLCHAIN='local')
        self.command('bash', '-c', recipe('tidy.yml', 'Tidy dependencies'))
        self.assertNotIn('require fixture.invalid/unused', (fixture / 'go.mod').read_text())
        self.publish('tidy')
        self.assert_core_unchanged()
        self.assertEqual(self.git('diff-tree', '--no-commit-id', '--name-only', '-r', 'HEAD').stdout.splitlines(), ['go.mod'])
        head = self.git('rev-parse', 'HEAD').stdout.strip()
        self.command('bash', '-c', recipe('tidy.yml', 'Tidy dependencies'))
        self.assertEqual(self.git('status', '--porcelain').stdout, '')
        self.assertEqual(self.git('rev-parse', 'HEAD').stdout.strip(), head)


if __name__ == '__main__':
    unittest.main(verbosity=2)
