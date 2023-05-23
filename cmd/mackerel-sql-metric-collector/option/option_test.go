package option

import (
	"io"
	"log"
	"net/url"
	"os"
	"reflect"
	"testing"

	collector "github.com/mackerelio-labs/mackerel-sql-metric-collector"
	"github.com/mackerelio-labs/mackerel-sql-metric-collector/exporter/stdout"
)

type urlValue url.URL

func (v *urlValue) String() string {
	return (*url.URL)(v).String()
}

func (v *urlValue) Set(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*(*url.URL)(v) = *u
	return nil
}

func TestFlags(t *testing.T) {
	var opts struct {
		Name     string   // no tag
		Age      int      `flag:"age"`
		Exported bool     `flag:"exported"`
		URL      urlValue `flag:"url"`
		Ignored  float64  `flag:"-"`
	}

	f, err := Flags("xx", &opts)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Parse([]string{"-Name", "option1", "-age", "20", "-exported", "-url", "https://example.com"}); err != nil {
		t.Fatal(err)
	}
	if want := "option1"; opts.Name != want {
		t.Errorf("Name = %q; want %q", opts.Name, want)
	}
	if want := 20; opts.Age != want {
		t.Errorf("Age = %d; want %d", opts.Age, want)
	}
	if want := true; opts.Exported != want {
		t.Errorf("Exported = %t; want %t", opts.Exported, want)
	}
	if want := "https://example.com"; opts.URL.String() != want {
		t.Errorf("URL = %v; want %s", opts.URL, want)
	}
}

func TestFlags_error(t *testing.T) {
	var opts struct {
		Age int      `flag:"age"`
		URL urlValue `flag:"url"`
	}
	f, err := Flags("xx", &opts)
	if err != nil {
		t.Fatal(err)
	}

	// Stop default usage and error outputs because it is pretty chatty.
	f.Usage = func() {}
	f.SetOutput(io.Discard)

	if err := f.Parse([]string{"-age", "2x"}); err == nil {
		t.Errorf("set '2x' to IntVar: should be a parsing error")
	}
	if err := f.Parse([]string{"-url", ".https:///"}); err == nil {
		t.Errorf("set '.https:///' to URL: should be a parsing error")
	}
}

func Example_usage() {
	opts := struct {
		S string `flag:"str" usage:"^s^tring"`
		N int    `flag:"num" usage:"^n^umber"`
		V bool   `flag:"value" usage:"bool"`
	}{
		N: 5,
	}
	f, err := Flags("xx", &opts)
	if err != nil {
		log.Fatal(err)
	}
	f.SetOutput(os.Stdout)
	f.PrintDefaults()
	// Output:
	//   -num n
	//     	number (default 5)
	//   -str s
	//     	string
	//   -value
	//     	bool
}

func TestHandlerOptions_ToConfig(t *testing.T) {
	opts := &HandlerOptions{
		DSNRef:             "host=127.1 port=123 user=root",
		DefaultServiceRef:  "s3://example/service",
		MaxConcurrency:     10,
		QueryFilePath:      "file",
		MackerelAPIKeyRef:  "ssm://mackerel/key",
		MackerelAPIBaseRef: "ssm://mackerel/base",
		Exporter:           stdout.Name,
		LogFormat:          "json",
		LogLevel:           "error",
	}
	c := opts.ToConfig()
	want := &Config{
		CollectorConfig: &collector.Config{
			MaxConcurrency: 10,
		},
		QueryFilePath: "file",
		Exporter:      stdout.Name,
		LogFormat:     "json",
		LogLevel:      "error",

		DSNRef:             "host=127.1 port=123 user=root",
		DefaultServiceRef:  "s3://example/service",
		MackerelAPIKeyRef:  "ssm://mackerel/key",
		MackerelAPIBaseRef: "ssm://mackerel/base",
	}
	if !reflect.DeepEqual(c, want) {
		t.Errorf("ToConfig() = %+v; but want %+v", c, want)
	}
}

func TestConfig_Merge(t *testing.T) {
	c := &Config{
		CollectorConfig: &collector.Config{
			DSN:            "host=127.1 port=123 user=root",
			DefaultService: "Service1",
			MaxConcurrency: 10,
		},
		QueryFilePath:   "file",
		MackerelAPIKey:  "xxxx",
		MackerelAPIBase: "https://mackerel.io",
		Exporter:        stdout.Name,
		LogFormat:       "json",
		LogLevel:        "error",

		DSNRef:             "host=127.1 port=123 user=root",
		DefaultServiceRef:  "s3://example/service",
		MackerelAPIKeyRef:  "ssm://mackerel/key",
		MackerelAPIBaseRef: "ssm://mackerel/base",
	}
	opts := &HandlerOptions{
		DSNRef:             "host=127.2 port=123 user=root",
		DefaultServiceRef:  "s3://example/service2",
		MaxConcurrency:     20,
		QueryFilePath:      "file2",
		MackerelAPIKeyRef:  "ssm://mackerel/key2",
		MackerelAPIBaseRef: "ssm://mackerel/base2",
		Exporter:           stdout.Name,
		LogFormat:          "console",
		LogLevel:           "info",
	}

	// Here makes a expected Config value.
	// There are some unaffected fields in the expected Config.
	// Thus these fields should be filled with original Config's fields.
	want := opts.ToConfig()
	want.CollectorConfig.DSN = c.CollectorConfig.DSN
	want.CollectorConfig.DefaultService = c.CollectorConfig.DefaultService
	want.MackerelAPIKey = c.MackerelAPIKey
	want.MackerelAPIBase = c.MackerelAPIBase
	// Special case: If there are any boolean flags, they should keeps original value because new value, false, is zero.
	// want.Xxx = true

	c.Merge(opts)
	if !reflect.DeepEqual(c, want) {
		t.Errorf("MergeTo() = %+v; but want %+v", c, want)
	}
}

func TestParse(t *testing.T) {
	const (
		dsn    = "host=127.1 port=123 user=root"
		apiKey = "a12345"
	)

	t.Setenv("MACKEREL_APIKEY", apiKey)
	c, err := Parse("", []string{"-dsn", dsn})
	if err != nil {
		t.Errorf("Parse: %v", err)
		return
	}
	if s := c.DSNRef; s != dsn {
		t.Errorf("DSNRef = %s; want %s", s, dsn)
	}
	if s := c.MackerelAPIKeyRef; s != apiKey {
		t.Errorf("MackerelAPIKeyRef = %s; want %s", s, apiKey)
	}
}
