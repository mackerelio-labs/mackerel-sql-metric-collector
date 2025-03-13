package s3

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/mackerelio-labs/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher"
)

const (
	optRegionHint     = "regionHint"
	defaultRegionHint = "ap-northeast-1"
)

func init() {
	fetcher.Register("s3", &Driver{})
}

// Driver represents ...
type Driver struct{}

// Fetch is ...
func (d *Driver) Fetch(u *url.URL) ([]byte, error) {
	return d.FetchWithContext(context.Background(), u)
}

// FetchWithContext is ...
func (d *Driver) FetchWithContext(ctx context.Context, u *url.URL) ([]byte, error) {
	rh := resolveRegionHint(u)

	c, err := createS3Client(ctx, u.Host, rh)
	if err != nil {
		return nil, err
	}

	return fetchFromS3(ctx, c, u.Host, u.Path)
}

func resolveRegionHint(u *url.URL) string {
	rh := strings.TrimSpace(u.Query().Get(optRegionHint))
	if rh != "" {
		return rh
	}
	return defaultRegionHint
}

func createS3Client(ctx context.Context, bucket, regionHint string) (*manager.Downloader, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(regionHint))
	if err != nil {
		return nil, err
	}

	r, err := manager.GetBucketRegion(ctx, s3.NewFromConfig(cfg), bucket)
	if err != nil {
		return nil, err
	}
	cfg.Region = r

	return manager.NewDownloader(s3.NewFromConfig(cfg)), nil
}

func fetchFromS3(ctx context.Context, client *manager.Downloader, bucket, key string) ([]byte, error) {
	buf := manager.NewWriteAtBuffer([]byte{})

	key = strings.TrimPrefix(key, "/")

	_, err := client.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
