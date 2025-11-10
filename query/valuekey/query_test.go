package valuekey

import (
	"io"
	"log"
	"os"
	"slices"
	"strings"
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
	t.Parallel()
	testCases := map[string]struct {
		query *Query
		want  []*mackerel.MetricValue
	}{
		"basic": {
			query: &Query{
				KeyPrefix: "agent",
				ValueKey: map[string]string{
					"versions.#{agent_version}": "host_num",
				},
				SQL: "SELECT * FROM dummy",
			},
			want: []*mackerel.MetricValue{
				{
					Name:  "agent.versions.0_1_0",
					Time:  nowFunc().Unix(),
					Value: int64(10),
				},
			},
		},
		"defaultValue": {
			query: &Query{
				KeyPrefix: "agent",
				ValueKey: map[string]string{
					"versions.#{agent_version}": "host_num",
				},
				DefaultValue: map[string]float64{
					"versions.0_1_1": 0.0,
				},
				SQL: "SELECT * FROM dummy",
			},
			want: []*mackerel.MetricValue{
				{
					Name:  "agent.versions.0_1_0",
					Time:  nowFunc().Unix(),
					Value: int64(10),
				},
				{
					Name:  "agent.versions.0_1_1",
					Time:  nowFunc().Unix(),
					Value: float64(0.0),
				},
			},
		},
		"defaultValue_template": {
			query: &Query{
				KeyPrefix: "agent",
				ValueKey: map[string]string{
					"versions.#{agent_version}": "host_num",
				},
				DefaultValue: map[string]float64{
					"versions.#{agent_version}": 0.0,
				},
				SQL: "SELECT * FROM dummy",
			},
			want: []*mackerel.MetricValue{
				{
					Name:  "agent.versions.0_1_0",
					Time:  nowFunc().Unix(),
					Value: int64(10),
				},
				{
					Name:  "agent.versions.0_1_1",
					Time:  nowFunc().Unix(),
					Value: float64(0.0),
				},
				{
					Name:  "agent.versions.0_1_2",
					Time:  nowFunc().Unix(),
					Value: float64(0.0),
				},
			},
		},
		"defaultValue_non_exist_key": {
			query: &Query{
				KeyPrefix: "agent",
				ValueKey: map[string]string{
					"versions.#{agent_version}": "host_num",
				},
				DefaultValue: map[string]float64{
					"versions.0_2_0": 1.0,
				},
				SQL: "SELECT * FROM dummy",
			},
			want: []*mackerel.MetricValue{
				{
					Name:  "agent.versions.0_1_0",
					Time:  nowFunc().Unix(),
					Value: int64(10),
				},
				{
					Name:  "agent.versions.0_2_0",
					Time:  nowFunc().Unix(),
					Value: float64(1.0),
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := stdr.New(log.New(io.Discard, "", 0))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("sqlmock.New: ", err)
			}
			t.Cleanup(func() {
				db.Close() // nolint
			})
			columns := []string{"agent_version", "host_num"}
			rows := sqlmock.NewRows(columns).AddRow("0.1.0", 10).AddRow("0.1.1", nil).AddRow("0.1.2", nil)
			mock.ExpectQuery("SELECT (.+) FROM (.+)").WillReturnRows(rows)

			values, err := tc.query.Execute(db, logger)
			if err != nil {
				t.Errorf("Execute: got %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("ExpectationsWereMet: got %v", err)
			}
			slices.SortStableFunc(values, func(a, b *mackerel.MetricValue) int {
				return strings.Compare(a.Name, b.Name)
			})
			if diff := cmp.Diff(tc.want, values); diff != "" {
				t.Errorf("Execute: (-want, +got)\n%s", diff)
			}
		})
	}
}
