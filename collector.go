package collector

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/exporter"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/query"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
	_ "github.com/speee/go-athena"  // AWS Athena driver
	"golang.org/x/sync/errgroup"
	_ "gorm.io/driver/bigquery/driver" // BigQuery driver
)

// Collector represents ...
type Collector struct {
	config   *Config
	exporter exporter.Exporter
	logger   logr.Logger
}

// NewCollector is ...
func NewCollector(conf *Config, exporter exporter.Exporter, logger logr.Logger) (*Collector, error) {
	return &Collector{
		config:   conf,
		exporter: exporter,
		logger:   logger,
	}, nil
}

// Run collect and post metrics.
func (c *Collector) Run(queries []query.Query) error {
	return c.RunWithContext(context.Background(), queries)
}

// RunWithContext collect and post metrics with context.Context.
func (c *Collector) RunWithContext(ctx context.Context, queries []query.Query) error {
	db, err := openDataSource(c.config.DSN)
	if err != nil {
		return err
	}
	defer db.Close() // nolint

	queue := make(chan struct{}, c.config.MaxConcurrency-1)

	eg := &errgroup.Group{} // Create *errgroup.Group as we want to run all queries.
	for _, q := range queries {
		queue <- struct{}{}
		q := q
		eg.Go(func() error {
			defer func() {
				<-queue
			}()
			metrics, err := q.ExecuteWithContext(ctx, db, c.logger)
			if err != nil {
				return err
			}
			return c.exporter.ExportWithContext(ctx, c.detectService(q), metrics)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (c *Collector) detectService(q query.Query) string {
	s := q.GetService()
	if s == "" {
		return c.config.DefaultService
	}
	return s
}

func openDataSource(dsn string) (*sql.DB, error) {
	driverName, dataSourceName, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

const (
	bigQueryDriverName = "bigquery"
	bigQueryDSNPrefix  = "bigquery://"
)

func parseDSN(dsn string) (string, string, error) {
	var driverName, dataSourceName string

	// see. https://github.com/go-sql-driver/mysql#dsn-data-source-name
	parts := strings.SplitN(dsn, "://", 2)
	if len(parts) != 2 {
		return driverName, dataSourceName, fmt.Errorf("invalid dsn: %s", dsn)
	}

	driverName = parts[0]
	dataSourceName = parts[1]

	// for gorm.io/driver/bigquery/driver
	if driverName == bigQueryDriverName {
		dataSourceName = bigQueryDSNPrefix + dataSourceName
	}

	return driverName, dataSourceName, nil
}
