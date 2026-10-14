package cost

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DeBoX85/Cloud-Assess/internal/azure"
)

func TestCostHTTPRequiresDocumentedQueryCompletionStatus(t *testing.T) {
	for _, status := range []int{200, 201, 202, 203, 204, 205, 206} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			options := azure.DefaultHTTPClientOptions(time.Second)
			options.MaxRetries = -1
			options.Scope = "https://management.azure.com/.default"
			options.Transport = costContractTransport(func(request *http.Request) (*http.Response, error) {
				payload := `{"properties":{"rows":[[1,"Storage","USD"]]}}`
				if status == http.StatusNoContent {
					payload = ""
				}
				return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(payload)), Header: http.Header{}, Request: request}, nil
			})
			scanner := NewWithClient(azure.NewHTTPClient(costContractCredential{}, options), "https://management.azure.com")
			result, err := scanner.scanAt(context.Background(), map[string]string{"sub": "Synthetic"}, time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC))
			if status == http.StatusOK {
				if err != nil || len(result.Records) != 1 {
					t.Fatalf("documented cost response failed: result=%+v error=%v", result, err)
				}
			} else if status == http.StatusNoContent {
				if err != nil || len(result.Records) != 0 {
					t.Fatalf("documented no-content cost response failed: result=%+v error=%v", result, err)
				}
			} else if err == nil || len(result.Records) != 0 {
				t.Fatalf("undocumented Cost status admitted as completed query: status=%d result=%+v error=%v", status, result, err)
			}
		})
	}
}
