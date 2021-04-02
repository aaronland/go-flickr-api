package upload

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	_ "github.com/sfomuseum/runtimevar"
	_ "gocloud.dev/runtimevar/constantvar"
	"log"
	"net/url"
)

var params multi.KeyValueString
var client_uri string

type UploadApplication struct {
}

func (app *UploadApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("upload")

	fs.StringVar(&client_uri, "client-uri", "", "...")
	fs.Var(&params, "param", "...")

	return fs
}

func (app *UploadApplication) Run(ctx context.Context) error {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

func (app *UploadApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	paths := fs.Args()

	/*
	client_uri, err := runtimevar.StringVar(ctx, client_uri)

	if err != nil {
		return err
	}
	*/
	
	cl, err := client.NewClient(ctx, client_uri)

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
