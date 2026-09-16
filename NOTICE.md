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

The source/rule license texts above do not replace a dependency inventory. Before any release or redistribution, generate the dependency and embedded-content inventory from the actual Cloud Assess release tree and dependency graph, then review `NOTICE.md` and `THIRD_PARTY_LICENSES.md` for completeness.
