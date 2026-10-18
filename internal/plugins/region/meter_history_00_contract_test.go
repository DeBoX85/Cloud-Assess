package region

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type historyOracleRequest struct {
	SubscriptionID string    `json:"SubscriptionID"`
	Months         int       `json:"Months"`
	Now            time.Time `json:"Now"`
	Currency       string    `json:"Currency"`
}
type historyOracleMeter struct {
	MeterID        string  `json:"MeterID"`
	HistoricalCost float64 `json:"HistoricalCost"`
}
type historyOracleCollection struct {
	SubscriptionID string               `json:"SubscriptionID"`
	Currency       string               `json:"Currency"`
	Status         string               `json:"Status"`
	FailureCode    string               `json:"FailureCode"`
	Meters         []historyOracleMeter `json:"Meters"`
}
type historyOracleResponse struct {
	Status      int    `json:"status"`
	Body        string `json:"body"`
	Error       string `json:"error"`
	NilResponse bool   `json:"nilResponse"`
}
type historyOracleCase struct {
	Name      string                   `json:"name"`
	Selected  map[string]string        `json:"selected"`
	Request   historyOracleRequest     `json:"request"`
	Responses []historyOracleResponse  `json:"responses"`
	Want      *historyOracleCollection `json:"want"`
	WantError string                   `json:"wantError"`
	WantURLs  []string                 `json:"wantURLs"`
}
type historyOracle struct {
	Format      string         `json:"format"`
	RequestBody map[string]any `json:"requestBody"`
	QueryCases  []struct {
		Name     string               `json:"name"`
		Request  historyOracleRequest `json:"request"`
		WantFrom string               `json:"wantFrom"`
		WantTo   string               `json:"wantTo"`
	} `json:"queryCases"`
	CollectorCases []historyOracleCase `json:"collectorCases"`
}

func historyReadOracle(t *testing.T) historyOracle {
	t.Helper()
	b, e := os.ReadFile("testdata/meter-history-contract-v1.json")
	if e != nil {
		t.Fatal(e)
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	var o historyOracle
	if e = d.Decode(&o); e != nil {
		t.Fatal(e)
	}
	var trailing any
	if d.Decode(&trailing) != io.EOF || o.Format != "region-meter-history-oracle-v1" || len(o.QueryCases) != 8 || len(o.CollectorCases) != 59 {
		t.Fatal("strict complete frozen oracle inventory")
	}
	seen := map[string]bool{}
	for _, c := range o.CollectorCases {
		if c.Name == "" || seen[c.Name] {
			t.Fatal("oracle case identity")
		}
		seen[c.Name] = true
	}
	return o
}
func historyOracleJSON(t *testing.T, v any) any {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	var out any
	if e = json.Unmarshal(b, &out); e != nil {
		t.Fatal(e)
	}
	return out
}
func historyRequest(c historyOracleRequest) MeterHistoryRequest {
	return MeterHistoryRequest{SubscriptionID: c.SubscriptionID, Months: c.Months, Now: c.Now, Currency: c.Currency}
}

// Published before BOTH target APIs exist. Every expectation is read from the
// unchanged independent literal manifest; compilation is preparation, not PASS.
func TestMeterHistoryFrozenQuery(t *testing.T) {
	o := historyReadOracle(t)
	for _, c := range o.QueryCases {
		t.Run(c.Name, func(t *testing.T) {
			b, e := MeterHistoryQuery(context.Background(), historyRequest(c.Request))
			if e != nil {
				t.Fatal(e)
			}
			var got map[string]any
			if e = json.Unmarshal(b, &got); e != nil {
				t.Fatal(e)
			}
			want := historyOracleJSON(t, o.RequestBody).(map[string]any)
			want["timePeriod"] = map[string]any{"from": c.WantFrom, "to": c.WantTo}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("complete frozen query: got=%s want=%v", b, want)
			}
		})
	}
}

type historyOraclePoster struct {
	t      *testing.T
	c      historyOracleCase
	want   map[string]any
	urls   []string
	bodies [][]byte
}

func (p *historyOraclePoster) PostBounded(_ context.Context, u string, body io.ReadSeekCloser, limit int64) ([]byte, *http.Response, error) {
	p.t.Helper()
	if limit != 1<<20 {
		p.t.Fatalf("frozen per-attempt body bound: %d", limit)
	}
	b, e := io.ReadAll(body)
	if e != nil {
		p.t.Fatal(e)
	}
	i := len(p.urls)
	p.urls = append(p.urls, u)
	p.bodies = append(p.bodies, append([]byte{}, b...))
	var query any
	if e = json.Unmarshal(b, &query); e != nil || !reflect.DeepEqual(query, p.want) {
		p.t.Fatalf("complete query before HTTP: %s %v", b, e)
	}
	if i >= len(p.c.Responses) {
		p.t.Fatal("unexpected request after terminal/unsafe page")
	}
	r := p.c.Responses[i]
	switch r.Error {
	case "ordinary":
		return nil, nil, errors.New("provider secret should not escape")
	case "cancelled":
		return nil, nil, fmt.Errorf("wrapped provider secret: %w", context.Canceled)
	case "deadline":
		return nil, nil, fmt.Errorf("wrapped provider secret: %w", context.DeadlineExceeded)
	}
	if r.NilResponse {
		return []byte(r.Body), nil, nil
	}
	return []byte(r.Body), &http.Response{StatusCode: r.Status}, nil
}

func TestMeterHistoryFrozenCollector(t *testing.T) {
	o := historyReadOracle(t)
	for _, c := range o.CollectorCases {
		t.Run(c.Name, func(t *testing.T) {
			before := historyOracleJSON(t, c)
			p := &historyOraclePoster{t: t, c: c, want: o.RequestBody, urls: []string{}}
			collector, e := NewMeterHistoryCollector("https://management.azure.com", p)
			if e != nil {
				t.Fatal(e)
			}
			got, e := collector.Collect(context.Background(), c.Selected, historyRequest(c.Request))
			switch c.WantError {
			case "cancelled":
				if !errors.Is(e, context.Canceled) {
					t.Fatalf("terminal cancel lost: %v", e)
				}
			case "deadline":
				if !errors.Is(e, context.DeadlineExceeded) {
					t.Fatalf("terminal deadline lost: %v", e)
				}
			default:
				message := ""
				if e != nil {
					message = e.Error()
				}
				if message != c.WantError {
					t.Fatalf("fixed error: got=%q want=%q", message, c.WantError)
				}
			}
			if !reflect.DeepEqual(historyOracleJSON(t, got), historyOracleJSON(t, c.Want)) {
				t.Fatalf("complete frozen collection: got=%+v want=%+v", got, c.Want)
			}
			if !reflect.DeepEqual(p.urls, c.WantURLs) || !reflect.DeepEqual(before, historyOracleJSON(t, c)) {
				t.Fatalf("request boundary/input mutation: urls=%v want=%v", p.urls, c.WantURLs)
			}
			for _, b := range p.bodies {
				if !bytes.Equal(b, p.bodies[0]) {
					t.Fatal("paging changed request bytes")
				}
			}
			if got != nil {
				evidence := got.Evidence()
				wantMeters := []historyOracleMeter{}
				if c.Want.Status == "complete" {
					wantMeters = c.Want.Meters
				}
				if evidence.Complete != (c.Want.Status == "complete") || !reflect.DeepEqual(historyOracleJSON(t, evidence.Meters), historyOracleJSON(t, wantMeters)) {
					t.Fatal("diagnostic prefix became complete cost weights")
				}
				snapshot := historyOracleJSON(t, got)
				if len(evidence.Meters) > 0 {
					evidence.Meters[0].MeterID = "mutated"
				}
				if !reflect.DeepEqual(snapshot, historyOracleJSON(t, got)) {
					t.Fatal("cost evidence shares result storage")
				}
			}
		})
	}
}
