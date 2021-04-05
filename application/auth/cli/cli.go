package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http/oauth1"
	"github.com/aaronland/go-http-server"
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

type AuthApplication struct{}

func (app *AuthApplication) DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("auth")

	fs.StringVar(&client_uri, "client-uri", "", "...")
	fs.StringVar(&server_uri, "server-uri", "", "...")
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

	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	svr, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %v", err)
	}

	token_ch := make(chan auth.AuthorizationToken)
	err_ch := make(chan error)

	auth_handler, err := oauth1.NewAuthorizationTokenHandlerWithChannels(token_ch, err_ch)

	if err != nil {
		return fmt.Errorf("Failed to create request handler, %v", err)
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
		return fmt.Errorf("Failed to create client, %v", err)
	}

	req_token, err := cl.GetRequestToken(ctx, svr.Address())

	if err != nil {
		return fmt.Errorf("Failed to create request token, %v", err)
	}

	auth_url, err := cl.GetAuthorizationURL(ctx, req_token, perms)

	if err != nil {
		return fmt.Errorf("Failed to create authorization URL, %v", err)
	}

	log.Printf("Authorize this application %s\n", auth_url)

	// Add a timeout Timer here

	var auth_token auth.AuthorizationToken

	for {
		select {
		case err := <-err_ch:
			return fmt.Errorf("Failed to authorize request, %v", err)
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
		return fmt.Errorf("Failed to get access token, %v", err)
	}

	cl, err = cl.WithAccessToken(ctx, access_token)

	if err != nil {
		return fmt.Errorf("Failed to assign client with access token, %v", err)
	}

	args := &url.Values{}
	args.Set("method", "flickr.test.login")

	_, err = cl.ExecuteMethod(ctx, args)

	if err != nil {
		return fmt.Errorf("Failed to test login, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(access_token)

	if err != nil {
		return fmt.Errorf("Failed to write access token, %v", err)
	}

	return nil
}
