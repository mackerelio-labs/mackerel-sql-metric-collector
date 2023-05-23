// Package fetcher provides functions to retrieve a string from arbitary URL.
package fetcher

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher/driver"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// Register is ...
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("register driver is nil")
	}
	if _, ok := isRegistered(name); ok {
		panic("register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Fetch is ...
func Fetch(u *url.URL) ([]byte, error) {
	return FetchWithContext(context.Background(), u)
}

// FetchWithContext is ...
func FetchWithContext(ctx context.Context, u *url.URL) ([]byte, error) {
	name := u.Scheme
	if name == "" {
		name = "file"
	}

	driver, ok := isRegistered(name)
	if !ok {
		return nil, fmt.Errorf("%s driver not registered", name)
	}

	return driver.FetchWithContext(ctx, u)
}

// IsRegistered is ...
func IsRegistered(name string) bool {
	_, ok := isRegistered(name)
	return ok
}

func isRegistered(name string) (driver.Driver, bool) {
	driver, ok := drivers[name]
	return driver, ok
}
