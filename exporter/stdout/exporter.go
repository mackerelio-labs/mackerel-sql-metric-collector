package stdout

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/mackerelio/mackerel-client-go"
)

// Name defines this exporter name.
const Name = "stdout"

// Exporter represents ...
type Exporter struct{}

// NewExporter is ...
func NewExporter() *Exporter {
	return &Exporter{}
}

// Export is ...
func (e *Exporter) Export(service string, metrics []*mackerel.MetricValue) error {
	return e.ExportWithContext(context.Background(), service, metrics)
}

// ExportWithContext is ...
func (e *Exporter) ExportWithContext(_ context.Context, _ string, metrics []*mackerel.MetricValue) error {
	for _, m := range metrics {
		err := printValue(os.Stdout, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func printValue(w io.Writer, metric *mackerel.MetricValue) error {
	var v float64

	switch i := metric.Value.(type) {
	case int32:
		v = float64(i)
	case uint32:
		v = float64(i)
	case float32:
		v = float64(i)
	case int64: // rounded
		v = float64(i)
	case uint64: // rounded
		v = float64(i)
	case float64:
		v = i
	default:
		v = math.NaN()

	}

	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("invalid metric.Value: key = %s, metric.Value = (%T)%v", metric.Name, metric.Value, metric.Value)
	}

	fmt.Fprintf(w, "%s\t%f\t%d\n", metric.Name, v, metric.Time) // nolint

	return nil
}
