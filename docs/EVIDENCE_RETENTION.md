# Development evidence retention and verification

[CURRENT_STATE](CURRENT_STATE.md) owns status; [action register](ACTION_REGISTER.md) tracks gaps. This procedure preserves reviewable development evidence, not regulatory immutability or release approval.

## Evidence classes and ownership

| Evidence | Location and owner | Required treatment |
| --- | --- | --- |
| Code/contracts and compact outcomes | Repository history and task PR; working agent | Publish/read back exact revisions; PR body is final QA index; preserve earlier snapshots. |
| Reviewed synthetic CI and local rehearsal | Versioned bundle/manifest under docs/evidence; working agent | Review sensitivity before publication; preserve complete decoded logs and run/job/step metadata, checksums, revision provenance and verification method. |
| Raw Azure/customer/private evidence | Approved private storage; authorized operator | Never place in public Git or public artifacts. Record only safe references/counts/hashes. Actual access/retention must be approved and verified. |
| Off-repository disaster recovery | Durable independent backup owner/location not yet established | DC-006 remains open until retention, access and a restore exercise are recorded. Repository history alone is not independent backup. |

GitHub's [October1 retention change](https://github.blog/changelog/2026-10-01-actions-retention-now-covers-checks-runs-and-statuses/) applies retention to checks, runs and statuses as well as logs/artifacts; public repositories have at most90days. Current product workflow explicitly retains its conditional synthetic failure artifact for14days. Actual repository-wide retention setting has not been verified through the available connection. Do not equate workflow artifact configuration with the log/run policy or infer unavailable administrator settings. Archive essential proof promptly rather than relying on expiring run URLs.

## Retained checkpoint

[Manifest](evidence/20261017/manifest.json) indexes [the base64 gzip JSON bundle](evidence/20261017/pr150-pr151.json.gz.base64). It retains complete connector-decoded logs for both PR150/151, required comparison logs from PR149, raw run/job metadata for PR150/151, per-revision review/source reconstruction summaries and local rollback commands/script. Scope is synthetic CI/local Git only; no customer Azure operations occurred. Common-token-pattern screening is a limited sanity check, not a general privacy certification. Original tool-returned decoded text is preserved, including Unicode/BOM where present; do not claim original transport archive hashes.

The base64 bundle is marked -text in .gitattributes so Windows checkout cannot change its hash through newline conversion. Do not normalize archived bytes.

The bundle contains44members; its manifest binds every UTF-8 member length/hash, compressed text hash and decoded JSON hash to explicit PR149/150/151 revisions. It is526101bytes, SHA256 df421cc9b97e26e218ac8ff4a16c77e44bec22fc02d69c295cf5da7ec08ced13. The local review scripts retain historical absolute workspace paths and document the executed method; adjust paths in a disposable copy for replay. The standalone integrity verifier below does not require those paths or a tool store.

## Verify without extraction

From a checkout containing these files, run:

```sh
python scripts/verify-evidence-bundle.py --bundle docs/evidence/20261017/pr150-pr151.json.gz.base64 --manifest docs/evidence/20261017/manifest.json
```

The verifier reads only the explicit files, bounds input/decompression, rejects duplicate JSON/member names, unsafe member paths and wrong inventories/lengths/hashes, and never extracts or executes members. Its tests include actual failing CLI exit checks and a valid committed-bundle check. Valid integrity does not mean contained commands were safe, tests passed, privacy is certified or the publisher is authenticated. Anyone able to replace both bundle and manifest can recompute hashes; remote Git/review provenance and backups remain separate controls.

For a later checkpoint use new dated filenames, record the responsible reviewer, source revision/run/job IDs and evidence class, then verify full readback. Do not overwrite an earlier checkpoint to match newer outcomes. Keep the actual evidence and manifest together; a hash without contents is not recovery. Archive only material completed evidence, not every transient tool call. Repeat full exact-head QA when code/docs change; this retained predecessor bundle does not certify its successor.
