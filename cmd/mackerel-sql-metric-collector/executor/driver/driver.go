package driver

import (
	"context"

	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/option"
)

// Driver represents ...
type Driver interface {
	Invoke(context.Context, *option.Config, func(context.Context, *option.Config) error) error
}
