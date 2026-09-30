"""Execute the real publisher against isolated local Git repos, never GitHub/Azure."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

PUBLISHER = Path(__file__).resolve().parents[1] / 'publish-maintenance.sh'
CORE = 'bootstrap/core-v1'
BOT = 'github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>'


class PublicationTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.remote = self.root / 'remote.git'
        self.repo = self.root / 'repo'
        self.summary = self.root / 'summary.txt'
        self.env = os.environ.copy()
        self.env.update(GIT_CONFIG_GLOBAL=os.devnull, GIT_CONFIG_NOSYSTEM='1',
                        GIT_TERMINAL_PROMPT='0', GITHUB_RUN_ID='12345',
                        GITHUB_REPOSITORY='fixture/synthetic',
                        GITHUB_STEP_SUMMARY=str(self.summary),
                        GIT_AUTHOR_NAME='fixture human', GIT_AUTHOR_EMAIL='fixture@example.invalid',
                        GIT_COMMITTER_NAME='fixture human', GIT_COMMITTER_EMAIL='fixture@example.invalid')
        self.command('git', 'init', '--bare', str(self.remote), cwd=self.root)
        self.command('git', 'init', '-b', CORE, str(self.repo), cwd=self.root)
        for name in ['go.mod', 'go.sum', '.gitmodules', 'unrelated.txt',
                     'internal/rules/custom/rule.yml', 'internal/rules/upstream/rule.yml']:
            path = self.repo / name
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text('' if name == '.gitmodules' else 'baseline\n')
        self.git('add', '.')
        self.git('commit', '-m', 'synthetic baseline')
        self.git('remote', 'add', 'origin', str(self.remote))
        self.git('push', 'origin', CORE)
        self.base = self.git('rev-parse', 'HEAD').stdout.strip()

    def command(self, *args, cwd=None, check=True):
        return subprocess.run(args, cwd=cwd or self.repo, env=self.env,
                              text=True, capture_output=True, check=check, timeout=15)

    def git(self, *args, **kwargs):
        return self.command('git', *args, **kwargs)

    def publish(self, kind, check=True):
        return self.command('bash', str(PUBLISHER), kind, check=check)

    def refs(self):
        return self.git('--git-dir='+str(self.remote), 'for-each-ref',
                        '--format=%(refname:short)').stdout.splitlines()

    def assert_core_unchanged(self):
        self.assertEqual(self.git('rev-parse', CORE).stdout.strip(), self.base)
        self.assertEqual(self.git('--git-dir='+str(self.remote), 'rev-parse', CORE).stdout.strip(), self.base)

    def change(self, kind):
        path = 'go.mod' if kind == 'tidy' else 'internal/rules/custom/rule.yml'
        (self.repo / path).write_text('changed\n')
        return path

    def test_no_change_creates_no_branch_or_commit(self):
        for kind in ['tidy', 'rule-import']:
            with self.subTest(kind=kind):
                self.publish(kind)
                self.assertEqual(self.git('rev-parse', 'HEAD').stdout.strip(), self.base)
                self.assertEqual(self.refs(), [CORE])
                self.assertFalse(self.summary.exists())
                self.assertEqual(self.git('status', '--porcelain').stdout, '')

    def test_changed_output_is_scoped_bot_proposal(self):
        for kind in ['tidy', 'rule-import']:
            with self.subTest(kind=kind):
                path = self.change(kind)
                self.publish(kind)
                branch = 'maintenance/'+kind+'-12345'
                head = self.git('rev-parse', 'HEAD').stdout.strip()
                self.assertEqual(self.git('rev-parse', 'HEAD^').stdout.strip(), self.base)
                self.assertEqual(self.git('show', '-s', '--format=%an <%ae>%n%cn <%ce>').stdout.splitlines(), [BOT, BOT])
                self.assertEqual(self.git('diff-tree', '--no-commit-id', '--name-only', '-r', 'HEAD').stdout.splitlines(), [path])
                self.assertEqual(self.git('--git-dir='+str(self.remote), 'rev-parse', branch).stdout.strip(), head)
                self.assert_core_unchanged()
                self.assertIn('...'+branch+'?expand=1', self.summary.read_text())
                self.git('switch', CORE)

    def test_rule_proposal_retains_submodule_gitlink(self):
        child = self.root / 'source-child'
        self.command('git', 'init', str(child), cwd=self.root)
        (child / 'synthetic-rule.txt').write_text('synthetic rule\n')
        self.command('git', 'add', '.', cwd=child)
        self.command('git', 'commit', '-m', 'synthetic child', cwd=child)
        pin = self.command('git', 'rev-parse', 'HEAD', cwd=child).stdout.strip()
        self.git('-c', 'protocol.file.allow=always', 'submodule', 'add',
                 str(child), 'internal/rules/upstream/aprl')
        self.publish('rule-import')
        self.assertEqual(self.git('rev-parse', 'HEAD:internal/rules/upstream/aprl').stdout.strip(), pin)
        self.assertTrue(self.git('ls-tree', 'HEAD', 'internal/rules/upstream/aprl').stdout.startswith('160000 commit '))
        self.assertEqual(self.git('diff-tree', '--no-commit-id', '--name-only', '-r', 'HEAD').stdout.splitlines(),
                         ['.gitmodules', 'internal/rules/upstream/aprl'])
        self.assert_core_unchanged()

    def test_unrelated_staged_change_is_rejected(self):
        self.change('tidy')
        (self.repo / 'unrelated.txt').write_text('unexpected\n')
        self.git('add', 'unrelated.txt')
        result = self.publish('tidy', check=False)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('Unexpected staged path', result.stderr)
        self.assertEqual(self.refs(), [CORE])
        self.assert_core_unchanged()

    def test_wrong_branch_and_invalid_metadata_are_rejected(self):
        self.change('tidy')
        self.git('switch', '-c', 'wrong-branch')
        self.assertNotEqual(self.publish('tidy', check=False).returncode, 0)
        self.git('switch', CORE)
        self.env['GITHUB_RUN_ID'] = '../../bad'
        self.assertNotEqual(self.publish('tidy', check=False).returncode, 0)
        self.assertEqual(self.refs(), [CORE])
        self.assert_core_unchanged()

    def test_rejected_push_is_failure_without_success_summary(self):
        hook = self.remote / 'hooks' / 'pre-receive'
        hook.write_text('#!/bin/sh\nexit 1\n')
        hook.chmod(0o755)
        self.change('rule-import')
        result = self.publish('rule-import', check=False)
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(self.refs(), [CORE])
        self.assertFalse(self.summary.exists())
        self.assert_core_unchanged()

    def test_existing_remote_proposal_is_not_force_replaced(self):
        branch = 'maintenance/tidy-12345'
        self.git('switch', '-c', branch)
        (self.repo / 'unrelated.txt').write_text('remote proposal\n')
        self.git('add', 'unrelated.txt')
        self.git('commit', '-m', 'existing proposal')
        old_head = self.git('rev-parse', 'HEAD').stdout.strip()
        self.git('push', 'origin', branch)
        self.git('switch', CORE)
        self.git('branch', '-D', branch)
        self.change('tidy')
        self.assertNotEqual(self.publish('tidy', check=False).returncode, 0)
        self.assertEqual(self.git('--git-dir='+str(self.remote), 'rev-parse', branch).stdout.strip(), old_head)
        self.assertFalse(self.summary.exists())
        self.assert_core_unchanged()


if __name__ == '__main__':
    unittest.main(verbosity=2)
