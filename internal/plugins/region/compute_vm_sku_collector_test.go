package region

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

const vmTestID = "11111111-1111-1111-1111-111111111111"
const vmTestType = "microsoft.compute/virtualmachines"
const vmTestURL = "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Compute/skus?%24filter=location+eq+%27eastus%27&api-version=2021-07-01"

type vmGetterFunc func(context.Context, string, int64) ([]byte, *http.Response, error)

func (f vmGetterFunc) GetBoundedWithResponse(ctx context.Context, u string, n int64) ([]byte, *http.Response, error) {
	return f(ctx, u, n)
}

type vmManifest struct {
	Primary []struct {
		File, SHA256 string
		Bytes        int
		Blob         string `json:"source_blob"`
	} `json:"primary_files"`
	Cases []struct {
		Name, Context string
		Request       struct {
			SubscriptionID string `json:"subscription_id"`
			Region         string
			ResourceType   string   `json:"resource_type"`
			RequiredNames  []string `json:"required_names"`
		}
		Pages []struct {
			Status      int
			Body        string
			BodyFile    string `json:"body_file"`
			BodyPointer string `json:"body_pointer"`
		}
		Expected struct {
			Status   *string `json:"collection_status"`
			Code     string  `json:"failure_code"`
			Groups   map[string]ComputeVMSKUGroup
			Evidence *SKUEvidence
			Calls    int `json:"getter_calls"`
			Error    string
			URLs     []string `json:"request_urls"`
		}
	}
}

func vmRead(t *testing.T, name string) []byte {
	t.Helper()
	b, e := os.ReadFile("testdata/" + name)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestComputeVMSKUFrozenManifest(t *testing.T) {
	var manifest vmManifest
	if e := json.Unmarshal(vmRead(t, "compute-vm-target-contract-v1.json"), &manifest); e != nil {
		t.Fatal(e)
	}
	if len(manifest.Cases) != 67 {
		t.Fatal("frozen scenario count", len(manifest.Cases))
	}
	seen := map[string]bool{}
	for _, fixture := range manifest.Cases {
		if seen[fixture.Name] {
			t.Fatal("duplicate scenario")
		}
		seen[fixture.Name] = true
		t.Run(fixture.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if fixture.Context == "cancelled-before" {
				cancel()
			}
			calls := 0
			getter := vmGetterFunc(func(_ context.Context, u string, limit int64) ([]byte, *http.Response, error) {
				if calls >= len(fixture.Pages) || calls >= len(fixture.Expected.URLs) {
					t.Fatalf("unexpected getter %d %s", calls, u)
				}
				if u != fixture.Expected.URLs[calls] || limit != 1048576 {
					t.Fatalf("literal request/bound mismatch %q/%d", u, limit)
				}
				p := fixture.Pages[calls]
				calls++
				if fixture.Context == "deadline-after-get" {
					return nil, nil, fmt.Errorf("opaque canary: %w", context.DeadlineExceeded)
				}
				body := []byte(p.Body)
				if p.BodyFile != "" {
					if !strings.HasPrefix(p.BodyFile, "compute-primary-2021-07-01-") || p.BodyPointer != "/responses/200/body" {
						t.Fatal("unexpected fixture path")
					}
					var primary struct {
						Responses map[string]struct{ Body json.RawMessage }
					}
					if e := json.Unmarshal(vmRead(t, p.BodyFile), &primary); e != nil {
						t.Fatal(e)
					}
					body = primary.Responses["200"].Body
				}
				return body, &http.Response{StatusCode: p.Status}, nil
			})
			collector, e := NewComputeVMSKUCollector("https://management.azure.com", getter)
			if e != nil {
				t.Fatal(e)
			}
			q := ComputeVMSKURequest{SubscriptionID: fixture.Request.SubscriptionID, Region: fixture.Request.Region, ResourceType: fixture.Request.ResourceType}
			got, e := collector.Collect(ctx, map[string]string{vmTestID: "selected"}, q)
			if calls != fixture.Expected.Calls {
				t.Fatalf("calls %d want %d", calls, fixture.Expected.Calls)
			}
			if fixture.Expected.Error != "" {
				if got != nil || e == nil {
					t.Fatalf("nil result/error expected: %#v/%v", got, e)
				}
				switch fixture.Expected.Error {
				case "context.Canceled":
					if !errors.Is(e, context.Canceled) {
						t.Fatal(e)
					}
				case "context.DeadlineExceeded":
					if !errors.Is(e, context.DeadlineExceeded) {
						t.Fatal(e)
					}
				default:
					if !strings.Contains(e.Error(), "["+fixture.Expected.Error+"]") {
						t.Fatal(e)
					}
				}
				return
			}
			if got == nil || e != nil || fixture.Expected.Status == nil || got.Status != *fixture.Expected.Status || got.FailureCode != fixture.Expected.Code || got.SubscriptionID != vmTestID || got.Region != q.Region || got.ResourceType != q.ResourceType {
				t.Fatalf("collection binding/status mismatch %#v/%v", got, e)
			}
			if !reflect.DeepEqual(got.Groups, fixture.Expected.Groups) {
				t.Fatalf("full groups %#v want %#v", got.Groups, fixture.Expected.Groups)
			}
			evidence, e := got.Evidence(context.Background(), fixture.Request.RequiredNames)
			if e != nil || fixture.Expected.Evidence == nil || !reflect.DeepEqual(evidence, *fixture.Expected.Evidence) {
				t.Fatalf("full evidence %#v/%v want %#v", evidence, e, fixture.Expected.Evidence)
			}
		})
	}
}

func TestComputeVMSKUPrimaryIntegrity(t *testing.T) {
	var m vmManifest
	if e := json.Unmarshal(vmRead(t, "compute-vm-target-contract-v1.json"), &m); e != nil {
		t.Fatal(e)
	}
	literal := map[string]string{
		"compute-primary-2021-07-01-skus.json":                                           "81f9bb91b33be77111ef82ffa1549c633cf0b50edb6907ba2ca9e4d9f484fe59",
		"compute-primary-2021-07-01-ListAvailableResourceSkus.json":                      "e4cbaed858756360cf8649e5d6c09006600c8df3ff0665b1bea8f0e02e6da314",
		"compute-primary-2021-07-01-ListAvailableResourceSkusForARegion.json":            "11fa69a6a6a3d062f32f780fb86a0ec6153433a05188d198c4216e2f3e69c083",
		"compute-primary-2021-07-01-ListAvailableResourceSkusWithExtendedLocations.json": "fe0dd2a9f44d64b4d2b9c7828a4c945d0827ed6f181b478d32e389635e43ed26",
	}
	if len(m.Primary) != 4 {
		t.Fatal("primary file count")
	}
	for _, p := range m.Primary {
		b := vmRead(t, p.File)
		h := sha256.Sum256(b)
		if hex.EncodeToString(h[:]) != literal[p.File] || p.SHA256 != literal[p.File] || len(b) != p.Bytes {
			t.Fatal("immutable primary bytes changed", p.File)
		}
	}
}
