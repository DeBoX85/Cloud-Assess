package equivalence

import (
	"encoding/json"
	"reflect"
	"testing"
)

func guardProjection(value string) Projection {
	return Projection{Comparable: true, Datasets: map[string]Dataset{
		DatasetFindings: {Enabled: true, Records: map[string]Record{"one": {Key: "one", Fields: map[string]string{"impact": value}}}},
	}}
}

func TestCompareGuards(t *testing.T) {
	for _, name := range []string{"missing", "extra", "changed", "coverage", "precondition"} {
		t.Run(name, func(t *testing.T) {
			left, right := guardProjection("High"), guardProjection("High")
			dataset := right.Datasets[DatasetFindings]
			switch name {
			case "missing":
				delete(dataset.Records, "one")
			case "extra":
				dataset.Records["two"] = Record{Key: "two", Fields: map[string]string{"impact": "High"}}
			case "changed":
				dataset.Records["one"].Fields["impact"] = "Low"
			case "coverage":
				dataset.Enabled = false
			case "precondition":
				right.Comparable = false
			}
			right.Datasets[DatasetFindings] = dataset
			before, _ := json.Marshal(right.Datasets[DatasetFindings].Records)
			got := Compare(left, right)
			if got.Equivalent {
				t.Fatalf("%s incorrectly passed", name)
			}
			if name == "precondition" {
				if len(got.Preconditions) != 1 {
					t.Fatal("failed precondition was not reported")
				}
			} else {
				diff := got.Datasets[0]
				if (name == "missing" && diff.MissingCount != 1) || (name == "extra" && diff.ExtraCount != 1) || (name == "changed" && diff.ChangedCount != 1) || (name == "coverage" && !diff.CoverageMismatch) {
					t.Fatalf("wrong classification: %+v", diff)
				}
			}
			after, _ := json.Marshal(right.Datasets[DatasetFindings].Records)
			if string(before) != string(after) {
				t.Fatal("comparator mutated input")
			}
		})
	}
	if !Compare(guardProjection("High"), guardProjection("High")).Equivalent {
		t.Fatal("identical input failed comparison")
	}
}

func TestDatasetNamesStableForSparseAndUnknownDatasets(t *testing.T) {
	left := Projection{Datasets: map[string]Dataset{"z-custom": {}, DatasetFindings: {}, "a-custom": {}}}
	want := []string{DatasetFindings, "a-custom", "z-custom"}
	for i := 0; i < 100; i++ {
		if got := datasetNames(left, Projection{}); !reflect.DeepEqual(got, want) {
			t.Fatalf("dataset order=%v, want %v", got, want)
		}
	}
}

func FuzzCompareFieldSymmetry(f *testing.F) {
	f.Add("High", "Low")
	f.Add("", "")
	f.Add("\x00", "é")
	f.Fuzz(func(t *testing.T, a, b string) {
		left, right := guardProjection(a), guardProjection(b)
		forward, backward := Compare(left, right), Compare(right, left)
		if forward.Equivalent != (a == b) || backward.Equivalent != forward.Equivalent {
			t.Fatalf("comparison disagrees with independent field equality: %q %q", a, b)
		}
		encoded, err := json.Marshal(forward)
		if err != nil || !json.Valid(encoded) {
			t.Fatalf("invalid report encoding: %v", err)
		}
	})
}
