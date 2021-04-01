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

	fh, err := cl.ExecuteMethod(ctx, args)

	if err != nil {
		log.Fatalf("Failed to execute method, %v", err)
	}

	_, err = io.Copy(os.Stdout, fh)

	if err != nil {
		log.Fatalf("Failed to write method results, %v", err)
	}
}
