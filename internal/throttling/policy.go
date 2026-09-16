// Portions of this file reproduce throttling behavior from Microsoft Azure Quick Review
// (MIT licensed). See NOTICE.md.

package throttling

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"golang.org/x/time/rate"
)

type Kind string

const (
	KindARM    Kind = "arm"
	KindGraph  Kind = "graph"
	KindCost   Kind = "cost"
	KindRetail Kind = "retail-prices"
)

var (
	armLimiter   = rate.NewLimiter(rate.Limit(20), 100)
	graphLimiter = rate.NewLimiter(rate.Limit(3), 10)
	costLimiter  = rate.NewLimiter(rate.Limit(0.2), 1)
)

type Policy struct{}

func NewPolicy() policy.Policy { return &Policy{} }

func KindForURL(url string) Kind {
	switch {
	case strings.Contains(url, "Microsoft.ResourceGraph/resources"):
		return KindGraph
	case strings.Contains(url, "Microsoft.CostManagement/query"):
		return KindCost
	case strings.Contains(url, "prices.azure.com"):
		return KindRetail
	default:
		return KindARM
	}
}

func (p *Policy) Do(request *policy.Request) (*http.Response, error) {
	var err error
	switch KindForURL(request.Raw().URL.String()) {
	case KindGraph:
		err = graphLimiter.Wait(request.Raw().Context())
	case KindCost:
		err = costLimiter.Wait(request.Raw().Context())
	case KindRetail:
		return request.Next()
	default:
		err = armLimiter.Wait(request.Raw().Context())
	}
	if err != nil {
		return nil, fmt.Errorf("throttling wait failed: %w", err)
	}
	return request.Next()
}
