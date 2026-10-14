# Disk collector preparation correction, 2026-10-14

PR140 cd9d0807 remains unaccepted. Focused preparation38497615288:
Windows115636893592 passed all collector cases, vet and21compiling controls
with restored named baselines (full210814log). Linux115636894116 passed the
collector baseline (including race), then FAILED while decoding intentionally
malformed mutation-test diagnostic bytes as UTF-8. This is a Python harness
failure, not evidence that the failing mutation passed or production was accepted.
Preserve the failed run. Use explicit UTF-8/backslashreplace only for subprocess
diagnostics and quoted malformed-body test output; ASCII named runtime-failure
and compile/restoration checks remain mandatory. No fixture result, collector
admission rule, source pin, module or caller budget is relaxed.
Adopt independently reconstructed native gofmt output. Both hosts and full
final-head native/source/review/merge/distinct-push evidence remain required.

---

# Corrective QA verified; first disk collector continuation, 2026-10-14

ARG fix459f73dd: complete focused Linux115626163022/full84049 and
Windows115626163109/full85385 both passed. Baseline compiled-old bug reconfirmed,
all ARG tests (Linux race), vet,4compiling controls and restored source passed.
Complete current core native38494207758 also passed: Linux115626163433/full97812/
59steps/291controls; Windows115626163757/full78227/44steps/286controls.
Every earlier287/282named identity preserved. Complete formatter files/source hash
matched authored core/test code. This completes the corrective QA phase before
continuation; no independent-human approval or final combined acceptance inferred.

Now implementing disk-only internal collector under published concrete contract.
Added15literal full states/maps,80+malformed fixtures, exact bounds, actual transport
canary/closures, snapshots and8concurrent owned results.21compiling production
controls require specific named runtime failures and restored named baselines.
WIP files and preparation are not proof until complete logs/formatter bytes read.
Final edited head must own all native/source/review/merge/push evidence; do not
transfer this core-only pass. Other review scope/deferred gates stay as recorded.

---

# Confirmed finding and corrective QA preparation, 2026-10-14

P2 ARG trailing-data bug is reproduced on Linux and Windows after successful
compilation of original canonical blob3c93b2bb. Run38493524807:
Linux115623988968/full37682 and Windows115623989223/full38665,18named buffered/split
cases each. Formatter2778bytes/c1fbc2a4 independently matches authored new tests.
Initial Windows115622579696/full23568FAILED raw checkout hash because of
CRLF conversion; canonical identity was unchanged. Corrected checker uses git show
and ONLY native CRLF reconciliation. Keep that failed preparation as history.

Correction validates both decoder-buffered bytes and remaining body as JSON
whitespace until EOF, streaming with1024byte scratch. It does not decode/allocate
a second arbitrary JSON value. Three independent test functions preserve one full
literal row/future metadata, reject18bad suffix cases including late extra data,
and verify body closure and wrapped read/context errors.
Four compiling production controls have separate named assertions and restored
baselines. Full ARG package/transport tests and vet on both hosts, Linux race,
formatter/fullfile proof and all source/native final gates remain required.
This source fix is a deliberate target protocol-framing correction; original
source captures/pins/module graph are untouched. No live/whole-repo safety claim.

---

# Windows reproduction preparation correction, 2026-10-14

Linux115622579941 reproduced the compiled original parser gap with18buffered/split
runtime failures; full37454log, formatter2778bytes/c1fbc2a4 matched authored test.
Windows115622579696 FAILED BEFORE compilation at raw checkout byte SHA assertion.
Its486532f9 was independently reproduced by LF-to-CRLF checkout conversion of the
original3c93b2bb blob. This is preparation failure, not Windows parser evidence.
Use git show HEAD:<path> for canonical Git blob verification; independently require
native checkout equal canonical bytes after ONLY CRLF conversion. No arbitrary
normalization, source edit or lowered expected identity. Re-run both reproductions.
An immediate PR read also lagged after a successful leased branch update; actual
original branch ref and repeated PR read matched. No duplicate mutation was made.

---

# Code review and QA checkpoint, 2026-10-14

Status: UNACCEPTED review work on PR140, before continuation into target disk code.
Accepted protected baseline remains PR1396f923c1d. The user explicitly requested
code review of work thus far, QA, then continuation. Earlier supplied AGENTS.md
instructions were revoked; this report records technical evidence and task choices.

## Review scope and confirmed finding

Inventoried all production Go paths and retrieved118current code/workflow/helper
files (899311characters) for inspection. Retrieval is not a claim of line-by-line
manual review of every file. Manual risk review covered HTTP/ARM/ARG parsing and
retry boundaries, selected scope and cancellation, region provider/quota/reservation
collection, result ownership/health, report staging/Windows privacy, JSON redaction,
CSV/XLSX/SARIF export policy and recent Windows cleanup/source capture changes.
Full existing native QA covers other packages; no universal absence-of-bugs claim.

P2: internal/arg/http_transport.go decodes the first JSON document, then io.Copy
drains only the underlying body without validating trailing bytes, including bytes
already buffered by json.Decoder. Thus a valid ARG prefix followed by a second JSON
value or malformed data can return successful query results. Regression cases
are declared independently before the production fix: buffered/split suffixes,
large whitespace then extra data, valid whitespace/unknown metadata, body close,
and retained I/O/cancellation/deadline errors. Hosted compile plus named expected
failure reproduction is pending; do not call expected failing WIP acceptance.

Correction will validate EOF through decoder.Buffered plus the original body,
using fixed-size streaming whitespace admission to avoid allocating a second
unbounded JSON value. Keep valid response rows, future metadata and body closure.
Preserve original source pins/oracles; this is a deliberate target framing correction.

SARIF retains resource identities by the explicit CLI flag contract. That behavior
was checked and is not reported as an undisclosed redaction bug. Existing deferred
AR-02 volume/load/general ARG response/paging bounds, module-only GO-2026-5932,
live/freshOS/integration/public-region/release gates remain open and are not erased
by a code-review pass. Laptop/Azure access remains unconfirmed.

## QA and exact next

PR139 full candidate and distinct accepted-push proof remains in its live body.
PR140 initial47e73257 documentation-only CI succeeded on Linux/Windows/source;
new regression edits require their own QA. The temporary hosted reproduction
workflow verifies the original ARG blob before running new tests and requires
successful compilation followed by the named failed runtime assertion on both
hosts. Standard gates may fail on this deliberate reproduction head. Preserve that
failure as evidence; fix the parser, verify positive and compiling mutation controls,
remove the helper, run complete current-head review/native/source QA, then continue
the Compute disk collector under its published concrete pre-edit contract.
