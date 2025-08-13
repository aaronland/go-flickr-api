package main

import (
	"context"
	"log"

	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"github.com/aaronland/go-flickr-api/application/upload"
)

func main() {

	ctx := context.Background()

	app := &upload.UploadApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run upload application, %v", err)
	}
}
