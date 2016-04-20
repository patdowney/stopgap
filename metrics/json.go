package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// JSONMetricDecoder ...
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

func filterMapNumbers(m map[string]interface{}) map[string]string {
	filteredMap := make(map[string]string)

	for k, v := range m {
		switch v.(type) {
		case int64, float64:
			filteredMap[k] = fmt.Sprintf("%v", v)
		}
	}
	return filteredMap
}

// FlattenMap ...
func (d *JSONMetricDecoder) flattenMap(reader io.Reader, pairs *[]Metric) error {
	jsonMap := make(map[string]interface{})
	return d.flattenAndFilter(reader, pairs, jsonMap)
}

func (d *JSONMetricDecoder) flattenAndFilter(reader io.Reader, pairs *[]Metric, data interface{}) error {
	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.UseNumber()

	err := jsonDecoder.Decode(&data)
	if err == nil {
		agg := filterMapNumbers(Flatten(d.KeyPrefix, data, d.ListItemKey))
		for k, v := range agg {
			*pairs = append(*pairs, d.metric(k, v))
		}
	}
	return err
}

// FlattenMapList ...
func (d *JSONMetricDecoder) flattenMapList(reader io.Reader, pairs *[]Metric) error {
	var jsonList []map[string]interface{}
	return d.flattenAndFilter(reader, pairs, jsonList)
}

// Decode ...
func (d *JSONMetricDecoder) Decode(pairs *[]Metric) error {
	b := new(bytes.Buffer)
	r := io.TeeReader(d.reader, b)
	err := d.flattenMap(r, pairs)
	if err != nil {
		err = d.flattenMapList(b, pairs)
		fmt.Printf("FlattenMapList: %s\n", err)
	}

	return err
}

// NewDecoder ...
func NewDecoder(reader io.Reader) *JSONMetricDecoder {
	metricDecoder := &JSONMetricDecoder{reader: reader}

	return metricDecoder
}
