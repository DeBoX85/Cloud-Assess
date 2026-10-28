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
	"sort"
	"strings"
	"testing"
)

const storageTestID = "11111111-1111-1111-1111-111111111111"

type storageGetterFunc func(context.Context, string, int64) ([]byte, *http.Response, error)

func (f storageGetterFunc) GetBoundedWithResponse(ctx context.Context, url string, limit int64) ([]byte, *http.Response, error) {
	return f(ctx, url, limit)
}

type storageTerminalContext struct {
	context.Context
	done chan struct{}
	err  error
}

func (c *storageTerminalContext) Done() <-chan struct{} { return c.done }
func (c *storageTerminalContext) Err() error            { return c.err }
func (c *storageTerminalContext) finish(err error) {
	c.err = err
	close(c.done)
}

type storageManifest struct {
	Cases []struct {
		Name    string
		Request struct {
			SubscriptionID  string              `json:"subscription_id"`
			Region          string              `json:"region"`
			LocationAliases map[string]string   `json:"location_aliases"`
			RequiredKinds   map[string][]string `json:"required_kinds"`
		}
		Response struct {
			Status      *int
			Body        string
			BodyFile    string `json:"body_file"`
			BodyPointer string `json:"body_pointer"`
			Error       string
		}
		Context  string
		Expected struct {
			CollectionStatus *string `json:"collection_status"`
			FailureCode      string  `json:"failure_code"`
			Groups           map[string]map[string]StorageSKUGroup
			UnqualifiedNames []string `json:"unqualified_names"`
			Evidence         *struct {
				Status string
				Values map[string]struct {
					State        int
					BlockedZones []string `json:"blocked_zones"`
				}
			}
			GetterCalls int `json:"getter_calls"`
			Error       string
		}
	}
}

func storageRead(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestStorageSKUFrozenManifest(t *testing.T) {
	var manifest storageManifest
	if err := json.Unmarshal(storageRead(t, "storage-target-contract-v1.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Cases) != 71 {
		t.Fatal("frozen scenario count changed", len(manifest.Cases))
	}
	seen := map[string]bool{}
	for _, fixture := range manifest.Cases {
		if seen[fixture.Name] {
			t.Fatal("duplicate scenario", fixture.Name)
		}
		seen[fixture.Name] = true
		t.Run(fixture.Name, func(t *testing.T) {
			body := []byte(fixture.Response.Body)
			if fixture.Response.BodyFile != "" {
				if fixture.Response.BodyFile != "storage-primary-2023-01-01-SKUList.json" || fixture.Response.BodyPointer != "/responses/200/body" {
					t.Fatal("unexpected primary fixture path")
				}
				var primary struct {
					Responses map[string]struct{ Body json.RawMessage }
				}
				if err := json.Unmarshal(storageRead(t, fixture.Response.BodyFile), &primary); err != nil {
					t.Fatal(err)
				}
				body = primary.Responses["200"].Body
			}
			ctx := &storageTerminalContext{Context: context.Background(), done: make(chan struct{})}
			if fixture.Context == "cancelled-before" {
				ctx.finish(context.Canceled)
			}
			calls := 0
			getter := storageGetterFunc(func(_ context.Context, target string, limit int64) ([]byte, *http.Response, error) {
				calls++
				if calls != 1 || target != "https://management.azure.com/subscriptions/"+storageTestID+"/providers/Microsoft.Storage/skus?api-version=2023-01-01" || limit != 1048576 {
					t.Fatalf("literal request/bound mismatch: %q limit=%d calls=%d", target, limit, calls)
				}
				if fixture.Context == "deadline-after-get" {
					ctx.finish(context.DeadlineExceeded)
				}
				var response *http.Response
				if fixture.Response.Status != nil {
					response = &http.Response{StatusCode: *fixture.Response.Status}
				}
				var err error
				switch fixture.Response.Error {
				case "":
				case "context.Canceled":
					err = fmt.Errorf("wrapped canary: %w", context.Canceled)
				default:
					err = errors.New(fixture.Response.Error)
				}
				return body, response, err
			})
			collector, err := NewStorageSKUCollector("https://management.azure.com", getter)
			if err != nil {
				t.Fatal(err)
			}
			request := StorageSKURequest{SubscriptionID: fixture.Request.SubscriptionID, Region: fixture.Request.Region, LocationAliases: fixture.Request.LocationAliases}
			got, err := collector.Collect(ctx, map[string]string{storageTestID: "selected"}, request)
			if calls != fixture.Expected.GetterCalls {
				t.Fatalf("getter count=%d want %d", calls, fixture.Expected.GetterCalls)
			}
			if fixture.Expected.Error != "" {
				if got != nil || err == nil {
					t.Fatalf("expected nil collection/error, got %#v / %v", got, err)
				}
				switch fixture.Expected.Error {
				case "context.Canceled":
					if !errors.Is(err, context.Canceled) {
						t.Fatal(err)
					}
				case "context.DeadlineExceeded":
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Fatal(err)
					}
				default:
					if !strings.Contains(err.Error(), "["+fixture.Expected.Error+"]") {
						t.Fatal(err)
					}
				}
				return
			}
			if got == nil || err != nil || fixture.Expected.CollectionStatus == nil || got.Status != *fixture.Expected.CollectionStatus || got.FailureCode != fixture.Expected.FailureCode || got.SubscriptionID != storageTestID || got.Region != request.Region {
				t.Fatalf("collection/status mismatch: %#v / %v", got, err)
			}
			if !reflect.DeepEqual(got.Groups, fixture.Expected.Groups) {
				t.Fatalf("complete grouped map mismatch: %#v want %#v", got.Groups, fixture.Expected.Groups)
			}
			names := []string{}
			for name, present := range got.UnqualifiedNames {
				if !present {
					t.Fatal("false unqualified name")
				}
				names = append(names, name)
			}
			sort.Strings(names)
			if !reflect.DeepEqual(names, fixture.Expected.UnqualifiedNames) {
				t.Fatalf("unqualified names=%v want %v", names, fixture.Expected.UnqualifiedNames)
			}
			evidence, err := got.Evidence(context.Background(), fixture.Request.RequiredKinds)
			want := SKUEvidence{Status: fixture.Expected.Evidence.Status, Values: map[string]SKUAvailability{}}
			for name, value := range fixture.Expected.Evidence.Values {
				want.Values[name] = SKUAvailability{State: value.State, BlockedZones: append([]string{}, value.BlockedZones...)}
			}
			if err != nil || !reflect.DeepEqual(evidence, want) {
				t.Fatalf("entire projection mismatch: %#v / %v want %#v", evidence, err, want)
			}
			if strings.Contains(got.FailureCode, "canary") {
				t.Fatal("opaque diagnostic leaked")
			}
		})
	}
}

func TestStorageSKUPrimaryFixtureIntegrity(t *testing.T) {
	data := storageRead(t, "storage-primary-2023-01-01-SKUList.json")
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != "aceb9c06f233a6db9b4adb726e8f7951037531887dd28c949a24dee5b354830e" {
		t.Fatal("immutable primary example bytes changed")
	}
	var primary struct {
		Responses map[string]struct {
			Body struct {
				Value []struct {
					Name, Kind string
					Locations  []string
				}
			}
		}
	}
	if err := json.Unmarshal(data, &primary); err != nil {
		t.Fatal(err)
	}
	rows := primary.Responses["200"].Body.Value
	names, groups, locations := map[string]bool{}, map[string]int{}, map[string]bool{}
	for _, row := range rows {
		names[row.Name] = true
		groups[row.Name+"/"+row.Kind]++
		for _, location := range row.Locations {
			locations[location] = true
		}
	}
	if len(rows) != 248 || len(names) != 5 || len(groups) != 8 || len(locations) != 31 || !locations["eastus2(stage)"] {
		t.Fatal("literal primary premises changed")
	}
	for _, count := range groups {
		if count != 31 {
			t.Fatal("primary group does not retain all31regional rows")
		}
	}
}
