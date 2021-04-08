package www

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/application"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http/oauth1"
	"github.com/aaronland/go-http-server"
	"github.com/mitchellh/go-wordwrap"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/runtimevar"
	"gocloud.dev/docstore"
	"log"
	"net/http"
	"net/url"
	"os"
)

var client_uri string
var server_uri string
var collection_uri string
var perms string
var use_runtimevar bool

// AuthApplication implements the application.Application interface as a commandline application to
// start an HTTP server for initiating a Flickr API autorization flow in a web browser.
type AuthApplication struct {
	application.Application
}

// Return the default FlagSet necessary for the AuthApplication to run.
func (app *AuthApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("auth")

	fs.StringVar(&client_uri, "client-uri", "", "A valid aaronland/go-flickr-api client URI.")
	fs.StringVar(&server_uri, "server-uri", "", "A valid aaronland/go-http-server URI.")
	fs.StringVar(&collection_uri, "collection-uri", "", "A valid gocloud.dev/docstore URI. The docstore is used to store token requests during the time a user is approving an authentication request.")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.")
	fs.StringVar(&perms, "permissions", "", "A valid Flickr API permissions flag.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "HTTP server for initiating a Flickr API autorization flow in a web browser.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nNotes:\n\n")		
		fmt.Fprintf(os.Stderr, wordwrap.WrapString("If you are running this application on localhost and are not using a 'tls://' server-uri flag (including your own TLS key and certificate) you will need to specify the 'mkcert://' server-uri flag and ensure that you have the https://github.com/FiloSottile/mkcert tool installed on your computer. This is because Flickr will automatically rewrite authorization callback URLs starting in 'http://' to 'https://' even if those URLs are pointing back to localhost.\n", 80))

		fmt.Fprintf(os.Stderr, "\n")
	}

	return fs
}

// Invoke the AuthApplication with its default FlagSet.
func (app *AuthApplication) Run(ctx context.Context) (interface{}, error) {
	fs := app.DefaultFlagSet()
	return app.RunWithFlagSet(ctx, fs)
}

// Invoke the AuthApplication with a custom FlagSet.
func (app *AuthApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) (interface{}, error) {

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

		runtime_uri, err = runtimevar.StringVar(ctx, server_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive runtime value for server URI, %v", err)
		}

		server_uri = runtime_uri

		runtime_uri, err = runtimevar.StringVar(ctx, collection_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive runtime value for collection URI, %v", err)
		}

		collection_uri = runtime_uri
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create client, %v", err)
	}

	col, err := docstore.OpenCollection(ctx, collection_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %v", err)
	}

	svr, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new server, %v", err)
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
		return nil, fmt.Errorf("Failed to create request handler, %v", err)
	}

	auth_opts := &oauth1.AuthorizationTokenHandlerOptions{
		Client:     cl,
		Collection: col,
	}

	auth_handler, err := oauth1.NewAuthorizationTokenHandler(auth_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create authorization handler, %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(req_uri, request_handler)
	mux.Handle(auth_uri, auth_handler)

	log.Printf("Listening for requests on %s\n", svr.Address())
	err = svr.ListenAndServe(ctx, mux)

	if err != nil {
		return nil, fmt.Errorf("Failed to start server, %v", err)
	}

	return nil, nil
}
