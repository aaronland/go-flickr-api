package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/sfomuseum/go-flags/multi"
	_ "io"
	"log"
	"net/url"
	"os"
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

		fh, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", err)
		}

		_, err = cl.Upload(ctx, fh, args)

		if err != nil {
			log.Printf("Failed to upload '%s', %v", err)
		}
	}

}
