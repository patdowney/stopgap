package metrics

import (
	"fmt"
	"io"
	"time"
)

type Metric struct {
	Key   string
	Value string
	Time  time.Time
}

func DumpMetrics(writer io.Writer, metrics []Metric) {
	for _, m := range metrics {
		fmt.Fprintf(writer, "%v: %v\n", m.Key, m.Value)
	}
}
