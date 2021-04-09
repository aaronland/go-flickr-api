package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/application"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http/oauth1"
	"github.com/aaronland/go-http-server"
	"github.com/mitchellh/go-wordwrap"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/runtimevar"
	"log"
	"net/http"
	"net/url"
	"os"
)

var client_uri string
var server_uri string
var perms string
var use_runtimevar bool

// AuthApplication implements the application.Application interface as a commandline application to
// initiate a Flickr API authorization flow. This application will launch a background HTTP process
// to receive authorization callback requests (from Flickr) and block execution, using channels, until
// an authorization request is approved or triggers and error.
type AuthApplication struct {
	application.Application
}

// Return the default FlagSet necessary for the AuthApplication to run.
func (app *AuthApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("auth")

	fs.StringVar(&client_uri, "client-uri", "", "A valid aaronland/go-flickr-api client URI.")
	fs.StringVar(&server_uri, "server-uri", "", "A valid aaronland/go-http-server URI.")
	fs.BoolVar(&use_runtimevar, "use-runtimevar", false, "Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.")
	fs.StringVar(&perms, "permissions", "", "A valid Flickr API permissions flag.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for initiating a Flickr API authorization flow.\n\n")
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
			return nil, fmt.Errorf("Failed to derive runtime value for client URI '%s', %v", client_uri, err)
		}

		client_uri = runtime_uri
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	svr, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new server, %v", err)
	}

	token_ch := make(chan auth.AuthorizationToken)
	err_ch := make(chan error)

	auth_handler, err := oauth1.NewAuthorizationTokenHandlerWithChannels(token_ch, err_ch)

	if err != nil {
		return nil, fmt.Errorf("Failed to create request handler, %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", auth_handler)

	go func() {

		log.Printf("Listening for requests on %s\n", svr.Address())
		err := svr.ListenAndServe(ctx, mux)

		if err != nil {
			panic(err)
		}
	}()

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create client, %v", err)
	}

	req_token, err := cl.GetRequestToken(ctx, svr.Address())

	if err != nil {
		return nil, fmt.Errorf("Failed to create request token, %v", err)
	}

	auth_url, err := cl.GetAuthorizationURL(ctx, req_token, perms)

	if err != nil {
		return nil, fmt.Errorf("Failed to create authorization URL, %v", err)
	}

	log.Printf("Authorize this application %s\n", auth_url)

	// Add a timeout Timer here

	var auth_token auth.AuthorizationToken

	for {
		select {
		case err := <-err_ch:
			return nil, fmt.Errorf("Failed to authorize request, %v", err)
		case t := <-token_ch:
			auth_token = t
		default:
			// pass
		}

		if auth_token != nil {
			break
		}
	}

	access_token, err := cl.GetAccessToken(ctx, req_token, auth_token)

	if err != nil {
		return nil, fmt.Errorf("Failed to get access token, %v", err)
	}

	cl, err = cl.WithAccessToken(ctx, access_token)

	if err != nil {
		return nil, fmt.Errorf("Failed to assign client with access token, %v", err)
	}

	args := &url.Values{}
	args.Set("method", "flickr.test.login")

	_, err = cl.ExecuteMethod(ctx, args)

	if err != nil {
		return nil, fmt.Errorf("Failed to test login, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(access_token)

	if err != nil {
		return nil, fmt.Errorf("Failed to write access token, %v", err)
	}

	return nil, nil
}
