package valuekey

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/mackerelio/mackerel-client-go"
)

var valueKeyRE = regexp.MustCompile(`#\{([a-z\_]+)\}`)
var invalidMackerelMetricKeyCharsRE = regexp.MustCompile(`[^-a-zA-Z0-9_]`)
var commandExecRE = regexp.MustCompile(`\A\$\((.*)\)\z`)

// Query represents ...
type Query struct {
	KeyPrefix    string             `yaml:"keyPrefix"`
	ValueKey     map[string]string  `yaml:"valueKey"`
	DefaultValue map[string]float64 `yaml:"defaultValue,omitempty"`
	SQL          string             `yaml:"sql"`
	Params       []interface{}      `yaml:"params"`
	Service      string             `yaml:"service,omitempty"`
	Time         string             `yaml:"time"`
}

// Execute is ...
func (q *Query) Execute(db *sql.DB, logger logr.Logger) ([]*mackerel.MetricValue, error) {
	return q.ExecuteWithContext(context.Background(), db, logger)
}

var nowFunc = time.Now

// ExecuteWithContext is ...
func (q *Query) ExecuteWithContext(ctx context.Context, db *sql.DB, logger logr.Logger) ([]*mackerel.MetricValue, error) {
	rows, err := q.queryDBWithContext(ctx, db)
	if err != nil {
		return nil, err
	}

	metrics := make([]*mackerel.MetricValue, 0, len(q.ValueKey))
	now := nowFunc().Unix()

	for _, r := range rows {
		defaults := make(map[string]float64, len(q.DefaultValue))
		for k, v := range q.DefaultValue {
			vk, err := replaceValueKey(k, r)
			if err != nil {
				logger.Info(err.Error(), "query", q)
				continue
			}
			defaults[vk] = v
		}

		for k, v := range q.ValueKey {
			var err error

			vk, err := replaceValueKey(k, r)
			if err != nil {
				logger.Info(err.Error(), "query", q)
				continue
			}

			value, ok := r[v]
			if !ok {
				return nil, fmt.Errorf("%q not exists in columns", v)
			}
			if value == nil {
				defaultValue, ok := defaults[vk]
				if !ok {
					continue
				}
				value = defaultValue
			}

			if q.KeyPrefix != "" {
				vk = fmt.Sprintf("%s.%s", q.KeyPrefix, vk)
			}

			t := now
			if q.Time != "" {
				value, ok := r[q.Time]
				if !ok {
					return nil, fmt.Errorf("%q not exists in columns", q.Time)
				}
				t, ok = value.(int64)
				if !ok {
					return nil, fmt.Errorf("failed to convert %q to int64", q.Time)
				}
			}

			mv := mackerel.MetricValue{
				Name:  vk,
				Value: value,
				Time:  t,
			}

			metrics = append(metrics, &mv)
		}
	}

	return metrics, nil
}

// GetService is ...
func (q *Query) GetService() string {
	return q.Service
}

type dbRow map[string]interface{}

func (q *Query) queryDBWithContext(ctx context.Context, db *sql.DB) ([]dbRow, error) {
	params, err := evalParams(q.Params)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, q.SQL, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []dbRow

	for rows.Next() {
		var row = make([]interface{}, len(cols))
		var rowp = make([]interface{}, len(cols)) // Slice pointer to each column of row (row[i]).
		for i := 0; i < len(row); i++ {
			rowp[i] = &row[i]
		}

		err := rows.Scan(rowp...) // Pass a slice of column pointers as a parameters.
		if err != nil {
			return nil, err
		}

		rowMap := make(dbRow)
		for i, col := range cols {
			switch row[i].(type) {
			case []byte:
				row[i] = string(row[i].([]byte))

				n, err := strconv.Atoi(row[i].(string))
				if err == nil {
					row[i] = n
				} else {
					f, err := strconv.ParseFloat(row[i].(string), 64)
					if err == nil {
						row[i] = f
					}
				}
			}
			rowMap[col] = row[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func evalParams(params []interface{}) ([]interface{}, error) {
	evaluated := make([]interface{}, len(params))

	for i, p := range params {
		v, ok := p.(string)
		if !ok {
			evaluated[i] = p
			continue
		}

		var err error

		v = commandExecRE.ReplaceAllStringFunc(v, func(cmd string) string {
			matches := commandExecRE.FindStringSubmatch(cmd)

			if len(matches) != 2 {
				err = errors.New("command not found")
				return ""
			}

			cmd = matches[1]

			out, e := exec.Command("/bin/sh", "-c", cmd).Output()
			if e != nil {
				if err == nil {
					err = e
				}
				return ""
			}

			return strings.TrimSpace(string(out))
		})

		if err != nil {
			return nil, err
		}

		evaluated[i] = v
	}

	return evaluated, nil
}

func replaceValueKey(key string, row dbRow) (string, error) {
	var err error

	replaced := valueKeyRE.ReplaceAllStringFunc(key, func(match string) string {
		matches := valueKeyRE.FindStringSubmatch(match)

		if len(matches) != 2 {
			err = errors.New("ValueKey not found")
			return ""
		}

		col := matches[1]

		v, ok := row[col]
		if !ok {
			err = fmt.Errorf("%q not exists in columns", col)
			return ""
		}
		if v == nil {
			err = fmt.Errorf("%q value is nil", col)
			return ""
		}

		// convert query result value to string.
		// string, int64, int32, float64, bool
		s := strings.TrimSpace(fmt.Sprintf("%v", v))
		if s == "" {
			err = fmt.Errorf("%q is empty", col)
			return ""
		}

		return invalidMackerelMetricKeyCharsRE.ReplaceAllString(s, "_")
	})

	return replaced, err
}
