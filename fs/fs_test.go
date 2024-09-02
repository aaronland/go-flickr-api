package fs

import (
	"context"
	"flag"
	"image"
	_ "image/jpeg"
	io_fs "io/fs"
	"log/slog"
	"net/url"
	"testing"

	"github.com/aaronland/go-flickr-api/client"
)

var client_uri = flag.String("client-uri", "", "...")

func TestFS(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	if *client_uri == "" {
		slog.Info("-client-uri flag not set, skipping test.")
		t.Skip()
	}

	ctx := context.Background()
	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		t.Fatalf("Failed to create new client, %v", err)
	}

	fs := New(ctx, cl)

	fl, err := fs.Open("53961664838")

	if err != nil {
		t.Fatalf("Failed to open , %v", err)
	}

	_, _, err = image.Decode(fl)

	if err != nil {
		t.Fatalf("Failed to decode image, %v", err)
	}

	walk_func := func(path string, d io_fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		slog.Info("Walk", "path", path, "d", d)
		return nil
	}

	u := url.Values{}
	u.Set("method", "flickr.photosets.getPhotos")
	u.Set("photoset_id", "72177720319945125")
	u.Set("user_id", "35034348999@N01")

	q := u.Encode()

	err = io_fs.WalkDir(fs, q, walk_func)

	if err != nil {
		t.Fatalf("Failed to walk '%s', %v", q, err)
	}
}
