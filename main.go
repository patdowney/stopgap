package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/marpaia/graphite-golang"
)

type GraphiteConfig struct {
	Host string
	Port int
}

type Config struct {
	DryRun      bool
	Prefix      string
	DefaultTime time.Time
	Graphite    GraphiteConfig
}

type MetricKey struct {
	Key string
}

func (k *MetricKey) Add(newKeyPart string) MetricKey {
	newKey := newKeyPart
	if k.Key != "" {
		newKey = fmt.Sprintf("%v.%v", k.Key, newKeyPart)
	}

	return MetricKey{Key: newKey}
}

func (k *MetricKey) String() string {
	return k.Key
}

func decodePair(key MetricKey, value interface{}) map[string]string {
	pair := make(map[string]string)

	switch value.(type) {
	case bool, string, []interface{}:
		break
	case json.Number:
		pair[key.String()] = value.(json.Number).String()
	case float64:
		pair[key.String()] = fmt.Sprintf("%f", value)
	default:
		pair = decodeDict(key, value.(map[string]interface{}))
	}
	return pair
}

func decodeDict(prefixKey MetricKey, dict map[string]interface{}) map[string]string {
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
	KeyPrefix   MetricKey
	DefaultTime time.Time
	jsonDecoder *json.Decoder
}

type Metric struct {
	Key   string
	Value string
	Time  time.Time
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

func DumpMetrics(writer io.Writer, metrics []Metric) {
	for _, m := range metrics {
		fmt.Fprintf(writer, "%v: %v\n", m.Key, m.Value)
	}
}

func PublishMetrics(gclient *graphite.Graphite, metrics []Metric) {
	for _, m := range metrics {
		gmetric := graphite.Metric{
			Name:      m.Key,
			Value:     m.Value,
			Timestamp: m.Time.Unix()}
		gclient.SendMetric(gmetric)
	}
}

func graphiteConfig(cfg *GraphiteConfig) {
	flag.StringVar(&cfg.Host, "graphite-host", "localhost", "graphite host")
	flag.IntVar(&cfg.Port, "graphite-port", 2003, "graphite port")
}

type TimeArg struct {
	time.Time
}

func (t *TimeArg) Set(value string) error {
	var err error
	t.Time, err = time.Parse(time.RFC3339, value)
	return err
}

func config() *Config {
	c := &Config{}

	flag.BoolVar(&c.DryRun, "dry-run", false, "dry run")
	flag.StringVar(&c.Prefix, "prefix", "", "metric prefix")

	var t TimeArg
	flag.Var(&t, "default-time", "default metric time")

	graphiteConfig(&c.Graphite)

	flag.Parse()

	c.DefaultTime = t.Time

	return c
}

func main() {
	cfg := config()
	metrics := make([]Metric, 0, 1)
	metricDecoder := NewDecoder(os.Stdin)
	metricDecoder.DefaultTime = cfg.DefaultTime
	metricDecoder.KeyPrefix = (&MetricKey{}).Add(cfg.Prefix)
	_ = metricDecoder.Decode(&metrics)

	gclient := graphite.NewGraphiteNop(cfg.Graphite.Host, cfg.Graphite.Port)
	if !cfg.DryRun {
		var err error
		gclient, err = graphite.NewGraphite(cfg.Graphite.Host, cfg.Graphite.Port)
		if err != nil {
			log.Fatal(err)
		}
	}

	PublishMetrics(gclient, metrics)
}
