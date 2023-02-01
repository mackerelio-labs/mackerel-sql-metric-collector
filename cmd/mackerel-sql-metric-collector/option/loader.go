package option

import (
	"context"
	"net/url"
	"strings"

	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher"
)

// Load resolves all indirect values in c.
func (c *Config) Load(ctx context.Context) error {
	var f errFetcher
	c.CollectorConfig.DSN = f.FetchString(ctx, c.DSNRef)
	c.CollectorConfig.DefaultService = f.FetchString(ctx, c.DefaultServiceRef)
	c.MackerelAPIKey = f.FetchString(ctx, c.MackerelAPIKeyRef)
	c.MackerelAPIBase = f.FetchString(ctx, c.MackerelAPIBaseRef)
	return f.err
}

type errFetcher struct {
	err error
}

func (f *errFetcher) FetchString(ctx context.Context, s string) string {
	if f.err != nil {
		return ""
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || !fetcher.IsRegistered(u.Scheme) {
		return s
	}

	d, err := fetcher.FetchWithContext(ctx, u)
	if err != nil {
		f.err = err
		return ""
	}
	return strings.TrimSpace(string(d))
}
