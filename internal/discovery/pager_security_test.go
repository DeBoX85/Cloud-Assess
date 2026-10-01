package discovery

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

type continuationTransport struct {
	link                 string
	calls, foreignBearer int
}

func (f *continuationTransport) Do(r *http.Request) (*http.Response, error) {
	f.calls++
	if r.URL.Host != "management.azure.com" && r.Header.Get("Authorization") != "" {
		f.foreignBearer++
	}
	body := `{"value":[]}`
	if f.calls == 1 {
		body = `{"value":[],"nextLink":"` + f.link + `"}`
	}
	return &http.Response{StatusCode: 200, Header: make(http.Header), Request: r, Body: io.NopCloser(strings.NewReader(body))}, nil
}
func TestScopePagersConstrainContinuations(t *testing.T) {
	for _, operation := range []string{"subscriptions", "group-subscriptions", "descendants"} {
		for _, link := range []string{"https://foreign.invalid/page", "https://management.azure.com:444/page", "https://user@management.azure.com/page", "https://management.azure.com/page#fragment", "https://MANAGEMENT.azure.com/page?skip=2"} {
			t.Run(operation+"/"+link, func(t *testing.T) {
				transport := &continuationTransport{link: link}
				client, err := NewAzureScopeClient(fixtureCredential{}, &arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}})
				if err != nil {
					t.Fatal(err)
				}
				switch operation {
				case "subscriptions":
					_, err = client.ListSubscriptions(context.Background())
				case "group-subscriptions":
					_, err = client.SubscriptionsUnderManagementGroup(context.Background(), "fixture")
				case "descendants":
					_, err = client.DescendantManagementGroups(context.Background(), "fixture")
				}
				safe := strings.Contains(link, "MANAGEMENT")
				if safe {
					if err != nil || transport.calls != 2 {
						t.Fatalf("safe continuation: error=%v calls=%d", err, transport.calls)
					}
					return
				}
				if err == nil || transport.calls != 1 || transport.foreignBearer != 0 {
					t.Fatalf("unsafe continuation: error=%v calls=%d foreignBearer=%d", err, transport.calls, transport.foreignBearer)
				}
				if strings.Contains(err.Error(), "foreign.invalid") {
					t.Fatal("error disclosed supplied URL")
				}
			})
		}
	}
}
