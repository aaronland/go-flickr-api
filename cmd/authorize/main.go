package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http"
	"github.com/aaronland/go-http-server"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "")
	client_uri := flag.String("client-uri", "", "...")
	perms := flag.String("perms", "read", "...")

	flag.Parse()

	ctx := context.Background()

	svr, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	token_ch := make(chan *auth.AuthorizationToken)
	err_ch := make(chan error)

	auth_handler, err := http.NewAuthorizationTokenHandler(token_ch, err_ch)

	if err != nil {
		log.Fatalf("Failed to create request handler, %v", err)
	}

	mux := gohttp.NewServeMux()
	mux.Handle("/", auth_handler)

	go func() {

		log.Printf("Listening for requests on %s\n", svr.Address())
		err := svr.ListenAndServe(ctx, mux)

		if err != nil {
			panic(err)
		}
	}()

	cl, err := client.NewHTTPClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	req_token, err := cl.GetRequestToken(ctx, svr.Address())

	if err != nil {
		log.Fatalf("Failed to create request token, %v", err)
	}

	auth_url, err := cl.AuthorizationURL(ctx, req_token, *perms)

	if err != nil {
		log.Fatalf("Failed to create authorization URL, %v", err)
	}

	log.Printf("Authorize this application %s\n", auth_url)

	// Add a timeout Timer here

	var auth_token *auth.AuthorizationToken

	for {
		select {
		case err := <-err_ch:
			log.Fatalf("Failed to authorize request, %v", err)
		case t := <-token_ch:
			auth_token = t
		default:
			// pass
		}

		if auth_token != nil {
			break
		}
	}

	access_token, err := cl.GetAccessToken(ctx, auth_token)

	if err != nil {
		log.Fatalf("Failed to get access token, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(access_token)

	if err != nil {
		log.Fatalf("Failed to write access token, %v", err)
	}

}
