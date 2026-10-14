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
