package replace

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/application"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
	"github.com/aaronland/go-flickr-api/response"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/runtimevar"
	_ "gocloud.dev/runtimevar/constantvar"
	"io"
	"log"
	"net/url"
	"os"
)

var params multi.KeyValueString
var client_uri string
var use_runtimevar bool

// ReplaceApplication implements the application.Application interface as a commandline application for
// replacing photos using the Flickr API
type ReplaceApplication struct {
	application.Application
}

// Return the default FlagSet necessary for the ReplaceApplication to run.
func (app *ReplaceApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("upload")

	fs.StringVar(&client_uri, "client-uri", "", "A valid aaronland/go-flickr-api client URI.")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.")
	fs.Var(&params, "param", "Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}

// Invoke the ReplaceApplication with its default FlagSet.
func (app *ReplaceApplication) Run(ctx context.Context) (interface{}, error) {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

// Invoke the ReplaceApplication with a custom FlagSet.
func (app *ReplaceApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) (interface{}, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "FLICKR")

	if err != nil {
		return nil, fmt.Errorf("Failed to set flags from environment variables, %v", err)
	}

	paths := fs.Args()

	if use_runtimevar {

		runtime_uri, err := runtimevar.StringVar(ctx, client_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive runtime value for client URI, %v", err)
		}

		client_uri = runtime_uri
	}

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create client, %v", err)
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}

	// Do this concurrently...

	for _, path := range paths {

		fh, err := reader.NewReader(ctx, path)

		if err != nil {
			return nil, fmt.Errorf("Failed to create reader for '%s', %v", path, err)
		}

		defer fh.Close()

		rsp, err := cl.Replace(ctx, fh, args)

		if err != nil {
			log.Printf("Failed to replace '%s', %v", err)
		}

		up, err := response.UnmarshalUploadResponse(rsp)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal upload response, %v", err)
		}

		if up.Error != nil {
			return nil, fmt.Errorf("Upload failed, %v", up.Error)
		}

		log.Println(up, up.Error)

		rsp.Seek(0, 0)
		io.Copy(os.Stdout, rsp)
	}

	return nil, nil
}
