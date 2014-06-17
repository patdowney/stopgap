package metrics

import "github.com/marpaia/graphite-golang"

func PublishMetrics(gclient *graphite.Graphite, metrics []Metric) {
	for _, m := range metrics {
		gmetric := graphite.Metric{
			Name:      m.Key,
			Value:     m.Value,
			Timestamp: m.Time.Unix()}
		gclient.SendMetric(gmetric)
	}
}
