package azure

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func TestARMOptionsDoNotAutoRegisterProvidersOnReadFailure(t *testing.T) {
	calls := 0
	writes := 0
	options := NewARMClientOptions()
	options.Retry.MaxRetries = -1
	options.Transport = operationTransportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		status := http.StatusConflict
		body := `{"error":{"code":"MissingSubscriptionRegistration","message":"synthetic registration required"}}`
		if r.Method != http.MethodGet {
			writes++
			status = 400
			body = `{"error":{"code":"SyntheticWriteRejected"}}`
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Request: r, Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	pipeline, err := armruntime.NewPipeline("offline-qa", "test", &mockCredential{token: "arm-registration-canary"}, runtime.PipelineOptions{}, options)
	if err != nil {
		t.Fatal(err)
	}
	request, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://management.azure.com/subscriptions/sub-fixture/resourceGroups/rg/providers/Microsoft.Test/widgets/one")
	if err != nil {
		t.Fatal(err)
	}
	response, err := pipeline.Do(request)
	if err != nil || response == nil || response.StatusCode != 409 || calls != 1 || writes != 0 {
		t.Fatalf("read failure triggered extra action: error=%v calls=%d writes=%d", err, calls, writes)
	}
	response.Body.Close()
}
