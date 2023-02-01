package query

import (
	"context"
	"database/sql"

	"github.com/go-logr/logr"
	"github.com/mackerelio/mackerel-client-go"
)

// Query represents ...
type Query interface {
	Execute(*sql.DB, logr.Logger) ([]*mackerel.MetricValue, error)
	ExecuteWithContext(context.Context, *sql.DB, logr.Logger) ([]*mackerel.MetricValue, error)
	GetService() string
}
