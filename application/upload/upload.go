package upload

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/reader"
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

type UploadResult struct {
	Path    string `json:"path,omitempty"`
	PhotoId int64  `json:"photoid,omitempty"`
	Error   error  `json:"error,omitempty"`
}

type UploadApplication struct{}

func (app *UploadApplication) DefaultFlagSet() *flag.FlagSet {

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

func (app *UploadApplication) Run(ctx context.Context) error {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

func (app *UploadApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

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

	done_ch := make(chan bool)
	rsp_ch := make(chan *UploadResult)

	for _, path := range paths {

		go func(path string) {

			defer func() {
				done_ch <- true
			}()

			rsp := &UploadResult{
				Path: path,
			}

			fh, err := reader.NewReader(ctx, path)

			if err != nil {
				rsp.Error = fmt.Errorf("Failed to create reader for '%s', %v", path, err)
				rsp_ch <- rsp
				return
			}

			defer fh.Close()

			photo_id, err := client.UploadAsyncWithClient(ctx, cl, fh, args)

			if err != nil {
				rsp.Error = fmt.Errorf("Failed to upload '%s', %v", err)
				rsp_ch <- rsp
				return
			}

			rsp.PhotoId = photo_id
			rsp_ch <- rsp
			return
		}(path)
	}

	remaining := len(paths)
	results := make([]*UploadResult, 0)

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case rsp := <-rsp_ch:
			results = append(results, rsp)
		default:
			// pass
		}
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(results)

	if err != nil {
		log.Fatalf("Failed to encode results, %v", err)
	}

	return nil
}
