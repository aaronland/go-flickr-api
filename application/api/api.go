package api

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/application"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/gocloud/runtimevar"
	"github.com/mitchellh/go-wordwrap"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"io"
	_ "log"
	"net/url"
	"os"
)

var params multi.KeyValueString
var client_uri string
var use_runtimevar bool
var paginated bool

// APIApplication implements the application.Application interface as a commandline application to invoke
// the Flickr API and output results to STDOUT. It does not support uploading or replacing photos.
type APIApplication struct {
	application.Application
}

// Return the default FlagSet necessary for the APIApplication to run.
func (app *APIApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("upload")

	fs.StringVar(&client_uri, "client-uri", "", "A valid aaronland/go-flickr-api client URI.")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that the -client-uri flag is encoded as a gocloud.dev/runtimevar string URI.")
	fs.BoolVar(&paginated, "paginated", false, "Automatically paginate (and iterate through) all API responses.")
	fs.Var(&params, "param", "Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.")

	fs.Usage = func() {
		fmt.Fprint(os.Stderr, wordwrap.WrapString("Command-line tool for invoking the Flickr API. Results are emitted to STDOUT.\n\n", 80))
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options]\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, "\nNotes:\n\n")
		fmt.Fprint(os.Stderr, wordwrap.WrapString("Uploading and replacing images are not supported by this tool. You can use the 'upload' and 'replace' tools, respectively, for those tasks.\n", 80))

		fmt.Fprintf(os.Stderr, "\n")
	}

	return fs
}

// Invoke the APIApplication with its default FlagSet.
func (app *APIApplication) Run(ctx context.Context) (interface{}, error) {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

// Invoke the APIApplication with a custom FlagSet.
func (app *APIApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) (interface{}, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "FLICKR")

	if err != nil {
		return nil, fmt.Errorf("Failed to set flags from environment variables, %v", err)
	}

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
			return nil, fmt.Errorf("Failed to write method results, %v", err)
		}

	} else {

		fh, err := cl.ExecuteMethod(ctx, args)

		err = cb(ctx, fh, err)

		if err != nil {
			return nil, fmt.Errorf("Failed to write method results, %v", err)
		}
	}

	return nil, nil
}
