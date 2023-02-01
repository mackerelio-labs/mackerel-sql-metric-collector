package collector

import "testing"

func TestParseDSN(t *testing.T) {
	testCases := []struct {
		dsn            string
		driverName     string
		dataSourceName string
		isValidDSN     bool
	}{
		{
			dsn:        "invalid dsn",
			isValidDSN: false,
		},
		{
			dsn:            "postgres://host=localhost port=5432 user=MYUSER password=MYPASSWORD dbname=MYAPP sslmode=disable",
			driverName:     "postgres",
			dataSourceName: "host=localhost port=5432 user=MYUSER password=MYPASSWORD dbname=MYAPP sslmode=disable",
			isValidDSN:     true,
		},
		{
			dsn:            "bigquery://project/location/dataset",
			driverName:     "bigquery",
			dataSourceName: "bigquery://project/location/dataset",
			isValidDSN:     true,
		},
	}

	for _, tc := range testCases {
		driverName, dataSourceName, err := parseDSN(tc.dsn)

		if tc.isValidDSN && err != nil {
			t.Errorf("%q should not raise error: %v", tc.dsn, err)
		} else if !tc.isValidDSN && err == nil {
			t.Errorf("%q should raise error", tc.dsn)
		}

		if expected := tc.driverName; driverName != expected {
			t.Errorf("driverName should be %q but got %q", expected, driverName)
		}

		if expected := tc.dataSourceName; dataSourceName != expected {
			t.Errorf("dataSourceName should be %q but got %q", expected, dataSourceName)
		}
	}
}
