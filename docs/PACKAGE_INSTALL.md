# Development candidate installation

These archives are development candidates, not approved releases. They cover the currently exercised Linux/Windows amd64, CGO-disabled configuration. First-release feature scope, full licensing review, successful installed Azure scans and security/operations approval remain open.

Keep the supplied LICENSE, NOTICE.md, THIRD_PARTY_LICENSES.md and DEPENDENCY_NOTICES.md with the executable. PACKAGE_MANIFEST.json records source/tree, source pins and payload checksums; BUILD_INFO.json records compiled modules/settings. Archive SHA-256 detects corruption when compared with a trusted independently obtained checksum; it does not authenticate a publisher or replace signatures/provenance approval.

Verify the ZIP checksum before extraction:

```sh
sha256sum -c cloud-assess-dev-linux-amd64.zip.sha256
```

For Windows PowerShell, compare both the digest and expected ZIP filename against the supplied sidecar:

```powershell
(Get-FileHash .\cloud-assess-dev-windows-amd64.zip -Algorithm SHA256).Hash.ToLowerInvariant()
Get-Content .\cloud-assess-dev-windows-amd64.zip.sha256
```

Extract into a new directory. On Linux, preserve the executable permission or set it explicitly with `chmod 755 cloud-assess`; on Windows use cloud-assess.exe. Run from the extracted package directory:

```sh
./cloud-assess --version
./cloud-assess scan --help
```

```powershell
.\cloud-assess.exe --version
.\cloud-assess.exe scan --help
```

The executable needs no Go or Python runtime for these commands. Optional addition of the package directory to your own PATH is manual; this package does not edit system PATH, register services or configure Azure credentials. Azure scans require separately configured authentication and stage-specific permissions. No Azure scan is part of installation validation here.

For upgrade/rollback, retain the old package directory, unpack the new candidate separately and verify it before switching your chosen executable path. Do not overwrite evidence/reports or the retained previous package. Removing the extracted package directory removes these candidate files; authentication caches or assessment reports are separate operator-owned data.
