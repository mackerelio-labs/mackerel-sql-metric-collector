package main

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher"
	"github.com/hatena/mackerel-sql-metric-collector/query"
	"github.com/hatena/mackerel-sql-metric-collector/query/valuekey"
	"gopkg.in/yaml.v2"
)

func loadQueryWithContext(ctx context.Context, path string) ([]query.Query, error) {
	var err error

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	var data []byte
	if u.Path == "" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = fetcher.FetchWithContext(ctx, u)
	}
	if err != nil {
		return nil, err
	}

	queries, err := buildValuekeyQueriesFromYAMLString(data)
	if err != nil {
		return nil, err
	}

	return convertToQuery(queries), nil
}

func convertToQuery(vkQueries []*valuekey.Query) []query.Query {
	queries := make([]query.Query, len(vkQueries))
	for i, v := range vkQueries {
		queries[i] = v
	}
	return queries
}

func buildValuekeyQueriesFromYAMLString(str []byte) ([]*valuekey.Query, error) {
	queries := []*valuekey.Query{}
	err := yaml.Unmarshal(str, &queries)
	if err != nil {
		return nil, err
	}
	return queries, nil
}
