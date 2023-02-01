package mackerel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestExporterExport(t *testing.T) {
	const (
		apiKey      = "xxx"
		serviceName = "Service1"
	)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/services/"+serviceName+"/tsdb" {
			http.Error(w, r.URL.Path, http.StatusBadRequest)
			return
		}
		w.Write([]byte("OK")) // nolint
	}))
	t.Cleanup(s.Close)

	e, err := NewExporter(apiKey, s.URL)
	if err != nil {
		t.Fatal("NewExporter: ", err)
	}
	err = e.Export(serviceName, []*mackerel.MetricValue{})
	if err != nil {
		t.Errorf("Export: got %v", err)
	}
}
