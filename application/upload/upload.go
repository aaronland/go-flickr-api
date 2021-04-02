package upload

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"net/url"
)

type UploadApplication struct {
}

func (app *UploadApplication) Run(ctx context.Context) error {

	var params multi.KeyValueString
	flag.Var(&params, "param", "...")

	client_uri := flag.String("client-uri", "oauth1://", "...")

	flag.Parse()

	paths := flag.Args()

	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		return fmt.Errorf("Failed to create client, %v", err)
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}

	// Do this concurrently...

	for _, path := range paths {

		fh, err := reader.NewReader(ctx, path)

		if err != nil {
			return fmt.Errorf("Failed to create reader for '%s', %v", path, err)
		}

		defer fh.Close()

		photo_id, err := client.UploadAsyncWithClient(ctx, cl, fh, args)

		if err != nil {
			return fmt.Errorf("Failed to upload '%s', %v", err)
		}

		log.Println("OK", photo_id)
	}

	return nil
}
