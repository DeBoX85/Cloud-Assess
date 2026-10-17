# Notices and Attribution

Cloud Assess contains code, recommendation metadata, queries, or behavior derived from third-party open-source projects.

The repository-level license for original Cloud Assess material is the Apache License 2.0. It does not replace license obligations for incorporated third-party material.

The complete license texts for the incorporated source/rule families currently known to the project are preserved in [THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md).

## Reference implementation

Portions of the implementation are derived from or behaviorally based on Azure Quick Review source code from Microsoft Corporation, licensed under the MIT License.

Reference baseline: `DeBoX85/azqr` commit `8e4f0577f3615e6c9014c031bcad079f235369cc`.

Copyright (c) Microsoft Corporation.

## Azure Proactive Resiliency Library v2

Cloud Assess incorporates recommendation material from the Azure Proactive Resiliency Library v2 (APRL), a Microsoft project distributed under the MIT License.

Pinned APRL revision: `60eaddda76541f6adbc1c5ffa686829807e55e29`.

Copyright (c) Microsoft Corporation.

## Azure Orphan Resources

Cloud Assess incorporates recommendation material derived from the Azure Orphan Resources project, distributed under the MIT License.

Copyright (c) 2023 Dolev Shor.

The imported snapshot revision is recorded in `internal/rules/provenance.go`.

## Generated dependency notices

The source/rule license texts above do not replace a dependency inventory. The current [generated inventory and notices](docs/DEPENDENCIES.md) are checked for freshness in CI. Before any release or redistribution, regenerate from the actual release tree and dependency graph, review platform/artifact contents and these attributions for completeness, and include the required license texts with the distribution.

## Azure REST API specifications test fixture

The Storage SKU primary example in internal/plugins/region/testdata/storage-primary-2023-01-01-SKUList.json is copied unchanged from Azure/azure-rest-api-specs commit `1cdcc5414764d20e8815554c8b89708c10747c18`.

Copyright (c) 2017 Microsoft. Licensed under the MIT License. The complete original notice and permission text are retained in [storage-primary-LICENSE.txt](internal/plugins/region/testdata/storage-primary-LICENSE.txt). This fixture is primary schema/example evidence, separate from the pinned AZQR source captures.

The Compute Resource SKUs schema and three complete primary examples in internal/plugins/region/testdata/compute-primary-2021-07-01-*.json are also copied unchanged from the same Azure REST API specifications commit. Their original Git blobs and SHA256 hashes are recorded in docs/REGION_VM_VMSS_SKU_CONTRACT_20261016.md. The same complete MIT license in storage-primary-LICENSE.txt applies. These are primary contract evidence, separate from immutable AZQR captures.

The Cognitive Services 2024-10-01 complete schema and GetSkus example in internal/plugins/region/testdata/cognitive-primary-2024-10-01-*.json are copied unchanged from the same Azure REST API specifications commit. Their exact Git blobs, UTF-8 byte counts and SHA256 hashes are recorded in docs/COGNITIVE_SKU_CONTRACT_20261017.md. Copyright (c) 2017 Microsoft, MIT; the same complete original notice and permission text in storage-primary-LICENSE.txt applies. These are account-tier primary contract evidence, separate from immutable AZQR execution captures and model/deployment eligibility.
