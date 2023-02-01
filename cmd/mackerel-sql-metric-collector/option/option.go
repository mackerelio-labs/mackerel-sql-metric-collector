// Package option contains options for the SQL Metric Collector handler.
package option

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	collector "github.com/hatena/mackerel-sql-metric-collector"
	"github.com/hatena/mackerel-sql-metric-collector/exporter/mackerel"
	"github.com/hatena/mackerel-sql-metric-collector/exporter/stdout"
)

// FIXME: It contains expected values to usage tag of both LogFormat and LogLevel.
var exporterNames = []string{mackerel.Name, stdout.Name} // nolint

// HandlerOptions is used to configure the handler.
type HandlerOptions struct {
	DSNRef            string `json:"dsn" flag:"dsn" usage:"datasource name"`
	DefaultServiceRef string `json:"default-service" flag:"default-service" usage:"default mackerel service ^name^"`
	MaxConcurrency    int    `json:"max-concurrency" flag:"max-concurrency" usage:"maximum ^number^ of concurrent queries"`

	QueryFilePath      string `json:"query-file" flag:"query-file" usage:"query yaml ^filename^"`
	MackerelAPIKeyRef  string `json:"mackerel-apikey" flag:"mackerel-apikey" usage:"mackerel ^apikey^"`
	MackerelAPIBaseRef string `json:"mackerel-apibase" flag:"mackerel-apibase" usage:"mackerel apibase ^url^"`
	Exporter           string `json:"exporter" flag:"exporter" usage:"exporter to ^backend^ service [mackerel, stdout]"`
	LogFormat          string `json:"log-format" flag:"log-format" usage:"log ^format^ [console, json]"`
	LogLevel           string `json:"log-level" flag:"log-level" usage:"log ^level^ [info, error]"`
}

var defaultHandlerOptions = HandlerOptions{
	MaxConcurrency: 5,
	Exporter:       mackerel.Name,
	LogFormat:      "console",
	LogLevel:       "info",
}

var methods = map[reflect.Kind]string{
	reflect.Bool:   "BoolVar",
	reflect.Int:    "IntVar",
	reflect.String: "StringVar",
}

// Flags returns a pointer to flag.FlagSet that sets option values to corresponding to opts.
func Flags(name string, opts interface{}) (*flag.FlagSet, error) {
	c := flag.NewFlagSet(name, flag.ContinueOnError)
	p := reflect.ValueOf(opts)
	for i, field := range reflect.VisibleFields(p.Elem().Type()) {
		tag := field.Tag.Get("flag")
		if tag == "-" {
			continue
		}
		name := field.Name
		if tag != "" {
			name = tag
		}

		// When an option is typed any primitive types,
		// it should use type specific functions such as flag.VarBool.
		if method, ok := methods[field.Type.Kind()]; ok {
			v := reflect.ValueOf(c)
			m := v.MethodByName(method)
			m.Call([]reflect.Value{
				p.Elem().Field(i).Addr(),
				reflect.ValueOf(name),
				p.Elem().Field(i),
				reflect.ValueOf(usage(field)),
			})
			continue
		}

		// Otherwise, if a type that is implements flag.Value interface, can use flag.Var function.
		iface := reflect.TypeOf((*flag.Value)(nil)).Elem()
		if reflect.PtrTo(field.Type).Implements(iface) {
			v := reflect.ValueOf(c)
			m := v.MethodByName("Var")
			m.Call([]reflect.Value{
				p.Elem().Field(i).Addr(),
				reflect.ValueOf(name),
				reflect.ValueOf(usage(field)),
			})
			continue
		}
		return nil, fmt.Errorf("field '%s': unsupported kind", field.Name)
	}
	return c, nil
}

func usage(f reflect.StructField) string {
	s := f.Tag.Get("usage")
	// The flag package uses back-quoted string in an usage as a name of an option value.
	// But Go's literal string cannot escape backquotes in a string.
	return strings.ReplaceAll(s, "^", "`")
}

// Config contains configuration values of the SQL Metric Collector.
type Config struct {
	CollectorConfig *collector.Config
	QueryFilePath   string
	MackerelAPIKey  string
	MackerelAPIBase string
	Exporter        string
	LogFormat       string
	LogLevel        string

	DSNRef             string
	DefaultServiceRef  string
	MackerelAPIKeyRef  string
	MackerelAPIBaseRef string
}

// ToConfig returns Config that is initialized with corresponding fields of opts.
func (opts *HandlerOptions) ToConfig() *Config {
	return &Config{
		CollectorConfig: &collector.Config{
			MaxConcurrency: opts.MaxConcurrency,
		},
		QueryFilePath: opts.QueryFilePath,
		Exporter:      opts.Exporter,
		LogFormat:     opts.LogFormat,
		LogLevel:      opts.LogLevel,

		DSNRef:             opts.DSNRef,
		DefaultServiceRef:  opts.DefaultServiceRef,
		MackerelAPIKeyRef:  opts.MackerelAPIKeyRef,
		MackerelAPIBaseRef: opts.MackerelAPIBaseRef,
	}
}

// DeepCopy returns a copy of c.
func (c *Config) DeepCopy() *Config {
	cc := *c.CollectorConfig
	nc := *c
	nc.CollectorConfig = &cc
	return &nc
}

// Merge updates each fields of c with corresponding field of opts if opts's field value is not zero.
func (c *Config) Merge(opts *HandlerOptions) {
	updateValue(&c.CollectorConfig.MaxConcurrency, opts.MaxConcurrency)
	updateValue(&c.QueryFilePath, opts.QueryFilePath)
	updateValue(&c.Exporter, opts.Exporter)
	updateValue(&c.LogFormat, opts.LogFormat)
	updateValue(&c.LogLevel, opts.LogLevel)
	updateValue(&c.DSNRef, opts.DSNRef)
	updateValue(&c.DefaultServiceRef, opts.DefaultServiceRef)
	updateValue(&c.MackerelAPIKeyRef, opts.MackerelAPIKeyRef)
	updateValue(&c.MackerelAPIBaseRef, opts.MackerelAPIBaseRef)
}

func updateValue[T comparable](p *T, v T) {
	var zero T
	if v != zero {
		*p = v
	}
}

// Parse parses the args; parsed values are set into corresponding fields of Config.
func Parse(name string, args []string) (*Config, error) {
	opts := defaultHandlerOptions
	flags, err := Flags(name, &opts)
	if err != nil {
		return nil, err
	}

	var setErr error
	flags.VisitAll(func(f *flag.Flag) {
		env := envVar(f.Name)
		if s := os.Getenv(env); s != "" {
			err := f.Value.Set(s)
			if err != nil {
				setErr = fmt.Errorf("failed to set flag %s via env(%s): %w", f.Name, env, err)
			}
		}
	})
	if setErr != nil {
		return nil, setErr
	}

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return opts.ToConfig(), nil
}

func envVar(s string) string {
	return strings.ReplaceAll(strings.ToUpper(s), "-", "_")
}
