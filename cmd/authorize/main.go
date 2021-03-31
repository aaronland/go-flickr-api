package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http"
	"log"
	gohttp "net/http"
	"net/url"
	"os"
)

func main() {

	client_uri := flag.String("client-uri", "", "...")

	flag.Parse()

	host := "localhost"
	port := 8080

	addr := fmt.Sprintf("%s:%d", host, port)

	cb_url, err := url.Parse(addr)

	if err != nil {
		log.Fatalf("Failed to parse '%s', %v", addr, err)
	}

	ctx := context.Background()

	cl, err := client.NewHTTPClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	req_token, err := cl.GetRequestToken(ctx, cb_url)

	if err != nil {
		log.Fatalf("Failed to create request token, %v", err)
	}

	auth_url, err := cl.AuthorizationURL(ctx, req_token)

	if err != nil {
		log.Fatalf("Failed to create authorization URL, %v", err)
	}

	token_ch := make(chan *auth.AuthorizationToken)
	err_ch := make(chan error)

	auth_handler, err := http.NewOAuth1AuthorizeTokenHandler(token_ch, err_ch)

	if err != nil {
		log.Fatalf("Failed to create request handler, %v", err)
	}

	mux := gohttp.NewServeMux()
	mux.Handle("/", auth_handler)

	go func() {

		err := gohttp.ListenAndServe(addr, mux)

		if err != nil {
			panic(err)
		}
	}()

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
