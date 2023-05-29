package main

import (
	_ "github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher/driver/file"
	_ "github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher/driver/s3"
	_ "github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher/driver/ssm"
)
