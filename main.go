package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marpaia/graphite-golang"
	"github.com/patdowney/stopgap/metrics"
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

func graphiteConfig(cfg *GraphiteConfig) {
	flag.StringVar(&cfg.Host, "graphite-host", "localhost", "graphite host")
	flag.IntVar(&cfg.Port, "graphite-port", 2003, "graphite port")
}

type timeArg struct {
	time.Time
}

func (t *timeArg) Set(value string) error {
	var err error
	t.Time, err = time.Parse(time.RFC3339, value)
	return err
}
func (t *timeArg) String() string {
	return fmt.Sprint(t.Time.Unix()) //Format(time.RFC3339)
}

func config() *Config {
	c := &Config{}

	flag.BoolVar(&c.DryRun, "dry-run", false, "dry run")
	flag.StringVar(&c.Prefix, "prefix", "", "metric prefix")

	var t timeArg
	t.Time = time.Now()
	flag.Var(&t, "metric-timestamp", "metric time stamp (defaults to now)")

	graphiteConfig(&c.Graphite)

	flag.Parse()

	c.DefaultTime = t.Time

	return c
}

func main() {
	cfg := config()
	m := make([]metrics.Metric, 0, 1)
	metricDecoder := metrics.NewDecoder(os.Stdin)
	metricDecoder.DefaultTime = cfg.DefaultTime
	metricDecoder.KeyPrefix = (metrics.Key{}).Add(cfg.Prefix)
	_ = metricDecoder.Decode(&m)

	gclient := graphite.NewGraphiteNop(cfg.Graphite.Host, cfg.Graphite.Port)
	if !cfg.DryRun {
		var err error
		gclient, err = graphite.NewGraphite(cfg.Graphite.Host, cfg.Graphite.Port)
		if err != nil {
			log.Fatal(err)
		}
	}

	metrics.PublishMetrics(gclient, m)
}
