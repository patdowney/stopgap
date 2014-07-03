package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func FlattenList(key Key, value []map[string]interface{}, keyID string) map[string]string {
	pair := make(map[string]string)
	for i, v := range value {
		newKey := key.Add(fmt.Sprintf("%v", i))
		if keyID != "" {
			newKey = key.Add(fmt.Sprintf("%v", v[keyID]))
		}
		p := FlattenMap(newKey, v)
		for k, nv := range p {
			pair[k] = nv
		}
	}
	return pair
}

func FlattenPair(key Key, value interface{}) map[string]string {
	var pair map[string]string
	switch value.(type) {
	case []map[string]interface{}:
		pair = FlattenList(key, value.([]map[string]interface{}), "")
	case map[string]interface{}:
		pair = FlattenMap(key, value.(map[string]interface{}))
	case json.Number:
		pair = make(map[string]string)
		pair[key.String()] = value.(json.Number).String()
	case float64:
		pair = make(map[string]string)
		pair[key.String()] = fmt.Sprintf("%f", value)
	case bool, string, nil, []interface{}:
		break
	}
	return pair
}

func FlattenMap(prefixKey Key, m map[string]interface{}) map[string]string {
	aggregate := make(map[string]string)
	for key, value := range m {
		dotKey := prefixKey.Add(key)
		for k, v := range FlattenPair(dotKey, value) {
			aggregate[k] = v
		}
	}
	return aggregate
}

type JSONMetricDecoder struct {
	KeyPrefix   Key
	ListItemKey string
	MetricTime  time.Time
	reader      io.Reader
}

func (d *JSONMetricDecoder) time() time.Time {
	var nullTime time.Time
	if d.MetricTime == nullTime {
		return time.Now()
	}

	return d.MetricTime
}

func (d *JSONMetricDecoder) metric(k, v string) Metric {
	return Metric{
		Key:   k,
		Value: v,
		Time:  d.time()}
}

func (d *JSONMetricDecoder) FlattenMap(reader io.Reader, pairs *[]Metric) error {
	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.UseNumber()
	jsonMap := make(map[string]interface{})
	err := jsonDecoder.Decode(&jsonMap)

	if err == nil {
		agg := FlattenMap(d.KeyPrefix, jsonMap)
		for k, v := range agg {
			*pairs = append(*pairs, d.metric(k, v))
		}
	}
	return err
}

func (d *JSONMetricDecoder) FlattenMapList(reader io.Reader, pairs *[]Metric) error {
	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.UseNumber()
	jsonList := make([]map[string]interface{}, 0)
	err := jsonDecoder.Decode(&jsonList)

	if err == nil {
		agg := FlattenList(d.KeyPrefix, jsonList, d.ListItemKey)
		for k, v := range agg {
			*pairs = append(*pairs, d.metric(k, v))
		}
	}
	return err
}

func (d *JSONMetricDecoder) Decode(pairs *[]Metric) error {
	b := new(bytes.Buffer)
	r := io.TeeReader(d.reader, b)
	err := d.FlattenMap(r, pairs)
	if err != nil {
		err = d.FlattenMapList(b, pairs)
	}

	return err
}

func NewDecoder(reader io.Reader) *JSONMetricDecoder {
	metricDecoder := &JSONMetricDecoder{reader: reader} //jsonDecoder: jsonDecoder}

	return metricDecoder
}
