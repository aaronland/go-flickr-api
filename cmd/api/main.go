package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-api/client"
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

	paginated := flag.Bool("paginated", false, "...")

	flag.Parse()

	if len(params) == 0 {
		log.Fatal("Missing one or more -param flags")
	}

	ctx := context.Background()

	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}

	cb := func(ctx context.Context, fh io.ReadSeekCloser, err error) error {

		if err != nil {
			return err
		}

		_, err = io.Copy(os.Stdout, fh)

		if err != nil {
			return err
		}

		return nil
	}

	if *paginated {

		err := client.ExecuteMethodPaginated(ctx, cl, args, cb)

		if err != nil {
			log.Fatalf("Failed to write method results, %v", err)
		}

	} else {

		fh, err := cl.ExecuteMethod(ctx, args)

		err = cb(ctx, fh, err)

		if err != nil {
			log.Fatalf("Failed to write method results, %v", err)
		}
	}
}
