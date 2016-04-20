package metrics

import "github.com/marpaia/graphite-golang"

// MetricPublisher ...
type MetricPublisher interface {
	Publish(metrics []Metric)
}

// GraphiteConfig ...
type GraphiteConfig struct {
	Host string
	Port int
}

// GraphitePublisher ...
type GraphitePublisher struct {
	Config         GraphiteConfig
	graphiteClient *graphite.Graphite
}

// Publish ...
func (p *GraphitePublisher) Publish(metrics []Metric) {
	for _, m := range metrics {
		gmetric := graphite.Metric{
			Name:      m.Key,
			Value:     m.Value,
			Timestamp: m.Time.Unix()}
		p.graphiteClient.SendMetric(gmetric)
	}
}

// NewGraphitePublisher ...
func NewGraphitePublisher(config GraphiteConfig, dryRun bool) (MetricPublisher, error) {
	p := &GraphitePublisher{Config: config}
	p.graphiteClient = graphite.NewGraphiteNop(
		config.Host, config.Port)

	if !dryRun {
		var err error
		p.graphiteClient, err = graphite.NewGraphite(
			config.Host, config.Port)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}
