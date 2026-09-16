package rules

const (
	ReferenceRepositoryCommit = "8e4f0577f3615e6c9014c031bcad079f235369cc"
	APRLCommit                = "60eaddda76541f6adbc1c5ffa686829807e55e29"
	AORSnapshotTree           = "a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5"
	CustomRulesSnapshotTree   = "674b9b3dcb443ce6dc445b48e1db47e4a0ca7082"
)

type Provenance struct {
	Source   string `json:"source"`
	Revision string `json:"revision"`
	Origin   string `json:"origin"`
}

// PinnedProvenance returns the rule-source revisions used as the Cloud Assess v1 equivalence baseline.
func PinnedProvenance() []Provenance {
	return []Provenance{
		{
			Source:   SourceAPRL,
			Revision: APRLCommit,
			Origin:   "Azure/Azure-Proactive-Resiliency-Library-v2",
		},
		{
			Source:   SourceAOR,
			Revision: AORSnapshotTree,
			Origin:   "embedded AOR snapshot from reference implementation",
		},
		{
			Source:   SourceCustom,
			Revision: CustomRulesSnapshotTree,
			Origin:   "custom rule snapshot from reference implementation",
		},
	}
}
