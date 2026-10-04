# Audit ARG response completion investigation

Date2026-10-04. Audit finding AUD003. Test-only reproduction candidate, production unchanged. Accepted base07011b63440e69f4a5acb128e0449943f759ee9c. Comprehensive audit PR114 and feature pause remain authoritative.

## Trigger and expected behavior

Actual internal/arg/http_transport.go decodes one JSON response and then ignores both values returned by io.Copy(io.Discard,response.Body). A valid empty or nonempty JSON value followed by a non-EOF body read error may therefore produce a successful ARG query. This remains a suspected completeness defect until the real compiling regression is observed failing against unchanged production code.

Independent fixtures provide literal valid empty/nonempty JSON, then separately signal a synthetic read error, unexpected EOF, cancellation or deadline. The original caller context is healthy, so result rejection must arise from the operation/body failure. Client.Query plus actual HTTPTransport.Do must return no successful result, retain errors.Is identity, and close the body once. Healthy terminal EOF must preserve exact zero/one rows and one closure. The fixture reader owns a named strings.Reader field rather than embedding it, preventing io.Copy from selecting a promoted WriterTo that would bypass the failure.

This tests the library transport boundary with an injected streamPoster. It is not a live Azure comparison or a claim about SDK middleware ordering. Follow with production shared-client/coordinator evidence if needed to establish reachability and report consequences.

## Plan and gates

First publish the regression against unchanged production code and require a compiling named assertion failure in native Linux/Windows tests, with healthy EOF fixtures passing. Classify unexpected compile/format/test results rather than accepting them as reproduction.

If confirmed, propagate terminal body errors with wrapping instead of discarding them; retain error identity, nil successful candidate and closure. Inspect Diagnostics' adjacent streaming decoder and inherited pinned-source behavior. Record a deliberate safety correction and avoid unrelated schema/normalization/volume changes.

Fresh full native/source checks on the corrected head, diff/source-contract review and protected expected-head merge remain required. Do not claim audit completion, Azure/live/Gate004/release closure or independent-person review. Rollback is a reviewed revert after dependency review.

Primary contracts: [io.Copy](https://pkg.go.dev/io#Copy) returns non-EOF copy failures, while [json.Decoder.Decode](https://pkg.go.dev/encoding/json#Decoder.Decode) reads one JSON value. TARGET_SPECIFICATION requires visible failed retrieval and HTTP operation budgets covering body consumption. Existing interrupted-request tests do not exercise a failure after a valid JSON value.
