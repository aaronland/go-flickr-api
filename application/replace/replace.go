package replace

import (
	"context"
	"flag"
	"fmt"
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

type ReplaceApplication struct{}

func (app *ReplaceApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("upload")

	fs.StringVar(&client_uri, "client-uri", "", "A valid aaronland/go-flickr-api client URI.")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.")
	fs.Var(&params, "param", "Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.")

	return fs
}

func (app *ReplaceApplication) Run(ctx context.Context) error {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

func (app *ReplaceApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "FLICKR")

	if err != nil {
		return err
	}

	paths := fs.Args()

	if use_runtimevar {

		runtime_uri, err := runtimevar.StringVar(ctx, client_uri)

		if err != nil {
			return fmt.Errorf("Failed to derive runtime value for client URI, %v", err)
		}

		client_uri = runtime_uri
	}

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

		rsp, err := cl.Replace(ctx, fh, args)

		if err != nil {
			log.Printf("Failed to replace '%s', %v", err)
		}

		up, err := response.UnmarshalUploadResponse(rsp)

		if err != nil {
			return fmt.Errorf("Failed to unmarshal upload response, %v", err)
		}

		if up.Error != nil {
			return fmt.Errorf("Upload failed, %v", up.Error)
		}

		log.Println(up, up.Error)

		rsp.Seek(0, 0)
		io.Copy(os.Stdout, rsp)
	}

	return nil
}
