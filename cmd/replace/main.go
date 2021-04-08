package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/replace"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"log"
)

func main() {

	ctx := context.Background()

	app := &replace.ReplaceApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run replace application, %v", err)
	}
}
