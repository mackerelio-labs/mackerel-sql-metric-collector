package exporter

import (
	"context"

	"github.com/mackerelio/mackerel-client-go"
)

// Exporter represents ...
type Exporter interface {
	Export(string, []*mackerel.MetricValue) error
	ExportWithContext(context.Context, string, []*mackerel.MetricValue) error
}
