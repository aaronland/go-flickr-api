package api

import (
	"context"
	"flag"
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

	f := NewFS(ctx, cl)

	_, err = f.Open("53961664838")

	if err != nil {
		t.Fatalf("Failed to open , %v", err)
	}

}
