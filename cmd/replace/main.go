package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
	"github.com/aaronland/go-flickr-api/response"
	"github.com/sfomuseum/go-flags/multi"
	"io"
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

		fh, err := reader.NewReader(ctx, path)

		if err != nil {
			log.Fatalf("Failed to create reader for '%s', %v", path, err)
		}

		defer fh.Close()

		rsp, err := cl.Replace(ctx, fh, args)

		if err != nil {
			log.Printf("Failed to replace '%s', %v", err)
		}

		up, err := response.UnmarshalUploadResponse(rsp)

		if err != nil {
			log.Fatalf("Failed to unmarshal upload response, %v", err)
		}

		if up.Error != nil {
			log.Fatalf("Upload failed, %v", up.Error)
		}

		log.Println(up, up.Error)

		rsp.Seek(0, 0)
		io.Copy(os.Stdout, rsp)
	}

}
