package ssm

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

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

	client, err := createSSMClient(ctx, r)
	if err != nil {
		return nil, err
	}

	return fetchFromSSM(ctx, client, u.Path, wd)
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

func createSSMClient(ctx context.Context, region string) (*ssm.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return ssm.NewFromConfig(cfg), nil
}

func fetchFromSSM(ctx context.Context, client *ssm.Client, name string, withDecryption bool) ([]byte, error) {
	res, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(withDecryption),
	})
	if err != nil {
		return nil, err
	}

	return []byte(*res.Parameter.Value), nil
}
