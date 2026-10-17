# PR153 final evidence snapshot

Artifact-only retention branch. This is not a feature successor, proposed protected merge, release or code certification for this different artifact tree. Continue from PR153's code candidate ccfee397c4c12da80b12d0b8c5d0043836d3f543, tree24089ad9240fa15eb41cdff37195190b838fe262, parent0067731b73cef1f0a9e201b680c5c85fc38f6672 and its live final evidence index.

[Bundle](pr153-final.json.gz.base64) and [manifest](pr153-final.manifest.json) retain 27 members / 2220745 UTF-8 member bytes: complete decoded native/control/gate/source logs and raw run/job metadata, predecessor comparison logs, source/control reconstruction review, local manifest structural/drift/restoration record and recovered review script. Bundle 311289 bytes; SHA256 0ad4871cc186c4ca51686fab6a0cab97f1e4f722435344ae9aa13f6bb8efdb98. Integrity verified with the existing bounded candidate verifier before publication; all artifact files read back after publication. Source log fingerprints are connector-decoded text, not original transport download bytes.

Verify without extraction:

```sh
python scripts/verify-evidence-bundle.py --bundle docs/evidence/20261017/pr153-final.json.gz.base64 --manifest docs/evidence/20261017/pr153-final.manifest.json
```

All385Linux/380Windows controls retained, all23source captures/315chunks/909865bytes matched. The new61-case oracle is future API preparation, not executed catalog scenarios. Its integrity guard passed on both native hosts; the local deliberately changed byte was rejected after compilation and exact restoration passed. No VMSS production API or scoring enablement; no Azure/live/protected/release proof.

The earlier temporary QA helper disappeared during development. It and predecessor logs were restored from the verified committed44-member archive, then the helper was adapted to this revision. Scripts retain historical absolute paths; replay requires a disposable workspace and explicit path adjustment. This demonstrates recovery from published evidence, not independent backup or immutable audit storage. Keep private Azure evidence out of this branch. Existing frozen heads/pins/captures remain unchanged.
