package valuekey

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-logr/stdr"
	"github.com/google/go-cmp/cmp"
	"github.com/mackerelio/mackerel-client-go"
)

func TestMain(m *testing.M) {
	nowFunc = func() time.Time {
		return time.Date(2022, 1, 2, 3, 4, 5, 6, time.Local)
	}
	os.Exit(m.Run())
}

func TestQueryExecute(t *testing.T) {
	logger := stdr.New(log.New(io.Discard, "", 0))
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("sqlmock.New: ", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	columns := []string{"agent_version", "host_num"}
	mock.ExpectQuery("SELECT (.+) FROM (.+)").WillReturnRows(sqlmock.NewRows(columns).AddRow("0.1.0", 10))

	q := &Query{
		KeyPrefix: "agent",
		ValueKey: map[string]string{
			"versions.#{agent_version}": "host_num",
		},
		SQL: "SELECT * FROM dummy",
	}
	values, err := q.Execute(db, logger)
	if err != nil {
		t.Errorf("Execute: got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ExpectationsWereMet: got %v", err)
	}
	want := []*mackerel.MetricValue{
		{
			Name:  "agent.versions.0_1_0",
			Time:  nowFunc().Unix(),
			Value: int64(10),
		},
	}
	if diff := cmp.Diff(want, values); diff != "" {
		t.Errorf("Execute: (-want, +got)\n%s", diff)
	}
}
