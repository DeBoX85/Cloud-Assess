# Optional-stage projection fixtures

`optional-stage-projections.json` is fabricated test data, not an Azure capture. All identifiers, resource names and descriptions are synthetic. Expected canonical JSON records were reviewed against the pinned reference's explicit row builders and Cloud Assess's documented field names:

- [Azure Policy row builder](https://github.com/DeBoX85/azqr/blob/8e4f0577f3615e6c9014c031bcad079f235369cc/internal/scanners/azure_policy.go): resource-ID helpers, policy fields and projected subscription display name.
- [Defender row builders](https://github.com/DeBoX85/azqr/blob/8e4f0577f3615e6c9014c031bcad079f235369cc/internal/scanners/defender.go): plan tier/name/display name; recommendation subscription-name lookup, projected resource type and `https://` prefix.
- [Reference tolerant row decoder](https://github.com/DeBoX85/azqr/blob/8e4f0577f3615e6c9014c031bcad079f235369cc/internal/graph/graph.go): skip a JSON row on unmarshal error while retaining valid rows.

The golden `expected` maps are literal fixture data. Tests compare every field of the persisted dataset against them, rather than producing expectations with the scanner or renderer under test. The differing projected/discovered subscription names deliberately verify their distinct source semantics. Numeric subscription IDs deliberately fail the string decoder.

The fixture exercises the production HTTP/ARG/scanner/coordinator/application/JSON path through a local TLS endpoint. It covers one representative row per stage, a mixed valid/malformed response, and an all-malformed response. Mixed/all-malformed responses retain explicit stage warnings and `complete_with_warnings`; all-malformed data must not look like a warning-free empty success. The existing 403/429/missing-data/valid-empty cases remain.

It does not execute KQL, prove Azure's live response shape, independently run the source scanner, test every row/filter variant, or close Gate 004's representative live coverage. Existing scanner unit tests cover additional field/filter/dedup cases. Arc SQL numeric `vcores` and live Policy/Defender Recommendations evidence remain open.
