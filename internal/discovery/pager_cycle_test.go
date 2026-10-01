package discovery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

type cycleTransport struct {
	links  []string
	calls  int
	repeat bool
}

func (f *cycleTransport) Do(r *http.Request) (*http.Response, error) {
	f.calls++
	if f.calls > 6 {
		return nil, fmt.Errorf("fixture runaway tripwire")
	}
	link := ""
	if f.repeat {
		link = f.links[(f.calls-1)%len(f.links)]
	} else if f.calls <= len(f.links) {
		link = f.links[f.calls-1]
	}
	body := `{"value":[{"subscriptionId":"sub-fixture","name":"sub-fixture","type":"Microsoft.Management/managementGroups"}]}`
	if link != "" {
		body = `{"value":[{"subscriptionId":"sub-fixture","name":"sub-fixture","type":"Microsoft.Management/managementGroups"}],"nextLink":"` + link + `"}`
	}
	return &http.Response{StatusCode: 200, Header: make(http.Header), Request: r, Body: io.NopCloser(strings.NewReader(body))}, nil
}
func TestScopePagerCyclesFailWithoutPartialRows(t *testing.T) {
	for _, op := range []string{"subscriptions", "group-subscriptions", "descendants"} {
		for _, links := range [][]string{{"https://management.azure.com/page?secret=a"}, {"https://management.azure.com/page?a=1", "https://management.azure.com/page?a=2"}} {
			t.Run(op+fmt.Sprint(len(links)), func(t *testing.T) {
				transport := &cycleTransport{links: links, repeat: true}
				options := &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
				options.Retry.MaxRetries = -1
				client, err := NewAzureScopeClient(fixtureCredential{}, options)
				if err != nil {
					t.Fatal(err)
				}
				var rows any
				switch op {
				case "subscriptions":
					var v []Subscription
					v, err = client.ListSubscriptions(context.Background())
					if v != nil {
						rows = v
					}
				case "group-subscriptions":
					var v []Subscription
					v, err = client.SubscriptionsUnderManagementGroup(context.Background(), "fixture")
					if v != nil {
						rows = v
					}
				case "descendants":
					var v []string
					v, err = client.DescendantManagementGroups(context.Background(), "fixture")
					if v != nil {
						rows = v
					}
				}
				if err == nil || !strings.Contains(err.Error(), "repeated a continuation") || transport.calls != len(links)+1 || rows != nil {
					t.Fatalf("error=%v calls=%d rows=%v", err, transport.calls, rows)
				}
				if strings.Contains(err.Error(), "secret=") {
					t.Fatal("continuation leaked in error")
				}
			})
		}
	}
}

func TestScopeContinuationHostCaseAndOpaqueQuery(t *testing.T) {
	seen := map[string]struct{}{}
	a := "https://management.azure.com/Page?token=A"
	b := "https://MANAGEMENT.AZURE.COM/Page?token=A"
	c := "https://management.azure.com/Page?token=a"
	if checkScopeContinuation(&a, seen) != nil || checkScopeContinuation(&b, seen) == nil || checkScopeContinuation(&c, seen) != nil {
		t.Fatal("host case or opaque query mishandled")
	}
}
func TestScopePagerHealthyFinitePages(t *testing.T) {
	for _, op := range []string{"subscriptions", "group-subscriptions", "descendants"} {
		t.Run(op, func(t *testing.T) {
			transport := &cycleTransport{links: []string{"https://management.azure.com/page?token=1", "https://management.azure.com/page?token=2"}}
			client, err := NewAzureScopeClient(fixtureCredential{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}})
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				transport.calls = 0
				count := 0
				switch op {
				case "subscriptions":
					var rows []Subscription
					rows, err = client.ListSubscriptions(context.Background())
					count = len(rows)
				case "group-subscriptions":
					var rows []Subscription
					rows, err = client.SubscriptionsUnderManagementGroup(context.Background(), "fixture")
					count = len(rows)
				case "descendants":
					var rows []string
					rows, err = client.DescendantManagementGroups(context.Background(), "fixture")
					count = len(rows)
				}
				if err != nil || count != 3 || transport.calls != 3 {
					t.Fatalf("healthy pagination: error=%v rows=%d calls=%d", err, count, transport.calls)
				}
			}
		})
	}
}
func TestScopePagerCanceledBeforeListingSendsNoRequest(t *testing.T) {
	transport := &cycleTransport{}
	client, err := NewAzureScopeClient(fixtureCredential{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, op := range []string{"subscriptions", "group-subscriptions", "descendants"} {
		switch op {
		case "subscriptions":
			_, err = client.ListSubscriptions(ctx)
		case "group-subscriptions":
			_, err = client.SubscriptionsUnderManagementGroup(ctx, "fixture")
		case "descendants":
			_, err = client.DescendantManagementGroups(ctx, "fixture")
		}
		if !errors.Is(err, context.Canceled) || transport.calls != 0 {
			t.Fatalf("canceled listing: error=%v calls=%d", err, transport.calls)
		}
	}
}
