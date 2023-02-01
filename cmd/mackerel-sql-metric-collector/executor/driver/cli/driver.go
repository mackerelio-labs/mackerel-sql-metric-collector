package cli

import (
	"context"

	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/executor"
	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/option"
)

const (
	// Name is ...
	Name = "cli"
)

func init() {
	executor.Register(Name, &Driver{})
}

// Driver represents ...
type Driver struct{}

// Invoke is ...
func (d *Driver) Invoke(ctx context.Context, c *option.Config, handler func(context.Context, *option.Config) error) error {
	return handler(ctx, c)
}
