package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func decodeArray(key Key, value []map[string]interface{}) map[string]string {
	pair := make(map[string]string)
	for i, v := range value {
		newKey := key.Add(fmt.Sprintf("%v", i))
		p := decodeDict(newKey, v)
		for k, nv := range p {
			pair[k] = nv
		}
	}
	return pair
}

func decodePair(key Key, value interface{}) map[string]string {
	var pair map[string]string
	switch value.(type) {
	case []map[string]interface{}:
		pair = decodeArray(key, value.([]map[string]interface{}))
	case map[string]interface{}:
		pair = decodeDict(key, value.(map[string]interface{}))
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

func decodeDict(prefixKey Key, dict map[string]interface{}) map[string]string {
	aggregate := make(map[string]string)
	for key, value := range dict {
		dotKey := prefixKey.Add(key)
		for k, v := range decodePair(dotKey, value) {
			aggregate[k] = v
		}
	}
	return aggregate
}

type JSONMetricDecoder struct {
	KeyPrefix  Key
	MetricTime time.Time
	reader     io.Reader
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

func (d *JSONMetricDecoder) decodeDict(reader io.Reader, pairs *[]Metric) error {
	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.UseNumber()
	jsonMap := make(map[string]interface{})
	err := jsonDecoder.Decode(&jsonMap)

	if err == nil {
		//log.Printf("bk: %v", d.KeyPrefix.Key)
		agg := decodeDict(d.KeyPrefix, jsonMap)
		for k, v := range agg {
			*pairs = append(*pairs, d.metric(k, v))
		}
	}
	return err
}

func (d *JSONMetricDecoder) decodeDictArray(reader io.Reader, pairs *[]Metric) error {
	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.UseNumber()
	jsonMap := make([]map[string]interface{}, 0)
	err := jsonDecoder.Decode(&jsonMap)

	if err == nil {
		for i, m := range jsonMap {
			k := d.KeyPrefix.Add(fmt.Sprint(i))
			agg := decodeDict(k, m)
			for k, v := range agg {
				*pairs = append(*pairs, d.metric(k, v))
			}
		}
	}
	return err
}

func (d *JSONMetricDecoder) Decode(pairs *[]Metric) error {
	b := new(bytes.Buffer)
	r := io.TeeReader(d.reader, b)
	err := d.decodeDict(r, pairs)
	if err != nil {
		err = d.decodeDictArray(b, pairs)
	}

	return err
}

func NewDecoder(reader io.Reader) *JSONMetricDecoder {
	metricDecoder := &JSONMetricDecoder{reader: reader} //jsonDecoder: jsonDecoder}

	return metricDecoder
}
