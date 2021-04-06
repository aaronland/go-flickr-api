package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/auth/www"
	_ "gocloud.dev/docstore/memdocstore"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"log"
)

func main() {

	ctx := context.Background()

	app := &www.AuthApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run auth application, %v", err)
	}
}
