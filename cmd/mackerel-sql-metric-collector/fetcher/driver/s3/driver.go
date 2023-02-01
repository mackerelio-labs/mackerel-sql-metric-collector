package s3

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hatena/mackerel-sql-metric-collector/cmd/mackerel-sql-metric-collector/fetcher"
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

func createS3Client(ctx context.Context, bucket, regionHint string) (*s3manager.Downloader, error) {
	sess := session.Must(session.NewSession())

	r, err := s3manager.GetBucketRegion(ctx, sess, bucket, regionHint)
	if err != nil {
		return nil, err
	}
	sess.Config.Region = aws.String(r)

	return s3manager.NewDownloader(sess), nil
}

func fetchFromS3(ctx context.Context, client *s3manager.Downloader, bucket, key string) ([]byte, error) {
	buf := &aws.WriteAtBuffer{}

	_, err := client.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
