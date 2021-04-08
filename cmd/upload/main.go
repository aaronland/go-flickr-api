package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/upload"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"log"
)

func main() {

	ctx := context.Background()

	app := &upload.UploadApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run upload application, %v", err)
	}
}
