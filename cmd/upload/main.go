package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/response"
	"github.com/sfomuseum/go-flags/multi"
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

		rsp, err := cl.Upload(ctx, fh, args)

		if err != nil {
			log.Fatalf("Failed to upload '%s', %v", err)
		}

		up, err := response.UnmarshalUploadResponse(rsp)

		if err != nil {
			log.Fatalf("Failed to unmarshal upload response, %v", err)
		}

		if up.Error != nil {
			log.Fatalf("Upload failed, %v", err)
		}

		log.Println(up.PhotoId)
	}

}
