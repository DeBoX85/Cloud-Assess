// Wire fields follow the pinned AZQR SDK versions. See AI_GOVERNANCE_REQUESTS.md.
package aigov

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxPageBytes       = 2 << 20
	MaxJSONDepth       = 16
	MaxJSONKeys        = 64
	MaxBatchAccounts   = 50
	MaxPageDeployments = 1000
	MaxNextLinkBytes   = 4096
)

// MetricPage must be combined with projection health. FailedIDs/Unreported are
// incomplete request coverage, not successful empty series. Errors retain only
// already decoded valid points; ambiguous whole JSON fails before projection.
type MetricPage struct {
	Points                []Point
	Malformed             int
	FailedIDs, Unreported []string
}
type DeploymentPage struct {
	Values    []Deployment
	NextLink  string
	Malformed int
}

func validateJSON(ctx context.Context, raw []byte) error {
	if len(raw) > MaxPageBytes || !utf8.Valid(raw) {
		return fmt.Errorf("invalid AI response size or encoding")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	tokens := 0
	var walk func(int) error
	walk = func(depth int) error {
		if e := ctx.Err(); e != nil {
			return e
		}
		tokens++
		if depth > MaxJSONDepth || tokens > 1<<20 {
			return fmt.Errorf("AI JSON structure limit")
		}
		token, e := d.Token()
		if e != nil {
			return fmt.Errorf("invalid AI JSON")
		}
		delimiter, compound := token.(json.Delim)
		if !compound {
			return nil
		}
		switch delimiter {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				if e := ctx.Err(); e != nil {
					return e
				}
				key, e := d.Token()
				s, ok := key.(string)
				if e != nil || !ok || !safe(s, 128) || len(seen) >= MaxJSONKeys || seen[strings.ToLower(s)] {
					return fmt.Errorf("ambiguous AI JSON object")
				}
				seen[strings.ToLower(s)] = true
				if e = walk(depth + 1); e != nil {
					return e
				}
			}
			end, e := d.Token()
			if e != nil || end != json.Delim('}') {
				return fmt.Errorf("invalid AI JSON object")
			}
		case '[':
			n := 0
			for d.More() {
				if n >= MaxPoints {
					return fmt.Errorf("AI JSON array limit")
				}
				n++
				if e = walk(depth + 1); e != nil {
					return e
				}
			}
			end, e := d.Token()
			if e != nil || end != json.Delim(']') {
				return fmt.Errorf("invalid AI JSON array")
			}
		default:
			return fmt.Errorf("invalid AI JSON delimiter")
		}
		return nil
	}
	if e := walk(0); e != nil {
		return e
	}
	if _, e := d.Token(); e != io.EOF {
		return fmt.Errorf("trailing AI JSON")
	}
	return nil
}
func wireObject(raw json.RawMessage, known ...string) (map[string]json.RawMessage, error) {
	if len(raw) == 0 || raw[0] != '{' {
		return nil, fmt.Errorf("invalid AI object")
	}
	var values map[string]json.RawMessage
	if json.Unmarshal(raw, &values) != nil || values == nil {
		return nil, fmt.Errorf("invalid AI object")
	}
	for k := range values {
		if strings.EqualFold(k, "error") {
			return nil, fmt.Errorf("AI provider error envelope")
		}
		for _, field := range known {
			if strings.EqualFold(k, field) && k != field {
				return nil, fmt.Errorf("aliased AI field")
			}
		}
	}
	return values, nil
}
func optionalObject(raw json.RawMessage, known ...string) (map[string]json.RawMessage, error) {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil, nil
	}
	return wireObject(raw, known...)
}
func wireArray(raw json.RawMessage, limit int, optional bool) ([]json.RawMessage, error) {
	if optional && (len(raw) == 0 || bytes.Equal(raw, []byte("null"))) {
		return nil, nil
	}
	if len(raw) == 0 || raw[0] != '[' {
		return nil, fmt.Errorf("missing AI array")
	}
	var values []json.RawMessage
	if json.Unmarshal(raw, &values) != nil || len(values) > limit {
		return nil, fmt.Errorf("invalid or excessive AI array")
	}
	return values, nil
}
func wireString(raw json.RawMessage, limit int) (*string, error) {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil, nil
	}
	var s string
	if json.Unmarshal(raw, &s) != nil || !safe(s, limit) {
		return nil, fmt.Errorf("invalid AI string")
	}
	return &s, nil
}
func DecodeMetrics(ctx context.Context, raw []byte, requested []string) (MetricPage, error) {
	out := MetricPage{Points: []Point{}, FailedIDs: []string{}, Unreported: []string{}}
	scope := map[string]string{}
	reported := map[string]bool{}
	failed := map[string]bool{}
	finish := func() {
		out.FailedIDs = nil
		out.Unreported = nil
		for _, id := range requested {
			key := strings.ToLower(id)
			if failed[key] {
				out.FailedIDs = append(out.FailedIDs, id)
			}
			if !reported[key] {
				out.Unreported = append(out.Unreported, id)
			}
		}
	}
	fail := func(e error) (MetricPage, error) { finish(); return out, e }
	if len(requested) == 0 || len(requested) > MaxBatchAccounts {
		return fail(fmt.Errorf("invalid AI metrics batch scope"))
	}
	for _, id := range requested {
		key := strings.ToLower(id)
		if !safe(id, 2048) || !accountID.MatchString(id) || scope[key] != "" {
			return fail(fmt.Errorf("invalid AI metrics batch identity"))
		}
		scope[key] = id
	}
	if e := validateJSON(ctx, raw); e != nil {
		return fail(e)
	}
	fields, e := wireObject(bytes.TrimSpace(raw), "values")
	if e != nil {
		return fail(e)
	}
	resources, e := wireArray(fields["values"], MaxBatchAccounts, false)
	if e != nil {
		return fail(e)
	}
	for _, r := range resources {
		if e := ctx.Err(); e != nil {
			return fail(e)
		}
		f, e := wireObject(r, "resourceid", "value", "namespace", "resourceregion", "interval", "starttime", "endtime")
		if e != nil {
			out.Malformed++
			continue
		}
		id, e := wireString(f["resourceid"], 2048)
		if e != nil || id == nil || *id == "" {
			out.Malformed++
			continue
		}
		key := strings.ToLower(*id)
		if scope[key] == "" || reported[key] {
			return fail(fmt.Errorf("unrequested or duplicate AI metrics resource"))
		}
		reported[key] = true
		namespace, e := wireString(f["namespace"], MaxLabelBytes)
		if e != nil || namespace != nil && !strings.EqualFold(*namespace, "Microsoft.CognitiveServices/accounts") {
			failed[key] = true
			continue
		}
		metrics, e := wireArray(f["value"], 16, false)
		if e != nil {
			failed[key] = true
			continue
		}
		for _, m := range metrics {
			if e := ctx.Err(); e != nil {
				return fail(e)
			}
			metric, e := wireObject(m, "timeseries", "errorCode", "errorMessage", "name", "unit", "id", "type", "displayDescription")
			if e != nil {
				out.Malformed++
				failed[key] = true
				continue
			}
			code, e := wireString(metric["errorCode"], MaxLabelBytes)
			if e != nil || code != nil && *code != "" && *code != "Success" {
				failed[key] = true
				continue
			}
			metricName, e := optionalObject(metric["name"], "value", "localizedValue")
			if e != nil {
				failed[key] = true
				continue
			}
			name, e := wireString(metricName["value"], MaxLabelBytes)
			if e != nil || name != nil && *name != "AzureOpenAIRequests" {
				failed[key] = true
				continue
			}
			unit, e := wireString(metric["unit"], MaxLabelBytes)
			if e != nil || unit != nil && *unit != "Count" {
				failed[key] = true
				continue
			}
			series, e := wireArray(metric["timeseries"], MaxPoints, true)
			if e != nil {
				failed[key] = true
				continue
			}
			for _, s := range series {
				if e := ctx.Err(); e != nil {
					return fail(e)
				}
				value, e := wireObject(s, "metadatavalues", "data")
				if e != nil {
					out.Malformed++
					continue
				}
				metadata, e := wireArray(value["metadatavalues"], 64, true)
				if e != nil {
					out.Malformed++
					continue
				}
				var deployment, model, status *string
				bad := false
				for _, v := range metadata {
					meta, e := wireObject(v, "name", "value")
					if e != nil {
						bad = true
						break
					}
					n, e := optionalObject(meta["name"], "value", "localizedValue")
					if e != nil {
						bad = true
						break
					}
					dimensionName, e := wireString(n["value"], MaxLabelBytes)
					if e != nil {
						bad = true
						break
					}
					dimensionValue, e := wireString(meta["value"], MaxLabelBytes)
					if e != nil {
						bad = true
						break
					}
					if dimensionName == nil || dimensionValue == nil {
						continue
					}
					switch strings.ToLower(*dimensionName) {
					case "modeldeploymentname":
						deployment = dimensionValue
					case "modelname":
						model = dimensionValue
					case "statuscode":
						status = dimensionValue
					}
				}
				if bad {
					out.Malformed++
					continue
				}
				data, e := wireArray(value["data"], MaxPoints, true)
				if e != nil {
					out.Malformed++
					continue
				}
				for _, v := range data {
					if e := ctx.Err(); e != nil {
						return fail(e)
					}
					if len(out.Points) >= MaxPoints {
						return fail(fmt.Errorf("AI metrics point limit"))
					}
					point := Point{ResourceID: scope[key], Deployment: copyString(deployment), Model: copyString(model), Status: copyString(status)}
					f, e := wireObject(v, "timeStamp", "count", "average", "maximum", "minimum", "total")
					if e != nil {
						out.Malformed++
						continue
					}
					timestamp, e := wireString(f["timeStamp"], 64)
					if e != nil {
						out.Malformed++
						continue
					}
					if timestamp != nil {
						t, e := time.Parse(time.RFC3339Nano, *timestamp)
						if e != nil {
							out.Malformed++
							continue
						}
						point.Timestamp = &t
					}
					count := f["count"]
					if len(count) > 0 && !bytes.Equal(count, []byte("null")) {
						var n float64
						if json.Unmarshal(count, &n) != nil || math.IsNaN(n) || math.IsInf(n, 0) || n < 0 {
							out.Malformed++
							continue
						}
						point.Count = &n
					}
					out.Points = append(out.Points, point)
				}
			}
		}
	}
	finish()
	return out, nil
}
func DecodeDeployments(ctx context.Context, raw []byte) (DeploymentPage, error) {
	out := DeploymentPage{Values: []Deployment{}}
	if e := validateJSON(ctx, raw); e != nil {
		return out, e
	}
	fields, e := wireObject(bytes.TrimSpace(raw), "value", "nextLink")
	if e != nil {
		return out, e
	}
	rows, e := wireArray(fields["value"], MaxPageDeployments, false)
	if e != nil {
		return out, e
	}
	for _, r := range rows {
		if e := ctx.Err(); e != nil {
			return out, e
		}
		f, e := wireObject(r, "name", "properties", "sku", "id", "type", "systemData", "etag")
		if e != nil {
			out.Malformed++
			continue
		}
		d := Deployment{}
		d.Name, e = wireString(f["name"], MaxLabelBytes)
		if e != nil || d.Name == nil {
			out.Malformed++
			continue
		}
		properties, e := optionalObject(f["properties"], "model", "versionUpgradeOption", "spilloverDeploymentName")
		if e != nil {
			out.Malformed++
			continue
		}
		model, e := optionalObject(properties["model"], "version", "format")
		if e != nil {
			out.Malformed++
			continue
		}
		bad := false
		for k, p := range map[string]**string{"version": &d.ModelVersion, "format": &d.ModelFormat} {
			*p, e = wireString(model[k], MaxLabelBytes)
			if e != nil {
				bad = true
			}
		}
		d.Upgrade, e = wireString(properties["versionUpgradeOption"], MaxLabelBytes)
		if e != nil {
			bad = true
		}
		d.Spillover, e = wireString(properties["spilloverDeploymentName"], MaxLabelBytes)
		if e != nil {
			bad = true
		}
		sku, e := optionalObject(f["sku"], "capacity", "name", "tier")
		if e != nil {
			bad = true
		}
		capacity := sku["capacity"]
		if len(capacity) > 0 && !bytes.Equal(capacity, []byte("null")) {
			var n int64
			if json.Unmarshal(capacity, &n) != nil || n < 0 {
				bad = true
			} else {
				d.Capacity = &n
			}
		}
		if bad {
			out.Malformed++
			continue
		}
		out.Values = append(out.Values, d)
	}
	link, e := wireString(fields["nextLink"], MaxNextLinkBytes)
	if e != nil {
		return out, e
	}
	if link != nil {
		out.NextLink = *link
	}
	return out, nil
}

func copyString(p *string) *string {
	if p == nil {
		return nil
	}
	s := *p
	return &s
}
