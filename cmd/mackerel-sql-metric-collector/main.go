package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"

	collector "github.com/mackerelio-labs/mackerel-sql-metric-collector"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/executor"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/executor/driver/cli"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/executor/driver/lambda"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/option"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/exporter"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/exporter/mackerel"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/exporter/stdout"
)

var revision string

// logger will replace by command-line args or environment vars.
var logger logr.Logger = stdr.New(log.Default())

func main() {
	exits(run(os.Args[0], os.Args[1:]))
}

func exits(err error) {
	if err != nil {
		logger.Error(err, "failed to execute")
		os.Exit(1)
	}
	os.Exit(0)
}

func run(name string, args []string) error {
	ctx := context.TODO() // TODO: signal handling.

	conf, err := option.Parse(name, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	return executor.ExecuteWithContext(ctx, detectExecutorName(), conf, handlerFunc(name))
}

func handlerFunc(name string) func(context.Context, *option.Config) error {
	return func(ctx context.Context, conf *option.Config) error {
		if err := conf.Load(ctx); err != nil {
			return err
		}
		logger, err := buildLogger(conf.LogFormat, conf.LogLevel)
		if err != nil {
			return err
		}

		var exp exporter.Exporter
		switch conf.Exporter {
		case stdout.Name:
			exp = stdout.NewExporter()
		case mackerel.Name:
			var err error
			exp, err = mackerel.NewExporter(conf.MackerelAPIKey, conf.MackerelAPIBase)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: unknown exporter", conf.Exporter)
		}

		c, err := collector.NewCollector(conf.CollectorConfig, exp, logger)
		if err != nil {
			return err
		}

		queries, err := loadQueryWithContext(ctx, conf.QueryFilePath)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("start %s", name), "revision", revision)

		return c.RunWithContext(ctx, queries)
	}
}

func detectExecutorName() string {
	if e := os.Getenv("EXECUTOR"); e != "" {
		return e

	}
	if os.Getenv("AWS_EXECUTION_ENV") != "" {
		return lambda.Name

	}
	return cli.Name
}
