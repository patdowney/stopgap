package metrics

import (
	"fmt"
	"io"
	"time"
)

// Metric ...
type Metric struct {
	Key   string
	Value string
	Time  time.Time
}

// DumpMetrics ...
func DumpMetrics(writer io.Writer, metrics []Metric) {
	for _, m := range metrics {
		fmt.Fprintf(writer, "%v: %v\n", m.Key, m.Value)
	}
}
