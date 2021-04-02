package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-flickr-api/http/oauth1"
	"github.com/aaronland/go-http-server"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/memdocstore"
	"log"
	"net/http"
	"net/url"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "")
	client_uri := flag.String("client-uri", "", "...")

	collection_uri := flag.String("collection-uri", "mem://collection/Token", "...")
	perms := flag.String("perms", "read", "...")

	flag.Parse()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	col, err := docstore.OpenCollection(ctx, *collection_uri)

	if err != nil {
		log.Fatalf("Failed to open collection, %v", err)
	}

	svr, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	req_uri := "/"
	auth_uri := "/auth"

	cb_url, _ := url.Parse(svr.Address())
	cb_url.Path = auth_uri

	auth_callback := cb_url.String()

	request_opts := &oauth1.RequestTokenHandlerOptions{
		Client:       cl,
		Collection:   col,
		Permissions:  *perms,
		AuthCallback: auth_callback,
	}

	request_handler, err := oauth1.NewRequestTokenHandler(request_opts)

	if err != nil {
		log.Fatalf("Failed to create request handler, %v", err)
	}

	auth_opts := &oauth1.AuthorizationTokenHandlerOptions{
		Client:     cl,
		Collection: col,
	}

	auth_handler, err := oauth1.NewAuthorizationTokenHandler(auth_opts)

	if err != nil {
		log.Fatalf("Failed to create authorization handler, %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(req_uri, request_handler)
	mux.Handle(auth_uri, auth_handler)

	log.Printf("Listening for requests on %s\n", svr.Address())
	err = svr.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
