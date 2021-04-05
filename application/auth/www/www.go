package www

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http/oauth1"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/runtimevar"
	"gocloud.dev/docstore"
	"log"
	"net/http"
	"net/url"
)

var collection_uri string

var client_uri string
var server_uri string
var perms string
var use_runtimevar bool

type AuthApplication struct{}

func (app *AuthApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("auth")

	fs.StringVar(&client_uri, "client-uri", "", "...")
	fs.StringVar(&server_uri, "server-uri", "", "...")
	fs.StringVar(&collection_uri, "collection-uri", "", "...")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "")
	fs.StringVar(&perms, "permissions", "", "")

	return fs
}

func (app *AuthApplication) Run(ctx context.Context) error {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

func (app *AuthApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

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

		runtime_uri, err = runtimevar.StringVar(ctx, server_uri)

		if err != nil {
			return fmt.Errorf("Failed to derive runtime value for server URI, %v", err)
		}

		server_uri = runtime_uri

		runtime_uri, err = runtimevar.StringVar(ctx, collection_uri)

		if err != nil {
			return fmt.Errorf("Failed to derive runtime value for collection URI, %v", err)
		}

		collection_uri = runtime_uri
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return fmt.Errorf("Failed to create client, %v", err)
	}

	col, err := docstore.OpenCollection(ctx, collection_uri)

	if err != nil {
		return fmt.Errorf("Failed to open collection, %v", err)
	}

	svr, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %v", err)
	}

	req_uri := "/"
	auth_uri := "/auth"

	cb_url, _ := url.Parse(svr.Address())
	cb_url.Path = auth_uri

	auth_callback := cb_url.String()

	request_opts := &oauth1.RequestTokenHandlerOptions{
		Client:       cl,
		Collection:   col,
		Permissions:  perms,
		AuthCallback: auth_callback,
	}

	request_handler, err := oauth1.NewRequestTokenHandler(request_opts)

	if err != nil {
		return fmt.Errorf("Failed to create request handler, %v", err)
	}

	auth_opts := &oauth1.AuthorizationTokenHandlerOptions{
		Client:     cl,
		Collection: col,
	}

	auth_handler, err := oauth1.NewAuthorizationTokenHandler(auth_opts)

	if err != nil {
		return fmt.Errorf("Failed to create authorization handler, %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(req_uri, request_handler)
	mux.Handle(auth_uri, auth_handler)

	log.Printf("Listening for requests on %s\n", svr.Address())
	err = svr.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %v", err)
	}

	return nil
}
