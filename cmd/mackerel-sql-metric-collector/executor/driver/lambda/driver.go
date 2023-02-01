package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/executor"
	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/option"
)

const (
	// Name is ...
	Name = "lambda"
)

func init() {
	executor.Register(Name, &Driver{})
}

// Driver represents ...
type Driver struct{}

// Invoke is ...
func (d *Driver) Invoke(ctx context.Context, c *option.Config, handler func(context.Context, *option.Config) error) error {
	originalConfig := c.DeepCopy()
	lambda.StartWithContext(ctx, func(ctx context.Context, opts *option.HandlerOptions) error {
		c := originalConfig.DeepCopy()
		c.Merge(opts)
		return handler(ctx, c)
	})
	return nil
}
