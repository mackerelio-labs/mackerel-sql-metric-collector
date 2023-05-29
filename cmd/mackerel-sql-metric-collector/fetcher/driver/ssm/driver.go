package ssm

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher"
)

const (
	optWithDecryption = "withDecryption"
	optRegion         = "region"
	defaultRegion     = "ap-northeast-1"
)

func init() {
	fetcher.Register("ssm", &Driver{})
}

// Driver represents ...
type Driver struct{}

// Fetch is ...
func (d *Driver) Fetch(u *url.URL) ([]byte, error) {
	return d.FetchWithContext(context.Background(), u)
}

// FetchWithContext is ...
func (d *Driver) FetchWithContext(ctx context.Context, u *url.URL) ([]byte, error) {
	r := resolveRegion(u.Query())

	wd, err := resolveWithDecryption(u.Query())
	if err != nil {
		return nil, err
	}

	return fetchFromSSM(ctx, createSSMClient(r), u.Path, wd)
}

func resolveRegion(values url.Values) string {
	r := values.Get(optRegion)
	if r == "" {
		r = defaultRegion
	}
	return r
}

func resolveWithDecryption(values url.Values) (bool, error) {
	var wd bool
	if v := values.Get(optWithDecryption); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return wd, err
		}
		wd = b
	}
	return wd, nil
}

func createSSMClient(region string) *ssm.SSM {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(region)
	return ssm.New(sess)
}

func fetchFromSSM(ctx context.Context, client *ssm.SSM, name string, withDecryption bool) ([]byte, error) {
	res, err := client.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(withDecryption),
	})
	if err != nil {
		return nil, err
	}

	return []byte(*res.Parameter.Value), nil
}
