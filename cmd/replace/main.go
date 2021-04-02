package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/replace"
	_ "gocloud.dev/runtimevar/constantvar"
	"log"
)

func main() {

	ctx := context.Background()

	app := &replace.ReplaceApplication{}
	err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run replace application, %v", err)
	}
}
