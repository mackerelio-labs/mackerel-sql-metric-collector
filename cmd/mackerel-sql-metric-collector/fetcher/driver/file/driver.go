package file

import (
	"context"
	"net/url"
	"os"

	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher"
)

func init() {
	fetcher.Register("file", &Driver{})
}

// Driver represents ...
type Driver struct{}

// Fetch is ...
func (d *Driver) Fetch(u *url.URL) ([]byte, error) {
	return d.FetchWithContext(context.Background(), u)
}

// FetchWithContext is ...
func (d *Driver) FetchWithContext(_ context.Context, u *url.URL) ([]byte, error) {
	return os.ReadFile(u.Path)
}
