package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func decodePair(key Key, value interface{}) map[string]string {
	var pair map[string]string

	switch value.(type) {
	case bool, string, []interface{}:
		break
	case json.Number:
		pair = make(map[string]string)
		pair[key.String()] = value.(json.Number).String()
	case float64:
		pair = make(map[string]string)
		pair[key.String()] = fmt.Sprintf("%f", value)
	default:
		pair = decodeDict(key, value.(map[string]interface{}))
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
	KeyPrefix   Key
	DefaultTime time.Time
	jsonDecoder *json.Decoder
}

func (d *JSONMetricDecoder) time() time.Time {
	var nullTime time.Time
	if d.DefaultTime == nullTime {
		return time.Now()
	}

	return d.DefaultTime
}

func (d *JSONMetricDecoder) metric(k, v string) Metric {
	return Metric{
		Key:   k,
		Value: v,
		Time:  d.time()}
}

func (d *JSONMetricDecoder) Decode(pairs *[]Metric) error {
	jsonMap := make(map[string]interface{})
	err := d.jsonDecoder.Decode(&jsonMap)

	if err == nil {
		//log.Printf("bk: %v", d.KeyPrefix.Key)
		agg := decodeDict(d.KeyPrefix, jsonMap)
		for k, v := range agg {
			*pairs = append(*pairs, d.metric(k, v))
		}
	}
	return err
}

func NewDecoder(reader io.Reader) *JSONMetricDecoder {
	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.UseNumber()

	metricDecoder := &JSONMetricDecoder{jsonDecoder: jsonDecoder}

	return metricDecoder
}

func NewRemoteDecoder(url url.URL) (*JSONMetricDecoder, error) {
	res, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	return NewDecoder(res.Body), nil
}
