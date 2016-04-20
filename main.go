package main

import (
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

type config struct {
	DryRun      bool
	Prefix      string
	ListItemKey string
	MetricTime  time.Time
	Graphite    metrics.GraphiteConfig
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

func fetchConfig() *config {
	c := &config{}

	flag.BoolVar(&c.DryRun, "dry-run", false, "dry run")
	flag.StringVar(&c.Prefix, "prefix", "", "metric prefix")
	flag.StringVar(&c.ListItemKey, "list-item-key", "", "use item key instead of item index")

	var t timeArg
	t.Time = time.Now()
	flag.Var(&t, "metric-timestamp", "metric time stamp (defaults to now)")

	graphiteConfig(&c.Graphite)

	flag.Parse()

	c.MetricTime = t.Time

	return c
}

func openHTTPReader(readerURL string) (io.ReadCloser, error) {
	res, err := http.Get(readerURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return res.Body, nil
	}

	return nil, fmt.Errorf("non-200 response from %v(%v)", readerURL, res.StatusCode)
}

func openURLReader(readerURL string) (io.ReadCloser, error) {
	u, err := url.Parse(readerURL)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "":
		return os.Open(u.String())
	case "http", "https":
		return openHTTPReader(u.String())
	default:
		return nil, fmt.Errorf("unsupported schema: %v", u.Scheme)
	}
}

func openReader(readerURI string) (io.ReadCloser, error) {
	if readerURI == "-" {
		return os.Stdin, nil
	}
	return openURLReader(readerURI)
}

func openReaders(readerURIs []string) ([]io.ReadCloser, error) {
	readers := make([]io.ReadCloser, 0, 1)

	if len(readerURIs) == 0 {
		readers = append(readers, os.Stdin)
	} else {
		for _, a := range readerURIs {
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
	cfg := fetchConfig()
	m := make([]metrics.Metric, 0, 1)

	readers, err := openReaders(flag.Args())
	defer closeReaders(readers)

	for _, reader := range readers {
		metricDecoder := metrics.NewDecoder(reader)
		metricDecoder.MetricTime = cfg.MetricTime
		metricDecoder.ListItemKey = cfg.ListItemKey
		metricDecoder.KeyPrefix = metrics.GraphiteKey{cfg.Prefix}
		err = metricDecoder.Decode(&m)
		if err != nil {
			log.Print(err.Error())
		}
	}

	publisher, err := metrics.NewGraphitePublisher(
		cfg.Graphite, cfg.DryRun)
	if err != nil {
		log.Fatal(err)
	}

	publisher.Publish(m)
}
