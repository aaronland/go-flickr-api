package fs

import (
	"context"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	io_fs "io/fs"
	"log/slog"
	"net/url"
	"testing"

	"github.com/aaronland/go-flickr-api/client"
)

var client_uri = flag.String("client-uri", "", "...")

func TestRePhoto(t *testing.T) {

	tests := []string{
		"53961664838",
		"/65535/53961664838_49a7d74e87_o.jpg",
		"method=flickr.photosets.getPhotos&photoset_id=72177720319945125&user_id=35034348999%40N01/#/65535/53961664838_49a7d74e87_o.jpg",
	}

	for _, str := range tests {

		if !re_photo.MatchString(str) {
			t.Fatalf("Failed to match '%s'", str)
		}
	}

}

func TestReURL(t *testing.T) {

	tests := map[string]string{
		"/65535/53961664838_49a7d74e87_o.jpg": "/65535/53961664838_49a7d74e87_o.jpg",
		"method=flickr.photosets.getPhotos&photoset_id=72177720319945125&user_id=35034348999%40N01/#/65535/53961664838_49a7d74e87_o.jpg": "/65535/53961664838_49a7d74e87_o.jpg",
	}

	for str, expected := range tests {

		if !re_url.MatchString(str) {
			t.Fatalf("Failed to match '%s'", str)
		}

		m := re_url.FindStringSubmatch(str)

		if m[1] != expected {
			t.Fatalf("Unexpected match for '%s', expected '%s' but got '%s'", str, expected, m[1])
		}
	}

}

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

		slog.Info("Walk", "path", path, "d", d.Name())

		r, err := fs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open '%s', %w", path, err)
		}

		defer r.Close()
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
