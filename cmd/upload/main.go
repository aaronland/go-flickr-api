package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"net/url"
)

func main() {

	var params multi.KeyValueString
	flag.Var(&params, "param", "...")

	client_uri := flag.String("client-uri", "oauth1://", "...")

	flag.Parse()

	paths := flag.Args()

	ctx := context.Background()

	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}

	// Do this concurrently...

	for _, path := range paths {

		fh, err := reader.NewReader(ctx, path)

		if err != nil {
			log.Fatalf("Failed to create reader for '%s', %v", path, err)
		}

		defer fh.Close()

		photo_id, err := client.UploadAsyncWithClient(ctx, cl, fh, args)

		if err != nil {
			log.Fatalf("Failed to upload '%s', %v", err)
		}

		log.Println("OK", photo_id)
	}

}
