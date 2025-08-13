package main

import (
	"context"
	"log"

	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"github.com/aaronland/go-flickr-api/application/auth/cli"
)

func main() {

	ctx := context.Background()

	app := &cli.AuthApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run auth application, %v", err)
	}
}
