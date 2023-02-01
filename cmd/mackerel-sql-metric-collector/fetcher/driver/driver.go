package driver

import (
	"context"
	"net/url"
)

// Driver represents ...
type Driver interface {
	Fetch(*url.URL) ([]byte, error)
	FetchWithContext(context.Context, *url.URL) ([]byte, error)
}
