package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/executor/driver"
	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/option"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// Register is ...
func Register(name string, d driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d == nil {
		panic("register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("register called twice for driver " + name)
	}
	drivers[name] = d
}

// Execute is ...
func Execute(driverName string, c *option.Config, handler func(context.Context, *option.Config) error) error {
	return ExecuteWithContext(context.Background(), driverName, c, handler)
}

// ExecuteWithContext is ...
func ExecuteWithContext(ctx context.Context, driverName string, c *option.Config, handler func(context.Context, *option.Config) error) error {
	d, ok := drivers[driverName]
	if !ok {
		return fmt.Errorf("%s driver not registered", driverName)
	}

	return d.Invoke(ctx, c, handler)
}
