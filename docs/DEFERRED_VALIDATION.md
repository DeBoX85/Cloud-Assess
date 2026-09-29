# Deferred Validation Register

This short register keeps environment-dependent evidence gaps visible outside the execution roadmap. A deferred item is not a PASS or an accepted release limitation. Review it at every Quality Gate 004 decision and before declaring a release candidate.

## DV-001: Live nested management-group traversal

**Status:** Deferred for lack of a suitable non-production nested hierarchy. Revisit before Quality Gate 004. Owner of scope selection: the operator of the Azure test environment.

**Why deferred:** The accessible `AdvisoryDev` group is a leaf. The known `Advisory` parent also contains a production child. The operator does not have a suitable non-production parent with nested child groups. A Dev subscription include filter is not a safe substitute because discovery visits child groups before applying returned-subscription filters. No parent scan or Azure hierarchy change is authorized by this register.

**Evidence already held:** Phase P compares the `AdvisoryDev` leaf with pinned AZQR. A deterministic test drives the production Azure SDK scope adapter through synthetic root, child and leaf responses, duplicate descendants, subscription filtering, disabled subscriptions and a denied child. It checks target implementation behavior, not live source-versus-target equivalence.

**Resume trigger:** A stable non-production management-group parent with at least one nested child and known accessible descendant subscriptions becomes available to the same identity for both scans. Verify the full resolved scope before executing. A production-containing parent requires a separate explicit scope decision.

**Completion evidence:** Record pinned source/target/APRL revisions; group hierarchy and permissions; requested, resolved and contributing subscriptions; stage health; retained private source/target JSON and logs; semantic comparator result; and classification of deltas and warnings. Update the development ledger and Quality Gate 004 matrix.

**If still unavailable at Gate 004:** Record a formal acceptance or block decision with impact and first-release scope. The current target specification advertises recursive management-group resolution, so synthetic tests alone must not be presented as live parity evidence. A narrowed first-release support claim requires corresponding updates to the target specification, implementation plan, README and ledger before the gate passes.
