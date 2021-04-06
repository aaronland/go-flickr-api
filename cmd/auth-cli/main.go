package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/auth/cli"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"log"
)

func main() {

	ctx := context.Background()

	app := &cli.AuthApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run auth application, %v", err)
	}
}
