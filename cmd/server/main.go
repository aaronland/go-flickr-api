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

	request_handler, err := oauth1.NewRequestTokenHandler(cl, col, *perms)

	if err != nil {
		log.Fatalf("Failed to create request handler, %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", request_handler)

	log.Printf("Listening for requests on %s\n", svr.Address())
	err = svr.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
