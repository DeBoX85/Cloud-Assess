// diagnostics-probe helps investigate non-successful ARM diagnostic-settings
// batch subrequests. It performs individual GET requests only.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/DeBoX85/Cloud-Assess/internal/assessment"
	"github.com/DeBoX85/Cloud-Assess/internal/azure"
	"github.com/DeBoX85/Cloud-Assess/internal/diagnostics"
)

const diagnosticSettingsAPI = "2021-05-01-preview"

type getter interface {
	Get(context.Context, string) ([]byte, error)
}

type failure struct {
	ResourceIDHash string `json:"resourceIdHash"`
	ResourceType   string `json:"resourceType"`
	HTTPStatus     int    `json:"httpStatus"`
	AzureErrorCode string `json:"azureErrorCode,omitempty"`
	ErrorClass     string `json:"errorClass,omitempty"`
}

type summary struct {
	OriginalHTTP400Warnings int       `json:"originalHTTP400Warnings"`
	EligibleResources       int       `json:"eligibleResources"`
	SuccessfulGETs          int       `json:"successfulGets"`
	Failures               []failure `json:"failures"`
}

var safeErrorCode = regexp.MustCompile(`^[A-Za-z0-9._-]{1,80}$`)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("diagnostics-probe", flag.ContinueOnError)
	flags.SetOutput(stderr)
	targetPath := flags.String("target", "", "Unredacted Cloud Assess target.json from a live equivalence pass")
	maxRequests := flags.Int("max-requests", 100, "Maximum individual read-only GET requests")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *targetPath == "" || flags.NArg() != 0 || *maxRequests < 1 {
		fmt.Fprintln(stderr, "provide --target and a positive --max-requests value, with no positional arguments")
		return 2
	}

	file, err := os.Open(*targetPath) //nolint:gosec // caller provides the local evidence path
	if err != nil {
		fmt.Fprintf(stderr, "open target report: %v\n", err)
		return 2
	}
	defer file.Close()
	var target struct {
		SchemaVersion string                      `json:"schemaVersion"`
		Resources     []assessment.Resource       `json:"resources"`
		Stages        []assessment.StageExecution `json:"stages"`
	}
	if err := json.NewDecoder(file).Decode(&target); err != nil || target.SchemaVersion == "" {
		fmt.Fprintln(stderr, "target is not a valid canonical Cloud Assess JSON report")
		return 2
	}
	eligible := eligibleResources(target.Resources)
	if len(eligible) == 0 {
		fmt.Fprintln(stderr, "target report contains no inventory resources supported by the Diagnostics scanner")
		return 2
	}
	if len(eligible) > *maxRequests {
		fmt.Fprintf(stderr, "%d eligible resources exceed --max-requests=%d; raise the limit deliberately if this scope is intended\n", len(eligible), *maxRequests)
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	credential, err := azure.NewCredential()
	if err != nil {
		fmt.Fprintln(stderr, "create Azure credential failed; use the same authenticated environment as the live scan")
		return 1
	}
	client := azure.NewHTTPClient(credential, azure.DefaultHTTPClientOptions(60*time.Second))
	result := probe(ctx, client, azure.ResourceManagerEndpoint(), eligible, originalHTTP400Warnings(target.Stages))
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintln(stderr, "write probe summary failed")
		return 1
	}
	return 0
}

func eligibleResources(resources []assessment.Resource) []assessment.Resource {
	eligible := make([]assessment.Resource, 0, len(resources))
	for _, resource := range resources {
		if diagnostics.Supports(resource.Type) {
			eligible = append(eligible, resource)
		}
	}
	slices.SortFunc(eligible, func(a, b assessment.Resource) int {
		return strings.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID))
	})
	return eligible
}

func originalHTTP400Warnings(stages []assessment.StageExecution) int {
	for _, stage := range stages {
		if stage.Name != "diagnostics" {
			continue
		}
		count := 0
		for _, warning := range stage.Warnings {
			if warning.Code == "diagnostics_subrequest_non_success" && strings.HasSuffix(warning.Message, "HTTP 400") {
				count++
			}
		}
		return count
	}
	return 0
}

func probe(ctx context.Context, client getter, endpoint string, resources []assessment.Resource, originalWarnings int) summary {
	result := summary{
		OriginalHTTP400Warnings: originalWarnings,
		EligibleResources:       len(resources),
		Failures:               []failure{},
	}
	for _, resource := range resources {
		url := strings.TrimSuffix(endpoint, "/") + resource.ID + "/providers/microsoft.insights/diagnosticSettings?api-version=" + diagnosticSettingsAPI
		_, err := client.Get(ctx, url)
		if err == nil {
			result.SuccessfulGETs++
			continue
		}
		sum := sha256.Sum256([]byte(strings.ToLower(resource.ID)))
		item := failure{
			ResourceIDHash: hex.EncodeToString(sum[:8]),
			ResourceType:   strings.ToLower(resource.Type),
		}
		var responseError *azcore.ResponseError
		if errors.As(err, &responseError) {
			item.HTTPStatus = responseError.StatusCode
			if safeErrorCode.MatchString(responseError.ErrorCode) {
				item.AzureErrorCode = responseError.ErrorCode
			}
		} else {
			item.ErrorClass = "transport_or_authentication_error"
		}
		result.Failures = append(result.Failures, item)
	}
	return result
}
