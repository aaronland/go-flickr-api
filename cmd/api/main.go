package main

import (
	"context"
	"log"

	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"github.com/aaronland/go-flickr-api/application/api"
)

func main() {

	ctx := context.Background()

	app := &api.APIApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run upload application, %v", err)
	}
}
