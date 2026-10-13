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

const computeSourceVM = "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Compute/locations/eastus/vmSizes?api-version=2023-03-01"
const computeSourceDisk = "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Compute/skus?api-version=2023-03-01&$filter=location%20eq%20%27eastus%27"

func computeSourceJSON(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	value := map[string]any{}
	if json.Unmarshal(raw, &value) != nil {
		t.Fatal("invalid Compute source JSON")
	}
	return value
}

func computeSourceLiteral(t *testing.T, text string) any {
	t.Helper()
	var value any
	if json.Unmarshal([]byte(text), &value) != nil {
		t.Fatal("invalid independent Compute source literal")
	}
	return value
}

func computeSourceMaps() map[string]string {
	result := map[string]string{}
	set := func(names, value string) {
		for _, name := range strings.Fields(names) {
			result[name] = value
		}
	}
	set("vm-single vmss-single vm-single-alias vm-duplicate-name-last vm-duplicate-value-last vm-duplicate-case-name", `{"standard_d2":{"State":0,"BlockedZones":null}}`)
	set("vm-no-trim", `{" standard_d2 ":{"State":0,"BlockedZones":null}," ":{"State":0,"BlockedZones":null}}`)
	set("vm-surrogate-text", `{"\ufffd":{"State":0,"BlockedZones":null}}`)
	for _, provider := range []string{"vm", "vmss", "disk"} {
		for _, suffix := range strings.Fields("empty missing null-list null-envelope null-items") {
			result[provider+"-"+suffix] = "{}"
		}
		for _, suffix := range strings.Fields("wrong-list wrong-item wrong-name malformed wrong-next trailing") {
			result[provider+"-"+suffix] = "null"
		}
	}
	set("disk-name-tier-and-ignored-type", `{"premium_lrs":{"State":0,"BlockedZones":null},"standard_lrs":{"State":0,"BlockedZones":null}}`)
	set("disk-location-precedence disk-unicode-field-alias disk-alias-last-state disk-paged", `{"d":{"State":1,"BlockedZones":null}}`)
	set("disk-location-unknown-reason disk-first-location-wins disk-empty-zone-capability disk-duplicate-last-state", `{"d":{"State":2,"BlockedZones":null}}`)
	set("disk-zone-order-duplicates", `{"d":{"State":3,"BlockedZones":["2","1","2"]}}`)
	set("disk-zone-empty-string", `{"d":{"State":3,"BlockedZones":[""]}}`)
	set("disk-capability-mixed-true disk-null-restriction disk-next-foreign", `{"d":{"State":0,"BlockedZones":null}}`)
	set("disk-untrimmed-name", `{" d ":{"State":0,"BlockedZones":null}}`)
	set("disk-invalid-zones", "null")
	for _, provider := range []string{"vm", "disk"} {
		for _, suffix := range strings.Fields("http400 http403 http404 http429 http500 later-http400 later-http403 later-http404 later-http429 later-http500 later-malformed next-relative next-malformed bounded-cycle cancelled") {
			result[provider+"-"+suffix] = "null"
		}
	}
	set("vm-paged", `{"first":{"State":0,"BlockedZones":null},"second":{"State":0,"BlockedZones":null}}`)
	set("vm-next-foreign", `{"first":{"State":0,"BlockedZones":null}}`)
	set("pager-empty-url", "{}")
	return result
}

func computeSourceRecord(t *testing.T, value any, entries string, failed bool, urls []string) {
	t.Helper()
	requests := []any{}
	for _, url := range urls {
		pathQuery := strings.TrimPrefix(url, "https://management.azure.com")
		if strings.HasPrefix(url, "https://synthetic.invalid") {
			pathQuery = strings.TrimPrefix(url, "https://synthetic.invalid")
		}
		path, query, _ := strings.Cut(pathQuery, "?")
		requests = append(requests, map[string]any{
			"URL": url, "Path": path, "RawPath": "", "RawQuery": query,
			"Method": "GET", "EmptyBody": true, "SyntheticAuthorization": true,
		})
	}
	want := map[string]any{"Entries": computeSourceLiteral(t, entries), "Failed": failed, "Requests": requests}
	if !reflect.DeepEqual(value, want) {
		t.Fatal("complete Compute source record mismatch")
	}
}

func TestComputeSKUSourceMetadata(t *testing.T) {
	input := computeSourceJSON(t, "source-compute-sku-inputs.json")
	blobs := `{"provider.go":"7f71b296d332001ae12e49a021b6fee44f88a970","cache.go":"5be9c158fd38c609b06ffc67537cd07b838fc9d5","restrictions.go":"8d053daf58c7c8a18f39b579e95d4a53f8ebd55a","compute.go":"ee9300e5f8652591fe43121d0973321742ee70e3","storage.go":"b5a3a8d020c8143c64d296d48e8c51ae3b0ca2cf","cognitiveservices.go":"7b3dc63af21f12213004c6e91d9eaaa06c8a2189","sql.go":"1cb1e6430f11476e66e8bbb13bc0cc2ad5e0a43f","../types/types.go":"f9b61fb74075e5a2297be1b898067b8cb047f9b3","../../../../az/http_client.go":"05d85c5a50ee81f088731160020a914af83e1ca5"}`
	if !reflect.DeepEqual(input["Blobs"], computeSourceLiteral(t, blobs)) {
		t.Fatal("Compute source blob authority mismatch")
	}
	cases, ok := input["Cases"].(map[string]any)
	if !ok || len(cases) != 89 {
		t.Fatal("Compute source input matrix mismatch")
	}
	for name := range computeSourceMaps() {
		if name != "pager-empty-url" && cases[name] == nil {
			t.Fatal("Compute source input missing", name)
		}
	}
	raw, err := os.ReadFile("testdata/source-compute-sku-inputs.json")
	if err != nil || fmt.Sprintf("%x", sha256.Sum256(raw)) != "75a69648317d4f53615107172255b30853b77b6f2185f421bcc6eb01e1c7ba12" {
		t.Fatal("original complete Compute source input bytes changed")
	}
}

func TestComputeSKUSourceCells(t *testing.T) {
	output := computeSourceJSON(t, "source-compute-sku-outputs.json")
	cases, ok := output["Cases"].(map[string]any)
	want := computeSourceMaps()
	if !ok || len(cases) != 90 || len(want) != 90 || len(output) != 9 {
		t.Fatal("complete Compute source output matrix mismatch")
	}
	for name, entries := range want {
		record, ok := cases[name].(map[string]any)
		if !ok || len(record) != 3 || !reflect.DeepEqual(record["Entries"], computeSourceLiteral(t, entries)) || record["Failed"] != (entries == "null") {
			t.Fatal("independent complete Compute source map/state mismatch", name)
		}
	}
}

func TestComputeSKUSourceRequests(t *testing.T) {
	output := computeSourceJSON(t, "source-compute-sku-outputs.json")
	cases, ok := output["Cases"].(map[string]any)
	if !ok || len(cases) != 90 {
		t.Fatal("Compute source request matrix missing")
	}
	for name, entries := range computeSourceMaps() {
		start := computeSourceVM
		if strings.HasPrefix(name, "disk-") {
			start = computeSourceDisk
		}
		urls := []string{start}
		if strings.HasSuffix(name, "-cancelled") || name == "pager-empty-url" {
			urls = []string{}
		} else if strings.HasSuffix(name, "-bounded-cycle") {
			urls = []string{start, start, start}
		} else if strings.Contains(name, "-later-") || strings.HasSuffix(name, "-paged") {
			urls = []string{start, start + "&$skiptoken=synthetic"}
		} else if strings.HasSuffix(name, "-next-foreign") {
			urls = []string{start, "https://synthetic.invalid/compute/next"}
		}
		t.Run(name, func(t *testing.T) {
			computeSourceRecord(t, cases[name], entries, entries == "null", urls)
		})
	}
}

func TestComputeSKUSourceCache(t *testing.T) {
	output := computeSourceJSON(t, "source-compute-sku-outputs.json")
	registry := `["microsoft.cognitiveservices/accounts","microsoft.compute/disks","microsoft.compute/virtualmachines","microsoft.compute/virtualmachinescalesets","microsoft.sql/managedinstances","microsoft.sql/servers/databases","microsoft.storage/storageaccounts"]`
	if !reflect.DeepEqual(output["Registry"], computeSourceLiteral(t, registry)) || output["DuplicateRegistrationPanics"] != true {
		t.Fatal("independent Compute registry mismatch")
	}
	cache, ok := output["Cache"].(map[string]any)
	if !ok || len(cache) != 7 {
		t.Fatal("Compute cache matrix missing")
	}
	initial := `{"standard_d2":{"State":0,"BlockedZones":null},"":{"State":0,"BlockedZones":null}}`
	injected := `{"standard_d2":{"State":0,"BlockedZones":null},"":{"State":0,"BlockedZones":null},"injected":{"State":1,"BlockedZones":null}}`
	allURLs := []string{computeSourceVM, computeSourceVM, strings.Replace(computeSourceVM, "11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222", 1), strings.Replace(computeSourceVM, "eastus", "westus", 1)}
	computeSourceRecord(t, cache["first"], initial, false, allURLs[:1])
	computeSourceRecord(t, cache["normalized-reuse"], initial, false, allURLs[:1])
	computeSourceRecord(t, cache["caller-map-mutation-shared"], injected, false, allURLs[:1])
	computeSourceRecord(t, cache["vmss-separate"], `{"standard_d2":{"State":0,"BlockedZones":null}}`, false, allURLs[:2])
	computeSourceRecord(t, cache["subscription-separate"], `{"other":{"State":0,"BlockedZones":null}}`, false, allURLs[:3])
	computeSourceRecord(t, cache["region-separate"], `{"west":{"State":0,"BlockedZones":null}}`, false, allURLs)
	computeSourceRecord(t, cache["unsupported"], "null", true, allURLs)
	retry, ok := output["Retry"].([]any)
	if !ok || len(retry) != 3 {
		t.Fatal("Compute retry matrix missing")
	}
	computeSourceRecord(t, retry[0], "null", true, []string{computeSourceVM})
	for _, index := range []int{1, 2} {
		computeSourceRecord(t, retry[index], `{"standard_d2":{"State":0,"BlockedZones":null}}`, false, []string{computeSourceVM, computeSourceVM})
	}
	if !reflect.DeepEqual(output["ZoneBefore"], computeSourceLiteral(t, `{"State":3,"BlockedZones":["2","1","2"]}`)) || !reflect.DeepEqual(output["ZoneAfter"], computeSourceLiteral(t, `{"State":3,"BlockedZones":["changed","1","2"]}`)) {
		t.Fatal("source zone slice sharing mismatch")
	}
	computeSourceRecord(t, output["ConcurrentSameKey"], `{"standard_d2":{"State":0,"BlockedZones":null}}`, false, []string{computeSourceVM})
	if !reflect.DeepEqual(output["ConcurrentSeparateKeys"], computeSourceLiteral(t, `{"east":{"east":{"State":0,"BlockedZones":null}},"west":{"west":{"State":0,"BlockedZones":null}}}`)) {
		t.Fatal("complete concurrent separate-key maps mismatch")
	}
}
