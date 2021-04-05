package api

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/runtimevar"
	"io"
	"log"
	"net/url"
	"os"
)

var params multi.KeyValueString
var client_uri string
var use_runtimevar bool
var paginated bool

type APIApplication struct{}

func (app *APIApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("upload")

	fs.StringVar(&client_uri, "client-uri", "", "A valid aaronland/go-flickr-api client URI.")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.")
	fs.BoolVar(&paginated, "paginated", false, "Automatically paginate (and iterate through) all API responses.")
	fs.Var(&params, "param", "Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.")

	return fs
}

func (app *APIApplication) Run(ctx context.Context) error {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

func (app *APIApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "FLICKR")

	if err != nil {
		return err
	}

	if use_runtimevar {

		runtime_uri, err := runtimevar.StringVar(ctx, client_uri)

		if err != nil {
			return fmt.Errorf("Failed to derive runtime value for client URI, %v", err)
		}

		client_uri = runtime_uri
	}

	cl, err := client.NewClient(ctx, client_uri)

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

	if paginated {

		err := client.ExecuteMethodPaginatedWithClient(ctx, cl, args, cb)

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

	return nil
}
