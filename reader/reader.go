package reader

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"io"
	_ "log"
	"net/url"
	"path/filepath"
)

// Return an io.ReadCloser instance using the gocloud.dev/blob NewReader method.
// If path is not fully qualified URI then assume the file:// scheme.
func NewReader(ctx context.Context, path string) (io.ReadCloser, error) {

	u, err := url.Parse(path)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse '%s', %v", path, err)
	}

	if u.Scheme == "" {

		u.Scheme = "file"
		abs_path, err := filepath.Abs(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive absolute path for '%s', %v", path, err)
		}

		u.Path = abs_path
		path = u.String()
	}

	root := filepath.Dir(path)
	fname := filepath.Base(path)

	b, err := blob.OpenBucket(ctx, root)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket for '%s', %v", root, err)
	}

	defer b.Close()

	return b.NewReader(ctx, fname, nil)
}
