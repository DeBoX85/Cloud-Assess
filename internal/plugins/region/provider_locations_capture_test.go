package region

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func providerCapturedJSON(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	value := map[string]any{}
	if json.Unmarshal(raw, &value) != nil {
		t.Fatal("invalid provider capture JSON")
	}
	return value
}

func providerLiteral(t *testing.T, text string) any {
	t.Helper()
	var value any
	if json.Unmarshal([]byte(text), &value) != nil {
		t.Fatal("invalid independent provider literal")
	}
	return value
}

func providerLiteralCaches() map[string]string {
	result := map[string]string{}
	set := func(names, value string) {
		for _, name := range strings.Fields(names) {
			result[name] = value
		}
	}
	set("basic normalization null-location case-field-alias duplicate-field", "{\"Data\":{\"microsoft.compute\":{\"virtualmachines\":{\"westeurope\":{}}}}}")
	set("global", "{\"Data\":{\"microsoft.compute\":{\"virtualmachines\":{\"global\":{}}}}}")
	set("empty missing-value null-value null-envelope missing-namespace null-namespace missing-resource-types null-resource-types", "{\"Data\":{}}")
	set("empty-resource-types missing-type", "{\"Data\":{\"microsoft.compute\":{}}}")
	set("empty-locations null-locations missing-locations duplicate-to-empty", "{\"Data\":{\"microsoft.compute\":{\"virtualmachines\":{}}}}")
	set("duplicate-type paged-overwrite", "{\"Data\":{\"microsoft.compute\":{\"virtualmachines\":{\"eastus\":{}}}}}")
	set("empty-namespace", "{\"Data\":{\"\":{\"virtualmachines\":{\"westeurope\":{}}}}}")
	set("empty-type", "{\"Data\":{\"microsoft.compute\":{\"\":{\"westeurope\":{}}}}}")
	set("key-whitespace", "{\"Data\":{\" microsoft.compute \":{\" virtualmachines \":{\"westeurope\":{}}}}}")
	set("location-whitespace", "{\"Data\":{\"microsoft.compute\":{\"virtualmachines\":{\"eastus\":{},\"west\\teurope\":{},\"west\u00a0europe\":{}}}}}")
	set("nested-type", "{\"Data\":{\"microsoft.sql\":{\"servers/databases\":{\"westeurope\":{}}}}}")
	set("unicode-namespace", "{\"Data\":{\"microsoft.k\":{\"widgets\":{\"westeurope\":{}}}}}")
	set("paged", "{\"Data\":{\"microsoft.compute\":{\"virtualmachines\":{\"westeurope\":{}}},\"microsoft.storage\":{\"storageaccounts\":{\"eastus\":{}}}}}")
	set("cancelled later-cancelled cycle-stopped relative-next malformed-next malformed wrong-value wrong-locations later-malformed http400 http403 http404 http405 later-http400 later-http403 later-http404 later-http405 null-provider null-type", "null")
	return result
}

func providerNameSet(names string) map[string]bool {
	result := map[string]bool{}
	for _, name := range strings.Fields(names) {
		result[name] = true
	}
	return result
}

func TestProviderSourceMetadata(t *testing.T) {
	input := providerCapturedJSON(t, "source-provider-locations-inputs.json")
	want := map[string]any{}
	want["SourceBlob"] = "55e59b9e061dd31d27cc3446ad00659bd6961174"
	want["TypesBlob"] = "f9b61fb74075e5a2297be1b898067b8cb047f9b3"
	want["OptionsBlob"] = "73693e3730f5ed51132a552b79d5e95078721a10"
	want["CloudBlob"] = "9677ae498c8d22ebbc77a1fe3a47e17579f3c4d5"
	want["SDKModule"] = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0"
	want["AzcoreModule"] = "github.com/Azure/azure-sdk-for-go/sdk/azcore v1.23.1"
	want["SDKProvidersBlob"] = "e1612d16090a9a66c4b2bb9fa5697e7376263992"
	want["SDKDefaultTransportBlob"] = "928e9bf92550dcd83180b7fecff40b7a50fd7ff9"
	want["ExecutionMode"] = "original-source-custom-cloud-loopback-TLS; separate-public-SDK-injected-transport"
	want["TLSVerificationEnabled"] = true
	want["DirectSDKDisableRPRegistration"] = true
	want["DirectSDKMaxRetries"] = float64(-1)
	want["PublicSDKRequests"] = []any{"GET https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers?api-version=2021-04-01"}
	for key, value := range want {
		if !reflect.DeepEqual(input[key], value) {
			t.Fatal("provider source metadata mismatch", key)
		}
	}
	cases, ok := input["Cases"].(map[string]any)
	if !ok || len(cases) != len(providerLiteralCaches()) {
		t.Fatal("provider source input matrix mismatch")
	}
	for name := range providerLiteralCaches() {
		if _, exists := cases[name]; !exists {
			t.Fatal("provider source input case missing", name)
		}
	}
	// Original complete-byte fingerprint complements the independent cell tests.
	raw, err := os.ReadFile("testdata/source-provider-locations-inputs.json")
	if err != nil || fmt.Sprintf("%x", sha256.Sum256(raw)) != "aba5e5eb67a149bdd5f297518688fc8cf8b0c8daa89d48238f1c0bbd5175c910" {
		t.Fatal("provider source input bytes changed")
	}
}

func TestProviderSourceCacheSemantics(t *testing.T) {
	output := providerCapturedJSON(t, "source-provider-locations-outputs.json")
	want := providerLiteralCaches()
	failed := providerNameSet("cancelled later-cancelled cycle-stopped relative-next malformed-next malformed wrong-value wrong-locations later-malformed http400 http403 http404 http405 later-http400 later-http403 later-http404 later-http405")
	panicked := providerNameSet("null-provider null-type")
	if len(want) != 48 || len(output) != len(want) {
		t.Fatal("provider source output matrix mismatch")
	}
	for name, cache := range want {
		value, ok := output[name].(map[string]any)
		if !ok || len(value) != 9 || !reflect.DeepEqual(value["Cache"], providerLiteral(t, cache)) {
			t.Fatal("provider complete cache mismatch", name)
		}
		if value["Failed"] != failed[name] || value["Panicked"] != panicked[name] || value["CallerPanicked"] != panicked[name] {
			t.Fatal("provider failure/panic mismatch", name)
		}
	}
}

func TestProviderSourceRequests(t *testing.T) {
	output := providerCapturedJSON(t, "source-provider-locations-outputs.json")
	first := "/subscriptions/11111111-1111-1111-1111-111111111111/providers?api-version=2021-04-01"
	second := "/subscriptions/11111111-1111-1111-1111-111111111111/providers?%24skiptoken=synthetic-second&api-version=2021-04-01"
	paged := providerNameSet("paged paged-overwrite later-malformed later-http400 later-http403 later-http404 later-http405")
	for name := range providerLiteralCaches() {
		record, ok := output[name].(map[string]any)
		if !ok {
			t.Fatal("missing provider request record", name)
		}
		paths := []any{first}
		if name == "cancelled" {
			paths = []any{}
		} else if name == "cycle-stopped" {
			paths = []any{first, first, first, first}
		} else if paged[name] {
			paths = []any{first, second}
		}
		urls := []any{}
		for _, path := range paths {
			urls = append(urls, "GET https://provider.capture.invalid"+path.(string))
		}
		for _, key := range []string{"Requests", "CallerRequests"} {
			if !reflect.DeepEqual(record[key], urls) {
				t.Fatal("complete provider method/URI sequence mismatch", name, key)
			}
		}
		for _, key := range []string{"RawRequestURIs", "CallerRawRequestURIs"} {
			if !reflect.DeepEqual(record[key], paths) {
				t.Fatal("complete provider raw request sequence mismatch", name, key)
			}
		}
	}
}

func TestProviderSourceCallerRows(t *testing.T) {
	output := providerCapturedJSON(t, "source-provider-locations-outputs.json")
	available := providerNameSet("basic normalization null-location global case-field-alias duplicate-field paged")
	panicked := providerNameSet("null-provider null-type")
	rowText := "{\"SubscriptionID\":\"11111111-1111-1111-1111-111111111111\",\"SubscriptionName\":\"Synthetic\",\"SourceRegion\":\"eastus\",\"TargetRegion\":\"westeurope\",\"SourceResourceTypeCount\":1,\"AvailableTypes\":0,\"UnavailableTypes\":1,\"AvailabilityPercent\":0,\"AvgCostDifference\":0,\"HasCostData\":false,\"AvgLatencyMs\":0,\"LatencyEstimated\":false,\"MissingResourceTypes\":[\"microsoft.compute/virtualmachines\"],\"MissingSKUs\":[],\"RestrictedSKUs\":[],\"ZoneRestrictedSKUs\":[],\"TotalSKUsChecked\":0,\"AvailableSKUs\":0,\"UnavailableSKUs\":0,\"UnknownSKUs\":0,\"SKUAvailabilityPercent\":0,\"SourceZoneCount\":0,\"TargetZoneCount\":0,\"TargetZoneMappings\":null,\"Score\":0}"
	for name := range providerLiteralCaches() {
		record, ok := output[name].(map[string]any)
		if !ok {
			t.Fatal("missing provider caller record", name)
		}
		var want any
		if !panicked[name] {
			row := providerLiteral(t, rowText).(map[string]any)
			if available[name] {
				row["AvailableTypes"] = float64(1)
				row["UnavailableTypes"] = float64(0)
				row["AvailabilityPercent"] = float64(100)
				row["MissingResourceTypes"] = []any{}
			}
			want = []any{row}
		}
		if !reflect.DeepEqual(record["CallerRows"], want) || record["CallerPanicked"] != panicked[name] {
			t.Fatal("complete original provider caller comparison mismatch", name)
		}
	}
}
