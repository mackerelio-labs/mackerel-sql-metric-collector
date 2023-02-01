package mackerel

import (
	"context"

	"github.com/mackerelio/mackerel-client-go"
)

const (
	// Name defines this exporter name.
	Name = "mackerel"

	userAgent = "mackerel-sql-metric-collector"
)

// Exporter represents ...
type Exporter struct {
	client *mackerel.Client
}

// NewExporter is ...
func NewExporter(apiKey, apiBase string) (*Exporter, error) {
	client, err := newMackerelClient(apiKey, apiBase)
	if err != nil {
		return nil, err
	}

	client.UserAgent = userAgent

	return &Exporter{
		client: client,
	}, nil

}

// Export is ...
func (e *Exporter) Export(service string, metrics []*mackerel.MetricValue) error {
	return e.ExportWithContext(context.Background(), service, metrics)
}

// ExportWithContext is ...
func (e *Exporter) ExportWithContext(ctx context.Context, service string, metrics []*mackerel.MetricValue) error {
	// TODO: mackerel.Client does not support context.
	return e.client.PostServiceMetricValues(service, metrics)
}

func newMackerelClient(apiKey, apiBase string) (*mackerel.Client, error) {
	if apiBase != "" {
		return mackerel.NewClientWithOptions(apiKey, apiBase, false)
	}

	return mackerel.NewClient(apiKey), nil
}
