package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/patdowney/stopgap/metrics"
)

type Config struct {
	DryRun     bool
	Prefix     string
	MetricTime time.Time
	Graphite   metrics.GraphiteConfig
}

func graphiteConfig(cfg *metrics.GraphiteConfig) {
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

	c.MetricTime = t.Time

	return c
}

func openHTTPReader(readerUrl string) (io.ReadCloser, error) {
	res, err := http.Get(readerUrl)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return res.Body, nil
	}

	return nil, errors.New(fmt.Sprintf("non-200 response from %v(%v)", readerUrl, res.StatusCode))
}

func openURLReader(readerUrl string) (io.ReadCloser, error) {
	u, err := url.Parse(readerUrl)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "":
		return os.Open(u.String())
	case "http", "https":
		return openHTTPReader(u.String())
	default:
		return nil, errors.New(fmt.Sprintf("unsupported schema: %v", u.Scheme))
	}
}

func openReader(readerUri string) (io.ReadCloser, error) {
	if readerUri == "-" {
		return os.Stdin, nil
	}
	return openURLReader(readerUri)
}

func openReaders(readerUris []string) ([]io.ReadCloser, error) {
	readers := make([]io.ReadCloser, 0, 1)

	if len(readerUris) == 0 {
		readers = append(readers, os.Stdin)
	} else {
		for _, a := range readerUris {
			r, err := openReader(a)
			if err != nil {
				return nil, err
			}
			readers = append(readers, r)
		}
	}
	return readers, nil
}

func closeReaders(readers []io.ReadCloser) {
	for _, r := range readers {
		r.Close()
	}
}

func main() {
	cfg := config()
	m := make([]metrics.Metric, 0, 1)

	readers, err := openReaders(flag.Args())

	for _, reader := range readers {
		metricDecoder := metrics.NewDecoder(reader)
		metricDecoder.MetricTime = cfg.MetricTime
		metricDecoder.KeyPrefix = (metrics.Key{}).Add(cfg.Prefix)
		err = metricDecoder.Decode(&m)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	closeReaders(readers)

	publisher, err := metrics.NewGraphitePublisher(
		cfg.Graphite, cfg.DryRun)
	if err != nil {
		log.Fatal(err)
	}

	publisher.Publish(m)
}
