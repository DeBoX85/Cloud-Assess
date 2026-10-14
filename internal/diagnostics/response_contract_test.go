package diagnostics

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
)

type batchFramingBody struct {
	*strings.Reader
	closes int
}

func (b *batchFramingBody) Close() error {
	b.closes++
	return nil
}

func TestDiagnosticsRejectsTrailingBatchResponseContent(t *testing.T) {
	for _, suffix := range []string{"", " \r\n\t", "{}", "null", "junk", "{"} {
		t.Run("suffix-"+suffix, func(t *testing.T) {
			body := &batchFramingBody{Reader: strings.NewReader(`{"responses":[{"httpStatusCode":200,"content":{"value":[]}}]}` + suffix)}
			client := &fakeBatchClient{handler: func(armBatchRequest) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
			}}
			resource := assessment.Resource{ID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/one", Type: "Microsoft.Storage/storageAccounts"}
			result, err := NewWithClient(client, "https://management.azure.com").Scan(context.Background(), []assessment.Resource{resource}, nil, nil)
			if body.closes != 1 {
				t.Fatalf("batch body closes=%d, want one", body.closes)
			}
			valid := suffix == "" || suffix == " \r\n\t"
			if valid {
				if err != nil || len(result.Findings) != 1 {
					t.Fatalf("valid batch framing rejected: result=%+v error=%v", result, err)
				}
			} else if err == nil || len(result.Findings) != 0 {
				t.Fatalf("trailing batch response content became assessment findings: suffix=%q result=%+v error=%v", suffix, result, err)
			}
		})
	}
}

func TestDiagnosticsRejectsSettingsOutsideRequestedBatch(t *testing.T) {
	const one = "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/one"
	const two = "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/two"
	for _, foreign := range []bool{false, true} {
		t.Run(map[bool]string{false: "owned", true: "foreign"}[foreign], func(t *testing.T) {
			client := &fakeBatchClient{handler: func(request armBatchRequest) (*http.Response, error) {
				id := strings.Split(request.Requests[0].RelativeURL, "/providers/microsoft.insights/diagnosticSettings?")[0]
				if foreign && id == one {
					id = two
				}
				content, err := json.Marshal(diagnosticSettingsPayload{Value: []diagnosticSetting{{ID: strings.ToUpper(id) + "/providers/Microsoft.Insights/diagnosticSettings/default"}}})
				if err != nil {
					t.Fatal(err)
				}
				return jsonResponse(t, armBatchResponse{Responses: []armBatchResponseItem{{HTTPStatusCode: http.StatusOK, Content: content}}}), nil
			}}
			scanner := NewWithClient(client, "https://management.azure.com")
			scanner.batchSize, scanner.maxWorkers = 1, 1
			result, err := scanner.Scan(context.Background(), []assessment.Resource{{ID: one, Type: "Microsoft.Storage/storageAccounts"}, {ID: two, Type: "Microsoft.Storage/storageAccounts"}}, nil, nil)
			if foreign {
				if err == nil || len(result.Findings) != 0 {
					t.Fatalf("foreign diagnostic setting accepted for requested batch: result=%+v error=%v", result, err)
				}
			} else if err != nil || len(result.Findings) != 0 {
				t.Fatalf("owned case-insensitive diagnostic setting rejected: result=%+v error=%v", result, err)
			}
		})
	}
}

var _ io.ReadCloser = (*batchFramingBody)(nil)
