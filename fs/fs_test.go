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

type apiTest struct {
	photo_id    string
	user_id     string
	photoset_id string
}

var client_uri = flag.String("client-uri", "", "...")

func TestMatchesPhotoId(t *testing.T) {

	tests := []string{
		"53961664838",
		"/65535/53961664838_49a7d74e87_o.jpg",
		"65535/53961664838_49a7d74e87_o.jpg",
		"method=flickr.photosets.getPhotos&photoset_id=72177720319945125&user_id=35034348999%40N01/#/65535/53961664838_49a7d74e87_o.jpg",
	}

	for _, str := range tests {

		if !MatchesPhotoId(str) {
			t.Fatalf("Failed to match photo ID '%s'", str)
		}
	}

}

func TestMatchesPhotoURL(t *testing.T) {

	tests := map[string]string{
		"/65535/53961664838_49a7d74e87_o.jpg": "/65535/53961664838_49a7d74e87_o.jpg",
		"65535/53961664838_49a7d74e87_o.jpg":  "65535/53961664838_49a7d74e87_o.jpg",
		"method=flickr.photosets.getPhotos&photoset_id=72177720319945125&user_id=35034348999%40N01/#/65535/53961664838_49a7d74e87_o.jpg": "/65535/53961664838_49a7d74e87_o.jpg",
	}

	for str, expected := range tests {

		if !MatchesPhotoURL(str) {
			t.Fatalf("Failed to match '%s'", str)
		}

		v, err := DerivePhotoURL(str)

		if err != nil {
			t.Fatalf("Failed to derive photo URL from '%s', %v", str, err)
		}

		if v != expected {
			t.Fatalf("Unexpected match for '%s', expected '%s' but got '%s'", str, expected, v)
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

	tests := []apiTest{
		apiTest{
			// https://flickr.com/photos/straup/6923069836
			photo_id: "6923069836",
			// https://flickr.com/photos/straup
			user_id: "35034348999@N01",
			// https://flickr.com/photos/straup/albums/72157629455113026/
			photoset_id: "72157629455113026",
		},
		apiTest{
			// https://flickr.com/photos/bees/65753018
			photo_id: "65753018",
			// https://flickr.com/photos/bees
			user_id: "12037949754@N01",
			// https://flickr.com/photos/bees/albums/1418449/
			photoset_id: "1418449",
		},
	}

	for _, fs_test := range tests {

		fl, err := fs.Open(fs_test.photo_id)

		if err != nil {
			t.Fatalf("Failed to open photo %s, %v", fs_test.photo_id, err)
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

		// https://www.flickr.com/services/api/flickr.photosets.getPhotos.html
		u.Set("method", "flickr.photosets.getPhotos")
		u.Set("photoset_id", fs_test.photoset_id)
		u.Set("user_id", fs_test.user_id)
		q := u.Encode()

		err = io_fs.WalkDir(fs, q, walk_func)

		if err != nil {
			t.Fatalf("Failed to walk '%s', %v", q, err)
		}
	}
}
