package metrics

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

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
