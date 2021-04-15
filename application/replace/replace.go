package replace

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/application"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
	"github.com/aaronland/go-flickr-api/response"
	"github.com/mitchellh/go-wordwrap"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/runtimevar"
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
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that the -client-uri flag is encoded as a gocloud.dev/runtimevar string URI.")
	fs.Var(&params, "param", "Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for replacing an image in Flickr.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] path(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nNotes:\n\n")
		fmt.Fprintf(os.Stderr, wordwrap.WrapString("Under the hood the replace tool is using the GoCloud blob abstraction layer for reading files. By default only local files the file:// URI scheme are supported. If you need to read files from other sources you will need to clone this application and import the relevant packages. As a convenience if no URI scheme is included then each path will be resolved to its absolute URI and prepended with file://.\n", 80))

		fmt.Fprintf(os.Stderr, "\n")
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

		enc := json.NewEncoder(os.Stdout)
		err = enc.Encode(up)

		if err != nil {
			log.Fatalf("Failed to encode results, %v", err)
		}
	}

	return nil, nil
}
