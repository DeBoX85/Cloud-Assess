package aigov

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func fixtureBytes(t *testing.T, name string) []byte {
	t.Helper()
	b, e := os.ReadFile("testdata/" + name + ".json")
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func metricWire(id, series string) string {
	b, _ := json.Marshal(id)
	return `{"values":[{"resourceid":` + string(b) + `,"value":[{"name":{"value":"AzureOpenAIRequests"},"unit":"Count","errorCode":"Success","timeseries":` + series + `}]}]}`
}

const oneSeries = `[{"metadatavalues":[{"name":{"value":"ModelDeploymentName"},"value":"fixture-deployment"}],"data":[{"timeStamp":"2026-09-30T12:05:00Z","count":1.5}]}]`

func TestDecodeCapturedWire(t *testing.T) {
	page, e := DecodeMetrics(context.Background(), fixtureBytes(t, "metrics-input"), []string{fixtureID, fixtureAccounts()[1].ID})
	if e != nil || len(page.Points) != 8 || page.Malformed != 1 || !reflect.DeepEqual(page.Unreported, []string{fixtureAccounts()[1].ID}) || len(page.FailedIDs) != 0 {
		t.Fatalf("capture coverage: %+v %v", page, e)
	}
	if page.Points[2].Timestamp.Format("2006-01-02 15:04 -07:00") != "2026-09-30 13:00 +02:00" || *page.Points[7].Status != "503" || page.Points[3].Count != nil || page.Points[4].Timestamp != nil {
		t.Fatal("source point semantics")
	}
	*page.Points[0].Deployment = "changed"
	if *page.Points[1].Deployment != "fixture-deployment" {
		t.Fatal("decoded points share mutable dimensions")
	}
	d, e := DecodeDeployments(context.Background(), fixtureBytes(t, "deployments-input"))
	if e != nil || len(d.Values) != 1 || d.Values[0].Spillover == nil || *d.Values[0].Spillover != "" || *d.Values[0].Capacity != 42 {
		t.Fatalf("capture deployments: %+v %v", d, e)
	}
}

func TestMetricCoverageCorrelationAndErrors(t *testing.T) {
	for _, tc := range []struct {
		name, body                        string
		fatal                             bool
		failed, absent, malformed, points int
	}{
		{"empty reported", metricWire(fixtureID, "[]"), false, 0, 0, 0, 0},
		{"empty missing", `{"values":[]}`, false, 0, 1, 0, 0},
		{"wrong metric", strings.Replace(metricWire(fixtureID, oneSeries), "AzureOpenAIRequests", "Other", 1), false, 1, 0, 0, 0},
		{"query error", strings.Replace(metricWire(fixtureID, oneSeries), `"Success"`, `"BadRequest"`, 1), false, 1, 0, 0, 0},
		{"wrong unit", strings.Replace(metricWire(fixtureID, oneSeries), `"Count"`, `"Bytes"`, 1), false, 1, 0, 0, 0},
		{"wrong namespace", strings.Replace(metricWire(fixtureID, oneSeries), `"value":[`, `"namespace":"Other","value":[`, 1), false, 1, 0, 0, 0},
		{"foreign", metricWire(strings.Replace(fixtureID, fixtureSub, "22222222-2222-4222-8222-222222222222", 1), oneSeries), true, 0, 1, 0, 0},
		{"child", metricWire(fixtureID+"/deployments/x", oneSeries), true, 0, 1, 0, 0},
		{"missing array", `{"values":[{"resourceid":"` + fixtureID + `"}]}`, false, 1, 0, 0, 0},
		{"alias", `{"values":[{"resourceId":"` + fixtureID + `","value":[]}]}`, false, 0, 1, 1, 0},
		{"negative", metricWire(fixtureID, strings.Replace(oneSeries, "1.5", "-1", 1)), false, 0, 0, 1, 0},
		{"invalid time", metricWire(fixtureID, strings.Replace(oneSeries, "2026-09-30T12:05:00Z", "invalid", 1)), false, 0, 0, 1, 0},
		{"overflow", metricWire(fixtureID, strings.Replace(oneSeries, "1.5", "1e999", 1)), false, 0, 0, 1, 0},
		{"null series", metricWire(fixtureID, "null"), false, 0, 0, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p, e := DecodeMetrics(context.Background(), []byte(tc.body), []string{fixtureID})
			if (e != nil) != tc.fatal || len(p.FailedIDs) != tc.failed || len(p.Unreported) != tc.absent || p.Malformed != tc.malformed || len(p.Points) != tc.points {
				t.Fatalf("unexpected health %+v %v", p, e)
			}
		})
	}
	good := metricWire(fixtureID, oneSeries)
	var root struct {
		Values []json.RawMessage `json:"values"`
	}
	if e := json.Unmarshal([]byte(good), &root); e != nil {
		t.Fatal(e)
	}
	duplicate := `{"values":[` + string(root.Values[0]) + `,` + string(root.Values[0]) + `]}`
	p, e := DecodeMetrics(context.Background(), []byte(duplicate), []string{fixtureID})
	if e == nil || len(p.Points) != 1 {
		t.Fatal("duplicate identity concealed or prior point lost")
	}
	for _, ids := range [][]string{nil, {fixtureID, strings.ToUpper(fixtureID)}, {"invalid"}, make([]string, 51)} {
		if _, e := DecodeMetrics(context.Background(), []byte(good), ids); e == nil {
			t.Fatal("invalid request scope accepted")
		}
	}
}

func TestStrictWireAndDeploymentDefaults(t *testing.T) {
	for _, body := range []string{`null`, `[]`, `{}`, `{"value":null}`, `{"value":[],"VALUE":[]}`, `{"value":[],"error":{"code":"secret"}}`, `{"value":[]} {}`, `{"Value":[]}`, strings.Repeat(" ", MaxPageBytes+1), "{\"value\":[],\"x\":\"\xff\"}"} {
		if _, e := DecodeDeployments(context.Background(), []byte(body)); e == nil {
			t.Fatalf("invalid wire accepted: %.80s", body)
		}
	}
	body := `{"value":[null,{}, {"name":"nil"}, {"name":"zero","sku":{"capacity":0},"properties":{"model":{"version":"","format":""},"versionUpgradeOption":"","spilloverDeploymentName":""}}, {"name":"negative","sku":{"capacity":-1}}, {"name":"bad","properties":{"model":{"Version":"secret"}}}],"nextLink":"https://arm.test/path"}`
	p, e := DecodeDeployments(context.Background(), []byte(body))
	if e != nil || len(p.Values) != 2 || p.Malformed != 4 || p.Values[0].Capacity != nil || p.Values[1].Capacity == nil || *p.Values[1].Capacity != 0 || p.Values[1].Spillover == nil || p.NextLink != "https://arm.test/path" {
		t.Fatalf("defaults/malformed: %+v %v", p, e)
	}
	p, e = DecodeDeployments(context.Background(), []byte(`{"value":[{"name":"kept"}],"nextLink":42}`))
	if e == nil || len(p.Values) != 1 {
		t.Fatal("invalid continuation discarded earlier metadata")
	}
	depth := `{"value":[],"unknown":` + strings.Repeat("[", MaxJSONDepth+1) + "0" + strings.Repeat("]", MaxJSONDepth+1) + "}"
	if _, e := DecodeDeployments(context.Background(), []byte(depth)); e == nil {
		t.Fatal("depth limit missing")
	}
	for _, decode := range []func(context.Context, []byte) error{
		func(c context.Context, b []byte) error { _, e := DecodeDeployments(c, b); return e },
		func(c context.Context, b []byte) error { _, e := DecodeMetrics(c, b, []string{fixtureID}); return e },
	} {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if !errors.Is(decode(ctx, []byte(`{}`)), context.Canceled) {
			t.Fatal("lost cancellation identity")
		}
	}
}

func FuzzBoundedAIWire(f *testing.F) {
	for _, s := range []string{metricWire(fixtureID, oneSeries), `{"value":[{"name":"x"}],"nextLink":null}`, `{"value":[],"Value":[]}`, `null`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 65536 {
			return
		}
		m, _ := DecodeMetrics(context.Background(), raw, []string{fixtureID})
		d, _ := DecodeDeployments(context.Background(), raw)
		if len(m.Points) > MaxPoints || len(d.Values) > MaxPageDeployments {
			t.Fatal("unbounded decoder output")
		}
		for _, p := range m.Points {
			if p.ResourceID != fixtureID {
				t.Fatal("decoder accepted foreign identity")
			}
		}
	})
}
