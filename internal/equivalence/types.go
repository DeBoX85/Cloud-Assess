package equivalence

import (
	"sort"
)

const (
	DatasetRecommendations         = "recommendations"
	DatasetFindings                = "findings"
	DatasetResourceTypes           = "resourceTypes"
	DatasetInventory               = "inventory"
	DatasetOutOfScope              = "outOfScope"
	DatasetAdvisor                 = "advisor"
	DatasetDefender                = "defender"
	DatasetDefenderRecommendations = "defenderRecommendations"
	DatasetAzurePolicy             = "azurePolicy"
	DatasetArcSQL                  = "arcSQL"
	DatasetCosts                   = "costs"
)

var datasetOrder = []string{
	DatasetRecommendations,
	DatasetFindings,
	DatasetResourceTypes,
	DatasetInventory,
	DatasetOutOfScope,
	DatasetAdvisor,
	DatasetAzurePolicy,
	DatasetArcSQL,
	DatasetDefenderRecommendations,
	DatasetDefender,
	DatasetCosts,
}

type Record struct {
	Key    string            `json:"key"`
	Fields map[string]string `json:"fields"`
}

type Dataset struct {
	Enabled bool              `json:"enabled"`
	Records map[string]Record `json:"-"`
}

type Projection struct {
	Comparable bool               `json:"comparable"`
	Notes      []string           `json:"notes,omitempty"`
	Datasets   map[string]Dataset `json:"-"`
}

type FieldDiff struct {
	Field     string `json:"field"`
	Reference string `json:"reference"`
	Target    string `json:"target"`
}

type ChangedRecord struct {
	Key    string      `json:"key"`
	Fields []FieldDiff `json:"fields"`
}

type DatasetDiff struct {
	Name             string          `json:"name"`
	ReferenceEnabled bool            `json:"referenceEnabled"`
	TargetEnabled    bool            `json:"targetEnabled"`
	ReferenceCount   int             `json:"referenceCount"`
	TargetCount      int             `json:"targetCount"`
	MissingCount     int             `json:"missingCount"`
	ExtraCount       int             `json:"extraCount"`
	ChangedCount     int             `json:"changedCount"`
	CoverageMismatch bool            `json:"coverageMismatch,omitempty"`
	Missing          []Record        `json:"missing,omitempty"`
	Extra            []Record        `json:"extra,omitempty"`
	Changed          []ChangedRecord `json:"changed,omitempty"`
}

type Report struct {
	Equivalent     bool          `json:"equivalent"`
	Preconditions  []string      `json:"preconditions,omitempty"`
	ReferenceNotes []string      `json:"referenceNotes,omitempty"`
	TargetNotes    []string      `json:"targetNotes,omitempty"`
	Datasets       []DatasetDiff `json:"datasets"`
}

func newProjection() Projection {
	datasets := make(map[string]Dataset, len(datasetOrder))
	for _, name := range datasetOrder {
		datasets[name] = Dataset{Records: map[string]Record{}}
	}
	return Projection{Comparable: true, Datasets: datasets}
}

func datasetNames(left, right Projection) []string {
	seen := map[string]struct{}{}
	names := make([]string, 0, len(datasetOrder))
	for _, name := range datasetOrder {
		if _, ok := left.Datasets[name]; ok {
			seen[name] = struct{}{}
			names = append(names, name)
			continue
		}
		if _, ok := right.Datasets[name]; ok {
			seen[name] = struct{}{}
			names = append(names, name)
		}
	}
	for name := range left.Datasets {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	for name := range right.Datasets {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	if len(names) > len(datasetOrder) {
		tail := append([]string(nil), names[len(datasetOrder):]...)
		sort.Strings(tail)
		copy(names[len(datasetOrder):], tail)
	}
	return names
}
