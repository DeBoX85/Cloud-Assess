package region

import (
	"context"
	"sort"
	"strings"
	"testing"
)

var _ = BuildRequiredInventory
var _ = AdmitRequiredEvidence

func TestRequiredInventoryFrozenLiterals(t *testing.T) {
	for _, c := range requiredFixtureLoad(t) {
		t.Run(c.Name, func(t *testing.T) {
			plan := requiredFixtureBuild(t, c)
			if plan == nil {
				return
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
