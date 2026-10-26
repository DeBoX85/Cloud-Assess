package region

// First consumer, published before either production API exists. Keep the frozen
// fixtures and this first revision as evidence; compiler failure is not test success.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

// Put both names before other future types so both native compiler logs expose
// the actual missing APIs, rather than stopping after unrelated missing shapes.
var _ = BuildRequiredInventory
var _ = AdmitRequiredEvidence

const requiredFixtureDescriptorSHA = "99fe17f4367ce3ec37fc8d7a0b3b954a62543ab9b620f8312714e0243e48e6f2"
const requiredFixtureBase = "docs/evidence/20261021/recovered-local/task8-provisional/required-inventory-contract.provisional.json"
const requiredFixtureExtra = "docs/evidence/20261026/inventory-output/additional-expectations.json"

type requiredFixtureDeclaration struct {
	SubscriptionID string `json:"subscription_id"`
	Status         string `json:"status"`
}
type requiredFixtureRequest struct {
	SubscriptionID string `json:"subscription_id"`
	SourceRegion   string `json:"source_region"`
	TargetRegion   string `json:"target_region"`
}
type requiredFixtureCollection struct {
	SubscriptionID string                    `json:"subscription_id"`
	SourceRegion   string                    `json:"source_region"`
	TargetRegion   string                    `json:"target_region"`
	ResourceType   string                    `json:"resource_type"`
	Kind           string                    `json:"kind"`
	Storage        *StorageSKUCollection     `json:"storage"`
	Disk           *ComputeDiskSKUCollection `json:"disk"`
	VM             *ComputeVMSKUCollection   `json:"vm"`
	Cognitive      *CognitiveSKUCollection   `json:"cognitive"`
	SQL            *SQLCapabilityCollection  `json:"sql"`
}
type requiredFixtureRequirement struct {
	SubscriptionID string              `json:"subscription_id"`
	SourceRegion   string              `json:"source_region"`
	ResourceType   string              `json:"resource_type"`
	Kind           string              `json:"kind"`
	SKUs           []string            `json:"skus"`
	Kinds          map[string][]string `json:"kinds"`
	ResourceIDs    []string            `json:"resource_ids"`
	Flags          []string            `json:"flags"`
}
type requiredFixtureValue struct {
	State        int      `json:"state"`
	BlockedZones []string `json:"blocked_zones"`
}
type requiredFixtureSKU struct {
	Status string                          `json:"status"`
	Values map[string]requiredFixtureValue `json:"values"`
}
type requiredFixtureCatalog struct {
	ResourceType          string         `json:"resource_type"`
	Kind                  string         `json:"kind"`
	Status                string         `json:"status"`
	DeploymentEligibility string         `json:"deployment_eligibility"`
	Keys                  []string       `json:"keys"`
	States                map[string]int `json:"states"`
}
type requiredFixtureWant struct {
	Error           string                        `json:"error"`
	ErrorPhase      string                        `json:"error_phase"`
	PlanStatus      string                        `json:"plan_status"`
	ScopeStatus     string                        `json:"scope_status"`
	AdmissionStatus string                        `json:"admission_status"`
	Requirements    []requiredFixtureRequirement  `json:"requirements"`
	SKUs            map[string]requiredFixtureSKU `json:"sku_evidence"`
	Catalogs        []requiredFixtureCatalog      `json:"catalogs"`
	PreserveRaw     bool                          `json:"preserve_raw_records"`
}
type requiredFixtureCase struct {
	Name         string                       `json:"name"`
	Selected     map[string]string            `json:"selected"`
	Resources    []assessment.Resource        `json:"resources"`
	Declarations []requiredFixtureDeclaration `json:"declarations"`
	Request      requiredFixtureRequest       `json:"request"`
	Collections  []requiredFixtureCollection  `json:"collections"`
	Want         requiredFixtureWant          `json:"want"`
}

func requiredFixtureRead(t *testing.T, root, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func requiredFixtureEqual(t *testing.T, label string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s mismatch\ngot: %#v\nwant: %#v", label, got, want)
	}
}

// Python's frozen independent grammar rejects duplicate/aliased/unknown keys,
// trailing JSON, nonfinite numbers, wrong presence/types and invalid Unicode
// BEFORE any typed Go decode. Go's default case-insensitive matching is not the guard.
func requiredFixtureLoad(t *testing.T) []requiredFixtureCase {
	t.Helper()
	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	descriptor := requiredFixtureRead(t, root, "docs/evidence/20261026/inventory-freeze/contract-descriptor.json")
	digest := sha256.Sum256(descriptor)
	if hex.EncodeToString(digest[:]) != requiredFixtureDescriptorSHA {
		t.Fatal("frozen descriptor changed")
	}
	var bound struct {
		Files []struct {
			Path   string
			Bytes  int
			SHA256 string
		}
	}
	if err := json.Unmarshal(descriptor, &bound); err != nil {
		t.Fatal(err)
	}
	if len(bound.Files) != 32 {
		t.Fatal("dependency inventory changed")
	}
	seen := map[string]bool{}
	for _, f := range bound.Files {
		if !strings.HasPrefix(f.Path, "docs/") || strings.Contains(f.Path, "..") || seen[f.Path] {
			t.Fatal("invalid dependency path")
		}
		seen[f.Path] = true
		b := requiredFixtureRead(t, root, f.Path)
		sum := sha256.Sum256(b)
		if len(b) != f.Bytes || hex.EncodeToString(sum[:]) != f.SHA256 {
			t.Fatalf("dependency changed: %s", f.Path)
		}
	}
	python := "python3"
	if runtime.GOOS == "windows" {
		python = "python"
	}
	cmd := exec.Command(python, filepath.Join(root, "docs/evidence/20261026/inventory-freeze/check-consumer-schema.py.txt"), root)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("strict fixture preflight: %v\n%s", err, output)
	}
	var proof struct {
		Original     int      `json:"original_cases_structurally_admitted"`
		Supplemental int      `json:"supplemental_cases_structurally_admitted"`
		Procedural   int      `json:"procedural_expectations_preserved"`
		Rejected     []string `json:"negative_controls_rejected"`
		Executed     bool     `json:"future_api_executed"`
	}
	if err := json.Unmarshal(output, &proof); err != nil {
		t.Fatal(err)
	}
	if proof.Original != 50 || proof.Supplemental != 26 || proof.Procedural != 13 || len(proof.Rejected) != 10 || proof.Executed {
		t.Fatal("strict preflight inventory changed")
	}
	var manifest struct {
		Cases []json.RawMessage `json:"cases"`
	}
	if err := json.Unmarshal(requiredFixtureRead(t, root, requiredFixtureBase), &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Cases) != 50 {
		t.Fatal("original case inventory changed")
	}
	known := map[string]json.RawMessage{}
	cases := []requiredFixtureCase{}
	for _, raw := range manifest.Cases {
		var c requiredFixtureCase
		if err := json.Unmarshal(raw, &c); err != nil {
			t.Fatal(err)
		}
		if _, exists := known[c.Name]; exists {
			t.Fatal("duplicate original name")
		}
		known[c.Name] = raw
		cases = append(cases, c)
	}
	var extra struct {
		Cases []struct {
			Name     string                     `json:"name"`
			Base     string                     `json:"base_case"`
			Set      map[string]json.RawMessage `json:"set"`
			Error    string                     `json:"want_error"`
			Phase    string                     `json:"want_phase"`
			Nil      bool                       `json:"want_nil"`
			Expected map[string]json.RawMessage `json:"want_same_as_base_except"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(requiredFixtureRead(t, root, requiredFixtureExtra), &extra); err != nil {
		t.Fatal(err)
	}
	if len(extra.Cases) != 26 {
		t.Fatal("supplemental case inventory changed")
	}
	aliases := map[string]string{"": "", "build": "inventory", "admit": "admission", "inventory": "inventory", "admission": "admission"}
	for _, patch := range extra.Cases {
		base, ok := known[patch.Base]
		if !ok {
			t.Fatal("unknown literal base")
		}
		var v any
		dec := json.NewDecoder(bytes.NewReader(base))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			t.Fatal(err)
		}
		m := v.(map[string]any)
		m["name"] = patch.Name
		for path, raw := range patch.Set {
			requiredFixtureSet(t, m, path, raw)
		}
		if patch.Error != "" {
			phase, ok := aliases[patch.Phase]
			if !ok || phase == "" || !patch.Nil {
				t.Fatal("invalid supplemental error expectation")
			}
			w := m["want"].(map[string]any)
			w["error"] = patch.Error
			w["error_phase"] = phase
		} else {
			for path, raw := range patch.Expected {
				requiredFixtureSet(t, m["want"], path, raw)
			}
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		var c requiredFixtureCase
		if err := json.Unmarshal(b, &c); err != nil {
			t.Fatal(err)
		}
		if _, exists := known[c.Name]; exists {
			t.Fatal("duplicate supplemental name")
		}
		known[c.Name] = b
		cases = append(cases, c)
	}
	if len(cases) != 76 || len(known) != 76 {
		t.Fatal("strict consumer did not retain all76cases")
	}
	return cases
}

// Apply literal SETs only, never calculate an expected result from product code.
func requiredFixtureSet(t *testing.T, doc any, path string, raw json.RawMessage) {
	t.Helper()
	if !strings.HasPrefix(path, "/") {
		t.Fatal("invalid literal path")
	}
	parts := strings.Split(path[1:], "/")
	for i := range parts {
		parts[i] = strings.ReplaceAll(strings.ReplaceAll(parts[i], "~1", "/"), "~0", "~")
	}
	var value any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		t.Fatal(err)
	}
	for i, k := range parts {
		last := i == len(parts)-1
		switch x := doc.(type) {
		case map[string]any:
			if last {
				x[k] = value
				return
			}
			next, ok := x[k]
			if !ok {
				t.Fatalf("absent literal parent %s", path)
			}
			doc = next
		case []any:
			n, err := strconv.Atoi(k)
			if err != nil || n < 0 || n >= len(x) {
				t.Fatalf("invalid literal index %s", path)
			}
			if last {
				x[n] = value
				return
			}
			doc = x[n]
		default:
			t.Fatalf("invalid literal parent %s", path)
		}
	}
}

func requiredFixtureSKUs(t *testing.T, in map[string]SKUEvidence) map[string]requiredFixtureSKU {
	t.Helper()
	if in == nil {
		t.Fatal("nil scoring map")
	}
	out := map[string]requiredFixtureSKU{}
	for typ, e := range in {
		if e.Values == nil {
			t.Fatal("nil scoring values")
		}
		v := map[string]requiredFixtureValue{}
		for k, row := range e.Values {
			if row.BlockedZones == nil {
				t.Fatal("nil zone slice")
			}
			v[k] = requiredFixtureValue{row.State, row.BlockedZones}
		}
		out[typ] = requiredFixtureSKU{e.Status, v}
	}
	return out
}

func TestRequiredInventoryFrozenLiterals(t *testing.T) {
	for _, c := range requiredFixtureLoad(t) {
		t.Run(c.Name, func(t *testing.T) {
			declarations := make([]InventoryScopeDeclaration, len(c.Declarations))
			for i, d := range c.Declarations {
				declarations[i] = InventoryScopeDeclaration{SubscriptionID: d.SubscriptionID, Status: d.Status}
			}
			plan, err := BuildRequiredInventory(context.Background(), c.Selected, c.Resources, declarations)
			if c.Want.ErrorPhase == "inventory" {
				if plan != nil || err == nil || err.Error() != c.Want.Error {
					t.Fatalf("inventory rejection: plan=%v err=%v want=%s", plan, err, c.Want.Error)
				}
				return
			}
			if err != nil || plan == nil {
				t.Fatalf("unexpected build result plan=%v err=%v", plan, err)
			}
			requiredFixtureEqual(t, "plan status", plan.Status, c.Want.PlanStatus)
			if plan.Resources == nil || plan.Requirements == nil || plan.Scopes == nil {
				t.Fatal("nil owned plan container")
			}
			reqs := []requiredFixtureRequirement{}
			for _, r := range plan.Requirements {
				if r.SKUs == nil || r.Kinds == nil || r.ResourceIDs == nil {
					t.Fatal("nil requirement container")
				}
				flags := []string{}
				for k, on := range map[string]bool{"catalog_only": r.CatalogOnly, "unsupported": r.Unsupported, "missing_sku": r.MissingSKU, "missing_kind": r.MissingKind, "unrepresentable_sku": r.UnrepresentableSKU, "unrepresentable_kind": r.UnrepresentableKind} {
					if on {
						flags = append(flags, k)
					}
				}
				sort.Strings(flags)
				reqs = append(reqs, requiredFixtureRequirement{r.SubscriptionID, r.SourceRegion, r.ResourceType, r.Kind, r.SKUs, r.Kinds, r.ResourceIDs, flags})
			}
			requiredFixtureEqual(t, "requirements", reqs, c.Want.Requirements)
			if len(plan.Scopes) != len(c.Selected) {
				t.Fatal("scope count mismatch")
			}
			for id, name := range c.Selected {
				status := "unknown"
				for _, d := range c.Declarations {
					if strings.EqualFold(d.SubscriptionID, id) {
						status = d.Status
					}
				}
				s, ok := plan.Scopes[strings.ToLower(id)]
				if !ok || s.Name != name || s.Status != status {
					t.Fatalf("scope diagnostics mismatch for %s", id)
				}
			}
			if !c.Want.PreserveRaw {
				t.Fatal("unrecognized raw record expectation")
			}
			if len(plan.Resources) != len(c.Resources) {
				t.Fatal("raw record count mismatch")
			}
			byID := map[string]assessment.Resource{}
			sourceByID := map[string]string{}
			for _, r := range c.Resources {
				r.Tags = nil
				r.SLA = ""
				byID[r.ID] = r
			}
			for _, r := range c.Want.Requirements {
				for _, id := range r.ResourceIDs {
					sourceByID[id] = r.SourceRegion
				}
			}
			previous := ""
			seen := map[string]bool{}
			for _, r := range plan.Resources {
				raw, ok := byID[r.Record.ID]
				if !ok || seen[r.Record.ID] {
					t.Fatal("unknown/duplicate raw record")
				}
				seen[r.Record.ID] = true
				requiredFixtureEqual(t, "retained raw record", r.Record, raw)
				key := strings.ToLower(raw.ID)
				if previous != "" && key <= previous {
					t.Fatal("resource order mismatch")
				}
				previous = key
				if r.SubscriptionID != strings.ToLower(raw.SubscriptionID) || r.ResourceType != strings.ToLower(raw.Type) || r.Kind != strings.ToLower(raw.Kind) || r.SourceRegion != sourceByID[raw.ID] || r.Physical != (sourceByID[raw.ID] != "") {
					t.Fatal("canonical resource coordinates mismatch")
				}
			}
			request := RequiredEvidenceRequest{SubscriptionID: c.Request.SubscriptionID, SourceRegion: c.Request.SourceRegion, TargetRegion: c.Request.TargetRegion}
			collections := make([]RequiredCollectionInput, len(c.Collections))
			for i, x := range c.Collections {
				collections[i] = RequiredCollectionInput{SubscriptionID: x.SubscriptionID, SourceRegion: x.SourceRegion, TargetRegion: x.TargetRegion, ResourceType: x.ResourceType, Kind: x.Kind, Storage: x.Storage, Disk: x.Disk, VM: x.VM, Cognitive: x.Cognitive, SQL: x.SQL}
			}
			admitted, err := AdmitRequiredEvidence(context.Background(), plan, request, collections)
			if c.Want.ErrorPhase == "admission" {
				if admitted != nil || err == nil || err.Error() != c.Want.Error {
					t.Fatalf("admission rejection: result=%v err=%v want=%s", admitted, err, c.Want.Error)
				}
				return
			}
			if c.Want.Error != "" || c.Want.ErrorPhase != "" {
				t.Fatal("unrecognized expected error")
			}
			if err != nil || admitted == nil {
				t.Fatalf("unexpected admission result=%v err=%v", admitted, err)
			}
			canonical := RequiredEvidenceRequest{SubscriptionID: strings.ToLower(c.Request.SubscriptionID), SourceRegion: strings.ReplaceAll(strings.ToLower(c.Request.SourceRegion), " ", ""), TargetRegion: strings.ReplaceAll(strings.ToLower(c.Request.TargetRegion), " ", "")}
			requiredFixtureEqual(t, "admitted request", admitted.Request, canonical)
			requiredFixtureEqual(t, "scope status", admitted.ScopeStatus, c.Want.ScopeStatus)
			requiredFixtureEqual(t, "admission status", admitted.Status, c.Want.AdmissionStatus)
			requiredFixtureEqual(t, "scoring", requiredFixtureSKUs(t, admitted.SKUs), c.Want.SKUs)
			if admitted.Catalogs == nil {
				t.Fatal("nil catalog slice")
			}
			catalogs := []requiredFixtureCatalog{}
			for _, b := range admitted.Catalogs {
				if b.SubscriptionID != canonical.SubscriptionID || b.SourceRegion != canonical.SourceRegion || b.TargetRegion != canonical.TargetRegion {
					t.Fatal("catalog binding mismatch")
				}
				keys := []string{}
				states := map[string]int{}
				payloads := 0
				if b.Cognitive != nil {
					payloads++
					for k, v := range b.Cognitive.Values {
						keys = append(keys, k)
						states[k] = v.State
					}
				}
				if b.SQL != nil {
					payloads++
					for k, v := range b.SQL.Values {
						keys = append(keys, k)
						states[k] = v.State
					}
				}
				if b.VMSS != nil {
					payloads++
					for k, v := range b.VMSS.Values {
						keys = append(keys, k)
						states[k] = v.State
					}
				}
				if payloads != 1 {
					t.Fatal("catalog diagnostic union mismatch")
				}
				sort.Strings(keys)
				catalogs = append(catalogs, requiredFixtureCatalog{b.ResourceType, b.Kind, b.Status, b.DeploymentEligibility, keys, states})
			}
			requiredFixtureEqual(t, "catalogs", catalogs, c.Want.Catalogs)
			bound, values := admitted.Evidence()
			requiredFixtureEqual(t, "accessor binding", bound, canonical)
			requiredFixtureEqual(t, "accessor scoring", requiredFixtureSKUs(t, values), c.Want.SKUs)
		})
	}
}
