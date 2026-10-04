# Audit ARG response completion correction

Current status: correction VERIFIED OFFLINE and accepted through combined PR114, mergec33fb2eee93a2e5730f72612b351d0b838737d03. Original individual proposal is closed/superseded; compiling failed controls and corrected proof below are retained as historical evidence. Final combined native37188397367/source37188397374 and matching tree/preview verified. Distinct accepted-push/final disposition evidence is indexed in PR114 and the disposition PR. No historical live incidence, independent-person review or release approval.

Date: 2026-10-04. Audit finding AUD003. Frozen accepted base: 07011b63440e69f4a5acb128e0449943f759ee9c. Comprehensive audit PR114 and feature pause remain active.

## Confirmed boundary and impact

The actual ARG HTTPTransport decoded one JSON response and then ignored io.Copy drain errors. Test-only head 730a29e6eed02d82283803931aa50ca599bb9706 reproduced successful results after terminal non-EOF body failure in all eight named assertions on both native Linux and Windows. Run [37183817468](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37183817468), jobs111381519563 and111381519469, compiled successfully; formatting passed. Literal empty/nonempty JSON followed by sentinel read error, unexpected EOF, cancellation or deadline returned success. Healthy EOF controls did not fail.

This is a confirmed library adapter contract gap with an injected streamPoster. It is not evidence of a default live Azure false-complete assessment. The pinned azcore v1.23.1 runtime.NewPipeline inserts bodyDownloadPolicy; this project's PostStream does not call SkipBodyDownload. That middleware reads the complete underlying body before returning and already rejects these errors. The additional independent shared HTTPClient fixture verifies this mitigation through the actual authenticated pipeline with one transport call and one underlying closure.

Pinned AZQR8e4f0577 internal/graph/graph.go contains the same ignored io.Copy pattern. The correction is an intentional defensive difference from that source, not a newly introduced fork regression. Adjacent Diagnostics uses the same default buffered pipeline; its inspection does not establish a corresponding live defect.

## Correction and independent regression

Propagate non-EOF drain errors with wrapping, retaining errors.Is identity and returning no successful response. Keep deferred closure and healthy EOF behavior. Do not change SDK body buffering, endpoints, query semantics, quotas, source pins or global timeout policy.

The independent reader owns a named strings.Reader field, preventing a promoted WriterTo from bypassing its terminal Read error. The original caller context remains healthy. Tests exercise actual Client.Query plus HTTPTransport.Do and cover empty/nonempty JSON, four errors, exact closure, and successful zero/one-row EOF controls. A separate fixture uses actual azure.HTTPClient, credentials and SDK policy pipeline; it verifies the default route's existing rejection, not production reachability of the injected-poster defect.

## Gates and remaining work

Fresh native Linux/Windows and pinned-source QA on the corrected head remain pending. Record complete logs, preview commit/tree/parents, provenance and capture evidence before considering acceptance. The compiling failing reproduction is retained at the immutable test-only commit. No audit-wide, independent-person, Azure/live, Gate004 or release closure is claimed. Rollback is a reviewed revert after dependency review.

Primary contracts: [io.Copy](https://pkg.go.dev/io#Copy), [json.Decoder.Decode](https://pkg.go.dev/encoding/json#Decoder.Decode). SDK mitigation inspected at official Azure/azure-sdk-for-go tag sdk/azcore/v1.23.1: runtime/pipeline.go and runtime/policy_body_download.go, plus sdk/internal/v1.12.0/exported/exported.go Payload.
