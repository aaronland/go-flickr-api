package api

import (
	"context"
	"flag"
	"image"
	_ "image/jpeg"
	"testing"

	"github.com/aaronland/go-flickr-api/client"
)

var client_uri = flag.String("client-uri", "", "...")

func TestFS(t *testing.T) {

	if *client_uri == "" {
		t.Skip()
	}

	ctx := context.Background()
	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		t.Fatalf("Failed to create new client, %v", err)
	}

	fs := NewFS(ctx, cl)

	fl, err := fs.Open("53961664838")

	if err != nil {
		t.Fatalf("Failed to open , %v", err)
	}

	_, _, err = image.Decode(fl)

	if err != nil {
		t.Fatalf("Failed to decode image, %v", err)
	}

}
